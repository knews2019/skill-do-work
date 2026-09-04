package lifecycleadvance

import (
	"fmt"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/finalization"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type advanceFinalizationInputs struct {
	requestPath  string
	manifestPath string
}

func executeAdvanceFinalization(executionContext commandruntime.ExecutionContext, projected resultmodel.CommandResult, arguments []string) resultmodel.CommandResult {
	advance := projected.Advance
	inputs, err := parseAdvanceFinalizationInputs(arguments)
	if err != nil {
		return advanceRefusal(advance.RequestID, []string{advance.RequestPath}, "ADVANCE-FINALIZATION-USAGE", err.Error(), advance)
	}
	if inputs.requestPath != "" && inputs.requestPath != advance.RequestPath {
		return advanceRefusal(advance.RequestID, []string{advance.RequestPath, inputs.requestPath}, "ADVANCE-EVIDENCE-MISMATCH", "supplied request path does not match the discovered request", advance)
	}
	if inputs.manifestPath == "" {
		return projected
	}
	advance.NextArgv = []string{}
	advance.MissingEvidence = []resultmodel.AdvanceMissingEvidence{}
	result := finalization.FinalizeBound(executionContext, inputs.manifestPath, advance.RequestID, advance.RequestPath)
	result.Advance = advance
	return result
}

func parseAdvanceFinalizationInputs(arguments []string) (advanceFinalizationInputs, error) {
	inputs := advanceFinalizationInputs{}
	requestPathSeen := false
	manifestPathSeen := false
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		var value string
		var err error
		switch {
		case argument == "--request-path" || strings.HasPrefix(argument, "--request-path="):
			if requestPathSeen {
				return inputs, fmt.Errorf("--request-path may be supplied only once")
			}
			requestPathSeen = true
			value, err = advanceOptionValue(arguments, &index, "--request-path")
			if err == nil && value == "" {
				err = fmt.Errorf("--request-path requires a non-empty value")
			}
			inputs.requestPath = value
		case argument == "--finalization-manifest" || strings.HasPrefix(argument, "--finalization-manifest="):
			if manifestPathSeen {
				return inputs, fmt.Errorf("--finalization-manifest may be supplied only once")
			}
			manifestPathSeen = true
			value, err = advanceOptionValue(arguments, &index, "--finalization-manifest")
			if err == nil && value == "" {
				err = fmt.Errorf("--finalization-manifest requires a non-empty value")
			}
			inputs.manifestPath = value
		default:
			return inputs, fmt.Errorf("unknown finalization input %q", argument)
		}
		if err != nil {
			return inputs, err
		}
	}
	return inputs, nil
}
