---
id: REQ-668
title: 'do-work status action and do-work-cli run-status report one class, ETA and remedy per open REQ'
status: cancelled
created_at: 2026-10-09T21:20:59Z
user_request: UR-147
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-10-10T12:52:21Z
status_changed_at: 2026-10-10T12:57:17Z
completed_at: 2026-10-10T12:59:20Z
---
# do-work status: One Class, ETA and Remedy per Open REQ
## What
Add one read-only action, `do-work status` (`skills/do-work/actions/status.md`), backed by one deterministic subcommand, `do-work-cli run-status [--run <dir>] [--req REQ-NNN] [--watch] --format json|text`. For every claimed, needs-input, blocked or waiting REQ it prints one row: phase, minutes since last activity, ETA, one class (C1 to C8) and one remedy. Route "status", "status and ETA", "is it stuck", "how do I unblock" and similar phrases to it. It reports ages of durable records and never says a run or builder is dead.
## Why
Report Request and evidence (UR-147 input): in one consumer repo, from 2026-10-02 to 2026-10-09, about 20 user prompts in 12 sessions on 7 days asked for status, stuck lanes, unblocking or a watchdog. `do-work forensics` was used 0 times in 30 days. One coordinator session made 1005 tool calls, 378 of them nudges. Each answer was rebuilt by hand from REQ files, the run manifest, `git log` per builder branch, `ls -l` of the run directory and `ps`. The facts exist on five surfaces, and two of them (last activity, remaining time) are only in the HTML board.
## Verified Facts (checked at 0.305.87)
- Routing: `skills/do-work/SKILL.md:42` routes forensics on `forensics`, `diagnose`, `health check` only. `SKILL.md:43` routes `roadmap`, `queue-status`, `where are we`, `what's left` to roadmap, which says "A planning aid, not a diagnostic" (`actions/roadmap.md:5`). `SKILL.md:38` routes the word `blocked` to clarify. `skills/do-work-board/actions/board.md:33` routes `summary`, `status`, `counts` to column counts.
- `skills/do-work/actions/fan-out-reference.md:165`: "It asks about progress or state. Answer from disk: the REQ files, the run manifest, the hand-backs." No command stands behind this rule.
- `skills/do-work-board/tools/queue-kanban/open_work.go:79` (`writeClaimedSection`) prints each claimed REQ as id and title only.
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go:40-42`: `LastActivityAt`, `LastActivityKind`, `LastActivityPhase` are computed per claimed REQ. Only the HTML board reads them.
- `skills/do-work-board/tools/queue-kanban/generate.go:215`: `p50_active_minutes` has "one client consumer", the UR summary's remaining-time arm (`web/board-user-request-summary.js:128`). There is no per-REQ ETA anywhere.
- `generate.go:142-148`: the PendingReady / PendingWaiting / PendingEarmarked partition. `board.md:92`: the Needs input · Blocked explanation.
- `skills/do-work-board/tools/queue-kanban/verify.go:80`: `staleClaimThreshold = 3 * time.Hour`. `verify.go:814-819` (report cited 861-868): "reported, not judged dead". `verify.go:1123-1124` (report cited 1170-1171): REQ-073 (the decision to ship no liveness machinery) "rules out a lock, heartbeat, PID check, mtime heuristic and time threshold alike".
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go:375` (`addStuckWorkFinding`): STUCK-WORK after 1 hour, `next_argv` is `git log --full-history -- <path>`, an inspection. `actions/forensics.md:54` points takeover to Crash Recovery.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go:47`: `CommandFinding` has `code`, `severity`, `affected_ids`, `affected_paths`, `observed_evidence`, `fixability`, `automation_stop_reason`, `next_argv`, `next_just_recipe`, `verification_argv`.
- `fan-out-reference.md:142`: the run `manifest.md` holds "REQ id → builder, `<operative_name>`, handback file, landed status, held dispatch instant".
- `actions/work-reference.md:15`: claims travel "with no lock, no lease, and nothing to acquire". `work-reference.md:435` (report cited 526): "The checkpoint grants no lock and no liveness claim." `work-reference.md:130`: `p50_active_minutes` is display only.
- `skills/do-work-board/tools/queue-kanban/walk.go:192` still prunes `runs`.
## Detailed Requirements
P1, `do-work-cli run-status` (read-only):
1. Pick the newest `do-work/runs/work-*` directory unless `--run <dir>` is given. `--req REQ-NNN` limits output to one REQ.
2. Emit one typed record for every claimed REQ and for every needs-input, blocked or waiting REQ, with: `status`, `claimed_at`, `blocked_by`, `blocked_at`, `depends_on`, `assigned_to`, open pending-answers.
3. Column and placement reason, reusing the board's partition and its explanation text.
4. `last_activity_at`, its kind and phase, reusing the board's activity correlation, plus minutes since.
5. From the run manifest: dispatch instant, hand-back file present or absent, landed status.
6. Builder branch tip age (`git log -1 --format=%ct` on the `worktree-agent-REQ-NNN-*` branch).
7. `p50_active_minutes`, minutes since claim, and `eta_minutes` = p50 minus elapsed. When elapsed passes p50, print "over estimate by N min", never a negative ETA. Display only.
8. Run-local files the suite does not define (lock files, progress logs, status files) as a plain list: path, age, first line, labelled "run-local, not interpreted".
9. One `class` and one `next_argv` per row, from this table, in the same `CommandFinding` shape doctor emits, so forensics and status share one result model:

| Code | Class | Shown when | Remedy |
| --- | --- | --- | --- |
| C1 | progressing | activity in the last 20 min | none |
| C2 | quiet | no new stamp or commit for 20 min or more (display boundary only, authorizes nothing) | action layer checks the builder (P2) |
| C3 | hand-back landed, not integrated | hand-back file exists, REQ still in `working/` | `do-work run` |
| C4 | needs operator | `pending-answers` or `blocked` with every dependency met | `do-work clarify` |
| C5 | waiting on dependencies | unmet `depends_on` | name the REQ it waits on |
| C6 | earmarked | `assigned_to` set | quote the board reason; `do-work run REQ-NNN` to run it by name |
| C7 | claim past 3 h | same rule as `verify.go` `staleClaimThreshold` | `do-work run-with-recovery REQ-NNN`, only if the user knows the run is gone |
| C8 | finalization pending | doctor or recover reports an unfinished finalization | the emitted `next_argv` |

P2, `actions/status.md`:
10. Run P1, then add only what the CLI must not do, each item labelled "this machine, this session, now".
11. For C2 rows owned by agents this session spawned: ListAgents and the subagent transcript's modified time, then a one-line nudge text the user can paste or approve.
12. `ps` for test runners and gate processes named in a run-local lock file's owner line. If the owner is not found, print the lock path, its age, and the `rm` command for the user to run. The action never deletes it, because the same lock may belong to another checkout.
13. Output: one table (REQ with a three-to-six-word gloss, phase, minutes since activity, ETA, class), then one remedy line per row that is not C1. The last line is the single most useful next command.

P3, routing:
14. Add a row to `skills/do-work/SKILL.md` before the clarify and roadmap rows: `status`, `status and ETA`, `eta`, `is it stuck`, `stuck`, `why is REQ-N taking so long`, `how do I unblock`, `why is REQ-N in` → `./actions/status.md`. Leave the bare word `blocked` on clarify. Status names clarify as the remedy when that is the cause.

P4, `--watch`:
15. Same table, compact, for a `/loop` body. Document the exact line in the action: `/loop 15m /do-work status --watch`. No state is written between ticks.

P5, `do-work status --fix` (action layer only; `run-status` has no `--fix`):
16. Run only remedies that are already canonical CLI commands from the row's `next_argv`, one at a time, after the user confirms each. Never delete a lock file, never reset a claim.

All:
17. Document the action in the user-facing docs the way sibling actions are documented, and release per `_dev/primes/prime-releases.md`.
## Constraints
- Out of scope (report, Out of scope): a lock, heartbeat, lease, PID check or process registry in shipped Go code (REQ-073); shipping the consumer's lock helper scripts; any change to scheduling, selection, `advance`, the 3-hour threshold, or what `estimate` means; new REQ fields or statuses; a board UI change. The HTML board keeps `last activity`; this adds the terminal surface only.
- Ages are reported; nothing is judged dead. Process checks happen only in the action layer, on this machine, labelled as such.
- `grep -nE 'os\.Remove|syscall\.Kill|FindProcess'` over the new `run-status` Go files finds nothing.
- The original ask behind P5 (release a lock whose owner process is dead) needs a PID check in shipped code, which REQ-073 forbids. If the maintainer wants it, it is its own REQ that revisits REQ-073 by name. Not part of this REQ.
## Assumptions (recorded at capture, no questions asked)
- **One REQ.** The maintainer asked for one REQ per numbered item in the Request section, and one REQ when it has one ask. The Request section has one ask (the action plus its subcommand); P1 to P5 sit in Proposed direction and are this REQ's parts. Capture's challenge, recorded here and not acted on: this is a large REQ, and P5 (`--fix`) and P4 (`--watch`) are the most separable parts. If the maintainer wants a split, `do-work abandon` this REQ and recapture P1 to P5 as five.
- **Cross-module reuse.** `run-status` lives in `do-work-cli` (the `skills/do-work` module), while the board partition, its placement text and activity correlation live in `queue-kanban` (the `skills/do-work-board` module, a required sibling per `suite/modules.tsv`). The two are separate Go modules. "Reuse" means one source of truth: the builder may call the board tool's machine output, or move the shared logic to a place both can use. Never two hand-maintained copies of the same predicate (lessons family `paired-predicate-drift`). If reuse forces a `queue-kanban` change, the board's versioning rules in `_dev/primes/prime-kanban-board.md` apply.
- **Activity correlation decision is open.** The maintainer has not decided whether to keep the board's activity correlation code or replace part of it (the ai-report for REQ-632, activity correlation determinism, holds the two options). `run-status` reads whatever the board computes and adds no second correlation, so either outcome flows through.
- **Class precedence.** Each row gets exactly one class. Suggested order, first match wins: C8, C3, C4, C5, C6, C7, C2, C1. The builder may change it and must document the order in the action.
- **Thresholds.** The 20-minute boundary is approximate, per the maintainer's rule that numeric targets are approximate unless stated strict. C7 reads the same 3-hour value as `staleClaimThreshold`, from one definition, not a copied constant.
- **Rows without an argv.** C1 and C5 have no command remedy. Their `next_argv` stays empty; C5's evidence names the REQ it waits on.
- **CommandFinding shape.** Each row is emitted as a `CommandFinding` with `code` set to the class code and the remedy in `next_argv`. Extra typed fields (ETA, last activity, manifest facts) go into `observed_evidence` or a sibling field in the same JSON result. The builder chooses; the shared struct changes only if that is the smaller design.
- **No run directory.** When `do-work/runs/` has no `work-*` directory, `run-status` still reports queue-side rows and marks run-manifest fields absent. It never fails for a missing run.
- **Builder branch.** A missing `worktree-agent-REQ-NNN-*` branch leaves the field absent. More than one match uses the newest tip and lists the others.
- **`--watch` does not loop.** It changes the output shape only (compact, at most 20 lines for 8 open REQs). Repetition comes from `/loop` or the harness scheduler. No sleep, no state between ticks.
- **The earlier run-directory suggestion.** The report says this REQ "can absorb" the two probes from a 2026-08-23 suggestion (a brief with no hand-back, a stalled run manifest). That suggestion is not in this repo's inbox or queue. Absorbed means `run-status` shows those facts (hand-back absent, dispatch age). `walk.go:192` and board verify stay unchanged, because a board change is out of scope.
- **Route collision.** `do-work-board`'s `status` row (`board.md:33`) stays on board summary. The new row is in `skills/do-work/SKILL.md`, a different skill, so `do-work status` and the board's `status` do not compete.
- **Lock `rm` line.** The action prints the command with the real path filled in. It never runs it.
## Dependencies
No `depends_on` edge. Related queued work: REQ-664 (the coordinated-run stall check) says it calls this run-status subcommand once it exists and reads progress logs inline until then. REQ-663 (the coordinated-run preflight) also looks at stale locks. Neither blocks this REQ, and this REQ does not block them.
## Builder Guidance
Medium-high certainty on what to show; the report gives the class table and acceptance checks. Latitude: Go package layout, how the cross-module reuse is done (see Assumptions), class precedence, text formatting, and where the extra typed fields live in the JSON. Keep the action short: it runs the CLI and adds only the session-local checks.
## Red-Green Proof
**RED prompt/case:** In a fixture repo with one claimed REQ carrying `estimate.p50_active_minutes: 60`, claimed 45 minutes ago, run `do-work-cli run-status --format text`. Ask a session "status and ETA please".
**Why RED now:** `run-status` does not exist, and "status and ETA" routes nowhere (`SKILL.md` has no row for it).
**GREEN when:** (from the report's Acceptance check)
- `do-work status`, "status and ETA please" and "is REQ-N stuck?" route to `actions/status.md`; "blocked" still routes to clarify.
- The fixture prints an ETA of about 15 min. At 75 minutes elapsed it prints "over estimate by 15 min".
- A claimed REQ whose hand-back file exists in the run directory is C3 with remedy `do-work run`.
- A blocked REQ with an unmet `depends_on` is C5 and quotes the board's Waiting reason. An earmarked REQ is C6 and quotes the Earmarked reason.
- A `full-gate.lock` in the run directory is listed as run-local with its age and first line. `run-status` never deletes it, and `--fix` never deletes it.
- `grep -nE 'os\.Remove|syscall\.Kill|FindProcess'` over the new `run-status` Go files finds nothing.
- The JSON output uses the `resultmodel.CommandFinding` shape doctor emits.
- `--watch` output fits in 20 lines for 8 open REQs.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: a new do-work-cli subcommand that emits `next_argv` (families `destructive-next-argv`, `time-is-not-authorization`).
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8334 tokens, over budget; `slugged: partial`). Matching reason: reuses the board partition and activity correlation (families `paired-predicate-drift`, `git-history-evidence`).
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over budget; `slugged: partial`). Matching reason: adds an action and a routing row.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: the action prescribes `ps` and prints an `rm` command.
## Full Context
See `do-work/user-requests/UR-147/input.md` for complete verbatim input. The source report stays at `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-status-action.md`. No queued candidate shares this root cause (the queue held REQ-654 to REQ-667 at capture; REQ-664 names this subcommand as a separate companion request).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-status-action.md`, Request: "Please add one read-only action, `do-work status`, backed by one deterministic subcommand, `do-work-cli run-status [--run <dir>] [--req REQ-NNN] [--watch]`, with `--format json|text`."*

## Cancelled

- **When:** 2026-10-10T12:59:20Z
- **Why:** folded into the simplified recapture of 2026-10-10 (maintainer chose the aggressive simplification)
- **Decided by:** user, via `do-work abandon`
