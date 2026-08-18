---
id: REQ-259
title: Retire the skill-root citation reading at its three unbackticked sites
status: claimed
claimed_at: 2026-08-18T21:16:24Z
route: B
created_at: 2026-08-18T18:07:48Z
status_changed_at: 2026-08-18T20:59:31Z
user_request: UR-055
addendum_to: REQ-249
domain: general
review_generated: true
sweep: true
sweep_key: retired-reading-outside-backtick-spans
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/SKILL.md
- skills/do-work/actions/commit.md
- skills/do-work/crew-members/security.md
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route B
    - 3-file write set
    - 3 subsystems involved
    - 3 acceptance criteria
    - full-suite verification
---

# Retire the Skill-Root Citation Reading at Its Three Unbackticked Sites

## What

REQ-249 retired the skill-root-relative citation reading and swept every **backticked** cross-package citation to the literal form — but the retired *reading* survives at three shipped sites that are bare text, which both the sweep and the new checker are structurally blind to. One of them states the retired resolution rule as prose in the core router, now contradicting the prime and the swept corpus.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Found by REQ-249's independent review (Important, gate: rule-change; reproduced by execution — paths verified non-resolving, checker verified passing over them). Not builder scope drift: the user-confirmed rule's letter covers backticked citations, and these are bare text — but the diff retired the reading, and these sites still state or use it. Created `pending-answers` per the generation-≥2 cascade stop (REQ-249 is itself review-generated).

## Instances

- [ ] **`skills/do-work/SKILL.md:16`** — "When a core action names a sibling path, resolve it from the parent directory containing these four skill roots" is the retired resolution rule stated as prose in the core router; applied to the new `../../` form it computes wrong paths. Restate to match the literal rule in `_dev/primes/prime-action-files.md` § Cross-Referencing.
- [ ] **`skills/do-work/actions/commit.md:17`** — bare `../do-work-toolbox/actions/inspect.md`, wrong depth from `actions/`.
- [ ] **`skills/do-work/crew-members/security.md:3`** — bare `../do-work-toolbox/actions/code-review.md` in the JIT_CONTEXT comment; the checker's comment scan covers backticked spans only, so this evades.

## Requirements

- No shipped text states or uses the retired skill-root-relative reading outside REQ-249's documented exemptions (fenced template/example blocks).
- Whether the two bare paths become backticked (and thereby checkable) is builder latitude; the SKILL.md prose must agree with the prime either way.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Open Questions

- [x] REQ-249's review found the retired citation reading surviving at three unbackticked shipped sites (listed under Instances). Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-18]:** User approved via `do-work clarify`. An external automated review independently flagged the same three sites while this question was open, which corroborates the finding.

---

## Triage

**Route: B** - Medium

**Reasoning:** Three named instances, but the requirement is a class claim (`no shipped text states or uses the retired reading`), so the builder must sweep the corpus for further bare-text sites the backtick-scoped checker is blind to before declaring it closed.

**Planning:** Not required

---

## Exploration

All three instances reproduce exactly as the REQ states:

- `skills/do-work/SKILL.md:16` — "When a core action names a sibling path, resolve it from the parent directory containing these four skill roots. Do not search the core package for an extension action." The first sentence is the retired reading; the second is a separate, still-true instruction.
- `skills/do-work/actions/commit.md:17` — `route to ../do-work-toolbox/actions/inspect.md instead`, bare, wrong depth from `actions/` (needs `../../`).
- `skills/do-work/crew-members/security.md:3` — `../do-work-toolbox/actions/code-review.md` inside the JIT_CONTEXT comment; that same comment already spells a same-package citation in backticks, so the file mixes both forms.

The rule to agree with is `_dev/primes/prime-action-files.md` § Cross-Referencing (line 91): literal relative path from the citing file's own directory, depth per-file, fenced template/example blocks exempt, and `_dev/tests/shipped-package-reference-contract.sh` checks the backticked ones in both topologies.

*Explored inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/SKILL.md` (modify) — restate line 16's resolution prose to match the prime's literal rule
- `skills/do-work/actions/commit.md` (modify) — line 17's bare sibling path
- `skills/do-work/crew-members/security.md` (modify) — line 3's bare sibling path in the JIT_CONTEXT comment

**Files I will NOT touch:** `_dev/primes/prime-action-files.md` (it is the rule this agrees with, not a site to change), `_dev/tests/shipped-package-reference-contract.sh` (widening the checker to bare text is its own REQ, not this one — record it as a Discovered Task if the sweep argues for it), any other package's files.

**Acceptance criteria (restated from REQ):**
- [ ] No shipped text states or uses the retired skill-root-relative reading outside REQ-249's fenced template/example exemptions — established by sweeping the corpus for the primitive, not by fixing the three listed instances
- [ ] `skills/do-work/SKILL.md`'s prose agrees with `_dev/primes/prime-action-files.md` § Cross-Referencing
- [ ] Any site found beyond the three is either fixed here or reported with its reason
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0
