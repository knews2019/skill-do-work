---
id: REQ-663
title: 'Coordinated run prints a four-line preflight for disk, leftover processes, stale locks and run policy before the first dispatch'
status: cancelled
created_at: 2026-10-09T21:16:07Z
user_request: UR-146
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-662]
related: [REQ-662, REQ-664, REQ-665, REQ-666, REQ-667]
batch: run-coordinate-mode
completed_at: 2026-10-10T12:58:54Z
---
# Coordinated Run Prints a Four-Line Preflight Before the First Dispatch
## What
Before the first builder is spawned in a `--coordinate` run, print one line each for free disk on the repo root, leftover test or browser processes, stale locks whose owner process is dead, and whether `do-work/run-policy.md` was read. Stop on a critical disk reading. Report leftover processes and dead-owner locks and ask before removing them. `do-work/run-policy.md` is a new optional, committed file of plain bullets the coordinator copies into every brief.
## Why
Report item A2 (UR-146 input): a disk probe exists because a fan-out run once filled the disk, but it is a board finding, and nothing runs it before dispatch. Concurrent full suites and leftover processes caused load-only failures in the 2026-10-07 to 10-08 run. Run rules such as "no live-ops" are re-typed per session today.
## Verified Facts (checked at 0.305.87)
- `skills/do-work-board/tools/queue-kanban/verify.go:47` defines category `low-disk-space`; `:217-224` sets fixed thresholds, warning below 10 GiB and critical below 3 GiB free (strict "below"); `disk_space_test.go:253` pins that only the repo root is measured (the v0.305.86 narrowing the report names).
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:23` records why the disk probe exists.
- No run-start check exists in `actions/work.md` or `actions/fan-out-reference.md`.
## Detailed Requirements
1. Before the first spawn of a `--coordinate` run, print exactly these four lines:
   ```
   disk      <free> on repo root (verify low-disk-space level)
   procs     <n> leftover test/browser processes (names, pids) or none
   locks     <path> owner pid <pid> dead|alive, or none
   policy    do-work/run-policy.md read (<n> rules) or absent
   ```
2. Stop the run on a critical disk level. A warning level prints and continues.
3. Report leftover processes and dead-owner locks, and ask before removing them. The preflight itself never kills or deletes.
4. `do-work/run-policy.md` is optional and committed: plain bullets the coordinator copies into every builder and integrator brief. Illustrative content from the report: a skip list of REQs, "one heavy gate at a time", "no browser QA while `full-gate.lock` is held", "no live-ops (deploy, publish, production writes)".
5. Document the preflight and the policy file in `actions/fan-out-reference.md` (or `actions/work.md` where the run starts) and in `docs/work-guide.md`.
6. Release per `_dev/primes/prime-releases.md`.
## Constraints
- The disk line reuses the board's `low-disk-space` reading and thresholds; it does not copy the numbers into a second place (report, Request A2).
- Disk cleanup is out of scope: the preflight reports, it does not delete (report, Out of scope).
- Plain `--fan-out` runs are unchanged (REQ-662 acceptance).
- No new REQ status and no new frontmatter field.
## Assumptions (recorded at capture, no questions asked)
- The four checks are mechanical, so per the maintainer rule "programs beat prose" they are one do-work-cli command (name is the builder's choice) that the action calls, with Go tests. That is why `tdd: true`. If the builder finds a reason to keep it prose, it records the reason in Decisions and the REQ moves to `tdd: false` with the same Red-Green Proof.
- Reuse of the disk reading: the board probe lives in the separate `do-work-board` module. The builder either calls the board's verify output or moves the threshold reading to a place both can import, whichever the module boundary allows. If the board is not installed, the disk line says the reading is unknown and the run does not stop; an unknown reading never prints as healthy.
- "Leftover test or browser processes" is a condition, not a name list: processes whose working directory is the repo root, one of its worktrees, or below them. The builder may refine the condition and records it.
- "Stale locks" covers the `full-gate.lock` files from REQ-667 (the coordinated-run rules) in run directories under `do-work/runs/`. The suite has no other lock files today; the builder adds any it finds.
- The policy file is read once at run start. The `<n> rules` count is the number of top-level bullets.
## Dependencies
Depends on REQ-662 (the `--coordinate` mode the preflight runs under). The lock check reads REQ-667's lock format, which REQ-662 already depends on.
## Builder Guidance
High certainty on the four lines and the stop rule. Lower on how the board's disk reading is reused across modules and on the process condition. Latitude: command name, output padding, policy-file location under `do-work/`.
## Red-Green Proof
**RED prompt/case:** In a fixture repo with less than 3 GiB free (faked measurer) and a `full-gate.lock` naming a dead pid, start `do-work run --coordinate --fan-out 2`.
**Why RED now:** no preflight exists; the run dispatches builders without any disk, process, lock or policy line.
**GREEN when:** the four lines print before any spawn, the critical disk reading stops the run, and the dead-owner lock is reported and not removed. With healthy disk and no policy file, the run continues and the policy line says absent.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8334 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: families `disk-space-blind-spot` and `unknown-reads-as-clean` fit the disk line.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over budget; `slugged: partial`). Matching reason: a new do-work-cli command; family `destructive-next-argv` fits a check that must never kill or delete.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: process and pid probing in shipped commands.
## Full Context
See `do-work/user-requests/UR-146/input.md` for complete verbatim input (sections Request A2, Where the behaviour lives today, Proposed direction A2, Acceptance check, Out of scope). No queued candidate shares this root cause (the queue held REQ-654 to REQ-661, the ai-report and do-work-cli batches, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-coordinate-mode.md`, Request item A2: "Preflight before the first dispatch: free disk, leftover test or browser processes, stale locks whose owner process is dead, and an optional `do-work/run-policy.md` read at start."*

## Cancelled

- **When:** 2026-10-10T12:58:54Z
- **Why:** folded into the simplified recapture of 2026-10-10 (maintainer chose the aggressive simplification)
- **Decided by:** user, via `do-work abandon`
