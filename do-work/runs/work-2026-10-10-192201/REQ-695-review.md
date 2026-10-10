# Review: REQ-695 (run-status names a missing --run folder and stops recommending do-work run while an integrator works)

**Approve** — both review gaps from REQ-690 (the run-status report) are closed with small, correct code, and four small follow-ups are cheap enough to do now.
Route B | range `d2991804..5d4e136c` (builder commit `9fd11e18`, merge `5d4e136c`) | orchestrated mode, independent reviewer

## What's built

- A mistyped or missing `--run` folder no longer passes silently. JSON `run_directory` reads `<path> (not found; run fields absent)`, hand-back and manifest fields are absent (not `false`), and the full text form prints `run: <folder>` just before `next:`. There is no refusal and no exit-code change, so D-19 (the coordinator ruling that a missing folder is reported, not refused) still holds.
- A C3 row (claimed, hand-back landed) whose last activity is in phase integration, review, remediation or re-review and less than 20 minutes old now has an empty `next_argv` and says an integrator is working on it. All other C3 rows still recommend `do-work run`.
- The C3 cells in `skills/do-work/actions/status.md:48` and `skills/do-work/docs/status-guide.md:26` name the exception.

## Decisions / risks for you

- `integratorPhases` is a copy of four board labels with no check keeping them equal. A cheap shell check following an existing precedent closes this (F2).
- No test pins the 20-minute boundary on the new arm. Without the boundary, a dead integrator's REQ would say "an integrator is working" forever and never offer `do-work run` (F1).

## Correctness of the two fixes

**F3 of REQ-690 (missing folder), judged correct.**
- `run_status.go` `runStatusAt` stats the resolved folder once. A stat error or a non-directory sets the reader folder to `""`, so `readManifestLines`, `buildRow` and `listRunLocalFiles` all take their existing "no run folder" branch. The absent-not-false contract therefore matches `TestRunStatusReportsQueueRowsWithoutARunDirectory`.
- An auto-picked folder comes from `os.ReadDir` with `entry.IsDir()` and a `work-` prefix, so it always passes the stat (D-06 is right: one code path, no cost). A symlinked run folder is skipped by `ReadDir`'s `IsDir`, which is unchanged behavior from before this REQ.
- The `run:` line sits in the full form only. The `--watch` branch returns earlier, so the 20-line budget test stays valid.
- An absolute `--run` outside the repository shows the absolute path (via `displayPath`). That is fine.
- A permission error on stat also reads "not found" (D-07 wording). This is a nit (F6).

**F4 of REQ-690 (C3 while an integrator works), judged correct.**
- The new arm comes after C8 and before the old C3 arm, so C8 (finalization pending) still wins. The arm needs `claimed`, `HandbackPresent` non-nil and true, a phase in `integratorPhases`, and `MinutesSinceActivity` non-nil and under `quietActivityBoundary` (20 min).
- When `MinutesSinceActivity` is nil, the arm does not fire and the row falls through to the old C3 arm with `do-work run`. This is the safe pre-REQ behavior. There is no nil dereference: the `Sprintf` reads `*record.MinutesSinceActivity` only inside the arm, after the nil check.
- Negative minutes (activity stamped in the future through clock skew) count as recent and fire the arm. That is acceptable.
- Phase spelling is verified. `skills/do-work-board/tools/queue-kanban/durations.go:264-277` `phaseMilestonesOf` has labels `Integration`, `Review`, `Remediation`, `Re-review`. `activity_correlation.go:206` applies `strings.ToLower(milestone.label)`, which gives `integration`, `review`, `remediation`, `re-review`. These exactly match the four keys in `integratorPhases`.
- `Builder handback` (lowercased `builder handback`) is correctly excluded. That phase lies between the hand-back stamp and `integration_at`, and it is also the phase a dead orchestrator leaves behind, which is when `do-work run` is the right remedy.
- `renderReport` takes `next:` from the first row with a non-empty `next_argv`, so the integrating row drops out of `next:` with no render change.

## Findings

**Important:** None

