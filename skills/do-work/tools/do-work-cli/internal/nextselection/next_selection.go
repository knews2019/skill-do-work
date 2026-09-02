package nextselection

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
)

func Select(snapshot *repositorymodel.RepositorySnapshot, graph *dependencygraph.DependencyGraph, options SelectionOptions, probeRunner ProbeRunner) resultmodel.CommandResult {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess}
	result.SelectionSummary = summarizeQueue(snapshot)
	candidates, targetExclusions := resolveTargets(snapshot, graph, options)
	result.Excluded = append(result.Excluded, targetExclusions...)
	eligible := []resultmodel.SelectionRecord{}
	for _, candidate := range candidates {
		record, exclusion, probed, probeSucceeded := evaluateCandidate(snapshot, candidate, graph, options, probeRunner)
		if probed {
			result.SelectionSummary.Probed++
		}
		if probeSucceeded {
			result.SelectionSummary.ProbeSucceeded++
		}
		appendSchemaWarnings(&result, candidate.RequestFile)
		if exclusion != nil {
			if exclusion.Code == "IMPACT-NEGLIGIBLE" {
				result.SelectionSummary.SkippedImpactNegligible++
			}
			result.Excluded = append(result.Excluded, *exclusion)
			continue
		}
		eligible = append(eligible, *record)
	}
	if len(options.TargetTokens) == 0 {
		sort.SliceStable(eligible, func(leftIndex, rightIndex int) bool {
			return priorityRank(eligible[leftIndex].SelectionPriority) < priorityRank(eligible[rightIndex].SelectionPriority)
		})
	}

	limit := len(eligible)
	if options.FanOutLimit != nil && *options.FanOutLimit < limit {
		limit = *options.FanOutLimit
	} else if options.FanOutLimit == nil && len(options.TargetTokens) == 0 && options.WaveDepth == nil && !options.SimpleOnly && limit > 1 {
		limit = 1
	}
	result.Selected = append(result.Selected, eligible[:limit]...)
	for _, record := range eligible[limit:] {
		exclusion := exclusionFor(record.RequestID, record.Title, record.Provenance,
			"FAN-OUT-LIMIT", "ready but outside this invocation's fan-out bound", []string{"do-work-cli", "next", record.RequestID})
		copySelectionEvidenceToExclusion(&exclusion, record)
		result.Excluded = append(result.Excluded, exclusion)
	}
	for _, record := range result.Selected {
		if record.EstimateKnown {
			result.SelectionSummary.TotalEstimatedMinutes += record.EstimateMinutes
		} else {
			result.SelectionSummary.UnknownEstimateCount++
		}
	}
	verificationArgv := nextVerificationArgv(options)
	for index := range result.Selected {
		result.Selected[index].VerificationArgv = append([]string(nil), verificationArgv...)
	}
	for index := range result.Excluded {
		result.Excluded[index].VerificationArgv = append([]string(nil), verificationArgv...)
	}
	return result
}

func nextVerificationArgv(options SelectionOptions) []string {
	arguments := []string{"do-work-cli", "--format", "json", "next"}
	if options.SimpleOnly {
		arguments = append(arguments, "--simple")
	}
	arguments = append(arguments, options.TargetTokens...)
	if options.WaveDepth != nil {
		arguments = append(arguments, "--wave", strconv.Itoa(*options.WaveDepth))
	}
	if options.SkipImpactNegligible {
		arguments = append(arguments, "--skip-impact-negligible")
	}
	if options.FanOutLimit != nil {
		arguments = append(arguments, "--fan-out", strconv.Itoa(*options.FanOutLimit))
	}
	return arguments
}

