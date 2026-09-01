package nextselection

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const CommandNext = "next"

var discoverRepository = repositorymodel.DiscoverRepository

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{CommandNext: handleNext}
}

func handleNext(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, parseError := parseNextOptions(arguments)
	if parseError != nil {
		return nextFailure("NEXT-USAGE", parseError.Error())
	}
	snapshot, discoveryError := discoverRepository(executionContext.RepositoryRoot)
	if discoveryError != nil {
		return nextFailure("NEXT-DISCOVERY-FAILED", discoveryError.Error())
	}
	graph := dependencygraph.BuildGraph(snapshot)
	result := Select(snapshot, graph, options, RunBlockedProbe)
	result.RepositoryRoot = executionContext.RepositoryRoot
	return result
}

func parseNextOptions(arguments []string) (SelectionOptions, error) {
	options := SelectionOptions{}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		switch {
		case argument == "--skip-impact-negligible":
			options.SkipImpactNegligible = true
		case argument == "--simple":
			options.SimpleOnly = true
		case argument == "--wave" || strings.HasPrefix(argument, "--wave="):
			value, consumed, err := optionInteger(argument, arguments[index+1:], "--wave", false)
			if err != nil {
				return options, err
			}
			index += consumed
			options.WaveDepth = &value
		case argument == "--fan-out" || strings.HasPrefix(argument, "--fan-out="):
			if argument == "--fan-out" && (index+1 >= len(arguments) || strings.HasPrefix(arguments[index+1], "-") || hasIDPrefix(arguments[index+1], "REQ-") || hasIDPrefix(arguments[index+1], "UR-")) {
				value := 2
				options.FanOutLimit = &value
				continue
			}
			value, consumed, err := optionInteger(argument, arguments[index+1:], "--fan-out", true)
			if err != nil {
				return options, err
			}
			index += consumed
			options.FanOutLimit = &value
		default:
			options.TargetTokens = append(options.TargetTokens, argument)
		}
	}
	if err := validateTargetTokens(options.TargetTokens); err != nil {
		return options, err
	}
	if options.WaveDepth != nil && len(options.TargetTokens) > 0 {
		return options, fmt.Errorf("--wave cannot be combined with REQ or UR targets")
	}
	if options.SimpleOnly && (len(options.TargetTokens) > 0 || options.WaveDepth != nil || options.FanOutLimit != nil) {
		return options, fmt.Errorf("--simple computes its own complete set and cannot be combined with targets, --wave, or --fan-out")
	}
	return options, nil
}

func optionInteger(argument string, remaining []string, name string, positive bool) (int, int, error) {
	valueText := ""
	consumed := 0
	if strings.HasPrefix(argument, name+"=") {
		valueText = strings.TrimPrefix(argument, name+"=")
	} else {
		if len(remaining) == 0 {
			return 0, 0, fmt.Errorf("%s requires an integer", name)
		}
		valueText = remaining[0]
		consumed = 1
	}
	value, err := strconv.Atoi(valueText)
	if err != nil || value < 0 || (positive && value == 0) {
		return 0, 0, fmt.Errorf("%s requires %s integer", name, map[bool]string{true: "a positive", false: "a non-negative"}[positive])
	}
	return value, consumed, nil
}

func nextFailure(code, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{evidence},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "queue selection could not run",
		NextArgv: []string{"do-work-cli", "--format", "text", "next"}, NextJustRecipe: "do-work-next",
		VerificationArgv: []string{"do-work-cli", "--format", "json", "next"},
	}}}
}
