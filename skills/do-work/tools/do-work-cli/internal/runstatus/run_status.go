// Package runstatus owns `run-status`: one row per claimed, blocked, waiting or
// earmarked REQ with its phase, ages, ETA, one class (C1 to C8) and one remedy.
//
// The board's partition, placement text and activity correlation have one home,
// the queue-kanban module, which this standard-library module may not import. So
// the caller (actions/status.md) writes `queue-kanban open-work --format json` to
// a file and passes it as --board-facts; this command quotes those facts and adds
// what only the core module reads: frontmatter, the frozen estimate, the run
// directory, builder branches and unfinished finalization.
//
// Read-only. Ages are reported; nothing here judges a run or builder dead, and no
// process is looked up or signalled (process checks belong to the action layer).
package runstatus

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/doctor"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/nextselection"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// quietActivityBoundary splits C1 (progressing) from C2 (quiet). Display only:
// it is approximate and authorizes nothing.
const quietActivityBoundary = 20 * time.Minute

const runLocalLabel = "run-local, not interpreted"

// classPrecedence is the order classes are tested in (first match wins) and the
// order rows are listed in. actions/status.md documents the same order.
var classPrecedence = []string{"C8", "C3", "C4", "C5", "C6", "C7", "C2", "C1"}

var classNames = map[string]string{
	"C1": "progressing", "C2": "quiet", "C3": "hand-back landed, not integrated", "C4": "needs operator",
	"C5": "waiting on dependencies", "C6": "earmarked", "C7": "claim past threshold", "C8": "finalization pending",
}

// reportedColumns are the board columns that get a row. A ready pending REQ is
// neither stuck nor waiting, so it has no class; `do-work roadmap` lists it.
var reportedColumns = map[string]bool{"claimed": true, "needs-input-or-blocked": true, "pending-waiting": true, "pending-earmarked": true}

