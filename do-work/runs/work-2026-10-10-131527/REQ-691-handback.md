# REQ-691 hand-back: do-work trace action and queue-kanban request-commits

- Branch: `worktree-agent-REQ-691-trace-coverage-action`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-691-trace-coverage-action`
- Base commit: `bd56c4b0` (recorded before the first edit)
- Commits: `e8e4e9c8` `[REQ-691] add do-work trace action and queue-kanban request-commits` (one commit)
- No path under `do-work/` was created, edited, staged or committed on the branch. This hand-back file is not staged.

## File Manifest

- `skills/do-work/actions/trace.md` (new): the action. Description with the core-package justification, When to Use, Input, six Steps (read, split, search, print, ask, hand off on go), and three Rules.
- `skills/do-work/actions/trace-reference.md` (new): Source Handling, Search Order, Commit Evidence (the `request-commits` call and the D-04 fallback), Coverage Table (template, three-row example, criteria count, date, four verdicts), Go-Ahead Payload.
- `skills/do-work/SKILL.md` (modified): one routing row directly above the capture fallback row, and `trace <url|image|text|UR-NNN>` in the argument hint right after `roadmap`.
- `skills/do-work/actions/help.md` (modified): one two-line menu entry directly after the `roadmap` line, in the same wrapped style as `run-with-recovery`.
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (modified): the per-commit attribution is extracted into `requestIdsCreditedByCommit`. `correlateCommitsToRequests` calls it, with no change in behaviour.
- `skills/do-work-board/tools/queue-kanban/request_commits.go` (new): `runRequestCommitsCommand(args, standardOut, standardErr, runner) int`. It makes one full-history `git log` call with no pathspec, reuses the shared helper, and prints TSV output. Exit codes are 0, 1 and 2.
- `skills/do-work-board/tools/queue-kanban/request_commits_test.go` (new): the two named tests.
- `skills/do-work-board/tools/queue-kanban/main.go` (modified): the switch case, the synopsis line and the unknown-subcommand list.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified): `request-commits REQ-NNN…` added to the subcommand list.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Followed the REQ's Plan D-01 to D-08 without change. Board: extract one attribution helper, add a read-only `request-commits` with the frontmatter command's testable shape (writers in, exit code out). The git runner is passed in so the tests stub it. Prose: keep `trace.md` short and put the long parts in `trace-reference.md`. Routing row above the capture fallback. One help line. Tests first for the Go RED. Then a hand dry run in a scratch fixture repo for the GREEN cases, including the go-ahead capture through the real `capture-files`.
- [x] **[APPLY]:** Wrote the two tests and ran them against a stub. Then wrote the helper, the command and the wiring, then the two action files, the routing row, the hint and the help line. The dry run showed one case the reference did not cover: an ask narrower than its REQ. I added one clause for it (D-10). No file outside Scope was touched.
- [x] **[UNIFY]:** `git diff bd56c4b0 --stat` (below) shows exactly the nine Scope files. Checked files: all nine. I read each diff. There is no debug output and no build artifact; the built `queue-kanban` binary is gitignored and was removed after the dry run. Check results:

| Check | Exit | Wall |
|---|---|---|
| `bash …/REQ-691-probe.sh` (worktree root) | 0 | 5 s (first run 4.5 s) |
| `go test -C skills/do-work-board/tools/queue-kanban -count=1 -run 'RequestCommits\|RequestActivity\|RequestPathPattern' .` | 0 | 4 s |
| `gofmt -l skills/do-work-board/tools/queue-kanban` | prints nothing | <1 s |
| `go vet -C skills/do-work-board/tools/queue-kanban ./...` | 0 | <1 s (cached) |
| `git diff --check` | 0 | <1 s |
| `bash _dev/tests/contract-regressions.sh` (router budgets, action template rules) | 0 | 31 s |
| `bash _dev/tests/staged-skills-contract.sh` | 2, refused: "heavy-only; run maintainer-verify.sh --heavy". Not run; it belongs to the integrator's gate | 0 s |

No wall-time budget failed, so nothing was rerun.

```text
 .../tools/queue-kanban/activity_correlation.go     | 28 +++++---
 skills/do-work-board/tools/queue-kanban/main.go    |  5 +-
 .../tools/queue-kanban/prime-do-kanban.md          |  2 +-
 .../tools/queue-kanban/request_commits.go          | 78 ++++++++++++++++++++++
 .../tools/queue-kanban/request_commits_test.go     | 72 ++++++++++++++++++++
 skills/do-work/SKILL.md                            |  3 +-
 skills/do-work/actions/help.md                     |  2 +
 skills/do-work/actions/trace-reference.md          | 69 +++++++++++++++++++
 skills/do-work/actions/trace.md                    | 49 ++++++++++++++
 9 files changed, 295 insertions(+), 13 deletions(-)
