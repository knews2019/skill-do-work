---
id: REQ-156
title: "Review fix: Handle Just multiline strings in collision scanning"
status: pending-answers
domain: general
created_at: 2026-08-08T19:55:50Z
user_request: UR-031
addendum_to: REQ-152
review_generated: true
effort_estimate: normal
---

# Review Fix: Handle Just Multiline Strings in Collision Scanning

## What
Make reserved-recipe collision scanning ignore recipe- and alias-shaped text inside valid triple-quoted Just variable strings.

## Context
Found during review of REQ-152. The line-based detector correctly handles normal definitions, but it does not retain multiline-string state and therefore rejects valid custom Justfiles whose string content begins with a reserved name.

This is a standalone user-visible parser defect. Collision ordering, exact recipe/alias detection, managed-span exclusion, and pre-mutation preservation are already correct.

## Requirements
- Track triple-single-quoted and triple-double-quoted Just string boundaries across lines before classifying definitions.
- Ignore recipe- and alias-shaped content inside those strings, including reserved names.
- Preserve detection immediately before/after strings, escaped content rules, CRLF handling, managed-span exclusion, and deterministic collision reporting.
- Add Just-parseable positive fixtures for both quote styles and nearby real-collision controls.
- Keep installer/updater/full contracts, paired identities, syntax/lint, and pre-mutation preservation passing.

## Open Questions

- [ ] The new no-Just safety gate rejects real reserved recipes correctly, but it can also reject a valid custom multiline string whose text looks like a recipe. The cascade-depth rule requires your consent before automatically working a follow-up created by the review of another review-generated task. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
