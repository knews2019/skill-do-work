---
id: REQ-615
title: 'Distinguish functional manifest edits from release metadata'
status: completed
route: A
heavy_verified_at: 2026-09-13T13:53:07Z
heavy_verified_revision: 5d1300d5d4aa8a8286232f987585e9242b52f8f0
re_review_at: 2026-09-13T13:06:29Z
kb_status: pending
remediation_at: 2026-09-13T13:02:42Z
review_at: 2026-09-13T13:02:03Z
builder_handback_at: 2026-09-13T12:57:49Z
integration_at: 2026-09-13T12:57:49Z
dispatch_at: 2026-09-13T12:52:56Z
estimate:
  p50_active_minutes: 10
  confidence: medium
  basis:
  - Route A
  - 3-file write set
  - 4 acceptance criteria
  calculated_at: 2026-09-13T12:51:00Z
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
claimed_at: 2026-09-13T12:49:50Z
commit: 4cded81079d242a511cbc22a4c091f9252473744
completed_at: 2026-09-13T13:53:07Z
release_at: 2026-09-13T13:53:07Z
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
- [x] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [x] **[APPLY]:** Implement the agreed scope and test-first regression.
- [x] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.

## Triage

**Route: A** — specific release-guard defect with named files and reproducible acceptance cases.

## Plan

