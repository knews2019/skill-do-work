---
id: REQ-496
title: '[impact-critical] Review fix: Add shared executable already-green repair validator'
status: claimed
priority: now
domain: backend
created_at: 2026-09-02T04:53:21Z
user_request: UR-095
addendum_to: REQ-494
review_generated: true
impact: impact-critical
effort_estimate: effort-substantive
tdd: true
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
depends_on: [REQ-494]
related: [REQ-492]
sweep: true
sweep_key: already-green-repair-shared-validator-missing
claimed_at: 2026-09-03T22:25:28Z
route: C
estimate:
  p50_active_minutes: 60
  confidence: low
  calculated_at: 2026-09-03T22:25:42Z
  basis:
    - Route C
    - shared executable authority across action, validation, completion, and selector seams
    - adversarial Git/staging fixtures and full-suite verification
---

# Review Fix: Add Shared Executable Already-Green Repair Validator

## What

Replace the duplicated prose/test decision for an already-green repository-gate repair with one executable validator consumed by both work and review.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- One executable validator is the sole decision authority for both TDD bypass and no-diff review.
- Match the expected fingerprint to actual repair intake evidence; never self-assert it.
- Validate staged paths against the exact successful canonical-completion result, not an archive prefix.
- Refuse an unrelated staged archive and every ordinary, malformed, nonempty, release-mutated, or over-staged neighbor.
- Drive real REQ/Git state through canonical completion, metadata, and selector behavior.

## Red-Green Proof

**RED prompt/case:** REQ-494's fixture can report eligibility from its own `action_decisions()` oracle while shipped guards or evidence are wrong.

**Why RED now:** Test and prose duplicate the decision instead of consuming one executable authority.

**GREEN when:** Removing or corrupting the shared validator breaks both action consumers; exact intake/result paths pass and an unrelated staged archive refuses.

**Validation:** REQ-494 re-review critical finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Instances

- [ ] `impact-critical` — TDD and review decisions still use a parallel test oracle.
- [ ] `impact-critical` — Fingerprint identity is not sourced from repair intake.
- [ ] `impact-critical` — Archive staging is prefix-authorized rather than exact-result-authorized.

## Triage

**Route: C** — Complex

**Reasoning:** This replaces duplicated action prose/test authority with one executable validator spanning repair intake, exact completion output, staged-path ownership, TDD bypass, no-diff review, metadata, and selector behavior. Planning and source exploration are required before declaring the write boundary.

**Planning:** Required
