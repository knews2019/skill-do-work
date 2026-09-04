---
id: UR-116
title: 'Append-only lifecycle stamps and board wall time from the first stamp'
created_at: 2026-09-04T23:52:00Z
requests: [REQ-575, REQ-576]
word_count: 148
---

# Append-Only Lifecycle Stamps and Board Wall Time From the First Stamp

## Summary

The board card for REQ-505 (moving selection and claim behind `advance`) said "wall time 1m 23s" for a request that was worked for over six hours. The card measures `completed_at` minus `claimed_at`. REQ-505 was claimed at 16:39 UTC, built until 17:26, held for heavy testing, and at 21:02 the green heavy result moved it back to `pending` and deleted the `claimed_at` line. The re-claim at 23:00 wrote a fresh stamp and the completion at 23:01 closed the span. The user asked whether a per-REQ timing log file would give proper timings. The agent answered that the frontmatter already is the timing log and one transition deleted one stamp, and proposed two smaller changes instead. The user asked to capture both.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-575 | No lifecycle transition may delete or overwrite an existing `*_at` frontmatter stamp; the two code paths that still do stop doing it |
| REQ-576 | The board card's wall time starts at the earliest lifecycle stamp the REQ carries, not only at `claimed_at` |

## Batch Constraints

- No new timing file, stream, or writer. REQ-562 (recording lightweight per-REQ lifecycle timings) owns command-level attribution; these two REQs only stop losing the stamps already recorded and read them correctly.
- REQ-570 (deleting the pending-heavy-testing status) removes the transition that damaged REQ-505. REQ-575 is the general guard so the next hold or requeue path cannot repeat it.

## Full Verbatim Input

> ```
> should we have a timing log file for each req so we can get the proper timings?
> 
> capture a req for append-only stamps and the board wall time change
> 
> [The two changes the user is pointing at, proposed by the agent earlier in the same conversation after tracing why REQ-505 shows "wall time 1m 23s":]
> Make stamps append-only. No lifecycle transition may delete or overwrite a `*_at` field once it exists. That is one rule in the state-apply code plus one test, and it makes a repeat of this impossible even if a new hold state appears later.
> Let the board measure from the first stamp, not from `claimed_at`. Wall time from the earliest phase stamp to `completed_at` would have shown about 6h 12m for REQ-505 even with the damaged claim. The drawer already has the phase list, so this is a change in what the card summarizes.
> ```

---
*Captured: 2026-09-04T23:52:00Z*
