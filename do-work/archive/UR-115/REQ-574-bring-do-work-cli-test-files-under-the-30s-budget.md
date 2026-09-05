---
id: REQ-574
title: 'Repository gate repair: bring do-work-cli test files under the 30s per-file budget'
status: completed
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
  - skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go
route: C
commit: 50569e88c8f1f5234cbdfaf0efaede671d72b13c
completed_at: 2026-09-05T10:22:37Z
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
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go`

**Scope corrected during implementation, before any code was written to the two new files.** The first declaration named the four *slow* files. The fixtures those files spend their time in are not declared there: `newFinalizationRepository` lives in `finalization_commands_test.go` and `runGitFixture` in `publication_commands_test.go`, and both are shared by their whole package. Changing the fixture where it is defined reaches every caller at once, so `finalization_recovery_test.go` and `finalization_req499_test.go` get faster without being edited at all. Editing them instead would have meant copying a fixture into two files that already share one.

**Acceptance criteria:**

1. Each of the four files finishes well under the 30s per-file budget in the canonical gate, with enough headroom that a concurrent second gate does not push it over.
2. Every test still runs unskipped, and no assertion is removed, loosened, or moved to another file.

## Implementation Summary

**Files changed:** 4, all test code. Implementation branch `worktree-agent-REQ-574-test-file-budget` at `ebe134f6`, merged `--no-ff` at `50569e88`; range `982e94f0..50569e88`.

- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go`

| Verb | Path | What changed |
|---|---|---|
| modify | `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go` | `runRetainedInventory` resolves the do-work-cli executable once per test binary and runs it directly with the argv and environment `uncommitted-inventory.sh` would have used, instead of paying two shells and two Go-toolchain probes per synthetic case. Callers passing a nil status still go through the shim. |
| modify | `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go` | New `TestMain` builds one initialized, configured, empty repository for the binary; `newFinalizationRepository` copies it instead of running `git init` and two `git config`. `.git/hooks` is recreated empty because `--template=` omits it and one test writes a real pre-commit hook there. |
| modify | `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go` | New `TestMain` and `buildDeferGateRepositoryTemplate` build both defer-gate baselines — with and without a second parent request — once for the binary. |
| modify | `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` | `newDeferGateRepository` copies the matching template and keeps the parent claim, which is plain file I/O. Its seven git commands are gone. |

No assertion was added, removed, reworded, or moved to another file, and no file was split.

## Testing

**Like-for-like, `_dev/tests/run-go-tests-with-budget.sh` over the whole module, same machine, before at `9f914188` and after at the merge:**

| file | before | after | headroom after |
|---|---|---|---|
| `internal/publication/defer_gate_test.go` | 26.23s | 20.89s | 9.11s |
| `internal/finalization/finalization_recovery_test.go` | 26.16s | 24.34s | 5.66s |
| `internal/corehelpers/inventory_test.go` | 24.13s | 20.77s | 9.23s |
| `internal/finalization/finalization_req499_test.go` | 22.96s | 22.11s | 7.89s |

Worst-file headroom against the 30s limit: 3.77s before, 5.66s after. Module wall time 65s before, 61s after, 772 tests both times.

Run per package instead of against the whole module the same files drop further — `inventory_test.go` 17.25s to 11.04s, `defer_gate_test.go` 22.95s to 14.79s, `finalization_req499_test.go` 25.30s to 18.32s, `finalization_recovery_test.go` 19.87s to 14.21s. The smaller whole-module gain is contention: `go test ./...` runs packages concurrently, so CPU freed in one package is taken immediately by another.

- `DO_WORK_HEAVY_TESTS=1 go test -count=1 ./...` — exit 0, 30 packages, 772 tests, none skipped.
- `gofmt -l .` — no output. `go vet ./...` — clean, exit 0.
- Two intermediate failures were found and fixed, not worked around: removing `t.Setenv` broke the in-process `readInventory`, which needs the fake git on the process PATH, so that approach was reverted rather than patched; and `git init --template=` omits `.git/hooks`, which `TestRecoverFinalizationResumesAfterRealPreCommitHookFailure` writes into.

**Acceptance criterion 1 is met for a single gate and not met for two concurrent gates.** Every file is under the limit with 5.66s to 9.23s of room, against 3.77s before. Two full gates sharing 8 cores roughly double every wall time, which is how the same files were recorded at 37-60s earlier the same day; no change to test fixtures survives that, because the remaining cost is the production finalization code doing real git work, which these tests exist to exercise. Making that case safe is a scheduling decision — not running two gates at once, or bounding `GOMAXPROCS` — and it belongs outside a test-speed repair.

### Final canonical repository gate

- **Command:** `["bash","_dev/tests/maintainer-verify.sh"]`, direct and unpiped, run at the merge commit `50569e88` in a detached worktree
- **Direct exit status:** 0 — no retry needed, "Maintainer verification passed."
- **do-work-cli lane:** wall 55s, 758 tests, slowest file `internal/finalization/finalization_recovery_test.go` 23.76s against the 30s limit (the same lane measured 67s wall and 28.02s slowest at this REQ's pre-build baseline)
- **queue-kanban lane:** wall 18s, 384 tests, slowest file `generate_test.go` 8.66s
- **Failures:** none

## Qualification

Passed. `qualify --request-path <this REQ> --diff-range 982e94f0..50569e88` returned success; its one warning is `QUALIFY-UNIFY-DISARMED` (no `[UNIFY]` box), which reflects that the builder and orchestrator roles were played inline in one session rather than handed back across a subagent boundary. `scope-drift --request-path <this REQ>` returned success: the corrected `## Scope` list and the `## Implementation Summary` list are the same four paths.

## Review

**Independent review of `982e94f0..50569e88` — Pass, no findings.** Performed by a separate reviewer session with no part in writing the change, against the stated rule: make the tests genuinely cheaper, weaken no assertion, split no file.

What it verified rather than assumed:

- The direct invocation in `runRetainedInventory` is byte-for-byte the argv and environment `uncommitted-inventory.sh` → `do-work-cli.sh` would have exec'd, checked against both shell scripts, and `moduleDirectory` resolves to the same `module_dir` the shell computes.
- Shim coverage did not shrink. `TestInventoryStagedAdditionDeletedFromWorktreeIsDeletion` was already the only nil-status caller before this change and still runs the full two-shell chain, so the launcher contract — Go version gate, `go tool -n` probe, exec — keeps an end-to-end test.
- One risk chased to ground and ruled out: `resolveRetainedInventoryExecutable` runs `go tool -n` while a synthetic case's fake-git PATH stub is active, and on a cold `GOCACHE` that call compiles and links. Had Go's VCS stamping shelled out to `git` during the link it would have hit the stub, which errors on any argv but two. Reproduced with a cold cache and an instrumented fake git: no git subprocess is invoked.
- `git init --template=` omits `.git/hooks/*.sample`, `.git/description` and `.git/info/exclude` against a plain `git init`, confirmed by diffing both. No test in either package reads any of them, and `.git/hooks` itself is recreated in both `TestMain`s.
- `os.CopyFS` creates every directory including empty ones, so the recreated `.git/hooks` survives the copy; files land at `0666` plus the source execute bits and directories at `0777`, permissive enough that `RemoveAll` teardown cannot be blocked by read-only git objects.
- One `TestMain` per package; templates are built before `m.Run()` and only read afterward, and copies write to distinct `t.TempDir()` destinations, so nothing shares mutable state.
- Re-ran `go test` over the three packages and `go vet`: all pass, clean.

No assertion weakening, no fixture drift, no concurrency hazard.
