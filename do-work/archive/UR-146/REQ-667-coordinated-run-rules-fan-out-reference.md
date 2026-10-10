---
id: REQ-667
title: 'Coordinated-run rules for heartbeat log, full-gate lock, focused builder suites, red-gate triage and manifest timing go into fan-out-reference'
status: cancelled
created_at: 2026-10-09T21:16:07Z
user_request: UR-146
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-662, REQ-663, REQ-664, REQ-665, REQ-666]
batch: run-coordinate-mode
claimed_at: 2026-10-10T12:52:20Z
status_changed_at: 2026-10-10T12:57:16Z
completed_at: 2026-10-10T12:59:17Z
---
# Coordinated-Run Rules Go Into fan-out-reference
## What
Add one subsection under `### Delegated integration — the coordinator shape` in `skills/do-work/actions/fan-out-reference.md` that holds five run rules: a per-REQ heartbeat log, a full-gate lock so only one integrator runs the full suite at a time, focused suites for builders, red-gate triage, and a manifest timing row. Add the two new run-directory files and the three manifest columns to the run-directory table.
## Why
Report item A6 (UR-146 input, "What happened"): on 2026-10-07 to 10-08 a coordinated run in one consumer repo integrated 22 REQs. Background integrators stalled for 3 hours and 2 hours because they ended their turn while a command ran and the completion notice did not wake them. Concurrent full test suites caused load-only failures. After the user set these five rules, each integration took 11 to 57 minutes, and a later 3-REQ run took 10 to 21 minutes per integration with every full gate green on its first run. The rules live today only in one per-machine memory note, so a second machine or repo starts without them.
## Verified Facts (checked at 0.305.87)
- `skills/do-work/actions/fan-out-reference.md:122` is `### Delegated integration — the coordinator shape`; `:132` is `### Run directory, briefs and hand-backs`, whose table lists run directory `do-work/runs/work-<YYYY-MM-DD-HHMMSS>/`, `REQ-NNN-brief.md`, `REQ-NNN-handback.md`, `REQ-NNN-integrate.md`, `manifest.md` and bounded waves. No progress log, no gate lock, and no timing columns.
- `grep -rn "progress.log\|full-gate.lock" skills/` finds nothing today.
- `skills/do-work/crew-members/background-agents.md:11-14` says the background pattern makes failures "survivable and recoverable", not prevented.
## Detailed Requirements
1. **Heartbeat.** Every builder and integrator appends `<UTC> <phase> <command>` to `do-work/runs/<run>/REQ-NNN-progress.log` before and after any command expected to take over a minute. An agent never ends its turn while a command it started is still running.
2. **Full-gate lock.** Only integrators run the full test suite, one at a time, holding `do-work/runs/<run>/full-gate.lock` with the owner pid and UTC start inside. Builders run focused suites for the files they touched.
3. **Red full gate.** Rerun the failed suites alone. Failures that pass alone are load-only: one full rerun. A real regression: fix it, then run one full gate.
4. **Manifest row per integration:** takeover-to-finalization minutes, full gates run, stall restarts (three new `manifest.md` columns).
5. **Hand-back coverage.** Before writing an integrator brief, the coordinator checks that the hand-back covers every item forwarded to the builder mid-run.
6. Add `REQ-NNN-progress.log` and `full-gate.lock` rows to the run-directory table at `fan-out-reference.md:132`, and the three columns to its `manifest.md` row.
7. Sweep restatements of the run-directory table and the builder test instruction across `actions/work.md`, `actions/work-reference.md`, `crew-members/background-agents.md` and `docs/work-guide.md`.
8. Release per `_dev/primes/prime-releases.md`.
## Constraints
- No new REQ status and no new frontmatter field (report, Request).
- Serial integration, the one-writer rule and dispatch-mechanism neutrality (`actions/work.md:37`) do not change (report, Out of scope).
- Do not describe the rules as a fix. They make stalls and load failures visible and shorter; the `background-agents.md` ceiling note stays.
- No push, deploy or live-ops step (report, Out of scope).
## Assumptions (recorded at capture, no questions asked)
- These rules apply whenever the coordinator shape is used, not only under `--coordinate`. REQ-662 (the `--coordinate` run mode) makes the shape mandatory and points at them. That is why this REQ is the root of the batch.
- Liveness history: REQ-069 and REQ-073 (fan-out dispatch, v0.161.0) deleted an orchestrator lock, `heartbeat_at` and a claim registry, and `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:45` still says no lock, heartbeat, PID check or mtime heuristic may be added to the board's worktree verify. This REQ does not reverse that. The progress log is an append-only run-directory text file the coordinator reads, not queue state, and the full-gate lock serializes one test command inside one run, not queue ownership. The builder states this distinction in the subsection so a later reader does not read it as the deleted machinery coming back. If the builder finds a live rule that forbids even this, it records the conflict in Decisions and stops that requirement.
- "Focused suites for the files they touched" is the builder's judgment per project; the rule names the condition, not a command list.
- The lock file is written and removed by the integrator. A lock whose owner pid is dead is stale; detecting and reporting it is REQ-663 (preflight) and REQ-665 (resume drift check), not this REQ.
## Dependencies
None. REQ-662, REQ-663, REQ-664 and REQ-665 depend on this REQ.
## Builder Guidance
High certainty: the report gives the five rules almost word for word. Latitude: subsection title, exact column names, and the wording of the liveness distinction.
## Red-Green Proof
**RED prompt/case:** `grep -n "progress.log\|full-gate.lock" skills/do-work/actions/fan-out-reference.md`
**Why RED now:** no match; the rules exist only in a consumer's memory note.
**GREEN when:** the grep finds the A6 subsection under Delegated integration, the run-directory table has both new rows, and `manifest.md`'s row names the three timing columns.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: changes a downstream reader contract (the run-directory table); family `restated-mechanism-unchecked` fits the restatement sweep.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8334 tokens, over budget; `slugged: partial`). Matching reason: its REQ-083 lesson holds the liveness boundary this REQ must not cross.
## Full Context
See `do-work/user-requests/UR-146/input.md` for complete verbatim input (sections Request A6, What happened, Where the behaviour lives today, Proposed direction A6, Acceptance check). No queued candidate shares this root cause (the queue held REQ-654 to REQ-661, the ai-report and do-work-cli batches, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-coordinate-mode.md`, Request item A6: "The run rules (heartbeat line, full-gate lock, focused suites for builders, red-gate triage, manifest timing row) are written into `actions/fan-out-reference.md`. None of them is in the suite today."*

## Cancelled

- **When:** 2026-10-10T12:59:17Z
- **Why:** folded into the simplified recapture of 2026-10-10 (maintainer chose the aggressive simplification)
- **Decided by:** user, via `do-work abandon`
