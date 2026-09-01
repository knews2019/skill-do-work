package knowledgecommands

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const memoryCharacterLimit = 2500

var credentialPattern = regexp.MustCompile(`(?i)(ghp_[A-Za-z0-9_]{12,}|github_pat_[A-Za-z0-9_]{12,}|sk-[A-Za-z0-9_-]{12,}|AKIA[A-Z0-9]{12,}|xox[a-z]-[A-Za-z0-9-]{8,}|eyJ[A-Za-z0-9_-]{12,}|Bearer\s+[A-Za-z0-9._~+/-]{8,}|(?:password|passwd|secret|token|api[_-]?key)\s*[:=]\s*\S+)`)

type memoryOptions struct {
	memoryRoot string
	section    string
	payload    string
	matches    []string
	demoteIDs  []string
	replaceID  string
	manifest   string
	engine     string
	kb         string
	dryRun     bool
	commit     bool
	confirm    bool
}

type memoryMatch struct {
	ID      string
	Path    string
	Line    int
	Content string
	Working bool
}

func parseMemoryOptions(arguments []string, mutable bool) (memoryOptions, error) {
	options := memoryOptions{memoryRoot: "memory", engine: "both", kb: "kb"}
	payload := []string{}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		value := func(name string) (string, error) {
			index++
			if index >= len(arguments) {
				return "", fmt.Errorf("%s requires a value", name)
			}
			return arguments[index], nil
		}
		var err error
		switch {
		case argument == "--memory-root" || argument == "--path":
			options.memoryRoot, err = value(argument)
		case strings.HasPrefix(argument, "--memory-root="):
			options.memoryRoot = strings.TrimPrefix(argument, "--memory-root=")
		case argument == "--section":
			options.section, err = value(argument)
		case strings.HasPrefix(argument, "--section="):
			options.section = strings.TrimPrefix(argument, "--section=")
		case argument == "--match":
			var item string
			item, err = value(argument)
			options.matches = append(options.matches, item)
		case strings.HasPrefix(argument, "--match="):
			options.matches = append(options.matches, strings.TrimPrefix(argument, "--match="))
		case argument == "--replace" || argument == "--replace-id":
			options.replaceID, err = value(argument)
		case argument == "--demote":
			var item string
			item, err = value(argument)
			options.demoteIDs = append(options.demoteIDs, item)
		case strings.HasPrefix(argument, "--demote="):
			options.demoteIDs = append(options.demoteIDs, strings.TrimPrefix(argument, "--demote="))
		case argument == "--manifest":
			options.manifest, err = value(argument)
		case strings.HasPrefix(argument, "--manifest="):
			options.manifest = strings.TrimPrefix(argument, "--manifest=")
		case argument == "--engine":
			options.engine, err = value(argument)
		case strings.HasPrefix(argument, "--engine="):
			options.engine = strings.TrimPrefix(argument, "--engine=")
		case argument == "--kb":
			options.kb, err = value(argument)
		case strings.HasPrefix(argument, "--kb="):
			options.kb = strings.TrimPrefix(argument, "--kb=")
		case (argument == "--text" || argument == "--query"):
			var item string
			item, err = value(argument)
			payload = append(payload, item)
		case argument == "--dry-run" && mutable:
			options.dryRun = true
		case argument == "--commit" && mutable:
			options.commit = true
		case argument == "--confirm" && mutable:
			options.confirm = true
		case strings.HasPrefix(argument, "-"):
			return options, fmt.Errorf("unknown option %q", argument)
		default:
			payload = append(payload, argument)
		}
		if err != nil {
			return options, err
		}
	}
	options.payload = strings.TrimSpace(strings.Join(payload, " "))
	options.section = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(options.section), "_", "-"))
	if options.dryRun && options.commit {
		return options, errors.New("--dry-run and --commit cannot be combined")
	}
	if options.memoryRoot == "" {
		return options, errors.New("memory root must not be empty")
	}
	return options, nil
}

