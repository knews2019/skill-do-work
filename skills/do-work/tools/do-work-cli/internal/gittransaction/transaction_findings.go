package gittransaction

import (
	"fmt"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// unmappedFailureCode is what a failure kind with no remediation template reports. It is a
// loud, complete error finding rather than a plausible-looking one with empty guidance,
// because a kind added without a template must be visible the first time it fires.
const unmappedFailureCode = "GIT-UNMAPPED-FAILURE"

const (
	changeKindCreated  = "created"
	changeKindModified = "modified"
)

// findingTemplate supplies only what cannot be derived from the failure itself. The finding
// code is deliberately absent: FindingCode computes it from the kind, so this table can go
// out of date without any kind losing its stable code.
type findingTemplate struct {
	severity             resultmodel.FindingSeverity
	fixability           resultmodel.FindingFixability
	automationStopReason string
	remediation          func(commandName string, result TransactionResult) (nextArgv []string, verificationArgv []string)
}

var failureTemplates = map[FailureKind]findingTemplate{
	FailureInvalidOptions: {
		severity:             resultmodel.SeverityError,
		fixability:           resultmodel.FixabilityManual,
		automationStopReason: "the transaction options are not valid, so nothing was attempted",
		remediation: func(commandName string, _ TransactionResult) ([]string, []string) {
			return []string{"do-work-cli", "--format", "text", commandName},
				[]string{"do-work-cli", "--format", "json", commandName}
		},
	},
	FailureNotGit: {
		severity:             resultmodel.SeverityError,
		fixability:           resultmodel.FixabilityManual,
		automationStopReason: "mutating commands require a Git repository",
		remediation: func(commandName string, _ TransactionResult) ([]string, []string) {
			return []string{"do-work-cli", "--repo-root", "<git-repository>", commandName},
				[]string{"git", "rev-parse", "--show-toplevel"}
		},
	},
	FailureDirtyTarget: {
		severity:             resultmodel.SeverityWarning,
		fixability:           resultmodel.FixabilityRefused,
		automationStopReason: "the target path already carries uncommitted work this transaction will not overwrite",
		remediation: func(_ string, result TransactionResult) ([]string, []string) {
			return gitPathArgv([]string{"git", "status", "--short"}, failurePaths(result)),
				gitPathArgv([]string{"git", "diff", "--quiet", "--exit-code"}, failurePaths(result))
		},
	},
	FailureDirtyIndex: {
		severity:             resultmodel.SeverityWarning,
		fixability:           resultmodel.FixabilityRefused,
		automationStopReason: "--commit needs an empty index so the commit can hold exactly the declared targets",
		remediation: func(string, TransactionResult) ([]string, []string) {
			return []string{"git", "diff", "--cached", "--name-only"},
				[]string{"git", "diff", "--cached", "--quiet", "--exit-code"}
		},
	},
	FailureMutation: {
		severity:             resultmodel.SeverityError,
		fixability:           resultmodel.FixabilityManual,
		automationStopReason: "the mutation failed and every declared target was restored",
		remediation: func(commandName string, result TransactionResult) ([]string, []string) {
			return []string{"do-work-cli", "--format", "json", commandName},
				gitPathArgv([]string{"git", "status", "--short"}, failurePaths(result))
		},
	},
	FailureCommit: {
		severity:             resultmodel.SeverityError,
		fixability:           resultmodel.FixabilityManual,
		automationStopReason: "the commit failed and every declared target was restored",
		remediation: func(commandName string, result TransactionResult) ([]string, []string) {
			return []string{"do-work-cli", "--format", "json", commandName},
				gitPathArgv([]string{"git", "status", "--short"}, failurePaths(result))
		},
	},
	FailureRollback: {
		severity:             resultmodel.SeverityError,
		fixability:           resultmodel.FixabilityManual,
		automationStopReason: "the rollback did not complete, so the worktree needs a person before any retry",
		remediation: func(_ string, result TransactionResult) ([]string, []string) {
			return gitPathArgv([]string{"git", "status", "--short"}, failurePaths(result)),
				[]string{"git", "status", "--porcelain=v1"}
		},
	},
	FailureCommittedRisk: {
		severity:             resultmodel.SeverityError,
		fixability:           resultmodel.FixabilityManual,
		automationStopReason: "the commit landed but its contents could not be verified; history is never rewritten",
		remediation: func(_ string, result TransactionResult) ([]string, []string) {
			return result.RevertArgv, []string{"git", "show", "--stat", result.CommitSHA}
		},
	},
}

// FindingCode derives a stable code from the failure kind rather than reading one out of a
// table, so a kind added later still gets a code that matches its name.
func FindingCode(kind FailureKind) string {
	return "GIT-" + strings.ToUpper(strings.ReplaceAll(string(kind), "_", "-"))
}

// BuildCommandResult turns a Git transaction into the one typed result the whole CLI
// renders. This is the only place a FailureKind becomes a finding, so later commands consume
// findings rather than re-deriving them from raw kinds.
//
// commandName is what the caller was invoked as. Every remediation template that names the
// CLI threads it through, so a finding's next_argv is a command line a reader can paste
// rather than a shape they have to fill in.
func BuildCommandResult(commandName string, result TransactionResult) resultmodel.CommandResult {
	commandResult := resultmodel.CommandResult{
		Command:        commandName,
		Outcome:        result.Outcome,
		RepositoryRoot: result.RepositoryRoot,
		Changes:        survivingChanges(result),
		Rollback:       result.Rollback,
	}
	if result.Failure != nil {
		commandResult.Findings = []resultmodel.CommandFinding{buildFinding(commandName, result)}
	}
	return commandResult
}

func buildFinding(commandName string, result TransactionResult) resultmodel.CommandFinding {
	failure := result.Failure
	template, mapped := failureTemplates[failure.Kind]
	if !mapped {
		return resultmodel.CommandFinding{
			Code:          unmappedFailureCode,
			Severity:      resultmodel.SeverityError,
			AffectedPaths: failurePaths(result),
			Evidence: []string{fmt.Sprintf("failure kind %q has no remediation template: %s",
				failure.Kind, failure.Reason)},
			Fixability:           resultmodel.FixabilityManual,
			AutomationStopReason: "this failure kind ships without remediation guidance, so no next step can be named",
			NextArgv:             []string{"git", "status", "--short"},
			VerificationArgv:     []string{"git", "status", "--porcelain=v1"},
		}
	}
	nextArgv, verificationArgv := template.remediation(commandName, result)
	return resultmodel.CommandFinding{
		Code:                 FindingCode(failure.Kind),
		Severity:             template.severity,
		AffectedPaths:        failurePaths(result),
		Evidence:             []string{failure.Reason},
		Fixability:           template.fixability,
		AutomationStopReason: template.automationStopReason,
		NextArgv:             nextArgv,
		VerificationArgv:     verificationArgv,
	}
}

// survivingChanges reports what is still on disk. A completed rollback undid every recorded
// change, so listing them there would describe a worktree that no longer exists.
func survivingChanges(result TransactionResult) []resultmodel.RecordedChange {
	if result.Rollback.Status == resultmodel.RollbackSucceeded {
		return nil
	}
	createdPaths := make(map[string]struct{}, len(result.CreatedPaths))
	for _, path := range result.CreatedPaths {
		createdPaths[path] = struct{}{}
	}
	detail := "left in the worktree"
	if result.CommitSHA != "" {
		detail = "committed in " + result.CommitSHA
	}
	changes := make([]resultmodel.RecordedChange, 0, len(result.ChangedPaths))
	for _, path := range result.ChangedPaths {
		kind := changeKindModified
		if _, created := createdPaths[path]; created {
			kind = changeKindCreated
		}
		changes = append(changes, resultmodel.RecordedChange{Path: path, Kind: kind, Detail: detail})
	}
	return changes
}

// failurePaths falls back to the declared targets so a finding always names something a
// reader can inspect, even for a failure recorded before any path was classified.
func failurePaths(result TransactionResult) []string {
	if result.Failure != nil && len(result.Failure.Paths) > 0 {
		return result.Failure.Paths
	}
	return result.ChangedPaths
}

func gitPathArgv(prefix, paths []string) []string {
	if len(paths) == 0 {
		return prefix
	}
	argv := make([]string, 0, len(prefix)+1+len(paths))
	argv = append(argv, prefix...)
	argv = append(argv, "--")
	return append(argv, paths...)
}
