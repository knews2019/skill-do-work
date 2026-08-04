---
id: REQ-096
title: "Execution-model re-grain: claim anywhere, one releaser; dispatch widened to any tree"
status: pending
created_at: 2026-08-04T19:44:17Z
user_request: UR-018
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-094]
maintenance: false
related: [REQ-094, REQ-097, REQ-099, REQ-101]
batch: parallel-building
write_set: [actions/work-reference.md]
---

# Execution-Model Re-Grain: Claim Anywhere, One Releaser

## What

Rewrite the Execution Model contract (`actions/work-reference.md:53–61`) from "one queue owner per checkout, cross-session ownership unsupported" to the user's chosen model: **any checkout may capture and claim/build; exactly one designated releaser checkout runs the release tail** (merge integration, version bump, `CHANGELOG.md` entry, archive moves, UR closure). Widen Worktree Dispatch (`:275–341`) so a builder tree may be a spawned worktree, a user workspace, a clone, or a remote/cloud sandbox.

## Detailed Requirements

**Execution Model rewrite (`:53–61`):**
- Any checkout captures and claims; claims and captures travel between checkouts via ordinary git sync.
- One releaser per queue owns the release tail. Two releasers = unspecified. Two sessions in one working tree = unspecified (unchanged). No prevention machinery for either — repair path is `actions/forensics.md` / `actions/cleanup.md`.
- The `:57` rule survives **verbatim**: never probe for a concurrent session, never ask the user to arbitrate one.
- Cross-checkout conflicts (double claims, duplicate REQ ids from concurrent capture) are ordinary merge artifacts, fixed when the branches meet; `queue-kanban verify` (`duplicate-req-id` probe) is the cheap detector. This philosophy sentence belongs in the contract so downstream prose can cite it.

**Worktree Dispatch widening (`:275–341`):**
- Builder tree generalization: worktree, user workspace, clone, or remote sandbox all satisfy the builder definition (own tree, own branch, hands back a branch).
- Remote hand-back travels on the branch itself — the absolute-main-tree-path handback mechanism is local-only.
- A non-releaser checkout treats its synced `do-work/` snapshot as potentially stale.
- New Red Flag: a second checkout running the **release tail** is the violation to watch — claiming/building/capturing elsewhere is now in contract.
- Keep intact: merge --no-ff hand-back sequence, serial integration, "the merge is the non-interference proof" (`:321`), worktrees-outside-the-repo (`:285`), run-directory pattern.

**Ripple check:** grep shipped files for restatements of the old exclusive-session claim ("one queue owner per checkout", "cross-session ownership") — per the Closed Enumerations Go Stale rule, update every echo (`docs/work-guide.md`'s summary is a known one; REQ-101 owns the full docs pass but inline one-liners that would become *false* get fixed here).

## Constraints

- Serial-only list (`:325`) keeps every item — the release tail stays serial and single-checkout.
- No reservation vocabulary (`reserve` verb, `reserved` status) — the forbidden-token ratchet stays green.
- Do-not-build list per UR-018 Batch Constraints.

## Red-Green Proof

**RED prompt/case:** Today `actions/work-reference.md:55` declares cross-session ownership outside the contract — an agent in a second clone reading it must refuse to claim.
**Why RED now:** The exclusive-session model (0.161.0) bounds ownership to one checkout.
**GREEN when:** The rewritten section licenses capture+claim from any checkout, names the single-releaser rule, and the `:57` never-probe sentence is unchanged; no shipped file still asserts the old boundary.
**Validation:** User confirmed (ask-tool answers: "Claim anywhere, one releaser"; "collisions are fixed by the agent as needed").

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 2, items 3–4).

## Addendum (2026-08-04)

From REQ-094's review (Minor finding, folded here because this REQ owns the lines): as of REQ-094, `actions/work-reference.md:55`'s "the pipeline does not detect, coordinate, or recover a second owner" and `docs/work-guide.md:91`'s "does not coordinate a second owner" are **partly false for a reason unrelated to this REQ's own scope** — crash recovery now *detects* another checkout's live claim by its `writer:` label and reports it (`claim held by <writer>, not touched`); it still doesn't coordinate or recover one. The Execution Model rewrite must account for that already-shipped behavior rather than rediscovering it. Also reword `actions/work-reference.md`'s Step-10 template line "recovery classifies each `working/` file by name" → by name *and label* if this REQ touches that paragraph.

---
*Source: approved plan, Phase 2*
