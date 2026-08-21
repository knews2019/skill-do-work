---
id: REQ-292
title: Move Durations label placement into the browser and delete the measured-face constants
status: completed
created_at: 2026-08-19T14:36:44Z
status_changed_at: 2026-08-19T14:36:44Z
claimed_at: 2026-08-21T01:15:48Z
completed_at: 2026-08-21T01:38:03Z
kb_status: pending
commit: ce28510
route: C
user_request: UR-061
domain: general
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-291]
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/generate.go
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Move Durations Label Placement into the Browser and Delete the Measured-Face Constants

## What

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

The user chose to close both by moving placement to where the font is (UR-061), accepting
the cost the solutions report priced.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- **Placement moves to the browser and measures real extents.** Each label's width comes from
  the rendered text, not from a character count times a constant. `durationsLabelCharacterWidthUnits`
  and the measured-face constant block in `durations_test.go` are deleted, not adjusted —
  the point of this REQ is that no such number survives.
- **The algorithm's behaviour is preserved, not reinvented.** The existing rules are
  deliberate and each has a REQ behind it: longest-span-first order (`durationLabelMagnitudeOrder`,
  stable over completion order so equal spans keep left-to-right precedence); first row where
  the text touches nothing already placed; a sample that fits nowhere is counted, never
  silently dropped; and the two-pass reserve, where the remainder sentence's room at the last
  row's right edge is only held back on a pass that actually dropped something. Port these,
  and say in the code which ones you kept and why.
- **The remainder count moves with the placement.** Go emits it today and the renderer reads
  `durations.labels`; once Go stops deciding what fits, a count computed there is a lie.
  Whatever places the labels produces the count.
- **Measure before paint.** Render the text to measure it, then position it — offscreen, or
  hidden, or in a pass the reader never sees. The chart must not visibly reflow after first
  paint. Name the approach chosen.
- **Every deleted Go test's property is re-pinned as a browser probe, or its loss is recorded
  with a reason.** These assert against the model being deleted:
  `TestDenseOverflowLabelsStayBoundedAndNeverOverlap`,
  `TestReversedLabelPlacementIsIndependentOfOverflowDensity`,
  `TestOverflowLabelsGoToTheLongestSpans`,
  `TestClusteredOverflowLabelsFillBothLabelRows`,
  `TestDurationLabelGeometryMatchesTheRenderer`,
  `TestDurationsLabelWidthEstimateCoversTheRenderedFace`,
  `TestDurationsLabelRowPitchClearsTheLabelTextBox`,
  `TestDurationsLabelRowsClearTheMarkBands`,
  `TestDurationsLastLabelRowClearsPanelBTitle`,
  `TestDurationsMeasuredConstantsNameTheirChromiumBuild`.
  Deleting one without a replacement is how the defect that started UR-051 comes back. A
  property that genuinely no longer exists — the width-model bounds, for instance, once there
  is no width model — is recorded as deliberately dropped, in the commit and in this file.
- **Any browser-measured number that survives the move names the build it was measured on.**
  Folded in from REQ-266, cancelled as superseded on 2026-08-20: that REQ asked for build
  provenance beside `DURATIONS_LABEL_ROW_HEIGHT` and `DURATIONS_LABEL_TEXT_ASCENT` in
  `web/board-durations.js`, two numbers this REQ deletes. The rule outlives the numbers —
  measuring at runtime is meant to leave no hardcoded face behind, so if one does survive,
  it carries its build exactly as REQ-252 requires on the Go side, and the mechanism that
  keeps that true (a JS-side check or a stated review convention) is recorded either way.
  If none survives, say so; that is the requirement met, not skipped.
- **The browser probes assert what a reader would see**, not what the code computed: no two
  drawn labels overlap; labels went to the longest spans; both rows fill when spans cluster;
  the remainder count equals the number actually not drawn; the last row clears the Panel B
  title; rows clear the mark bands.
- `bash _dev/tests/maintainer-verify.sh` exits 0, including on a machine with no browser,
  where the lane skips per REQ-291.

## Red-Green Proof

**RED:** a board rendered where `--font-sans` resolves to a face wider than 7.15 units per
character draws labels past the slots placement assigned them, and nothing in the suite
notices. Confirmed by measurement: Arial Black draws 7.34.

**GREEN:** placement asks the engine how wide the text actually is, so the same board packs
correctly in any face, and a browser probe fails if two drawn labels ever intersect.

