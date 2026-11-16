# Run Manifest — work-2026-11-16-171725

Run dir: do-work/runs/work-2026-11-16-171725/
Concurrency: `do-work run --coordinate`: the main session coordinates; two REQs of UR-180 claimed with one `advance UR-180 --fan-out 2` call; one pre-dispatch agent per REQ, parallel worktree builders, one background integrator per REQ in series (REQ-738 first, REQ-737 last).
Test policy: focused probe per REQ, then `bash _dev/tests/maintainer-verify.sh` at the merge revision (integrator only, under full-gate.lock), then the planned heavy lanes.
Status: in progress
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------|------|------|
| REQ-737 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-737-relation-map-phone-titles | do-work/runs/work-2026-11-16-171725/REQ-737-handback.md | pass 1 handed back 17:45:01Z (tip 8ef2084e); maintainer kept form 1 and asked for the page-top fix; pass 2 building | 2026-11-16T17:31:30Z | | | |
| REQ-738 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-738-fold-second-opinion-direction | do-work/runs/work-2026-11-16-171725/REQ-738-handback.md | handed back 17:43:57Z (tip 6432bb22); captures approved by the maintainer; integrating | 2026-11-16T17:31:30Z | | | |
