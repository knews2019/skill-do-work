---
id: REQ-206
title: Finish active publication delegation
status: pending-answers
domain: frontend
created_at: 2026-08-15T20:15:07Z
user_request: UR-042
addendum_to: REQ-201
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: completed-work-publication-active-delegation
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
maintenance: true
---

# Review Fix: Finish Active Publication Delegation

## What

Complete the shared-publication boundary by removing the final paraphrased whole-artifact algorithm from `present-video` and making regression tests require an active application directive at each consumer's execution step.

## Context

REQ-201 centralized the generic mechanics, but review showed that one local “one final path for every source file” sentence remains and that a passive checklist mention can satisfy the current pointer test.

## Instances

- [ ] `present-video` Step 5: retain preferred naming and resolved-path reporting, but delete the one-path/every-file algorithm paraphrase.
- [ ] Presentation contract tests: require an active `apply ... Collision-Safe Publication` directive at the publication step and reject the paraphrased algorithm.

## Requirements

- Keep only the preferred video directory and consumer-specific result reporting local.
- Require an affirmative shared-section application directive before output creation in both live consumers.
- A checklist or Rules mention alone must not satisfy delegation.
- Reject local wording equivalent to “one final path for every file.”
- Preserve the canonical section, output behavior, and all existing presentation contracts.

## Red-Green Proof

**RED prompt/case:** Remove the active Step 5 publication directive while leaving its checklist heading, and insert `one final path selected by that contract for every source file`; the current tests remain green.
**Why RED now:** Passive cross-reference presence and an uncaught paraphrase allow consumer-local mechanics to regrow while the canonicalization suite reports success.
**GREEN when:** Both mutations fail, unmodified consumers pass, and only preferred path/content/result concerns remain outside the shared publication section.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] Publication mechanics are centralized, but one paraphrase and one passive-pointer test escape remain. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle.
