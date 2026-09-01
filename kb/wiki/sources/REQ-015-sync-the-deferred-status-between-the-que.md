---
title: "Lessons from REQ-015: Sync the `deferred` status between the queue-kanban parser and the Schema Read Contract"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-015-sync-the-deferred-status-between-the-que.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-015: Sync the `deferred` status between the queue-kanban parser and the Schema Read Contract

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

`tools/queue-kanban/model.go` treats `deferred` as a recognized Needs-input/Blocked status (`isNeedsInputOrBlockedStatus`, ~line 410, plus the column comment ~line 116), but the Schema Read Contract in `actions/work-reference.md` does not list `deferred` in its status enum, and nothing in the skill ever writes `status: deferred`. Make the two vocabularies agree — recommended direction: remove `deferred` from `model.go`.

## Solution summary

Removed the producer-less `deferred` status from the queue-kanban parser's recognized-status set so it exactly matches the Schema Read Contract enum (`actions/work-reference.md`); a hand-edited `status: deferred` ticket now surfaces in Needs-input/Blocked with an unrecognized-status warning instead of being silently blessed.

## What worked

- Anchoring RED directly on the REQ's captured `## Red-Green Proof` (the failing `isNeedsInputOrBlockedStatus("deferred")` assertion) made the TDD cycle mechanical; repointing existing `deferred` test expectations at the warning path preserved coverage instead of deleting it.

## What didn't work

- Nothing — no dead ends on this one.

## Worth knowing

- Synthetic tickets in `model_test.go` must set `OriginalStatus` (not just `Status`) for unrecognized-status warning assertions to exercise the real code path — the warning text is built from `ticket.OriginalStatus`, so omitting it makes a warning assertion pass trivially. Also, `TestNormalizeStatus` tests generic lowercase/trim passthrough, not status recognition — don't read its cases as a sanctioned-status list.

## Back-reference

See `do-work/archive/UR-003/REQ-015-deferred-status-vocabulary-sync.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `27f1005`.