## Requirements Deliberately Excluded

- The tick labels' `--font-mono` stack and the axis-title face carry the same open-generic
  exposure. Out of scope here, on purpose — this REQ is the mark labels.
- No re-budgeting of Panel A's or Panel B's vertical geometry. REQ-265 narrowed the Panel B
  ceiling to 0.10 model units and that constraint is untouched; if measured placement needs
  more vertical room, that is a separate REQ with both constraints visible at once.

## Context

Second of two under UR-061, and unbuildable before REQ-291 lands.

The direction was chosen by the user against the solutions report's own recommendation. The
report preferred bundling a font (O1) purely on cost — O2 was called "the textbook-correct
answer" and rejected only because it puts a browser in a test path that has never had one.
The user's instruction supplied that browser rather than disputing the reasoning, so the
cost line is answered and the recommendation no longer applies. Do not relitigate this in
the plan step.

Full argument, both options, and the live measurements:
`ai-reports/2026-08-19_1345_durations-label-face-robustness/index.html`.

---

## Triage

**Route: C** - Complex

**Reasoning:** An algorithm moves between two languages, a width model and its constants are deleted, ten Go tests lose their subject, and the payload contract changes shape. The direction is settled by the REQ and is not to be relitigated, but the sequencing — what can be deleted before what still compiles, and which properties survive the move — had to be worked out before editing.

**Planning:** Required

## Plan

### The port

`board-durations.js` gains the placement algorithm, sized from `getComputedTextLength()` instead of a character count. Every rule is ported with its reason, since each has a REQ behind it: magnitude order stable over completion order (REQ-237), first-row-that-fits with a separation gap consulted against **every** interval on the row (not a high-water mark, because the walk is by magnitude not by x), anchor-before-the-mark with after-the-mark as the plot-edge fallback (REQ-231), fits-nowhere is counted never dropped, and the two-pass reserve where only a pass that dropped something is redone.

### Measure before paint

The SVG is appended to `chartHost` **before** marks and labels are drawn, so text nodes created afterwards are in the document and measurable. The whole render is one synchronous call, and a browser does not paint until the task yields — so measuring every label, then positioning them all, produces no visible reflow without needing an offscreen clone or a hidden pass. That is the approach, and it is named in the code: "before paint" is a property of the task here, not a visibility trick.

### The deletion, in an order that keeps the build alive

Sample fields and the aggregate's label bands first, then the placement functions, then the constants, then the payload fields. The width model is the target; the x-axis domain is not, so `durationLabelTimeRange`/`durationLabelPlotX` move into `durations_test.go` rather than being deleted — the parity pin they serve (REQ-248's UTC-day anchoring against the renderer's buckets) is still real and still needs a second side.

### The ten tests

Six properties re-pin as browser probes. Three genuinely lose their subject and are recorded as dropped. One is replaced by a guard that keeps the requirement true going forward. Detail in `## Decisions` D-03 and in the new file's header.

### Plan validation

- **Requirement coverage:** placement-measures-real-extents → the port; behaviour-preserved → the ported rules, each named; remainder-moves-with-placement → produced by the placer; measure-before-paint → named; every-deleted-test-re-pinned-or-recorded → D-03; surviving-number-names-its-build → D-04; probes-assert-what-a-reader-sees → the browser probes; verify-exits-0-with-and-without-a-browser → both run.
- **No orphan tasks.**
- **Scope sanity:** 5 tasks. Under the flag.

*Planned inline by the orchestrator*

## Exploration

