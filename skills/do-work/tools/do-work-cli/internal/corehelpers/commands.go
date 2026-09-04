package corehelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/doctor"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/nextselection"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandPreflight           = "preflight"
	CommandQualify             = "qualify"
	CommandScopeDrift          = "scope-drift"
	CommandInventory           = "uncommitted-inventory"
	CommandAssociate           = "associate-files"
	CommandProtectedInventory  = "protected-inventory"
	CommandRecordCommit        = "record-commit-hash"
	CommandCaptureScreenshot   = "capture-screenshot"
	CommandAtomicDownload      = "atomic-download"
	CommandAddExclude          = "add-local-git-exclude"
	CommandBlockedCheck        = "run-blocked-check"
	CommandShowCommitDiff      = "show-commit-diff"
	CommandStageDeletion       = "stage-exact-deletion"
	CommandCleanupReservations = "cleanup-req-reservations"
	CommandRepairTimestamps    = "repair-req-timestamps"
	CommandAuditTimestamps     = "audit-archive-timestamps"
	CommandHandoffSurvey       = "handoff-state-survey"
	CommandArchiveCollision    = "archive-collision"
	CommandEstimateP50         = "estimate-p50"
	CommandNow                 = "now"
	CommandFrontmatter         = "frontmatter"
)

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandPreflight: handlePreflight, CommandQualify: handleQualify,
		CommandScopeDrift: handleScopeDrift, CommandInventory: handleInventory,
		CommandAssociate: handleAssociate, CommandProtectedInventory: handleProtectedInventory,
		CommandRecordCommit: handleRecordCommit, CommandCaptureScreenshot: handleCaptureScreenshot,
		CommandAtomicDownload: handleAtomicDownload, CommandAddExclude: handleAddExclude,
		CommandBlockedCheck: handleBlockedCheck, CommandShowCommitDiff: handleShowCommitDiff,
		CommandStageDeletion: handleStageDeletion, CommandCleanupReservations: handleCleanupReservations,
		CommandRepairTimestamps: handleRepairTimestamps, CommandAuditTimestamps: handleAuditTimestamps,
		CommandHandoffSurvey:    handleHandoffSurvey,
		CommandArchiveCollision: handleArchiveCollision, CommandEstimateP50: handleEstimateP50,
		CommandNow: handleNow, CommandFrontmatter: handleFrontmatter,
	}
}

func handleNow(_ commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	if len(arguments) != 0 {
		return usageResult(CommandNow, "now accepts no arguments")
	}
	output := time.Now().UTC().Truncate(time.Second).Format("2006-01-02T15:04:05Z") + "\n"
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: &output}
}

func handleFrontmatter(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	if len(arguments) < 3 || arguments[0] != "get" {
		return usageResult(CommandFrontmatter, "usage: frontmatter get <file> <field> [--normalize] [--in-set SET]")
	}
	filePath, field := arguments[1], arguments[2]
	normalize, setName := false, ""
	for index := 3; index < len(arguments); index++ {
		switch arguments[index] {
		case "--normalize":
			normalize = true
		case "--in-set":
			index++
			if index >= len(arguments) {
				return usageResult(CommandFrontmatter, "--in-set requires a value")
			}
			setName = arguments[index]
		default:
			return usageResult(CommandFrontmatter, "unknown option "+arguments[index])
		}
	}
	contents, err := os.ReadFile(absoluteFromRoot(executionContext.RepositoryRoot, filePath))
	if err != nil {
		return usageResult(CommandFrontmatter, err.Error())
	}
	value, found, err := flatFrontmatterValue(string(contents), field)
	if err != nil {
		return usageResult(CommandFrontmatter, err.Error())
	}
	if !found {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("FRONTMATTER-FIELD-MISSING", resultmodel.SeverityWarning, []string{filePath}, "field is absent: "+field, resultmodel.FixabilityManual, "the requested frontmatter fact is unavailable", nil, nil)}}
	}
	if normalize {
		value = strings.ToLower(strings.TrimSpace(value))
	}
	if setName != "" {
		inSet := false
		switch setName {
		case "terminal-success":
			inSet = value == "completed" || value == "completed-with-issues"
		case "terminal-resolved":
			inSet = value == "completed" || value == "completed-with-issues" || value == "cancelled" || value == "failed"
		default:
			return usageResult(CommandFrontmatter, "unknown set "+setName)
		}
		if inSet {
			value = "true"
		} else {
			value = "false"
		}
	}
	output := value + "\n"
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: &output}
}

func flatFrontmatterValue(contents, field string) (string, bool, error) {
	lines := strings.Split(strings.ReplaceAll(contents, "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimPrefix(lines[0], "\ufeff") != "---" {
		return "", false, fmt.Errorf("file has no frontmatter block")
	}
	for _, line := range lines[1:] {
		if line == "---" {
			return "", false, nil
		}
		if strings.HasPrefix(line, field+":") {
			value := strings.TrimSpace(strings.TrimPrefix(line, field+":"))
			value = strings.Trim(value, "'\"")
			return value, value != "", nil
		}
	}
	return "", false, fmt.Errorf("unterminated frontmatter block")
}

func handleArchiveCollision(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	if len(arguments) != 1 || !regexp.MustCompile(`^REQ-[0-9]+$`).MatchString(arguments[0]) {
		return usageResult(CommandArchiveCollision, "usage: archive-collision REQ-NNN")
	}
	archiveRoot := filepath.Join(executionContext.RepositoryRoot, "do-work", "archive")
	entries, err := os.ReadDir(archiveRoot)
	if os.IsNotExist(err) {
		output := ""
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: &output}
	}
	if err != nil {
		return usageResult(CommandArchiveCollision, err.Error())
	}
	paths := []string{}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && (name == arguments[0]+".md" || strings.HasPrefix(name, arguments[0]+"-")) && strings.HasSuffix(name, ".md") {
			paths = append(paths, filepath.ToSlash(filepath.Join("do-work/archive", name)))
		}
	}
	sort.Strings(paths)
	output := ""
	if len(paths) > 0 {
		output = strings.Join(paths, "\n") + "\n"
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, ExactTextOutput: &output}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: &output}
}

