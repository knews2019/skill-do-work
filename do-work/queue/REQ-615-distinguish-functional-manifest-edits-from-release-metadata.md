---
id: REQ-615
title: 'Distinguish functional manifest edits from release metadata'
status: pending
created_at: 2026-09-07T15:44:32Z
user_request: UR-129
domain: backend
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
related: ["REQ-616", "REQ-617", "REQ-618", "REQ-619", "REQ-620", "REQ-621"]
batch: validated-finalization-feedback
depends_on: []
write_set: ["skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go", "skills/do-work/tools/do-work-cli/internal/releaseownership/release_ownership.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go"]
required_lessons: [_dev/primes/lessons-releases.md]
---
# Distinguish functional manifest edits from release metadata

## What
Allow a release whose shipped implementation changes functional dependencies, entry points, or build configuration in package.json, Cargo.toml, or pyproject.toml; exclude only edits limited to release metadata.

## Detailed Requirements
- [ ] Compare changed manifest content before classifying the implementation as release-metadata-only.
- [ ] Retain rejection of version-only/release-only edits and maintainer-only releases.
- [ ] Cover supplied_commit and primary_commit provenance; preserve the consumer exemption when there are no declared maintainer module roots.
- [ ] Do not change the ownership predicate into a content predicate for its unrelated callers.

## Finding Provenance
Original severity: P2. Verdict: Accept. Original comments: F1, F5. Duplicate mapping: F1, F5 are one defect.

Source: `do-work/user-requests/UR-129/assets/source-review.txt`, copied byte-for-byte from `/Users/t2/.codex/attachments/87105c57-a948-4858-b928-6e1c45a94733/pasted-text.txt`. External provider not identified. Original installed-path citations map to this repository's canonical `skills/do-work/` tree.

### Original Claim F1 (Verbatim)
> ```
>   - [P2] Distinguish metadata-only edits from functional manifest changes --- [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:46-48
>     When a supplied implementation changes only a shipped package.json, Cargo.toml, or pyproject.toml, this check refuses the release even if the change updates runtime dependencies, entry points, or build configuration. IsReleaseMetadataPath identifies files that can
>     contain release metadata, not changes limited to release metadata. Inspect the changed content before excluding these manifests so functional dependency/configuration fixes remain releasable.
> 
> ```

### Original Claim F5 (Verbatim)
> ```
>   - [P2] Preserve dependency changes when excluding release metadata — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:46-47
>     When a shipped module changes only dependencies in package.json, Cargo.toml, or pyproject.toml, this check rejects its release because IsReleaseMetadataPath excludes the entire file. These files contain executable configuration and dependency requirements, not just
>     version fields. Distinguish version-only edits from substantive manifest changes rather than rejecting by basename, consistent with the prime's condition-based classification rule (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L35).
> 
> ```

## Validation Evidence
Validated against HEAD `d4ea51d7` through source, surrounding checks, independent reviewers, and Git history. No execution probe was run during triage; no staged or unstaged fix was present.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:46-48 excludes each entire metadata-bearing path before root matching.
- skills/do-work/tools/do-work-cli/internal/releaseownership/release_ownership.go:42-46 defines IsReleaseMetadataPath as whether the file carries release metadata at all.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go:129 applies this guard before release planning.

History and scope correction: Introduced by 6ecff1bc7aad6ef8bc80f0ecb515270b3d20c45f. Ordinary consumer repositories bypass the guard; no such manifests currently exist under this checkout's skills tree. This is a supported-input defect verified from source, not an observed failed release here.

## Surface-cost
N/A — direct correction of the existing overbroad guard. No second validation layer is needed.

## Red-Green Proof
**RED prompt/case:** In a declared maintainer module, change only runtime dependencies or executable configuration in one of the three named manifests and request release; the whole-file exclusion refuses it.

**Why RED now:** skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:46-48 excludes each entire metadata-bearing path before root matching.

**GREEN when:** A focused regression in finalization_release_guard_test.go accepts substantive edits and rejects version-only edits for both provenance modes, including package.json, Cargo.toml, and pyproject.toml.

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
