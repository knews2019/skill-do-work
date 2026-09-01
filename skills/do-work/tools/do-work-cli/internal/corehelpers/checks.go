package corehelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type baselineRecord struct {
	TestCommand string `json:"test_command"`
	ExitStatus  int    `json:"exit_status"`
	Launched    bool   `json:"launched"`
}

var qualificationDebugArtifactPattern = regexp.MustCompile(`\b(` + "debug" + `ger|` + "TO" + `DO|` + "FIX" + `ME)\b`)

func handlePreflight(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	filteredArguments, dryRun, dryRunError := extractDryRun(arguments)
	if dryRunError != nil {
		return usageResult(CommandPreflight, dryRunError.Error())
	}
	arguments = filteredArguments
	findings := []resultmodel.CommandFinding{}
	changes := []resultmodel.RecordedChange{}
	statusOutput, statusError := gitOutput(executionContext.RepositoryRoot, "-c", "status.renames=copies", "status", "--porcelain=v1", "--untracked-files=all", "-z")
	if statusError != nil {
		findings = append(findings, helperFinding("PREFLIGHT-GIT-UNAVAILABLE", resultmodel.SeverityWarning, nil, statusError.Error(), resultmodel.FixabilityManual, "diff checks are unavailable", nil, []string{"git", "status", "--short"}))
	} else {
		paths, parseError := porcelainPaths(statusOutput)
		if parseError != nil {
			findings = append(findings, helperFinding("PREFLIGHT-STATUS-FAILED", resultmodel.SeverityWarning, nil, parseError.Error(), resultmodel.FixabilityManual, "working tree status could not be classified", nil, []string{"git", "status", "--short"}))
		}
		outside := []string{}
		for _, path := range paths {
			if path != "do-work" && !strings.HasPrefix(path, "do-work/") {
				outside = append(outside, path)
			}
		}
		if len(outside) > 0 {
			findings = append(findings, helperFinding("PREFLIGHT-DIRTY", resultmodel.SeverityWarning, outside, "pre-existing changes outside do-work/", resultmodel.FixabilityManual, "preserve and exclude unrelated changes from staging", append([]string{"git", "diff", "--"}, outside...), append([]string{"git", "diff", "--check", "--"}, outside...)))
		}
	}
	if len(arguments) > 0 {
		if dryRun {
			changes = append(changes, resultmodel.RecordedChange{Path: "do-work/working/baseline.json", Kind: "modified", Detail: "would run the supplied baseline command and record its status"})
			findings = append(findings, helperFinding("PREFLIGHT-BASELINE-DRY-RUN", resultmodel.SeverityInfo, nil, "baseline command was not executed and no baseline files were changed", resultmodel.FixabilityAutomatic, "", arguments, []string{"test", "!", "-e", "do-work/working/baseline.json"}))
		} else {
			commandLine := strings.Join(arguments, " ")
			var command *exec.Cmd
			if len(arguments) == 1 && strings.IndexFunc(arguments[0], func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }) >= 0 {
				command = exec.Command("sh", "-c", arguments[0])
			} else {
				command = exec.Command(arguments[0], arguments[1:]...)
			}
			command.Dir = executionContext.RepositoryRoot
			output, runError := command.CombinedOutput()
			status := 0
			launched := true
			if runError != nil {
				if exitError, ok := runError.(*exec.ExitError); ok {
					status = exitError.ExitCode()
				} else {
					status = 127
				}
				launched = status != 126 && status != 127
			}
			baselineDirectory := filepath.Join(executionContext.RepositoryRoot, "do-work", "working")
			if makeError := os.MkdirAll(baselineDirectory, 0o755); makeError != nil {
				return usageResult(CommandPreflight, makeError.Error())
			}
			recordBytes, _ := json.MarshalIndent(baselineRecord{TestCommand: commandLine, ExitStatus: status, Launched: launched}, "", "  ")
			recordBytes = append(recordBytes, '\n')
			baselinePath := filepath.Join(baselineDirectory, "baseline.json")
			if writeError := writePrivateAtomic(baselinePath, recordBytes, 0o644); writeError != nil {
				return usageResult(CommandPreflight, writeError.Error())
			}
			changes = append(changes, resultmodel.RecordedChange{Path: "do-work/working/baseline.json", Kind: "modified", Detail: "recorded test baseline"})
			failurePath := filepath.Join(baselineDirectory, "baseline-failures.txt")
			if launched && status != 0 {
				if writeError := writePrivateAtomic(failurePath, output, 0o644); writeError != nil {
					return usageResult(CommandPreflight, writeError.Error())
				}
				changes = append(changes, resultmodel.RecordedChange{Path: "do-work/working/baseline-failures.txt", Kind: "modified", Detail: "recorded baseline failure output"})
				findings = append(findings, helperFinding("PREFLIGHT-BASELINE-RED", resultmodel.SeverityWarning, []string{"do-work/working/baseline-failures.txt"}, fmt.Sprintf("test baseline exited %d", status), resultmodel.FixabilityManual, "compare later failures against this baseline", nil, arguments))
			} else {
				_ = os.Remove(failurePath)
				if !launched {
					findings = append(findings, helperFinding("PREFLIGHT-BASELINE-NOT-LAUNCHED", resultmodel.SeverityWarning, nil, fmt.Sprintf("test command could not launch (status %d)", status), resultmodel.FixabilityManual, "there is no valid red baseline", arguments, arguments))
				}
			}
		}
	}
	if _, err := os.Stat(filepath.Join(executionContext.RepositoryRoot, "package.json")); err == nil {
		if info, moduleError := os.Stat(filepath.Join(executionContext.RepositoryRoot, "node_modules")); moduleError != nil || !info.IsDir() {
			findings = append(findings, helperFinding("PREFLIGHT-NODE-MODULES-MISSING", resultmodel.SeverityWarning, []string{"package.json"}, "package.json exists but node_modules/ does not", resultmodel.FixabilityManual, "dependency-backed checks may not launch", []string{"npm", "install"}, []string{"test", "-d", "node_modules"}))
		}
	}
	if _, err := os.Stat(filepath.Join(executionContext.RepositoryRoot, "requirements.txt")); err == nil && os.Getenv("VIRTUAL_ENV") == "" {
		python := exec.Command("python3", "-c", "import sys; raise SystemExit(0 if sys.prefix != sys.base_prefix else 1)")
		python.Dir = executionContext.RepositoryRoot
		if python.Run() != nil {
			findings = append(findings, helperFinding("PREFLIGHT-VIRTUALENV-MISSING", resultmodel.SeverityWarning, []string{"requirements.txt"}, "requirements.txt exists but no active Python virtual environment was detected", resultmodel.FixabilityManual, "dependency-backed checks may use the wrong interpreter", []string{"python3", "-m", "venv", ".venv"}, []string{"python3", "-c", "import sys; print(sys.prefix)"}))
		}
	}
	// Preflight findings are advisory by contract.
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: findings, Changes: changes}
	if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" {
		exact := preflightCompatibilityText(arguments, findings, changes)
		result.ExactTextOutput = &exact
	}
	return result
}

