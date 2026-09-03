---
id: REQ-543
title: '[impact-critical] Reap the commit hook with its own parent'
status: completed
created_at: 2026-09-02T23:58:00Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
write_set:
  - skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group.go
  - skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix.go
  - skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unsupported.go
  - skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix_test.go
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_cancellation_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_unix.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_windows.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
related: [REQ-457]
claimed_at: 2026-09-03T13:36:18Z
route: B
review_at: 2026-09-03T15:42:47Z
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-03T13:36:30Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - async lifecycle behavior
    - cross-route regression gates
completed_at: 2026-09-03T16:35:08Z
commit: 1cc3beb
release_at: 2026-09-03T16:35:08Z
---

# Kill the Owned Commit Process Group on Cancellation

## What

Make cancellation of a media Git transaction terminate the whole owned process group, not just the direct `git` child. A `pre-commit` hook that ignores `SIGTERM` currently outlives the cancelled transaction and keeps running after the command has returned.

## Instances

- `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback` fails with `media commit hook survived cancellation` (`report_image_process_test.go:85`). The transaction rolls back correctly and returns within the deadline; the orphaned hook process is the only unmet assertion.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-do-work-cli.md` (including Package direction), `lessons-do-work-cli.md` (the `single-exit-owner` and `interruptible-blocking-io` entries), `CLAUDE.md`, `crew-members/communication-style.md`, and `_dev/lessons/validated-runtime-boundaries.md` (found by grep; it supplied D-07). **Did not read `general.md`, `coding-guardrails.md`, `testing.md` or `debugging.md` — the builder reported them absent from the tree, which is false; all 18 crew-member files are present and tracked.** See the Qualification note; the orchestrator checked the diff against those guardrails instead.
- [x] **[APPLY]:** Four new files in a new `internal/ownedprocess/` package, three modified, two deleted (their contents moved into the shared package), plus one prime edit the REQ's own Traps section required. Scope drift recorded in Qualification rather than absorbed.
- [x] **[UNIFY]:** Orchestrator independently re-ran `go build ./...`, `go vet ./...` (clean), `gofmt -l .` (silent), `GOOS=windows GOARCH=amd64 go build ./...` (OK), `GOOS=plan9 go build ./internal/ownedprocess/` (OK), and `go test -count=1 ./...` → **exit 0, zero failures**. Read the new package (551 lines across four files), the `git_transaction.go` diff including the D-05 guard, and the prime edit. No debug artifacts. Removed one untracked probe artifact (`do-work/working/baseline.json`) the builder flagged and correctly declined to touch.

## Finding Provenance

- **Discovered by:** REQ-457's Step 5.75 pre-flight, as one of two pre-existing failures in the `bash _dev/tests/maintainer-verify.sh` baseline at `008f3d3`. Unrelated to REQ-457's created-path ownership invariant, so it is captured here rather than folded into that REQ.
- **Evidence:** the test writes a `pre-commit` hook that traps and ignores `TERM` then sleeps 30s, cancels the context once the hook is running, and asserts `syscall.Kill(pid, 0)` fails afterwards.

## Diagnosis

Measured by process-table sampling at 20 ms during a live failing run, not inferred. Three plausible causes are ruled out: `Setpgid` **is** applied (`configureCancellableProcessGroup`, `internal/gittransaction/git_transaction.go:1344`), the negative-PID group signal **is** delivered (`pidfd_open(-pid)` returns `EINVAL`, so Go falls back to `syscall.Kill(-pid, sig)`), and the hook's `trap '' TERM` sets `SIG_IGN`, which is inherited across fork and exec, so `sleep` is TERM-proof too.

The real gap is ordering plus reaping:

- **The escalation is detached and nothing waits for it.** `Cancel` returns straight after the group SIGTERM; the SIGKILL lives in an orphan goroutine (`git_transaction.go:1362`) with a 1 s timer. `runGit` has no path that proves the group is dead before returning. `runOwnedProcess` (`internal/toolboxcommands/report_image_process.go:62-64`) already has exactly that invariant — `runGit` is the one that lacks it.
- **Group-SIGTERM kills the leader first, orphaning the hook; the later SIGKILL then leaves an unreaped zombie.** `git` dies immediately, the hook is reparented to PID 1, and `Cmd.Wait` stays blocked because the hook still holds the inherited stderr pipe. Measured: the transaction returns ~20 ms after the SIGKILL, and `syscall.Kill(pid, 0)` **succeeds on a zombie**. Orphan reap latency in this environment is ~1.5-2.0 s against the test's 3 s budget.

**This is why merely making the escalation synchronous is not sufficient.** A runner that TERMs the group, waits out grace, SIGKILLs the group, then polls `ownedProcessGroupAlive` — which deliberately treats `Z` as not-alive — would return with the hook a zombie and the assertion would still read alive.

Second launch path in the same file: `indexIsEmpty` (`git_transaction.go:860`) uses a bare `exec.CommandContext` with no process group, bypassing the seam the REQ constrains to one.

## Required Shape

1. Signal **descendants only** first (every group member whose pid is not the leader), so the leader survives to reap them. `git` `waitpid()`s its own hook, which removes the orphan and zombie window entirely — and `git commit` then exits non-zero on its own, which the existing `FailureCommit` path already turns into a rollback.
2. Escalate to SIGKILL on the descendants after a grace window.
3. Only then, if the leader is still alive, signal and escalate on the leader.
4. **Block until the group is gone before returning.** No detached goroutine.
5. Keep `command.WaitDelay` (currently 2 s, `git_transaction.go:1352`) as a last-resort backstop only. The grace budget must stay well under it, or it must be raised in tandem, or `WaitDelay` will fire first and re-create a return that outruns the kill.

## Traps

- **Windows asymmetry cuts the opposite way from `toolboxcommands`.** `configureOwnedProcess` (`report_image_process_windows.go:12`) fails closed. `runGit` must **not** — every Git transaction in the module goes through it, and `configureCancellableProcessGroup` deliberately probes `Setpgid` by reflection and degrades to default cancellation where it is absent. One shared runner needs both behaviours behind one API. `internal/gittransaction` has no build-tagged files today and cross-compiles clean for `GOOS=windows`; introducing a split there is new surface and needs a matching `GOOS=windows` compile line in the prime's Verify section.
- **Three divergent group runners already exist** with three grace policies: this inline `Cancel` closure (1 s), `report_image_process.go` (1 s), and `internal/nextselection/blocked_probe_unix.go` (500 ms + 500 ms). The last one also forwards the received signal and returns `128+signal` statuses — a contract a shared runner must not flatten. Consolidating the seam and sweeping the per-caller status and fail-closed contracts belong in the same change.
- **`prctl(PR_SET_CHILD_SUBREAPER)` is not available.** It lives in `x/sys/unix`, and the prime requires standard-library-only dependencies. The Go floor is 1.25.0, enforced by `_dev/tests/do-work-cli-go125-compatibility.sh`. This is a second reason the descendants-first shape is the right one.
- `reportImageGracePeriod` (`report_image_process.go:16`) is a test seam reassigned by `report_image_process_test.go:22` and `:104`; a refactor must preserve it.

**Regression net:** `report_image_process_test.go:19,92,101` (three currently-passing owned-process tests), `git_transaction_test.go:928-1001,1146` (three `pre-commit` hook tests, all exercising a *failing* hook rather than a cancelled one), `last30days_test.go:135`, `suiteinstall/{install,update}_transaction_test.go`.

## Detailed Requirements

- Cancelling a Git-committing transaction must terminate every process the transaction launched, including grandchildren a hook spawns, not only the direct child.
- Escalate past a process that ignores the graceful signal; a hook trapping `SIGTERM` must not survive.
- Preserve the existing return-within-deadline behavior and the existing rollback-to-preimage behavior on cancellation.
- Do not leave a killed process group's exit status misreported as success.

## Constraints

- Keep the owned-process-group runner the single seam for launched subprocesses; do not add a second ad-hoc kill path.
- Preserve exact typed result and rollback contracts in `prime-do-work-cli.md`.

## Dependencies

No request prerequisite.

## Red-Green Proof

**RED prompt/case:** `go test ./internal/toolboxcommands/ -run TestRemediationCancellationReachesMediaGitCommitAndRollback` — currently fails with `media commit hook survived cancellation`.
**Why RED now:** cancellation signals only the direct `git` child, so a hook that ignores `SIGTERM` keeps running as an orphan.
**GREEN when:** that test passes without modification, the transaction still returns inside its deadline, and the target still rolls back to its preimage bytes.

---
*Source: REQ-457 pre-flight baseline, captured during the work run.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The failure is one assertion in one test, but the diagnosis already recorded on this REQ shows the fix spans a shared subprocess seam with three divergent implementations and a Windows asymmetry that cuts the opposite way for `runGit` than for the image backends. Where the shared runner may live is a package-direction question, not a mechanical one.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation. The `## Diagnosis` and `## Required Shape` sections above are the exploration, measured by process-table sampling during a live failing run before this REQ was claimed.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify) — replace the detached-goroutine escalation in `configureCancellableProcessGroup` with a blocking descendants-first termination that returns only once the group is gone
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify) — a cancellation lock-in for a TERM-deaf hook, and the `indexIsEmpty` seam
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process.go` (modify) — source of the termination half being shared, keeping `reportImageGracePeriod` as its test seam
- a new shared package under `skills/do-work/tools/do-work-cli/internal/` (new) — **only if** the CLI prime's **Package direction** rule permits it; if it does not, say so and keep the seam where the rule allows

**Files I will NOT touch:** `internal/nextselection/blocked_probe_unix.go`'s `128+signal` status contract must not be flattened if that runner is consolidated; `report_image_process_windows.go`'s fail-closed behavior for image backends must stay fail-closed.

**Acceptance criteria (restated from REQ):**
- [x] Cancelling a Git-committing transaction terminates every process it launched, including a hook's grandchildren, not only the direct child
- [x] A hook that ignores the graceful signal does not survive — escalation is not optional
- [x] Return-within-deadline and rollback-to-preimage behavior on cancellation are both preserved
- [x] A killed process group's exit status is not misreported as success
- [x] The owned-process-group runner stays the single seam for launched subprocesses; no second ad-hoc kill path is added
- [x] `TestRemediationCancellationReachesMediaGitCommitAndRollback` passes without modification

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group.go` (new)
- `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix.go` (new)
- `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unsupported.go` (new)
- `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_cancellation_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_unix.go` (deleted)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_windows.go` (deleted)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**What was done:** A new `internal/ownedprocess` package holds the whole owned-process-group seam — `ConfigureGroup`, `TerminateGroup`, `DefaultGracePeriod` — with the platform split in one place. Teardown is **descendants-first, levelled from a single snapshot**: every group member that is not the leader is signalled bottom-up, escalated to SIGKILL after the grace window, and only then is the leader signalled if it is still alive. That ordering is what removes the zombie entirely — `git` survives long enough to `waitpid()` its own hook, so nothing is orphaned to init and `kill(pid, 0)` has nothing stale to report. `TerminateGroup` blocks until the group is gone. `configureCancellableProcessGroup`'s detached escalation goroutine and its `Setpgid` reflection probe are deleted; `indexIsEmpty` is folded onto the same seam; and `toolboxcommands`' two build-tagged files are absorbed into the shared package.

## Decisions

- **D-01** — Shared package at `internal/ownedprocess/`. **DECIDE & STATE.** The prime's Package direction rule constrains import *edges*, not new packages; the new edges (`gittransaction → ownedprocess`, `toolboxcommands → ownedprocess`) are acyclic, stdlib-only, and register no command. Recorded in the prime in the same change. **Value:** the build-tagged platform split lives in one package, so `gittransaction` keeps zero build-tagged production files and the `Setpgid` reflection probe is deleted. **Risk:** a third import edge to sweep if the API changes.
- **D-02** — One API, two behaviours via a returned boolean. **DECIDE & STATE.** `ConfigureGroup` reports whether ownership was established and refuses to decide what that means: `toolboxcommands` fails closed on `false`, `runGit` returns early and leaves `os/exec`'s default cancellation. **Value:** the asymmetry is stated once at each call site with its reason, rather than duplicated as two runners. **Risk:** a future caller could ignore the boolean; the doc comment names the obligation.
- **D-03** — `nextselection`'s blocked probe deliberately **not** consolidated. **DECIDE & STATE.** Its runner couples signal forwarding, `128+signal` status derivation, and a post-*success* orphan cleanup to its teardown, and its `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits` case has no live leader — so the descendants-first reaping property does not even apply there. Consolidating would need a mode flag used by exactly one caller. Two runners remain, down from three, and the REQ's protected contract is untouched.
- **D-04** — Snapshot-levelled sweep, not "signal whatever is childless now". **DECIDE & STATE.** The first implementation re-derived terminable leaves each round and **hung the `generate-report-image` shell probes forever**, because a backend that respawns a helper in a loop always has a live child, so the sweep never climbs to the parent. Levels now come from one snapshot, bounded by construction with no deadline. **Risk:** a child forked during the sweep is orphaned; a second pass ends it, and orphans are proved not-running rather than reaped, since only init can reap them.
- **D-05** — A HEAD-advanced guard added to `ExecuteTransaction`. **ESCALATE.** Descendants-first lets `git` reach its own error path — which also means it can reach a *successful* ref update when the hook happens to exit zero. `ExecuteTransaction` had no such guard (`ExecuteExactCommit` already did), so it rolled the worktree back from an advanced HEAD and reported `rolled_back` with the committed bytes still in place. **Value:** satisfies "do not misreport a killed group's status" in the direction that actually bites, and closes a latent misreport this REQ's own fix would otherwise have made reachable. **Risk:** one extra `rev-parse` per committing transaction.
- **D-06** — The leader gets the grace window to exit on its own before being signalled. **DECIDE & STATE.** After its descendants die, `git` releases `index.lock` and reports the hook failure, so the rollback does not trip over a stale lock. **Risk:** widens the window in which the commit can land, which is exactly what D-05 closes.
- **D-07** — An isolation boundary check before any group signal. **DECIDE & STATE.** `Getpgid(leader) != leader` falls back to bare-pid escalation, per `_dev/lessons/validated-runtime-boundaries.md` and REQ-220. Unreachable from inside the module, so no test reddens it.
- **D-08** — No `GOMAXPROCS=1` pin, despite REQ-525's precedent. **DECIDE & STATE.** That pin sharpens a goroutine race to `os.Exit`; here the ordering is enforced by process death and `Cmd.Wait`, not the Go scheduler. Every neuter below is already 100% deterministic unpinned, and `GOMAXPROCS=1 go test -count=5` on the new tests is green — the pin would add cost and no determinism. Recorded because declining a just-established precedent deserves a reason.

## Discovered Tasks

- `nextselection/blocked_probe_unix.go` still TERMs its whole group leader-first, so `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits` depends on init's reap latency inside a 2s window. It passes, but it is the same class this REQ fixed, and could move onto `ownedprocess` if the `128+signal` contract is preserved behind an explicit initial-signal parameter.
- `ownedprocess.terminateWholeGroup` — the fallback for a Unix host where `ps` cannot be executed — has no test; no cheap way to simulate a missing `ps`.
- `exact_commit.go`'s HEAD-advanced guard fails the whole transaction if the *pre*-commit `rev-parse` errors, so it cannot handle an unborn branch. `currentHeadDespiteCancellation` handles that case and could replace it.
- **A builder reported four always-load crew-member files as absent from the tree when all 18 are present and tracked.** Worth understanding: if the check was a relative `ls` run from the Go module directory, every builder dispatched into a subdirectory could silently skip the always-on guardrails. That is a pipeline-reliability question, not a code one.
- **impact-user-visible — `finalize` refuses any REQ that has no `## In Progress (interrupted)` entry in `do-work/CHECKPOINT.md`.** Found while closing this REQ, measured at the seam rather than inferred. `planTargets` adds the checkpoint to `plan.TargetPaths` **only when its bytes would change** (`internal/requeststate/state_plan.go:398`), while `PlannedPostimages` emits a checkpoint postimage **unconditionally** whenever `plan.CheckpointPath != ""` (`internal/requeststate/state_apply.go:156-158`). `finalization` takes its lifecycle *pre*images from `TargetPaths` and its *post*images from `PlannedPostimages`, so for a REQ with nothing to remove the two sets differ by exactly one path and `imageSetState` (`internal/finalization/finalization_apply.go:296`) refuses with `FINALIZATION-LIFECYCLE-CONFLICT: journal image sets have different path counts`. Verified from the refused journal: 3 preimage paths against 4 postimage paths, the extra one being `do-work/CHECKPOINT.md`. The refusal is safe (phase `prepared`, no commit, no release mutation applied) but it is unrecoverable in place — `recover-finalization` replays the same journal into the same refusal, and there is no discard verb, so the Git-private journal has to be removed by hand before a corrected manifest can be prepared. **Fold-first scan run, not assumed:** the six `sweep: true` pending REQs were checked; `REQ-502` (`sweep_key: checkpoint-section-blind-line-editing`) is the closest and is a **different** root cause — it is about the cleanup mover removing only an entry's header line, not about pre/postimage path symmetry in the finalize planner — so this is not a fold, and it needs its own REQ through the `pending-answers` consent flow rather than being minted silently here. Worked around for this closure by giving REQ-543 the checkpoint entry the Session Checkpoint template says it should already have had (it is claimed by this checkout and sitting in `do-work/working/`), inserted so that the planner's own removal reproduces the current bytes exactly and the committed checkpoint is byte-unchanged.

