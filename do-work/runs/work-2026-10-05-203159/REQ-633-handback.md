# REQ-633 hand-back — Panel B and the calibration log exclude by largest stamp gap

- Branch: `worktree-agent-REQ-633-panel-b-largest-stamp-gap-exclusion`
- Commit: `43024f7e` (one commit on top of `9dacf899`)
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-633-panel-b-largest-stamp-gap-exclusion
- Nothing under `do-work/` written in either tree except this file. No release files touched. `do-work/calibration-log.tsv` not backfilled.

## File manifest

New:
- `skills/do-work/tools/do-work-cli/internal/requestmodel/calibration_row.go` — the shared formatter: `CalibrationLogHeader` (6 columns), `CalibrationGapStampFields`, `CalibrationLogCarriesGapColumn`, `FormatCalibrationRow`.
- `_dev/backfill-calibration-gap-column.py` — one-off idempotent backfill with `--repo-root` and `--dry-run`.

Modified:
- Board: `durations.go` (activityGapCeiling 2h replaces analysisOutlierCeiling; `largestLifecycleStampGap`; verdict `idle-gap`; `dayMedianExclusionRule` text; header comment), `generate.go` (`durations.exclusionRule` field + comments), `timeline.go` (comment), `web/board-durations.js` (summary, axis title, mark label, table cell read the shipped text), `web/board-timeline.js` (forecast sentence), `web/board-user-request-summary.js` (comment), `docs/board-guide.md`.
- Board tests: `durations_test.go`, `generate_test.go`, `timeline_test.go` (comments), `javascript_behavior_a_test.go`, `javascript_behavior_d_test.go`.
- Core docs: `actions/estimate-reference.md` (Calibration rule, column list, header-as-version, second-reader sentence, 188-of-190 note, stamps-only vs drawer sentence), `actions/work-reference.md` (phase stamps now feed `max_stamp_gap_minutes`).
- Core CLI: `requeststate/state_plan.go`, `requeststate/state_apply.go`, `finalization/finalization_discovery.go` (all three call `FormatCalibrationRow`).
- Core CLI tests: `requeststate/state_apply_test.go`, `finalization/finalization_recovery_test.go` (see D-04).
- `_dev/tests/contracts/queue-kanban.sh` (stamp-list pin, see D-02).

## Red-green evidence

RED (board, stub with `activityGapCeiling` and empty rule text, old verdict logic):
- `TestImplementationSpanVerdictReadsTheLargestStampGap`: "a phase gap one minute over the ceiling is excluded, though the span is only 3h01m: verdict "", want "idle-gap""; "no phase stamps, span one minute over the ceiling: verdict "", want "idle-gap""; "a 4h21m span whose largest stamp gap is 67 minutes is kept: verdict "paused", want """; rewritten-claim case: verdict "paused", want "idle-gap".
- `TestGeneratedRequestCarriesTheDoneCardImplementationSpan`: REQ-909 "paused", want ""; REQ-902 "paused", want "idle-gap"; `durations.exclusionRule = "", want "largest idle gap over 2h"`.
- `TestJavaScriptBehaviorDurationsHeadlineRollingMedianAndCadenceTicks`: summary still said "over four hours is an assumed pause".
- `TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions`: wanted "Idle-gap and reversed spans are excluded".
- Also `TestDurationDayMedianAppliesTheReadTimeOutlierRule`, `TestImplementationSpanOpensAtTheEarliestLifecycleStamp`, `TestImplementationSpanAgreesWithTheDurationsAggregate` (vacuity guard: never produced "idle-gap").

RED (CLI, builders not yet wired to the formatter):
- `TestCompleteCalibrationRowFollowsTheLogHeader`: new log got the 5-column header and no gap; six-column log got a five-column row; phaseless case had no `\t240`. (The five-column-log case passed before and after: it is the old-header guard.)
- `TestArchivedCalibrationVerifierReadsTheHeaderShape`: "six-column row under a six-column header rejected".
- `TestCalibrationAppendProofFollowsTheLogHeader`: six-column row → proves false (want true); five-column row under six-column header → proves true (want false).

RED (backfill): fixture test against an empty stub script failed "FAIL: first run" (log unchanged).

GREEN:
- `cd skills/do-work-board/tools/queue-kanban && QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...` → ok 65.5s (the live-archive calibration pins still hold).
- `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/requeststate/ ./internal/requestmodel/ ./internal/finalization/` → ok (6.7s / 0.3s / 45.5s package; the new finalization test adds about 0.13s to `finalization_recovery_test.go`, measured 18.5s before).
- Whole CLI module `go test -count=1 ./...` → 32 packages ok, exit 0.
- `gofmt -l` empty and `go vet ./...` clean in both modules.
- `bash _dev/tests/contracts/queue-kanban.sh` → passed; with `re_review_at` removed from the CLI list it fails with both lists printed (it catches a mismatch).
- `bash _dev/tests/shipped-package-reference-contract.sh` → PASS (new cross-package citation in estimate-reference.md).
- Backfill fixture test (scratchpad `req633/backfill-fixture-test.sh`): `--dry-run` writes nothing; the first run adds the column (121 for a phased REQ whose release_at would have made 359; 45 for a phaseless one; empty for a row with no archive match) and leaves wall_minutes as it was (999 in the fixture); the second run changes nothing → PASS.
- `grep -rn '"paused"' queue-kanban/*.go queue-kanban/web/*.js` → no matches.