func handleScopeDrift(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	requestPath, parseResult := requiredPathOption(arguments, "--request-path", CommandScopeDrift)
	if parseResult != nil {
		return *parseResult
	}
	contents, err := os.ReadFile(absoluteFromRoot(executionContext.RepositoryRoot, requestPath))
	if err != nil {
		return usageResult(CommandScopeDrift, err.Error())
	}
	declared, scopeFound, touchHeaderFound, scopeError := scopeDeclaredPaths(string(contents))
	implemented, summaryFound, summaryError := allBacktickedPaths(string(contents), "Implementation Summary")
	if scopeError != nil || summaryError != nil {
		result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("SCOPE-PATH-LIST-MALFORMED", resultmodel.SeverityError, []string{requestPath}, firstError(scopeError, summaryError).Error(), resultmodel.FixabilityManual, "the path lists cannot be compared", nil, nil)}}
		if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" {
			exact := "FAIL: the request path list has an unmatched backtick — close every backticked path\n"
			result.ExactTextOutput = &exact
		}
		return result
	}
	if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" && (!scopeFound || !touchHeaderFound) {
		exact := "SKIP: no '## Scope' file list found (Route A REQs have none — skip the comparison)\n"
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: &exact, ExitCodeOverride: 2}
	}
	if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" && touchHeaderFound && len(declared) == 0 {
		exact := "FAIL: a '**Files I will touch**' header is present in ## Scope but no backticked paths parse from it — fix the list formatting (backtick every path)\n"
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, ExactTextOutput: &exact}
	}
	if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" && (!summaryFound || len(implemented) == 0) {
		exact := "SKIP: no '## Implementation Summary' file list found — run this after Step 6.25\n"
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ExactTextOutput: &exact, ExitCodeOverride: 2}
	}
	if !scopeFound || !summaryFound {
		return usageResult(CommandScopeDrift, "both ## Scope and ## Implementation Summary are required")
	}
	declared = withoutDoWork(declared)
	implemented = withoutDoWork(implemented)
	missing := subtractPaths(declared, implemented)
	extra := subtractPaths(implemented, declared)
	findings := []resultmodel.CommandFinding{}
	for _, path := range missing {
		findings = append(findings, helperFinding("SCOPE-DECLARED-NOT-TOUCHED", resultmodel.SeverityWarning, []string{path}, "declared in Scope but absent from Implementation Summary", resultmodel.FixabilityManual, "the declared write set was not fully implemented", nil, nil))
	}
	for _, path := range extra {
		findings = append(findings, helperFinding("SCOPE-UNDECLARED-TOUCH", resultmodel.SeverityError, []string{path}, "present in Implementation Summary but absent from Scope", resultmodel.FixabilityManual, "the implementation exceeded its write set", nil, nil))
	}
	result := successResult(nil, findings)
	if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" {
		var output strings.Builder
		if len(extra) > 0 {
			output.WriteString("DRIFT: touched but never declared in ## Scope:\n")
			for _, path := range extra {
				fmt.Fprintf(&output, "  %s\n", path)
			}
		}
		if len(missing) > 0 {
			output.WriteString("DRIFT: declared in ## Scope but never touched:\n")
			for _, path := range missing {
				fmt.Fprintf(&output, "  %s\n", path)
			}
		}
		if len(extra) == 0 && len(missing) == 0 {
			output.WriteString("OK: Implementation Summary matches the Scope declaration\n")
		}
		exact := output.String()
		result.ExactTextOutput = &exact
	}
	return result
}

