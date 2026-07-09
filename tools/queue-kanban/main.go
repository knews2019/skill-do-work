package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// queue-kanban parses the do-work/ Markdown tree into a board model and renders
// it. The model + parser (LoadBoard) is the foundation; on top of it sit three
// subcommands — `summary` (column counts), `generate` (a self-contained static
// board), and `serve` (a live local board that re-walks the tree per request).
//
// Dispatch is a minimal hand-rolled subcommand switch over os.Args[1] — no
// external CLI library — with each subcommand owning its own flag.FlagSet:
//
//	queue-kanban summary  [--repo-root DIR] [--recent-window DUR]
//	queue-kanban generate --out DIR [--repo-root DIR]
//	queue-kanban serve    [--port PORT] [--repo-root DIR] [--open]
//
// Invoking the binary with no subcommand prints the model summary.
//
// Only summary exposes --recent-window: the HTML board picks its visible
// Recently-done window client-side (the 24h/48h/7d toggle, default 24h), so a
// server-side window flag on generate/serve would be advertised but inert.
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
	default:
		fmt.Fprintf(os.Stderr, "queue-kanban: unknown subcommand %q (want summary | generate | serve)\n", subcommand)
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

	fmt.Printf("queue-kanban board summary\n")
	fmt.Printf("  total REQ tickets   : %d\n", len(board.AllRequests))
	fmt.Printf("  user requests       : %d\n", len(board.UserRequests))
	fmt.Printf("  pending             : %d\n", len(board.Columns.Pending))
	fmt.Printf("    ready to work     : %d\n", len(board.Columns.PendingReady))
	fmt.Printf("    waiting on deps   : %d\n", len(board.Columns.PendingWaiting))
	fmt.Printf("  claimed             : %d\n", len(board.Columns.Claimed))
	fmt.Printf("  needs-input/blocked : %d\n", len(board.Columns.NeedsInputOrBlocked))
	fmt.Printf("  recently-done       : %d\n", len(board.Columns.RecentlyDone))
	fmt.Printf("  calendar entries    : %d\n", len(board.Calendar))
	fmt.Printf("  dependency edges    : %d\n", len(board.DependencyGraph.Edges))
	if len(board.Warnings) > 0 {
		fmt.Printf("  warnings            : %d\n", len(board.Warnings))
		for _, warningText := range board.Warnings {
			fmt.Printf("    ! %s\n", warningText)
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
	fmt.Printf("queue-kanban: wrote static board to %s/ (index.html + board-data.js, %d REQs, %d URs, %d calendar entries)\n",
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
