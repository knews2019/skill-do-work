package lifecycleadvance

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/finalization"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const CommandRecover = "recover"

type recoverOptions struct {
	assumeSoleAuthority bool
	takeOverRequestID   string
}

func handleRecover(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, parseError := parseRecoveryArguments(arguments)
	if parseError != nil {
		return recoveryFailure(executionContext.RepositoryRoot, "RECOVERY-USAGE", parseError.Error())
	}
	mode := "observe"
	if options.assumeSoleAuthority {
		mode = "sole-authority"
	} else if options.takeOverRequestID != "" {
		mode = "take-over"
	}
	recovery := &resultmodel.RecoveryResult{
		AuthorityMode: mode, TakeOverRequestID: options.takeOverRequestID,
		NextArgv:         []string{"do-work-cli", "--format", "json", "next"},
		VerificationArgv: []string{"do-work-cli", "--format", "json", CommandRecover},
	}
	finalizationArguments := []string{"--discover"}
	if options.assumeSoleAuthority {
		finalizationArguments = append(finalizationArguments, "--assume-sole-releaser")
		recovery.VerificationArgv = append(recovery.VerificationArgv, "--assume-sole-authority")
	} else if options.takeOverRequestID != "" {
		recovery.VerificationArgv = append(recovery.VerificationArgv, "--take-over", options.takeOverRequestID)
	}
	finalizationResult := finalization.Handlers()[finalization.CommandRecoverFinalization](executionContext, finalizationArguments)
	aggregate := resultmodel.CommandResult{
		Command: CommandRecover, Outcome: finalizationResult.Outcome, RepositoryRoot: executionContext.RepositoryRoot,
		Findings:      append([]resultmodel.CommandFinding(nil), finalizationResult.Findings...),
		Changes:       append([]resultmodel.RecordedChange(nil), finalizationResult.Changes...),
		Finalization:  finalizationResult.Finalization,
		Finalizations: append([]resultmodel.FinalizationResult(nil), finalizationResult.Finalizations...),
		Recovery:      recovery,
	}
	if finalizationResult.Outcome != resultmodel.OutcomeSuccess {
		return aggregate
	}
	recovery.FinalizationPassed = true
	setAsideRequests := setAsideRequestIDs(finalizationResult.Finalizations)
	snapshot, discoveryError := repositorymodel.DiscoverRepository(executionContext.RepositoryRoot)
	if discoveryError != nil {
		return recoveryFailureWithState(executionContext.RepositoryRoot, recovery, "RECOVERY-DISCOVERY-FAILED", discoveryError.Error())
	}
	working := recoverableWorkingRequests(snapshot)
	if options.takeOverRequestID != "" {
		filtered := working[:0]
		for _, request := range working {
			if request.TypedRecord.RequestID == options.takeOverRequestID {
				filtered = append(filtered, request)
			}
		}
		working = filtered
		if len(working) != 1 {
			return recoveryFailureWithState(executionContext.RepositoryRoot, recovery, "RECOVERY-TAKEOVER-NOT-FOUND", "take-over requires exactly one recoverable working request")
		}
	}
	for _, request := range working {
		requestID := request.TypedRecord.RequestID
		requestPath := advanceRequestPath(request)
		evidence := recoveryCheckpointEvidence(snapshot.CheckpointClaimsByID[requestID])
		claimResult := resultmodel.RecoveryClaimResult{RequestID: requestID, RequestPath: requestPath, CheckpointEvidence: evidence}
		if setAsideRequests[requestID] {
			claimResult.Decision = "finalization set aside; claim preserved"
			recovery.Claims = append(recovery.Claims, claimResult)
			continue
		}
		authorized := options.assumeSoleAuthority || options.takeOverRequestID == requestID
		if !authorized {
			claimResult.Decision = "takeover available; claim preserved"
			recovery.Claims = append(recovery.Claims, claimResult)
			aggregate.Findings = append(aggregate.Findings, resultmodel.CommandFinding{
				Code: "RECOVERY-TAKEOVER-AVAILABLE", Severity: resultmodel.SeverityWarning,
				AffectedIDs: []string{requestID}, AffectedPaths: []string{requestPath},
				Evidence: recoveryEvidenceLabels(evidence), Fixability: resultmodel.FixabilityManual,
				AutomationStopReason: "working claim requires explicit authority",
				NextArgv:             []string{"do-work-cli", CommandRecover, "--take-over", requestID},
				VerificationArgv:     []string{"do-work-cli", "--format", "json", CommandRecover},
			})
			continue
		}
		if heldForHeavyLanes(executionContext.RepositoryRoot, request) {
			claimResult.Decision = "held for heavy lanes; claim preserved"
			claimResult.HeldForHeavyLanes = true
			recovery.Claims = append(recovery.Claims, claimResult)
			aggregate.Findings = append(aggregate.Findings, resultmodel.CommandFinding{
				Code: "RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES", Severity: resultmodel.SeverityInfo,
				AffectedIDs: []string{requestID}, AffectedPaths: []string{requestPath},
				Evidence: []string{"section: " + heavyVerificationPlanSection,
					"commit: " + strings.TrimSpace(request.TypedRecord.ImplementationCommit) + " is an ancestor of HEAD"},
				Fixability:           resultmodel.FixabilityManual,
				AutomationStopReason: "heavy lanes run from the work action's drain, which owns the lane manifest path",
				NextArgv:             []string{"do-work-cli", "--format", "json", CommandAdvance, requestID},
				VerificationArgv:     []string{"do-work-cli", "--format", "json", CommandRecover},
			})
			continue
		}
		stateOptions := requeststate.StateOptions{
			Transition: requeststate.TransitionRecover, RequestID: requestID, RequestPath: requestPath,
			AssumeSoleWriter: true, Commit: true, Now: time.Now().UTC().Truncate(time.Second),
		}
		if len(evidence) == 0 {
			stateOptions.CheckpointAbsent = true
		} else {
			stateOptions.CheckpointAllEntries = true
		}
		plan := requeststate.BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), stateOptions)
		applied := requeststate.ApplyPlan(context.Background(), plan)
		aggregate.Findings = append(aggregate.Findings, applied.Findings...)
		aggregate.Changes = append(aggregate.Changes, applied.Changes...)
		if applied.Outcome != resultmodel.OutcomeSuccess {
			claimResult.Decision = "recovery refused"
			recovery.Claims = append(recovery.Claims, claimResult)
			aggregate.Outcome = applied.Outcome
			return aggregate
		}
		claimResult.Decision = "recovered to queue"
		claimResult.Recovered = true
		recovery.Claims = append(recovery.Claims, claimResult)
		aggregate.Findings = append(aggregate.Findings, resultmodel.CommandFinding{
			Code: "RECOVERY-CLAIM-RECOVERED", Severity: resultmodel.SeverityInfo,
			AffectedIDs: []string{requestID}, AffectedPaths: append([]string(nil), plan.TargetPaths...),
			Evidence: recoveryEvidenceLabels(evidence), Fixability: resultmodel.FixabilityAutomatic,
			NextArgv:         []string{"do-work-cli", "--format", "json", "next", requestID},
			VerificationArgv: []string{"do-work-cli", "--format", "json", CommandAdvance, requestID},
		})
		refreshed, refreshError := repositorymodel.DiscoverRepository(executionContext.RepositoryRoot)
		if refreshError != nil {
			aggregate.Outcome = resultmodel.OutcomeFailure
			aggregate.Findings = append(aggregate.Findings, recoveryFinding("RECOVERY-DISCOVERY-FAILED", refreshError.Error()))
			return aggregate
		}
		snapshot = refreshed
	}
	if len(working) == 0 {
		aggregate.Findings = append(aggregate.Findings, resultmodel.CommandFinding{
			Code: "RECOVERY-NONE", Severity: resultmodel.SeverityInfo,
			Evidence: []string{"no recoverable working claims"}, Fixability: resultmodel.FixabilityAutomatic,
			NextArgv: append([]string(nil), recovery.NextArgv...), VerificationArgv: append([]string(nil), recovery.VerificationArgv...),
		})
	}
	aggregate.Outcome = resultmodel.OutcomeSuccess
	return aggregate
}

