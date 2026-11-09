# Run Manifest — 2026-11-09-layout-queue

Run dir: do-work/runs/2026-11-09-layout-queue/
Concurrency: `do-work run` on the five-REQ folder-layout queue (serial loop, dispatch bound 1). The orchestrator plays every orchestrator role; each REQ is built in its own worktree.
Test policy: focused probe per REQ, then `bash _dev/tests/maintainer-verify.sh` at the merge revision from a detached checkout, then the planned heavy lanes.
Status: running
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant |
|-----|---------|----------------|---------------|--------|------------------|
| REQ-713 | orchestrator as builder (Route A, worktree) | worktree-agent-REQ-713-folder-structure-check-default-param-recipes | do-work/runs/2026-11-09-layout-queue/REQ-713-handback.md | merge 218404ff, review-fix re-merge 013ea54e; gate green x2 (125s, 122s); review 92% then re-check 96%; four heavy lanes executed green; released 0.314.1, UR-165 closed | 2026-11-09T21:52:25Z |
| REQ-714 | builder agent (Route C, worktree) | worktree-agent-REQ-714-suite-writes-kb-and-decisions-under-docs | do-work/runs/2026-11-09-layout-queue/REQ-714-handback.md | dispatched | 2026-11-09T22:20:14Z |
