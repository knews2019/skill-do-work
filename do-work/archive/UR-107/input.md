---
id: UR-107
title: 'A priority field the selector orders by and the board shows'
created_at: 2026-09-03T20:38:55Z
requests: [REQ-561]
word_count: 90
---

# A Priority Field the Selector Orders By and the Board Shows

## Summary

The maintainer wants a priority list for the pending queue that the kanban board reflects. Nothing carries a rank today: the schema has no priority field, the selector orders ready work only by three classes (gate repairs, deferred parents, everything else in queue order), and the board sorts nothing by rank. One REQ adds a closed three-value `priority` field, makes the selector honour it inside the ordinary class below `depends_on`, makes the board sort the pending column by it and tag the card, and stamps the current queue per the 23:20 triage table in the velocity report. REQ-530 (order ready work by the newest REQ it unblocks) is the same need in a narrower form and is superseded; it is cancelled with the landing hash, not folded, because its rule would become the tie-break inside a priority class rather than the order itself.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-561 | Add `priority: now | next | later` to the REQ schema; the selector orders ready work by it inside the ordinary class, the board sorts pending by it and shows a tag; stamp the current queue per the triage table |

## Batch Constraints

- Parser and board move in lock-step (prime-kanban-board.md); the schema normalizer and the selector are one commit with the board change so a stamped queue is never malformed to either reader.
- `depends_on` stays the hard order; priority is soft and never lets a dependent run before its dependency.
- One field, three values, absent reads as `next`; no numeric scale, no per-UR priority, no second selector flag.

## Full Verbatim Input

> ```
> can we make a priority list that would reflect in the kanban board as well, or is that a new req?
> 
> [assistant: it is a new REQ; proposed one frontmatter field `priority` with three values (now, next, later), read by the selector within the ordinary class and by the board as a sort key and card tag, subsuming REQ-530; offered to capture it as UR-107 with the triage table's build-now set marked now and the deferred set later]
> 
> capture UR-107 and tell me what todo with the current claimed REQs
> ```

---
*Captured: 2026-09-03T20:38:55Z*
