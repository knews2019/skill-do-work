---
id: REQ-147
title: "Addendum: reserve request numbers during allocation"
status: completed
completed_at: 2026-08-07T22:49:55Z
commit:
claimed_at: 2026-08-07T22:41:26Z
status_changed_at: 2026-08-07T22:41:26Z
created_at: 2026-08-07T19:15:15Z
user_request: UR-033
addendum_to: REQ-072
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-134]
maintenance: false
route: B
kb_status: pending
kb_entry:
---

# Addendum: Reserve Request Numbers During Allocation

## What

Change `queue-kanban next-req` from a read-only maximum scan into an atomic reservation operation so every successful invocation owns a distinct REQ number before it prints that number.

## Context

REQ-134 was briefly observed as a duplicate while UR-031 and UR-032 were captured concurrently. The committed metadata now has a unique mapping—UR-032 owns REQ-134 and UR-031 owns REQ-135 through REQ-146—but the collision exposed a real allocator defect: calculating `max(existing REQ)+1` does not reserve the result, so two callers can receive the same identifier before either creates its request file.

## Prior Implementation

REQ-072 introduced `next-req` in commit `5db22ea` as a read-only allocator that scans existing request filenames and frontmatter. That implementation deliberately avoided writes, which means sequential or concurrent calls made before request creation can return the same number.

## Requirements

- Atomically reserve a REQ number before returning it from `queue-kanban next-req`.
- Count both existing request records and prior reservations when selecting the next number.
- Ensure sequential calls and concurrent processes receive distinct numbers without relying on a long-lived global lock.
- Use an exclusive filesystem operation so competing allocators retry safely instead of overwriting another reservation.
- Keep abandoned reservations as accepted gaps; never recycle a number merely because its request file has not appeared.
- Keep command output semantics unchanged: on success, print exactly one decimal request number.
- Keep reservation paths contained under the repository metadata tree and reject unsafe symlink or path-escape conditions.
- Update the capture workflow contract so a successful capture stages its reservation marker with the UR and REQ records.
- Add regression tests for sequential allocation, concurrent allocation, existing reservations, and unsafe reservation paths.
- Run `go test ./...`, `go vet ./...`, and the relevant contract regression checks.

## Constraints

- Limit implementation changes to the queue-kanban allocator and directly coupled capture contract.
- Do not modify `contact_processor`, generated assets, unrelated skill documentation, cards, warnings, or pipeline fields.
- Do not introduce a changelog write or change any other queue-kanban command semantics.

## Red-Green Proof

**RED prompt/case:** Invoke `queue-kanban next-req` twice before writing a request, and run multiple allocator processes concurrently against the same repository.

**Why RED now:** A maximum-only scan has no durable claim step, so callers can observe identical state and all print the same next number.

**GREEN when:** Every successful call leaves a durable, exclusively created reservation; sequential and concurrent calls return unique numbers; existing REQs and markers are both honored; unsafe marker paths fail closed; and the Go plus capture-contract verification suites pass.

**Validation:** User explicitly required the Go allocator to reserve numbers so the next call receives a different ID.

## Assets

None.

---
*Source: the go app should reserve the numbers, so the next call gets a different id*

## Triage

**Route: B** — the desired behavior and allocator entry point are explicit, while the existing scan, capture contract, mirrored runtime copy, and safe filesystem pattern require repository discovery.

## Scope

**Files I will touch:**