type estimateGraphNode struct {
	minutes      int
	dependencies []string
}

func handleEstimateP50(_ commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	output, outcome, err := estimateP50Output(arguments)
	if err != nil {
		return usageResult(CommandEstimateP50, err.Error())
	}
	return resultmodel.CommandResult{Outcome: outcome, ExactTextOutput: &output}
}

func estimateP50Output(arguments []string) (string, resultmodel.CommandOutcome, error) {
	if len(arguments) > 0 && arguments[0] == "critical-path" {
		if len(arguments) == 1 {
			return "", resultmodel.OutcomeFailure, fmt.Errorf("critical-path needs at least one graph entry")
		}
		nodes := map[string]estimateGraphNode{}
		total := 0
		for _, entry := range arguments[1:] {
			parts := strings.SplitN(entry, ":", 3)
			if len(parts) < 2 || parts[0] == "" {
				return "", resultmodel.OutcomeFailure, fmt.Errorf("malformed graph entry %q", entry)
			}
			minutes, err := strconv.Atoi(parts[1])
			if err != nil || minutes < 0 {
				return "", resultmodel.OutcomeFailure, fmt.Errorf("minutes for %s must be a non-negative integer", parts[0])
			}
			dependencies := []string{}
			if len(parts) == 3 && parts[2] != "" {
				dependencies = strings.Split(parts[2], ",")
			}
			nodes[parts[0]] = estimateGraphNode{minutes: minutes, dependencies: dependencies}
			total += minutes
		}
		memo, visiting := map[string]int{}, map[string]bool{}
		var longest func(string) (int, error)
		longest = func(identifier string) (int, error) {
			if value, ok := memo[identifier]; ok {
				return value, nil
			}
			node, ok := nodes[identifier]
			if !ok {
				return 0, nil
			}
			if visiting[identifier] {
				return 0, fmt.Errorf("dependency cycle involving %s", identifier)
			}
			visiting[identifier] = true
			best := 0
			for _, dependency := range node.dependencies {
				candidate, err := longest(dependency)
				if err != nil {
					return 0, err
				}
				if candidate > best {
					best = candidate
				}
			}
			delete(visiting, identifier)
			memo[identifier] = best + node.minutes
			return memo[identifier], nil
		}
		critical := 0
		for identifier := range nodes {
			candidate, err := longest(identifier)
			if err != nil {
				return "", resultmodel.OutcomeFindings, err
			}
			if candidate > critical {
				critical = candidate
			}
		}
		return fmt.Sprintf("total_estimated_effort_minutes: %d\ncritical_path_minutes: %d\n", total, critical), resultmodel.OutcomeSuccess, nil
	}

	route := ""
	counts := map[string]int{"write-set": 0, "new-files": 0, "subsystems": 1, "acceptance": 0, "deps-depth": 0}
	flags := map[string]bool{}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--trivial" {
			flags["trivial"] = true
			continue
		}
		if argument == "--route" {
			index++
			if index >= len(arguments) || !strings.Contains("ABCabc", arguments[index]) || len(arguments[index]) != 1 {
				return "", resultmodel.OutcomeFailure, fmt.Errorf("--route must be A, B, or C")
			}
			route = strings.ToUpper(arguments[index])
			continue
		}
		name := strings.TrimPrefix(argument, "--")
		if _, ok := counts[name]; ok {
			index++
			if index >= len(arguments) {
				return "", resultmodel.OutcomeFailure, fmt.Errorf("--%s needs a non-negative integer argument", name)
			}
			value, err := strconv.Atoi(arguments[index])
			if err != nil || value < 0 {
				return "", resultmodel.OutcomeFailure, fmt.Errorf("--%s needs a non-negative integer argument", name)
			}
			counts[name] = value
			continue
		}
		switch name {
		case "browser", "persistence", "async-behavior", "performance", "regression", "full-suite":
			flags[name] = true
		default:
			return "", resultmodel.OutcomeFailure, fmt.Errorf("unrecognized argument %q", argument)
		}
	}
	if flags["trivial"] {
		return "p50_active_minutes: 5\nconfidence: high\nbasis:\n- trivial short-circuit\n", resultmodel.OutcomeSuccess, nil
	}
	if route == "" {
		return "", resultmodel.OutcomeFailure, fmt.Errorf("--route is required (or use --trivial)")
	}
	raw := map[string]int{"A": 5, "B": 10, "C": 20}[route]
	raw += counts["write-set"] + counts["new-files"]*2 + counts["acceptance"] + counts["deps-depth"]*2
	if counts["subsystems"] > 1 {
		raw += (counts["subsystems"] - 1) * 3
	}
	weights := map[string]int{"browser": 8, "persistence": 6, "async-behavior": 6, "performance": 4, "regression": 4, "full-suite": 4}
	for name, weight := range weights {
		if flags[name] {
			raw += weight
		}
	}
	minutes := ((raw + 2) / 5) * 5
	if minutes < 5 {
		minutes = 5
	}
	confidence := "medium"
	if route == "A" && raw <= 10 {
		confidence = "high"
	}
	if route == "C" && (counts["write-set"] >= 15 || counts["subsystems"] >= 3 || raw >= 75) {
		confidence = "low"
	}
	var output strings.Builder
	fmt.Fprintf(&output, "p50_active_minutes: %d\nconfidence: %s\nbasis:\n- Route %s\n", minutes, confidence, route)
	labels := []struct{ key, suffix string }{{"write-set", "-file write set"}, {"new-files", " new files"}, {"subsystems", " subsystems involved"}, {"acceptance", " acceptance criteria"}, {"deps-depth", "dependency depth "}}
	for _, label := range labels {
		value := counts[label.key]
		if label.key == "subsystems" && value <= 1 || label.key != "subsystems" && value == 0 {
			continue
		}
		if label.key == "deps-depth" {
			fmt.Fprintf(&output, "- %s%d\n", label.suffix, value)
		} else {
			fmt.Fprintf(&output, "- %d%s\n", value, label.suffix)
		}
	}
	flagLabels := []struct{ key, label string }{{"browser", "browser evidence"}, {"persistence", "persistence changes"}, {"async-behavior", "async lifecycle behavior"}, {"performance", "performance instrumentation"}, {"regression", "cross-route regression gates"}, {"full-suite", "full-suite verification"}}
	for _, item := range flagLabels {
		if flags[item.key] {
			fmt.Fprintf(&output, "- %s\n", item.label)
		}
	}
	return output.String(), resultmodel.OutcomeSuccess, nil
}

