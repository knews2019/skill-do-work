# Topic: Platform Portability

The do-work skill targets **any agentic coding tool** that can read/write files and run shell commands — not Claude Code specifically. That ambition constrains how action files are written and how actions are dispatched.

## The Portability Contract

Two invariants define "platform-agnostic":

1. **Every action file is a standalone prompt.** Paste it into a bare chat interface and it should work — no hidden tool-specific APIs, no required subagent machinery, no Claude-Code-only metadata.
2. **Dispatch has a floor.** If subagents and background execution are available, use them. If they aren't, the action still runs — just inline, in the current context.

## ADRs in This Cluster

| ADR  | Title                                     | Status   |
|------|-------------------------------------------|----------|
| [[../adr/0004-platform-agnostic-action-files\|0004]] | Platform-agnostic action files                   | accepted |
| [[../adr/0005-subagent-dispatch-pattern\|0005]]      | Subagent dispatch with direct-read fallback      | accepted |

## How They Relate

- [[../adr/0004-platform-agnostic-action-files|0004]] states the *content* rule — what action files may and may not reference.
- [[../adr/0005-subagent-dispatch-pattern|0005]] states the *execution* rule — how those files get run when tools vary.

Both are enforced by language in `CLAUDE.md` ("Agent Compatibility") and `SKILL.md` ("Action Dispatch"). Breaking either invariant is how the skill stops working on tool X.

## Recurring Temptations

These are the shortcuts this cluster is designed to block:

- Baking `Co-Authored-By: Claude <...>` into commit templates (reverted in v0.12.1).
- Referencing the Task tool, `run_in_background`, or `mcp__*` tools directly in action prose (reverted in v0.11.1).
- Assuming subagents exist without a fallback path (reverted in v0.11.1).

Every time one of these slipped in, a follow-up release had to pull it back out. The cluster is a reminder of why.
