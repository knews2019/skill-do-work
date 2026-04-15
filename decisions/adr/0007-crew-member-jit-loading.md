---
id: 0007
title: Crew Members — JIT-Loaded Domain Rules
status: accepted
decided: 2026-04-07
version: 0.50.0
topic: content-structure
supersedes: []
superseded_by: null
related:
  - adr: 0009
    rel: complements
---

# ADR-0007: Crew Members — JIT-Loaded Domain Rules

## Context

Domain-specific guidance — "how a backend engineer would think about this," "how a frontend engineer would check accessibility," "how a security reviewer would spot a misuse" — used to live inline inside `work.md`. Every REQ processed through the work action got the full dose of every domain's rules loaded into the subagent's context, whether the REQ was a UI tweak, a database migration, or a README edit.

The cost was both bloat and irrelevance. Work.md grew past 1,200 lines, most of which didn't apply to any given REQ. Subagents spent a portion of their context on rules for domains they weren't touching, and the relevant rules were diluted by the irrelevant ones.

An early fix (`agent-rules/` with always-loaded files) moved the rules out of work.md but didn't address the relevance problem — every file still loaded for every REQ.

## Decision

Domain rules live in individual files under `crew-members/`. They load based on what the REQ actually is:

- **Always loaded**: `general.md` (cross-domain baseline), `karpathy.md` (behavioral guardrails) — loaded during Step 6 of the work action regardless of REQ type.
- **Loaded by `domain:` frontmatter**: `backend.md`, `frontend.md`, `performance.md`, `security.md`, `ui-design.md` — loaded when the REQ's `domain:` frontmatter matches and the file exists.
- **Loaded by other triggers**:
  - `testing.md` on `tdd: true` or `domain: testing`, and alongside `debugging.md` after 2+ test failures.
  - `caveman.md` when `caveman:` frontmatter is set (truthy or intensity level).
  - `debugging.md` during remediation (review fail → retry) or after 2+ test failures.
  - `approach-directives.md` when dispatching multiple sub-agents for parallel/sequential work on related REQs.
- **Graceful absence**: If a rules file is missing, the work action proceeds without it — never blocks on a missing rules file.

Each crew-member file documents its own loading conditions in a top-of-file `JIT_CONTEXT` comment.

## Alternatives Considered

- **Monolithic rules file.** Everything in one always-loaded file. Rejected — the bloat problem that motivated the change.
- **Rules per action.** Each action has its own domain rules inlined. Rejected — duplication across actions, and the rules are about *domains* (backend, security), not about *actions* (review, build).
- **No domain rules, only general guidance.** Trust the agent's training to know domain specifics. Rejected — the rules exist to close specific recurring gaps (security, accessibility, performance) where general guidance consistently misses.

## Consequences

- **Context scales with complexity.** A trivial REQ loads `general.md` + `karpathy.md` and nothing else. A security-sensitive REQ loads `security.md` on top of that. The context budget tracks the work.
- **Adding a domain is cheap.** New crew-member file + one-line addition to the work action's loading rules + a `JIT_CONTEXT` comment at the top of the file. No changes to other crew files.
- **Capture has to tag REQs correctly.** The `domain:` frontmatter gets set during capture. If it's wrong, the wrong rules load — so capture quality matters. The verify-requests action checks for this.
- **Crew members have titles.** v0.50.1 named each role (The Compass for general, The Engineer for backend, etc.) to make it easier to talk about who's handling what.
- **Always-loaded files are privileged.** Adding a new always-loaded crew member is a meaningful change — `karpathy.md` was promoted to always-loaded in v0.62.0 after specific evidence that its principles were broadly useful. Promotions should be deliberate.

## References

- **CHANGELOG**: v0.50.0 — The Crew (2026-04-07) rename from `agent-rules/`; v0.50.1 — The Roll Call (titles); v0.54.0 — The Test Bench (testing.md); v0.62.0 — The Karpathy Nod (always-loaded promotion); v0.62.5 — The Few Words (caveman.md)
- **Action files**: `actions/work.md` (Step 6 loading rules), `crew-members/*.md` (each has `JIT_CONTEXT`)
- **Related ADRs**: [[0009-companion-reference-files]] (the same "load only when relevant" instinct drives companion files)