func handleMemoryRemember(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	operationTime := nowUTC()
	options, err := parseMemoryOptions(arguments, true)
	if err != nil {
		return usageResult(CommandMemoryRemember, err)
	}
	if options.payload == "" {
		return usageResult(CommandMemoryRemember, errors.New("memory text is required"))
	}
	if strings.ContainsAny(options.payload, "\r\n\x00") {
		return memoryFindingResult(CommandMemoryRemember, "MEMORY-TEXT-MALFORMED", resultmodel.SeverityError, options.memoryRoot, "memory text must be one NUL-free line", resultmodel.OutcomeRefused)
	}
	if credentialPattern.MatchString(options.payload) {
		return memoryFindingResult(CommandMemoryRemember, "MEMORY-CREDENTIAL-REFUSED", resultmodel.SeverityError, options.memoryRoot, "credential-shaped text was refused before persistence", resultmodel.OutcomeRefused)
	}
	memoryAbsolute, memoryRelative, resolveError := resolveMemoryRoot(executionContext.RepositoryRoot, options.memoryRoot, true)
	if resolveError != nil {
		return memoryFailure(CommandMemoryRemember, "MEMORY-ROOT-INVALID", options.memoryRoot, resolveError)
	}
	workingPath := filepath.ToSlash(filepath.Join(memoryRelative, "working-memory.md"))
	workingBytes, readError := os.ReadFile(filepath.Join(memoryAbsolute, "working-memory.md"))
	if readError != nil {
		return memoryFailure(CommandMemoryRemember, "MEMORY-NOT-SET-UP", workingPath, readError)
	}
	lines := strings.Split(strings.ReplaceAll(string(workingBytes), "\r\n", "\n"), "\n")
	originalLines := append([]string(nil), lines...)
	normalizedPayload := normalizeMemoryText(options.payload)
	for index, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "-") && normalizeMemoryText(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "-"))) == normalizedPayload {
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, SkippedWork: []resultmodel.SkippedWork{{Code: "MEMORY-EXACT-DUPLICATE", Reason: fmt.Sprintf("the fact already exists at %s:%d", workingPath, index+1)}}}
		}
	}
	sectionHeading := map[string]string{"active-threads": "## Active Threads", "notes": "## Notes", "pending-decisions": "## Pending Decisions"}[options.section]
	if sectionHeading == "" {
		return memoryFindingResult(CommandMemoryRemember, "MEMORY-SECTION-REQUIRED", resultmodel.SeverityWarning, workingPath, "choose --section active-threads|notes|pending-decisions; semantic placement remains caller-owned", resultmodel.OutcomeRefused)
	}
	demoteSet := map[string]bool{}
	for _, id := range options.demoteIDs {
		demoteSet[id] = true
	}
	if options.replaceID != "" && demoteSet[options.replaceID] {
		return memoryFindingResult(CommandMemoryRemember, "MEMORY-PLAN-CONFLICT", resultmodel.SeverityWarning, workingPath, "one content-bound entry cannot be replaced and demoted in the same plan", resultmodel.OutcomeRefused)
	}
	if options.replaceID != "" {
		replaced := false
		for index, line := range originalLines {
			if strings.HasPrefix(strings.TrimSpace(line), "-") && memoryLineID(workingPath, index+1, line) == options.replaceID {
				lines[index] = "- " + options.payload
				replaced = true
				break
			}
		}
		if !replaced {
			return memoryFindingResult(CommandMemoryRemember, "MEMORY-REPLACE-STALE", resultmodel.SeverityWarning, workingPath, "the content-bound replacement ID no longer matches", resultmodel.OutcomeRefused)
		}
	}
	demoted := []string{}
	if len(demoteSet) > 0 {
		kept := make([]string, 0, len(lines))
		for index, line := range lines {
			id := memoryLineID(workingPath, index+1, originalLines[index])
			if demoteSet[id] && strings.HasPrefix(strings.TrimSpace(line), "-") {
				demoted = append(demoted, strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "-")))
				delete(demoteSet, id)
				continue
			}
			kept = append(kept, line)
		}
		if len(demoteSet) > 0 {
			return memoryFindingResult(CommandMemoryRemember, "MEMORY-DEMOTION-STALE", resultmodel.SeverityWarning, workingPath, "a content-bound demotion ID no longer matches", resultmodel.OutcomeRefused)
		}
		lines = kept
	}
	if options.replaceID == "" {
		inserted := false
		for index, line := range lines {
			if strings.TrimSpace(line) == sectionHeading {
				insertAt := index + 1
				for insertAt < len(lines) && strings.TrimSpace(lines[insertAt]) == "" {
					insertAt++
				}
				updated := append([]string{}, lines[:insertAt]...)
				updated = append(updated, "- "+options.payload)
				updated = append(updated, lines[insertAt:]...)
				lines = updated
				inserted = true
				break
			}
		}
		if !inserted {
			return memoryFailure(CommandMemoryRemember, "MEMORY-SECTION-MISSING", workingPath, fmt.Errorf("working memory is missing %s", sectionHeading))
		}
	}
	for index, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "updated:") {
			lines[index] = "updated: " + operationTime.Format("2006-01-02")
			break
		}
	}
	newWorking := []byte(strings.Join(lines, "\n"))
	if len(newWorking) > memoryCharacterLimit {
		return memoryFindingResult(CommandMemoryRemember, "MEMORY-CAP-PLAN-REQUIRED", resultmodel.SeverityWarning, workingPath, fmt.Sprintf("planned working memory is %d bytes; explicit consolidation/demotion is required to stay at %d", len(newWorking), memoryCharacterLimit), resultmodel.OutcomeRefused)
	}
	logRelative := filepath.ToSlash(filepath.Join(memoryRelative, "logs", operationTime.Format("2006-01-02")+".md"))
	logAbsolute := filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(logRelative))
	logBytes, logError := os.ReadFile(logAbsolute)
	if logError != nil && !os.IsNotExist(logError) {
		return memoryFailure(CommandMemoryRemember, "MEMORY-LOG-READ-FAILED", logRelative, logError)
	}
	if len(logBytes) > 0 && logBytes[len(logBytes)-1] != '\n' {
		logBytes = append(logBytes, '\n')
	}
	for _, demotedText := range demoted {
		logBytes = append(logBytes, []byte(fmt.Sprintf("\n## %s UTC note\n%s\n", operationTime.Format("15:04"), demotedText))...)
	}
	logBytes = append(logBytes, []byte(fmt.Sprintf("\n## %s UTC note\n%s\n", operationTime.Format("15:04"), options.payload))...)
	targets := []string{workingPath, logRelative}
	createdDirs := absentDirectories(executionContext.RepositoryRoot, []string{filepath.ToSlash(filepath.Join(memoryRelative, "logs"))})
	transaction := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{RepositoryRoot: executionContext.RepositoryRoot, TargetPaths: targets, PrivateUntrackedTargetPaths: []string{logRelative}, CreatedDirectoryPaths: createdDirs, DryRun: options.dryRun, Commit: options.commit, CommitMessage: "Remember curated project fact"}, func(recorder *gittransaction.MutationRecorder) error {
		if err := createTransactionDirectories(executionContext.RepositoryRoot, recorder, createdDirs); err != nil {
			return err
		}
		if err := publishTransactionFile(executionContext.RepositoryRoot, recorder, workingPath, newWorking, 0o644, false); err != nil {
			return err
		}
		return publishTransactionFile(executionContext.RepositoryRoot, recorder, logRelative, logBytes, 0o644, true)
	})
	result := gittransaction.BuildCommandResult(CommandMemoryRemember, transaction)
	if options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		result.Changes = memoryPlannedChanges(targets)
	}
	if !options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		appendMemoryLedger(memoryAbsolute, "write", normalizedPayload, 1, "remember")
	}
	return result
}

