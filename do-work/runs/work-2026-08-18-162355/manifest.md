# Run manifest — work-2026-08-18-162355

Integration branch: `claude/restart-prompt-handoff-nprlgd` · branch point `67dae6b` · `--fan-out 3`

Wave 1 — the three dependency chain-heads, taken in numeric id order from the ready set.

| REQ | Route | Estimate | Branch / worktree basename | Brief | Hand-back | Landed |
|---|---|---|---|---|---|---|
| REQ-246 | C | 45 min | `worktree-agent-REQ-246-repair-wrong-queue-and-working-timestamps-from-the-session-hook` | `REQ-246-brief.md` | `REQ-246-handback.md` | pending |
| REQ-248 | B | 30 min | `worktree-agent-REQ-248-anchor-durations-day-buckets-to-utc-midnight` | `REQ-248-brief.md` | `REQ-248-handback.md` | pending |
| REQ-249 | B | 45 min | `worktree-agent-REQ-249-decide-the-cross-package-citation-path-form` | `REQ-249-brief.md` | `REQ-249-handback.md` | pending |

Total estimated effort: 120 active minutes · estimated critical path: 45 active minutes.

Worktrees live at `/home/user/skill-do-work-worktrees/<basename>`, outside the repo working tree.

Integration is serial: merge → qualify → test → review → changelog → archive, one REQ at a time.
