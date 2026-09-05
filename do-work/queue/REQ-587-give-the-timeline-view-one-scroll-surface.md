---
id: REQ-587
title: 'Give the Timeline view one scroll surface, in the same style as the Activity view'
status: pending
created_at: 2026-09-05T12:50:59Z
user_request: UR-123
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: [REQ-585]
related: [REQ-585, REQ-586]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
---

# Give the Timeline View One Scroll Surface, in the Same Style as the Activity View

## What

The Timeline view scrolls twice, like the Activity view did before REQ-585 (give the Activity view one scroll surface): the chart rows live in `.timeline-scroll`, a box capped at `max-height: 58vh` with `overflow-y: auto`, inside `.board-main`, which is the board's scroll container. Leave exactly one scroll surface on this view, in the same style as the Activity fix: the board scrolls, the rows are content, and what the reader needs to keep in view (the time axis, as the column header is on Activity) sticks to the top edge.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

This is not the CSS-only change REQ-585 was. The Timeline rows are virtualized: `board-timeline.js` reads `scrollHost.scrollTop` and `scrollHost.clientHeight` (the element with id `timeline-scroll`) to decide which rows carry nodes, to keep the same REQ under the pointer when the window changes, to scroll a focused row into view, and to jump to the open work (roughly lines 1618, 1757 to 1797, 2077 to 2117, 2803 to 2806, 3186). The block comment above the timeline CSS (board.css near line 2579) states that `.timeline-scroll` is the scroll container by design. Two ways to get one scroll surface:

- **Re-point the scroll host to `.board-main`.** The chart becomes full height as ordinary content; the visible-row math reads the board's scroll position minus the chart's offset inside it; the axis (`.timeline-axis`) becomes sticky under the board's top edge. This is the same reading as REQ-585's M1 and is the recommended shape. Wheel-to-zoom (Ctrl or Cmd plus wheel) and drag-to-pan must keep working on the chart.
- **Let the chart fill the remaining viewport and stop the board scrolling on this view** (REQ-585's M2 shape). Cheaper, but the forecast line, the excluded list, the hint paragraph and the "REQs in this window, as a table" details sit below the chart today and would need a new home, and the Timeline would have a scroll model no other view has.

The Activity mockups and the side-by-side comparison of the two shapes are in `ai-reports/2026-09-05_1520_activity-view-double-scroll-mockups/index.html`. REQ-585's merged diff shows how the padding move was scoped to one view; reuse that mechanism rather than inventing a second one. Keyboard behaviour named in the view's hint paragraph (Tab to the chart, arrows move between rows and pan, plus and minus zoom) must survive.

## Red-Green Proof
**RED prompt/case:** On the served board, open the Timeline with the All days window and run in the console: `const m=document.querySelector('.board-main'), t=document.getElementById('timeline-scroll'); [m.scrollHeight>m.clientHeight, t.scrollHeight>t.clientHeight]`.
**Why RED now:** It returns `[true, true]` and two scrollbars are visible (screenshot under Assets: the chart's own scrollbar at the right edge of the rows, the board's behind it).
**GREEN when:** Exactly one element scrolls, the axis stays visible while scrolling through the rows, and the rows still render correctly at every scroll position (the virtualized range follows the scroll: a row 400 px below the fold appears after scrolling there, and the row under the pointer stays put when the window chips change). Measured in a real engine and recorded in the hand-back; a browser probe in the existing `*_browser_probe_test.go` lane pins it.
**Validation:** Inferred during capture

## Open Questions

- [ ] Which shape: re-point the scroll host to the board (recommended; same style as the Activity fix, axis sticky) or let the chart fill the remaining viewport with the board not scrolling on this view. **Recommended:** re-point the scroll host to the board. The builder proceeds with that if this is still open at claim time.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`, so bare only): matches on "Changing queue-kanban UI, timeline, or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views or timeline behavior". Over the budget on its own.

## Assets

- `do-work/user-requests/UR-123/assets/REQ-587-screenshot-1-timeline-double-scroll.png`: the live board at 127.0.0.1:8090, Timeline view, All days window, generated 12:38 UTC on 2026-09-05. Above the chart, the verify-findings strip with five worktree cards, the "How long each REQ waited, and how long it took" heading, the legend, and the period controls (Last day to All days, From 05/27/2026 to 09/07/2026). The chart lists REQs grouped by UR (UR-120 down to UR-114) with bars at the right edge near the now-line. Two scrollbars on the right: the chart box's own, starting at the axis, and the board's thin one behind it starting under the top bar.

*Source: "BTW: this view also have a scroll issue, add a REQ for it as well: - the style should be similar"*
