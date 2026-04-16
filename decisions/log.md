---
title: "Decision Timeline"
type: timeline-log
status: reference
sources:
  - CHANGELOG.md
  - decisions/records/
related:
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: evidence-for
  - page: adr-002-load-reusable-spec-templates-during-work
    rel: evidence-for
  - page: adr-003-always-load-karpathy-guardrails
    rel: evidence-for
  - page: adr-004-canonicalize-pending-reqs-under-do-work-queue
    rel: evidence-for
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: evidence-for
  - page: adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles
    rel: evidence-for
  - page: adr-007-close-the-pipeline-with-present-and-a-technical-debrief
    rel: evidence-for
  - page: adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats
    rel: evidence-for
  - page: adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki
    rel: evidence-for
  - page: adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb
    rel: evidence-for
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# Decision Timeline

Append-only timeline. Historical entries use the original decision dates from `CHANGELOG.md`; the final entry records the ADR-log bootstrap on 2026-04-15.

## [2026-04-05] Knowledge-base foundation

- Accepted [[adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki]] from `0.43.0` through `0.43.5`: BKB is a compiled wiki with a master index, topic indexes, stable processed-source paths, and append-only logs.

## [2026-04-06] Knowledge-base connective tissue

- Accepted [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]] from `0.45.0`, `0.46.0`, and `0.47.0`: typed relationships, retrieval memory, and the BKB crew model are part of the operating design.

## [2026-04-07] Skill modularization pattern

- Accepted [[adr-001-modular-action-prompts-and-companion-references]] from `0.49.0` and later follow-up entries: keep runnable action prompts lean and move bulky stable material into companion references.

## [2026-04-08] Pipeline becomes a first-class orchestrator

- Accepted [[adr-005-pipeline-is-stateful-and-resumable]] from `0.51.5`: the pipeline owns resumable macro-state in `do-work/pipeline.json` and dispatches existing actions.

## [2026-04-10] Queue and workflow coordination harden

- Accepted [[adr-004-canonicalize-pending-reqs-under-do-work-queue]] from `0.60.3`: pending REQs live in `do-work/queue/`.
- Accepted [[adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles]] from `0.56.0` through `0.56.2`: post-pipeline queue continuation happens in bounded run-review loops.
- Accepted [[adr-002-load-reusable-spec-templates-during-work]] from `0.59.0`: `specs/` templates are a reusable scaffold for recurring task types.

## [2026-04-12] Always-on quality guardrails

- Accepted [[adr-003-always-load-karpathy-guardrails]] from `0.62.0` and `0.62.1`: Karpathy-style behavior rules are always loaded during implementation and audited during review.

## [2026-04-13] Pipeline deliverables become durable, multi-surface artifacts

- Accepted [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] from `0.63.0` and `0.63.1`: the pipeline ends with present-work and a persisted debrief.
- Accepted [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] from `0.63.2`, `0.64.0`, and `0.64.1`: one dataset, three renderings, and cross-links for both stakeholder and developer audiences.

## [2026-04-15] ADR log bootstrap

- Created `decisions/` with [[_master_index]], four topic indexes, ten retroactive ADRs, this timeline, and [[_progress]].
- Scope rule for this pass: capture only load-bearing decisions that remain in force as of the current repo state; defer superseded or short-lived experiments to a future expansion pass.
