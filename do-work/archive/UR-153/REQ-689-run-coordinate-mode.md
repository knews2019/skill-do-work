---
id: REQ-689
title: 'do-work run --coordinate: the mandatory coordinator shape, its run rules, a prose preflight, a stall loop, and a handoff that resumes coordinated'
status: completed
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
required_lessons: [_dev/primes/lessons-releases.md]
write_set: [skills/do-work/SKILL.md, skills/do-work/actions/work.md, skills/do-work/actions/fan-out-reference.md, skills/do-work/actions/restart-with-parallel-handoff.md, skills/do-work/actions/work-reference.md, skills/do-work/docs/work-guide.md, skills/do-work/crew-members/background-agents.md, skills/do-work/docs/standing-preferences.md]
builder_handback_at: 2026-10-10T13:27:27Z
integration_at: 2026-10-10T18:22:38Z
review_at: 2026-10-10T18:38:00Z
kb_status: pending
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-10-10T13:18:58Z
  basis:
    - Route B
    - 8-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
commit: 9f815f154c8c206638d5d595ee15f85e3961c3f6
heavy_verified_at: 2026-10-10T18:39:03Z
heavy_verified_revision: 9f815f154c8c206638d5d595ee15f85e3961c3f6
completed_at: 2026-10-10T18:39:58Z
release_at: 2026-10-10T18:39:58Z
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
- [x] **[PLAN]:** (from the builder hand-back) Read the brief, REQ (all sections incl. Exploration PD-1 to PD-6 and Scope), UR-153, archived REQ-662/666/667 (as data), crew members general, coding-guardrails, shared-principles, communication-style, primes prime-action-files and prime-releases, lessons-releases. Verified before restating: queue-mode `advance` refuses unknown tokens (per the Exploration's cited `next_targets.go`), `verify.go` `diskSpaceLevelFor` returns critical/warning/neutral and emits a `low-disk-space` finding only for the first two, `actions/forensics.md` heading is `### 14. Release and Queue Invariants (board-owned)`, REQ-690 defines `--watch` as an output shape for a scheduler body. Plan: insert-only edits at the exact anchors in the brief; one subsection under Delegated integration holding the rules, preflight, run policy and stall loop; teardown in work.md Step 10 only (fan-out-reference's stall loop points there).
- [x] **[APPLY]:** (from the builder hand-back) Edits as planned, all insertions beside existing lines, no reflow. Six files changed, two named files left unchanged after the sweep.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff bd56c4b0 --stat`:
  ```
   skills/do-work/SKILL.md                            |  1 +
   skills/do-work/actions/fan-out-reference.md        | 29 +++++++++++++++++++++-
   .../actions/restart-with-parallel-handoff.md       |  2 ++
   skills/do-work/actions/work.md                     | 13 +++++-----
   skills/do-work/docs/standing-preferences.md        |  1 +
   skills/do-work/docs/work-guide.md                  |  2 +-
   6 files changed, 40 insertions(+), 8 deletions(-)
  ```
  Checks (from the worktree root):
  - RED probe at base: exit 1, 0.14 s wall.
  - GREEN probe (runs shipped-package-reference-contract.sh and action-shell-blocks.sh): exit 0, 3.96 s wall before commit; exit 0, 8.90 s wall after commit.
  - `grep -rniw pid skills/do-work/actions/fan-out-reference.md`: no output (exit 1 = no match).
  - `git diff --check`: exit 0, clean.
  - `_dev/tests/contract-regressions.sh`: exit 0, 34 s wall.
  - `_dev/tests/staged-skills-contract.sh`: refuses to run outside `maintainer-verify.sh --heavy` ("heavy-only"); not run, the gate is the integrator's.
  - Citation check is live: temporarily renaming my forensics citation to a nonexistent section made `shipped-package-reference-contract.sh` FAIL on that line; restored before commit.
  - Files checked by reading the diff: all six changed files; `work-reference.md:396`, `background-agents.md` (whole-file grep for argument list, run-directory table, builder test rule), `run-simple-reqs.md` (PD-6, not edited).
  - No debug artifacts; nothing under `do-work/` staged or committed.
*Source: maintainer decision of 2026-10-10 (UR-153) folding REQ-662 to REQ-667: "fold 662/664/666 into one, cancel 665, 663 reduced to a prose checklist".*

---

## Triage

**Route: B** - Medium

**Reasoning:** The what is fully specified (one flag, one route row, five run rules, a prose preflight, a stall loop, a resume line, a teardown line) and the files are named, but the change is prose spread over seven shipped files whose restatements must be found and kept in lock-step, so exploration of the exact anchors is needed. No new Go code and no design choice that needs a plan agent.

**Planning:** Not required


## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Required-lessons consult (Step 5).** `do-work/lessons-index.md` matched three rows. Kept: `_dev/primes/lessons-releases.md` (666 tokens, `slugged: full`; the REQ ships a release and `prime_files` names prime-releases), now in `required_lessons`. Still dropped for budget, unchanged from capture: `_dev/primes/lessons-action-files.md` (7756 tokens) and `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (8815 tokens), both `slugged: partial`, so no targeted form. Two bullets from the dropped satellites were read here and go into the builder brief by quote: lessons-action-files `restated-mechanism-unchecked` (REQ-639: Verified Facts are claims; read the command before shipped prose restates what it does) and lessons-do-kanban `disk-space-blind-spot` (the `low-disk-space` probe reports and never deletes; it measures the repo root only since 0.305.86).

