---
id: REQ-336
title: "[impact-critical] Timeline clicks open the detail drawer again"
status: pending
created_at: 2026-08-23T18:30:26Z
user_request: UR-067
addendum_to: REQ-333
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-337, REQ-338]
batch: timeline-click-regression
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
---

# Timeline Clicks Open the Detail Drawer Again

## What

Mouse clicks on Timeline (Gantt) view REQ bars no longer open the detail drawer — every mouse
click inside `#timeline-scroll` is swallowed. Restore clicking while keeping the pan behaviour
REQ-333 and REQ-324 established.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

A primary interaction of a shipped view is broken: no mouse path from the Timeline to a REQ/UR's
details. It also breaks REQ-324's contract that a press which does not move is still a click.
User-judged `impact-critical` (broken production path).

## Context

**Cause (user-diagnosed):** the `setPointerCapture(downEvent.pointerId)` call in the pointerdown
handler of `web/board-timeline.js` (currently at line ~2580). Pointer capture retargets subsequent
pointer events AND the synthesized click to the capturing element (`#timeline-scroll`), so
`clickEvent.target.closest("[data-detail-kind]")` in the delegated click handler at
`web/board-controls.js:183` finds nothing and the drawer never opens. Introduced by REQ-333,
commit 36c4518, version 0.236.28.

**Blast radius (already audited — do not re-widen):** mouse clicks inside `#timeline-scroll`
only. Enter and Space on a focused row still open the drawer. A mouse click on a Board-view card
still opens it. Hover readout, ctrl+wheel zoom over plot and axis, and view switching are all
unaffected.

## Detailed Requirements

Reproduce RED first, with real browser input — not synthetic PointerEvents (the user's diagnosis
came from instrumented real input; synthetic PointerEvents with an engine-unknown `pointerId`
make `setPointerCapture` throw and never establish capture, so they cannot reproduce this).

Two candidate fixes, both already validated per-gesture on fresh pages by the user:

- **B** — call `setPointerCapture` only once the pan engages past `TIMELINE_PAN_THRESHOLD_PX`,
  not on pointerdown
- **C** — drop capture entirely, bind the pointerup/pointercancel release to window

**Prefer B**: it keeps REQ-333's release guarantee and makes "a drag is not a click" an explicit
property of the capture instead of an accident of node rebuilding. Take C only if B turns out to
leave a release path open.

## Builder Guidance

Certainty: Firm — cause, blast radius, and preferred fix are user-diagnosed and pre-validated.
`tdd: false` because the existing probe lane cannot drive the real input this RED requires
(REQ-333's UNIFY notes: a synthetic `pointerId` cannot be captured in this lane); the durable
automated check is REQ-337's deliverable — do not duplicate it here. Prove RED/GREEN with real
browser input against a generated board; what counts as real input in the probe lane is
REQ-337's call to make.

## Red-Green Proof

**RED prompt/case:** Generate a board, open the Timeline view, click a REQ bar with a real mouse
click. No drawer opens. Instrumented: pointerdown targets the `<rect>`, but the click event
targets `DIV#timeline-scroll`, so `closest("[data-detail-kind]")` returns NONE and the delegated
handler at `web/board-controls.js:183` never fires. A 2px sub-threshold nudge is swallowed too.
**Why RED now:** `setPointerCapture` on pointerdown retargets the synthesized click to
`#timeline-scroll` for every press, drag or not.
**GREEN when:** All four hold, with real browser input:
1. plain click on a bar → drawer opens with that REQ/UR
2. 2px nudge then release → drawer opens (REQ-324's contract)
3. 150px drag → pans, drawer does NOT open
4. drag released outside the chart → pan ends, no stuck pan state
**Validation:** User confirmed (RED, GREEN, cause, and candidates all stated verbatim in the input)

## Prior Implementation

REQ-333 ("Keep the timeline pointer and keyboard paths alive", archived, `completed-with-issues`,
commit 36c4518) added the `setPointerCapture` call to make a drag released outside the chart a
guaranteed release rather than a hope, fixing the latched `.is-panning` state. It wrote
`web/board-timeline.js` and `timeline_browser_probe_test.go`; the capture call is feature-guarded
and degrades to the ordinary release path when unavailable. That release guarantee is GREEN
condition 4 — keep it.

## Dependencies

None to start. REQ-337 (the lock-in check) depends on this REQ; REQ-338 also writes
`web/board-timeline.js`, declared in both write sets.

---
*Source: UR-067 — see `do-work/user-requests/UR-067/input.md` for complete verbatim input.*
