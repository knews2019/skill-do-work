package nextselection

import (
	"fmt"
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
		record, exclusion, probed, probeSucceeded := evaluateCandidate(candidate, graph, options, probeRunner)
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

	limit := len(eligible)
	if options.FanOutLimit != nil && *options.FanOutLimit < limit {
		limit = *options.FanOutLimit
	} else if options.FanOutLimit == nil && len(options.TargetTokens) == 0 && options.WaveDepth == nil && !options.SimpleOnly && limit > 1 {
		limit = 1
	}
	result.Selected = append(result.Selected, eligible[:limit]...)
	for _, record := range eligible[limit:] {
		result.Excluded = append(result.Excluded, exclusionFor(record.RequestID, record.Title, record.Provenance,
			"FAN-OUT-LIMIT", "ready but outside this invocation's fan-out bound", []string{"do-work-cli", "next", record.RequestID}))
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

func evaluateCandidate(candidate selectionCandidate, graph *dependencygraph.DependencyGraph, options SelectionOptions, probeRunner ProbeRunner) (*resultmodel.SelectionRecord, *resultmodel.SelectionExclusion, bool, bool) {
	requestFile := candidate.RequestFile
	record := requestFile.TypedRecord
	identifier := candidate.RequestID
	if identifier == "" || record.RequestTitle == "" || record.OriginalStatus == "" || requestFile.ParseFailure != "" {
		reason := "required frontmatter id, title, or status is missing"
		if requestFile.ParseFailure != "" {
			reason = requestFile.ParseFailure
		}
		exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "INVALID-REQUEST", reason, []string{"do-work", "doctor"})
		return nil, &exclusion, false, false
	}

	status := record.RequestStatus
	probed := false
	probeSucceeded := false
	if status == "blocked" {
		blockedCheck := record.FieldEvidenceByName["blocked_check"].ScalarValue
		if strings.TrimSpace(blockedCheck) == "" || probeRunner == nil {
			reason := "blocked request has no runnable blocked_check; confirm its external condition"
			if blockedBy := record.FieldEvidenceByName["blocked_by"].ScalarValue; blockedBy != "" {
				reason = "blocked by " + blockedBy
			}
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "BLOCKED", reason, []string{"do-work", "clarify"})
			return nil, &exclusion, false, false
		}
		probed = true
		exitCode, probeError := probeRunner([]byte(blockedCheck), 30)
		if probeError != nil || exitCode != 0 {
			reason := fmt.Sprintf("blocked_check failed this run with exit %d", exitCode)
			if probeError != nil {
				reason = "blocked_check could not launch: " + probeError.Error()
			}
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "BLOCKED-PROBE-FAILED", reason, []string{"do-work-cli", "next", identifier})
			return nil, &exclusion, true, false
		}
		probeSucceeded = true
		status = "pending"
	}
	if status != "pending" {
		exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "STATUS-NOT-PENDING", "status is "+status, []string{"do-work", "roadmap"})
		return nil, &exclusion, probed, probeSucceeded
	}

	explicit := candidate.Provenance == ProvenanceExplicit
	if !explicit && record.AssignedTo != "" {
		exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "ASSIGNED-ELSEWHERE", "assigned to "+record.AssignedTo, []string{"do-work", "run", identifier})
		return nil, &exclusion, probed, probeSucceeded
	}
	if options.SkipImpactNegligible && !explicit && record.ImpactValue == "impact-negligible" {
		exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "IMPACT-NEGLIGIBLE", "impact-negligible and --skip-impact-negligible is set", []string{"do-work", "run", identifier})
		return nil, &exclusion, probed, probeSucceeded
	}
	if options.SimpleOnly {
		if record.EffortEstimateValue != "effort-mechanical" {
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "NOT-MECHANICAL", "effort is not effort-mechanical", []string{"do-work", "run", identifier})
			return nil, &exclusion, probed, probeSucceeded
		}
		if record.MaintenanceValue == "true" {
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "MAINTENANCE-JUDGMENT", "maintenance: rule prose has no objective test", []string{"do-work", "run", identifier})
			return nil, &exclusion, probed, probeSucceeded
		}
		if record.DomainValue == "security" {
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "SECURITY-RISK", "security: cost of a miss is unbounded", []string{"do-work", "run", identifier})
			return nil, &exclusion, probed, probeSucceeded
		}
		if record.ImpactValue == "impact-critical" {
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "IMPACT-CRITICAL", "impact-critical work requires a full-strength session", []string{"do-work", "run", identifier})
			return nil, &exclusion, probed, probeSucceeded
		}
	}

	node := graph.NodesByID[identifier]
	depth := -1
	if node != nil {
		depth = queueDependencyDepth(graph, identifier, map[string]int{}, map[string]bool{})
	}
	if options.WaveDepth != nil && depth != *options.WaveDepth {
		exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "WAVE-MISMATCH", fmt.Sprintf("dependency depth is %d, not requested wave %d", depth, *options.WaveDepth), []string{"do-work-cli", "next", "--wave", strconv.Itoa(depth)})
		return nil, &exclusion, probed, probeSucceeded
	}
	if !explicit {
		switch {
		case node == nil:
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "DEPENDENCY-UNKNOWN", "dependency graph has no node for this request", []string{"do-work", "doctor"})
			return nil, &exclusion, probed, probeSucceeded
		case node.IsAmbiguous || len(node.AmbiguousTargets) > 0:
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "DEPENDENCY-AMBIGUOUS", "dependency identity is ambiguous: "+strings.Join(node.AmbiguousTargets, ", "), []string{"do-work", "doctor"})
			return nil, &exclusion, probed, probeSucceeded
		case node.IsCyclic:
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "DEPENDENCY-CYCLE", "dependency cycle must be broken before selection", []string{"do-work", "roadmap"})
			return nil, &exclusion, probed, probeSucceeded
		case len(node.MissingTargets) > 0:
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "DEPENDENCY-MISSING", "missing dependencies: "+strings.Join(node.MissingTargets, ", "), []string{"do-work", "roadmap"})
			return nil, &exclusion, probed, probeSucceeded
		case !node.DependenciesSatisfied:
			exclusion := exclusionFor(identifier, record.RequestTitle, candidate.Provenance, "DEPENDENCIES-UNMET", "waits on "+strings.Join(node.UnmetDependencies, ", "), []string{"do-work", "run", node.UnmetDependencies[0]})
			return nil, &exclusion, probed, probeSucceeded
		}
	}

	estimateMinutes, estimateKnown := frozenEstimate(record.FieldEvidenceByName)
	if options.SimpleOnly && !estimateKnown {
		estimateMinutes, estimateKnown = 5, true
	}
	selected := resultmodel.SelectionRecord{
		RequestID: identifier, Title: record.RequestTitle, Provenance: candidate.Provenance,
		DependencyDepth: depth, Dependencies: append([]string(nil), record.DependsOn...),
		EstimateMinutes: estimateMinutes, EstimateKnown: estimateKnown,
		NextArgv: []string{"do-work", "run", identifier}, NextJustRecipe: "do-work-run " + identifier,
		VerificationArgv: []string{"do-work-cli", "--format", "json", "next", identifier},
	}
	return &selected, nil, probed, probeSucceeded
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

func exclusionFor(identifier, title, provenance, code, reason string, nextArgv []string) resultmodel.SelectionExclusion {
	return resultmodel.SelectionExclusion{
		RequestID: identifier, Title: title, Provenance: provenance, Code: code, Reason: reason,
		NextArgv: nextArgv, NextJustRecipe: justRecipeFor(nextArgv),
		VerificationArgv: []string{"do-work-cli", "--format", "json", "next", identifier},
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
