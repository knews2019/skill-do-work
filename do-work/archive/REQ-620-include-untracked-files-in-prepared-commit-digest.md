---
id: REQ-620
title: '[impact-critical] Include untracked files in prepared commit digest'
status: completed
route: A
review_at: 2026-09-13T13:41:19Z
kb_status: pending
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-09-13T13:39:22Z
created_at: 2026-09-07T15:44:32Z
user_request: UR-129
domain: backend
impact: impact-critical
effort_estimate: effort-substantive
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
related: ["REQ-615", "REQ-616", "REQ-617", "REQ-618", "REQ-619", "REQ-621"]
batch: validated-finalization-feedback
depends_on: []
write_set: ["skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go"]
claimed_at: 2026-09-13T13:39:22Z
completed_at: 2026-09-13T13:41:19Z
commit: 04d550423ab7ec89b1270869b09386c43a11ec41
---
# Include untracked files in prepared commit digest

## What
Bind prepared and recovered commit identities to the same complete content, including untracked archive destinations and implementation additions.

## Detailed Requirements
- [ ] Compute the intended commit image using all exact allowlisted changes, including untracked additions and deletions.
- [ ] Use an isolated temporary index initialized from the recorded HEAD and hash its cached binary diff, or an equally exact approach that leaves the real index intact.
- [ ] Recognize the already-created commit after interruption between commit creation and primary_committed persistence.
- [ ] Preserve lifecycle archival recovery and prevent duplicate commits or erroneous rollback.

## Finding Provenance
Original severity: P1. Verdict: Accept. Original comments: F9. Duplicate mapping: F9 is a distinct defect.

Source: `do-work/user-requests/UR-129/assets/source-review.txt`, copied byte-for-byte from `/Users/t2/.codex/attachments/87105c57-a948-4858-b928-6e1c45a94733/pasted-text.txt`. External provider not identified. Original installed-path citations map to this repository's canonical `skills/do-work/` tree.

### Original Claim F9 (Verbatim)
> ```
>   - [P1] Include untracked files in the prepared commit digest — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go:1583-1584
>     git diff HEAD omits untracked files, including a newly created archive destination or implementation file. If execution stops after committing but before persisting primary_committed, matchingHeadCommit compares the full committed diff against this incomplete digest
>     and cannot recognize the transaction. A tracked-edit-plus-new-file probe reproduces the mismatch. Restore temporary-index staging so preparation and recovery compare the same bytes, as required by the journal contract (.claude/skills/do-work/tools/do-work-cli/lessons-
>     do-work-cli.md#L158).
> 
> ```

## Validation Evidence
Validated against HEAD `d4ea51d7` through source, surrounding checks, independent reviewers, and Git history. No execution probe was run during triage; no staged or unstaged fix was present.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go:1581-1591 hashes git diff HEAD, omitting untracked files.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:114-120 persists the digest before committing.
- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:37-42 requires an empty index; line 58 stages additions only later.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:540-543 hashes the committed diff including additions when trying to recognize recovery.

History and scope correction: This implementation dates to e191b266b9d501ece52d1c9888f7711aaf72f281. Available local history does not support the original claim to "restore" an earlier temporary-index implementation.

## Surface-cost
N/A — direct repair of the existing commit identity computation.

## Red-Green Proof
**RED prompt/case:** Prepare a tracked edit plus a new implementation file or archive destination, commit, then interrupt before primary_committed persists. Recovery cannot match the committed diff to the incomplete prepared digest.

**Why RED now:** skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go:1581-1591 hashes git diff HEAD, omitting untracked files.

**GREEN when:** A focused finalization_recovery_test.go regression recognizes the tracked-plus-new-file commit after the interruption, preserves the real index, and recovers without duplicate commit or pre-primary rollback.

**Validation:** User confirmed — the user explicitly requested capture of the accepted triage findings and their evidence, including the proof cases above. These are targets for the builder to reproduce, not tests already executed.

