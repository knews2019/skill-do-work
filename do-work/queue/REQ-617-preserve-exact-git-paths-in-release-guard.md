---
id: REQ-617
title: 'Preserve exact Git paths in release guard'
status: pending
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
- [ ] **[PLAN]:** Read listed prime files and agent rules. Write a brief technical approach before coding.
- [ ] **[APPLY]:** Implement the agreed scope and test-first regression.
- [ ] **[UNIFY]:** Review every changed file and run the appropriate native checks; record what was verified.

## Full Context
See `do-work/user-requests/UR-129/input.md` for the full input and batch mapping.
