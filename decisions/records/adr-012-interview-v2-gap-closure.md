---
title: "ADR-012: Interview v2 — Gap Closure Patches"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-04-17
sources:
  - decisions/imported-specs/2026-04-12_close-gaps-in-interview.md
  - actions/interview.md
  - actions/interview-reference.md
  - interviews/work-operating-model.md
  - crew-members/interviewer.md
related:
  - page: adr-011-interview-framework-with-prescriptive-templates
    rel: extends
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: depends-on
  - page: adr-002-load-reusable-spec-templates-during-work
    rel: depends-on
created: 2026-04-17
updated: 2026-04-17
confidence: high
---

# ADR-012: Interview v2 — Gap Closure Patches

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-011-interview-framework-with-prescriptive-templates]] (extends), [[adr-001-modular-action-prompts-and-companion-references]] (depends-on), [[adr-002-load-reusable-spec-templates-during-work]] (depends-on)

## Context

V1 of the `interview` action (ADR-011) shipped the framework, session schema, first template, crew persona, docs, and imported spec. Five gaps remained: (1) export schemas were described in prose, not as mechanical render templates — so different runs could produce different file shapes, breaking the "feed `USER.md` into any agent platform" guarantee; (2) the placement of `crew-members/interviewer.md` assumed a convention that required empirical validation; (3) `update` re-run mode did not specify whether revalidation was per-layer or per-entry, which affected CHANGELOG granularity and the semantics of `last_validated_at`; (4) mid-layer recovery was undefined — if a user quit partway through a layer, resume behavior was implicit; (5) the `ingest` sub-command's file mapping was underspecified (one file per export with no layer summaries, and its frontmatter diverged from BKB's canonical schema).

The v2 spec at `decisions/imported-specs/2026-04-12_close-gaps-in-interview.md` proposed surgical patches for each gap, explicitly rejecting a rewrite.

## Decision

Close the five gaps as surgical patches. V1's core architecture stays intact.

1. **Export templates move into the template file.** `interviews/work-operating-model.md` gains an `## Export Templates` section containing the five verbatim handlebars-style render templates (`USER.md`, `SOUL.md`, `HEARTBEAT.md`, `operating-model.json`, `schedule-recommendations.json`). `actions/interview-reference.md` trims its per-template prose to framework-level invariants only (narrative tone, source-confidence filtering, cadence requirements, traceability to layer entries). This reverses v1's "do not duplicate in template body" note, which is explicitly overridden by this ADR.

2. **Crew placement: leave `interviewer.md` in `crew-members/`.** Empirical audit of the directory (12 files: `general`, `karpathy`, `approach-directives`, domain personas, `debugging`, `caveman`, `interviewer`) showed it is a generic persona pool used across `work`, `pipeline`, and `interview`, not scoped to one action. No move required.

3. **`update` re-run mode is entry-level.** Each entry in a stored layer receives a `[confirm / edit / mark-stale / delete / skip]` prompt; new entries can be added after the walk; the layer-level approval gate still applies. The v1 note "Per-entry edit friction is intentional … Do not invent a per-entry patch path" is explicitly overridden. CHANGELOG entries for update runs summarize `N confirmed, N edited, N marked stale, N deleted, N added`.

4. **Mid-layer recovery uses draft checkpoints.** The interview writes a draft checkpoint at `./do-work/interview/<template>/checkpoints/.draft-<layer-id>.md` after producing candidate entries but before user approval. On resume, the action checks for the draft, offers pick-up vs. start-over, and deletes the draft after approval or discard.

5. **`ingest` writes 10 files per run for `work-operating-model`.** 5 export files (`<template>-<export-name>.md`) + 5 layer summaries (`<template>-<layer-id>.md`), plus a manifest row per file in `kb/raw/_inbox_queue.md`. Frontmatter aligns with BKB's schema: `sources:` as a list, `related:` with `rel`, valid `type` values (`source-summary` for exports, `concept` for summaries). The HHMMSS collision prefix is honored.

A plan-and-approve checkpoint is now the expected workflow norm for future prompt-level changes to this skill — future imported specs should require a written plan before edits begin.

## Alternatives

1. Full rewrite of v1 to embed every gap fix.
Rejected as wasteful — v1's core architecture is sound, and the gaps are localized to specific sections.

2. Per-layer granularity for update mode (v1's shipped behavior).
Rejected because it is coarser than users need — fixing one typo forces re-interviewing the whole layer. Entry-level revalidation preserves long-lived entries that remain accurate.

3. Reconstruct mid-layer state from chat history on resume.
Rejected as unreliable — draft checkpoint files are explicit, inspectable, and idempotent.

4. Move `crew-members/interviewer.md` inline into `interview-reference.md`.
Rejected because the audit showed `crew-members/` is a generic persona pool, not a `work`-scoped directory. Consistency with the existing layout wins.

5. Keep export schemas in the reference and duplicate handlebars templates there.
Rejected because it creates two sources of truth. The template file owns mechanical rendering; the reference owns framework-level invariants.

## Consequences

Exports are now reproducible across runs and across implementations — an agent consuming the templates produces the same file shape the next agent will. The `ingest` sub-command produces a predictable shape in `kb/raw/inbox/`, making the cross-action seam between `interview` and `bkb` reliable. `update` mode supports partial revalidation, preserving entries that remain accurate while letting others be marked stale or deleted individually. Mid-layer quits are recoverable without re-answering. Crew placement reflects the repo's actual pattern rather than a misread analogy. Future prompt-level patches to this skill are expected to follow the plan-and-approve workflow demonstrated here.

The cost is a small ongoing maintenance tax on the template file: when new templates are authored, each one must define its own `## Export Templates` section using the same handlebars syntax. The upside is that the reference stays slim and new templates inherit a clear authoring contract.

## References

- [decisions/imported-specs/2026-04-12_close-gaps-in-interview.md](../imported-specs/2026-04-12_close-gaps-in-interview.md) — v2 spec this ADR implements
- [decisions/imported-specs/expand-skill-do-work-interview.md](../imported-specs/expand-skill-do-work-interview.md) — v1 spec
- [actions/interview.md](../../actions/interview.md) — action (draft checkpoint step, ingest sub-command update)
- [actions/interview-reference.md](../../actions/interview-reference.md) — reference (Export Schemas trim, update mode rewrite, Mid-layer recovery, Ingest File Mapping)
- [interviews/work-operating-model.md](../../interviews/work-operating-model.md) — template (Export Templates section added)
- [crew-members/interviewer.md](../../crew-members/interviewer.md) — persona (unchanged; audit confirmed placement)