- `tools/queue-kanban/allocate.go`
- `tools/queue-kanban/allocate_test.go`
- `tools/queue-kanban/main.go`
- `tools/queue-kanban/frontmatter_cli.go`
- `tools/queue-kanban/prime-do-kanban.md`
- `skills/do-work-board/tools/queue-kanban/allocate.go`
- `skills/do-work-board/tools/queue-kanban/allocate_test.go`
- `skills/do-work-board/tools/queue-kanban/main.go`
- `skills/do-work-board/tools/queue-kanban/frontmatter_cli.go`
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`
- `actions/capture.md`
- `skills/do-work/actions/capture.md`
- `CLAUDE.md`
- `_dev/tests/contract-regressions.sh`
- `VERSION`
- `actions/version.md`
- `CHANGELOG.md`
- `skills/do-work/VERSION`
- `skills/do-work/actions/version.md`
- `skills/do-work/CHANGELOG.md`

**Files I will not touch:** generated board assets, pipeline semantics, application code, unrelated documentation, REQ-144's gated cutover, or the unrelated dirty REQ-146/UR-031 planning edits.

## Plan

1. Prove the current maximum-only allocator returns duplicate ids before either request exists.
2. Reserve ids with contained exclusive-create markers and retry collisions without a global lock.
3. Prove sequential and cross-process uniqueness, prior-marker handling, unsafe-path rejection, unchanged decimal output, and capture staging.
4. Mirror the shipped runtime sources, qualify the complete change, release it, and archive UR-033.

## Pre-Flight

The archive collision check passed. Existing REQ files are currently the allocator's only durable input; repeated calls against an unchanged queue return the same number and leave no claim marker. The worktree also contains unrelated user edits, which will be preserved and excluded from this REQ.

## Decisions

- Use one empty, fixed-width marker per number with `O_CREATE|O_EXCL`; a losing process advances rather than waiting on a global lock.
- Keep every successful marker permanently. A stopped capture consumes an id, making a harmless gap instead of a potentially reused public identifier.
- Open the resolved repository and reservation directory through `os.Root`. This preserves containment across path changes after validation, while explicit `Lstat` checks reject a symlinked reservation store.
- Keep the no-Go manual capture fallback, but require it to count existing markers and re-scan before each write because it cannot make an atomic reservation.

## Implementation Summary

Implemented durable, process-safe REQ id reservation and kept the bridge and staged modular distributions synchronized.

- `tools/queue-kanban/allocate.go` (modified) — counts REQ records plus markers, opens a contained reservation root, and retries exclusive per-number marker creation.
- `tools/queue-kanban/allocate_test.go` (modified) — covers sequential calls, sixteen concurrent OS processes, existing markers, symlink/path escapes, queue preservation, and one-decimal subprocess output.
- `tools/queue-kanban/main.go` (modified) — documents the reservation command contract.
- `tools/queue-kanban/frontmatter_cli.go` (modified) — records the complete write-surface count.
- `tools/queue-kanban/prime-do-kanban.md` (modified) — documents reservation behavior and traps.
- `skills/do-work-board/tools/queue-kanban/allocate.go` (modified) — mirrors the allocator implementation.
- `skills/do-work-board/tools/queue-kanban/allocate_test.go` (modified) — mirrors the allocator regression coverage.
- `skills/do-work-board/tools/queue-kanban/main.go` (modified) — mirrors the command contract.
- `skills/do-work-board/tools/queue-kanban/frontmatter_cli.go` (modified) — mirrors the write-surface contract.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified) — mirrors the behavior while retaining modular sibling paths.
- `actions/capture.md` (modified) — counts and stages reservation markers.
- `skills/do-work/actions/capture.md` (modified) — mirrors the modular capture contract.
- `CLAUDE.md` (modified) — records the third queue-kanban write surface.
- `_dev/tests/contract-regressions.sh` (modified) — ratchets the capture and write-surface contracts.
- `VERSION` (modified) — bumps the bridge release to v0.183.22.
- `actions/version.md` (modified) — reports v0.183.22.
- `CHANGELOG.md` (modified) — describes atomic request number reservations.
- `skills/do-work/VERSION` (modified) — keeps the staged core release synchronized.
- `skills/do-work/actions/version.md` (modified) — keeps staged version reporting synchronized.
- `skills/do-work/CHANGELOG.md` (modified) — mirrors the release notes.

## Testing

**RED:** A binary built from the pre-change `HEAD` returned `43` for two sequential calls against a queue ending at REQ-042 and created zero reservation markers.

**GREEN:** Focused allocator tests pass for sequential and cross-process allocation, prior markers, a symlinked reservation directory, and a `do-work` path escaping the repository. The concurrency case re-enters the compiled test binary in sixteen separate OS processes and accepts only output parseable as exactly one decimal number.

Full root and staged queue-kanban `go test ./...`, `go vet ./...`, and `go build ./...` passes; the Windows amd64 test binary cross-compiles; warning-level ShellCheck, Bash syntax, Just parsing, all repository contract suites, queue-kanban verification, gofmt, scope drift, and `git diff --check` pass.

## Qualification

Passed. The implementation is substantive, every requirement traces to allocator logic, a regression, or the capture contract, and the only new data flow is an empty marker created through a repository-rooted handle. The staged board copy is byte-identical for Go sources, and its prime retains the required modular sibling paths.

## Review

**Acceptance: Pass.** `O_CREATE|O_EXCL` is the single concurrency boundary; no process can overwrite another marker, abandoned markers remain authoritative, and `os.Root` prevents a validated path from escaping during the write. Ordinary queue files and every other command remain unchanged. No Important or Minor findings remain.

## Lessons Learned

- A read-only max scan cannot allocate an identifier; allocation needs a durable ownership event before output.
- Per-id exclusive markers avoid stale-lock recovery entirely: a stopped caller creates a safe gap, not a lock that another process must guess how to break.
- Path containment should survive the interval after validation, so rooted filesystem handles are stronger than joining an already-checked absolute path.

## Orientation

[MAP CHANGED] `queue-kanban next-req` is now the third write surface. Its only mutation is an empty marker under `do-work/.req-reservations/`, which capture retains and stages with the matching REQ/UR records.
