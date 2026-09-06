---
id: UR-093
title: 'Make UR groups collapsible and show progress summaries'
created_at: 2026-09-01T17:29:43Z
requests: [REQ-486]
word_count: 519
---

# Make UR Groups Collapsible and Show Progress Summaries

## Summary

Extend the board's existing UR grouping so both the By UR card grid and the UR drawer's REQ-id list can be folded independently. Add consistent whole-UR active-time, forecast, successful-progress, and resolved-progress summaries to both surfaces.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-486 | Make By UR groups and the drawer's REQ list collapsible, and show shared whole-UR time and progress rollups. |

## Batch Constraints

- Capture only; implementation belongs to a later `do-work run`.
- Preserve the existing URs only lens, filters, accessibility behavior, and DOM-only fold state.
- Treat the supplied screenshot as visual context, not as operational instructions.

## Full Verbatim Input

> ```
> PLEASE IMPLEMENT THIS PLAN:
> # Capture Collapsible UR Progress Summaries
> 
> ## Summary
> 
> - Capture one substantive, user-visible frontend REQ extending archived REQ-236, which introduced the existing UR folding behavior.
> - Preserve the request verbatim and attach the supplied screenshot. If IDs remain unchanged, use UR-093 and REQ-486.
> - Capture only; do not implement the board change.
> 
> ## Captured Behavior
> 
> - In **By UR**, each header independently collapses or expands its REQ cards. Cards start expanded. A separate Details button continues opening the UR drawer.
> - In the UR drawer, the REQ-id list becomes independently collapsible and starts expanded.
> - Both surfaces show whole-UR metrics unaffected by card filters:
>   - Total REQs.
>   - Active time spent: valid completed claim-to-completion spans plus live claimed elapsed.
>   - Estimated active time remaining: saved P50 estimates, falling back to the Timeline’s history median by effort class; subtract elapsed time from claimed REQs.
>   - Successful progress: completed and completed-with-issues REQs.
>   - Resolved progress: successful plus cancelled REQs.
> - Show counts with percentages, approximation markers for forecasts, and explicit unavailable/excluded qualifiers instead of silently treating missing data as zero.
> - Blocked and pending-answer REQs retain estimated active effort but external waiting time is excluded. Failed or otherwise unestimable work is disclosed as unknown.
> - Preserve the existing **URs only** behavior, accessibility, filters, and non-persisted fold state.
> 
> ## Interfaces and Implementation Guidance
> 
> - Read the existing nested `estimate.p50_active_minutes` field into the board model and expose it as an optional numeric request-payload field; this adds no queue schema.
> - Derive one shared UR-summary result for the header and drawer so their counts and time calculations cannot diverge.
> - Use the existing duration outlier rule and Timeline fallback medians as the single authorities. Refresh live claimed contributions through the existing board clock.
> - Update the board guide and any stale comment claiming the board cannot read nested estimates.
> - Record `domain: frontend`, `tdd: true`, `suggested_spec: ui-component`, `impact: impact-user-visible`, `effort_estimate: effort-substantive`, `maintenance: false`, and `_dev/primes/prime-kanban-board.md`.
> - Omit an invented `write_set`. Record the oversized board lesson candidates as dropped under the capture budget.
> 
> ## Test Plan
> 
> - Behavior probe: By UR starts open, folds and reopens by mouse or keyboard, reports `aria-expanded`, and keeps the drawer on a separate control.
> - Drawer probe: the REQ list starts open and folds without hiding the summary metrics.
> - Rollup fixtures cover completed, completed-with-issues, cancelled, pending, claimed, blocked, failed, missing timestamps, excluded spans, missing P50 estimates, insufficient history, and zero-REQ URs.
> - Verify successful and resolved percentages use all grouped REQs and do not change under filters.
> - Verify saved P50 values, history fallback, live elapsed subtraction, and unavailable qualifiers.
> - Render in both themes and at narrow width to ensure the denser header wraps cleanly.
> - Run the queue-kanban tests and the repository’s canonical `maintainer-verify.sh`.
> 
> ## Capture Record
> 
> - Link the REQ with `addendum_to: REQ-236` and summarize that REQ’s existing shared renderer, separate Details control, DOM-only fold state, and tests.
> - Save the screenshot as `do-work/user-requests/UR-093/assets/REQ-486-screenshot-1-ur-view.png` if those IDs remain available.
> - Use the atomic `capture-files` transaction and commit only its declared UR, REQ, asset, and reservation paths, preserving the unrelated dirty worktree.
> ```

---
*Captured: 2026-09-01T17:29:43Z*
