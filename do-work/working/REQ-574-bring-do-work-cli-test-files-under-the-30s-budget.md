---
id: REQ-574
title: 'Repository gate repair: bring do-work-cli test files under the 30s per-file budget'
status: claimed
created_at: 2026-09-04T23:50:46Z
user_request: UR-115
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: do-work-cli-test-file-budget
depends_on: []
related: [REQ-572]
estimate:
  p50_active_minutes: 40
  confidence: low
  calculated_at: 2026-09-05T00:00:27Z
  basis:
    - Route C
    - 4-file write set
    - 3 subsystems involved
    - 2 acceptance criteria
    - cross-route regression gates
    - full-suite verification
status_changed_at: 2026-09-05T09:44:26Z
claimed_at: 2026-09-05T09:44:40Z
write_set:
  - skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go
route: C
---

# Repository gate repair: bring do-work-cli test files under the 30s per-file budget

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: go-test-file-budget:do-work-cli:publication-defer-gate-test affecting REQ-572 (found by REQ-572 / UR-115)

## Repository Gate Repair Intake

- **Parent:** REQ-572
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** go-test-file-budget:do-work-cli:publication-defer-gate-test
- **Repair dependency:** REQ-574
- **Diagnostic evidence:** "post-merge run at 4adcff4e (fbdcd35e merged): FAIL: internal/corehelpers/inventory_test.go accumulated 38.92s; internal/publication/defer_gate_test.go 37.01s; internal/finalization/finalization_recovery_test.go 35.65s; internal/finalization/finalization_req499_test.go 30.85s; each test file must finish under 30s; every test passed"
- **Diagnostic evidence:** "pre-build run at f6c43d22: the same four files over budget (38.61s, 37.93s, 35.73s, 33.72s)"
- **Diagnostic evidence:** "detached diagnostic worktree at base 7ad53bff (clean tree): FAIL: internal/publication/defer_gate_test.go accumulated 32.52s; each test file must finish under 30s; every test passed; queue-kanban 24s"
- **Implementation base:** 7ad53bff1d867f1453e1e7765e988dedb308e7e1
- **Implementation merge:** fbdcd35e0908aca6a01f554cc9b7fd7c85347a49

---

---

## Triage

**Route: C** - Complex

**Reasoning:** A repository-gate repair minted by defer-gate with a `route: C` preset. Four test files across three do-work-cli packages (publication, finalization, corehelpers) sit at 18-20s each when the gate runs alone and cross the 30s per-file budget whenever a second gate process shares the machine, so the work is to find where the time goes and cut it without weakening what the tests pin.

**Planning:** Required

## Plan

**Where the time goes.** Measured directly on an otherwise idle machine, per source test file, summing the top-level test wall times exactly as `_dev/tests/run-go-tests-with-budget.sh` does:

| file | cold cache | warm cache | dominant test |
|---|---|---|---|
| `internal/corehelpers/inventory_test.go` | 39.10s | 17.25s | `TestInventoryMatchesRetainedPorcelainXYMatrix` 13.19-31.28s |
| `internal/publication/defer_gate_test.go` | 37.81s | 22.95s | `TestDeferGateRollsBackUntrackedCreateAndFoldTopologies` 9.08-15.85s |
| `internal/finalization/finalization_recovery_test.go` | 31.68s | 19.87s | `TestRecoverFinalizationResumesEveryDurablePhaseExactlyOnce` 7.13-12.05s |
| `internal/finalization/finalization_req499_test.go` | 25.84s | 25.30s | `TestRecoverFinalizationRequiresWorkspaceMemberLockMirrors` 7.09-8.59s |

The cost is process spawning, not computation. Every case builds a real git repository in its own `t.TempDir()` and drives it with individual `git` subprocesses — `newDeferGateRepository` alone runs seven of them per case and is called fifteen times inside one test. The 45-case porcelain matrix additionally writes a fake `git` shell script per case and runs the retained `uncommitted-inventory.sh` through it.

**The measured quantity is the parent test's wall time.** The budget script sums `Elapsed` only for tests whose name has no `/` in it, so subtests are already excluded and a parent's number is simply how long its subtests took end to end. Every case in these matrices owns a separate temporary repository and shares nothing, so running them concurrently cuts the parent's wall time without weakening a single assertion. The machine has 8 logical CPUs.

**Tasks.**

1. **`internal/corehelpers/inventory_test.go` — make the 45-case porcelain matrix concurrent.** `runRetainedInventory` calls `t.Setenv` for `DO_WORK_INVENTORY_STATUS_FIXTURE` and `PATH`, and Go panics if a test calls `t.Setenv` and `t.Parallel`. Pass both variables on the child's own `exec.Cmd.Env` instead — they exist only to reach that one subprocess, so process-wide mutation was never needed — then mark the matrix subtests parallel.

2. **`internal/publication/defer_gate_test.go` — make the rollback topologies concurrent.** This file touches no package-level variable and no environment, so its subtests are already independent. Give each `failureIndex` iteration its own named subtest rather than a bare loop, and mark them parallel.

3. **`internal/finalization/finalization_req499_test.go` — parallelize the table-driven cases that own no global.** `TestRecoverFinalizationRequiresWorkspaceMemberLockMirrors` builds one repository per case and reassigns nothing package-level.

4. **`internal/finalization/finalization_recovery_test.go` — parallelize everything that does not reassign `afterFinalizationPhase`.** Three tests swap that package-level hook (`finalization_apply.go:25`) and must stay serial; the rest of the file's 25 tests may run concurrently. `finalization_req499_test.go` has the same constraint at its line 288 (`enumerateTrackedReleasePaths`) and 510.

**Acceptance.** Each of the four files finishes under 30s in the canonical gate, with every test still executed and unskipped, and no assertion removed or weakened. Splitting a file to spread its total across two names is explicitly rejected: it would move the number without making anything faster.

*Written inline by the orchestrator; no Plan agent was dispatched.*

## Pre-Build Repository Gate

- **Expected diagnostic fingerprint:** go-test-file-budget:do-work-cli:publication-defer-gate-test
- **Gate command:** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 0 (green), run at 9f914188 in a detached worktree
- **Observed per-file result:** `inventory_test.go` 28.02s, `finalization_recovery_test.go` 27.85s, `defer_gate_test.go` 27.25s, `finalization_req499_test.go` 21.79s, against a 30s limit
- **Two earlier runs at this REQ's claim were red on other fingerprints and never reached this lane:** SC2148 on a newly tracked probe script (fixed at 4d47c821), then 23 broken archive references in `lessons-do-work-cli.md` (fixed at 9f914188). Neither is the recorded budget fingerprint.

**Why this run does not take the already-green no-op branch.** That branch's stated premise is that "the defect was repaired elsewhere". Nothing repaired it. The three named files sit 1.98s, 2.15s and 2.75s under a hard 30s limit, the same files were recorded at 37-60s earlier the same day, and `do-work/test-durations.tsv` shows the numbers rising across the day as tests were added. A budget check passing by two seconds is not a repaired budget, and closing this REQ on that reading would return REQ-572 and REQ-573 to a gate that tips over the next time two sessions share the machine. The deviation is deliberate and recorded here rather than silently taken.

## Scope

**Files I will touch:**

- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go`
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go`

**Acceptance criteria:**

1. Each of the four files finishes well under the 30s per-file budget in the canonical gate, with enough headroom that a concurrent second gate does not push it over.
2. Every test still runs unskipped, and no assertion is removed, loosened, or moved to another file.
