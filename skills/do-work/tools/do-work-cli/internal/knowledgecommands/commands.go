package knowledgecommands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandBKBInit             = "bkb-init"
	CommandBKBStatus           = "bkb-status"
	CommandBKBLintStructure    = "bkb-lint-structure"
	CommandDreamScan           = "dream-scan"
	CommandInterviewList       = "interview-list"
	CommandInterviewStatus     = "interview-status"
	CommandInterviewExport     = "interview-export"
	CommandInterviewIngest     = "interview-ingest"
	CommandInterviewReset      = "interview-reset"
	CommandInterviewVersions   = "interview-versions"
	CommandMemoryRemember      = "memory-remember"
	CommandMemoryForget        = "memory-forget"
	CommandMemoryRecall        = "memory-recall"
	CommandMemoryStatus        = "memory-status"
	CommandMemoryBootstrap     = "memory-bootstrap"
	CommandMemoryAudit         = "memory-audit"
	CommandInstallMemoryHooks  = "install-memory-hooks"
	CommandLexicalMemoryRecall = "lexical-memory-recall"
)

var nowUTC = func() time.Time { return time.Now().UTC() }

type bkbOptions struct {
	target   string
	fillGaps bool
	dryRun   bool
	commit   bool
}

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandBKBInit:             handleBKBInit,
		CommandBKBStatus:           handleBKBStatus,
		CommandBKBLintStructure:    handleBKBLint,
		CommandDreamScan:           handleDreamScan,
		CommandInterviewList:       handleInterviewList,
		CommandInterviewStatus:     handleInterviewStatus,
		CommandInterviewExport:     handleInterviewExport,
		CommandInterviewIngest:     handleInterviewIngest,
		CommandInterviewReset:      handleInterviewReset,
		CommandInterviewVersions:   handleInterviewVersions,
		CommandMemoryRemember:      handleMemoryRemember,
		CommandMemoryForget:        handleMemoryForget,
		CommandMemoryRecall:        handleMemoryRecall,
		CommandMemoryStatus:        handleMemoryStatus,
		CommandMemoryBootstrap:     handleMemoryBootstrap,
		CommandMemoryAudit:         handleMemoryAudit,
		CommandInstallMemoryHooks:  handleInstallMemoryHooks,
		CommandLexicalMemoryRecall: handleLexicalMemoryRecall,
	}
}

func parseBKBOptions(arguments []string, mutable bool) (bkbOptions, error) {
	options := bkbOptions{target: "kb"}
	seenTarget := false
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		switch {
		case argument == "--kb" || argument == "--path":
			index++
			if index >= len(arguments) {
				return options, fmt.Errorf("%s requires a path", argument)
			}
			options.target, seenTarget = arguments[index], true
		case strings.HasPrefix(argument, "--kb="):
			options.target, seenTarget = strings.TrimPrefix(argument, "--kb="), true
		case strings.HasPrefix(argument, "--path="):
			options.target, seenTarget = strings.TrimPrefix(argument, "--path="), true
		case argument == "--fill-gaps" && mutable:
			options.fillGaps = true
		case argument == "--dry-run" && mutable:
			options.dryRun = true
		case argument == "--commit" && mutable:
			options.commit = true
		case !strings.HasPrefix(argument, "-") && !seenTarget:
			options.target, seenTarget = argument, true
		default:
			return options, fmt.Errorf("unknown option %q", argument)
		}
	}
	if options.dryRun && options.commit {
		return options, fmt.Errorf("--dry-run and --commit cannot be combined")
	}
	if strings.TrimSpace(options.target) == "" {
		return options, fmt.Errorf("target path must not be empty")
	}
	options.target = filepath.Clean(options.target)
	return options, nil
}

func usageResult(commandName string, err error) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
		knowledgeFinding(commandName, "KNOWLEDGE-USAGE", resultmodel.SeverityError, nil, err.Error(), resultmodel.FixabilityManual, "the command line is invalid"),
	}}
}

func knowledgeFinding(commandName, code string, severity resultmodel.FindingSeverity, paths []string, evidence string, fixability resultmodel.FindingFixability, stopReason string) resultmodel.CommandFinding {
	commandArgv := []string{"do-work-cli", commandName}
	verification := []string{"do-work-cli", "--format", "json", commandName}
	recipe := commandName
	if len(paths) > 0 && paths[0] != "." {
		target := paths[0]
		commandArgv = append(commandArgv, "--path", target)
		verification = append(verification, "--path", target)
		recipe += " " + quoteRecipeArgument(target)
	}
	return resultmodel.CommandFinding{Code: code, Severity: severity, AffectedPaths: append([]string(nil), paths...), Evidence: []string{evidence}, Fixability: fixability, AutomationStopReason: stopReason, NextArgv: commandArgv, NextJustRecipe: recipe, VerificationArgv: verification}
}

func quoteRecipeArgument(value string) string {
	if value != "" && !strings.ContainsAny(value, " \t\n'\"\\$`;&|<>()") {
		return value
	}
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func ensureSafeTarget(root, relative string) (string, error) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	if physicalRoot, resolveError := filepath.EvalSymlinks(absoluteRoot); resolveError == nil {
		absoluteRoot = physicalRoot
	}
	absoluteTarget := relative
	if !filepath.IsAbs(absoluteTarget) {
		absoluteTarget = filepath.Join(absoluteRoot, relative)
	}
	return physicalPath(filepath.Clean(absoluteTarget))
}

func physicalPath(path string) (string, error) {
	missingParts := []string{}
	current := path
	for {
		_, err := os.Lstat(current)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", err
		}
		missingParts = append(missingParts, filepath.Base(current))
		current = parent
	}
	physical, err := filepath.EvalSymlinks(current)
	if err != nil {
		return "", err
	}
	for index := len(missingParts) - 1; index >= 0; index-- {
		physical = filepath.Join(physical, missingParts[index])
	}
	return filepath.Clean(physical), nil
}

func sortFindings(findings []resultmodel.CommandFinding) {
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].Code != findings[j].Code {
			return findings[i].Code < findings[j].Code
		}
		return strings.Join(findings[i].AffectedPaths, "\x00") < strings.Join(findings[j].AffectedPaths, "\x00")
	})
}
