---
title: "Lessons from REQ-119: An off-vocabulary route warns on the board like domain does"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-119-an-off-vocabulary-route-warns-on-the-boa.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-119: An off-vocabulary route warns on the board like domain does

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

REQ-116 made the board normalize `route` and REQ-117 gave `domain` an unrecognized flag plus a `board.Warnings` entry. `route` got the normalization and not the warning, so `route: z` now reaches the card as `Z` with no footprint anywhere — the same silence REQ-117 was written to remove, one field over.

## Solution summary

`RequestTicket` gained `OriginalRoute` and `RouteUnrecognized`; the route read site derives the flag with `isKnownSchemaFieldValue` while keeping the case-folded value. `collectDomainWarnings` became `collectSchemaFieldWarnings`, iterating both contract fields the board reads, so the contract's warning wording still exists in exactly one place.

## What worked

**What worked:** Taking the reviewer's finding at face value only after checking it. Codex's claim was that the board was out of lock-step with the contract for `route`, and reading the read site confirmed it exactly — normalization without a recognition result. The finding was also *structurally* predictable: REQ-116 and REQ-117 were captured from the same review round and split by field, so the first shipped normalization before the second had established the channel. Splitting one contract leg across two REQs leaves a window where the fields disagree.

**What didn't:** Nothing failed, but the first instinct — derive the flag with a second `resolveSchemaField` call for symmetry with domain — would have been wrong. For route that call returns the empty-string default, so the flag would have been right and the value destroyed. `isKnownSchemaFieldValue` on the already-normalized value is the correct pairing when a field has no default.

**Worth knowing:** The board's two contract-field warnings now share one collector, so adding a third field is a slice entry plus its read-site flag — but the two halves are useless apart, and nothing enforces that pairing except the doc comment. A field given a flag and no `Original*` value would warn naming an empty string; a field given `Original*` and no flag would never warn at all.

## Back-reference

See `do-work/archive/UR-025/REQ-119-route-unrecognized-warning.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c327f24`.