var (
	requestIDPattern        = regexp.MustCompile(`^REQ-[0-9]+$`)
	requestPrefixedPattern  = regexp.MustCompile(`^REQ-[0-9]+-`)
	utcInstantPattern       = regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}(:[0-9]{2})?Z`)
	openQuestionItemPattern = regexp.MustCompile(`^- \[ \]`)
)

// boardFacts mirrors `queue-kanban open-work --format json`.
type boardFacts struct {
	GeneratedAt                string         `json:"generated_at"`
	StaleClaimThresholdMinutes int            `json:"stale_claim_threshold_minutes"`
	Requests                   []boardRequest `json:"requests"`
}

type boardRequest struct {
	ID                string   `json:"id"`
	Title             string   `json:"title"`
	Status            string   `json:"status"`
	Column            string   `json:"column"`
	PlacementReason   string   `json:"placement_reason"`
	UnmetDependencies []string `json:"unmet_dependencies"`
	AssignedTo        string   `json:"assigned_to"`
	LastActivityAt    string   `json:"last_activity_at"`
	LastActivityKind  string   `json:"last_activity_kind"`
	LastActivityPhase string   `json:"last_activity_phase"`
}

type commandOptions struct {
	boardFactsPath string
	runDirectory   string
	requestID      string
	watch          bool
}

// statusRow is one report row plus the finding that carries its class.
type statusRow struct {
	record  resultmodel.RunStatusRow
	finding resultmodel.CommandFinding
	remedy  string
}

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{"run-status": func(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
		return runStatusAt(executionContext, arguments, time.Now().UTC())
	}}
}

// runStatusAt is the command with the clock injected, so tests pin ages.
func runStatusAt(executionContext commandruntime.ExecutionContext, arguments []string, now time.Time) resultmodel.CommandResult {
	repositoryRoot := executionContext.RepositoryRoot
	options, parseError := parseCommandOptions(arguments)
	if parseError != nil {
		return refusal("RUN-STATUS-USAGE", parseError.Error(), "the command line does not name a valid report")
	}
	facts, factsError := readBoardFacts(options.boardFactsPath)
	if factsError != nil {
		return refusal("RUN-STATUS-BOARD-FACTS", factsError.Error(), fmt.Sprintf(
			"run-status needs the board's facts: build queue-kanban and run `queue-kanban open-work --format json --repo-root %s > FILE`, then pass --board-facts FILE (actions/status.md does both)", repositoryRoot))
	}
	runDirectory := resolveRunDirectory(repositoryRoot, options.runDirectory)
	snapshot, discoveryError := repositorymodel.DiscoverRepository(repositoryRoot)
	if discoveryError != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{{
			Code: "RUN-STATUS-DISCOVERY-FAILED", Severity: resultmodel.SeverityError, Evidence: []string{discoveryError.Error()},
			Fixability: resultmodel.FixabilityManual, AutomationStopReason: "the repository snapshot could not be built",
		}}}
	}

	finalizationByID := map[string]resultmodel.CommandFinding{}
	for _, finding := range doctor.FinalizationTailFindings(context.Background(), snapshot) {
		for _, requestID := range finding.AffectedIDs {
			if _, seen := finalizationByID[requestID]; !seen {
				finalizationByID[requestID] = finding
			}
		}
	}
	manifestLines := readManifestLines(runDirectory)
	threshold := time.Duration(facts.StaleClaimThresholdMinutes) * time.Minute

	rows := []statusRow{}
	for _, request := range facts.Requests {
		if !reportedColumns[request.Column] || (options.requestID != "" && request.ID != options.requestID) {
			continue
		}
		row := buildRow(repositoryRoot, snapshot, request, runDirectory, manifestLines, now)
		finalization, hasFinalization := finalizationByID[request.ID]
		classifyRow(&row, request, hasFinalization, finalization, threshold)
		rows = append(rows, row)
	}
	sort.SliceStable(rows, func(left, right int) bool {
		return slices.Index(classPrecedence, rows[left].record.Class) < slices.Index(classPrecedence, rows[right].record.Class)
	})

	report := &resultmodel.RunStatusResult{Rows: []resultmodel.RunStatusRow{}, RunLocalFiles: listRunLocalFiles(repositoryRoot, runDirectory, now)}
	if runDirectory != "" {
		report.RunDirectory = displayPath(repositoryRoot, runDirectory)
	}
	findings := []resultmodel.CommandFinding{}
	for _, row := range rows {
		report.Rows = append(report.Rows, row.record)
		findings = append(findings, row.finding)
	}
	text := renderReport(rows, report, options.watch, now)
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: findings, RunStatus: report, ExactTextOutput: &text}
}

func parseCommandOptions(arguments []string) (commandOptions, error) {
	options := commandOptions{}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--watch" {
			options.watch = true
			continue
		}
		var target *string
		switch argument {
		case "--board-facts":
			target = &options.boardFactsPath
		case "--run":
			target = &options.runDirectory
		case "--req":
			target = &options.requestID
		default:
			return options, fmt.Errorf("unknown argument %q: run-status takes --board-facts FILE [--run DIR] [--req REQ-NNN] [--watch]", argument)
		}
		index++
		if index >= len(arguments) || arguments[index] == "" {
			return options, fmt.Errorf("%s requires a value", argument)
		}
		*target = arguments[index]
	}
	if options.boardFactsPath == "" {
		return options, fmt.Errorf("--board-facts FILE is required")
	}
	if options.requestID != "" && !requestIDPattern.MatchString(options.requestID) {
		return options, fmt.Errorf("--req %q is not a REQ-NNN id", options.requestID)
	}
	return options, nil
}

func readBoardFacts(path string) (boardFacts, error) {
	var facts boardFacts
	contents, readError := os.ReadFile(path)
	if readError != nil {
		return facts, fmt.Errorf("board facts unreadable: %w", readError)
	}
	if decodeError := json.Unmarshal(contents, &facts); decodeError != nil {
		return facts, fmt.Errorf("board facts %s are not open-work JSON: %w", path, decodeError)
	}
	return facts, nil
}

// resolveRunDirectory returns --run, or the newest do-work/runs/work-* directory
// by name, or "" when there is none (not an error: run fields stay absent).
func resolveRunDirectory(repositoryRoot, override string) string {
	if override != "" {
		if !filepath.IsAbs(override) {
			override = filepath.Join(repositoryRoot, override)
		}
		return override
	}
	entries, _ := os.ReadDir(filepath.Join(repositoryRoot, "do-work", "runs"))
	newest := ""
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "work-") && entry.Name() > newest {
			newest = entry.Name()
		}
	}
	if newest == "" {
		return ""
	}
	return filepath.Join(repositoryRoot, "do-work", "runs", newest)
}

func readManifestLines(runDirectory string) []string {
	if runDirectory == "" {
		return nil
	}
	contents, readError := os.ReadFile(filepath.Join(runDirectory, "manifest.md"))
	if readError != nil {
		return nil
	}
	return strings.Split(string(contents), "\n")
}

func buildRow(repositoryRoot string, snapshot *repositorymodel.RepositorySnapshot, request boardRequest, runDirectory string, manifestLines []string, now time.Time) statusRow {
	record := resultmodel.RunStatusRow{
		RequestID: request.ID, Title: request.Title, Status: request.Status, Column: request.Column, PlacementReason: request.PlacementReason,
		UnmetDependencies: append([]string{}, request.UnmetDependencies...), DependsOn: []string{}, BlockedBy: []string{}, AssignedTo: request.AssignedTo,
		LastActivityAt: request.LastActivityAt, LastActivityKind: request.LastActivityKind, LastActivityPhase: request.LastActivityPhase,
	}
	evidence := []string{fmt.Sprintf("column %s: %s", request.Column, request.PlacementReason)}
	if len(request.UnmetDependencies) > 0 {
		evidence = append(evidence, "waits on "+strings.Join(request.UnmetDependencies, ", ")+" (unmet)")
	}
	requestPath := ""
	if requestFile := openRequestFile(snapshot, request.ID); requestFile != nil {
		requestPath = "do-work/" + requestFile.RelativePath
		typed := requestFile.TypedRecord
		record.Status, record.ClaimedAt, record.AssignedTo = typed.RequestStatus, typed.ClaimedAt, typed.AssignedTo
		record.DependsOn = append(record.DependsOn, typed.DependsOn...)
		if blockedBy, found := typed.FieldEvidenceByName["blocked_by"]; found {
			record.BlockedBy = append(record.BlockedBy, blockedBy.ListValues...)
			if len(blockedBy.ListValues) == 0 && strings.TrimSpace(blockedBy.ScalarValue) != "" {
				record.BlockedBy = append(record.BlockedBy, strings.TrimSpace(blockedBy.ScalarValue))
			}
		}
		record.BlockedAt = strings.TrimSpace(typed.FieldEvidenceByName["blocked_at"].ScalarValue)
		if requestFile.ParsedDocument != nil {
			record.OpenQuestions = countOpenQuestions(string(requestFile.ParsedDocument.BodyBytes()))
		}
		if p50, known := nextselection.FrozenEstimate(typed.FieldEvidenceByName); known {
			record.P50ActiveMinutes = &p50
		}
	}
	if len(record.BlockedBy) > 0 {
		evidence = append(evidence, "blocked_by: "+strings.Join(record.BlockedBy, ", "))
	}
	if record.OpenQuestions > 0 {
		evidence = append(evidence, fmt.Sprintf("%d open question(s) under ## Open Questions", record.OpenQuestions))
	}
	if claimedAt, parseError := requestmodel.ParseTimestamp(record.ClaimedAt); record.ClaimedAt != "" && parseError == nil {
		record.MinutesSinceClaim = minutesBetween(claimedAt, now)
		evidence = append(evidence, fmt.Sprintf("claimed %d min ago", *record.MinutesSinceClaim))
	}
	if activityAt, parseError := time.Parse(time.RFC3339, request.LastActivityAt); parseError == nil {
		record.MinutesSinceActivity = minutesBetween(activityAt, now)
		evidence = append(evidence, fmt.Sprintf("last activity %d min ago (%s, phase %s)", *record.MinutesSinceActivity, request.LastActivityKind, fallbackDash(request.LastActivityPhase)))
	} else if request.Column == "claimed" {
		evidence = append(evidence, "no activity record")
	}
	record.EtaMinutes, record.EtaText = estimateRemaining(record.P50ActiveMinutes, record.MinutesSinceClaim)
	evidence = append(evidence, "ETA: "+record.EtaText)

	if runDirectory != "" {
		handbackPath := filepath.Join(runDirectory, request.ID+"-handback.md")
		_, statError := os.Stat(handbackPath)
		handbackPresent := statError == nil
		record.HandbackPresent = &handbackPresent
		if handbackPresent {
			evidence = append(evidence, "hand-back landed: "+displayPath(repositoryRoot, handbackPath))
		}
		record.ManifestRow, record.DispatchedAt = manifestRowFor(manifestLines, request.ID)
		if record.ManifestRow != "" {
			evidence = append(evidence, "manifest row: "+record.ManifestRow)
		}
	}
	record.BuilderBranch, record.BuilderBranchTipAt, record.OtherBuilderBranches = builderBranches(repositoryRoot, request.ID)
	if record.BuilderBranch != "" {
		evidence = append(evidence, fmt.Sprintf("builder branch %s, tip %s", record.BuilderBranch, record.BuilderBranchTipAt))
	}

	finding := resultmodel.CommandFinding{Severity: resultmodel.SeverityInfo, AffectedIDs: []string{request.ID}, Evidence: evidence, Fixability: resultmodel.FixabilityManual}
	if requestPath != "" {
		finding.AffectedPaths = []string{requestPath}
	}
	return statusRow{record: record, finding: finding}
}

// classifyRow sets the class, next_argv and remedy. First match wins, in
// classPrecedence order.
func classifyRow(row *statusRow, request boardRequest, hasFinalization bool, finalization resultmodel.CommandFinding, threshold time.Duration) {
	record, finding := &row.record, &row.finding
	claimed := request.Column == "claimed"
	switch {
	case hasFinalization:
		record.Class, finding.NextArgv = "C8", append([]string{}, finalization.NextArgv...)
		finding.Evidence = append(finding.Evidence, finalization.Evidence...)
		row.remedy = "an unfinished finalization holds it; run `" + strings.Join(finding.NextArgv, " ") + "`"
	case claimed && record.HandbackPresent != nil && *record.HandbackPresent:
		record.Class, finding.NextArgv = "C3", []string{"do-work", "run"}
		row.remedy = "the build is done and waits for integration; run `do-work run`"
	case request.Column == "needs-input-or-blocked" && (request.Status == "pending-answers" || request.Status == "blocked"):
		record.Class, finding.NextArgv = "C4", []string{"do-work", "clarify"}
		row.remedy = "it waits on a person; run `do-work clarify`"
	case request.Column == "needs-input-or-blocked":
		record.Class, finding.NextArgv = "C4", []string{"do-work", "forensics"}
		row.remedy = "status " + request.Status + " is not a question clarify answers; run `do-work forensics` to see what holds it"
	case request.Column == "pending-waiting":
		record.Class = "C5"
		row.remedy = "it waits on " + strings.Join(request.UnmetDependencies, ", ") + "; nothing to run until that is source-ready"
	case request.Column == "pending-earmarked":
		record.Class, finding.NextArgv = "C6", []string{"do-work", "run", request.ID}
		row.remedy = request.PlacementReason + " Run `do-work run " + request.ID + "` to run it by name"
	case claimed && record.MinutesSinceClaim != nil && time.Duration(*record.MinutesSinceClaim)*time.Minute >= threshold:
		// destructive-next-argv: next_argv is followed literally, so it is the
		// read-only inspection doctor's STUCK-WORK offers; the takeover is named
		// only here, with what it destroys and when it is safe.
		record.Class = "C7"
		finding.NextArgv = []string{"git", "log", "--full-history", "--", finding.AffectedPaths[0]} // claimed_at came from this REQ file
		row.remedy = fmt.Sprintf("claimed %d min ago, past the %d min threshold; inspect `%s`. Run `do-work run-with-recovery %s` only if you know the run that claimed it is gone: it requeues the claim and strips generated sections",
			*record.MinutesSinceClaim, int(threshold/time.Minute), strings.Join(finding.NextArgv, " "), request.ID)
	case claimed && (record.MinutesSinceActivity == nil || time.Duration(*record.MinutesSinceActivity)*time.Minute >= quietActivityBoundary):
		record.Class = "C2"
		row.remedy = "no activity record"
		if record.MinutesSinceActivity != nil {
			row.remedy = fmt.Sprintf("no new stamp or commit for %d min", *record.MinutesSinceActivity)
		}
		row.remedy += "; the status action checks the builder on this machine"
	default:
		record.Class = "C1"
		row.remedy = "progressing"
	}
	finding.Code = record.Class
	finding.AutomationStopReason = row.remedy
}

// listRunLocalFiles lists every run-directory entry the suite does not define:
// neither manifest.md nor named REQ-NNN-* (a condition, not a list of names).
func listRunLocalFiles(repositoryRoot, runDirectory string, now time.Time) []resultmodel.RunLocalFile {
	localFiles := []resultmodel.RunLocalFile{}
	if runDirectory == "" {
		return localFiles
	}
	entries, _ := os.ReadDir(runDirectory)
	for _, entry := range entries {
		if entry.Name() == "manifest.md" || requestPrefixedPattern.MatchString(entry.Name()) {
			continue
		}
		info, infoError := entry.Info()
		if infoError != nil {
			continue
		}
		entryPath := filepath.Join(runDirectory, entry.Name())
		localFile := resultmodel.RunLocalFile{Path: displayPath(repositoryRoot, entryPath), AgeMinutes: *minutesBetween(info.ModTime(), now), Label: runLocalLabel}
		if !entry.IsDir() {
			localFile.FirstLine = firstLineOf(entryPath)
		}
		localFiles = append(localFiles, localFile)
	}
	return localFiles
}

func firstLineOf(path string) string {
	file, openError := os.Open(path)
	if openError != nil {
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return ""
	}
	line := strings.TrimSpace(scanner.Text())
	if len(line) > 160 {
		line = line[:160] + "…"
	}
	return line
}

// manifestRowFor returns the manifest table line whose first cell is the REQ id,
// verbatim, and the first UTC instant in it as the dispatch instant. The column
// set varies between runs, so no other cell is interpreted.
func manifestRowFor(manifestLines []string, requestID string) (string, string) {
	for _, line := range manifestLines {
		trimmed := strings.TrimSpace(line)
		cells := strings.Split(trimmed, "|")
		if len(cells) < 3 || cells[0] != "" || strings.TrimSpace(cells[1]) != requestID {
			continue
		}
		return trimmed, utcInstantPattern.FindString(trimmed)
	}
	return "", ""
}

// builderBranches reads the REQ's worktree-agent branches, newest tip first. Any
// git failure (no repository, no git) leaves the fields absent.
func builderBranches(repositoryRoot, requestID string) (string, string, []string) {
	output, gitError := exec.Command("git", "-C", repositoryRoot, "for-each-ref", "--sort=-committerdate",
		"--format=%(refname:short) %(committerdate:iso-strict)", "refs/heads/worktree-agent-"+requestID+"-*").Output()
	if gitError != nil {
		return "", "", nil
	}
	newestBranch, newestTip, otherBranches := "", "", []string(nil)
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		branchName, tipAt, found := strings.Cut(strings.TrimSpace(line), " ")
		if !found {
			continue
		}
		if newestBranch == "" {
			newestBranch, newestTip = branchName, tipAt
		} else {
			otherBranches = append(otherBranches, branchName)
		}
	}
	return newestBranch, newestTip, otherBranches
}

// openRequestFile prefers the REQ's queue or working copy over an archived one.
func openRequestFile(snapshot *repositorymodel.RepositorySnapshot, requestID string) *repositorymodel.RequestFile {
	var archived *repositorymodel.RequestFile
	for _, requestFile := range snapshot.RequestsByID[requestID] {
		if requestFile.TreeSection != "archive" {
			return requestFile
		}
		if archived == nil {
			archived = requestFile
		}
	}
	return archived
}

// countOpenQuestions counts unchecked `- [ ]` items under `## Open Questions`.
func countOpenQuestions(body string) int {
	count, inSection := 0, false
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			inSection = trimmed == "## Open Questions"
			continue
		}
		if inSection && openQuestionItemPattern.MatchString(trimmed) {
			count++
		}
	}
	return count
}

// estimateRemaining is p50 minus minutes since claim. Past p50 it says how far
// over, never a negative number. Display only.
func estimateRemaining(p50 *int, minutesSinceClaim *int) (*int, string) {
	switch {
	case p50 == nil:
		return nil, "no estimate"
	case minutesSinceClaim == nil:
		return nil, fmt.Sprintf("%d min estimate, not claimed", *p50)
	case *minutesSinceClaim > *p50:
		return nil, fmt.Sprintf("over estimate by %d min", *minutesSinceClaim-*p50)
	}
	remaining := *p50 - *minutesSinceClaim
	return &remaining, fmt.Sprintf("%d min", remaining)
}

func refusal(code, evidence, stopReason string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{evidence}, Fixability: resultmodel.FixabilityManual,
		AutomationStopReason: stopReason, NextArgv: []string{"do-work", "status"},
	}}}
}

func minutesBetween(earlier, later time.Time) *int {
	minutes := int(later.Sub(earlier) / time.Minute)
	return &minutes
}

func displayPath(repositoryRoot, path string) string {
	if relative, relError := filepath.Rel(repositoryRoot, path); relError == nil && !strings.HasPrefix(relative, "..") {
		return filepath.ToSlash(relative)
	}
	return path
}
