---
id: 0004
title: Platform-Agnostic Action Files
status: accepted
decided: 2026-02-03
version: 0.8.0
topic: platform-portability
supersedes: []
superseded_by: null
related:
  - adr: 0005
    rel: complements
---

# ADR-0004: Platform-Agnostic Action Files

## Context

The skill was initially written with Claude Code in mind. Action files referenced tool-specific APIs directly — the Task tool for subagent dispatch, `run_in_background` for fire-and-forget execution, `Co-Authored-By: Claude <noreply@anthropic.com>` trailers baked into commit templates.

As soon as anyone tried running the skill on another agentic coding tool, it fell apart. Generic agents would stamp Claude-specific commit metadata onto their commits just by following the template verbatim. Agents without a Task tool would try to invoke one and fail. The skill that claimed to be about task queues turned out to be a task queue **for Claude Code** — a lock-in nobody had intentionally chosen.

The underlying question was whether the skill's audience is "one specific tool's users" or "anyone with a shell and a Markdown editor." Every decision that embedded tool-specific assumptions made the second audience impossible.

## Decision

Action files are written with **three portability rules**:

1. **Generalized language in prose.** No tool-specific API names, no proprietary metadata, no assumptions about background execution or subagent availability. Where an action needs to talk about "dispatching," it uses phrases like "spawn a subagent" or "use your environment's ask-user prompt."
2. **Each action file stands alone.** It must work as a prompt pasted into a bare chat interface. If reading the file requires knowing about a specific runtime feature, that's a defect.
3. **Design for the floor, not the ceiling.** The simplest agent that can read/write files and run shell commands must be able to follow the instructions. Subagents, parallel execution, and background dispatch are nice-to-haves layered on top, not load-bearing assumptions.

These rules are documented as the "Agent Compatibility" section of `CLAUDE.md` so future edits stay honest.

## Alternatives Considered

- **Target Claude Code only, let others fork.** The path of least resistance — ship what works for the most-used tool, and hope other tools' users port it. Rejected — in practice, nobody forks maintenance projects; portability has to be a first-class commitment from day one.
- **Conditional branches for each tool.** "If running in Claude Code, do X; else, do Y." Rejected — N × 2 complexity for every tool added, and the conditionals rot quickly when either tool changes.
- **Build an adapter layer.** A translation shim between the skill and each target tool. Rejected — the only consistent adapter is prose, and the skill's primary deliverable already is prose.

## Consequences

- **The skill works across tools.** Claude Code, other LLM CLIs, and even bare chat interfaces can follow the action files.
- **No tool-specific optimizations baked in.** Can't lean on one tool's advanced features without either (a) making them optional fallbacks or (b) violating the contract. See [[0005-subagent-dispatch-pattern]] for how the subagent-vs-direct-read split handles this.
- **The floor sets the ceiling.** If a new feature can only work with Claude Code's specific APIs, it probably doesn't land in the skill. This is a real constraint — it rules out some useful ideas.
- **Historical tool-specific leaks have to be cleaned up.** v0.11.1 removed Claude Code language from dispatch; v0.12.1 pulled out the hardcoded `Co-Authored-By` trailer. Each leak is a reminder that portability isn't self-enforcing.

## References

- **CHANGELOG**: v0.8.0 — The Bright Light (agent compatibility section added to CLAUDE.md), v0.11.1 — The Soft Landing (removed Claude Code dispatch language), v0.12.1 — The Passport Check (removed Co-Authored-By trailer)
- **Documents**: `CLAUDE.md` ("Agent Compatibility" section)
- **Related ADRs**: [[0005-subagent-dispatch-pattern]] (the how — specifically for subagent availability)
