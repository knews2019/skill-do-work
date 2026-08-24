# do-work fan-out wave

Status: in-progress

| REQ | Builder | Operative name | Hand-back | Landed |
| --- | --- | --- | --- | --- |
| REQ-365 | `/root/build_req354` | `worktree-agent-REQ-365-a-tdd-req-must-name-a-test-file-in-its-write-set` | pending | no |
| REQ-366 | `/root/build_req360` | `worktree-agent-REQ-366-keep-dependency-gated-blocked-reqs-out-of-needs-input` | pending | no |
| REQ-367 | `/root/review_req360_premerge` | `worktree-agent-REQ-367-copy-all-reqs-per-board-column` | pending | no |

Implementation runs concurrently in isolated worktrees. Integration, verification, review, release,
and archival remain serial.
