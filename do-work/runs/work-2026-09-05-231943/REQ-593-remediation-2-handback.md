# REQ-593 Second Remediation Hand-back — the empty-lane false green, and three guards that could not tell fixed from unfixed

**Branch:** `worktree-agent-REQ-593-heavy-tier-verdict`
**Head:** `4298719b6ab04d72a8044eed65c5518a02487f6d`
(one commit on top of `b496589`, which already carried both merged parts and the review)
**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-593-heavy-tier-verdict`
Working tree is clean. Nothing outside the worktree was written except this file, which is neither
staged nor committed. Nothing was pushed. The repository-wide SIGPIPE generalization was not
attempted, as instructed.

## What changed

Five files, `+193 / -15`.

- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go` — the empty-lane guard
  (F1); the red finding keeps the announced skip line (F6).
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go` — three assertions:
  the empty-lane refusal, the nested-output leak, and the announcement on the red finding.
- `_dev/tests/update-script-behavior.sh` — the source scanner works on the defect's condition instead
  of one spelling, plus a fixture that pins the widening (F3); both marker archives asserted readable
  (F4).
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — the `silent-skip-reads-as-red` trap states
  the new two-part rule (F5). **Scope widening, see below.**
- `do-work/working/REQ-593-make-the-heavy-tier-verdict-honest.md` — the prime added to `write_set`.

## SCOPE WIDENING — declare in Scope

`skills/do-work/tools/do-work-cli/prime-do-work-cli.md` was in the request's `prime_files` and not in
its `write_set`. Its trap `[family: silent-skip-reads-as-red]` still ended "the runner keys `skipped`
on that prefix, never on a lane-name list", which is the contract this request replaced. It now reads
"keys `skipped` on that prefix **and** a zero exit status", names the red finding's second evidence
line, and adds the nested-runner half. The file was added to the request's `write_set` frontmatter.
That is the only widening; no other file outside the declared set was touched.

## F1 — `RunLanes` with no lanes is refused, in `RunLanes`

`heavy_run.go:91`. The guard is the first statement in the function, before the manifest resolve, the
HEAD resolve and the dirty-tree refusal, so an invalid request cannot come back wearing another
refusal's code (a dirty tree would otherwise mask it). New typed code `HEAVY-RUN-NO-LANE-REQUESTED`,
carried by the existing `refuseLaneRun` / `LaneRunRefusalCode` path, so the CLI reports it the same way
it reports `HEAVY-RUN-LANE-UNKNOWN`. `heavy_commands.go` was not touched: its parser already refuses an
empty `--lane` set, and that refusal is now a second line of defence rather than the only one.

The reviewer asked for it "beside the `LaneTimeout` guard". It is in the same function and not on the
same lines, because the `LaneTimeout` guard sits after four git operations that an empty request should
never reach. Stated here rather than silently deviating.

**Assertion:** `TestRunLanesRefusesARunThatNamesNoLane` calls `RunLanes` exactly as the reviewer's probe
did — `LaneRunRequest{RepositoryRoot, ManifestPath, LaneOutputWriter: io.Discard}` with no `LaneIDs`.

**RED by ablation** (guard deleted):

```
--- FAIL: TestRunLanesRefusesARunThatNamesNoLane (0.05s)
    heavy_run_test.go:194: a run naming no lane succeeded: lanes=0 findings=0
```

GREEN with the guard restored.

## F2 — an assertion that discriminates on the writer seam

The old evidence was "0 `^SKIP:` lines in the lane run", and the review is right that it is 0 either
way, because `run-go-tests-with-budget.sh` re-prints `go test -json` Output events only when the
package fails. The new test does not look at the lane's output at all. It swaps this process's
`os.Stderr` for a temporary file for the duration of the test, runs the `skip-lane` fixture through the
same helper every other test uses, then fails if any captured line starts with `SKIP:`.

It swaps the `os.Stderr` *variable*, not file descriptor 2, on purpose: that variable is the seam.
`handleRunHeavyVerification` passes it as the lane output writer, and the reverted helper passes it too.

**RED by ablation** (`runHeavyLanes` reverted to pass `os.Stderr` instead of `io.Discard`):

```
--- FAIL: TestNestedLaneOutputNeverReachesThisProcessStderr (0.06s)
    heavy_run_test.go:222: a fixture lane's announcement reached this process's stderr:
    "SKIP: no browser is available"; the enclosing heavy lane reads that as its own skip