func handleMemoryForget(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	operationTime := nowUTC()
	options, err := parseMemoryOptions(arguments, true)
	if err != nil {
		return usageResult(CommandMemoryForget, err)
	}
	if options.payload == "" {
		return usageResult(CommandMemoryForget, errors.New("forget query is required"))
	}
	memoryAbsolute, memoryRelative, resolveError := resolveMemoryRoot(executionContext.RepositoryRoot, options.memoryRoot, true)
	if resolveError != nil {
		return memoryFailure(CommandMemoryForget, "MEMORY-ROOT-INVALID", options.memoryRoot, resolveError)
	}
	matches, scanError := findMemoryMatches(memoryAbsolute, memoryRelative, options.payload)
	if scanError != nil {
		return memoryFailure(CommandMemoryForget, "MEMORY-FORGET-SCAN-FAILED", memoryRelative, scanError)
	}
	if len(matches) == 0 {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, SkippedWork: []resultmodel.SkippedWork{{Code: "MEMORY-NO-MATCH", Reason: "no standing-memory bullet or daily-log line matched"}}}
	}
	if len(options.matches) == 0 || !options.confirm {
		findings := make([]resultmodel.CommandFinding, 0, len(matches))
		for _, match := range matches {
			finding := memoryFinding(CommandMemoryForget, "MEMORY-FORGET-MATCH", resultmodel.SeverityWarning, match.Path, fmt.Sprintf("line=%d content=%s", match.Line, match.Content), resultmodel.FixabilityManual, "the user must confirm the exact content-bound match ID")
			finding.AffectedIDs = []string{match.ID}
			finding.NextArgv = []string{"do-work-cli", CommandMemoryForget, "--confirm", "--match", match.ID, "--memory-root", options.memoryRoot, options.payload}
			finding.VerificationArgv = []string{"do-work-cli", "--format", "json", CommandMemoryForget, "--memory-root", options.memoryRoot, options.payload}
			finding.NextJustRecipe = CommandMemoryForget + " " + quoteRecipeArgument(options.payload) + " --confirm --match " + quoteRecipeArgument(match.ID) + " --memory-root " + quoteRecipeArgument(options.memoryRoot)
			findings = append(findings, finding)
		}
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: findings}
	}
	byID := map[string]memoryMatch{}
	for _, match := range matches {
		byID[match.ID] = match
	}
	selected := []memoryMatch{}
	selectedIDs := map[string]bool{}
	for _, id := range options.matches {
		if selectedIDs[id] {
			return memoryFindingResult(CommandMemoryForget, "MEMORY-FORGET-SELECTION-DUPLICATE", resultmodel.SeverityWarning, memoryRelative, "a confirmed match ID was supplied more than once: "+id, resultmodel.OutcomeRefused)
		}
		selectedIDs[id] = true
		match, exists := byID[id]
		if !exists {
			return memoryFindingResult(CommandMemoryForget, "MEMORY-FORGET-SELECTION-STALE", resultmodel.SeverityWarning, memoryRelative, "a confirmed content-bound match no longer exists: "+id, resultmodel.OutcomeRefused)
		}
		selected = append(selected, match)
	}
	selectedByPath := map[string]map[int]memoryMatch{}
	for _, match := range selected {
		if selectedByPath[match.Path] == nil {
			selectedByPath[match.Path] = map[int]memoryMatch{}
		}
		selectedByPath[match.Path][match.Line] = match
	}
	writes := map[string][]byte{}
	private := []string{}
	for path, lineMatches := range selectedByPath {
		data, readErr := os.ReadFile(filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(path)))
		if readErr != nil {
			return memoryFailure(CommandMemoryForget, "MEMORY-FORGET-SELECTION-STALE", path, readErr)
		}
		lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
		updated := []string{}
		for index, line := range lines {
			match, chosen := lineMatches[index+1]
			if !chosen {
				updated = append(updated, line)
				continue
			}
			if memoryLineID(path, index+1, line) != match.ID {
				return memoryFindingResult(CommandMemoryForget, "MEMORY-FORGET-SELECTION-STALE", resultmodel.SeverityWarning, path, "selected line changed before mutation", resultmodel.OutcomeRefused)
			}
			if strings.HasPrefix(strings.TrimSpace(line), "##") {
				return memoryFindingResult(CommandMemoryForget, "MEMORY-FORGET-HEADING-REFUSED", resultmodel.SeverityError, path, "heading lines are never forget targets", resultmodel.OutcomeRefused)
			}
			if match.Working {
				continue
			}
			prefix := ""
			if strings.HasPrefix(strings.TrimSpace(line), ">") {
				leading := line[:strings.Index(line, ">")+1]
				prefix = leading + " "
			}
			updated = append(updated, prefix+"[forgotten — redacted by memory forget "+operationTime.Format("2006-01-02")+"]")
		}
		if path == filepath.ToSlash(filepath.Join(memoryRelative, "working-memory.md")) {
			for index, line := range updated {
				if strings.HasPrefix(strings.TrimSpace(line), "updated:") {
					updated[index] = "updated: " + operationTime.Format("2006-01-02")
					break
				}
			}
		}
		writes[path] = []byte(strings.Join(updated, "\n"))
		if path != filepath.ToSlash(filepath.Join(memoryRelative, "working-memory.md")) {
			private = append(private, path)
		}
	}
	targets := uniqueSorted(mapKeysBytes(writes))
	if options.commit && !containsMemoryString(targets, filepath.ToSlash(filepath.Join(memoryRelative, "working-memory.md"))) {
		return memoryFindingResult(CommandMemoryForget, "MEMORY-COMMIT-NO-TRACKED-TARGET", resultmodel.SeverityWarning, memoryRelative, "confirmed matches are private-only; there is no committable target", resultmodel.OutcomeRefused)
	}
	transaction := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{RepositoryRoot: executionContext.RepositoryRoot, TargetPaths: targets, PrivateUntrackedTargetPaths: private, DryRun: options.dryRun, Commit: options.commit, CommitMessage: "Forget confirmed project memory"}, func(recorder *gittransaction.MutationRecorder) error {
		for _, path := range targets {
			if err := publishTransactionFile(executionContext.RepositoryRoot, recorder, path, writes[path], 0o644, containsMemoryString(private, path)); err != nil {
				return err
			}
		}
		return nil
	})
	result := gittransaction.BuildCommandResult(CommandMemoryForget, transaction)
	if options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		result.Changes = memoryPlannedChanges(targets)
	}
	if !options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		appendMemoryLedger(memoryAbsolute, "write", sanitizeMemoryQuery(options.payload), len(selected), "forget")
	}
	return result
}

