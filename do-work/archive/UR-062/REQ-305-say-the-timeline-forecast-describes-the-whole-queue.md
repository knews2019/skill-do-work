---
id: REQ-305
title: "Say the timeline forecast describes the whole queue, not the filtered rows"
status: completed
created_at: 2026-08-20T08:37:41Z
status_changed_at: 2026-08-21T08:30:38Z
claimed_at: 2026-08-21T08:30:38Z
completed_at: 2026-08-21T08:40:45Z
kb_status: promoted
kb_entry: REQ-305-say-the-timeline-forecast-describes-the-.md
commit: ef0cc55
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
estimate:
  p50_active_minutes: 20
  confidence: high
  calculated_at: 2026-08-21T08:36:08Z
  basis:
    - Route B
    - 2-file write set
    - existing probe extended
    - 4 acceptance criteria
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
- [x] **[PLAN]:** Replace the unused `rows` parameter with `filtersActive`, computed at the call site from `rows.length < timeline.rows.length`; lead both the forecast paragraph and the excluded heading with the whole-queue label when it is true. Extend the existing forecast probe, and add a second probe that drives the whole `renderTimelineView` so the wiring is covered too.
- [x] **[APPLY]:** `web/board-timeline.js` and `generate_test.go` only. The Go projection builder is untouched; nothing is recomputed client-side.
- [x] **[UNIFY]:** `git diff --stat` reviewed (2 files); `gofmt -l .` clean; `go vet ./...` clean; `go test -count=1 ./...` green in 70s; no debug artifacts; both changed files re-read at every edit site.

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

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — the `filtersActive` parameter, the whole-queue label, the call site
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the extended forecast probe and the new whole-view probe

**Files I will NOT touch:** `timeline.go` and `generate.go` — the projection is built Go-side and this REQ forbids re-deriving it client-side.

**Acceptance criteria (restated from REQ):**
- [x] With a filter active and a nonempty subset showing, the forecast text and the excluded list say they describe the whole queue
- [x] No filters: copy unchanged. Zero matches: forecast still cleared entirely
- [x] The unused `rows` parameter is resolved
- [x] The projection is not recomputed, re-filtered, or re-derived client-side
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified) — `renderTimelineForecast(projection, filtersActive)`; the `wholeQueueNote` lead on all three forecast branches; " from the whole queue" in the excluded heading; the call site computes the one bit
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — three filtered cases added to `TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions`; new `TestJavaScriptBehaviorTimelineForecastLabelsAFilteredView`; the REQ-304 stub's filter made injectable

**What was done:** The ignored `rows` parameter became `filtersActive`, the single bit the caller can
answer and the function cannot. When it is true the forecast paragraph and the excluded heading lead
with "Filters are on; this covers the whole queue, not the rows shown." and " from the whole queue".
Nothing about the projection changed — it is still built Go-side and consumed verbatim.

## Decisions

- **D-01 — The parameter became a boolean rather than being used as a row list.** DECIDE & STATE.
  The REQ allowed either using `rows` or dropping it. Using it would mean the function comparing the
  filtered rows against something to learn whether a filter is on, and only the caller holds the
  unfiltered list. Passing the answer keeps the filter logic at the one site that already owns it.
- **D-02 — The label leads the paragraph instead of trailing it.** DECIDE & STATE. The REQ's
  constraint is that this sentence survives being screenshotted alone. A trailing caveat is the part
  a crop removes; a leading one is the part a crop keeps.
- **D-03 — A second probe drives the whole `renderTimelineView`, not just the forecast function.**
  DECIDE & STATE. The first mutation round proved the gap: hard-coding the call site to `false`
  passed every extended assertion, because they all called `renderTimelineForecast` directly. The
  defect was in the wiring, so a test that never exercises the wiring cannot hold it.

## Testing

