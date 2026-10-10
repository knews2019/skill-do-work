# REQ-690 hand-back (do-work status action and do-work-cli run-status)

- Branch: `worktree-agent-REQ-690-run-status-action`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-690-run-status-action`
- Base commit: `bd56c4b0` (checked with `git rev-parse --short HEAD` before the first edit)
- Commits (oldest first), every subject starts with `[REQ-690]`:
  - `b5f1d587` add queue-kanban open-work --format json with placement reasons
  - `135ac34f` add do-work-cli run-status: one class, ETA and remedy per open REQ
  - `516dacea` add do-work status action, routing row, guide and doc pointers
  - `972b2956` split run-status text rendering into its own file
- Nothing under `do-work/` was created, edited, staged or committed on the branch. No CHANGELOG, VERSION or mirror touched. No built binary in `git status`.

## File manifest

| File | Change | What changed |
| --- | --- | --- |
| `skills/do-work-board/tools/queue-kanban/model.go` | modified | `RequestTicket.PlacementReason`, set inside each `bucketColumns` arm (waiting, earmarked, ready, claimed, needs-input, unrecognized). Terminal arms set none. |
| `skills/do-work-board/tools/queue-kanban/open_work.go` | modified | `openWorkFacts` / `openWorkRequestFacts` and `writeOpenWorkJSON`: open tickets by column, reason, unmet deps, assigned_to, and the board's own last activity via `attachRequestActivity` over one `readWorktreeAgentGitState` read. Threshold from `staleClaimThreshold`. |
| `skills/do-work-board/tools/queue-kanban/main.go` | modified | `open-work --format text|json` (default text, other values exit 2), synopsis comment updated. Text digest unchanged. |
| `skills/do-work-board/tools/queue-kanban/open_work_test.go` | modified | `TestOpenWorkJSONPlacesWaitingAndEarmarkedTicketsWithTheirReasons` through `buildBoard`. |
| `skills/do-work-board/actions/board.md` | modified | One clause in the open-work bullet naming `--format json` and its consumer `../../do-work/actions/status.md`. |
| `skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go` | new | The `run-status` command: options, board facts, run directory, frontmatter facts, builder branches, run-local files, classes C1 to C8, remedies. |
| `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_render.go` | new | Text and `--watch` rendering (split out, see D-05). |
| `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go` | new | The nine probe-named tests, one fixture builder. |
| `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` | modified | One import and one registration loop (placed after the doctor loop). |
| `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` | modified | `RunStatusResult`, `RunStatusRow`, `RunLocalFile`; optional `CommandResult.RunStatus` (`run_status`). `CommandFinding` unchanged. |
| `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` | modified | `frozenEstimate` exported in place as `FrozenEstimate` (one caller updated). |
| `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` | modified | One Package-routing bullet for `internal/runstatus/`, appended at the end of that list. |
| `skills/do-work/actions/status.md` | new | The action: build board, 4-line shell (mktemp, open-work JSON, run-status, rm), class table, this-session checks, `/loop 15m /do-work status --watch`, two earned Rules. |
| `skills/do-work/SKILL.md` | modified | Row above clarify; `status [REQ] [--watch]` before `| roadmap` in the hint. |
| `skills/do-work/actions/help.md` | modified | One menu line between forensics and roadmap. |
| `skills/do-work/docs/status-guide.md` | new | Short user guide in the roadmap-guide shape. |
| `skills/do-work/actions/fan-out-reference.md` | modified | Mid-Run Messages progress bullet: "read through `actions/status.md`". |

## P-A-U

**[PLAN]** Followed the pre-dispatch Plan (D-01 to D-04) exactly. Board side first (test, then `PlacementReason` + JSON writer + flag), then core side (tests, result field, `FrozenEstimate` export, package, registration), then action, routing and docs. The action's shell stays at 4 lines, so no launcher script (D-14).

**[APPLY]** As planned, plus one split of the renderer into `run_status_render.go` after `run_status.go` reached 612 lines (D-05). End-to-end run against the main tree's real queue (read-only, binaries built into a `mktemp -d`): 12 rows, 9 C3 for the siblings whose hand-backs had landed, 1 C5 (REQ-656 waiting on REQ-655), 2 C1; run-local lines for `INTEGRATOR-GUIDE.md`, `PREDISPATCH-GUIDE.md`, `helpers`; last line `next: do-work run`; `--watch` printed 14 lines.

**[UNIFY]** `git diff bd56c4b0 --stat`:

```
 skills/do-work-board/actions/board.md              |   2 +-
 skills/do-work-board/tools/queue-kanban/main.go    |  14 +-
 skills/do-work-board/tools/queue-kanban/model.go   |  13 +
 .../do-work-board/tools/queue-kanban/open_work.go  |  72 +++
 .../tools/queue-kanban/open_work_test.go           |  62 +++
 skills/do-work/SKILL.md                            |   3 +-
 skills/do-work/actions/fan-out-reference.md        |   2 +-
 skills/do-work/actions/help.md                     |   1 +
 skills/do-work/actions/status.md                   |  78 +++
 skills/do-work/docs/status-guide.md                |  58 +++
 .../tools/do-work-cli/cmd/do-work-cli/main.go      |   4 +
 .../internal/nextselection/next_selection.go       |   6 +-
 .../internal/resultmodel/result_model.go           |  48 ++
 .../do-work-cli/internal/runstatus/run_status.go   | 522 +++++++++++++++++++++
 .../internal/runstatus/run_status_render.go        | 102 ++++
 .../internal/runstatus/run_status_test.go          | 284 +++++++++++
 .../tools/do-work-cli/lessons-do-work-cli.md       |   1 +
 17 files changed, 1266 insertions(+), 6 deletions(-)
