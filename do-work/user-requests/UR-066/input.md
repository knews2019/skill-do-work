---
id: UR-066
title: Audit every state of the Timeline Gantt chart
created_at: 2026-08-23T11:12:00Z
requests: [REQ-326, REQ-327, REQ-328, REQ-329, REQ-330, REQ-331, REQ-332, REQ-333, REQ-334, REQ-335]
word_count: 53
---

# Audit Every State of the Timeline Gantt Chart

## Verbatim Input

> [Screenshot of the board's Timeline view at http://127.0.0.1:8090 — project `g1w-game-find-the-difference`,
> generated 2026-08-23 11:08 UTC. The Week chip is lit. The From/to fields read 06/01/2026 → 06/07/2026 and the
> readout `2026-06-01 00:00 UTC → 2026-06-08 00:00 UTC`. The summary reads "Nothing was drawn between
> 2026-06-01 00:00 UTC and 2026-06-08 00:00 UTC. Widen the window, step to another period, or press Fit all —
> 795 REQs are outside it." The axis reads `1 Jun | 2 Jun | 3 Jun | 4 Jun | 5 Jun | 6 Jun | 8 Jun` — 7 Jun is
> absent. The chart area below the axis is empty.]
>
> audit every state of this gantt chart and make it make sense.
> Let me start when I click on week, every entry disappears because it goes far in the past instead of be
> sensible and show me the now.

Second message, same session:

> you are on ultracode, capture requests, verify them, run them, first ask me things if you really need, but
> make sure to use common sense.

## Assets

The pasted screenshot was inline in the session and has no file on disk. The same states were reproduced
locally against this repo's own queue (317 REQs, now 2026-08-23 11:13 UTC) and captured under the session
scratchpad; the reproduction commands and the observed values are recorded in each REQ's Red-Green Proof.

## Notes on Scope

"Audit every state" is the whole request; the Week-chip complaint is the example the user led with, not the
boundary. Findings were produced by eight parallel state-cluster audits of `board-timeline.js` driven against
a real Chromium render, each finding then adversarially verified before it was allowed to become a REQ.
