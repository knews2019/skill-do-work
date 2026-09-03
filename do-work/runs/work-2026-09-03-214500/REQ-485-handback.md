# REQ-485 builder handback

## Result

- Branch: `worktree-agent-REQ-485-canonicalize-req-reservation-marker-filenames`
- Commit: `271aa8ae28660cfaa3a23b5181478e9e103895e4`
- Subject: `[REQ-485] Canonicalize reservation marker filenames`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-485-canonicalize-req-reservation-marker-filenames`
- Builder worktree is clean after the commit.

New writers now use the stored request ID's minimum-three-digit marker basename (`REQ-001`, `REQ-042`, `REQ-482`, `REQ-1207`). Every reservation reader changed by this request accepts an exact `REQ-` plus positive all-digit basename, so fixed-six legacy markers continue to block allocation, participate in publication checks/folds, and reap through cleanup without admitting suffix junk.

## Changed files

The commit matches the declared 15-file write set exactly:

- `skills/do-work-board/tools/queue-kanban/allocate.go` — canonical `%03d` writer and exact width-agnostic reservation parser for max scans.
- `skills/do-work-board/tools/queue-kanban/allocate_test.go` — literal basename contract plus canonical/legacy max-scan and malformed-suffix cases.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` — canonical marker guidance and legacy read compatibility.
- `skills/do-work/actions/capture.md` — canonical capture manifest marker path and explicit fixed-six compatibility text.
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` — canonical allocator writer and exact reservation discovery parser.
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` — literal writer paths and exact canonical/legacy discovery coverage.
- `skills/do-work/tools/do-work-cli/internal/publication/reservations.go` — shared canonical path derivation and rooted numeric-alias discovery.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` — canonical manifest enforcement and legacy/canonical alias collision refusal.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go` — true queue-kanban-to-capture collision lock-in plus legacy and noncanonical-path regressions.
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` — canonical manifest/create behavior and width-compatible collision/fold reads.
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` — legacy create collision and matching legacy fold coverage.
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` — delegated override-capture canonical-path refusal and fixture updates.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go` — exact all-digit positive cleanup matching while preserving removal revalidation.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations_test.go` — canonical/legacy coexistence, timeout, and malformed preservation coverage.
- `skills/do-work/tools/do-work-cli/internal/hookcommands/session_start_test.go` — canonical marker fixture at the registered cleanup seam.

No lifecycle metadata, run manifest, version, changelog, release file, cleanup launcher, board parser/UI, or other REQ was changed in the builder worktree.

## RED evidence

The new literal, independent regressions were run against the pre-change implementation before production edits:

- Board targeted tests failed because `requestReservationFileName` emitted `REQ-000042` rather than `REQ-042`, and `REQ-999-copy` incorrectly advanced allocation to 1000 instead of 43.
- Repository-model targeted tests failed because suffix junk was discovered as a reservation and allocator basenames remained `REQ-000013` / `REQ-000014`.
- Publication targeted tests failed because capture did not collide with the board-created spelling, a fixed-six legacy alias did not block canonical capture/defer create, and defer fold misclassified a present legacy alias as absent (`DEFER-GATE-REPAIR-COMMIT-AUTHORITY-MISSING`).
- Cleanup targeted tests failed because coexistence cleanup removed only `REQ-000482` and left `REQ-482`.

These failures were caused by literal fixture names and a cross-module process invocation, not expectations derived from the production formatter.

## GREEN and verification evidence

- `go test ./...` in `skills/do-work-board/tools/queue-kanban`: PASS.
- `go vet ./...` in `skills/do-work-board/tools/queue-kanban`: PASS.
- `go test ./...` in `skills/do-work/tools/do-work-cli`: PASS.
- `go vet ./...` in `skills/do-work/tools/do-work-cli`: PASS.
- Direct changed-package rerun, `go test ./internal/corehelpers ./internal/publication`: PASS.
- `_dev/tests/contract-regressions.sh`: PASS.
- `git diff --check`: PASS.
- Final staged diff and commit audit: exactly the 15 declared files; no debug artifacts; clean worktree after commit.

The first contract-regression run exposed an existing literal documentation guard requiring the legacy `REQ-NNNNNN` token. `capture.md` now names that spelling explicitly as read-only legacy compatibility while keeping `REQ-NNN` as the canonical writer/manifest contract; the unchanged contract suite then passed.

## Compatibility and risk notes

- Existing marker files are not renamed. They drain naturally under the existing committed-request or 48-hour cleanup policy.
- Readers accept any exact positive all-digit width, including canonical minimum-three-digit and fixed-six legacy spellings. Names with suffixes such as `REQ-482-copy`, zero identity, overflow, or non-digits are not reservation authority.
- If canonical and legacy aliases coexist, cleanup evaluates and revalidates each concrete file independently; both reap when eligible. Defer fold refuses ambiguous multiple aliases and still validates exact marker bytes when one alias exists.
- Rooted/exclusive-create allocation protections, reservation-directory identity checks, and cleanup's final Git authority/object identity/age revalidation remain intact.
- Current writers converge on one canonical `O_EXCL` target. A concurrently running pre-fix binary can still create a different legacy filename after a compatibility scan; atomicity against old software cannot be provided by two filesystem basenames, so deployments should upgrade all writers together.
- The board card/parser/UI does not consume reservation markers because hidden directories are excluded from its tree walk. No parser lock-step or browser-harness change was needed.
- Marker contents remain intentionally unchanged by this request: allocators may write empty markers while publication writes `REQ-ID\n`; cleanup remains name/mtime/Git based.

## Merge guidance

Merge or cherry-pick `271aa8ae28660cfaa3a23b5181478e9e103895e4` into the orchestrator-owned integration branch. After integration, rerun both Go modules' `go test ./...` and `go vet ./...`, then `_dev/tests/contract-regressions.sh`. Lifecycle completion, main-tree request movement, release/version/changelog decisions, and any worktree cleanup remain orchestrator-owned.
