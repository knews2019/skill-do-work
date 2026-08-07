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
// board subcommands — `summary` (column counts), `open-work` (a per-ticket
// digest of what is in flight), `generate` (a self-contained static board), and
// `serve` (a live local board that re-walks the tree per
// request) — plus the read-only `frontmatter` field reader and three
// release-ritual subcommands: `next-req` atomically reserves a number,
// `next-version` allocates a version, and `verify` checks the cross-file
// invariants otherwise checked by hand.
//
// Dispatch is a minimal hand-rolled subcommand switch over os.Args[1] — no
// external CLI library — with each subcommand owning its own flag.FlagSet:
//
//	queue-kanban summary      [--repo-root DIR] [--recent-window DUR]
//	queue-kanban open-work    [--repo-root DIR]
//	queue-kanban generate     --out DIR [--repo-root DIR]
//	queue-kanban serve        [--port PORT] [--repo-root DIR] [--open]
//	queue-kanban next-req     [--repo-root DIR]
//	queue-kanban next-version <patch|minor|major> [--repo-root DIR] [--version-file PATH]
//	queue-kanban verify       [--repo-root DIR]
//	queue-kanban frontmatter get FILE FIELD [--normalize] [--in-set SET]
//	queue-kanban now
//
// Invoking the binary with no subcommand prints the model summary.
//
// next-version is the one subcommand mixing a positional with flags, and the
// flags may appear on either side of it — the synopsis above is conventional
// notation, not a required order (see parseNextVersionArguments for why that
// needed saying). Every subcommand rejects leftover tokens rather than ignoring
// them; silently discarding an argument is how next-version shipped bumping the
// wrong repo.
//
// Only summary exposes --recent-window: the HTML board picks its visible
// Recently-done window client-side (the 24h/48h/7d toggle, default 24h), so a
// server-side window flag on generate/serve would be advertised but inert. The
// same reasoning keeps it off open-work, which shows only open work and so has no
// recently-done section for a window to size.
//
// Write surfaces, in full: the board's testing view (serve; see testing.go)
// writes the testing-track frontmatter fields plus do-work/testers.md;
// `next-version` writes one line in one version file; and `next-req` creates one
// durable marker under do-work/.req-reservations/. Everything else here is
// read-only, and no subcommand ever writes CHANGELOG.md, which stays an
// owner-only, human-authored file.
//
// `now` takes no --repo-root: it reads a clock, not a tree, so it is the one
// subcommand that works outside a project entirely.
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
	case "open-work":
		runOpenWorkCommand(subcommandArgs)
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
	case "now":
		runNowCommand(subcommandArgs)
	case "frontmatter":
		os.Exit(runFrontmatterCommand(subcommandArgs, os.Stdout, os.Stderr))
	default:
		fmt.Fprintf(os.Stderr, "queue-kanban: unknown subcommand %q (want summary | generate | serve | next-req | next-version | verify | now | frontmatter | open-work)\n", subcommand)
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
	exitOnLeftoverArguments("summary", flagSet.Args())

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

// runOpenWorkCommand prints the open-work digest — the fast terminal answer to
// "what is in flight?" (open count, claimed titles, needs-input statuses). The
// renderer lives in open_work.go; see its header for why this is a separate
// subcommand rather than more lines in summary.
func runOpenWorkCommand(args []string) {
	flagSet := flag.NewFlagSet("open-work", flag.ExitOnError)
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("open-work", flagSet.Args())

	// defaultRecentWindow is passed because LoadBoard requires a window, not
	// because the digest has anything windowed to show.
	board := loadBoardOrExit(*repoRootOverride, defaultRecentWindow)
	writeOpenWorkDigest(os.Stdout, board)
}

// runGenerateCommand writes the self-contained static board into --out.
func runGenerateCommand(args []string) {
	flagSet := flag.NewFlagSet("generate", flag.ExitOnError)
	outputDirectory := flagSet.String("out", "", "output directory for the self-contained static board (required)")
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("generate", flagSet.Args())

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

// runNextRequestCommand atomically reserves the next free REQ number, then
// prints the number and nothing else so a caller can use it directly:
// `REQ-$(queue-kanban next-req)`.
func runNextRequestCommand(args []string) {
	flagSet := flag.NewFlagSet("next-req", flag.ExitOnError)
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("next-req", flagSet.Args())

	allocatedNumber, allocateError := nextRequestNumber(*repoRootOverride)
	if allocateError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", allocateError)
		os.Exit(1)
	}
	fmt.Printf("%d\n", allocatedNumber)
}

// nextVersionArguments is the fully-resolved argument set for `next-version`.
// It exists so the parse can be asserted directly: runNextVersionCommand calls
// os.Exit, so anything left inline in it is untestable — which is exactly how
// the flags-after-positional defect shipped and stayed invisible.
type nextVersionArguments struct {
	BumpSize            string
	RepoRootOverride    string
	VersionFileOverride string
}