**Anchors at main bd56c4b0 (0.305.101):**
- `skills/do-work/SKILL.md:33` is the plain run row; `:32` the run-with-recovery row above it. Routing is first-match-wins, and `run --coordinate` contains the word `run`, so the new row must sit between `:32` and `:33`. `:4` `argument-hint` lists no run flags; leave it (REQ-690/691/692 edit that line).
- `skills/do-work/actions/work.md:37` Architecture paragraph on `--fan-out`; `:99-116` `## Input` (`:104` `--fan-out`, `:108` stripped-token list, `:111` usage line); `:470` Step 10 teardown (checkpoint then cleanup after the last integrator); `:481` Step 0 checklist line; `:501` Step 10 checklist line.
- `skills/do-work/actions/fan-out-reference.md:101-118` Fan-Out Dispatch and Auto-wave; `:124-132` Delegated integration (five bullets: Entry, Never recover, Never Step 10, One writer, The integrator brief); `:134-151` run-directory table (`:144` is the `manifest.md` row). `:3` file header lists the sections it holds.
- `skills/do-work/actions/restart-with-parallel-handoff.md:11-13` When to Use (the `phandoff` trigger); `:61-64` resume command list (`:63` is `do-work run --fan-out N`); `:55` REMOVABLE rule (lists, never removes).
- `skills/do-work/docs/work-guide.md:132-134` Building several REQs at once (`:134` ends with the optional coordinator sentence).
- `skills/do-work/crew-members/background-agents.md:11-14` ceiling note (stays); `:180-188` isolation and the pointer to fan-out-reference; no run-argument list and no run-directory table, so the sweep likely changes nothing there.
- `skills/do-work/actions/work-reference.md:396` Worktree Dispatch Mode summary restates the integrator rules (enter by `advance REQ-NNN`, never recover, never Step 10); it gains at most one clause naming `--coordinate`.
- `skills/do-work/docs/standing-preferences.md:13` maps the pasted "use background agents" preference to built-in behavior; the pasted coordinator directive (the REQ's Why) has no row.

**Verified fact the REQ does not state.** Queue-mode `advance` shares `next`'s argument grammar (`skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go:39-83`): any token other than `--skip-impact-negligible`, `--simple`, `--wave`, `--fan-out` becomes a target token, and `next_targets.go:203` refuses it ("unrecognized argument ... expected REQ-NNN or UR-NNN"). So `--coordinate` must be consumed by the action and never forwarded to `advance`; bare `--coordinate` forwards a bare `--fan-out` (the CLI defaults that to 2, `next_commands.go:56-59`). The continuation argv `advance` returns never carries the mode (`lifecycleadvance/queue_commands.go:306-326`), so the coordinator holds it in session context and the handoff restates it.

**Other facts checked.** `grep -rn "progress.log\|full-gate.lock\|run-policy" skills/` finds nothing. No test reads the SKILL.md routing rows except `_dev/tests/staged-skills-contract.sh:289-296` (pipeline/full routes must be absent), so no routing test can show the bare word `coordinator` is collision-free. The board's disk reading is reached through `queue-kanban verify` exactly as `skills/do-work/actions/forensics.md:68-77` prescribes (build, then `verify --repo-root`; skip and say unverified when `go` is absent). `skills/do-work/actions/run-simple-reqs.md:24,67` forwards only `--fan-out` to `do-work run`.

**Decisions taken before dispatch (decide and state; reversible):**
- PD-1: Route phrases only (`drive the queue`, `use the main session as a coordinator`, `run --coordinate`), no bare `coordinator`: the REQ allows the bare word only when routing tests show no collision, and no such test exists.
- PD-2: `--coordinate` is an action-level flag. The action strips it before queue-mode `advance` and passes bare `--fan-out` when the user gave none. No Go change (Constraints forbid shipped Go).
- PD-3: REQ-690 (do-work status) is built in parallel and integrates before this REQ in the planned order. The stall loop names `do-work status --watch` as the reader and keeps the REQ's own fallback (progress-log tails plus each worktree's last commit time) for when the status action is absent. The builder records that both paths were documented.
- PD-4: Under `--coordinate` without worktree or agent-dispatch support, print one line naming the degradation and run the serial loop; do not refuse. Refusing would strand the queue on the simplest harness, which `actions/work.md:33` says must be able to follow the action.
- PD-5: `docs/standing-preferences.md` gains one row for the pasted coordinator directive pointing at `do-work run --coordinate`, because the REQ's Why is 13 sessions pasting that directive and that table is where a user checks whether a paste is still needed. One row, no other edit there.
- PD-6: `run-simple-reqs.md` keeps forwarding only `--fan-out`; extending it to `--coordinate` is not asked for and goes to Discovered Tasks.

