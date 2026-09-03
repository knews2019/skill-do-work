---
id: REQ-546
title: '[impact-rule-change] Make the owned-group teardown contract match what it implements'
status: pending
created_at: 2026-09-03T15:42:47Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
addendum_to: REQ-543
review_generated: true
related: [REQ-543, REQ-534]
sweep: true
sweep_key: owned-group-teardown-contract-gaps
write_set:
  - skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group.go
  - skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix.go
  - skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix_test.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_test.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
---

# Make the Owned-Group Teardown Contract Match What It Implements

## What

**Root cause: the owned-group teardown contract is stated and tested more strongly than it is implemented.** REQ-543 (reaping a cancelled commit's hook with its own parent instead of with init) built `internal/ownedprocess` as the module's single seam for killing a launched process tree. The seam works for the shape it was built for. Three separate statements about it are wider than the code, and three guards that the code keeps are pinned by no test at all — so the next person to refactor this seam is told the wrong invariant by the doc comment, the prime, and a test's own name, and nothing reds when they act on it.

This sweep closes the class rather than patching one spot: after it, the teardown has one stated contract that names all three of its real outcomes, the one implementation gap behind that contract is fixed, and every guard the seam keeps is held by a test that reds when the guard is removed. **Done means the class cannot recur** — a future edit that restores the old detached shape, drops the zombie predicate, or re-opens the reaped-leader hole fails a test instead of passing review.

## Instances

- **[ ] F1 — a reaped leader is treated as an empty group, so a live group is never swept.** `internal/ownedprocess/owned_process_group_unix.go:50`: `syscall.Getpgid(leaderPID)` returning `ESRCH` is routed into the D-07 "process group cannot be proved isolated" branch, and `escalateOnProcess` (`:255`) then sees the leader gone and returns `os.ErrProcessDone`. `ESRCH` on the leader does not mean the group is empty — the leader can be reaped while its descendants run. Measured at the seam: returns in 2.2 us while `kill(-pgid, 0)` still succeeds and a TERM-deaf descendant stays in state `R`; an old-style `kill(-pgid, TERM/KILL)` — exactly what the deleted `terminateOwnedProcess`/`killOwnedProcess` did — ends it. Reachable through `runOwnedProcess` (`internal/toolboxcommands/report_image_process.go:35`) when both select channels are ready: **1 live worker left behind in 16 teardown-branch runs**, while the code's own new comment claims "scratch cleanup cannot race a backend that is still writing". Not an aggregate regression at that caller (baseline leaves 26/40, new 25/40, dominated by the pre-existing `<-done` select branch), which is why it was not blocking. **The reference behaviour already exists in this module:** `internal/nextselection/blocked_probe_unix.go:54`'s `cleanupReapedProcessGroup` sweeps the group precisely when the leader has been reaped. (found by REQ-543 review / UR-085)
- **[ ] F2 — the exported contract over-claims, and the prime restates it.** `internal/ownedprocess/owned_process_group.go:30-34` says `TerminateGroup` "blocks until they are gone, so a caller that returns afterwards has proved the group is dead rather than merely signalled", and `skills/do-work/tools/do-work-cli/prime-do-work-cli.md:30` says "`TerminateGroup` blocks until the group is gone". There are **three** real outcomes, not one: the group is genuinely dead where the dying process's parent is a live group member (the `git` shape); an **orphan** is left as a zombie that still answers `kill(pid, 0)`, proved not-running rather than reaped, because only init can reap it; and a **`setsid` escapee** leaves the group and is never signalled, so it survives running. The internal `escalateOnMembers` comment (`owned_process_group_unix.go:118-124`) already states the orphan distinction correctly — only the two sentences a caller and a future builder actually read are wrong. Both limits are inherent to group ownership without `PR_SET_CHILD_SUBREAPER`, which REQ-543's own Traps rules out (it lives in `x/sys/unix`; this module is standard-library-only), so the fix is to state the contract accurately, not to widen the mechanism. (found by REQ-543 review / UR-085)
- **[ ] F3 — the one test that claims to pin the escalation grace does not pin it.** `internal/ownedprocess/owned_process_group_unix_test.go:16-21` says the pause between SIGTERM and SIGKILL "is the only thing that makes the SIGTERM mean anything ... This pins the window with a process that can only exit gracefully." Zeroing every grace budget in the implementation — the member loop and the leader loop both — leaves `TestTerminateGroupLetsTheGracefulSignalRunFirst` and both caller packages green at `-count=3`. It reds only when SIGTERM is dropped from the escalation entirely. What actually gives the signal handler its window is the `ps` fork inside `leaderRuns()` (`owned_process_group_unix.go:239`) between the two signals. The test pins TERM-before-KILL, which is worth pinning; either its comment tells the truth about what it pins and a second test pins the budget, or the budget stops being described as a pinned invariant. (found by REQ-543 review / UR-085)
- **[ ] F12 — two kept guards are unearned; neutering either reds nothing.** Setting the tree sweep's `requireReaped` to `false` (`owned_process_group_unix.go:124`, `anyMember`'s `includeZombies` at `:216`) leaves all three packages green at `-count=3`. Detaching the escalation entirely, answers preserved, also leaves all three green — because `runGit` blocks on `command.Run()` and the leader cannot exit before reaping its hook, so no current case distinguishes blocking from detaching. **Keep both and earn both cheaply**, as REQ-543's Qualification already recorded them as unearned rather than claiming them as coverage: blocking is earned by a five-line seam test asserting zero live members immediately on return (it reds under the detach neuter), and the zombie predicate by a unit test of `anyMember`'s `includeZombies` against one synthetic unreaped child. Removing them instead would restore the detached shape REQ-543 deleted. (found by REQ-543 review / UR-085)

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- Findings **F1, F2, F3** (Important) and **F12** (Minor) from the independent review of REQ-543 (reaping the commit hook with its own parent), verdict Approve at 83% with acceptance Pass. Every one was measured at the seam by the reviewer, not read off the diff.
- **Fold-first scan, re-verified at creation rather than taken from the review report:** `grep -rl "^sweep: true" do-work/queue/` returns six REQs, all `status: pending` — `REQ-475` (`memory-configured-tree-readers-not-rooted`), `REQ-496` (`already-green-repair-shared-validator-missing`), `REQ-502` (`checkpoint-section-blind-line-editing`), `REQ-512` (`legacy-finalization-semantic-ownership-incomplete`), `REQ-526` (`transaction-created-path-rollback-identity`), `REQ-544` (`answer-line-marker-position-spoofing`). None shares this root cause, and no non-sweep pending REQ in any UR does either. This is the justified-exception new file, not a missed fold.
- **F4 from the same review is not here.** It is prose-only (`_dev/lessons/validated-runtime-boundaries.md:7` states the superseded teardown shape) and went to `do-work/prose-backlog.md` per the Fold-First Rule's destination 3.