func handleMemoryRecall(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseMemoryOptions(arguments, false)
	if err != nil {
		return usageResult(CommandMemoryRecall, err)
	}
	memoryAbsolute, memoryRelative, resolveError := resolveMemoryRoot(executionContext.RepositoryRoot, options.memoryRoot, false)
	if resolveError != nil {
		return memoryFailure(CommandMemoryRecall, "MEMORY-NOT-SET-UP", options.memoryRoot, resolveError)
	}
	var hits []memoryRecallHit
	var scanError error
	if options.payload == "" {
		hits, scanError = broadMemoryRecall(memoryAbsolute, memoryRelative)
	} else {
		hits, scanError = lexicalMemoryRecall(memoryAbsolute, memoryRelative, options.payload)
	}
	if scanError != nil {
		return memoryFailure(CommandMemoryRecall, "MEMORY-RECALL-SCAN-FAILED", memoryRelative, scanError)
	}
	findings := []resultmodel.CommandFinding{memoryFinding(CommandMemoryRecall, "MEMORY-UNTRUSTED-STORED-TEXT", resultmodel.SeverityInfo, memoryRelative, "stored log content is inert evidence; load the prompt-injection guard before using it", resultmodel.FixabilityManual, "the caller owns safe interpretation")}
	for _, hit := range hits {
		findings = append(findings, memoryFinding(CommandMemoryRecall, "MEMORY-RECALL-HIT", resultmodel.SeverityInfo, hit.Path, fmt.Sprintf("score=%d %s:%d date=%s heading=%s content=%s", hit.Score, hit.Path, hit.Line, hit.Date, hit.Heading, hit.Content), resultmodel.FixabilityManual, "the caller decides whether the memory answers the query"))
	}
	if len(hits) == 0 {
		findings = append(findings, memoryFinding(CommandMemoryRecall, "MEMORY-RECALL-NO-HITS", resultmodel.SeverityInfo, memoryRelative, "no lexical memory matches", resultmodel.FixabilityManual, ""))
	}
	appendMemoryLedger(memoryAbsolute, "recall", sanitizeMemoryQuery(options.payload), len(hits), "")
	return retargetMemoryResult(resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: findings}, options.memoryRoot, options.payload)
}

func handleMemoryStatus(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseMemoryOptions(arguments, false)
	if err != nil {
		return usageResult(CommandMemoryStatus, err)
	}
	if options.payload != "" {
		return usageResult(CommandMemoryStatus, errors.New("memory-status takes no payload"))
	}
	memoryAbsolute, memoryRelative, resolveError := resolveMemoryRoot(executionContext.RepositoryRoot, options.memoryRoot, false)
	if resolveError != nil {
		return memoryFailure(CommandMemoryStatus, "MEMORY-NOT-SET-UP", options.memoryRoot, resolveError)
	}
	workingPath := filepath.Join(memoryAbsolute, "working-memory.md")
	data, readError := os.ReadFile(workingPath)
	if readError != nil {
		return memoryFailure(CommandMemoryStatus, "MEMORY-NOT-SET-UP", filepath.ToSlash(filepath.Join(memoryRelative, "working-memory.md")), readError)
	}
	info, _ := os.Stat(workingPath)
	front, _, _ := parseFrontmatter(string(data))
	logs, _ := filepath.Glob(filepath.Join(memoryAbsolute, "logs", "*.md"))
	sort.Strings(logs)
	newest := "none"
	lastCapture := "never"
	if len(logs) > 0 {
		newest = strings.TrimSuffix(filepath.Base(logs[len(logs)-1]), ".md")
		bytes, _ := os.ReadFile(logs[len(logs)-1])
		for _, line := range strings.Split(string(bytes), "\n") {
			if strings.HasPrefix(line, "## ") && strings.Contains(line, " UTC session capture ") {
				lastCapture = newest + " " + strings.TrimPrefix(line, "## ")
			}
		}
	}
	ledgerSummary := memoryLedgerTail(filepath.Join(memoryAbsolute, "usage-ledger.jsonl"), 5)
	evidence := fmt.Sprintf("bytes=%d cap=%d updated=%s mtime=%s log_days=%d newest=%s last_capture=%s ledger_tail=%s", len(data), memoryCharacterLimit, front["updated"], info.ModTime().UTC().Format(time.RFC3339), len(logs), newest, lastCapture, ledgerSummary)
	outcome := resultmodel.OutcomeSuccess
	findings := []resultmodel.CommandFinding{memoryFinding(CommandMemoryStatus, "MEMORY-STATUS", resultmodel.SeverityInfo, memoryRelative, evidence, resultmodel.FixabilityManual, "")}
	if len(data) > memoryCharacterLimit {
		outcome = resultmodel.OutcomeFindings
		findings = append(findings, memoryFinding(CommandMemoryStatus, "MEMORY-CAP-EXCEEDED", resultmodel.SeverityWarning, filepath.ToSlash(filepath.Join(memoryRelative, "working-memory.md")), fmt.Sprintf("working memory is %d bytes", len(data)), resultmodel.FixabilityManual, "consolidation requires semantic judgment"))
	}
	return retargetMemoryResult(resultmodel.CommandResult{Outcome: outcome, Findings: findings}, options.memoryRoot, options.payload)
}

type bootstrapEntry struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Source  string `json:"source"`
	Summary string `json:"summary"`
}

