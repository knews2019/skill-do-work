---
id: REQ-292
title: Move Durations label placement into the browser and delete the measured-face constants
status: pending
created_at: 2026-08-19T14:36:44Z
status_changed_at: 2026-08-19T14:36:44Z
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
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/model.go
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