## Detailed Requirements

- A reaped leader must not be read as an empty group. When `Getpgid(leaderPID)` fails with `ESRCH`, the group may still be alive: sweep it. Keep the genuine not-isolated case (a leader whose `Getpgid` succeeds but does not equal `leaderPID`) on the bare-pid path, because `-leaderPID` could otherwise name the caller's own group — that boundary is REQ-543's D-07 and REQ-220's rule and must survive intact.
- `TerminateGroup`'s doc comment and the prime's `internal/ownedprocess/` entry must state the same contract, and it must name all three outcomes: dead where the parent is a live group member; an orphan proved not-running but left unreaped; a `setsid` escapee out of reach. One statement, restated once in the prime, with no third wording.
- Every guard the seam keeps has a test that reds when the guard is removed: the escalation grace budget, blocking-rather-than-detaching, and the zombie predicate. A guard that cannot be earned this way is deleted rather than kept as decoration.
- `TestRemediationLeaderExitStillKillsTermDeafDescendant` (`internal/toolboxcommands/report_image_process_test.go:101`) asserts only `result.Interrupted` and never that the descendant died — its name promises exactly the invariant F1 breaks, so it cannot red on it (review finding F11). Give it the assertion its name promises.
- Preserve the typed result and rollback contracts REQ-543 preserved, including D-05's `committed_state_risk` guard, and keep `reportImageGracePeriod` as a test seam.

## Constraints

- Standard library only. `PR_SET_CHILD_SUBREAPER` is not available and is not the fix for F2 — the contract changes to match the mechanism, not the reverse.
- Do not add a second kill path. `internal/ownedprocess` stays the single seam for `gittransaction` and `toolboxcommands`, and `internal/nextselection`'s runner stays separate per REQ-543's D-03 (its signal forwarding and `128+signal` status contract must not be flattened).
- Honour the prime's **Package direction** rule: `ownedprocess` imports only the standard library and never a command package.
- Editing `prime-do-work-cli.md` and any file under `skills/` makes this a shipped-file change, so it carries a release commit with a changelog entry and the mirror. Do not fold an unrelated release into it.

## Dependencies

No request prerequisite: REQ-543's implementation is already committed at `1cc3beb`, so this REQ edits code that exists.

**Sequence with REQ-534, do not fold into it.** REQ-534 (running blocked probes from the repository root and propagating interruptions) is `status: pending` over the same file family — `internal/nextselection/blocked_probe*` — with a **different root cause**: probe working directory, interruption propagation, and a signal-handler installation window. One fix would not close both, so they are two REQs. They touch each other in one place worth knowing: `blocked_probe_unix.go:54`'s `cleanupReapedProcessGroup` is the reference behaviour F1 lacks, and REQ-543's own discovered task proposes moving that probe onto the `ownedprocess` seam later. That move is blocked until F1 is fixed — it would otherwise regress `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits` — so this REQ runs first and neither REQ should attempt the migration.

## Red-Green Proof

**RED prompt/case:** At the seam, reap the leader and leave one live TERM-deaf member in its group, then call `TerminateGroup(leaderPID, grace)` and assert nothing in the group is running on return. Separately: zero both grace budgets and run `go test -count=3 ./internal/{ownedprocess,gittransaction,toolboxcommands}`; detach the escalation and run the same; set `requireReaped` to `false` and run the same.
**Why RED now:** the reaped-leader case returns `os.ErrProcessDone` in 2.2 us with a group member still in state `R` (1 live worker left behind in 16 teardown-branch runs through `runOwnedProcess`), and all three neuters leave every package green, so no test holds any of the three guards.
**GREEN when:** the reaped-leader case sweeps the group and the new seam test passes; each of the three neuters reds at least one named test; `TerminateGroup`'s doc comment and `prime-do-work-cli.md:30` state the same three-outcome contract; `TestRemediationLeaderExitStillKillsTermDeafDescendant` asserts the descendant died; and `go test -count=1 ./...` in `skills/do-work/tools/do-work-cli/` stays exit 0.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---
*Source: REQ-543 independent review, findings F1, F2, F3 and F12.*
