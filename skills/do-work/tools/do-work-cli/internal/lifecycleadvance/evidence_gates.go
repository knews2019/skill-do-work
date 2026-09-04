package lifecycleadvance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gateevidence"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type advanceGateInputs struct {
	requestPath    string
	diffRange      string
	gateArgv       []string
	gateExitStatus *int
	phaseArgv      []string
	separatorSeen  bool
}

func executeAdvanceEvidenceGates(executionContext commandruntime.ExecutionContext, target *repositorymodel.RequestFile, projected resultmodel.CommandResult, arguments []string) resultmodel.CommandResult {
	advance := projected.Advance
	inputs, parseError := parseAdvanceGateInputs(arguments)
	if parseError != nil {
		return advanceRefusal(advance.RequestID, []string{advance.RequestPath}, "ADVANCE-GATE-USAGE", parseError.Error(), advance)
	}
	if inputs.requestPath != "" && inputs.requestPath != advance.RequestPath {
		return advanceRefusal(advance.RequestID, []string{advance.RequestPath, inputs.requestPath}, "ADVANCE-EVIDENCE-MISMATCH", "supplied request path does not match the discovered request", advance)
	}
	if inputs.gateExitStatus != nil && len(inputs.gateArgv) == 0 {
		return advanceRefusal(advance.RequestID, []string{advance.RequestPath}, "ADVANCE-GATE-USAGE", "--gate-exit-status requires at least one --gate-arg", advance)
	}

	records := []resultmodel.AdvanceGateRecord{}
	switch advance.Phase {
	case corehelpers.CommandEstimateP50:
		if inputs.diffRange != "" || len(inputs.gateArgv) > 0 || inputs.gateExitStatus != nil {
			return irrelevantAdvanceGateInput(advance, "estimate accepts only estimator argv after --")
		}
		if !inputs.separatorSeen || len(inputs.phaseArgv) == 0 {
			records = append(records, missingAdvanceGateInput(advance, corehelpers.CommandEstimateP50, "estimator argv after --"))
		} else {
			records = append(records, composeCoreGate(executionContext, advance, corehelpers.CommandEstimateP50, inputs.phaseArgv, resultmodel.AdvanceGateExecuted))
		}
	case corehelpers.CommandPreflight:
		if inputs.diffRange != "" {
			return irrelevantAdvanceGateInput(advance, "preflight does not accept --diff-range")
		}
		if !inputs.separatorSeen {
			records = append(records, missingAdvanceGateInput(advance, corehelpers.CommandPreflight, "resolved test argv after --; use a bare -- when no test command exists"))
		} else {
			records = append(records, composeCoreGate(executionContext, advance, corehelpers.CommandPreflight, inputs.phaseArgv, resultmodel.AdvanceGateExecuted))
		}
		records = append(records, composeGreenGate(executionContext, advance, inputs))
	case corehelpers.CommandQualify:
		if inputs.separatorSeen || len(inputs.gateArgv) > 0 || inputs.gateExitStatus != nil {
			return irrelevantAdvanceGateInput(advance, "qualification accepts only --diff-range and the discovered request path")
		}
		if inputs.diffRange == "" {
			records = append(records, missingAdvanceGateInput(advance, corehelpers.CommandQualify, "exact --diff-range <pre>..<merge>"))
		} else {
			qualifyArguments := []string{"--request-path", advance.RequestPath, "--diff-range", inputs.diffRange}
			records = append(records, composeCoreGate(executionContext, advance, corehelpers.CommandQualify, qualifyArguments, resultmodel.AdvanceGateMergedRange))
			if target.TypedRecord.RouteValue != "A" {
				records = append(records, composeCoreGate(executionContext, advance, corehelpers.CommandScopeDrift, []string{"--request-path", advance.RequestPath}, resultmodel.AdvanceGateMergedRange))
			}
		}
	case "test-gate":
		if inputs.diffRange != "" {
			return irrelevantAdvanceGateInput(advance, "test gate reads the saved baseline and accepts probe argv after --")
		}
		if target.TypedRecord.RouteValue != "A" {
			records = append(records, composeCoreGate(executionContext, advance, corehelpers.CommandScopeDrift, []string{"--request-path", advance.RequestPath}, resultmodel.AdvanceGateMergedRange))
		}
		if !inputs.separatorSeen || len(inputs.phaseArgv) == 0 {
			records = append(records, missingAdvanceGateInput(advance, corehelpers.CommandBlockedCheck, "run-blocked-check argv after --"))
		} else {
			probeArguments := append([]string(nil), inputs.phaseArgv...)
			if !hasAdvanceOption(probeArguments, "--baseline-json") && !hasAdvanceOption(probeArguments, "--baseline-failures") {
				probeArguments = append(probeArguments,
					"--baseline-json", "do-work/working/baseline.json",
					"--baseline-failures", "do-work/working/baseline-failures.txt",
				)
			}
			records = append(records, composeCoreGate(executionContext, advance, corehelpers.CommandBlockedCheck, probeArguments, resultmodel.AdvanceGateBaselineRecord))
		}
		records = append(records, composeGreenGate(executionContext, advance, inputs))
	default:
		return irrelevantAdvanceGateInput(advance, "the current lifecycle phase requires agent judgment, not gate argv")
	}
	advance.GateRecords = records
	advance.NextArgv = []string{}
	advance.MissingEvidence = []resultmodel.AdvanceMissingEvidence{}
	return aggregateAdvanceGateResult(advance, records)
}

