---
id: REQ-690
title: 'do-work status action and do-work-cli run-status report one class, ETA and remedy per open REQ'
status: claimed
created_at: 2026-10-10T13:05:57Z
user_request: UR-153
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-10-10T12:52:21Z
status_changed_at: 2026-10-10T13:14:51Z
related: [REQ-689, REQ-691]
route: C
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work-board/tools/queue-kanban/model.go", "skills/do-work-board/tools/queue-kanban/open_work.go", "skills/do-work-board/tools/queue-kanban/main.go", "skills/do-work-board/tools/queue-kanban/open_work_test.go", "skills/do-work-board/actions/board.md", "skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go", "skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go", "skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go", "skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go", "skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md", "skills/do-work/actions/status.md", "skills/do-work/SKILL.md", "skills/do-work/actions/help.md", "skills/do-work/docs/status-guide.md", "skills/do-work/actions/fan-out-reference.md"]
planning_at: 2026-10-10T13:23:24Z
estimate:
  p50_active_minutes: 50
  confidence: low
  basis:
  - Route C
  - 16-file write set
  - 3 subsystems involved
  - 8 acceptance criteria
  calculated_at: 2026-10-10T13:22:28Z
builder_handback_at: 2026-10-10T13:42:46Z
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

All:
17. Document the action in the user-facing docs the way sibling actions are documented, and release per `_dev/primes/prime-releases.md`.
## Constraints
- Out of scope (report, Out of scope): a lock, heartbeat, lease, PID check or process registry in shipped Go code (REQ-073); shipping the consumer's lock helper scripts; any change to scheduling, selection, `advance`, the 3-hour threshold, or what `estimate` means; new REQ fields or statuses; a board UI change. The HTML board keeps `last activity`; this adds the terminal surface only.
- Ages are reported; nothing is judged dead. Process checks happen only in the action layer, on this machine, labelled as such.
- `grep -nE 'os\.Remove|syscall\.Kill|FindProcess'` over the new `run-status` Go files finds nothing.
- The source report's `--fix` part (run a row's remedy after confirmation) is dropped by the maintainer's 2026-10-10 decision: the remedy line is already printed and the user can paste it. Its original ask (release a lock whose owner process is dead) needs a PID check in shipped code, which REQ-073 forbids; if wanted, it is its own REQ that revisits REQ-073 by name.
## Assumptions (recorded at capture, no questions asked)
- **One REQ.** The source report had one ask (the action plus its subcommand) in five parts; the maintainer kept it as one REQ on 2026-10-09 and on 2026-10-10 dropped the fifth part (`--fix`). P1 to P4 remain.
- **Cross-module reuse.** `run-status` lives in `do-work-cli` (the `skills/do-work` module), while the board partition, its placement text and activity correlation live in `queue-kanban` (the `skills/do-work-board` module, a required sibling per `suite/modules.tsv`). The two are separate Go modules. "Reuse" means one source of truth: the builder may call the board tool's machine output, or move the shared logic to a place both can use. Never two hand-maintained copies of the same predicate (lessons family `paired-predicate-drift`). If reuse forces a `queue-kanban` change, the board's versioning rules in `_dev/primes/prime-kanban-board.md` apply.
- **Activity correlation is settled.** The maintainer decided on 2026-10-10 (shipped 0.305.88): the merge-range scaffold was deleted and nothing stamps commits; the board's last-activity correlation stays as code. `run-status` reads whatever the board computes and adds no second correlation.
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
No `depends_on` edge. Related queued work: REQ-689 (do-work run --coordinate) arms a stall loop that reads `do-work status --watch` once it exists and reads progress logs until then; REQ-691 (do-work trace) and REQ-692 (validate-feedback --capture) add rows to the same routing table. None blocks this REQ.
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
- A `full-gate.lock` in the run directory is listed as run-local with its age and first line. `run-status` never deletes it.
- `grep -nE 'os\.Remove|syscall\.Kill|FindProcess'` over the new `run-status` Go files finds nothing.
- The JSON output uses the `resultmodel.CommandFinding` shape doctor emits.
- `--watch` output fits in 20 lines for 8 open REQs.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (19046 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: a new do-work-cli subcommand that emits `next_argv` (families `destructive-next-argv`, `time-is-not-authorization`).
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8815 tokens, over budget; `slugged: partial`). Matching reason: reuses the board partition and activity correlation (families `paired-predicate-drift`, `git-history-evidence`).
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over budget; `slugged: partial`). Matching reason: adds an action and a routing row.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: the action prescribes `ps` and prints an `rm` command.
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (6122 tokens, over budget; `slugged: partial`). Matching reason: adds a machine-readable mode to `queue-kanban open-work` (added at claim time by the pre-dispatch consult).
## Full Context
See `do-work/user-requests/UR-153/input.md` for the 2026-10-10 decision record. The cancelled original, `do-work/archive/UR-147/REQ-668-run-status-action.md`, carries the complete body; the source report stays at `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-status-action.md`, with verbatim input in `do-work/archive/UR-147/input.md`.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-status-action.md`, Request: "Please add one read-only action, `do-work status`, backed by one deterministic subcommand, `do-work-cli run-status [--run <dir>] [--req REQ-NNN] [--watch]`, with `--format json|text`."*

---

## Triage

**Route: C** - Complex

**Reasoning:** Four parts across three subsystems (a new do-work-cli Go command, a machine-readable mode on the separate queue-kanban module, and a new action plus routing and docs), with one open architectural choice the REQ leaves to the builder: how a standard-library-only do-work-cli reuses facts computed in a different Go module. That choice is settled in the Plan below so the builder does not reopen it.

**Planning:** Required

## Plan

**Architecture (decided at pre-dispatch, do not reopen).** do-work-cli is standard-library only and may not import the separate `queue-kanban` module (lessons-do-work-cli.md § Package direction); queue-kanban depends on goldmark and yaml.v3. So the two modules cannot share Go code, and the board's partition and activity correlation must not be copied (Assumptions: one source of truth). The reuse is a file handoff:

1. **Board side (queue-kanban), one machine-readable mode.** `queue-kanban open-work --format json --repo-root R` prints one JSON document; `--format text` (default) keeps today's digest byte for byte. The document carries, for every open ticket the board buckets (pending ready, pending waiting, pending earmarked, claimed, needs input or blocked): id, title, status, column, placement reason, unmet dependencies, assigned_to, and for claimed tickets the board's own `last_activity_at`, `last_activity_kind`, `last_activity_phase` (computed by calling the existing `attachRequestActivity` over one `readWorktreeAgentGitState` read, exactly as `attachVerifyFindingsAndRequestActivity` does, never a second correlation). Top level: `generated_at` and `stale_claim_threshold_minutes` read from `staleClaimThreshold` (one definition). The placement reason is set inside the same `bucketColumns` switch arm that picks the column (a new `RequestTicket.PlacementReason` string), so column and reason cannot drift apart; its sentences restate `skills/do-work-board/actions/board.md`'s Needs input · Blocked paragraph and the Earmarked tooltip in plain words. No web/ change, no payload (`generatedRequest`) change, no new write surface.
2. **Core side (do-work-cli), one new package and command.** `internal/runstatus/` registers `run-status` in `cmd/do-work-cli/main.go`. Flags: `--board-facts FILE` (required; the JSON from step 1), `--run DIR` (default: the newest `do-work/runs/work-*` directory by name; none is not an error), `--req REQ-NNN`, `--watch`, `--now RFC3339` only if the builder needs it for tests (prefer an injected clock in the package). It reads the repository with `repositorymodel.DiscoverRepository` for frontmatter facts (status, claimed_at, blocked_by, blocked_at, depends_on, assigned_to, unchecked `- [ ]` items under `## Open Questions` for pending-answers) and the frozen `estimate.p50_active_minutes` through the same parse `nextselection.frozenEstimate` uses (export it in place as `FrozenEstimate`; never a second parser). Run directory facts: hand-back present when `RUN/REQ-NNN-handback.md` exists; the manifest row is the `manifest.md` table row whose first cell is the REQ id, kept verbatim, with the first UTC ISO-8601 instant in that row taken as the dispatch instant (absent when none; skip and report, no column grammar). Builder branch: `git for-each-ref --sort=-committerdate --format=%(refname:short) %(committerdate:iso-strict) refs/heads/worktree-agent-REQ-NNN-*`; newest tip used, others listed; absent when none. Run-local files: every entry of the run directory that is neither `manifest.md` nor named `REQ-NNN-*` (a condition, not a list), each with path, modification age and first line, labelled "run-local, not interpreted". Unfinished finalization: `doctor.FinalizationTailFindings` for the REQ id (exported already, read-only).
3. **Result shape.** One `resultmodel.CommandFinding` per row: `code` = the class code (`C1` … `C8`), `severity: info`, `affected_ids` = the REQ, `affected_paths` = its file, `observed_evidence` = short plain lines (column and placement reason, ages, ETA), `next_argv` = the remedy argv or empty. The typed per-row fields and the run-local list go in one new sibling field `run_status` on `resultmodel.CommandResult` (the pattern `audit_metrics`, `checkpoint` and the others already follow); `CommandFinding` itself does not change. Outcome `success` (exit 0) whenever the report was produced; `refused` only for a usage error or an unreadable `--board-facts` file. Text format uses `ExactTextOutput`: one table (REQ, title gloss, phase, minutes since activity, ETA, class), one remedy line per non-C1 row, last line the single most useful next command. `--watch` is the same data, compact, at most 20 lines for 8 open REQs, no state written.
4. **Classes, first match wins: C8, C3, C4, C5, C6, C7, C2, C1** (the REQ's suggested order). C4 = the board column is needs-input-or-blocked; C5 = column pending-waiting (evidence names the unmet REQ ids and quotes the placement reason); C6 = column pending-earmarked (quotes the placement reason); C7 = claimed and minutes since claim ≥ `stale_claim_threshold_minutes`; C2 = claimed with last activity 20 min or more ago (approximate display boundary); C1 = claimed and active within 20 min. A claimed REQ with no activity record falls to C2 with evidence "no activity record", never C1. ETA = p50 minus minutes since claim; when elapsed passes p50 the row says "over estimate by N min", never a negative number; no estimate prints "no estimate".
5. **Remedies (`next_argv`).** C1, C2, C5: empty. C3: `do-work run`. C4: `do-work clarify`. C6: `do-work run REQ-NNN`. C8: the doctor finding's own `next_argv`. **C7: the read-only inspection `git log --full-history -- <REQ path>` (doctor STUCK-WORK's argv), with `do-work run-with-recovery REQ-NNN` named only in `automation_stop_reason` and in the text remedy line, together with what it does (requeues the claim and strips generated sections) and the condition "only if you know the run that claimed it is gone".** See D-01.
6. **Action and routing.** New `skills/do-work/actions/status.md` (prime-action-files template): build the board binary the way `skills/do-work-board/actions/board.md` Step 4 does, write its `open-work --format json` output to a `mktemp` file, run `do-work-cli run-status --board-facts <file>` with the user's `--run`, `--req` and `--watch`, then add only the session-local checks P2 names (items 11 and 12), each labelled "this machine, this session, now"; the lock `rm` line is printed with the real path and never run; document `/loop 15m /do-work status --watch` and the class precedence. If that shell exceeds about five lines, move it into one committed script with a one-line launcher. `skills/do-work/SKILL.md`: one routing row placed directly above the clarify row, which puts it above both the clarify and roadmap rows as item 14 asks, triggers per item 14, `blocked` stays on clarify; add `status [REQ] [--watch]` to the argument hint. Docs: one `skills/do-work/actions/help.md` menu line, a short `skills/do-work/docs/status-guide.md` like the roadmap and forensics guides, one `do-work-board/actions/board.md` mention of `open-work --format json`, and one clause at `skills/do-work/actions/fan-out-reference.md:167` pointing progress questions at `actions/status.md`. One Package-routing bullet for `internal/runstatus/` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`.
7. **Testing approach (TDD, RED first).** do-work-cli: table tests over a temp fixture repo plus a hand-written board-facts JSON and an injected clock: ETA 15 at 45 of 60 min, "over estimate by 15 min" at 75; C3 when the hand-back file exists; C5 with unmet depends_on quoting the placement reason; C6 earmarked; a `full-gate.lock` listed as run-local with age and first line and still present after the run; JSON finding shape; `--watch` ≤ 20 lines for 8 open REQs; no run directory still reports queue rows. queue-kanban: one test through `buildBoard` that the JSON mode puts a blocked REQ with an unmet dependency in pending-waiting with its reason and an earmarked one in pending-earmarked (paired-predicate-drift lesson: pin through the full board build), and that the text digest is unchanged.

