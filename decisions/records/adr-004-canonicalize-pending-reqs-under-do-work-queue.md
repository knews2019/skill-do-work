---
title: "ADR-004: Canonicalize Pending REQs Under do-work/queue"
type: architecture-decision-record
status: accepted
topic_cluster: workflow-orchestration
decided: 2026-04-10
sources:
  - CHANGELOG.md (0.60.3 The Paved Path)
  - CHANGELOG.md (0.60.5 The Honest Mirror)
  - CLAUDE.md
  - actions/capture.md
  - actions/work.md
  - actions/cleanup.md
  - docs/work-guide.md
related:
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: complements
  - page: adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-004: Canonicalize Pending REQs Under do-work/queue

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-005-pipeline-is-stateful-and-resumable]] (complements), [[adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles]] (complements)

## Context

Earlier revisions of the skill stored pending REQs at the `do-work/` root, which created repeated stale-path bugs as contributors instinctively referenced `do-work/queue/` instead. The changelog records a deliberate path migration to pave that cow path, followed by contradiction cleanup to catch remaining references.

The present repository still depends on that decision. `CLAUDE.md` has an explicit Queue Path Convention section, and capture, work, cleanup, forensics, review, versioning, and docs all reference `do-work/queue/` as the canonical pending-work location.

## Decision

Pending `REQ-*.md` files belong in `do-work/queue/`, not in the `do-work/` root. All queue scanning, creation, relocation, cleanup, verification, and review flows are written against that canonical path.

Legacy root-level paths may still be mentioned as a fallback during migration or verification, but they are compatibility shims, not the intended steady-state layout.

## Alternatives

1. Keep pending REQs at the `do-work/` root.
This was rejected because the team kept naturally using `do-work/queue/`, and the mismatch caused recurring documentation and implementation drift.

2. Support multiple canonical pending locations forever.
This was rejected because ambiguity in a queue path weakens tooling, cleanup, and user guidance.

3. Hide the queue path behind undocumented heuristics.
This was rejected because the repo wants explicit, inspectable filesystem conventions.

## Consequences

The benefit is a single mental model for pending work. Queue tooling, cleanup logic, and help text can all assume one stable path.

The cost is migration baggage. Some actions still mention the old root path as a legacy fallback, and reviewers need to watch for stale references when new docs or prompts land.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.60.3 The Paved Path`, `0.60.5 The Honest Mirror`
- [CLAUDE.md](../../CLAUDE.md) — Queue Path Convention
- [actions/capture.md](../../actions/capture.md)
- [actions/work.md](../../actions/work.md)
- [actions/cleanup.md](../../actions/cleanup.md)
- [docs/work-guide.md](../../docs/work-guide.md)
