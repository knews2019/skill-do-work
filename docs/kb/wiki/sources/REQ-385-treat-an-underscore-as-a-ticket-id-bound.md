---
title: "Lessons from REQ-385: Treat an underscore as a ticket-id boundary on both surfaces"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-385-treat-an-underscore-as-a-ticket-id-bound.md]
related:
  - page: REQ-381-index-cited-ticket-ids-and-let-the-filte
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-385: Treat an underscore as a ticket-id boundary on both surfaces

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

`\b` counts `_` as a word character, so the mention pattern behaves differently around an underscore
than a reader expects. Change the boundary on both the Go and the client side together, in one
commit, so the agreement test stays green.

## Solution summary

Underscores now delimit ticket IDs. Boundary checks run after candidate matching, preventing the regular expression from falling back to a shorter UR alternative. Ticket candidates remain consumed even when resolution intentionally suppresses them; only non-ticket runs keep the previous retry behavior.

## Worth knowing

Compound-first alternation does not guarantee compound-first behavior when a failing boundary permits fallback. Consume before checking boundaries, and keep intentionally suppressed ticket candidates consumed too: a caller's retry can recreate a fallback the regular expression no longer performs.

## Back-reference

See `do-work/archive/UR-075/REQ-385-treat-underscore-as-a-ticket-id-boundary.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `259b1479`.
