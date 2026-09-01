package corehelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func handlePreflight(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
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
			findings = append(findings, helperFinding("PREFLIGHT-DIRTY", resultmodel.SeverityWarning, outside, "pre-existing changes outside do-work/", resultmodel.FixabilityManual, "preserve and exclude unrelated changes from staging", []string{"git", "status", "--short", "--untracked-files=all"}, []string{"git", "diff", "--check"}))
		}
	}
	if len(arguments) > 0 {
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
	// Preflight findings are advisory by contract.
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: findings, Changes: changes}
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
	declared, scopeFound, scopeError := firstBacktickedPaths(string(contents), "Scope", true)
	implemented, summaryFound, summaryError := allBacktickedPaths(string(contents), "Implementation Summary")
	if scopeError != nil || summaryError != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("SCOPE-PATH-LIST-MALFORMED", resultmodel.SeverityError, []string{requestPath}, firstError(scopeError, summaryError).Error(), resultmodel.FixabilityManual, "the path lists cannot be compared", nil, nil)}}
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
	return successResult(nil, findings)
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
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{helperFinding("QUALIFY-SUMMARY-MISSING", resultmodel.SeverityError, []string{requestPath}, evidence, resultmodel.FixabilityManual, "qualification has no claimed files", nil, nil)}}
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
	findings := []resultmodel.CommandFinding{}
	nonDoWork := 0
	for _, path := range withoutDoWork(paths) {
		nonDoWork++
		if _, err := os.Lstat(absoluteFromRoot(executionContext.RepositoryRoot, path)); err != nil && os.IsNotExist(err) && !changedSet[path] {
			findings = append(findings, helperFinding("QUALIFY-CLAIMED-PATH-MISSING", resultmodel.SeverityError, []string{path}, "claimed path is absent and has no deletion evidence", resultmodel.FixabilityManual, "the implementation claim is unsupported", nil, nil))
		} else if !changedSet[path] {
			findings = append(findings, helperFinding("QUALIFY-PATH-NOT-IN-DIFF", resultmodel.SeverityWarning, []string{path}, "claimed path is not in the selected diff", resultmodel.FixabilityManual, "verify the path belongs to this REQ", nil, nil))
		}
	}
	if nonDoWork == 0 {
		findings = append(findings, helperFinding("QUALIFY-NO-PROJECT-FILES", resultmodel.SeverityError, nil, "Implementation Summary contains only do-work paths", resultmodel.FixabilityManual, "no implementation was claimed", nil, nil))
	}
	if !strings.Contains(string(contents), "[UNIFY]") {
		findings = append(findings, helperFinding("QUALIFY-UNIFY-DISARMED", resultmodel.SeverityWarning, []string{requestPath}, "no [UNIFY] box exists", resultmodel.FixabilityManual, "P-A-U audit is not armed", nil, nil))
	}
	return successResult(nil, findings)
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
