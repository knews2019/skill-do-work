---
id: REQ-409
title: 'Implement safe cleanup passes and explicit destructive repairs'
status: completed-with-issues
claimed_at: 2026-08-30T19:30:57Z
completed_at: 2026-08-30T20:35:44Z
created_at: 2026-08-29T20:28:26Z
route: C
estimate:
  p50_active_minutes: 90
  confidence: low
  calculated_at: 2026-08-30T19:30:57Z
  basis:
    - Route C
    - 16-file write set
    - 8 new files
    - 5 subsystems involved
    - 5 acceptance criteria
    - dependency depth 3
    - persistence changes
    - cross-route regression gates
    - full-suite verification
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-408]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
kb_status: pending
---

# Implement Safe Cleanup Passes and Explicit Destructive Repairs

## What
Make cleanup a canonical no-LLM command that applies only provably safe repairs by default.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ, UR, action/CLI primes, cleanup policy, shared transaction/result/repository packages, and crew rules; mapped Passes 0–4, explicit repairs, registration, and delegation to 16 files.
- [x] **[APPLY]:** Implemented the cleanup planner, guarded applier, Git/worktree/recovery helpers, handler, registration, shared evidence/rollback extensions, and delegating docs after runnable REDs.
- [x] **[UNIFY]:** Reviewed all 16 scoped files, preserved required contract tokens, confirmed no dependency or later-command drift, and ran focused/race/vet/full/contract/Windows/format checks.

## Detailed Requirements
- Implement cleanup Passes 0–4, including stranded terminal REQ archival, documentation-link repointing, and merged-worktree cleanup.
- Support text/JSON, dry-run, optional commit, and shared rollback/target guards.
- Apply only provably safe repairs by default and report conflicts with evidence and exact next actions.
- Require explicit destructive flags for blanked-record restoration and deletion of unmerged worktrees.

## Constraints
- Default `do-work-cleanup` must be safe and mechanical.
- The screenshot’s `STRANDED-FINISHED-REQ` case must be repairable directly without an LLM.

## Dependencies
Depends on REQ-408 (shared repository model).

## Builder Guidance
Certainty level: Firm. Preserve the existing cleanup pass meanings while moving deterministic decisions and mutations into Go.

## Red-Green Proof
**RED prompt/case:** Put a completed REQ in `do-work/queue/` and run `do-work-cleanup --dry-run`, then run the applying form in a clean Git fixture.
**Why RED now:** Verification reports the stranded REQ but no direct no-LLM cleanup command performs the repair.
**GREEN when:** Dry-run reports the exact archive move, apply performs it safely, JSON gives actionable evidence, and destructive cases remain refused until their explicit flags are supplied.
**Validation:** User confirmed via the supplied implementation plan and original stranded-REQ example.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Triage

**Route: C** - Complex

**Reasoning:** This command family combines five cleanup passes, Git-guarded mutations and rollback, archive/path rewriting, worktree lifecycle handling, stable text/JSON evidence, and explicit destructive boundaries.

**Planning:** Required

## Plan

Build the cleanup family in four ordered units, keeping mutation policy in the existing transaction/result layers.

1. **Repository evidence and operation planning:** Extend the shared snapshot for active/archived URs and run manifests, then add a cleanup package that plans Passes 0–4 as exact, conflict-aware operation groups. Membership comes from scanned `user_request:` fields; symlinks are never followed; consumed runs alone are deletable; link rewrites stay coupled to their moves.
2. **Grouped Git-guarded application:** Extend `gittransaction` with public target preflight and exact created-directory rollback. Recheck every eligible group immediately before applying; use rename/atomic replace/exact removal; never overwrite. Dry-run produces the same proposed changes without writes, and commit reuses empty-index, exact-stage, rollback, and post-commit-risk behavior.
3. **Explicit repairs and command handler:** Register `cleanup` with `--dry-run`, `--commit`, repeatable exact-target `--restore-blanked`, and repeatable exact-name `--discard-worktree`. Default discovery reports blanked records and unsafe worktrees without deleting them. Restore blanked content atomically from Git history; remove merged clean worktrees mechanically; require an exact explicit selection for unmerged discard. Render through `resultmodel.CommandResult` only.
4. **Delegating action and docs:** Make the natural-language cleanup action invoke the canonical CLI and stop on tooling failure, update its guide and the CLI prime, retain the legacy blank scanner only as a compatibility surface until REQ-420, and leave flat Just wiring to REQ-419.

