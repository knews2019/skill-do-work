---
id: REQ-536
title: 'Review fix: record real file modes in lifecycle postimages'
status: completed
priority: now
created_at: 2026-09-03T12:20:21Z
user_request: UR-103
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-534, REQ-535]
batch: validate-feedback-2026-09-03
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
route: A
estimate:
  p50_active_minutes: 5
  rationale: "Fix PlannedPostimages to stat source files for accurate postimage mode projection, add unit and recovery lock-in tests"
write_set: [skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go, skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go]
claimed_at: 2026-09-04T12:55:13Z
completed_at: 2026-09-04T13:01:00Z
commit: 2f9b4b5bcfbd9fa2bb526338c041f98fe84a0cd5
---

# Record Real File Modes in Lifecycle Postimages

## What

`PlannedPostimages` in `requeststate` projects every lifecycle postimage (the REQ file, moved UR files, checkpoint, calibration) with mode `0644`. The mutation path preserves whatever mode the file already had: atomic replace copies the original mode and moves use `os.Rename`. The finalization journal compares bytes and mode when it classifies recovery state, so after a crash between the lifecycle mutation and the phase write, a file whose real mode is not 0644 matches neither the preimage nor the postimage. Recovery and rollback both refuse and the journal stays stuck.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this root cause. REQ-512 (complete legacy finalization semantic ownership) and REQ-507 (hand the archive and commit tails to finalize) reorganize who owns finalization; neither changes how postimage modes are projected.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** In `PlannedPostimages`, query the actual mode of existing source files (`plan.SourcePath`, `move.SourcePath`, existing checkpoint, existing calibration) using `os.Lstat` instead of hardcoded `0o644`. For newly created files fallback to `0o644`. Add unit test in `state_apply_test.go` and lock-in crash recovery test in `finalization_recovery_test.go`.
- [x] **[APPLY]:** Code written exactly as planned in `state_apply.go`, unit tested in `state_apply_test.go`, and lock-in tested in `finalization_recovery_test.go`.
- [x] **[UNIFY]:** Ran `git diff` and verified changes strictly within scope. Ran full test suite across entire package (`go test -count=1 ./...`), passing with 0 failures.

## Triage
- **Route:** A (single component bug fix, clear TDD red-green verification, no architectural changes)
- **Primary files:** `internal/requeststate/state_apply.go`
- **Tests:** `internal/requeststate/state_apply_test.go`, `internal/finalization/finalization_recovery_test.go`

## Plan
1. In `PlannedPostimages` in `state_apply.go`:
   - Inspect existing source mode for `plan.SourcePath` via helper `sourceMode(plan.SourcePath, 0o644)` using `os.Lstat` within `plan.RepositoryRoot`.
   - Apply `sourceMode` to destination postimages where the file was moved or updated in place (`plan.SourcePath == plan.DestinationPath` or `plan.DestinationPath`).
   - For `move` in `plan.AdditionalMoves`, use `sourceMode(move.SourcePath, 0o644)` for `move.DestinationPath`.
   - For `plan.CheckpointPath`, if `plan.CheckpointExisted` use `sourceMode(plan.CheckpointPath, 0o644)`, else `0o644`.
   - For `plan.CalibrationPath`, if `plan.CalibrationExisted` use `sourceMode(plan.CalibrationPath, 0o644)`, else `0o644`.
2. Add Red test in `finalization_recovery_test.go`:
   - Test `0600` REQ file crash recovery: prepare journal, execute lifecycle mutation directly (`ApplyPlan`), simulate crash by not writing `PhaseLifecycleApplied`, run `handleRecoverFinalization`.
   - Before fix: fails with `FINALIZATION-LIFECYCLE-CONFLICT` (matches neither preimage nor postimage).
   - After fix: succeeds with `PhaseCleanupComplete`.
3. Add unit tests in `state_apply_test.go` testing `PlannedPostimages` preservation of modes.

