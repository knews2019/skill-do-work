---
id: REQ-438
title: '[impact-critical] Refuse mismatched Git transaction roots'
status: claimed
route: A
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T21:36:03Z
  basis:
    - trivial short-circuit
related: [REQ-437, REQ-439, REQ-440, REQ-441, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
claimed_at: 2026-09-01T21:35:07Z
---

# Refuse Mismatched Git Transaction Roots

## What

Require a mutating command's supplied repository root to identify the same physical worktree root returned by Git. Refuse nested roots before discovery callbacks can mutate files that preflight, rollback, change detection, and `--commit` inspect at a different path.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reuse root resolution for both preflight and execution, but require the supplied directory and Git worktree root to identify the same physical directory. Add test-first coverage for nested mismatch refusal, exact-root success, and a symlink alias.
- [x] **[APPLY]:** Added one shared physical-root identity guard used by transaction preflight and execution, plus regression coverage at the transaction and cleanup caller seams.
- [x] **[UNIFY]:** Reviewed the complete diff for `git_transaction.go`, `git_transaction_test.go`, and `cleanup_apply_test.go`; verified shared wiring, failure classification, exact-root/alias compatibility, pre-mutation refusal, and no debug artifacts. Focused package tests pass.

## Finding Provenance

- **Verbatim claim / severity:** `[P1] Require the supplied root to match the Git worktree root.`
- **Evidence:** `git_transaction.go:333-349` silently replaces the supplied root with `git rev-parse --show-toplevel`, while cleanup discovery and mutation retain the supplied root.
- **Origin / earned by:** Introduced by transaction foundation `329c55a9`. An isolated nested-root replay returned cleanup success, moved files below `nested/do-work/`, left HEAD unchanged, and left uncommitted changes.
- **Surface-cost:** Earned. The false-success `--commit` incident justifies one transaction-boundary identity check and error; that is cheaper than rebasing every generic callback and target path. Test exact and nested roots, and do not falsely reject a physical alias of the same root.

## Detailed Requirements

- Reject a supplied root that is physically below or otherwise distinct from Git's worktree root.
- Apply the same identity contract to preflight and transaction execution before any mutation callback runs.
- Leave target bytes, index, and HEAD unchanged on refusal and return actionable evidence naming the mismatch.
- Preserve exact-root behavior and legitimate path aliases that resolve to the same physical directory.

## Constraints

- Prefer one shared root-identity check. Do not add per-command target rebasing.
- Preserve the existing installer precedent that mutations require the actual Git project root.

## Red-Green Proof

**RED prompt/case:** In a Git repository with `nested/do-work/`, run cleanup with `--repo-root <repo>/nested --commit` and a mutatable nested REQ.
**Why RED now:** Git guards inspect `<repo>/do-work/...` while cleanup mutates `<repo>/nested/do-work/...`, allowing false success without a commit.
**GREEN when:** The nested-root invocation refuses before its callback, leaves bytes/index/HEAD unchanged, exact-root transactions remain green, and a physical alias of the exact root is handled deliberately.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm on rejection and physical identity; implementation details belong to the builder.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 12 from the validated external feedback.*

## Triage

**Route: A** - Simple

**Reasoning:** The failing boundary and exact shared helper are known, and the requested repair is one identity predicate plus focused tests.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modified)

**What was done:** Git transaction preflight and execution now share a physical-directory identity check that rejects a supplied root distinct from Git's worktree root with typed mismatch evidence. Regression tests cover nested-root refusal without mutation, exact-root success, physical aliases, and cleanup's `--commit` caller seam.

## Qualification

Passed — 3 files verified, 4 requirements traced, P-A-U confirmed. The implementation is substantive, the shared helper is wired into both preflight and execution, and the diff contains no debug artifacts or unrelated changes.

## Testing

**Tests run:** `go test ./internal/gittransaction ./internal/cleanup -count=1`; `go vet ./...`; `go test ./... -count=1`; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing; canonical maintainer verification exited 0.

**Red-green validation:**
- `TestTransactionRootsMustMatchGitWorktreePhysicalIdentity`: ✗ before implementation (nested-root preflight incorrectly succeeded) → ✓ after
- `TestApplyPlanRefusesNestedRepositoryRootBeforeMutation`: ✓ after, proving the captured cleanup `--commit` case refuses with source bytes, index, and HEAD unchanged

**New tests added:**
- Transaction preflight/execution refusal for a nested physical root, with exact-root and symlink-alias acceptance
- Cleanup caller-seam refusal before move or commit

*Verified by work action*

## Review

**Overall: 98%** | 2026-09-01T21:46:23Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — both transaction entry points reject distinct physical roots before mutation, while exact roots and physical aliases remain accepted.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

Git-backed do-work mutations now require the supplied repository root to identify the actual worktree root; the contract lives at the shared `gittransaction` boundary. The do-work CLI prime remains current.
