---
id: REQ-637
title: '[impact-rule-change] Board Panel B compares the stamp gap in whole minutes like the calibration log, and the drawer row is relabelled to mean the gap between events'
status: completed
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-07T13:59:57Z
created_at: 2026-10-06T22:45:18Z
user_request: UR-137
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-635, REQ-636, REQ-638]
batch: review-0305-69-findings
required_lessons: [_dev/primes/lessons-releases.md]
write_set: [skills/do-work-board/tools/queue-kanban/durations.go, skills/do-work-board/tools/queue-kanban/durations_test.go, skills/do-work-board/tools/queue-kanban/web/board-detail.js, skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go, skills/do-work-board/docs/board-guide.md, skills/do-work-board/tools/queue-kanban/activity_correlation.go, skills/do-work-board/tools/queue-kanban/web/board-cards.js, skills/do-work/actions/estimate-reference.md]
claimed_at: 2026-10-07T13:58:57Z
dispatch_at: 2026-10-07T13:16:27Z
builder_handback_at: 2026-10-07T13:19:58Z
integration_at: 2026-10-07T14:30:02Z
review_at: 2026-10-07T14:42:26Z
kb_status: pending
commit: 9882b505c89811f41723a3c53bbeb90f6444076d
heavy_verified_at: 2026-10-07T14:59:19Z
heavy_verified_revision: 9882b505c89811f41723a3c53bbeb90f6444076d
completed_at: 2026-10-07T14:59:39Z
release_at: 2026-10-07T14:59:39Z
---
# Board Panel B Compares the Stamp Gap in Whole Minutes Like the Calibration Log, and the Drawer Row Is Relabelled to Mean the Gap Between Events

## What
Two small changes to how the queue-kanban board reads activity gaps. Review finding F4: `dayMedianExclusionReason` (`skills/do-work-board/tools/queue-kanban/durations.go:354-362`) compares the exact `time.Duration` against the 2h ceiling, while the calibration log column `max_stamp_gap_minutes` is written in whole minutes rounded down (`skills/do-work/tools/do-work-cli/internal/requestmodel/calibration_row.go:52`). For a 2h00m40s gap the board excludes the REQ and a re-fit reading the log with `> 120` keeps it. Make the board truncate to whole minutes before comparing. Review finding F6 (optional, accepted as a relabel): the drawer row "Largest idle gap" (`web/board-detail.js:491`) measures only between past events and ignores the open gap from the last event to now on a claimed card. Rename the row so it clearly means the gap between events.

## Why
`skills/do-work/actions/estimate-reference.md:94,96` defines the logged value as whole minutes rounded down and names the board's Panel B as "its second reader, not a second definition", so the two readers must agree at the boundary. For F6, relabelling rather than including the open gap keeps one meaning on claimed and done cards, avoids a third gap definition next to the Panel B rule (`durations.go:44` uses the same words "largest idle gap") and the calibration column, and the card already shows "last activity Nh ago" for the open gap. That reasoning goes into the CHANGELOG, as the review asks.

## Detailed Requirements
1. F4: in `dayMedianExclusionReason`, compare `int(largestStampGap.Minutes()) > int(activityGapCeiling.Minutes())`, or an equivalent whole-minute truncation, so a gap is excluded only when its whole-minute value exceeds 120. The existing boundary tests at `durations_test.go:504-506` use exact-minute gaps and must stay green.
2. F4 boundary test, written first and confirmed failing: a largest stamp gap of 2h00m40s is kept (not `idle-gap`), and 2h01m00s is excluded.
3. Fix the comment at `durations.go:27`: "a long continuous session still counts" is wrong for a REQ with no phase stamps, whose single `claimed_at` to `completed_at` gap over 2h is excluded by the documented rule. Say so. The rule itself does not change (the review rejected changing it).
4. F6: change the label at `web/board-detail.js:491` from "Largest idle gap" to a label that says it is the gap between recorded events, for example "Largest gap between events". The measurement in `activity_correlation.go:247-256` does not change. Update any test or fixture that asserts the old label text.
5. Run the full Go test suite for the queue-kanban module and whatever browser or JS check covers `board-detail.js` in this repo.
6. Release per `_dev/primes/prime-releases.md`. The CHANGELOG entry states why the row was relabelled rather than widened to the open gap.
7. Lessons: one entry in `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` under the existing `paired-predicate-drift` family (two readers of one documented rule compared different units at the boundary) and the matching token refresh in `do-work/lessons-index.md`, in the same commit.

