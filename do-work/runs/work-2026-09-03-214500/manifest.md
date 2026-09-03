# Work run 2026-09-03-214500

Mode: `run-with-recovery --fan-out 3 --skip-impact-negligible`
Wave 1 selector order: REQ-475, REQ-483, REQ-485.
Integration owner: main checkout, serial releaser.

| Lane | REQ / description | Worktree | Handback | Status |
|---|---|---|---|---|
| lane-1 | REQ-475 — Confine all configured Memory tree readers | `worktree-agent-REQ-475-confine-all-configured-memory-tree-readers` | `REQ-475-handback.md` | implementing |
| lane-2 | REQ-483 — Bound the architecture bundle-claim loop and restore `--commit` | `worktree-agent-REQ-483-bound-architecture-bundle-claim-loop` | `REQ-483-handback.md` | parked: pending-heavy-testing at `3eb87519` |
| lane-3 | REQ-485 — Canonicalize REQ reservation marker filenames across allocation flows | `worktree-agent-REQ-485-canonicalize-req-reservation-marker-filenames` | `REQ-485-handback.md` | remediation after review; cumulative pre `6b07c546` |
