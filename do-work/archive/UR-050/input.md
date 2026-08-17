---
id: UR-050
title: Durations view on the Kanban board
created_at: 2026-08-17T17:17:22Z
requests: [REQ-219]
word_count: 51
---

# Durations View on the Kanban Board

## Summary

Add the duration/trend/cadence chart — built and design-validated in the UR-048 AI report — to the Go-served Kanban board as a new dedicated view, sourcing its data from the archive scan the board already performs.

## Extracted Requests

| REQ | Title | Depends on |
|-----|-------|-----------|
| REQ-219 | Durations view on the Kanban board | — |

## Batch Constraints

Resolved with the user during capture (three questions, all answered with the recommended option):

- **Data source: scan the archive at build time.** `queue-kanban` already parses every archived REQ's `claimed_at`/`completed_at`, so the 195 samples available today cost no new file reads and no new parser fields. Full history immediately, and it stays correct when a stamp is later repaired. Explicitly rejected: sourcing from `do-work/calibration-log.tsv`, which holds 8 rows today, only grows as new REQs archive, and is append-only so a repaired stamp never propagates.
- **Placement: a new dedicated view**, alongside the existing board and Testing views. The Kanban columns stay untouched, and the shared calendar axis gets full width.
- **Panels: all three**, as validated in the report — duration per REQ, median per active day, REQs completed per day.
- **Read-only.** The board has exactly three write surfaces (CLAUDE.md § Kanban Board Write Surfaces). This view adds none, and that sentence must not need amending.
- Board versioning is folded into the skill: normal CHANGELOG entry plus suite version bump, per `_dev/primes/prime-kanban-board.md`.

## Design Source

The chart already exists and was render-judged in light and dark at
`ai-reports/2026-08-17_1401_UR-048-estimator-calibration-and-anomaly-surfacing/index.html`
(section 8, commit `6e79932`). It is the reference implementation: same three panels,
same shared axis, same validated route ramp, same outlier handling. The board port
should reuse those decisions rather than re-derive them.

## Capture Note — UR Number

This capture first wrote itself as UR-049 and lost that number to a concurrent capture
(the codeload-429 batch, REQ-216 through REQ-218) that claimed UR-049 during the write.
Renumbered to UR-050; the other batch was left untouched. REQ-219 was reserved through
`queue-kanban next-req`, so the REQ id never collided.

## Full Verbatim Input

add a duration chart to the ai-reports/2026-08-17_1401_UR-048-estimator-calibration-and-anomaly-surfacing/index.html because that way I can identify if there is an overtime response time degradation, outliers, how often are these reqs executed

and also add capture-request prompt to have such chart added to the go served kanban board

[Capture-session context: the first message produced the chart in the AI report; this UR
covers only the second — porting it to the board. The user's stated purpose from the first
message carries over verbatim and is the acceptance frame for the port: detect over-time
response-time degradation, see outliers, and see how often REQs are executed.]

---
*Captured: 2026-08-17T17:17:22Z*
