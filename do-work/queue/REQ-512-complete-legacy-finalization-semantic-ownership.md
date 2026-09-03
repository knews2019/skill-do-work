---
id: REQ-512
title: '[impact-critical] Review fix: Complete legacy finalization semantic ownership'
status: pending
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