func parseAdvanceGateInputs(arguments []string) (advanceGateInputs, error) {
	inputs := advanceGateInputs{}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--" {
			inputs.separatorSeen = true
			inputs.phaseArgv = append([]string(nil), arguments[index+1:]...)
			return inputs, nil
		}
		var value string
		var valueError error
		switch {
		case argument == "--request-path" || strings.HasPrefix(argument, "--request-path="):
			value, valueError = advanceOptionValue(arguments, &index, "--request-path")
			if inputs.requestPath != "" {
				return inputs, fmt.Errorf("--request-path may be supplied only once")
			}
			inputs.requestPath = value
		case argument == "--diff-range" || strings.HasPrefix(argument, "--diff-range="):
			value, valueError = advanceOptionValue(arguments, &index, "--diff-range")
			if inputs.diffRange != "" {
				return inputs, fmt.Errorf("--diff-range may be supplied only once")
			}
			inputs.diffRange = value
		case argument == "--gate-arg" || strings.HasPrefix(argument, "--gate-arg="):
			value, valueError = advanceOptionValue(arguments, &index, "--gate-arg")
			inputs.gateArgv = append(inputs.gateArgv, value)
		case argument == "--gate-exit-status" || strings.HasPrefix(argument, "--gate-exit-status="):
			value, valueError = advanceOptionValue(arguments, &index, "--gate-exit-status")
			if inputs.gateExitStatus != nil {
				return inputs, fmt.Errorf("--gate-exit-status may be supplied only once")
			}
			parsed, parseError := strconv.Atoi(value)
			if parseError != nil {
				return inputs, fmt.Errorf("--gate-exit-status requires an integer")
			}
			inputs.gateExitStatus = &parsed
		default:
			return inputs, fmt.Errorf("unknown advance gate option %q", argument)
		}
		if valueError != nil {
			return inputs, valueError
		}
		if strings.TrimSpace(value) == "" {
			return inputs, fmt.Errorf("%s requires a nonblank value", strings.SplitN(argument, "=", 2)[0])
		}
	}
	return inputs, nil
}

func advanceOptionValue(arguments []string, index *int, name string) (string, error) {
	argument := arguments[*index]
	if strings.HasPrefix(argument, name+"=") {
		return strings.TrimPrefix(argument, name+"="), nil
	}
	*index = *index + 1
	if *index >= len(arguments) {
		return "", fmt.Errorf("%s requires a value", name)
	}
	return arguments[*index], nil
}

