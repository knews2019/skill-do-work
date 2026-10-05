## Review

**Overall: 94%** | 2026-10-05T21:49:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 93% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- Rounding at the boundary differs between the two readers: the board excludes when the exact gap is over 2h (`durations.go:284`, `largestStampGap > activityGapCeiling`), while the CLI writes the gap floored to whole minutes (`calibration_row.go:51`) and `estimate-reference.md:94-96` gives no comparison for the column. A re-fit using `max_stamp_gap_minutes > 120` keeps a gap in [120m01s, 121m) that Panel B excludes. On this repo all 309 backfilled rows agree with the board's verdict, so it is theoretical today. Fix with one doc clause (`> 120` on the floored column means "121 minutes or more"), or compare in seconds. — impact-negligible → report only
- `estimate-reference.md` tells the next re-fit to apply the gap rule, but by design an existing five-column log never gets the column: consumer repos keep the old header and the backfill is in `_dev/`, so it does not ship. The doc does not say what a re-fit does with five-column rows (for example, fall back to the span rule or skip them). The header-decides-shape choice is documented in the doc and the formatter comment, but it is not among the hand-back's D-01..D-05. — impact-rule-change → report only
- The REQ's Constraint asked for the stamps-only vs stamps-plus-commits sentence "once in the doc". It appears in two shipped docs, `skills/do-work-board/docs/board-guide.md:57` and `skills/do-work/actions/estimate-reference.md:96`, plus the `durations.go` header comment. Keep one doc copy and point to it from the other. — impact-negligible → report only
- Operational trap that nothing documents: `calibrationAppendProves` (`finalization_discovery.go:470-479`) requires the working-tree log to start with the HEAD bytes. The backfill rewrites every line, so if it runs and is not committed on its own before the next REQ completes, that REQ's finalization cannot prove its calibration append and refuses. The script's docstring should say: run it, then commit it alone before the next completion. — impact-user-visible → report only
- Stale restatement: the comment at `skills/do-work-board/tools/queue-kanban/generate_test.go:2194` still says the sweep covers "everything under the four-hour read-time ceiling that admits a sample at all". Under the gap rule, spans over 4h are admitted, for example REQ-909 at 261 minutes. Behaviour is unaffected because Panel B clamps medians at 45 minutes. — impact-negligible → report only
- Nit: the comment at `durations_browser_probe_test.go:797` still calls its fixtures "paused spans". — impact-negligible → report only

**Acceptance:** Pass — I built the board at 0951df36 and ran generate against the repo. Result: 540 samples, 51 `idle-gap`, 489 kept, `durations.exclusionRule` = "largest idle gap over 2h", and no `"paused"` in board-data.js. The backfill dry run wrote nothing (md5 unchanged). On a scratch copy the backfill is idempotent and leaves columns 1-5 byte-identical. Its gaps agree with the board's verdict on 309 of 309 rows.
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

### Review detail (orchestrated report)

**Approve** — The gap rule, the rename and the new log column all behave as specified. The two gap definitions agree on real data. The six findings are minor documentation and comment gaps.
Route B | 43024f7e merged as 0951df36

#### Check-by-check
1. Gap definitions agree:
   - Stamp set: board `phaseMilestonesOf` minus release_at; CLI `CalibrationGapStampFields` plus completed_at. The contract check pins the two lists and passes.
   - release_at is excluded on both sides.
   - Both sides sort by instant and take the largest consecutive difference.
   - Fallback: a phaseless REQ falls back to the whole claimed→completed span on both sides.
   - Both use the same parser layouts.
   - The only difference is the boundary rounding (Minor 1).
   - Cross-check: the backfilled column agrees with board verdicts on all 309 log rows.
2. One formatter: `requestmodel.FormatCalibrationRow` is called by the appender (`state_plan.go:408,425`), the verifier (`state_apply.go:497`) and the finalization proof (`finalization_discovery.go:483`).
   - Header-decides-shape is tested for a new log, a 6-column log and a 5-column log through the command seam (`TestCompleteCalibrationRowFollowsTheLogHeader`, which includes the verifier), plus `TestArchivedCalibrationVerifierReadsTheHeaderShape`.
   - The prefix proof is tested for an existing 6-column and 5-column log (`TestCalibrationAppendProofFollowsTheLogHeader`).
   - The new-log proof gap predates this REQ and is reported in the hand-back.
3. Restatement Sweep: I grepped skills/, _dev/primes, _dev/tests and the docs for paused / assumed pause / four-hour / 4h / analysisOutlierCeiling / implementationSpanReason / excludedReason / calibration-log columns.
   - Updated: Go, all three web readers, board-guide, estimate-reference (rule, columns, second reader, 188-of-190 note) and work-reference:263.
   - Correctly left alone: the remaining "paused" hits are unrelated (estimate scope, hold folder, the four-hour heavy-lane reuse) or quote history (the REQ-374 lesson, the timeline_test Red-Green quote).
   - Stale: two test comments (Minor 5 and the Nit).
   - Board `verify.go:1616` reads at least 4 columns, so it is safe with 6.
4. The Constraint sentence is present, but twice (Minor 3).
5. Backfill: idempotent through the header check, wall_minutes untouched (verified), the dry run writes nothing, and it uses an atomic replace.
6. Tests: each named test states its pin. The straddling pair is derived from `activityGapCeiling` (−1m / at / +1m). Also covered: 4h21m with a 67-minute gap is kept, the phaseless fallback, release_at ignored, and the rewritten-claim sort. The vacuity guard covers idle-gap. The builder recorded RED for all of them.
7. Release (changelog and version): left to finalization, as the hand-back states.

Focused tests rerun: requestmodel ok, requeststate ok (6.4s), the finalization proof test ok, board Duration/ImplementationSpan/Timeline/JS subsets with JavaScript probes ok (7.3s), and `_dev/tests/contracts/queue-kanban.sh` passed.

#### Suggested Additional Testing
- After finalization, run the backfill without `--dry-run` and commit it alone. Then complete one REQ and confirm the new row has 6 columns and finalization proves the append.
- Open the Durations page and check the axis title width and the mark tooltip text at phone width.
- Optional: one board test that pins the agreement between the board verdict and a floored CLI gap at 120m30s, if Minor 1 is resolved in favour of exact comparison.

#### Scores
| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Requirements 1-6 delivered; requirement 7 (release) is owned by finalization |
| Code Quality | 93% | One formatter, the rule ships as text, clear comments; the rounding contract is not stated |
| Test Adequacy | 90% | Straddling pair, seam tests at all three writers; no cross-reader test at the boundary |
| Scope | 95% | Three added files recorded in Scope (D-02, D-04) |
| Risk | Low | Fail-closed finalization if the backfill is left uncommitted |
| Acceptance | Pass | Generated board shows 51 excluded and the shipped rule text |
