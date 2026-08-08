---
id: REQ-149
title: "Review fix: Complete moved-command shim mapping"
status: pending
domain: general
created_at: 2026-08-08T15:38:44Z
user_request: UR-031
addendum_to: REQ-144
review_generated: true
effort_estimate: normal
---

# Review Fix: Complete Moved-Command Shim Mapping

## What
Make every legacy command routed by the modular core print one concrete, correct sibling invocation and stop. Close the whole router-to-shim mapping class, not only the aliases observed by this review.

## Context
Found during review of REQ-144. Several routed board and knowledge aliases fall through to help, `install run-do-work-update` prints the wrong Just recipe, and the current content-only test cannot detect either defect.

This is a standalone user-visible compatibility defect rather than part of a broader sweep: one complete table-driven router/shim contract closes it.

## Requirements
- Map every legacy alias routed to `moved-command-shim.md` to the exact canonical board, knowledge, or toolbox invocation while preserving trailing arguments.
- Correct the managed updater recipe replacement so it prints `just run-do-work-update`, not `just run-kanban`.
- Keep the shim print-only: it must never forward, load sibling actions, or mutate state.
- Replace placeholder-string greps with a table-driven contract that exercises every routed alias and rejects unmatched or ambiguous mappings.
