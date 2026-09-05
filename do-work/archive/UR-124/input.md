---
id: UR-124
title: 'Verify findings rows read as a paragraph, pick a row layout from mock-ups'
created_at: 2026-09-05T14:45:41Z
requests: [REQ-588]
word_count: 141
---

# Verify Findings Rows Read as a Paragraph, Pick a Row Layout From Mock-ups

## Summary

The user sent a board screenshot and said the Verify Findings strip's styling is broken, asking which REQ should fix it. No pending REQ covers the strip: REQ-579 (render verify findings and skipped probes as compact rows in one list) shipped the current rows in commit b169396e and is completed, and REQ-580 (stop the redundant committed-queue-state probe) is completed too. The user then asked to capture it and fix it, and to make a do-work ai-report with mock-ups so they have options to choose from. This UR captures one addendum REQ against REQ-579; the mock-up report is produced in the same session and the chosen layout is recorded as the REQ's answer.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-588 | Make each verify-finding row read as one warning line: remedy visually separated from the detail, chip and text on one grid, one type scale across subject, chip and row; layout picked by the user from the mock-up report |

## Batch Constraints

- REQ-579's approved design (one list, weight only from `fixable`, grouping by producer subject, hide rules unchanged) stays as is. This addendum changes how a row is laid out, not what the list contains.
- The strip is not a REQ card: no `.board-request*` classes, no `board-anomalies-cards` grid.

## Full Verbatim Input

> ```
> [Screenshot 1: the queue-kanban board at 127.0.0.1:8090, Board view, light theme, 14:25 UTC on 2026-09-05. Under the top bar, the Verify Findings strip reads "VERIFY FINDINGS 2 findings queue and process problems queue-kanban verify detects — each names what to do about it". Below it two subject headings, worktree-agent-REQ-573-activity-drawer and worktree-agent-REQ-582-arrow-citations, each followed by one row: an uppercase chip (UNMERGED-WORKTREE-LEFTOVER, WORKTREE-PRESENT-RUN-IN-FLIGHT) and then the detail sentence, an arrow, and the long remedy sentence flowing together as one grey paragraph that wraps to a second line with one orphaned word ("copy", "archived, not before"). The four kanban columns (Pending 9, Claimed 1, Needs input 0, Recently done 41) sit below the strip.]
> 
> <- verify findings styling is broken, which req should fix it?
> 
> capture it and fix it
> 
> also it's make do-work ai-report with mock-ups so I have options to choose from
> ```
