---
id: UR-084
title: 'Phase-level timing and a consumer-field plan-validation check'
created_at: 2026-08-31T20:38:40Z
requests: [REQ-448, REQ-449]
word_count: 103
---

# Phase-Level Timing and a Consumer-Field Plan-Validation Check

## Summary

Validated feedback proposing two independent orchestration improvements: (1) record per-phase timestamps through the work pipeline so total wall time is not mislabeled as implementation time, keeping `claimed_at` → `completed_at` for calibration; (2) add a plan-validation check that plans name the per-record fields consumed by each command whose output drives an action-owned mutation.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-448 | Record per-phase timestamps (planning, dispatch, builder handback, integration, review, remediation, re-review, release); display the phase breakdown; keep historical REQs compatible |
| REQ-449 | Consumer-field check in Route C plan validation |

## Capture-Time Decisions

The second item was reshaped with the user during capture: captured as a **warning-level** fourth plan-validation check, on the same warnings-not-blockers footing as the existing three. The original's mandatory end-to-end consumer test and the "preflight" name were explicitly excluded (the word is taken by Step 5.75), per the user's choice of "Warning-level check" over "As written".

## Full Verbatim Input

> ````text
> do-work validate-feedback: Improve do-work orchestration with phase-level timing and a Route C producer/consumer contract preflight. Record timestamps for planning, dispatch, builder handback, integration, review, remediation, re-review, and release; retain claimed_at → completed_at for calibration, but display the phase breakdown so total wall time is not mislabeled as implementation time. During plan validation, require every command whose output drives an action-owned mutation to identify the exact per-record identity, provenance, state, and outcome fields required by its consumer, plus one end-to-end consumer test. Keep historical REQs compatible when phase timestamps are absent.
>
> talk with me, does this make sense to be built using do-work capture-request?
> ````

---
*Captured: 2026-08-31T20:38:40Z*
