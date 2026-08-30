# Work Run Manifest — 2026-08-29 21:35:39Z

Mode: `do-work run --fan-out 2` (auto-wave, worktree dispatch).
Owner checkout: `vm:/home/user/skill-do-work` on `claude/do-work-skill-loop-iif68v`.
Worktrees parent: `/home/user/skill-do-work-worktrees/`.

## Wave 1

| REQ | Operative name | Brief | Hand-back | Landed |
|---|---|---|---|---|
| REQ-390 | `worktree-agent-REQ-390-timeline-trailing-window-periods` | `REQ-390-brief.md` | `REQ-390-handback.md` | pending |
| REQ-406 | `worktree-agent-REQ-406-create-do-work-cli-foundation` | `REQ-406-brief.md` | `REQ-406-handback.md` | pending |

## Later waves

REQ-407 through REQ-420 form a serial `depends_on` chain and each becomes its own
single-REQ wave once its predecessor integrates. Wave membership is recomputed by
Step 1 after every integration, never carried forward.