*Generated by pre-dispatch exploration (coordinator session, run work-2026-10-10-131527)*

## Scope

**Files I will touch:**
- `skills/do-work/SKILL.md` (modify) — one route row above the plain run row
- `skills/do-work/actions/work.md` (modify) — Input flag entry, stripped-token list, usage line, Architecture paragraph, Step 10 teardown, two checklist lines
- `skills/do-work/actions/fan-out-reference.md` (modify) — mandatory-shape paragraph, run-rules subsection, preflight, stall loop, teardown additions, two run-directory rows, three manifest columns
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modify) — coordinated resume command and the automatic context-threshold trigger
- `skills/do-work/actions/work-reference.md` (modify) — one clause in the Worktree Dispatch Mode summary, only if the sweep finds the restatement stale
- `skills/do-work/docs/work-guide.md` (modify) — the coordinator sentence and the run-flag description
- `skills/do-work/crew-members/background-agents.md` (modify) — only if the sweep finds a restated run-directory table or builder test rule; expected unchanged
- `skills/do-work/docs/standing-preferences.md` (modify) — one row for the pasted coordinator directive (PD-5)
- `skills/do-work/docs/status-guide.md` (modify, added by the integrator as a merge seam) — drop progress logs from the run-local file example, because run-status skips every REQ-NNN name
- `skills/do-work/actions/status.md` (modify, added by the integrator for review F6) — full-gate.lock names a writer label, so no rm line
- `skills/do-work-toolbox/actions/ai-report.md` (modify, added by the integrator for the wave-end sweep F12) — Input sentence names the kind form
- `skills/do-work-toolbox/actions/architecture-report.md` (modify, added by the integrator for the wave-end sweep F13) — ai-report input contract restatement
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify, added by the integrator for the wave-end sweep F14) — opening sentence names the kind form

**Files I will NOT touch:** any Go source or test under `skills/do-work/tools/` or `skills/do-work-board/`; `skills/do-work/actions/run-simple-reqs.md` (PD-6); `skills/do-work/crew-members/communication-style.md` (the `phandoff` alias stays as is); the release paths (`CHANGELOG.md`, `VERSION`, mirrors), which the integrator writes at finalization.