## Constraints
- Do NOT change the no-phase-stamp exclusion rule itself; `estimate-reference.md` documents it and the review rejected changing it.
- `activityGapCeiling` stays the single constant the rule is stated from; the whole-minute comparison derives from it rather than restating 120.
- Batch constraint: each REQ in this batch is its own release and its own commit.

## Dependencies
None. Independent of REQ-635, REQ-636 and REQ-638; REQ-636 edits `activity_correlation.go` and `verify.go`, which this REQ does not touch.

## Builder Guidance
Certainty is high; both changes are a few lines. Builder latitude: the exact relabel wording and whether the CHANGELOG reasoning for F6 is one sentence or two.

## Red-Green Proof
**RED prompt/case:** A completed REQ with `claimed_at` 2026-08-03T01:00:00Z, one phase stamp at 03:00:40Z and `completed_at` at 03:30:00Z (largest stamp gap 2h00m40s), passed through the Panel B exclusion. Separately, open a claimed card's drawer with two past events.
**Why RED now:** Panel B marks the REQ `idle-gap` while a re-fit reading the calibration log value 120 keeps it. The drawer row says "Largest idle gap" although the open gap to now is not included.
**GREEN when:** Panel B keeps the 2h00m40s REQ and still excludes a 2h01m00s one; the drawer row label states it measures the gap between events; the comment at `durations.go:27` mentions the no-phase-stamp case.
**Validation:** Inferred during capture. Confirmed by reading `durations.go:354-362`, `calibration_row.go:52`, `estimate-reference.md:94,96` and `board-detail.js:491` during triage.

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (7451 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matched: family `paired-predicate-drift` is exactly F4's failure shape.
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`). Matched: it governs finding-detail and view text, which the drawer relabel changes.

## Full Context
See `do-work/user-requests/UR-137/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** F4: truncate both sides to whole minutes in `dayMedianExclusionReason`, deriving 120 from `activityGapCeiling` (no restated number). RED cases: the REQ's 2h00m40s phase-gap fixture and a phaseless ceiling+59s span; the existing ceiling+1m case pins the exclusion side. Fix the header comment for the phaseless case. F6: relabel the drawer row in `board-detail.js` to "Largest gap between events", update the JS probe's expected label and the guide paragraph that names the row. (Ticked by the orchestrator from the builder hand-back.)
- [x] **[APPLY]:** Done as planned: tests first (both RED confirmed), then `durations.go` and `board-detail.js`; guide paragraph in its own commit a49fde43. Orchestrator commit 7ee3ad48 renamed the row in three restating comments/docs (D-06). (Ticked by the orchestrator from the hand-back and its own commit.)
- [x] **[UNIFY]:** `git diff --stat 670428d1..ecb86826`: 8 files, all in Scope. `gofmt -l` empty, `go vet ./...` clean (builder); full module suite with JS probes on passed on the branch; no debug artifacts. (Ticked by the orchestrator from the hand-back, cross-checked against the merge range.)
*Source: review findings F4 and F6 from the pasted consumer code review of 0.305.58 to 0.305.69, triaged with `do-work validate-feedback` on 2026-10-06 and accepted (F6 as a relabel).*

## Triage

**Route: B** - Medium

**Reasoning:** The REQ names both defects, their lines and the intended change (whole-minute comparison derived from `activityGapCeiling`, one label string), so no Plan agent is needed. The relabel has reach the REQ does not list: shipped docs and comments that restate the drawer row by name, so exploration finds every restatement before the Scope is fixed.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The build already existed when this run claimed the REQ (a parallel builder worked it on its branch; hand-back `do-work/runs/work-2026-10-07-135835/REQ-637-handback.md`), so exploration was taken from that hand-back and its real file list, checked against the source:

- **Required lessons consult:** `_dev/primes/lessons-releases.md` read whole (666 tokens, `slugged: full`; neither family touches this change). `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` stays dropped for budget (now 7451 tokens, `slugged: partial`, so no targeted form), but its `paired-predicate-drift` bullets were read by the builder and the orchestrator because they are F4's failure shape. `_dev/primes/lessons-kanban-board.md` stays dropped for budget. No listed file was missing.
- **Where the change goes:** `durations.go` `dayMedianExclusionReason` is the only Panel B comparison against `activityGapCeiling`; `calibration_row.go:52` writes `max_stamp_gap_minutes` as `int(gap.Minutes())`, so whole-minute truncation of both sides matches it. The drawer row label is one string in `web/board-detail.js`, asserted by `TestJavaScriptBehaviorDetailStatesTheLargestIdleGap` in `javascript_behavior_a_test.go`.
- **Restatements of the drawer row by name (the relabel's reach):** `skills/do-work-board/docs/board-guide.md` (bold row name, user-facing), the `activity_correlation.go` header comment ("states its largest idle gap"), the `web/board-cards.js` span comment ("the drawer's largest-idle-gap row") and `skills/do-work/actions/estimate-reference.md` § Calibration ("the board's drawer idle-gap row"). The Panel B rule text `dayMedianExclusionRule` ("largest idle gap over 2h") names the Panel B rule, not the drawer row, and stays (builder D-04).
- **Overlap with REQ-636 (0.305.71, landed first):** it rewrote the `activity_correlation.go` header comment right next to the stale line and edited `board-guide.md` one paragraph above this REQ's paragraph, so both files meet at the merge.

*Generated from the builder hand-back by the work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — whole-minute comparison derived from activityGapCeiling; header comment names the phaseless case and the whole-minute agreement
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — the 2h00m40s and ceiling-plus-59s boundary cases
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modify) — drawer row label only
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modify) — the expected label in the drawer gap-row probe
- `skills/do-work-board/docs/board-guide.md` (modify) — the guide names the relabelled row (builder D-03, orchestrator D-05)
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (modify) — header comment names the row by its new label, no code (D-06)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — one comment names the row by its new label, no code (D-06)
- `skills/do-work/actions/estimate-reference.md` (modify) — § Calibration names the drawer row by its new label (D-06)

