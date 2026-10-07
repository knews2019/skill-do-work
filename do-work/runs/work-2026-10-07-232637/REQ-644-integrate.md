# Integrator brief — REQ-644 ([impact-rule-change] Capture distinguishes a session earmark from work that needs the user as operator)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-644.

- REQ id: REQ-644. Working REQ path: `do-work/working/REQ-644-capture-earmark-vs-operator-blocked.md`. UR: UR-140 (`do-work/user-requests/UR-140/input.md`).
- Run directory: `do-work/runs/work-2026-10-07-232637/` (run id `work-2026-10-07-232637`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-07-232637/REQ-644-handback.md`. Builder reported one commit `b789ca74` on top of `531dd00e`, two files (`skills/do-work/actions/capture.md`, `skills/do-work/docs/work-guide.md`), both contract tests green, three decisions, one report-only discovered task.
- Operative name (branch = worktree basename): `worktree-agent-REQ-644-capture-earmark-vs-operator-blocked`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-644-capture-earmark-vs-operator-blocked`.
- Dispatch instant (for the builder-work timing event): `2026-10-07T23:30:13Z`.
- Route A, tdd: false, impact-rule-change, effort-mechanical. Prime: `_dev/primes/prime-action-files.md`. Probe: `do-work/runs/work-2026-10-07-232637/REQ-644-probe.sh` (citation contract + contract regressions, ≈22 s).
- Wave membership: REQ-643 (board: earmarked pending REQ under Pending → Earmarked) and REQ-644 (this). Set aside: none. REQ-643 is still building, so **REQ-644 is NOT the wave's last successful integration**; do not pass the wave-end fact to the reviewer.
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, cleanup pass); the coordinator does after the last integrator (work.md Step 10 says why).
- Lessons satellite for the lesson bullet: `_dev/primes/lessons-action-files.md` (then refresh its `do-work/lessons-index.md` row). Link the bullet to the flat archive path `../../do-work/archive/REQ-644-capture-earmark-vs-operator-blocked.md` (UR-140 stays open because REQ-643 is still in flight) — mirror the relative-link form the file's existing bullets use.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-644-capture-earmark-vs-operator-blocked.md`. Expected release: 0.305.77 → 0.305.78 (re-read VERSION first).
- Review depth: Route A quick scan; the reviewer should check the Red-Green Proof by reading the edited sentences as a cold capture agent, the never-invent rule, that no bold label or heading changed, and that the text does not state the board's Pending placement of an earmark (REQ-643 changes it).
- Changelog title direction: say what shipped in plain words, for example "Capture Tells a Session Earmark Apart From Work That Needs You at the Keyboard".
