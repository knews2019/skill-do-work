---
id: REQ-240
title: Stop the Timeline axis printing a fake minute
status: completed
completed_at: 2026-08-18T12:03:00Z
commit: 664b269
claimed_at: 2026-08-18T11:42:03Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T11:42:03Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
created_at: 2026-08-18T11:37:10Z
user_request: UR-052
addendum_to: REQ-235
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Stop the Timeline Axis Printing a Fake Minute

## What

`timelineFormatAxisTick` in `web/board-timeline.js` formats any window of three days or less as `"<day> <Mon> <HH>:00"` — where `:00` is a **string literal**, not the tick's actual minute. So a tick at 11:55 prints `11:00`. Once the window is short enough that several ticks fall inside one hour, they all print the same label.

Measured on the merged tree, immediately after clicking `Now`: **seven ticks, two distinct labels** — `18 Aug 11:00` five times, then `18 Aug 12:00` twice.

## Context

Found by REQ-235's review, by rendering the board and looking at the axis rather than by any assertion — the ticks are correctly *positioned*, so nothing in the suite notices that they are incorrectly *labelled*.

The defect is REQ-227's, not REQ-235's: the literal has been there since the view was built. What REQ-235 changed is how often you meet it. Before, reaching a sub-hour window took deliberate zooming. Now the `Now` button lands you in one by design — it sizes the window to cover the now-line and the forecast, which on a healthy queue is well under an hour — and the new `Day` level is one step away from it. The single most-used new control in REQ-235 lands the reader on an axis that reads as five identical labels.

That is also why it matters more than its size suggests: UR-052's complaint was "I can not jump to the remaining work", REQ-235 built the jump, and the jump's landing state is the one place the axis is least readable.

## Requirements

- An axis tick's label states the tick's real instant. No component of the label may be a literal that the instant does not carry.
- At any window span, two ticks at different instants must not render identical labels — that is the property the current code violates, and it is what makes the axis unreadable rather than merely imprecise.
- Tick labels stay legible at the existing tick count and font size; a longer format that overlaps its neighbour trades one unreadable axis for another.
- No change to tick *positions* — this is a formatting defect, not a layout one.
- The day/week/month period windows and the `Fit all` window keep the labels they render today, where those are already correct.

## Red-Green Proof

**RED prompt/case:** a Node behaviour probe driving `timelineFormatAxisTick` over the tick instants of a sub-hour window (the shape `Now` produces), asserting that the rendered labels are pairwise distinct and that a tick at a non-zero minute does not render `:00`.
**Why RED now:** the `:00` is a literal, so a 1-hour window's ticks collapse to at most two distinct labels; measured at 7 ticks → 2 distinct on the live board.
**GREEN when:** the probe passes and a rendered board, after clicking `Now`, shows an axis whose labels are all different.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Notes for the builder

Two adjacent observations from the same review, both recorded on REQ-235 and neither in this REQ's scope — do not fix them here, but know they exist so you do not design against them:

- A week window's six evenly-spaced ticks land at 1.167-day intervals, so interior labels skip a day (`7, 8, 9, 10, 11` with 12 missing). A period-aware tick *generator* would fix that and this together, but it is an axis redesign and wants its own decision.
- `Fit all` spanning the whole capture history is what crushes the bars into the right-hand fraction of the plot in the first place.

---

## Triage

**Route: B** - Medium

**Reasoning:** One named function and a stated defect, but the fix had to pick a format that stays legible at the existing tick count, and the branch condition turned out to encode an unstated assumption about tick spacing — both needed reading the axis renderer first.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — the tick formatter and the tick-count constant it keys on
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — Node behaviour probe over several window spans

**Files I will NOT touch:** `web/board.css` and `web/template.html` (no markup or styling change — this is a string), `durations.go` / `durations_test.go` (a sibling builder holds them), and the axis tick *generator* (REQ-235 recorded two findings there; positions are explicitly out of scope).

**Acceptance criteria (restated from REQ):**
- [ ] A tick's label states the tick's real instant; no component is a literal the instant does not carry
- [ ] At any window span, two ticks at different instants do not render identical labels
- [ ] Labels stay legible at the existing tick count and font size
- [ ] Tick positions are unchanged
- [ ] The day/week/month windows and `Fit all` keep the labels they render today

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** `timelineFormatAxisTick` now reads the minute off the instant instead of printing the literal `":00"`, giving `18 Aug 11:36` where it previously gave `18 Aug 11:00`. The chosen form is character-for-character the same width as the string it replaces, so the axis geometry does not move.