**Files I will NOT touch:** the Panel B rule text in durations.go (dayMedianExclusionRule), the gap measurement in activity_correlation.go, the no-phase-stamp exclusion rule, calibration_row.go. Written by finalization in the main tree, not on the branch: every release path, skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md and do-work/lessons-index.md.

**Acceptance criteria (restated from REQ):**
- [ ] Panel B compares whole minutes rounded down, derived from activityGapCeiling: a 2h00m40s largest stamp gap is kept and 2h01m00s is excluded; the existing exact-minute boundary tests stay green
- [ ] The boundary test was written first and seen failing
- [ ] The durations.go header comment says a REQ with no phase stamps has one claimed_at to completed_at gap, excluded over 2h; the rule itself is unchanged
- [ ] The drawer row label says it measures the gap between recorded events; the measurement is unchanged; tests asserting the old label are updated
- [ ] The full queue-kanban Go suite and the JS/browser checks covering board-detail.js pass
- [ ] Release with a CHANGELOG entry that says why the row was relabelled rather than widened to the open gap
- [ ] One paired-predicate-drift lesson entry with its lessons-index token refresh, in the same commit

## Pre-Flight

**Git:** ✓ Integration tip 186202f5 on `main` (this REQ's claim commit, after REQ-636's finalization e29c9fd5); the only dirt is this REQ's own trail (working REQ, run directory) and `do-work/working/baseline.json`, which this pre-flight rewrites
**Tests baseline:** ✓ focused board tests green (`do-work/runs/work-2026-10-07-135835/REQ-637-probe.sh`, launched by advance)
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 186202f5 from the detached checkout `.git/work-run-2026-10-07-135835/drain-head` (14:27:34Z to 14:29:26Z, 112s; queue-kanban-fast-tests reused from 14:03:13Z, do-work-cli-fast-tests executed, slowest file finalization_recovery_test.go 29.34s < 30s); green-gate record satisfied. Two earlier runs at the same revision (14:00:51Z, 14:15:36Z) were red only on the per-file 30s budget of do-work-cli files this REQ does not touch (finalization_req499_test.go 53.83s/32.01s, defer_gate_test.go 40.80s/32.87s, finalization_recovery_test.go 34.93s/37.67s) at load averages 37 to 64 from other processes; the third run waited for a load under 8
**Dependencies:** ✓ Go toolchain and node present; no new module dependency

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified)
- `skills/do-work/actions/estimate-reference.md` (modified)