**Acceptance criteria (restated from REQ):**
- [ ] `do-work run --coordinate --fan-out 3` parses: `--coordinate` is in `## Input`, the stripped-token list, the usage line and the Step 0 checklist line; bare `--coordinate` implies `--fan-out`
- [ ] "drive the queue", "use the main session as a coordinator" and "run --coordinate" route to `./actions/work.md` with `--coordinate`, from a row above the plain run row
- [ ] Plain `do-work run --fan-out 3` behaves exactly as today (the optional coordinator shape stays optional)
- [ ] Under Delegated integration: the mandatory-shape paragraph (dispatch, briefs in gaps, report; never builds, merges or runs Step 6 to Step 9; no integrator pushes; blocking question becomes `pending-answers`)
- [ ] One run-rules subsection with heartbeat, full-gate lock, red-gate triage, manifest timing row, hand-back coverage, and the statement that the progress log and the lock are not queue state, a liveness claim, or the machinery REQ-069 and REQ-073 deleted; `grep -n "progress.log\|full-gate.lock\|run-policy" skills/do-work/actions/fan-out-reference.md` finds the subsection, both run-directory rows and the three `manifest.md` columns
- [ ] Preflight paragraph (disk level from the board's verify, leftover processes, stale `full-gate.lock`, `do-work/run-policy.md` bullet count or absent) and the optional `do-work/run-policy.md` copied into every brief
- [ ] Stall-loop paragraph: one recurring check every 15 to 20 minutes, reads `do-work status --watch` (fallback: progress-log tails and worktree last commit), restarts an integrator silent 20 minutes or more through `advance REQ-NNN`, never `recover`, reports a silent builder only, one-line turn status, deleted at teardown
- [ ] A handoff written from a coordinated run resumes with `do-work run --coordinate --fan-out N`; the automatic context-threshold handoff writes and commits only
- [ ] Teardown stops idle background agents and lists merged worktrees for removal without removing them
- [ ] `grep -rn "pid" skills/do-work/actions/fan-out-reference.md` finds no liveness check
- [ ] Restatements swept across `actions/work.md`, `actions/work-reference.md`, `docs/work-guide.md`, `crew-members/background-agents.md`

## Implementation Summary

**Files changed:**
- `skills/do-work/SKILL.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/fan-out-reference.md` (modified)
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)
- `skills/do-work/docs/standing-preferences.md` (modified)
- `skills/do-work/docs/status-guide.md` (modified, integrator seam and review F6)
- `skills/do-work/actions/status.md` (modified, review F6)
- `skills/do-work-toolbox/actions/ai-report.md` (modified, wave-end sweep F12)
- `skills/do-work-toolbox/actions/architecture-report.md` (modified, wave-end sweep F13)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified, wave-end sweep F14)