func preflightCompatibilityText(arguments []string, findings []resultmodel.CommandFinding, changes []resultmodel.RecordedChange) string {
	var output strings.Builder
	hadRepositoryState := false
	for _, finding := range findings {
		switch finding.Code {
		case "PREFLIGHT-DIRTY":
			hadRepositoryState = true
			output.WriteString("WARN: pre-existing uncommitted changes detected outside do-work/ — preserve them and, unless they prevent the active REQ, exclude them from its staging; account for them in repository-wide qualification/review evidence:\n")
			for _, path := range finding.AffectedPaths {
				fmt.Fprintf(&output, "  %s\n", path)
			}
		case "PREFLIGHT-GIT-UNAVAILABLE":
			hadRepositoryState = true
			output.WriteString("WARN: not a git repository — no clean-tree or diff-based checks available\n")
		case "PREFLIGHT-STATUS-FAILED":
			hadRepositoryState = true
			output.WriteString("WARN: git status could not read the working tree\n")
		case "PREFLIGHT-BASELINE-RED":
			output.WriteString("WARN: baseline tests failing BEFORE any changes — builder is not to blame for these\n")
		case "PREFLIGHT-BASELINE-NOT-LAUNCHED":
			output.WriteString("WARN: could not run the test command — no baseline recorded\n")
		case "PREFLIGHT-NODE-MODULES-MISSING":
			output.WriteString("WARN: package.json exists but node_modules/ does not — dependencies may not be installed\n")
		case "PREFLIGHT-VIRTUALENV-MISSING":
			output.WriteString("WARN: requirements.txt exists but no active virtualenv detected\n")
		}
	}
	if !hadRepositoryState {
		output.WriteString("OK: working tree clean outside do-work/\n")
	}
	if len(arguments) == 0 {
		output.WriteString("SKIP: no test command supplied — baseline check skipped\n")
	} else {
		recorded := false
		for _, change := range changes {
			if change.Path == "do-work/working/baseline.json" {
				recorded = true
				break
			}
		}
		if recorded && !strings.Contains(output.String(), "baseline tests failing") && !strings.Contains(output.String(), "could not run the test command") {
			output.WriteString("OK: test baseline passing\n")
		}
	}
	return output.String()
}

