## Review

**Overall: 97%** | 2026-10-07T14:42:26Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- M1: `skills/do-work-board/tools/queue-kanban/durations.go:43-44` — the `activityGapCeiling` doc comment still says "a gap between two consecutive lifecycle stamps longer than this is idle time", the exact-duration reading this REQ retired; a 2h00m40s gap is now kept. The header comment at lines 32-36 states the whole-minute rule correctly, so the risk is a reader of the constant alone. — impact-negligible → report only
- M2 (Nit): `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go:2376` — the test `TestJavaScriptBehaviorDetailStatesTheLargestIdleGap` and its probe label "detail largest idle gap" keep the old row name. The builder documented this as D-02 (archived REQ-632 records cite the name). — impact-negligible → report only

**Acceptance:** Pass — focused duration tests and the drawer JS probe pass at ecb86826; reverting the comparison to the exact Duration turns the two new cases RED.
**Suggested testing:** 2 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

---

## Review: REQ-637

**Approve** — Panel B now reads the stamp gap in whole minutes like the calibration log, and the drawer row has a name that says what it measures.
Route B | merge range 670428d1..ecb86826 (worktree dispatch)

### What's built
- `dayMedianExclusionReason` compares `int(largestStampGap.Minutes()) > int(activityGapCeiling.Minutes())`. This matches the CLI writer `calibration_row.go:52` (`int(gap.Minutes())`) and the backfill script `_dev/backfill-calibration-gap-column.py:81` (`seconds // 60`). All three readers and writers now floor to whole minutes.
- The drawer row is "Largest gap between events". The measurement in `activity_correlation.go` is unchanged. The board guide adds one sentence saying the open stretch to now is not counted.
- The `durations.go` header comment now names the phaseless case (one gap, excluded over 2h even when continuous) and the whole-minute agreement.
- Release, CHANGELOG and the lessons entry are finalization's work and were not reviewed here.

### Decisions / risks for you
- None. Builder decisions D-01 to D-04 and orchestrator D-05/D-06 are recorded and match the diff.

### Findings

**Important:** None.

**Minor:**
- M1: `durations.go:43-44` `activityGapCeiling` comment says "longer than this is idle time", which is the old exact-duration reading. Suggested wording: "a gap of more than this many whole minutes". — impact-negligible → report only

**Nit:**
- M2: test name `TestJavaScriptBehaviorDetailStatesTheLargestIdleGap` and probe label keep the old row name (D-02, deliberate). — impact-negligible → report only

### Requirements Checklist

- [x] R1 Whole-minute comparison in `dayMedianExclusionReason`, derived from `activityGapCeiling` (no restated 120) — delivered (`durations.go:363`)
- [x] R2 Boundary test written first and seen RED: 2h00m40s kept, 2h01m excluded — delivered (`durations_test.go:509-517`; the existing ceiling+1m cases at 506/508 pin the exclusion side). RED evidence in the hand-back; reviewer mutation confirms.
- [x] R3 `durations.go` header comment names the phaseless case; rule unchanged — delivered (lines 27-30)
- [x] R4 Drawer label changed, measurement unchanged, test asserting the old label updated — delivered (`board-detail.js:491`, `javascript_behavior_a_test.go:2387-2389`)
- [x] R5 Full queue-kanban suite with JS probes — delivered by the builder on the branch (pass, 0 failures, per hand-back); the gate and heavy lanes run at the merge by the orchestrator. Reviewer ran focused tests only, per instruction.
- [ ] R6 Release + CHANGELOG reasoning for the relabel — N/A (finalization). The hand-back's proposed CHANGELOG text states the reason.
- [ ] R7 Lessons entry + index refresh — N/A (finalization). Proposed lesson text is in the hand-back.
- [x] Constraint: no-phase-stamp exclusion rule unchanged — confirmed (only the comment changed)
- [x] Constraint: `activityGapCeiling` stays the single constant — confirmed
- [x] UR batch constraint: struct field names in `activity_correlation.go` and `board-controls.js` untouched — confirmed (only a comment changed in `activity_correlation.go`; `board-controls.js` not in the diff)

