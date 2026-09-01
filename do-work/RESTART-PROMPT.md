```
do-work run --fan-out 2
This command is sufficient; everything below it is context.
```

---

## Reference

- **Integration state:** `main` is at `82534d36` before this handoff commit. Suite version is `0.256.1`. REQ-418 builder `a7c975c5` merged as `94560fde`; its sole remediation `a43b2587` merged as `82534d36`. Full builder range: `b713fa8b..a43b2587`.
- **In-flight REQ:** REQ-418 remains `claimed` in `do-work/working/REQ-418-migrate-toolbox-absorb-audit-metrics.md`. Implementation and remediation are merged and verified. Initial review `do-work/runs/work-2026-08-31-165510/REQ-418-review.md` recorded nine Important findings; remediation closed all nine in named tests. The required fresh re-review was interrupted for this handoff before it wrote `REQ-418-rereview.md`.
- **What remains, in order:** (1) assign a new independent reviewer to reproduce all nine closure cases on `82534d36`; (2) write `do-work/runs/work-2026-08-31-165510/REQ-418-rereview.md`; (3) if clean, record REQ evidence, complete/archive with implementation hash `82534d36`, release, and clean the worktree; (4) if any Important residual remains, remediation is exhausted, so mark `completed-with-issues`, route every residual via fold-first, then archive/release; (5) run the selector again with fan-out 2.
- **Verification:** remediation passed focused and race tests for toolboxcommands plus gittransaction; full do-work-cli tests and vet; exact Go 1.25; Windows compile/fail-closed behavior; standalone audit valid/malformed differential; retained toolbox fixtures and all 118 prescribed-shell cases; contracts, staged skills, install/update; exact documented 32-path ceiling; diff checks; and final unpiped `bash _dev/tests/maintainer-verify.sh`. Only the optional browser lane skipped because no browser was configured.
- **Worktree verdict:** ACTIVE — `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-418-toolbox-migration`, branch `worktree-agent-REQ-418-toolbox-migration`, HEAD `a43b2587`. The branch is merged and the worktree is clean, but the claim is not archived. Do not remove it until the next session completes REQ-418, then re-check merge ancestry, cleanliness, and archival before removal.
- **Uncommitted REQ-418 evidence:** `do-work/runs/work-2026-08-31-165510/REQ-418-exploration.md`, `REQ-418-plan.md`, `REQ-418-handback.md`, `REQ-418-review.md`, and `REQ-418-remediation-handback.md` are intentional untracked scratch. There is no `REQ-418-rereview.md` yet. Never stage the run directory.
- **Queue state:** 31 pending, 0 pending-answers, 0 blocked, 1 in progress. No `do-work clarify` step is needed. REQ-419 waits on REQ-418; REQ-420 waits on REQ-419. Those dependency gates encode the critical path. After REQ-418 closes, fan-out 2 remains the useful bound for the selector's dependency-ready pair.
- **Concurrent main commits:** `593c5145` archived REQ-454, REQ-455, REQ-464, and REQ-465 after queue analysis and folded their work into REQ-420/REQ-460; `19a75452` published its report. These commits are user-owned, already preserved on main, and unrelated to REQ-418's 32 product paths.
- **First-ten-minutes heads-up:** do not rebuild or run a second remediation. Start at fresh independent re-review. The initial report is not the terminal verdict, and the two interrupted re-review agents wrote no report. Keep the existing scratch untracked and preserve unrelated user-owned commits.