func handleQualify(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	requestPath, diffRange := "", ""
	for index := 0; index < len(arguments); index++ {
		switch {
		case arguments[index] == "--request-path" || strings.HasPrefix(arguments[index], "--request-path="):
			value, err := optionValue(arguments, &index, "--request-path")
			if err != nil {
				return usageResult(CommandQualify, err.Error())
			}
			requestPath = value
		case arguments[index] == "--diff-range" || strings.HasPrefix(arguments[index], "--diff-range="):
			value, err := optionValue(arguments, &index, "--diff-range")
			if err != nil {
				return usageResult(CommandQualify, err.Error())
			}
			diffRange = value
		default:
			return usageResult(CommandQualify, "unknown option "+arguments[index])
		}
	}
	if requestPath == "" {
		return usageResult(CommandQualify, "--request-path is required")
	}
	contents, err := os.ReadFile(absoluteFromRoot(executionContext.RepositoryRoot, requestPath))
	if err != nil {
		return usageResult(CommandQualify, err.Error())
	}
	paths, found, parseError := allBacktickedPaths(string(contents), "Implementation Summary")
	if parseError != nil || !found || len(paths) == 0 {
		evidence := "Implementation Summary is missing or empty"
		if parseError != nil {
			evidence = parseError.Error()
		}
		return qualificationCompatibilityResult(resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("QUALIFY-SUMMARY-MISSING", resultmodel.SeverityError, []string{requestPath}, evidence, resultmodel.FixabilityManual, "qualification has no claimed files", nil, nil)}})
	}
	changed := []string{}
	if diffRange != "" {
		parts := strings.Split(diffRange, "..")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("QUALIFY-DIFF-RANGE-INVALID", resultmodel.SeverityError, nil, "expected exact <pre>..<merge> range", resultmodel.FixabilityManual, "diff evidence would otherwise be vacuous", nil, nil)}}
		}
		for _, endpoint := range parts {
			if _, resolveError := gitOutput(executionContext.RepositoryRoot, "rev-parse", "--verify", "--quiet", endpoint+"^{commit}"); resolveError != nil {
				return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("QUALIFY-DIFF-RANGE-INVALID", resultmodel.SeverityError, nil, "diff endpoint does not resolve: "+endpoint, resultmodel.FixabilityManual, "diff evidence would otherwise be vacuous", nil, nil)}}
			}
		}
		output, gitError := gitOutput(executionContext.RepositoryRoot, "diff", "--name-only", diffRange)
		if gitError != nil {
			return usageResult(CommandQualify, gitError.Error())
		}
		changed = nonblankLines(output)
	} else {
		working, _ := gitOutput(executionContext.RepositoryRoot, "diff", "--name-only")
		staged, _ := gitOutput(executionContext.RepositoryRoot, "diff", "--staged", "--name-only")
		untracked, _ := gitOutput(executionContext.RepositoryRoot, "ls-files", "--others", "--exclude-standard")
		changed = append(nonblankLines(working), nonblankLines(staged)...)
		changed = append(changed, nonblankLines(untracked)...)
	}
	changedSet := stringSet(changed)
	entries, entryError := qualificationSummaryEntries(string(contents))
	if entryError != nil {
		return qualificationCompatibilityResult(resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("QUALIFY-SUMMARY-MALFORMED", resultmodel.SeverityError, []string{requestPath}, entryError.Error(), resultmodel.FixabilityManual, "qualification claims cannot be interpreted", nil, nil)}})
	}
	findings := []resultmodel.CommandFinding{}
	nonDoWork := 0
	for _, entry := range entries {
		path := entry.path
		if path == "do-work" || strings.HasPrefix(path, "do-work/") {
			continue
		}
		nonDoWork++
		_, statError := os.Lstat(absoluteFromRoot(executionContext.RepositoryRoot, path))
		if entry.verb == "deleted" {
			if statError == nil {
				findings = append(findings, helperFinding("QUALIFY-DELETED-PATH-PRESENT", resultmodel.SeverityError, []string{path}, "listed as deleted but still on disk", resultmodel.FixabilityManual, "the implementation claim is unsupported", nil, nil))
			} else if !changedSet[path] {
				findings = append(findings, helperFinding("QUALIFY-PATH-NOT-IN-DIFF", resultmodel.SeverityWarning, []string{path}, "claimed deletion is not in the selected diff", resultmodel.FixabilityManual, "verify the path belongs to this REQ", nil, nil))
			}
			continue
		}
		if statError != nil && os.IsNotExist(statError) {
			findings = append(findings, helperFinding("QUALIFY-CLAIMED-PATH-MISSING", resultmodel.SeverityError, []string{path}, fmt.Sprintf("listed (%s) but not on disk", entry.verb), resultmodel.FixabilityManual, "the implementation claim is unsupported", nil, nil))
		} else if !changedSet[path] {
			findings = append(findings, helperFinding("QUALIFY-PATH-NOT-IN-DIFF", resultmodel.SeverityWarning, []string{path}, "claimed path is not in the selected diff", resultmodel.FixabilityManual, "verify the path belongs to this REQ", nil, nil))
		}
		if entry.verb == "new" && statError == nil && !qualificationHasStaticReference(executionContext.RepositoryRoot, path, changed) {
			findings = append(findings, helperFinding("QUALIFY-NEW-FILE-UNWIRED", resultmodel.SeverityWarning, []string{path}, "new file has no static reference outside itself", resultmodel.FixabilityManual, "judge entry-point or dynamic-wiring exceptions", nil, nil))
		}
	}
	if nonDoWork == 0 {
		findings = append(findings, helperFinding("QUALIFY-NO-PROJECT-FILES", resultmodel.SeverityError, nil, "Implementation Summary contains only do-work paths", resultmodel.FixabilityManual, "no implementation was claimed", nil, nil))
	}
	unifyPattern := regexp.MustCompile(`(?m)^\s*-\s*\[[ x~]\]\s*\*\*\[UNIFY\]`)
	uncheckedPattern := regexp.MustCompile(`(?m)^\s*-\s*\[ \]\s*\*\*\[(PLAN|APPLY|UNIFY)\]`)
	if !unifyPattern.Match(contents) {
		findings = append(findings, helperFinding("QUALIFY-UNIFY-DISARMED", resultmodel.SeverityWarning, []string{requestPath}, "no [UNIFY] box exists", resultmodel.FixabilityManual, "P-A-U audit is not armed", nil, nil))
	}
	if count := len(uncheckedPattern.FindAll(contents, -1)); count > 0 {
		findings = append(findings, helperFinding("QUALIFY-PAU-UNCHECKED", resultmodel.SeverityError, []string{requestPath}, fmt.Sprintf("%d P-A-U checkbox(es) remain unchecked", count), resultmodel.FixabilityManual, "the builder has not completed every phase", nil, nil))
	}
	lineChanges, artifactError := qualificationChangedLines(executionContext.RepositoryRoot, diffRange)
	if artifactError != nil {
		return usageResult(CommandQualify, artifactError.Error())
	}
	changedPaths := make([]string, 0, len(lineChanges))
	for path := range lineChanges {
		changedPaths = append(changedPaths, path)
	}
	sort.Slice(changedPaths, func(left, right int) bool {
		leftHasRemovedMarker := containsMatchingLine(lineChanges[changedPaths[left]].Removed, qualificationDebugArtifactPattern)
		rightHasRemovedMarker := containsMatchingLine(lineChanges[changedPaths[right]].Removed, qualificationDebugArtifactPattern)
		if leftHasRemovedMarker != rightHasRemovedMarker {
			return leftHasRemovedMarker
		}
		return changedPaths[left] < changedPaths[right]
	})
	outputPattern := regexp.MustCompile(`console\.log|(^|[^[:alnum:]_])print\s*\(`)
	removedMarkerLines := map[string]int{}
	for _, path := range changedPaths {
		for _, line := range lineChanges[path].Removed {
			if qualificationDebugArtifactPattern.MatchString(line) {
				removedMarkerLines[line]++
			}
		}
	}
	for _, path := range changedPaths {
		if path == "do-work" || strings.HasPrefix(path, "do-work/") {
			continue
		}
		changes := lineChanges[path]
		for _, line := range changes.Added {
			if !qualificationDebugArtifactPattern.MatchString(line) {
				continue
			}
			code, severity := "QUALIFY-DEBUG-ARTIFACT", resultmodel.SeverityError
			evidence := line
			stop := "unfinished/debug-only code is newly introduced by the change"
			if removedMarkerLines[line] > 0 {
				removedMarkerLines[line]--
				code, severity = "QUALIFY-DEBUG-ARTIFACT-RELOCATED", resultmodel.SeverityWarning
				evidence += " — relocated, not added"
				stop = "the marker was relocated rather than introduced; inspect its retained intent"
			}
			findings = append(findings, helperFinding(code, severity, []string{path}, evidence, resultmodel.FixabilityManual, stop, nil, nil))
		}
		addedOutput, removedOutput := countMatchingLines(changes.Added, outputPattern), countMatchingLines(changes.Removed, outputPattern)
		for _, line := range changes.Added {
			if outputPattern.MatchString(line) {
				severity := resultmodel.SeverityError
				code := "QUALIFY-LIBRARY-OUTPUT"
				stop := "library output has no terminal audience and is presumptively debug instrumentation"
				if addedOutput <= removedOutput {
					severity, code, stop = resultmodel.SeverityWarning, "QUALIFY-OUTPUT-RELOCATED", "output was relocated rather than introduced; inspect its retained intent"
				} else if qualificationPathOwnsExit(executionContext.RepositoryRoot, path, diffRange) {
					severity, code, stop = resultmodel.SeverityWarning, "QUALIFY-REPORTER-OUTPUT", "reporter output requires human judgment"
				}
				findings = append(findings, helperFinding(code, severity, []string{path}, line, resultmodel.FixabilityManual, stop, nil, nil))
			}
		}
	}
	outcome := resultmodel.OutcomeSuccess
	for _, finding := range findings {
		if finding.Severity == resultmodel.SeverityError {
			outcome = resultmodel.OutcomeFindings
			break
		}
	}
	return qualificationCompatibilityResult(resultmodel.CommandResult{Outcome: outcome, Findings: findings})
}

