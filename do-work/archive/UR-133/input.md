---
id: UR-133
title: 'UR-132 follow-ups: board guide update, and handoff resume without resetting claims'
created_at: 2026-10-02T19:49:03Z
requests: [REQ-628, REQ-629]
word_count: 4
---
# UR-132 Follow-ups: Board Guide Update, and Handoff Resume Without Resetting Claims

## Summary
At the end of the UR-132 run (releases 0.305.61 and 0.305.62) the session offered to capture two report-only items; the user answered with the verbatim input below. "Both" refers to:
1. The board guide lists only three of the six pages, does not mention the new page links, and does not mention the Testing page's disk line (REQ-626 review finding F3, REQ-627 review finding M1).
2. The handoff action says to continue claimed REQs, but `recover --take-over`, the next step `recover` offers for them, resets them to the queue and strips their sections.

## Extracted Requests
- R1 → REQ-628 (Board guide names all six pages, their links, and the disk line)
- R2 → REQ-629 (Resuming a handoff must not reset its own claimed REQs)

## Batch Constraints
- Independent; either order.
- Both change shipped files, so each is a release.

## Full Verbatim Input
> ```
> capture both as REQs
> ```
