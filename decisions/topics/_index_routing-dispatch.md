# Topic: Routing & Dispatch

How a user utterance like `do work capture login form` or `do work review` is turned into an executed action. The routing table in `SKILL.md` is the single source of truth.

## The Routing Invariant

- **First match wins.** The routing table is numbered. The agent walks top-to-bottom, matches the first row whose pattern applies, and dispatches. It does **not** score multiple candidates or ask the user which one to run.
- **Keywords beat content.** Explicit action keywords (`capture`, `work`, `review`, `version`, `cleanup`, …) take priority over descriptive-content interpretation, so `do work version` doesn't get filed as a feature request.
- **Descriptive content is the catch-all.** The very last row handles free-text input by routing to the capture action.

## ADRs in This Cluster

| ADR  | Title                               | Status   |
|------|-------------------------------------|----------|
| [[../adr/0006-priority-ordered-routing\|0006]] | Priority-ordered routing table | accepted |

## Why a Single ADR

The priority-ordering decision is the whole cluster. Everything else about routing (help menu rendering, per-command help, specific keyword mappings) follows from it mechanically — adding an action to the skill is a matter of picking a priority slot and writing the dispatch row.

Future ADRs in this cluster would cover things like: a radically different dispatch model (e.g., tool-use-first instead of keyword-first), or replacing SKILL.md as the routing authority.

## Downstream Behaviors

Every action file assumes it was invoked through the routing table, not called directly. The priority ordering is also what makes `do work code review` (no hyphen, no scope) deterministic — see v0.53.2 for the specific reordering that closed a silent-fallthrough footgun.
