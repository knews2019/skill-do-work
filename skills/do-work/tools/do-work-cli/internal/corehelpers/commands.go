package corehelpers

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/archivefetch"
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
		CommandHandoffSurvey: handleHandoffSurvey,
	}
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
	switch {
	case strings.HasPrefix(code, "SCOPE-"):
		return withPaths("git", "status", "--short"), withPaths("git", "diff", "--name-status")
	case strings.HasPrefix(code, "QUALIFY-"):
		return withPaths("git", "diff"), []string{"git", "diff", "--check"}
	case strings.HasPrefix(code, "ASSOCIATION-"):
		return withPaths("git", "status", "--short"), []string{"do-work-cli", "--format", "json", CommandInventory}
	case strings.HasPrefix(code, "RESERVATION-"):
		return []string{"do-work-cli", CommandCleanupReservations, "--dry-run"}, []string{"do-work-cli", "--format", "json", CommandCleanupReservations, "--dry-run"}
	case strings.HasPrefix(code, "HANDOFF-"):
		return []string{"do-work-cli", CommandHandoffSurvey}, []string{"do-work-cli", "--format", "json", CommandHandoffSurvey}
	case strings.HasPrefix(code, "PREFLIGHT-"):
		return []string{"do-work-cli", CommandPreflight, "--dry-run"}, []string{"do-work-cli", "--format", "json", CommandPreflight, "--dry-run"}
	case strings.HasPrefix(code, "INVENTORY-"):
		return []string{"do-work-cli", CommandProtectedInventory, "start", "--dry-run"}, []string{"do-work-cli", "--format", "json", CommandInventory}
	case strings.HasPrefix(code, "TIMESTAMP-"):
		return []string{"do-work-cli", CommandAuditTimestamps, "--fix", "--dry-run"}, withPaths("git", "diff")
	case strings.HasPrefix(code, "DELETION-"):
		return withPaths("git", "status", "--short"), withPaths("git", "diff", "--cached", "--name-status")
	case strings.HasPrefix(code, "SCREENSHOT-"):
		return withPaths("git", "status", "--short"), withPaths("git", "diff", "--check")
	case strings.HasPrefix(code, "DOWNLOAD-"):
		if len(paths) > 0 {
			return []string{"test", "!", "-e", paths[0]}, []string{"ls", "-ld", paths[0]}
		}
		return []string{"do-work-cli", CommandAtomicDownload, "--dry-run"}, []string{"do-work-cli", "--format", "json", CommandAtomicDownload, "--dry-run"}
	default:
		return withPaths("git", "status", "--short"), withPaths("git", "diff", "--check")
	}
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

func handleBlockedCheck(_ commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	probeFile := ""
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
		default:
			return usageResult(CommandBlockedCheck, "unknown option "+arguments[index])
		}
	}
	if probeFile == "" {
		return usageResult(CommandBlockedCheck, "--probe-file is required")
	}
	probeBytes, err := os.ReadFile(probeFile)
	if err != nil {
		return usageResult(CommandBlockedCheck, err.Error())
	}
	status, runError := nextselection.RunBlockedProbe(probeBytes, timeoutSeconds)
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
	evidence := fmt.Sprintf("raw probe status %d", status)
	if runError != nil {
		evidence += ": " + runError.Error()
	}
	return resultmodel.CommandResult{Outcome: outcome, Findings: []resultmodel.CommandFinding{helperFinding(code, severity, []string{probeFile}, evidence,
		resultmodel.FixabilityManual, map[bool]string{true: "the probe did not clear the blocked condition", false: ""}[status != 0],
		[]string{"do-work-cli", CommandBlockedCheck, "--probe-file", probeFile}, []string{"do-work-cli", "--format", "json", CommandBlockedCheck, "--probe-file", probeFile})}}
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
	if len(result.Findings) > 0 && result.Outcome == resultmodel.OutcomeSuccess {
		result.Outcome = resultmodel.OutcomeFindings
	}
	return result
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
	transfer := archivefetch.DownloadAtomic(context.Background(), sourceURL, targetPath)
	if transfer.Err != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{helperFinding("DOWNLOAD-FAILED", resultmodel.SeverityError, []string{targetPath}, fmt.Sprintf("raw transfer status %d: %s", transfer.StatusCode, transfer.Err), resultmodel.FixabilityManual, "no target was published", []string{"do-work-cli", CommandAtomicDownload, "--source-url", sourceURL, "--target-path", targetPath}, []string{"test", "!", "-e", targetPath})}}
	}
	return successResult([]resultmodel.RecordedChange{{Path: targetPath, Kind: "created", Detail: fmt.Sprintf("published %d bytes with mode 0600", transfer.BytesWritten)}}, []resultmodel.CommandFinding{helperFinding("DOWNLOAD-PUBLISHED", resultmodel.SeverityInfo, []string{targetPath}, fmt.Sprintf("HTTP status %d", transfer.StatusCode), resultmodel.FixabilityAutomatic, "", nil, []string{"test", "-s", targetPath})})
}

func absoluteFromRoot(repositoryRoot, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(repositoryRoot, filepath.FromSlash(path))
}
