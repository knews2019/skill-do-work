---
id: UR-033
title: Reserve REQ numbers during allocation
created_at: 2026-08-07T19:15:15Z
requests: [REQ-147]
word_count: 25
---

# Reserve REQ Numbers During Allocation

## Summary

Make the queue-kanban allocator reserve every returned REQ number so another capture cannot receive the same identifier before the first request file is written.

## Full Verbatim Input

```text
the go app should reserve the numbers, so the next call gets a different id

fix this req and make a new req for it
```

---
*Captured: 2026-08-07T19:15:15Z*
