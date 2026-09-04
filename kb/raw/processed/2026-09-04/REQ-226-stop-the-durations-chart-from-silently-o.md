---
source_type: req_lesson
req_id: REQ-226
req_path: do-work/archive/UR-051/REQ-226-stop-durations-chart-overprinting-and-clipping.md
date: 2026-08-18
domain: general
module: _dev/primes
tags: [general, stop, durations, chart]
---

# Lessons from REQ-226: Stop the Durations chart from silently overprinting and clipping

## What the REQ was about

Fix two defects in the Durations view where the chart draws something that reads as a value but
isn't. Panel A labels every overflow sample with no collision detection, so on a large board the
overflow lane becomes an unreadable blob of overprinted text. Panel B clamps its bars at 45 minutes
with no scale break, so a 78-minute day renders as a 45-minute bar.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)

## What worked

- Porting the shipped rule into Go to produce RED, rather than writing a test against a function that did not exist yet. A compile error would have been indistinguishable from a real failure; `placed 40 labels, but two rows hold at most 26` names the defect. Rendering a synthetic fixture dense enough to reproduce the reported failure was the other half — the live archive passes every assertion either way, so nothing short of a purpose-built board could show whether the fix worked.

## What didn't work

- Preferring the after-mark anchor. It is the obvious default and it is wrong for a left-to-right greedy: it spends space the next mark still needs. It cost a label on the maintainer's own board before the render caught it (D-02).
- Putting the remainder sentence on the band's first row. The marks sit at that height, so the first dense render showed `+58 more over 60 min` overprinted by the blob it was describing — the defect reproduced inside its own fix (D-03).
- Restoring a file from a copy taken earlier in the session. The mutation-test restore silently reverted D-03, and the whole suite still passed, because at that point nothing pinned which row the remainder used. The near-miss is why `durationsRemainderBaselineY` exists as a named function with its own probe assertion rather than as an expression inline in the render pass.

## Worth knowing

- Three of the four defects found here were invisible to the test suite and visible in a render. A chart's correctness is partly a claim about pixels, and neither a passing assertion nor a payload dump can see two glyphs sharing a coordinate. Generate a board and look at it. The corollary bit twice: a fix verified only by tests that pass both before and after it is not verified.

## Back-reference

See `do-work/archive/UR-051/REQ-226-stop-durations-chart-overprinting-and-clipping.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `787c846`.