**What was done:** Panel B's `dayMedianExclusionReason` now compares `int(largestStampGap.Minutes()) > int(activityGapCeiling.Minutes())`, so a 2h00m40s largest stamp gap is kept, the same verdict a re-fit reading the logged `max_stamp_gap_minutes` 120 reaches, and 2h01m is still excluded. The durations.go header comment now says a REQ with no phase stamps has one gap, its whole claimed_at → completed_at span, so it is excluded over 2h even when the work was continuous, and that this reader compares whole minutes like the log. The drawer row "Largest idle gap" is now "Largest gap between events"; the measurement is unchanged. The board guide names the new row and says the open stretch to now is not counted. Three restatements of the old row name (activity_correlation.go header comment, board-cards.js span comment, estimate-reference.md § Calibration) now use the new label (D-06). After review, the `activityGapCeiling` comment says the ceiling is compared in whole minutes rounded down (review M1, D-07). Merge range 670428d1..9882b505 (builder commits 0c948168 and a49fde43, orchestrator commits 7ee3ad48 and ffc166e7, first merge ecb86826 with the activity_correlation.go header-comment conflict against REQ-636 resolved by keeping both changes, second merge 9882b505).

## Qualification

**Diff range:** 670428d1..ecb86826 (builder commits 0c948168 and a49fde43, orchestrator commit 7ee3ad48, merge ecb86826)
**Gate records:** qualify satisfied; scope-drift satisfied (the eight changed files equal the declared Scope).
**Warnings judged:** none raised.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. 1: `dayMedianExclusionReason` compares `int(largestStampGap.Minutes()) > int(activityGapCeiling.Minutes())`, derived from the one constant, matching `calibration_row.go`'s `int(gap.Minutes())`; the exact-minute boundary cases stay in the table. 2: the 2h00m40s phase-gap case (the REQ's literal fixture) and the ceiling+59s phaseless case were seen RED by the builder; ceiling+1m stays excluded. 3: the header comment names the phaseless one-gap case and the rule is unchanged. 4: one label string in `board-detail.js`, the measurement untouched, the JS probe's expected label updated. 5: the board module suite with JS probes passed on the branch; the gate and heavy lanes run at the merge. 6 and 7 are finalization's (release, lesson, index refresh).
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back and the orchestrator's own commit; APPLY cross-checked against git diff --stat 670428d1..ecb86826 (8 files, all in Scope).
**After review:** the post-review fix ffc166e7 (one durations.go comment, no code) was merged with the same pre as 9882b505; qualify's record stays on the first range, and the cumulative range 670428d1..9882b505 has the same 8 files, all in Scope.
**Live data flow:** `measureImplementationSpan` calls `dayMedianExclusionReason` for every Panel B span, and `board-detail.js` renders the row from the payload's `largestActivityGap`; no payload shape changed.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 9882b505 (detached checkout `.git/work-run-2026-10-07-135835/drain-head`)
**Result:** ✓ All passing — exit 0 on the first run at this revision, 14:54:11Z to 14:56:31Z, gate wall 140s; stage queue-kanban-fast-tests executed (423 tests, wall 45s, slowest file strict_behavior_regression_test.go 21.14s < 30s); stage do-work-cli-fast-tests executed (876 tests, wall 65s, slowest file internal/finalization/finalization_recovery_test.go 23.55s < 30s). Green-gate record satisfied by advance. The review fix ffc166e7 (a comment) was merged before this gate ran, so no gate ran at the first merge ecb86826.

**Focused tests:** `do-work/runs/work-2026-10-07-135835/REQ-637-probe.sh` (Panel B span verdict tests, and the drawer gap-row JS probe with QUEUE_KANBAN_JAVASCRIPT_PROBES=on) → exit 0, advance probe record satisfied. The builder ran the full queue-kanban module with JS probes on at a49fde43: pass, 0 failures.