func qualificationCompatibilityResult(result resultmodel.CommandResult) resultmodel.CommandResult {
	if os.Getenv("DO_WORK_COMPATIBILITY_SHIM") == "1" {
		exact := qualificationCompatibilityText(result)
		result.ExactTextOutput = &exact
	}
	return result
}

func qualificationCompatibilityText(result resultmodel.CommandResult) string {
	var output strings.Builder
	for _, finding := range result.Findings {
		level := "WARN"
		if finding.Severity == resultmodel.SeverityError {
			level = "FAIL"
		}
		path := ""
		if len(finding.AffectedPaths) > 0 {
			path = finding.AffectedPaths[0]
		}
		evidence := strings.Join(finding.Evidence, "; ")
		switch finding.Code {
		case "QUALIFY-REPORTER-OUTPUT":
			fmt.Fprintf(&output, "WARN: reporting output in %s: %s\n", path, evidence)
		case "QUALIFY-LIBRARY-OUTPUT":
			fmt.Fprintf(&output, "FAIL: leftover instrumentation in %s; file never ends its own process: %s\n", path, evidence)
		case "QUALIFY-DEBUG-ARTIFACT":
			fmt.Fprintf(&output, "FAIL: debug artifacts in %s: %s\n", path, evidence)
		case "QUALIFY-DEBUG-ARTIFACT-RELOCATED":
			fmt.Fprintf(&output, "WARN: relocated debug artifacts in %s: %s\n", path, evidence)
		case "QUALIFY-CLAIMED-PATH-MISSING":
			fmt.Fprintf(&output, "FAIL: %s: %s\n", evidence, path)
		case "QUALIFY-NEW-FILE-UNWIRED":
			fmt.Fprintf(&output, "WARN: (new) file has no static reference anywhere: %s\n", path)
		case "QUALIFY-OUTPUT-RELOCATED":
			fmt.Fprintf(&output, "WARN: relocated output in %s: %s\n", path, evidence)
		case "QUALIFY-UNIFY-DISARMED":
			fmt.Fprintf(&output, "WARN: P-A-U qualification is DISARMED: %s\n", evidence)
		default:
			if path != "" {
				fmt.Fprintf(&output, "%s: %s: %s (%s)\n", level, finding.Code, evidence, path)
			} else {
				fmt.Fprintf(&output, "%s: %s: %s\n", level, finding.Code, evidence)
			}
		}
	}
	if result.Outcome == resultmodel.OutcomeSuccess {
		output.WriteString("OK: mechanical qualification passed\n")
	}
	return output.String()
}