## Scope
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`

## Implementation Summary
- Updated `PlannedPostimages` in `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`:
  - Replaced hardcoded `0o644` with a dynamic `sourceMode` helper querying `os.Lstat` on the source path within `plan.RepositoryRoot`.
  - Propagated real source file mode to retained REQ paths, destination REQ paths, moved UR closure paths, and existing checkpoint/calibration paths.
  - Retained `0o644` default for newly created files (matching `writeCoupledFile` exclusive creation semantics).
- Added `TestPlannedPostimagesPreservesRealFileModes` in `internal/requeststate/state_apply_test.go`.
- Added lock-in test `TestRecoverFinalizationPreservesAndMatchesPrivateFileModeInLifecyclePostimages` in `internal/finalization/finalization_recovery_test.go`.

## Testing
- `go test -v -run TestRecoverFinalizationPreservesAndMatchesPrivateFileModeInLifecyclePostimages ./internal/finalization`: verified RED before fix (failed with `FINALIZATION-LIFECYCLE-CONFLICT`), verified GREEN after fix.
- `go test -v -run TestPlannedPostimagesPreservesRealFileModes ./internal/requeststate`: verified GREEN for all permutations.
- `go test -count=1 ./...`: verified full package test suite passes.

## Review
- **Independent Verification:** Completed against all requirements.
  - Requirement 1 (retained REQ, moved UR, existing checkpoint/calibration retain source mode, fresh keeps 0644): Verified.
  - Requirement 2 (single owner): Verified in `PlannedPostimages`.
  - Requirement 3 (recovery converges for 0600 REQ): Verified in `TestRecoverFinalizationPreservesAndMatchesPrivateFileModeInLifecyclePostimages`.
  - Requirement 4 (lock-in test): Verified.
- **Defects:** 0 critical, 0 non-critical.
- **Score:** 100%.

## Lessons Learned
- When journal recovery compares files by exact `Mode` bits, projections of postimages must accurately reflect the mutations performed by `atomicfile` and `os.Rename`, which preserve pre-existing file permissions.

## Context

Finding provenance, carried per the Finding-Closure Ratchet:

- **Source:** external review comment, severity P1, adjudicated by `do-work validate-feedback` on 2026-09-03 with verdict Accept; full block preserved in UR-103 input.md.
- **Verbatim claim:** "[P1] Record lifecycle postimage modes accurately — state_apply.go:146-150. When a lifecycle source, moved UR, checkpoint, or calibration file has a mode other than 0644, ApplyPlan preserves that existing mode while these journal postimages claim 0644. A crash after lifecycle mutation but before the phase write then leaves files matching neither the preimage nor postimage, so recovery and rollback both refuse and the finalization journal becomes stuck."
- **Evidence:** `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:147-160` hardcodes `0o644` for every postimage. `internal/atomicfile/atomic_file.go:55` chmods the replacement to the original's mode; `state_apply.go:64,79` move with `os.Rename`, which keeps the mode. Preimages record the real mode (`internal/finalization/finalization_prepare.go:222`); `equalImage` compares mode (`internal/finalization/finalization_apply.go:419`); `imageSetState` refuses a file that matches neither image (`finalization_apply.go:311-312`). Release postimages already copy the preimage mode when a path existed (`finalization_prepare.go:236-238`). Verified by reading the chain; not reproduced live.
- **Surface-cost:** N/A, direct bug fix; no guard, fallback, or validation layer is added.

## Requirements

- A lifecycle postimage for a file that already exists (retained REQ file, moved UR file, existing checkpoint or calibration file) records the mode the file will actually have after `ApplyPlan`: the existing file's mode. A file `ApplyPlan` creates fresh keeps `0644`, matching `writeCoupledFile`'s create mode.
- Either `PlannedPostimages` stats the source and carries the real mode, or `finalization_prepare.go` fills the lifecycle postimage mode from the matching preimage when the preimage existed, the way release postimages already do. One place owns the rule; do not do both.
- Recovery after a simulated crash between the lifecycle mutation and the phase write converges for a REQ file with mode `0600`: `imageSetState` classifies it as `post`, not "matches neither".
- Closing check: one lock-in test naming this failure (REQ file with mode 0600, journal prepared, lifecycle applied, phase not yet persisted; recovery must classify the file as post and complete instead of refusing).

## Red-Green Proof
**RED prompt/case:** Prepare a finalization for a REQ whose queue file has mode `0600`, apply the lifecycle plan, and stop before the journal phase write (the crash point). Run finalization recovery.
**Why RED now:** The archived REQ file has mode 0600 while the journal postimage claims 0644, so recovery fails with "matches neither the journal preimage nor postimage" and rollback refuses for the same reason; the journal is stuck.
**GREEN when:** Recovery classifies the file as matching the postimage and completes the finalization; a lock-in test under `internal/requeststate` or `internal/finalization` pins the 0600 case.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 3124 tokens, over the 2000-token budget; `slugged: partial`, so the targeted `#final-boundary-identity` form is not eligible. Matched on the final-boundary-identity family (identity of a file at a transaction boundary). The owning prime is listed in `prime_files` instead.

## Assets
None.

---
*Source: `do-work validate-feedback` triage of 2026-09-03, Finding 5 (Accept); full block preserved in UR-103 input.md.*

