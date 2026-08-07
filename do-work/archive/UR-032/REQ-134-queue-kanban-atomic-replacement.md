---
id: REQ-134
title: "Addendum: make queue-kanban atomic replacement cross-platform and symlink-safe"
status: completed
completed_at: 2026-08-07T22:38:55Z
commit: 7e0536a
claimed_at: 2026-08-07T22:32:02Z
route: B
created_at: 2026-08-07T18:58:52Z
user_request: UR-032
addendum_to: REQ-072
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
kb_status: pending
kb_entry:
---

# Addendum: Make Queue-Kanban Atomic Replacement Cross-Platform and Symlink-Safe

## What

Correct the shared queue-kanban atomic replacement path so its crash-safety contract is valid on every supported platform and `next-version` cannot silently replace a symlink with a regular file.

## Context

Addendum to completed REQ-072. A code review of `tools/queue-kanban` found two Important issues in the crash-safe `next-version` follow-up:

- `writeFileAtomically` ends with `os.Rename`, but Go explicitly does not guarantee that operation is atomic on non-Unix platforms. Queue-kanban contains an explicit Windows path, so the current helper overstates its cross-platform guarantee.
- The former `os.WriteFile` behavior followed a version-file symlink. Renaming the temporary file over that path replaces the symlink itself; the read-back succeeds against the new regular file and does not reveal that the intended target was left unchanged.

## Prior Implementation

REQ-072 introduced `next-version`, including the single-version-line rewrite and read-back verification, in commit `5db22ea`. A later direct fix in commit `5e5f338` replaced the truncating `os.WriteFile` call in `release.go` with the shared `writeFileAtomically` helper from `testing.go` and added a regression test covering normal bumps, distinct file identity, permission preservation, and temporary-file cleanup.

Key files are `tools/queue-kanban/release.go`, `release_test.go`, `testing.go`, and `testing_test.go`. The existing testing API rejects symlinked REQ targets before calling the helper; the release path currently has no equivalent symlink policy.

## Requirements

- Provide a genuinely atomic replace operation on every platform queue-kanban supports, including Windows if Windows support remains.
- Preserve the existing `next-version` command semantics: update the intended version file without silently replacing a symlink entry. If safe target resolution is impossible, fail clearly before writing rather than reporting success against the wrong file.
- Preserve existing file permission behavior, the single-line-only rewrite, and the current read-back verification.
- Keep the testing-field write path crash-safe and compatible with its existing regular-file and containment guards.
- Add regression coverage for a symlinked version-file fixture and for the platform-specific replacement abstraction where practical.
- Keep normal patch, minor, and major bump behavior unchanged.
- Run `go test ./...` and `go vet ./...` from `tools/queue-kanban`.

## Constraints

- Scope only `tools/queue-kanban`; do not modify `contact_processor`, unrelated skill documentation, generated assets, or command semantics.
- Do not weaken existing path-containment, loopback, or testing-write protections.
- Do not claim a cross-platform atomicity guarantee that the selected primitive does not provide.

## Red-Green Proof

**RED prompt/case:** Point `allocateNextVersion` at a symlink to a version file and run a patch bump. Today the helper replaces the symlink entry with a regular file; separately, the shared helper relies on an operation whose atomicity is not guaranteed on non-Unix platforms.

**Why RED now:** The current symlink case reports success while changing filesystem identity and leaving the original target unchanged, and the Windows crash-safety claim exceeds Go's documented `os.Rename` guarantee.

**GREEN when:** The symlink remains a symlink and its intended target contains the bumped version (or the command fails before any write under an explicitly safe policy), the replacement implementation has a valid atomic guarantee for each supported platform, all existing mode/read-back tests remain green, and `go test ./...` plus `go vet ./...` pass.

**Validation:** User confirmed by requesting capture of the code-review findings.

## Assets

None.

---
*Source: do-work capture-request: Make queue-kanban atomic replacement cross-platform and symlink-safe  Capture the review findings*

## Triage

**Route: B** — the defect is isolated to one shared write helper and its release caller, but the replacement primitive and build coverage differ by operating system.

## Plan

1. Add a symlinked-version RED test and pin the platform replacement seam with per-OS source files.
2. Reject non-regular targets before creating a temporary file; use atomic rename on Unix and Windows `ReplaceFileW`, with a fail-closed implementation for unsupported platforms.
3. Run native tests/vet/build plus a Windows cross-compile, review every queue-kanban diff, and document the root cause.

## Scope

**Files I will touch:** queue-kanban atomic-write implementation/build-tag files, `release_test.go`, this REQ, and release metadata. **Files I will not touch:** the unrelated REQ-147 allocation changes already present in the same module or any non-queue-kanban application.

## Pre-Flight

The native Go suite passes. The current shared helper calls `os.Stat` (which follows a symlink) and then `os.Rename`; the symlink fixture therefore replaces the link entry, while Go's own contract says the rename is not atomic on non-Unix systems. Microsoft documents `ReplaceFileW` as the single-file replacement alternative for document-like atomicity.

## Root Cause

The shared helper treated Go's portable API surface as a portable atomicity guarantee and used `os.Stat`, which erased the distinction between a regular path and a symlink before replacement. As a result, Unix happened to provide the desired rename semantics, Windows did not promise them, and `next-version` could report success after replacing the link rather than its intended target.

## Implementation Summary

- `tools/queue-kanban/atomic_write.go` — validates a stable existing regular target, prepares and syncs the same-directory temporary file, preserves mode, and delegates the final replacement.
- `tools/queue-kanban/atomic_replace_unix.go`, `atomic_replace_windows.go`, `atomic_replace_unsupported.go` — use the substantiated Unix rename and Windows `ReplaceFileW` primitives, and fail closed elsewhere.
- `tools/queue-kanban/testing.go` — removes the generic `os.Rename` implementation while retaining both existing callers.
- `tools/queue-kanban/release_test.go` — proves a symlinked version path fails without changing the link, target, or temporary-file state.
- `skills/do-work-board/tools/queue-kanban/` mirrored files — keep the staged board package byte-equivalent for the modular cutover.
- Release metadata and this request record — publish the patch-level behavior and its evidence.

## Qualification

Passed. Every implementation line maps to target validation, OS-specific replacement, the symlink regression, or the required staged distribution mirror. Existing Testing-view containment remains upstream of the stricter shared helper, normal version bumps retain mode/read-back behavior, and no other `os.Rename` replacement path exists in queue-kanban.

## Testing

**RED:** `TestAllocateNextVersionRejectsSymlinkWithoutWriting` succeeded on the old helper and failed because the symlink entry became a regular file. **GREEN:** the focused regression, both native queue-kanban suites, `go test ./...`, `go vet ./...`, `go build ./...`, Windows amd64 test-binary cross-compilation, repository contracts, formatting, and `git diff --check` pass.

## Review

**Acceptance: Pass.** The supported Unix and Windows paths now use documented single-operation replacement primitives; unsupported platforms fail before claiming atomicity. No Important or Minor findings remain.

## Lessons Learned

- A cross-platform API name does not imply a cross-platform atomicity contract; the last filesystem mutation needs an OS-specific proof.
- Validate write-target identity with `Lstat` before temporary-file creation when replacing a pathname rather than following it.

## Orientation

[MAP CHANGED] Queue-kanban complete-file writes now live in `atomic_write.go`; the final primitive is selected by the `atomic_replace_*` build-tag files and is shared by Testing frontmatter writes and `next-version`.