**Pre-dispatch decisions (recorded, the builder continues from D-05):**
- D-01: C7's `next_argv` is the read-only `git log --full-history -- <path>`, not `do-work run-with-recovery REQ-NNN`. Reasoning: lessons-do-work-cli.md [family: destructive-next-argv] (REQ-629) says a finding's `next_argv` is followed literally and must never be the destructive command; run-with-recovery takes over the claim. The REQ's remedy text still appears, in the stop reason and the text remedy line. Value: no automation can reset a live claim from a status row. Risk: one extra step for a user who knows the run is gone; reversible.
- D-02: `--board-facts FILE` is required and the action supplies it; do-work-cli never builds or executes the board itself. Reasoning: the core binary can run prebuilt with no Go toolchain and does not know where the sibling package is installed, while the action resolves sibling paths literally (SKILL.md). Value: no new cross-module process launch in Go. Risk: `run-status` alone, without the action, needs one extra command; a host without Go cannot build the board and gets the board action's own "Go toolchain required" stop.
- D-03: run-local means "not `manifest.md` and not `REQ-NNN-*`" in the run directory. Reasoning: a condition, not a list (CLAUDE.md, closed enumerations go stale); REQ-689 makes `full-gate.lock` a documented run-directory file but its name still is not REQ-prefixed, so it stays listed, matching the GREEN check. Value: no list to maintain. Risk: coordinator guide files (for example PREDISPATCH-GUIDE.md) also appear as run-local lines; harmless, and labelled not interpreted.
- D-04: rows are the board's open tickets minus pending-ready (claimed, needs-input-or-blocked, pending-waiting, pending-earmarked). Reasoning: the REQ's class table has a class for each of those and none for a ready pending REQ, which is neither stuck nor waiting. Value: the table stays short. Risk: a user who wants the whole queue uses `do-work roadmap`.

