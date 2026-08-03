package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

// queue-kanban parses the do-work/ Markdown tree into a board model and renders
// it. The model + parser (LoadBoard) is the foundation; on top of it sit the
// board subcommands — `summary` (column counts), `generate` (a self-contained
// static board), and `serve` (a live local board that re-walks the tree per
// request) — plus three release-ritual subcommands: `next-req` and
// `next-version` allocate numbers, and `verify` checks the cross-file invariants
// that are otherwise verified by hand on every commit.
//
// Dispatch is a minimal hand-rolled subcommand switch over os.Args[1] — no
// external CLI library — with each subcommand owning its own flag.FlagSet:
//
//	queue-kanban summary      [--repo-root DIR] [--recent-window DUR]
//	queue-kanban generate     --out DIR [--repo-root DIR]
//	queue-kanban serve        [--port PORT] [--repo-root DIR] [--open]
//	queue-kanban next-req     [--repo-root DIR]
//	queue-kanban next-version <patch|minor|major> [--repo-root DIR] [--version-file PATH]
//	queue-kanban verify       [--repo-root DIR]
//
// Invoking the binary with no subcommand prints the model summary.
//
// Only summary exposes --recent-window: the HTML board picks its visible
// Recently-done window client-side (the 24h/48h/7d toggle, default 24h), so a
// server-side window flag on generate/serve would be advertised but inert.
//
// Write surfaces, in full: the board's testing view (serve; see testing.go)
// writes the testing-track frontmatter fields plus do-work/testers.md, and
// `next-version` writes one line in one version file. Nothing else here writes
// anything — `next-req` and `verify` are read-only, and no subcommand ever writes
// CHANGELOG.md, which stays an owner-only, human-authored file.
func main() {
	subcommand := ""
	subcommandArgs := os.Args[1:]
	if len(subcommandArgs) > 0 && !isFlagArgument(subcommandArgs[0]) {
		subcommand = subcommandArgs[0]
		subcommandArgs = subcommandArgs[1:]
	}

	switch subcommand {
	case "", "summary":
		runSummaryCommand(subcommandArgs)
	case "generate":
		runGenerateCommand(subcommandArgs)
	case "serve":
		runServeCommand(subcommandArgs)
	case "next-req":
		runNextRequestCommand(subcommandArgs)
	case "next-version":
		runNextVersionCommand(subcommandArgs)
	case "verify":
		runVerifyCommand(subcommandArgs)
	default:
		fmt.Fprintf(os.Stderr, "queue-kanban: unknown subcommand %q (want summary | generate | serve | next-req | next-version | verify)\n", subcommand)
		os.Exit(2)
	}
}

// isFlagArgument reports whether an argument is a flag (leading "-") rather than
// a subcommand name, so `queue-kanban --recent-window …` still routes to summary.
func isFlagArgument(argument string) bool {
	return len(argument) > 0 && argument[0] == '-'
}

// defaultRecentWindow is the Recently-done horizon used to bucket the board
// model's RecentlyDone column. Summary lets the user override it via
// --recent-window; generate and serve always use this default because their
// visible window is chosen client-side in board.js.
const defaultRecentWindow = 7 * 24 * time.Hour

// runSummaryCommand prints the parsed board model's counts — the REQ-1207 smoke.
func runSummaryCommand(args []string) {
	flagSet := flag.NewFlagSet("summary", flag.ExitOnError)
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	recentWindow := flagSet.Duration("recent-window", defaultRecentWindow, "window for the Recently-done column")
	_ = flagSet.Parse(args)

	board := loadBoardOrExit(*repoRootOverride, *recentWindow)
	writeBoardSummary(os.Stdout, board)
}

// writeBoardSummary renders the summary block. Split from runSummaryCommand so
// tests can assert the headless output — the summary is the one mode with no
// browser, so it must expose completion anomalies on its own.
func writeBoardSummary(outputWriter io.Writer, board *Board) {
	fmt.Fprintf(outputWriter, "queue-kanban board summary\n")
	fmt.Fprintf(outputWriter, "  total REQ tickets   : %d\n", len(board.AllRequests))
	fmt.Fprintf(outputWriter, "  user requests       : %d\n", len(board.UserRequests))
	fmt.Fprintf(outputWriter, "  pending             : %d\n", len(board.Columns.Pending))
	fmt.Fprintf(outputWriter, "    ready to work     : %d\n", len(board.Columns.PendingReady))
	fmt.Fprintf(outputWriter, "    waiting on deps   : %d\n", len(board.Columns.PendingWaiting))
	fmt.Fprintf(outputWriter, "  claimed             : %d\n", len(board.Columns.Claimed))
	fmt.Fprintf(outputWriter, "  needs-input/blocked : %d\n", len(board.Columns.NeedsInputOrBlocked))
	fmt.Fprintf(outputWriter, "  recently-done       : %d\n", len(board.Columns.RecentlyDone))
	fmt.Fprintf(outputWriter, "  completion anomalies : %d\n", len(board.Columns.CompletionAnomalies))
	for _, anomalousTicket := range board.Columns.CompletionAnomalies {
		fmt.Fprintf(outputWriter, "    ! %s — %s\n", anomalousTicket.RequestId, anomalousTicket.CompletionAnomalyReason)
	}
	fmt.Fprintf(outputWriter, "  calendar entries    : %d\n", len(board.Calendar))
	fmt.Fprintf(outputWriter, "  dependency edges    : %d\n", len(board.DependencyGraph.Edges))
	if len(board.Warnings) > 0 {
		fmt.Fprintf(outputWriter, "  warnings            : %d\n", len(board.Warnings))
		for _, warningText := range board.Warnings {
			fmt.Fprintf(outputWriter, "    ! %s\n", warningText)
		}
	}
}