```

Checks, run from the worktree root after the last commit (machine under load from sibling builders; each ran once and passed):

| Check | Exit | Wall |
| --- | --- | --- |
| `bash .../REQ-690-probe.sh` (prints `REQ-690 GREEN probe: PASS`) | 0 | 2 s |
| `gofmt -l` over runstatus, resultmodel, nextselection, cmd/do-work-cli, queue-kanban/*.go | empty | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/runstatus/ ./internal/resultmodel/ ./internal/nextselection/ ./cmd/do-work-cli/` | 0 | <1 s |
| `go vet -C skills/do-work-board/tools/queue-kanban .` | 0 | <1 s |
| `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/resultmodel/` | 0 | 1 s |
| `go test -C skills/do-work/tools/do-work-cli -count=1 -run 'FrozenEstimate|Simple' ./internal/nextselection/` | 0 | <1 s |
| `go test -count=1 -v -run '^TestRunStatus' ./internal/runstatus/` (9 PASS, 0 SKIP) | 0 | 1 s |
| `go test -count=1 -run '^(TestOpenWork|TestBucketColumns$|TestAssignedPendingRequestIsEarmarkedNotReady$)'` queue-kanban (10 PASS) | 0 | <1 s |
| queue-kanban broader subset `-run 'Bucket|Column|Pending|Earmark|Waiting|BuildBoard|LoadBoard|Ticket|Generated'` (after the model change) | 0 | 11 s |
| `bash _dev/tests/action-shell-blocks.sh` (76 fences, ShellCheck enabled) | 0 | 2 s |
| `bash _dev/tests/shipped-package-reference-contract.sh` | 0 | 1 s |
| `git diff --check bd56c4b0` | 0 | <1 s |
| `git status --short` (clean, no `queue-kanban` binary) | clean | <1 s |

Files reviewed in the final diff: all 17 above. Checked: no debug output, no `os.Remove` / `syscall.Kill` / `FindProcess` in `internal/runstatus/`, comments match behavior, text digest code path untouched.

