package lifecycleadvance

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/nextselection"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	queueFrozenMemberOption  = "--frozen-member"
	queueDispatchBoundOption = "--dispatch-bound"
)

type queueAdvanceOptions struct {
	selection     nextselection.SelectionOptions
	dispatchBound int
	frozenMembers []resultmodel.QueueAdvanceMember
	continuation  bool
}

func handleQueueAdvance(executionContextRoot string, arguments []string) resultmodel.CommandResult {
	options, parseError := parseQueueAdvanceOptions(arguments)
	if parseError != nil {
		return advanceFailure("ADVANCE-USAGE", parseError.Error())
	}
	snapshot, discoveryError := discoverAdvanceRepository(executionContextRoot)
	if discoveryError != nil {
		return advanceFailure("ADVANCE-DISCOVERY-FAILED", discoveryError.Error())
	}
	unboundedObservation := options.continuation || len(options.selection.TargetTokens) > 0 || options.selection.WaveDepth != nil || options.selection.FanOutLimit != nil || options.selection.SimpleOnly
	selectionResult := selectQueueAdvance(snapshot, options.selection, unboundedObservation)
	selectionResult.RepositoryRoot = executionContextRoot
	if selectionResult.Outcome == resultmodel.OutcomeFailure {
		return selectionResult
	}
	if !options.continuation {
		options.frozenMembers = freezeQueueMembers(snapshot, selectionResult, options.selection)
	}
	queueResult := &resultmodel.QueueAdvanceResult{
		TargetTokens: append([]string(nil), options.selection.TargetTokens...),
		WaveDepth:    options.selection.WaveDepth, SkipImpactNegligible: options.selection.SkipImpactNegligible,
		SimpleOnly: options.selection.SimpleOnly, DispatchBound: options.dispatchBound,
		FrozenMembers:    append([]resultmodel.QueueAdvanceMember(nil), options.frozenMembers...),
		VerificationArgv: []string{"git", "status", "--short", "--", "do-work"},
	}
	result := resultmodel.CommandResult{
		Command: CommandAdvance, Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: executionContextRoot,
		Selected:         append([]resultmodel.SelectionRecord(nil), selectionResult.Selected...),
		Excluded:         append([]resultmodel.SelectionExclusion(nil), selectionResult.Excluded...),
		SelectionSummary: selectionResult.SelectionSummary, Findings: append([]resultmodel.CommandFinding(nil), selectionResult.Findings...),
		QueueAdvance: queueResult,
	}
	if len(options.frozenMembers) == 0 {
		queueResult.ContinuationArgv = continuationArgv(options)
		return result
	}

	selectedByID := map[string]resultmodel.SelectionRecord{}
	excludedByID := map[string]resultmodel.SelectionExclusion{}
	for _, record := range selectionResult.Selected {
		selectedByID[record.RequestID] = record
	}
	for _, exclusion := range selectionResult.Excluded {
		excludedByID[exclusion.RequestID] = exclusion
	}

	// Collision and cycle holds are terminal work for this queue pass. They are
	// handled before dispatch so one bad record never livelocks the next scan.
	holdRefused := false
	for memberIndex := range queueResult.FrozenMembers {
		member := &queueResult.FrozenMembers[memberIndex]
		if member.Consumed {
			continue
		}
		freshSnapshot, freshError := discoverAdvanceRepository(executionContextRoot)
		if freshError != nil {
			appendQueueFailure(&result, queueResult, member.RequestID, "discovery", advanceFailure("ADVANCE-DISCOVERY-FAILED", freshError.Error()))
			break
		}
		queueTarget := queueRequestAtPath(freshSnapshot, member.RequestID, member.RequestPath)
		archivePaths := archiveCollisionPaths(freshSnapshot, member.RequestID)
		transition := requeststate.Transition("")
		phase := ""
		if queueTarget != nil && len(archivePaths) > 0 {
			transition = requeststate.TransitionHoldArchiveCollision
			phase = "archive-collision-hold"
		} else if exclusion, found := excludedByID[member.RequestID]; found && exclusion.Code == "DEPENDENCY-CYCLE" && queueTarget != nil {
			transition = requeststate.TransitionHoldDependencyCycle
			phase = "dependency-cycle-hold"
		}
		if transition == "" {
			continue
		}
		holdResult := applyQueueState(freshSnapshot, queueTarget, requeststate.StateOptions{
			Transition: transition, RequestID: member.RequestID, RequestPath: member.RequestPath,
			ResolvedTarget: queueTarget, OriginalStatus: "pending", Commit: true, Now: time.Now().UTC(),
		})
		if len(archivePaths) > 0 && holdResult.Outcome == resultmodel.OutcomeSuccess {
			holdResult.Findings = append(holdResult.Findings, resultmodel.CommandFinding{
				Code: "QUEUE-ARCHIVE-COLLISION-HELD", Severity: resultmodel.SeverityWarning,
				AffectedIDs: []string{member.RequestID}, AffectedPaths: append([]string{member.RequestPath}, archivePaths...),
				Evidence:   []string{"recursive repository discovery found an archived record with the same request id"},
				Fixability: resultmodel.FixabilityManual, AutomationStopReason: "the duplicate needs user judgment",
			})
		}
		appendQueuePhase(&result, queueResult, member.RequestID, phase, holdResult)
		if holdResult.Outcome == resultmodel.OutcomeSuccess {
			member.Consumed = true
			continue
		}
		holdRefused = true
		break
	}

	claimedCount := 0
	for memberIndex := range queueResult.FrozenMembers {
		if holdRefused {
			break
		}
		if claimedCount >= options.dispatchBound {
			break
		}
		member := &queueResult.FrozenMembers[memberIndex]
		if member.Consumed {
			continue
		}
		selection, selected := selectedByID[member.RequestID]
		if !selected {
			if exclusion, excluded := excludedByID[member.RequestID]; excluded {
				appendQueueSelectionBlocker(&result, queueResult, *member, exclusion)
			}
			break
		}
		freshSnapshot, freshError := discoverAdvanceRepository(executionContextRoot)
		if freshError != nil {
			appendQueueFailure(&result, queueResult, member.RequestID, "discovery", advanceFailure("ADVANCE-DISCOVERY-FAILED", freshError.Error()))
			break
		}
		if selection.UnblockRequired {
			queueTarget := queueRequestAtPath(freshSnapshot, member.RequestID, member.RequestPath)
			unblockResult := applyQueueState(freshSnapshot, queueTarget, requeststate.StateOptions{
				Transition: requeststate.TransitionUnblock, RequestID: member.RequestID, RequestPath: member.RequestPath,
				ResolvedTarget: queueTarget, OriginalStatus: selection.OriginalStatus, ProbeStatus: selection.ProbeStatus,
				UnblockRequired: true, UnblockSource: requeststate.UnblockProbe, Commit: true, Now: time.Now().UTC(),
			})
			appendQueuePhase(&result, queueResult, member.RequestID, "unblock", unblockResult)
			if unblockResult.Outcome != resultmodel.OutcomeSuccess {
				break
			}
			freshSnapshot, freshError = discoverAdvanceRepository(executionContextRoot)
			if freshError != nil {
				appendQueueFailure(&result, queueResult, member.RequestID, "discovery", advanceFailure("ADVANCE-DISCOVERY-FAILED", freshError.Error()))
				break
			}
		}
		queueTarget := queueRequestAtPath(freshSnapshot, member.RequestID, member.RequestPath)
		claimResult := applyQueueState(freshSnapshot, queueTarget, requeststate.StateOptions{
			Transition: requeststate.TransitionClaim, RequestID: member.RequestID, RequestPath: member.RequestPath,
			ResolvedTarget: queueTarget, OriginalStatus: "pending", Provenance: stateProvenance(member.Provenance),
			WriterLabel: queueWriterLabel(executionContextRoot), Commit: true, Now: time.Now().UTC(),
		})
		appendQueuePhase(&result, queueResult, member.RequestID, "claim", claimResult)
		if claimResult.Outcome != resultmodel.OutcomeSuccess {
			break
		}
		member.Consumed = true
		member.RequestPath = filepath.ToSlash(filepath.Join("do-work", "working", filepath.Base(member.RequestPath)))
		queueResult.Claimed = append(queueResult.Claimed, *member)
		claimedCount++
	}
	queueResult.Partial = len(queueResult.Claimed) > 0 && result.Outcome != resultmodel.OutcomeSuccess
	queueResult.ContinuationArgv = continuationArgv(queueAdvanceOptions{
		selection: options.selection, dispatchBound: options.dispatchBound,
		frozenMembers: queueResult.FrozenMembers, continuation: true,
	})
	if !hasUnconsumedQueueMember(queueResult.FrozenMembers) {
		queueResult.ContinuationArgv = []string{}
	}
	return result
}