func successResult(changes []resultmodel.RecordedChange, findings []resultmodel.CommandFinding) resultmodel.CommandResult {
	outcome := resultmodel.OutcomeSuccess
	for _, finding := range findings {
		if finding.Severity != resultmodel.SeverityInfo {
			outcome = resultmodel.OutcomeFindings
			break
		}
	}
	return resultmodel.CommandResult{Outcome: outcome, Changes: changes, Findings: findings}
}

func helperFinding(code string, severity resultmodel.FindingSeverity, paths []string, evidence string, fixability resultmodel.FindingFixability, stop string, next, verify []string) resultmodel.CommandFinding {
	if severity != resultmodel.SeverityInfo && len(next) == 0 {
		next, _ = findingSpecificCommands(code, paths)
	}
	if severity != resultmodel.SeverityInfo && len(verify) == 0 {
		_, verify = findingSpecificCommands(code, paths)
	}
	return resultmodel.CommandFinding{Code: code, Severity: severity, AffectedPaths: paths,
		Evidence: []string{evidence}, Fixability: fixability, AutomationStopReason: stop,
		NextArgv: next, VerificationArgv: verify}
}

func findingSpecificCommands(code string, paths []string) ([]string, []string) {
	withPaths := func(arguments ...string) []string {
		if len(paths) == 0 {
			return arguments
		}
		return append(append(arguments, "--"), paths...)
	}
	switch code {
	case "SCOPE-UNDECLARED-TOUCH":
		return withPaths("git", "diff"), withPaths("git", "diff", "--name-only")
	case "SCOPE-DECLARED-NOT-TOUCHED":
		return withPaths("git", "log", "-1"), withPaths("git", "status", "--short")
	case "SCOPE-PATH-LIST-MALFORMED":
		return withPaths("sed", "-n", "/^## Scope$/,/^## /p"), withPaths("git", "diff", "--check")
	case "QUALIFY-DEBUG-ARTIFACT", "QUALIFY-DEBUG-TOKEN-INTRODUCED", "QUALIFY-DEBUG-ARTIFACT-RELOCATED":
		return withPaths("git", "diff", "-U0"), withPaths("git", "diff", "--check")
	case "QUALIFY-LIBRARY-OUTPUT", "QUALIFY-OUTPUT-RELOCATED", "QUALIFY-REPORTER-OUTPUT":
		return withPaths("git", "diff", "-U0"), withPaths("git", "diff", "--word-diff=porcelain")
	case "QUALIFY-SUMMARY-MISSING", "QUALIFY-SUMMARY-MALFORMED", "QUALIFY-UNIFY-DISARMED", "QUALIFY-PAU-UNCHECKED":
		return withPaths("sed", "-n", "/^## Implementation Summary/,$p"), withPaths("git", "diff", "--check")
	case "QUALIFY-DELETED-PATH-PRESENT", "QUALIFY-CLAIMED-PATH-MISSING", "QUALIFY-PATH-NOT-IN-DIFF", "QUALIFY-NEW-FILE-UNWIRED":
		return withPaths("git", "diff", "--name-status"), withPaths("git", "status", "--short")
	case "QUALIFY-DIFF-RANGE-INVALID":
		return []string{"git", "rev-parse", "--verify", "HEAD"}, []string{"git", "merge-base", "HEAD", "HEAD"}
	case "QUALIFY-NO-PROJECT-FILES":
		return []string{"git", "diff", "--name-only", "--", ".", ":(exclude)do-work/"}, []string{"git", "status", "--short", "--", ".", ":(exclude)do-work/"}
	case "ASSOCIATION-UNOWNED":
		return withPaths("git", "diff"), withPaths("git", "status", "--short")
	case "RESERVATION-MALFORMED", "RESERVATION-RACED", "RESERVATION-REMOVE-FAILED", "RESERVATION-ROOT-RACED", "RESERVATION-ROOT-UNSAFE", "RESERVATION-GIT-AUTHORITY-UNAVAILABLE":
		return []string{"do-work-cli", CommandCleanupReservations, "--dry-run"}, []string{"do-work-cli", "--format", "json", CommandCleanupReservations, "--dry-run"}
	case "HANDOFF-WORKTREE-MISSING":
		return []string{"git", "worktree", "prune", "--dry-run"}, []string{"git", "worktree", "list", "--porcelain"}
	case "PREFLIGHT-GIT-UNAVAILABLE", "PREFLIGHT-STATUS-FAILED":
		return []string{"git", "--version"}, []string{"git", "status", "--short"}
	case "PREFLIGHT-DIRTY":
		return withPaths("git", "diff"), withPaths("git", "status", "--short")
	case "PREFLIGHT-BASELINE-RED", "PREFLIGHT-BASELINE-NOT-LAUNCHED":
		return withPaths("sed", "-n", "1,160p"), withPaths("test", "-s")
	case "PREFLIGHT-NODE-MODULES-MISSING":
		return []string{"npm", "install"}, []string{"test", "-d", "node_modules"}
	case "PREFLIGHT-VIRTUALENV-MISSING":
		return []string{"python3", "-m", "venv", ".venv"}, []string{"test", "-x", ".venv/bin/python"}
	case "TIMESTAMP-ARCHIVE-WALK-FAILED", "TIMESTAMP-REPAIR-PENDING", "TIMESTAMP-REPAIR-REFUSED", "TIMESTAMP-REPLACE-FAILED", "TIMESTAMP-CONTENT-LOSS", "TIMESTAMP-TRUNCATION-FLOOR-UNAVAILABLE":
		return []string{"do-work-cli", CommandAuditTimestamps, "--fix", "--dry-run"}, withPaths("git", "diff")
	case "DELETION-DRY-RUN-NOT-DELETED", "DELETION-STAGE-VERIFY-FAILED":
		return withPaths("git", "status", "--short"), withPaths("git", "diff", "--cached", "--name-status")
	case "SCREENSHOT-DESTINATION-OCCUPIED", "SCREENSHOT-PUBLISH-FAILED", "SCREENSHOT-SOURCE-CLEANUP-WARNING", "SCREENSHOT-STAGING-CLEANUP-WARNING":
		return withPaths("ls", "-ld"), withPaths("test", "-e")
	case "DOWNLOAD-TARGET-OCCUPIED", "DOWNLOAD-FAILED":
		if len(paths) > 0 {
			return []string{"ls", "-ld", paths[0]}, []string{"test", "!", "-e", paths[0]}
		}
	case "GIT-EXCLUDE-NOT-A-REPOSITORY":
		return []string{"git", "init"}, []string{"git", "rev-parse", "--git-dir"}
	case "FRONTMATTER-FIELD-MISSING":
		return withPaths("sed", "-n", "1,80p"), withPaths("git", "diff", "--check")
	default:
		return nil, nil
	}
	return nil, nil
}

