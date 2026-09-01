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

## Addendum (2026-09-01)

User added (v4 revision, validate-feedback Findings 4 and 6 — Accept):

> ```
> Today's touch-conditional rule at work.md:404 stays in force for all REQs, stamped or not; the mandatory read is additive, not a replacement, and Step 6 says so in one sentence so a builder never guesses which regime applies. A missing listed file is proceed-without-it, per the existing missing-rules-file convention, and the miss is recorded in the hand-back. [...] Extend skills/do-work-toolbox/actions/prime.md audit to flag, all mechanically: [...] an index estimate more than ~25% off the recomputed value; an index hook whose slug set does not match the slug set actually present in the satellite (either direction); a `slugged: full` flag on a satellite that still has un-slugged bullets.
> ```

- Resolved conflict: the Detailed Requirements sentence "which stays in force for unstamped REQs" is superseded — the touch-conditional rule at work.md:404 stays in force for ALL REQs, stamped or not. The mandatory `required_lessons` read is additive, never a replacement, and Step 6 states this in one sentence. (The original wording would have dropped the conditional read for satellites the budget excluded from a stamped REQ.)
- A missing listed file still proceeds-without-it, and the miss is recorded in the hand-back.
- The audit gains three further mechanical checks: an index estimate more than ~25% off the recomputed value; an index hook whose slug set does not match the satellite's actual slug set (either direction); a `slugged: full` flag on a satellite that still has un-slugged bullets.
- Provenance: validate-feedback 2026-09-01, Findings 4 and 6. Surface-cost: additive-read fix N/A (removes an accidental narrowing); audit checks Earned — the index becomes a routing and budgeting authority, and all three checks are read-only recompute/grep comparisons inside the existing milestone audit.

## Addendum (2026-09-01, claim-time consult)

User added (approved UR-081 improvement plan, 2026-09-01 session):

> ```
> At claim, work.md Step 5 consults the lessons index for EVERY REQ — stamped or
> not: grep the index's family hooks against the REQ's scope/prime files and treat
> matches as required reading, recording the resolved list in the run scratch or
> refreshing required_lessons on the REQ (builder's choice). The consult uses the
> same budget constant, entry forms, cost rule, and full-only targeting as
> REQ-478's stamping, counting entries already stamped on the REQ against the
> budget first; budget-dropped matches are recorded, never silent. This covers
> REQ-457 and every other pre-existing REQ at claim, and it is what covers a
> serial UR whose later REQs were stamped before their siblings' lessons existed —
> capture-time stamping cannot see a lesson written at an earlier REQ's archive
> (work.md Step 7), and Trap promotion fires only on the second same-family lesson
> write, catching recurrence #3, not #2 (the REQ-414-to-415 shape).
> ```

- Resolved conflict: this supersedes the body's capture-stamped-only scope for the
  mandatory read — the enforced read now has two triggers: stamped `required_lessons`
  entries (body) and claim-time index matches (this addendum). The first addendum's
  rule that the touch-conditional read at work.md:404 stays in force for ALL REQs is
  unchanged; the claim-time consult is a third additive layer, and Step 6's
  one-sentence regime statement must name all three.
- Depends on REQ-478's budget constant, entry forms, and cost rule (already this
  REQ's `depends_on`) and on REQ-477's index hooks.
- REQ-481 (the one-time stamping pass) is superseded by this addendum and was
  cancelled 2026-09-01: every pending REQ now receives the same decision at claim,
  so a one-shot pass adds nothing it would not redo.
- Provenance: maintainer-approved plan, 2026-09-01 session (UR-081 never-ending-story
  analysis). Surface-cost: Earned — the 2026-08-31 run's recurrences (REQ-414→415)
  happened between REQs captured in one batch before any lesson existed, the exact
  window capture-time stamping cannot reach.
