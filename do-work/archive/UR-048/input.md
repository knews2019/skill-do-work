---
id: UR-048
title: Calibrate the estimator from archive actuals; surface negative-duration anomalies
created_at: 2026-08-17T08:05:43Z
requests: [REQ-211, REQ-212, REQ-213]
word_count: 97
---

# Calibrate the Estimator from Archive Actuals; Surface Negative-Duration Anomalies

## Summary

Apply the archive-derived calibration to the P50 estimator (188 samples, medians by route, >4h outlier rule), start recording estimate-vs-wall actuals at archive time so future re-fits need no archaeology, and make queue-kanban surface the negative-duration anomaly (completed_at earlier than claimed_at) the mining pass uncovered.

## Extracted Requests

| REQ | Title | Depends on |
|-----|-------|-----------|
| REQ-211 | Calibrate estimator scoring table to archive actuals | — |
| REQ-212 | Record estimate-vs-wall calibration log at archive time | — |
| REQ-213 | Board surfaces negative claimed→completed duration as a completion anomaly | — |

## Batch Constraints

- Calibration basis (recorded for provenance): 190 archived REQs carried both `claimed_at` and `completed_at`; 2 excluded by the >4h/negative outlier rule (assumed paused / broken stamp); medians by route — A 4.7 min (n=50), B 9.2 (n=53), C 21.4 (n=45); p80 — A 8.7, B 17.8, C 37.5. Wall≈active assumption holds for this corpus because runs were autonomous; mild upward bias acknowledged and acceptable (keeps estimates slightly conservative).
- New table: bases A=5, B=10, C=20; floor 5; weights ≈ old ÷2.5 (write-set +1/file, new file +2, subsystem +3, acceptance +1, deps-depth +2, browser +8, persistence +6, async +6, performance +4, regression +4, full-suite +4). Determinism, rounding-to-5, no-P80, and confidence-rubric structure unchanged.
- Frozen estimates on archived REQs are never rewritten by recalibration.
- The calibration log is bookkeeping appended by the work action's archive step — not a queue-kanban write surface; the board's three-write-surface sentence in the maintainer doc must NOT need amending (the board change is read-only reporting).
- The negative-duration check fires only when both stamps parse and completed < claimed; it joins the existing CompletionAnomaly strip/verify plumbing rather than inventing a parallel channel.

## Full Verbatim Input

Yes, apply it, also just kanban should report the negative duration anomaly so it can be surfaced and fixed

[Context from the same conversation: the user directed estimating calibration ratios from git history — "It should be clear on the difference between claimed at and finished at you can remove outliers cause you can assume that it was paused in between for example, an outlier would be if the difference is more than four hours" — and accepted the proposed application: bases to route medians, weights ÷2.5, provenance recorded in estimate-reference.md, plus the previously offered archive-time calibration log.]

---
*Captured: 2026-08-17T08:05:43Z*