Not run (by the brief): the repository gate, the whole queue-kanban suite (browser lanes), `contract-regressions.sh` (writes the test-duration log into the tree; left to the integrator's gate).

## Proof record

- **RED probe** (2026-10-10T13:33:27Z, before any edit): `missing skills/do-work/actions/status.md`, exit 1.
- **RED board test** (before production code): `go test -run TestOpenWorkJSON` → build failed: `undefined: writeOpenWorkJSON`, `undefined: openWorkFacts`, `undefined: openWorkRequestFacts`.
- **RED runstatus tests** (before production code): `go test -run '^TestRunStatus' ./internal/runstatus/` → build failed: `undefined: runStatusAt`, `resultmodel.CommandResult has no field or method RunStatus`, `undefined: resultmodel.RunLocalFile`, `undefined: Handlers`.
- **First GREEN attempt**: 8 of 9 passed; `TestRunStatusStaleClaimNamesRecoveryOnlyInTheStopReason` failed on wording ("Only if" vs the asserted "only if"); the stop reason was reworded, then 9 of 9 passed.
- **Routing walk** (SKILL.md, first match wins, rows checked top to bottom):

| Input | Rows passed over | First match |
| --- | --- | --- |
| `do-work status` | help, run-simple-reqs, run-with-recovery, run, check-for-updates, verify, review, stakeholder-answers: no trigger word | status row (`status`) → `./actions/status.md` |
| "status and ETA please" | same rows, none match | status row (`status and ETA`) → `./actions/status.md` |
| "is REQ-12 stuck?" | same rows, none match | status row (`stuck`) → `./actions/status.md` |
| "blocked" | same rows, then the status row (it carries no `blocked` trigger) | clarify row (`blocked`) → `./actions/clarify.md` |

SKILL.md diff: the hint gains `status [REQ] [--watch] |` right before `roadmap`; one new row directly above the clarify row (`status`, `status and ETA`, `eta`, `is it stuck`, `stuck`, `why is REQ-N taking so long`, `how do I unblock`, `why is REQ-N in` → `./actions/status.md`). The probe checks the same order mechanically.
- **GREEN probe** (2026-10-10T13:41:56Z and again after the split commit): `REQ-690 GREEN probe: PASS`, exit 0; 9 runstatus tests PASS, 0 SKIP; board tests PASS, 0 SKIP; `run-status` listed by `do-work-cli --format json help`.

## Decisions

Class precedence shipped as item 4: C8, C3, C4, C5, C6, C7, C2, C1. D-01 to D-04 followed as written.

- **D-05 DECIDE & STATE: renderer split into a 17th file.** `run_status.go` reached 612 lines; brief item 3 says split when one file passes about 500. The text/`--watch` renderer and its four display helpers moved to `run_status_render.go` (102 lines); `run_status.go` is 522. Value: the brief's size rule holds. Risk: one path beyond the Scope list of 16; reversible by folding it back.
- **D-06 DECIDE & STATE: builder-branch read follows the brief's `for-each-ref` without `--no-merged`.** It is display only (`builder_branch`, `builder_branch_tip_at`, `other_builder_branches`); no class reads it. Classes use the board's activity, which does apply `--no-merged` (git-history-evidence). Risk: a builder branch with no commits of its own shows the integration tip's date. Value: no integration-ref guess in core. See Discovered Tasks.
- **D-07 DECIDE & STATE: ETA text for a row with an estimate but no claim** is `N min estimate, not claimed`. With no estimate it is `no estimate`. A claimed row past p50 is `over estimate by N min` and `eta_minutes` is absent, never negative.
- **D-08 DECIDE & STATE: phase column** shows the board's `last_activity_phase`, or the board column name for an unclaimed row, or `-`.
- **D-09 DECIDE & STATE: rows are listed in class precedence order** (stable within a class, in board order), so the most urgent rows come first.
- **D-10 DECIDE & STATE: the last line** is `next: <argv>` of the first row (in that order) with a non-empty `next_argv`; with none, `next: nothing to run now; check again with \`do-work status\``.
- **D-11 DECIDE & STATE: two refusals beyond the brief.** Board facts without a positive `stale_claim_threshold_minutes` are refused (a zero threshold would class every claim C7), and a `--run` that is not a directory is refused (a typo would otherwise silently drop every run field). Both use outcome `refused`, `next_argv` `["do-work","status"]`. Not earned by an incident; the integrator may delete either (3 lines each).
- **D-12 DECIDE & STATE: frontmatter wins over the board** for `status`, `claimed_at`, `assigned_to` when the REQ file is found (the queue copy before an archived one); board values are the fallback.
- **D-13 DECIDE & STATE: placement sentences.** Waiting and Earmarked restate `board.md` and the badge tooltip as the brief asked; Ready, Claimed, Needs input and the unrecognized-status arm also get one plain sentence so every open ticket in the JSON carries a reason. No `web/`, payload or `verify.go` change.
- **D-14 DECIDE & STATE: no launcher script.** The action's shell is 4 lines (`mktemp`; open-work JSON `&&` run-status; `rm -f`).
- **D-15 DECIDE & STATE: `do-work/lessons-index.md` not refreshed.** The added line is a Package-routing bullet, not a `[family:]` lesson, and `do-work/` is outside my write boundary. The integrator may refresh the satellite's size estimate.
- **D-16 DECIDE & STATE: `--watch`** prints one header line (time, row count, class counts, run-local file count), one line per row and the `next:` line: 10 lines for 8 rows. Run-local files appear only as a count there.
- **D-17 DECIDE & STATE: `verification_argv`** is `do-work-cli run-status --board-facts <the given file> [--run DIR] --req REQ-NNN`. Through the action that file is a deleted temp file, so re-verifying means re-running `do-work status`. Kept because the brief names the run-status argv.

## Discovered Tasks

- impact-negligible: `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` Read-first lines for `main.go` (flag list) and `open_work.go` ("the headless in-flight digest") do not mention `open-work --format json`; restatement drift per paired-predicate-drift → report only
- impact-negligible: `run-status`'s builder-branch field could adopt the board's `--no-merged=<integration ref>` rule so an unowned branch tip never shows as builder activity (D-06) → report only
- impact-negligible: `do-work/lessons-index.md` row for `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` may need its size estimate refreshed for one added Package-routing bullet (D-15) → report only

## Lessons read

- Required: `_dev/primes/lessons-releases.md` (whole).
- Families from the brief: [destructive-next-argv] `lessons-do-work-cli.md:87` (drives C7's read-only `next_argv`); [paired-predicate-drift] `lessons-do-kanban.md:18`, `:61` (JSON placement pinned through `buildBoard`; column and reason set in the same switch arm); [git-history-evidence] `lessons-do-kanban.md:22` (activity reused from the board, not re-correlated).
- `lessons-do-work-cli.md` § Package routing and § Package direction (no `queue-kanban` import; `runstatus` imports `doctor`, `nextselection`, `repositorymodel`, `requestmodel`, `resultmodel`, `commandruntime`, all acyclic).
- `lessons-action-files.md` § Template only.

## Anti-bloat check

`git diff --stat` is in [UNIFY] above. Added beyond what the brief names, each with its reason:

- `run_status_render.go` (file): the brief's split rule (D-05).
- `run_status.go` helpers: `parseCommandOptions`, `readBoardFacts`, `resolveRunDirectory`, `readManifestLines`, `manifestRowFor`, `builderBranches`, `listRunLocalFiles`, `firstLineOf`, `openRequestFile`, `countOpenQuestions`, `estimateRemaining`, `buildRow`, `classifyRow`, `verificationArgv`, `refusal`, `minutesBetween`, `displayPath`: each implements one fact or step item 3 to 6 names.
- `run_status_render.go` helpers: `renderReport`, `phaseOf`, `minutesText`, `shortTitle` (the six-word title cut), `fallbackDash`: item 5's text and watch shapes.
- Package values: `quietActivityBoundary` (the named 20-minute constant), `runLocalLabel`, `classPrecedence`, `classNames`, `reportedColumns` (D-04), four regexes (REQ id, `REQ-NNN-` prefix, UTC instant, `- [ ]` item).
- Types: `boardFacts`, `boardRequest` (decode the board JSON), `commandOptions`, `statusRow`; board `openWorkFacts`, `openWorkRequestFacts`; resultmodel `RunStatusResult`, `RunStatusRow`, `RunLocalFile`.
- Two refusals (D-11). Flagged as the only guards not earned by an incident.
- Tests: exactly the ten the probe names; helpers `writeRunStatusFixture` (the one fixture builder), `writeFixtureFile`, `claimedFixture`, `runFixture`, `findingFor`, `textOf`.

No new flags beyond `--format` (board) and `--board-facts`, `--run`, `--req`, `--watch` (core). No just recipe, no `--fix`, no lock or PID logic in Go.

## Proposed CHANGELOG entry

**do-work status: one class, ETA and remedy per in-flight REQ**

Asking "is it stuck?" used to mean rebuilding the answer by hand from REQ files, the run manifest, builder branches and `ls` of the run directory. `do-work status` now answers it in one table, from the same facts the board shows, and never calls a run dead.

- New action `do-work status [REQ] [--watch]` (`actions/status.md`, `docs/status-guide.md`), routed from "status", "status and ETA", "eta", "is it stuck", "stuck", "how do I unblock" and "why is REQ-N taking so long / in". The bare word "blocked" still routes to clarify.
- New `do-work-cli run-status --board-facts FILE [--run DIR] [--req REQ-NNN] [--watch]`: one `CommandFinding` per claimed, blocked, waiting or earmarked REQ with class C1 to C8 (first match wins: C8, C3, C4, C5, C6, C7, C2, C1), ETA from the frozen p50 ("over estimate by N min", never negative), and a typed `run_status` block with ages, manifest row, hand-back presence, builder branch and run-local files labelled "run-local, not interpreted".
- A claim past the 3-hour threshold offers the read-only `git log --full-history`; `run-with-recovery` is named only in the stop reason, with what it destroys.
- `queue-kanban open-work --format json` prints the board's open tickets with column, placement reason and last activity. The text digest is unchanged.
- `/loop 15m /do-work status --watch` gives a compact tick of about 10 lines for 8 open REQs and writes nothing.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`.

- [family: cross-module-fact-handoff] [REQ-690: when do-work-cli needs a fact only queue-kanban computes (a column, its reason, last activity), have the board print it as JSON and pass the file in; the core module never imports, builds or launches the board, and never re-derives the predicate](../../../../do-work/archive/UR-153/REQ-690-run-status-action.md)

## Integration seams touched

- `skills/do-work/SKILL.md`: hint token right before `| roadmap` (shared line with REQ-691's `trace` token); row directly above clarify (separate hunk from REQ-689's row above run and REQ-691/692's rows above capture).
- `skills/do-work/actions/fan-out-reference.md`: one clause in Mid-Run Messages only (REQ-689 edits the run-directory table and coordinator rules).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: one bullet appended at the end of § Package routing (REQ-660/659/658 also append bullets; expect an adjacent-line conflict there, resolve by keeping all bullets).
- `cmd/do-work-cli/main.go`: import inserted after `requeststate`; loop inserted after the `doctor` loop. `result_model.go`: three new types before `CommandResult`, one field after `LifecycleTiming`.
- REQ-689 contract: `do-work status --watch` output is header, one line per row, `next:` line; `full-gate.lock` is listed as run-local (not REQ-prefixed), matching D-03.

Test wall times: runstatus 0.5 s, board JSON subset 0.3 s, resultmodel 0.2 s, nextselection subset 0.2 s, probe 2 s.
