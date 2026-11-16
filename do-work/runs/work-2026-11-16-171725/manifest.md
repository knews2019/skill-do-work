# Run Manifest — work-2026-11-16-171725

Run dir: do-work/runs/work-2026-11-16-171725/
Concurrency: `do-work run --coordinate`: the main session coordinates; two REQs of UR-180 claimed with one `advance UR-180 --fan-out 2` call; one pre-dispatch agent per REQ, parallel worktree builders, one background integrator per REQ in series (REQ-738 first, REQ-737 last).
Test policy: focused probe per REQ, then `bash _dev/tests/maintainer-verify.sh` at the merge revision (integrator only, under full-gate.lock), then the planned heavy lanes.
Status: done (both members finalized; UR-180 closed; REQ-739 captured as UR-181 for a later run)
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------|------|------|
| REQ-737 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-737-relation-map-phone-titles | do-work/runs/work-2026-11-16-171725/REQ-737-handback.md | released 0.324.2, UR-180 closed; three build passes (tips 8ef2084e, 4b3d6fad, 9375501d), maintainer kept form 1, asked for the page-top fix, approved the final captures; merge 8e6d94e9; catalog rebuild 3f6ebf38; finalization d6fad8a7; gate green (164 s); review 95% pass with the restatement sweep; six heavy lanes green | 2026-11-16T17:31:30Z | 30 | 1 | 0 |
| REQ-738 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-738-fold-second-opinion-direction | do-work/runs/work-2026-11-16-171725/REQ-738-handback.md | released 0.324.1; merge 943c4385; finalization ad8842ea; gate green (129 s); review 96% approve; four heavy lanes green; flat archive, UR-180 open | 2026-11-16T17:31:30Z | 16 | 1 | 0 |
