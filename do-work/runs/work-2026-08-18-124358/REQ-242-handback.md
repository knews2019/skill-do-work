# REQ-242 — Stop Panel B's Slowest-Day Annotation Colliding With Its Own Title

## Branch

`worktree-agent-REQ-242-stop-panel-b-annotation-colliding-with-its-title` — one commit, `b139805`.

## P-A-U

- [x] **PLAN** — Read `_dev/primes/prime-kanban-board.md`, `CLAUDE.md`, and the three always-on crew members, then the whole of `web/board-durations.js` and the geometry tests in `durations_test.go` (`TestDurationsLabelRowsClearTheMarkBands`, `TestDurationsLabelRowPitchClearsTheLabelTextBox`) for the shape a geometric assertion takes here.

  Confirmed the defect reproduces on my own fixture with REQ-241 already in the tree (measurements below), so the REQ is live.

  Approach: **move the annotation, not the title** — `DURATIONS_MEDIAN_TITLE_Y` and `describeAtPointer` are untouched. Above the bar is not available at all: the slowest day is by definition the tallest bar, an over-ceiling one tops out at `DURATIONS_MEDIAN_TOP` (368), and everything above that is the title's box. Inside the plot is not available either: `barWidth` bottoms out at 4 units, so a label there overprints its neighbours on a dense board. That leaves the strip below the panel's baseline, which is reachable at **every** x and every bar height — the clearance stops depending on which day is slowest, which is the actual defect. Extract the drawing into a named function so a behavior probe can drive it at worst-case x, then assert its text box against the renderer's own constants.

- [x] **APPLY** — Exactly the two planned files. `web/board-durations.js`: new `DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y = 467`, new `DURATIONS_TICK_BASELINE_DROP = 4` (replacing three literal `+ 4`s so a test can read the tick row), new `drawDurationsSlowestDayAnnotation(svg, slowestDay, dayCentreX)` replacing the inline block, and the bar-height-dependent y arithmetic deleted. `generate_test.go`: four measured face constants, the six probe cases, and `TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours`.

- [x] **UNIFY** — `git diff --stat` is the two files below and nothing else; `git diff --name-only <merge-base>...HEAD -- do-work/` is empty; `VERSION`, `skills/do-work/VERSION`, `actions/version.md` and `CHANGELOG.md` are untouched. Reviewed both files line by line: `gofmt -l .` clean, `go vet` clean via maintainer-verify, no debug prints, no leftover scratch (all fixtures, boards and screenshots live under `/tmp/qk-242` and `/tmp/board-242-*`). Re-generated the board from the committed source and re-measured it in the browser as the last step. `bash _dev/tests/maintainer-verify.sh` exits 0.

## Files Changed

```
 .../tools/queue-kanban/generate_test.go            | 175 +++++++++++++++++++++
 .../tools/queue-kanban/web/board-durations.js      |  69 +++++---
 2 files changed, 225 insertions(+), 19 deletions(-)
```

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` — the annotation's placement moves from "7 units above the bar's top" (a baseline of 355 for an over-ceiling day) to a single fixed baseline of 467, below panel B's own baseline and centred between the `0` tick row and panel C's title. The drawing moves into `drawDurationsSlowestDayAnnotation` so a probe can call it; `DURATIONS_TICK_BASELINE_DROP` replaces the three literal `+ 4` tick offsets so the test can read the tick row instead of copying the number.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` — `TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours` plus the four browser-measured face constants it needs. Six cases span both dimensions that decided whether the defect showed (where the slowest day sits, how tall its bar is) and assert one unchanging baseline whose text box clears panel B's title, panel B's `0` axis tick, panel C's title, and the plot the bars occupy — and that the annotation still says what it said and stays centred on its day.

## Red-Green Evidence

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

## Verification

```
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	15.069s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (4.71s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	4.905s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.301s
Maintainer verification passed.
0
```

Run unpiped from the worktree root; `echo $?` on its own line printed `0`.

## Integration Seams

None. Both files are mine alone this wave, no shared registry or cross-REQ text was touched, and `durations.go` / `durations_test.go` are untouched. One thing to note for whoever writes the changelog entry: the visible change is *where panel B's slowest-day figure sits*, which a reader of the board will notice.

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

## Pushback