func handleMemoryBootstrap(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	operationTime := nowUTC()
	options, err := parseMemoryOptions(arguments, true)
	if err != nil {
		return usageResult(CommandMemoryBootstrap, err)
	}
	if options.commit {
		return memoryFindingResult(CommandMemoryBootstrap, "MEMORY-BOOTSTRAP-COMMIT-REFUSED", resultmodel.SeverityWarning, options.memoryRoot, "bootstrap logs and sentinel are machine-local and must never be committed", resultmodel.OutcomeRefused)
	}
	if !options.confirm {
		return memoryFindingResult(CommandMemoryBootstrap, "MEMORY-BOOTSTRAP-CONSENT-REQUIRED", resultmodel.SeverityWarning, options.memoryRoot, "bootstrap requires --confirm after summaries are reviewed", resultmodel.OutcomeRefused)
	}
	if options.manifest == "" {
		return usageResult(CommandMemoryBootstrap, errors.New("--manifest is required"))
	}
	memoryAbsolute, memoryRelative, resolveError := resolveMemoryRoot(executionContext.RepositoryRoot, options.memoryRoot, true)
	if resolveError != nil {
		return memoryFailure(CommandMemoryBootstrap, "MEMORY-ROOT-INVALID", options.memoryRoot, resolveError)
	}
	sentinelRelative := filepath.ToSlash(filepath.Join(memoryRelative, ".bootstrap-imported"))
	if data, statError := os.ReadFile(filepath.Join(memoryAbsolute, ".bootstrap-imported")); statError == nil {
		return memoryFindingResult(CommandMemoryBootstrap, "MEMORY-BOOTSTRAP-ALREADY-RUN", resultmodel.SeverityWarning, sentinelRelative, "bootstrap already ran: "+strings.TrimSpace(string(data)), resultmodel.OutcomeRefused)
	}
	manifestBytes, readError := os.ReadFile(options.manifest)
	if readError != nil {
		return memoryFailure(CommandMemoryBootstrap, "MEMORY-BOOTSTRAP-MANIFEST-READ", options.manifest, readError)
	}
	if !utf8.Valid(manifestBytes) || bytes.IndexByte(manifestBytes, 0) >= 0 {
		return memoryFindingResult(CommandMemoryBootstrap, "MEMORY-BOOTSTRAP-MANIFEST-MALFORMED", resultmodel.SeverityError, options.manifest, "manifest must be valid UTF-8 without NUL bytes", resultmodel.OutcomeRefused)
	}
	if credentialPattern.Match(manifestBytes) {
		return memoryFindingResult(CommandMemoryBootstrap, "MEMORY-CREDENTIAL-REFUSED", resultmodel.SeverityError, options.manifest, "credential-shaped manifest text was refused before persistence", resultmodel.OutcomeRefused)
	}
	var entries []bootstrapEntry
	if err := json.Unmarshal(manifestBytes, &entries); err != nil {
		return memoryFailure(CommandMemoryBootstrap, "MEMORY-BOOTSTRAP-MANIFEST-MALFORMED", options.manifest, err)
	}
	if len(entries) == 0 {
		return memoryFindingResult(CommandMemoryBootstrap, "MEMORY-BOOTSTRAP-MANIFEST-EMPTY", resultmodel.SeverityWarning, options.manifest, "manifest has no approved summaries", resultmodel.OutcomeRefused)
	}
	writes := map[string][]byte{}
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timePattern := regexp.MustCompile(`^\d{2}:\d{2}$`)
	for _, entry := range entries {
		_, dateError := time.Parse("2006-01-02", entry.Date)
		_, timeError := time.Parse("15:04", entry.Time)
		if !datePattern.MatchString(entry.Date) || !timePattern.MatchString(entry.Time) || dateError != nil || timeError != nil || strings.TrimSpace(entry.Source) == "" || strings.TrimSpace(entry.Summary) == "" || strings.ContainsAny(entry.Source+entry.Summary, "\r\n\x00") {
			return memoryFindingResult(CommandMemoryBootstrap, "MEMORY-BOOTSTRAP-MANIFEST-MALFORMED", resultmodel.SeverityError, options.manifest, "each entry requires date, HH:MM time, source, and summary", resultmodel.OutcomeRefused)
		}
		path := filepath.ToSlash(filepath.Join(memoryRelative, "logs", entry.Date+".md"))
		if _, exists := writes[path]; !exists {
			data, err := os.ReadFile(filepath.Join(executionContext.RepositoryRoot, filepath.FromSlash(path)))
			if err != nil && !os.IsNotExist(err) {
				return memoryFailure(CommandMemoryBootstrap, "MEMORY-LOG-READ-FAILED", path, err)
			}
			writes[path] = data
		}
		data := writes[path]
		if len(data) > 0 && data[len(data)-1] != '\n' {
			data = append(data, '\n')
		}
		data = append(data, []byte(fmt.Sprintf("\n## %s UTC bootstrap import\nSource transcript: %s\n%s\n", entry.Time, entry.Source, entry.Summary))...)
		writes[path] = data
	}
	writes[sentinelRelative] = []byte(operationTime.UTC().Format("2006-01-02") + "\n")
	targets := uniqueSorted(mapKeysBytes(writes))
	createdDirs := absentDirectories(executionContext.RepositoryRoot, []string{filepath.ToSlash(filepath.Join(memoryRelative, "logs"))})
	transaction := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{RepositoryRoot: executionContext.RepositoryRoot, TargetPaths: targets, PrivateUntrackedTargetPaths: targets, CreatedDirectoryPaths: createdDirs, DryRun: options.dryRun}, func(recorder *gittransaction.MutationRecorder) error {
		if err := createTransactionDirectories(executionContext.RepositoryRoot, recorder, createdDirs); err != nil {
			return err
		}
		for _, path := range targets {
			if err := publishTransactionFile(executionContext.RepositoryRoot, recorder, path, writes[path], 0o600, true); err != nil {
				return err
			}
		}
		return nil
	})
	result := gittransaction.BuildCommandResult(CommandMemoryBootstrap, transaction)
	if options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		result.Changes = memoryPlannedChanges(targets)
	}
	if !options.dryRun && result.Outcome == resultmodel.OutcomeSuccess {
		appendMemoryLedger(memoryAbsolute, "write", "", len(entries), "bootstrap")
	}
	return result
}