**Minor:**
- **F1** The 20-minute boundary on the new arm is untested. In mutation M1 I dropped the `< quietActivityBoundary` condition, and all 12 `^TestRunStatus` tests still passed. Without the boundary, a REQ whose integrator died in the integration phase would say "an integrator is working on it" forever and never offer `do-work run`. Fix: add a second row to `TestRunStatusIntegratingHandbackDoesNotRecommendRun`, last activity 25 min ago in phase `integration` with a hand-back, and assert `next_argv == do-work run`. I built this as a scratch test and it kills M1. — impact-user-visible → recommend fix now
- **F2** `integratorPhases` copies the board's labels without a lock-step check. A Go import cannot work: queue-kanban is `package main` in a separate module (`github.com/knews2019/skill-do-work/queue-kanban`). A cheap shell check can. `_dev/tests/contracts/queue-kanban.sh:90-98` already pins `CalibrationGapStampFields` against `phaseMilestonesOf` in the same way. About 5 lines would do it: sed the `var integratorPhases` line, extract the `"x": true` keys, then `grep -qi 'label: "x"' durations.go` for each one. I tested a prototype in scratch. It passes on HEAD and fails when `re-review` is mutated to `rereview`. This file is outside the declared write set. — impact-negligible → recommend fix now
- **F3** Stale restatement. `skills/do-work/actions/fan-out-reference.md:163` says "A C3 row means 'brief the next integrator', never a `do-work run` to execute." After this REQ, a C3 row whose phase is an integrator phase with activity in the last 20 minutes means an integrator is already on it. A coordinator following that sentence literally (for example after context compaction) would brief a second integrator, which is the double-orchestrator failure this REQ targets. Fix: one clause, for example "a C3 row in an integrator phase with recent activity is already being integrated; any other C3 row means brief the next integrator". This file is outside the declared write set. — impact-rule-change → recommend fix now
- **F4** `skills/do-work/actions/status.md:67` (Step 4 Report) lists what to relay (table, remedy lines, Step 3 lines, `next:` last) but does not name the new `run:` line. An agent that follows the list literally can drop the "not found" note, and that note is the whole point of the missing-folder fix. Fix: one clause, "keep the `run:` line, which names the run folder read and says when it was not found". This file is in the write set. — impact-user-visible → recommend fix now
- **F5** The `--watch` form still never shows the run folder, so a mistyped `--run` under `/loop` stays silent there. I agree with the builder's discovered task: a one-word marker on the header line would fit the 20-line budget. — impact-user-visible → report only

**Nit:**
- **F6** `run_directory` now carries prose in a path field (accepted in D-03, and no reader exists in the repository). A stat permission error also reads "not found" (D-07). Mutations M2 (dropping the `!info.IsDir()` check) and M3 (keeping only `integration` in the set) both survive. All three are accepted trade-offs. — impact-negligible → report only
- **F7** `status.md:48` says "(phase integration or later)". In the board's order, "later" also includes `completed` and `release`, which the code does not match. Naming the four phases would be exact. — impact-negligible → report only

## Requirements Checklist

- [x] A missing `--run` folder: JSON `run_directory` carries the "not found; run fields absent" note, the text prints `run: <dir>`, and there is no refusal and no exit-code change. Delivered (`TestRunStatusMissingRunFolderIsNamedInTextAndJSON`).
- [x] The text names the run folder it read whenever there is one. Delivered (full form only, by design D-03).
- [x] A C3 row with recent activity in phase `integration` does not recommend `do-work run`, says an integrator is working on it, and has an empty `next_argv`. Delivered (`TestRunStatusIntegratingHandbackDoesNotRecommendRun`).
- [x] Other C3 rows are unchanged. Delivered: `TestRunStatusLandedHandbackIsClassC3` is untouched and passes.
- [x] One test per failure. `next_argv` is empty on the new arm, never a command. Delivered.

## Acceptance Testing

**Result: Pass** (stages covered: unit tests and mutation, local only)
- `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^TestRunStatus' ./internal/runstatus/` on main: exit 0, 12/12 PASS, including both new tests.
- Mutations ran on a module copy in the scratch directory. M4 (new arm disabled) is killed by the F4 test. M1 (boundary dropped), M2 (`IsDir` dropped) and M3 (later phases dropped) survive (F1, F6).
- Lock-step prototype (F2) and boundary test prototype (F1) both behave as described, in scratch only.

## Suggested Additional Testing

- Live acceptance: unassessed. Run `do-work status --run do-work/runs/typo` on this repository and confirm the `run:` line shows `(not found; run fields absent)`.
- Live acceptance: unassessed. During the next coordinated integration, run `do-work status` while the integrator is in review, and confirm the C3 row has no `do-work run` on its `next:` line.
- Edge case: run `--run` against a regular file and against an unreadable folder. Both should read "not found" with run fields absent.

## Scores (on the record, not the headline)

