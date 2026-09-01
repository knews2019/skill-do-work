---
id: REQ-444
title: 'Refuse untracked consumed scratch in cleanup commit mode'
status: completed
route: A
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T22:42:49Z
  basis:
    - one existing preflight exception plus focused regressions
related: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-441, REQ-442, REQ-443]
batch: accepted-feedback-regressions
claimed_at: 2026-09-01T22:42:19Z
completed_at: 2026-09-01T22:51:43Z
commit: 2311e76e
---

# Refuse Untracked Consumed Scratch in Cleanup Commit Mode

## What

Refuse entirely untracked consumed-scratch deletion whenever cleanup runs with `--commit`. Commit mode must not report success after a scratch-only deletion that leaves HEAD unchanged or silently perform scratch deletion outside an otherwise valid tracked cleanup commit.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** At the existing untracked-consumed-scratch preflight exception, refuse the group when `--commit` is active and report its exact inventory plus a non-commit remediation. Add test-first scratch-only and mixed tracked/scratch cases; retain the existing non-commit deletion and dirty-index guards.
- [x] **[APPLY]:** Refused untracked consumed-scratch groups at the existing preflight exception whenever commit mode is active, with exact inventory evidence and non-commit remediation; added scratch-only and mixed regressions.
- [x] **[UNIFY]:** Reviewed `cleanup_apply.go` and `cleanup_apply_test.go` in scoped diff context; verified the refusal is confined to commit mode, tracked groups remain independently eligible, scratch bytes stay intact, HEAD behavior is truthful, and no debug artifacts remain.

## Finding Provenance

- **Verbatim claim / severity:** `[P2] Honor --commit before deleting untracked scratch.`
- **Evidence:** `cleanup_apply.go:44-49` admits scratch after preflight failure; when no tracked group is eligible, the transaction is skipped and scratch deletion still runs, so success can leave HEAD unchanged.
- **Origin / earned by:** The explicit non-rollback exception in `a57bf51e` did not define a commit-mode boundary. Scratch-only and mixed tracked/scratch replays show deletion outside the requested exact-path commit.
- **Surface-cost:** Earned. One `options.Commit` refusal and two regressions are cheaper than a successful destructive side effect that the requested commit and rollback cannot cover.
- **Fold-first result:** REQ-432 (Enforce the Commit Guard for Consumed Scratch Cleanup) owns the related nonempty-index incident but is dependency-gated through REQ-431, so the canonical fold-first rule forbids widening it and requires this independently runnable REQ.

## Detailed Requirements

- Refuse every entirely untracked consumed-scratch group when `options.Commit` is true, even with an empty index.
- Cover both a scratch-only run and a mixed run containing otherwise eligible tracked cleanup work.
- Leave the scratch inventory byte-identical and report truthful refusal evidence.
- Preserve the existing exact-inventory, rooted-containment, and consumed-manifest checks.
- Preserve non-commit cleanup's narrow, explicitly labeled non-rollback scratch deletion.
- Leave REQ-432's separate global empty-index invariant intact.

## Constraints

- Add one commit-mode boundary to the existing exception; do not add a new transaction class for untracked scratch.
- Shared files with REQ-432 do not create a dependency; this request must remain independently selectable because REQ-432 is currently gated.

## Red-Green Proof

**RED prompt/case:** With an empty index, run `cleanup --commit` first on an entirely untracked consumed run alone, then on that scratch beside one eligible tracked cleanup group.
**Why RED now:** The scratch exception runs outside the transaction; scratch-only cleanup succeeds without a commit, and a mixed run deletes scratch outside its tracked commit.
**GREEN when:** Both commit-mode fixtures refuse and preserve the scratch bytes, tracked commit behavior remains truthful, and the same scratch is still deleted by an otherwise identical non-commit cleanup run.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. The accepted remedy is the narrow `options.Commit` refusal, not an attempt to make Git roll back untracked scratch.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 16 from the validated external feedback.*

## Triage

**Route: A** - Simple

**Reasoning:** The cleanup engine already identifies the exact untracked consumed-scratch exception; commit mode needs one refusal branch there, plus two focused execution regressions.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

`ApplyPlan` deliberately admitted an entirely untracked consumed-scratch group after Git target preflight reported it dirty. That exception did not check `options.Commit`, so the group bypassed the transaction in both scratch-only and mixed runs and was deleted afterward even though no Git commit could contain or roll back the deletion.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modified)

**What was done:** Commit mode now refuses the existing untracked consumed-scratch exception before eligibility resolution. The finding names every scratch path, explains why Git cannot include the deletion, and points to the non-commit cleanup command. Mixed plans still commit independently eligible tracked cleanup while preserving the refused scratch group; non-commit mode retains its explicitly labeled non-rollback deletion.

## Qualification

Passed — 2 files verified, 6 requirements traced, P-A-U confirmed. Exact inventory, rooted containment, consumed-manifest revalidation, dirty-index behavior, and the non-commit exception are unchanged. No unrelated paths were included.

## Testing

**Tests run:** focused consumed-scratch tests; full cleanup package; `go vet ./... && go test ./... -count=1` in the CLI module; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing, including queue-board ordinary/strict JavaScript lanes and canonical maintainer verification. The optional external-browser lane was unavailable and skipped; this backend cleanup boundary has no browser acceptance condition.

**Red-green validation:**
- `TestConsumedUntrackedRunCommitRefusesScratchOnlyAndMixedDeletion`: ✗ before implementation (both cases returned success and deleted scratch; scratch-only left HEAD unchanged) → ✓ after (exact scratch inventory refused and preserved; mixed tracked work commits independently)

**New tests added:**
- Scratch-only `--commit` refusal with byte-identical inventory and unchanged HEAD
- Mixed tracked/scratch `--commit` behavior with a truthful tracked commit and preserved refused scratch

*Verified by work action*

## Review

**Overall: 99%** | 2026-09-01T22:50:38Z

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
**Acceptance:** Pass — commit mode cannot report or perform an untracked scratch deletion outside its Git transaction, while non-commit cleanup preserves the established narrow exception.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

Cleanup commit mode now refuses destructive work that its exact Git commit cannot represent or roll back. The do-work CLI prime remains current.
