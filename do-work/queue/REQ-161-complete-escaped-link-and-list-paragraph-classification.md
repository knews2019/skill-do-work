---
id: REQ-161
title: "Review fix: Complete escaped-link and list-paragraph classification"
status: pending-answers
domain: general
created_at: 2026-08-10T09:46:11Z
user_request: UR-031
addendum_to: REQ-158
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: markdown-rendered-region-classification
---

# Review Fix: Complete Escaped-Link and List-Paragraph Classification

## What

Complete the shared rendered-versus-ignored classification for the remaining escaped inline-link delimiter shapes and list-item paragraph continuations. Done means the approved Markdown class cannot recur through a different structural delimiter or block context; this remains a test-only release-guard correction and does not change the downstream publication policy.

## Context

Found during review of REQ-158. Its named regressions and full contracts pass, but production-helper probes show that closing-bracket/opening-parenthesis escapes still disagree with the bare-URL fallback and that list-item paragraph continuations are treated as indented code.

## Requirements

- Apply odd/even escape parity consistently at every inline-link structural delimiter so ignored escaped links cannot re-enter through bare first-party URL discovery.
- Recognize bullet and ordered-list paragraph continuation context before classifying four effective columns as indented code.
- Add exact positive/negative production-helper fixtures and retain all REQ-158, topology, raw/blob, path-containment, changelog-identity, focused, full-contract, and distribution behavior.

## Instances

- [ ] `_dev/tests/shipped-package-reference-contract.sh`: `[label\](first-party URL)` is rejected structurally but its URL is re-added by the bare fallback.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: `[label]\(first-party URL)` is rejected structurally but its URL is re-added by the bare fallback.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: a four-space continuation inside a bullet-list paragraph is rendered but masked as indented code.
- [ ] `_dev/tests/shipped-package-reference-contract.sh`: a four-space continuation inside an ordered-list paragraph is rendered but masked as indented code.

## Open Questions

- [ ] REQ-158 fixes escaped opening brackets and top-level paragraph continuation, but escaped closing brackets or opening parentheses still let first-party URLs re-enter through the bare-URL scan, and four-space continuations inside bullet or ordered-list paragraphs are still hidden as code. This is another review of review-generated work, so the cascade-depth rule requires your consent before completing the already-approved rendered-region behavior. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
