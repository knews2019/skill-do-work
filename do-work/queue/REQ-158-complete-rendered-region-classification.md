---
id: REQ-158
title: "Review fix: Complete rendered-region classification in shipped Markdown references"
status: pending-answers
domain: general
created_at: 2026-08-09T18:38:08Z
user_request: UR-031
addendum_to: REQ-154
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: markdown-rendered-region-classification
---

# Review Fix: Complete Rendered-Region Classification in Shipped Markdown References

## What

Make one rendered-versus-ignored decision govern every target-discovery path in the shipped-package Markdown guard for the already-approved escaped-link, indented-code, and inline-code classes. Done means these classes cannot recur as either release-blocking false positives or broken-link false negatives; this does not add a documentation policy or broaden the downstream reference contract.

## Context

Found during review of REQ-154. Its ordinary fixtures and all full contracts pass, but bounded production-helper probes show that independent masking and bare-URL scans disagree at valid Markdown boundaries.

## Requirements

- Apply the ignored escaped-link decision to relative and bare first-party URL discovery alike.
- Recognize Markdown indentation by effective columns and block context so tab-expanded indented code is ignored without hiding a four-space continuation in an active paragraph.
- Do not treat backslash-escaped backticks as inline-code delimiters; preserve ordinary exact-run code-span behavior.
- Add exact positive/negative production-helper fixtures for every instance below and retain the current topology, raw/blob, path-containment, changelog-identity, and full-contract behavior.

## Instances

- [ ] `_dev/tests/shipped-package-reference-contract.sh`: an escaped link containing a first-party URL is ignored by inline extraction but re-added by the bare-URL fallback.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: one-to-three leading spaces followed by a tab can form four columns of indented code but is scanned as live Markdown.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: a four-space continuation inside an active paragraph is rendered Markdown but is masked as code.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: backslash-escaped backticks are treated as a code span and hide a published link between them.

## Open Questions

- [ ] REQ-154 improved the release guard, but its parser can still reject valid documentation or miss a broken published link in the approved escaped-link, indented-code, and inline-code classes. This is a review of review-generated work, so your consent is required before another implementation task enters the work loop. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
