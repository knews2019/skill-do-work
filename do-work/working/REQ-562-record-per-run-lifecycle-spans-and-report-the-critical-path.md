---
id: REQ-562
title: 'Addendum: record lightweight per-REQ lifecycle timings'
status: claimed
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
claimed_at: 2026-09-05T00:33:24Z
---

# Addendum: Record Lightweight Per-REQ Lifecycle Timings

## What

Extend REQ-448's phase milestones with a lightweight, flat timing stream so a completed `do-work run` attributes one REQ's elapsed time to its major workflow stages and material external commands. Keep the raw stream in Git-private local state, create no instrumentation commits, and fold one concise timing summary into the REQ during finalization.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-531 took 1 hour 27 minutes, but its durable evidence could attribute that time only to broad milestone intervals. The maintainer needs command- and wait-level evidence to decide whether gates, agent latency, merge conflicts, retries, or lifecycle bookkeeping are the highest-value optimization targets.

## Context

REQ-448 already records optional flat timestamps for planning, dispatch, builder handback, integration, review, remediation, re-review, and release, renders a phase breakdown, preserves `claimed_at` to `completed_at` calibration, and keeps historical REQs compatible. This addendum deepens that observation model only enough to answer how long a serial REQ took and why; it does not discard those fields or introduce a general tracing system.

Related pending work owns adjacent optimizations: REQ-539 records per-test durations, REQ-542 moves gate execution into a background runner, and REQ-559 retries a transiently red repository gate before creating repair work. This REQ measures those paths and consumes their evidence where available; it does not reimplement them.

## Prior Implementation

REQ-448 added eight optional successful-event milestone fields across the canonical lifecycle schema, work action, board parser, duration model, generated payload, and card/detail UI. Recovery strips stale phase observations, absent fields remain backwards-compatible, and claim-to-completion stays the calibration span. Its implementation commit is `01a167697743c505e0b77e21a83db23fdc255cc0`.

Key implementation surfaces were `skills/do-work/actions/work.md`, `skills/do-work/actions/work-reference.md`, and the queue-kanban model, duration, generator, test, and card/detail files under `skills/do-work-board/tools/queue-kanban/`.

## Detailed Requirements

- Persist one flat, append-safe timing stream per run in Git-private state shared by the repository's worktrees, associated with the REQ and run identity. JSONL is an acceptable initial representation.
- Record the existing major lifecycle boundaries under stable categories: recovery/selection, planning, exploration/preflight, builder work, hand-back/merge, verification/gate, review, remediation, finalization, and cleanup. Record delegated-agent waits and canonical external commands that materially contribute wall time; do not trace file reads, individual shell primitives, or nested implementation details.
- Every event records its operation/category, wall-clock start and end, elapsed duration measured in-process when one process observes both ends, outcome, relevant revision when known, and responsible agent or process when known. Command events also record exit status and a safely redacted command identity. No parent span or nesting model is required.
- At finalization, calculate total observed elapsed time, time by category, the slowest stage and command events, and unattributed wall time. Do not calculate exclusive time, a critical path, overlap, or parallelism savings.
- Fold one compact, structured `## Timing` summary into the archived REQ in the existing finalization commit, then remove the Git-private raw stream. Instrumentation must not create an agent turn, tracked run artifact, or Git commit of its own.
- Reuse or ingest REQ-539's per-test duration evidence when available rather than introducing a second persistent test-duration system.
- Preserve REQ-448's milestone fields and `claimed_at` to `completed_at` calibration for compatibility. Historical REQs and runs with no structured timing stream continue to parse and display normally.
- Keep timing evidence bounded by recording metadata only, never command output or raw arguments that could contain secrets or user-controlled content.
- Use deterministic synthetic-clock fixtures for duration, category aggregation, command attribution, final-summary generation, raw-stream cleanup, and unattributed-time assertions.

## Constraints

- One canonical timing writer owns timestamp and duration mechanics; do not scatter `date` calls or independent timing formats through action prose and scripts.
- Use a monotonic clock where one process observes both ends of an event and UTC wall timestamps for human correlation. Do not add a resident process merely to manufacture cross-process monotonic spans.
- Concurrent runs and agents must remain separable by run, REQ, and span identity.
- A hierarchy, parent/child spans, exclusive-time accounting, critical-path analysis, parallelism metrics, Chrome Trace export, board visualization, and unavailable model-internal token/latency telemetry are explicitly out of scope.
- Test-duration ownership remains with REQ-539; retry policy remains with REQ-559; background gate scheduling remains with REQ-542.

## Builder Guidance

Certainty: Firm on the flat evidence and compact final summary, exploratory on the exact Git-private storage schema and canonical writer. Prefer one small executable primitive with deterministic tests over procedural timing prose. Optimize for zero additional agent turns and zero instrumentation-only commits.

## Red-Green Proof

**RED prompt/case:** Run a synthetic serial Route C lifecycle containing planning, exploration, a delegated builder wait, a timed gate command, review, and finalization. Today the completed REQ has milestone timestamps but no durable command timings, category totals, slowest-event report, or measure of unattributed wall time.
**Why RED now:** REQ-448 records broad phase boundaries only, so a long interval cannot distinguish agent work, material commands, waits, and lifecycle bookkeeping.
**GREEN when:** The deterministic fixture emits flat Git-private events and one final `## Timing` summary whose total observed time, category totals, slowest stage/command events, and unattributed time match the synthetic clock; finalization removes the raw stream without an instrumentation-only commit, and a historical no-trace fixture remains unchanged.
**Validation:** User narrowed the captured request after reviewing the cost evidence; serial per-REQ attribution is the goal, not general distributed tracing.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, over the 2000-token budget and `slugged: partial`; matched because this addendum extends lifecycle fields and action-owned evidence readers.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over budget and `slugged: partial`; matched because external command timing and a canonical execution primitive may touch shipped shell boundaries.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5924 tokens, over budget and `slugged: partial`; matched because the request adds structured evidence projection and lifecycle identity.

## Addendum (2026-09-04)

User added:

> ````text
> do it
> ````

- Narrow the requested instrumentation to the minimum needed to explain one REQ's serial wall time.
- Keep flat major-stage and material-command timings; write raw events to Git-private state and fold one summary into the REQ at finalization.
- Drop nested spans, parent relationships, exclusive-time and critical-path calculation, overlap and parallelism metrics, Chrome Trace output, and board visualization.
- Require zero extra agent turns, tracked run artifacts, or instrumentation-only commits.

Resolved conflict: the original request required a general nested tracing model with critical-path and parallelism analysis → the user chose lightweight serial per-REQ attribution with a flat timing stream and one final summary.

## Full Context

See `do-work/user-requests/UR-108/input.md` for the complete verbatim invocation and its capture-time summary. This addendum extends REQ-448 rather than replacing it.

---
*Source: "ok, capture it using do-work capture-request" — referring to the agreed structured lifecycle timing and critical-path proposal.*
