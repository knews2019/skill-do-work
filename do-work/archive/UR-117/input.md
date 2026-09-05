---
id: UR-117
title: 'Hide the verify-findings strip on the Activity view'
created_at: 2026-09-04T23:58:59Z
requests: [REQ-578]
word_count: 42
---

# Hide the Verify-Findings Strip on the Activity View

## Summary

Looking at the Activity view after REQ-572 (one row per lifecycle stamp) landed, the user asked for the Verify Findings strip to be removed from that view, and repeated that the REQ id must be clickable there. The clickable-row request is already REQ-573 (open the detail drawer from an Activity row and highlight every row of the same REQ), which is queued behind REQ-572; no new REQ is minted for it.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-578 | The Verify Findings strip is not shown while the Activity view is active |

## Folded Requests

- REQ-573 — "I asked for the REQ to be clickable": already captured as opening the detail drawer from an Activity row and highlighting its sibling rows

## Full Verbatim Input

> ```
> [Screenshot 3: the board's Activity view after REQ-572 landed, showing "175 transitions across 38 REQs in the last 24 hours" under a two-card Verify Findings strip]
> 
> <- remove verify finding from this view, also I asked for the REQ to be clickable
> ```
