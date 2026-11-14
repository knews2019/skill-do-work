# Run Manifest — work-2026-11-15-012821

Run dir: do-work/runs/work-2026-11-15-012821/
Concurrency: `do-work run --coordinate`: the main session coordinates; ten REQs claimed with two `advance --fan-out 5` calls; parallel worktree builders, one background integrator per REQ in series. REQ-734 (depends on REQ-733) is claimed after REQ-733 finalizes. REQ-727 stays blocked on the maintainer.
Test policy: focused probe per REQ, then `bash _dev/tests/maintainer-verify.sh` at the merge revision (integrator only, under full-gate.lock), then the planned heavy lanes.
Status: in-progress
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------|------|------|
| REQ-723 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-723-archive-collision-skips-assets | do-work/runs/work-2026-11-15-012821/REQ-723-handback.md | building | 2026-11-14T23:32:28Z | | | |
| REQ-724 | pending | pending | do-work/runs/work-2026-11-15-012821/REQ-724-handback.md | pre-dispatch | | | | |
| REQ-725 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-725-memory-note-destination-report | do-work/runs/work-2026-11-15-012821/REQ-725-handback.md | building | 2026-11-14T23:32:30Z | | | |
| REQ-726 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-726-no-review-canvas-decision | do-work/runs/work-2026-11-15-012821/REQ-726-handback.md | building | 2026-11-14T23:32:56Z | | | |
| REQ-728 | pending | pending | do-work/runs/work-2026-11-15-012821/REQ-728-handback.md | pre-dispatch | | | | |
| REQ-729 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-729-judge-relative-out-rule | do-work/runs/work-2026-11-15-012821/REQ-729-handback.md | building | 2026-11-14T23:32:29Z | | | |
| REQ-730 | pending | pending | do-work/runs/work-2026-11-15-012821/REQ-730-handback.md | pre-dispatch | | | | |
| REQ-731 | pending | pending | do-work/runs/work-2026-11-15-012821/REQ-731-handback.md | pre-dispatch | | | | |
| REQ-732 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-732-allow-docs-run-prompts | do-work/runs/work-2026-11-15-012821/REQ-732-handback.md | building | 2026-11-14T23:33:30Z | | | |
| REQ-733 | pending | pending | do-work/runs/work-2026-11-15-012821/REQ-733-handback.md | pre-dispatch | | | | |
