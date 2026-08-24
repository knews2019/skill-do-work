# do-work fan-out wave

Status: in-progress

| REQ | Builder | Operative name | Hand-back | Landed |
| --- | --- | --- | --- | --- |
| REQ-365 | `/root/build_req354` | `worktree-agent-REQ-365-a-tdd-req-must-name-a-test-file-in-its-write-set` | ready (`b7d2362`) | no |
| REQ-366 | `/root/build_req360` | `worktree-agent-REQ-366-keep-dependency-gated-blocked-reqs-out-of-needs-input` | ready (`b909815`) | no |
| REQ-367 | `/root/review_req360_premerge` | `worktree-agent-REQ-367-copy-all-reqs-per-board-column` | pending | no |
| REQ-371 | `/root/build_req354` | `worktree-agent-REQ-371-keep-timeline-bars-inside-the-plot-after-the-drawer-opens` | pending | no |

Implementation runs concurrently in isolated worktrees. Integration, verification, review, release,
and archival remain serial.
