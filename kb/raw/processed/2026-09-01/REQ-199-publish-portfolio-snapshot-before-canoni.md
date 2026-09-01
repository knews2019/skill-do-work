---
source_type: req_lesson
req_id: REQ-199
req_path: do-work/archive/UR-042/REQ-199-publish-portfolio-snapshot-before-canonical-refresh.md
date: 2026-08-15
domain: general
module: _dev/primes
tags: [general, publish, portfolio, snapshot, canonical]
---

# Lessons from REQ-199: Publish portfolio snapshot before canonical refresh

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

## Back-reference

See `do-work/archive/UR-042/REQ-199-publish-portfolio-snapshot-before-canonical-refresh.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `74f2220`.
