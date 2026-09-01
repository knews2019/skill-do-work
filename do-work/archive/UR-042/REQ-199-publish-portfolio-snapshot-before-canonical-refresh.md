---
id: REQ-199
title: Publish portfolio snapshot before canonical refresh
status: completed
claimed_at: 2026-08-15T19:41:49Z
route: B
completed_at: 2026-08-15T20:01:55Z
commit: 74f2220
kb_status: promoted
kb_entry: REQ-199-publish-portfolio-snapshot-before-canoni.md
domain: general
created_at: 2026-08-15T17:09:44Z
user_request: UR-042
addendum_to: REQ-190
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
maintenance: true
write_set: [skills/do-work-toolbox/actions/present-work.md, skills/do-work-toolbox/docs/present-work-guide.md, skills/do-work-toolbox/scripts/publish-portfolio-summary.sh, skills/do-work/docs/prescribed-shell-primitives.md, _dev/tests/prescribed-shell-canonicalization.sh, _dev/tests/prescribed-shell-scripts-behavior.sh, _dev/tests/staged-skills-contract.sh, _dev/tests/contract-regressions.sh]
---

# Review Fix: Publish Portfolio Snapshot Before Canonical Refresh

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Explore existing exclusive-create and atomic-replacement patterns, promote portfolio publication into one helper, and replay every branch/failure before changing action prose.
- [x] **[APPLY]:** Added helper/inventory/semantic/behavior assertions first, confirmed RED, then implemented snapshot-first and canonical-only transactions within the eight-file scope.
- [x] **[UNIFY]:** Reviewed the complete helper/caller/guide/registry/test diff and passed all four focused suites, shell syntax/lint, and diff checks.

## What

Make the Yes/unavailable snapshot branch publish its exclusive, no-clobber snapshot from the retained bytes before atomically refreshing the canonical portfolio file. A snapshot publication failure must leave the prior canonical summary unchanged instead of partially completing the promised two-output branch.

This is a standalone user-visible publication-order contract and cannot fold into the existing image-generation or ID-normalization follow-ups: its fix applies only to the portfolio writer's canonical-plus-snapshot transaction.

## Context

Found during review of REQ-190. The new action resolves the snapshot name before writes but tells the agent to refresh the canonical file before publishing the exclusive snapshot, allowing a late collision/publication failure to leave only the canonical output.

## Requirements

- In the Yes/unavailable branch, exclusively publish the snapshot first from the retained byte sequence.
- Only after snapshot success, atomically refresh the canonical file from those same bytes.
- In the No branch, atomically refresh only the canonical file.
- Preserve collision suffixing, no-clobber snapshots, byte identity, artifact immutability, and no automatic cleanup.
- Add or identify replayable contract assertions for snapshot-publication failure ordering and both successful branches.

## Red-Green Proof

**RED prompt/case:** Inspect or replay the Yes/unavailable branch with a snapshot candidate that becomes occupied or cannot be exclusively published after path resolution; current ordering can refresh the canonical file before snapshot failure.
**Why RED now:** The branch promises canonical plus snapshot output but can partially complete while destroying the previous canonical state.
**GREEN when:** Snapshot publication happens first, failure leaves the prior canonical file unchanged, success atomically refreshes canonical from identical bytes, and the No branch still refreshes canonical only.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Triage

**Route: B — Medium**

**Reasoning:** The broken ordering is clear, but exploration must identify the canonical exclusive-create and atomic-refresh primitives plus the executable replay seam before scope is declared.

**Planning:** Not required — exploration-guided implementation.

## Exploration

