package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// audit-metrics computes the maintainability audit's deterministic numbers —
// tracked-file inventory, size distributions, band flags, git churn, and the
// churn × size hotspot join — as pasteable markdown tables, so the audit
// action pastes tool output instead of prescribing fragile find/wc/awk
// pipelines to an LLM. Judgment (calibrating thresholds, reading the numbers)
// stays in the action's prose; this tool only measures. It is strictly
// read-only: it prints, and writes nothing.
//
// Dispatch is a minimal hand-rolled subcommand switch over os.Args[1] — no
// external CLI library — with each subcommand owning its own flag.FlagSet:
//
//	audit-metrics inventory [--repo-root DIR] [--exclude-path PREFIX]... [--top-count N]
//	                        [--watch-lines N] [--flag-lines N] [--watch-words N] [--flag-words N]
//	audit-metrics folders   [--repo-root DIR] [--exclude-path PREFIX]... [--top-count N]
//	                        [--watch-files N] [--flag-files N]
//	audit-metrics churn     [--repo-root DIR] [--exclude-path PREFIX]... [--top-count N]
//	                        [--since-window WINDOW]
//	audit-metrics hotspots  [--repo-root DIR] [--exclude-path PREFIX]... [--top-count N]
//	                        [--since-window WINDOW]
//
// Band thresholds come ONLY from flags: no flag, no band output. A value
// strictly greater than a threshold is flagged; equal is not. --exclude-path
// is a repeatable repo-relative path PREFIX and defaults to empty — the caller
// owns the exclude list (any built-in "ceremony files" list would go stale).
// Every subcommand rejects leftover tokens rather than ignoring them — the
// queue-kanban lesson where a silently discarded argument shipped a wrong-repo
// bump.
func main() {
	subcommand := ""
	subcommandArgs := os.Args[1:]
	if len(subcommandArgs) > 0 && !isFlagArgument(subcommandArgs[0]) {
		subcommand = subcommandArgs[0]
		subcommandArgs = subcommandArgs[1:]
	}

	switch subcommand {
	case "inventory":
		runInventoryCommand(subcommandArgs)
	case "folders":
		runFoldersCommand(subcommandArgs)
	case "churn":
		runChurnCommand(subcommandArgs)
	case "hotspots":
		runHotspotsCommand(subcommandArgs)
	default:
		fmt.Fprintf(os.Stderr, "audit-metrics: unknown subcommand %q (want inventory | folders | churn | hotspots)\n", subcommand)
		os.Exit(2)
	}
}

// isFlagArgument reports whether an argument is a flag (leading "-") rather
// than a subcommand name.
func isFlagArgument(argument string) bool {
	return len(argument) > 0 && argument[0] == '-'
}

// rejectLeftoverArguments turns unconsumed tokens into an error — silently
// discarding an argument is how a tool ships answering about the wrong repo.
func rejectLeftoverArguments(subcommandName string, leftoverArguments []string) error {
	if len(leftoverArguments) == 0 {
		return nil
	}
	return fmt.Errorf("audit-metrics %s: unexpected argument(s): %s", subcommandName, strings.Join(leftoverArguments, " "))
}

// exitOnLeftoverArguments is rejectLeftoverArguments for the command wrappers,
// which own os.Exit.
func exitOnLeftoverArguments(subcommandName string, leftoverArguments []string) {
	if leftoverError := rejectLeftoverArguments(subcommandName, leftoverArguments); leftoverError != nil {
		fmt.Fprintln(os.Stderr, leftoverError)
		os.Exit(2)
	}
}

// exitOnCommandError reports a compute failure (usually git declining to
// answer) and exits non-zero.
func exitOnCommandError(subcommandName string, commandError error) {
	fmt.Fprintf(os.Stderr, "audit-metrics %s: %v\n", subcommandName, commandError)
	os.Exit(1)
}

// repeatablePathList collects every occurrence of a repeatable path flag.
type repeatablePathList []string

func (pathList *repeatablePathList) String() string { return strings.Join(*pathList, ",") }

func (pathList *repeatablePathList) Set(flagValue string) error {
	*pathList = append(*pathList, flagValue)
	return nil
}

// defaultSinceWindow is the churn/hotspots history horizon when the caller
// passes no --since-window; the rendered output always names the window in
// use, so a defaulted run is never ambiguous.
const defaultSinceWindow = "12 months"

// defaultTopCount is the top-N size for every offender table.
const defaultTopCount = 10

