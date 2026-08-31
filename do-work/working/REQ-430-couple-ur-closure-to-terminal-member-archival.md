---
id: REQ-430
title: 'Review fix: Couple UR closure to terminal member archival'
status: claimed
domain: general
created_at: 2026-08-30T20:35:44Z
claimed_at: 2026-08-31T14:50:42Z
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-31T14:50:42Z
  basis:
    - Route B
    - 3-file write set
    - 3 acceptance criteria
    - full-suite verification
user_request: UR-081
addendum_to: REQ-409
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Couple UR Closure to Terminal Member Archival

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Make UR closure depend on the successful archival of every live terminal member in the same cleanup run, so a refused member can never leave an archived UR input beside an active REQ.

## Context
Found during re-review of REQ-409. A dirty terminal REQ group was refused while its independently planned `CLOSE-UR` group still archived the active UR input. Fold-first scan found no pending REQ or sweep in any UR that shares this cleanup-dependency root cause.

## Requirements
- Express UR closure as dependent on every member move required for that closure.
- Refuse or defer closure whenever any required member group is refused or fails.
- Preserve independent progress for cleanup operations unrelated to that UR.

## Red-Green Proof
**RED prompt/case:** Put a dirty terminal member in `do-work/queue/` under an otherwise closable UR and run cleanup; assert neither the member nor the UR input moves.
**Why RED now:** The planner emits the UR-close group independently from the member archival group.
**GREEN when:** The named fixture keeps both inputs active, reports the blocking member evidence, and still applies unrelated safe groups.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: B** - Medium

**Reasoning:** The required cleanup dependency behavior is clear, but the planner's existing operation-group structure and tests must be explored before declaring the exact files and acceptance boundary.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*
