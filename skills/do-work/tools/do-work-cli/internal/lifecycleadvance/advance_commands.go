// Package lifecycleadvance advances queue selection and executes the current
// mechanical evidence phase for already-working requests.
package lifecycleadvance

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/finalization"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const CommandAdvance = "advance"

var (
	discoverAdvanceRepository = repositorymodel.DiscoverRepository
	advanceRequestIDPattern   = regexp.MustCompile(`^REQ-[0-9]+$`)
	advanceHeadingPattern     = regexp.MustCompile(`(?m)^## ([^\r\n]+?)[ \t]*\r?$`)
	unresolvedQuestionPattern = regexp.MustCompile(`(?m)^- \[ \]`)
)

type sectionEvidence struct {
	start int
	end   int
	count int
}

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{CommandAdvance: handleAdvance, CommandRecover: handleRecover}
}

func handleAdvance(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	if len(arguments) == 1 && arguments[0] == "--checkpoint" {
		return handleAdvanceCheckpoint(executionContext)
	}
	if len(arguments) >= 1 && advanceRequestIDPattern.MatchString(arguments[0]) {
		snapshot, discoveryError := discoverAdvanceRepository(executionContext.RepositoryRoot)
		if discoveryError != nil {
			return advanceFailure("ADVANCE-DISCOVERY-FAILED", discoveryError.Error())
		}
		if candidates := snapshot.RequestsByID[arguments[0]]; len(candidates) == 1 && candidates[0].TreeSection != "queue" {
			if candidates[0].ParseFailure != "" || candidates[0].TypedRecord.RequestID != arguments[0] {
				return advanceRefusal(arguments[0], []string{advanceRequestPath(candidates[0])}, "ADVANCE-EVIDENCE-MISSING", "request identity or frontmatter is malformed", nil)
			}
			projected := classifyAdvance(candidates[0])
			if len(arguments) == 1 || projected.Advance == nil || projected.Outcome == resultmodel.OutcomeRefused {
				return projected
			}
			if projected.Advance.Phase == finalization.CommandFinalize {
				return executeAdvanceFinalization(executionContext, projected, arguments[1:])
			}
			return executeAdvanceEvidenceGates(executionContext, candidates[0], projected, arguments[1:])
		}
	}
	return handleQueueAdvance(executionContext.RepositoryRoot, arguments)
}

func parseAdvanceArguments(arguments []string) (string, error) {
	if len(arguments) != 1 || !advanceRequestIDPattern.MatchString(arguments[0]) {
		return "", fmt.Errorf("usage: advance REQ-NNN")
	}
	return arguments[0], nil
}

func classifyAdvance(target *repositorymodel.RequestFile) resultmodel.CommandResult {
	requestPath := advanceRequestPath(target)
	record := target.TypedRecord
	verificationArgv := []string{"do-work-cli", "--format", "json", CommandAdvance, record.RequestID}
	advance := &resultmodel.AdvanceLifecycleResult{
		RequestID: record.RequestID, RequestPath: requestPath, TreeSection: target.TreeSection,
		Status: record.RequestStatus, Route: record.RouteValue, VerificationArgv: verificationArgv,
	}

	if duplicateField(record.FieldEvidenceByName, "status") || duplicateField(record.FieldEvidenceByName, "route") || duplicateField(record.FieldEvidenceByName, "estimate") {
		return advanceRefusal(record.RequestID, []string{requestPath}, "ADVANCE-EVIDENCE-MISSING", "duplicate lifecycle field", advance)
	}
	if !record.StatusEvidence.IsRecognized {
		return advanceRefusal(record.RequestID, []string{requestPath}, "ADVANCE-PHASE-UNKNOWN", "unrecognized request status", advance)
	}

	switch target.TreeSection {
	case "queue":
		return classifyQueueAdvance(target, advance)
	case "working":
		return classifyWorkingAdvance(target, advance)
	case "archive":
		return classifyArchiveAdvance(target, advance)
	default:
		return advanceRefusal(record.RequestID, []string{requestPath}, "ADVANCE-PHASE-UNKNOWN", "request is outside a lifecycle tree", advance)
	}
}

