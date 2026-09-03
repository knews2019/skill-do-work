---
id: REQ-543
title: '[impact-critical] Reap the commit hook with its own parent'
status: claimed
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
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process.go
related: [REQ-457]
claimed_at: 2026-09-03T13:36:18Z
route: B
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