```

(The stat above was taken before the commit. The commit holds the same nine files.)

## Proof Record

### 1. RED at base (bd56c4b0, before any edit)

```text
test -f skills/do-work/actions/trace.md        -> exit 1
grep -n 'actions/trace.md' skills/do-work/SKILL.md -> no output, exit 1
bash …/REQ-691-probe.sh                         -> "FAIL: skills/do-work/actions/trace.md is missing", exit 1
```

At base, "is this captured? …" matches no row before `SKILL.md`'s last row, so it falls through to the capture fallback ("unmatched descriptive multi-word input").

### 2. RED, Go

The first run, before `request_commits.go` existed, was a build failure (`undefined: runRequestCommitsCommand`). For a real assertion failure, I ran again against a stub that returned 0 and printed nothing:

```text
--- FAIL: TestRequestCommitsListsOnlyCommitsTheBoardCredits
    stdout = "" want  header + REQ-007 c4 (2 paths) / c3 (0) / c1 (0)
    git arguments = [], want ["log" "--format=%H%x00%cI%x00%s" "--name-only"] (no --since window, no pathspec)
--- FAIL: TestRequestCommitsRefusesAnArgumentThatIsNotARequestId
    exit code = 0, want 2
FAIL
```

After the implementation, both pass, together with the three existing correlation tests the helper extraction touches (`…WithoutAPrefix`, `…EveryPrefixTokenInASubject`, `TestRequestPathPatternMatchesOnlyRequestFilesAndRunArtifacts`).

### 3. GREEN, fixture dry run

- **Fixture:** a scratch repo made with `mktemp -d` under the session scratchpad, outside both trees, and deleted after this record. It has three commits: `seed`; `[REQ-007] Export orders as CSV` (`src/export.py`, `tests/test_export.py`, with tests for header row and one row per order, which is 2 of the 3 Detailed Requirements); and `[REQ-007] archive request lifecycle` (`do-work/archive/UR-001/input.md`, `do-work/archive/UR-001/REQ-007-export-orders-csv.md`, status completed, `commit:` set, three Detailed Requirements bullets, and a `**GREEN when:**` line that needs comma quoting).
- **Spec:** pasted text for invented ticket `SHOP-42`: (1) a CSV with a header row; (2) export every order, with a comma in a name kept in one column; (3) the owner gets the CSV by email every Monday.

I acted as the agent and followed the worktree's `trace.md`:

- **Step 1:** loaded prompt-injection, then wrote the transcription to `mktemp -d …/trace-source.XXXXXX/pm-ticket-source.md` (outside the repo) and printed its path. Its first lines are `Ticket: SHOP-42` / `Source: pasted text` / `Fetched: 2026-10-10`.
- **Step 2:** asks A1 to A3 keep the source's numbers 1 to 3.
- **Step 3:** fixed-string search over `do-work/`. `SHOP-42`, `email` and `Monday` gave exit 1 (no match). `CSV` and `export` matched UR-001 `input.md` and the REQ-007 file.

`request-commits` was built from the worktree and run with `--repo-root <fixture>`. It printed this and exited 0:

```text
request_id	commit	committed_at	paths_outside_do_work	subject
REQ-007	480c54b662ea6230fe1b90783203f1ca87a116db	2026-10-10T16:31:09+03:00	0	[REQ-007] archive request lifecycle
REQ-007	170e7aca09df96f589a11468465e135d085abbb3	2026-10-10T16:31:09+03:00	2	[REQ-007] Export orders as CSV
```

Code and tests: the two fixture tests pass. The GREEN-when check with "Lee, Ann" gives column counts `[3, 4, 3]`, so the GREEN line is not met. `grep -rniF email src tests` gives exit 1.

The printed table (Step 4):

| Ask | REQ (plain title, status) | Commit | Criteria evidenced | Verdict (as of) | Gap | Open question |
| --- | --- | --- | --- | --- | --- | --- |
| A1 | REQ-007 (export orders as CSV), completed | `170e7ac`, 2026-10-10 | 1 of 1, bullet 1 only | complete (2026-10-10) | none | none |
| A2 | REQ-007 (export orders as CSV), completed | `170e7ac`, 2026-10-10 | 2 of 3 + GREEN no | partial (2026-10-10) | bullet 3: a comma in a name splits the row, so the GREEN line fails | none |
| A3 | none found (searched: "SHOP-42", "email", "Monday") | none | n/a | not started (2026-10-10) | all | none |

The read-only evidence is `git -C <fixture> status --porcelain` before the trace and after the no-go stop. Both printed nothing. `git -C <worktree> status --porcelain` listed only my own source edits. The build's `queue-kanban` binary did not appear because it is gitignored.

**Go-ahead branch.** I followed the worktree's `capture.md` Step 5 in the same fixture. The raw input was the A2 and A3 rows plus the transcription path. Payloads and the manifest were built by script in a `mktemp -d` scratch directory. `do-work-cli.sh --repo-root <fixture> --format json capture-files --manifest … --dry-run` gave `success` / `PUBLICATION-DRY-RUN`. Then `--commit` gave `success` / `PUBLICATION-APPLIED`, with fixture commit `e087dd5 [UR-002] capture SHOP-42 order export gaps from trace`. The resulting records:

- `UR-002/input.md`: `title: 'SHOP-42: order export gaps (comma quoting, Monday email)'`, `requests: [REQ-008, REQ-009]`. The UR holds only the partial and not-started asks. A1 appears only in the Summary line that says it is complete and was not captured again.
- `REQ-008`: `addendum_to: REQ-007`. Its `## Prior Implementation` names commit `170e7aca09df96f589a11468465e135d085abbb3`.
- `REQ-009`: the Monday email, a plain new REQ.
- Asset `UR-002/assets/REQ-008-shop-42-source-transcription.md`: first lines `Ticket: SHOP-42`, `Source: pasted text`, `Fetched: 2026-10-10`.

