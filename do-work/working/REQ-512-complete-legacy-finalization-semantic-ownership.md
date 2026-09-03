---
id: REQ-512
title: '[impact-critical] Review fix: Complete legacy finalization semantic ownership'
status: claimed
priority: now
domain: backend
created_at: 2026-09-02T18:08:38Z
user_request: UR-097
addendum_to: REQ-499
review_generated: true
impact: impact-critical
effort_estimate: effort-substantive
tdd: true
suggested_spec: bug-fix
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
sweep: true
sweep_key: legacy-finalization-semantic-ownership-incomplete
claimed_at: 2026-09-03T22:25:28Z
route: C
planning_at: 2026-09-03T22:37:36Z
exploration_at: 2026-09-03T22:37:36Z
dispatch_at: 2026-09-03T22:39:03Z
write_set:
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_req512_test.go
estimate:
  p50_active_minutes: 70
  confidence: low
  calculated_at: 2026-09-03T22:25:42Z
  basis:
    - Route C
    - legacy append provenance and three workspace ecosystems
    - recovery and fail-closed enumeration matrix
---

# Review Fix: Complete Legacy Finalization Semantic Ownership

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Replace heuristic legacy-follow-up and workspace-release association with provenance that bounds the entire owned append and identifies the actually changed workspace source before selecting mirrors, so neither foreign bytes nor equal-version neighbors can be absorbed.

## Context

Found during REQ-499's post-remediation re-review. The single remediation closed fail-open enumeration, public recovery-to-claim, explicit foreign-section, and stale npm/Cargo/uv cases, but independent boundary fixtures still reproduced an unsafe fold admission and a valid workspace false refusal.

## Instances

- [ ] A valid `## Review Fold — REQ-NNN` or `## Recovery Fold — REQ-NNN` followed by an unrelated unheaded paragraph is accepted because the predicate rejects only another top-level heading.
- [ ] A member-only npm release is tied to an unchanged root manifest and root lock copies when root and member merely share the same old version.

## Requirements

- Prove the exact whole append from a durable originating preimage or an equivalently bounded format; do not infer ownership from heading shape plus absence of another heading.
- Reject foreign bytes before or after the one owned fold, including unheaded paragraphs, comments, malformed headings, and delimiter-shaped additions.
- Identify which workspace manifest actually changed from old to new before deriving required locks and member entries.
- For npm member releases, require only `packages["<member-path>"].version`; unchanged root `package.json`, root lock `version`, and `packages[""].version` must not be selected even when values coincide.
- Apply the same changed-source-first ownership rule to Cargo and uv workspace members and keep enumeration errors typed and fail-closed.
- Preserve strict behavior without `--assume-sole-releaser`, exact protected-path refusal, unrelated-work preservation, and public recovery-to-claim success.

## Red-Green Proof

**RED prompt/case:** Run strict discovery with (a) a valid named fold followed by an unrelated unheaded paragraph and (b) a complete npm member-only release whose unchanged root package shares the old member version.
**Why RED now:** Case (a) reaches owned association and can commit the foreign tail; case (b) refuses on unchanged root copies instead of recognizing the member's complete mirror set.
**GREEN when:** Case (a) refuses byte-identically, case (b) reaches `cleanup_complete` without selecting root copies, equivalent Cargo/uv fixtures pass, and the full existing recovery matrix remains green.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

None.

---
*Source: REQ-499 post-remediation review findings.*

## Triage

**Route: C** — Complex

**Reasoning:** The fix must replace heuristic ownership with bounded preimage provenance and derive release mirrors from the actually changed workspace member across npm, Cargo, and uv while preserving recovery and strict-mode behavior. The current discovery/association matrix needs explicit planning and exploration.

**Planning:** Required

## Plan

1. Replace the open-ended tracked follow-up regex with a closed, line-bounded fold envelope whose opening and terminal marker both bind kind and request ID to the durable HEAD preimage.
2. Record structured dirty manifest transitions first, then derive only those changed sources' declared npm, Cargo, and uv workspace mirrors; treat clean roots and siblings as topology rather than release participants.
3. Keep root lock copies conditional on a changed root source, merge shared-lock descriptors deterministically, and preserve exact replacement plus typed fail-closed enumeration/read/parse behavior.
4. Add tracked-fold rejection/atomicity coverage and member-only, stale-mirror, changed-root, and shared-lock matrices across all three ecosystems, then run focused, race, vet, full-module, and public recovery checks.

**Plan validation:** Every requirement maps to the bounded fold parser, changed-source-first release association, ecosystem mirror derivation, or the strict recovery matrix. CLI schema, mutation/commit layout, unrelated metadata, and assume-mode widening are excluded.

*Generated from delegated exploration; full evidence: `do-work/runs/work-2026-09-03-214500/REQ-512-exploration.md`.*

## Exploration

`followupPathProves` anchors the durable prefix but accepts an unrestricted suffix after a matching heading, so foreign unheaded bytes can be absorbed. Separately, release association scans clean manifests and locks for the common old version before identifying the dirty source, falsely requiring unchanged workspace roots. Ownership must instead flow from the exact bounded append and from a proven changed owned source through declared workspace topology to its exact mirror.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req512_test.go` (new)

**Acceptance criteria:** Foreign bytes outside one closed tracked fold refuse atomically; member-only npm/Cargo/uv releases select only the changed source and its exact mirror; unchanged equal-version roots remain untouched; changed-root controls and all typed strict recovery protections remain green.

## Pre-Flight

**Git:** The shared wave baseline was clean at `b051879c` after claims, briefs, and both Route C exploration artifacts were committed.

**Tests:** Direct canonical fast gate passed and was recorded at the shared wave baseline before source dispatch.

**Dependencies:** REQ-499 is completed and supplies the existing recovery/finalization authority and regression matrix this addendum tightens.
