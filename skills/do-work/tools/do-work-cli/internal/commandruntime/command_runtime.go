package commandruntime

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type ExecutionContext struct {
	RepositoryRoot string
	Format         resultmodel.OutputFormat
}

type CommandHandler func(ExecutionContext, []string) resultmodel.CommandResult

type CommandRuntime struct {
	output   io.Writer
	handlers map[string]CommandHandler
}

func NewRuntime(output io.Writer, handlers map[string]CommandHandler) *CommandRuntime {
	if output == nil {
		output = io.Discard
	}
	if handlers == nil {
		handlers = map[string]CommandHandler{}
	}
	return &CommandRuntime{output: output, handlers: handlers}
}

func (runtime *CommandRuntime) Run(arguments []string) int {
	context, command, commandArgs, parseFinding := parseGlobalOptions(arguments)
	if parseFinding != nil {
		return runtime.writeResult(context.Format, resultmodel.CommandResult{
			Command:        command,
			Outcome:        resultmodel.OutcomeFailure,
			RepositoryRoot: context.RepositoryRoot,
			Findings:       []resultmodel.CommandFinding{*parseFinding},
		})
	}
	handler, exists := runtime.handlers[command]
	if !exists {
		return runtime.writeResult(context.Format, resultmodel.CommandResult{
			Command:        command,
			Outcome:        resultmodel.OutcomeFailure,
			RepositoryRoot: context.RepositoryRoot,
			Findings: []resultmodel.CommandFinding{usageFinding(
				"UNKNOWN-COMMAND", fmt.Sprintf("command %q is not available", command),
			)},
		})
	}
	result := handler(context, commandArgs)
	result.Command = command
	result.RepositoryRoot = context.RepositoryRoot
	return runtime.writeResult(context.Format, result)
}

func (runtime *CommandRuntime) writeResult(outputFormat resultmodel.OutputFormat, result resultmodel.CommandResult) int {
	output, err := resultmodel.RenderResult(result, outputFormat)
	if err != nil {
		fallback := resultmodel.CommandResult{
			Command:        result.Command,
			Outcome:        resultmodel.OutcomeFailure,
			RepositoryRoot: result.RepositoryRoot,
			Findings: []resultmodel.CommandFinding{usageFinding(
				"OUTPUT-RENDER-FAILED", err.Error(),
			)},
		}
		output, _ = resultmodel.RenderResult(fallback, resultmodel.FormatText)
		result = fallback
	}
	_, _ = runtime.output.Write(output)
	return resultmodel.ExitCode(result.Outcome)
}

func parseGlobalOptions(arguments []string) (ExecutionContext, string, []string, *resultmodel.CommandFinding) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		workingDirectory = "."
	}
	context := ExecutionContext{RepositoryRoot: workingDirectory, Format: resultmodel.FormatText}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if !strings.HasPrefix(argument, "-") || argument == "-" {
			context.RepositoryRoot = absolutePath(context.RepositoryRoot)
			return context, argument, arguments[index+1:], nil
		}
		switch {
		case argument == "--repo-root":
			index++
			if index >= len(arguments) {
				context.RepositoryRoot = absolutePath(context.RepositoryRoot)
				finding := usageFinding("MISSING-OPTION-VALUE", "--repo-root requires a path")
				return context, "", nil, &finding
			}
			context.RepositoryRoot = arguments[index]
		case strings.HasPrefix(argument, "--repo-root="):
			context.RepositoryRoot = strings.TrimPrefix(argument, "--repo-root=")
			if context.RepositoryRoot == "" {
				context.RepositoryRoot = absolutePath(workingDirectory)
				finding := usageFinding("MISSING-OPTION-VALUE", "--repo-root requires a path")
				return context, "", nil, &finding
			}
		case argument == "--format":
			index++
			if index >= len(arguments) {
				context.RepositoryRoot = absolutePath(context.RepositoryRoot)
				finding := usageFinding("MISSING-OPTION-VALUE", "--format requires text or json")
				return context, "", nil, &finding
			}
			context.Format = resultmodel.OutputFormat(arguments[index])
		case strings.HasPrefix(argument, "--format="):
			context.Format = resultmodel.OutputFormat(strings.TrimPrefix(argument, "--format="))
		default:
			context.RepositoryRoot = absolutePath(context.RepositoryRoot)
			finding := usageFinding("UNKNOWN-GLOBAL-OPTION", fmt.Sprintf("global option %q is not available", argument))
			return context, "", nil, &finding
		}
		if context.Format != resultmodel.FormatText && context.Format != resultmodel.FormatJSON {
			context.Format = resultmodel.FormatText
			context.RepositoryRoot = absolutePath(context.RepositoryRoot)
			finding := usageFinding("INVALID-OUTPUT-FORMAT", "--format requires text or json")
			return context, "", nil, &finding
		}
	}
	context.RepositoryRoot = absolutePath(context.RepositoryRoot)
	finding := usageFinding("MISSING-COMMAND", "a command is required")
	return context, "", nil, &finding
}

func usageFinding(code, evidence string) resultmodel.CommandFinding {
	return resultmodel.CommandFinding{
		Code:                 code,
		Severity:             resultmodel.SeverityError,
		Evidence:             []string{evidence},
		Fixability:           resultmodel.FixabilityManual,
		AutomationStopReason: "the command line is not valid",
		NextArgv:             []string{"do-work-cli", "--format", "text", "<command>"},
	}
}

func absolutePath(path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return filepath.Clean(absolute)
}
