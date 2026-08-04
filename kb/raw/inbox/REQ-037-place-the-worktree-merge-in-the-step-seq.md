---
source_type: req_lesson
req_id: REQ-037
req_path: do-work/archive/UR-007/REQ-037-worktree-merge-placement-evidence-repointing.md
date: 2026-07-29
domain: general
module: actions
tags: [actions, re-point, evidence-consuming, worktree, merge]
---

# Lessons from REQ-037: Place the worktree merge in the step sequence and re-point the evidence-consuming steps

## What the REQ was about

REQ-033's worktree dispatch mode says who merges and how (`git merge --no-ff`, dependency order) but never *when* — and every evidence-consuming step downstream assumes uncommitted main-tree work. State the merge point explicitly, and make the diff-based checks read the right evidence when the mode is on.

## Solution summary

Placed the worktree merge in the pipeline (orchestrator merges at hand-back — end of Step 6, before Step 6.25) and made every evidence-consuming step read the merged diff in worktree mode. `actions/work-reference.md`'s Worktree Dispatch Mode gained a single source-of-truth block defining the merge point, the `pre`/`merge_hash` capture, and the merge range `<pre>..HEAD` (with its merge-base-collapse rationale, the hold-in-memory-never-`HEAD^1` rule, and the consumer list); its "Post-merge verification" default was aligned to per-merge (matching work.md — the folded REQ-035 gap); and its Commit & Metadata Procedure gained worktree-mode staging + validation branches. `tools/checks/qualify.sh` reads an optional `DO_WORK_DIFF_RANGE` and branches both the file-list check and the debug-artifact scan on it — the unset (serial) branch is byte-identical to before, the exec bit and exit-code contract are unchanged, and the header now documents the env var plus `Exit 2`. `actions/review-work.md` reads `<pre>..HEAD` in worktree pipeline mode (Two-Modes row, nothing-to-review exit, Step 4 Get-the-Diff); standalone `git show <commit>` left unchanged (recorded as a Discovered Task). `actions/work.md` Step 6.3, Step 7, Step 8, Step 9, and the "one commit per request" Rules bullet all reference the merge point/range and reconcile Step 9 for already-committed merged work (stage only changelog/version/metadata; write the merge commit's hash). Serial/floor path behaviorally unchanged; no ratchet added; no schema change.

## What worked

- Defining the merge range once (single source of truth in Worktree Dispatch Mode) and pointing every consumer at it made the re-pointing auditable — the hunter could grep every worktree/range site and confirm the consumer list matched the re-pointed sites. Gating every change on "worktree dispatch mode" / `[ -n "$diff_range" ]` kept the serial floor path provably byte-identical.

## What didn't work

- The plan captured a stable anchor (`merge_hash`) for exactly the "HEAD moves" reason, then defined the range with live `HEAD` anyway — a self-inconsistency that only bites in per-batch mode, which the requirements-walk missed and two independent finders caught. Lesson: when you capture a stable handle *because* a live ref moves, use that handle at every reference, not just the one that motivated capturing it.

## Worth knowing

- A self-check that greps a diff for `console.log|TODO|FIXME|debugger` cannot cleanly qualify a change to *its own* detection pattern (or to any file that carries those tokens as data) — qualify.sh false-positived on this REQ. Recorded as a `[low]` follow-up; until fixed, the orchestrator must eyeball such a FAIL rather than bounce the builder.

## Back-reference

See `do-work/archive/UR-007/REQ-037-worktree-merge-placement-evidence-repointing.md` for the full REQ — triage, implementation, review, and lessons. Commit `1348d11`.