**Planned files (16; 8 new):**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/actions/cleanup.md` (modified)
- `skills/do-work/docs/cleanup-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**Requirement coverage:** Passes 0–4, link repointing, and safe groups map to tasks 1–2; text/JSON, dry-run, commit, rollback, explicit repairs, and worktree behavior map to tasks 2–3; the action-delegation requirement maps to task 4. Flat recipes and wholesale shim conversion remain in REQ-419/REQ-420.

**Testing:** RED-before-GREEN fixtures cover stranded terminal records, aliases/checkpoint ownership, UR closure, archive consolidation, misplaced trees, consumed runs, link anchors, dirty groups, rollback, blank restoration, merged/unmerged worktrees, exact destructive targeting, commit path verification, and text/JSON parity. Then run focused, race, vet, full-module, Windows compile-only, and the unpiped canonical maintainer gate.

**Plan validation:** All Detailed Requirements map to at least one task; every task maps back to the REQ. The exact destructive flag spelling and `--commit` interaction require exploration, but neither changes the safe-default boundary.

*Generated by Plan agent*

## Exploration

The 16-file boundary is viable with several load-bearing corrections. `repositorymodel` already owns active and archived UR inputs, so it only needs rooted run-manifest evidence. Cleanup must filter membership to the documented live/archive locations, skip `archive/hold`, and exclude canonical/shipped/nested-repository `do-work` directories plus symlinked trees from Pass 3a.

Pass 0 couples a lossless terminal-alias rewrite and working-file move with removal of only this checkout's exact checkpoint writer entry. Each move is modeled as an old tracked target plus a new absent target; operation groups preflight independently, while the union of eligible exact file targets runs through one transaction and at most one commit. Created directories are recorded explicitly and rolled back empty/deepest-first, never recursively.

Blanked recovery evidence is resolved at the original path before any same-run archive move. Default blank/worktree findings remain report-only; repeatable exact-target flags are the consent tokens. Automatic worktree removal requires a clean tree and ancestry proof against the integration HEAD. The binary command token is `cleanup`; the `do-work-cleanup` Just recipe belongs to REQ-419. Existing contract tokens in the action/guide remain intact.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/actions/cleanup.md` (modified)
- `skills/do-work/docs/cleanup-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**Files I will NOT touch:** flat Just recipes (REQ-419), retained shell-to-shim conversion (REQ-420), doctor/forensics behavior (REQ-410), request-state commands (REQ-412), `queue-kanban`, module dependencies, and unrelated queue files.

**Acceptance criteria:**
- [ ] Dry-run and apply cover cleanup Passes 0–4 with exact evidence, conflicts, and stable ordering.
- [ ] Stranded terminal REQs are normalized and archived safely; working moves remove only this checkout's checkpoint entry.
- [ ] UR closure derives membership from `user_request:` evidence and respects terminal-resolved semantics.
- [ ] Documentation links are repointed with anchors preserved and bare prose untouched.
- [ ] Dirty/conflicting groups are refused without blocking unrelated safe groups; rollback removes only invocation-created files/directories.
- [ ] Text/JSON, optional commit, exact path verification, and actionable next commands use the shared result/transaction layers.
- [ ] Blanked restoration and unmerged/dirty worktree discard require repeatable exact-target flags; default cleanup reports only.
- [ ] Merged clean `worktree-agent-*` leftovers are removed only with cleanliness and ancestry proof.

**Portability/safety boundary:** destination absence, rooted containment, regular-file identity, slash-normalized evidence, `git ... --porcelain -z`, and Windows compile-only coverage are mandatory. No recursive rollback deletion or constructed shell command is permitted.

## Decisions

### D-01: Move resolved UR members directly to their final archive folder

**Decision:** DECIDE & STATE — when all scanned members are terminally resolved, stranded and loose members target `archive/UR-NNN/` directly rather than using an intermediate loose-archive move in the same transaction.

**Reasoning:** A single source cannot safely participate in two ordered moves under one exact-target snapshot; direct final destinations preserve the pass outcome without an unstable intermediate path.