func evaluateCandidate(snapshot *repositorymodel.RepositorySnapshot, candidate selectionCandidate, graph *dependencygraph.DependencyGraph, options SelectionOptions, probeRunner ProbeRunner) (*resultmodel.SelectionRecord, *resultmodel.SelectionExclusion, bool, bool) {
	requestFile := candidate.RequestFile
	record := requestFile.TypedRecord
	identifier := candidate.RequestID
	evidence := selectionEvidence{
		RequestPath: pathForSelection(requestFile), OriginalStatus: record.RequestStatus,
		ProbeStatus: resultmodel.ProbeNotApplicable, ProbeExitCode: -1,
	}
	newExclusion := func(code, reason string, nextArgv []string) resultmodel.SelectionExclusion {
		exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, code, reason, nextArgv)
		exclusion.SelectionPriority = candidate.Priority
		applySelectionEvidenceToExclusion(&exclusion, evidence)
		return exclusion
	}
	if identifier == "" || record.RequestTitle == "" || record.OriginalStatus == "" || requestFile.ParseFailure != "" {
		reason := "required frontmatter id, title, or status is missing"
		if requestFile.ParseFailure != "" {
			reason = requestFile.ParseFailure
		}
		exclusion := newExclusion("INVALID-REQUEST", reason, []string{"do-work", "doctor"})
		return nil, &exclusion, false, false
	}
	claimEvidence := selectionClaimEvidence(requestFile, identifier, snapshot.CheckpointClaimsByID[identifier])
	if len(claimEvidence) > 0 {
		exclusion := newExclusion("ALREADY-CLAIMED", alreadyClaimedReason(claimEvidence), []string{"do-work", "doctor"})
		exclusion.ClaimEvidence = claimEvidence
		return nil, &exclusion, false, false
	}

	status := record.RequestStatus
	probed := false
	probeSucceeded := false
	if status == "blocked" {
		blockedCheck := record.FieldEvidenceByName["blocked_check"].ScalarValue
		if strings.TrimSpace(blockedCheck) == "" {
			evidence.ProbeStatus = resultmodel.ProbeMissing
			reason := "blocked request has no runnable blocked_check; confirm its external condition"
			if blockedBy := record.FieldEvidenceByName["blocked_by"].ScalarValue; blockedBy != "" {
				reason = "blocked by " + blockedBy
			}
			exclusion := newExclusion("BLOCKED", reason, []string{"do-work", "clarify"})
			return nil, &exclusion, false, false
		}
		if probeRunner == nil {
			evidence.ProbeStatus = resultmodel.ProbeLaunchFailed
			evidence.ProbeAttempted = true
			evidence.ProbeExitCode = 125
			exclusion := newExclusion("BLOCKED-PROBE-FAILED", "blocked_check could not launch: probe runner is unavailable", []string{"do-work-cli", "next", identifier})
			return nil, &exclusion, true, false
		}
		probed = true
		evidence.ProbeAttempted = true
		exitCode, probeError := probeRunner([]byte(blockedCheck), 30)
		evidence.ProbeExitCode = exitCode
		if probeError != nil || exitCode != 0 {
			reason := fmt.Sprintf("blocked_check failed this run with exit %d", exitCode)
			if probeError != nil {
				evidence.ProbeStatus = resultmodel.ProbeLaunchFailed
				reason = "blocked_check could not launch: " + probeError.Error()
			} else if exitCode == 124 {
				evidence.ProbeStatus = resultmodel.ProbeTimedOut
			} else {
				evidence.ProbeStatus = resultmodel.ProbeFailed
			}
			exclusion := newExclusion("BLOCKED-PROBE-FAILED", reason, []string{"do-work-cli", "next", identifier})
			return nil, &exclusion, true, false
		}
		probeSucceeded = true
		evidence.ProbeStatus = resultmodel.ProbeSucceeded
		evidence.UnblockRequired = true
		status = "pending"
	}
	if status != "pending" {
		exclusion := newExclusion("STATUS-NOT-PENDING", "status is "+status, []string{"do-work", "roadmap"})
		return nil, &exclusion, probed, probeSucceeded
	}

	explicit := candidate.Provenance == ProvenanceExplicit
	if !explicit && record.AssignedTo != "" {
		exclusion := newExclusion("ASSIGNED-ELSEWHERE", "assigned to "+record.AssignedTo, []string{"do-work", "run", identifier})
		return nil, &exclusion, probed, probeSucceeded
	}
	if options.SkipImpactNegligible && !explicit && record.ImpactValue == "impact-negligible" {
		exclusion := newExclusion("IMPACT-NEGLIGIBLE", "impact-negligible and --skip-impact-negligible is set", []string{"do-work", "run", identifier})
		return nil, &exclusion, probed, probeSucceeded
	}
	if options.SimpleOnly {
		if record.EffortEstimateValue != "effort-mechanical" {
			exclusion := newExclusion("NOT-MECHANICAL", "effort is not effort-mechanical", []string{"do-work", "run", identifier})
			return nil, &exclusion, probed, probeSucceeded
		}
		if record.MaintenanceValue == "true" {
			exclusion := newExclusion("MAINTENANCE-JUDGMENT", "maintenance: rule prose has no objective test", []string{"do-work", "run", identifier})
			return nil, &exclusion, probed, probeSucceeded
		}
		if record.DomainValue == "security" {
			exclusion := newExclusion("SECURITY-RISK", "security: cost of a miss is unbounded", []string{"do-work", "run", identifier})
			return nil, &exclusion, probed, probeSucceeded
		}
		if record.ImpactValue == "impact-critical" {
			exclusion := newExclusion("IMPACT-CRITICAL", "impact-critical work requires a full-strength session", []string{"do-work", "run", identifier})
			return nil, &exclusion, probed, probeSucceeded
		}
	}

	node := graph.NodesByID[identifier]
	depth := -1
	if node != nil {
		depth = queueDependencyDepth(graph, identifier, map[string]int{}, map[string]bool{})
	}
	if options.WaveDepth != nil && depth != *options.WaveDepth {
		exclusion := newExclusion("WAVE-MISMATCH", fmt.Sprintf("dependency depth is %d, not requested wave %d", depth, *options.WaveDepth), []string{"do-work-cli", "next", "--wave", strconv.Itoa(depth)})
		return nil, &exclusion, probed, probeSucceeded
	}
	if !explicit {
		switch {
		case node == nil:
			exclusion := newExclusion("DEPENDENCY-UNKNOWN", "dependency graph has no node for this request", []string{"do-work", "doctor"})
			return nil, &exclusion, probed, probeSucceeded
		case node.IsAmbiguous || len(node.AmbiguousTargets) > 0:
			exclusion := newExclusion("DEPENDENCY-AMBIGUOUS", "dependency identity is ambiguous: "+strings.Join(node.AmbiguousTargets, ", "), []string{"do-work", "doctor"})
			return nil, &exclusion, probed, probeSucceeded
		case node.IsCyclic:
			exclusion := newExclusion("DEPENDENCY-CYCLE", "dependency cycle must be broken before selection", []string{"do-work", "roadmap"})
			return nil, &exclusion, probed, probeSucceeded
		case len(node.MissingTargets) > 0:
			exclusion := newExclusion("DEPENDENCY-MISSING", "missing dependencies: "+strings.Join(node.MissingTargets, ", "), []string{"do-work", "roadmap"})
			return nil, &exclusion, probed, probeSucceeded
		case !node.DependenciesSatisfied:
			exclusion := newExclusion("DEPENDENCIES-UNMET", "waits on "+strings.Join(node.UnmetDependencies, ", "), []string{"do-work", "run", node.UnmetDependencies[0]})
			return nil, &exclusion, probed, probeSucceeded
		}
	}

	estimateMinutes, estimateKnown := frozenEstimate(record.FieldEvidenceByName)
	if options.SimpleOnly && !estimateKnown {
		estimateMinutes, estimateKnown = 5, true
	}
	selected := resultmodel.SelectionRecord{
		RequestID: identifier, RequestPath: evidence.RequestPath, Title: record.RequestTitle, Provenance: candidate.Provenance,
		SelectionPriority: candidate.Priority,
		OriginalStatus:    evidence.OriginalStatus, ProbeStatus: evidence.ProbeStatus,
		ProbeAttempted: evidence.ProbeAttempted, ProbeExitCode: evidence.ProbeExitCode, UnblockRequired: evidence.UnblockRequired,
		DependencyDepth: depth, Dependencies: append([]string(nil), record.DependsOn...),
		EstimateMinutes: estimateMinutes, EstimateKnown: estimateKnown,
		NextArgv: []string{"do-work", "run", identifier}, NextJustRecipe: "do-work-run " + identifier,
		VerificationArgv: []string{"do-work-cli", "--format", "json", "next", identifier},
	}
	return &selected, nil, probed, probeSucceeded
}

