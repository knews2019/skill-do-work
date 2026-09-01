---
title: "Lessons from REQ-114: The three remaining shell-logic extraction candidates, restated decay-free"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-114-the-three-remaining-shell-logic-extracti.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-114: The three remaining shell-logic extraction candidates, restated decay-free

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Carries forward the three extraction candidates the shell-logic census ranked but never captured, restated so they survive without the census's line-number table. Each is a consolidation of a primitive that is currently copy-pasted across several action files. Candidate B was split out, approved, and delivered as REQ-121; **Candidates A and C are not approved work** and must each clear the floor constraints before becoming a change.

Every candidate below is described by **what to grep for**, not by line numbers. That is the point: the census's citations went stale within hours of a single merge to `actions/work-reference.md`, so a durable record has to name the search rather than the coordinate.

## Solution summary

Recorded the durable current disposition: Candidate B is delivered by REQ-121, while Candidates A and C remain separate, unapproved candidates. Corrected this request's now-false original "None of these" status wording; no extraction implementation or action-prose change was made.

## What worked

Keeping the candidate record grep-based made it possible to verify the current state without reviving the census's stale line-number table.

## What didn't work

The inventory's original blanket statement that none of its candidates was approved became stale after Candidate B split into REQ-121; a disposition close-out needs to update that statement in both the audit and the REQ.

## Worth knowing

A queue run authorizes processing the inventory REQ, not selecting an unapproved candidate for implementation. Candidate A and Candidate C remain separate decisions.

## Back-reference

See `do-work/archive/UR-022/REQ-114-residual-extraction-candidates.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `b857933`.