**Red-green validation:** traced to `## Red-Green Proof`; RED taken by the builder before the production change, GREEN at 0c948168 and again at the merge:
- TestImplementationSpanVerdictReadsTheLargestStampGap, case "a 2h00m40s phase gap is kept" (the captured fixture: claimed 01:00:00Z, dispatch 03:00:40Z, completed 03:30:00Z): ✗ verdict "idle-gap" → ✓ kept
- the same test, case "no phase stamps, span 59 seconds over the ceiling": ✗ verdict "idle-gap" → ✓ kept; the existing case "a phase gap one minute over the ceiling" (2h01m) stays excluded
- TestJavaScriptBehaviorDetailStatesTheLargestIdleGap: ✗ label "Largest idle gap" → ✓ "Largest gap between events"
- The reviewer reverted the comparison to the exact Duration at ecb86826 and both new cases went RED again.

**New tests added:**
- two table cases in `TestImplementationSpanVerdictReadsTheLargestStampGap` (`durations_test.go`)

**Existing tests updated (cross-REQ impact):**
- javascript_behavior_a_test.go (REQ-632): the expected drawer row label — intentional; the test name is kept (builder D-02)

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 670428d1..9882b505
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — files changed under skills/do-work-board/tools/queue-kanban
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

*Verified by work action*

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
- M1 `skills/do-work-board/tools/queue-kanban/durations.go:43-44` — the `activityGapCeiling` doc comment still says "a gap between two consecutive lifecycle stamps longer than this is idle time", the exact-duration reading this REQ retired; a 2h00m40s gap is now kept. The header comment at lines 32-36 states the whole-minute rule correctly, so the risk is a reader of the constant alone. — impact-negligible → report only (fixed before finalization in ffc166e7, merged as 9882b505)
- M2 (Nit): `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go:2376` — the test `TestJavaScriptBehaviorDetailStatesTheLargestIdleGap` and its probe label "detail largest idle gap" keep the old row name. The builder documented this as D-02 (archived REQ-632 records cite the name). — impact-negligible → report only

**Acceptance:** Pass — focused duration tests and the drawer JS probe pass at ecb86826; reverting the comparison to the exact Duration turns the two new cases RED.
**Suggested testing:** 2 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*

Full reviewer report: `do-work/runs/work-2026-10-07-135835/REQ-637-review.md`.

## Decisions

Builder decisions, copied from the hand-back (`do-work/runs/work-2026-10-07-135835/REQ-637-handback.md`).

- D-01 (DECIDE & STATE): the label is "Largest gap between events", the REQ's example. It is short, and "between events" says the open gap to now is not included.
- D-02 (DECIDE & STATE): the test function name `TestJavaScriptBehaviorDetailStatesTheLargestIdleGap` is kept. Archived REQ-632 records cite it by name, and renaming gains little.
- D-03 (DECIDE & STATE): the builder edited `skills/do-work-board/docs/board-guide.md`, outside the captured write_set, in its own commit a49fde43, because the guide names the row in bold and would be stale otherwise.
- D-04 (DECIDE & STATE): the Panel B rule text `dayMedianExclusionRule` ("largest idle gap over 2h") is unchanged. It names the Panel B rule, not the drawer row, and REQ-633 chose it on purpose.

