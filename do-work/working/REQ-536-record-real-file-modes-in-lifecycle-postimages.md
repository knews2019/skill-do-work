---
id: REQ-536
title: 'Review fix: record real file modes in lifecycle postimages'
status: claimed
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
write_set: [skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go, skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go]
claimed_at: 2026-09-04T12:55:13Z
---

# Record Real File Modes in Lifecycle Postimages

## What

`PlannedPostimages` in `requeststate` projects every lifecycle postimage (the REQ file, moved UR files, checkpoint, calibration) with mode `0644`. The mutation path preserves whatever mode the file already had: atomic replace copies the original mode and moves use `os.Rename`. The finalization journal compares bytes and mode when it classifies recovery state, so after a crash between the lifecycle mutation and the phase write, a file whose real mode is not 0644 matches neither the preimage nor the postimage. Recovery and rollback both refuse and the journal stays stuck.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this root cause. REQ-512 (complete legacy finalization semantic ownership) and REQ-507 (hand the archive and commit tails to finalize) reorganize who owns finalization; neither changes how postimage modes are projected.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