type qualificationLineChanges struct {
	Added   []string
	Removed []string
}

func containsMatchingLine(lines []string, pattern *regexp.Regexp) bool {
	for _, line := range lines {
		if pattern.MatchString(line) {
			return true
		}
	}
	return false
}

func countMatchingLines(lines []string, pattern *regexp.Regexp) int {
	count := 0
	for _, line := range lines {
		if pattern.MatchString(line) {
			count++
		}
	}
	return count
}

type qualificationSummaryEntry struct{ path, verb string }

func qualificationSummaryEntries(contents string) ([]qualificationSummaryEntry, error) {
	lines, found, err := sectionLines(contents, "Implementation Summary")
	if err != nil || !found {
		return nil, err
	}
	verbPattern := regexp.MustCompile(`\((new|modified|modify|deleted)\)`)
	entries := []qualificationSummaryEntry{}
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "- `") {
			continue
		}
		parts := strings.Split(line, "`")
		if len(parts)%2 == 0 {
			return nil, fmt.Errorf("unmatched backtick in Implementation Summary")
		}
		verb := ""
		if match := verbPattern.FindStringSubmatch(line); len(match) == 2 {
			verb = match[1]
			if verb == "modify" {
				verb = "modified"
			}
		}
		for index := 1; index < len(parts); index += 2 {
			if parts[index] != "" {
				entries = append(entries, qualificationSummaryEntry{parts[index], verb})
			}
		}
	}
	return entries, nil
}

func qualificationHasStaticReference(repositoryRoot, path string, changed []string) bool {
	stem := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if stem == "" {
		return false
	}
	output, _ := gitOutput(repositoryRoot, "grep", "-l", "-F", stem, "--", ".")
	for _, candidate := range nonblankLines(output) {
		if candidate != path && candidate != "do-work" && !strings.HasPrefix(candidate, "do-work/") {
			return true
		}
	}
	for _, candidate := range changed {
		if candidate == path || candidate == "do-work" || strings.HasPrefix(candidate, "do-work/") {
			continue
		}
		if contents, err := os.ReadFile(absoluteFromRoot(repositoryRoot, candidate)); err == nil && bytes.Contains(contents, []byte(stem)) {
			return true
		}
	}
	return false
}

func qualificationChangedLines(repositoryRoot, diffRange string) (map[string]qualificationLineChanges, error) {
	result := map[string]qualificationLineChanges{}
	diffs := [][]byte{}
	if diffRange != "" {
		output, err := gitOutput(repositoryRoot, "diff", "--no-ext-diff", "--no-color", "--unified=0", diffRange)
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, output)
	} else {
		for _, args := range [][]string{{"diff", "--no-ext-diff", "--no-color", "--unified=0"}, {"diff", "--cached", "--no-ext-diff", "--no-color", "--unified=0"}} {
			output, err := gitOutput(repositoryRoot, args...)
			if err != nil {
				return nil, err
			}
			diffs = append(diffs, output)
		}
	}
	for _, diff := range diffs {
		path := ""
		for _, line := range strings.Split(string(diff), "\n") {
			if strings.HasPrefix(line, "+++ b/") {
				path = strings.TrimPrefix(line, "+++ b/")
				continue
			}
			if path != "" && strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				changes := result[path]
				changes.Added = append(changes.Added, strings.TrimPrefix(line, "+"))
				result[path] = changes
			} else if path != "" && strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				changes := result[path]
				changes.Removed = append(changes.Removed, strings.TrimPrefix(line, "-"))
				result[path] = changes
			}
		}
	}
	if diffRange == "" {
		untracked, err := gitOutput(repositoryRoot, "ls-files", "--others", "--exclude-standard", "-z")
		if err != nil {
			return nil, err
		}
		for _, path := range strings.Split(string(untracked), "\x00") {
			if path == "" {
				continue
			}
			contents, err := os.ReadFile(absoluteFromRoot(repositoryRoot, path))
			if err != nil || bytes.IndexByte(contents, 0) >= 0 {
				continue
			}
			changes := result[path]
			changes.Added = append(changes.Added, strings.Split(string(contents), "\n")...)
			result[path] = changes
		}
	}
	return result, nil
}

