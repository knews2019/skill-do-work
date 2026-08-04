---
source_type: req_lesson
req_id: REQ-035
req_path: do-work/archive/UR-007/REQ-035-lock-multi-claim-representation.md
date: 2026-07-29
domain: general
module: actions
tags: [actions, concurrent, claims, orchestrator, lock]
---

# Lessons from REQ-035: Represent concurrent claims in the orchestrator lock and Crash Recovery gate

## What the REQ was about

REQ-032's parallel-dispatch gate says every concurrently dispatched REQ runs Steps 2–9 "including the orchestrator lock's `claimed_req` bookkeeping" — but `claimed_req` is a single string per session, and Crash Recovery's per-file gate only protects the one REQ it names (and only against *other* sessions). Give the lock a way to represent N concurrent claims by one orchestrator, and make the Crash Recovery gate honor them (including on the same session's own Step 10 → Step 1 loop iteration).

## Solution summary

Added a canonical list-valued `claimed_reqs` field to the orchestrator lock (holder and each `coexisting_sessions[]` entry), with `claimed_req` retained as an additive derived legacy mirror (`claimed_reqs[0]`, or `null` when empty; never made array-shaped). Every reader/writer now tells one story: the parallel-dispatch gate bullet gained the only-live-claimant mixed-version precondition; Step 2 appends the claimed id (re-deriving the mirror); Step 8 substep 6, the Step 8 failure path, and the mid-run blocked-flip each remove *only* this REQ's id after the move; the Crash Recovery per-file gate now reads each entry's `claimed_reqs` list and skips any file in any fresh (≤45m) claim set *including this session's own* (the "session other than this one" clause deleted, so a Step 10 → Step 1 loop no longer re-queues a co-dispatched sibling); Crash Recovery step 2 additionally clears `write_set`; both heartbeat recompute rules (holder + coexisting-entry) recompute `claimed_reqs` scoped to this session's own dispatched `working/` files (never a bare listing of everything in `working/`); the JSON examples, schema prose, warn/prompt/refusal templates, acquire-time initializer, and cleanup.md Pass 0 gate were all updated in lock-step. Four contract-regression ratchets pin the field's presence in work.md/work-reference.md/cleanup.md and the Crash Recovery gate's same-session-inclusion phrasing. Per-merge post-merge verification became the default whenever >1 REQ is in flight (discharging REQ-033's dormant per-batch finding). No frontmatter shape change, so `tools/queue-kanban/model.go` is untouched; the serial/floor path is behaviorally unchanged.

## What worked

- The `claimed_reqs`-list + derived-`claimed_req`-mirror design (additive back-compat, one owning writer per field, recompute-scoped-to-own-claims) handled every reader/writer cleanly and kept the serial/floor path behaviorally identical. The adversarial contradiction-hunter earned its cost here: it found two stale gate restatements that the reviewer's requirements-walk also caught but that the plan's edit-site list structurally could not.

## What didn't work

- The plan discovered its edit sites by grepping the literal token `claimed_req`. That grep **cannot** find the places that state the same gate rule *in prose without the token* — "skips any file another live session still actively claims" (work.md:119) and "not actively claimed by another live session" (work-reference.md:207). A token-driven edit list will always miss the behavior-phrased restatements of the same rule. This is the recurring failure mode for "one story everywhere" coherence edits, and it's the second time this class of miss has surfaced in UR-007 (the plan-missed *write* sites were caught by the explorer; these prose restatements slipped past both plan and explorer and needed the review).

## Worth knowing

- When a rule lives in prose in N places, pin BOTH the token presence AND the behavior phrasing in the primary files — a token-only ratchet (`claimed_reqs` present) stays green while a behavior restatement silently tells the old story. The fix added a phrasing ratchet (`this session's own co-dispatched claims` in work.md) alongside the token ratchets. Ratchets are still file-granular presence checks, so they pin that the right phrasing exists *somewhere* in the file, not that a specific line is correct — a second restatement in the same file could still drift; the residual is documented, not eliminated.

## Back-reference

See `do-work/archive/UR-007/REQ-035-lock-multi-claim-representation.md` for the full REQ — triage, implementation, review, and lessons. Commit `fd56267`.