**What was done:** `do-work run --coordinate [--fan-out N] [REQ-NNN ...]` now exists. `work.md` `## Input` defines the flag (composes with `--fan-out`, `--wave` and targets; bare `--coordinate` implies `--fan-out`; the action consumes it and never passes it to queue-mode `advance`; the session keeps the mode because the continuation does not carry it; without worktree or agent dispatch it prints one line and runs the serial loop), and the stripped-token list, usage line, Architecture paragraph, Step 10 teardown and both checklist lines name it. `SKILL.md` routes `drive the queue`, `use the main session as a coordinator` and `run --coordinate` to `./actions/work.md` with `--coordinate`, from a row between run-with-recovery and the plain run row. `fan-out-reference.md` → Delegated integration gains the mandatory-shape paragraph and a `#### Coordinated run rules` subsection: progress log, full-gate lock, red-gate triage, manifest timing columns, hand-back coverage, the REQ-069/REQ-073 boundary sentence, a four-line prose preflight (board disk level, leftover processes, stale locks, optional `do-work/run-policy.md`) and the stall loop; the run-directory table gains the `REQ-NNN-progress.log` and `full-gate.lock` rows and three `manifest.md` columns. `restart-with-parallel-handoff.md` adds the automatic context-threshold trigger under `--coordinate` and the `do-work run --coordinate --fan-out N` resume line. `work-guide.md` describes the flag in one sentence; `standing-preferences.md` maps the pasted coordinator directive to it. `work-reference.md` and `background-agents.md` were swept and left unchanged (D-07, D-08). Integrator seams in the merge commit: REQ-690's `do-work status` has shipped, so the stall loop reads `do-work status --watch` (`actions/status.md`) plus the progress-log tails (status does not read them) and keeps the tails-and-commit reading only for when status cannot report (no board or no Go toolchain); `status-guide.md` no longer says progress logs show as run-local files, because `run-status` skips every `REQ-NNN-*` name. Review fixes (integrator, on the builder branch `927f61d8` and `209e0bb5`, re-merged as `9cd669eb` and `9f815f15`): a builder may also write its own `REQ-NNN-progress.log` by the absolute main-tree path, and the hand-back merge allows but never stages progress logs and `full-gate.lock` (F1, N2); the stall loop does not stop an integrator whose last line starts a command still inside its usual run time, and frees the lock of an integrator it stopped (F2); C2 rows are acted on only for this session's agents, and a C3 row means "brief the next integrator" (F3); a coordinated handoff writes claim lines as "hand REQ-NNN to an integrator" and the automatic handoff waits for a writing gap, runs once per run and still ends with its two announcement lines (F4, F5, N1); `status` never offers `rm` for `full-gate.lock` (F6); the coordinator never runs the span from the hand-back merge through Step 9 (F7); the argument hint shows `run [REQ|UR] [--fan-out [N]] [--coordinate]` (F8); the work-guide trigger aliases name the two phrases (F9). Wave-end sweep fixes: Step 9 cleanup pointers in `work.md` and `fan-out-reference.md` (F10, F11) and the ai-report `--kind` input wording in `ai-report.md`, `architecture-report.md` and `ai-report-guide.md` (F12 to F14).

## Decisions

*(from the builder hand-back)*

- D-01 DECIDE & STATE: route phrasing is exactly `drive the queue`, `use the main session as a coordinator`, `run --coordinate`; no bare `coordinator` (PD-1). Route cell `./actions/work.md` with `--coordinate`.
- D-02 DECIDE & STATE: subsection title `#### Coordinated run rules`. Preflight, run policy and stall loop sit inside it as bold-led paragraphs, each stating its own condition (`--coordinate` only), while the five rules apply whenever the coordinator shape is used. Keeps the new heading count at one.
- D-03 DECIDE & STATE: the heartbeat rule is titled **Progress log**, not "Heartbeat", so it does not read as the `heartbeat_at` machinery REQ-069/REQ-073 deleted, which the next sentence of the same subsection rules out.
- D-04 DECIDE & STATE: degradation wording under `--coordinate`: it prints one line, `--coordinate: no worktree or agent dispatch here, running the serial loop`, and runs the serial loop (PD-4). Plain `--fan-out` still degrades silently.
- D-05 DECIDE & STATE: the Input entry adds one clause that `advance`'s continuation does not carry the flag, so the session keeps the mode for the whole run (fact from the Exploration, `lifecycleadvance/queue_commands.go:306-326`). Without it a coordinator following the continuation argv could drop the mode.
- D-06 DECIDE & STATE: "Each tick re-runs the selector" gained "in a writing gap", because a queue-mode selection commits a claim and the existing one-writer rule forbids coordinator writes while an integrator runs.
- D-07 DECIDE & STATE: `work-reference.md:396` left unchanged. It restates the integrator rules (entry, never recover, never Step 10), none of which changed; it says nothing about the shape being optional, so it is not stale.
- D-08 DECIDE & STATE: `background-agents.md` left unchanged. It restates neither the run argument list nor the run-directory table nor a builder test rule; its manifest example is a generic code-review manifest. Ceiling note kept.
- D-09 DECIDE & STATE: the preflight disk line cites `actions/forensics.md` → **14. Release and Queue Invariants (board-owned)** (the heading's full text; the brief's shorter form would fail the section-citation check). It quotes the level as warning, critical or none: `verify.go` emits a `low-disk-space` finding only at warning or critical, and "neutral" otherwise.
- D-10 DECIDE & STATE: the manifest-row rule names the coordinator as the writer of the three columns, in the gap after each integration, consistent with the existing "the orchestrator's, never written by a builder" and the one-writer rule.
- D-11 DECIDE & STATE: the full-gate lock is removed by the integrator after its suite; a lock left behind is what the preflight reports as stale and never deletes.

