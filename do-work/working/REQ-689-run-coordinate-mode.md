---
id: REQ-689
title: 'do-work run --coordinate: the mandatory coordinator shape, its run rules, a prose preflight, a stall loop, and a handoff that resumes coordinated'
status: claimed
created_at: 2026-10-10T13:05:57Z
user_request: UR-153
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-690]
batch: run-coordinate-mode
claimed_at: 2026-10-10T13:14:51Z
---
# do-work run --coordinate: Mandatory Coordinator Shape, Run Rules, Preflight, Stall Loop and Coordinated Handoff
## What
One prose change set across `skills/do-work/actions/work.md`, `actions/fan-out-reference.md`, `actions/restart-with-parallel-handoff.md`, `docs/work-guide.md` and `crew-members/background-agents.md`. It adds the run flag `do-work run --coordinate [--fan-out N] [REQ-NNN ...]` with a SKILL.md route, makes the coordinator shape mandatory under that flag, writes the five coordinated-run rules into `fan-out-reference.md`, adds a prose preflight checklist and an optional `do-work/run-policy.md`, adds a stall loop built on `do-work status --watch`, and makes a handoff written from a coordinated run resume with `--coordinate`. No new Go code. This REQ replaces the cancelled REQ-662 to REQ-667 (UR-146), which the maintainer folded on 2026-10-10.
## Why
About 13 sessions from 2026-09-27 to 2026-10-09 started with a pasted coordinator directive, 5 of them identical. Sessions that missed the paste integrated in the main session and grew to about 670k tokens before a handoff. In the 2026-10-07 to 10-08 run, background integrators stalled for 3 and 2 hours, concurrent full suites caused load-only failures, and after the user set five run rules each integration took 11 to 57 minutes. Those rules live only in one per-machine memory note. A handoff written from a coordinated run resumes uncoordinated today (`restart-with-parallel-handoff.md` writes `do-work run --fan-out N`).
## Why this is one REQ, and what was dropped
- REQ-662 (the flag), REQ-666 (handoff and teardown) and REQ-667 (the run rules) edited the same two prose files; three releases for one change set.
- REQ-663 (a four-line preflight as a Go command) and REQ-665 (a GO/NO-GO resume check as a Go command) both printed whether a lock's owner pid is dead. REQ-073 rules out a PID check in shipped code (`skills/do-work-board/tools/queue-kanban/verify.go`, comment near `staleClaimThreshold` use: "REQ-073 rules out a lock, heartbeat, PID check, mtime heuristic and time threshold alike"), and REQ-690 (do-work status) refuses it for the same reason. REQ-665 also parsed the prose handoff Reference, against the maintainer's no-grammar ruling, for three observed cases. The preflight survives here as a prose checklist; the resume drift check is dropped.
- REQ-664 (a stall check reading progress logs) overlapped REQ-690, which classifies quiet and stale rows. The stall loop here reads `do-work status --watch`.
## Verified Facts (checked at 0.305.101)
- `skills/do-work/SKILL.md` run route: `run`, `go`, `start`, `work`, `begin`, `process`, `execute`, `build`, `continue`, `resume`. No coordinator phrase routes anywhere.
- `actions/work.md` → `## Input` defines `--fan-out [N]`; unrecognized arguments are rejected after stripping `--wave N`, `--fan-out [N]` and `--skip-impact-negligible`; a usage line and a Step 0 checklist line restate the argument list. Under fan-out "integration stays serial" and "the dispatch mechanism stays unspecified".
- `actions/fan-out-reference.md` → `### Delegated integration — the coordinator shape`: "the orchestrator may hand each REQ's integration ... to one agent at a time"; entry by `advance REQ-NNN`, never `recover`; one writer under the project root at a time. `### Run directory, briefs and hand-backs` lists `do-work/runs/work-<stamp>/`, `REQ-NNN-brief.md`, `REQ-NNN-handback.md`, `REQ-NNN-integrate.md`, `manifest.md`. `grep -rn "progress.log\|full-gate.lock" skills/` finds nothing.
- `actions/work.md` teardown: the checkpoint, then cleanup, after the last integrator returns. Nothing stops idle agents.
- `actions/restart-with-parallel-handoff.md`: the build resume command is `do-work run --fan-out N`; the `phandoff` alias in `crew-members/communication-style.md` is the only trigger.
- `skills/do-work-board/tools/queue-kanban/verify.go` defines `low-disk-space` with fixed thresholds (warning below 10 GiB, critical below 3 GiB free on the repo root).
- `skills/do-work/crew-members/background-agents.md` says the pattern makes failures "survivable and recoverable", not prevented.
- REQ-639 (delegated integration, archived in UR-138, commit c40c6460) decided "There is no flag". This REQ reverses that at the maintainer's request; the optional shape stays for plain `--fan-out`.
## Detailed Requirements
A. The flag and route:
1. `actions/work.md` → `## Input` accepts `--coordinate`. It composes with `--fan-out [N]`, `--wave N` and targeting tokens; bare `--coordinate` implies `--fan-out`. Add it to the stripped tokens, the usage line and the Step 0 checklist line.
2. Add a SKILL.md route row above the plain run row: `drive the queue`, `use the main session as a coordinator`, `run --coordinate` → `./actions/work.md` with `--coordinate`. The bare word `coordinator` is a trigger only if routing tests show no collision; otherwise phrases only.
3. In `fan-out-reference.md` → Delegated integration add: "Under `--coordinate` this shape is mandatory. The main session dispatches, writes briefs in writing gaps, and reports. It never builds, never merges and never runs Step 6 to Step 9 itself. No integrator pushes. A blocking question becomes a `pending-answers` follow-up, never a wait."
4. Plain `do-work run --fan-out 3` behaves exactly as today.

