# Status Action

> **Part of the do-work skill.** Invoked when the user asks for the status, ETA or unblocking of in-flight work. Read-only: one row per claimed, blocked, waiting or earmarked REQ with its phase, ages, ETA, one class and one remedy. It belongs in core because it completes `actions/forensics.md` (broken state) and `actions/roadmap.md` (intended state) for work that is in flight now, through the core `run-status` command. User-facing walkthrough: [`docs/status-guide.md`](../docs/status-guide.md).

Ages are reported, never verdicts: this action never says a run or builder is dead. Claims carry no lock and no liveness record, so an age is the only fact there is.

## When to Use

**Use when:**
- The user asks "status", "status and ETA", "eta", "is it stuck", "why is REQ-N taking so long", "how do I unblock", or "why is REQ-N in" a column.
- A coordinator answers a mid-run progress question, or a `/loop` body needs a compact tick.

**Do NOT use when:**
- The user wants the whole queue, including ready work → `actions/roadmap.md`.
- The user suspects damaged records, collisions or hollow completions → `actions/forensics.md`.
- The user wants to answer pending questions → `actions/clarify.md`.

## Input

`$ARGUMENTS`: an optional `REQ-NNN` (passed as `--req`), an optional run directory (passed as `--run`; default is the newest `do-work/runs/work-*`), and an optional `--watch`.

## Steps

### Step 1: Build the board tool

Resolve the project root as `../../do-work-board/actions/board.md` Step 3 does, then build the board as its Step 4 does. Without the Go toolchain, stop with that action's Step 2 message: the board's facts are the input to every row.

### Step 2: Run the report

Write the board's open-work facts to a temporary file, hand it to `run-status`, and remove it. Add `--run DIR`, `--req REQ-NNN` and `--watch` from the input to the `run-status` line.

```bash
board_facts="$(mktemp)"
"<skill-root>/../do-work-board/tools/queue-kanban/queue-kanban" open-work --format json --repo-root "<project-root>" >"$board_facts" \
  && "<skill-root>/tools/do-work-cli.sh" --repo-root "<project-root>" run-status --board-facts "$board_facts"
rm -f "$board_facts"
```

Exit 0 means the report was produced. A `queue-kanban` error, or exit 1 with outcome `refused`, means no report: relay the error and stop. Do not rebuild the rows by hand from REQ files.

Each row is one `CommandFinding` (`code` is the class, `next_argv` the remedy) plus a typed `run_status.rows[]` entry; `run_status.run_local_files[]` lists run-directory files the suite does not define, labelled "run-local, not interpreted". Classes, first match wins:

| Order | Class | Shown when | Remedy (`next_argv`) |
| --- | --- | --- | --- |
| 1 | C8 finalization pending | doctor reports an unfinished finalization for the REQ | doctor's own `next_argv` |
| 2 | C3 hand-back landed | claimed and `<run>/REQ-NNN-handback.md` exists | `do-work run` |
| 3 | C4 needs operator | board column Needs input · Blocked | `do-work clarify` |
| 4 | C5 waiting on dependencies | board column Pending → Waiting | none; the row names the REQ it waits on |
| 5 | C6 earmarked | board column Pending → Earmarked | `do-work run REQ-NNN` |
| 6 | C7 claim past threshold | claimed at least the board's stale-claim threshold ago | read-only `git log --full-history -- <REQ path>` |
| 7 | C2 quiet | claimed, no new stamp or commit for about 20 min, or no activity record | none; Step 3 checks the builder |
| 8 | C1 progressing | claimed and active within about 20 min | none |

The 20-minute boundary is display only. ETA is the frozen `p50_active_minutes` minus minutes since claim; past it the row says "over estimate by N min".

### Step 3: Add this-session checks

Only these two checks, and only on this machine. Label each line "this machine, this session, now".

- **C2 rows owned by an agent this session spawned.** Read the harness's agent list (for example a `ListAgents` tool) and the modified time of that agent's transcript file. Print the transcript age, then one nudge line the user can paste or approve, such as `Status check: REQ-NNN shows no stamp or commit for N min. Reply with your current step, or hand back if done.` Skip rows owned by other sessions.
- **Run-local lock files.** For a run-local file whose first line names an owner process, look it up: `ps -p <pid> -o pid=,etime=,command=` when the line carries a pid, otherwise `pgrep -fl -- '<command named in the line>'`. When no owner is found, print the lock's path, its age, and `rm -- '<absolute lock path>'` with the real path filled in. Never run that `rm`.

### Step 4: Report

Relay the command's text: the table (REQ, short title, phase, minutes since activity, ETA, class), then one remedy line per row that is not C1. Put the Step 3 lines after the remedy lines and keep the command's `next:` line last: it is the single most useful next command. With `--watch`, relay the compact output as is.

## Watching

`--watch` prints the same data in at most about 20 lines for 8 open REQs and writes nothing. Repetition comes from the harness, never from the command:

```text
/loop 15m /do-work status --watch
```

## Rules

- **C7's remedy is an inspection.** `next_argv` is followed literally, so it is the read-only `git log`. `do-work run-with-recovery REQ-NNN` requeues the claim and strips generated sections; name it only as the stop reason says, and only if the user knows the run that claimed it is gone.
- **Never delete a run-local file.** The same lock may belong to another checkout or session, so the action prints the `rm` line and the user decides.
