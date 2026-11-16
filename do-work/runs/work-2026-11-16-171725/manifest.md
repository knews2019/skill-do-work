# Run Manifest — work-2026-11-16-171725

Run dir: do-work/runs/work-2026-11-16-171725/
Concurrency: `do-work run --coordinate`: the main session coordinates; two REQs of UR-180 claimed with one `advance UR-180 --fan-out 2` call; one pre-dispatch agent per REQ, parallel worktree builders, one background integrator per REQ in series (REQ-738 first, REQ-737 last).
Test policy: focused probe per REQ, then `bash _dev/tests/maintainer-verify.sh` at the merge revision (integrator only, under full-gate.lock), then the planned heavy lanes.
Status: in progress
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------|------|------|
| REQ-737 | builder agent (Fable 5.2, worktree) | (set at dispatch) | do-work/runs/work-2026-11-16-171725/REQ-737-handback.md | claimed; pre-dispatch running | | | | |
| REQ-738 | builder agent (Fable 5.2, worktree) | (set at dispatch) | do-work/runs/work-2026-11-16-171725/REQ-738-handback.md | claimed; pre-dispatch running | | | | |