B. The run rules, one subsection under Delegated integration, applying whenever the coordinator shape is used:
5. Heartbeat: every builder and integrator appends `<UTC> <phase> <command>` to `do-work/runs/<run>/REQ-NNN-progress.log` before and after any command expected to take over a minute. An agent never ends its turn while a command it started is still running.
6. Full-gate lock: only integrators run the full test suite, one at a time, holding `do-work/runs/<run>/full-gate.lock` with the owner label and UTC start inside. Builders run focused suites for the files they touched.
7. Red full gate: rerun the failed suites alone; failures that pass alone are load-only, one full rerun; a real regression is fixed, then one full gate.
8. Manifest row per integration: takeover-to-finalization minutes, full gates run, stall restarts (three new `manifest.md` columns).
9. Hand-back coverage: before writing an integrator brief, the coordinator checks that the hand-back covers every item forwarded to the builder mid-run.
10. Add `REQ-NNN-progress.log` and `full-gate.lock` rows to the run-directory table and the three columns to its `manifest.md` row. State in the subsection that the progress log is an append-only run-directory text file and the lock serializes one test command inside one run; neither is queue state, a liveness claim, or the machinery REQ-069 and REQ-073 deleted.

C. Preflight, prose, before the first spawn of a `--coordinate` run:
11. Four lines, each answered from an existing surface: disk (run the board's verify and quote its `low-disk-space` level; stop the run on critical, continue on warning, say "unknown" when the board is not installed and never print unknown as healthy); leftover test or browser processes under the repo root or its worktrees (report, ask before killing, never kill unasked); stale `full-gate.lock` files under `do-work/runs/` (report path and age, never delete); `do-work/run-policy.md` read with its bullet count, or absent.
12. `do-work/run-policy.md` is a new optional, committed file of plain bullets the coordinator copies into every builder and integrator brief (illustrative content: a skip list of REQs, "one heavy gate at a time", "no browser QA while `full-gate.lock` is held", "no live-ops").

D. Stall loop:
13. After the first dispatch, the coordinator arms one recurring check every 15 to 20 minutes on the harness's scheduler (CronCreate and `/loop` are examples, not the only routes); without a scheduler it checks at the start of each of its own turns and says so. Each tick runs `do-work status --watch` (REQ-690) once it exists, and reads each `REQ-NNN-progress.log` tail plus each worktree's last commit time until then.
14. An integrator silent 20 minutes or more (no new progress line, no new commit; the value is approximate) is stopped, and a fresh integrator resumes from its landed phase through `advance REQ-NNN`, never `recover`; it never redoes a merge already on the integration branch. A silent builder is reported only, unless the user says restart. Each tick re-runs the selector so newly captured REQs join the next wave. Each coordinator turn ends with one line: done/total, active lanes, ETA (or unknown). Teardown deletes the check; exactly one is armed per run.

E. Handoff and teardown:
15. `restart-with-parallel-handoff.md`: a run started with `--coordinate` resumes with `do-work run --coordinate --fan-out N`. Under `--coordinate` the coordinator follows that action on its own when the harness reports context usage above a threshold (the harness supplies the reading; a harness with no reading never triggers it), as well as at `phandoff`. The automatic handoff writes and commits the handoff only; it does not end the session.
16. Teardown also stops idle background agents (agents with no assigned REQ left; never one still building or integrating) and lists merged worktrees for removal without removing them.

F. Sweep and release:
17. Sweep restatements of the run argument list, the coordinator paragraph, the run-directory table and the builder test instruction across `actions/work.md`, `actions/work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md`.
18. Release per `_dev/primes/prime-releases.md`.
## Constraints
- No new REQ status and no new frontmatter field. New surface for the whole change: one flag, one optional policy file, two run-directory files, three manifest columns.
- Serial integration, the one-writer rule and dispatch-mechanism neutrality (`actions/work.md`) do not change. The flag names who integrates, not how an agent is spawned.
- No shipped Go code: no PID check, no lock helper, no preflight command, no resume parser. Every check here is coordinator-side judgment during a live run, on this machine, and the prose says so.
- Never `--force`, never `recover --take-over`, never a push, deploy or live-ops step during a run.
- Do not describe the rules as a fix; they make stalls and load failures visible and shorter. The `background-agents.md` ceiling note stays.
## Assumptions (recorded at capture)
- A harness without worktree or agent-dispatch support: plain `--fan-out` degrades silently to the serial loop; under `--coordinate` the degradation prints one line naming it (the builder may refuse instead and records why).
- "Writes briefs in writing gaps" is the existing one-writer rule; the coordinator's run-directory writes happen only in gaps.
- The lock holder label is the same writer label the checkpoint uses; no pid is recorded, because nothing may read it.
- If `do-work status` (REQ-690) is not yet shipped when this REQ is built, item 13's fallback wording stands and no edge is added; the builder records which path it documented.
## Dependencies
No `depends_on` edge. REQ-690 (do-work status) is the preferred liveness reader for item 13; the fallback wording covers its absence.
## Builder Guidance
High certainty on the flag, the five rules, the preflight lines and the resume command; they come from the source reports almost word for word. Latitude: subsection titles, column names, route phrasing, the degradation and stall-line wording. Keep each addition short; the reference already specifies entry, gaps and merge mechanics.
## Red-Green Proof
**RED prompt/case:** `do-work run --coordinate --fan-out 3`; the message "drive the queue"; `grep -n "progress.log\|full-gate.lock\|run-policy" skills/do-work/actions/fan-out-reference.md`; a `phandoff` typed inside a coordinated run.
**Why RED now:** `actions/work.md` rejects `--coordinate` as unrecognized; no SKILL.md row routes "drive the queue"; the grep finds nothing; the handoff's resume line is `do-work run --fan-out N`.
**GREEN when:** the flag parses and the phrase routes to the same mode; plain `--fan-out 3` is unchanged; the grep finds the rules subsection, both run-directory rows and the three manifest columns; the preflight and stall-loop paragraphs exist under Delegated integration; the handoff written from a coordinated run carries `--coordinate` in its resume command; `grep -rn "pid" skills/do-work/actions/fan-out-reference.md` finds no liveness check.
**Validation:** Inferred during capture (from the source report's Acceptance checks and the maintainer's 2026-10-10 decision).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: changes action routing and a downstream argument contract; families `restated-mechanism-unchecked` and `cross-action-exception-closure` fit the restatement sweep.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8815 tokens, over budget; `slugged: partial`). Matching reason: its REQ-083 lesson holds the liveness boundary this REQ must not cross.
## Full Context
See `do-work/user-requests/UR-153/input.md` for the decision record. The cancelled originals with their full bodies are under `do-work/archive/UR-146/` (REQ-662 to REQ-667); their source report is `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-coordinate-mode.md`.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: maintainer decision of 2026-10-10 (UR-153) folding REQ-662 to REQ-667: "fold 662/664/666 into one, cancel 665, 663 reduced to a prose checklist".*