func handleMemoryAudit(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseMemoryOptions(arguments, false)
	if err != nil {
		return usageResult(CommandMemoryAudit, err)
	}
	if options.engine != "both" && options.engine != "memory" && options.engine != "bkb" {
		return usageResult(CommandMemoryAudit, errors.New("--engine requires bkb, memory, or both"))
	}
	findings := []resultmodel.CommandFinding{}
	if options.engine == "both" || options.engine == "memory" {
		memoryAbsolute, memoryRelative, resolveError := resolveMemoryRoot(executionContext.RepositoryRoot, options.memoryRoot, false)
		if resolveError != nil {
			findings = append(findings, memoryFinding(CommandMemoryAudit, "MEMORY-AUDIT-ENGINE", resultmodel.SeverityWarning, options.memoryRoot, "memory classification=Absent", resultmodel.FixabilityManual, "the user decides whether setup is warranted"))
		} else {
			classification, evidence := auditMemoryEngine(memoryAbsolute)
			findings = append(findings, memoryFinding(CommandMemoryAudit, "MEMORY-AUDIT-ENGINE", resultmodel.SeverityInfo, memoryRelative, "memory classification="+classification+" "+evidence, resultmodel.FixabilityManual, "the action owns the comparative recommendation"))
		}
	}
	if options.engine == "both" || options.engine == "bkb" {
		kbAbsolute, kbRelative, locateError := locateKnowledgeBase(executionContext.RepositoryRoot, options.kb)
		if locateError != nil {
			findings = append(findings, memoryFinding(CommandMemoryAudit, "MEMORY-AUDIT-BKB", resultmodel.SeverityWarning, options.kb, "bkb classification=Absent", resultmodel.FixabilityManual, "the user decides whether setup is warranted"))
		} else {
			classification, evidence := auditBKBEngine(executionContext.RepositoryRoot, kbAbsolute, kbRelative)
			findings = append(findings, memoryFinding(CommandMemoryAudit, "MEMORY-AUDIT-BKB", resultmodel.SeverityInfo, kbRelative, "bkb classification="+classification+" "+evidence+"; missing pre-ledger events are not evidence of disuse", resultmodel.FixabilityManual, "the action owns the comparative recommendation"))
		}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: findings}
}

type memoryRecallHit struct {
	Score                  int
	Path                   string
	Line                   int
	Date, Heading, Content string
}

func lexicalMemoryRecall(memoryAbsolute, memoryRelative, query string) ([]memoryRecallHit, error) {
	tokens := memoryTokens(query)
	if len(tokens) == 0 {
		return []memoryRecallHit{}, nil
	}
	sources := []string{}
	working := filepath.Join(memoryAbsolute, "working-memory.md")
	if info, err := os.Lstat(working); err == nil {
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("working-memory.md is not regular")
		}
		sources = append(sources, working)
	}
	logs, _ := filepath.Glob(filepath.Join(memoryAbsolute, "logs", "*.md"))
	sort.Strings(logs)
	sources = append(sources, logs...)
	hits := []memoryRecallHit{}
	for _, source := range sources {
		info, err := os.Lstat(source)
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("memory source %s is not regular", source)
		}
		data, err := os.ReadFile(source)
		if err != nil {
			return nil, err
		}
		heading := "(no heading)"
		relative, _ := filepath.Rel(filepath.Dir(memoryAbsolute), source)
		relative = filepath.ToSlash(relative)
		date := "working memory"
		weight := 4
		if source != working {
			date = strings.TrimSuffix(filepath.Base(source), ".md")
			parsed, parseErr := time.Parse("2006-01-02", date)
			if parseErr != nil {
				weight = 1
			} else {
				days := int(nowUTC().Sub(parsed).Hours() / 24)
				if days <= 7 {
					weight = 3
				} else if days <= 30 {
					weight = 2
				} else {
					weight = 1
				}
			}
		}
		for index, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
			if strings.HasPrefix(line, "## ") {
				heading = line
			}
			lower := strings.ToLower(line)
			distinct := 0
			for _, token := range tokens {
				if strings.Contains(lower, token) {
					distinct++
				}
			}
			if distinct > 0 {
				hits = append(hits, memoryRecallHit{Score: distinct * weight, Path: relative, Line: index + 1, Date: date, Heading: heading, Content: strings.ReplaceAll(line, "\t", " ")})
			}
		}
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].Score != hits[j].Score {
			return hits[i].Score > hits[j].Score
		}
		if hits[i].Path != hits[j].Path {
			return hits[i].Path < hits[j].Path
		}
		return hits[i].Line < hits[j].Line
	})
	if len(hits) > 8 {
		hits = hits[:8]
	}
	return hits, nil
}
func broadMemoryRecall(memoryAbsolute, memoryRelative string) ([]memoryRecallHit, error) {
	hits := []memoryRecallHit{}
	working := filepath.Join(memoryAbsolute, "working-memory.md")
	if data, err := os.ReadFile(working); err == nil {
		relative := filepath.ToSlash(filepath.Join(memoryRelative, "working-memory.md"))
		heading := "(no heading)"
		for index, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "## ") {
				heading = line
			}
			if strings.TrimSpace(line) != "" {
				hits = append(hits, memoryRecallHit{Score: 4, Path: relative, Line: index + 1, Date: "working memory", Heading: heading, Content: line})
			}
		}
	}
	logs, _ := filepath.Glob(filepath.Join(memoryAbsolute, "logs", "*.md"))
	sort.Sort(sort.Reverse(sort.StringSlice(logs)))
	if len(logs) > 3 {
		logs = logs[:3]
	}
	for _, path := range logs {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		relative := filepath.ToSlash(filepath.Join(memoryRelative, "logs", filepath.Base(path)))
		date := strings.TrimSuffix(filepath.Base(path), ".md")
		heading := ""
		capture := false
		for index, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "## ") {
				heading = line
				capture = strings.Contains(line, " UTC session capture ")
				continue
			}
			if capture {
				continue
			}
			if strings.TrimSpace(line) != "" {
				hits = append(hits, memoryRecallHit{Score: 3, Path: relative, Line: index + 1, Date: date, Heading: heading, Content: line})
			}
		}
	}
	return hits, nil
}

