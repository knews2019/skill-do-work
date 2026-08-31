---
id: REQ-435
title: 'Review fix: Complete the doctor-forensics delegation contract'
status: claimed
domain: general
created_at: 2026-08-30T22:01:41Z
claimed_at: 2026-08-31T16:23:46Z
route: B
user_request: UR-081
addendum_to: REQ-410
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-31T14:43:25Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
tdd: true
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Complete the Doctor-Forensics Delegation Contract

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Make the natural-language forensics action executable end-to-end from its declared authorities, including report counts, actionable remedies, and valid cross-references.

## Context
Found during terminal re-review of REQ-410 after its single remediation. The action forbids deterministic rescanning and says to report only doctor evidence, while its output contract still requires queue/archive/working counts and remediation details absent from the typed doctor result. It also removed the former stuck-work reset procedure while `work-reference.md` and other consumers still point to numbered checks that no longer exist. An agent must currently rescan, omit required output, or follow stale pointers. Fold-first scan found no pending REQ or sweep in any UR that shares this delegation-contract root cause.

## Requirements
- Choose and document one complete ownership contract for queue/archive/working counts and every required deterministic report field without ad hoc rescanning.
- Ensure emitted doctor findings or the action's remaining judgment steps provide actionable remedies for the required report, including stuck-work handling.
- Replace stale numbered-check references in `work-reference.md` and every other in-scope consumer with stable canonical anchors or commands.
- Keep recurring-correction judgment and board-owned release verification outside doctor without duplicating their mechanics.
- Add contract coverage proving the action can produce its documented report using only its declared authorities.

## Constraints
- Delete unused queue/archive/working count requirements by default. Add typed counts only when a concrete consumer and regression test justify the lasting schema surface.

## Red-Green Proof
**RED prompt/case:** Execute the forensics action contract against a fixture with stuck work and mixed queue/archive/working states; assert every required report field and remedy has an authoritative source and every referenced anchor resolves.
**Why RED now:** Doctor's current result omits the required state counts and manual-reset procedure, while remaining documentation still references deleted numbered checks.
**GREEN when:** The action produces the documented report without independent mechanical scans, all remedies advance the user, and a reference contract finds no stale Check-number consumers.
**Validation:** Terminal review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: B** - Medium

**Reasoning:** The default constraint narrows the ownership choice, but the action, reference consumers, doctor result, and contract tests must be explored before freezing the exact four-file cross-subsystem scope.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*
