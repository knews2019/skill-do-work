---
title: "Lessons from REQ-199: Publish portfolio snapshot before canonical refresh"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-199-publish-portfolio-snapshot-before-canoni.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-199: Publish portfolio snapshot before canonical refresh

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Make the Yes/unavailable snapshot branch publish its exclusive, no-clobber snapshot from the retained bytes before atomically refreshing the canonical portfolio file. A snapshot publication failure must leave the prior canonical summary unchanged instead of partially completing the promised two-output branch.

This is a standalone user-visible publication-order contract and cannot fold into the existing image-generation or ID-normalization follow-ups: its fix applies only to the portfolio writer's canonical-plus-snapshot transaction.

## Solution summary

Promoted portfolio output into a shipped helper that verifies one retained source, publishes an optional snapshot exclusively before canonical replacement, handles numeric collisions, and atomically refreshes canonical. Wired the action/guide, registered the primitive, and added branch/failure behavior replays plus semantic/inventory contracts.

## What worked

- Promoting publication into a shipped helper made branch order and failure outcomes directly executable and replayable.
- Exclusive snapshot publication before canonical replacement closes the original partial-update defect on ordinary paths.

## What didn't work

- Hard-link identity proved byte equality at publication but coupled later mutable canonical writes back into the durable snapshot.
- Two-operand `ln` and `mv` interpret directory destinations as containers; without exact-path type guards, successful status does not prove the requested path was published.

## Worth knowing

Immutable evidence needs independent file contents, not merely a second directory entry for the same inode. Publication helpers must test exact target type/identity because core utilities treat directories differently from file destinations.

## Back-reference

See `do-work/archive/UR-042/REQ-199-publish-portfolio-snapshot-before-canonical-refresh.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `74f2220`.
