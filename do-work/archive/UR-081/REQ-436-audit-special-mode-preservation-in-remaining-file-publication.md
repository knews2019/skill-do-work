---
id: REQ-436
title: '[impact-negligible] Audit special-mode preservation in remaining file publication'
status: completed
claimed_at: 2026-08-31T19:48:00Z
route: B
completed_at: 2026-08-31T20:30:00Z
commit: f0715c41
created_at: 2026-08-31T10:56:05Z
status_changed_at: 2026-08-31T13:52:30Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
depends_on: [REQ-426]
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
addendum_to: REQ-426
sweep: true
sweep_key: preserve-special-mode-bits-in-file-publication
write_set:
  - skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go
  - skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file_test.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-31T19:48:00Z
  basis:
    - trivial short-circuit
---

# Audit Special-Mode Preservation in Remaining File Publication

## What

REQ-426 fixed two managed install paths that silently narrowed Unix modes to the low nine permission bits. The same `Mode().Perm()` publication shape remains in atomic replacement and cleanup moves. Audit both contracts, preserve setuid/setgid/sticky where they promise to preserve a source file's mode, and pin the class so it cannot recur.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Classified both named publication seams and the hidden downstream `CreateExclusiveAt` mask; froze the four-file production/test scope and RED/GREEN matrix before coding.
- [x] **[APPLY]:** Added RED-first replacement, exclusive-create, and real cleanup-move fixtures, then fixed both production masks inside the exact four-file scope.
- [x] **[UNIFY]:** Reviewed all four files; passed gofmt, focused/full CLI tests, vet, exact Go 1.25, Windows atomic compile, canonical verification, mask inventory, and diff/scope hygiene with no debug artifacts.

## Instances

- [x] `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go:55` narrows the original target mode before publishing its replacement.
- [x] `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go:245` narrows the source mode before publishing a cleanup move destination.

## Red-Green Proof

**RED prompt/case:** Replace and move regular-file fixtures carrying setuid, setgid, and sticky bits through the two named public seams.
**Why RED now:** Both production calls pass `Mode().Perm()`, the same low-nine-bit mask that REQ-426 proved strips all three special bits.
**GREEN when:** Each contract either preserves the complete mode with RED/GREEN regression proof, or explicitly documents and tests why narrowing is intentional; the audit contains no unclassified `Mode().Perm()` publication path.
**Validation:** Discovered during REQ-426 implementation; apply the finding-closure ratchet.

## Open Questions

- [x] I found two remaining file-publication paths with the same special-mode-bit narrowing pattern fixed by REQ-426. Should I process this low-reach audit as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to `pending`) so every mode-preservation promise uses the same complete-mode contract.
  Also: No, discard it; these bits are rarely set on the affected files and REQ-426 already closes the reported installer paths.
  Value: prevents the same silent metadata loss in atomic replacement and cleanup moves.
  Risk: low and reversible; the work is a focused two-path audit with regression tests, but it adds queue work for an uncommon filesystem edge case.

  **Answered 2026-08-31** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the recommendation via `do-work clarify`: add the focused audit to the queue
  so the two remaining publication paths preserve complete special mode bits and carry regression
  proof. Discarding the uncommon edge case is out of scope.

---
*Source: discovered while implementing REQ-426 (UR-081).*

---

## Triage

**Route: B** - Medium

**Reasoning:** The change is mechanically small but crosses two publication subsystems and needs focused RED/GREEN mode-bit fixtures plus a class audit. Exploration should freeze the exact production and test files before implementation.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The atomic replacement seam writes before chmod, but masks the original mode to its low nine bits. Cleanup loses the same special bits twice: `moveWithoutOverwrite` passes `sourceInfo.Mode().Perm()`, then `CreateExclusiveAt` applies `fileMode.Perm()` again. The complete-mode contract should preserve ordinary permissions plus `ModeSetuid`, `ModeSetgid`, and `ModeSticky`; exclusive creation must apply the final mode after writing and before sync so content writes cannot clear setuid/setgid.

REQ-426 supplies the numeric-Unix-to-Go-mode test convention. The real cleanup move must be exercised end to end; fixing or testing only one mask cannot pass. Existing managed-section and suite-install paths are already compliant, and literal reservation-marker modes are not source-mode preservation. A similar queue-kanban site belongs to a separate module and is outside this named CLI audit.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go` (modify) — preserve complete modes in replacement and exclusive publication.
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file_test.go` (modify) — RED/GREEN replacement and exclusive-create mode fixtures.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (modify) — pass the complete source mode into publication.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modify) — prove the real safe move retains contents and all special bits.

**Files I will NOT touch:** already-compliant managedsection/suiteinstall paths, literal-mode repository markers, action/prime files, queue state from the builder, release metadata, or the separate queue-kanban module.

