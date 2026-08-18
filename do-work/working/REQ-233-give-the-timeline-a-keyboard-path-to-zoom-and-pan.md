---
id: REQ-233
title: Give the Timeline a keyboard path to zoom and pan
status: claimed
claimed_at: 2026-08-18T11:00:00Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-18T11:00:00Z
  basis:
    - Route B
    - 3-file write set
    - browser evidence
    - cross-route regression gates
status_changed_at: 2026-08-18T10:26:34Z
domain: general
created_at: 2026-08-18T01:18:57Z
user_request: UR-051
addendum_to: REQ-227
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- skills/do-work-board/tools/queue-kanban/web/template.html
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Discovered Task: Give the Timeline a Keyboard Path to Zoom and Pan

## What

The Timeline view's zoom and pan are pointer-only. The three zoom buttons are reachable by keyboard, but panning the time axis is not, and both affordances are described only in a hint line of prose beside the chart rather than anywhere assistive technology reads.

## Context

Found while implementing REQ-227, which built the view. Every value the chart draws is also in the table below it, so no *data* is unreachable — what is unreachable is the navigation. A keyboard user can read every row but cannot move the window they are read in.

**Narrowed since capture.** This REQ originally recorded that rows "already take focus and open the detail drawer with Enter". That was wrong: the rows carried `role="button"` and `tabindex="0"` but the drawer opened only on a delegated `click`, which a non-native `<g>` never synthesizes from Enter or Space. PR #144's review caught it and it was fixed there — row activation now works from the keyboard and is pinned by `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard`. What remains for this REQ is only what its title says: arrow-key panning and `+`/`-` zoom, plus stating the interaction in the panel's accessible name.

Not a regression: no other board view has zoom or pan at all, so REQ-227 added an un-keyboarded capability rather than removing a keyboarded one. That is why it is a follow-up rather than a fix inside REQ-227.

## Requirements

- With the chart focused, arrow keys pan the time axis and `+`/`-` zoom it, using the same `timelineZoomedWindow` transform the pointer path uses, so the two cannot diverge.
- The panel's accessible name states the interaction, rather than leaving it to the visual hint line beside the chart.
- Focus is visible on whatever element takes the keyboard interaction — a focus ring that exists only in the default user-agent style is not enough on a dark surface.
- No change to the pointer path's behavior.

## Red-Green Proof

**RED prompt/case:** A Node behavior probe driving the keyboard handler through the same `timelineZoomedWindow` transform, asserting that an arrow-key pan shifts the window by a bounded step and clamps at the range edges, and that `+`/`-` reach the same floor and ceiling the pointer path reaches.
**Why RED now:** there is no keyboard handler, so there is nothing to drive.
**GREEN when:** the probe passes and a headless run confirms the chart takes focus and responds to the keys.
**Validation:** Discovered during REQ-227; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-227: the new Timeline chart can be zoomed and dragged with a mouse, but there is no way to do either from the keyboard — the three zoom buttons can be tabbed to, dragging cannot be done at all, and the instructions for both sit in a line of text next to the chart rather than anywhere a screen reader announces. Nothing is unreadable: every row can be focused, opens its detail panel with Enter, and every number is repeated in the table underneath. What a keyboard user cannot do is move the time window they are reading in. Adding arrow-key panning and `+`/`-` zoom is a small, self-contained change to the one file. It is your call rather than mine because no other board view has zoom or pan at all, so this is a new capability to finish rather than a regression to repair, and you may prefer it batched with a wider accessibility pass over the board instead of done piecemeal here. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  *Answered 2026-08-18 via `do-work clarify` — user approved building the keyboard path now rather than batching it into a wider accessibility pass.*