func usageResult(commandName, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
		helperFinding("HELPER-USAGE", resultmodel.SeverityError, nil, evidence, resultmodel.FixabilityManual,
			"the command line is invalid", []string{"do-work-cli", commandName}, []string{"do-work-cli", "--format", "json", commandName}),
	}}
}

func optionValue(arguments []string, index *int, name string) (string, error) {
	argument := arguments[*index]
	if strings.HasPrefix(argument, name+"=") {
		value := strings.TrimPrefix(argument, name+"=")
		if value == "" {
			return "", fmt.Errorf("%s requires a value", name)
		}
		return value, nil
	}
	*index++
	if *index >= len(arguments) || arguments[*index] == "" {
		return "", fmt.Errorf("%s requires a value", name)
	}
	return arguments[*index], nil
}

func extractDryRun(arguments []string) ([]string, bool, error) {
	filtered := make([]string, 0, len(arguments))
	dryRun := false
	for _, argument := range arguments {
		if argument != "--dry-run" {
			filtered = append(filtered, argument)
			continue
		}
		if dryRun {
			return nil, false, fmt.Errorf("--dry-run may be supplied only once")
		}
		dryRun = true
	}
	return filtered, dryRun, nil
}

func handleRecordCommit(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	requestPath, implementationHash := "", ""
	verifyOnly, dryRun := false, false
	for index := 0; index < len(arguments); index++ {
		switch {
		case arguments[index] == "--verify":
			verifyOnly = true
		case arguments[index] == "--dry-run":
			dryRun = true
		case arguments[index] == "--request-path" || strings.HasPrefix(arguments[index], "--request-path="):
			value, err := optionValue(arguments, &index, "--request-path")
			if err != nil {
				return usageResult(CommandRecordCommit, err.Error())
			}
			requestPath = value
		case arguments[index] == "--implementation-hash" || strings.HasPrefix(arguments[index], "--implementation-hash="):
			value, err := optionValue(arguments, &index, "--implementation-hash")
			if err != nil {
				return usageResult(CommandRecordCommit, err.Error())
			}
			implementationHash = value
		default:
			return usageResult(CommandRecordCommit, "unknown option "+arguments[index])
		}
	}
	if verifyOnly && dryRun {
		return usageResult(CommandRecordCommit, "--verify and --dry-run cannot be combined")
	}
	return requeststate.RecordCommitProvenance(context.Background(), executionContext.RepositoryRoot, requestPath, implementationHash, verifyOnly, dryRun)
}