func classifyQueueAdvance(target *repositorymodel.RequestFile, advance *resultmodel.AdvanceLifecycleResult) resultmodel.CommandResult {
	requestPath := advance.RequestPath
	switch advance.Status {
	case "pending":
		return advancePhase(advance, "claim", resultmodel.AdvancePhaseMechanical,
			advanceEvidence("field", requestPath, "status", "", "claimed"),
			[]string{"do-work-cli", "claim", advance.RequestID, "--request-path", requestPath, "--provenance", "explicit-req", "--commit"})
	case "blocked":
		return advancePhase(advance, "blocked-check", resultmodel.AdvancePhaseMechanical,
			advanceEvidence("field", requestPath, "status", "", "pending after successful probe and unblock"),
			[]string{"do-work-cli", "--format", "json", "next", advance.RequestID})
	case "pending-answers", "blocked-archive-collision", "blocked-dependency-cycle":
		return advancePhase(advance, "agent judgment: resolve "+advance.Status, resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("field", requestPath, "status", "", "pending"), nil)
	default:
		return advanceRefusal(advance.RequestID, []string{requestPath}, "ADVANCE-PHASE-UNKNOWN", "status cannot occur in the queue lifecycle phase", advance)
	}
}

func classifyWorkingAdvance(target *repositorymodel.RequestFile, advance *resultmodel.AdvanceLifecycleResult) resultmodel.CommandResult {
	requestPath := advance.RequestPath
	record := target.TypedRecord
	if advance.Status != "claimed" && advance.Status != "completed-with-issues" {
		return advanceRefusal(advance.RequestID, []string{requestPath}, "ADVANCE-PHASE-UNKNOWN", "working request is not claimed", advance)
	}
	sections, sectionError := advanceSections(target.ParsedDocument.BodyBytes())
	if sectionError != "" {
		return advanceRefusal(advance.RequestID, []string{requestPath}, "ADVANCE-EVIDENCE-MISSING", sectionError, advance)
	}

	if record.RouteValue == "" {
		if hasAnySection(sections, "Plan", "Exploration", "Scope", "Pre-Flight", "Implementation Summary", "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "route")
		}
		return advancePhase(advance, "agent judgment: triage and open questions", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("field", requestPath, "route", "", "A, B, or C"), nil)
	}
	if !record.RouteEvidence.IsRecognized || (record.RouteValue != "A" && record.RouteValue != "B" && record.RouteValue != "C") {
		return advanceRefusal(advance.RequestID, []string{requestPath}, "ADVANCE-PHASE-UNKNOWN", "route is not A, B, or C", advance)
	}
	if !hasSection(sections, "Triage") {
		if hasAnySection(sections, "Plan", "Exploration", "Scope", "Pre-Flight", "Implementation Summary", "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "Triage")
		}
		return advancePhase(advance, "agent judgment: triage and open questions", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("section", requestPath, "", "Triage", "route decision"), nil)
	}
	if sectionContains(target.ParsedDocument.BodyBytes(), sections, "Open Questions", unresolvedQuestionPattern) {
		if hasAnySection(sections, "Plan", "Exploration", "Scope", "Pre-Flight", "Implementation Summary", "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "resolved Open Questions")
		}
		return advancePhase(advance, "agent judgment: triage and open questions", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("section", requestPath, "", "Open Questions", "all questions resolved or deferred"), nil)
	}
	if !validEstimate(record.FieldEvidenceByName["estimate"].NestedValues["p50_active_minutes"]) {
		if hasAnySection(sections, "Plan", "Exploration", "Scope", "Pre-Flight", "Implementation Summary", "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "estimate.p50_active_minutes")
		}
		nextArgv := []string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", requestPath, "--", "--route", record.RouteValue, "--write-set", "<count>", "--subsystems", "<count>", "--acceptance", "<count>"}
		if record.EffortEstimateValue == "effort-mechanical" && record.RouteValue == "A" {
			nextArgv = []string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", requestPath, "--", "--trivial"}
		}
		return advancePhase(advance, "estimate-p50", resultmodel.AdvancePhaseMechanical,
			advanceEvidence("field", requestPath, "estimate.p50_active_minutes", "", "positive integer"), nextArgv)
	}

	if record.RouteValue == "A" && hasAnySection(sections, "Exploration", "Scope", "Pre-Flight") {
		return advanceRefusal(advance.RequestID, []string{requestPath}, "ADVANCE-PHASE-UNKNOWN", "Route A contains a Route B/C-only phase", advance)
	}
	if !hasSection(sections, "Plan") {
		if hasAnySection(sections, "Exploration", "Scope", "Pre-Flight", "Implementation Summary", "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "Plan")
		}
		phase := "agent judgment: record planning not required"
		expected := "planning not required"
		if record.RouteValue == "C" {
			phase = "agent judgment: planning"
			expected = "validated Route C plan"
		}
		return advancePhase(advance, phase, resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("section", requestPath, "", "Plan", expected), nil)
	}
	if record.RouteValue == "C" && strings.TrimSpace(record.FieldEvidenceByName["planning_at"].ScalarValue) == "" {
		if hasAnySection(sections, "Exploration", "Scope", "Pre-Flight", "Implementation Summary", "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "planning_at")
		}
		return advancePhase(advance, "agent judgment: planning", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("field", requestPath, "planning_at", "", "canonical timestamp after plan validation"), nil)
	}

	if record.RouteValue != "A" && !hasSection(sections, "Exploration") {
		if hasAnySection(sections, "Scope", "Pre-Flight", "Implementation Summary", "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "Exploration")
		}
		return advancePhase(advance, "agent judgment: exploration", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("section", requestPath, "", "Exploration", "exploration findings"), nil)
	}
	if record.RouteValue != "A" && !hasSection(sections, "Scope") {
		if hasAnySection(sections, "Pre-Flight", "Implementation Summary", "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "Scope")
		}
		return advancePhase(advance, "agent judgment: scope declaration", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("section", requestPath, "", "Scope", "declared files and acceptance criteria"), nil)
	}
	if record.RouteValue != "A" && len(record.WritePaths) == 0 {
		return advancePhase(advance, "agent judgment: scope declaration", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("field", requestPath, "write_set", "", "Scope file-list mirror"), nil)
	}

	if !hasSection(sections, "Implementation Summary") {
		if hasAnySection(sections, "Qualification", "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "Implementation Summary")
		}
		if record.RouteValue == "A" {
			return advancePhase(advance, "agent judgment: implementation and summary", resultmodel.AdvancePhaseAgentJudgment,
				advanceEvidence("section", requestPath, "", "Implementation Summary", "implemented file manifest"), nil)
		}
		return advancePhase(advance, "preflight", resultmodel.AdvancePhaseMechanical,
			advanceEvidence("section", requestPath, "", "Implementation Summary", "preflight complete and implementation summarized"),
			[]string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", requestPath, "--gate-arg", "<canonical-gate-argv-token>", "--", "<resolved-test-argv>"})
	}
	if !hasSection(sections, "Qualification") {
		if hasAnySection(sections, "Testing", "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "Qualification")
		}
		return advancePhase(advance, "qualify", resultmodel.AdvancePhaseMechanical,
			advanceEvidence("section", requestPath, "", "Qualification", "typed qualification result"),
			[]string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", requestPath, "--diff-range", "<pre>..<merge_hash>"})
	}
	if !hasSection(sections, "Testing") {
		if hasAnySection(sections, "Review", "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "Testing")
		}
		return advancePhase(advance, "test-gate", resultmodel.AdvancePhaseMechanical,
			advanceEvidence("section", requestPath, "", "Testing", "tests and gate evidence"),
			[]string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", requestPath, "--gate-arg", "<canonical-gate-argv-token>", "--", "--probe-file", "<focused-test-probe>"})
	}
	if !hasSection(sections, "Review") {
		if hasAnySection(sections, "Lessons Learned", "Orientation") {
			return missingBeforeLaterRefusal(advance, "Review")
		}
		return advancePhase(advance, "agent judgment: review", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("section", requestPath, "", "Review", "independent review verdict"), nil)
	}
	if !hasSection(sections, "Lessons Learned") {
		if hasSection(sections, "Orientation") {
			return missingBeforeLaterRefusal(advance, "Lessons Learned")
		}
		return advancePhase(advance, "agent judgment: lessons and orientation", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("section", requestPath, "", "Lessons Learned", "captured lessons"), nil)
	}
	if !hasSection(sections, "Orientation") {
		return advancePhase(advance, "agent judgment: lessons and orientation", resultmodel.AdvancePhaseAgentJudgment,
			advanceEvidence("section", requestPath, "", "Orientation", "subsystem-level handback"), nil)
	}
	return advancePhase(advance, finalization.CommandFinalize, resultmodel.AdvancePhaseMechanical,
		resultmodel.AdvanceMissingEvidence{Kind: "file", Path: "<action-authored-finalization-manifest>", Expected: "action-authored finalization manifest"},
		[]string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", requestPath, "--finalization-manifest", "<action-authored-finalization-manifest>"})
}

func classifyArchiveAdvance(target *repositorymodel.RequestFile, advance *resultmodel.AdvanceLifecycleResult) resultmodel.CommandResult {
	requestPath := advance.RequestPath
	switch advance.Status {
	case "completed", "completed-with-issues", "failed", "cancelled":
	default:
		return advanceRefusal(advance.RequestID, []string{requestPath}, "ADVANCE-PHASE-UNKNOWN", "archive request is not terminal", advance)
	}
	if strings.TrimSpace(target.TypedRecord.CompletedAt) == "" {
		return advancePhase(advance, "recover-finalization", resultmodel.AdvancePhaseMechanical,
			advanceEvidence("field", requestPath, "completed_at", "", "canonical terminal timestamp"),
			[]string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"})
	}
	if strings.TrimSpace(target.TypedRecord.ImplementationCommit) == "" {
		advance.Phase = "recover-finalization"
		advance.PhaseKind = resultmodel.AdvancePhaseMechanical
		advance.MissingEvidence = []resultmodel.AdvanceMissingEvidence{advanceEvidence("field", requestPath, "commit", "", "recorded implementation or primary commit")}
		advance.NextArgv = []string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"}
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Advance: advance, Findings: []resultmodel.CommandFinding{{
			Code: "ADVANCE-EVIDENCE-MISSING", Severity: resultmodel.SeverityWarning,
			AffectedIDs: []string{advance.RequestID}, AffectedPaths: []string{requestPath},
			Evidence:   []string{"archived terminal request has no recorded commit provenance"},
			Fixability: resultmodel.FixabilityManual, AutomationStopReason: "finalization provenance is incomplete",
			NextArgv: append([]string(nil), advance.NextArgv...), VerificationArgv: append([]string(nil), advance.VerificationArgv...),
		}}}
	}
	return advancePhase(advance, "complete", resultmodel.AdvancePhaseComplete, resultmodel.AdvanceMissingEvidence{}, nil)
}

func advancePhase(advance *resultmodel.AdvanceLifecycleResult, phase string, phaseKind resultmodel.AdvancePhaseKind, missing resultmodel.AdvanceMissingEvidence, nextArgv []string) resultmodel.CommandResult {
	advance.Phase = phase
	advance.PhaseKind = phaseKind
	if missing.Kind != "" {
		advance.MissingEvidence = []resultmodel.AdvanceMissingEvidence{missing}
	}
	advance.NextArgv = append([]string(nil), nextArgv...)
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Advance: advance}
}

func advanceEvidence(kind, path, field, section, expected string) resultmodel.AdvanceMissingEvidence {
	return resultmodel.AdvanceMissingEvidence{Kind: kind, Path: path, Field: field, Section: section, Expected: expected}
}

func advanceFailure(code, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{evidence},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "advance could not classify the request",
		NextArgv: []string{"do-work-cli", "help"}, VerificationArgv: []string{"do-work-cli", "--format", "json", CommandAdvance, "REQ-NNN"},
	}}}
}