- `durations.go:252` `durationsLabelCharacterWidthUnits = 7.15` and `:299` `durationLabelWidthUnits` are the width model; `:408-539` is placement.
- `DURATIONS_PLOT_WIDTH` in the renderer is `1200 - 54 - 18 = 1128`, **identical to Go's `durationsPlotWidthUnits`**, so the port is coordinate-compatible; the only difference is that Go worked plot-local (0..1128) and the renderer works in absolute SVG units.
- `board-durations.js:258` appends the svg to `chartHost` before any mark is drawn — the fact that makes measurement real rather than zero, and the whole reason no offscreen pass is needed.
- The renderer already draws the label text (`sample.id + " " + formatDurationMinutes(...)`), so the port needs no second copy of the view's copy — `formatDurationLabelMinutes` in Go existed only to size a string it did not draw.
- Ten tests reference the model. `TestJavaScriptBehaviorDurationsDayBucketsStayInsideThePlot` also uses `durationLabelTimeRange` — for the axis domain, not for placement, which is why that function survives in the test file.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — the ported placement, the measurement pass, the remainder produced by the placer
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — width model, placement, constants and label fields deleted
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — `labelRow`/`labelAnchor`/`labels` removed from the payload
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — the ten tests and their helpers removed; the axis-domain reference kept here
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go` (new) — the re-pinned properties
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the payload-verdict probe rewritten to the properties that survive

**Files I will NOT touch:** `model.go` — declared in the REQ's write set, but the label fields turned out to live in `durations.go` and `generate.go`; nothing in `model.go` needed changing. Panel A's or Panel B's vertical geometry, and the tick/axis faces, are excluded by the REQ.

**Acceptance criteria (restated from the REQ):**
1. Placement measures real extents; the width model and measured-face constants are deleted, not adjusted.
2. The algorithm's behaviour is preserved and the kept rules are named in the code.
3. The remainder count is produced by whatever places the labels.
4. Measure before paint; the chart does not visibly reflow; the approach is named.
5. Every deleted test's property is re-pinned or its loss recorded with a reason.
6. Any surviving browser-measured number names its build — or none survives, said so.
7. Probes assert what a reader would see.
8. `maintainer-verify.sh` exits 0, including with no browser.

## Decisions

- **D-01** (DECIDE & STATE): Measurement happens inline in the render pass, with no offscreen clone and no hidden pass. Reasoning: the SVG is already in the document when labels are built, and the render is one synchronous task — the browser cannot paint an intermediate state it never gets a turn to paint. Naming it matters because a future reader might "fix" the apparent missing hide/reveal and add a reflow that is not there today.
- **D-02** (ESCALATE): `durationLabelTimeRange` and `durationLabelPlotX` were **moved into `durations_test.go`** rather than deleted with the rest. Reasoning: they are the x-axis domain (REQ-248's UTC-day anchoring), not the width model, and `TestJavaScriptBehaviorDurationsDayBucketsStayInsideThePlot` uses them as the second side of a real drift pin — a domain that moves in the renderer alone would push Panel B off canvas at one or two active days, which is the defect REQ-248 fixed. Keeping them in `durations.go` would have left unused production code; deleting them would have left the parity assertion with nothing to compare against. Value: the pin survives with no dead shipped code. Risk: a test-file definition is easier to let drift from a renderer change than a production one; the test that uses it fails in either direction, which is the mitigation.
- **D-03** (ESCALATE): **Three of the ten tests are recorded as deliberately dropped**, and the REQ requires that be justified rather than silent:
  - `TestDurationLabelGeometryMatchesTheRenderer` held Go's placement constants against the renderer's. **There is no second set of constants left** — placement lives only in the renderer, which is the REQ's point. A property with one side is not a property.
  - `TestDurationsLabelWidthEstimateCoversTheRenderedFace` pinned the width model above the face's supremum and below a slack ceiling. **There is no width model**; the width is the drawn width. Its subject is gone.
  - `TestDurationsMeasuredConstantsNameTheirChromiumBuild` required every hand-transcribed measured constant to name its build. **No such constant survives**, so there is nothing left to name one. This is REQ-266's requirement *met*, not skipped — see D-04.
  The other seven are re-pinned as browser probes, each named against the test it replaces in `durations_browser_probe_test.go`'s header.
- **D-04** (DECIDE & STATE): REQ-266's rule outlives the numbers it was written about, so it is enforced by a **guard rather than by provenance**: `TestDurationsCarriesNoMeasuredFaceConstants` fails if any retired measured-face token reappears in `durations.go` or `board-durations.js`. Reasoning: the REQ says if none survives, say so — and a statement that is only true on the day it ships is worth less than a check. It caught a real leftover on its first run (a comment in `board-durations.js` still citing the deleted `durationsMeasuredLabelBoxHeightUnits`), which is the case for having it.
- **D-05** (DECIDE & STATE): `TestJavaScriptBehaviorDurationsLabelsFollowTheShippedVerdict` was renamed to `...DurationsLabelRowsAndRemainders` and narrowed. Reasoning: half of what it asserted — "draw the payload's verdict, do not re-derive it" — has no payload verdict left to obey. The surviving half (row-to-baseline mapping on the sample's own band, out-of-range row is no label, zero remainder prints nothing while a nonzero one states the count) still catches the original defect it was written for.

## Implementation Summary

**What was done:** Moved Durations direct-label placement out of Go and into the renderer, where it sizes each label from the text the engine actually draws. Deleted the per-character width model and every hand-transcribed measured-face constant. The remainder count is now produced by the code that decides what fits. Ten Go tests lost their subject: seven are re-pinned as browser probes asserting what a reader would see, three are recorded as deliberately dropped, and a guard keeps the no-measured-constants property true going forward.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified) — `durationsLabelBandOf`, `durationsBandRowY`, `durationsLabelSpan`, `durationsSpanIsBlocked`, `placeDurationsLabelBand`, `packDurationsLabelBand`, `composeDurationsRemainderText`; the render pass now measures every label before positioning any; `durationsLabelBaselineY` takes the row it decided rather than reading one from the payload; the row-pitch comment no longer cites a deleted constant.
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified) — deleted `durationsLabelCharacterWidthUnits`, `durationsLabelGapUnits`, `durationsLabelSeparationUnits`, `durationsLabelRowCount`, `durationsLabelRemainderReserveUnits`, `durationsLabelRowUnplaced`, `durationLabelWidthUnits`, `durationLabelBandOf`, `planDurationLabels`, `packDurationLabelBand`, `durationLabelMagnitudeOrder`, `durationLabelSpanIsBlocked`, `placeDurationLabelBand`, `durationLabelSpan`, `DurationLabelBand`, `durationLabelInterval`, and the `LabelRow`/`LabelAnchor` sample fields.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — `generatedDurationLabels`, the `labels` payload field, and `labelRow`/`labelAnchor` on each sample removed.
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified) — the ten tests and their helpers removed; `durationLabelTimeRange`/`durationLabelPlotX` kept here as the axis-domain parity reference (D-02).
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go` (new) — seven probes plus the no-measured-constants guard.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — the payload-verdict probe rewritten and renamed (D-05).