func handleBlockedCheck(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	probeFile, baselineJSONPath, baselineFailuresPath := "", "", ""
	timeoutSeconds := 30
	for index := 0; index < len(arguments); index++ {
		switch {
		case arguments[index] == "--probe-file" || strings.HasPrefix(arguments[index], "--probe-file="):
			value, err := optionValue(arguments, &index, "--probe-file")
			if err != nil {
				return usageResult(CommandBlockedCheck, err.Error())
			}
			probeFile = value
		case arguments[index] == "--timeout-seconds" || strings.HasPrefix(arguments[index], "--timeout-seconds="):
			value, err := optionValue(arguments, &index, "--timeout-seconds")
			if err != nil {
				return usageResult(CommandBlockedCheck, err.Error())
			}
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed <= 0 {
				return usageResult(CommandBlockedCheck, "--timeout-seconds requires a positive integer")
			}
			timeoutSeconds = parsed
		case arguments[index] == "--baseline-json" || strings.HasPrefix(arguments[index], "--baseline-json="):
			value, err := optionValue(arguments, &index, "--baseline-json")
			if err != nil {
				return usageResult(CommandBlockedCheck, err.Error())
			}
			baselineJSONPath = value
		case arguments[index] == "--baseline-failures" || strings.HasPrefix(arguments[index], "--baseline-failures="):
			value, err := optionValue(arguments, &index, "--baseline-failures")
			if err != nil {
				return usageResult(CommandBlockedCheck, err.Error())
			}
			baselineFailuresPath = value
		default:
			return usageResult(CommandBlockedCheck, "unknown option "+arguments[index])
		}
	}
	if probeFile == "" {
		return usageResult(CommandBlockedCheck, "--probe-file is required")
	}
	if (baselineJSONPath == "") != (baselineFailuresPath == "") {
		return usageResult(CommandBlockedCheck, "--baseline-json and --baseline-failures must be supplied together")
	}
	probeBytes, err := os.ReadFile(absoluteFromRoot(executionContext.RepositoryRoot, probeFile))
	if err != nil {
		return usageResult(CommandBlockedCheck, err.Error())
	}
	probeEvidence, runError := nextselection.RunBlockedProbeEvidenceAtRoot(executionContext.RepositoryRoot, probeBytes, timeoutSeconds)
	status := probeEvidence.ExitStatus
	code, outcome, severity := "BLOCKED-PROBE-SUCCEEDED", resultmodel.OutcomeSuccess, resultmodel.SeverityInfo
	if status != 0 {
		code, outcome, severity = "BLOCKED-PROBE-FAILED", resultmodel.OutcomeFindings, resultmodel.SeverityWarning
	}
	if status == 124 {
		code = "BLOCKED-PROBE-TIMED-OUT"
	}
	if status == 125 || runError != nil {
		code, outcome, severity = "BLOCKED-PROBE-LAUNCH-FAILED", resultmodel.OutcomeFailure, resultmodel.SeverityError
	}
	focusedTest := &resultmodel.FocusedTestResult{
		ProbeFile: probeFile, ExitStatus: status, Launched: probeEvidence.Launched,
		TimedOut: probeEvidence.TimedOut, Diagnostic: probeEvidence.Diagnostic,
		DiagnosticSHA256: probeEvidence.DiagnosticSHA256, BaselineState: resultmodel.FocusedBaselineNotCompared,
		CommandText: strings.TrimSpace(string(probeBytes)),
	}
	baselineFinding := resultmodel.CommandFinding{}
	if baselineJSONPath != "" {
		baselineFinding = compareFocusedBaseline(executionContext.RepositoryRoot, baselineJSONPath, baselineFailuresPath, focusedTest)
		if baselineFinding.Code != "" && baselineFinding.Severity == resultmodel.SeverityError && outcome == resultmodel.OutcomeSuccess {
			outcome = resultmodel.OutcomeFindings
		}
	}
	evidence := fmt.Sprintf("raw probe status %d; diagnostic sha256 %s", status, probeEvidence.DiagnosticSHA256)
	if runError != nil {
		evidence += ": " + runError.Error()
	}
	findings := []resultmodel.CommandFinding{helperFinding(code, severity, []string{probeFile}, evidence,
		resultmodel.FixabilityManual, map[bool]string{true: "the probe did not clear the blocked condition", false: ""}[status != 0],
		[]string{"do-work-cli", CommandBlockedCheck, "--probe-file", probeFile}, []string{"do-work-cli", "--format", "json", CommandBlockedCheck, "--probe-file", probeFile})}
	if baselineFinding.Code != "" {
		findings = append(findings, baselineFinding)
	}
	return resultmodel.CommandResult{Outcome: outcome, ExitCodeOverride: status, Findings: findings, FocusedTest: focusedTest}
}

