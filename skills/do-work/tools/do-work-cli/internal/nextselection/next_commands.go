package nextselection

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	runnerPath, runnerError := blockedRunnerPath(executionContext.RepositoryRoot)
	var runner ProbeRunner = func([]byte, int) (int, error) { return 125, runnerError }
	if runnerError == nil {
		runner = processProbeRunner(runnerPath)
	}
	result := Select(snapshot, graph, options, runner)
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

func blockedRunnerPath(repositoryRoot string) (string, error) {
	candidates := []string{
		filepath.Join(repositoryRoot, "skills", "do-work", "scripts", "run-blocked-check.sh"),
		filepath.Join(repositoryRoot, ".claude", "skills", "do-work", "scripts", "run-blocked-check.sh"),
		filepath.Join(repositoryRoot, ".codex", "skills", "do-work", "scripts", "run-blocked-check.sh"),
	}
	if executablePath, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Clean(filepath.Join(filepath.Dir(executablePath), "..", "..", "scripts", "run-blocked-check.sh")))
	}
	for _, candidate := range candidates {
		if fileInfo, err := os.Stat(candidate); err == nil && fileInfo.Mode().IsRegular() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("shipped run-blocked-check.sh was not found")
}

func processProbeRunner(runnerPath string) ProbeRunner {
	return func(probeBytes []byte, timeoutSeconds int) (int, error) {
		probeFile, err := os.CreateTemp("", "do-work-blocked-check-*.sh")
		if err != nil {
			return 125, err
		}
		probePath := probeFile.Name()
		defer os.Remove(probePath)
		if _, err := probeFile.Write(probeBytes); err != nil {
			probeFile.Close()
			return 125, err
		}
		if err := probeFile.Close(); err != nil {
			return 125, err
		}
		command := exec.Command(runnerPath, probePath, strconv.Itoa(timeoutSeconds))
		err = command.Run()
		if err == nil {
			return 0, nil
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), nil
		}
		return 125, err
	}
}
