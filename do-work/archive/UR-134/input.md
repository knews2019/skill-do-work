---
id: UR-134
title: 'UR-133 leftovers: board intros name every page, and recover''s takeover message is accurate'
created_at: 2026-10-02T21:51:57Z
requests: [REQ-630, REQ-631]
word_count: 8
---
# UR-133 Leftovers: Board Intros Name Every Page, and Recover's Takeover Message Is Accurate

## Summary
After the UR-133 run (releases 0.305.63 and 0.305.64) the session listed three report-only leftovers and recommended fixing the first two; the user answered with the verbatim input below. "Board intros" and "the recover message" refer to:
1. Two intro lines still name only some of the board's pages: `skills/do-work-board/docs/board-guide.md` line 3 and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` line 3 (REQ-628 review finding M1).
2. `recover`'s takeover stop reason says the reset "requeues it as pending", but the reset writes `pending-answers` when Open Questions has an unchecked item and keeps `blocked` for a blocked claim (REQ-629 review finding M2).
The third leftover (renaming the `RECOVERY-TAKEOVER-AVAILABLE` code) was recommended against and is not requested.

## Extracted Requests
- R1 → REQ-630 (Board intros describe the board by what its switcher covers)
- R2 → REQ-631 (Recover's takeover message says where a reset claim goes)

## Batch Constraints
- Independent; either order.
- Both change shipped files, so each is a release.

## Full Verbatim Input
> ```
> ok, do board intros and the recover message
> ```