The branch condition changed from a span threshold (`spanMs <= 3 days`) to a tick-*gap* threshold (`spanMs / TIMELINE_AXIS_TICK_COUNT < TIMELINE_DAY_MS`). The old threshold silently encoded "six ticks", so it was wrong for spans between three and six days — a second instance of the same defect that the REQ did not name. Expressing the rule as *show the time once two ticks can land on the same day* fixes that band and stays true if the tick count ever changes; `renderAxis`'s local `var tickCount = 6` was hoisted into `TIMELINE_AXIS_TICK_COUNT` so the formatter and the generator read one number. A third branch appends the year once the span reaches a year, closing the same collision at the other end as `Fit all`'s range keeps growing.

## Qualification

Passed — 2 files verified in the merge range `950b7e3..664b269`, 5 acceptance criteria traced.

Judgment checks, run against the merged tree rather than taken from the builder's report:
- **Tick positions are provably untouched.** `renderAxis`'s loop arithmetic is unchanged; the only edit there replaces the literal `6` with the constant holding `6`.
- **Requirement 5 verified by comparison, not by claim.** The `Fit all`, `Month`, `Week` and `Day` label strings on the merged board are identical to the ones I measured on the pre-fix board earlier in this session (`28 May…20 Aug`, `1 Jul…1 Aug`, `13 Jul…20 Jul`, `16 Jul 00:00…17 Jul 00:00`). Only `Now` and the sub-day free-zoom band changed, which is the requirement.
- Nothing hollow: every label component is read off the instant.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0 on the merged tree, run unpiped

**Red-green validation:**
- `TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant`: ✗ → ✓. The first RED was already a real assertion failure rather than a reference error, because the function existed — `the Now window draws 7 ticks with only 2 distinct labels: ["18 Aug 11:00" "18 Aug 11:00" "18 Aug 11:00" "18 Aug 11:00" "18 Aug 12:00" "18 Aug 12:00" "18 Aug 12:00"]`, reproducing the review's measurement exactly. A second RED put the old span threshold back with everything else fixed, proving the newly-covered four-day case bites: `the free zoom, four days window draws 7 ticks with only 5 distinct labels`.
- REQ-233's `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard` and `TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer`, and REQ-235's `TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow`: all still pass, unmodified.

**New tests added:**
- `TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant` — drives the formatter over seven windows (`Now`, `Day`, a four-day free zoom, `Week`, `Month`, `Fit all`, and a two-year `Fit all`), asserting per window that labels are pairwise distinct **and** that the whole label matches a pattern built from that tick's own day, hour, minute and year. The second assertion is what makes the test unable to pass on a literal. It also pins the four already-correct windows to the labels they render today, so requirement 5 is a test rather than a promise.

**Render evidence — measured on the merged tree, in the state the defect lives in.** Six states driven in a browser, labels read from the live DOM:

| State | ticks / distinct | labels |
|---|---|---|
| Fit all | 7 / 7 | `28 May … 20 Aug` (unchanged) |
| Month | 7 / 7 | `1 Jul … 1 Aug` (unchanged) |
| Week | 7 / 7 | `13 Jul … 20 Jul` (unchanged) |
| Day | 7 / 7 | `16 Jul 00:00 … 17 Jul 00:00` (unchanged) |
| Day + 3 zoom-out (~4 days) | 7 / 7 | `14 Jul 10:50 … 18 Jul 13:09` — **was 5 distinct** |
| **Now** | **7 / 7** | `18 Aug 11:30 · 11:40 · 11:50 · 12:00 · 12:10 · 12:20 · 12:30` — **was 2 distinct** |

Legibility measured rather than assumed: at a 2352px axis the seven labels are 72–73px wide with gaps of 262–298px and **zero overlapping pairs**.

*Verified by work action*

## Decisions

