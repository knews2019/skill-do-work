# REQ-371 builder brief

- Base: to be created from the claim/bookkeeping commit
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-371-keep-timeline-bars-inside-the-plot-after-the-drawer-opens`
- Branch: `worktree-agent-REQ-371-keep-timeline-bars-inside-the-plot-after-the-drawer-opens`
- Request: `do-work/working/REQ-371-keep-timeline-bars-inside-the-plot-after-the-drawer-opens.md`
- Declared scope: `web/board-timeline.js`, `timeline_browser_probe_test.go`

Read `CLAUDE.md`, the request, `_dev/primes/prime-kanban-board.md`, and archived REQ-331. Reproduce the
isolated Chromium 1228 failure first. Measure the ResizeObserver → frame scheduling boundary under
the real generated page; plan before editing. Restore condition-based remeasurement without
enumerating drawer callers, keep axis/row width alignment and drawer close recovery, and strengthen
the existing trusted probe only where needed for current-browser non-vacuity. Preserve one
measurement per render and zero-size behavior. Replay the required-invalidation mutation, run the
named probe repeatedly, Timeline/strict browser tests, focused/static gates, and canonical verification.
Do not edit `do-work/` or release files in the worktree. Commit and write the hand-back to
`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-201759/REQ-371-handback.md`
with P-A-U, plan, exploration, exact scope, decisions, browser measurements, mutations, hash, and clean status.
