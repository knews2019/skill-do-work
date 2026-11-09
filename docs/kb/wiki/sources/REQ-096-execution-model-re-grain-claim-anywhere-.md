---
title: "Lessons from REQ-096: Execution-model re-grain: claim anywhere, one releaser; dispatch widened to any tree"
type: source-summary
topic_cluster: worktree-and-parallel-dispatch
sources: [raw/processed/2026-09-01/REQ-096-execution-model-re-grain-claim-anywhere-.md]
related:
  - page: REQ-094-checkpoint-writer-label-crash-recovery-i
    rel: depends-on
  - page: REQ-095-two-clone-acceptance-run-checkpoint-pois
    rel: complements
  - page: REQ-097-assigned-to-advisory-field-schema-line-s
    rel: complements
  - page: REQ-099-automatic-wave-dispatch-the-work-loop-co
    rel: complements
  - page: REQ-101-docs-adr-multi-checkout-guide-and-the-se
    rel: complements
  - page: REQ-104-label-less-checkpoint-entries-locally-mo
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-096: Execution-model re-grain: claim anywhere, one releaser; dispatch widened to any tree

Part of the [[concept-worktree-isolation-and-parallelism]] cluster.

## What the REQ was about

Rewrite the Execution Model contract (`actions/work-reference.md:53–61`) from "one queue owner per checkout, cross-session ownership unsupported" to the user's chosen model: **any checkout may capture and claim/build; exactly one designated releaser checkout runs the release tail** (merge integration, version bump, `CHANGELOG.md` entry, archive moves, UR closure). Widen Worktree Dispatch (`:275–341`) so a builder tree may be a spawned worktree, a user workspace, a clone, or a remote/cloud sandbox.

## Solution summary

Re-grained the ownership contract from one-owner-per-checkout to claim-anywhere/one-releaser, widened builder trees to any own-tree-own-branch shape, and folded in all three addenda. Every echo of the old boundary was found by grepping the condition rather than a phrase list, and both suite assertions that pinned the rewritten sentences moved with them in this commit. No check was loosened: one was retargeted, one was retitled, and one was added.

## What worked

- Grepping for the *condition* rather than the phrase. `one queue owner per checkout` had one hit;
  the boundary it states had nine, across five files and two suite rationales. A phrase-list sweep would
  have shipped a renamed section still cited by five dangling anchors.
- Checking the suite for pins **before** writing, not after. The heading pin at `:326` was not in this
  REQ's brief or the batch handdown — only the invariant count was — and finding it first turned a
  surprise failure into a planned two-assertion change.

## What didn't work

- The first instinct was to leave the section heading alone to avoid the five-citation ripple. That would
  have shipped `## Execution Model — Exclusive Session` above a paragraph beginning "Any checkout may
  capture and claim" — precisely the stale-name drift this repo's maintenance rules exist to stop. The
  ripple turned out to be five mechanical edits verifiable by one grep.
- Editing the suite's reservation *rationale* looked out of scope at first. It is not: with claim-anywhere
  in contract, "cross-session REQ allocation is outside the product contract" reads as banning exactly
  what REQ-097 is about to ship, and a reviewer following the message would reject a valid change. A
  ratchet's justification is part of the ratchet.

## Worth knowing

- The suite extracts `Current-REQ relevance` by its **bold markers**
  (`sed -n '/^\*\*Current-REQ relevance\./,/^\*\*Three-attempt stop\./p'`), not by the section heading —
  which is why renaming the heading was safe for it. Both bold lead-ins must keep their exact text and
  their order, or that whole block of assertions silently extracts nothing.
- `tools/queue-kanban/prime-do-kanban.md` is the one shipped file allowed to name the retired
  exclusive-session premise, because its lesson entry quotes it as history. Any future sweep of that
  phrase must exempt it, and the suite comment says so.
- `one releaser per queue` is now a counted invariant: state it once, point at it everywhere else.

## Back-reference

See `do-work/archive/UR-018/REQ-096-execution-model-regrain.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `7024c4a`.