One refusal on the way: the first dry run refused with `PUBLICATION-MANIFEST-INVALID: manifest operation "capture" does not match command "capture-files"`. The fix was `"operation": "capture-files"`. This is not in trace's scope; see Discovered Tasks.

### 4. GREEN, source rules (walk-through)

- **URL that returns only a site name.** `trace.md` Step 1: "If a URL fetch returned only a shell, say so, ask for a paste or a screenshot, and stop." `trace-reference.md` Source Handling: "When the result is a shell (only the site name, or almost no text), say so and ask the user for a paste or a screenshot. Never judge coverage from it; a row that rests on it is `unverified`." `trace.md` Rules: "No fetched row is `complete` without fetched text."
- **Screenshot.** `trace-reference.md` Source Handling: "Image: read it directly and write a thorough text description (…), as `actions/capture.md` Step 4 does for screenshots. Describe it before the split; the asks come from the description."

### 5. GREEN, routing (walk-through, first matching `SKILL.md` row)

- "is this captured? <url>": the trace row, through `is this captured`. No earlier row has a trigger at the start of this phrase.
- "do-work trace UR-12": the trace row, through `trace`.
- "how much of it is implemented": the trace row, through `how much of it is implemented`.
- "capture-request: add X": the capture row, through `capture-request:`. It has none of the trace row's triggers, and the capture row is still the last row.

## Decisions