Orchestrator decisions:
- D-05 (DECIDE & STATE, scope extension): a49fde43 is kept. The user-facing guide is the place a reader learns what the drawer row means, so a relabel that left it stale would ship a wrong doc. Scope and write_set gained `skills/do-work-board/docs/board-guide.md` and `javascript_behavior_a_test.go` (the REQ's Requirement 4 already named "any test that asserts the old label"). The guide paragraph merged cleanly beside REQ-636's edit of the row above it; both changes are in the merged guide.
- D-06 (DECIDE & STATE, scope extension): the three remaining restatements of the old row name were renamed in 7ee3ad48 before the first merge: the `activity_correlation.go` header comment, the `web/board-cards.js` span comment, and `skills/do-work/actions/estimate-reference.md` § Calibration. A stale restatement is a review finding anyway, and the coordinator asked for it. The builder named two of them; the exploration sweep found board-cards.js. The activity_correlation.go comment conflicted with REQ-636's rewrite of the same paragraph; the merge commit ecb86826 keeps REQ-636's owned-tip wording and this REQ's label. Scope and write_set gained all three files.
- D-07 (DECIDE & STATE, scope extension inside a declared file): review finding M1 (the `activityGapCeiling` comment still described the exact-duration comparison) was fixed in ffc166e7 and re-merged with the same pre as 9882b505 before the gate and heavy lanes ran, so they ran once, at the final merge.
- D-08 (DECIDE & STATE): dispatch_at (13:16:27Z, the branch creation in the reflog) and builder_handback_at (13:19:58Z, the builder commits) come from the branch, because the coordinator dispatched the builder before this run claimed the REQ at 13:58:57Z. They therefore precede claimed_at and the estimate. No builder-work timing event was recorded, because the recorder times from a start instant to now and would charge this run's work to the builder.
- D-09 (DECIDE & STATE): review finding M2 (the JS test keeps the old name) stays report only; it is builder D-02.

## Discovered Tasks

From the builder's hand-back and the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- `strict_behavior_regression_test.go` took 34.27s in the builder's JS-probe run of the board module, over the 30s per-file budget, with sibling sessions loading the machine; it took 21.14s inside the gate at 9882b505. — impact-negligible → report only
- The JS test `TestJavaScriptBehaviorDetailStatesTheLargestIdleGap` keeps the old row name (builder D-02, review M2). — impact-negligible → report only

## Lessons Learned

**What worked:** Deriving the whole-minute bound from `activityGapCeiling` with `int(…Minutes())`, the same expression the calibration-log writer uses, so both readers share one truncation and no number is restated. A sub-minute fixture (2h00m40s) was the only kind of test that could see the defect; the reviewer's mutation confirmed it goes RED.
**What didn't:** The first pass of the relabel left the old row name in four places outside `board-detail.js` (the guide, two code comments, the calibration doc), and the constant's own comment still described the exact comparison; all were caught by sweeps (builder, exploration, review M1) rather than by any test.
**Worth knowing:** A label rename is a restatement sweep: grep the old text across `skills/` before merging. The Panel B rule text "largest idle gap over 2h" deliberately keeps the old words, because it names the rule, not the drawer row. The new lesson is the third `paired-predicate-drift` bullet in `lessons-do-kanban.md`, so the generalized trap in `prime-do-kanban.md` § Traps was amended to cover a reader of a logged, rounded form.

## Orientation

Now the Durations page's Panel B keeps or drops a REQ at the 2h boundary exactly as a re-fit of the calibration log would (whole minutes, rounded down), and the drawer row reads "Largest gap between events"; lives in the queue-kanban Durations view and request drawer (`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`). No map change. Prime spot-check: no prime names the old row label or the exact-duration comparison; `prime-do-kanban.md` § Traps amended its paired-predicate-drift line.

## Heavy Verification Plan

- Base revision: 670428d1e7b267dcf7aaf60e4ced294f7af6ac56
- Target revision: 9882b505c89811f41723a3c53bbeb90f6444076d (landed in `commit:`)
- queue-kanban-javascript — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — board files changed
- queue-kanban-browser — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — board files changed
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed (board code, board-guide.md, estimate-reference.md)

## Heavy Verification Result

- Target revision: 9882b505c89811f41723a3c53bbeb90f6444076d
- Execution revision: 9882b505c89811f41723a3c53bbeb90f6444076d (detached drain checkout `.git/work-run-2026-10-07-135835/drain-head`, run 14:56:50Z to 14:59:08Z, QUEUE_KANBAN_BROWSER set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 7s
- queue-kanban-browser: exit 0, executed (not skipped), 82s
- staged-skills: exit 0, executed, 47s

Green: every selected lane present, exit 0, none skipped, none reused; no HEAVY-RUN-LANE-SKIPPED finding in the log.

## Timing

Observed 2026-10-07T14:29:53Z to 2026-10-07T14:56:38Z: 26m 45s total, 5m 59s attributed across 3 events, 20m 46s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 3m 17s | 1 |
| verification-gate | 2m 27s | 1 |
| handback-merge | 15s | 1 |

Slowest stage: review / independent review agent and review-fix M1 re-merge, 3m 17s, outcome success.
