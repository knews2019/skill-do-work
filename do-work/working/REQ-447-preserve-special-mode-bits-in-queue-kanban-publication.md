---
id: REQ-447
title: '[impact-rule-change] Preserve special mode bits in queue-kanban publication'
status: claimed
created_at: 2026-08-31T20:30:00Z
status_changed_at: 2026-08-31T20:30:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: [REQ-436]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
addendum_to: REQ-436
sweep: true
sweep_key: preserve-special-mode-bits-in-file-publication
claimed_at: 2026-09-01T04:10:28Z
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-09-01T04:18:00Z
  basis:
    - Route B
    - 4-file write set
    - 1 new file
    - 2 modules involved
    - filesystem metadata and rollback behavior
    - cross-platform compile and canonical gates
---

# Preserve Special Mode Bits in Queue-Kanban Publication

## What

Extend REQ-436's complete-mode publication contract to the separate queue-kanban atomic writer at `skills/do-work-board/tools/queue-kanban/atomic_write.go`. It still narrows the original regular target with `Mode().Perm()` before replacement.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Accepted an exploration-guided four-file plan that closes both sweep instances with real-writer and forced-rollback full-mode fixtures.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Preserve ordinary permissions plus setuid, setgid, and sticky through queue-kanban's existing-target atomic replacement.
- Add RED/GREEN tests for each special bit using the real writer and full-mode assertions.
- Retain existing symlink/special-file refusal, changed-target detection, atomic replacement, and sync/close behavior.

## Instances

- [ ] `skills/do-work-board/tools/queue-kanban/atomic_write.go`: existing-target publication narrows the original complete mode with `Mode().Perm()`. (found by REQ-436 / UR-081)
- [ ] `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`: existing-untracked rollback snapshots only `Mode().Perm()` and recreates with `os.WriteFile`, stripping setuid, setgid, and sticky bits. (found by REQ-412 / UR-081)

## Triage

**Route: B** — Medium

**Reasoning:** The production delta is narrow but crosses two standalone Go modules and must retain distinct atomic-publication and transaction-rollback safety envelopes.

**Estimate:** p50 35 active minutes, medium confidence.

## Plan

1. Add real-seam RED tables for `04640`, `02640`, and `01640` through queue-kanban atomic replacement and forced existing-untracked transaction rollback, asserting bytes, full mode, outcome, and untracked state.
2. Preserve the sanitized complete regular-file mode (`Perm` plus setuid/setgid/sticky) in queue-kanban after writing and before publication; snapshot the same mode in gittransaction and apply it after rollback bytes are recreated.
3. Review the exact four-file diff, inventory remaining narrowing expressions, and run both focused/full modules, vet, exact Go 1.25, cross-platform compile, scope/diff, and canonical gates.

## Decisions

- Both enumerated Instances share the sweep root and must close here.
- Preserve ordinary permissions plus setuid, setgid, and sticky only; never copy file-type bits.
- Final chmod follows content writes so writes cannot clear special bits and rollback is umask-independent.
- Keep package-local projections rather than introducing a cross-module abstraction.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/atomic_write.go`
- `skills/do-work-board/tools/queue-kanban/atomic_write_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go`

**Files I will NOT touch:** queue-kanban release/platform replacement files, CLI atomicfile/requeststate/result/registration surfaces, actions/primes, release metadata, or any REQ-416 path.

**Acceptance criteria:**
- [ ] Queue-kanban atomic replacement preserves exact bytes and `04640`, `02640`, and `01640` modes.
- [ ] Existing-untracked rollback restores exact bytes, untracked state, and each complete mode.
- [ ] Existing refusal, identity, atomic replacement, cleanup, sync/close, transaction, and unrelated-dirt behavior remains green.
- [ ] Exact four-file scope and both module/cross-platform/canonical gates pass.

## Red-Green Proof

**RED prompt/case:** Replace queue-kanban fixtures carrying `04640`, `02640`, and `01640`, then assert exact contents and full low-twelve-bit modes.
**Why RED now:** The writer passes `originalInfo.Mode().Perm()` to chmod, the same narrowing shape REQ-426 and REQ-436 proved strips all three special bits.
**GREEN when:** All three modes and contents survive the real publication seam without weakening its existing safety fixtures.

---
*Source: discovered while implementing REQ-436 (UR-081).*
