---
id: REQ-479
title: '[impact-rule-change] Enforce required-lessons reads and audit un-promoted families'
status: pending
created_at: 2026-09-01T10:47:44Z
user_request: UR-088
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-478]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-477, REQ-478]
batch: lessons-transfer-routing
write_set: [skills/do-work/actions/work.md, skills/do-work/crew-members/general.md, skills/do-work-toolbox/actions/prime.md]
---

# Enforce Required-Lessons Reads and Audit Un-Promoted Families

## What

Make the work pipeline read every stamped lessons reference before implementation, and extend `prime audit` to catch missed promotions and index drift.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- **Enforced read.** `skills/do-work/actions/work.md` Step 5 (explore context) and Step 6 (builder instructions) make reading every `required_lessons` entry mandatory before implementation — a bare path means the whole satellite, a targeted path-plus-slug reference means the matching bullets — unconditional, unlike today's touch-conditional rule at `work.md:404`, which stays in force for unstamped REQs.
- **Missing file:** proceed-without-it, per the existing missing-rules-file convention. Never block on a missing lessons file.
- **Audit safety net.** Extend `skills/do-work-toolbox/actions/prime.md` `audit` to flag: a satellite with 2+ same-family entries and no corresponding Trap line in its prime; a satellite missing from the index; an index line whose path is dead.
- Keep `skills/do-work/crew-members/general.md` § Lessons Discipline consistent with the new mandatory-read rule in the same change.

## Constraints

- The enforcement is instruction-level, floor-agent compatible: reading files before coding, nothing more.
- The touch-conditional rule stays for REQs with no stamp — this REQ adds a stronger path, it removes nothing.

## Dependencies

Depends on REQ-478 (the `required_lessons` field this REQ enforces; transitively REQ-477's slugs and index that the audit checks).

## Builder Guidance

Certainty level: Firm. Latitude: exact wording and where in Steps 5/6 the mandate lands; audit finding names and severities.

## Red-Green Proof

**RED prompt/case:** (1) Dispatch a builder for a REQ stamped with a lessons reference: nothing in Step 5/6 requires the read. (2) Run `prime audit` against a satellite holding two same-family entries with no Trap line and a stale index: the audit reports healthy.
**Why RED now:** Enforcement and the audit checks do not exist; the only lessons-read rule is touch-conditional (`work.md:404`).
**GREEN when:** Steps 5 and 6 name the mandatory read (bare-path and path-plus-slug forms); `prime audit` reports the un-promoted family, the missing index line, and the dead index path as findings.
**Validation:** User confirmed (approved plan, 2026-09-01 session).

## Full Context

See `do-work/user-requests/UR-088/input.md` for complete verbatim input.

---
*Source: UR-088 (Lessons routing with token-budgeted mandatory reads and a fold-gate fix)*
