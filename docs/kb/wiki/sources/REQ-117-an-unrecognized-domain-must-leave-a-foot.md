---
title: "Lessons from REQ-117: An unrecognized domain must leave a footprint on the board, not become general in silence"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-117-an-unrecognized-domain-must-leave-a-foot.md]
related:
  - page: REQ-116-normalize-route-at-the-board-s-read-site
    rel: complements
  - page: REQ-118-the-normalize-flag-must-stop-calling-voc
    rel: complements
  - page: REQ-119-an-off-vocabulary-route-warns-on-the-boa
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-117: An unrecognized domain must leave a footprint on the board, not become general in silence

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Since REQ-111 wired `domain` through `resolveSchemaField`, a typo'd `domain: quantum` renders on the board as a `general` badge with no warning anywhere — the recognized flag is discarded at the call site in `tools/queue-kanban/model.go`. Before that change the typo was at least visible verbatim, so the typo is now *better hidden* than it was. Add the flag and the data warning that `testing_status` already has, keep the `general` fallback, and delete the comment claiming the board has no channel for it.

## Solution summary

`RequestTicket` gained `OriginalDomain` and `DomainUnrecognized`; the domain read site keeps `resolveSchemaField`'s recognized flag instead of discarding it; `collectDomainWarnings` mirrors `collectTestingWarnings` and is appended to `board.Warnings` in `buildBoard` right after it. The `general` fallback and the absent-domain guard are unchanged. The comment claiming the board has no warning channel for domain is replaced with what is actually true, including why the old reasoning was wrong.

## What worked

Reading the sibling field before designing anything. `testing_status` had already solved this exact contract leg — flag on the ticket, collector over all tickets, appended in `buildBoard` — so there was no design decision left to make, just a pattern to copy. The tests were a copy too, which is the strongest signal the shapes really match.

## What didn't work

The premise in the code comment was false, and it had been reviewed at 98% when it shipped. "The board has no warning channel for it" was three greps from being disproven: `board.Warnings` is declared with the words "surfaced, never silently dropped" in the same file, the sibling field feeds it, and the frontend renders it in a banner. A stated *reason* for a design choice is a factual claim and a review has to check it like any other — a plausible-sounding justification in a comment is exactly where a wrong premise survives longest.

## Worth knowing

`board.Warnings` is a free UI channel. Anything appended to it prints in `summary` and renders in the board's data-warnings banner with no frontend work, because `web/board.js` reads `boardData.warnings` generically. Two consequences: a new warning class costs one `append` line, and noise is cheap to introduce — which is why the recognized-alias and absent-field cases have their own test. The contract's absent-field carve-out (`resolveSchemaField` returns recognized=true for an empty value) is what keeps a real queue from warning on nearly every REQ.

## Back-reference

See `do-work/archive/UR-024/REQ-117-unrecognized-domain-warning.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `42f71e2`.
