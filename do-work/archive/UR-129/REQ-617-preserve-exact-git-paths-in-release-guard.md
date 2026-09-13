---
id: REQ-617
title: 'Preserve exact Git paths in release guard'
status: completed
route: A
kb_status: pending
heavy_verified_at: 2026-09-13T13:50:31Z
heavy_verified_revision: 5d1300d5d4aa8a8286232f987585e9242b52f8f0
review_at: 2026-09-13T13:22:16Z
builder_handback_at: 2026-09-13T13:15:24Z
integration_at: 2026-09-13T13:15:24Z
dispatch_at: 2026-09-13T13:11:36Z
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-09-13T13:10:15Z
created_at: 2026-09-07T15:44:32Z
user_request: UR-129
domain: backend
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md", "_dev/primes/prime-shell-commands.md"]
tdd: true
maintenance: false
related: ["REQ-615", "REQ-616", "REQ-618", "REQ-619", "REQ-620", "REQ-621"]
batch: validated-finalization-feedback
depends_on: []
write_set: ["skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go", "skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go"]
required_lessons: [_dev/primes/lessons-releases.md]
claimed_at: 2026-09-13T13:09:50Z
commit: 24b52cacdad9eb6bf0c863e1c65c6739ea54e8b0
completed_at: 2026-09-13T13:51:43Z
release_at: 2026-09-13T13:51:43Z
---
# Preserve exact Git paths in release guard

## What
Read supplied implementation commit filenames as exact NUL-delimited paths so Git display quoting does not reject valid shipped changes.

## Detailed Requirements
- [ ] Request diff-tree -z output and split the raw result on NUL.
- [ ] Avoid TrimSpace or other normalization that destroys legitimate filename bytes.
- [ ] Preserve first-parent comparison for merge commits and root-commit handling.
- [ ] Cover non-ASCII, tab, newline, and quote-bearing filenames under declared shipped roots.

## Finding Provenance
Original severity: P2. Verdict: Accept. Original comments: F3, F6. Duplicate mapping: F3, F6 are one defect.

Source: `do-work/user-requests/UR-129/assets/source-review.txt`, copied byte-for-byte from `/Users/t2/.codex/attachments/87105c57-a948-4858-b928-6e1c45a94733/pasted-text.txt`. External provider not identified. Original installed-path citations map to this repository's canonical `skills/do-work/` tree.

### Original Claim F3 (Verbatim)
> ```
>   - [P2] Read changed filenames using NUL-delimited Git output --- [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:78-82
>     With Git's default quoting, a supplied commit touching only skills/do-work/docs/café.md returns a quoted, escaped pathname. The subsequent prefix comparison therefore misses the shipped root and incorrectly refuses the release. Filenames containing tabs, newlines, or
>     quotes have the same problem. Request diff-tree -z output and split on NUL rather than interpreting display-formatted lines as exact paths.
> 
> ```

### Original Claim F6 (Verbatim)
> ```
>   - [P2] Read changed filenames using NUL-delimited Git output — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:78-81
>     With Git's default core.quotePath=true, a supplied commit changing only a shipped filename such as skills/do-work/docs/café.md is incorrectly rejected. diff-tree --name-only returns a quoted, escaped pathname, which this parser preserves literally, so it no longer
>     matches the shipped-root prefix. Request -z output and split on NUL to preserve non-ASCII filenames and filenames containing control characters.
> 
> ```

## Validation Evidence
Validated against HEAD `d4ea51d7` through source, surrounding checks, independent reviewers, and Git history. No execution probe was run during triage; no staged or unstaged fix was present.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:68 omits -z; lines 79-81 split display-formatted output on newline after TrimSpace.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:45-50 does not unquote the path before comparing the shipped-root prefix.

History and scope correction: Introduced by 6ecff1bc7aad6ef8bc80f0ecb515270b3d20c45f. Current HEAD has no special filenames; validation established the code path, not an executed end-to-end reproduction.

## Surface-cost
N/A — direct parser repair.

## Red-Green Proof
**RED prompt/case:** A supplied implementation commit changes only skills/do-work/docs/café.md with Git default quoting; the returned quoted path fails shipped-root matching.

**Why RED now:** skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:68 omits -z; lines 79-81 split display-formatted output on newline after TrimSpace.

**GREEN when:** Focused finalization_release_guard_test.go cases accept shipped non-ASCII/control-character/quote filenames while preserving existing ordinary, merge, root-commit, and maintainer-only behavior.

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
- `_dev/primes/lessons-shell-commands.md` — 8480 indexed tokens; argv/quoting or migration parity matches, but `slugged: partial` prevents narrowing and the whole satellite exceeds budget.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [x] **[APPLY]:** Implement the agreed scope and test-first regression.
- [x] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.

## Triage

**Route: A** — direct exact-path parser correction with named code and behavioral cases.

## Plan

Planning not required. Consume NUL-delimited raw Git paths and retain root/merge provenance behavior; test exact bytes and guard outcomes. Required lesson index consulted; captured release lessons retained and partial large satellites remain dropped from the budget. Builder still loads touch-required satellites.