### Restatement Sweep

Redefined elements: (a) the drawer row name, (b) Panel B's gap comparison (exact Duration to whole minutes).

- (a) `git grep -i 'idle[- ]gap'` and `'drawer.{0,40}gap'` over the repo, excluding `do-work/` and `CHANGELOG.md`, at ecb86826. No remaining restatement of the drawer row as "Largest idle gap" or "idle-gap row". Remaining "idle gap" / "idle-gap" hits all name the Panel B rule or its verdict token (`dayMedianExclusionRule`, the `"idle-gap"` reason, `board-durations.js`, `board-timeline.js`, `generate.go`, `timeline.go`, the contract probe comment at `_dev/tests/contracts/queue-kanban.sh:89`), which is intentional (D-04). The `ai-reports/` hits are about axis idle time, unrelated.
- (b) Grepped `activityGapCeiling`, `exceeds 2h`, `over 2h`, `max_stamp_gap`. `estimate-reference.md:94,96,100` states "exceeds 2h" against a value already defined as whole minutes rounded down, so it agrees. `board-guide.md:41,57` says "over 2h", consistent. The one stale statement is the `activityGapCeiling` doc comment (M1). No other Go or JS code compares a stamp gap to the ceiling; `timeline.go` and the JS consume the Go verdict.

### Acceptance Testing

**Result: Pass**
- Scratch worktree at ecb86826 (detached, removed after).
- `go test -count=1 -run 'ImplementationSpan|DayMedian|Durations' .` — ok (3.6s).
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 -run TestJavaScriptBehaviorDetailStatesTheLargestIdleGap -v .` — PASS.
- `gofmt -l` on the four changed Go files — empty.
- Mutation 1: comparison reverted to `largestStampGap > activityGapCeiling` — `TestImplementationSpanVerdictReadsTheLargestStampGap` FAIL.
- Mutation 2: comparison changed to `math.Round(minutes) > ceiling minutes` (rounding instead of flooring) — same test FAIL, so the ceiling+59s case also pins floor-versus-round.
- `TestImplementationSpanAgreesWithTheDurationsAggregate` (card verdict equals Panel B verdict) stays green, so the card and Panel B still share one predicate.

### Suggested Additional Testing

- Open a claimed card's drawer in the served board and confirm the row reads "Largest gap between events" (browser probes were not run by the builder; no browser probe asserts this row).
- On the next estimator re-fit, confirm a REQ with a sub-minute-over-120 gap reads the same verdict in Panel B and in the re-fit.

### Scores (on the record — not the headline)

**Overall: 97%** (mean of 100, 95, 95, 100 = 97.5, rounded down; Risk None, Acceptance Pass)

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | R1-R5 delivered; R6-R7 finalization-owned, N/A |
| Code Quality | 95% | One-line fix derived from the constant; one stale constant comment (M1) |
| Test Adequacy | 95% | Sub-minute boundary cases on both phase and phaseless paths; both mutants caught; no cross-module test of CLI and board on one fixture (different Go modules, acceptable) |
| Scope | 100% | 8 files equal the declared Scope; D-03/D-06 out-of-write_set edits are documented restatement fixes |
| Risk | None | Verdict changes only for gaps in (120m, 121m); payload shape unchanged |
| Acceptance | Pass | Focused tests, JS probe, two mutations |

### Self-validation
- Checked the floor-versus-round edge with a second mutant, not only the exact-Duration revert.
- Checked the third reader of the gap unit (the Python backfill script) floors too.
- Did not run the full module suite or the gate, per the orchestrator's instruction; that evidence comes from the builder and the orchestrator's gate run.
- Builder's Discovered Task about `activity_correlation.go:18` and `estimate-reference.md` was resolved by orchestrator commit 7ee3ad48 (verified in the diff).

### Follow-ups created
None (2 findings report only)