func compareFocusedBaseline(repositoryRoot, baselineJSONPath, baselineFailuresPath string, focusedTest *resultmodel.FocusedTestResult) resultmodel.CommandFinding {
	baselineBytes, readError := os.ReadFile(absoluteFromRoot(repositoryRoot, baselineJSONPath))
	if os.IsNotExist(readError) {
		focusedTest.BaselineState = resultmodel.FocusedBaselineMissing
		return helperFinding("FOCUSED-BASELINE-MISSING", resultmodel.SeverityWarning, []string{baselineJSONPath}, "no baseline record exists", resultmodel.FixabilityManual, "every failing test remains a candidate regression", nil, nil)
	}
	if readError != nil {
		focusedTest.BaselineState = resultmodel.FocusedBaselineUnusable
		return helperFinding("FOCUSED-BASELINE-UNREADABLE", resultmodel.SeverityError, []string{baselineJSONPath}, readError.Error(), resultmodel.FixabilityManual, "baseline comparison is unavailable", nil, nil)
	}
	var baseline baselineRecord
	if decodeError := json.Unmarshal(baselineBytes, &baseline); decodeError != nil {
		focusedTest.BaselineState = resultmodel.FocusedBaselineUnusable
		return helperFinding("FOCUSED-BASELINE-INVALID", resultmodel.SeverityError, []string{baselineJSONPath}, decodeError.Error(), resultmodel.FixabilityManual, "baseline comparison is unavailable", nil, nil)
	}
	focusedTest.BaselineStatus = baseline.ExitStatus
	focusedTest.BaselineLaunched = baseline.Launched
	if !baseline.Launched {
		focusedTest.BaselineState = resultmodel.FocusedBaselineUnusable
		return helperFinding("FOCUSED-BASELINE-NOT-LAUNCHED", resultmodel.SeverityError, []string{baselineJSONPath}, "baseline launched=false", resultmodel.FixabilityManual, "a record for a command that never ran cannot exclude failures", nil, nil)
	}
	if focusedTest.ExitStatus == 0 {
		focusedTest.BaselineState = resultmodel.FocusedBaselineGreen
		return resultmodel.CommandFinding{}
	}
	baselineFailureBytes, failureError := os.ReadFile(absoluteFromRoot(repositoryRoot, baselineFailuresPath))
	if failureError != nil {
		focusedTest.BaselineState = resultmodel.FocusedBaselineNewRed
		return helperFinding("FOCUSED-BASELINE-FAILURE-EVIDENCE-MISSING", resultmodel.SeverityError, []string{baselineFailuresPath}, failureError.Error(), resultmodel.FixabilityManual, "the current failure cannot be excluded", nil, nil)
	}
	_, baselineDiagnosticSHA256 := nextselection.BlockedProbeDiagnosticIdentity(string(baselineFailureBytes), repositoryRoot)
	if baseline.ExitStatus != 0 && baseline.ExitStatus == focusedTest.ExitStatus && strings.TrimSpace(baseline.TestCommand) == focusedTest.CommandText && baselineDiagnosticSHA256 == focusedTest.DiagnosticSHA256 {
		focusedTest.BaselineState = resultmodel.FocusedBaselineMatchingRed
		return helperFinding("FOCUSED-BASELINE-MATCH", resultmodel.SeverityInfo, []string{baselineFailuresPath}, "same command, status, and normalized diagnostic identity", resultmodel.FixabilityAutomatic, "", nil, nil)
	}
	focusedTest.BaselineState = resultmodel.FocusedBaselineNewRed
	return helperFinding("FOCUSED-NEW-RED", resultmodel.SeverityError, []string{probeFileOrFallback(focusedTest.ProbeFile)}, "current failure does not match the saved command and normalized diagnostic identity", resultmodel.FixabilityManual, "the failure is a candidate regression", nil, nil)
}

func probeFileOrFallback(probeFile string) string {
	if probeFile == "" {
		return "<focused-test-probe>"
	}
	return probeFile
}

func handleRepairTimestamps(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	filtered, dryRun, err := extractDryRun(arguments)
	if err != nil || len(filtered) != 0 {
		return usageResult(CommandRepairTimestamps, "repair-req-timestamps accepts no options")
	}
	return runTimestampCommand(executionContext.RepositoryRoot, doctor.TimestampScopeActive, true, dryRun)
}

func handleAuditTimestamps(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	filtered, dryRun, err := extractDryRun(arguments)
	if err != nil {
		return usageResult(CommandAuditTimestamps, err.Error())
	}
	applyFixes := false
	for _, argument := range filtered {
		if argument == "--fix" {
			applyFixes = true
		} else {
			return usageResult(CommandAuditTimestamps, "unknown option "+argument)
		}
	}
	return runTimestampCommand(executionContext.RepositoryRoot, doctor.TimestampScopeArchive, applyFixes, dryRun)
}

