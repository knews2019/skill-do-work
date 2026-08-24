# do-work fan-out wave

Status: in-progress

| REQ | Builder | Operative name | Hand-back | Landed |
| --- | --- | --- | --- | --- |
| REQ-365 | `/root/build_req354` | `worktree-agent-REQ-365-a-tdd-req-must-name-a-test-file-in-its-write-set` | ready (`b7d2362`) | yes (`6265f1c`) |
| REQ-366 | `/root/build_req360` | `worktree-agent-REQ-366-keep-dependency-gated-blocked-reqs-out-of-needs-input` | ready (`1961590`) | yes (`c18deb8`) |
| REQ-367 | `/root/review_req360_premerge` | `worktree-agent-REQ-367-copy-all-reqs-per-board-column` | ready (`5dea6e5`) | yes (`66e8de5`) |
| REQ-371 | `/root/build_req354` | `worktree-agent-REQ-371-keep-timeline-bars-inside-the-plot-after-the-drawer-opens` | ready (`d7cef79`) | yes (`5fda05b`) |
| REQ-368 | `/root/build_req360` | `worktree-agent-REQ-368-ur-copy-all-includes-its-reqs` | ready (`cedaa14`) | yes (`61dddb2`) |
| REQ-369 | `/root/build_req354` | `worktree-agent-REQ-369-wait-for-the-timeline-table-rebuild-condition` | ready (`3a114d7`) | yes (`ed15507`) |
| REQ-370 | `/root/build_req354` | `worktree-agent-REQ-370-restore-a-falsifiable-timeline-pointer-capture-mutation` | ready (`dc75446`) | yes (`46ed690`) |

Implementation runs concurrently in isolated worktrees. Integration, verification, review, release,
and archival remain serial.
