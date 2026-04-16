---
title: "Decision Log Progress"
type: progress-log
status: reference
sources:
  - CHANGELOG.md
  - decisions/_master_index.md
  - decisions/log.md
related:
  - page: _master_index
    rel: depends-on
  - page: log
    rel: complements
created: 2026-04-15
updated: 2026-04-16
confidence: high
---

# Decision Log Progress

## Status

- State: complete for the bootstrap pass plus ADR-011 (interview framework).
- Last updated: 2026-04-16.
- Next ADR number: `ADR-012`.

## Completed Scope

- Built `decisions/_master_index.md`, `decisions/log.md`, four topic indexes, and ten ADR files under `decisions/records/`.
- Mined only `CHANGELOG.md` entries that looked load-bearing and still in force after cross-checking current files.
- Grouped the ADRs into four topic clusters: skill architecture, workflow orchestration, pipeline deliverables, and knowledge base.

## ADR Inventory

- [x] ADR-001 — [[adr-001-modular-action-prompts-and-companion-references]]
- [x] ADR-002 — [[adr-002-load-reusable-spec-templates-during-work]]
- [x] ADR-003 — [[adr-003-always-load-karpathy-guardrails]]
- [x] ADR-004 — [[adr-004-canonicalize-pending-reqs-under-do-work-queue]]
- [x] ADR-005 — [[adr-005-pipeline-is-stateful-and-resumable]]
- [x] ADR-006 — [[adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles]]
- [x] ADR-007 — [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]]
- [x] ADR-008 — [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]]
- [x] ADR-009 — [[adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki]]
- [x] ADR-010 — [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]]
- [x] ADR-011 — [[adr-011-interview-framework-with-prescriptive-templates]]

## Resume Notes

- If the log needs expansion, start by scanning newer `CHANGELOG.md` entries above `0.64.1` or older foundational entries that were intentionally left out of this first pass.
- Prefer adding only decisions that are both load-bearing and still in force. Superseded or short-lived experiments should either be omitted or clearly marked as historical-only in a future pass.
- When adding a new ADR, update the relevant topic index, `_master_index.md`, and `log.md` in the same change so navigation stays coherent.
- Preserve the BKB-style typed `related:` frontmatter and `[[wiki-links]]` body references for new pages.

## Deferred Candidates

- `Trail of Intent` as a first-class philosophy layer (`0.51.3`) if the log expands beyond the first 10 ADRs.
- Action-file template governance and guardrail sections (`0.61.0`) if the team wants a denser process-governance cluster.
- Testing/security crew expansion (`0.51.0`, `0.54.0`) if crew-member policy gets its own topic cluster.
