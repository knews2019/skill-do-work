---
id: UR-112
title: 'Board window for recently touched REQs regardless of status'
created_at: 2026-09-04T17:54:05Z
requests: [REQ-568]
word_count: 5
---

# Board Window for Recently Touched REQs Regardless of Status

## Summary

The maintainer asked how to see recent activity on the Kanban board, then noticed a gap: the claimed card (REQ-505, moving selection and claim behind advance) was 20 minutes old while the newest Recently done card (REQ-485, canonicalizing reservation marker filenames) was two hours old. Git history showed the gap was full: REQ-567, REQ-503, and REQ-504 were each claimed, built, merged, and held as `pending-heavy-testing` between 14:57 and 16:39 UTC, and a release 0.275.3 shipped at 15:54. None of that appears in any "recent" surface, because Recently done only lists terminal states and the Timeline plots spans, not the hold event.

The assistant proposed: "If you want a board surface for that, a 'recently touched' window keyed on the newest stamp on a ticket, `updated_at` or the hold time, would be the shape. Say so and I will capture it as a REQ." The verbatim input below is the maintainer's acceptance of that proposal.

## Full Verbatim Input

> ```
> capture that as a REQ
> ```

---
*Captured: 2026-09-04T17:54:05Z*
