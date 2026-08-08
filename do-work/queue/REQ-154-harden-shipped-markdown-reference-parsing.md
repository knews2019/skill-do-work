---
id: REQ-154
title: "Review fix: Harden shipped Markdown reference parsing"
status: pending-answers
domain: general
created_at: 2026-08-08T19:16:22Z
user_request: UR-031
addendum_to: REQ-150
review_generated: true
effort_estimate: normal
---

# Review Fix: Harden Shipped Markdown Reference Parsing

## What
Make the shipped-package reference guard distinguish published Markdown links from standards-valid code, comments, escapes, and escaped destinations so legitimate documentation cannot falsely block a release.

## Context
Found during review of REQ-150. The guard repairs and protects the current ten broken references, but its handwritten parser also treats links inside four-space indented code, HTML comments, and escaped link syntax as live links, and misreads escaped-parenthesis destinations.

This is a standalone user-visible release-gate defect: it changes only the parser and adversarial fixtures, not the canonical reference policy or current link repairs.

## Requirements
- Ignore Markdown link-shaped text inside fenced code, four-space indented code, HTML comments, inline code, and escaped link syntax.
- Decode standards-valid escaped destination punctuation, including parentheses, before resolution.
- Preserve inline/reference-style published link detection, source/install topology checks, first-party raw/blob policy, path-escape rejection, and changelog identity.
- Add adversarial positive and negative fixtures for every parser class above.
- Keep the current four-module reference scan and full contract regressions passing.

## Open Questions

- [ ] The current link guard fixes all known broken package references, but its parser can reject valid future documentation that contains link-shaped examples or escapes. The cascade-depth rule requires your consent before automatically working a follow-up created by the review of another review-generated task. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
