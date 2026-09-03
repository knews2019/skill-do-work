---
id: UR-101
title: 'Pick the newest REQ first, not the oldest'
created_at: 2026-09-03T10:59:01Z
requests: [REQ-530]
word_count: 105
---

# Pick the Newest REQ First, Not the Oldest

## Full Verbatim Input

> ```
> [maintainer]
> do-work capture-request when picking up the next REQuest it should be the latest LIFO not the oldest REQuest. Of course don't break the dependencies. The principle is that the one that I just captured I'm more interested in it.
> 
> [capture-time answer, 2026-09-03]
> Q1 When the newest REQ depends on an older REQ that is still pending, which runs first? -> Pull the prerequisite forward: ready work is ordered by the newest REQ it unblocks.
> 
> [addendum, 2026-09-03]
> update ur-101 so that the kanban report also is in sync with the choosen task ordering
> BTW: is there a go tool that is doing the picking of the next task? I imagine there should be one to do that and the LLM should use that.
> also the kanban board should use that so that the go tool is the point of truth in the ordering
> ```

---
*Captured: 2026-09-03T10:59:01Z*
