---
id: UR-129
title: 'Capture seven validated finalization and output defects'
created_at: 2026-09-07T15:44:32Z
requests: ["REQ-615", "REQ-616", "REQ-617", "REQ-618", "REQ-619", "REQ-620", "REQ-621"]
word_count: 22
---

## Summary
Capture seven distinct accepted defects from eleven original review comments, retaining all claims, source citations, severities, duplicate mappings, evidence, and surface-cost results.

## Extracted Requests
| Request | Original comments | Severity | Work |
|---|---|---|---|
| REQ-615 | F1, F5 | P2 | Distinguish functional manifest edits from release metadata |
| REQ-616 | F2, F4, F10 | P2 | Require pending shipped changes for primary release |
| REQ-617 | F3, F6 | P2 | Preserve exact Git paths in release guard |
| REQ-618 | F7 | P1 | Preserve committed-risk before finalization rollback |
| REQ-619 | F8 | P1 | Verify primary commit content against prepared identity |
| REQ-620 | F9 | P1 | Include untracked files in prepared commit digest |
| REQ-621 | F11 | P2 | Emit diagnostics for exact-text command results |

## Batch Constraints
- No candidate exists in the queue across any UR: the queue directory was absent at capture. Working is empty. The archived filename scan found no existing home for these exact defects. Each new request is explicitly authorized by this capture.
- Keep the three P1 commit-integrity/recovery fixes coordinated. REQ-619 (Verify primary commit content against prepared identity) depends on REQ-618 (Preserve committed-risk before finalization rollback) and REQ-620 (Include untracked files in prepared commit digest). Other defects have no prerequisite.
- Preserve severity in provenance without inventing a scheduling priority from severity.
- All seven are behavior changes with runnable Go regression targets; TDD applies. The user adopted the prior validated findings and proof targets by requesting this capture.
- Preserve history corrections: the available history does not establish recent removal for the P1 gaps or a recent Go stderr-writer removal.
- No implementation is authorized by capture alone.

## Source Assets
- `assets/source-review.txt` preserves the original attachment bytes.
- `assets/validated-findings.md` preserves the accepted evidence, surface-cost judgments, historical corrections, and proof targets.

## Original Review Context (Verbatim)
> ```
> Full review comments:
> 
>   - [P2] Distinguish metadata-only edits from functional manifest changes --- [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:46-48
>     When a supplied implementation changes only a shipped package.json, Cargo.toml, or pyproject.toml, this check refuses the release even if the change updates runtime dependencies, entry points, or build configuration. IsReleaseMetadataPath identifies files that can
>     contain release metadata, not changes limited to release metadata. Inspect the changed content before excluding these manifests so functional dependency/configuration fixes remain releasable.
> 
>   - [P2] Check actual pending changes for primary-commit provenance --- [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65
>     For primary_commit, an unchanged shipped file in commit_paths is sufficient to authorize a release containing only maintainer changes and release metadata. This list is an allowlist, and CommitExactPaths explicitly filters out clean entries before committing, so
>     membership does not prove the implementation changes that file. Intersect the allowlist with actual pending changes before checking shipped roots.
> 
>   - [P2] Read changed filenames using NUL-delimited Git output --- [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:78-82
>     With Git's default quoting, a supplied commit touching only skills/do-work/docs/café.md returns a quoted, escaped pathname. The subsequent prefix comparison therefore misses the shipped root and incorrectly refuses the release. Filenames containing tabs, newlines, or
>     quotes have the same problem. Request diff-tree -z output and split on NUL rather than interpreting display-formatted lines as exact paths.
> 
>   - [P2] Check actual changes instead of trusting the commit allowlist — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65
>     For primary_commit provenance, an unchanged shipped file in commit_paths satisfies this guard. CommitExactPaths later drops clean paths, so the resulting commit can contain only maintainer files and release metadata despite passing the new check. Intersect the
>     allowlist with actual pending changes before deciding that the implementation ships something. This follows the prime's affirmative-evidence rule (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L35).
> 
>   - [P2] Preserve dependency changes when excluding release metadata — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:46-47
>     When a shipped module changes only dependencies in package.json, Cargo.toml, or pyproject.toml, this check rejects its release because IsReleaseMetadataPath excludes the entire file. These files contain executable configuration and dependency requirements, not just
>     version fields. Distinguish version-only edits from substantive manifest changes rather than rejecting by basename, consistent with the prime's condition-based classification rule (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L35).
> 
>   - [P2] Read changed filenames using NUL-delimited Git output — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:78-81
>     With Git's default core.quotePath=true, a supplied commit changing only a shipped filename such as skills/do-work/docs/café.md is incorrectly rejected. diff-tree --name-only returns a quoted, escaped pathname, which this parser preserves literally, so it no longer
>     matches the shipped-root prefix. Request -z output and split on NUL to preserve non-ASCII filenames and filenames containing control characters.
> 
>    - [P1] Preserve committed-risk evidence before handling commit failure — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:134-138
>     If a hook stages an extra path, CommitExactPaths returns both a failure and the SHA of an already-created commit. This branch now returns before recording that SHA, and the deferred rollback restores lifecycle/release files and resets the journal to prepared despite
>     HEAD having advanced. A regression probe reproduced a restored working request alongside committed archival. Record the SHA before handling failure and prohibit pre-primary rollback once a commit exists, preserving the journal contract (.claude/skills/do-work/tools/
>     do-work-cli/lessons-do-work-cli.md#L158).
> 
>   - [P1] Restore primary-commit content verification — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:133-133
>     When a pre-commit hook rewrites or deletes an allowlisted implementation file, path-only verification can still pass. Passing nil removes the check against the prepared bytes, while verifyFinalState checks request/release state rather than implementation content. Re-
>     running the deleted hook-deletion test now reports successful finalization with the implementation file removed. Restore content verification before accepting the primary commit.
> 
>   - [P1] Include untracked files in the prepared commit digest — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go:1583-1584
>     git diff HEAD omits untracked files, including a newly created archive destination or implementation file. If execution stops after committing but before persisting primary_committed, matchingHeadCommit compares the full committed diff against this incomplete digest
>     and cannot recognize the transaction. A tracked-edit-plus-new-file probe reproduces the mismatch. Restore temporary-index staging so preparation and recovery compare the same bytes, as required by the journal contract (.claude/skills/do-work/tools/do-work-cli/lessons-
>     do-work-cli.md#L158).
> 
>   - [P2] Check changed paths rather than the commit allowlist — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65
>     Under primary_commit provenance, an unchanged shipped file in commit_paths authorizes a release even when the only actual changes are maintainer tests and release metadata. CommitExactPaths filters unchanged entries out, so the resulting commit ships no implementation
>     change despite passing this guard. Intersect the allowlist with actual pending changes, including additions and deletions. This is the distinction required by the shipped-change rule (.claude/skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#L5), and a focused
>     probe reproduces the bypass.
> 
>   - [P2] Retain stderr diagnostics for exact-text results — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go:94-95
>     For text-mode generate-report-image-batch, ExactTextOutput contains only the generated directory, or an empty string when every image fails. The command intentionally returns success with per-image warning findings. Removing stderr rendering therefore makes partial
>     failures invisible and all-failed runs exit zero without any explanation. Keep compatibility stdout unchanged while emitting warning/error findings to stderr, preserving the output-evidence contract (.claude/skills/do-work/tools/do-work-cli/prime-do-work-cli.md#L36).
> ```

## Full Verbatim Input
> ````
> ```sql
> do-work capture-request: Capture the seven distinct accepted defects above, preserving every original claim, severity/source, evidence, duplicate mapping, and surface-cost result.
> ```
> ````
