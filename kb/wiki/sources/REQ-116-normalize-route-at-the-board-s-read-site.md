---
title: "Lessons from REQ-116: Normalize route at the board's read site and correct 0.174.15's board-wide claim"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-116-normalize-route-at-the-board-s-read-site.md]
related:
  - page: REQ-117-an-unrecognized-domain-must-leave-a-foot
    rel: complements
  - page: REQ-118-the-normalize-flag-must-stop-calling-voc
    rel: complements
  - page: REQ-119-an-off-vocabulary-route-warns-on-the-boa
    rel: complements
  - page: REQ-120-shipped-files-stop-citing-the-export-ign
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-116: Normalize route at the board's read site and correct 0.174.15's board-wide claim

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

The board parses `route` verbatim (`tools/queue-kanban/model.go` — `Route: coerceScalarToString(fields["route"])`), so a REQ written `route: a` reaches the card as lowercase `a`. Wire that read through the normalizer REQ-111 already added, and correct the 0.174.15 changelog entry's claim that the board honors the Schema Read Contract for all nine fields — it does not, and only `domain` was ever wired.

## Solution summary

The board's `route` read in `parseRequestTicket` now goes through `normalizeSchemaField("route", …)` behind a present-value-only guard, replacing the bare `coerceScalarToString`. Four parse-level test cases pin the behaviour: lowercase uppercases, a canonical letter is unchanged, a padded lowercase letter uppercases, an absent route stays empty, and an unrecognized letter is reported case-folded rather than blanked. `CLAUDE.md`'s display-parsed field enumeration gained `route` (see D-01). The changelog correction for 0.174.15's board-wide claim is written in the Commit Phase entry, not as a rewrite of the shipped entry.

## What worked

Writing the test at parse level rather than at the normalizer. The normalizer's own table test (`{"route", "a", "A"}`) was green the entire time the board was wrong — the test existed, passed, and proved nothing about the field the user sees. Choosing the altitude of the assertion mattered more than the number of assertions.

## What didn't work

Nothing failed, but one instinct was wrong: `resolveSchemaField` is the natural sibling call and would have been a defect here. Its contract substitutes the field's documented default, and route's default is the empty string, so `route: z` would have arrived as absent — indistinguishable from a REQ with no route at all, in the exact field re-triage reads to find the problem. The two helpers are one word apart and differ in whether the caller may invent a value.

## Worth knowing

The reason route drifted is recorded in the wrong place to prevent it. `CLAUDE.md`'s lock-step sentence names the fields the board parses for display and obliges any contract change to be mirrored in `model.go` — and `route` was never in that list, so a field that was parsed, badged and drawer-rowed carried no mirroring obligation. When a "keep these in sync" rule is expressed as a field enumeration, a field's absence from the list is silent permission to drift. Five of the contract's nine fields remain unread by the board on purpose; the guard against re-drift is the enumeration, not the code.

## Back-reference

See `do-work/archive/UR-024/REQ-116-board-route-normalization.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `2a2cd59`.