- The defect is isolated to portfolio publication, but safe executable ownership belongs in a shipped helper rather than increasingly detailed action prose.
- Use one retained source, a verified private temp adjacent to canonical, hard-link exclusive snapshot publication with numeric collision suffixing, then atomic canonical replacement. In the snapshot branch the immutable output must succeed first; the No branch performs only canonical replacement.
- Replay canonical-only success, snapshot success/byte identity, a late occupied candidate that advances to a suffix, exclusive-publication failure that preserves old canonical, and canonical-replacement failure that retains the already-published snapshot.
- Register the helper in the canonical prescribed-shell home/inventory and staged package contract; keep the action/guide at intent, invocation, and observable outcome altitude.
- Preserve exact dispatch/question/fallback, artifact immutability, no automatic cleanup, and unrelated presentation ownership. Exclude router/archive/schema/ai-report/video changes.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/present-work.md` (modify) — invoke the helper and define snapshot-before-canonical outcomes
- `skills/do-work-toolbox/docs/present-work-guide.md` (modify) — document observable publication order/failure behavior
- `skills/do-work-toolbox/scripts/publish-portfolio-summary.sh` (new) — own exclusive snapshot and atomic canonical publication
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — register helper mechanics
- `_dev/tests/prescribed-shell-canonicalization.sh` (modify) — include the executable home
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — replay all branch and failure outcomes
- `_dev/tests/staged-skills-contract.sh` (modify) — verify shipped helper/caller resolution
- `_dev/tests/contract-regressions.sh` (modify) — ratchet semantic order and reject canonical-first regrowth

**Files I will NOT touch:** routers; archive/schema; ai-report/video; prior artifacts; cleanup; version/lifecycle until integration; or unrelated helpers/tests.

**Acceptance criteria:**
- [x] Yes/unavailable publishes an exclusive snapshot before any canonical refresh; failure leaves old canonical unchanged.
- [x] No atomically refreshes canonical only.
- [x] Successful snapshot and canonical outputs are byte-identical to the one retained source.
- [x] Collision suffixing, no-clobber immutability, and no automatic cleanup remain intact.
- [x] Replayable RED/GREEN focused gates pass; canonical verification remains the integration gate.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `skills/do-work-toolbox/docs/present-work-guide.md` (modified)
- `skills/do-work-toolbox/scripts/publish-portfolio-summary.sh` (new)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Promoted portfolio output into a shipped helper that verifies one retained source, publishes an optional snapshot exclusively before canonical replacement, handles numeric collisions, and atomically refreshes canonical. Wired the action/guide, registered the primitive, and added branch/failure behavior replays plus semantic/inventory contracts.

## Testing

**Tests run:** baseline four focused suites; test-only RED four-suite run; GREEN prescribed-shell canonicalization, prescribed-shell behavior, staged-skills contract, and contract regressions; helper/test `bash -n`; shell-block lint; `git diff --check`

**Result:** ✓ All focused checks passed; behavior suite now reports 27 named cases.

**Red-green validation:** ✗ current tree lacked the helper, executable inventory, order contract, and branch behavior → ✓ canonical-only, snapshot-first success, numeric collision, snapshot failure, and canonical-refresh failure cases all pass with byte/immutability checks.

*Verified by work action*

## Qualification

Passed — eight scoped files are substantive and wired; Scope matches Implementation Summary; P-A-U is complete; helper mode/argument/data flow is explicit; every required branch and failure is behavior-tested; executable/inventory/pointer contracts resolve; syntax, ShellCheck-backed lint, and diff checks pass.

## Review

**Overall: 65%** | 2026-08-15T20:00:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 72% |
| Code Quality | 64% |
| Test Adequacy | 63% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings:**
- Hard-linking the snapshot and renaming the same inode to canonical means a later in-place canonical edit mutates the supposedly immutable snapshot; the current test explicitly locks same-inode identity. — gate: user-visible → consolidated into generation-2 pending-answers REQ-205.
- Unvalidated directory operands let two-operand `ln`/`mv` nest private files while reporting the directory itself as the published exact path. — gate: user-visible → consolidated into generation-2 pending-answers REQ-205.

**Minor findings:** 0
**Acceptance:** Partial — original snapshot-first ordering and failure preservation are credible for ordinary file paths; independent immutable bytes and exact regular-file destinations remain.
**Suggested testing:** post-publication canonical mutation and pre-existing directory-destination probes.
**Follow-ups created:** REQ-205; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Promoting publication into a shipped helper made branch order and failure outcomes directly executable and replayable.
- Exclusive snapshot publication before canonical replacement closes the original partial-update defect on ordinary paths.

**What didn't:**
- Hard-link identity proved byte equality at publication but coupled later mutable canonical writes back into the durable snapshot.
- Two-operand `ln` and `mv` interpret directory destinations as containers; without exact-path type guards, successful status does not prove the requested path was published.

**Worth knowing:** Immutable evidence needs independent file contents, not merely a second directory entry for the same inode. Publication helpers must test exact target type/identity because core utilities treat directories differently from file destinations.

**Knowledge handoff:** Pending human consent. No knowledge-base file was written automatically.

## Orientation

Portfolio publication now has one shipped helper and executes the preservation branch in the safe high-level order: snapshot first, canonical second. REQ-205 holds the consent-gated hardening needed for independent immutable files and exact destination semantics.