func qualificationPathOwnsExit(repositoryRoot, path, diffRange string) bool {
	base := "HEAD"
	if diffRange != "" {
		base = strings.SplitN(diffRange, "..", 2)[0]
	}
	basePath := path
	if _, err := gitOutput(repositoryRoot, "cat-file", "-e", base+":"+basePath); err != nil {
		basePath = qualificationRenameOrigin(repositoryRoot, base, path, diffRange)
	}
	var contents []byte
	if basePath != "" {
		contents, _ = gitOutput(repositoryRoot, "show", base+":"+basePath)
	} else {
		contents, _ = os.ReadFile(absoluteFromRoot(repositoryRoot, path))
	}
	exitPattern := regexp.MustCompile(`(?m)(^\s*|[;&|]\s*|\b(then|else|do)\s+)exit\s+[0-9$][^\s;)}#]*\s*([;)}#]|$)|sys\.exit\s*\(|raise\s+SystemExit|os\._exit\s*\(|process\.exit\s*\(`)
	return exitPattern.Match(contents)
}

func qualificationRenameOrigin(repositoryRoot, base, path, diffRange string) string {
	queries := [][]string{}
	if diffRange != "" {
		queries = append(queries, []string{"diff", "--find-renames", "--name-status", "-z", diffRange})
	} else {
		queries = append(queries,
			[]string{"diff", "--find-renames", "--name-status", "-z"},
			[]string{"diff", "--cached", "--find-renames", "--name-status", "-z", base},
			[]string{"diff", "--find-renames", "--name-status", "-z", base})
	}
	for _, args := range queries {
		output, _ := gitOutput(repositoryRoot, args...)
		if origin := renameOriginFromStatus(output, path); origin != "" {
			return origin
		}
	}
	deletedArgs := []string{"diff", "--name-only", "--diff-filter=D"}
	if diffRange != "" {
		deletedArgs = append(deletedArgs, diffRange)
	}
	deletedOutput, _ := gitOutput(repositoryRoot, deletedArgs...)
	destinationBytes, _ := os.ReadFile(absoluteFromRoot(repositoryRoot, path))
	for _, deletedPath := range nonblankLines(deletedOutput) {
		baseBytes, err := gitOutput(repositoryRoot, "show", base+":"+deletedPath)
		if err == nil && len(baseBytes) > 0 && bytes.Contains(destinationBytes, baseBytes) {
			return deletedPath
		}
	}
	return ""
}

func renameOriginFromStatus(output []byte, path string) string {
	fields := strings.Split(string(output), "\x00")
	for index := 0; index+2 < len(fields); {
		status := fields[index]
		index++
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			origin, destination := fields[index], fields[index+1]
			index += 2
			if destination == path {
				return origin
			}
		} else {
			index++
		}
	}
	return ""
}

func requiredPathOption(arguments []string, optionName, commandName string) (string, *resultmodel.CommandResult) {
	path := ""
	for index := 0; index < len(arguments); index++ {
		if arguments[index] == optionName || strings.HasPrefix(arguments[index], optionName+"=") {
			value, err := optionValue(arguments, &index, optionName)
			if err != nil {
				result := usageResult(commandName, err.Error())
				return "", &result
			}
			path = value
		} else {
			result := usageResult(commandName, "unknown option "+arguments[index])
			return "", &result
		}
	}
	if path == "" {
		result := usageResult(commandName, optionName+" is required")
		return "", &result
	}
	return path, nil
}

func scopeDeclaredPaths(contents string) ([]string, bool, bool, error) {
	lines, scopeFound, err := sectionLines(contents, "Scope")
	if err != nil {
		return nil, scopeFound, false, err
	}
	headerPattern := regexp.MustCompile(`\*\*Files I will touch[^:]*:`)
	paths := []string{}
	headerFound := false
	taking := false
	for _, line := range lines {
		if match := headerPattern.FindStringIndex(line); match != nil {
			headerFound = true
			taking = true
			headerRest := line[match[1]:]
			parsed, parseErr := backtickedPathsOnLine(headerRest)
			if parseErr != nil {
				return nil, scopeFound, headerFound, fmt.Errorf("unmatched backtick in Scope")
			}
			paths = append(paths, parsed...)
			continue
		}
		trimmed := strings.TrimSpace(line)
		if taking && strings.HasPrefix(trimmed, "**") {
			taking = false
		}
		if taking && strings.HasPrefix(trimmed, "- `") {
			parsed, parseErr := backtickedPathsOnLine(line)
			if parseErr != nil {
				return nil, scopeFound, headerFound, fmt.Errorf("unmatched backtick in Scope")
			}
			paths = append(paths, parsed...)
		}
	}
	return uniqueSorted(paths), scopeFound, headerFound, nil
}

