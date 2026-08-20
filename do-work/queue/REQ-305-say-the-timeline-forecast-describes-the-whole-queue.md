---
id: REQ-305
title: "Say the timeline forecast describes the whole queue, not the filtered rows"
status: pending
created_at: 2026-08-20T08:37:41Z
user_request: UR-062
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
related: [REQ-303, REQ-304]
batch: upstream-consumer-report-2026-08-20
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Say the Timeline Forecast Describes the Whole Queue, Not the Filtered Rows

## What

The Timeline view filters its rows (`web/board-timeline.js:486-488`) but reads the projection
unfiltered (`:497`), then hands both to `renderTimelineForecast` (`:554`) — whose `rows` parameter it
never reads. With a filter active and a nonempty subset showing, the view can say it holds three REQs
above a forecast that schedules the entire queue and an excluded list naming IDs that are not on
screen. Label the forecast for what it is instead of letting it read as a statement about the visible
rows.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The code already states the rule it breaks. At `:508-510` the zero-match branch clears the forecast
with the comment "The forecast describes the rows; with none on screen it must go too." The
partial-filter case is the same situation with a subset instead of an empty set, and it was missed.

## Context

From the 2026-08-20 consumer review, filed P2. Verified: the `rows` argument is genuinely unused
inside `renderTimelineForecast` (`:392-478`), and the excluded list at `:439-460` renders
`projection.excluded` verbatim.

The reviewer offered two remedies. **Recomputing or filtering the projection is rejected**: the
medians, the bucket split, the chain scheduling and the exclusion reasons are all built Go-side
(`timeline.go`, `generate.go`), and a per-filter client-side recomputation would fork the definition
of the forecast into a second implementation that drifts. `:483-485` already documents why filters
apply to the Gantt but not to the Durations distribution — the same reasoning applies to a
distribution-derived forecast.

Note that the projected **segments** are already filter-consistent: they are drawn per row through
`projectedById[row.id]` (`:643-646`), so only the prose and the excluded list are wrong.

## Detailed Requirements

- While any filter is active and a nonempty subset is showing, the forecast text and the excluded
  list must say plainly that they describe the whole queue, not the rows on screen.
- Preserve the existing behavior in both settled cases: no filters (unlabelled forecast, as today)
  and zero matches (forecast cleared entirely, as today).
- Resolve the unused `rows` parameter — either use it for the filter-active determination or drop
  it. Do not leave a parameter the function ignores.
- Do not recompute, re-filter, or re-derive the projection client-side.

## Constraints

- Renderer-only; the Go projection builder stays untouched.
- The forecast sentence is one of the artifacts the REQ-241 note calls out as screenshot-and-quoted
  material — whatever wording is chosen has to survive being read out of context.

## Red-Green Proof

**RED prompt/case:** Filter the Timeline view to a single domain that matches 3 of 25 REQs. Today the
summary says 3 REQs while the forecast schedules all 25 and the excluded list names IDs no row on
screen carries.
**Why RED now:** `rows` is filtered, `projection` is not, and nothing in the forecast copy says which
population it describes.
**GREEN when:** With that filter active the forecast is explicitly labelled as covering the whole
queue rather than the visible rows; with no filter the copy is unchanged; with zero matches the
forecast is still cleared entirely.
**Validation:** Inferred during capture — from the filter path and the zero-match comment read
together.

## Builder Guidance

Certainty: Firm on the diagnosis, mixed on the wording — pick a phrasing that reads correctly when
the forecast paragraph is screenshotted alone. `generate_test.go`'s
`TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions` (~line 2696) already stubs the DOM and
runs the real `renderTimelineForecast`; extend that probe rather than building a new harness.

## Full Context
See `do-work/user-requests/UR-062/input.md` for complete verbatim input.

---
*Source: "[P2] Keep filtered timelines and forecasts in sync — … rows and the summary/table are filtered, but projection remains global and renderTimelineForecast does not use its rows argument. The view can therefore say it contains three REQs while forecasting and listing exclusions for the entire queue, including hidden IDs; recompute/filter the projection or suppress and clearly label the global forecast while filters are active."*
