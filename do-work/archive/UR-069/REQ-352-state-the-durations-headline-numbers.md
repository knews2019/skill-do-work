---
id: REQ-352
title: "State the Durations view's own headline numbers"
status: completed
claimed_at: 2026-08-24T16:43:37Z
completed_at: 2026-08-24T17:22:38Z
commit: a522484997a7c07c1cc5d6875f7be943c7d72c26
status_changed_at: 2026-08-24T17:22:38Z
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
- [x] **[APPLY]:** Added window-scoped headline statistics, a trailing seven-active-day median series, exact cadence ticks, semantic markup, responsive styling, and focused renderer/browser lock-ins in the declared five-file scope.
- [x] **[UNIFY]:** Reviewed the complete five-file merge diff and integration seams; `node --check`, focused Go/browser tests, the full module suite, `go vet`, and canonical maintainer verification passed with no debug artifacts.

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

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified): computes window-scoped raw median, p90, active-day coverage, and REQs per active day; draws the trailing seven-eligible-active-day Panel B median; and adds Panel C zero, midpoint, and peak ticks.
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified): adds the semantic four-item Durations definition list.
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified): adds responsive stat tiles and rolling-series presentation.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified): pins headline values, 6/7/8-day rolling boundaries, active-day gaps, draw order, and odd midpoint geometry.
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go` (modified): proves complete generated-board semantics, responsive layout, theme contrast, finite geometry, window transitions, tick separation, and clean console state.

## Decisions

- **D-01 — Headline duration statistics describe every plotted Panel A span.** Median and p90 use signed raw `wallMinutes`, including paused and reversed spans, and the visible label says “all plotted spans.”
- **D-02 — Cadence uses the projected UTC-day axis.** Active days include excluded-only completion days; the denominator is the complete plotted UTC-day span, and REQs per active day uses every projected sample.
- **D-03 — Exactly seven eligible Panel B days draw one honest marker.** Six or fewer draw nothing; eight or more add the connecting path, while idle and excluded-only days never enter the window.
- **D-04 — Share the existing quantile and y-scale rules.** The headline, daily distribution, and rolling series use one R-7 implementation, and high rolling values clip to Panel B's established 45-minute ceiling.
- **D-05 — Contrast against the real surface.** Rolling ink uses `--ink-strong` because the SVG is transparent over the page body in both themes.

## Discovered Tasks

- The pre-existing dense Panel A browser probe placed its result node after the script that writes to it. The concurrent REQ-351 branch already repairs that harness seam, so no duplicate follow-up was created.

## Testing

- RED proved the semantic statistic row and all four values were absent and that seven/eight eligible-day fixtures drew neither the required rolling marker nor path.
- Focused renderer tests passed for exact median/p90/cadence values, six/seven/eight-day rolling boundaries, active-day gaps, draw order, and odd midpoint ticks.
- Builder browser proof passed at 320/768/1280 in light and dark with semantic focus behavior, responsive tile wrapping, finite chart geometry, distinct window values, separated cadence ticks, and no application console errors.
- Builder full module tests and `go vet ./...` passed; canonical maintainer verification passed. The post-merge `GOTOOLCHAIN=go1.26.1 go test ./... -count=1 -timeout=10m` suite passed in 73.589s.
- The orchestrator generated this repository's real queue and inspected desktop and 320px renders. The four tiles showed `15.0 min`, `55.5 min`, `25 / 30`, and `11.9`, remained readable without collision, and the only console item was the unrelated missing favicon.

## Qualification

- Exact merge range `0d933971..a522484` passed mechanical qualification.
- Scope drift passed: the five-file Implementation Summary exactly matches the declared Scope.
- Orchestrator judgment confirmed substantive implementation, requirement coverage, one shared quantile/data path, preserved exclusion wording, and no generated/debug artifacts in the merge.

## Review

Independent review approved with no Important, Minor, or Nit findings. Requirements 100%, Code Quality 98%, Test Adequacy 99%, Scope Discipline 100%, overall 99%, low risk, acceptance pass. The reviewer independently reran focused renderer tests and the four-case Chromium matrix and inspected both fixture and real-queue evidence.

## Lessons Learned

Headline statistics are trustworthy only when their visible labels name the population they summarize, and rolling trends remain honest only when their eligibility, minimum window, and real-calendar positioning share the bars' existing data rules.

## Orientation

Released in 0.236.46. Durations now states its window-scoped median, p90, active-day coverage, and cadence; Panel B shows a trailing seven-active-day median; and Panel C labels zero, midpoint, and peak.