**Plan validation:** every Detailed Requirement maps to a step (P1 items 1 to 9 → steps 1 to 5; P2 items 10 to 13 → step 6; P3 item 14 → step 6; P4 item 15 → steps 3 and 6; item 17 docs → step 6, the release is the integrator's). No orphan step. Seven steps, above the three-task comfort line, but each is small and they are one feature; splitting would put the board mode and its only reader in different releases. Consumer field contract: the action consumes `run_status.rows[].class`, `next_argv`, `run_status.run_local_files[]` (path, age, first line), all named above.

*Generated by the pre-dispatch agent (planning done inline, not by a separate Plan agent)*

## Exploration

Read at main `bd56c4b0` by the pre-dispatch agent (inline, not a separate Explore agent).

**Lessons consult.** `required_lessons` now holds `_dev/primes/lessons-releases.md` (666 tokens, the only matching satellite under budget; it applies to the integrator's release). Every other match is a whole `slugged: partial` satellite over the 2000-token budget and stays in the Dropped section above; `_dev/primes/lessons-kanban-board.md` was added there. The pre-dispatch agent read the families the capture named and copies the three that change the design into the builder brief: [family: destructive-next-argv] (lessons-do-work-cli.md:87, REQ-629: a finding's `next_argv` must never be the destructive command; drives D-01), [family: git-history-evidence] (lessons-do-kanban.md:22, branch tips count only when `--no-merged` against the integration ref, which the board's read already does), [family: paired-predicate-drift] (lessons-do-kanban.md:18 and :61: pin a placement through the full `buildBoard`, and grep every restatement when a label's meaning changes).

**Board module (`skills/do-work-board/tools/queue-kanban/`, own `go.mod`, deps goldmark + yaml.v3).**
- `main.go:66-87` hand-rolled subcommand switch; `runOpenWorkCommand` at `main.go:147` parses only `--repo-root` and calls `writeOpenWorkDigest`. The synopsis comment at `main.go:23-32` lists each subcommand's flags and must gain `--format`. Every subcommand rejects leftover tokens (`exitOnLeftoverArguments`).
- `open_work.go:57-80` `writeOpenWorkDigest` (text digest, assertable writer split from the command); `open_work_test.go` holds the digest tests (`TestOpenWorkDigest…`), the style to copy.
- `model.go:1718-1765` `bucketColumns`: the partition. Pending holds `pending` plus `blocked` with unmet dependencies; inside it, unmet or non-pending → PendingWaiting, then `AssignedTo != ""` → PendingEarmarked, else PendingReady. `claimed` → Claimed; `isNeedsInputOrBlockedStatus` → NeedsInputOrBlocked; unrecognized status → NeedsInputOrBlocked with a warning. `RequestTicket` (`model.go` about :184-300) has `UnmetDependencies` (:203), `AssignedTo` (:193), `ClaimedAt`. No Go-side placement text exists today: the explanations live in `actions/board.md:92` (Needs input · Blocked) and the client tooltip at `web/board-cards.js:234-240` (Earmarked). Web files are out of scope.
- `activity_correlation.go:262-298` `attachRequestActivity(data, board, now, runner, gitState)` fills `data.RequestActivity[id]` (`lastActivityAt`, `lastActivityKind`, `lastActivityPhase`) for claimed and recently-done tickets; `generate.go:682` `attachVerifyFindingsAndRequestActivity` shows the call with `readWorktreeAgentGitState(repoRoot, runGitCommand)` (`verify.go:1417`). Calling `attachRequestActivity` on a zero `generatedBoardData` and reading `RequestActivity` reuses it with no refactor.
- `verify.go:80` `const staleClaimThreshold = 3 * time.Hour` (comment: reports, never authorizes). `verify.go:1401-1440` `worktreeAgentGitState` already reads owned builder-branch tips with `for-each-ref --no-merged`, by REQ id but without branch names.
- No machine-readable output exists on any board subcommand today; `do-work-cli` never invokes the board (only comments mention it).

**Core CLI (`skills/do-work/tools/do-work-cli/`, standard library only).**
- `cmd/do-work-cli/main.go:25-80`: registration is one import plus one `for name, handler := range pkg.Handlers()` loop.
- `internal/doctor/doctor_commands.go:19-42`: handler shape (`parseCommandOptions`, `repositorymodel.DiscoverRepository`, injected now). `internal/doctor/doctor_scan.go:96` exported `FinalizationTailFindings(ctx, snapshot)` (C8 source). `doctor_scan.go:375-404` `addStuckWorkFinding` (STUCK-WORK after 1 h, `next_argv` = `git log --full-history -- <path>`; the C7 read-only argv to reuse).
- `internal/resultmodel/result_model.go:47-58` `CommandFinding`; `:609-644` `CommandResult` with optional typed sibling fields and `ExactTextOutput *string` (text override, used by `commandruntime.writeResult`, `internal/commandruntime/command_runtime.go:95`). Outcome `findings` exits 1, `success` exits 0 (`ExitCode`, `:648`).
- `internal/nextselection/next_selection.go:471-480` `frozenEstimate(fields)` reads `estimate.p50_active_minutes` (unexported; the one core parse of that field). `next_targets.go:179` precedent for action argv in `next_argv`: `["do-work", "roadmap"]`.
- `internal/cleanup/cleanup_git.go:230,352` and `internal/corehelpers/handoff.go:55` already list `worktree-agent-*` branches (prefix `worktree-agent-REQ-`).

**Run directory (`skills/do-work/actions/fan-out-reference.md:138-145`).** Suite-defined names: `manifest.md`, `REQ-NNN-brief.md`, `REQ-NNN-handback.md`, `REQ-NNN-integrate.md`. Real manifests (`do-work/runs/work-2026-10-10-100748/manifest.md`) are a Markdown table `| REQ | Builder | Operative name | Handback file | Status | Dispatch instant |`, but the column set is prose-defined and varies between runs, so only the REQ cell and an ISO instant are safe to read. Run directories are named `work-YYYY-MM-DD-HHMMSS` and sort by name; `do-work/runs/` also holds a non-`work-` directory (`REQ-531-review-findings-critical-only`).

**Routing and docs.** `skills/do-work/SKILL.md:4` argument hint; `:29-47` routing table; clarify row `:38` carries `blocked`; forensics `:42`; roadmap `:43`. `skills/do-work-board/actions/board.md:33` routes the board's own `status` to summary (a different skill; no collision). `actions/help.md:7-31` full menu. `docs/roadmap-guide.md` (51 lines) and `docs/forensics-guide.md` are the guide pattern. `actions/fan-out-reference.md:167` "It asks about progress or state. Answer from disk…" has no command behind it.

**Concerns.**
- Sibling seams in this run: REQ-691 (trace) edits the same `SKILL.md:4` argument-hint line and adds a row above capture; REQ-692 adds a row above capture; REQ-689 adds a row above the run row and edits `fan-out-reference.md` (run-directory table, coordinator rules). REQ-660, REQ-659 and REQ-658 add bullets to `lessons-do-work-cli.md`. Keep each edit local.
- REQ-689 documents `full-gate.lock` and `REQ-NNN-progress.log` as run-directory files. D-03's condition still lists the lock as run-local; the progress log is REQ-prefixed and is not listed. If REQ-689 lands first, status.md may cite its run-directory table for the lock's owner line.
- The guide's sibling table names `work-reference.md` for this REQ; no requirement needs it, so it is not in Scope.

*Generated by the pre-dispatch agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — RequestTicket.PlacementReason, set inside the bucketColumns arms
- `skills/do-work-board/tools/queue-kanban/open_work.go` (modify) — JSON renderer for the open-work facts
- `skills/do-work-board/tools/queue-kanban/main.go` (modify) — the format flag on open-work and the synopsis comment
- `skills/do-work-board/tools/queue-kanban/open_work_test.go` (modify) — JSON placement test through the board build; text digest unchanged
- `skills/do-work-board/actions/board.md` (modify) — one mention of the JSON mode in the open-work bullet
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go` (new) — the run-status command, classes, remedies, text and watch rendering
- `skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go` (new) — the RED/GREEN fixture tests
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify) — register the package's handlers
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — the typed run_status sibling field
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modify) — export the one estimate parse in place as FrozenEstimate
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modify) — one Package routing bullet
- `skills/do-work/actions/status.md` (new) — the action
- `skills/do-work/SKILL.md` (modify) — routing row above clarify, argument hint
- `skills/do-work/actions/help.md` (modify) — one menu line
- `skills/do-work/docs/status-guide.md` (new) — short user guide
- `skills/do-work/actions/fan-out-reference.md` (modify) — one clause at the progress-question rule

**Files I will NOT touch:** `skills/do-work-board/tools/queue-kanban/web/` (board UI out of scope), `generate.go` payload structs, `verify.go` (threshold read, not changed), `walk.go`, `activity_correlation.go` (called, not changed), `actions/work.md`, `actions/work-reference.md`, `actions/roadmap.md`, `actions/forensics.md`, `actions/clarify.md`, any `CHANGELOG`, `VERSION` or version mirror (integrator), anything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] `do-work status`, "status and ETA please" and "is REQ-N stuck?" route to `actions/status.md`; "blocked" still routes to clarify.
- [ ] A claimed REQ with p50 60 claimed 45 min ago prints an ETA of about 15 min; at 75 min it prints "over estimate by 15 min".
- [ ] A claimed REQ whose hand-back file exists in the run directory is C3 with remedy `do-work run`.
- [ ] A blocked REQ with an unmet `depends_on` is C5 and quotes the board's Waiting reason; an earmarked REQ is C6 and quotes the Earmarked reason.
- [ ] A `full-gate.lock` in the run directory is listed as run-local with its age and first line; `run-status` never deletes it.
- [ ] `grep -nE 'os\.Remove|syscall\.Kill|FindProcess'` over the new run-status Go files finds nothing.
- [ ] The JSON output uses the `resultmodel.CommandFinding` shape doctor emits.
- [ ] `--watch` output fits in 20 lines for 8 open REQs.
