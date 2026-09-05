---
id: UR-118
title: 'Verify findings as compact rows in one list'
created_at: 2026-09-05T00:19:58Z
requests: [REQ-579, REQ-580]
word_count: 101
---

# Verify Findings as Compact Rows in One List

## Summary

Looking at the Verify Findings strip on the board, the user asked whether the two finding cards and the two "probe could not run" lines under them were the same information. They were two different checks with one root cause (the worktree's branch is gone). The user then said that for a warning the small one-line form is the one they like, that the cards take too much space for something not that important, and asked how to make the strip good, beautiful and non-contradictory. The assistant proposed five points (D1 one list of rows, D2 two weights taken only from the producer's `fixable` flag, D3 rows grouped by subject, D4 stop the redundant probe at the source, D5 keep the hiding rules as they are) plus cancelling REQ-482 (stack verify-finding cards full width), whose body asked for wider cards. The user answered "ok, do it". REQ-482 was cancelled in the same session with a reason pointing here.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-579 | Render every verify finding and every skipped probe as a compact row in one list, weighted only by the producer's `fixable` flag, grouped by the subject the producer names (D1, D2, D3, D5) |
| REQ-580 | When a leftover worktree's merge state is already undetermined, do not run the committed-queue-state probe for it; fold the "not checked" fact into the undetermined finding instead (D4) |

## Batch Constraints

- The strip keeps its current hide rules: hidden when there is nothing to say, and hidden on the Activity view once REQ-578 (hide the verify-findings strip on the Activity view) lands. Nothing else conditional.
- No severity scale is invented in the client. The producer knows `fixable` or not; that is the only weight.
- A skipped probe must still never read as "checked and clean" (the rule that placed the disclosure under the strip in the first place). REQ-580 removes a redundant line only by moving its fact into the finding that already covers the same worktree.

## Full Verbatim Input

> ```
> [Screenshot 1: the board's Verify Findings strip showing two WORKTREE-MERGE-STATE-UNDETERMINED cards, one for the REQ-506 worktree and one for the REQ-577 worktree, above a "2 probe(s) could not run — unverified, not clean" disclosure listing the committed-queue-state probe for each of the same two worktrees]
> 
> <- are these finding duplicated? the big box and the small line is the same info?
> 
> for a warning I like the small lines, the big boxes are using up too much space for something not that important, is there a req talking about this?
> 
> how to make it good, beautiful and non-contradictory?
> 
> ok, do it
> ```
