---
id: REQ-242
title: Stop Panel B's slowest-day annotation colliding with its own title
status: claimed
status_changed_at: 2026-08-18T13:05:12Z
created_at: 2026-08-18T12:09:46Z
user_request: UR-051
addendum_to: REQ-237
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T13:05:12Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
claimed_at: 2026-08-18T13:05:12Z
route: B
---

# Stop Panel B's Slowest-Day Annotation Colliding With Its Own Title

## What

In the Durations view, Panel B's slowest-day annotation is drawn at `y = 355` while Panel B's own title sits at `y = 350` (`DURATIONS_MEDIAN_TITLE_Y`). The two overlap: on a synthetic fixture the annotation `209 min` renders directly through the words "paused and broken spans excluded".

## Context

Found while reviewing REQ-237 by rendering a dense fixture and looking at it. **It is not a REQ-237 regression** — the same annotation sits at the identical `x = 357.2, y = 355.0` on a board built from the pre-REQ-237 binary, checked side by side. It is pre-existing and was simply never looked at on a fixture whose slowest day lands under the title text.

It is invisible on this repository's own board because the annotation's x-position depends on which day is slowest, and here that day falls clear of the title's width. That is luck, not design — which is why it wants pinning rather than nudging.

The annotation reuses the `durations-mark-label` class, so it is not part of either label band's row packing and is not covered by REQ-231's mark-band geometry test or by REQ-237's row-fill test. Nothing in the suite looks at it at all.

## Requirements

- Panel B's slowest-day annotation does not overlap Panel B's title at any x-position the annotation can take, including when the slowest day is the leftmost one.
- The annotation stays associated with the bar it describes — moving it somewhere it no longer reads as belonging to that day is not a fix.
- No change to Panel A, Panel C, or the label bands; no change to `describeAtPointer`'s panel boundary.
- A test pins the separation, so the next person to move `DURATIONS_MEDIAN_TITLE_Y` finds out.

## Red-Green Proof

**RED prompt/case:** an assertion that the slowest-day annotation's text box and Panel B's title text box do not intersect, read from the renderer's own constants the way `TestDurationsLabelRowsClearTheMarkBands` reads them — evaluated at the annotation's worst-case x, not at whichever x this repository's data happens to produce.
**Why RED now:** the title's baseline is 350 and the annotation's is 355, so their boxes intersect wherever their x-ranges do; reproduced on a fixture as `209 min` drawn through the title.
**GREEN when:** the assertion passes and a rendered fixture whose slowest day sits under the title shows the two clear of each other.
**Validation:** Review finding on REQ-237; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-237 I rendered the Durations chart on a test board and found the little "209 min" note that marks the slowest day printed straight through the heading above it — the note sits five units below a heading that is taller than five units. It has been like this since before any of today's work; it does not show on your own board only because the slowest day happens to fall to the right of where the heading's text ends, which is chance rather than design. The fix is small — move the note, or move the heading, and add an assertion so the next person who shifts either one finds out. I am asking rather than doing it because "move which one, and where" is a look-and-feel choice about a chart you read regularly, and there is more than one reasonable answer. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — it only shows on contrived data and the note is legible enough where it is.
  → Confirmed: Yes, add to queue (builder picks placement, pinned by a test). [2026-08-18, via do-work clarify]
