---
id: REQ-352
title: "State the Durations view's own headline numbers"
status: claimed
claimed_at: 2026-08-24T16:43:37Z
status_changed_at: 2026-08-24T16:43:37Z
route: B
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-350]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-346, REQ-347, REQ-348, REQ-349, REQ-350, REQ-351, REQ-353, REQ-354]
batch: durations-panel-improvement
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-24T16:43:37Z
  basis:
    - Route B
    - 3-file seeded write set
    - 7 acceptance criteria
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go
---

# State the Durations View's Own Headline Numbers

## What

The view makes a reader hover to learn the numbers it exists to report. Add a stat tile row above
panel A, a rolling median line on panel B, and the ticks panel C is missing.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route B exploration traced the existing projected samples/days and R-7 quantile helper. Add four semantic tiles, a trailing seven-Panel-B-active-day median with a one-point marker case, and exact zero/midpoint Panel C ticks, then pin complete-renderer and live responsive/theme behavior.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The summary line gives counts and the exclusion rule and nothing else. The median, the p90 and the
completions-per-active-day rate — the numbers a reader opens a durations view for — appear nowhere.
Panel B annotates exactly one day (its slowest) and draws no trend; panel C labels only its peak,
with no zero tick and no gridline. Every one of those figures is derivable from the payload already
shipped.

## Detailed Requirements

- **A small stat tile row above panel A**: median REQ duration, p90, active-day count against the
  axis span, and REQs per active day.
- **Each figure is stated for the window actually plotted**, not for all history.
- **A rolling median line on panel B** over its per-day bars, so "is it getting slower" is answerable
  without reading 29 bars one at a time. Specified so two implementations cannot draw materially
  different lines and both claim to satisfy this:
  - **Window: the trailing 7 active days** — the bars panel B actually draws, not 7 calendar days.
    The axis is 66 percent idle here, so a calendar window collapses to one or two samples through
    every quiet stretch and the line stops meaning anything exactly where the reader looks hardest.
  - **Trailing, not centered.** The point at a given day summarizes that day and the six drawn days
    before it, so the line never depends on days to its right — a centered window would make the
    most recent points the least stable, which is the wrong end to be unstable.
  - **Edge: draw nothing before the seventh active day.** No partial windows. A median of two days
    plotted on the same line as a median of seven invites reading noise as trend, and the left edge
    is where that misreads worst.
  - **Gaps: idle calendar days are skipped, never zero-filled.** They carry no samples, and a zero
    would be a measurement the archive does not contain. The line's x positions stay on the real
    calendar axis, so a 7-active-day window may span far more than 7 days of chart width — that is
    the honest shape and it must not be compressed.
- **Panel C gets its zero and midpoint ticks.** Today only its peak is labelled.
- **Derive every figure from the samples already in the payload.** No new aggregate in `durations.go`.
- **Keep one measure per scale**, and keep the existing summary sentence's exclusion-rule wording
  intact.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- Panels A and B disagree about which samples count on purpose, and that rule lives in one place.
  Whichever panel a tile describes, it uses that panel's rule and says so if the two could be
  confused.
- Generate a board and look at it.

## Dependencies

`depends_on: REQ-350` — the tiles are stated for the plotted window, so the window state must exist
first. This resolves the report's own "whether the tiles follow A3's window or all history": they
follow the window.

## Builder Guidance

**Certainty: firm.** Four tiles, one line, two ticks, all from existing data. Keep it small — this is
the cheapest legibility gain in the batch and it should not grow a new aggregation layer.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repository and open Durations. Nowhere on the view does
it say the median is 13.2 minutes or the p90 is 42.7; panel B draws 29 bars with one annotation and no
trend; panel C's only labelled value is its peak.

**Why RED now:** The summary composes counts and the exclusion rule only, and neither panel B nor
panel C draws an aggregate or a full tick set.

**GREEN when:** a tile row above panel A states median, p90, active days against axis span and REQs
per active day for the plotted window; panel B carries a rolling median line; panel C has zero and
midpoint ticks; and the summary sentence's exclusion-rule wording is unchanged.

**Validation:** User confirmed (bundled invocation).

---
*Source: prompt A7, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html` (finding F7).*

## Triage

**Route: B** — The tile measures, exact seven-active-day trailing rule, tick outcome, payload source, and three-file surface are specified. Exploration must locate existing quantile, day-bar, axis, and rendered-probe conventions.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

- Median and p90 tiles use Panel A's rule: all plotted raw spans, including excluded paused/reversed samples, visibly labelled “all plotted spans.” Reuse `durationQuantile` and `formatDurationMinutes`.
- Cadence tiles use Panel C's rule: active days are `days.length` against `timeSpan / DURATIONS_DAY_MS`; REQs per active day is `samples.length / days.length` to one decimal. Excluded-only completion days remain active.
- Panel B's chronological input is exactly `days.filter(day.hasMedian)`. Starting at index six, take the trailing seven drawn-day medians, compute their R-7 median, and place it at the current day's real-calendar noon. Draw after bars, add point markers so exactly seven active days has a visible result, and draw nothing for six or fewer.
- Idle and excluded-only days never enter the rolling window and are never zero-filled. The visible title and SVG copy state “trailing 7-active-day median.”
- Panel C keeps its peak/baseline, adds zero and exact arithmetic `peakCount / 2` ticks/gridline, and formats an odd midpoint honestly with one decimal.
- A semantic four-item `<dl>` above Panel A uses an auto-fitting grid. Live proof must cover window updates, light/dark line contrast against body, 320/768/1280 wrapping, finite/non-overlapping geometry, exact tick placement, and console state.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`
