---
id: REQ-242
title: Stop Panel B's slowest-day annotation colliding with its own title
status: completed
status_changed_at: 2026-08-18T14:00:46Z
created_at: 2026-08-18T12:09:46Z
user_request: UR-051
addendum_to: REQ-237
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T13:05:12Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
claimed_at: 2026-08-18T13:05:12Z
route: B
completed_at: 2026-08-18T14:00:46Z
commit: 48263dd
kb_status: promoted
kb_entry: REQ-242-stop-panel-b-s-slowest-day-annotation-co.md
---

# Stop Panel B's Slowest-Day Annotation Colliding With Its Own Title

## What

In the Durations view, Panel B's slowest-day annotation is drawn at `y = 355` while Panel B's own title sits at `y = 350` (`DURATIONS_MEDIAN_TITLE_Y`). The two overlap: on a synthetic fixture the annotation `209 min` renders directly through the words "paused and broken spans excluded".

## Context

Found while reviewing REQ-237 by rendering a dense fixture and looking at it. **It is not a REQ-237 regression** — the same annotation sits at the identical `x = 357.2, y = 355.0` on a board built from the pre-REQ-237 binary, checked side by side. It is pre-existing and was simply never looked at on a fixture whose slowest day lands under the title text.

It is invisible on this repository's own board because the annotation's x-position depends on which day is slowest, and here that day falls clear of the title's width. That is luck, not design — which is why it wants pinning rather than nudging.

The annotation reuses the `durations-mark-label` class, so it is not part of either label band's row packing and is not covered by REQ-231's mark-band geometry test or by REQ-237's row-fill test. Nothing in the suite looks at it at all.

## Requirements

- Panel B's slowest-day annotation does not overlap Panel B's title at any x-position the annotation can take, including when the slowest day is the leftmost one.
- The annotation stays associated with the bar it describes — moving it somewhere it no longer reads as belonging to that day is not a fix.
- No change to Panel A, Panel C, or the label bands; no change to `describeAtPointer`'s panel boundary.
- A test pins the separation, so the next person to move `DURATIONS_MEDIAN_TITLE_Y` finds out.

## Red-Green Proof

**RED prompt/case:** an assertion that the slowest-day annotation's text box and Panel B's title text box do not intersect, read from the renderer's own constants the way `TestDurationsLabelRowsClearTheMarkBands` reads them — evaluated at the annotation's worst-case x, not at whichever x this repository's data happens to produce.
**Why RED now:** the title's baseline is 350 and the annotation's is 355, so their boxes intersect wherever their x-ranges do; reproduced on a fixture as `209 min` drawn through the title.
**GREEN when:** the assertion passes and a rendered fixture whose slowest day sits under the title shows the two clear of each other.
**Validation:** Review finding on REQ-237; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-237 I rendered the Durations chart on a test board and found the little "209 min" note that marks the slowest day printed straight through the heading above it — the note sits five units below a heading that is taller than five units. It has been like this since before any of today's work; it does not show on your own board only because the slowest day happens to fall to the right of where the heading's text ends, which is chance rather than design. The fix is small — move the note, or move the heading, and add an assertion so the next person who shifts either one finds out. I am asking rather than doing it because "move which one, and where" is a look-and-feel choice about a chart you read regularly, and there is more than one reasonable answer. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — it only shows on contrived data and the note is legible enough where it is.
  → Confirmed: Yes, add to queue (builder picks placement, pinned by a test). [2026-08-18, via do-work clarify]

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — move the slowest-day annotation to a fixed baseline clear of both neighbours
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the geometric assertion, driven at the extremes of x