## Implementation Summary

**Files changed:**

- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (modified): request git diff-tree -z, split raw stdout on NUL, omit only empty records, and remove implementation-path normalization before shipped-root classification. Existing root/first-parent selection and functional-manifest interpretation remain unchanged.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go` (modified): add one nine-case regression matrix asserting exact supplied-commit path bytes and release guard verdicts for ordinary, non-ASCII, tab, newline, quote, surrounding leaf spaces, leading root space, initial commit, and merge commit fixtures.

The commit contains exactly these two files, with 66 insertions and 5 deletions. No queue or release files were edited in the worktree. `git status --porcelain=v1` was empty after committing.

## Decisions

- D1: Git already emits canonical repository-relative path records under `-z`; retain the bytes and remove the redundant normalization at their classification boundary. No unquoting helper or new abstraction is needed.
- D2: Initial-commit fixture commits only the special shipped path, then adds maintainer declarations later. This prevents an ordinary `suite/modules.tsv` path from making release acceptance pass accidentally.
- D3: Merge fixture includes unrelated main-branch work and asserts the sole returned path is the special feature-branch file. Retain the existing maintainer-only merge refusal to prove first-parent semantics in both directions.

### Builder Discovered Tasks

None. The leading-space false acceptance is another manifestation of this request's exact-path parsing defect and remains inside its stated no-normalization acceptance criterion.

## Implementation Evidence

Range5aed3c4c4bae6e557c59327b79416a13de313776..24b52cacdad9eb6bf0c863e1c65c6739ea54e8b0. Genuine RED test (8of9 failures,1.731s) preceded code; GREEN9cases (1.757s), full release-guard matrix (5.873s), vet and diff checks passed. Required release/CLI/shell lessons and all crew rules read; no missing files. All P-A-U phases completed and exact two-file implementation committed.

## Qualification

Canonical request-bound merged-range qualifier satisfied with no findings after correcting summary markup. Actual two-file diff matches the claimed implementation and completed P-A-U evidence. Both false refusal and whitespace false acceptance are covered.

## Testing

Genuine builder RED: eight of nine exact-path cases failed before production edits (1.731s). Identical GREEN passed (1.757s). Adjacent release matrix passed (5.873s), including functional manifests, first-parent rejection, primary and consumer controls. Independent merged review reran exact paths/merge/root (2.484s) and actual consumer/primary controls (0.533s), all exit 0.

Canonical timed _dev/tests/maintainer-verify.sh exited 0 in 128s on merged revision 24b52cacdad9eb6bf0c863e1c65c6739ea54e8b0: ShellCheck, gofmt, contracts, vet and native suites passed; slowest CLI file 20.33s, within 30s budget. Request-bound advance recorded both focused and green gates satisfied; focused release matrix passed in 8.372s. Heavy lanes await the shared drain below.

## Review

Independent review: **100% — Acceptance: Pass**. No important, minor or nit findings; no follow-ups. Read-only reviewer inspected exact diff and genuine RED/GREEN evidence, tested real Git path bytes plus actual consuming guard, and swept related contracts. All reviewer commands completed.

## Lessons Learned

The existing exact-basename-authority lesson applies: whitespace is path identity before presentation. Exact-byte assertions catch trailing-space loss that a prefix-based acceptance assertion misses; the leading-space case also proves prevention of false acceptance. No new lesson satellite needed.

## Orientation

Release checks now recognize supplied commits with Unicode, tabs, newlines, quotes and spaces in filenames while preserving exact shipped-root boundaries and first-parent semantics.

## Discovered Tasks

None outside this request.

## Heavy Verification Plan

```json
{
  "manifest_path": "_dev/tests/heavy-lanes.json",
  "base_revision": "5aed3c4c4bae6e557c59327b79416a13de313776",
  "target_revision": "24b52cacdad9eb6bf0c863e1c65c6739ea54e8b0",
  "forced_all": false,
  "uncertain": false,
  "changed_paths": [
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

Stored plan recomputed at its exact base and target with no drift. Target 24b52cacdad9eb6bf0c863e1c65c6739ea54e8b0; shared execution 5d1300d5d4aa8a8286232f987585e9242b52f8f0. Every selected lane exited0 and none was skipped.

- do-work-cli-integrations: executed, exit0, 66s (fingerprint_mismatch).
- staged-skills: executed, exit0, 32s (fingerprint_mismatch).
- updater: executed, exit0, 64s (fingerprint_mismatch).
- installer: executed, exit0, 29s (fingerprint_mismatch).

All heavy runner processes finished. No additional lane drain is needed for this implementation commit.

## Timing

Observed 2026-09-13T13:16:12Z to 2026-09-13T13:18:20Z: 2m 08s total, 2m 08s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 2m 08s | 1 |

Slowest command: verification-gate / maintainer-verify.sh (1 argv tokens), 2m 08s, exit 0, maintainer-verify.sh (1 argv tokens).
