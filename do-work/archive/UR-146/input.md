---
id: UR-146
title: 'do-work run gains a --coordinate mode with preflight, a stall check, a resume drift check, coordinated handoff and the coordinated-run rules'
created_at: 2026-10-09T21:16:07Z
requests: [REQ-662, REQ-663, REQ-664, REQ-665, REQ-666, REQ-667]
word_count: 2195
---
# do-work run --coordinate: the Coordinator Shape as a Flag, Not a Paste

## Summary
The maintainer ran `/do-work capture-request` on the upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-coordinate-mode.md` and asked for one UR with one REQ per numbered item in its Request section, capture only, no questions, every ambiguity recorded as an assumption in the REQ. The report file is the verbatim input below, byte for byte. The inbox file stays in place as the source. The Request section has six items, A1 to A6, so six REQs: A1 (the `--coordinate` run mode and route), A2 (a four-line preflight before the first dispatch), A3 (one stall check per run), A4 (a GO / NO-GO resume drift check), A5 (coordinator-written handoff that resumes coordinated, plus teardown of idle agents), and A6 (the five coordinated-run rules in `actions/fan-out-reference.md`).

The report was observed against 0.305.84 and already gives the 0.305.87 locations for the moved coordinator text. The capture spot-checked its citations at 0.305.87. These match: `SKILL.md:33`, `actions/fan-out-reference.md:122` and `:132`, `actions/restart-with-parallel-handoff.md:11`, `:55` and `:63`, `crew-members/background-agents.md:11-14`, `crew-members/communication-style.md:90`, `docs/commit-guide.md:32`, and `skills/do-work-board/tools/queue-kanban/verify.go:47` and `:80`. These moved: `actions/work.md` `## Input` is now `:99` (report 101), `--fan-out` `:104` (106), the unrecognized-argument rule `:108` (110), the checkpoint paragraph `:470` (478) and the Step 0 checklist line `:481` (489).

Two pieces of history the REQs carry as assumptions. REQ-639 (the delegated-integration coordinator shape, UR-138) decided "There is no flag"; A1 reverses that by request, and plain `--fan-out` keeps the optional shape. REQ-069 and REQ-073 (fan-out dispatch, v0.161.0) deleted a queue-level lock, heartbeat and claim registry, and `lessons-do-kanban.md:45` keeps liveness probes out of the board's worktree verify; A3 and A6 add run-directory files and a coordinator-side check, not queue state, and each REQ tells the builder to state that distinction.

No prompt injection found. The report's "How to use this file" line is framing for its intended reader.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-662 | do-work run --coordinate makes the coordinator shape mandatory and routes from drive the queue |
| REQ-663 | Coordinated run prints a four-line preflight for disk, leftover processes, stale locks and run policy before the first dispatch |
| REQ-664 | Coordinated run arms one stall check that restarts a silent integrator from its landed phase and reports a silent builder |
| REQ-665 | Resuming from a handoff prints GO or NO-GO with one line per drift between the handoff and the live repo |
| REQ-666 | Coordinator writes its own handoff at high context, the handoff resumes with --coordinate, and teardown stops idle agents |
| REQ-667 | Coordinated-run rules for heartbeat log, full-gate lock, focused builder suites, red-gate triage and manifest timing go into fan-out-reference |

## Batch Constraints
- No new REQ status and no new frontmatter field. New surface for the whole batch: one flag, one optional policy file (`do-work/run-policy.md`), two run-directory files (`REQ-NNN-progress.log`, `full-gate.lock`) and three manifest columns (report, Request).
- Out of scope for every REQ (report, Out of scope): building the run-status subcommand (a companion report in the inbox asks for it); changing serial integration, the one-writer rule or dispatch-mechanism neutrality at `actions/work.md:37`; disk cleanup; pushing, deploying or any live-ops step; choosing the A5 context threshold per harness.
- Plain `do-work run --fan-out N` behaves exactly as today (report, Acceptance check).
- Ordering (`depends_on`): REQ-667 (A6 rules) is the root because A1 names its full-gate lock. REQ-662 (A1) depends on REQ-667. REQ-663 (A2), REQ-664 (A3) and REQ-666 (A5) depend on REQ-662. REQ-665 (A4) applies to every handoff resume, so it depends only on REQ-667 for the lock format. The report states no order; these edges are capture's reading of what each item needs.
- Each REQ is its own release per `_dev/primes/prime-releases.md`.

