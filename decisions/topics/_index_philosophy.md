# Topic: Philosophy

The mental models the rest of the skill hangs on. These ADRs aren't structural — they're about *what the skill is trying to be*. Remove them and the other decisions stop cohering.

## Core Framing

> The skill produces a **trail of intent**, not just code. Capture elicits and validates user intent; the queue holds validated intent; the builder implements against validated intent; review checks that built code matches captured intent; archive preserves the intent trail so later readers can reconstruct *why*.

Two framing commitments make this possible:

1. **A hard boundary between capture and execution.** Capture slows things down on purpose — it asks questions, surfaces assumptions, writes REQ files. Only when the user explicitly invokes the work action does anything get built. Without this boundary, capture becomes execution and the intent trail becomes just commit history.
2. **REQs are validated intent, not drafts.** A captured REQ represents something the user has seen and implicitly agreed to. Addenda can extend it, but silent contradictions are forbidden — the skill flags them for user resolution.

## ADRs in This Cluster

| ADR  | Title                                 | Status   |
|------|---------------------------------------|----------|
| [[../adr/0001-capture-execute-boundary\|0001]]  | Capture ≠ Execute boundary                 | accepted |
| [[../adr/0010-reqs-as-validated-intent\|0010]]  | REQs as validated intent                   | accepted |

## How They Relate

- [[../adr/0001-capture-execute-boundary|0001]] is the **behavioral** commitment — what the agent *does*.
- [[../adr/0010-reqs-as-validated-intent|0010]] is the **epistemic** commitment — what the artifacts *mean*.

They support each other: the hard boundary from 0001 is what creates the validation window that 0010 then trusts.

## Why These Are Not Just "Niceties"

Every temptation the skill is designed to resist — helpfully starting to build after capture, letting users edit in-flight REQs, inferring requirements instead of asking — violates one of these two principles. When you see a guardrail in an action file that says "STOP," "do not proceed," or "ask the user," it's almost certainly enforcing one of these ADRs.