- **D-09, DECIDE & STATE: wording of the routing row.** The row reads: `` | `trace`, `is this captured`, `already implemented`, `how much of it is implemented`, `didn't we have a request`, or `coverage` together with a spec, URL, image or UR target | `./actions/trace.md` | ``. It is placed directly above the capture fallback. The sentence about an unknown single word, after the table, is unchanged. A bare `coverage` therefore still gets the "captured or another command?" question.
- **D-10, DECIDE & STATE: an ask narrower than its REQ.** The acceptance check says one archived REQ yields both `complete` and `partial 2 of 3`. That is only possible when the done ask covers part of REQ-007. `trace-reference.md` now says: "When an ask covers only part of its REQ, count only the criteria inside the ask and name them (`1 of 1, bullet 1 only`)." This follows requirement 15's "known scope" wording. Reversible: one sentence.
- **D-11, DECIDE & STATE: the format of the "Criteria evidenced" cell is `N of M + GREEN yes|no`.** Requirement 16 counts the Detailed Requirements bullets plus the GREEN-when line. The acceptance text says "partial with 2 of 3". This format keeps both: the bullets as `2 of 3`, the GREEN line as its own yes or no. I changed the source report's example table to this format (`3 of 3 + GREEN yes`, `2 of 4 + GREEN no`). A REQ with neither section counts its `## What` as the one criterion.
- **D-12, DECIDE & STATE: the git runner is a fourth parameter.** `runRequestCommitsCommand(args, standardOut, standardErr, runner gitCommandRunner) int` is the frontmatter command's shape plus the existing injectable runner. `main.go` passes `runGitCommand`. This avoids a package-level variable just for tests.
- **D-13, DECIDE & STATE: flags first, then ids.** The command takes the synopsis order `[--repo-root DIR] REQ-NNN…`. A flag placed after an id is read as an id. It fails the `REQ-` plus digits check and exits 2, so it is never silently dropped. That follows the board's leftover-argument rule.
- **D-14, DECIDE & STATE: "read-only until go" is stated as unchanged status.** It is not stated as an empty status. The Rules line says `git status --porcelain` is the same after a no-go trace as before it, because a user's tree can already be dirty. In the clean fixture, that means empty, which is the acceptance check.
- **D-15, DECIDE & STATE: the Commit column cites the newest `request-commits` row above 0.** A row of 0 is bookkeeping or a merge. When no row is above 0, use the REQ's `commit:` value. Any non-zero exit, a build failure included, takes the D-04 path, so a git failure never reads as "no commits".
- **D-16, DECIDE & STATE: how the content is split.** `trace.md` (49 lines) holds the triggers, the six steps and three rules. Each step points to a named section of `trace-reference.md` (Source Handling, Search Order, Commit Evidence, Coverage Table, Go-Ahead Payload), which holds every detail. The verdict definitions and the table template live only in the reference, so they are stated once. `trace.md` names the four verdict words and does not define them.
- **D-17, DECIDE & STATE: the go-ahead hand-off is one `capture-request:` invocation.** Its verbatim input is the partial and not-started rows plus the transcription file as an asset. The fixture showed that `capture-files` holds the mixed batch (one addendum REQ and one new REQ under one UR) in one transaction. No change was needed.

## Discovered Tasks

- `skills/do-work/actions/capture.md` Step 5 never shows the manifest's `operation` value. My first hand-built manifest used `"capture"` and was refused with `PUBLICATION-MANIFEST-INVALID`; `"capture-files"` is required. REQ-688 (`capture-files --example`) covers this gap. impact-low → report only
- The narrative header comment in `skills/do-work-board/tools/queue-kanban/main.go` (the paragraph before the synopsis) lists the subcommands in prose and was already missing `now`. I updated only the synopsis line and the unknown-subcommand list, as the brief says. impact-negligible → report only
- `_dev/tests/staged-skills-contract.sh` is heavy-only, so it could not run here. The integrator's gate runs it. The Plan's exploration found that no test lists core routes, so no contract edit should be needed. impact-low → report only

## Lessons Read

- Whole: `_dev/primes/lessons-releases.md`.
- Families, read with `grep -F`: `paired-predicate-drift` and `git-history-evidence` (`skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`); `example-path-read-as-citation` and `read-only-action-tool-side-effects` (`_dev/primes/lessons-action-files.md`).
- How they were applied: `git-history-evidence` means no pathspec and no merge range, and the test pins the git arguments. `paired-predicate-drift` means one shared helper. `example-path-read-as-citation` means the example file name `pm-ticket-source.md` is kept bare. `read-only-action-tool-side-effects` means the reference says the build writes only the gitignored binary in the board tool directory, and the transcription goes to `mktemp -d` outside the tree.
- Also read: the action template (`_dev/primes/lessons-action-files.md` § Template), and the primes `prime-action-files.md`, `prime-kanban-board.md`, `prime-shell-commands.md`, `prime-do-kanban.md` and `prime-releases.md`.

