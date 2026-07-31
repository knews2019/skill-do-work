---
id: REQ-038
title: Handle worktree/branch name collision when re-dispatching a crash-recovered REQ
status: completed
claimed_at: 2026-07-29T08:48:16Z
completed_at: 2026-07-29T08:51:00Z
commit: efb6300
route: A
created_at: 2026-07-28T22:41:30Z
user_request: UR-007
addendum_to: REQ-033
depends_on: []
related: [REQ-033, REQ-037]
batch: parallel-dispatch
domain: general
prime_files: []
tdd: false
review_generated: true
write_set:
  - actions/work-reference.md
maintenance: false
---

# Handle worktree/branch name collision when re-dispatching a crash-recovered REQ

## What

The `worktree-agent-REQ-NNN-<suffix>` name is deterministic with no uniqueness component. After a crash, the sweep correctly reports (never deletes) an unmerged leftover — and the recovered REQ then re-dispatches into exactly that occupied name: `git worktree add` / branch creation fails, and nothing in the text says what to do.

## Why (if provided)

Review of REQ-033 (confirmed by 2 independent adversarial verifiers): the report-only sweep and the deterministic naming rule compose into a livelock — the leftover blocks the name, consent-gated deletion may not run for a while, and the re-dispatched REQ cannot start.

## Detailed Requirements

- In the Worktree Dispatch Mode section's Naming paragraph (`actions/work-reference.md`): on a name collision at creation, do not delete or force — dispatch under a fresh unique variant (e.g. an incrementing `-2`, `-3` or short timestamp token appended to the suffix), report the coexistence, and leave the original leftover to its owners (crash sweep if merged, cleanup Pass 5 if unmerged). The REQ-id prefix must stay intact so sweeps still correlate both names to the REQ.
- One sentence in the crash sweep: a reported unmerged leftover does not block re-dispatch — the naming rule's collision variant covers it.
- Keep the naming ratchet green (`worktree-agent-REQ-` correlation text must survive).

## Constraints

- `actions/work-reference.md` only; a couple of sentences, not a new subsection.

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the exact location (the Worktree Dispatch Mode **Naming** paragraph, `actions/work-reference.md:229`) and the exact change (a collision-variant naming rule + one sentence in the crash sweep). Single file, a couple of sentences, no new subsection — well-specified with obvious scope. No planning or exploration needed.

**Planning:** Not required (Route A)

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified)

**What was done:** Added a name-collision rule to the Worktree Dispatch Mode **Naming** paragraph — on a collision at creation (`git worktree add`/branch creation fails because the deterministic `worktree-agent-REQ-NNN-<suffix>` is held by a reported-but-not-deleted leftover), dispatch under a fresh unique variant (incrementing `-2`/`-3` or a short timestamp token appended to the suffix), keeping the `worktree-agent-REQ-NNN-` prefix intact so sweeps still correlate both names, report the coexistence, and leave the leftover to its owners. Added one sentence to the crash sweep's "Unmerged or dirty" bullet: a reported unmerged leftover does not block re-dispatch — the collision variant covers it, so the two coexist until Pass 5 resolves the leftover. This breaks the report-only-sweep + deterministic-name livelock REQ-033's review flagged.

*Summary written by work action (orchestrator)*

## Review

**Overall: 95%** (Route A — single-reviewer quick scan per the session's calibration) | 2026-07-29T08:50:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor. Both requirements delivered — the Naming collision variant (fresh `-2`/`-3`/timestamp token, prefix kept) and the crash-sweep "does not block re-dispatch" sentence; cross-references between the two resolve; `worktree-agent-REQ-` and `git branch -d` ratchets green; single file, no new subsection (honors the Constraint). Instruction prose gated to worktree mode; serial path untouched.
**Acceptance:** Pass — `contract-regressions.sh` green; change verified by read.
**Follow-ups created:** None.

*Reviewed by work action (orchestrator, single-reviewer pass)*

## Orientation

Worktree dispatch mode's deterministic `worktree-agent-REQ-NNN-<suffix>` naming now has a collision rule: a crash-recovered REQ re-dispatching into the name its own reported-but-undeleted leftover still holds gets a fresh unique variant instead of a failed `git worktree add`, breaking the report-only-sweep + deterministic-name livelock. One-paragraph addition to `actions/work-reference.md` → Worktree Dispatch Mode (Naming + crash sweep). Leaf change, no map impact.