func advanceRefusal(requestID string, paths []string, code, evidence string, advance *resultmodel.AdvanceLifecycleResult) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Advance: advance, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, AffectedIDs: []string{requestID}, AffectedPaths: paths,
		Evidence: []string{evidence}, Fixability: resultmodel.FixabilityRefused,
		AutomationStopReason: "the lifecycle phase cannot be derived safely",
		VerificationArgv:     []string{"do-work-cli", "--format", "json", CommandAdvance, requestID},
	}}}
}

func missingBeforeLaterRefusal(advance *resultmodel.AdvanceLifecycleResult, evidenceName string) resultmodel.CommandResult {
	return advanceRefusal(advance.RequestID, []string{advance.RequestPath}, "ADVANCE-EVIDENCE-MISSING", "later lifecycle evidence exists before "+evidenceName, advance)
}

func advanceRequestPath(requestFile *repositorymodel.RequestFile) string {
	return filepath.ToSlash(filepath.Join("do-work", filepath.FromSlash(requestFile.RelativePath)))
}

func duplicateField(fields map[string]requestmodel.FieldEvidence, fieldName string) bool {
	return fields[fieldName].DuplicateCount > 1
}

func validEstimate(value string) bool {
	minutes, parseError := strconv.Atoi(strings.TrimSpace(value))
	return parseError == nil && minutes > 0
}

func advanceSections(body []byte) (map[string]sectionEvidence, string) {
	matches := advanceHeadingPattern.FindAllSubmatchIndex(body, -1)
	sections := map[string]sectionEvidence{}
	for matchIndex, match := range matches {
		name := string(body[match[2]:match[3]])
		section := sections[name]
		section.count++
		if section.count == 1 {
			section.start = match[0]
			section.end = len(body)
			if matchIndex+1 < len(matches) {
				section.end = matches[matchIndex+1][0]
			}
		}
		sections[name] = section
	}
	for name, section := range sections {
		if section.count > 1 {
			return sections, "duplicate lifecycle section " + name
		}
	}
	return sections, ""
}

func hasSection(sections map[string]sectionEvidence, name string) bool {
	return sections[name].count == 1
}

func hasAnySection(sections map[string]sectionEvidence, names ...string) bool {
	for _, name := range names {
		if hasSection(sections, name) {
			return true
		}
	}
	return false
}

func sectionContains(body []byte, sections map[string]sectionEvidence, name string, pattern *regexp.Regexp) bool {
	section, found := sections[name]
	return found && section.start >= 0 && section.end <= len(body) && pattern.Match(body[section.start:section.end])
}