### D-02: Preflight groups independently and commit their eligible union once

**Decision:** DECIDE & STATE — dirty or conflicting groups are reported and omitted while unrelated safe groups continue; all eligible exact targets share one transaction and at most one optional commit.

**Reasoning:** Cleanup's default promise is “apply only provably safe repairs,” not “one unsafe item blocks every repair” or “one commit per pass.”

### D-03: Publish move destinations exclusively through rooted handles

**Decision:** DECIDE & STATE — require destination absence, create the destination through a validated rooted directory handle, verify source content again, then delete the source.

**Reasoning:** `os.Rename` overwrite behavior varies by platform and does not itself express the no-overwrite contract.

### D-04: Restore blanked records atomically from Git evidence

**Decision:** DECIDE & STATE — resolve the newest recoverable nonblank version at the original path, assemble its provenance field before publication, and replace the blank record as one guarded transaction.

**Reasoning:** Separating content recovery from provenance repair recreates the partial-repair state the legacy scanner warns against.

### D-05: Treat exact destructive targets as consent tokens

**Decision:** DECIDE & STATE — blank restoration and forced worktree discard accept repeatable exact paths/names only; default discovery remains report-only.

**Reasoning:** An “all” flag would silently broaden destructive authority between planning and application. Exact current evidence keeps consent reviewable and pasteable.

### R-D-01: Separate recovery source from implementation provenance

**Decision:** DECIDE & STATE — full-history recovery tracks both the commit supplying nonblank bytes and the implementation hash recorded by the later metadata commit; only the latter is written to `commit:`.

**Reasoning:** The recovery blob's commit identifies where bytes were found, not necessarily the implementation those bytes are meant to attribute.

### R-D-02: Keep held records outside active UR membership

**Decision:** DECIDE & STATE — `archive/hold` remains repository evidence but does not participate in cleanup's four-location active membership predicate; zero active members satisfy the every-member closure rule.

**Reasoning:** Held records are intentionally outside ordinary closure, while an empty captured UR has no unresolved member that can justify keeping it open.

### R-D-03: Give misplaced items independent conflict domains

**Decision:** DECIDE & STATE — each misplaced file or UR folder plans independently, including legacy `CONTEXT-*`, so one occupied destination cannot block unrelated safe relocations.

**Reasoning:** The cleanup contract says conflicts stop at the conflicting item, not at the nearest shared source directory.

### R-D-04: Preserve transaction and revert truth in command results

**Decision:** DECIDE & STATE — proposed changes are relabeled only after the transaction outcome is known, and committed-risk results carry the transaction's exact `git revert <sha>` argv.

**Reasoning:** A rolled-back move labeled “applied” or a committed-risk result without its recovery command is actionable evidence corruption.

### R-D-05: Prove detached worktrees by recorded HEAD

**Decision:** DECIDE & STATE — parse NUL-delimited worktree records, retain detached HEADs for ancestry checks, and reject discard tokens that do not match the current enumeration; `--commit` and discard are separate runs.

**Reasoning:** Branch absence is not proof of unmerged work, and a typo cannot count as destructive consent.

### R-D-06: Isolate the consumed-scratch non-rollback exception

**Decision:** DECIDE & STATE — general dirty-target guards remain unchanged. Only an entirely untracked Pass-4 group whose exact rooted inventory and `Status: consumed` manifest still match is deleted outside the rollback transaction and labeled as non-rollback spent scratch.

**Reasoning:** Consumed run output is explicitly expendable scratch, but pretending Git can restore untracked bytes would be dishonest. Mixed or durable dirty targets still fail closed.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (new)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/actions/cleanup.md` (modified)
- `skills/do-work/docs/cleanup-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**What was done:** Added the canonical `cleanup` command and its safe planning/application layer for Passes 0–4, link repointing, exact-target Git transactions, blank-record recovery, and worktree evidence/repairs. Registered the command and changed the natural-language action, guide, and CLI prime to delegate deterministic cleanup to it.

## Qualification

**Attempt 1: Passed mechanically.** All 16 files exist, are substantive, and match Scope exactly; P-A-U is complete and the command is registered through the shared runtime. The static-reference warnings are expected for new same-package source/test files. Requirements trace to cleanup planning/application, shared repository/transaction seams, and delegation surfaces.