// heavyVerificationPlanSection is the request-body section the work action
// appends when it holds a reviewed request for its selected heavy lanes.
const heavyVerificationPlanSection = "Heavy Verification Plan"

// heldForHeavyLanes reports a claimed working request that is waiting for this
// session's heavy-lane drain rather than an interrupted build. The evidence is
// the held phase's own two marks: the plan section the hold appended, and an
// implementation commit that is already an ancestor of HEAD. The section alone
// is not proof — a request whose commit never landed on this history is
// ordinary interrupted work. No clock is read.
func heldForHeavyLanes(repositoryRoot string, request *repositorymodel.RequestFile) bool {
	implementationCommit := strings.TrimSpace(request.TypedRecord.ImplementationCommit)
	if implementationCommit == "" || request.ParsedDocument == nil {
		return false
	}
	sections, sectionError := advanceSections(request.ParsedDocument.BodyBytes())
	if sectionError != "" || !hasSection(sections, heavyVerificationPlanSection) {
		return false
	}
	return exec.Command("git", "-C", repositoryRoot, "merge-base", "--is-ancestor", implementationCommit, "HEAD").Run() == nil
}

// setAsideRequestIDs collects the REQs this run's finalization pass excluded
// from selection. Their claims must survive recovery: releasing one puts the
// REQ back in the queue where the same run selects it again, and its unfinished
// journal then refuses the rebuilt work at the finalize tail (REQ-515). The
// exclusion rides in the record's reason codes, so nothing new is parsed here.
func setAsideRequestIDs(records []resultmodel.FinalizationResult) map[string]bool {
	setAside := map[string]bool{}
	for _, record := range records {
		for _, reasonCode := range record.ReasonCodes {
			if reasonCode == finalization.SetAsideReasonCode && record.RequestID != "" {
				setAside[record.RequestID] = true
			}
		}
	}
	return setAside
}

