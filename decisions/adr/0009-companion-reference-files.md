---
id: 0009
title: Companion Reference File Pattern
status: accepted
decided: 2026-04-11
version: 0.61.1
topic: content-structure
supersedes: []
superseded_by: null
related:
  - adr: 0007
    rel: complements
---

# ADR-0009: Companion Reference File Pattern

## Context

Several action files grew past the point where an agent could read them in a single pass. Modern agentic tools enforce a per-read token ceiling (around 10k tokens for a single file read). When an action file exceeds that limit, the reader either has to make multiple reads or work with a truncated view — both of which introduce subtle bugs ("I wrote the first half correctly, but I didn't see the rules in the second half").

`work.md` was the first to hit this wall, followed by `build-knowledge-base.md`, `deep-explore.md`, and `pipeline.md`. Each one contained a mix of *flow* (what the action does, in what order) and *reference material* (templates, schemas, detailed error handling, persona prompts) — the reference material was what pushed each file over the limit.

The tempting reaction — "trim it" — would have cost real signal. The templates and schemas were there for specific reasons; cutting them would just trade size bugs for correctness bugs.

## Decision

When an action file exceeds the single-read token limit, split it into two:

1. **Primary file** (`{action}.md`) — holds the action's flow: description, When to Use, Input, Steps, Output, Rules, Common Rationalizations, Red Flags, Verification Checklist.
2. **Companion reference** (`{action}-reference.md`) — holds the reference material: templates, schema definitions, persona prompts, long error-handling tables, detailed state schemas.

The two are cross-linked: the primary file tells the agent exactly when to consult the companion ("see the 'Schema File Content' section of the bkb-reference action"), and the companion opens with a pointer back to the primary. The companion is loaded **on demand** from the specific step that needs it — never always-loaded.

Current primary/companion pairs:

- `work.md` + `work-reference.md`
- `build-knowledge-base.md` + `bkb-reference.md`
- `deep-explore.md` + `deep-explore-reference.md`
- `pipeline.md` + `pipeline-reference.md`

## Alternatives Considered

- **Keep the files monolithic and accept truncation.** Rejected — the truncation risk is exactly the bug this decision is designed to prevent.
- **Summarize aggressively.** Cut the reference material. Rejected — the templates and schemas contained content-dependent detail (persona voice examples, slide-number sequences) that couldn't be usefully summarized without losing signal.
- **Many small files.** Split an action into five or six topical files with a coordinating index. Rejected — too much fragmentation for the reader; the primary/companion split is the smallest partition that solves the problem.
- **Generate companions automatically.** Extract reference sections via tooling. Rejected — the skill is platform-agnostic ([[0004-platform-agnostic-action-files]]); tooling would be tool-specific.

## Consequences

- **Primary files stay readable in one pass.** Agents can load the flow without loading the templates they don't need yet.
- **Companion files load only when needed.** A work action that doesn't hit its error-handling table never loads the companion.
- **Two files to audit for consistency.** Changes that touch the action's flow and its templates have to update both. Reviewers watch for drift.
- **Cross-references are prose, not paths.** Per the CLAUDE.md convention, one file refers to another by short name ("the bkb-reference action"), not by file path. SKILL.md owns the file-path mappings.
- **The pattern scales.** New actions that start small can stay monolithic; the split only happens when the single-read limit is crossed. A recent example is v0.64.1, when `pipeline.md` crossed the limit and was split the same way.

## References

- **CHANGELOG**: v0.61.1 — The Lean Cut (2026-04-11) split `build-knowledge-base.md`; v0.64.1 — The Companion Split (2026-04-13) split `pipeline.md`
- **Documents**: `CLAUDE.md` ("Cross-reference other actions by short name")
- **Action files**: `actions/work.md` + `actions/work-reference.md`; `actions/build-knowledge-base.md` + `actions/bkb-reference.md`; `actions/deep-explore.md` + `actions/deep-explore-reference.md`; `actions/pipeline.md` + `actions/pipeline-reference.md`
- **Related ADRs**: [[0007-crew-member-jit-loading]] (both express the same "load only what's relevant" instinct at different granularities)