Planning not required — focused correction using existing manifest content comparison and genuine test-first regressions.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (modified) — compare functional manifest content at the supplied or pending provenance boundary.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go` (modified) — real Git cases for dependency, configuration and version-only edits.

**What was done:** Manifest dependency and executable configuration changes qualify as shipped implementation; release-only edits remain excluded. Ownership API unchanged. Implemented by e74113f32c56d19bcb0025762f5c35e5086e0d6a, integrated in 4cded81079d242a511cbc22a4c091f9252473744; range b490d5ff5704ba925e70c2e29704d10b3cddaf40..4cded81079d242a511cbc22a4c091f9252473744.

## Decisions

D-01 (DECIDE & STATE): Keep edit classification in finalization and preserve ownership semantics. Normalize only the project version, retaining dependency versions even when equal.

## Implementation Evidence

Builder read both listed primes, required release lessons and the touch-required CLI satellite; none missing. PLAN/APPLY/UNIFY completed. RED before code: TestReleaseGuardDistinguishesFunctionalManifestEdits failed in 4.254s; GREEN passed in 3.304s. Full finalization package passed in 48.361s, go vet and diff check passed. The repository gate below measures each test file separately against the 30s budget.

## Remediation

Strengthen real-Git regression for commented dependency/project headers and correct section recognition. Existing tests did not cover valid trailing section comments.

## Review History

First review Partial (76.25%): valid TOML dependency header comments confused section tracking. Remediation 922e59448b7a6b60b5eab1129f8a110aa34dd65d added four real Git RED cases and corrected section recognition; release guard tests GREEN in 6.915s, focused vet and diff checks pass. First merged gate passed in 137s. Cumulative range retains initial b490d5ff base; metadata-only cf4bd1bc inside it is explicitly orchestration evidence, not implementation.

## Qualification

Canonical cumulative merged-range qualifier satisfied without findings. Two source/test files remain the implementation scope; root bookkeeping inside cumulative range is excluded by contract. Header comments now transition TOML sections correctly.

## Testing

Canonical focused and green-gate records both satisfied at final merged revision. `_dev/tests/maintainer-verify.sh` executed directly via timed runner: PASS, 147s; every test file under30s, slowest CLI file24.14s. Focused release regressions and independent reviewer check passed (5.166s reviewer). RED/GREEN: original six functional cases failed before code; 18 expanded controls pass. Four commented-header cases failed before remediation and pass now. Heavy lanes selected below.

## Review

Overall: 97.5%
Acceptance: Pass
Independent re-review approves all four detailed requirements, verifies original TOML comment finding closed, and finds no open issues. Requirements100%, Code95%, Tests95%, Scope100%; Risk Low. Ownership restatement sweep clean. No follow-ups.

## Lessons Learned

**What worked:** Real Git cases bind manifest bytes to each provenance mode.
**What did not:** Bare TOML header recognition lost section boundaries when headers had trailing comments; the paired dependency/project controls now pin it.
**Worth knowing:** Ownership of a metadata-bearing file does not imply every edit in it is release-only. Dedicated metadata paths remain excluded; manifest changes require content comparison.

## Orientation

The release guard now accepts shipped dependency and executable-configuration changes in package manifests while continuing to refuse version-only or maintainer-only releases.

## Discovered Tasks

None outside the request. The review finding was fixed in scope.

## Heavy Verification Plan

```json
{
  "manifest_path": "_dev/tests/heavy-lanes.json",
  "base_revision": "b490d5ff5704ba925e70c2e29704d10b3cddaf40",
  "target_revision": "4cded81079d242a511cbc22a4c091f9252473744",
  "forced_all": false,
  "uncertain": false,
  "changed_paths": [
    "do-work/working/REQ-615-distinguish-functional-manifest-edits-from-release-metadata.md",
    "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go",
    "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go"
  ],
  "uncovered_paths": [],
  "selected_lanes": [
    {
      "lane_id": "do-work-cli-integrations",
      "command_argv": [
        "env",
        "GIT_CONFIG_NOSYSTEM=1",
        "GIT_CONFIG_GLOBAL=/dev/null",
        "bash",
        "_dev/tests/maintainer-verify.sh",
        "--heavy-lane",
        "do-work-cli-integrations"
      ],
      "reasons": [
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go matched subtree skills/do-work/tools/do-work-cli"
      ]
    },
    {
      "lane_id": "staged-skills",
      "command_argv": [
        "env",
        "GIT_CONFIG_NOSYSTEM=1",
        "GIT_CONFIG_GLOBAL=/dev/null",
        "bash",
        "_dev/tests/maintainer-verify.sh",
        "--heavy-lane",
        "staged-skills"
      ],
      "reasons": [
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go matched subtree skills",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go matched subtree skills"
      ]
    },
    {
      "lane_id": "updater",
      "command_argv": [
        "env",
        "GIT_CONFIG_NOSYSTEM=1",
        "GIT_CONFIG_GLOBAL=/dev/null",
        "bash",
        "_dev/tests/maintainer-verify.sh",
        "--heavy-lane",
        "updater"
      ],
      "reasons": [
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go matched subtree skills/do-work/tools/do-work-cli"
      ]
    },
    {
      "lane_id": "installer",
      "command_argv": [
        "env",
        "GIT_CONFIG_NOSYSTEM=1",
        "GIT_CONFIG_GLOBAL=/dev/null",
        "bash",
        "_dev/tests/maintainer-verify.sh",
        "--heavy-lane",
        "installer"
      ],
      "reasons": [
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go matched subtree skills/do-work/tools/do-work-cli",
        "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go matched subtree skills/do-work/tools/do-work-cli"
      ]
    }
  ]
}
```

## Heavy Verification Result

Stored plan recomputed at its exact base and target with no drift. Target 4cded81079d242a511cbc22a4c091f9252473744; shared execution 5d1300d5d4aa8a8286232f987585e9242b52f8f0. Every selected lane exited0 and none was skipped.

- do-work-cli-integrations: executed, exit0, 66s (fingerprint_mismatch).
- staged-skills: executed, exit0, 32s (fingerprint_mismatch).
- updater: executed, exit0, 64s (fingerprint_mismatch).
- installer: executed, exit0, 29s (fingerprint_mismatch).

All heavy runner processes finished. No additional lane drain is needed for this implementation commit.

## Deferred Lesson Publication

Added the manifest-ownership-vs-edit-content family to _dev/primes/lessons-releases.md and refreshed its exact family set/size in do-work/lessons-index.md. The lesson links the stable UR-129 archive destination chosen by finalization. First occurrence specific to manifest release classification, so no broad prime promotion.

## Timing

Observed 2026-09-13T12:58:21Z to 2026-09-13T13:05:32Z: 7m 11s total, 4m 44s attributed across 2 events, 2m 27s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 4m 44s | 2 |

Slowest command: verification-gate / maintainer-verify.sh (1 argv tokens), 2m 27s, exit 0, maintainer-verify.sh (1 argv tokens).