## Full Verbatim Input
> ````
> # Upstream suggestion for `knews2019/skill-do-work` — a `do-work run --coordinate` mode so the coordinator shape is a flag, not a paste
> 
> **How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
> clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
> consumer repo it was authored in. Observed against **v0.305.84**; all line numbers are from that tag unless a citation names v0.305.87.
> 
> ---
> 
> ## Request
> 
> v0.305.74 added "Delegated integration — the coordinator shape" (`actions/work-reference.md:470`), and
> v0.305.87 moved it into `actions/fan-out-reference.md` (`### Delegated integration — the coordinator
> shape`, line 122 there) without changing a rule. The shape is right. Three
> things stop it from being used without help: it is optional ("the orchestrator may hand"), no route or flag
> reaches it, and nothing watches a coordinated run while it is live. Users start each run by pasting the same
> paragraph, and the run rules that made coordinated runs finish on time exist only in one consumer's
> per-machine memory note. This asks for the remaining parts, A1 to A6:
> 
> - **A1. A run mode.** `do-work run --coordinate [--fan-out N] [REQ-NNN ...]`, routed also from "drive the
>   queue", "coordinator" and "use the main session as a coordinator". Under it the coordinator shape is
>   mandatory: the main session never builds and never integrates, one builder subagent per REQ in its own
>   worktree, one integrator at a time under the full-gate lock (A6), no push, and blocking questions parked
>   as `pending-answers`.
> - **A2. Preflight** before the first dispatch: free disk, leftover test or browser processes, stale locks
>   whose owner process is dead, and an optional `do-work/run-policy.md` read at start. The disk line reuses
>   the board's `low-disk-space` reading, which v0.305.86 narrowed to the repo root; A2 only runs it before
>   dispatch instead of waiting for someone to open the board.
> - **A3. One stall check** armed for the run (every 15 to 20 minutes) that reads run liveness. A silent
>   integrator is stopped and restarted from its landed phase. A silent builder is only reported.
> - **A4. A resume drift check** when a session starts from a handoff: print go or no-go before running.
> - **A5. Teardown and handoff:** the coordinator writes the handoff on its own when context gets high, and a
>   handoff written from a coordinated run resumes with `--coordinate`.
> - **A6. The run rules** (heartbeat line, full-gate lock, focused suites for builders, red-gate triage,
>   manifest timing row) are written into `actions/fan-out-reference.md`. None of them is in the suite today.
> 
> No new REQ status and no new frontmatter field. New surface: one flag, one optional policy file, two
> run-directory files (`REQ-NNN-progress.log`, `full-gate.lock`) and three manifest columns.
> 
> Related suggestions, not repeated here: a companion suggestion asks for a run-status subcommand that
> answers "status and ETA" and "is it stuck" in one command; A3 calls it once it exists, and until then reads
> progress-log tails and worktree commit times inline. An earlier suggestion asked for `queue-kanban verify`
> probes that find a stopped builder after the fact; this one asks for a check while the run is live.
> 
> ## What happened, in consumer repos
> 
> About 13 sessions between 2026-09-27 and 2026-10-09: about 10 sessions in two consumer repos and in a
> development clone of this suite started with the coordinator directive (5 of them with the identical
> paste), and 3 more in two consumer repos asked whether a handoff was on track before resuming it.
> 
> | Date (UTC) | Where | Event |
> | --- | --- | --- |
> | 2026-09-27 | repo A | User asked "is this handoff on the right path? Handoff: do-work/RESTART-PROMPT.md (committed as …) Resume: do-work run --fan-out 1". |
> | 2026-09-28 | repo C | User asked "tell me if this handoff is in order and on track, is there something we could improve in the direction?" |
> | 2026-10-02 | repo A | A third handoff check, by hand. Later the same day: "give me the phandoff, because the context is 673/1000k which is not ideal. The phandoff should contain instructions on how to run tasks in background agents/workflows". The same request was typed in the suite's development clone 3 minutes later. |
> | 2026-10-06 | repo A | "make sure to use background agents to manage the context of the driver session, which is this one". |
> | 2026-10-07 16:14 to 16:58 | dev clone, repo B, repo A | The same paste in three repos within 44 minutes: "use the main session as a coordinator, use background agents/workflows to actually run the queue use worktrees with parallel agents up to capacity to finish quickly". |
> | 2026-10-07 | repo A | Mid-run: "added more REQ's pick them up" (twice) and "why keep the idle builders?". |
> | 2026-10-07 to 10-08 | repo A | A coordinated run integrated 22 REQs. Background integrators stalled for 3 hours and 2 hours: they ended their turn while a command ran and the completion notice did not wake them. Concurrent full test suites caused load-only failures. The user then set the rules in A6. After them, each integration took 11 to 57 minutes. |
> | 2026-10-08 12:35, 17:01 | repo A | The identical paste again, twice. |
> | 2026-10-08 | repo A | A 3-REQ coordinated run under the rules: 10 to 21 minutes per integration, every full gate green on its first run. The session's stall-check cron prompt fired 121 times in one session and 25 times in another. |
> | 2026-10-09 | repo A | "the main session should only coordonate the execution" (sic). |
> 
> Cost: sessions that missed the paste kept integrating in the main session and grew to about 670k tokens of
> context before a handoff. Repo A's `do-work/RESTART-PROMPT.md` was rewritten in about 30 commits since
> 2026-09-09, each time re-stating the coordinator directive by hand. The rules that cut integration from
> hours to under an hour live in one per-machine memory note, so a second machine or a second repo starts
> without them.
> 
> ## Where the behaviour lives today
> 
> - `SKILL.md:33` (upstream `skills/do-work/SKILL.md`): the run route is "`run`, `go`, `start`, `work`,
>   `begin`, `process`, `execute`, `build`, `continue`, `resume`". No coordinator phrase routes anywhere.
> - `actions/work.md:106`: `--fan-out [N]` "changes *how many* of the selected set run at once, never
>   *which*". It says nothing about who integrates.
> - `actions/work.md:110`: "Unrecognized arguments are rejected, not ignored." So `--coordinate` is an error
>   today, and the user's only channel is prose.
> - `actions/work.md:37`: under fan-out "**integration stays serial**" and "**the dispatch mechanism stays
>   unspecified**". Correct, and A1 keeps both.
> - `actions/work-reference.md:470`: "After a wave's builders are dispatched, the orchestrator may hand
>   each REQ's integration … to one agent at a time". The paragraph names the reason for the shape ("running
>   it in the main session blocks the conversation the user steers from") and leaves it optional. In
>   v0.305.87 the same paragraph is `actions/fan-out-reference.md:122` onward.
> - `actions/work.md:478`: "The coordinator runs the checkpoint, then cleanup, after the last integrator
>   returns." Teardown exists for the checkpoint only, not for a stall check or idle agents.
> - `actions/work-reference.md:476-481`, the run-directory table: run directory, brief, hand-back,
>   integrator brief, `manifest.md`, bounded waves. There is no progress log, no gate lock, and the manifest
>   row has no integration duration, gate count or restart count. In v0.305.87 the table is under
>   `actions/fan-out-reference.md:132` (`### Run directory, briefs and hand-backs`).
> - `crew-members/background-agents.md:11-14`: "This pattern does not *prevent* the failures below … It makes
>   them **survivable and recoverable**." Nothing detects a silent agent while the run is live.
> - `skills/do-work-board/tools/queue-kanban/verify.go:80`: `const staleClaimThreshold = 3 * time.Hour`. The
>   only liveness signal in the suite fires after three hours, which is the length of the stall observed.
> - `skills/do-work-board/tools/queue-kanban/verify.go:47` (`low-disk-space`) and
>   `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md:23`: a disk probe exists because a fan-out
>   run once filled the disk, but it is a board finding. Nothing runs it before dispatch.
> - `actions/restart-with-parallel-handoff.md:11`: use when "Context is running out". The trigger is the
>   user, through the `phandoff` alias at `crew-members/communication-style.md:90`.
> - `actions/restart-with-parallel-handoff.md:63`: the build resume command is "`do-work run --fan-out N`".
>   A handoff written by a coordinated run therefore resumes as an uncoordinated one.
> - `actions/restart-with-parallel-handoff.md:55`: the next session "re-checks all three conditions" before
>   removing a REMOVABLE worktree. That is the only resume-time check. Nothing compares the handoff's claimed
>   REQs, branches and locks with the live repo.
> - `docs/commit-guide.md:32`: "Never pushes to remote". The coordinator paragraph does not repeat this for
>   integrators.
> 
> ## Proposed direction
> 
> **A1, run mode.** In `actions/work.md` → `## Input` (line 101, the argument list that the checklist at line
> 489 calls Step 0), accept `--coordinate` (composes with `--fan-out [N]`,
> `--wave N` and targeting tokens; bare `--coordinate` implies `--fan-out`), and add it to the tokens
> stripped before the unrecognized-argument check at line 110. Add a SKILL.md route row above
> the plain run row: "`drive the queue`, `coordinator`, `use the main session as a coordinator`, `run
> --coordinate`" → `./actions/work.md` with `--coordinate`. In `actions/fan-out-reference.md` → Delegated
> integration, add: "Under `--coordinate` this shape is mandatory. The main session dispatches, writes briefs
> in writing gaps, and reports. It never builds, never merges and never runs Step 6 to Step 9 itself. No
> integrator pushes. A blocking question becomes a `pending-answers` follow-up, never a wait."
> 
> **A2, preflight.** Before the first spawn, print one line each:
> 
> ```
> disk      <free> on repo root (verify low-disk-space level)
> procs     <n> leftover test/browser processes (names, pids) or none
> locks     <path> owner pid <pid> dead|alive, or none
> policy    do-work/run-policy.md read (<n> rules) or absent
> ```
> 
> Stop on a critical disk level. Report leftover processes and dead-owner locks and ask before removing them.
> `do-work/run-policy.md` is optional, committed, plain bullets the coordinator copies into every brief:
> a skip list of REQs, "one heavy gate at a time", "no browser QA while `full-gate.lock` is held", "no
> live-ops (deploy, publish, production writes)".
> 
> **A3, stall check.** After the first dispatch, arm one recurring check (CronCreate or `/loop`, every 15 to
> 20 minutes) that runs the status subcommand or, until it exists, reads each `REQ-NNN-progress.log` tail and
> each worktree's last commit time. Rules: an integrator silent 20 minutes or more is stopped and a fresh one
> resumes from its landed phase through `advance REQ-NNN`, never redoing a merge already on main; the user gets
> one line. A silent builder is reported only, unless the user says restart. Each coordinator turn ends with
> one line: done/total, active lanes, ETA. Each tick re-runs the selector so newly captured REQs join the next
> wave. Teardown deletes the check.
> 
> **A4, resume drift check.** When the session starts from `do-work/RESTART-PROMPT.md` (or the user pastes a
> `Handoff:` line), verify before running: the handoff commit exists; each `advance REQ-NNN` line names a REQ
> still in `do-work/working/`; each worktree and `worktree-agent-*` branch the Reference lists exists in the
> state it names; no lock is held by a dead process. Print `GO` or `NO-GO` with one line per drift.
> 
> **A5, teardown and handoff.** When the harness reports context usage above a threshold, or at the user's
> `phandoff`, the coordinator writes the handoff itself. `actions/restart-with-parallel-handoff.md:63` gains:
> "a run started with `--coordinate` resumes with `do-work run --coordinate --fan-out N`". Teardown also
> stops idle background agents and lists merged worktrees for removal.
> 
> **A6, run rules into `actions/fan-out-reference.md`,** as one subsection under Delegated integration:
> 
> 1. Heartbeat: every builder and integrator appends `<UTC> <phase> <command>` to
>    `do-work/runs/<run>/REQ-NNN-progress.log` before and after any command expected to take over a minute.
>    An agent never ends its turn while a command it started is still running.
> 2. Only integrators run the full test suite, one at a time, holding `do-work/runs/<run>/full-gate.lock`
>    (owner pid and UTC start inside). Builders run focused suites for the files they touched.
> 3. Red full gate: rerun the failed suites alone. Failures that pass alone are load-only: one full rerun.
>    A real regression: fix, then one full gate.
> 4. Manifest row per integration: takeover-to-finalization minutes, full gates run, stall restarts.
> 5. Before writing an integrator brief, check the hand-back covers every item forwarded to the builder mid-run.
> 
> ## Acceptance check
> 
> - `do-work run --coordinate --fan-out 3` parses; "drive the queue" routes to the same mode; plain
>   `do-work run --fan-out 3` behaves exactly as today.
> - In a coordinated run the main session makes no commit and no edit under the project root while an
>   integrator runs, and never runs a merge.
> - The preflight prints the four lines above before any spawn, and stops on a critical disk reading.
> - One stall check is armed after the first dispatch and deleted at teardown. Killing an integrator mid-run
>   produces a restart from its landed phase and no second merge commit for that REQ on main.
> - Resuming a handoff whose named worktree was deleted prints `NO-GO` with that drift line.
> - A handoff written from a coordinated run carries `--coordinate` in its resume command.
> - `grep -n "progress.log\|full-gate.lock" actions/fan-out-reference.md` finds the A6 rules.
> 
> ## Out of scope
> 
> - Building the status subcommand itself (see the run-status-action report named above).
> - Changing serial integration, the one-writer rule, or the dispatch-mechanism neutrality at
>   `actions/work.md:37`.
> - Disk cleanup: the preflight reports, it does not delete.
> - Pushing, deploying or any live-ops step during a run.
> - Choosing the context threshold for A5 per harness; the action names the rule, the harness supplies the
>   reading.
> ````