**Attempt 2 after review remediation: Passed.** The same 16 files remain exactly scoped and substantive. Regression coverage now traces every review finding to a corrected behavior: provenance, UR/hold/empty semantics, legacy/partial Pass 3, transaction truth, worktree consent/enumeration, consumed scratch, composed links, and canonical-root preservation.

## Testing

**Tests run:** focused cleanup/repositorymodel/gittransaction/cmd tests; focused race tests; `go vet ./...`; uncached full-module tests; contract regressions; Windows cleanup and atomicfile compile-only; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All executed gates pass, including the canonical repository gate. Its optional browser lane remained in the documented default-skipped state.

**Red-green validation:**
- Caller seam: `cleanup --dry-run` initially returned `UNKNOWN-COMMAND`/exit 2 → ✓ registered command returns a shared result.
- Planner: a tracked doc link was initially absent from its move group → ✓ single-move link repointing passes.
- Repository discovery: a REQ-shaped file under run scratch was initially classified as stray → ✓ rooted run-manifest evidence ignores run contents as queue state.
- Publication: an existing destination was initially overwritten by the planned move path → ✓ exclusive no-overwrite publication passes.
- Review remediation: live blank restores were missed and used the wrong provenance; zero-member/held URs planned incorrectly; CONTEXT and partial relocations were absent; rollback said “applied”; worktree records truncated/ignored exact targets; untracked consumed scratch was refused; two link rewrites competed; and canonical roots were pruned → ✓ each named regression now passes with truthful result evidence.

**New tests added:** four cleanup package test files; repository-model run-manifest coverage; Git-transaction preflight and exact created-directory rollback coverage.

*Verified by work action*

## Review

**Overall: 50%** | 2026-08-30T20:35:44Z

| Dimension | Score |
|-----------|-------|
| Requirements | 68% |
| Code Quality | 68% |
| Test Adequacy | 74% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- UR closure is planned independently from required member archival, so a refused member can leave an archived active UR — impact-user-visible → REQ-430 created
- Composed documentation rewrites are attached to the first move group instead of the move that makes each destination valid — impact-user-visible → REQ-431 created
- The consumed-scratch exception bypasses `--commit`'s empty-index guard — impact-user-visible → REQ-432 created
- Pass 3b gives a whole misplaced UR folder one conflict domain, preventing required partial merges — impact-user-visible → REQ-433 created

**Minor findings:** 1 (report only)
**Acceptance:** Fail — remediation resolved the original eight findings, but four end-to-end coupling defects still violate safe-cleanup and exact-guard requirements.
**Suggested testing:** 4 items
**Follow-ups created:** REQ-430, REQ-431, REQ-432, REQ-433; **sweeps appended to:** None

*Reviewed by review-work action*

## Remediation

**Attempt 1:** The initial review failed at 50% with eight Important findings spanning provenance, UR membership, misplaced-tree coverage, transaction truth, worktree guards, consumed scratch, composed links, and canonical-root preservation.

**Attempt 2 (single allowed remediation):** Corrected and regression-tested those eight findings, then reran qualification, focused/race/vet/full/Windows/contract/canonical gates. Re-review still failed at 50% because operation dependencies remain incorrectly coupled in four user-visible cases: UR closure vs member archival, documentation rewrites vs owning moves, consumed scratch vs the commit guard, and per-item partial merges. REQ-430 through REQ-433 record the mandatory remaining work.

## Lessons Learned

**What worked:** Exact rooted evidence, transaction-result relabeling after outcome, and focused adversarial fixtures closed the original safety defects without broadening the write set.

**What didn't:** Modeling each move as locally safe was insufficient; safety also depends on explicit prerequisites between operation groups. Unit coverage of isolated groups missed refusal combinations until end-to-end re-review.

**Worth knowing:** Cleanup planners need a dependency graph, not just deterministic ordering: every derived mutation must name the successful operation that makes it valid, while unrelated groups remain independently eligible.

## Orientation

[MAP CHANGED] The do-work CLI now contains the canonical cleanup command for safe Passes 0–4, explicit blank/worktree repairs, shared Git-guarded application, and delegating action/docs. It is archived with known interaction defects tracked by REQ-430 through REQ-433; the action-file and CLI primes still point to existing canonical paths.
