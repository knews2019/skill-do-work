---
id: REQ-134
title: "Addendum: make queue-kanban atomic replacement cross-platform and symlink-safe"
status: pending
created_at: 2026-08-07T18:58:52Z
user_request: UR-032
addendum_to: REQ-072
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
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
