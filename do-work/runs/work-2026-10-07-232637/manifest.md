# Run Manifest — work-2026-10-07-232637

Run dir: do-work/runs/work-2026-10-07-232637/
Concurrency: 2 (UR-140 wave REQ-643 + REQ-644, disjoint write sets, no depends_on; coordinator shape: coordinator pre-dispatch, one builder per REQ in its own worktree, one integrator per REQ in series)
Status: building (REQ-643) / integrating (REQ-644)
Set aside: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant |
|-----|---------|----------------|---------------|--------|------------------|
| REQ-643 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-643-board-earmarked-pending-subgroup | do-work/runs/work-2026-10-07-232637/REQ-643-handback.md | builder dispatched (pre-flight green: gate 125s, probe 3.6s) | 2026-10-07T23:31:21Z |
| REQ-644 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-644-capture-earmark-vs-operator-blocked | do-work/runs/work-2026-10-07-232637/REQ-644-handback.md | handback received b789ca74 (2 files); integrator starting | 2026-10-07T23:30:13Z |
