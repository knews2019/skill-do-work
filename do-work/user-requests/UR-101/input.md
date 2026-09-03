---
id: UR-101
title: 'Pick the newest REQ first, not the oldest'
created_at: 2026-09-03T10:59:01Z
requests: [REQ-530]
word_count: 39
---

# Pick the Newest REQ First, Not the Oldest

## Full Verbatim Input

> ```
> [maintainer]
> do-work capture-request when picking up the next REQuest it should be the latest LIFO not the oldest REQuest. Of course don't break the dependencies. The principle is that the one that I just captured I'm more interested in it.
> 
> [capture-time answer, 2026-09-03]
> Q1 When the newest REQ depends on an older REQ that is still pending, which runs first? -> Pull the prerequisite forward: ready work is ordered by the newest REQ it unblocks.
> ```

---
*Captured: 2026-09-03T10:59:01Z*