- The brief's third point — "test at the annotation's worst-case x" — is satisfied, but by making x irrelevant rather than by pinning a worst-case x. Six cases spanning the leftmost, mid-plot and rightmost day are in the test and all assert the same baseline; if a future change ties the baseline back to the bar, the case set fails on the disagreement rather than on any one coordinate. Flagging it because "pin the worst-case x" and "prove x cannot matter" are different tests and I chose the second.
- **Unrelated pre-existing quirk, not fixed here:** the leftmost day's bar is drawn at x = 38.9, i.e. *left of* `DURATIONS_MARGIN_LEFT` (54) and outside the plot area. `xOfEpoch` maps the day bucket's midnight while `timeStart` is the first completion instant, so a day whose first REQ completed at 09:00 renders 15 hours' worth of axis to the left of the plot. Visible on the live board too (the `29 May` bar sits in the gutter). It is a placement bug of its own — worth a REQ, outside this write set.

---

# Addendum — Review Remediation

Branch merged with `main` first (bringing REQ-241's remediation, REQ-245, and the orchestrator's seam edit to the `durationsMeasuredAxisTitleAscentUnits` comment block — preserved verbatim, not reverted). Remediation commit: **`ff842ee`**.

```
 .../tools/queue-kanban/generate_test.go            | 186 ++++++++++++++++++---
 .../tools/queue-kanban/web/board-durations.js      |  26 ++-
 2 files changed, 181 insertions(+), 31 deletions(-)
```

`git diff --name-only HEAD~1 HEAD -- do-work/ VERSION CHANGELOG.md skills/do-work/VERSION skills/do-work/actions/version.md` is empty. `web/board.css` was mutated temporarily to prove one assertion bites and restored; `git status` is clean apart from the two write-set files.

## F1 — the six-point sample is now a sweep plus an exact structural check

Finding accepted in full. D-02 claimed an x-free property and the test pinned six points; those are not the same statement, and the gap was where the real data lives.

**Reproduced first, on the merged branch, before changing anything.** The M5 mutant, applied to the `y:` expression alone:

```
y: (dayCentreX > 700 && dayCentreX < 1100) ? BASELINE - 112 : BASELINE
→ ok  github.com/knews2019/skill-do-work/queue-kanban  0.825s
```

Two mechanisms now close it, and each was shown to bite on its own.

**1. A deterministic sweep.** The six named extremes plus 200 pseudo-random positions across `[-400, 1400]` crossed with 50 medians across `[0, 240]` — 10 006 cases, crossed rather than paired so a rule keyed on position *and* height together has nowhere to hide. Seeded LCG, so a failure names coordinates the next run reproduces. Positions are generated at one decimal place because the renderer writes x through `toFixed(1)`, and a case the probe cannot echo back exactly would be testing the rounding.

**M5 failing against the sweep alone** (structural check temporarily commented out to isolate it):

```
--- FAIL: TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours (0.30s)
    generate_test.go:2128: swept day at x=936.8 with a 233.50-minute median: the annotation's
    text box [344.00, 357.80] intersects panel B's title's box [337.90, 352.80] — the two
    overprint wherever their x ranges meet, and x follows whichever day is slowest
```

The median-banded mutant, same isolation:

```
--- FAIL: … generate_test.go:2128: swept day at x=936.8 with a 40.27-minute median: the
    annotation's text box [344.00, 357.80] intersects panel B's title's box [337.90, 352.80] — …
```

**2. An exact structural check.** A sweep is still a sample: a band narrower than its spacing slips through. `assertDurationsAnnotationBaselineIgnoresItsInputs` reads the shipped function out of the generated page and requires its baseline expression to mention neither `dayCentreX` nor `medianMinutes` — which is precisely what makes one measurement at one x a statement about every x.

**M5 failing against it**, with the sweep in place (this is the test as committed):

```
--- FAIL: TestJavaScriptBehaviorDurationsSlowestDayAnnotationClearsItsNeighbours (0.28s)
    generate_test.go:2052: the annotation's baseline expression "y: (dayCentreX > 700 &&
    dayCentreX < 1100) ? DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y - 112 :
    DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y," reads dayCentreX — a baseline that depends on
    where the slowest day sits or how tall its bar is puts the clearance back at the mercy
    of the data, which is the defect REQ-242 fixed
```

And the median-banded mutant against it:

```
    generate_test.go:2052: the annotation's baseline expression "y: (slowestDay.medianMinutes > 1
    && slowestDay.medianMinutes < 44) ? … " reads medianMinutes — …
```

Clean run after both mutants were reverted: `PASS (0.28s)`. The sweep costs no measurable runtime — the whole test is under a third of a second, unchanged from the six-point version.

## F2 — the month rule, named and accepted

