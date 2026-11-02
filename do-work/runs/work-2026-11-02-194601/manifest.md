# Run Manifest — work-2026-11-02-194601

Run dir: do-work/runs/work-2026-11-02-194601/
Concurrency: `do-work run` on a one-REQ queue (serial loop, dispatch bound 1). The orchestrator plays every orchestrator role; the builder is a background agent in its own worktree.
Test policy (user decision 2026-11-02): focused tests only. The gate argv is `bash do-work/runs/work-2026-11-02-194601/REQ-709-gate.sh`; `_dev/tests/maintainer-verify.sh` and heavy lanes are not run ("heavy lanes not run: user has not granted heavy-lane permission for this run; focused tests only (2026-11-02)").
Status: in-progress
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant | Takeover-to-finalization min | Full gates | Stall restarts |
|-----|---------|----------------|---------------|--------|------------------|------------------------------|------------|----------------|
| REQ-709 | builder-709 (background agent, Route A, worktree) | worktree-agent-REQ-709-recover-ignores-unowned-do-work-notes | do-work/runs/work-2026-11-02-194601/REQ-709-handback.md | handback 4634893c (2026-11-02T20:01:32Z) | 2026-11-02T19:51:31Z | | 0 | 0 |
