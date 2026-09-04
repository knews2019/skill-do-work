---
id: REQ-513
title: '[impact-rule-change] Commit the claim footprint in every mode'
status: completed
created_at: 2026-09-02T20:35:18Z
user_request: UR-099
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
route: C
planning_at: 2026-09-02T20:47:05Z
dispatch_at: 2026-09-02T21:01:09Z
builder_handback_at: 2026-09-02T21:09:42Z
integration_at: 2026-09-02T21:09:42Z
review_at: 2026-09-02T21:25:40Z
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-09-02T20:45:56Z
  basis:
    - Route C
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
related: [REQ-514, REQ-515, REQ-516, REQ-517]
batch: recovery-never-traps
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/docs/work-guide.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go]
claimed_at: 2026-09-02T20:45:12Z
completed_at: 2026-09-02T21:33:21Z
release_at: 2026-09-02T21:33:21Z
commit: 33852cb4a9c0e8af197d789fea6f2624beb68ffe
kb_status: promoted
kb_entry: REQ-513-commit-the-claim-footprint-in-every-mode.md
---

# Commit the claim footprint in every mode

## What

Make Step 2's claim commit its own footprint in every mode, serial included: the queue-to-working move plus the checkpoint entry land as one bookkeeping commit at claim time. The CLI's `claim --commit` already exists; `actions/work.md` Step 2 does not use it, and worktree dispatch Step 0 stages the same moves by hand later.

The fold-first scan found no pending or pending-answers REQ in any UR that owns this claim-commit asymmetry; REQ-505 (Move selection and claim behind advance) moves the claim later but keeps its current commit behavior.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add a failing action-contract assertion for claim-time commit, pin the existing CLI transaction's exact Git footprint, replace deferred-claim prose with the canonical `claim --commit` contract, and remove the obsolete hand-back staging path.
- [x] **[APPLY]:** Updated the action contracts and user guide, then added the planned action-seam plus request-state behavior coverage within the declared five-file scope.
- [x] **[UNIFY]:** Reviewed `git diff --stat` and the full diff for both action files, the contract suite, and `state_apply_test.go`; ran `git diff --check`, `go vet ./...`, focused request-state tests, and the contract suite; found no introduced debug artifacts.

## Why

Serial mode defers the claim commit to Step 9 while worktree mode commits it at dispatch. That asymmetry is the whole bug behind the REQ-456 trap: the journaled finalizer only ever meets a dirty checkpoint in serial mode, and its complete transition treats that dirt as foreign. Once claim always commits, complete never meets dirt it made, and no special-case acceptance is needed. This is the delete-before-add fix.

## Context

Origin: REQ-456 (Wait for theme transitions before contrast measurement) finished build and review, then `complete` prepared its journal and refused at lifecycle apply with `FINALIZATION-LIFECYCLE-APPLY`, "target path do-work/CHECKPOINT.md is already dirty". The previous Route A completion, REQ-440, avoided it with a hand-made "[do-work] Record REQ-440 claim checkpoint" commit. REQ-499 through REQ-501 were Route C and never saw it because dispatch committed the claim first. The unblock for REQ-456 was the same hand-made commit, `cd9b01b0`.

## Detailed Requirements

- `actions/work.md` Step 2 invokes `claim` with `--commit` in every mode; the commit message shape is the CLI's own.
- Worktree dispatch Step 0 (`actions/work-reference.md`) stops staging claim moves and the checkpoint by hand, since the claim commit already landed; keep its guard for unrelated staged paths.
- A dirty checkpoint at claim time is shared-target dirt: the claim's refusal stands, and its finding names a verb other than `claim` (REQ-514 owns the general invariant; this REQ only makes sure the claim path has one).
- Contract predicates that pin the deferred-commit wording are deleted with the prose; a behavior test on `claim --commit` covers the serial shape.
- `git log` after a serial run shows one claim commit per REQ, followed by the implementation commit, as worktree mode already does.

## Constraints

- Do not add a second commit surface; use the existing `claim --commit` flag.
- Serial and worktree modes end with the same commit shape for the claim.
- Do not touch the finalizer here; REQ-516 owns what recovery accepts under the sole-authority assertion.

## Batch Constraints

- Judgment stays prose; mechanics stay in the Go CLI. No new prose that walks a shell sequence.
- A guard may still refuse. What it may not do is refuse for a REQ-scoped reason in a way that stops unrelated REQs, or name itself as the fix.
- Nothing here widens recovery to secret-classified or project paths; only dirt the pipeline itself wrote earlier in the run is in scope.
- Every REQ carries a behavior test on the command or a contract predicate on the action, never a sentence pin alone.

## Dependencies

REQ-517 (Pin the serial claim-to-recovery trap) depends on this REQ. Related to REQ-505, which later moves the claim behind `advance` and must inherit the commit.

## Builder Guidance

Certainty level: Firm. The flag exists; the work is wiring, predicate cleanup, and one behavior test. Read `_dev/primes/prime-action-files.md` before touching an action file.

## Red-Green Proof

**RED prompt/case:** Run a serial Route A REQ through Step 2 and inspect `git status` before Step 3.
**Why RED now:** The queue file shows as deleted, the working file as untracked, and `do-work/CHECKPOINT.md` as modified; nothing committed the claim.
**GREEN when:** After Step 2 the tree is clean for those three paths and `git log -1` is the claim's bookkeeping commit; the same is true when dispatch runs the claim in worktree mode.
**Validation:** User confirmed (verify-requests, 2026-09-02).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3539 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing action routing and status contracts.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on semantic recovery completeness and structured evidence projection in do-work-cli internals.

## Full Context

