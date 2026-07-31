---
id: REQ-056
title: "Worktree-mode blocked flip: judge \"did edits land\" from the builder's branch, not the main tree"
status: completed
claimed_at: 2026-07-29T14:46:44Z
commit: 1b7d824
status_changed_at: 2026-07-29T15:07:51Z
route: B
domain: general
tdd: false
maintenance: false
prime_files: []
created_at: 2026-07-29T09:32:07Z
user_request: UR-007
depends_on: []
write_set: ["actions/work.md", "actions/work-reference.md"]
---

# Worktree-mode blocked flip: judge "did edits land" from the builder's branch, not the main tree

## What

In the Step 8 blocked-flip procedure (the guard that parks a REQ as `blocked` when a builder hits a missing external precondition *without* having done real work), add a worktree-mode branch: decide "did edits land this attempt?" by consulting the builder's `worktree-agent-*` branch (or the hand-back manifest) instead of `git diff` on the main tree.

## Why

In worktree dispatch mode a builder commits on its own branch, so the main tree reads clean even after substantial work — a worktree builder that hits a missing precondition *after* real work is wrongly flipped to `blocked` instead of archived as `failed` with its follow-up REQ. Rare path; pre-dates REQ-037. Approved by the user via `do-work clarify` on 2026-07-29 (surfaced by REQ-037's review, queued via REQ-042).

## Constraints

- The serial path is unchanged — the new check applies only when the REQ was dispatched to a worktree builder.

## Acceptance

- A worktree builder failing after real committed work classifies as `failed` (with follow-up), not `blocked`; a worktree builder failing before any work still flips to `blocked`; serial-mode behavior is untouched.

## Implementation Summary

Added a worktree-mode bullet to `actions/work.md` Step 8's mid-run blocked-flip procedure, between the "Both must hold to flip" test and the "If both hold" action. It substitutes evidence for clause (1) only (clause (2) and the serial working-tree check are untouched) and splits by how far Step 6's hand-back got: holding a `<merge_hash>` proves the builder committed — the hand-back stops on `Already up to date.` rather than fabricating a merge commit — so edits landed and the flip is refused. With no merge, branch existence is probed first (`git rev-parse --verify -q '<operative_name>'`): a branch that was never created is the before-any-work case and flips to `blocked`, and the probe has to precede the count because `rev-list` on a missing branch exits fatal and prints no number. Only a branch that resolves goes to `git rev-list --count HEAD..<operative_name>` — `HEAD` is the integration branch, since the orchestrator never leaves it — where `0` also still flips. The bullet judges from git rather than the handed-back manifest, and states that uncommitted edits inside the builder's worktree are not "landed" since the main tree stays pristine and the existing worktree sweeps own the leftover.

**Files changed:**
- `actions/work.md` — new worktree-dispatch bullet in the Step 8 mid-run blocked-flip procedure judging "did edits land this attempt" from the builder's branch/merge instead of the main tree.

## Review

Adversarial review workflow (4 Opus lenses -> 2 diverse refuters per Important+ finding; 10 agents): verdict FIX-THEN-PASS -> fixes applied -> PASS.

- Upheld + fixed: `<integration_branch>` was an undefined placeholder (occurred exactly once in the shipped skill; resolving it to the repo default branch inverts the classification on a feature-branch integration) — replaced with self-defining `HEAD` + inline gloss.
- Adopted minor (flagged by three lenses): `git rev-list --count` on a never-created branch exits fatal and prints no count — branch-existence probe (`git rev-parse --verify -q '<operative_name>'`) now precedes the count and is itself the before-any-work signal.
- Skipped (recorded, deliberate): a `do-work/`-exclusion analogue for the branch probe — worktree builders never write `do-work/` per work-reference.md's State-stays-home contract, so the qualifier holds by construction.
- Killed by refutation: count-0-after-merge ambiguity (the `<merge_hash>` case-split governs before the count is consulted); a claimed stale closed enumeration in work-reference.md (misread of what the list enumerates).