// parseNextVersionArguments accepts the bump size and the flags in ANY order.
//
// flag.FlagSet.Parse stops at the first non-flag argument, so a single Parse of
// `["patch", "--repo-root", "X"]` consumes nothing and silently discards both
// overrides — the command then writes whatever tree it was launched from and
// exits 0. That is the shape the skill's own prescribed invocation had.
//
// The fix parses twice rather than lifting the positional out of the slice by
// index. Two reasons: the positional is not always at index 0 (a flag may
// precede it), so an index-based lift needs its own mini-parser to find it; and
// double-parsing keeps flag.FlagSet the single authority on what a flag looks
// like, including `--flag=value`, `-flag value`, and `--`. First Parse consumes
// any leading flags and halts on the bump size; Args()[1:] is then everything
// after it, which the second Parse consumes.
//
// Leftovers are an error, never silence. A stray second positional
// (`next-version patch minor`) and an unknown flag both fail with a message
// naming the offending token, because ignoring them is what let this defect
// survive a release.
func parseNextVersionArguments(args []string) (nextVersionArguments, error) {
	parsed := nextVersionArguments{}

	flagSet := flag.NewFlagSet("next-version", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard) // the caller renders errors and usage; tests must stay quiet
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	versionFileOverride := flagSet.String("version-file", "", "file carrying the `**Current version**: X.Y.Z` line (default: <repo-root>/actions/version.md)")

	if parseError := flagSet.Parse(args); parseError != nil {
		return parsed, parseError
	}
	if flagSet.NArg() == 0 {
		return parsed, fmt.Errorf("name the bump size — patch | minor | major")
	}

	parsed.BumpSize = flagSet.Arg(0)
	if parseError := flagSet.Parse(flagSet.Args()[1:]); parseError != nil {
		return parsed, parseError
	}
	if leftoverError := rejectLeftoverArguments("next-version", flagSet.Args()); leftoverError != nil {
		return parsed, leftoverError
	}

	parsed.RepoRootOverride = *repoRootOverride
	parsed.VersionFileOverride = *versionFileOverride
	return parsed, nil
}

// rejectLeftoverArguments turns unconsumed tokens into an error.
//
// The condition, not a list of today's subcommands: ANY subcommand that finishes
// parsing with tokens left over must fail rather than ignore them. Silence is
// how `next-version` shipped writing the wrong repo — its discarded
// `--repo-root` sat unread in Arg(1) and the command exited 0. The same shape is
// reachable on a flags-only subcommand too: a stray token placed first halts
// Parse, so every flag after it is discarded exactly the same way. Any
// subcommand added later inherits the rule by calling this.
func rejectLeftoverArguments(subcommandName string, leftoverArguments []string) error {
	if len(leftoverArguments) == 0 {
		return nil
	}
	return fmt.Errorf("unrecognized argument(s) for %s: %v", subcommandName, leftoverArguments)
}

// exitOnLeftoverArguments is rejectLeftoverArguments for the command wrappers,
// which report and exit rather than returning an error.
func exitOnLeftoverArguments(subcommandName string, leftoverArguments []string) {
	if leftoverError := rejectLeftoverArguments(subcommandName, leftoverArguments); leftoverError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", leftoverError)
		os.Exit(2)
	}
}

// runNextVersionCommand bumps the version file by an explicitly named size and
// prints the new version. The size is a positional argument on purpose: patch vs
// minor vs major is a human judgment about what the change did to consumers, and
// a default would quietly make it for them.
//
// It never writes CHANGELOG.md. Composing the entry — and deciding the bump size
// this command is told — stays with the human running the release ritual.
func runNextVersionCommand(args []string) {
	parsedArguments, parseError := parseNextVersionArguments(args)
	if parseError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban next-version:", parseError)
		fmt.Fprintln(os.Stderr, "usage: queue-kanban next-version <patch|minor|major> [--repo-root DIR] [--version-file PATH]")
		os.Exit(2)
	}

	repoRoot, resolveError := resolveRepoRootOrDefault(parsedArguments.RepoRootOverride)
	if resolveError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", resolveError)
		os.Exit(1)
	}

	bumpSize := parsedArguments.BumpSize
	versionFilePath := resolveVersionFilePath(repoRoot, parsedArguments.VersionFileOverride)
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

// runNowCommand prints the current UTC instant in the Timestamp rule's shape, so
// a caller can use it directly: `STAMP=$(queue-kanban now)`. Read-only — it
// reads a clock and touches no file. It parses a FlagSet purely so an
// unrecognized flag is rejected rather than silently ignored; there are no flags
// to accept, because there is no tree to point at.
func runNowCommand(args []string) {
	flagSet := flag.NewFlagSet("now", flag.ExitOnError)
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("now", flagSet.Args())

	writeCanonicalTimestamp(os.Stdout, time.Now())
}

// runVerifyCommand prints the verify report and exits non-zero when it found
// anything. Read-only: it reports and routes, and repairs nothing — fixes belong
// to `do-work cleanup`, which asks before it acts.
func runVerifyCommand(args []string) {
	flagSet := flag.NewFlagSet("verify", flag.ExitOnError)
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	_ = flagSet.Parse(args)
	exitOnLeftoverArguments("verify", flagSet.Args())

	report, verifyError := runVerifyProbes(*repoRootOverride, time.Now())
	if verifyError != nil {
		fmt.Fprintln(os.Stderr, "queue-kanban:", verifyError)
		os.Exit(1)
	}
	fmt.Print(renderVerifyReport(report))
	os.Exit(report.ExitCode())
}
