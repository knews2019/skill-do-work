---
id: REQ-637
title: '[impact-rule-change] Board Panel B compares the stamp gap in whole minutes like the calibration log, and the drawer row is relabelled to mean the gap between events'
status: pending
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
write_set: [skills/do-work-board/tools/queue-kanban/durations.go, skills/do-work-board/tools/queue-kanban/durations_test.go, skills/do-work-board/tools/queue-kanban/web/board-detail.js, skills/do-work/CHANGELOG.md, skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md, do-work/lessons-index.md]
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
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (7253 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matched: family `paired-predicate-drift` is exactly F4's failure shape.
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`). Matched: it governs finding-detail and view text, which the drawer relabel changes.

## Full Context
See `do-work/user-requests/UR-137/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: review findings F4 and F6 from the pasted consumer code review of 0.305.58 to 0.305.69, triaged with `do-work validate-feedback` on 2026-10-06 and accepted (F6 as a relabel).*
