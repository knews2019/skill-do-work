---
id: REQ-616
title: 'Require pending shipped changes for primary release'
status: claimed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-09-13T13:06:35Z
created_at: 2026-09-07T15:44:32Z
user_request: UR-129
domain: backend
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
related: ["REQ-615", "REQ-617", "REQ-618", "REQ-619", "REQ-620", "REQ-621"]
batch: validated-finalization-feedback
depends_on: []
write_set: ["skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go"]
required_lessons: [_dev/primes/lessons-releases.md]
claimed_at: 2026-09-13T13:06:29Z
---
# Require pending shipped changes for primary release

## What
Require actual pending shipped implementation changes for primary_commit release provenance instead of treating commit_paths membership as evidence of a change.

## Detailed Requirements
- [ ] Intersect the exact commit allowlist with actual pending changes before evaluating shipped roots.
- [ ] Include additions and deletions, and exclude clean shipped entries.
- [ ] Retain rejection when only maintainer files and release metadata actually change.
- [ ] Correct the existing evidence source without adding a parallel admission layer.

## Finding Provenance
Original severity: P2. Verdict: Accept. Original comments: F2, F4, F10. Duplicate mapping: F2, F4, F10 are one defect.

Source: `do-work/user-requests/UR-129/assets/source-review.txt`, copied byte-for-byte from `/Users/t2/.codex/attachments/87105c57-a948-4858-b928-6e1c45a94733/pasted-text.txt`. External provider not identified. Original installed-path citations map to this repository's canonical `skills/do-work/` tree.

### Original Claim F2 (Verbatim)
> ```
>   - [P2] Check actual pending changes for primary-commit provenance --- [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65
>     For primary_commit, an unchanged shipped file in commit_paths is sufficient to authorize a release containing only maintainer changes and release metadata. This list is an allowlist, and CommitExactPaths explicitly filters out clean entries before committing, so
>     membership does not prove the implementation changes that file. Intersect the allowlist with actual pending changes before checking shipped roots.
> 
> ```

### Original Claim F4 (Verbatim)
> ```
>   - [P2] Check actual changes instead of trusting the commit allowlist — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65
>     For primary_commit provenance, an unchanged shipped file in commit_paths satisfies this guard. CommitExactPaths later drops clean paths, so the resulting commit can contain only maintainer files and release metadata despite passing the new check. Intersect the
>     allowlist with actual pending changes before deciding that the implementation ships something. This follows the prime's affirmative-evidence rule (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L35).
> 
> ```

### Original Claim F10 (Verbatim)
> ```
>   - [P2] Check changed paths rather than the commit allowlist — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65
>     Under primary_commit provenance, an unchanged shipped file in commit_paths authorizes a release even when the only actual changes are maintainer tests and release metadata. CommitExactPaths filters unchanged entries out, so the resulting commit ships no implementation
>     change despite passing this guard. Intersect the allowlist with actual pending changes, including additions and deletions. This is the distinction required by the shipped-change rule (.claude/skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#L5), and a focused
>     probe reproduces the bypass.
> 
> ```

## Validation Evidence
Validated against HEAD `d4ea51d7` through source, surrounding checks, independent reviewers, and Git history. No execution probe was run during triage; no staged or unstaged fix was present.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65 returns manifest.CommitPaths directly; lines 49-51 accept any listed shipped path.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go:158-174 checks normalization and required-target coverage, not pending implementation changes.
- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:45-58 filters clean allowlisted entries before staging.

History and scope correction: Introduced by 6ecff1bc7aad6ef8bc80f0ecb515270b3d20c45f. The existing primary-provenance guard test even accepts a shipped path without creating it; no later preparation check compensates.

## Surface-cost
N/A — direct correction of the existing guard's evidence source.

## Red-Green Proof
**RED prompt/case:** List an unchanged shipped file beside dirty maintainer files and release metadata in commit_paths; the release guard accepts although the final commit ships no implementation change.

**Why RED now:** skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65 returns manifest.CommitPaths directly; lines 49-51 accept any listed shipped path.

**GREEN when:** A focused regression in finalization_release_guard_test.go refuses the clean-shipped-path bypass and accepts actual shipped additions, edits, and deletions.

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
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (modified or added in 04d55042).
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_review_release_test.go` (modified or added in 04d55042).

**What was done:** Actual pending shipped changes are intersected with the exact commit allowlist. Clean listed paths cannot authorize a release.

Implementation predates this run: 04d550423ab7ec89b1270869b09386c43a11ec41, already released as 0.305.39. Current run verifies and reconciles the durable request only. Its original multi-fix commit includes unrelated changes, which are not attributed to this request.

## Decisions

D-01 (DECIDE & STATE): Preserve the existing implementation and reuse its supplied commit provenance; no redundant fix or release. Historical RED/GREEN verifies the unchanged regression against the fix parent, and does not claim this run authored tests before code.

## Implementation Evidence

Required lessons consulted from current lessons index. Release lesson entries retained/read where captured; large partial CLI satellite remains budget-dropped but was read by independent verifier under touch-conditional discipline. Original builder handback unavailable, so no authorship/timing reconstructed. Current P-A-U means plan existing-fix reconciliation, apply no source edits, unify by reviewing actual diff and executing focused/historical evidence.