// runGenerateCommand writes the self-contained static board into --out.
func runGenerateCommand(args []string) {
	flagSet := flag.NewFlagSet("generate", flag.ExitOnError)
	outputDirectory := flagSet.String("out", "", "output directory for the self-contained static board (required)")
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	_ = flagSet.Parse(args)

	if *outputDirectory == "" {
		fmt.Fprintln(os.Stderr, "queue-kanban generate: --out DIR is required")
		os.Exit(2)
	}

	board := loadBoardOrExit(*repoRootOverride, defaultRecentWindow)
	if generateError := generateStaticSite(*outputDirectory, board); generateError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", generateError)
		os.Exit(1)
	}
	fmt.Printf("queue-kanban: wrote static board to %s/ (index.html + board-data.js + lazy board-markdown.js, %d REQs, %d URs, %d calendar entries)\n",
		*outputDirectory, len(board.AllRequests), len(board.UserRequests), len(board.Calendar))
}

// loadBoardOrExit builds the board against the live tree or exits non-zero with a
// diagnostic — the shared front half of every subcommand.
func loadBoardOrExit(repoRootOverride string, recentWindow time.Duration) *Board {
	board, loadError := LoadBoard(repoRootOverride, time.Now(), recentWindow)
	if loadError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", loadError)
		os.Exit(1)
	}
	return board
}

// runNextRequestCommand prints the next free REQ number and nothing else, so a
// caller can use it directly: `REQ-$(queue-kanban next-req)`. Read-only toward
// the queue — it allocates a number, it does not reserve one (see allocate.go on
// why that is accepted).
func runNextRequestCommand(args []string) {
	flagSet := flag.NewFlagSet("next-req", flag.ExitOnError)
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	_ = flagSet.Parse(args)

	allocatedNumber, allocateError := nextRequestNumber(*repoRootOverride)
	if allocateError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", allocateError)
		os.Exit(1)
	}
	fmt.Printf("%d\n", allocatedNumber)
}

// runNextVersionCommand bumps the version file by an explicitly named size and
// prints the new version. The size is a positional argument on purpose: patch vs
// minor vs major is a human judgment about what the change did to consumers, and
// a default would quietly make it for them.
//
// It never writes CHANGELOG.md. Composing the entry — and deciding the bump size
// this command is told — stays with the human running the release ritual.
func runNextVersionCommand(args []string) {
	flagSet := flag.NewFlagSet("next-version", flag.ExitOnError)
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	versionFileOverride := flagSet.String("version-file", "", "file carrying the `**Current version**: X.Y.Z` line (default: <repo-root>/actions/version.md)")
	_ = flagSet.Parse(args)

	bumpSize := ""
	if flagSet.NArg() > 0 {
		bumpSize = flagSet.Arg(0)
	}
	if bumpSize == "" {
		fmt.Fprintln(os.Stderr, "queue-kanban next-version: name the bump size — patch | minor | major")
		os.Exit(2)
	}

	repoRoot, resolveError := resolveRepoRootOrDefault(*repoRootOverride)
	if resolveError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", resolveError)
		os.Exit(1)
	}

	versionFilePath := resolveVersionFilePath(repoRoot, *versionFileOverride)
	allocatedVersion, allocateError := allocateNextVersion(versionFilePath, bumpSize)
	if allocateError != nil {
		// Exit 1 without writing. The caller falls back to its own version-source
		// resolution (actions/work-reference.md → Changelog Entry Procedure) —
		// this command only serves a repo that keeps a `**Current version**:` line.
		fmt.Fprintln(os.Stderr, "queue-kanban:", allocateError)
		os.Exit(1)
	}
	fmt.Printf("%s\n", allocatedVersion)
}

// runVerifyCommand prints the verify report and exits non-zero when it found
// anything. Read-only: it reports and routes, and repairs nothing — fixes belong
// to `do-work cleanup`, which asks before it acts.
func runVerifyCommand(args []string) {
	flagSet := flag.NewFlagSet("verify", flag.ExitOnError)
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	_ = flagSet.Parse(args)

	report, verifyError := runVerifyProbes(*repoRootOverride, time.Now())
	if verifyError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", verifyError)
		os.Exit(1)
	}
	fmt.Print(renderVerifyReport(report))
	os.Exit(report.ExitCode())
}
