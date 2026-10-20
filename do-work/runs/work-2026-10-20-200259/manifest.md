# Run Manifest — work-2026-10-20-200259

Run dir: do-work/runs/work-2026-10-20-200259/
Concurrency: `do-work run UR-160 --coordinate` (bare `--fan-out`, bound 2). Main session coordinates only; pre-dispatch, builder and integrator run as background agents. REQ-704 depends on REQ-703 (same Go file), so members run one at a time.
Status: in-flight
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------------------------------|------------|----------------|
| REQ-703 | builder-703 (background agent, Route C, worktree) | worktree-agent-REQ-703-second-opinion-presets-and-direction | do-work/runs/work-2026-10-20-200259/REQ-703-handback.md | handback f06e82e0; merge dd3639bf, review-fix re-merge 1d68e403; gate green x2; review 93% Approve; heavy 4 lanes green; released 0.306.6 (a6c00290) | 2026-10-20T20:17:40Z | 25 | 2 | 0 |
| REQ-705 | builder-705 (background agent, Route B, worktree) | worktree-agent-REQ-705-bundle-local-second-opinion-presets | do-work/runs/work-2026-10-20-200259/REQ-705-handback.md | answered by the user 2026-10-20 (presets only inside the bundle); dispatched from 11780ac0; pre-flight green (gate 135 s) | 2026-10-20T21:31:30Z | | | |
| REQ-704 | (waits on REQ-705 build) | | do-work/runs/work-2026-10-20-200259/REQ-704-handback.md | queued; dependency-deferred | | | | |
