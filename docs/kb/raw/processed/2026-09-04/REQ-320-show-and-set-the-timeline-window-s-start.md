---
source_type: req_lesson
req_id: REQ-320
req_path: do-work/archive/UR-065/REQ-320-show-and-set-the-timeline-windows-start-and-end.md
date: 2026-08-23
domain: frontend
module: _dev/primes
tags: [frontend, show, timeline, window]
---

# Lessons from REQ-320: Show and set the timeline window's start and end

## What the REQ was about

State the visible window's start and end instants in text, and add two date fields that set
them.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

## What worked

- Routing the new control through the settle every other mover uses, so it has
- no floor, ceiling or clamp of its own — the review confirmed that at both bounds and far
- outside them without finding a window the other controls cannot reach. And rendering the
- board before believing the code: the first render produced D-01, and the review's render
- produced the rest. Two of the three defect classes in this REQ were found by looking.

## What didn't work

- **One function name for two conversions that are not inverses.** `timelineEpochToDateField`
- read as a general converter and was correct only on a window start; using it on an exclusive
- end put every period window a day out. The name is what hid it — a start and an end are
- different types wearing the same one. Two names, and the asymmetry becomes visible at the
- **Re-reading both fields on every commit.** Editing one field then re-applying the other's
- day-truncated value is how a control silently moves something the reader did not touch.
- Apply the field that changed; leave the other endpoint at its exact instant.
- **"Skip the focused field" as the mid-edit rule.** Focus is not editing. A reader who clicks
- into a field and then zooms leaves it stale, and committing it later undoes their zoom.
- Compare against the value the code last wrote instead.
- **Thirteen unit cases and none of them a round trip.** Every case built a window from typed
- text; none rendered a window into the fields and parsed it back. That single assertion kills
- two of the four defects on its own, and it is the assertion the shape of the feature was
- asking for.

## Worth knowing

- the window is a half-open interval and the date fields are inclusive days.
- Anything that converts between them belongs in the pair
- `timelineStartEpochToDateField` / `timelineEndEpochToDateField`, and the round-trip case in
- `TestJavaScriptBehaviorTimelineTypedDatesMoveTheWindow` is what keeps them inverses. The
- period chips depend on that exactness: `timelinePeriodLevelOfWindow` compares instants for
- equality, so a window 1 ms off a calendar period reads as `custom span`.

## Back-reference

See `do-work/archive/UR-065/REQ-320-show-and-set-the-timeline-windows-start-and-end.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c0382e5`.
