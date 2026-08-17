---
id: UR-047
title: P50 active-duration estimation for REQs
created_at: 2026-08-16T23:52:07Z
requests: [REQ-208, REQ-209, REQ-210]
word_count: 811
---

# P50 Active-Duration Estimation for REQs

## Summary

Add an informational P50 estimate of active agent wall-time to REQs: a deterministic, explainable estimator computes `p50_active_minutes` (rounded to 5, floored, with confidence and a `basis` list), the work action ensures/prints the estimate before planning and freezes it once execution begins, verify-requests recalculates after material repairs, and multi-REQ output distinguishes total effort from critical-path duration. Estimation never blocks execution. No P80, no calendar promises.

## Extracted Requests

| REQ | Title | Depends on |
|-----|-------|-----------|
| REQ-208 | Deterministic P50 estimator: script, reference file, and schema | — |
| REQ-209 | Wire P50 estimation into the work action | REQ-208 |
| REQ-210 | Recalculate P50 estimates in verify-requests after material repairs | REQ-208 |

## Batch Constraints

Agreed in the capture session (analysis validated against the pipeline before capture):

- **Determinism via a program, not prose.** The agent extracts judgment signals; a shipped script maps signals → minutes. "Same result for the same normalized inputs" holds over the extracted signals, and `_dev/tests/` lock-in tests test the script. Rationale: an agent following prose is not deterministic and cannot be regression-tested.
- **Route is assigned at work Step 3 (triage), not at capture.** The spec's lifecycle rule 1 (capture calculates an initial estimate using route) is amended: v1 does **not** wire estimation into capture. The spec's own legacy rule becomes the universal rule — an estimate is derived when the REQ is next verified or selected for work. This keeps the interactive capture window fast and gives the estimator its strongest signal (the real route).
- **Primary wire point is work, immediately after triage, before planning.** Satisfies "print the estimate before planning starts" and gives legacy REQs an estimate for free at first selection.
- **`effort_estimate` stays a two-value triage chip.** The `estimate:` block is a separate field; the bridge is `effort_estimate: trivial` ⇒ emit the floor estimate via a cheap short-circuit (no reference-file load). The `work-reference.md` schema comment fencing off "an estimation system" must be amended in the same commit that adds the `estimate:` block (co-location rule).
- **`actual_active_minutes` and history-based calibration are OUT of v1.** do-work has no pause-aware execution timing; `claimed_at`/`completed_at` subtraction is forbidden by the spec itself. Deferred until pause-aware instrumentation exists. v1 is a pure-prior estimator and documents that honestly.
- **Board display is OUT of v1** (would trigger the parser lock-step rule; separate REQ later if wanted). Verified: the board's yaml.v3 permissive-map parser tolerates the nested `estimate:` block on cards today — backwards compatibility is free.
- **Nothing lands in queue-kanban.** The board tool's three-write-surface rule stays intact; the estimator script ships in the do-work skill's own tools.
- Frozen once execution begins; never rewritten with knowledge gained during implementation.
- Estimation must never block execution or require user clarification.
- No P80 or other percentile fields; no calendar-time promises.

## Full Verbatim Input

talk with me, does this make sese?

consider using do-work validate-feedback, but I want also an analysis of the feasibility, how much slower is this? Where should we wire it in, perhaps in lazy loading mode to get this information, then again if it's just 1% of total execution time, then perhaps we can just add it for every req.

See below the `Add P50 active-duration estimation to do-work REQs.`

# Add P50 active-duration estimation to do-work REQs.

Objective

Before a REQ begins execution, calculate and display a P50 estimate of the active agent wall time required to complete it.

Definition

p50_active_minutes is the median estimate of active wall-clock minutes while do-work and its agents are working, including:

- planning and exploration;
- implementation;
- tests, browser automation, builds, and other tool execution;
- independent review;
- the expected cost of ordinary remediation.

Exclude:

- time waiting for user input or approval;
- paused, suspended, or stopped sessions;
- overnight/user-controlled gaps;
- queue wait time;
- calendar completion dates.

