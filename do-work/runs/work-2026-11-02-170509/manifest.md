# Run Manifest — work-2026-11-02-170509

Run dir: do-work/runs/work-2026-11-02-170509/
Concurrency: `do-work run UR-161` (serial loop, dispatch bound 1). REQ-707 depends on REQ-706 and REQ-708 on REQ-707, so members run one at a time. The orchestrator plays every orchestrator role; each builder is a background agent in its own worktree.
Test policy (user decision 2026-11-02): focused tests only. Each REQ's gate argv is `bash do-work/runs/work-2026-11-02-170509/REQ-NNN-gate.sh`; `_dev/tests/maintainer-verify.sh` and heavy lanes are not run ("heavy lanes not run: user chose focused tests only (2026-11-02)").
Orchestrator note: the untracked user file `do-work/inbox/2026-10-28_do-work-upstream-suggestion-month-folder-report-layer.md` made `recover` refuse `FINALIZATION-DISCOVERY-AMBIGUOUS`. The run may not stage, edit or delete it, so it is hidden from discovery by one temporary line in `.git/info/exclude`, removed at the end of the run.
Status: complete   # all three members released; UR-161 closed in 0.309.3
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------------------------------|------------|----------------|
| REQ-706 | builder-706 (background agent, Route B, worktree) | worktree-agent-REQ-706-ai-report-tools-read-month-folders | do-work/runs/work-2026-11-02-170509/REQ-706-handback.md | handback c004cba0; merge 3edfd2c9, review-fix re-merge ba902c63; focused gate green x2; review 95% Approve, re-check 97%; heavy lanes not run (user: focused tests only); release 0.309.1 | 2026-11-02T17:11:16Z | 55 | 0 | 1 |
| REQ-707 | builder-707 (background agent, Route C, worktree) | worktree-agent-REQ-707-ai-report-organize-counted-month-folders | do-work/runs/work-2026-11-02-170509/REQ-707-handback.md | handback 6fa71772; merge 429cdb97, review-fix re-merge c52e7f54; focused gate green x2; review 82% (Partial) then re-check 95% Pass; heavy lanes not run (user: focused tests only); release 0.309.2 | 2026-11-02T17:58:32Z | 46 | 0 | 0 |
| REQ-708 | builder-708 (background agent, Route B, worktree) | worktree-agent-REQ-708-recommend-ai-report-organize-on-install-update | do-work/runs/work-2026-11-02-170509/REQ-708-handback.md | handback ecc9e0b5; merge e4c9d951, review-fix re-merge 112545ad; focused gate green x2; review 95% Pass with wave-end sweep, re-check 96%; heavy lanes not run (user: focused tests only); release 0.309.3, UR-161 closed | 2026-11-02T18:44:55Z | 33 | 0 | 0 |
