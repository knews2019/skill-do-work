# Builder brief: REQ-695 (run-status names a missing --run folder and stops recommending do-work run while an integrator works)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-695-run-status-missing-run-integrator
- Branch: worktree-agent-REQ-695-run-status-missing-run-integrator, created by the coordinator from main HEAD with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch, never create the worktree yourself).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-695-run-status-missing-run-and-integrating-c3.md. Read it fully: What, Prior Implementation, Detailed Requirements, Red-Green Proof, Constraints, Builder Guidance, Required Lessons Dropped for Budget, and the orchestrator's sections below it (Triage with D-01 to D-04, Exploration with D-05, Scope). The release belongs to the integrator; you PROPOSE the changelog entry and a lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-154/input.md (review follow-ups R1 to R5; this REQ is R4).
- Source review (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-131527/REQ-690-review.md (F3 has the reviewer's suggested text; F4 the live sighting). Prior REQ: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/UR-153/REQ-690-run-status-action.md (D-19 and D-20 at lines 287-288).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-192201/REQ-695-handback.md
- Route B, tdd: true, impact-user-visible, effort-mechanical, domain backend. Go change in the standard-library do-work-cli module plus two prose table cells.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, backend.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md (Traps: `closed-enumeration-for-a-condition`, `silent-skip-reads-as-red`), /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md (read; the release itself is the integrator's). Lessons: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-releases.md (whole file, small). The dropped satellite `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (20802 tokens) was read at pre-dispatch; its rules for this change are: a finding's `next_argv` is followed literally, so it is read-only or empty, never destructive (`destructive-next-argv`); a silent fallback reads as a verdict (`silent-skip-reads-as-red`, which is F3 itself). Grep that file for `[family: destructive-next-argv]` only if you want the incidents.

## The change (decided at pre-dispatch, do not reopen)
All Go edits are in `skills/do-work/tools/do-work-cli/internal/runstatus/` of your worktree. Line numbers are at main HEAD `fd9a2378`.

1. **F3, display only (D-03, D-05).** Keep D-19: no refusal, no new finding, no exit-code change, no class change.
   - In `runStatusAt` (`run_status.go:104-160`), after `resolveRunDirectory` (`:115`), stat the folder once. When `runDirectory != ""` and the stat fails or it is not a directory, keep the typed path for display but pass an empty folder to the three readers: `readManifestLines` (`:132`), `buildRow` (`:140`) and `listRunLocalFiles` (`:149`). That makes the run fields absent, not false, which is the contract `TestRunStatusReportsQueueRowsWithoutARunDirectory` (`run_status_test.go:275-284`) states for no folder (D-05).
   - At `:150-152`, when the folder was not found, `report.RunDirectory` is the display path plus ` (not found; run fields absent)` (the reviewer's text in REQ-690-review.md F3). No new JSON field and no `internal/resultmodel/` edit (D-03).
   - In `renderReport` (`run_status_render.go:15-69`), in the full text form only, print `run: <report.RunDirectory>` on its own line just before the final `next:` line (`:67`) when `report.RunDirectory != ""`. Do not add it to `--watch` (`:24-43`): the 20-line budget is pinned by `TestRunStatusWatchFitsTwentyLinesForEightOpenRequests`.
2. **F4, a C3 variant, not a new class (D-01, D-02).**
   - In `classifyRow` (`run_status.go:317-361`), add one arm just before the C3 arm at `:325`: a claimed row with a landed hand-back whose `LastActivityPhase` is one the integrator opens (`integration`, `review`, `remediation`, `re-review`) and whose `MinutesSinceActivity` is non-nil and below `quietActivityBoundary` (`:39`). It sets `record.Class = "C3"`, leaves `finding.NextArgv` empty, and its remedy says an integrator is working on it now, names the phase and the minutes since last activity, and says not to start a second run. Keep it to one short sentence.
   - The phase set is the board's fixed pipeline, not a sample (`closed-enumeration-for-a-condition`): put it in one small package-level set with a comment that names its source, `skills/do-work-board/tools/queue-kanban/durations.go` `phaseMilestonesOf` (labels lowercased by `activity_correlation.go:206`), and says these are the phases stamped at or after the hand-back merge (`skills/do-work/actions/work.md:310`, `:376`).
   - Do not touch `classPrecedence`, `classNames`, the C3 arm itself, or the render code for this part: `renderReport` already skips a row with an empty `next_argv` when it picks the `next:` line (`run_status_render.go:17-23`).
3. **Tests, one per failure, in `run_status_test.go`.** Add a `lastActivityPhase string` field to `fixtureRequest` (`:24-29`) and use it at `:54`, falling back to today's `"implementation"` when empty, so no existing test changes. Add exactly these two tests (the GREEN probe names them):
   - `TestRunStatusMissingRunFolderIsNamedInTextAndJSON`: a claimed REQ claimed 200 min ago, a hand-back at `do-work/runs/work-2026-10-10-113000/REQ-NNN-handback.md`, run with `--run do-work/runs/typo`. Assert `RunStatus.RunDirectory` contains `do-work/runs/typo` and `not found`, the text contains `run: do-work/runs/typo (not found`, the row's `HandbackPresent` is nil, and the row is still C7 (D-19 kept). Comment one line on the failure it pins.
   - `TestRunStatusIntegratingHandbackDoesNotRecommendRun`: a claimed REQ with a landed hand-back, `lastActivityPhase: "integration"`, last activity 5 min ago (the `claimedFixture` default). Assert code `C3`, `len(NextArgv) == 0`, `AutomationStopReason` contains `integrator`, and the text's `next:` line does not contain `do-work run`. Comment one line on the failure it pins.
   - Use a REQ number not used by other tests in the file (they use REQ-801 to REQ-817).
4. **Prose (D-04), one cell each.** `skills/do-work/actions/status.md:48` remedy cell becomes `` `do-work run`; none while an integrator worked on it in the last 20 min (phase integration or later) ``. `skills/do-work/docs/status-guide.md:26` last cell becomes `` `do-work run` integrates it; nothing while an integrator is already working on it ``. The GREEN probe checks that the C3 line in each file contains the word "integrator". No other prose edit: `status.md` Step 4 already says to relay the command's text, which carries the new `run:` line.

## Anti-bloat (YAGNI), verbatim from the batch constraints; the maintainer asked for this to be watched
- Smallest change that fixes the named failure. Prefer deleting or reusing existing code over adding.
- No new helpers, flags, options, config, abstractions or files unless this REQ names them. A tiny helper for "integrator is working" is acceptable only if the arm's condition does not fit on one readable line.
- Tests only pin the named failure (one RED case per failure); no decorative or "while we're here" tests.
- Do not fix adjacent things you notice; record them as discovered tasks instead.
In the hand-back, paste `git diff --stat` and list every function, constant, option, file or test you added that the brief did not name (expected: the phase set and the fixture field only). Each one gets a one-line reason, or you remove it.

## Write boundary
Exactly these five files:
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go`
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_render.go`
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go`
- `skills/do-work/actions/status.md`
- `skills/do-work/docs/status-guide.md`
Anything else: stop and say so in the hand-back.

## Hard rules
- Every commit subject on your branch starts with `[REQ-695]` (for example `[REQ-695] name a missing run folder in run-status`). The board credits commits by that prefix only.
- Never create, edit, stage or commit any path under `do-work/`, in the worktree or anywhere else, except the one hand-back file above (which you write but never stage or commit). The integrator's queue guard refuses a branch that touches `do-work/`.
- Never merge, rebase, push, or check out another branch. Never run `recover`, `advance`, any `do-work-cli` lifecycle command, `just do-work-update`, or the repository gate (`_dev/tests/maintainer-verify.sh`); the integrator owns the gate.
- Do not touch `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` or any version mirror: the release belongs to the integrator.
- Go: run focused `go test -run` and `go vet` on `./internal/runstatus/` only, never the whole module and never the repository gate.
- Scratch files and fixtures go in `mktemp -d` directories outside both trees and are never committed. Clean up anything you start.
- Load is high while four builders run at once: if a wall-time budget fails once under load, rerun it once and record both runs.

## Integration seam
No file overlap with the other three builders. REQ-693 (worktree status/cleanup survive a missing worktree folder) edits `internal/cleanup/`, `internal/corehelpers/handoff.go` and the board's `verify.go`. REQ-694 (`req append-section` and `frontmatter set` guards) edits `internal/corehelpers/request_writers.go` and its test. REQ-696 (stale-wording sweep) edits `skills/do-work/actions/forensics.md`, `README.md`, `clarify.md`, `capture.md`, `_dev/tests/staged-skills-contract.sh` and one fixture. Shared seam: all four may propose CHANGELOG entries and lesson bullets for `lessons-do-work-cli.md`; the integrator writes those one at a time. You edit no shared file.

## Proof to run and record (from the REQ's Red-Green Proof)
1. RED: write the two tests first and run them against the unchanged production code: `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestRunStatusMissingRunFolderIsNamedInTextAndJSON|TestRunStatusIntegratingHandbackDoesNotRecommendRun)$' ./internal/runstatus/`. Both must fail: the first because nothing says "not found" and `HandbackPresent` is false, the second because `next_argv` is `do-work run`. Record each failure message. Commit the tests together with the fix (one commit preferred), but record the RED run in the hand-back.
2. GREEN: the same command passes after the change, then `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^TestRunStatus' ./internal/runstatus/` passes (every existing run-status test still green, including `TestRunStatusLandedHandbackIsClassC3`, the "other C3 rows are unchanged" pin).

## Verify before hand-back (from the worktree root, wall times recorded)
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-192201/REQ-695-probe.sh` run from your worktree root (it reads paths relative to `git rev-parse --show-toplevel`, so it checks your tree): exit 0. At base it fails with "no PASS line for TestRunStatusMissingRunFolderIsNamedInTextAndJSON"; that is expected.
- `gofmt -l skills/do-work/tools/do-work-cli/internal/runstatus` prints nothing; `go vet -C skills/do-work/tools/do-work-cli ./internal/runstatus/` is clean.
- `git diff <base> --stat` shows only the five files above; `git diff --check` clean.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path, base commit (`git -C <worktree> rev-parse --short HEAD` before your first edit) and every commit hash (one commit preferred).
- File manifest: each file with (new)/(modified)/(deleted) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists `git diff --stat`, every check's exit code and wall time, and each file checked.
- Proof record: the RED run (both failure messages), the GREEN run (test names and PASS lines), the full `^TestRunStatus` run.
- `## Decisions` (continue from D-06; each DECIDE & STATE, or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellites and families).
- Anti-bloat check (see above): the `git diff --stat` and the list of anything added that the brief did not name.
- A proposed CHANGELOG entry: a descriptive title that says what shipped, in plain words, one or two sentences on why it matters, then specific bullets. The integrator writes it with the version.
- A proposed lesson bullet, naming the satellite you think it belongs in, in that file's bullet shape (`- [family: <slug>] [REQ-695: <one-line lesson>](<relative archive link; the integrator fixes the path>)`). Say "none" if nothing was learned that a later builder would need.
- Integration seams (see the note above) and test wall times.