func runTimestampCommand(repositoryRoot string, scope doctor.TimestampScope, apply, dryRun bool) resultmodel.CommandResult {
	commandName := CommandRepairTimestamps
	if scope == doctor.TimestampScopeArchive {
		commandName = CommandAuditTimestamps
		if walkFailure := archiveWalkFailure(repositoryRoot); walkFailure != "" {
			output := "audit-archive-timestamps: the archive walk failed — nothing was inspected.\n"
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, RepositoryRoot: repositoryRoot, ExactTextOutput: &output,
				Findings: []resultmodel.CommandFinding{helperFinding("TIMESTAMP-ARCHIVE-WALK-FAILED", resultmodel.SeverityError, []string{"do-work/archive"}, walkFailure, resultmodel.FixabilityManual, "the archive was not inspected", nil, nil)}}
		}
	}
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		return usageResult(commandName, err.Error())
	}
	plans, findings := doctor.BuildTimestampPlanForScope(context.Background(), snapshot, time.Now().UTC(), scope)
	if !apply {
		result := successResult(nil, findings)
		for _, plan := range plans {
			result.Findings = append(result.Findings, helperFinding("TIMESTAMP-REPAIR-PENDING", resultmodel.SeverityWarning, []string{plan.RelativePath}, "timestamp repair is pending", resultmodel.FixabilityAutomatic, "report-only audit does not mutate", []string{"do-work-cli", CommandAuditTimestamps, "--fix"}, []string{"git", "diff", "--", plan.RelativePath}))
		}
		if len(plans) > 0 {
			result.Outcome = resultmodel.OutcomeFindings
		}
		if !dryRun {
			output := timestampCompatibilityText(scope, false, plans, findings, resultmodel.CommandResult{}, archiveRequestCount(snapshot))
			result.ExactTextOutput = &output
			if scope == doctor.TimestampScopeArchive && len(plans) == 0 {
				result.Outcome = resultmodel.OutcomeSuccess
			}
		}
		return result
	}
	if dryRun {
		result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: repositoryRoot, Findings: findings}
		for _, plan := range plans {
			for _, change := range plan.Changes {
				result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: plan.RelativePath, Kind: "modified", Detail: fmt.Sprintf("would repair %s line %d: %s -> %s (%s)", change.FieldName, change.LineNumber, change.OldValue, change.NewValue, change.Source)})
			}
		}
		if len(findings) > 0 {
			result.Outcome = resultmodel.OutcomeFindings
		}
		return result
	}
	result := doctor.ApplyUncommittedTimestampPlans(snapshot, plans)
	result.Findings = append(findings, result.Findings...)
	output := timestampCompatibilityText(scope, true, plans, findings, result, archiveRequestCount(snapshot))
	result.ExactTextOutput = &output
	if len(result.Findings) == len(findings) && !hasActionableTimestampRefusal(findings) {
		result.Outcome = resultmodel.OutcomeSuccess
	} else if hasActionableTimestampRefusal(findings) && result.Outcome == resultmodel.OutcomeSuccess {
		result.Outcome = resultmodel.OutcomeFindings
	}
	return result
}

func archiveWalkFailure(repositoryRoot string) string {
	archiveRoot := filepath.Join(repositoryRoot, "do-work", "archive")
	info, err := os.Stat(archiveRoot)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil || !info.IsDir() {
		return fmt.Sprint(err)
	}
	command := exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")
	if output, err := command.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(output) + " " + err.Error())
	}
	return ""
}

func archiveRequestCount(snapshot *repositorymodel.RepositorySnapshot) int {
	count := 0
	for _, requestFile := range snapshot.RequestFiles {
		if strings.HasPrefix("do-work/"+requestFile.RelativePath, "do-work/archive/") {
			count++
		}
	}
	return count
}

func timestampCompatibilityText(scope doctor.TimestampScope, apply bool, plans []doctor.TimestampRepairPlan, findings []resultmodel.CommandFinding, applied resultmodel.CommandResult, scanned int) string {
	var output strings.Builder
	verb := "would repair"
	if apply {
		verb = "repaired"
	}
	for _, plan := range plans {
		for _, change := range plan.Changes {
			if apply && !recordedTimestampChange(applied.Changes, plan.RelativePath, change) {
				continue
			}
			fmt.Fprintf(&output, "do-work: %s %s %s: %s -> %s (%s)\n", verb, plan.RelativePath, change.FieldName, change.OldValue, change.NewValue, change.Source)
		}
	}
	refusalCount, actionableRefusalCount := 0, 0
	for _, finding := range findings {
		if finding.Code != "TIMESTAMP-REPAIR-REFUSED" {
			continue
		}
		refusalCount++
		path := "request"
		if len(finding.AffectedPaths) > 0 {
			path = finding.AffectedPaths[0]
		}
		evidence := strings.Join(finding.Evidence, "; ")
		if strings.Contains(evidence, "git blame") || strings.Contains(evidence, "could not derive") {
			actionableRefusalCount++
			fmt.Fprintf(&output, "do-work: FAILED to repair %s: %s\n", path, evidence)
		} else if scope == doctor.TimestampScopeArchive {
			fmt.Fprintf(&output, "do-work: %s refused: %s\n", path, evidence)
		}
	}
	for _, finding := range applied.Findings {
		if !strings.HasPrefix(finding.Code, "TIMESTAMP-PREIMAGE-") && finding.Code != "TIMESTAMP-REPLACE-FAILED" && finding.Code != "TIMESTAMP-TRUNCATION-FLOOR-UNAVAILABLE" && finding.Code != "TIMESTAMP-CONTENT-LOSS" && !strings.HasPrefix(finding.Code, "TIMESTAMP-POST-WRITE-") {
			continue
		}
		path := "request"
		if len(finding.AffectedPaths) > 0 {
			path = finding.AffectedPaths[0]
		}
		fmt.Fprintf(&output, "do-work: FAILED to repair %s: %s\n", path, strings.Join(finding.Evidence, "; "))
	}
	if scope != doctor.TimestampScopeArchive {
		return output.String()
	}
	if !apply && len(plans) > 0 {
		fmt.Fprintf(&output, "do-work: %d archived correction(s) pending — rerun with --fix to write them.\n", len(plans))
	} else if apply && len(applied.Changes) > 0 {
		fmt.Fprintf(&output, "do-work: repaired %d archived timestamp(s) — review and commit the correction(s) through the normal flow.\n", len(applied.Changes))
	} else if refusalCount > 0 && actionableRefusalCount == 0 {
		fmt.Fprintf(&output, "do-work: archive audit complete (%d file(s) scanned) — %d value(s) refused and left byte-identical, listed above. Not clean: those values were never inspected for defects.\n", scanned, refusalCount)
	} else if refusalCount == 0 {
		fmt.Fprintf(&output, "do-work: archive audit clean (%d file(s) scanned).\n", scanned)
	}
	return output.String()
}

