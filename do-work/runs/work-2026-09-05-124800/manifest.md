# Run manifest — work-2026-09-05-124800

Orchestrator: main checkout at /Users/t2/Desktop/e1-experimental-repos/skill-do-work2
Mode: `do-work run REQ-585 REQ-586` (serial, explicit targets; REQ-586 depends on REQ-585 and is frozen for the continuation)
Claimed REQ-585 at 1f43608b; REQ-585 held for heavy lanes at 7bd34648 with commit c08ac2b4; claimed REQ-586 at 59f169d0.

## Concurrency note

A sibling orchestrator (run work-2026-09-05-120117, same checkout) was mid-integration when this run started: REQ-579's merge landed at c78a0d3d one minute before the claim, and REQ-573/REQ-581/REQ-582 remain in flight there. `recover` refused with FINALIZATION-DISCOVERY-AMBIGUOUS naming that run's live paths (its REQ-579 hand-back and its modified working REQ-573 file). Judged as the sibling's live work, not dirt: left every byte in place, checked the index was empty before each commit, and continued. Every commit this run makes stages only its own declared paths.

## Builders

| REQ | Route | Operative name | Worktree | Handback |
| --- | --- | --- | --- | --- |
| REQ-585 | A | worktree-agent-REQ-585-one-scroll-surface | .git/work-run-20260905-1248/worktree-agent-REQ-585-one-scroll-surface | REQ-585-handback.md |
| REQ-586 | A | worktree-agent-REQ-586-top-bar-one-line | .git/work-run-20260905-1248/worktree-agent-REQ-586-top-bar-one-line | REQ-586-handback.md |

## Notes

- The repository gate for REQ-585 ran from a detached worktree at c08ac2b4 (`.git/work-run-20260905-1248/gate-c08ac2b4`) because the main tree carried a sibling session's uncommitted edits to the board package. First run red on the per-file budget of two unrelated do-work-cli finalization test files under load 23; rerun green at load 12.
- The frozen continuation argv `advance REQ-585 REQ-586 --dispatch-bound 1 ...` was refused (ADVANCE-FINALIZATION-USAGE) once REQ-585 sat in working/, because advance enters per-request mode when its first token names a non-queue REQ. Re-issued with the queued member first (`advance REQ-586 REQ-585 ...`), same set, same frozen members. Report-only finding for the CLI: a continuation whose first token has left the queue should still parse as queue mode.
