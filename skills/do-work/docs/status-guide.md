# Status

Read-only answer to "where is my in-flight work, when will it finish, and what do I do about the stuck one?" One row per claimed, blocked, waiting or earmarked REQ: phase, minutes since its last activity, ETA, one class and one remedy. Never modifies anything.

> **Sister actions:** `do-work roadmap` surveys the whole queue, including ready work (see `docs/roadmap-guide.md`). `do-work forensics` looks for *broken* state such as damaged records or hollow completions (see `docs/forensics-guide.md`). Status covers only work that is in flight now.

Status reports ages, never verdicts. A claim carries no lock and no heartbeat, so it never says a run or a builder is dead. It says how long ago something last happened and lets you decide.

## What each row shows

| Column | Where it comes from |
|--------|---------------------|
| **REQ** and short title | The REQ file |
| **Phase** | The board's last activity: the newest phase stamp or `[REQ-NNN]` commit |
| **Since activity** | Minutes since that last activity |
| **ETA** | The REQ's frozen `p50_active_minutes` minus minutes since claim; "over estimate by N min" once it passes, "no estimate" when none was frozen |
| **Class** | One of C1 to C8, below |

## Classes

First match wins, in this order:

| Class | Meaning | What to do |
|-------|---------|------------|
| **C8** finalization pending | An earlier finalization stopped partway | The command doctor prints |
| **C3** hand-back landed | The builder's hand-back file is in the run directory | `do-work run` integrates it |
| **C4** needs operator | Needs input · Blocked on the board | `do-work clarify` |
| **C5** waiting on dependencies | Waits for the REQ it names | Nothing; it moves when that REQ is done |
| **C6** earmarked | Reserved for another session, which a default run skips | `do-work run REQ-NNN` runs it by name |
| **C7** claim past 3 hours | Claimed longer than the board's stale-claim threshold | Read `git log --full-history -- <REQ file>`. Use `do-work run-with-recovery REQ-NNN` only if you know the run that claimed it is gone: it requeues the claim and strips generated sections |
| **C2** quiet | No new stamp or commit for about 20 minutes | The action checks agents this session started |
| **C1** progressing | Activity in the last 20 minutes | Nothing |

The report also lists files in the run directory that the suite does not define (lock files, progress logs), labelled "run-local, not interpreted", with their age and first line. When a lock's owner process is not running on this machine, the action prints an `rm` line with the real path. You decide whether to run it: the lock may belong to another checkout.

## Usage

```
do-work status
do-work status REQ-412
do-work status --watch
do-work is it stuck
do-work how do I unblock REQ-412
```

To watch a long run, let the harness repeat it:

```
/loop 15m /do-work status --watch
```

`--watch` prints the same rows in about 20 lines for 8 open REQs and writes nothing between ticks.

## When NOT to use

- Want the whole queue, including what is ready next → `do-work roadmap`
- Suspect damaged records or broken history → `do-work forensics`
- Want to answer pending questions → `do-work clarify`