**Tests touched:** ten deleted, eight added (seven browser probes plus the guard), one rewritten. Each deletion is accounted for in D-03 or by a named replacement.

## Qualification

Passed — 6 files verified, 8 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `gofmt -l .` clean, `go vet ./...` clean, `node --check` clean on the changed JS, `go test ./...` ok. No debug artifacts. `maintainer-verify.sh` exits 0 both with and without a browser.
- **Substantive:** the diff deletes a 130-line algorithm and its constants and adds a measured one; the real board renders 5 labels with 0 overlaps and a correct remainder.
- **Requirements traced:** AC1 → the constants are gone and the guard proves it; AC2 → each ported rule is named in the code with its REQ; AC3 → `packDurationsLabelBand` returns the count and the renderer draws that; AC4 → D-01, named in the code; AC5 → D-03 plus the probe file's header table; AC6 → D-04; AC7 → the probes assert drawn geometry; AC8 → both runs below.
- **Flowing:** placement reads real measured widths — proved by the face-response probe, which fails if the geometry is invariant under a letter-spacing change.

## Testing

- `go test ./...` — ok, 42.2s. `gofmt -l .` clean, `go vet ./...` clean, `node --check` clean.
- `bash _dev/tests/maintainer-verify.sh` — **exit 0 with a browser** (browser lane runs) **and exit 0 without one** (lane prints its SKIP), which is AC8 in both directions.

**Red-green validation** — the REQ's RED is a face wider than the model:

| | Before | After |
|---|---|---|
| A board whose face is wider than 7.15 units/char draws labels past their slots, and nothing notices | the width model returns the same number for any face, so the slots never move and the labels simply overprint | `TestBrowserBehaviorDurationsPlacementRespondsToTheRenderedFace` runs the same fixture at two letter-spacings and **fails if the geometry is identical** — which is exactly what a width model produces |

That probe is the RED turned into a standing assertion rather than a one-off measurement: a future regression to any fixed width model makes the two runs agree and fails it.

**The seven re-pinned properties, all passing in a real engine:**

