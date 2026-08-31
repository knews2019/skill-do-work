---
id: REQ-431
title: 'Review fix: Couple documentation rewrites to their owning moves'
status: claimed
domain: general
created_at: 2026-08-30T20:35:44Z
claimed_at: 2026-08-31T15:44:18Z
route: B
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-31T15:44:18Z
  basis:
    - Route B
    - 3 acceptance criteria
    - full-suite verification
user_request: UR-081
addendum_to: REQ-409
review_generated: true
depends_on: [REQ-430]
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Couple Documentation Rewrites to Their Owning Moves

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Compose documentation-link rewrites while retaining each rewrite's dependency on the move that makes its destination valid.

## Context
Found during re-review of REQ-409. All composed rewrites for one document are attached to the first move group, so refusal can either suppress a valid later rewrite or publish a link to a destination whose move was refused. Fold-first scan found no pending REQ or sweep in any UR that shares this conditional-rewrite root cause.

## Requirements
- Associate every planned link rewrite with the exact move that justifies its destination.
- Compose multiple successful rewrites to one document without losing their independent eligibility.
- Never rewrite a link when its owning move is refused or fails.

## Red-Green Proof
**RED prompt/case:** Plan two moves referenced by one document, refuse each move in turn, and assert only the link owned by the successful move is rewritten in both cases.
**Why RED now:** The composed document edit is owned solely by the first move group.
**GREEN when:** Both refusal directions preserve valid links, avoid nonexistent destinations, and retain the other safe rewrite.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: B** - Medium

**Reasoning:** The conditional-rewrite defect and acceptance behavior are clear, but the cleanup planner's rewrite composition and operation ownership seams must be explored before freezing the write set.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*
