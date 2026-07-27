---
title: "ADR-014: Considered & Declined — Autonomous Loop-Until-Done / Ultracode-Fable Workflow"
type: architecture-decision-record
status: declined
topic_cluster: workflow-orchestration
decided: 2026-06-29
sources:
  - CHANGELOG.md (0.98.0 The Delete Key)
  - commit ecbe59f (removed prompts/ultracode-fable-workflow.md, added the contract test)
  - _dev/tests/contract-regressions.sh
  - crew-members/maintenance.md
  - crew-members/background-agents.md
  - crew-members/karpathy.md
  - actions/work.md
  - SKILL.md
related:
  - page: adr-003-always-load-karpathy-guardrails
    rel: complements
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: complements
  - page: adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles
    rel: complements
  - page: adr-013-harden-the-vendored-skill-distribution-model
    rel: complements
created: 2026-06-29
updated: 2026-06-29
confidence: high
---

# ADR-014: Considered & Declined — Autonomous Loop-Until-Done / Ultracode-Fable Workflow

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-003-always-load-karpathy-guardrails]] (complements), [[adr-005-pipeline-is-stateful-and-resumable]] (complements), [[adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles]] (complements), [[adr-013-harden-the-vendored-skill-distribution-model]] (complements)

## Context

A recurring proposal is to add an "autonomous loop-until-done" runner to the skill: the `/loop do-work run` pattern (keep churning the queue unattended) bundled with four standing reminders — use background agents/workflows to manage context, commit via a background `do-work commit` agent, follow YAGNI, and an "ultracode mode" model-economy (a Fable/Opus/Sonnet/Haiku tier table that assigns audit/orchestrate/build roles to model tiers).

This capability already existed and was deliberately removed. `prompts/ultracode-fable-workflow.md` (206 lines) was deleted in commit `ecbe59f` (2026-06-15). The **same commit** added `_dev/tests/contract-regressions.sh`, a guard that fails if that file returns *or* if the strings `ultracode`/`fable` reappear in any of six active runtime docs (SKILL.md, README.md, next-steps.md, actions/work.md, actions/work-reference.md, prompts/README.md). The deletion was reaffirmed during the 0.98.0 "The Delete Key" release, whose crew rule `crew-members/maintenance.md` ("The Subtractor") codifies delete-before-you-add: a new rule is the most expensive fix, and any addition must clear a replay bar (a concrete case that fails without it and passes with it).

A subsequent investigation re-opened the question — *should the workflow be re-added?* — which is itself the recurrence this record exists to stop. An independent re-verification of the gap analysis (read-only) confirmed the prior conclusion.

## Decision

**Declined. No functional change.** The autonomous loop-until-done / ultracode-fable workflow will not be re-added, because every model-agnostic capability it bundled already survives as canon, and the only piece that did not survive is model-specific and intentionally out of scope for a model-agnostic skill:

- background agents / context durability (disk-as-source-of-truth, bounded waves, the "Native orchestration engine" rung, reasoning-block-corruption recovery) → `crew-members/background-agents.md`
- YAGNI → `crew-members/karpathy.md` §Simplicity First (always loaded during implementation; see [[adr-003-always-load-karpathy-guardrails]])
- the queue loop with between-REQ context wipe → `actions/work.md` Step 10
- background dispatch of `work` and `commit` → `SKILL.md` Action Dispatch
- explicit staging (never `git add -A`; validate staged files against the Implementation Summary) → `actions/work.md`, `actions/work-reference.md`, `actions/commit.md`
- fixture faithfulness / caller-seam testing → `crew-members/testing.md`; fresh-context review → `actions/review-work.md`

The lost piece — the Fable/Opus/Sonnet/Haiku tier table — is a model-economy policy tuned to one model family. It does not belong in a skill that must run on any agentic coding tool. Parallel-orchestrator guards (a full-suite gate, pre-commit collision guards) were likewise deliberately scoped out (`actions/work.md` records the single-orchestrator scope decision). No genuinely missing, model-agnostic rail clears `maintenance.md`'s replay bar, so no addition is warranted.

This ADR is documentation, not a runtime instruction: the contract test *enforces* the deletion mechanically but cannot *explain* it. This record supplies the rationale so the decision is not re-litigated; its replay case is the recurring re-investigation it short-circuits.

## Alternatives

1. **Re-add the full `prompts/ultracode-fable-workflow.md`.** Rejected: it would fight the contract test and the 0.98.0 delete-before-you-add posture, and it re-introduces model-specific policy (`ultracode`/`fable`) into a model-agnostic skill — duplicating capabilities that already live in the crew/action files above.

2. **Re-add only the model-agnostic parts as a new prompt or action.** Rejected: those parts are already canon. A new file would be a second, drifting copy of `background-agents.md` + `karpathy.md` + `work.md` Step 10 — the exact bloat `maintenance.md` warns against, with no replay case showing a gap.

3. **Add the model tier table alone (keep it model-specific but documented).** Rejected: it ties the skill to one model family and violates the agent-compatibility constraint (action files must work with any agentic coding tool). Model-economy tiering is a consumer-side policy, not skill canon.

4. **Record nothing; rely on the contract test alone.** Rejected: the test blocks the file's return but offers no rationale, so the investigation recurs (as it just did). A short declined-decision record is the cheaper long-run fix.

## Consequences

The skill stays lean and model-agnostic; the contract test continues to enforce the deletion, and this ADR carries the "why" a failing test cannot. Future proposals to re-add the workflow have a single citable answer, ending the re-investigation loop.

The trade-off is a small amount of new documentation surface (one ADR plus its index/log/topic updates) recording a *non-change*. This is accepted because the cost is one-time and the recurring re-investigation it prevents is not. If a future harness surfaces a concrete, model-agnostic capability gap with a replay case, this decision should be revisited via a new ADR — declining today is not declining forever.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.98.0 The Delete Key`
- commit `ecbe59f` — removed `prompts/ultracode-fable-workflow.md`, added the contract test
- [_dev/tests/contract-regressions.sh](../../_dev/tests/contract-regressions.sh)
- [crew-members/maintenance.md](../../crew-members/maintenance.md)
- [crew-members/background-agents.md](../../crew-members/background-agents.md)
- [crew-members/karpathy.md](../../crew-members/karpathy.md)
- [actions/work.md](../../actions/work.md)
- [SKILL.md](../../SKILL.md)