func findMemoryMatches(memoryAbsolute, memoryRelative, query string) ([]memoryMatch, error) {
	tokens := memoryTokens(query)
	if len(tokens) == 0 {
		return nil, nil
	}
	sources := []struct {
		absolute, relative string
		working            bool
	}{{filepath.Join(memoryAbsolute, "working-memory.md"), filepath.ToSlash(filepath.Join(memoryRelative, "working-memory.md")), true}}
	logs, _ := filepath.Glob(filepath.Join(memoryAbsolute, "logs", "*.md"))
	sort.Strings(logs)
	for _, path := range logs {
		sources = append(sources, struct {
			absolute, relative string
			working            bool
		}{path, filepath.ToSlash(filepath.Join(memoryRelative, "logs", filepath.Base(path))), false})
	}
	matches := []memoryMatch{}
	for _, source := range sources {
		data, err := os.ReadFile(source.absolute)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		for index, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "##") || trimmed == "" {
				continue
			}
			if source.working && !strings.HasPrefix(trimmed, "-") {
				continue
			}
			lower := strings.ToLower(line)
			matched := true
			for _, token := range tokens {
				if !strings.Contains(lower, token) {
					matched = false
					break
				}
			}
			if matched {
				matches = append(matches, memoryMatch{ID: memoryLineID(source.relative, index+1, line), Path: source.relative, Line: index + 1, Content: line, Working: source.working})
			}
		}
	}
	return matches, nil
}
func memoryLineID(path string, line int, content string) string {
	digest := sha256.Sum256([]byte(fmt.Sprintf("%s\x00%d\x00%s", path, line, content)))
	return hex.EncodeToString(digest[:6])
}
func memoryTokens(value string) []string {
	seen := map[string]bool{}
	result := []string{}
	var token strings.Builder
	flush := func() {
		text := strings.ToLower(token.String())
		token.Reset()
		if len(text) >= 3 && !seen[text] {
			seen[text] = true
			result = append(result, text)
		}
	}
	for _, runeValue := range value {
		if unicode.IsLetter(runeValue) || unicode.IsDigit(runeValue) {
			if runeValue < 128 {
				token.WriteRune(unicode.ToLower(runeValue))
			}
		} else {
			flush()
		}
	}
	flush()
	return result
}
func normalizeMemoryText(value string) string { return strings.Join(memoryTokens(value), " ") }
func sanitizeMemoryQuery(value string) string { return strings.Join(memoryTokens(value), " ") }

