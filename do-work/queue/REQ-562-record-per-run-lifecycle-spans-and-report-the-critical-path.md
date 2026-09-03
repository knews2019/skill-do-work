---
id: REQ-562
title: 'Addendum: record per-run lifecycle spans and report the critical path'
status: pending
created_at: 2026-09-03T21:28:31Z
user_request: UR-108
addendum_to: REQ-448
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-448, REQ-531, REQ-539, REQ-542, REQ-559]
---

# Addendum: Record Per-Run Lifecycle Spans and Report the Critical Path

## What

Extend REQ-448's phase milestones with structured per-run span evidence so a completed `do-work run` attributes elapsed time to workflow phases, commands, agents, waits, retries, conflicts, and finalization. Produce an automatic textual summary of the critical path without double-counting concurrent work.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-531 took 1 hour 27 minutes, but its durable evidence could attribute that time only to broad milestone intervals. The maintainer needs command- and wait-level evidence to decide whether gates, agent latency, merge conflicts, retries, or lifecycle bookkeeping are the highest-value optimization targets.

## Context

REQ-448 already records optional flat timestamps for planning, dispatch, builder handback, integration, review, remediation, re-review, and release, renders a phase breakdown, preserves `claimed_at` to `completed_at` calibration, and keeps historical REQs compatible. This addendum deepens that observation model; it does not discard those fields.

Related pending work owns adjacent optimizations: REQ-539 records per-test durations, REQ-542 moves gate execution into a background runner, and REQ-559 retries a transiently red repository gate before creating repair work. This REQ measures those paths and consumes their evidence where available; it does not reimplement them.

## Prior Implementation

REQ-448 added eight optional successful-event milestone fields across the canonical lifecycle schema, work action, board parser, duration model, generated payload, and card/detail UI. Recovery strips stale phase observations, absent fields remain backwards-compatible, and claim-to-completion stays the calibration span. Its implementation commit is `01a167697743c505e0b77e21a83db23fdc255cc0`.

Key implementation surfaces were `skills/do-work/actions/work.md`, `skills/do-work/actions/work-reference.md`, and the queue-kanban model, duration, generator, test, and card/detail files under `skills/do-work-board/tools/queue-kanban/`.

## Detailed Requirements

- Persist one structured timing stream per run, associated with the REQ and run identity. JSONL is an acceptable initial representation; the builder may choose another append-safe structured format if it preserves the same evidence.
- Record spans for workflow phases, external commands, delegated agents, explicit waiting, retries, merge/conflict handling, verification, review, release/finalization, and cleanup whenever those events are observable.
- Every span records a stable span id, optional parent span id, operation/category, responsible agent or process when known, wall-clock start/end timestamps, monotonic duration, outcome, and relevant revision. Command spans also record exit status and a safely redacted command identity.
- Calculate total elapsed time, exclusive time, critical-path time, parallelism saved, slowest spans, and unattributed time without summing overlapping child or agent spans twice.
- Emit a concise timing summary automatically when the run completes and retain the structured evidence for later comparison across runs.
- Reuse or ingest REQ-539's per-test duration evidence when available rather than introducing a second persistent test-duration system.
- Preserve REQ-448's milestone fields and `claimed_at` to `completed_at` calibration for compatibility. Historical REQs and runs with no structured timing stream continue to parse and display normally.
- Keep stored timing evidence bounded and redact or hash command arguments that could contain secrets or user-controlled content.
- Use deterministic synthetic-clock fixtures for duration, nesting, overlap, and critical-path assertions.

## Constraints

- One canonical span writer owns timestamp and duration mechanics; do not scatter `date` calls or independent timing formats through action prose and scripts.
- Use a monotonic clock for durations and UTC wall timestamps for human correlation.
- Concurrent runs and agents must remain separable by run, REQ, and span identity.
- The first delivery is structured evidence plus a textual summary. A Chrome Trace export, board visualization, and unavailable model-internal token/latency telemetry are explicitly out of scope.
- Test-duration ownership remains with REQ-539; retry policy remains with REQ-559; background gate scheduling remains with REQ-542.

## Builder Guidance

Certainty: Firm on the evidence and report outcomes, exploratory on the exact storage schema and canonical writer. Prefer a small executable primitive with deterministic tests over procedural timing prose.

## Red-Green Proof

**RED prompt/case:** Run a synthetic Route C lifecycle containing a timed command, two overlapping child spans, an agent wait, and a retry. Today the completed REQ has milestone timestamps, but no durable command/agent span stream and no way to distinguish summed work from the wall-clock critical path.
**Why RED now:** REQ-448 records phase boundaries only. It cannot attribute time inside a phase, represent overlap, calculate exclusive or critical-path duration, or separate productive work from waiting and rework.
**GREEN when:** The same deterministic fixture emits structured nested spans and a completion summary whose total elapsed, exclusive time, overlap savings, critical path, slowest spans, and unattributed time match the synthetic clock exactly; a historical no-trace fixture still parses and displays unchanged.
**Validation:** User confirmed by approving the proposed capture after reviewing the intended measurements, scope, and related work.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, over the 2000-token budget and `slugged: partial`; matched because this addendum extends lifecycle fields and action-owned evidence readers.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over budget and `slugged: partial`; matched because external command timing and a canonical execution primitive may touch shipped shell boundaries.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5924 tokens, over budget and `slugged: partial`; matched because the request adds structured evidence projection and lifecycle identity.

## Full Context

See `do-work/user-requests/UR-108/input.md` for the complete verbatim invocation and its capture-time summary. This addendum extends REQ-448 rather than replacing it.

---
*Source: "ok, capture it using do-work capture-request" — referring to the agreed structured lifecycle timing and critical-path proposal.*
