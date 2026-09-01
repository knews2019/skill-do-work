---
source_type: req_lesson
req_id: REQ-305
req_path: do-work/archive/UR-062/REQ-305-say-the-timeline-forecast-describes-the-whole-queue.md
date: 2026-08-21
domain: frontend
module: _dev/primes
tags: [frontend, timeline, forecast, describes, whole]
---

# Lessons from REQ-305: Say the timeline forecast describes the whole queue, not the filtered rows

## What the REQ was about

The Timeline view filters its rows (`web/board-timeline.js:486-488`) but reads the projection
unfiltered (`:497`), then hands both to `renderTimelineForecast` (`:554`) — whose `rows` parameter it
never reads. With a filter active and a nonempty subset showing, the view can say it holds three REQs
above a forecast that schedules the entire queue and an excluded list naming IDs that are not on
screen. Label the forecast for what it is instead of letting it read as a statement about the visible
rows.

## Solution summary

The ignored `rows` parameter became `filtersActive`, the single bit the caller can
answer and the function cannot. When it is true the forecast paragraph and the excluded heading lead
with "Filters are on; this covers the whole queue, not the rows shown." and " from the whole queue".
Nothing about the projection changed — it is still built Go-side and consumed verbatim.

## Worth knowing

- **A test that calls the function under test directly cannot hold its call site.** Five mutations
  of the copy were caught and the sixth — the one that reverted the actual defect — passed clean.
  When the bug is an argument the caller computes, the probe has to start above the caller.
- **Changing a parameter's type silently re-points every existing call.** `renderTimelineForecast(p, [])`
  kept compiling and kept running, and `[]` is truthy, so the two existing probe calls flipped to
  asserting the filtered copy. Nothing failed loudly; the negative assertion on the unfiltered text
  is what surfaced it.
- **Put a caveat where a crop keeps it.** The constraint "must survive being screenshotted alone"
  decides position, not just wording — a trailing qualifier is exactly what a crop removes.

## Back-reference

See `do-work/archive/UR-062/REQ-305-say-the-timeline-forecast-describes-the-whole-queue.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ef0cc55`.
