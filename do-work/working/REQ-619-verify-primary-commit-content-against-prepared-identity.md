---
id: REQ-619
title: '[impact-critical] Verify primary commit content against prepared identity'
status: claimed
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
- [ ] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [ ] **[APPLY]:** Implement the agreed scope and test-first regression.
- [ ] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.
