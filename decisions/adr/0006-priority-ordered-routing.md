---
id: 0006
title: Priority-Ordered Routing Table
status: accepted
decided: 2026-02-04
version: 0.9.1
topic: routing-dispatch
supersedes: []
superseded_by: null
related: []
---

# ADR-0006: Priority-Ordered Routing Table

## Context

SKILL.md has a routing table that maps user utterances to actions. Early versions of the table treated every row as equally valid — the agent was expected to inspect the input and pick the best match.

This produced silent misroutes. Someone typing `do work version` got their literal words treated as a task description and captured as a feature request ("add versioning"). `do work changelog` had the same problem. The catch-all "descriptive content → capture" pattern was happily swallowing keywords that should have routed to the version action.

The deeper issue was that the routing logic was ambiguous by design: with no priority ordering, two rows could both match, and which one won depended on the agent's interpretation. Subtle phrasing differences would route differently across sessions, making the skill feel unreliable.

## Decision

`SKILL.md`'s routing table is **numbered**. Routing follows three rules:

1. **First match wins.** The agent walks the table top-to-bottom and dispatches to the first row whose pattern matches.
2. **Keyword rows come before descriptive-content rows.** Specific action keywords (`version`, `capture`, `review`, `cleanup`, …) are at the top; the "descriptive content → capture" catch-all is at the very bottom.
3. **A single keyword matching the table is a routed action, not content.** If the user types `do work version`, that's the version action, not a feature request about versioning.

The table is flat (no nesting), numbered (priorities are explicit), and documented as "first match wins" in the routing intro so even a fresh agent reading it sees the rule.

## Alternatives Considered

- **Keyword blocklist.** Maintain a list of reserved words that are never interpreted as content. Rejected — brittle; every new action adds two entries (the action + the blocklist) instead of one.
- **Ask the user to disambiguate.** When two rows could match, prompt. Rejected — interrupts the default flow; users hit this on almost every invocation.
- **Fuzzy matching with confidence scores.** Let the agent score candidates and pick the highest. Rejected — non-deterministic routing is exactly what caused the original problem. Two sessions, same input, different routes.

## Consequences

- **Routing is deterministic.** Same input, same route, every time. Easy to test, easy to explain.
- **Priority slots are a finite resource.** Adding an action means deciding where in the numbered table it goes. Most new actions slot near the end (before the catch-all), but occasionally a new action needs to preempt an older one — see v0.53.2, where bare "code review" was moved from priority 9 (review-work) to priority 7 (code-review) to close a silent fallthrough.
- **The catch-all is the last row on purpose.** Descriptive content is only interpreted as a capture when nothing else matches. This is what makes it safe for users to type free-form requests without worrying about keyword collisions.
- **The table is auditable.** A reviewer can scan the numbered rows and confirm that priorities make sense without running the skill.

## References

- **CHANGELOG**: v0.9.1 — The Gatekeeper (2026-02-04), v0.53.2 — The Short Circuit (reordered "code review" to close fallthrough)
- **Action files**: `SKILL.md` (Routing Table)
- **Related ADRs**: none direct — routing-dispatch is a single-ADR cluster for now.
