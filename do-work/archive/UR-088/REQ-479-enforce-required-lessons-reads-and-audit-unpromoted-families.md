---
id: REQ-479
title: '[impact-rule-change] Enforce required-lessons reads and audit un-promoted families'
status: completed
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
write_set: [skills/do-work/actions/work.md, skills/do-work/crew-members/general.md, skills/do-work-toolbox/crew-members/general.md, skills/do-work-toolbox/actions/prime.md]
claimed_at: 2026-09-01T16:15:11Z
route: C
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-09-01T16:16:30Z
  basis:
    - Route C
    - 4-file write set
    - 2 subsystems involved
    - 9 acceptance criteria
    - dependency depth 2
    - cross-route regression gates
kb_status: pending
completed_at: 2026-09-01T16:28:16Z
commit: 9a1b7bfb
---

# Enforce Required-Lessons Reads and Audit Un-Promoted Families

## What

Make the work pipeline read every stamped lessons reference before implementation, and extend `prime audit` to catch missed promotions and index drift.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Defined the claim-time consult, three additive read layers, and six mechanical audit predicates before implementation.
- [x] **[APPLY]:** Updated the work action, both general-crew copies, and prime audit within the declared four-file scope.
- [x] **[UNIFY]:** Reviewed all four diffs, checked the shipped index against the new algorithm, ran contract searches and `git diff --check`, and passed the canonical clean gate.

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

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3259 tokens; owning prime governs this REQ, but the partial-coverage satellite cannot be narrowed and does not fit the shared 2000-token budget. It was still read under the always-on touch-conditional prime rule.

## Triage

**Route: C** — the enforced-read path changes exploration and builder dispatch, while the safety net widens the mutating prime audit's health contract.

## Plan

1. Add a claim-time index consult that counts captured stamps first, applies the shared budget contract to new matches, refreshes `required_lessons`, and records every drop.
2. Make Steps 5 and 6 enforce captured plus claim-time reads while retaining the touch-conditional satellite read for every REQ and recording missing files in the hand-back.
3. Extend prime audit with six mechanical lesson/index drift predicates and aligned report/checklist guidance, then mirror the three-layer read rule in both general crew files.

**Plan validation:** Three tasks cover the claim-time gap, builder enforcement, every original/addendum audit predicate, and both live general-crew copies.

## Exploration

- Step 5 currently passes prime lesson satellites to exploration only when a prime has a Lessons section; it does not inspect `required_lessons` or the lessons index.
- Step 6 currently repeats only the touch-conditional prime rule, so the builder never receives an unconditional stamped-read mandate.
- `prime audit` already discovers every source `lessons-*.md` satellite in Step 4 and owns read-only health findings, making it the smallest place to compare recurrence, index paths, estimates, family sets, and coverage.
- The two general crew files are semantic mirrors with only package-relative citations intentionally different; both load in shipped workflows.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — claim-time consult, required exploration reads, and three-layer builder mandate
- `skills/do-work/crew-members/general.md` (modify) — canonical Lessons Discipline enforcement
- `skills/do-work-toolbox/crew-members/general.md` (modify) — toolbox semantic mirror
- `skills/do-work-toolbox/actions/prime.md` (modify) — mechanical recurrence/index audit and report contract

**Acceptance criteria:** Captured entries are budgeted first; every claim consults current index hooks; resolved entries are read before implementation; missing files do not block and are reported; touch-conditional reads remain additive for all REQs; audit flags unpromoted recurrence, missing/dead rows, estimate drift over 25%, family-set mismatch, and false full coverage.

## Decisions

- **D-01 — DECIDE & STATE:** Refresh `required_lessons` frontmatter at claim rather than create a new run artifact; Step 6 already reads the REQ, and one field avoids a second source of truth.
- **D-02 — DECIDE & STATE:** Treat the index row's `Families` cell as the hook's canonical slug set. It is the exact machine-comparable set established by REQ-477; the prose `When it applies` cell remains the human match hook.
- **D-03 — DECIDE & STATE:** Keep prime-audit enforcement instruction-level and floor-agent compatible; fixed-string grep, `wc -c`, integer comparisons, and path existence checks need no new tool or parser.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified) — all-route claim-time consult, bounded frontmatter refresh, enforced exploration/builder reads, and missing/drop evidence
- `skills/do-work/crew-members/general.md` (modified) — canonical three-layer Lessons Discipline
- `skills/do-work-toolbox/crew-members/general.md` (modified) — toolbox semantic mirror of the same rule
- `skills/do-work-toolbox/actions/prime.md` (modified) — six named mechanical lesson/index health findings and report/checklist coverage

**What was done:** Every claimed REQ now re-evaluates the current lessons index before implementation, so pre-existing and serially captured work can inherit lessons created after capture. Builders read captured and claim-time pointers unconditionally while retaining the broader touched-prime rule, and prime audit can identify every accepted promotion/index drift class.

## Qualification

Passed — all four declared files are substantive and wired into exploration, builder dispatch, always-loaded crew guidance, and the prime audit report. Every original and addendum criterion maps to an explicit instruction or named audit predicate.

## Testing

- Current-index replay: PASS — six source satellites equal six index rows; all paths live; estimates are within 25%; family sets match; coverage is honest; repeated slugged families have matching prime traps.
- Contract searches: PASS — all six `LESSONS-*` finding names, the all-route consult, the three additive layers, and missing-file hand-back evidence are present.
- General-crew mirror comparison: PASS — new Lessons Discipline prose is identical; the only file difference remains the intentional package-relative prime citation.
- `git diff --check` on the REQ scope: PASS.
- Clean isolated `bash _dev/tests/maintainer-verify.sh`: PASS (strict browser lane skipped because no browser was available, per the gate's own contract).

## Review

**Verdict: Approve.** Route A receives the consult even though it skips exploration; existing stamps consume budget before new matches; capture-time, claim-time, and touch-conditional reads remain additive; audit checks both index directions and never guesses pre-slug families. No Important findings.

**Acceptance:** Pass

## Lessons Learned

**What worked:** Treating claim as a second projection point closes the exact time-order gap capture cannot see, while reusing the same named budget contract avoids a competing rule.

**What didn't:** Nesting the consult under the original “Routes B and C” exploration heading initially excluded Route A; checking the trigger against every route exposed it before release.

**Worth knowing:** Context routing that runs at capture only is stale by construction for serial batches; enforce it again at claim, and make that consult independent of whether the route performs exploration.

## Orientation

[MAP CHANGED] Lessons now enter implementation through three additive paths: capture stamps, an all-route claim-time index consult, and touched-prime discovery. Prime audit is the read-only safety net for the index and recurrence-to-Trap contract.