## Qualification

**Passed, with scope drift that is mostly mine** — 10 files verified, 6 requirements traced, P-A-U confirmed.

Mechanical: `tools/checks/qualify.sh` → `OK: mechanical qualification passed`, with five `WARN: (new) file has no static reference anywhere` lines. Those are the expected exception: a new Go package is referenced by *import path*, not by filename, and the two `_test.go` files are test files — both named exceptions in the check's own rules.

`tools/checks/scope-drift.sh` → exit 1 in both directions. Assessed rather than waved through:

**Touched but not declared — three causes, only one of them the builder's:**
- The four `internal/ownedprocess/` files and `git_transaction_cancellation_test.go`. **My Scope authoring.** I wrote "a new shared package under `.../internal/` (new)" without literal paths, and the check compares literal paths. The package was explicitly permitted, so this is an undeclared *file list*, not undeclared work.
- `report_image_process_unix.go` and `report_image_process_windows.go`, deleted. A direct consequence of moving their contents into the shared package I permitted. Reasonable, and it should have been declared.
- `prime-do-work-cli.md`. **REQ-mandated, correctly handled.** This REQ's own `## Traps` says "the prime's Verify section will need a matching `GOOS=windows` compile line for whatever package the runner lands in." `actions/work.md`'s builder rule covers exactly this: when the REQ's own requirements require a file class the Scope declaration contradicts, flag it and proceed with the required class. The builder flagged it.

