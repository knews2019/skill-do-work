# Work run 2026-08-18-211613 — wave 1

Orchestrator: main tree `/home/user/skill-do-work`, integration branch `claude/restart-prompt-handoff-nprlgd`.
Mode: `do-work run --fan-out 3` (auto-wave, worktree dispatch, serial integration).

| REQ | Route | Builder | Operative name | Hand-back | Landed |
|---|---|---|---|---|---|
| REQ-257 | B | builder-257 | `worktree-agent-REQ-257-decide-whether-the-repairer-learns-offset-and-fractional-stamps` | `REQ-257-handback.md` | — |
| REQ-259 | B | builder-259 | `worktree-agent-REQ-259-retire-the-skill-root-reading-outside-backtick-spans` | `REQ-259-handback.md` | — |
| REQ-260 | A | builder-260 | `worktree-agent-REQ-260-gofmt-the-durations-day-domain-truncation` | `REQ-260-handback.md` | — |

## Wave selection

Ready set (pending, dependency-ready, unclaimed, unassigned): REQ-257, 258, 259, 260, 261, 262, 263, 264, 265, 266 — ten, not the twelve the restart prompt claimed. **REQ-267 and REQ-268 are `pending-answers`**, so auto-wave's first condition excludes them; the handoff's "12 pending, 0 pending-answers" line was wrong, and its suggested first wave named REQ-267.

Numeric order would take 257, 258, 259. **REQ-258 was skipped** and 260 taken in its place: REQ-258 restructures `_dev/tests/prescribed-shell-scripts-behavior.sh` wholesale while REQ-257 must add lock-in cases to that same file, so the two collide by construction and REQ-257's cases would land in a file REQ-258 has dissolved. That is the one deviation from numeric order; everything else follows it. The resulting three write sets are disjoint.
