package commandruntime

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	context, command, commandArgs, parseFinding := runtime.parseGlobalOptions(arguments)
	if parseFinding != nil {
		return runtime.writeResult(context.Format, resultmodel.CommandResult{
			Command:        command,
			Outcome:        resultmodel.OutcomeFailure,
			RepositoryRoot: context.RepositoryRoot,
			Findings:       []resultmodel.CommandFinding{*parseFinding},
		})
	}
	if command == "help" {
		if len(commandArgs) > 0 {
			return runtime.writeResult(context.Format, resultmodel.CommandResult{
				Command: command, Outcome: resultmodel.OutcomeFailure, RepositoryRoot: context.RepositoryRoot,
				Findings: []resultmodel.CommandFinding{runtime.usageFinding(command, "HELP-ARGUMENTS", "help accepts no command arguments")},
			})
		}
		return runtime.writeResult(context.Format, resultmodel.CommandResult{
			Command: command, Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: context.RepositoryRoot,
			AvailableCommands: runtime.availableCommands(),
		})
	}
	handler, exists := runtime.handlers[command]
	if !exists {
		return runtime.writeResult(context.Format, resultmodel.CommandResult{
			Command:        command,
			Outcome:        resultmodel.OutcomeFailure,
			RepositoryRoot: context.RepositoryRoot,
			Findings: []resultmodel.CommandFinding{runtime.usageFinding(
				command, "UNKNOWN-COMMAND", fmt.Sprintf("command %q is not available: available commands are %s",
					command, strings.Join(runtime.availableCommands(), ", ")),
			)},
		})
	}
	result := handler(context, commandArgs)
	result.Command = command
	result.RepositoryRoot = context.RepositoryRoot
	exitCode := runtime.writeResult(context.Format, result)
	if result.ExitCodeOverride > 0 && result.ExitCodeOverride <= 255 {
		return result.ExitCodeOverride
	}
	return exitCode
}

func (runtime *CommandRuntime) writeResult(outputFormat resultmodel.OutputFormat, result resultmodel.CommandResult) int {
	output, err := resultmodel.RenderResult(result, outputFormat)
	if err != nil {
		fallback := resultmodel.CommandResult{
			Command:        result.Command,
			Outcome:        resultmodel.OutcomeFailure,
			RepositoryRoot: result.RepositoryRoot,
			Findings: []resultmodel.CommandFinding{runtime.usageFinding(
				result.Command, "OUTPUT-RENDER-FAILED", err.Error(),
			)},
		}
		output, _ = resultmodel.RenderResult(fallback, resultmodel.FormatText)
		result = fallback
	}
	_, _ = runtime.output.Write(output)
	return resultmodel.ExitCode(result.Outcome)
}

// parseGlobalOptions scans ahead for the command token before reporting an option error, so
// even a finding raised while parsing the global options names the command the user was
// reaching for rather than an unpasteable placeholder.
func (runtime *CommandRuntime) parseGlobalOptions(arguments []string) (ExecutionContext, string, []string, *resultmodel.CommandFinding) {
	intendedCommand := scanForCommandToken(arguments)
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
				finding := runtime.usageFinding(intendedCommand, "MISSING-OPTION-VALUE", "--repo-root requires a path")
				return context, "", nil, &finding
			}
			context.RepositoryRoot = arguments[index]
		case strings.HasPrefix(argument, "--repo-root="):
			context.RepositoryRoot = strings.TrimPrefix(argument, "--repo-root=")
			if context.RepositoryRoot == "" {
				context.RepositoryRoot = absolutePath(workingDirectory)
				finding := runtime.usageFinding(intendedCommand, "MISSING-OPTION-VALUE", "--repo-root requires a path")
				return context, "", nil, &finding
			}
		case argument == "--format":
			index++
			if index >= len(arguments) {
				context.RepositoryRoot = absolutePath(context.RepositoryRoot)
				finding := runtime.usageFinding(intendedCommand, "MISSING-OPTION-VALUE", "--format requires text or json")
				return context, "", nil, &finding
			}
			context.Format = resultmodel.OutputFormat(arguments[index])
		case strings.HasPrefix(argument, "--format="):
			context.Format = resultmodel.OutputFormat(strings.TrimPrefix(argument, "--format="))
		default:
			context.RepositoryRoot = absolutePath(context.RepositoryRoot)
			finding := runtime.usageFinding(intendedCommand, "UNKNOWN-GLOBAL-OPTION", fmt.Sprintf("global option %q is not available", argument))
			return context, "", nil, &finding
		}
		if context.Format != resultmodel.FormatText && context.Format != resultmodel.FormatJSON {
			context.Format = resultmodel.FormatText
			context.RepositoryRoot = absolutePath(context.RepositoryRoot)
			finding := runtime.usageFinding(intendedCommand, "INVALID-OUTPUT-FORMAT", "--format requires text or json")
			return context, "", nil, &finding
		}
	}
	context.RepositoryRoot = absolutePath(context.RepositoryRoot)
	finding := runtime.usageFinding(intendedCommand, "MISSING-COMMAND",
		"a command is required: available commands are "+strings.Join(runtime.availableCommands(), ", "))
	return context, "", nil, &finding
}

// usageFinding names a runnable argv. A valid intended command remains the remedy for a bad
// global option; missing and unknown commands route to the built-in help surface.
func (runtime *CommandRuntime) usageFinding(commandName, code, evidence string) resultmodel.CommandFinding {
	nextCommand := runtime.nameableCommand(commandName)
	if code == "UNKNOWN-COMMAND" || code == "MISSING-COMMAND" {
		nextCommand = "help"
	}
	return resultmodel.CommandFinding{
		Code:                 code,
		Severity:             resultmodel.SeverityError,
		Evidence:             []string{evidence},
		Fixability:           resultmodel.FixabilityManual,
		AutomationStopReason: "the command line is not valid",
		NextArgv:             []string{"do-work-cli", "--format", "text", nextCommand},
		VerificationArgv:     []string{"do-work-cli", "--format", "json", nextCommand},
	}
}

// nameableCommand resolves the argv token a finding should suggest. A registered command the
// user actually named wins; otherwise the built-in help surface is always runnable.
func (runtime *CommandRuntime) nameableCommand(commandName string) string {
	if commandName == "help" {
		return commandName
	}
	if _, registered := runtime.handlers[commandName]; registered {
		return commandName
	}
	return "help"
}

func (runtime *CommandRuntime) availableCommands() []string {
	names := runtime.registeredCommands()
	for _, name := range names {
		if name == "help" {
			return names
		}
	}
	names = append(names, "help")
	sort.Strings(names)
	return names
}

// registeredCommands lists what this runtime can actually run, sorted so the listing does
// not depend on map iteration order.
func (runtime *CommandRuntime) registeredCommands() []string {
	names := make([]string, 0, len(runtime.handlers))
	for name := range runtime.handlers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// scanForCommandToken finds the command the argv was reaching for. Global options are known
// here, so their VALUES are skipped rather than mistaken for the command.
func scanForCommandToken(arguments []string) string {
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if !strings.HasPrefix(argument, "-") || argument == "-" {
			return argument
		}
		if argument == "--repo-root" || argument == "--format" {
			index++
		}
	}
	return ""
}

func absolutePath(path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return filepath.Clean(absolute)
}