func parseQueueAdvanceOptions(arguments []string) (queueAdvanceOptions, error) {
	options := queueAdvanceOptions{dispatchBound: 1}
	selectionArguments := []string{}
	for argumentIndex := 0; argumentIndex < len(arguments); argumentIndex++ {
		argument := arguments[argumentIndex]
		switch argument {
		case queueDispatchBoundOption:
			argumentIndex++
			if argumentIndex >= len(arguments) {
				return options, fmt.Errorf("%s requires a positive integer", queueDispatchBoundOption)
			}
			bound, parseError := strconv.Atoi(arguments[argumentIndex])
			if parseError != nil || bound < 1 {
				return options, fmt.Errorf("%s requires a positive integer", queueDispatchBoundOption)
			}
			options.dispatchBound = bound
			options.continuation = true
		case queueFrozenMemberOption:
			if argumentIndex+4 >= len(arguments) {
				return options, fmt.Errorf("%s requires id, path, provenance, and consumed state", queueFrozenMemberOption)
			}
			member := resultmodel.QueueAdvanceMember{
				RequestID: arguments[argumentIndex+1], RequestPath: arguments[argumentIndex+2], Provenance: arguments[argumentIndex+3],
			}
			consumed, parseError := strconv.ParseBool(arguments[argumentIndex+4])
			if parseError != nil || !canonicalQueueMember(member) {
				return options, fmt.Errorf("invalid frozen member evidence")
			}
			member.Consumed = consumed
			options.frozenMembers = append(options.frozenMembers, member)
			options.continuation = true
			argumentIndex += 4
		default:
			selectionArguments = append(selectionArguments, argument)
		}
	}
	selectionOptions, parseError := nextselection.ParseOptions(selectionArguments)
	if parseError != nil {
		return options, parseError
	}
	options.selection = selectionOptions
	if !options.continuation && selectionOptions.FanOutLimit != nil {
		options.dispatchBound = *selectionOptions.FanOutLimit
	}
	if options.continuation && len(options.frozenMembers) == 0 {
		return options, fmt.Errorf("continuation requires frozen members")
	}
	seen := map[string]bool{}
	for _, member := range options.frozenMembers {
		if seen[member.RequestID] {
			return options, fmt.Errorf("duplicate frozen member %s", member.RequestID)
		}
		seen[member.RequestID] = true
	}
	return options, nil
}

