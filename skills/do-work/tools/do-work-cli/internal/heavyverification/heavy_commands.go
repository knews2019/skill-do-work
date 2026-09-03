package heavyverification

import (
	"fmt"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const CommandPlanHeavyVerification = "plan-heavy-verification"

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{CommandPlanHeavyVerification: handlePlanHeavyVerification}
}

func handlePlanHeavyVerification(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	manifestPath, baseRevision, targetRevision, forceAll, err := parsePlanArguments(arguments)
	if err != nil {
		return planFailure("HEAVY-PLAN-USAGE", err)
	}
	plan, err := Plan(executionContext.RepositoryRoot, manifestPath, baseRevision, targetRevision, forceAll)
	if err != nil {
		return planFailure("HEAVY-PLAN-UNVERIFIABLE", err)
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, HeavyVerification: &plan}
}

func parsePlanArguments(arguments []string) (string, string, string, bool, error) {
	manifestPath := "_dev/tests/heavy-lanes.json"
	baseRevision := ""
	targetRevision := ""
	forceAll := false
	seen := map[string]bool{}
	for argumentIndex := 0; argumentIndex < len(arguments); argumentIndex++ {
		argument := arguments[argumentIndex]
		if argument == "--force-all" {
			if seen[argument] {
				return "", "", "", false, fmt.Errorf("--force-all may be supplied only once")
			}
			seen[argument] = true
			forceAll = true
			continue
		}
		optionName, optionValue, hasInlineValue := strings.Cut(argument, "=")
		switch optionName {
		case "--manifest", "--base-revision", "--target-revision":
			if seen[optionName] {
				return "", "", "", false, fmt.Errorf("%s may be supplied only once", optionName)
			}
			seen[optionName] = true
			if !hasInlineValue {
				argumentIndex++
				if argumentIndex >= len(arguments) {
					return "", "", "", false, fmt.Errorf("%s requires a value", optionName)
				}
				optionValue = arguments[argumentIndex]
			}
			if strings.TrimSpace(optionValue) == "" {
				return "", "", "", false, fmt.Errorf("%s requires a value", optionName)
			}
			switch optionName {
			case "--manifest":
				manifestPath = optionValue
			case "--base-revision":
				baseRevision = optionValue
			case "--target-revision":
				targetRevision = optionValue
			}
		default:
			return "", "", "", false, fmt.Errorf("unknown plan-heavy-verification option %q", argument)
		}
	}
	if baseRevision == "" || targetRevision == "" {
		return "", "", "", false, fmt.Errorf("usage: %s [--manifest <path>] --base-revision <revision> --target-revision <revision> [--force-all]", CommandPlanHeavyVerification)
	}
	return manifestPath, baseRevision, targetRevision, forceAll, nil
}

func planFailure(code string, err error) resultmodel.CommandResult {
	finding := resultmodel.CommandFinding{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{err.Error()},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "heavy verification cannot be planned safely",
		NextArgv: []string{"git", "status", "--short"},
		VerificationArgv: []string{"do-work-cli", "--format", "json", CommandPlanHeavyVerification,
			"--base-revision", "<base>", "--target-revision", "<target>"},
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{finding}}
}
