---
title: "Lessons from REQ-241: Reconcile the Durations label metrics with the face actually rendered"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-04/REQ-241-reconcile-the-durations-label-metrics-wi.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-241: Reconcile the Durations label metrics with the face actually rendered

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Two constants that describe the Durations label face disagree with what the browser actually draws. Neither causes a visible collision today, but both are now load-bearing in a way they were not before, because REQ-237 made the label rows actually fill up.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)

## What worked

- **A per-character width model has a right answer, and it is the worst case over the label space, not the average over one board.** The 6.2 constant was probably fitted to a plausible-looking string; the comment then described the *intent* (generous) rather than the *result* (7% short), and nothing could catch the gap because no test knew what the face draws. Two things fix this class permanently: enumerate the label space when the format function makes it enumerable, and record the browser's measurement as a named constant so a Go test can hold the model to it without a browser in the loop.
- **A constant that clears the box the code declares can still fail the box the browser draws, and vice versa — assert both.** Rounding an ascent up (11 from 10.43) is safe and worth keeping, but it means the declared box and the drawn box are two independent claims. A single assertion against whichever happens to be larger today rots the moment someone trims the other. RED #2 was worth producing precisely because the second branch is unreachable at today's values and would otherwise read as dead code.
- **A REQ can name a consequence that measurement then shows is not required — say so rather than performing it.** The brief expected Panels B and C to shift with the row pitch. Shifting them would have broken the *other* stated requirement (`describeAtPointer` resolving the same panel for the same pointer). Measuring the actual clearance showed the pitch increase fit inside existing headroom, which satisfies both requirements at once. The useful artifact is not "no change needed" but the number: 1.36 units of remaining budget, and the specific pitch (15) at which the trade becomes real.
- **Isolate a multi-constant change to attribute its effects, instead of reasoning about them.** I first wrote that the three lost labels split two-to-the-width-retune and one-to-the-pitch, which sounded plausible and was wrong: building one binary per constant showed the pitch costs zero labels and the width costs all three. The reason is structural and was sitting in the code the whole time — `DURATIONS_LABEL_ROW_HEIGHT` exists only in the renderer, and `durations.go` packs each row purely horizontally, so a pitch change *cannot* alter which labels are drawn. When a change moves two numbers, one build per number is cheap and turns a guess into a fact.
- **A "0 overlaps" result can be true and still be hiding a halved margin.** Both before and after renders show zero same-row overlaps; the difference is that the tightest gap went from 3.08 to 11.35 units against a rule claiming 6. When a guarantee is checked only as a boolean, a margin can erode most of the way to failure without any test noticing. Where a rule states a *quantity*, measure the quantity, not just its sign.
- **Follow-up candidate (not filed):** the remainder reservation over-reserves by ~39 units at the new width (D-04), and Panel B's slowest-day annotation collides with Panel B's own title on a dense board — visible in both the before and after fixture screenshots at `/tmp/board-241/shot-{before,after}-fixture.png`, where `233 min` overprints `…spans excluded`. Both are pre-existing, neither is in this write set, and the second is a genuine visual defect on any board whose slowest day is near the ceiling.

## Back-reference

See `do-work/archive/UR-051/REQ-241-reconcile-durations-label-metrics-with-the-rendered-face.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `90c74b7`.