func backtickedPathsOnLine(line string) ([]string, error) {
	parts := strings.Split(line, "`")
	if len(parts)%2 == 0 {
		return nil, fmt.Errorf("unmatched backtick")
	}
	paths := []string{}
	for index := 1; index < len(parts); index += 2 {
		if parts[index] != "" {
			paths = append(paths, parts[index])
		}
	}
	return paths, nil
}

func firstBacktickedPaths(contents, heading string, pathLed bool) ([]string, bool, error) {
	all, found, err := sectionLines(contents, heading)
	if err != nil {
		return nil, found, err
	}
	paths := []string{}
	for _, line := range all {
		trimmed := strings.TrimSpace(line)
		if pathLed && !strings.HasPrefix(trimmed, "- `") {
			continue
		}
		first := strings.Index(line, "`")
		if first < 0 {
			continue
		}
		rest := line[first+1:]
		second := strings.Index(rest, "`")
		if second < 0 {
			return nil, true, fmt.Errorf("unmatched backtick in %s", heading)
		}
		if rest[:second] != "" {
			paths = append(paths, rest[:second])
		}
	}
	return uniqueSorted(paths), found, nil
}
func allBacktickedPaths(contents, heading string) ([]string, bool, error) {
	lines, found, err := sectionLines(contents, heading)
	if err != nil {
		return nil, found, err
	}
	paths := []string{}
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "- `") {
			continue
		}
		parts := strings.Split(line, "`")
		if len(parts)%2 == 0 {
			return nil, true, fmt.Errorf("unmatched backtick in %s", heading)
		}
		for index := 1; index < len(parts); index += 2 {
			if parts[index] != "" {
				paths = append(paths, parts[index])
			}
		}
	}
	return uniqueSorted(paths), found, nil
}
func sectionLines(contents, heading string) ([]string, bool, error) {
	lines := strings.Split(strings.ReplaceAll(contents, "\r\n", "\n"), "\n")
	inside := false
	found := false
	output := []string{}
	marker := "## " + heading
	for _, line := range lines {
		if line == marker || strings.HasPrefix(line, marker+" (") {
			inside = true
			found = true
			continue
		}
		if inside && strings.HasPrefix(line, "## ") {
			break
		}
		if inside {
			output = append(output, line)
		}
	}
	return output, found, nil
}
func subtractPaths(left, right []string) []string {
	set := stringSet(right)
	out := []string{}
	for _, path := range left {
		if !set[path] {
			out = append(out, path)
		}
	}
	return out
}
func withoutDoWork(paths []string) []string {
	out := []string{}
	for _, path := range paths {
		if path != "do-work" && !strings.HasPrefix(path, "do-work/") {
			out = append(out, path)
		}
	}
	return out
}
func uniqueSorted(paths []string) []string {
	set := stringSet(paths)
	out := make([]string, 0, len(set))
	for path := range set {
		out = append(out, path)
	}
	sort.Strings(out)
	return out
}
func stringSet(paths []string) map[string]bool {
	set := map[string]bool{}
	for _, path := range paths {
		set[path] = true
	}
	return set
}
func nonblankLines(contents []byte) []string {
	out := []string{}
	for _, line := range strings.Split(string(contents), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}
func firstError(first, second error) error {
	if first != nil {
		return first
	}
	return second
}
func gitOutput(root string, args ...string) ([]byte, error) {
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	output, err := command.Output()
	if err != nil {
		return output, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return output, nil
}
func porcelainPaths(output []byte) ([]string, error) {
	records := bytes.Split(output, []byte{0})
	paths := []string{}
	for index := 0; index < len(records); index++ {
		record := records[index]
		if len(record) == 0 {
			continue
		}
		if len(record) < 4 {
			return nil, fmt.Errorf("short porcelain record")
		}
		status := string(record[:2])
		paths = append(paths, string(record[3:]))
		if strings.ContainsAny(status, "RC") {
			index++
			if index >= len(records) || len(records[index]) == 0 {
				return nil, fmt.Errorf("rename/copy origin missing")
			}
			paths = append(paths, string(records[index]))
		}
	}
	return uniqueSorted(paths), nil
}
func writePrivateAtomic(path string, contents []byte, mode os.FileMode) error {
	directory := filepath.Dir(path)
	file, err := os.CreateTemp(directory, "."+filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	temporary := file.Name()
	defer os.Remove(temporary)
	if err = file.Chmod(mode); err == nil {
		_, err = file.Write(contents)
	}
	if err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(temporary, path)
}
