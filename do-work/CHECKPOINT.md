---
session_ended: 2026-08-17T08:35:00Z
last_completed: REQ-215
queue_state: 0 pending, 4 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 5
session_depth: moderate
---

# Session Checkpoint

## Completed This Session

- REQ-211: Calibrate estimator scoring table to archive actuals (Route B, 0.196.0)
- REQ-212: Record estimate-vs-wall calibration log at archive time (Route B, 0.197.0)
- REQ-213: Board surfaces negative claimed→completed duration anomaly (Route B, 0.198.0; review 94% Pass, 2 important → REQ-214/215)
- REQ-214: verify surfaces completion anomalies as findings (Route B, 0.199.0)
- REQ-215: Sync completion-anomaly prose with the reversed-span class (Route A sweep, 0.199.1)

## In Progress (interrupted)

## Still Queued

- REQ-203: Harden presentation target-ID source-seam tests (pending-answers)
- REQ-204: Harden ai-report generated-batch lifecycle (pending-answers)
- REQ-205: Make portfolio publication independent and exact (pending-answers)
- REQ-206: Finish active publication delegation (pending-answers)

## Session Notes

- UR-048 archived: the estimator now prices from measured history (bases = archive route medians 5/10/20, floor 5, weights ÷2.5; provenance in `estimate-reference.md` → Calibration), every archived estimated REQ self-records into `do-work/calibration-log.tsv` (7 lines seeded), and the board flags reversed claimed→completed spans.
- **`queue-kanban verify` now exits 1 on this tree by design**: it surfaces 10 pre-existing completion anomalies — REQ-091's reversed span (repairable stamp) and 9 archived REQs whose `commit:` hashes git cannot date (likely pre-rewrite history). Repairing archived frontmatter is an owner decision; until then forensics/verify consumers will see these findings.
- Anomaly prose now reads true for all four classes; the never-silent warning routes to per-class fix text.
- UR-042's REQ-203–206 remain `pending-answers` (run `do-work clarify`).
- maintainer-verify still carries the 41 container-environment FAILs (missing `just`, tar/gzip exec, stat probes); FAIL-set-vs-baseline was the regression gate all session.
