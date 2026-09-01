---
title: "Lessons from REQ-292: Move Durations label placement into the browser and delete the measured-face constants"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-01/REQ-292-move-durations-label-placement-into-the-.md]
related:
  - page: REQ-266-name-builds-beside-the-js-renderer-s-mea
    rel: complements
  - page: REQ-277-state-the-mark-label-face-constant-s-rea
    rel: complements
  - page: REQ-291-browser-behavior-probe-lane-beside-the-n
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-292: Move Durations label placement into the browser and delete the measured-face constants

Part of the [[concept-duration-estimation-and-breaks]] cluster.

## What the REQ was about

Durations label placement runs in Go, at generate time, before a browser exists. It reserves
`durationsLabelCharacterWidthUnits = 7.15` units per character (`durations.go:250`) and the
renderer stacks rows on a 13-unit pitch (`board-durations.js:60`). Both numbers describe one
Linux container's answer to `--font-sans`, and `board.css:65` ends that stack in the open
`sans-serif` generic — so the constants describe a face they never meet.

Measured live on a Mac while writing the solutions report: the board's own stack draws a
**13.0000**-unit line box against the 13-unit pitch, which is past the **12.97** that
`durationsMeasuredLabelBoxHeightUnits` pins as the recorded maximum. Eight of the nineteen
faces available on that machine exceed the pitch outright, and Arial Black exceeds the width
model at 7.34 units per character. The vertical half is a false clearance claim; the
horizontal half is the overprinting defect UR-051 was raised to fix, arriving by another
route.

## Solution summary

Moved Durations direct-label placement out of Go and into the renderer, where it sizes each label from the text the engine actually draws. Deleted the per-character width model and every hand-transcribed measured-face constant. The remainder count is now produced by the code that decides what fits. Ten Go tests lost their subject: seven are re-pinned as browser probes asserting what a reader would see, three are recorded as deliberately dropped, and a guard keeps the no-measured-constants property true going forward.

## What worked

Deleting in an order that kept the build compiling — sample fields, then the functions that used them, then the constants, then the payload — so every step's breakage named exactly what still referenced the thing being removed. `go vet` became the worklist. And running the real board in a browser at the end: the probes were all green before that, and only the rendered chart shows that the remainder sentence actually clears its neighbour on the row.

## What didn't work

A regex-driven "drop this function and its doc comment" helper, used to delete sixteen declarations, repeatedly clipped one character off the following declaration — leaving `f// comment`, `ffunc`, and a bare `/` at end of file. Every one was caught by the compiler, so nothing shipped wrong, but it cost several rounds. For a bulk deletion of Go declarations, deleting by explicit start/end line after listing them beats a comment-walking heuristic; the heuristic's failure mode is silent damage to the *neighbour*, not to the target.

Also: the first port measured the remainder reserve from a string that did not match the one drawn ("with a reversed span" against "reversed"). Caught by reading, not by a test, and it would have under-reserved by ~100 units. When two call sites must compose the same sentence, the composer belongs in one function — which is what `composeDurationsRemainderText` now is.

## Worth knowing

The reason this defect was undetectable, not merely undetected. A width model multiplies a character count by a constant, so it returns the *same* number for every face — which means the slots it assigns never move, and a wider face draws past them silently. No amount of pinning the constant could catch that, because the constant was never wrong in its own terms; it was answering a question about a face nobody had. That is why the replacement test asserts the geometry *changes* under a wider face rather than asserting any particular width: the property that matters is responsiveness, not accuracy.

## Back-reference

See `do-work/archive/UR-061/REQ-292-measure-durations-labels-in-the-browser.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ce28510`.