## Discovered Tasks

*(from the builder hand-back)*

- `skills/do-work/actions/run-simple-reqs.md:24,64,67` forwards only `--fan-out` to `do-work run`, so a simple-REQs run cannot be coordinated (PD-6). impact-user-visible → report only
- `skills/do-work/SKILL.md:4` `argument-hint` lists no run flags at all (`run [REQ|UR]`), so neither `--fan-out` nor `--coordinate` is discoverable there; left alone because REQ-690/691/692 edit that line. impact-negligible → report only

## Qualification

**Gate records.** `advance REQ-689 --diff-range cc049a01..f8f32682`: `qualify` satisfied (P-A-U boxes ticked from the hand-back, no debug artifacts; `git diff cc049a01..f8f32682` has no `console.log`, `TODO`, `XXX` or `DEBUG`). `scope-drift` first reported `SCOPE-UNDECLARED-TOUCH` for `skills/do-work/docs/status-guide.md`, the integrator's merge seam; I added that file to Scope with its reason and the error cleared. Two `SCOPE-DECLARED-NOT-TOUCHED` warnings remain, for `actions/work-reference.md` and `crew-members/background-agents.md`: Scope declared both as "only if the sweep finds a restatement", and the builder's D-07 and D-08 record that it did not. Declared 9, touched 7, nothing touched outside Scope.

