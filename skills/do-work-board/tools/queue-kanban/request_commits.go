package main

import (
	"flag"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// request-commits lists, for each named REQ, every commit in the full history
// that the board's activity correlation credits to it (requestIdsCreditedByCommit:
// a touched REQ path or a bracketed [REQ-NNN] subject token). It is the commit
// evidence `do-work trace` reads (../../../do-work/actions/trace-reference.md),
// so trace never keeps a second copy of the attribution rule. Read-only.
//
// Output is tab-separated with a header line, one row per (requested id,
// commit) in git log order (newest first). paths_outside_do_work counts touched
// paths not under do-work/: 0 means lifecycle bookkeeping (claim, archive, run
// artifacts) or a merge, because plain `--name-only` lists no paths for a merge.

// requestIdArgumentPattern is the only argument shape accepted after the flags.
var requestIdArgumentPattern = regexp.MustCompile(`^REQ-\d+$`)

// runRequestCommitsCommand runs the subcommand against the given writers and
// git runner and returns the exit code: 0 success (the header alone when
// nothing matches), 1 git or the repo root failed, 2 a usage error.
func runRequestCommitsCommand(args []string, standardOut io.Writer, standardErr io.Writer, runner gitCommandRunner) int {
	flagSet := flag.NewFlagSet("request-commits", flag.ContinueOnError)
	flagSet.SetOutput(standardErr)
	repoRootOverride := flagSet.String("repo-root", "", "repo root containing do-work/ (default: walk up from the working directory)")
	if flagSet.Parse(args) != nil {
		return 2
	}
	requestedIds := flagSet.Args()
	if len(requestedIds) == 0 {
		fmt.Fprintln(standardErr, "queue-kanban request-commits: name at least one REQ-NNN")
		return 2
	}
	for _, requestedId := range requestedIds {
		if !requestIdArgumentPattern.MatchString(requestedId) {
			fmt.Fprintf(standardErr, "queue-kanban request-commits: %q is not a request id (want REQ-NNN)\n", requestedId)
			return 2
		}
	}

	repoRoot, resolveError := resolveRepoRootOrDefault(*repoRootOverride)
	if resolveError != nil {
		fmt.Fprintln(standardErr, "queue-kanban request-commits:", resolveError)
		return 1
	}
	// No --since and no pathspec: shipped work can be old, and a do-work/
	// pathspec drops every merge.
	logOutput, logError := runner(repoRoot, "log", "--format=%H%x00%cI%x00%s", "--name-only")
	if logError != nil {
		fmt.Fprintln(standardErr, "queue-kanban request-commits: git log failed:", logError)
		return 1
	}

	fmt.Fprintln(standardOut, "request_id\tcommit\tcommitted_at\tpaths_outside_do_work\tsubject")
	for _, commit := range parseCorrelationLog(logOutput) {
		creditedIds := requestIdsCreditedByCommit(commit)
		pathsOutsideDoWork := 0
		for _, path := range commit.touchedPaths {
			if !strings.HasPrefix(path, "do-work/") {
				pathsOutsideDoWork++
			}
		}
		for _, requestedId := range requestedIds {
			if creditedIds[requestedId] {
				fmt.Fprintf(standardOut, "%s\t%s\t%s\t%d\t%s\n", requestedId, commit.hash,
					commit.committedAt.Format(time.RFC3339), pathsOutsideDoWork, commit.subject)
			}
		}
	}
	return 0
}