func selectionClaimEvidence(requestFile *repositorymodel.RequestFile, identifier string, checkpointClaims []repositorymodel.CheckpointClaimEvidence) []resultmodel.SelectionClaimEvidence {
	evidence := []resultmodel.SelectionClaimEvidence{}
	claimedAtField := requestFile.TypedRecord.FieldEvidenceByName["claimed_at"]
	if strings.TrimSpace(requestFile.TypedRecord.ClaimedAt) != "" {
		evidence = append(evidence, resultmodel.SelectionClaimEvidence{
			Source: "request-frontmatter", ClaimedAt: claimedAtField.RawValue,
			Path: pathForSelection(requestFile), SourceLine: claimedAtField.LineNumber,
		})
	}
	for _, checkpointClaim := range checkpointClaims {
		if checkpointClaim.RequestID != identifier || strings.TrimSpace(checkpointClaim.Writer) == "" {
			continue
		}
		evidence = append(evidence, resultmodel.SelectionClaimEvidence{
			Source: "checkpoint", ClaimedAt: checkpointClaim.ClaimedAt, Writer: checkpointClaim.Writer,
			Path:       filepath.ToSlash(filepath.Join("do-work", filepath.FromSlash(checkpointClaim.RelativePath))),
			SourceLine: checkpointClaim.SourceLine, HeaderText: checkpointClaim.HeaderText,
		})
	}
	return evidence
}

