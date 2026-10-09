---
id: UR-147
title: 'do-work status action and do-work-cli run-status report class, ETA and one remedy per open REQ'
created_at: 2026-10-09T21:20:59Z
requests: [REQ-668]
word_count: 2090
---
# do-work status: Activity Age, ETA and One Remedy per Open REQ

## Summary
The maintainer ran `/do-work capture-request` on the upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-status-action.md` and asked for one UR with one REQ per numbered item in its Request section, or one REQ when the section has one ask. Capture only, no questions, every ambiguity recorded as an assumption in the REQ. The report file is the verbatim input below, byte for byte, and it stays in the inbox as the source.

The Request section has one ask: one read-only action, `do-work status`, backed by one deterministic subcommand, `do-work-cli run-status`. Its five parts, P1 to P5, sit in the report's Proposed direction section, so they became the requirement groups of one REQ, REQ-668 (the status action and run-status subcommand). The REQ records the challenge that this is large and that P4 and P5 are the most separable parts.

The report was observed against 0.305.84. The capture spot-checked its citations at 0.305.87. These match: `SKILL.md:38`, `:42` and `:43`, `actions/roadmap.md:5`, `skills/do-work-board/actions/board.md:33` and `:92`, `open_work.go:79`, `activity_correlation.go:40-42`, `generate.go:142-148` and `:215`, `web/board-user-request-summary.js:128`, `verify.go:80`, `walk.go:192`, `doctor_scan.go:375`, `actions/forensics.md:54`, `actions/run-with-recovery.md:10`, `actions/work-reference.md:15` and `:130`. These moved: `verify.go` "reported, not judged dead" is now `:814-819` (report 861-868), the REQ-073 comment is `:1123-1124` (1170-1171), the checkpoint no-lock line is `work-reference.md:435` (526), "Answer from disk" is `fan-out-reference.md:165`, and the run manifest row is `fan-out-reference.md:142`.

No prompt injection found. The report's "How to use this file" line is framing for its intended reader.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-668 | do-work status action and do-work-cli run-status report one class, ETA and remedy per open REQ |

## Batch Constraints
- No new REQ field, no new status, and no lock, heartbeat, PID check or process registry in shipped Go code (REQ-073, the no-liveness-machinery decision). Ages are reported; nothing is judged dead.
- No board UI change; no change to scheduling, selection, `advance`, the 3-hour threshold or what `estimate` means.
- Related queued work, no edge: REQ-664 (coordinated-run stall check) calls this subcommand once it exists.
- Release per `_dev/primes/prime-releases.md`.

## Full Verbatim Input
> ```
> # Upstream suggestion for `knews2019/skill-do-work` — a `do-work status` action: activity age, ETA and one remedy per open REQ
> 
> **How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
> clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
> consumer repo it was authored in. Observed against **v0.305.84**; all line numbers are from that tag.
> 0.305.87 later moved the fan-out, run-directory and mid-run message prose from `actions/work.md` and
> `actions/work-reference.md` into `actions/fan-out-reference.md`. The two citations marked "now in
> `actions/fan-out-reference.md`" below must be re-found there by quoted text, not by line.
> 
> ---
> 
> ## Request
> 
> During a run, users ask the same questions again and again: "status and ETA", "is it stuck", "why does
> REQ-N take so long", "how do I unblock REQ-N", "why is REQ-N in the pending column". No action routes
> those phrases, and no single command answers them. The facts exist, but they are spread over five
> surfaces, and two of them (last activity, remaining time) are only in the HTML board.
> 
> Please add one read-only action, `do-work status`, backed by one deterministic subcommand,
> `do-work-cli run-status [--run <dir>] [--req REQ-NNN] [--watch]`, with `--format json|text`.
> 
> Already shipped, and reused rather than rebuilt:
> 
> - v0.305.67 added `last activity` (newest lifecycle stamp or linked commit, with its phase) to claimed
>   cards, but only on the served board.
> - v0.305.79 to v0.305.84 fixed where earmarked and operator-gated REQs sit, and the board explains the
>   Ready / Waiting / Earmarked split. That answers "why is REQ-N in the pending column" once a user opens
>   the board. The earlier suggestion `2026-10-08_do-work-upstream-suggestion-unblock-keeps-earmark`
>   covered that part.
> - The earlier suggestion `2026-08-23_upstream-run-directory-visibility` asked for two read-only `verify`
>   probes over `do-work/runs/` (a brief with no hand-back, a stalled run manifest). They have not shipped:
>   `walk.go:192` still prunes `runs`. This request can absorb them.
> 
> This asks for the remaining part: one terminal answer per open REQ, with a class, an ETA and one remedy.
> No new REQ field, no new status, and no lock, heartbeat or PID check in shipped Go code (see Out of scope).
> 
> ## What happened, in one consumer repo (evidence)
> 
> All evidence comes from one consumer repo over 2026-10-02 to 2026-10-09. A recount of its prompt log found
> about 20 user prompts in 12 separate sessions on 7 different days, all asking for status, stuck lanes,
> unblocking or a watchdog. In the same 30 days `do-work forensics` was invoked 0 times.
> 
> | Date | What the user did or asked |
> | --- | --- |
> | Oct 2 | Wrote a watchdog prompt by hand: "/loop 15m Run watchdog for do-work run <run>: check the heavy lock (... status) for a stale owner ... nudge any per-REQ agent that has been idle". Two sessions that day ran 24 and 64 Bash calls that only checked a lock file and per-REQ status lines. |
> | Oct 2 to 3 | One coordinator session ran 1005 tool calls. 378 of them were SendMessage nudges to builders. |
> | Oct 4 | "REQ-N <- what is the status, why does it take so long, over 3h? is it stuck?" |
> | Oct 5 | "we have quite a number of tickets that take over 4h now, we should not assume by start that they got paused, and in fact we can check it against the git commit history" (v0.305.67 shipped the board half of this the same day). |
> | Oct 6 | "some of these tasks take a long time, is there any good reason for it? Can we help them out?", then "give me a prompt that will unblock it". A pasted prompt, "Lanes are stalled. Check for a stale lock and dead lanes, then fix both. Do not assume; verify.", appears in three sessions. |
> | Oct 7 | "how to unblock these 2?" and "how to unblock this: REQ-N" |
> | Oct 7 to 8 | "why keep the idle builders?", "status and ETA please", "why does the merge take so long?" |
> | Oct 8 | "why are the earmarked reqs ... in the pending column? why are they not on the `need input - blocked` column?" |
> | Oct 9 | "please provide status and ETA" |
> 
> Each answer was rebuilt from scratch: read the REQ files, read the run manifest, run `git log` on each
> builder branch, `ls -l` the run directory, `ps` for test runners, then write a nudge by hand. The lock
> files the user asked about are not suite artifacts. The consumer's run orchestrators created them
> (`full-gate.lock`, `locks/<name>/owner`) inside the run directory.
> 
> ## Where the behaviour lives today
> 
> **Routing.** No row catches the user's words.
> 
> - `SKILL.md:42`: forensics routes on "`forensics`, `diagnose`, `health check`" only. "stuck", "status",
>   "ETA" and "unblock" are not there.
> - `SKILL.md:43`: "`roadmap`, `queue-status`, `where are we`, `what's left`" goes to roadmap, which says
>   it is "a planning aid, not a diagnostic" (`actions/roadmap.md:5`).
> - `SKILL.md:38`: "`blocked`" routes to clarify, so "how do I unblock REQ-N" lands in the answer flow
>   before anyone has said why the REQ is stuck.
> - `skills/do-work-board/actions/board.md:33`: "`summary`, `status`, `counts`" routes the word `status` to
>   column counts.
> 
> **What each existing surface gives.**
> 
> - `actions/work.md:169` (now in `actions/fan-out-reference.md`): "It asks about progress or state.
>   Answer from disk: the REQ files, the run manifest, the hand-backs." Correct rule, with no command
>   behind it.
> - `skills/do-work-board/tools/queue-kanban/open_work.go:79-89` (`writeClaimedSection`): each claimed REQ
>   is printed as id and title only. No phase, no age, no ETA.
> - `skills/do-work-board/tools/queue-kanban/activity_correlation.go:40-42`: `LastActivityAt`,
>   `LastActivityKind`, `LastActivityPhase` are already computed per claimed REQ. Only the HTML board
>   reads them.
> - `skills/do-work-board/tools/queue-kanban/generate.go:214-217`: the saved `p50_active_minutes` has "one
>   client consumer", the UR progress summary's remaining-time arm
>   (`skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js:128`). There is no per-REQ
>   ETA anywhere, and none in the terminal.
> - `skills/do-work-board/actions/board.md:92`: the Waiting-on-dependencies explanation, and
>   `generate.go:142-148`: the Ready / Waiting / Earmarked partition. This is the placement reason a status
>   row should quote.
> - `skills/do-work-board/tools/queue-kanban/verify.go:861-868`: a claim older than `staleClaimThreshold`
>   (3 hours, `verify.go:80`) is reported "reported, not judged dead", with the remedy "use
>   actions/work-reference.md -> Crash Recovery (Step 1) to decide whether to take it over".
> - `tools/do-work-cli/internal/doctor/doctor_scan.go:375-404` (`addStuckWorkFinding`): STUCK-WORK fires
>   on a claim older than 1 hour whose file was not modified in the last hour. Its `next_argv` is
>   `git log --full-history -- <path>`, an inspection, not a remedy. Forensics then points to Crash
>   Recovery (`actions/forensics.md:54`).
> - `actions/run-with-recovery.md:10`: the takeover verb, for when "Ordinary `run` refused ... or left
>   stale claims for a human ownership decision".
> - `actions/work-reference.md:480` (now in `actions/fan-out-reference.md`): the run directory's
>   `manifest.md` holds "REQ id → builder, `<operative_name>`, handback file, landed status, held dispatch
>   instant". Dispatch time and hand-back presence are durable and readable.
> 
> **The constraint any design must keep.**
> 
> - `actions/work-reference.md:15`: claims travel "with no lock, no lease, and nothing to acquire".
> - `actions/work-reference.md:526`: "The checkpoint grants no lock and no liveness claim."
> - `skills/do-work-board/tools/queue-kanban/verify.go:1170-1171`: REQ-073 (the suite's earlier decision
>   to ship no liveness machinery) "rules out a lock, heartbeat, PID check, mtime heuristic and time
>   threshold alike".
> 
> The deterministic subcommand below therefore reports ages of durable records and never says "dead".
> Process checks happen only in the action layer, on this machine, labelled as such.
> 
> ## Proposed direction
> 
> **P1. `do-work-cli run-status`.** Read-only. Picks the newest `do-work/runs/work-*` directory unless
> `--run` is given. For every claimed REQ, and for every needs-input, blocked or waiting REQ, it emits one
> typed record:
> 
> - `status`, `claimed_at`, `blocked_by`, `blocked_at`, `depends_on`, `assigned_to`, open pending-answers.
> - Column and placement reason, reusing the board's partition and its explanation text.
> - `last_activity_at`, kind and phase, reusing `activity_correlation.go`, plus minutes since.
> - From the run manifest: dispatch instant, hand-back file present or absent, landed status.
> - Builder branch tip age (`git log -1 --format=%ct worktree-agent-REQ-NNN-*`).
> - `p50_active_minutes`, minutes since claim, and `eta_minutes` = p50 minus elapsed. When elapsed passes
>   p50, print "over estimate by N min", never a negative ETA. The ETA is display only, like the field it
>   reads (`actions/work-reference.md:130`).
> - Run-local files the suite does not define (lock files, progress logs, status files) as a plain list:
>   path, age, first line. Labelled "run-local, not interpreted".
> - One `class` and one `next_argv` per row, from the table below. Same `CommandFinding` shape as
>   doctor, so forensics and status share one result model.
> 
> | Code | Class | Shown when | Remedy |
> | --- | --- | --- | --- |
> | C1 | progressing | activity in the last 20 min | none |
> | C2 | quiet | no new stamp or commit for 20 min or more (display boundary only, authorizes nothing) | action layer checks the builder (P2) |
> | C3 | hand-back landed, not integrated | hand-back file exists, REQ still in `working/` | `do-work run` |
> | C4 | needs operator | `pending-answers` or `blocked` with every dependency met | `do-work clarify` |
> | C5 | waiting on dependencies | unmet `depends_on` | name the REQ it waits on |
> | C6 | earmarked | `assigned_to` set | quote the board reason; `do-work run REQ-NNN` to run it by name |
> | C7 | claim past 3 h | same rule as `verify.go:861` | `do-work run-with-recovery REQ-NNN`, only if the user knows the run is gone |
> | C8 | finalization pending | doctor or recover reports an unfinished finalization | the emitted `next_argv` |
> 
> **P2. `actions/status.md`.** Runs P1, then adds only what the CLI must not do, each item labelled "this
> machine, this session, now":
> 
> - For C2 rows owned by agents this session spawned: ListAgents and subagent transcript modified time,
>   then a one-line nudge text the user can paste or approve.
> - `ps` for test runners and gate processes named in a run-local lock file's owner line. If the owner is
>   not found, print the lock path, its age, and the `rm` command for the user to run. The action never
>   deletes it, because the same lock may belong to another checkout.
> - Output: one table (REQ with a three-to-six-word gloss, phase, minutes since activity, ETA, class),
>   then one remedy line per row that is not C1. The last line is the single most useful next command.
> 
> **P3. Routing.** Add a row to `SKILL.md` before clarify and roadmap: `status`, `status and ETA`, `eta`,
> `is it stuck`, `stuck`, `why is REQ-N taking so long`, `how do I unblock`, `why is REQ-N in` →
> `./actions/status.md`. Leave the bare word `blocked` on clarify. Status names clarify as the remedy
> when that is the cause.
> 
> **P4. `--watch`.** Same table, compact, for a `/loop` body. Document the exact line in the action:
> `/loop 15m /do-work status --watch`. No state is written between ticks.
> 
> **P5. `do-work status --fix` (action layer only; `run-status` has no `--fix`).** Recommended: run only
> remedies that are already canonical CLI commands from the row's `next_argv`, one at a time, after the
> user confirms each. Never delete a lock file, never reset a claim.
> The original ask (release a lock whose owner process is dead) needs a PID check in shipped code, which
> REQ-073 (the no-liveness-machinery decision) forbids. If the maintainer wants it, it should be its own
> REQ that revisits REQ-073 by name.
> 
> ## Acceptance check
> 
> - `do-work status`, "status and ETA please" and "is REQ-N stuck?" all route to `actions/status.md`;
>   "blocked" still routes to clarify.
> - With one claimed REQ carrying `estimate.p50_active_minutes: 60` and claimed 45 min ago, `run-status
>   --format text` prints an ETA of about 15 min. At 75 min it prints "over estimate by 15 min".
> - A claimed REQ whose hand-back file exists in the run directory is class C3 with remedy `do-work run`.
> - A blocked REQ with an unmet `depends_on` is C5 and quotes the board's Waiting reason; an earmarked
>   REQ is C6 and quotes the Earmarked reason.
> - A `full-gate.lock` file in the run directory is listed as run-local with its age and first line;
>   `run-status` never deletes it, and `--fix` never deletes it.
> - `grep -nE 'os\.Remove|syscall\.Kill|FindProcess'` over the new `run-status` Go files finds nothing.
> - The JSON output uses the `resultmodel.CommandFinding` shape that doctor already emits.
> - `--watch` output fits in 20 lines for 8 open REQs.
> 
> ## Out of scope
> 
> - A lock, heartbeat, lease, PID check or process registry in shipped Go code
>   (REQ-073, the no-liveness-machinery decision). Ages are reported; nothing is judged dead.
> - Shipping the consumer's lock helper scripts. They were improvised inside run directories and have
>   already changed shape once.
> - Any change to scheduling, selection, `advance`, the 3-hour threshold, or what `estimate` means.
> - New REQ fields or statuses.
> - A board UI change. The HTML board keeps `last activity`; this adds the terminal surface only.
> ```