func runInventoryCommand(args []string) {
	flagSet := flag.NewFlagSet("inventory", flag.ExitOnError)
	repoRoot := flagSet.String("repo-root", ".", "git repo root to measure")
	var excludePaths repeatablePathList
	flagSet.Var(&excludePaths, "exclude-path", "repo-relative path prefix to exclude (repeatable; default none)")
	topCount := flagSet.Int("top-count", defaultTopCount, "how many top offenders to list")
	watchLines := flagSet.Int("watch-lines", bandThresholdUnset, "WATCH threshold for file lines (omit for no band)")
	flagLines := flagSet.Int("flag-lines", bandThresholdUnset, "FLAG threshold for file lines (omit for no band)")
	watchWords := flagSet.Int("watch-words", bandThresholdUnset, "WATCH threshold for file words (omit for no band)")
	flagWords := flagSet.Int("flag-words", bandThresholdUnset, "FLAG threshold for file words (omit for no band)")
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("inventory", flagSet.Args())

	resolvedRoot, rootError := resolveRepositoryRoot(*repoRoot)
	if rootError != nil {
		exitOnCommandError("inventory", rootError)
	}
	report, computeError := computeInventoryReport(resolvedRoot, excludePaths)
	if computeError != nil {
		exitOnCommandError("inventory", computeError)
	}
	thresholds := fileBandThresholds{WatchLines: *watchLines, FlagLines: *flagLines, WatchWords: *watchWords, FlagWords: *flagWords}
	writeInventoryReport(os.Stdout, report, thresholds, *topCount)
}

func runFoldersCommand(args []string) {
	flagSet := flag.NewFlagSet("folders", flag.ExitOnError)
	repoRoot := flagSet.String("repo-root", ".", "git repo root to measure")
	var excludePaths repeatablePathList
	flagSet.Var(&excludePaths, "exclude-path", "repo-relative path prefix to exclude (repeatable; default none)")
	topCount := flagSet.Int("top-count", defaultTopCount, "how many top offenders to list")
	watchFiles := flagSet.Int("watch-files", bandThresholdUnset, "WATCH threshold for files per folder (omit for no band)")
	flagFiles := flagSet.Int("flag-files", bandThresholdUnset, "FLAG threshold for files per folder (omit for no band)")
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("folders", flagSet.Args())

	resolvedRoot, rootError := resolveRepositoryRoot(*repoRoot)
	if rootError != nil {
		exitOnCommandError("folders", rootError)
	}
	report, computeError := computeInventoryReport(resolvedRoot, excludePaths)
	if computeError != nil {
		exitOnCommandError("folders", computeError)
	}
	writeFoldersReport(os.Stdout, report, *watchFiles, *flagFiles, *topCount)
}

func runChurnCommand(args []string) {
	flagSet := flag.NewFlagSet("churn", flag.ExitOnError)
	repoRoot := flagSet.String("repo-root", ".", "git repo root to measure")
	var excludePaths repeatablePathList
	flagSet.Var(&excludePaths, "exclude-path", "repo-relative path prefix to exclude (repeatable; default none)")
	topCount := flagSet.Int("top-count", defaultTopCount, "how many top offenders to list")
	sinceWindow := flagSet.String("since-window", defaultSinceWindow, "history window passed to git log --since")
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("churn", flagSet.Args())

	resolvedRoot, rootError := resolveRepositoryRoot(*repoRoot)
	if rootError != nil {
		exitOnCommandError("churn", rootError)
	}
	report, computeError := computeChurnReport(resolvedRoot, *sinceWindow, excludePaths)
	if computeError != nil {
		exitOnCommandError("churn", computeError)
	}
	writeChurnReport(os.Stdout, report, *topCount)
}

func runHotspotsCommand(args []string) {
	flagSet := flag.NewFlagSet("hotspots", flag.ExitOnError)
	repoRoot := flagSet.String("repo-root", ".", "git repo root to measure")
	var excludePaths repeatablePathList
	flagSet.Var(&excludePaths, "exclude-path", "repo-relative path prefix to exclude (repeatable; default none)")
	topCount := flagSet.Int("top-count", defaultTopCount, "how many top offenders to list")
	sinceWindow := flagSet.String("since-window", defaultSinceWindow, "history window passed to git log --since")
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("hotspots", flagSet.Args())

	resolvedRoot, rootError := resolveRepositoryRoot(*repoRoot)
	if rootError != nil {
		exitOnCommandError("hotspots", rootError)
	}
	report, computeError := computeChurnReport(resolvedRoot, *sinceWindow, excludePaths)
	if computeError != nil {
		exitOnCommandError("hotspots", computeError)
	}
	hotspots, joinError := computeHotspotEntries(resolvedRoot, report)
	if joinError != nil {
		exitOnCommandError("hotspots", joinError)
	}
	writeHotspotsReport(os.Stdout, report, hotspots, *topCount)
}
