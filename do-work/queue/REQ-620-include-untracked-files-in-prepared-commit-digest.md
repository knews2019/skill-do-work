---
id: REQ-620
title: '[impact-critical] Include untracked files in prepared commit digest'
status: pending
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
- [ ] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [ ] **[APPLY]:** Implement the agreed scope and test-first regression.
- [ ] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.
