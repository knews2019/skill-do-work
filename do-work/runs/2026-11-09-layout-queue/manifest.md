# Run Manifest — 2026-11-09-layout-queue

Run dir: do-work/runs/2026-11-09-layout-queue/
Concurrency: `do-work run` on the five-REQ folder-layout queue (serial loop, dispatch bound 1, for REQ-713 to REQ-715; then `do-work run --fan-out 2` for REQ-716 and REQ-717: parallel builders, serial integration). The orchestrator plays every orchestrator role; each REQ is built in its own worktree.
Test policy: focused probe per REQ, then `bash _dev/tests/maintainer-verify.sh` at the merge revision from a detached checkout, then the planned heavy lanes.
Status: running
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant |
|-----|---------|----------------|---------------|--------|------------------|
| REQ-713 | orchestrator as builder (Route A, worktree) | worktree-agent-REQ-713-folder-structure-check-default-param-recipes | do-work/runs/2026-11-09-layout-queue/REQ-713-handback.md | merge 218404ff, review-fix re-merge 013ea54e; gate green x2 (125s, 122s); review 92% then re-check 96%; four heavy lanes executed green; released 0.314.1, UR-165 closed | 2026-11-09T21:52:25Z |
| REQ-714 | builder agent (Route C, worktree) | worktree-agent-REQ-714-suite-writes-kb-and-decisions-under-docs | do-work/runs/2026-11-09-layout-queue/REQ-714-handback.md | handback 00eac450; merge fae09f99, review-fix re-merge 409a3e58; gate green x2 (142s, 146s); review 89% then re-check 95%; six heavy lanes executed green; released 0.315.0, UR-166 closed | 2026-11-09T22:20:14Z |
| REQ-715 | orchestrator as builder (Route B, worktree) | worktree-agent-REQ-715-changelog-archives-to-docs-changelog | do-work/runs/2026-11-09-layout-queue/REQ-715-handback.md | handback 742256f3; merge f2ccb2a7, review-fix re-merge 1a000d08; gate green x2 (130s, 141s); review 82% then re-check 96%; six heavy lanes executed green; release 0.315.1, UR-167 closes (tail finished by the resumed session) | 2026-11-09T23:13:02Z |
| REQ-716 | builder agent (Route B, worktree, wave 2 of 2 parallel) | worktree-agent-REQ-716-installer-scripts-tools-to-cmd | do-work/runs/2026-11-09-layout-queue/REQ-716-handback.md | pre-dispatch revised after the maintainer chose to delete the root copies (D-04); dispatched from dc5fd2df after REQ-717 released | 2026-11-10T00:58:12Z |
| REQ-717 | builder agent (Route B, worktree, wave 2 of 2 parallel) | worktree-agent-REQ-717-justfile-recipes-call-scripts | do-work/runs/2026-11-09-layout-queue/REQ-717-handback.md | handback b8b07228; merge 36012976; gate green (162s); review 95% Pass, no re-merge; six heavy lanes executed, installer red once under load then green alone; released 0.315.1, UR-169 closed | 2026-11-10T00:35:00Z |
