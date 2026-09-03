---
id: REQ-531
title: 'Let a consumer release a skill it authored itself'
status: pending-answers
created_at: 2026-09-03T11:10:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-461]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-461]
review_generated: true
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/release.go
  - skills/do-work/tools/do-work-cli/internal/publication/release_test.go
---

# Let a Consumer Release a Skill It Authored Itself

## What

REQ-461 replaced a directory-name denylist with proof of project ownership, and one of its two proofs over-refuses: a `SKILL.md` between the repository root and the target's parent is read as "an installed package owns this subtree". That is true of an installed skill and equally true of a skill the consumer **wrote**. So a consumer who authors `skills/my-skill/` in their own repository can no longer version it in place.

The old predicate accepted this, because it denylisted only the literal `skills/do-work/` prefix. This is a behavior regression introduced by REQ-461, recorded there as residual gap R2 rather than absorbed silently.

## Instances

- A consumer repository containing an authored `skills/my-skill/SKILL.md` and `skills/my-skill/VERSION`, releasing without `maintainer_release`. Before REQ-461: accepted. After: refused with `RELEASE-TARGET-OWNERSHIP-UNVERIFIED` naming their own marker.
- The same shape for an authored skill's `CHANGELOG.md`.
- **In-repo sibling suite packages (folded from REQ-461's review).** `skills/do-work-board/`, `skills/do-work-knowledge/` and `skills/do-work-toolbox/` each carry their own `SKILL.md`, so any release-metadata file under them now needs `maintainer_release` where the old prefix-only predicate accepted them. Not currently live — `skills/do-work-board/tools/queue-kanban/VERSION` is `0.236.20`, not the suite version, so discovery skips it — but it is the same over-refusal inside this very repository, which makes it a useful test bed for whatever declaration this REQ settles on.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- REQ-461's decision **D-04 / R2**, escalated by its builder rather than discovered later. The filesystem genuinely cannot distinguish a skill you wrote from one you installed: both are tracked, and both carry the marker.
- Worth noting what is *not* in question. R1 (a committed vendored copy with no `SKILL.md` is accepted) is the mirror-image gap and stays open deliberately: for the class the guard protects — installed suite metadata — `suitemanifest.ValidateSuite` requires a non-empty `SKILL.md` for every module, so it cannot be marker-free.

## Open Questions

- [ ] What should distinguish a skill the consumer authored from one they installed? The filesystem cannot, so the answer has to come from a declaration somewhere.
  - Recommended: the suite installer already knows which paths it wrote. If an install records its own managed paths (it already writes a managed section elsewhere), an authored skill is simply one no install claims — an affirmative declaration in the same spirit as REQ-461, rather than a new heuristic.
  - Also: let the release manifest declare an authored-package target explicitly, so the consumer states ownership and the planner verifies the declaration is not contradicted by an install record.
  - Also: accept a marked subtree when the marker itself is tracked *and* no do-work suite module manifest sits beside it — narrower, but leans on suite-specific structure and would not generalize to other package kinds.
  - Also: leave it as is and require `maintainer_release: true` for an authored skill — rejected as a default, because that flag exists to permit mutating *the suite's own* metadata and overloading it would make the exemption mean two different things.

## Detailed Requirements

- A consumer must be able to release a skill their own repository authored, without asserting `maintainer_release`.
- An installed or vendored suite package must remain refused, exactly as REQ-461 established.
- The distinction must rest on a declaration the planner can verify mechanically, not on a directory name and not on a new heuristic about file contents.
- `maintainer_release: true` must keep meaning only "may mutate the suite's own metadata", and must not become the general-purpose override for ownership questions.
- Refusals must keep naming the target and the missing or contradicted evidence.

## Constraints

- Do not reintroduce a directory-name test in any form; REQ-461's constraint still binds.
- Do not weaken the Git-index half of the proof.

## Dependencies

Depends on REQ-461, which introduced both the proof and this gap.

## Red-Green Proof

**RED prompt/case:** In a consumer repository, commit `skills/my-skill/SKILL.md` and `skills/my-skill/VERSION`, then plan a release of that VERSION without `maintainer_release`.
**Why RED now:** `enclosingPackageMarker` finds the consumer's own marker and reports that an installed package owns the subtree.
**GREEN when:** that release plans successfully; an installed suite package under the same shape still refuses; and both refusals still name the target and the missing evidence.

---
*Source: REQ-461 decision D-04, residual gap R2.*
