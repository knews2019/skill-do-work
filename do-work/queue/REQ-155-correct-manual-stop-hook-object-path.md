---
id: REQ-155
title: "Review fix: Correct the manual Stop-hook object path"
status: pending-answers
domain: general
created_at: 2026-08-08T19:28:25Z
user_request: UR-031
addendum_to: REQ-151
review_generated: true
effort_estimate: normal
---

# Review Fix: Correct the Manual Stop-Hook Object Path

## What
Make the no-JSON-tool instruction identify individual nested Stop-hook objects exactly, so following it cannot delete a whole hooks array containing custom neighbors.

## Context
Found during review of REQ-151. The instruction says `hooks.Stop[*].hooks objects`, which selects each inner hooks array rather than the individual objects at `hooks.Stop[*].hooks[*]`; its exact-output test preserves the same ambiguity.

This is a standalone user-visible wording defect in the manual fallback. Automated jq/Python reconciliation and the current custom-hook preservation behavior are already correct.

## Requirements
- Name the matching object path as `hooks.Stop[*].hooks[*]` in both byte-identical installer copies.
- State that an outer Stop wrapper is removed only if targeted-object removal leaves its `hooks` array empty, matching jq/Python behavior.
- Keep explicit preservation of unrelated/custom hooks, including same-wrapper neighbors.
- Update the exact-output regression and compare the described result for mixed and guard-only wrappers with automated reconciliation semantics.
- Keep focused/full installer contracts, syntax/lint, and byte identity passing.

## Open Questions

- [ ] The manual fallback now targets the retired guard, but its published JSON path can still be read as deleting an entire array that contains custom hooks. The cascade-depth rule requires your consent before automatically working a follow-up created by the review of another review-generated task. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
