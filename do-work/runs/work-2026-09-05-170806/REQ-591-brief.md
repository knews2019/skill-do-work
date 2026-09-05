# Builder brief — REQ-591

## Where you work

- **Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-591-reduce-repeated-setup-and-unaffected-reruns-in-the-fast-gate`
- **Branch:** `worktree-agent-REQ-591-reduce-repeated-setup-and-unaffected-reruns-in-the-fast-gate`, checked out there at the integration tip `e0bdf8bf`. Clean.
- **Commit on that branch.** Do not touch the main tree — with exactly one exception, the hand-back file below.

## Never touch

- Any path under `do-work/` except the one hand-back file.
- Anything outside your declared `## Scope` list. If you need a file outside it, stop and report with the exact line and where it goes.

## Hand-back

Write `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-170806/REQ-591-handback.md` (absolute path, main tree). Never stage or commit it.

Sections: `## File manifest`, `## Measured Evidence` (every number naming the load it was taken under), `## P-A-U`, `## Decisions` (continue from D-12; the plan used D-01 through D-11), `## Discovered Tasks`, `## Lesson evidence`, `## Integration seams`.

## The measurement condition is part of the work

This checkout is shared. Load has ranged 2 to 59 today, and five canonical gate runs earlier in this run failed purely on per-file wall-clock budgets under that load. **Every measurement must be gated on a checked quiet window and must record the load beside the number.** Check with `uptime` and with `ps -Ao args= | grep -c '[m]aintainer-verify'` before and after each timed run, and discard any run whose window was not quiet. A loaded-window comparison cannot establish the improvement this request is about, and the request explicitly refuses a one-off noisy sample.

## The uncontended before-number you already have

The orchestrator ran the canonical gate to completion at your exact branch point `e0bdf8bf`, in an isolated detached worktree, load 1.98 at start:

    real 103.39   user 103.25   sys 114.13   exit 0

Use it as the baseline for the whole-gate half of the comparison, and take your own repetitions per the plan's protocol.
