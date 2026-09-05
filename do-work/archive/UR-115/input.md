---
id: UR-115
title: 'Show every lifecycle transition of a REQ on the Activity view, open its detail on click, and highlight its sibling rows'
created_at: 2026-09-04T23:16:00Z
requests: [REQ-572, REQ-573]
word_count: 154
---

# Show Every Lifecycle Transition of a REQ on the Activity View

## Summary

The Activity view keeps one row per REQ, its newest lifecycle stamp only. The user asked to see the whole path a REQ took (captured, claimed, dispatched, merged, reviewed, released, completed) as rows on the same surface, to open the same detail drawer the Board opens when a row is clicked, and to highlight every row that belongs to the clicked REQ so its history can be scanned by eye.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-572 | Emit one Activity row per lifecycle stamp instead of one per REQ, so every state a REQ went through appears in the window |
| REQ-573 | Clicking an Activity row opens the REQ detail drawer (as on the Board) and highlights every row carrying the same REQ id |

## Batch Constraints

- Go decides which stamps exist and what each records (`lifecycleTimestampFields` in `model.go`); the client draws and filters. No second definition of a stamp's meaning in JavaScript.
- The existing window, filter chips, and empty-state messages keep working; only the row set and row behavior change.
- REQ-573 depends on REQ-572: highlighting sibling rows only means something once a REQ can have several rows.

## Full Verbatim Input

> ```
> [Screenshot 1: the board's Activity view, filtered by the browser's find box on "570"]
> 
> <- is this only showing the last status of a REQ? how about if I want to see when it went through all the states of it?
> 
> [Agent answer, summarized: yes, the Activity view keeps only the newest lifecycle stamp per REQ and no board surface shows a REQ's full state history; the cheapest change is to emit one Activity row per stamp instead of one per ticket, with a toggle and grouping by REQ.]
> 
> ok, do-work capture-request, and also make sure that if I click on a req it does show up on the left side just as it shows up in the board [Screenshot 2: the Board view with the detail drawer open for REQ-570] and also highlights all similar REQ-ID entries, so it can be visually scanned all the status the REQ went through to arrive there.
> ```