func alreadyClaimedReason(evidence []resultmodel.SelectionClaimEvidence) string {
	requestClaims := 0
	checkpointClaims := 0
	for _, claim := range evidence {
		if claim.Source == "request-frontmatter" {
			requestClaims++
		} else if claim.Source == "checkpoint" {
			checkpointClaims++
		}
	}
	switch {
	case requestClaims > 0 && checkpointClaims > 0:
		return "request carries claimed_at and checkpoint retains a writer claim"
	case requestClaims > 0:
		return "request carries claimed_at"
	default:
		return "checkpoint retains a writer claim"
	}
}

func summarizeQueue(snapshot *repositorymodel.RepositorySnapshot) resultmodel.SelectionSummary {
	summary := resultmodel.SelectionSummary{}
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.TreeSection != "queue" {
			continue
		}
		switch requestFile.TypedRecord.RequestStatus {
		case "pending":
			summary.Pending++
		case "pending-answers":
			summary.PendingAnswers++
		case "blocked":
			summary.Blocked++
		case "blocked-archive-collision":
			summary.BlockedArchiveCollision++
		case "blocked-dependency-cycle":
			summary.BlockedDependencyCycle++
		default:
			if schemanormalization.IsTerminalResolved(requestFile.TypedRecord.RequestStatus) {
				summary.FinishedAwaitingArchive++
			}
		}
	}
	return summary
}

type selectionEvidence struct {
	RequestPath     string
	OriginalStatus  string
	ProbeStatus     resultmodel.SelectionProbeStatus
	ProbeAttempted  bool
	ProbeExitCode   int
	UnblockRequired bool
}

func pathForSelection(requestFile *repositorymodel.RequestFile) string {
	if requestFile == nil || requestFile.RelativePath == "" {
		return ""
	}
	return filepath.ToSlash(filepath.Join("do-work", filepath.FromSlash(requestFile.RelativePath)))
}

func applySelectionEvidenceToExclusion(exclusion *resultmodel.SelectionExclusion, evidence selectionEvidence) {
	exclusion.RequestPath = evidence.RequestPath
	exclusion.OriginalStatus = evidence.OriginalStatus
	exclusion.ProbeStatus = evidence.ProbeStatus
	exclusion.ProbeAttempted = evidence.ProbeAttempted
	exclusion.ProbeExitCode = evidence.ProbeExitCode
	exclusion.UnblockRequired = evidence.UnblockRequired
}