## Anti-Bloat Check

The `git diff bd56c4b0 --stat` is in [UNIFY] above: nine files, +295 / −13.

Things I added that the REQ or Plan did not name:

- `requestIdsCreditedByCommit` (function): this is the one extracted helper the brief expects.
- `requestIdArgumentPattern` (package-level regexp): it carries D-05's exit-2 rule for "an argument that is not `REQ-` plus digits". The board's path and subject patterns are unanchored fragments, so neither can validate a whole argument.
- The `runner gitCommandRunner` parameter (D-12): it is what lets the two named tests stub git without a package-level variable.
- One sentence in `trace-reference.md` for asks narrower than their REQ (D-10): the acceptance fixture cannot reach `complete` and `partial 2 of 3` from one REQ without it.

Nothing else was added: no new flag beyond `--repo-root`, no other helper, no third test, no guide, no `next-steps.md` row.

## Proposed CHANGELOG Entry (the integrator writes it with the version)

**Trace Action: Check How Much of a Spec Is Already Captured and Built**

`do-work trace <url | image | pasted text | UR-NNN>` answers "is this captured, and how much of it is built?" before anything is captured again. It prints one table with a dated verdict per ask, and on your go it captures only the gaps.

- New core action `actions/trace.md`, with its long parts in `actions/trace-reference.md`. It splits the source into asks A1, A2, …, searches URs, REQs, commits, code and tests for each ask, and prints one row per ask with the verdict `complete`, `partial`, `not started` or `unverified`, plus the date and one evidence pointer.
- It is read-only until go. The transcription stays in a temporary directory outside the repo. A URL fetch that returns only a site name is never judged, and trace asks for a paste or a screenshot instead.
- On go, it makes one `capture-request:` call with only the `partial` and `not started` rows. A partial match on an archived REQ becomes an addendum REQ with the commit hash in `## Prior Implementation`. A ticket id goes in the UR title and in the transcription's first lines.
- New read-only board subcommand `queue-kanban request-commits [--repo-root DIR] REQ-NNN…`. It lists every commit the board credits to each REQ (a bracketed subject token or a touched REQ file), with the hash, the date, the count of paths outside `do-work/`, and the subject. It uses the same attribution helper as the board's activity lines, so the two cannot disagree.
- Routing: the phrases "is this captured", "already implemented", "how much of it is implemented", "didn't we have a request" and `trace` go to the new action, above the capture fallback. It is also added to the argument hint and the help menu.

## Proposed Lesson Bullet

Satellite: `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`

- [family: paired-predicate-drift] [REQ-691: a second reader of the commit-to-REQ rule (`request-commits`, used by `do-work trace`) calls the one extracted `requestIdsCreditedByCommit` instead of copying the path and subject patterns, and its test pins the exact `git log` arguments, because a `--since` window or a `do-work/` pathspec is how a reuse silently narrows what it credits while every row it does print still looks right](../../../../do-work/archive/UR-153/REQ-691-trace-coverage-action.md#lessons-learned)

## Integration Seams

- `skills/do-work/SKILL.md`: my row is directly above the capture fallback. REQ-692's row goes in the same place; the triggers do not overlap, so either order works. In the argument hint I added only ` | trace <url|image|text|UR-NNN>` right after `roadmap`. A conflict with REQ-690 or REQ-689 resolves as a union.
- `skills/do-work/actions/help.md`: two lines directly after the `roadmap` line. REQ-690's `status` line near `roadmap` resolves as a union.
- `skills/do-work-board/tools/queue-kanban/main.go`: one switch case, one synopsis line, and the unknown-subcommand list ending `| open-work | request-commits)`. If REQ-690 adds a case, union the cases and the list.
- The board has no version of its own, and the write-surface count in `_dev/primes/prime-kanban-board.md` is unchanged because the command is read-only.
- Test wall times: probe 5 s, focused `go test` 4 s, `contract-regressions.sh` 31 s.