This is an informational forecast, not a deadline or execution budget.

REQ metadata

Add backwards-compatible optional estimation metadata:

estimate:
  p50_active_minutes: 75
  confidence: medium
  calculated_at: 2026-08-16T12:00:00Z
  basis:
    - Route C
    - 12-file write set
    - real-browser acceptance
    - persistence changes
    - full-suite qualification

Do not add P80 or other percentile estimates.

Estimation lifecycle

1. Capture should calculate an initial estimate after the REQ has requirements, route, dependencies, acceptance criteria, and expected scope.
2. verify-requests should recalculate it after repairing or materially changing a REQ.
3. work should ensure an estimate exists immediately before claiming the REQ.
4. Freeze the estimate once execution begins. Do not rewrite it using knowledge gained during implementation.
5. Existing REQs without estimation metadata must remain valid. Derive and persist an estimate when they are next verified or selected for work.
6. Estimation must never block execution or require user clarification.

Estimation inputs

Build an explainable deterministic estimator using signals already present in the REQ and repository:

- Route A, B, or C;
- write-set size and file types;
- number of runtime subsystems involved;
- new files, assets, or dependencies;
- acceptance-criteria count;
- dependency depth and serialization;
- browser, responsive, visual, accessibility, or screenshot requirements;
- persistence, migration, API, or schema changes;
- asynchronous lifecycle, teardown, race, or retry behavior;
- performance instrumentation;
- focused versus full-suite verification;
- lint, deploy, asset-integrity, and cross-route regression requirements;
- independent review and expected remediation cost.

Use historical completed REQs when available for calibration, but never let missing history prevent estimation.

Avoid false precision:

- round estimates to the nearest five minutes;
- enforce a reasonable minimum;
- attach low, medium, or high confidence;
- record the dominant sizing factors in basis;
- produce the same result for the same normalized REQ inputs.

Presentation

When work selects a REQ, print the estimate before planning starts:

Starting REQ-1459 — Add SD-39 review links and QA gates
Estimated active duration: approximately 125 minutes (P50, medium confidence)
Dominant factors: Route C, browser evidence matrix, performance gate, storage auditing

For a UR or queue containing multiple REQs:

- show the P50 estimate for each REQ;
- show total estimated agent effort;
- show critical-path active time using the dependency graph rather than simply summing parallel branches;
- clearly label both values.

Example:

REQ-1454  75 min
REQ-1455  70 min  depends on REQ-1454
REQ-1456  80 min  depends on REQ-1455
REQ-1459  125 min depends on prior loop REQs

Total estimated effort: 350 active minutes
Estimated critical path: 350 active minutes

Actual-time calibration

If do-work already has reliable pause-aware execution timing, record this separately after completion:

actual_active_minutes: 128

Never derive actual time by simply subtracting claimed_at from completed_at, because that includes user pauses and suspended sessions.

Use completed estimates and actuals to improve future calibration without modifying historical estimates.

Acceptance criteria

- Every newly captured REQ receives p50_active_minutes before it can run.
- verify-requests refreshes estimates when scope changes materially.
- work displays the estimate before claiming/executing a REQ.
- Old REQs remain valid.
- Estimates are deterministic, rounded, explainable, and non-blocking.
- User wait and suspended time are excluded from the definition.
- Dependency-aware UR output distinguishes total effort from critical-path duration.
- No P80 field or calendar-time promise is introduced.
- Automated tests cover small Route A, focused Route B, integrated Route C, browser-heavy QA, dependency graphs, backwards compatibility, and deterministic output.
- Documentation explains that P50 means roughly a 50% chance of completing within the estimated active minutes.

[Follow-up message, same session:]

Do I understand correctly that you find Find the correctly that verify request is the correct target for this estimation

[Agent clarified: primary wire point is work (post-triage), verify-requests is the secondary recalculation point, capture demoted/dropped for v1. User then instructed:]

OK, capture the request for this and also estimate the time of each request that you capture use to work capture request

---
*Captured: 2026-08-16T23:52:07Z*
