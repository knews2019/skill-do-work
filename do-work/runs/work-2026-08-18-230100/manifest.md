# Work run 2026-08-18-230100 — wave 2

Orchestrator: main tree `/home/user/skill-do-work`, integration branch `claude/restart-prompt-handoff-nprlgd`.
Mode: `do-work run --fan-out 3` (auto-wave, worktree dispatch, serial integration).

| REQ | Route | Operative name | Hand-back | Landed |
|---|---|---|---|---|
| REQ-267 | B | `worktree-agent-REQ-267-close-the-two-remaining-repairer-shape-divergences` | `REQ-267-handback.md` | — |
| REQ-265 | A | `worktree-agent-REQ-265-raise-the-two-under-bounding-mark-label-constants` | `REQ-265-handback.md` | — |
| REQ-261 | A | `worktree-agent-REQ-261-decide-whether-queue-kanban-grows-a-date-only-mode` | `REQ-261-handback.md` | — |

## Wave selection

Ready set after wave 1: REQ-258, 261, 262, 263, 264, 265, 266, 267, 268, 269, 270, 271, 272 — thirteen.

`_dev/tests/prescribed-shell-scripts-behavior.sh` is the scheduling bottleneck: 258, 263, 264, 267, 268 and 271 all write to it, so exactly one can run per wave. **REQ-267 takes that slot** rather than the numerically-first REQ-258, because it is the only queued item carrying a live defect that can wedge the SessionStart hook into failing every session with no self-heal.

The other two are the largest disjoint pair available: REQ-265 (`durations_test.go`) and REQ-261 (`work-reference.md`). REQ-272 was passable on value but collides with REQ-267 on `repair-req-timestamps.sh` *and* with REQ-261 on `work-reference.md`; REQ-270 collides with REQ-261 on the same file.
