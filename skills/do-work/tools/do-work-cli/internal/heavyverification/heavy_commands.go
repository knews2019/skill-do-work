package heavyverification

import (
	"fmt"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandPlanHeavyVerification = "plan-heavy-verification"
	CommandPlanHeavyRevalidation = "plan-heavy-revalidation"
)

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandPlanHeavyVerification: handlePlanHeavyVerification,
		CommandPlanHeavyRevalidation: handlePlanHeavyRevalidation,
	}
}

func handlePlanHeavyVerification(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	manifestPath, baseRevision, targetRevision, forceAll, err := parsePlanArguments(arguments)
	if err != nil {
		return planFailure(CommandPlanHeavyVerification, "HEAVY-PLAN-USAGE", err)
	}
	plan, err := Plan(executionContext.RepositoryRoot, manifestPath, baseRevision, targetRevision, forceAll)
	if err != nil {
		return planFailure(CommandPlanHeavyVerification, "HEAVY-PLAN-UNVERIFIABLE", err)
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, HeavyVerification: &plan}
}

func handlePlanHeavyRevalidation(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	manifestPath, sourceRanges, executionRevision, forceAll, err := parseRevalidationArguments(arguments)
	if err != nil {
		return planFailure(CommandPlanHeavyRevalidation, "HEAVY-REVALIDATION-USAGE", err)
	}
	plan, err := PlanRevalidation(executionContext.RepositoryRoot, manifestPath, sourceRanges, executionRevision, forceAll)
	if err != nil {
		return planFailure(CommandPlanHeavyRevalidation, "HEAVY-REVALIDATION-UNVERIFIABLE", err)
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

func parseRevalidationArguments(arguments []string) (string, []resultmodel.HeavySourceRange, string, bool, error) {
	manifestPath := "_dev/tests/heavy-lanes.json"
	executionRevision := ""
	forceAll := false
	sourceRanges := []resultmodel.HeavySourceRange{}
	seen := map[string]bool{}
	for argumentIndex := 0; argumentIndex < len(arguments); argumentIndex++ {
		argument := arguments[argumentIndex]
		if argument == "--force-all" {
			if seen[argument] {
				return "", nil, "", false, fmt.Errorf("--force-all may be supplied only once")
			}
			seen[argument] = true
			forceAll = true
			continue
		}
		optionName, optionValue, hasInlineValue := strings.Cut(argument, "=")
		switch optionName {
		case "--manifest", "--execution-revision", "--source-range":
			if optionName != "--source-range" && seen[optionName] {
				return "", nil, "", false, fmt.Errorf("%s may be supplied only once", optionName)
			}
			seen[optionName] = true
			if !hasInlineValue {
				argumentIndex++
				if argumentIndex >= len(arguments) {
					return "", nil, "", false, fmt.Errorf("%s requires a value", optionName)
				}
				optionValue = arguments[argumentIndex]
			}
			if strings.TrimSpace(optionValue) == "" {
				return "", nil, "", false, fmt.Errorf("%s requires a value", optionName)
			}
			switch optionName {
			case "--manifest":
				manifestPath = optionValue
			case "--execution-revision":
				executionRevision = optionValue
			case "--source-range":
				baseRevision, targetRevision, found := strings.Cut(optionValue, "..")
				if !found || strings.TrimSpace(baseRevision) == "" || strings.TrimSpace(targetRevision) == "" || strings.Contains(targetRevision, "..") {
					return "", nil, "", false, fmt.Errorf("--source-range requires <base>..<target>")
				}
				sourceRanges = append(sourceRanges, resultmodel.HeavySourceRange{BaseRevision: baseRevision, TargetRevision: targetRevision})
			}
		default:
			return "", nil, "", false, fmt.Errorf("unknown plan-heavy-revalidation option %q", argument)
		}
	}
	if len(sourceRanges) == 0 || executionRevision == "" {
		return "", nil, "", false, fmt.Errorf("usage: %s [--manifest <path>] --source-range <base>..<target> [--source-range ...] --execution-revision <revision> [--force-all]", CommandPlanHeavyRevalidation)
	}
	return manifestPath, sourceRanges, executionRevision, forceAll, nil
}

func planFailure(command, code string, err error) resultmodel.CommandResult {
	finding := resultmodel.CommandFinding{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{err.Error()},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "heavy verification cannot be planned safely",
		NextArgv:         []string{"git", "status", "--short"},
		VerificationArgv: []string{"do-work-cli", "--format", "json", command},
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{finding}}
}