**Files I will NOT touch:** `durations.go`, `durations_test.go` (REQ-241's files, already merged), `DURATIONS_MEDIAN_TITLE_Y`, anything under `do-work/`, the version files, `CHANGELOG.md`.

**Acceptance criteria (restated from REQ):**
- [x] The annotation does not overlap Panel B's title at any x it can take, including the leftmost day
- [x] It stays associated with the bar it describes
- [x] No change to Panel A, Panel C, the label bands, or `describeAtPointer`'s panel boundary
- [x] A test pins the separation, so the next person to move `DURATIONS_MEDIAN_TITLE_Y` finds out

## Pre-Flight

**Git:** ✓ clean at claim (`2ad71eb`), with REQ-241's retuned constants already merged
**Tests baseline:** ✓ `maintainer-verify.sh` exit 0 before dispatch
**Dependencies:** ✓ Go 1.26.1, Playwright/Chromium for live-DOM measurement

*Checked by work action*

## AI Execution State (P-A-U Loop)

<!-- Filled from the builder's hand-back; a builder may not write do-work/ in worktree dispatch
     mode. Source of record: do-work/runs/work-2026-08-18-124358/REQ-242-handback.md -->


- [x] **PLAN** — Read `_dev/primes/prime-kanban-board.md`, `CLAUDE.md`, and the three always-on crew members, then the whole of `web/board-durations.js` and the geometry tests in `durations_test.go` (`TestDurationsLabelRowsClearTheMarkBands`, `TestDurationsLabelRowPitchClearsTheLabelTextBox`) for the shape a geometric assertion takes here.

  Confirmed the defect reproduces on my own fixture with REQ-241 already in the tree (measurements below), so the REQ is live.

  Approach: **move the annotation, not the title** — `DURATIONS_MEDIAN_TITLE_Y` and `describeAtPointer` are untouched. Above the bar is not available at all: the slowest day is by definition the tallest bar, an over-ceiling one tops out at `DURATIONS_MEDIAN_TOP` (368), and everything above that is the title's box. Inside the plot is not available either: `barWidth` bottoms out at 4 units, so a label there overprints its neighbours on a dense board. That leaves the strip below the panel's baseline, which is reachable at **every** x and every bar height — the clearance stops depending on which day is slowest, which is the actual defect. Extract the drawing into a named function so a behavior probe can drive it at worst-case x, then assert its text box against the renderer's own constants.

- [x] **APPLY** — Exactly the two planned files. `web/board-durations.js`: new `DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y = 467`, new `DURATIONS_TICK_BASELINE_DROP = 4` (replacing three literal `+ 4`s so a test can read the tick row), new `drawDurationsSlowestDayAnnotation(svg, slowestDay, dayCentreX)` replacing the inline block, and the bar-height-dependent y arithmetic deleted. `generate_test.go`: four measured face constants, the six probe cases, and `TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours`.

- [x] **UNIFY** — `git diff --stat` is the two files below and nothing else; `git diff --name-only <merge-base>...HEAD -- do-work/` is empty; `VERSION`, `skills/do-work/VERSION`, `actions/version.md` and `CHANGELOG.md` are untouched. Reviewed both files line by line: `gofmt -l .` clean, `go vet` clean via maintainer-verify, no debug prints, no leftover scratch (all fixtures, boards and screenshots live under `/tmp/qk-242` and `/tmp/board-242-*`). Re-generated the board from the committed source and re-measured it in the browser as the last step. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** Moved Panel B's slowest-day annotation to a single fixed baseline below the median line, clear of both Panel B's title and the axis tick row, and named `DURATIONS_TICK_BASELINE_DROP` so a test can read the tick row's offset instead of a literal buried in an attribute. Panel B's title did not move.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 (unpiped; `echo $?` printed `0` on its own line)

**Red-green validation:**
- **RED 0 — the defect itself, in a live browser, measured with REQ-241 already in the tree.** On a fixture whose *leftmost* day is the slowest, `getBoundingClientRect()` intersection is real in both axes and the screenshot reads `209 miB · Median minutes per active day …`. **The collision survives REQ-241** — so this REQ was not resolved by it.
- **RED 1 —** the assertion failing for the reason it exists, with the code in place and one rule broken (`DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y` set back to the shipped `355`): ✗ → ✓.
- **A second collision, found by rendering rather than by reasoning.** The first fix passed a fresh, purpose-built geometric assertion — and the render showed it printing through Panel B's `0` axis tick, which lives in the y-axis gutter and is therefore invisible as a neighbour at every x except the extreme left. Exactly the luck-of-x this REQ exists to remove, reproduced a second time inside the same REQ.

*Verified by work action*

### Red-Green Evidence (verbatim from hand-back)


### RED 0 — the defect itself, in a live browser

Fixture: 27 archived REQs over 29 days whose **leftmost** day is the slowest, with a median of 209 min (over the 45-min ceiling, under the 4-hour read-time exclusion). Built with the worktree's own binary; measured in an isolated headless Chromium via Playwright, with `location.href` returned from the same `evaluate` call as every box:

```
"href": "file:///tmp/board-242-before/index.html",
"titleText": "B · Median minutes per active day · paused and broken spans excluded",
"titleBox":  { "left": 110.33, "right": 559.39, "top": 615.75, "bottom": 631.75 },
"annotations": [{
  "text": "209 min", "x": "38.9", "y": "355.0", "anchor": "middle",
  "box": { "left": 70.91, "right": 117.13, "top": 623.15, "bottom": 637.15 },
  "intersectsPanelBTitle": true
}]
```

`getBoundingClientRect()` intersection is real in both axes (x overlap 110.33–117.13, y overlap 623.15–631.75). The screenshot shows the title rendering as `209 miB · Median minutes per active day …`. **The collision survives REQ-241** — measured with `DURATIONS_LABEL_ROW_HEIGHT = 13` and `durationsLabelCharacterWidthUnits = 6.75` already in the tree.

### RED 1 — the assertion failing for the reason it exists

Code and test in place, one rule broken: `DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y` set to `355`, the shipped geometry for an over-ceiling slowest day.

```
--- FAIL: TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours (0.25s)
    generate_test.go:2053: leftmost day, over the ceiling: the annotation's text box [344.00, 357.80]
    intersects panel B's title's box [337.90, 352.80] — the two overprint wherever their x ranges
    meet, and x follows whichever day is slowest
```

(Re-run against the committed test, so the quote is the message the shipped assertion prints. The RED 2 quotes below were captured mid-build, before later comment edits shifted line numbers by one.)

(I did not run the test against the pre-change blob: that source has no `drawDurationsSlowestDayAnnotation`, so the probe would fail on a missing anchor — a reference error, which per the brief's rule 1 proves nothing. Reproducing the shipped baseline inside the new code is the meaningful RED, and RED 0 is the defect itself.)

### RED 2 — the render catching what the arithmetic did not

My first fix put the annotation at baseline **461**. The suite passed. The render did not:

```
"annotations": [{ "text": "209 min", "y": "461",
  "box": { "top": 737.62, "bottom": 751.62, "left": 70.91, "right": 117.13 },
  "intersectsPanelBTitle": false,
  "intersectingTicks": [ { "text": "0", "y": "452",
      "box": { "top": 727.9, "bottom": 741.9, "left": 94.53, "right": 101.69 } } ] }]
```

Panel B's `0` axis tick lives in the y-axis gutter, so it is only ever in the annotation's way when the slowest day is the leftmost — the identical luck-of-x that hid the original defect. Extending the test to the tick row and re-running at 461 fails for that reason:

```
--- FAIL: TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours (0.25s)
    generate_test.go:2052: leftmost day, over the ceiling: the annotation's text box [450.00, 463.80]
    intersects panel B's "0" axis tick's box [441.00, 454.80] — …
```

Both other neighbour assertions were confirmed non-vacuous the same way: at 470 it fails on panel C's title box `[471.90, 486.80]`, and at 449 it fails with *"the annotation's text box starts at 438.00, above panel B's baseline at 448.00"*.

### GREEN

```
=== RUN   TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours
--- PASS: TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours (0.30s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.685s
```

And the render, regenerated from the committed source, on the same leftmost-slowest fixture:

```
file:///tmp/board-242-final/index.html
B · Median minutes per active day · paused and broken spans excluded
{'text': '209 min', 'x': '38.9', 'y': '467',
 'box': {'left': 70.91, 'right': 117.13, 'top': 744.1, 'bottom': 758.1},
 'intersectsPanelBTitle': False, 'intersectsPanelCTitle': False,
 'intersectingTicks': [], 'intersectsAnyBar': False}
```

The annotation's x range (70.91–117.13) still overlaps the title's (110.33–559.39) — that is the point: the separation is now purely vertical, so no x can bring it back. I looked at both renders; the title reads clean and "209 min" sits directly under its bar. Also rendered this worktree's own board (slowest day mid-plot, "42 min"): same placement, no intersection with anything, still legibly attached to its bar.

## Decisions


- **D-01 — The annotation moved; the title did not.** `DURATIONS_MEDIAN_TITLE_Y` is unchanged, so `describeAtPointer`'s A/B boundary (`pointerY <= DURATIONS_MEDIAN_TITLE_Y - 12`) is unchanged and REQ-241's 1.364 units of headroom above the title are untouched. No proof about pointer probes is owed, because no pointer geometry moved.
- **D-02 — The guarantee is vertical, therefore x-free.** The annotation now sits at one fixed baseline for every day and every bar height. This is deliberate over any scheme that dodges the title horizontally: the defect existed *because* the clearance depended on x, and a horizontal dodge would only move the dependence. The test states it by driving the extremes of x and asserting the baseline never varies.
- **D-03 — Below the baseline, not beside the bar's top.** Placing the label beside the bar's top edge would also clear both titles, but at a dense day count `barWidth` bottoms out at 4 units and a ~43-unit label would print across a dozen neighbouring bars. The strip under the baseline is the only region free at every density.
- **D-04 — `DURATIONS_TICK_BASELINE_DROP` is a named constant now.** The tick row is a neighbour other text has to clear, and a test cannot read a literal buried in an attribute. It replaced three copies of `+ 4` that all meant the same thing.
- **D-05 — The measured face constants live in `generate_test.go`, beside the test that uses them,** not in `durations_test.go` with REQ-241's. That file was out of my write set this wave. If someone later consolidates the measured-face constants into one block, `durations_test.go` is the natural home and the mark-label pair is a near-duplicate of what is already there.

## Lessons Learned


- **A chart's empty space is only empty at the x you looked at.** The first fix passed a test suite that included a fresh, purpose-built geometric assertion — and the render showed it printing through panel B's `0` axis tick. The tick lives in the y-axis gutter, so it is invisible as a neighbour at every x except the extreme left, which is precisely the luck-of-x this REQ existed to remove. Reproducing that failure twice in one REQ is the strongest case yet for `prime-kanban-board.md`'s "generate a board and look at it": the second collision was found the same way as the first, and neither was reachable by reasoning over the constants the fix was about.
- **When a defect is "invisible because of where the data happens to fall", the fix is to remove the dependence, not to widen the margin.** A larger gap above the title would still have been a clearance that held for some slowest days and not others.
- **A measured face is per-Chromium.** The 11px label face measured 10.4278 ascent for REQ-241 and 10.1853 here on a different build. Both constants round up and away from the model, so a test written that way survives the difference — but a test asserting an exact measured number would not.

## Builder Pushback


- The brief's third point — "test at the annotation's worst-case x" — is satisfied, but by making x irrelevant rather than by pinning a worst-case x. Six cases spanning the leftmost, mid-plot and rightmost day are in the test and all assert the same baseline; if a future change ties the baseline back to the bar, the case set fails on the disagreement rather than on any one coordinate. Flagging it because "pin the worst-case x" and "prove x cannot matter" are different tests and I chose the second.
- **Unrelated pre-existing quirk, not fixed here:** the leftmost day's bar is drawn at x = 38.9, i.e. *left of* `DURATIONS_MARGIN_LEFT` (54) and outside the plot area. `xOfEpoch` maps the day bucket's midnight while `timeStart` is the first completion instant, so a day whose first REQ completed at 09:00 renders 15 hours' worth of axis to the left of the plot. Visible on the live board too (the `29 May` bar sits in the gutter). It is a placement bug of its own — worth a REQ, outside this write set.

**Orchestrator resolution.** The first point is accepted, and the builder's choice is the better one. The brief said "test at the annotation's worst-case x"; the builder instead made x irrelevant and asserted the baseline never varies across six cases spanning the leftmost, mid-plot and rightmost day. Those are different tests, and proving x cannot matter is stronger than pinning the x that happens to be worst today — the defect existed *because* clearance depended on x. If a future change ties the baseline back to the bar, the case set fails on the disagreement rather than on any one coordinate.

The second point — the leftmost day's bar drawn at x = 38.9, left of `DURATIONS_MARGIN_LEFT` and outside the plot area, because `xOfEpoch` maps the day bucket's midnight while `timeStart` is the first completion instant — is a real pre-existing placement bug visible on this repository's own board. Correctly left alone and routed as a follow-up.

## Review

**Reviewer:** independent subagent, read-only, judged from git state, from a 310-node SVG dump diff, and from 11 hand-built fixtures measured in a freshly launched isolated Chromium per fixture.
**Score:** 88% — **PASS-WITH-FINDINGS** (first pass, pre-remediation)

### Findings returned to the builder

1. **The test is a 6-point sample, not the x-free / height-free proof D-02 claims — and a passing mutant reproduces the defect on this repository's own board.** Mutating only the `y:` expression to `(dayCentreX > 700 && dayCentreX < 1100) ? BASELINE - 112 : BASELINE` leaves the suite green while putting the annotation back at y=355, inside Panel B's title, for any slowest day in that band. **The real board's slowest day sits at x=881.4.** A second mutant banded on the median (`> 1 && < 44`) also passes; the sampled medians are only {0, 45, 209} and the real board's slowest median is **42 min**. A smooth x-dependence control correctly fails — so the net catches continuous drift and misses banded drift, and the bands that slip through are where the real data lives. The shipped code *is* x-free by inspection; this is a regression-net gap, not a live bug.
2. **The strip has three occupants, not two, and `board-durations.js:87-89` says two.** `.durations-month-line` spans y=84 to y=572 and therefore crosses y=467; on a 90-day fixture whose slowest day falls on a month boundary it intersects the annotation's box by 12.963 units², passing between the "9" and the " min". Not a regression — it crossed the old y=355 position too — but the enumeration is incomplete in exactly the way this REQ's own lesson warns about, and `neighbourBoxes` omits it.

### Findings routed elsewhere

3. **The `xOfEpoch` bug the builder reported is worse than reported, and is being filed at above-cosmetic priority.** On the real board the leftmost Panel B bar spans x=37.1–49.1, entirely left of `DURATIONS_MARGIN_LEFT` (54). On a **one-day** board the annotation renders at x=-3330 and the bar at x=-3342 — completely off-canvas, so Panel B renders empty. Two-day boards likewise (x=-336.5). Root cause is `timeStart` being the first completion *instant* while the day buckets are *midnights*; the suggested fix is flooring `timeStart` to its UTC midnight before computing `timeSpan`.
4. **REQ-241's recorded 1.364-unit headroom measures 0.185 units on Chromium 146** (title ascent 12.0372 / descent 2.7778 there, against 11.23 / 2.41 on REQ-241's build). Still positive, still zero intersections, and **not caused by REQ-242** — the SVG output for that region is byte-identical across the range. But D-03's "1.36 units is the whole budget" is per-Chromium and roughly 7× optimistic on at least one current build. This is what the integration seam on `durationsMeasuredAxisTitleAscentUnits` resolved conservatively, and it is why the larger value won.
5. `medianMinutes.toFixed(0)` prints "0 min" for a 0.4-minute median — pre-existing, and the annotation exists to state a value a clipped bar cannot.

### Upheld against the builder's claims

- **`DURATIONS_TICK_BASELINE_DROP` proven value-preserving rather than argued:** all 310 SVG child nodes were dumped from the pre- and post- boards and diffed. **Exactly one line differs** — the annotation's `y=366.8 → y=467`. Every tick baseline, gridline, bar, mark and label is byte-identical. All three former copies of the literal were replaced; the one remaining bare `4` is a different quantity, correctly left alone.
- **`describeAtPointer` unchanged from git and by driving events.** The extracted function body is byte-identical across the range (same `shasum` at both endpoints), and synthesized `mousemove` sweeps at three x positions on pre- and post- boards resolve `A` at y ≤ 338 and `B/C` at y ≥ 339 on both — matching REQ-241's recorded boundary exactly.
- The REQ's own failing case (leftmost slowest, 209 min) now shows **zero collisions** and zero Panel B title intersections; the defect is gone in the render.
- REQ-231's and REQ-237's guarantees hold **by construction** via the 310-node diff, and were re-measured anyway: 0 same-row label overlaps, 0 label/mark overlaps.
- Association with the bar survives at 70 and 90 days, where the label is 3–4× the bar width — nearest-bar-centre distance 0.00–0.05 units in every fixture.
- Scope exact: two files. The merge introduced nothing — `diff <(git diff 2ad71eb b139805) <(git diff 51beffc a8ef062)` is empty.

*Reviewed by work action*

## Review Remediation

Both findings closed on the builder's own branch and re-merged; the cumulative merge range for this REQ is `51beffc..48263dd`.

**F1 — the six-point sample is now a sweep plus an exact structural check.** The builder reproduced the M5 mutant on the merged branch before changing anything (`ok … 0.825s`, green while reproducing the defect), then closed it two ways and showed each biting on its own:

1. A **deterministic sweep** of 10 006 cases — the six named extremes plus 200 pseudo-random positions across `[-400, 1400]` **crossed** with 50 medians across `[0, 240]`, crossed rather than paired so a rule keyed on position *and* height together has nowhere to hide. Seeded, so a failure names coordinates the next run reproduces, and positions are generated at one decimal place because the renderer writes x through `toFixed(1)` — a case the probe cannot echo back exactly would be testing the rounding rather than the geometry.
2. An **exact structural check**, `assertDurationsAnnotationBaselineIgnoresItsInputs`, which reads the shipped function out of the generated page and requires its baseline expression to mention neither `dayCentreX` nor `medianMinutes`. This is the one that matters: a sweep is still a sample, and a band narrower than its spacing slips through. The structural check is what makes one measurement at one x a statement about every x — which is what D-02 claimed all along and did not previously pin.

Both mutants now fail, each against each mechanism in isolation, with the failure naming the offending expression verbatim. Runtime is unchanged at under a third of a second.

**F2 — the month rule is named as an accepted crossing rather than omitted.** The builder measured it independently rather than taking the review's number: 96 REQs over 90 days with the slowest day pinned to 1 July, a month boundary mid-plot. **12.963 units² of overlap, matching the review exactly.** Method note worth keeping: a vertical `<line>`'s client rect is zero wide because the stroke is not in the geometric bbox, so the overlap has to be computed in user units from `getBBox()` and the rule's `x1` ± `strokeWidth/2`.

The comment now says the strip has three occupants, that two are cleared and the third is an accepted crossing, and why. The test carries the same distinction: `neighbourBoxes` stays the clearance list, and the month rule gets its own block asserting that the annotation's box is *inside* the rule's span — so the crossing genuinely cannot be avoided by any legal baseline, which is why it is not a clearance failure — plus the two properties that make it acceptable, **read from `board.css` rather than claimed**. Setting the rule to `stroke-width: 3` makes that assertion fail with "it is allowed to cross the slowest-day annotation only because it is a hairline".

**Verification:** `bash _dev/tests/maintainer-verify.sh` exit 0, unpiped, on the builder's branch and again on the integrated tree.

*Remediation verified by work action*