- **D-01**: Print the real minute keeping the existing format shape (`18 Aug 11:36`), rejecting seconds (the tick gap can never be under ten minutes, so seconds carry no information and cost width) and an ISO/locale form (40% wider; the view is UTC-only by design). Same character count as the string replaced, so requirement 3 holds by construction rather than by luck. DECIDE & STATE.
- **D-02**: Key the format on the gap between ticks rather than on the window span. The old `spanMs <= 3 days` encoded "six ticks" without saying so, and was therefore wrong in the 3-to-6-day band — measured live at 7 ticks / 5 distinct. **This fixes an instance the REQ did not name**, which the builder flagged explicitly as the one place the diff exceeds the headline. Judged in scope: requirement 2 states the property *at any window span*, and the fix is the same one-line condition in the same function with no change to tick generation. DECIDE & STATE.
- **D-03**: Append the year once the span reaches a year, closing the same collision at the long end as `Fit all`'s range keeps growing. Changes nothing today — the live `Fit all` labels are byte-identical — and is a self-contained three-line branch. The builder considered dropping it on scope-discipline grounds and kept it because shipping a fix that is knowingly wrong one zoom level out is the shape of defect this REQ exists to close. **ESCALATE-grade reasoning, DECIDE & STATE outcome.** Value: the formatter is correct at every span the view can reach, not just the ones it reaches today. Risk: low and cheaply reversible — the branch is unreachable on any current board, so a revert is provably a no-op until the archive spans a year.
- **D-04**: One test over seven windows rather than a test per window. The interesting failure is the same one in every case; per-window tests would be method symmetry, not coverage. DECIDE & STATE.
- **D-05**: Tick *generation* deliberately untouched. REQ-235 recorded two findings there (week ticks at 1.167-day intervals; `Fit all` spanning the whole history) and both stay its. Keying on the tick gap is neutral to how ticks are placed, so a future period-aware generator can replace the loop without touching the formatter. DECIDE & STATE.

## Discovered Tasks

- None new. The two adjacent findings the builder was told about are already recorded on REQ-235 and were left alone. One note for whoever picks up the tick generator: the date-only label branch is now reachable only at gaps of a day or more, so a period-aware generator landing ticks on calendar boundaries would let that branch drop the year tier again.

## Review

**Overall: 97%** | 2026-08-18T12:03:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 96% |
| Scope | 94% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 1 (report only)
- Scope sits at 94% because D-02 and D-03 both reach past the REQ's named instance — one to a band the REQ did not mention, one to a span no current board can reach. Both are defensible and both were flagged rather than slipped in, which is the difference between judgment and drift. Recorded so the reach is visible in the trail: the REQ asked for one literal to stop lying, and the diff makes the formatter correct at every span. That is the right call for a formatter, and it is still more than was asked.

**Restatement sweep:** the diff redefines when the axis shows a time, which is stated in the formatter and nowhere else — the `.timeline-hint` line and the panel `aria-label` describe interaction, not label format, and neither mentions tick text. `TIMELINE_AXIS_TICK_COUNT` now has two readers (the formatter's threshold and `renderAxis`'s loop), which is the point of hoisting it; a third reader would need to join the same constant. `_dev/primes/prime-kanban-board.md` is unaffected. No stale restatement.

**Acceptance:** Pass — the state this REQ exists for went from 7 ticks / 2 distinct to 7 ticks / 7 distinct, measured in a browser on the merged tree, with the four already-correct windows byte-identical and zero label overlaps.

**Suggested testing:** 1 item
- Read the axis on the consuming board (677 REQs, four months) at each period level and after `Now`. This repo's own board is three months of history; the format's behaviour near the year boundary has a test but has never been seen.

**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Reading every function that assigns `timelineViewState` before choosing a format. That sweep — `timelineZoomedWindow`, `timelinePannedWindow`, `timelineKeyboardWindow`, `timelinePeriodWindow`, `timelineNowJump`, and the initial fit — established that the window span can never fall below `TIMELINE_MIN_SPAN_MS`, so the tick gap can never be under ten minutes, so **minute resolution is provably sufficient**. Without it, adding seconds would have been a defensible guess; with it, omitting them is a fact. REQ-233's single-window-model constraint paid a dividend here that had nothing to do with why it was built.

**What didn't:** The original threshold's form, which is the actual root cause rather than the literal. `spanMs <= 3 days` reads like a statement about the window and is really a statement about tick spacing — it silently assumed six ticks. That assumption was invisible, untested, and wrong for a whole band of spans nobody had looked at. The literal `":00"` was the visible symptom; a condition that encoded a number it never named was the reason the bug had a second instance.

**Worth knowing:** A threshold that depends on a count should read that count from the same place the count is defined, or it rots the next time someone edits the other. `TIMELINE_AXIS_TICK_COUNT` exists for exactly that reason, and it is the smaller half of this fix that will matter longer than the minute.

## Orientation

The Timeline's axis labels now state the instant each tick actually falls on, at every zoom level the view can reach — the minute is read from the tick instead of printed as `:00`, the time appears once two ticks could land on the same day rather than at an arbitrary three-day threshold, and the year appears once a window is long enough for one day-and-month to come round twice. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`).

Not `[MAP CHANGED]` — one formatting function and one hoisted constant; no contract, payload, or structure moved. It does close the loop opened by REQ-235: that REQ's `Now` button made a sub-hour window the routine landing state, and this makes that window readable. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path resolves and the three-write-surface count is unchanged — this REQ adds none. The prime is not stale.
