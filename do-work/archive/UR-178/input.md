---
id: UR-178
title: 'Build the proposal as a table and compare it with the map'
created_at: 2026-11-14T23:21:51Z
requests: [REQ-734]
word_count: 61
---
## Summary
The maintainer quoted the assistant's doubt that a map may be the wrong form for the Tinka proposal, and asked to build the table form too and to compare the two.

## Extracted Requests

| REQ | Request |
| --- | --- |
| REQ-734 | Build a table page with one row per option and the columns "supports" and "limits", then show it side by side with the improved map page. |

## Folded Requests

- REQ-727 (ship a relation-map preset for ai-report second-opinion) — the preset now waits for the comparison, because the maintainer picks map or table after seeing both

## Batch Constraints
- REQ-734 runs after REQ-733 (rebuild the relation map page).
- Nothing under `skills/` changes.

## Full Verbatim Input
> ```
> build this one too: "For this report, a map may be the wrong form. The content is one decision with reasons for and against. A table with one row per option and columns "supports" and "limits" would show the same relations with no zooming." <- basically compare the two, so I can see which one is better to tell the story
> ```