## P-A-U

- [x] **[PLAN]** One gap definition in two places: board `largestLifecycleStampGap` over `phaseMilestonesOf` minus release_at; CLI `largestCalibrationStampGap` over `CalibrationGapStampFields` plus completed_at. Sorted by time, largest consecutive difference, unparseable stamps skipped. A shared CLI formatter decides the row shape from the first line of the post-append log bytes, which all three builders already hold. The rule's words ship as `durations.exclusionRule`.
- [x] **[APPLY]** As planned. Write set = Scope plus `calibration_row.go` (the named formatter file), `finalization_recovery_test.go` (D-04) and `_dev/tests/contracts/queue-kanban.sh` (D-02). `state_plan_test.go` was not needed: the header cases run through the command seam in `state_apply_test.go`.
- [x] **[UNIFY]** `git diff --stat 9dacf899..43024f7e`: 22 files, 548 insertions and 123 deletions. Each file was reviewed. gofmt and vet are clean, and there are no debug artifacts. The only "paus" strings in the rendered page are REQ-632's own title.

## Decisions

- **D-01 (DECIDE & STATE)** — The verifier and finalization now parse the estimate with `strconv.Atoi` through the shared formatter. Before, they compared the raw string. Value: one formatter, and a padded estimate like `040` can no longer make the appender (which writes `40`) and the verifier disagree. Risk: none seen. The appender already refused a non-integer estimate.
- **D-02 (DECIDE & STATE)** — The stamp-list pin is a `_dev/tests/contracts/queue-kanban.sh` grep and not a Go test. Value: it is export-ignored, so no shipped test reads a sibling module's source or skips silently in a consumer (trap shipped-module-test-self-containment). Risk: it runs only in the maintainer gate. The backfill script holds a third copy of the list, which is acceptable for a one-off.
- **D-03 (DECIDE & STATE)** — The Panel B text is "largest idle gap over 2h, or a negative span from a broken stamp" in the summary and "excluded from day median — largest idle gap over 2h" in the table and mark labels. The axis title stays short: "… · idle-gap and broken spans excluded" (2 characters longer than before; the layout suite passed). The timeline forecast says "Idle-gap and reversed spans are excluded from both medians." and the UR group reads "1 idle-gap excluded". The Go field `ExclusionReason` keeps its name; only its values changed.
- **D-04 (ESCALATE, low)** — I added `TestCalibrationAppendProofFollowsTheLogHeader` to `finalization/finalization_recovery_test.go`, which is outside the REQ's write set. The brief's TDD list asks for "finalization proof accepts the new row", and lesson alternate-writer-contract-drift asks for seam coverage at every writer. No existing in-scope file could hold that test. Value: the third builder is pinned under both header shapes. Risk: one test file outside scope. It is easy to revert if you disagree.
- **D-05 (DECIDE & STATE)** — Backfill rows with no single archive match get an empty value and are reported, never guessed. On this repo, 0 rows are unresolved.

## Panel B exclusion counts (generated against the worktree root, same tree, both binaries)

- Before (`9dacf899` binary): 539 samples, 36 excluded `paused`, 503 kept.
- After (`43024f7e` binary): 539 samples, 51 excluded `idle-gap`, 488 kept. No REQ was kept that the old rule excluded. 15 REQs are newly excluded: REQ-043, 248, 346, 378, 412, 413, 414, 417, 418, 419, 420, 457, 543, 552, 594. These numbers match the exploration.
- Rendered (headless Chrome, `#durations`, 30-day window): "Panel B excludes 32 spans from its medians (largest idle gap over 2h, or a negative span from a broken stamp); panel A still plots them." 32 marks are labelled "excluded from day medians (largest idle gap over 2h)" and 32 table cells say "excluded from day median — largest idle gap over 2h".

## Backfill dry run on the main tree

`python3 _dev/backfill-calibration-gap-column.py --repo-root /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 --dry-run` → 309 rows, 309 filled, 43 with a gap over 120 min, 0 unresolved. The file's md5 was the same before and after. The orchestrator should run it without `--dry-run` in the main tree after finalization.

## Discovered Tasks

- impact-minor: Before REQ-633, `calibrationAppendProves` could not prove the creation of a brand-new log, because `after[len(before):]` included the header while `want` was the row alone. The same is still true for a new six-column log. It does not matter in a repo that already has the log, and no test covers creation. → report only

## Lessons read

REQ-374 (straddling pair: every boundary case is derived from `activityGapCeiling` ± 1 minute), REQ-219 (ship the verdict and text, never the inputs: the client gets `exclusionRule` text, no number), alternate-writer-contract-drift (all three writers go through one formatter, with seam tests at each). The primes and crew members listed in the brief were also read.

## Integration seams

None expected. Release files are untouched: VERSION/CHANGELOG belong to finalization. Shipped packages changed are do-work-board (Go, web, docs) and do-work (actions, CLI).
