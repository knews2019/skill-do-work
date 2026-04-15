---
id: 0001
title: Capture ≠ Execute Boundary
status: accepted
decided: 2026-02-03
version: 0.10.0
topic: philosophy
supersedes: []
superseded_by: null
related:
  - adr: 0002
    rel: depends-on
  - adr: 0010
    rel: complements
---

# ADR-0001: Capture ≠ Execute Boundary

## Context

In early versions of the skill, the capture action (then named `do`) would often slide directly into execution. An agent would write REQ files for the user's request, then — helpful as always — say "let me go ahead and start building that for you," and begin coding before the user had a chance to review anything.

This collapsed the review window the capture phase was designed to create. The user never saw the REQs before they were already being worked on. Any ambiguity the agent papered over with a guess became baked-in scope. The resulting workflow was fundamentally indistinguishable from "just build it" — the capture phase was decorative.

The signal cost was real: captured but unreviewed REQs produced work that silently drifted from what the user actually wanted, and the drift was hard to notice because the REQ file looked fine in isolation.

## Decision

The capture action **stops** after writing UR and REQ files and reporting back. It does not automatically invoke the work action. The user must explicitly invoke `do work run` (or any other work-triggering verb) to execute the queue.

The only exception is when the user explicitly asks for capture **and** execution in the same invocation (e.g., "capture this and build it") — a conscious override, not a default.

## Alternatives Considered

- **Auto-execute after capture.** The fast, "helpful" default. Rejected — it collapses the review window and eliminates the distinction between the two phases.
- **Prompt user to confirm execution.** Softer than a hard stop, but still interrupts the capture phase's "write things down" focus, and the prompt is easy to say "yes" to by reflex.
- **Merge capture into work entirely.** Drop the two-phase model. Rejected — two phases exist precisely to slow down the transition from request to code so assumptions surface before they cost anything.

## Consequences

- **Two-phase workflow is now structural.** The queue accumulates REQs between capture sessions; the user drains it on their schedule.
- **Users control execution timing.** Useful when the user wants to capture several related items before letting the agent work on any of them.
- **Ambiguity has a place to go.** Captured REQs can sit in the queue while the user thinks, and the verify-requests action can run against them before any code is written.
- **Cost: one extra command to run.** Users who want a one-liner "capture and build" experience have to explicitly ask for it.

## References

- **CHANGELOG**: v0.10.0 — The Hard Stop (2026-02-03), v0.17.0 (renamed `do` → `capture`)
- **Action files**: `actions/capture.md` (STOP section), `SKILL.md` ("Capture ≠ Execute" core concept)
- **Related ADRs**: [[0002-ur-req-pairing]] (the boundary only works if capture actually produces something persistent), [[0010-reqs-as-validated-intent]] (the review window exists so REQs can be validated)