```

GREEN with the seam restored. Acceptance criterion 5 is now met for this fix.

## F3 — the scanner scans the condition, and the widening is pinned

Confirmed the review's claim first: all five rewrites passed the old regex, `old-scanner-hits=0`.

`quiet_grep_pipeline_offenders` is now a named function taking the file to scan, so the same scan runs
over the script itself and over a fixture. It works on the defect's three ingredients rather than one
spelling:

- **logical lines, not physical** — a line ending in `|` or `\` is joined with the next before matching,
  which catches the pipe-at-end-of-line rewrite;
- **any stage after a pipe** — the joined line is split on `|` (with `||` neutralised first, since a
  logical OR is not a pipe) and every stage from the second on is examined, which catches
  `LC_ALL=C grep`, `command grep` and `/usr/bin/grep` without a prefix list;
- **any early-leaving option** — `-q`, a bundled short flag containing `q`, `--quiet`, `--silent`, and
  `-m` / `--max-count`, which leaves early for the same reason.

`verify_the_quiet_grep_scanner_catches_every_evasion` writes a five-line fixture in a temp directory and
requires all five to be caught. The fixture's pipe character is interpolated from a variable so this
script's own source does not carry the shape its own scanner looks for.

**Evidence, widened scanner:** real file 0 offenders; fixture 5 of 5 —

```
1: tar tzf archive.tgz |  grep -q marker
3: tar tzf archive.tgz | grep --quiet marker
4: tar tzf archive.tgz | grep --silent marker
5: tar tzf archive.tgz | LC_ALL=C grep -q marker
6: tar tzf archive.tgz | command grep -q marker
```

**RED by ablation** (scanner restored to the pre-widening regex):

```
the quiet-grep scanner caught 0 of 5 evasions:
FAIL: the quiet-grep scanner misses an ordinary spelling of the pipeline it exists to forbid
```

**Regression, the nine original sites still covered:** reverting the real site at line 735 to
`tar tzf … | grep -q 'requested-branch-marker\.txt'` produces exactly one offender naming that line;
restoring it returns 0.

**What a source scan still cannot catch, named as asked.** Two shapes are out of reach of any regex over
this file's text, and neither is a spelling of `grep -q`:

1. **A reader that leaves early without saying so.** `writer | head -1`, `writer | sed -n '1p;q'`,
   `writer | read -r line` and `writer | awk 'NR==1{print; exit}'` have the same SIGPIPE mechanism with
   no quiet flag to key on. Flagging every `| head` would flag correct code; the scanner names the
   grep family because that is the family this file uses and the one the defect was found in. The
   general rule stays prose in `_dev/primes/prime-shell-commands.md`.
2. **A pipeline assembled at runtime** — `$matcher_command`, `eval`, or a here-doc body executed later.
   The text the scanner reads is not the text the shell runs.

The known false-positive direction (the shape inside a string, a here-doc body, or a trailing comment)
is unchanged from the review's minor finding and was not addressed; the widening does not make it worse,
and the whole-file scan is still 0 offenders.

## F4 — both marker archives asserted readable

Two lines, each the same explicit check the fallback archive already carries at line 619, placed
immediately before the listing is captured. The comment states why capturing the listing is not free:
it discards the exit status the old pipeline form reported through `pipefail`.

**Truncated-archive probe reproduced** — marker committed first, 200 pad files, last 200 bytes cut:

```
full=2844 truncated=2644
tar tzf truncated.tar.gz | head -1  ->  requested-branch-marker.txt
tar tzf truncated.tar.gz            ->  status=2
without-readability: fail_count=0     <- the review's demonstrated false pass
with-readability:    fail_count=1     <- FAIL: the requested branch archive is not readable
```

## F6 — the red verdict keeps the lane's announcement

`runOneLane` no longer erases `announcedSkipLine` when the status is non-zero; `Skipped` is now
`exitStatus == 0 && announcedSkipLine != ""`, so the record is unchanged and the announcement survives.
`laneRedFinding` takes the line and adds a second evidence entry, `the lane also announced: …`, only
when there is one.

**RED by ablation** (announcement dropped again):

```
--- FAIL: TestRunLanesReportsARedLaneThatAlsoPrintedASkipLine (0.06s)
    heavy_run_test.go:177: the red finding dropped the lane's own announcement:
    []string{"heavy lane skip-then-fail-lane exited 4 after 0s"}
```

## Verification

Every command below ran under the prescribed environment.

| Check | Result |
|---|---|
| `go test ./internal/heavyverification/` | ok, 14.2s |
| `gofmt -l` on the package | empty |
| `go vet` on the package | exit 0 |
| `bash -n _dev/tests/update-script-behavior.sh` | exit 0 |
| `shellcheck --severity=warning` on the same file | exit 0 |
| heavy lane `updater` | **exit 0**, 70s, 0 FAIL lines |
| heavy lane `do-work-cli-integrations` | **exit 0**, 798 tests, 27s |

The full canonical gate was not run, as instructed. `RunLanes` has one production caller,
`runHeavyVerificationLanes`, and it always supplies `LaneIDs`; no shipped script or action passes an
empty lane set, so the new refusal cannot fire on an existing path.

## One observation, not acted on

`_dev/tests/maintainer-verify.sh:41` carries a comment saying "the lane runner keys `skipped` on the
SKIP: prefix". It is incomplete in the same way the prime was, but it is a maintainer script comment
rather than a contract statement, and the file is outside the write set even after the widening. Left
alone rather than widening scope twice.
