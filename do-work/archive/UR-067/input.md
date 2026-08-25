---
id: UR-067
title: Timeline click regression plus two audit leftovers
created_at: 2026-08-23T18:30:26Z
requests: [REQ-336, REQ-337, REQ-338]
word_count: 490
---

# Timeline Click Regression Plus Two Audit Leftovers

## Summary

A Timeline (Gantt) view regression in the board tool (`skills/do-work-board/tools/queue-kanban/`) plus two leftovers from the same audit line. Pointer capture on pointerdown (introduced by REQ-333, commit 36c4518) retargets the synthesized click so no mouse click inside `#timeline-scroll` reaches the delegated drawer handler. The capture also asked two questions before capturing two candidate items; both were answered during capture (see Capture Decisions).

## Extracted Requests

| REQ | Title | Source in input |
| --- | --- | --- |
| REQ-336 | [impact-critical] Timeline clicks open the detail drawer again | "REQ 1" section |
| REQ-337 | A check that can catch Timeline click retargeting | "REQ 2" section |
| REQ-338 | Cut the Timeline row list to one Tab stop | Q1, captured on the user's answer during capture |

## Capture Decisions

Resolved interactively during capture (ask tool):

- **Q1 (keyboard stop count)** — user chose **capture it**: one Tab stop for the Timeline row list, arrow keys move between rows (roving tabindex). Became REQ-338.
- **Q2 (rightmost axis tick can name a day the "to" field excludes)** — user chose **leave as is**: REQ-334's judgment stands (a tick names an instant, the fields name days of coverage). No REQ; this line is the record.
- **REQ 2's impact tag** — the input's `[impact-high]` is not a canonical `impact:` value; normalize-and-warn fired and the user confirmed **impact-user-visible** for REQ-337 (default tier, so its title carries no tag).

## Batch Constraints

- Board tool root: `skills/do-work-board/tools/queue-kanban/`.
- REQ-336's blast radius is already audited — do not re-widen it (mouse clicks inside `#timeline-scroll` only).
- REQ-337 must be mutation-tested: a guard no mutation can break is dead code.
- REQ-337 depends on REQ-336 so the committed suite never carries a red check; its RED is proven by reintroducing the broken behaviour, not by committing against broken HEAD.
- REQ-336 and REQ-338 both write `web/board-timeline.js` — overlapping write sets, declared on both.

## Full Verbatim Input

do-work capture

Timeline (Gantt) view regression plus two audit leftovers. Board tool:
skills/do-work-board/tools/queue-kanban/.

--- REQ 1 — [impact-critical] Timeline clicks open the detail drawer again ---

RED (reproduce first, with real browser input — not synthetic PointerEvents):
generate a board, open the Timeline view, click a REQ bar with a real mouse
click. No drawer opens. Instrumented: pointerdown targets the <rect>, but the
click event targets DIV#timeline-scroll, so
clickEvent.target.closest("[data-detail-kind]") returns NONE and the delegated
handler at web/board-controls.js:183 never fires. Every mouse click inside
#timeline-scroll is swallowed, including a 2px sub-threshold nudge — so this
also breaks REQ-324's contract that a press which does not move is still a
click.

Cause: the setPointerCapture(downEvent.pointerId) call in the pointerdown
handler of web/board-timeline.js. Pointer capture retargets subsequent pointer
events AND the synthesized click to the capturing element. Introduced by
REQ-333, commit 36c4518, version 0.236.28.

Blast radius (already audited, do not re-widen): mouse clicks inside
#timeline-scroll only. Enter and Space on a focused row still open the drawer.
A mouse click on a Board-view card still opens it. Hover readout, ctrl+wheel
zoom over plot and axis, and view switching are all unaffected.

GREEN, all four must hold:
  1. plain click on a bar → drawer opens with that REQ/UR
  2. 2px nudge then release → drawer opens (REQ-324's contract)
  3. 150px drag → pans, drawer does NOT open
  4. drag released outside the chart → pan ends, no stuck pan state

Two candidates, both already validated per-gesture on fresh pages:
  B — call setPointerCapture only once the pan engages past
      TIMELINE_PAN_THRESHOLD_PX, not on pointerdown
  C — drop capture entirely, bind the pointerup/pointercancel release to window
Prefer B: it keeps REQ-333's release guarantee and makes "a drag is not a
click" an explicit property of the capture instead of an accident of node
rebuilding. Take C only if B turns out to leave a release path open.

--- REQ 2 — [impact-high] A check that can catch click retargeting ---

TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving
(timeline_browser_probe_test.go:626) is REQ-324's lock-in test. It passes
against the broken build. Reason: it dispatches synthetic PointerEvents with
pointerId: 1, an id the engine does not know, so setPointerCapture throws and
capture is never established — the probe never exercises the code path that
breaks real clicks.

Give the lane a check that fails on the current HEAD behaviour and passes after
REQ 1. Either drive real input, or assert the structural property directly
(where the capture call sits relative to the engage threshold). Mutation-test
whichever you pick: a guard no mutation can break is dead code.

--- Ask me before capturing these two ---

  Q1  The Timeline rows' keyboard contract: Tab escapes the row list, but only
      after 29 presses. I refuted it as a trap under REQ-333 and left it. Worth
      a REQ to cut the stop count, or leave it?
  Q2  The rightmost axis tick can name a day the "to" field excludes (a tick
      names an instant, the fields name days of coverage). Judged not a defect
      under REQ-334. Want it changed anyway?

---
*Captured: 2026-08-23T18:30:26Z*