**Declared but never touched — both parser artifacts of my own prose, for the second time:**
- `git_transaction_test.go`: the new tests went into a new build-tagged `git_transaction_cancellation_test.go` instead. Equivalent, arguably better — `//go:build unix` on a whole file rather than guards inside a shared one.
- `configureCancellableProcessGroup`, `indexIsEmpty`, `reportImageGracePeriod`: backticked identifiers in my Scope bullets' trailing descriptions, which `scope-drift.sh` reads as declared paths. **This is the second time I have done this** — REQ-457 hit it and recorded it as a discovered task. Recording it again is not the fix; the fix is that a Scope bullet's description must carry no backticks, and the discovered task on the script stands.

Independent (orchestrator-run, not the builder's report):
- `go build ./... && go vet ./...` clean; `gofmt -l .` silent.
- `GOOS=windows GOARCH=amd64 go build ./...` → OK. `GOOS=plan9 GOARCH=amd64 go build ./internal/ownedprocess/` → OK. Both matter: the builder's own first attempt named the fallback `*_windows.go`, whose implicit constraint left every other non-unix target with **no implementation** (`undefined: configureOwnedGroup`). It caught that itself by cross-building for plan9 and renamed the file to `_unsupported.go`.
- **`go test -count=1 ./...` → exit 0, zero failures.** Verified here, not taken from the report. This is the first fully green full-module run of the session.
- Read the new package (551 lines across four files), the `git_transaction.go` diff including D-05's guard, and the five-line prime edit.
- **The builder did not read four always-load crew-member files**, reporting them absent when all 18 are present and tracked (`git ls-files skills/do-work/crew-members/ | wc -l` → 18). I therefore read the diff against those guardrails myself: the change *reduces* surface (a detached goroutine, a reflection probe and two build-tagged files deleted; three runners down to two), every added comment states an invariant, and names are greppable and multi-word — `ConfigureGroup`, `TerminateGroup`, `currentHeadDespiteCancellation`, `terminateWholeGroup`. It is consistent with those rules despite not having been read against them, and the underlying reliability question is filed as a discovered task.

**Two guards the builder declared unearned, which I am recording rather than quietly accepting as coverage:** the zombie-tree predicate (`requireReaped=false` reddens nothing, because with correct ordering the parent reaps within one `ps` round-trip) and blocking-rather-than-detaching (`Cmd.Wait` cannot return before the leader is reaped, and the leader cannot exit before reaping its hook, so no case distinguishes them). Both are kept as defence-in-depth against a recycled pid; neither is claimed as tested.

## Testing

**Tests run** (from `skills/do-work/tools/do-work-cli/`): `go test -count=1 ./...`; `go test -count=5 ./internal/{ownedprocess,gittransaction,toolboxcommands,nextselection,suiteinstall}`; `go build ./...`; `go vet ./...`; `gofmt -l .`; `GOOS=windows GOARCH=amd64 go build ./...` plus `go vet` on the three packages; `GOOS=plan9 GOARCH=amd64 go build ./internal/ownedprocess/`.
**Result:** ✓ full module **exit 0, 26 packages ok**, working tree unchanged by the run. `-count=5` on the five affected packages exit 0. `go vet` exit 0, `gofmt -l .` silent, both cross-builds OK. The full-module run was repeated once more when this record was closed: exit 0, 26 packages ok.

**Provenance of this section, stated because it matters for how much it is worth:** the builder's hand-back never appended a Testing section, so this one is assembled from the orchestrator's own runs (recorded in `## Qualification`) and the independent review's re-verification at `1cc3beb`. Nothing below is the builder's self-report.

**Red-green validation** — traces `## Red-Green Proof`:
- `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback`: ✗ `media commit hook survived cancellation` before → ✓ after, **unmodified** (the test file is absent from the implementation diff, which the review checked rather than assumed).

**Neuter table** (each guard tested by breaking it, in a scratch copy; the working tree was never modified):

| Neuter | Result |
|--------|--------|
| Teardown reverted to leader-first whole-group | **reds 4 tests across 3 packages**, including this REQ's original RED |
| D-05's HEAD-advanced guard removed | **reds** `TestCancelledCommitThatLandsReportsCommittedRisk` — `outcome = "rolled_back", want "committed_state_risk"` |
| Leader's SIGTERM dropped from the escalation | **reds** `TestTerminateGroupLetsTheGracefulSignalRunFirst` |
| Every grace budget zeroed (member loop and leader loop) | **reds nothing** at `-count=3` in all three packages → finding F3 |
| Escalation detached again (answers preserved) | **reds nothing** → finding F12 |
| Zombie predicate neutered (`requireReaped=false`) | **reds nothing** → finding F12 |

D-05's branch was measured as well as neutered: 15/15 runs took the committed-risk branch, so the guard is exercised deterministically even though its test tolerates both outcomes.

**D-04 boundedness, checked rather than argued:** the level walk is bounded by construction — each pid is `placed` once, the loop exits on an empty level, orphans are handled outside the walk. A no-sleep respawn loop finished in 711 ms leaving 14 zombies and zero running survivors. No looping input could be constructed; the sweep misses a process only when it leaves the group (`setsid`) or under F1.

**New tests added:**
- `internal/ownedprocess/owned_process_group_unix_test.go` — `TestTerminateGroupLetsTheGracefulSignalRunFirst`, `TestTerminateGroupReportsAnAlreadyFinishedGroup`, `TestTerminateGroupEndsParentsThatKeepForkingChildren`.
- `internal/gittransaction/git_transaction_cancellation_test.go` — `TestCancelledCommitKillsTermDeafHookBeforeReturning`, `TestCancelledCommitThatLandsReportsCommittedRisk`.

**Existing tests updated:** none. No assertion anywhere was weakened, which is what let the original RED stand as proof.

**Pre-existing red, not caused here:** `internal/nextselection` → `TestBlockedProbeTimeoutKillsDescendantGroup` fails under load. Answered in full in `## Review` below: latent flake, unrelated to this REQ. Cross-compilation stays red for plan9, js, wasip1, aix and solaris; only aix and solaris changed, which is finding F8.

*Verified by work action*

## Review

**Overall: 83%** | 2026-09-03T15:42:47Z | **Verdict: Approve. Route B.** Re-verified against merged HEAD at implementation commit `1cc3beb`; the module is unchanged since.

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 82% |
| Test Adequacy | 70% |
| Scope | 88% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- **F1** — a reaped leader is treated as an empty group, so a live group is never swept. `internal/ownedprocess/owned_process_group_unix.go:50`: `syscall.Getpgid(leaderPID)` returning `ESRCH` is routed into the D-07 "not isolated" branch, then `escalateOnProcess` sees the leader gone and returns `os.ErrProcessDone`. Measured at the seam: returns in 2.2 us while `kill(-pgid, 0)` still succeeds and a TERM-deaf descendant stays in state `R`. Left 1 live worker behind in 16 teardown-branch runs — `impact-rule-change` → REQ-546 created (sweep `owned-group-teardown-contract-gaps`)
- **F2** — the exported contract over-claims and the prime restates it. `internal/ownedprocess/owned_process_group.go:30-34` ("blocks until they are gone, so a caller that returns afterwards has proved the group is dead") and `skills/do-work/tools/do-work-cli/prime-do-work-cli.md:30` ("`TerminateGroup` blocks until the group is gone"). Three measured counterexamples: an orphan left as `Z` and visible to `kill(pid, 0)`, a `setsid` escapee never signalled, and F1's reaped-leader no-op. The internal `escalateOnMembers` comment states the orphan distinction correctly, so only the two statements a caller and a future builder actually read are wrong — `impact-rule-change` → REQ-546
- **F3** — the one test that claims to pin the escalation grace does not pin it. `internal/ownedprocess/owned_process_group_unix_test.go:16-21` says the pause "is the only thing that makes the SIGTERM mean anything ... This pins the window". Zeroing every grace budget leaves that test and both caller packages green at `-count=3`; it reds only when SIGTERM is dropped entirely. What gives the handler its window is the `ps` fork inside `leaderRuns()` between the two signals. The test pins TERM-before-KILL, not the budget — `impact-rule-change` → REQ-546
- **F4** — `_dev/lessons/validated-runtime-boundaries.md:7` states the superseded shape: "signal that group on timeout, escalate when needed, and reap its leader. If isolation cannot be proved, fail closed". The canonical seam now signals descendants level-by-level rather than the group as a unit, and `runGit` deliberately degrades instead of failing closed (D-02). D-07 cites that file as the authority, so a future builder would read the stale text as the rule. Prose-only — `impact-rule-change` → appended to `do-work/prose-backlog.md`

**Minor findings:** 8 (F5-F12) plus 2 nits, report only — with three exceptions handled at hand-back: F9 (missing lesson-family bullet) and F10 (`write_set` listed 2 files while 10 were touched) are bookkeeping this closure owes and does here, and F12 (blocking and the zombie predicate are both unearned — neutering either reds nothing) is folded into REQ-546 as an instance.

**Acceptance:** Pass — full module exit 0, `-count=5` exit 0, and the neuter table confirms the fix rather than the tests merely passing.
**Suggested testing:** 5 items — earn F12's two guards; a seam lock-in for F1; give `TestRemediationLeaderExitStillKillsTermDeafDescendant` the descendant assertion its name promises (F11); a cross-compile gate lane, since the prime's `GOOS=windows` Verify lines are run by no gate; real hook shapes (`husky`/`lint-staged`, a daemonizing hook, macOS BSD `ps`).
**Follow-ups created:** REQ-546 (sweep: the owned-group teardown contract is stated and tested more strongly than it is implemented); **sweeps appended to:** None — the fold-first scan found six `pending` `sweep: true` REQs in `do-work/queue/` and no root-cause match.

**Priority question that closure was waiting on — answered, and it does not change the verdict.** `internal/nextselection` → `TestBlockedProbeTimeoutKillsDescendantGroup` is a **pre-existing latent flake, neither caused nor meaningfully worsened by this REQ**. Four independent legs:

1. Main's REQ-534 did not fix it. REQ-534 (running blocked probes from the repository root and propagating interruptions) is capture-only — `do-work/queue/REQ-534-...`, `status: pending`, never implemented.
2. `internal/nextselection/` is **byte-identical from `546b7a3` through HEAD**: `git diff 546b7a3 HEAD -- .../internal/nextselection/` is empty, and so is the same diff from `8d0f994`. This REQ's pre-merge and post-merge code are identical too — `1cc3beb..HEAD` touches none of the three packages.
3. The test is timing-dependent by construction: it polls `kill(pid, 0)` on an orphan zombie against a **2 s budget** (`blocked_probe_test.go:35`), while measured orphan reap latency on this host is **1.04 s idle and 1.44 s worst case** under 10 churn loops — roughly 0.6 s of headroom.
4. It did not fail in **28 runs**: `-count=10` clean, `-count=10` under 6 CPU hogs, `-count=8` under 8 fork/`ps` churn loops, `-count=5` in the real tree.

So the only variable between the failing gate run and the passing one is load. It is correctly filed as the builder's own discovered task and needs no REQ of its own. It will red again on a slower or differently-inited host, which is what that discovered task is for.

**Scope limit of the fix, recorded so nobody reads the contract as wider than it is.** The zombie is gone only where the dying process's parent is a live group member. For the `git` shape it is genuinely gone. Two cases remain, and both are acknowledged rather than fixed:

- an **orphan** (a hook forks a grandchild, then exits) is left as a zombie that still answers `kill(pid, 0)` — it is proved not-running, not reaped, because only init can reap it;
- a **`setsid` escapee** leaves the group and is never signalled at all, so it survives running.

Both are inherent to group ownership without `PR_SET_CHILD_SUBREAPER`, which this REQ's own `## Traps` rules out (it lives in `x/sys/unix`, and the prime requires standard-library-only dependencies). F2 is exactly the gap between these three real outcomes and the two sentences that claim one.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Measuring before diagnosing. Process-table sampling at 20 ms during a live failing run refuted three plausible causes (`Setpgid` is applied, the negative-PID group signal is delivered, `trap '' TERM` is inherited by `sleep`) and produced the descendants-first shape directly. A guessed fix here would have been "make the escalation synchronous", which does not work.
- Testing every guard by breaking it. The neuter table is what turned "the tests pass" into "these four tests red when the fix is removed" — and it is the same instrument that found three guards nothing reds on (F3, F12). A green suite would have hidden both.
- Cross-building for plan9. The builder's first attempt named the non-Unix fallback `*_windows.go`, whose implicit constraint left every other non-Unix target with no implementation. One cross-build caught it before merge; `_unsupported.go` is the fix.

**What didn't:**
- The builder reported four always-load crew-member files as absent from the tree. All 18 are present and tracked. The guardrails were never read; the orchestrator read the diff against them afterwards instead. Filed as a discovered task, because if the check was a relative `ls` from the Go module directory then every builder dispatched into a subdirectory silently skips the always-on guardrails.
- The Scope declaration. Prose paths ("a new shared package under `.../internal/`") instead of literal ones, and backticked identifiers in the trailing descriptions, which `scope-drift.sh` reads as declared paths. **Second occurrence** — REQ-457 hit it first. The fix is that a Scope bullet's description carries no backticks; recording it a third time would not be a fix.
- `write_set` was never amended when the file list grew from 2 to 10. `## Scope` drift was reconciled in prose and the frontmatter guard was not, and under fan-out it is the guard, not the prose, that prevents collisions (F10). Corrected in this record's frontmatter at closure.
- The first sweep implementation re-derived terminable leaves each round and **hung the `generate-report-image` shell probes forever**: a backend that respawns a helper in a loop always has a live child, so the sweep never climbed to the parent. D-04's single-snapshot levelling is the replacement.

**Worth knowing:**
- **A zombie satisfies `kill(pid, 0)`.** That one fact is why making the escalation synchronous is not sufficient and why the ordering had to change: only letting `git` survive to `waitpid()` its own hook removes the window, because a wait on our own child ends in milliseconds while a wait on init is measured in seconds.
- `ESRCH` from `Getpgid(leader)` does **not** mean "the group is not isolated" — the leader can be reaped while the group is still alive and running. `internal/nextselection`'s `cleanupReapedProcessGroup` (`blocked_probe_unix.go:54`) already handles exactly this case and is the reference behaviour F1 lacks.
- Without `PR_SET_CHILD_SUBREAPER` — barred here, since it lives in `x/sys/unix` and this module is standard-library-only — an orphan can be proved not-running but never reaped, and a `setsid` escapee cannot be reached at all. Any future wording of this contract has to carry all three outcomes.
- The prime's `GOOS=windows` Verify lines are run by no gate: `grep -n GOOS _dev/tests/*.sh` is empty. The portability contract is written down and unenforced.
- The REQ's `## Traps` names `_dev/tests/do-work-cli-go125-compatibility.sh` as the Go 1.25 floor enforcer. Main deleted that script in `5e0e166` and the merge removed its prime Verify line, so the floor now rests on `go.mod`'s `go 1.25.0` alone. Not this REQ's doing, but the Traps text is stale for the next reader.
- New lesson family **`reaped-by-its-own-parent`**, promoted into the prime's `## Traps` (`prime-do-work-cli.md:56`) by the implementation and now carried by its satellite bullet in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as well. The gap the review recorded as finding F9 — the index row at `do-work/lessons-index.md:11` listing a family the satellite did not carry, which contradicts the index's own header rule ("`families` is the exact sorted set of `[family: <slug>]` markers present in lesson bullets") — is closed in the release commit below, together with the `tokens` recount the header rule prescribes (`ceil(bytes / 4)`: 5127 to 5660). The disagreement recorded at hand-back stands as written: the two halves belonged in one commit and were split only because that closure could not write under `skills/`.

## Orientation

Cancelling a Git-committing transaction now ends the whole process tree it launched, including the grandchildren a `pre-commit` hook spawns, and it blocks until they are gone instead of returning over a detached SIGKILL. A cancelled commit that nevertheless lands is now reported as `committed_state_risk` rather than as a rollback that silently left the bytes in place. Lives in the CLI's new `internal/ownedprocess` package, which both `gittransaction` and `toolboxcommands` now call.

`prime_files`: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — spot-checked at closure; every path it names still exists, and the builder's own edit (the package entry at `:30`, the `Package direction` edge at `:34`, the `reaped-by-its-own-parent` trap at `:56`, and the Windows Verify line at `:76`) is present and accurate except the one clause F2 names.

**[MAP CHANGED]** — owned-process-group launch and teardown is now a single named seam with the platform split in one package, down from three divergent runners to two (`nextselection`'s blocked probe keeps its own, per D-03). `gittransaction` holds zero build-tagged production files and the `Setpgid` reflection probe is gone. Why it matters: the teardown ordering — descendants before parents, level by level, from one snapshot — is now the module's single answer to "kill the tree", and the two outcomes it cannot deliver (an unreaped orphan, a `setsid` escapee) are properties of that one seam rather than of each caller.

## Closure

**Released as 0.270.1 — patch.** The implementation at `1cc3beb` changed shipped Go code under `skills/do-work/tools/do-work-cli/` and the prime beside it, which makes it a release under `_dev/primes/prime-releases.md`, and it landed without a version bump. This closure carries that release: `VERSION`, `skills/do-work/VERSION` and `skills/do-work/actions/version.md` go 0.270.0 to 0.270.1, and one entry lands in `CHANGELOG.md` and its byte-identical mirror `skills/do-work/CHANGELOG.md`. Patch and not minor because nothing user-invocable exists that did not exist before: this is a defect fix in an internal package plus a lessons bullet, no new command, flag, format or reversed default. The bump table's tie-breaker ("when genuinely torn between two levels, pick the smaller one") points the same way.

**Finding F9 is closed here** — `lessons-do-work-cli.md` gains the `reaped-by-its-own-parent` bullet the prime and the index already pointed at, and `do-work/lessons-index.md`'s token count for that satellite is recomputed from `wc -c` per the index's own header rule.

**This REQ does not close with a clean bill of health, and should not be read as one.** Its four Important review findings are filed, not fixed. They live on **REQ-546** (making the owned-group teardown contract match what it implements, `sweep_key: owned-group-teardown-contract-gaps`, `status: pending`), which carries all four as instances:

- **F1** — a reaped leader is read as an empty group, so a live group is never swept (1 live worker left behind in 16 teardown-branch runs).
- **F2** — the exported doc comment and the prime both claim one outcome where the mechanism has three.
- **F3** — the test that says it pins the escalation grace window pins TERM-before-KILL instead; zeroing every budget reds nothing.
- **F12** — blocking-rather-than-detaching and the zombie predicate are both kept and both unearned; neutering either reds nothing.

REQ-546 also gained instance **F13** from REQ-537's closure: the prescribed-shell interrupt fixture TERMs the wrapper pid only and never asserts the backend stopped, which is the shell-side counterpart of F11.

So what shipped is real and green — the module is `go test -count=1 ./...` exit 0 across 26 packages, re-run immediately before this release — while the contract's wording and three of its guards are known gaps with an owner.