func composeCoreGate(executionContext commandruntime.ExecutionContext, advance *resultmodel.AdvanceLifecycleResult, gateID string, arguments []string, provenance resultmodel.AdvanceGateProvenance) resultmodel.AdvanceGateRecord {
	handler := corehelpers.Handlers()[gateID]
	subordinate := handler(executionContext, append([]string(nil), arguments...))
	record := boundAdvanceGateRecord(advance, gateID, provenance, subordinate)
	if gateID == corehelpers.CommandBlockedCheck && subordinate.FocusedTest != nil {
		record.FocusedTest = subordinate.FocusedTest
		switch subordinate.FocusedTest.BaselineState {
		case resultmodel.FocusedBaselineGreen, resultmodel.FocusedBaselineMatchingRed:
			record.State = resultmodel.AdvanceGateSatisfied
		case resultmodel.FocusedBaselineMissing, resultmodel.FocusedBaselineUnusable, resultmodel.FocusedBaselineNewRed:
			if record.State != resultmodel.AdvanceGateFailed {
				record.State = resultmodel.AdvanceGateFindings
			}
		}
	}
	return record
}

func composeGreenGate(executionContext commandruntime.ExecutionContext, advance *resultmodel.AdvanceLifecycleResult, inputs advanceGateInputs) resultmodel.AdvanceGateRecord {
	if len(inputs.gateArgv) == 0 {
		return missingAdvanceGateInput(advance, gateevidence.CommandCheckGreenGate, "one --gate-arg per canonical repository-gate argv token")
	}
	handlerName := gateevidence.CommandCheckGreenGate
	handlerArguments := append([]string{"--"}, inputs.gateArgv...)
	provenance := resultmodel.AdvanceGateExistingEvidence
	if inputs.gateExitStatus != nil {
		handlerName = gateevidence.CommandRecordGreenGate
		handlerArguments = append([]string{"--gate-exit-status", strconv.Itoa(*inputs.gateExitStatus), "--"}, inputs.gateArgv...)
		provenance = resultmodel.AdvanceGateExecuted
	}
	subordinate := gateevidence.Handlers()[handlerName](executionContext, handlerArguments)
	record := boundAdvanceGateRecord(advance, "green-gate", provenance, subordinate)
	record.GreenGate = subordinate.GateEvidence
	if inputs.gateExitStatus == nil && subordinate.Outcome == resultmodel.OutcomeSuccess && subordinate.GateEvidence != nil && !subordinate.GateEvidence.Matches {
		record.State = resultmodel.AdvanceGateNeedsInput
		record.Outcome = resultmodel.OutcomeFindings
		record.NextArgv = append([]string(nil), inputs.gateArgv...)
		record.VerificationArgv = advanceGateVerificationArgv(advance, inputs, true)
		record.Findings = append(record.Findings, boundAdvanceFinding(advance, resultmodel.CommandFinding{
			Code: "ADVANCE-GREEN-GATE-DIRECT-RUN-REQUIRED", Severity: resultmodel.SeverityWarning,
			Evidence:   []string{"no reusable green record matches the exact gate argv at the current revision"},
			Fixability: resultmodel.FixabilityManual, AutomationStopReason: "a direct gate run is required",
			NextArgv: append([]string(nil), inputs.gateArgv...), VerificationArgv: append([]string(nil), record.VerificationArgv...),
		}))
	}
	return record
}

func boundAdvanceGateRecord(advance *resultmodel.AdvanceLifecycleResult, gateID string, provenance resultmodel.AdvanceGateProvenance, subordinate resultmodel.CommandResult) resultmodel.AdvanceGateRecord {
	state := resultmodel.AdvanceGateSatisfied
	switch subordinate.Outcome {
	case resultmodel.OutcomeFindings:
		state = resultmodel.AdvanceGateFindings
	case resultmodel.OutcomeFailure, resultmodel.OutcomeRefused, resultmodel.OutcomeRolledBack, resultmodel.OutcomeRisk:
		state = resultmodel.AdvanceGateFailed
	}
	findings := make([]resultmodel.CommandFinding, 0, len(subordinate.Findings))
	for _, finding := range subordinate.Findings {
		findings = append(findings, boundAdvanceFinding(advance, finding))
	}
	outputLines := []string{}
	if subordinate.ExactTextOutput != nil {
		for _, line := range strings.Split(strings.TrimSpace(*subordinate.ExactTextOutput), "\n") {
			if strings.TrimSpace(line) != "" {
				outputLines = append(outputLines, line)
			}
		}
	}
	return resultmodel.AdvanceGateRecord{
		RequestID: advance.RequestID, RequestPath: advance.RequestPath, GateID: gateID,
		Provenance: provenance, State: state, Outcome: subordinate.Outcome,
		Findings: findings, Changes: append([]resultmodel.RecordedChange(nil), subordinate.Changes...),
		OutputLines: outputLines, VerificationArgv: append([]string(nil), advance.VerificationArgv...),
	}
}

