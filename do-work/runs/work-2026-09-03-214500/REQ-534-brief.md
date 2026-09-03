# REQ-534 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-534-run-blocked-probes-from-the-repository-root-and-propagate-interruptions`
Branch: `worktree-agent-REQ-534-run-blocked-probes-from-the-repository-root-and-propagate-interruptions`
Handback: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-03-214500/REQ-534-handback.md`

Implement the Route C plan captured for REQ-534. Run relative probes from the selected repository root. Install Unix signal notification before child launch. After owned-group cleanup/reaping, propagate SIGINT/SIGHUP/SIGTERM as a typed interruption that stops selection, emits no later selected REQ, returns non-success, and sets `128+signal`; retain timeout 124 and launch-failure 125. Keep `ProbeRunner` stable through a root-capturing handler closure and update Windows signature parity.

Read `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, its full lessons satellite as the prime requires, and the always-on crew files. Stay within the declared seven-file Scope. Do not read or write any `do-work/` path in the worktree. Commit source/test changes on the branch. Write the absolute handback file with branch, commit, exact manifest, tests and durations, RED/GREEN evidence, lesson reads, P-A-U evidence, decisions, discovered tasks, seams, and scope issues.
