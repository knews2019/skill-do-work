---
id: REQ-662
title: 'do-work run --coordinate makes the coordinator shape mandatory and routes from drive the queue'
status: pending
created_at: 2026-10-09T21:16:07Z
user_request: UR-146
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-667]
related: [REQ-663, REQ-664, REQ-665, REQ-666, REQ-667]
batch: run-coordinate-mode
---
# do-work run --coordinate Makes the Coordinator Shape Mandatory
## What
Add a run flag `do-work run --coordinate [--fan-out N] [REQ-NNN ...]` and a SKILL.md route for "drive the queue", "coordinator" and "use the main session as a coordinator". Under the flag the coordinator shape is mandatory: the main session never builds and never integrates, one builder subagent per REQ runs in its own worktree, one integrator runs at a time under the full-gate lock, nothing is pushed, and blocking questions are parked as `pending-answers`.
## Why
Report item A1 (UR-146 input, "What happened"): about 13 sessions from 2026-09-27 to 2026-10-09 in two consumer repos and a development clone started with a pasted coordinator directive, 5 of them with the identical paste ("use the main session as a coordinator, use background agents/workflows to actually run the queue use worktrees with parallel agents up to capacity to finish quickly"). Sessions that missed the paste integrated in the main session and grew to about 670k tokens of context before a handoff. Today the shape is optional ("the orchestrator may hand"), and `--coordinate` is rejected as an unrecognized argument.
## Verified Facts (checked at 0.305.87; the report cites 0.305.84 line numbers)
- `skills/do-work/SKILL.md:33`: the run route is `run`, `go`, `start`, `work`, `begin`, `process`, `execute`, `build`, `continue`, `resume`. No coordinator phrase routes anywhere.
- `skills/do-work/actions/work.md:99` is `## Input` (report: 101). `:104` defines `--fan-out [N]` (report: 106). `:108` says unrecognized arguments are rejected after stripping `--wave N`, `--fan-out [N]` and `--skip-impact-negligible` (report: 110); `:111` is the usage line. `:481` is the Step 0 checklist line (report: 489).
- `actions/work.md:37`: under fan-out "integration stays serial" and "the dispatch mechanism stays unspecified".
- `actions/fan-out-reference.md:122-130`: Delegated integration, "the orchestrator may hand each REQ's integration ... to one agent at a time". It already states entry by `advance REQ-NNN`, never `recover`, never Step 10, and one writer under the project root at a time.
- `actions/work.md:470`: the coordinator runs the checkpoint, then cleanup, after the last integrator returns.
- `skills/do-work/docs/commit-guide.md:32`: "Never pushes to remote". The coordinator paragraph does not repeat it for integrators.
- REQ-639 (the delegated-integration coordinator shape, archived in UR-138, commit c40c6460) decided "There is no flag" and kept the shape a placement. This REQ reverses that decision at the maintainer's request by capturing this report; the optional shape stays for plain `--fan-out`.
## Detailed Requirements
1. In `actions/work.md` → `## Input`, accept `--coordinate`. It composes with `--fan-out [N]`, `--wave N` and targeting tokens. Bare `--coordinate` implies `--fan-out`.
2. Add `--coordinate` to the tokens stripped before the unrecognized-argument check (`:108`), to the usage line (`:111`) and to the Step 0 checklist line (`:481`).
3. Add a SKILL.md route row above the plain run row: `drive the queue`, `coordinator`, `use the main session as a coordinator`, `run --coordinate` → `./actions/work.md` with `--coordinate`.
4. In `actions/fan-out-reference.md` → Delegated integration, add: "Under `--coordinate` this shape is mandatory. The main session dispatches, writes briefs in writing gaps, and reports. It never builds, never merges and never runs Step 6 to Step 9 itself. No integrator pushes. A blocking question becomes a `pending-answers` follow-up, never a wait."
5. State that integrators hold the full-gate lock from REQ-667 (the coordinated-run rules) when they run the full suite.
6. Plain `do-work run --fan-out 3` behaves exactly as today.
7. Sweep restatements of the run argument list and the coordinator paragraph across `actions/work.md`, `actions/work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md`.
8. Release per `_dev/primes/prime-releases.md`.
## Constraints
- No new REQ status and no new frontmatter field. New surface for the whole batch: one flag, one optional policy file, two run-directory files and three manifest columns (report, Request).
- Serial integration, the one-writer rule and dispatch-mechanism neutrality at `actions/work.md:37` do not change (report, Out of scope). The flag names who integrates, not how an agent is spawned.
- No push, deploy or live-ops step during a run (report, Out of scope).
## Assumptions (recorded at capture, no questions asked)
- The flag is parsed by the agent from `actions/work.md`, like `--fan-out`. If a do-work-cli command also validates run arguments, the builder adds the token there too and records it in Decisions.
- A harness without worktree or agent-dispatch support: plain `--fan-out` degrades silently to the serial loop (`work.md:104`). Under `--coordinate` the user asked for the coordinator explicitly, so the degradation prints one line naming it instead of staying silent. The builder may choose to refuse instead and records why.
- "Writes briefs in writing gaps" means the existing gap rule at `fan-out-reference.md` (one writer under the project root at a time). The acceptance line "the main session makes no commit and no edit under the project root while an integrator runs" is that same rule; the coordinator's own run-directory writes happen only in gaps.
- The bare word `coordinator` as a route trigger may catch unrelated sentences. The builder may narrow it to phrases if routing tests or the SKILL.md conventions show a collision, and records the choice.
- REQ-663 (preflight), REQ-664 (stall check) and REQ-666 (handoff and teardown) attach to this mode. This REQ only adds the mode and makes the shape mandatory.
## Dependencies
Depends on REQ-667 (the run rules, which hold the full-gate lock this mode requires). REQ-663, REQ-664 and REQ-666 depend on this REQ.
## Builder Guidance
High certainty on the flag and its composition. Lower on the exact route phrases and the degradation wording. Latitude: wording of the mandatory paragraph, as long as each clause in requirement 4 is present.
## Red-Green Proof
**RED prompt/case:** `do-work run --coordinate --fan-out 3`, and separately the message "drive the queue".
**Why RED now:** `actions/work.md:108` rejects `--coordinate` as unrecognized, and no SKILL.md row routes "drive the queue".
**GREEN when:** `do-work run --coordinate --fan-out 3` parses, "drive the queue" routes to the same mode, plain `do-work run --fan-out 3` behaves exactly as today, and in a coordinated run the main session makes no commit, no edit under the project root while an integrator runs, and no merge.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: changes action routing and a downstream argument contract; families `restated-mechanism-unchecked` and `cross-action-exception-closure` fit the restatement sweep.
## Full Context
See `do-work/user-requests/UR-146/input.md` for complete verbatim input (sections Request A1, What happened, Where the behaviour lives today, Proposed direction A1, Acceptance check). No queued candidate shares this root cause (the queue held REQ-654 to REQ-661, the ai-report and do-work-cli batches, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-coordinate-mode.md`, Request item A1: "A run mode. `do-work run --coordinate [--fan-out N] [REQ-NNN ...]`, routed also from "drive the queue", "coordinator" and "use the main session as a coordinator"."*