func resolveMemoryRoot(repositoryRoot, supplied string, mutating bool) (string, string, error) {
	physicalRepositoryRoot, rootError := physicalPath(filepath.Clean(repositoryRoot))
	if rootError != nil {
		return "", "", rootError
	}
	absolute := supplied
	if !filepath.IsAbs(absolute) {
		absolute = filepath.Join(physicalRepositoryRoot, supplied)
	}
	physical, err := physicalPath(filepath.Clean(absolute))
	if err != nil {
		return "", "", err
	}
	info, err := os.Stat(physical)
	if err != nil || !info.IsDir() {
		return "", "", fmt.Errorf("memory root does not exist: %s", supplied)
	}
	relative, relErr := filepath.Rel(physicalRepositoryRoot, physical)
	if relErr != nil {
		return "", "", relErr
	}
	if mutating && (relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator))) {
		return "", "", errors.New("mutating memory root must be inside the repository")
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		relative = physical
	}
	return physical, filepath.ToSlash(relative), nil
}
func appendMemoryLedger(memoryAbsolute, event, query string, hits int, note string) {
	entry := map[string]any{"ts": nowUTC().UTC().Format(time.RFC3339), "engine": "memory", "event": event, "query": query, "hits": hits, "source": "do-work-cli", "note": note}
	data, _ := json.Marshal(entry)
	root, err := os.OpenRoot(memoryAbsolute)
	if err != nil {
		return
	}
	defer root.Close()
	if info, statError := root.Lstat("usage-ledger.jsonl"); statError == nil && (!info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0) {
		return
	} else if statError != nil && !os.IsNotExist(statError) {
		return
	}
	file, err := root.OpenFile("usage-ledger.jsonl", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.Write(append(data, '\n'))
}
func memoryLedgerTail(path string, limit int) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "none"
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) > limit {
		lines = lines[len(lines)-limit:]
	}
	return strings.Join(lines, " | ")
}
func auditMemoryEngine(root string) (string, string) {
	events, newest, hitCited := 0, time.Time{}, 0
	data, _ := os.ReadFile(filepath.Join(root, "usage-ledger.jsonl"))
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		var entry map[string]any
		if json.Unmarshal(scanner.Bytes(), &entry) != nil {
			continue
		}
		event := stringValue(entry["event"])
		stamp, _ := time.Parse(time.RFC3339, stringValue(entry["ts"]))
		if stamp.After(newest) {
			newest = stamp
		}
		if event != "inject" && event != "capture" {
			if !stamp.Before(nowUTC().AddDate(0, 0, -14)) {
				events++
			}
			if event == "hit_cited" {
				hitCited++
			}
		}
	}
	if workingBytes, err := os.ReadFile(filepath.Join(root, "working-memory.md")); err == nil {
		if fields, _, parseErr := parseFrontmatter(string(workingBytes)); parseErr == nil {
			if date, dateErr := time.Parse("2006-01-02", fields["updated"]); dateErr == nil && date.After(newest) {
				newest = date
			}
		}
	}
	logs, _ := filepath.Glob(filepath.Join(root, "logs", "*.md"))
	for _, path := range logs {
		if date, err := time.Parse("2006-01-02", strings.TrimSuffix(filepath.Base(path), ".md")); err == nil && date.After(newest) {
			newest = date
		}
	}
	classification := "Idle"
	if newest.IsZero() || newest.Before(nowUTC().AddDate(0, 0, -30)) {
		classification = "Stale"
	} else if events >= 3 && hitCited >= 1 {
		classification = "Active"
	}
	return classification, fmt.Sprintf("non_automatic_14d=%d hit_cited=%d ledger_newest=%s machine_local=true", events, hitCited, displayAuditTime(newest))
}
func auditBKBEngine(repositoryRoot, kbAbsolute, kbRelative string) (string, string) {
	ledger := filepath.Join(kbAbsolute, "usage-ledger.jsonl")
	events, newest, hitCited := 0, time.Time{}, 0
	data, _ := os.ReadFile(ledger)
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		var entry map[string]any
		if json.Unmarshal(scanner.Bytes(), &entry) != nil {
			continue
		}
		stamp, _ := time.Parse(time.RFC3339, stringValue(entry["ts"]))
		if stamp.After(newest) {
			newest = stamp
		}
		event := stringValue(entry["event"])
		if event != "inject" && event != "capture" && !stamp.Before(nowUTC().AddDate(0, 0, -14)) {
			events++
		}
		if event == "hit_cited" {
			hitCited++
		}
	}
	gitCount := 0
	gitNewest := time.Time{}
	command := exec.Command("git", "-C", repositoryRoot, "log", "--format=%cI", "--", kbRelative)
	if output, err := command.Output(); err == nil {
		lines := strings.Fields(string(output))
		gitCount = len(lines)
		if len(lines) > 0 {
			gitNewest, _ = time.Parse(time.RFC3339, lines[0])
		}
	}
	logNewest := time.Time{}
	if logBytes, err := os.ReadFile(filepath.Join(kbAbsolute, "wiki", "log.md")); err == nil {
		for _, match := range datePattern.FindAllString(string(logBytes), -1) {
			if date, parseErr := time.Parse("2006-01-02", match); parseErr == nil && date.After(logNewest) {
				logNewest = date
			}
		}
	}
	activityNewest := newest
	if gitNewest.After(activityNewest) {
		activityNewest = gitNewest
	}
	if logNewest.After(activityNewest) {
		activityNewest = logNewest
	}
	classification := "Idle"
	if activityNewest.IsZero() || activityNewest.Before(nowUTC().AddDate(0, 0, -30)) {
		classification = "Stale"
	} else if events >= 3 && hitCited >= 1 || !gitNewest.IsZero() && !gitNewest.Before(nowUTC().AddDate(0, 0, -14)) || !logNewest.IsZero() && !logNewest.Before(nowUTC().AddDate(0, 0, -14)) {
		classification = "Active"
	}
	return classification, fmt.Sprintf("non_automatic_14d=%d hit_cited=%d ledger_newest=%s git_commits=%d git_newest=%s log_newest=%s ledger_committed=true", events, hitCited, displayAuditTime(newest), gitCount, displayAuditTime(gitNewest), displayAuditTime(logNewest))
}
func displayAuditTime(value time.Time) string {
	if value.IsZero() {
		return "none"
	}
	return value.UTC().Format(time.RFC3339)
}
func memoryFinding(command, code string, severity resultmodel.FindingSeverity, path, evidence string, fixability resultmodel.FindingFixability, stop string) resultmodel.CommandFinding {
	finding := knowledgeFinding(command, code, severity, nil, evidence, fixability, stop)
	if path != "" {
		finding.AffectedPaths = []string{path}
		target := memoryTargetFromAffected(path)
		finding.NextArgv = []string{"do-work-cli", command, "--memory-root", target}
		finding.VerificationArgv = []string{"do-work-cli", "--format", "json", command, "--memory-root", target}
	}
	return finding
}

func retargetMemoryResult(result resultmodel.CommandResult, target, payload string) resultmodel.CommandResult {
	for index := range result.Findings {
		finding := &result.Findings[index]
		finding.NextArgv = appendFindingOption(finding.NextArgv, "--memory-root", target)
		finding.VerificationArgv = appendFindingOption(finding.VerificationArgv, "--memory-root", target)
		if len(finding.NextArgv) > 1 {
			finding.NextJustRecipe = finding.NextArgv[1]
			if finding.NextArgv[1] == CommandMemoryRecall {
				finding.NextJustRecipe += " " + quoteRecipeArgument(payload)
			}
			finding.NextJustRecipe += " --memory-root " + quoteRecipeArgument(target)
		}
	}
	return result
}

func memoryTargetFromAffected(path string) string {
	normalized := filepath.ToSlash(path)
	if marker := strings.Index(normalized, "/logs/"); marker >= 0 {
		return normalized[:marker]
	}
	base := filepath.Base(filepath.FromSlash(normalized))
	if base == "working-memory.md" || base == "usage-ledger.jsonl" || base == ".bootstrap-imported" {
		return filepath.ToSlash(filepath.Dir(filepath.FromSlash(normalized)))
	}
	return path
}
func memoryFailure(command, code, path string, err error) resultmodel.CommandResult {
	return memoryFindingResult(command, code, resultmodel.SeverityError, path, err.Error(), resultmodel.OutcomeFailure)
}
func memoryFindingResult(command, code string, severity resultmodel.FindingSeverity, path, evidence string, outcome resultmodel.CommandOutcome) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: outcome, Findings: []resultmodel.CommandFinding{memoryFinding(command, code, severity, path, evidence, resultmodel.FixabilityManual, "the deterministic memory operation cannot continue")}}
}
func memoryPlannedChanges(paths []string) []resultmodel.RecordedChange {
	changes := make([]resultmodel.RecordedChange, 0, len(paths))
	for _, path := range paths {
		changes = append(changes, resultmodel.RecordedChange{Path: path, Kind: "planned", Detail: "would publish deterministic memory bytes"})
	}
	return changes
}
func mapKeysBytes(values map[string][]byte) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}
func containsMemoryString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
func countLines(reader *bufio.Scanner) int {
	count := 0
	for reader.Scan() {
		count++
	}
	return count
}

var _ = strconv.Itoa