func copySelectionEvidenceToExclusion(exclusion *resultmodel.SelectionExclusion, selection resultmodel.SelectionRecord) {
	exclusion.SelectionPriority = selection.SelectionPriority
	applySelectionEvidenceToExclusion(exclusion, selectionEvidence{
		RequestPath: selection.RequestPath, OriginalStatus: selection.OriginalStatus,
		ProbeStatus: selection.ProbeStatus, ProbeAttempted: selection.ProbeAttempted,
		ProbeExitCode: selection.ProbeExitCode, UnblockRequired: selection.UnblockRequired,
	})
}

func exclusionFor(identifier, title, provenance, code, reason string, nextArgv []string) resultmodel.SelectionExclusion {
	return resultmodel.SelectionExclusion{
		RequestID: identifier, Title: title, Provenance: provenance, SelectionPriority: PriorityOrdinary, Code: code, Reason: reason,
		NextArgv: nextArgv, NextJustRecipe: justRecipeFor(nextArgv),
		VerificationArgv: []string{"do-work-cli", "--format", "json", "next", identifier},
	}
}

func priorityRank(priority string) int {
	switch priority {
	case PriorityRepositoryGateRepair:
		return 0
	case PriorityDeferredParent:
		return 1
	default:
		return 2
	}
}

func justRecipeFor(nextArgv []string) string {
	if len(nextArgv) >= 2 && nextArgv[0] == "do-work" {
		recipe := "do-work-" + nextArgv[1]
		if len(nextArgv) > 2 {
			recipe += " " + strings.Join(nextArgv[2:], " ")
		}
		return recipe
	}
	if len(nextArgv) >= 2 && nextArgv[0] == "do-work-cli" && nextArgv[1] == "next" {
		recipe := "do-work-next"
		if len(nextArgv) > 2 {
			recipe += " " + strings.Join(nextArgv[2:], " ")
		}
		return recipe
	}
	return ""
}

func frozenEstimate(fields map[string]requestmodel.FieldEvidence) (int, bool) {
	estimate, found := fields["estimate"]
	if !found || estimate.NestedValues == nil {
		return 0, false
	}
	minutes, err := strconv.Atoi(strings.TrimSpace(estimate.NestedValues["p50_active_minutes"]))
	if err != nil || minutes < 0 {
		return 0, false
	}
	return minutes, true
}

func appendSchemaWarnings(result *resultmodel.CommandResult, requestFile *repositorymodel.RequestFile) {
	record := requestFile.TypedRecord
	for _, evidence := range []struct {
		name       string
		recognized bool
		warning    string
	}{
		{"status", record.StatusEvidence.IsRecognized, record.StatusEvidence.WarningMessage},
		{"effort_estimate", record.EffortEstimateEvidence.IsRecognized, record.EffortEstimateEvidence.WarningMessage},
		{"impact", record.ImpactEvidence.IsRecognized, record.ImpactEvidence.WarningMessage},
		{"domain", record.DomainEvidence.IsRecognized, record.DomainEvidence.WarningMessage},
		{"maintenance", record.MaintenanceEvidence.IsRecognized, record.MaintenanceEvidence.WarningMessage},
	} {
		if !evidence.recognized {
			result.SkippedWork = append(result.SkippedWork, resultmodel.SkippedWork{Code: "SCHEMA-FALLBACK", Reason: requestID(requestFile) + " " + evidence.name + ": " + evidence.warning})
		}
	}
}

func queueDependencyDepth(graph *dependencygraph.DependencyGraph, identifier string, memo map[string]int, visiting map[string]bool) int {
	if depth, found := memo[identifier]; found {
		return depth
	}
	node := graph.NodesByID[identifier]
	if node == nil || node.IsCyclic || visiting[identifier] {
		return -1
	}
	visiting[identifier] = true
	maximum := -1
	for _, dependencyID := range node.DependencyIDs {
		dependency := graph.NodesByID[dependencyID]
		depth := 0
		if dependency != nil && schemanormalization.IsTerminalSuccess(dependency.RequestStatus) {
			depth = -1
		} else if dependency != nil && dependency.RequestStatus == "pending" {
			depth = queueDependencyDepth(graph, dependencyID, memo, visiting)
		}
		if depth > maximum {
			maximum = depth
		}
	}
	delete(visiting, identifier)
	memo[identifier] = maximum + 1
	return maximum + 1
}
