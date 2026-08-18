# Run Manifest — work-2026-08-18-124358

Orchestrator: t2s-Virtual-Machine:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2
Worktrees parent: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees

Two waves. REQ-241 and REQ-242 are serialised deliberately: both write
`web/board-durations.js`, and REQ-241 may move `DURATIONS_LABEL_ROW_HEIGHT`,
which shifts every panel below it — including the Panel B title REQ-242 is
separating an annotation from. REQ-242's collision is re-measured against the
integrated REQ-241 result before it is dispatched, and closed as
resolved-by-REQ-241 if it no longer reproduces.

Auto-wave (`do-work run --fan-out 5`) is deliberately NOT used: it computes its
set from `depends_on`, claim state and `assigned_to` and does not read
`write_set`, so it would run REQ-241 and REQ-242 together.

## Wave 1 — dispatched

| REQ | Route | Est (p50 min) | Operative name | Handback | Landed |
| --- | --- | --- | --- | --- | --- |
| REQ-241 | B | 35 | worktree-agent-REQ-241-reconcile-durations-label-metrics-with-the-rendered-face | REQ-241-handback.md | pending |
| REQ-243 | B | 20 | worktree-agent-REQ-243-check-that-shipped-markdown-pointers-resolve | REQ-243-handback.md | pending |
| REQ-245 | A | 5 | worktree-agent-REQ-245-name-fabricated-stamps-in-the-boards-future-stamp-warnings | REQ-245-handback.md | pending |

Wave 1 is disjoint by inspection: durations chart (241), one shell test script
(243), queue-kanban model/verify plus two named test files (245). REQ-245's
`*_test.go` glob is narrowed in its brief to `timestamp_test.go` and
`completion_anomaly_test.go` so it cannot collide with `generate_test.go`.

## Wave 2 — not yet dispatched

| REQ | Route | Est (p50 min) | Operative name | Handback | Landed |
| --- | --- | --- | --- | --- | --- |
| REQ-242 | — | — | — | — | not dispatched |
| REQ-244 | — | — | — | — | not dispatched |

REQ-244 has no `write_set` by design — its requirement *is* the sweep. Its
builder declares its own scope at Step 5.5 before writing anything.
If REQ-243 lands in wave 1, its link checker is live when REQ-244 rewrites
shipped markdown; a new failure there reads as the checker working.

## Out of scope for this run

REQ-246 and REQ-247 (UR-056, captured 2026-08-18T12:38:26Z) are pending in the
queue and were not part of this run's assignment. REQ-247 depends on REQ-246.