**Requirement trace against `git diff cc049a01..f8f32682 --stat` (7 files, +41 -9) and the merged files:**
- 1 (flag): `work.md` `## Input` has the `--coordinate` entry (composes with `--fan-out [N]`, `--wave N`, targets; bare form implies `--fan-out`); the stripped-token list, the usage line and the Step 0 checklist line name it.
- 2 (route): `SKILL.md:33` routes the three phrases to `./actions/work.md` with `--coordinate`, above the plain run row at `:34`; no bare `coordinator` (PD-1).
- 3 (mandatory shape): `fan-out-reference.md` Delegated integration carries the paragraph word for word.
- 4 (plain `--fan-out`): the `--fan-out [N]` entry and "degrades silently to the serial loop" are byte-unchanged; only the stripped list gains `--coordinate`.
- 5 to 10 (run rules): `#### Coordinated run rules` holds Progress log, Full-gate lock, Red full gate, Manifest row, Hand-back coverage and the REQ-069/REQ-073 boundary sentence; the run-directory table has the `REQ-NNN-progress.log` and `full-gate.lock` rows and the three `manifest.md` columns. `grep -niw pid` on the file finds nothing.
- 11 and 12 (preflight, run policy): four numbered lines (disk level from the board's verify via `actions/forensics.md` section 14, leftover processes, stale locks, run policy) and the optional `do-work/run-policy.md` paragraph.
- 13 and 14 (stall loop): one recurring check every 15 to 20 minutes, turn-start fallback, `do-work status --watch` plus progress-log tails (integrator seam, REQ-690 shipped), restart of an integrator silent about 20 minutes through `advance REQ-NNN`, never `recover`, silent builder reported only, selector re-run in a writing gap, one-line turn status, deleted at teardown.
- 15 (handoff): `restart-with-parallel-handoff.md` has the context-threshold trigger (writes and commits only) and the `do-work run --coordinate --fan-out N` resume line.
- 16 (teardown): `work.md` Step 10 and its checklist line stop idle agents, delete the stall check and list merged worktrees without removing them.
- 17 (sweep): `work-reference.md`, `background-agents.md` unchanged with stated reasons; `work-guide.md` gains one sentence; the wave-end Restatement Sweep goes to the reviewer.
- 18 (release): at finalization.

**Scope comparison.** Every touched file is in Scope. No Go, no new action file, no frontmatter field, no status. The extras the builder named (standing-preferences row, the continuation clause, "in a writing gap", "removing it after") are each one line with a recorded reason (PD-5, D-05, D-06, D-11).

## Testing

**Tests run:**
- Repository gate `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` in ROOT, after waiting for no other gate and 1-minute load under 5: exit 0 at `f8f32682` (115 s, load 2.27), exit 0 at `9cd669eb` (117 s, load 4.97), exit 0 at `9f815f15` (117 s, load 4.65). No retry needed.
- `advance REQ-689 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-689-probe.sh` at `9f815f15`: green-gate satisfied, run-blocked-check satisfied, the GREEN probe exit 0 (`BLOCKED-PROBE-SUCCEEDED`). The two scope-drift warnings are the conditional files from Qualification.
- The GREEN probe runs `shipped-package-reference-contract.sh` and `action-shell-blocks.sh` and checks the flag, route, rules, rows, columns, preflight, stall loop, resume line and teardown text.

**Result:** pass at all three merges.

**Red-green validation:** `tdd: false`. Builder proof from the hand-back: RED probe at base `bd56c4b0` exit 1 with 19 failed checks; GREEN probe exit 0 at the builder commit. The integrator's GREEN probe runs at `f8f32682` (reviewer, 2.6 s), `9cd669eb` (reviewer) and `9f815f15` (advance) all exit 0.

**New/updated tests:** none. The REQ is prose only; the GREEN probe is a run artifact, not a shipped test.

**Heavy verification plan:** range `cc049a01ecf60b538745a154d2940b99140331e6..9f815f154c8c206638d5d595ee15f85e3961c3f6` (cumulative, three merges). One lane: `staged-skills`, argv `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`, reason: the 11 changed files match subtree `skills`. Early drain at `f8f32682`: executed, exit 0, 36 s (a first try failed in 0 s on a wrong `test-durations.tsv` header I wrote; rerun with ROOT's header line).

## Review

**Overall: 86%** | 2026-10-10T18:38:00Z

**Verdict:** Approve after fixes F1, F2, F4, F5 (one sentence each; exact text above).

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 75% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 progress log contradicts the builder's one-exception write rule (`fan-out-reference.md:61`, `:73`, `:146`, `:173`); fix text in report — impact-user-visible → fix
- F2 stall loop stops an integrator mid full gate and strands `full-gate.lock` (`fan-out-reference.md:163`); fix text in report — impact-user-visible → fix
- F4 coordinated handoff's claim lines make the next session integrate in the main session (`restart-with-parallel-handoff.md:68`) — impact-user-visible → fix
- F5 automatic handoff writes during an integrator's span, refires each turn (`restart-with-parallel-handoff.md:14`) — impact-user-visible → fix
- F6 `status` offers `rm` for a live `full-gate.lock` (`status.md:63`, `status-guide.md:34`); Go comment `run_status.go:363` report only — impact-user-visible → fix

**Minor findings:** F3 stall loop lacks the this-machine label and C3 handling (`fan-out-reference.md:163`) — impact-negligible → fix; F7 "never runs Step 6 to Step 9" vs dispatch in Step 6 (`:134`) — impact-negligible → fix; F8 argument-hint lacks run flags (`SKILL.md:4`) — impact-negligible → fix; F9 trigger-alias list omits the new phrases (`work-guide.md:146`) — impact-negligible → fix; F10 "Step 8 substep 8" (`work.md:575-576`) — impact-negligible → fix; F11 "Step 8's" cleanup (`fan-out-reference.md:47`) — impact-negligible → fix; F12 `ai-report.md:32` Input sentence — impact-negligible → fix; F13 `architecture-report.md:135` — impact-negligible → fix; F14 (nit) `ai-report-guide.md:3` — impact-negligible → fix; F15 (nit) anti-bloat count, all within budget — impact-negligible → report only
**Acceptance:** Pass — GREEN probe at HEAD exit 0 (flag, route, rules, rows, columns, preflight, stall loop, resume line, teardown present; no `pid` in fan-out-reference.md)
**Restatement sweep:** redefined the run argument set (`--coordinate`), the coordinator shape (mandatory under the flag), the run-directory file set (`REQ-NNN-progress.log`, `full-gate.lock`) and three manifest columns, the builder's main-tree writes, the handoff triggers and resume command, the Step 10 teardown, and the core run route phrases; stale: F1, F4, F5, F6, F8, F9; inherited elements of all 11 siblings swept (none unread): F10, F11, F12, F13, F14 new, earlier recorded REQ-688 F2/F3, REQ-690 F8/F9, REQ-692 M2, REQ-660 F3 still stale
**Suggested testing:** 2 items — (1) read-through of a coordinated handoff paste block with one leftover claim after F4; (2) `do-work status` against a run directory holding a `full-gate.lock` with a label first line after F6
**Follow-ups created:** None (15 findings report only or fix-in-integration)

**Integrator disposition (2026-10-10T18:38:00Z):** F1 to F14 fixed with the reviewer's text on the builder branch (`927f61d8`), re-merged as `9cd669eb`; F5's last sentence changed so the two announcement lines still end the message (`restart-with-parallel-handoff.md` Step 6). Delta re-check of `f8f32682..9cd669eb`: 95%, Approve; N1 and N2 fixed (`209e0bb5`, re-merged as `9f815f15`); N3 (a remote builder gets no main-tree path, so it keeps no progress log; older than this REQ) → report only. F15 → report only. Still-stale items recorded by earlier siblings (REQ-688 F2/F3, REQ-690 F8/F9, REQ-692 M2, REQ-660 F3, the `run_status.go:363` comment) → report only. Gate and heavy lane re-run at `9f815f15`. Full report: `do-work/runs/work-2026-10-10-131527/REQ-689-review.md`.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** the wave-end reviewer, launched right after qualify, found five seams between the new run rules and older rules (the builder's one main-tree write, the one-writer rule, the handoff's claim lines, a long gate inside the 20-minute silence window, and `status`'s `rm` line for locks) plus five stale sibling restatements; each fix was one sentence and the delta re-check took under a minute.

**What didn't:** a review fix written against one file broke a rule in the same file's Step 6 (the handoff message must end with its announcement lines). Reading the fix target's own Red Flags before applying the text caught it; the re-check then caught the matching sentence in `fan-out-reference.md` (N1).

**Worth knowing:** a new run-directory file is also a new write path. Before adding one, check who may write the main tree (`fan-out-reference.md` → Sole integrator), what the hand-back merge stages (step 0), and how `run-status` lists it. A brief's section citation is a claim too: cite a heading's full text, number included, and let `shipped-package-reference-contract.sh` prove it (builder lesson).

## Orientation

Coordinated runs: `actions/work.md` `## Input` owns the `--coordinate` flag; `actions/fan-out-reference.md` → Delegated integration → Coordinated run rules owns the run rules, preflight and stall loop; `actions/restart-with-parallel-handoff.md` owns the coordinated resume.

## Heavy Verification Plan

- Base: `cc049a01ecf60b538745a154d2940b99140331e6`
- Target: `9f815f154c8c206638d5d595ee15f85e3961c3f6`
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`; reasons: the 11 changed files under `skills/do-work/` and `skills/do-work-toolbox/` match subtree `skills`.

## Heavy Verification Result

- Target: `9f815f154c8c206638d5d595ee15f85e3961c3f6`; execution revision `9f815f154c8c206638d5d595ee15f85e3961c3f6` (detached checkout `.git/work-run-work-2026-10-10-131527/drain-head-REQ-689`, removed after).
- `staged-skills`: executed, exit 0, 32 s; no `HEAVY-RUN-LANE-SKIPPED`. Earlier drain at `f8f32682`: executed, exit 0, 36 s.

## Timing

Observed 2026-10-10T18:21:44Z to 2026-10-10T18:38:54Z: 17m 10s total, 21m 52s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 14m 30s | 1 |
| verification-gate | 6m 28s | 4 |
| handback-merge | 54s | 1 |

Slowest stage: review / review, delta re-check and review fixes, 14m 30s, outcome success.

Notes: no builder-work event; the hand-back (builder commit 2026-10-10T13:27:27Z) landed long before integration began. The review event runs from the reviewer spawn to the recorded Review and overlaps the three gate runs, so attributed time exceeds the observed span.
