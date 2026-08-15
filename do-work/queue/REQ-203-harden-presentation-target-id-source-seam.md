---
id: REQ-203
title: Harden presentation target-ID source-seam tests
status: pending-answers
domain: general
created_at: 2026-08-15T19:20:09Z
user_request: UR-042
addendum_to: REQ-197
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: presentation-target-id-source-seam
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
maintenance: true
---

# Review Fix: Harden Presentation Target-ID Source-Seam Tests

## What

Make the completed-work presentation ID contracts prove active, ordered inheritance from the canonical Target ID Resolution source without copying its grammar into callers. Close the entire mutation class, not only the examples found in one review.

## Context

REQ-197's product instructions now cite and apply the right shared contract, but its single remediation left the regression guard able to accept “read without applying,” an omitted pre-dispatch order, and copied membership grammar.

## Instances

- [ ] Shared presentation resolver: reject semantic negations that retain an `apply` substring and keep application before target lookup.
- [ ] `present-work` item dispatch: require shared-contract application before the item branch and reject copied membership grammar as well as copied token examples.
- [ ] Regression block: use a replayable mutation matrix with safe positive controls rather than keyword-presence alone.

## Requirements

- Require a word-bounded, affirmative application directive rather than a matching substring.
- Reject semantic negations including “without applying” for both callers.
- Enforce Target ID Resolution before each caller's lookup or item-dispatch boundary.
- Reject caller-local copies of token and UR-membership grammar.
- Preserve the current correct caller instructions and canonical source contract.

## Red-Green Proof

**RED prompt/case:** Mutate each caller's active directive to “read without applying Target ID Resolution,” move the `present-work` citation below its item branch, and add a caller-local membership definition; the current assertions accept those invalid states.
**Why RED now:** A caller can silently stop applying or fork the canonical grammar while the suite remains green.
**GREEN when:** Every semantic-negation, ordering, and copied-membership mutation fails; the unmodified callers and canonical source pass; and the focused and canonical suites remain green.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] The current product instructions are correct, but their regression test still misses one semantic mutation family. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle.
