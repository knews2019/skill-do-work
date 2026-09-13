---
id: REQ-619
title: '[impact-critical] Verify primary commit content against prepared identity'
status: completed
route: A
review_at: 2026-09-13T13:43:48Z
kb_status: pending
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-09-13T13:41:44Z
created_at: 2026-09-07T15:44:32Z
user_request: UR-129
domain: backend
impact: impact-critical
effort_estimate: effort-substantive
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
related: ["REQ-615", "REQ-616", "REQ-617", "REQ-618", "REQ-620", "REQ-621"]
batch: validated-finalization-feedback
depends_on: ["REQ-618", "REQ-620"]
write_set: ["skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply_test.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go"]
claimed_at: 2026-09-13T13:41:44Z
completed_at: 2026-09-13T13:43:48Z
commit: 04d550423ab7ec89b1270869b09386c43a11ec41
---
# Verify primary commit content against prepared identity

## What
Verify committed implementation content against the complete prepared identity before accepting the primary commit, including hook-staged rewrites and deletions.

## Detailed Requirements
- [ ] Use the existing post-commit callback to compare the committed bytes with the complete prepared commit identity.
- [ ] Detect a hook rewriting or deleting and staging an already-dirty tracked implementation file, even though the changed-path set is unchanged.
- [ ] Include new files through the complete prepared identity fix.
- [ ] On verification failure preserve the already-created SHA and committed-risk journal semantics from the prerequisite fix.
- [ ] Keep ordinary tracked-plus-new-file commits successful.

## Finding Provenance
Original severity: P1. Verdict: Accept. Original comments: F8. Duplicate mapping: F8 is a distinct defect.

Source: `do-work/user-requests/UR-129/assets/source-review.txt`, copied byte-for-byte from `/Users/t2/.codex/attachments/87105c57-a948-4858-b928-6e1c45a94733/pasted-text.txt`. External provider not identified. Original installed-path citations map to this repository's canonical `skills/do-work/` tree.

### Original Claim F8 (Verbatim)
> ```
>   - [P1] Restore primary-commit content verification — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:133-133
>     When a pre-commit hook rewrites or deletes an allowlisted implementation file, path-only verification can still pass. Passing nil removes the check against the prepared bytes, while verifyFinalState checks request/release state rather than implementation content. Re-
>     running the deleted hook-deletion test now reports successful finalization with the implementation file removed. Restore content verification before accepting the primary commit.
> 
> ```

## Validation Evidence
Validated against HEAD `d4ea51d7` through source, surrounding checks, independent reviewers, and Git history. No execution probe was run during triage; no staged or unstaged fix was present.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:133 passes nil as the post-commit verifier.
- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:88-99 compares filenames only and skips a nil verifier.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:583-625 verifies archive identity/status/provenance and release images, not implementation bytes.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:530-559 compares the prepared digest only during recovery detection before making a new commit.

History and scope correction: The nil callback is present from foundation commit 761d8e6ab47dd940db6d686ef641ea919ee9bb31. No recent removal is established. Deleting a newly created file can already cause a path mismatch; the verified gap specifically includes tracked files.

## Surface-cost
Earned — a staged hook rewrite/deletion preserves the path set while violating prepared content. Reuse the existing callback/digest mechanisms rather than adding a new state system; paired negative tests and a normal commit control prove the value.

## Red-Green Proof
**RED prompt/case:** Preparation expects to modify a tracked implementation file; a pre-commit hook rewrites it or runs git rm and stages the result. The filename-only check accepts the changed bytes or deletion.

**Why RED now:** skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:133 passes nil as the post-commit verifier.

**GREEN when:** Focused finalization tests fail on staged hook rewrite and deletion with retained SHA/journal risk, and pass for the prepared tracked-plus-new-file content.

**Validation:** User confirmed — the user explicitly requested capture of the accepted triage findings and their evidence, including the proof cases above. These are targets for the builder to reproduce, not tests already executed.

## Constraints
- Preserve original claims and severity as source data; use the validated scope/history corrections when explaining the fix.
- Limit implementation to this defect and its meaningful regression coverage. Follow the existing journal, release, and output contracts.
- Capture does not execute: this request remains pending for a separate work invocation.

## Dependencies
- REQ-618 (Preserve committed-risk before finalization rollback) is a prerequisite: complete content verification needs both complete digest coverage and safe handling of already-created commits.
- REQ-620 (Include untracked files in prepared commit digest) is a prerequisite: complete content verification needs both complete digest coverage and safe handling of already-created commits.