**Tests run:** `QUEUE_KANBAN_BROWSER=<chromium> bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 — `Maintainer verification passed.` Board module: `go test -count=1 ./...` ok in 70s.

**Red-green validation:**
- **RED, on the untouched renderer, all three arms at once.** With filters on, the forecast read
  `"Queue empties around 2026-06-20 14:30 UTC — 2 REQs run one at a time…"` with no statement of
  which population; the excluded heading read `"1 REQ is not in that estimate…"` above an ID no
  visible row carried; the declined branch read `"No end estimate: …"` unlabelled.
- **GREEN.** All three now lead with the whole-queue label, and the unfiltered copy is byte-identical
  to before — asserted in the negative, so an always-on label fails.

**Mutation-tested — eight reverts, in two rounds:**

Round one, the copy (`renderTimelineForecast`):
- M1 remove the label → filtered forecast assertion, caught
- M2 label always on → unfiltered assertion, caught
- M3 label replaces the body instead of leading it → "does not state the end instant", caught
- M4 excluded heading unlabelled → excluded assertion, caught
- M5 declined branch unlabelled → declined assertion, caught

Round two, the wiring. **M6 — hard-code the call site to `false` — passed every assertion above.**
Every one of them called `renderTimelineForecast` directly, so none could see the call site. That
gap is what `TestJavaScriptBehaviorTimelineForecastLabelsAFilteredView` closes; with it in place:
- M6 call site always unfiltered → caught
- M7 call site always filtered → caught
- M8 compare the rows against themselves → caught

**New tests added:**
- Three filtered cases in `TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions`, plus a
  negative assertion pinning the unfiltered copy.
- `TestJavaScriptBehaviorTimelineForecastLabelsAFilteredView` — drives the whole
  `renderTimelineView` twice, once unfiltered and once filtered to one of three rows, and asserts the
  summary count differs before asserting anything about the label, so the fixture cannot silently
  stop filtering.

**Existing tests updated (cross-REQ impact):** the two existing
`renderTimelineForecast(…, [])` calls in the same probe became `(…, false)`. The parameter changed
from an ignored row list to a boolean; `[]` is truthy, so leaving them would have asserted the
filtered copy against the unfiltered expectations. No assertion changed meaning.

*Verified by work action*

## Review — 2026-08-21T08:40:28Z

**Overall: 96%**

| Dimension | Score |
|---|---|
| Requirements Compliance | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope Discipline | 100% |
| Risk | None |

**Acceptance: Pass.**

**Requirements Compliance.** All four hold. The zero-match path is untouched and still clears both
nodes (its existing assertion still passes); the no-filter copy is pinned negatively; the parameter
is used, not ignored; and no projection value is recomputed — `renderTimelineForecast` reads the
same fields it always read.

**Findings**

- **F1 — Important, found and fixed inside this REQ.** The first eight assertions could not
  distinguish a correct call site from one hard-coded to `false`. Caught by mutation, not by
  reading. The second probe closes it.
- **F2 — Minor, accepted.** `rows.length < (timeline.rows || []).length` infers "a filter is
  active" from a count rather than asking the filter state directly. It is correct for every filter
  the view has (they can only remove rows), and it is one expression at one site. A filter that
  could add rows would break it, and none can.

**Constraint check.** Renderer-only holds. `git diff --stat` touches `web/board-timeline.js` and
`generate_test.go`. The wording was chosen against the screenshot constraint: the label is the
first clause of the paragraph, so a crop keeps it.

## Lessons Learned

- **A test that calls the function under test directly cannot hold its call site.** Five mutations
  of the copy were caught and the sixth — the one that reverted the actual defect — passed clean.
  When the bug is an argument the caller computes, the probe has to start above the caller.
- **Changing a parameter's type silently re-points every existing call.** `renderTimelineForecast(p, [])`
  kept compiling and kept running, and `[]` is truthy, so the two existing probe calls flipped to
  asserting the filtered copy. Nothing failed loudly; the negative assertion on the unfiltered text
  is what surfaced it.
- **Put a caveat where a crop keeps it.** The constraint "must survive being screenshotted alone"
  decides position, not just wording — a trailing qualifier is exactly what a crop removes.

## Orientation

**What changed in the map.** The Timeline's forecast now knows which population it describes. It
still never computes that population — the projection is built Go-side and consumed verbatim — but
the one bit it needs, whether the rows above it are a subset, now reaches it from the caller that
already knows.

**What this makes true.** Filtering the Timeline to one domain no longer produces a view that says
"3 REQs" above a forecast scheduling twenty-five and an excluded list naming IDs that are not on
screen. The board also gained a probe that exercises `renderTimelineView` end to end, which is where
the next wiring defect will be caught.

**Subsystem:** the queue-kanban board tool's Timeline view. Prime: `_dev/primes/prime-kanban-board.md`.