See `do-work/user-requests/UR-099/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-02, item A1 of "how can I update the orchestrator to not end up in a trap like this?", captured by UR-099.*

---

## Triage

**Route: C** - Complex

**Reasoning:** This changes the claim lifecycle contract across serial and worktree modes, removes stale orchestration prose, and adds command-level regression coverage.

**Planning:** Required

## Plan

1. Add RED coverage for the action-to-CLI seam and exact `claim --commit` Git effects, while removing predicates that preserve delayed claim bookkeeping.
2. Make Step 2 invoke the existing commit mode and state its clean-tree postcondition.
3. Update the execution model and worktree hand-back instructions so they consume the already-committed claim instead of staging it later.
4. Run focused contract and Go tests, then the canonical repository gate and independent review.

Validation: all five Detailed Requirements map to these tasks; every task traces to the claim-commit asymmetry; four tasks stay below the split warning threshold; no command output drives a new action-owned mutation.

*Generated inline under the Plan-agent fallback*

## Exploration

- `internal/requeststate` already accepts `--commit` and commits only the request move plus checkpoint target with message `[REQ-NNN] claim request lifecycle`.
- `actions/work.md` Step 2 omits the flag, while its condensed worktree hand-back still stages the resulting dirt later.
- `actions/work-reference.md` has two stale restatements: the execution model says claims never commit at claim time, and worktree Step 0 treats the claim as a staged rename.
- `_dev/tests/contract-regressions.sh` already isolates both Step 2 and hand-back Step 0, making it the narrow action-to-command seam for RED coverage.

*Generated inline under the Explore-agent fallback*

## Pre-Flight

**Git:** ✓ Working tree clean outside `do-work/`
**Tests baseline:** ⚠ The pre-flight wrapper observed transient `internal/corehelpers` failures; the required direct unpiped `bash _dev/tests/maintainer-verify.sh` rerun passed in full before implementation.
**Dependencies:** ✓ Go and ShellCheck floors available

*Checked by work action*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — require committed claims and simplify hand-back Step 0
- `skills/do-work/actions/work-reference.md` (modify) — align execution-model and worktree contracts
- `skills/do-work/docs/work-guide.md` (modify) — remove the user-facing delayed-claim restatement
- `_dev/tests/contract-regressions.sh` (modify) — pin the action seam without preserving deferred bookkeeping
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modify) — prove the existing commit flag's exact serial Git footprint

**Files I will NOT touch:** finalization internals, request-state production code, or sibling recovery actions owned by later UR-099 REQs.

**Acceptance criteria (restated from REQ):**
- [x] Step 2 invokes `claim --commit` in serial and worktree modes.
- [x] Claim atomically commits the queue move and checkpoint entry with the CLI-owned message.
- [x] Worktree hand-back does not stage claim moves or the checkpoint.
- [x] Shared-target dirt still refuses the claim and points away from retrying `claim`.
- [x] Serial history contains one claim commit before implementation/finalization commits.

## Decisions

### D-01 — Extend scope to the user-facing work guide

**Decision:** Add `skills/do-work/docs/work-guide.md` to Scope after the review restatement sweep found two active delayed-claim statements.

**Reasoning:** The REQ requires deleting the deferred-commit prose with its predicates. Leaving the primary user guide on the old timing would ship a split contract, so this is requirement completion rather than adjacent documentation cleanup.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)

**What was done:** Step 2 now commits the canonical claim transaction in every mode, the execution model and user guide describe the same immediate claim footprint, the worktree hand-back no longer stages claim/checkpoint dirt later, and tests pin both the action seam and the CLI transaction's exact commit/refusal behavior.

## Qualification

Passed — 5 files verified, 5 requirements traced, P-A-U confirmed. The action command, canonical execution model, user guide, hand-back staging boundary, exact Git footprint, commit ordering, and dirty-checkpoint refusal all have direct diff or test evidence.

## Testing

**Tests run:** `cd skills/do-work/tools/do-work-cli && go test -count=1 ./internal/requeststate`; `bash _dev/tests/contract-regressions.sh`; `cd skills/do-work/tools/do-work-cli && go vet ./...`; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing

**Red-green validation:**
- Step 2 action-to-CLI seam: ✗ contract suite failed with `Step 2 must invoke ... --commit` before the action edit → ✓ full contract suite passed after the edit.
- Claim commit behavior: ✓ request-state tests prove the exact three-path commit, clean post-claim tree, CLI-owned subject, preserved foreign checkpoint entries, and dirty-checkpoint refusal with a non-`claim` remedy.

**New tests added:**
- Action contract assertion for the canonical `claim --commit` invocation and removal of deferred claim/checkpoint staging.
- `TestClaimCommitLandsExactFootprintAndLeavesCleanTree`.
- `TestClaimCommitRefusesDirtyCheckpointWithExternalRemedy`.

*Verified by work action*

## Review

**Overall: 100%** | 2026-09-02T21:25:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — the action invokes the existing atomic claim commit, both execution modes share the clean boundary, and focused plus repository-wide tests verify the exact Git and refusal behavior.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Testing the action-to-CLI seam and the CLI's exact Git footprint together caught both wiring drift and transaction regressions.
- The restatement sweep found the user guide's delayed-commit description before release.

**What didn't:**
- The first guide assertion passed an absolute path to a repo-relative test helper, so the canonical gate rejected the doubled path; using the helper's expected relative path fixed it.

**Worth knowing:**
- `claim --commit` already owned the correct transaction. The fix was to invoke that owner everywhere and delete the later hand-back staging contract, not add another commit path.

## Orientation

[MAP CHANGED] Every work run now begins implementation from a committed request claim: the queue move and checkpoint entry land atomically in the request-lifecycle subsystem, and serial plus worktree execution consume the same clean boundary. This removes a mode-specific hand-back stage and keeps the finalization subsystem from encountering dirt created earlier by its own pipeline.