**Acceptance criteria (restated from REQ):**
- [x] Atomic replacement retains ordinary permissions plus setuid, setgid, and sticky from the original regular target.
- [x] A cleanup move retains ordinary permissions plus all three special bits at the actual destination without weakening no-overwrite or source-safety behavior.
- [x] The downstream `CreateExclusiveAt` mask cannot erase the complete mode passed by cleanup.
- [x] Every scoped `Mode().Perm()` publication use is fixed or explicitly classified as intentional.

## Pre-Flight

**Git:** ⚠ Five unrelated concurrent main-tree edits are preserved and excluded (`_dev/tests/select-simple-reqs-behavior.sh`, `_dev/tests/update-script-behavior.sh`, `internal/archivefetch/archive_fetch.go`, `internal/archivefetch/archive_fetch_test.go`, and `internal/dependencygraph/dependency_graph.go`); the builder worktree starts from clean committed `de66214f`.
**Tests baseline:** ✓ The canonical gate passed at the immediately preceding implementation state; the only later commit contains claim bookkeeping.
**Dependencies:** ✓ Go module dependencies and exact Go 1.25 compatibility tooling are available.

*Checked by work action*

## Implementation Summary

**What was done:** Atomic replacement and cleanup moves now preserve ordinary permissions plus setuid, setgid, and sticky bits. Exclusive creation uses ordinary permissions for the rooted create call, then applies the sanitized complete mode after writing and before sync; the real cleanup move passes the source's complete `FileMode` through that boundary.

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file.go` (modified) — shared complete-mode projection and post-write publication.
- `skills/do-work/tools/do-work-cli/internal/atomicfile/atomic_file_test.go` (modified) — setuid, setgid, and sticky replacement/exclusive-create fixtures.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (modified) — complete source mode passed to atomic publication.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modified) — end-to-end real move proof for every special bit.

**Integration range:** `7207daef..f0715c41`

*Generated by work action from the builder hand-back*

## Decisions

### D-01: Centralize the complete Unix publication subset

**Decision:** DECIDE & STATE — use one unexported atomicfile helper that combines ordinary permissions with only setuid, setgid, and sticky.

**Reasoning:** Replacement and exclusive publication now share one explicit contract without admitting file-type bits. Value: the two named seams cannot drift. Risk: future special-mode support must extend the helper and fixtures together.

### D-02: Apply special bits after content write

**Decision:** DECIDE & STATE — pass low-nine permissions to `os.Root.OpenFile`, then chmod the complete mode after `Write` and before `Sync`.

**Reasoning:** The rooted API rejects high special flags, and content writes may clear setuid/setgid. Value: portable creation plus correct final Unix metadata. Risk: reordering publication must preserve the tested write/chmod/sync sequence.

## Qualification

Passed — all four declared files are substantive in `7207daef..f0715c41`, the helper's mode projection reaches both publication seams, the real cleanup path crosses the exclusive-create boundary, requirements trace to RED/GREEN fixtures, and P-A-U/debug/scope checks are clean.

## Testing

**Red-green validation:** Before production changes, the replacement, exclusive-create, and real cleanup-move cases each collapsed `04640`, `02640`, and `01640` to `0640`. The same fixtures pass after the complete-mode and post-write chmod changes.

**Merged-state checks:**
- `go test -count=1 ./internal/atomicfile ./internal/cleanup` — PASS.
- Full do-work-cli `go test -count=1 ./...` and `go vet ./...` — PASS.
- Exact Go 1.25 compatibility — PASS.
- Windows atomicfile compile — PASS after rerunning from the Go module directory; the initial repository-root invocation found no module and changed nothing.
- Qualification, scope drift, and `git diff --check` over `7207daef..f0715c41` — PASS.
- `bash _dev/tests/maintainer-verify.sh` on the merged tree — PASS; the optional browser lane skipped because no browser was available.

## Review

**Overall: 98%** | 2026-08-31T20:28:08Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 94% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.
**Minor findings:** None.
**Acceptance:** Pass — both publication seams preserve complete special modes without weakening exclusive creation, target/source revalidation, rollback, or deletion ordering.
**Suggested testing:** Optionally repeat the focused Unix fixtures under a restrictive umask; current post-write chmod ordering is designed to make the result independent of it.
**Follow-ups created:** REQ-447 for the analogous queue-kanban publication mask; **sweeps appended to:** `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** End-to-end fixtures across replacement, exclusive creation, and the real cleanup move exposed both masks and proved the final published metadata rather than an intermediate argument.
**What didn't:** Fixing the caller alone would still have lost the bits in `CreateExclusiveAt`, and applying the high flags directly to rooted creation was rejected by the API.
**Worth knowing:** For source-mode-preserving publication, sanitize to ordinary permissions plus the three special flags, create with ordinary permissions, write content, then apply the complete mode before syncing.

## Orientation

Atomic CLI publication now preserves complete Unix permission metadata through replacement and cleanup moves. The shared guarantee lives in `internal/atomicfile`; the separate queue-kanban atomic writer is queued as REQ-447 rather than folded across module boundaries.
