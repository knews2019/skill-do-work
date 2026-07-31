---
id: REQ-039
title: "Confirm: mention worktree isolation in background-agents.md harness tiers"
status: completed
claimed_at: 2026-07-29T08:52:18Z
completed_at: 2026-07-29T08:53:15Z
commit: 9c8ef45
route: A
status_changed_at: 2026-07-28T22:58:50Z
created_at: 2026-07-28T22:41:30Z
user_request: UR-007
addendum_to: REQ-033
builder_decided: true
domain: general
prime_files: []
review_generated: false
maintenance: false
---

# Confirm: mention worktree isolation in background-agents.md harness tiers

## What

While building REQ-033 (the worktree dispatch mode), the builder noticed `crew-members/background-agents.md` tiers harness capability (tiers 1–3) around manual parallel/background spawns and never mentions worktree isolation — a fourth, orthogonal capability axis (an isolated working directory per builder). Nothing is wrong today: that file's JIT_CONTEXT is condition-stated and the tiers still work. But a reader sizing a harness against those tiers won't learn worktree isolation exists or where it's documented (`actions/work-reference.md` → Worktree Dispatch Mode). This follow-up asks whether you want a brief cross-mention added.

## What the Builder Chose

Nothing — the file was outside REQ-033's declared write_set, so the builder recorded it as a [low] discovered task instead of editing it. This is a "should we also do this" question, not a review of something already built.

## What Would Change

If approved: one or two sentences in `crew-members/background-agents.md` near its harness-tier list, naming worktree isolation as a separate axis and pointing at `actions/work-reference.md` → Worktree Dispatch Mode. If declined: nothing changes; the tiers remain accurate for what they cover.

## Open Questions

- [x] Add a brief worktree-isolation cross-mention to `crew-members/background-agents.md`'s harness-tier area? → Confirmed: Yes, add the mention (one or two sentences plus a pointer to `actions/work-reference.md` → Worktree Dispatch Mode)
  Recommended: Yes — one sentence plus a pointer; cheap, and it closes a discoverability gap for harness sizing.
  Value: A reader evaluating "what can my harness do" finds the isolation mode without already knowing it exists.
  Risk: Near zero and fully reversible — one sentence; the only cost is a line of drift surface if the mode's home ever moves.
  Also: No — keep background-agents.md focused on spawn/durability tiers and let the work pipeline docs own all worktree knowledge.

---

## Triage

**Route: A** - Simple

**Reasoning:** The single Open Question was already resolved by the user (`- [x]` Confirmed: Yes, add the mention) before this run, so there is no ambiguity — implement the confirmed answer: one short paragraph in `crew-members/background-agents.md` naming worktree isolation as a separate axis plus a pointer to `actions/work-reference.md` → Worktree Dispatch Mode. Single file, no new machinery. Implemented inline by the orchestrator (floor-agent path) given the trivial, fully-specified scope.

**Planning:** Not required (Route A)

## Implementation Summary

**Files changed:**
- `crew-members/background-agents.md` (modified)

**What was done:** Added a short "Worktree isolation is a separate axis" paragraph immediately after the three harness rungs in the **Match the Pattern to the Harness** section. It states that the rungs measure how much orchestration the harness hands you, that per-builder git-worktree isolation is an orthogonal capability available at any rung, and points to `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** — closing the discoverability gap the REQ named (a reader sizing a harness against the tiers now learns worktree isolation exists and where it's documented). The file's JIT_CONTEXT and the tier list itself are unchanged.

*Summary written by work action (orchestrator)*

## Review

**Overall: 96%** (Route A — single-reviewer quick scan per the session's calibration) | 2026-07-29T08:53:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor. Delivers exactly the confirmed answer — one paragraph, a pointer to the canonical worktree home (no duplication of worktree knowledge into background-agents.md, honoring the "Also: No" concern's spirit), tiers untouched. `contract-regressions.sh` green.
**Acceptance:** Pass — cross-reference resolves to the real heading in `actions/work-reference.md`.
**Follow-ups created:** None (the Open Question was user-answered; no builder-deferred `- [~]` decision to confirm).

*Reviewed by work action (orchestrator, single-reviewer pass)*

## Orientation

A reader sizing a harness against `crew-members/background-agents.md`'s three fan-out rungs now also learns that per-builder git-worktree isolation is a separate, orthogonal capability and where it lives (`actions/work-reference.md` → Worktree Dispatch Mode). One-paragraph discoverability fix; leaf change, no map impact.