Finding accepted. `board-durations.js` now says the strip has **three** occupants, that two are cleared and the third is an **accepted crossing**, and why the crossing is acceptable. The test carries the same distinction: `neighbourBoxes` is the clearance list, and the month rule gets its own block at the end asserting (a) that the annotation's box is *inside* the rule's span — i.e. the crossing genuinely cannot be avoided by any legal baseline, which is why it is not in the clearance list — and (b) the two properties that make it acceptable, read from `board.css` rather than claimed.

**Measured it myself** rather than taking the number on trust. Fixture: 96 archived REQs over 90 days (15 May – 13 Aug 2026), slowest day pinned to **1 July**, a month boundary that falls mid-plot. Vertical `<line>` client rects are zero wide — the stroke is not in the geometric bbox — so the overlap is computed in user units from `getBBox()` and the rule's `x1` ± `strokeWidth/2`, with `location.href` returned from the same `evaluate` call:

```
"href": "file:///tmp/board-242-month/index.html",
"annotationUserBox": { "left": 617.35, "right": 658.25, "top": 456.815, "bottom": 469.778 },
"monthRulesUser": [
  { "x1": "262.2",  … "overlapUserUnits2": 0 },
  { "x1": "637.8",  "y1": "84", "y2": "572", "strokeWidth": "1px",
    "stroke": "rgba(17, 24, 39, 0.08)", "overlapUserUnits2": 12.963 },
  { "x1": "1026.0", … "overlapUserUnits2": 0 } ],
"panelBTitleOverlap": 0, "panelCTitleOverlap": 0, "tickOverlaps": []
```

**12.963 units², matching the review exactly.** The annotation is centred at x=637.8 and the 1 Jul rule sits at x=637.8, so it passes through the middle of "210 min". I looked at the render: at `rgba(17, 24, 39, 0.08)` and one unit wide the rule is barely perceptible where it crosses the label, and the same rule crosses panel A's reversed-band labels and every bar in panels B and C the same way. Accepted, and now asserted to stay that way:

```
--- FAIL: … generate_test.go:2171: the month rule's stroke-width is "3", not "1" — it is
    allowed to cross the slowest-day annotation only because it is a hairline
```

(produced by temporarily setting the rule to `stroke: var(--line-firm); stroke-width: 3` in `board.css`, then restoring it.)

## Verification

```
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	16.980s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (6.00s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	6.213s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.359s
Maintainer verification passed.
0
```

Run unpiped from the worktree root; `echo $?` on its own line printed `0`.

## Decisions

- **D-06 — Both mechanisms, not one.** The brief said either was sufficient. They fail for different reasons and neither subsumes the other: the sweep is a statement about what the renderer *does* and survives any refactor, while the structural check is exact where the sweep is dense-but-sampled. Given that a mutant literally reproduced the shipped defect on the maintainer's own board, "dense enough" did not feel like the right place to stop. The structural check is fifteen lines and its failure message tells the next person what to do if the expression legitimately becomes multi-line.
- **D-07 — The sweep is deterministic, not random.** A fixed LCG seed, so a failure prints coordinates that reproduce on the next run. A test that fails on one run in fifty and cannot be re-run into the same failure is worse than no test.
- **D-08 — The month rule is asserted as an accepted crossing, not as a clearance.** No baseline in the strip can clear a line spanning y=84 to y=572, so a clearance assertion there would be unsatisfiable and its absence would read as an oversight — which is exactly what F2 caught. The test instead asserts the crossing is unavoidable *and* that the rule stays one unit and soft. If someone later moves the month rules behind the panels or shortens their span, the first of those assertions fires and says the crossing belongs in the clearance list now.
- **D-09 — `durationsStyleDeclaration` is new and general.** Reading one property out of one rule in `board.css` is how a test holds a claim about how something is DRAWN against the sheet that draws it. `board.css` stays outside this REQ's write set; the helper only reads it.
- **D-10 — D-02 stands, but it is now earned rather than asserted.** The x-free framing was the right one; my mistake was banking a general claim on six points. The claim is unchanged; what changed is that the test now proves it.

## Lessons Learned (addendum)

- **A sample dressed as a proof is worse than an honest sample.** The six-point version was not wrong about the shipped code — it was wrong about what it *guaranteed*, and the hand-back's D-02 sold the stronger reading. Mutation testing is what told the difference; reading the test could not have.
- **Where a regression net has holes, they are rarely uniformly distributed.** Both surviving mutants banded exactly where the real board's data lives. That is not coincidence: a hand-picked case list gets picked from the *interesting* extremes, so the ordinary middle — which is where production data sits — is the part left uncovered.
- **When an enumeration cannot be completed, say the item is accepted rather than dropping it.** Two occupants named and a third silently omitted reads as a finished list. "Three occupants; two cleared, one accepted, here is why" is the same length and cannot mislead.