func selectQueueAdvance(snapshot *repositorymodel.RepositorySnapshot, options nextselection.SelectionOptions, unbounded bool) resultmodel.CommandResult {
	selector := nextselection.Select
	if unbounded {
		selector = nextselection.SelectUnbounded
	}
	return selector(snapshot, dependencygraph.BuildGraph(snapshot), options, func(probeBytes []byte, timeoutSeconds int) (int, error) {
		return nextselection.RunBlockedProbeAtRoot(snapshot.RepositoryRoot, probeBytes, timeoutSeconds)
	})
}

func freezeQueueMembers(snapshot *repositorymodel.RepositorySnapshot, selection resultmodel.CommandResult, options nextselection.SelectionOptions) []resultmodel.QueueAdvanceMember {
	includeExcluded := len(options.TargetTokens) > 0 || options.WaveDepth != nil || options.FanOutLimit != nil || options.SimpleOnly
	members := []resultmodel.QueueAdvanceMember{}
	seen := map[string]bool{}
	appendMember := func(requestID, requestPath, provenance string) {
		if requestID == "" || seen[requestID] {
			return
		}
		if requestPath == "" {
			requestPath = queuePathForID(snapshot, requestID)
		}
		if requestPath == "" {
			return
		}
		seen[requestID] = true
		members = append(members, resultmodel.QueueAdvanceMember{RequestID: requestID, RequestPath: requestPath, Provenance: provenance})
	}
	for _, record := range selection.Selected {
		appendMember(record.RequestID, record.RequestPath, record.Provenance)
	}
	if includeExcluded {
		for _, exclusion := range selection.Excluded {
			if exclusion.Code != "TARGET-NOT-FOUND" && exclusion.Code != "FAN-OUT-LIMIT" {
				appendMember(exclusion.RequestID, exclusion.RequestPath, exclusion.Provenance)
			}
		}
	} else {
		for _, exclusion := range selection.Excluded {
			if exclusion.Code == "DEPENDENCY-CYCLE" {
				appendMember(exclusion.RequestID, exclusion.RequestPath, exclusion.Provenance)
			}
		}
	}
	return members
}

func continuationArgv(options queueAdvanceOptions) []string {
	if len(options.frozenMembers) == 0 {
		return []string{}
	}
	arguments := []string{"do-work-cli", "--format", "json", CommandAdvance}
	if options.selection.SimpleOnly {
		arguments = append(arguments, "--simple")
	}
	arguments = append(arguments, options.selection.TargetTokens...)
	if options.selection.WaveDepth != nil {
		arguments = append(arguments, "--wave", strconv.Itoa(*options.selection.WaveDepth))
	}
	if options.selection.SkipImpactNegligible {
		arguments = append(arguments, "--skip-impact-negligible")
	}
	arguments = append(arguments, queueDispatchBoundOption, strconv.Itoa(options.dispatchBound))
	for _, member := range options.frozenMembers {
		arguments = append(arguments, queueFrozenMemberOption, member.RequestID, member.RequestPath, member.Provenance, strconv.FormatBool(member.Consumed))
	}
	return arguments
}