func parseRecoveryArguments(arguments []string) (recoverOptions, error) {
	options := recoverOptions{}
	for argumentIndex := 0; argumentIndex < len(arguments); argumentIndex++ {
		switch arguments[argumentIndex] {
		case "--assume-sole-authority":
			if options.assumeSoleAuthority {
				return recoverOptions{}, fmt.Errorf("--assume-sole-authority may be supplied only once")
			}
			options.assumeSoleAuthority = true
		case "--take-over":
			argumentIndex++
			if argumentIndex >= len(arguments) || options.takeOverRequestID != "" || !advanceRequestIDPattern.MatchString(arguments[argumentIndex]) {
				return recoverOptions{}, fmt.Errorf("--take-over requires one REQ-NNN")
			}
			options.takeOverRequestID = arguments[argumentIndex]
		default:
			return recoverOptions{}, fmt.Errorf("unknown recover option %q", arguments[argumentIndex])
		}
	}
	if options.assumeSoleAuthority && options.takeOverRequestID != "" {
		return recoverOptions{}, fmt.Errorf("--assume-sole-authority and --take-over are mutually exclusive")
	}
	return options, nil
}

func recoverableWorkingRequests(snapshot *repositorymodel.RepositorySnapshot) []*repositorymodel.RequestFile {
	requests := []*repositorymodel.RequestFile{}
	for _, request := range snapshot.RequestFiles {
		status := request.TypedRecord.RequestStatus
		blockedBy := strings.TrimSpace(request.TypedRecord.FieldEvidenceByName["blocked_by"].ScalarValue)
		if request.TreeSection == "working" && (status == "claimed" || status == "blocked" && blockedBy != "") {
			requests = append(requests, request)
		}
	}
	sort.Slice(requests, func(leftIndex, rightIndex int) bool {
		if requests[leftIndex].TypedRecord.RequestID == requests[rightIndex].TypedRecord.RequestID {
			return requests[leftIndex].RelativePath < requests[rightIndex].RelativePath
		}
		return requests[leftIndex].TypedRecord.RequestID < requests[rightIndex].TypedRecord.RequestID
	})
	return requests
}

func recoveryCheckpointEvidence(entries []repositorymodel.CheckpointClaimEvidence) []resultmodel.SelectionClaimEvidence {
	evidence := make([]resultmodel.SelectionClaimEvidence, 0, len(entries))
	for _, entry := range entries {
		evidence = append(evidence, resultmodel.SelectionClaimEvidence{
			Source: "checkpoint", ClaimedAt: entry.ClaimedAt, Writer: entry.Writer,
			Path: "do-work/" + entry.RelativePath, SourceLine: entry.SourceLine, HeaderText: entry.HeaderText,
		})
	}
	return evidence
}

func recoveryEvidenceLabels(evidence []resultmodel.SelectionClaimEvidence) []string {
	if len(evidence) == 0 {
		return []string{"checkpoint evidence: absent"}
	}
	labels := make([]string, 0, len(evidence))
	for _, entry := range evidence {
		label := "unlabelled"
		if entry.Writer != "" {
			label = entry.Writer
		}
		labels = append(labels, fmt.Sprintf("checkpoint line %d writer: %s", entry.SourceLine, label))
	}
	return labels
}

func recoveryFailure(repositoryRoot, code, evidence string) resultmodel.CommandResult {
	return recoveryFailureWithState(repositoryRoot, nil, code, evidence)
}

func recoveryFailureWithState(repositoryRoot string, recovery *resultmodel.RecoveryResult, code, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Command: CommandRecover, Outcome: resultmodel.OutcomeFailure, RepositoryRoot: repositoryRoot, Recovery: recovery,
		Findings: []resultmodel.CommandFinding{recoveryFinding(code, evidence)}}
}

func recoveryFinding(code, evidence string) resultmodel.CommandFinding {
	return resultmodel.CommandFinding{Code: code, Severity: resultmodel.SeverityError, Evidence: []string{evidence},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "recovery could not continue safely",
		NextArgv: []string{"do-work-cli", CommandRecover}, VerificationArgv: []string{"do-work-cli", "--format", "json", CommandRecover}}
}
