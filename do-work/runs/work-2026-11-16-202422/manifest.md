# Run Manifest — work-2026-11-16-202422

Run dir: do-work/runs/work-2026-11-16-202422/
Concurrency: `do-work run REQ-739` in the coordinator shape: the main session coordinates; one pre-dispatch agent, one worktree builder, one background integrator.
Test policy: focused probe, then `bash _dev/tests/maintainer-verify.sh` at the merge revision (integrator only, under full-gate.lock), then the planned heavy lanes.
Status: done (REQ-739 finalized; UR-181 closed; the queue is empty)
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------|------|------|
| REQ-739 | builder agent (Opus 5.5, worktree) | worktree-agent-REQ-739-header-phone-side-margin | do-work/runs/work-2026-11-16-202422/REQ-739-handback.md | released 0.324.3, UR-181 closed; merge ef5dc795; finalization 9bd69ce6; gate green (128 s); review 98% pass with the restatement sweep; four heavy lanes green | 2026-11-16T20:30:25Z | 15 | 1 | 0 |
