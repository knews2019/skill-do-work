package heavyverification

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandPlanHeavyVerification = "plan-heavy-verification"
	CommandPlanHeavyRevalidation = "plan-heavy-revalidation"
	CommandRunHeavyVerification  = "run-heavy-verification"
)

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandPlanHeavyVerification: handlePlanHeavyVerification,
		CommandPlanHeavyRevalidation: handlePlanHeavyRevalidation,
		CommandRunHeavyVerification:  handleRunHeavyVerification,
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

func handleRunHeavyVerification(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	manifestPath, laneIDs, laneTimeout, evidenceReuse, err := parseRunArguments(arguments)
	if err != nil {
		return runFailure("HEAVY-RUN-USAGE", err)
	}
	// Lane output is teed to stderr because stdout carries the command result.
	run, findings, err := RunLanes(LaneRunRequest{
		RepositoryRoot: executionContext.RepositoryRoot, ManifestPath: manifestPath,
		LaneIDs: laneIDs, LaneTimeout: laneTimeout, LaneOutputWriter: os.Stderr,
		EvidenceReuse: evidenceReuse,
	})
	if err != nil {
		return runFailure(LaneRunRefusalCode(err), err)
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, HeavyVerificationRun: &run, Findings: findings}
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

// parseRunArguments accepts the lanes to verify, the per-lane time bound, and
// whether stored evidence may stand in for an execution. --lane repeats; a
// repeated id would run the lane twice and produce two records under one id,
// which per-lane evidence cannot carry. Evidence reuse is on by default, which
// is what lets a drain skip lanes nothing changed under; --no-evidence-reuse
// forces every named lane to execute and refresh its record.
func parseRunArguments(arguments []string) (string, []string, time.Duration, bool, error) {
	manifestPath := "_dev/tests/heavy-lanes.json"
	laneIDs := []string{}
	laneTimeoutSeconds := defaultLaneTimeoutSeconds
	evidenceReuse := true
	seen := map[string]bool{}
	seenLaneIDs := map[string]bool{}
	for argumentIndex := 0; argumentIndex < len(arguments); argumentIndex++ {
		argument := arguments[argumentIndex]
		if argument == "--no-evidence-reuse" {
			if seen[argument] {
				return "", nil, 0, false, fmt.Errorf("--no-evidence-reuse may be supplied only once")
			}
			seen[argument] = true
			evidenceReuse = false
			continue
		}
		optionName, optionValue, hasInlineValue := strings.Cut(argument, "=")
		switch optionName {
		case "--manifest", "--lane", "--lane-timeout-seconds":
			if optionName != "--lane" && seen[optionName] {
				return "", nil, 0, false, fmt.Errorf("%s may be supplied only once", optionName)
			}
			seen[optionName] = true
			if !hasInlineValue {
				argumentIndex++
				if argumentIndex >= len(arguments) {
					return "", nil, 0, false, fmt.Errorf("%s requires a value", optionName)
				}
				optionValue = arguments[argumentIndex]
			}
			if strings.TrimSpace(optionValue) == "" {
				return "", nil, 0, false, fmt.Errorf("%s requires a value", optionName)
			}
			switch optionName {
			case "--manifest":
				manifestPath = optionValue
			case "--lane":
				if seenLaneIDs[optionValue] {
					return "", nil, 0, false, fmt.Errorf("--lane %s may be supplied only once", optionValue)
				}
				seenLaneIDs[optionValue] = true
				laneIDs = append(laneIDs, optionValue)
			case "--lane-timeout-seconds":
				parsedSeconds, parseError := strconv.Atoi(optionValue)
				if parseError != nil || parsedSeconds <= 0 {
					return "", nil, 0, false, fmt.Errorf("--lane-timeout-seconds requires a positive whole number of seconds")
				}
				laneTimeoutSeconds = parsedSeconds
			}
		default:
			return "", nil, 0, false, fmt.Errorf("unknown %s option %q", CommandRunHeavyVerification, argument)
		}
	}
	if len(laneIDs) == 0 {
		return "", nil, 0, false, fmt.Errorf("usage: %s [--manifest <path>] --lane <id> [--lane <id>...] [--lane-timeout-seconds <seconds>] [--no-evidence-reuse]", CommandRunHeavyVerification)
	}
	return manifestPath, laneIDs, time.Duration(laneTimeoutSeconds) * time.Second, evidenceReuse, nil
}

func runFailure(code string, err error) resultmodel.CommandResult {
	finding := resultmodel.CommandFinding{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{err.Error()},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "heavy lanes cannot be run safely",
		NextArgv:         []string{"git", "status", "--short"},
		VerificationArgv: []string{"do-work-cli", "--format", "json", CommandRunHeavyVerification},
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{finding}}
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