## Constraints
- Preserve original claims and severity as source data; use the validated scope/history corrections when explaining the fix.
- Limit implementation to this defect and its meaningful regression coverage. Follow the existing journal, release, and output contracts.
- Capture does not execute: this request remains pending for a separate work invocation.

## Dependencies
No prerequisite request. Related requests retain separate acceptance criteria. Shared files alone do not impose sequencing.

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
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modified or added in 04d55042).
- `skills/do-work/tools/do-work-cli/internal/finalization/review_regressions_test.go` (modified or added in 04d55042).

**What was done:** Prepared commit identity includes untracked additions and deletions through a private index. Interrupted primary commits are recognized without duplicates or lifecycle rollback.

Implementation predates this run: 04d550423ab7ec89b1270869b09386c43a11ec41, already released as 0.305.39. Current run verifies and reconciles the durable request only. Its original multi-fix commit includes unrelated changes, which are not attributed to this request.

## Decisions

D-01 (DECIDE & STATE): Preserve the existing implementation and reuse its supplied commit provenance; no redundant fix or release. Historical RED/GREEN verifies the unchanged regression against the fix parent, and does not claim this run authored tests before code.

## Implementation Evidence

Required lessons consulted from current lessons index. Release lesson entries retained/read where captured; large partial CLI satellite remains budget-dropped but was read by independent verifier under touch-conditional discipline. Original builder handback unavailable, so no authorship/timing reconstructed. Current P-A-U means plan existing-fix reconciliation, apply no source edits, unify by reviewing actual diff and executing focused/historical evidence.

## Qualification

Canonical qualification satisfied for original implementation range 0c0d9933..04d55042. Only the request-relevant files/hunks are attributed here.

## Reconciliation Heavy Evaluation

No new implementation delta: canonical plan-heavy-verification for 3bb066e26fc8e0b713ef97f8b6ad5f7a13a083f1..3bb066e26fc8e0b713ef97f8b6ad5f7a13a083f1 selects no lanes. Historical range above is used for behavioral verification and qualification.

## Testing

Request-bound focused test and green-gate records satisfied. Canonical `_dev/tests/maintainer-verify.sh` was executed directly for this revision; its successful status was supplied to advance. All per-test-file budgets passed.

**Red-green validation (historical differential, not current test authorship):** TestReviewPreparedIdentityIncludesAdditionsAndDeletionsWithoutChangingIndex; TestReviewRecoveryRecognizesTrackedAndNewFilesBeforePrimaryPhasePersisted. Existing tests copied unchanged to04d55042 parent (0c0d9933) fail on behavior (0.12s /0.40s); current passes (0.21s /0.64s). Historical preparation omits new bytes and interruption recovery fails to reuse its SHA. Current preserves real index, matches complete image and recovers once.

No new implementation changes were made by reconciliation. The typed heavy plan for this reconciliation range contains no selected lanes. The active batch also verifies all new CLI implementation through its selected shared heavy lanes.

## Review

Overall: 97.5%
Acceptance: Pass
Independent review confirms detailed requirements and historical finding closure. Requirements100%, Code95%, Tests95%, Scope100%; Risk Low. Existing implementation and tests belong to04d550423ab7ec89b1270869b09386c43a11ec41, released0.305.39. No new source is attributed to this run. Original handback unavailable; current reconciliation P-A-U and historical evidence are explicitly distinguished from authorship. No unresolved acceptance findings or follow-up requests.

## Lessons Learned

**Worth knowing:** A queued finding can outlive an independently delivered fix. Compare the existing regression with its actual pre-fix revision before adding code; source history is stronger than stale queue wording. No new subsystem lesson beyond existing finalization/evidence contracts.

## Orientation

Finalization recovery recognizes committed additions and deletions without recreating commits or changing the real index. This run verifies and closes the preexisting fix.

## Discovered Tasks

None.

## Timing

Observed 2026-09-13T13:39:51Z to 2026-09-13T13:40:56Z: 1m 05s total, 1m 05s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 1m 05s | 1 |

Slowest command: verification-gate / maintainer-verify.sh (1 argv tokens), 1m 05s, exit 0, maintainer-verify.sh (1 argv tokens).
