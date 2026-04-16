---
title: "ADR-011: Interview Framework with Prescriptive Templates"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-04-16
sources:
  - CHANGELOG.md (0.67.0 The Open Ear)
  - actions/interview.md
  - actions/interview-reference.md
  - interviews/work-operating-model.md
  - crew-members/interviewer.md
  - docs/interview-guide.md
related:
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: depends-on
  - page: adr-002-load-reusable-spec-templates-during-work
    rel: depends-on
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: depends-on
  - page: adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb
    rel: complements
created: 2026-04-16
updated: 2026-04-16
confidence: high
---

# ADR-011: Interview Framework with Prescriptive Templates

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-001-modular-action-prompts-and-companion-references]] (depends-on), [[adr-002-load-reusable-spec-templates-during-work]] (depends-on), [[adr-005-pipeline-is-stateful-and-resumable]] (depends-on), [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]] (complements)

## Context

The do-work skill includes actions for knowledge compilation (`bkb`), code review, commits, and pipeline orchestration, but has no action for extracting the operator's own tacit knowledge into a form agents can act on. Without this, users who want agents to take work off their plate must write their own `USER.md` / `SOUL.md` / `HEARTBEAT.md` by hand, which is the step most people fail at. The Work Operating Model by Nate B. Jones and Jonathan Edwards defines a five-layer elicitation interview that produces these files; it was published with an Open Brain / MCP dependency the skill does not want to take.

## Decision

Add an `interview` action implementing a generalized elicitation framework with prescriptive templates. Templates live in `interviews/<name>.md` and declare layers, per-layer prompts, canonical entry contract, and export schemas. The action enforces contracts, handles session state, runs checkpoint approval gates, and produces exports. The first template is `work-operating-model`, using the five layers and canonical entry contract from the source prescription verbatim. Persistence is local-file only under `./interview/<template>/`. Templates are prescriptive (not minimal, not fully executable) to balance consistency of output against authoring cost for future templates.

## Alternatives

1. Minimal template shape.
Rejected because a template that only names its layers would let every instance diverge into its own question flow, entry shape, and export format. The whole point of the framework is producing `USER.md` / `SOUL.md` / `HEARTBEAT.md` that downstream agents can consume predictably.

2. Executable template shape (each template owns its own full workflow).
Rejected because it pushes too much per-template authoring work and risks every template diverging into its own action. The prescriptive shape lets new templates (post-mortem, new-hire-onboarding, project-kickoff) be authored against a fixed schema with no action changes.

3. Profiles within the same repo.
Rejected in favor of multi-repo installation. Single instance per template per repo keeps paths simple (`./interview/<template>/`) and avoids a profile-selection UX layer the skill does not otherwise need.

4. Direct Open Brain / MCP integration.
Rejected because it adds infrastructure the skill does not otherwise require, locks users into a specific backend, and violates the skill's local-files-only constraint.

## Consequences

Users can produce agent-ready operating artifacts by running a ~45-minute interview. Exports flow into `bkb` via the `ingest` sub-command, making the operating model queryable alongside other knowledge. Adding new templates requires writing a template file against the prescriptive schema; no changes to the action. Session state is per-CWD, which means multi-context users install the skill in multiple repos — no profile concept is introduced. Re-runs support `fresh`, `update`, and `version` modes; versions are immutable. No external service dependencies.

The cost is a new surface area to maintain: the canonical entry contract, session schema, export schemas, and crew persona all need to stay in sync with the action behavior. The prescriptive template shape also means future template authors must conform to the declared layer/export structure — freedom is bounded.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.67.0 The Open Ear`
- [actions/interview.md](../../actions/interview.md)
- [actions/interview-reference.md](../../actions/interview-reference.md)
- [interviews/work-operating-model.md](../../interviews/work-operating-model.md)
- [crew-members/interviewer.md](../../crew-members/interviewer.md)
- [docs/interview-guide.md](../../docs/interview-guide.md)