**Overall: 93%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All five acceptance items delivered |
| Code Quality | 92% | Small, readable, matches existing patterns; label copy without a check (F2) |
| Test Adequacy | 80% | RED-to-GREEN proven for both cases; the boundary is unpinned (F1) |
| Scope | 100% | Exactly the five declared files; all decisions recorded (D-06 to D-09) |
| Risk | Low | Next_argv stays empty or read-only; fallback is the pre-REQ behavior |
| Acceptance | Pass | Unit plus mutation, local |

## Follow-ups created

None (7 findings report only)

## Review

**Overall: 93%** | 2026-10-10T20:38:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. Both REQ-690 review gaps are closed with correct code: the missing `--run` folder fix and the integrator-phase C3 arm.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1 The 20-min boundary on the integrator C3 arm is unpinned: mutation M1 survives. A dead integrator would never get `do-work run`. Add a quiet-integration row asserting `do-work run` to the F4 test. — impact-user-visible → recommend fix now
- F2 `integratorPhases` has no lock-step check. A Go import is impossible (queue-kanban is a separate-module `package main`). A 5-line check in `_dev/tests/contracts/queue-kanban.sh` following the `CalibrationGapStampFields` precedent (lines 90-98) works; the prototype catches the `rereview` mutation. — impact-negligible → recommend fix now
- F3 `skills/do-work/actions/fan-out-reference.md:163` still says every C3 row means "brief the next integrator". Add the integrator-phase exception. — impact-rule-change → recommend fix now
- F4 `skills/do-work/actions/status.md:67` Step 4 relay list omits the new `run:` line. — impact-user-visible → recommend fix now
- F5 `--watch` still never shows a missing run folder (the builder's discovered task). — impact-user-visible → report only
- F6 (nit) Prose in the `run_directory` path field (D-03); permission errors read "not found" (D-07); `IsDir` and later-phase membership unpinned (M2, M3). — impact-negligible → report only
- F7 (nit) `status.md:48` says "phase integration or later", but the code matches four named phases. — impact-negligible → report only

**Acceptance:** Pass. Unit and mutation stages only: focused `^TestRunStatus` 12/12 exit 0; M4 killed, M1-M3 survive.
**Restatement sweep:** redefined the C3 remedy (integrator exception), the run-status text shape (new `run:` line) and the `run_directory` value (suffix). Results: `status.md:48` and `status-guide.md:26` updated; `fan-out-reference.md:163` stale (F3); `status.md:67` incomplete (F4); `skills/do-work/CHANGELOG.md:55` is history; no reader of `run_directory` outside `internal/runstatus` and `resultmodel`. No shipped prose claims the text lacks a run line.
**Anti-bloat:** 0 unnamed additions. `integratorPhases`, the `lastActivityPhase` fixture field and the two tests were all named in the brief. `runDirectoryDisplay` is a local that replaced the old `if` block. No new options, files or decorative tests.
**Suggested testing:** 3 items
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*

## Re-review (delta 5d4e136c..bf9ed3d6)

**Overall: 95%** | 2026-10-10T20:41:26Z | **Approve**. Cumulative range `d2991804..bf9ed3d6` (fix builder commit `3ea1e4ae`).

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

- **F1 resolved.** The REQ-822 sub-check in `TestRunStatusIntegratingHandbackDoesNotRecommendRun` asserts that a row in phase `integration` with activity 25 minutes ago gets C3 with `next_argv == [do-work run]`. I re-ran mutation M1 (boundary dropped) on a `git archive` copy of `bf9ed3d6`: that test fails, so the mutation is killed. Focused `^TestRunStatus` on HEAD `bf9ed3d6`: exit 0.
- **F3 resolved.** `fan-out-reference.md:163` now names the exception: integrator phase plus activity in the last 20 minutes means the REQ already has an integrator. The coordinator reads `--watch`, which shows the phase and the activity age but not the remedy text. The parenthetical gives exactly those two facts, so the rule can be checked from `--watch` alone.
- **F4 resolved.** `status.md:67` Step 4 names the `run:` line and its "not found" note.
- **F7 resolved.** `status.md:48` names the four phases, matching `integratorPhases`.
- **F2 accepted as report only.** The coordinator's reason holds. When a label drifts, the row falls back to the pre-REQ `do-work run` advice. That is a safe failure and not a destructive one, and the check would add a cross-module contract outside the write set. → report only
- **F5, F6** unchanged. → report only
- **Delta findings:** None. The delta touches 3 files (+15 -3), all edits that F1, F3, F4 and F7 asked for. There are no unnamed additions.
- **Restatement sweep (delta):** the C3 exception and the `run:` line now read the same in `status.md`, `status-guide.md` and `fan-out-reference.md`. No other restatement found.

**Follow-ups created:** None (3 findings report only: F2, F5, F6)
