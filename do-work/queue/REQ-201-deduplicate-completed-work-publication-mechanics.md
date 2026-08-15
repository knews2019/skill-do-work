---
id: REQ-201
title: Deduplicate completed-work presentation publication mechanics
status: pending
domain: frontend
created_at: 2026-08-15T17:43:17Z
user_request: UR-042
addendum_to: REQ-191
review_generated: true
effort_estimate: trivial
sweep: true
sweep_key: completed-work-publication-mechanics-duplicated
maintenance: true
tdd: true
---

# Deduplicate Completed-Work Presentation Publication Mechanics

## What

Keep completed-work presentation consumers limited to their consumer-specific preferred output and delegate the collision/no-overwrite algorithm and verification to `actions/completed-work-presentation-reference.md`. Sweep every current consumer of that reference so this duplicated contract class cannot recur as parallel wording.

## Context

Found during review of REQ-191. `present-video` declares the shared reference the sole publication contract, then restates the existing-directory branch and related verification, creating a second rule that can drift.

## Instances

- [ ] `skills/do-work-toolbox/actions/present-video.md` publication step, Rules, and Verification Checklist: retain the preferred `<canonical-ID>-video` directory and immutable-output intent, but remove duplicated collision/suffix mechanics and point their execution and verification to the shared reference.

## Requirements

- Sweep all current consumers of `completed-work-presentation-reference.md` for duplicated target-resolution, safety, evidence, or publication mechanics.
- Keep consumer-specific eligibility, content shape, preferred output naming, and result reporting local.
- Replace duplicated collision/no-overwrite algorithm and verification wording with a named pointer to the shared contract.
- Preserve existing output behavior and never mutate prior deliverables.

## Red-Green Proof

**RED prompt/case:** Search shared-reference consumers for local existing-path collision branches, suffix selection algorithms, or duplicated no-overwrite verification after they declare the reference canonical.
**Why RED now:** `present-video.md` currently repeats the preserve-and-suffix branch and related verification at its publication boundary.
**GREEN when:** Every current consumer names only its preferred output and consumer-specific result contract while collision/no-overwrite mechanics and their verification have one canonical definition in `completed-work-presentation-reference.md`; focused presentation contracts and the maintainer suite pass.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
