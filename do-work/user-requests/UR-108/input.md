---
id: UR-108
title: 'Capture structured lifecycle timing and critical-path reporting'
created_at: 2026-09-03T21:28:31Z
requests: [REQ-562]
word_count: 6
---

# Capture Structured Lifecycle Timing and Critical-Path Reporting

## Summary

Capture the agreed follow-up to REQ-531's timing analysis: extend REQ-448's phase milestones into structured per-run spans that attribute command, agent, waiting, retry, conflict, and finalization time and automatically report the true critical path.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-562 | Extend REQ-448 with structured lifecycle spans, command and agent attribution, and critical-path summaries |

## Batch Constraints

- This is an addendum to completed REQ-448, not a replacement for its backwards-compatible milestone timestamps or claim-to-completion calibration.
- Interoperate with REQ-539's per-test duration records instead of introducing a second test-timing store.
- Keep the first delivery to structured evidence and a textual summary; visualization and inaccessible model-internal timing are out of scope.

## Full Verbatim Input

> ```
> ok, capture it using do-work capture-request
> ```

---
*Captured: 2026-09-03T21:28:31Z*