## Builder Guidance
Use the existing Go test harness for a genuine test-first regression. The accepted remedy defines the outcome; choose the smallest implementation that satisfies it. Do not claim a recently removed mechanism was restored without additional history evidence.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 15905 indexed tokens; owning prime and relevant release/recovery/evidence failure families match, but the index marks the satellite `slugged: partial`, so targeted loading is ineligible and the whole file exceeds the 2000-token budget.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [x] **[APPLY]:** Implement the agreed scope and test-first regression.
- [x] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.

## Triage

**Route: A** — reconcile already-delivered focused fix against its historical implementation and existing regression tests.

## Plan

Planning not required. Verify existing implementation commit 04d550423ab7ec89b1270869b09386c43a11ec41 and close the stale request without duplicate production edits or a duplicate release.

## Implementation Summary

**Files changed by the existing implementation:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modified or added in 04d55042).
- `skills/do-work/tools/do-work-cli/internal/finalization/review_regressions_test.go` (modified or added in 04d55042).

**What was done:** Primary commit content verification rejects hook-staged rewrites and deletions against the complete prepared identity while retaining committed risk.

Implementation predates this run: 04d550423ab7ec89b1270869b09386c43a11ec41, already released as 0.305.39. Current run verifies and reconciles the durable request only. Its original multi-fix commit includes unrelated changes, which are not attributed to this request.

## Decisions

D-01 (DECIDE & STATE): Preserve the existing implementation and reuse its supplied commit provenance; no redundant fix or release. Historical RED/GREEN verifies the unchanged regression against the fix parent, and does not claim this run authored tests before code.

## Implementation Evidence

Required lessons consulted from current lessons index. Release lesson entries retained/read where captured; large partial CLI satellite remains budget-dropped but was read by independent verifier under touch-conditional discipline. Original builder handback unavailable, so no authorship/timing reconstructed. Current P-A-U means plan existing-fix reconciliation, apply no source edits, unify by reviewing actual diff and executing focused/historical evidence.

## Qualification

Canonical qualification satisfied for original implementation range 0c0d9933..04d55042. Only the request-relevant files/hunks are attributed here.

## Reconciliation Heavy Evaluation

No new implementation delta: canonical plan-heavy-verification for afcb6a02cd34f42ed287d85c7d6a828ef9e4f8f8..afcb6a02cd34f42ed287d85c7d6a828ef9e4f8f8 selects no lanes. Historical range above is used for behavioral verification and qualification.

## Testing

Request-bound focused test and green-gate records satisfied. Canonical `_dev/tests/maintainer-verify.sh` was executed directly for this revision; its successful status was supplied to advance. All per-test-file budgets passed.

**Red-green validation (historical differential, not current test authorship):** TestReviewPrimaryCommitRejectsHookContentChanges. Existing tests copied unchanged to04d55042 parent (0c0d9933) fail on behavior (1.95s); current passes (1.61s). Historical staged tracked rewrite/deletion both return success. Current rejects both, preserving SHA and unresolved journal; normal primary creation also passes (0.91s). REQ618 adds the typed-risk/actionable-revert projection.

No new implementation changes were made by reconciliation. The typed heavy plan for this reconciliation range contains no selected lanes. The active batch also verifies all new CLI implementation through its selected shared heavy lanes.

## Review

Overall: 97.5%
Acceptance: Pass
Independent review confirms detailed requirements and historical finding closure. Requirements100%, Code95%, Tests100%, Scope95%; Risk Low. Existing implementation and tests belong to04d550423ab7ec89b1270869b09386c43a11ec41, released0.305.39. No new source is attributed to this run. Original handback unavailable; current reconciliation P-A-U and historical evidence are explicitly distinguished from authorship. No unresolved acceptance findings or follow-up requests.

## Lessons Learned

**Worth knowing:** A queued finding can outlive an independently delivered fix. Compare the existing regression with its actual pre-fix revision before adding code; source history is stronger than stale queue wording. No new subsystem lesson beyond existing finalization/evidence contracts.

## Orientation

Commit verification detects hook rewrites or deletions that preserve the changed filenames. This run verifies and closes the preexisting fix.

## Discovered Tasks

None.

## Prerequisite Review Closure

Independent reviewer verified merged REQ-618 hook tests2.733s and closed its set-aside prose correction at0df80349. Compatibility with external70c16925 passed the combined finalization6.025s and gittransaction1.121s runs. Prior conditional acceptance and missing current P-A-U observations are resolved by actual implementation/review and this reconciliation log. Original04d55042 attribution remains unchanged.

## Timing

Observed 2026-09-13T13:41:54Z to 2026-09-13T13:42:59Z: 1m 05s total, 1m 05s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 1m 05s | 1 |

Slowest command: verification-gate / maintainer-verify.sh (1 argv tokens), 1m 05s, exit 0, maintainer-verify.sh (1 argv tokens).
