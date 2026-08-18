---
id: UR-052
title: Timeline needs horizontal navigation and a jump to remaining work
created_at: 2026-08-18T10:22:05Z
requests: [REQ-235]
word_count: 22
---

# Timeline Needs Horizontal Navigation and a Jump to Remaining Work

## Full Verbatim Input

do-work capture-request [screenshot of the board's Timeline view] <- this is not working well, can not scroll horizontally, I can not jump to the remaining work

## Screenshot Description

The attached screenshot shows the queue-kanban board at `http://127.0.0.1:8090`, Timeline tab, for the project `g1w-game-find-the-difference`, generated 2026-08-18 09:54 UTC. Header: "How long each REQ waited, and how long it took — 677 REQs in capture order, oldest at the top. 97 still open, measured to the now-line at 2026-08-18 09:54 UTC." Legend: Waiting to be claimed / Being worked on / Still open, running to now / Projected, not measured / Broken stamps. Toolbar at the right of the chart: `−`, `+`, `Now`, `Fit all`.

The axis runs 16 Apr → 23 Jul across the full width. Every drawn bar (rows REQ-1020 through REQ-1063) is crushed into a narrow band at the far right of the plot — roughly the last 12% of the width — while the left 88% is empty. Most bars are one to three pixels wide and unreadable; a few grey waiting bars extend to the right edge and are cut off there.

Below the chart: "Queue empties around 2026-08-18 11:55 UTC — 6 REQs run one at a time from 2026-08-18 10:06 UTC…", then the hint line "Scroll to move down the queue. Hold ⌘ or Ctrl and scroll to zoom the time axis; drag to pan. Click a row for its full detail.", the hovered-row readout for REQ-1035, and the "Every REQ, as a table" disclosure listing REQ / Title / Status / Waited / Worked.

## Clarifications Answered During Capture

1. **What should "jump to the remaining work" land on?** → *Now + the forecast*: a control that moves the time window to the now-line and the projected queue-empty time, and also scrolls the row list there — the existing `Now` button moves only the window.
2. **How should sideways scrolling move the time axis?** → User redirected the question: *"consider also the performance if the time is long, so perhaps jump previous and next period using the current time zoom. Time zoom should be daily/weekly/monthly"* — i.e. discrete day/week/month zoom levels plus previous/next period stepping at the chosen level, chosen partly because a long time range is expensive to draw.
3. **Do day/week/month replace the current free zoom, or sit alongside it?** → *Alongside it*: keep the existing `−`/`+`, ⌘-scroll zoom and drag-pan exactly as they are; add day/week/month as presets plus prev/next period stepping.
4. **Should a picked period sit on calendar boundaries?** → *Calendar-aligned*: day = midnight to midnight UTC, week = Mon–Sun, month = the 1st to the last day, so prev/next always lands on a clean boundary.

---
*Captured: 2026-08-18T10:22:05Z*