func appendQueuePhase(result *resultmodel.CommandResult, queueResult *resultmodel.QueueAdvanceResult, requestID, phase string, phaseResult resultmodel.CommandResult) {
	queueResult.Phases = append(queueResult.Phases, resultmodel.QueueAdvancePhase{
		RequestID: requestID, Phase: phase, Outcome: phaseResult.Outcome,
		Changes:  append([]resultmodel.RecordedChange(nil), phaseResult.Changes...),
		Findings: append([]resultmodel.CommandFinding(nil), phaseResult.Findings...),
	})
	result.Changes = append(result.Changes, phaseResult.Changes...)
	result.Findings = append(result.Findings, phaseResult.Findings...)
	if phaseResult.Outcome != resultmodel.OutcomeSuccess {
		if len(queueResult.Claimed) > 0 {
			result.Outcome = resultmodel.OutcomeFindings
		} else {
			result.Outcome = phaseResult.Outcome
		}
	}
}

func appendQueueFailure(result *resultmodel.CommandResult, queueResult *resultmodel.QueueAdvanceResult, requestID, phase string, failure resultmodel.CommandResult) {
	appendQueuePhase(result, queueResult, requestID, phase, failure)
}

func appendQueueSelectionBlocker(result *resultmodel.CommandResult, queueResult *resultmodel.QueueAdvanceResult, member resultmodel.QueueAdvanceMember, exclusion resultmodel.SelectionExclusion) {
	finding := resultmodel.CommandFinding{
		Code: exclusion.Code, Severity: resultmodel.SeverityWarning,
		AffectedIDs: []string{member.RequestID}, AffectedPaths: []string{member.RequestPath}, Evidence: []string{exclusion.Reason},
		Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "a frozen target member is not currently claimable",
		NextArgv: append([]string(nil), exclusion.NextArgv...), VerificationArgv: append([]string(nil), exclusion.VerificationArgv...),
	}
	appendQueuePhase(result, queueResult, member.RequestID, "selection", resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{finding}})
}

func applyQueueState(snapshot *repositorymodel.RepositorySnapshot, target *repositorymodel.RequestFile, options requeststate.StateOptions) resultmodel.CommandResult {
	if target == nil {
		return advanceRefusal(options.RequestID, []string{options.RequestPath}, "REQUEST-SNAPSHOT-STALE", "selected queue request is no longer present", nil)
	}
	return requeststate.ApplyPlan(context.Background(), requeststate.BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), options))
}

func canonicalQueueMember(member resultmodel.QueueAdvanceMember) bool {
	if !advanceRequestIDPattern.MatchString(member.RequestID) || !strings.HasPrefix(filepath.ToSlash(filepath.Clean(member.RequestPath)), "do-work/") {
		return false
	}
	switch member.Provenance {
	case nextselection.ProvenanceDefault, nextselection.ProvenanceExplicit, nextselection.ProvenanceUserRequest, nextselection.ProvenanceSimple:
		return true
	default:
		return false
	}
}

func queueRequestAtPath(snapshot *repositorymodel.RepositorySnapshot, requestID, requestPath string) *repositorymodel.RequestFile {
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.TreeSection == "queue" && requestFile.TypedRecord.RequestID == requestID && queueRequestPath(requestFile) == requestPath {
			return requestFile
		}
	}
	return nil
}

func queuePathForID(snapshot *repositorymodel.RepositorySnapshot, requestID string) string {
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.TreeSection == "queue" && requestFile.TypedRecord.RequestID == requestID {
			return queueRequestPath(requestFile)
		}
	}
	return ""
}

func queueRequestPath(requestFile *repositorymodel.RequestFile) string {
	return filepath.ToSlash(filepath.Join("do-work", filepath.FromSlash(requestFile.RelativePath)))
}

func archiveCollisionPaths(snapshot *repositorymodel.RepositorySnapshot, requestID string) []string {
	paths := []string{}
	for _, requestFile := range snapshot.RequestsByID[requestID] {
		if requestFile.TreeSection == "archive" {
			paths = append(paths, queueRequestPath(requestFile))
		}
	}
	sort.Strings(paths)
	return paths
}

func stateProvenance(provenance string) requeststate.SelectionProvenance {
	switch provenance {
	case nextselection.ProvenanceExplicit:
		return requeststate.ProvenanceExplicit
	case nextselection.ProvenanceUserRequest:
		return requeststate.ProvenanceURExpanded
	default:
		return requeststate.ProvenanceDefault
	}
}

func queueWriterLabel(repositoryRoot string) string {
	hostname, hostnameError := os.Hostname()
	if hostnameError != nil || strings.TrimSpace(hostname) == "" {
		hostname = "unknown-host"
	}
	return hostname + ":" + repositoryRoot
}

func hasUnconsumedQueueMember(members []resultmodel.QueueAdvanceMember) bool {
	for _, member := range members {
		if !member.Consumed {
			return true
		}
	}
	return false
}