func missingAdvanceGateInput(advance *resultmodel.AdvanceLifecycleResult, gateID, expected string) resultmodel.AdvanceGateRecord {
	next := []string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", advance.RequestPath, "--", "<" + expected + ">"}
	finding := boundAdvanceFinding(advance, resultmodel.CommandFinding{
		Code: "ADVANCE-GATE-INPUT-REQUIRED", Severity: resultmodel.SeverityWarning,
		Evidence: []string{gateID + " requires " + expected}, Fixability: resultmodel.FixabilityManual,
		AutomationStopReason: "the action must supply its judgment-owned input", NextArgv: next,
		VerificationArgv: append([]string(nil), advance.VerificationArgv...),
	})
	return resultmodel.AdvanceGateRecord{
		RequestID: advance.RequestID, RequestPath: advance.RequestPath, GateID: gateID,
		Provenance: resultmodel.AdvanceGateExecuted, State: resultmodel.AdvanceGateNeedsInput,
		Outcome: resultmodel.OutcomeFindings, Findings: []resultmodel.CommandFinding{finding}, NextArgv: next,
		VerificationArgv: append([]string(nil), advance.VerificationArgv...),
	}
}

func irrelevantAdvanceGateInput(advance *resultmodel.AdvanceLifecycleResult, evidence string) resultmodel.CommandResult {
	return advanceRefusal(advance.RequestID, []string{advance.RequestPath}, "ADVANCE-GATE-INPUT-IRRELEVANT", evidence, advance)
}

func boundAdvanceFinding(advance *resultmodel.AdvanceLifecycleResult, finding resultmodel.CommandFinding) resultmodel.CommandFinding {
	if !containsAdvanceString(finding.AffectedIDs, advance.RequestID) {
		finding.AffectedIDs = append([]string{advance.RequestID}, finding.AffectedIDs...)
	}
	if !containsAdvanceString(finding.AffectedPaths, advance.RequestPath) {
		finding.AffectedPaths = append([]string{advance.RequestPath}, finding.AffectedPaths...)
	}
	return finding
}

func aggregateAdvanceGateResult(advance *resultmodel.AdvanceLifecycleResult, records []resultmodel.AdvanceGateRecord) resultmodel.CommandResult {
	outcome := resultmodel.OutcomeSuccess
	findings := []resultmodel.CommandFinding{}
	changes := []resultmodel.RecordedChange{}
	for _, record := range records {
		findings = append(findings, record.Findings...)
		changes = append(changes, record.Changes...)
		switch record.State {
		case resultmodel.AdvanceGateFailed:
			outcome = resultmodel.OutcomeFailure
		case resultmodel.AdvanceGateNeedsInput, resultmodel.AdvanceGateFindings:
			if outcome != resultmodel.OutcomeFailure {
				outcome = resultmodel.OutcomeFindings
			}
		}
	}
	return resultmodel.CommandResult{Outcome: outcome, Findings: findings, Changes: changes, Advance: advance}
}

func advanceGateVerificationArgv(advance *resultmodel.AdvanceLifecycleResult, inputs advanceGateInputs, recordGreen bool) []string {
	arguments := []string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", advance.RequestPath}
	if inputs.diffRange != "" {
		arguments = append(arguments, "--diff-range", inputs.diffRange)
	}
	if recordGreen {
		arguments = append(arguments, "--gate-exit-status", "0")
	}
	for _, argument := range inputs.gateArgv {
		arguments = append(arguments, "--gate-arg", argument)
	}
	if inputs.separatorSeen {
		arguments = append(arguments, "--")
		arguments = append(arguments, inputs.phaseArgv...)
	}
	return arguments
}

func hasAdvanceOption(arguments []string, option string) bool {
	for _, argument := range arguments {
		if argument == option || strings.HasPrefix(argument, option+"=") {
			return true
		}
	}
	return false
}

func containsAdvanceString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