| Probe | Replaces | Result |
|---|---|---|
| drawn labels never overlap, and none reaches outside the plot | `TestDenseOverflowLabelsStayBoundedAndNeverOverlap` | PASS (40-label saturated fixture) |
| the longest span always carries a label | `TestOverflowLabelsGoToTheLongestSpans` | PASS |
| clustered spans fill both rows | `TestClusteredOverflowLabelsFillBothLabelRows` | PASS |
| saturating one band changes nothing in the other | `TestReversedLabelPlacementIsIndependentOfOverflowDensity` | PASS |
| the remainder equals the number not drawn | (the count's own property) | PASS |
| row pitch clears the **measured** line box; rows clear the mark band; the last row clears Panel B's title | `TestDurationsLabelRowPitchClearsTheLabelTextBox`, `...RowsClearTheMarkBands`, `...LastLabelRowClearsPanelBTitle` | PASS |
| no measured-face constant reappears | `TestDurationsMeasuredConstantsNameTheirChromiumBuild` | PASS — **and it caught a real leftover on its first run** |

The probes run the **shipped** placement code, sliced out of the generated page rather than re-implemented, so they cannot drift into testing a copy.

**End-to-end on the real board**, driven in Chromium with the Durations view selected:

| Observation | Value |
|---|---|
| mark labels drawn | 5, across both rows |
| overlapping pairs | **0** |
| remainder sentence | `+8 more over 60 min`, clear of the nearest label on that row |
| page errors | none |
| measured label box | 12.9631 units, against the 13-unit pitch |

The remainder's clearance is the two-pass reserve working with a **measured** sentence width: the nearest row-1 label ends at 1038.8 and the reserved sentence begins near 1052.

## Review

**Overall: 91%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Placement measures real extents; width model and measured-face constants deleted not adjusted | ✅ and guarded against return |
| Behaviour preserved, not reinvented; kept rules named in the code | ✅ each rule carries its REQ |
| Remainder count produced by whatever places the labels | ✅ |
| Measure before paint; no visible reflow; approach named | ✅ D-01 |
| Every deleted test re-pinned, or its loss recorded with a reason | ✅ seven re-pinned, three recorded (D-03) |
| Any surviving measured number names its build; if none survives, say so | ✅ none survives; D-04 turns that into a check |
| Probes assert what a reader would see | ✅ drawn geometry, in a real engine |
| `maintainer-verify.sh` exits 0, including with no browser | ✅ both verified |

### Findings

**Important — none.**

**Minor:**

- **M1:** The browser probes drive placement through a synthetic SVG rather than through the real `renderDurationsView`. They run the shipped placement functions, sliced from the generated page, so the algorithm under test is the real one — but the wiring between it and the render pass is covered only by the end-to-end Chromium run recorded above, not by an assertion. Closing that would mean driving the view itself, which needs a real driver rather than `--dump-dom` (REQ-291's stated trade).
- **M2:** Three tests were dropped rather than replaced. Each is argued in D-03 and each genuinely has no subject left, but a reviewer counting tests will see the suite lose three assertions.
- **M3:** `model.go` was in the declared write set and needed no change; the label fields were in `durations.go` and `generate.go`. Recorded rather than quietly dropped.

**Nit:**

- **N1:** `runDurationsPlacementProbeWithFaceOverride` mutates a package-level variable and restores it with `defer`, which is not parallel-test-safe. No test in this package runs parallel today; a `t.Parallel()` added here later would need it threaded as a parameter instead.

### Restatement Sweep

Redefined elements: where placement happens, what the payload carries, and whether a label's width is modeled or measured. That is a lot of restatement surface, so the sweep was the longest part of this review.

- `board-durations.js`'s own header comment ("WHICH marks get one is decided in durations.go and arrives in the payload as labelRow/labelAnchor") — **rewritten**; it described the exact arrangement this REQ removed.
- `board-durations.js`'s `DURATIONS_LABEL_ROW_HEIGHT` comment cited `durationsMeasuredLabelBoxHeightUnits` and `TestDurationsLabelRowPitchClearsTheLabelTextBox`, both deleted — **rewritten**, and it was `TestDurationsCarriesNoMeasuredFaceConstants` that caught it rather than a human reading.
- `generate.go`'s `generatedDurationSample` doc said both fields "are decided in durations.go so the renderer reads one answer instead of becoming a second definition of the rule" — **rewritten** to say the opposite, which is now true.
- `generate_test.go`'s verdict probe doc — **rewritten** (D-05).
- `durations_test.go`'s remaining helpers — the ones describing the deleted model were removed with it; `durationsWindowStart` and the fixture builders survive and are still accurate.
- Root `CLAUDE.md` § Kanban Board Write Surfaces — **checked: still exactly three.** This REQ adds no write surface.
- `_dev/primes/prime-kanban-board.md` — grepped for statements about where placement runs or what the payload carries; it defers to the tool and to the action files rather than restating either, so nothing there went stale.

No stale restatement remains. Two of the four found were in comments no test could have caught, which is why this sweep was read rather than grepped alone.

### Acceptance Testing

Three levels, because a port like this can pass its own unit probes and still render wrongly. The probes exercise the shipped placement functions against measured widths in a real engine. The face-response probe proves the measurement actually feeds placement, by requiring the geometry to *change* when the drawn text widens. And the whole board was generated and driven in Chromium with the Durations view selected: five labels, zero overlaps, a correct remainder clear of its neighbour, no page errors, and a screenshot read by eye.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope Discipline | 95% |
| Risk | Low |
| Acceptance | Pass |

Test Adequacy 90% for M1 and M2. Risk Low rather than Medium: the change is large and touches the most visually complex view, but the property it exists to fix is now asserted against the real engine rather than against a constant, and the end-to-end render was inspected rather than assumed. The failure mode it removes — labels overprinting on any face wider than one container's — was previously undetectable by construction.

### Follow-up REQs Created

None from this REQ's own work. One queue observation for the hand-back: **REQ-266's body says it was "cancelled as superseded on 2026-08-20" and folded in here, but REQ-266 is still `pending` in `do-work/queue/` with `depends_on: [REQ-292]`.** Its requirement is met by D-04. Flagged rather than acted on, because cancelling another REQ is `do-work abandon`'s job and the user's decision.

## Lessons Learned

**What worked:** Deleting in an order that kept the build compiling — sample fields, then the functions that used them, then the constants, then the payload — so every step's breakage named exactly what still referenced the thing being removed. `go vet` became the worklist. And running the real board in a browser at the end: the probes were all green before that, and only the rendered chart shows that the remainder sentence actually clears its neighbour on the row.

**What didn't:** A regex-driven "drop this function and its doc comment" helper, used to delete sixteen declarations, repeatedly clipped one character off the following declaration — leaving `f// comment`, `ffunc`, and a bare `/` at end of file. Every one was caught by the compiler, so nothing shipped wrong, but it cost several rounds. For a bulk deletion of Go declarations, deleting by explicit start/end line after listing them beats a comment-walking heuristic; the heuristic's failure mode is silent damage to the *neighbour*, not to the target.

Also: the first port measured the remainder reserve from a string that did not match the one drawn ("with a reversed span" against "reversed"). Caught by reading, not by a test, and it would have under-reserved by ~100 units. When two call sites must compose the same sentence, the composer belongs in one function — which is what `composeDurationsRemainderText` now is.

**Worth knowing:** The reason this defect was undetectable, not merely undetected. A width model multiplies a character count by a constant, so it returns the *same* number for every face — which means the slots it assigns never move, and a wider face draws past them silently. No amount of pinning the constant could catch that, because the constant was never wrong in its own terms; it was answering a question about a face nobody had. That is why the replacement test asserts the geometry *changes* under a wider face rather than asserting any particular width: the property that matters is responsiveness, not accuracy.

## Orientation

The Durations view now decides where its labels go in the browser, sizing each one from the text the engine actually draws. The per-character width model and every hand-transcribed measured-face constant are gone — the numbers that described one Linux container's answer to `--font-sans` while the stylesheet ended in the open `sans-serif` generic, so they described a face they never met. Labels can no longer overprint on a face wider than the model assumed, because there is no model to be wider than. Lives in the board's Durations subsystem (`skills/do-work-board/tools/queue-kanban/web/board-durations.js`, with the Go side reduced accordingly), indexed by `_dev/primes/prime-kanban-board.md`.

**[MAP CHANGED]** — the board payload no longer carries `labelRow`, `labelAnchor`, or `labels`; anything reading those is reading fields that no longer exist. Placement is renderer-side, so `durations.go` no longer answers "which marks get a label" at all. The board tool's write-surface count is unchanged at three.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` and `_dev/primes/prime-shell-commands.md` — referenced paths still resolve. The board prime's parser-lock-step rule is untouched (no frontmatter field involved); neither prime restates where label placement runs, so neither was made stale.
