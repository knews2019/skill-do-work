---
id: REQ-163
title: "Review fix: Complete remaining inline-link and list-fence classification"
status: pending-answers
domain: general
created_at: 2026-08-10T11:30:40Z
user_request: UR-031
addendum_to: REQ-161
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: markdown-rendered-region-classification
write_set:
  - _dev/tests/shipped-package-reference-contract.sh
---

# Review Fix: Complete Remaining Inline-Link and List-Fence Classification

## What

Finish the test-only rendered-region classifier so delimiter parity, label context, and list-item fences cannot disagree with target discovery. Done means relative and first-party targets classify consistently and fenced examples inside list items remain ignored, without changing publication, topology, raw/blob, path-containment, runtime, or documentation policy.

## Context

Found during independent review of REQ-161. Its four reproduced cases and declared suites pass, but actual production-helper probes found nearby variants that still hide live relative references or scan fenced code as published Markdown. Because REQ-161 is itself review-generated, the cascade-depth rule requires user consent before this non-critical continuation becomes claimable.

## Requirements

- Make structural extraction preserve zero/even-parity destination-opening-parenthesis forms for relative targets, not only first-party URLs recovered by the bare fallback.
- Keep escaped opening brackets used as label content from causing the containing live link to be masked.
- Recognize backtick and tilde list-item fences with attached info strings while preserving approved list-paragraph continuations.
- Add exact production-helper fixtures for every instance and retain all publication-policy, offset, focused, full-contract, and distribution behavior.

## Instances

- [ ] Even two/four-backslash destination-opening-parenthesis forms retain relative targets, not only first-party URLs recovered by fallback.
- [ ] An escaped opening bracket used as label content does not cause the containing live link to be masked.
- [ ] Bullet and one-to-nine-digit ordered list items recognize backtick/tilde fences with attached info strings, including nested tilde fences.

## Open Questions

- [ ] REQ-161 fixes its four reproduced cases, but the remaining parity, label-content, and list-fence variants can still hide broken relative references or scan code as published Markdown. This is a non-critical follow-up to review-generated work, so the cascade-depth rule requires your approval instead of silently extending the queue. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
