# Run manifest — work-2026-09-05-120117

Orchestrator: main checkout at /Users/t2/Desktop/e1-experimental-repos/skill-do-work2
Mode: `do-work run --fan-out 3`
Wave 1 claimed at cd63690f: REQ-572, REQ-575, REQ-578.

## Pre-build gate

- `bash _dev/tests/maintainer-verify.sh` exited 1 twice at cd63690f (mandatory one retry). Fingerprint: shipped-package-reference-contract, 3 broken blob targets in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` pointing at the pre-UR-098 flat archive paths for REQ-503/504/505. No claimed REQ owns those bytes.
- Resolved in place as standalone maintainer commit `09a13839` (release 0.294.1), the 2d140f63 shape, instead of a minted repair REQ. Rationale recorded in that commit message.

## Builders

| REQ | Route | Operative name | Worktree | Handback |
| --- | --- | --- | --- | --- |
| REQ-572 | B | worktree-agent-REQ-572-activity-rows | .git/work-run-20260905-1201/worktree-agent-REQ-572-activity-rows | REQ-572-handback.md |
| REQ-575 | A | worktree-agent-REQ-575-append-only-stamps | .git/work-run-20260905-1201/worktree-agent-REQ-575-append-only-stamps | REQ-575-handback.md |
| REQ-578 | A | worktree-agent-REQ-578-findings-strip | .git/work-run-20260905-1201/worktree-agent-REQ-578-findings-strip | REQ-578-handback.md |

## Notes

- REQ-572 arrived `gate_deferred: true` with a saved implementation pair `7ad53bff..fbdcd35e`. The saved-range resume proof was rejected: three later commits (a55f24ce, 4c76c332 for REQ-571; ae184a7b for REQ-576) touch protected paths. The pair was deleted from the claimed record, prior qualification/testing/review evidence was demoted, and the REQ returned to Step 6.
- Two leftover worktrees from run work-2026-09-05 (REQ-506, REQ-577) predate this run and belong to `do-work cleanup` Pass 5.
