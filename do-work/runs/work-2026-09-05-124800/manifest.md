# Run manifest — work-2026-09-05-124800

Orchestrator: main checkout at /Users/t2/Desktop/e1-experimental-repos/skill-do-work2
Mode: `do-work run REQ-585 REQ-586` (serial, explicit targets; REQ-586 depends on REQ-585 and is frozen for the continuation)
Claimed REQ-585 at 1f43608b.

## Concurrency note

A sibling orchestrator (run work-2026-09-05-120117, same checkout) was mid-integration when this run started: REQ-579's merge landed at c78a0d3d one minute before the claim, and REQ-573/REQ-581/REQ-582 remain in flight there. `recover` refused with FINALIZATION-DISCOVERY-AMBIGUOUS naming that run's live paths (its REQ-579 hand-back and its modified working REQ-573 file). Judged as the sibling's live work, not dirt: left every byte in place, checked the index was empty before each commit, and continued. Every commit this run makes stages only its own declared paths.

## Builders

| REQ | Route | Operative name | Worktree | Handback |
| --- | --- | --- | --- | --- |
| REQ-585 | A | worktree-agent-REQ-585-one-scroll-surface | .git/work-run-20260905-1248/worktree-agent-REQ-585-one-scroll-surface | REQ-585-handback.md |
