---
title: "Lessons from REQ-231: Keep Panel A's direct labels clear of the mark band"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-04/REQ-231-keep-panel-a-s-direct-labels-clear-of-th.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-231: Keep Panel A's direct labels clear of the mark band

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

In the Durations view's overflow lane, the first row of direct labels sits at the same height as the marks themselves, so in a dense lane a label can be crossed by a *neighbouring* mark. REQ-226 stopped labels from overprinting each other; this is the remaining overlap, between a label and a dot that is not its own.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)

## What worked

- Reproducing RED by restoring the pre-fix geometry while exposing only the two constants the new test reads by name. A straight revert gave `undefined: durationsLabelTopCount` — a compile error, which proves nothing about geometry. Exposing the names and leaving the positions wrong produced `row 0's text box [33, 46] intersects the mark band [35, 45]`, which is the defect stated in the test's own words. The same trick isolated the second test: keep the whole geometry fix, neuter only the selection skip, and exactly one of the two tests fails.

## What didn't work

- Trusting the label *count* as a legibility proxy. The dense fixture drew 2 labels where the pre-fix render drew 27, which looked like a regression until a second fixture with scattered magnitudes drew 5 of 6. The first fixture correlated magnitude with x-position perfectly, so every top-N candidate landed in the same crowded corner — an adversarial input for top-N, not a defect in it. One fixture cannot tell a design's tail case from its common case.
- Reading the geometry test as sufficient. It reads renderer constants, so it proves the *declared* rows clear the *declared* band; it cannot see a leader tick crossing a glyph or a divider drawn in a colour nobody can see. Measuring `getBoundingClientRect()` intersections in the live DOM is the assertion that actually answers "do any two things touch" — 55 pairs before, 0 after, from the rendered document rather than from the source.

## Worth knowing

- Two reasons a sample can now go unlabelled — selection passed it over, or placement could not fit it — where there used to be one. Three comments still said "could not fit", and the restatement sweep is what caught them; the whole suite passed either way, because no test asserts on prose. When a change adds a *second* cause for an existing outcome, every sentence naming the first cause is now a half-truth.

## Back-reference

See `do-work/archive/UR-051/REQ-231-keep-durations-labels-clear-of-the-mark-band.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `720f23c`.