func hasActionableTimestampRefusal(findings []resultmodel.CommandFinding) bool {
	for _, finding := range findings {
		if finding.Code != "TIMESTAMP-REPAIR-REFUSED" {
			continue
		}
		evidence := strings.Join(finding.Evidence, "; ")
		if strings.Contains(evidence, "git blame") || strings.Contains(evidence, "could not derive") {
			return true
		}
	}
	return false
}

func recordedTimestampChange(changes []resultmodel.RecordedChange, path string, change doctor.TimestampFieldChange) bool {
	needle := "repaired " + change.FieldName + " line "
	for _, recorded := range changes {
		if recorded.Path == path && strings.Contains(recorded.Detail, needle) {
			return true
		}
	}
	return false
}

func handleAtomicDownload(_ commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	filtered, dryRun, dryRunError := extractDryRun(arguments)
	if dryRunError != nil {
		return usageResult(CommandAtomicDownload, dryRunError.Error())
	}
	arguments = filtered
	sourceURL, targetPath := "", ""
	for index := 0; index < len(arguments); index++ {
		switch {
		case arguments[index] == "--source-url" || strings.HasPrefix(arguments[index], "--source-url="):
			value, err := optionValue(arguments, &index, "--source-url")
			if err != nil {
				return usageResult(CommandAtomicDownload, err.Error())
			}
			sourceURL = value
		case arguments[index] == "--target-path" || strings.HasPrefix(arguments[index], "--target-path="):
			value, err := optionValue(arguments, &index, "--target-path")
			if err != nil {
				return usageResult(CommandAtomicDownload, err.Error())
			}
			targetPath = value
		default:
			return usageResult(CommandAtomicDownload, "unknown option "+arguments[index])
		}
	}
	if sourceURL == "" || targetPath == "" {
		return usageResult(CommandAtomicDownload, "--source-url and --target-path are required")
	}
	parsedURL, parseURLError := url.ParseRequestURI(sourceURL)
	if parseURLError != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return usageResult(CommandAtomicDownload, "--source-url must be an absolute HTTP(S) URL")
	}
	if dryRun {
		if _, err := os.Lstat(targetPath); err == nil {
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{helperFinding("DOWNLOAD-TARGET-OCCUPIED", resultmodel.SeverityError, []string{targetPath}, "target already exists", resultmodel.FixabilityManual, "dry-run preserves the occupied target", nil, []string{"test", "-e", targetPath})}}
		}
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Changes: []resultmodel.RecordedChange{{Path: targetPath, Kind: "created", Detail: "would fetch and publish private HTTP bytes after transfer validation"}}, Findings: []resultmodel.CommandFinding{helperFinding("DOWNLOAD-DRY-RUN", resultmodel.SeverityInfo, []string{targetPath}, "request and target arguments validated; network and filesystem were not mutated", resultmodel.FixabilityAutomatic, "", nil, []string{"test", "!", "-e", targetPath})}}
	}
	if info, err := os.Lstat(targetPath); err == nil && info.IsDir() {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("DOWNLOAD-TARGET-OCCUPIED", resultmodel.SeverityError, []string{targetPath}, "target is a directory", resultmodel.FixabilityRefused, "existing directory left unchanged", nil, []string{"test", "-d", targetPath})}}
	}
	temporary, err := os.CreateTemp(filepath.Dir(targetPath), filepath.Base(targetPath)+".download.*")
	if err != nil {
		return usageResult(CommandAtomicDownload, err.Error())
	}
	temporaryPath := temporary.Name()
	_ = temporary.Close()
	defer os.Remove(temporaryPath)
	curlArguments := []string{"-fsSL", "--retry", "3", "--retry-delay", "2", "--retry-max-time", "60"}
	if token := firstNonempty(os.Getenv("GH_TOKEN"), os.Getenv("GITHUB_TOKEN")); token != "" {
		curlArguments = append(curlArguments, "-H", "Authorization: Bearer "+token)
	}
	curlArguments = append(curlArguments, "-o", temporaryPath, sourceURL)
	command := exec.Command("curl", curlArguments...)
	if output, runError := command.CombinedOutput(); runError != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{helperFinding("DOWNLOAD-FAILED", resultmodel.SeverityError, []string{targetPath}, strings.TrimSpace(string(output))+": "+runError.Error(), resultmodel.FixabilityManual, "no target was published", []string{"do-work-cli", CommandAtomicDownload, "--source-url", sourceURL, "--target-path", targetPath}, []string{"test", "!", "-e", targetPath})}}
	}
	if err := os.Rename(temporaryPath, targetPath); err != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{helperFinding("DOWNLOAD-FAILED", resultmodel.SeverityError, []string{targetPath}, err.Error(), resultmodel.FixabilityManual, "downloaded bytes were discarded", nil, []string{"test", "!", "-e", targetPath})}}
	}
	info, _ := os.Stat(targetPath)
	return successResult([]resultmodel.RecordedChange{{Path: targetPath, Kind: "created", Detail: fmt.Sprintf("published %d bytes with mode 0600", info.Size())}}, nil)
}

func firstNonempty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func absoluteFromRoot(repositoryRoot, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(repositoryRoot, filepath.FromSlash(path))
}
