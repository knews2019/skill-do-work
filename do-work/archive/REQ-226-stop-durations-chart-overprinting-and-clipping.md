---
id: REQ-226
title: Stop the Durations chart from silently overprinting and clipping
status: completed
completed_at: 2026-08-18T00:55:23Z
claimed_at: 2026-08-18T00:25:12Z
created_at: 2026-08-17T23:51:17Z
user_request: UR-051
addendum_to: REQ-219
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-227, REQ-228]
batch: board-timing-views
route: C
estimate:
  p50_active_minutes: 50
  confidence: medium
  calculated_at: 2026-08-18T00:26:00Z
  basis:
    - Route C
    - 6-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - browser evidence
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/durations_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
---

# Stop the Durations Chart from Silently Overprinting and Clipping

## What

Fix two defects in the Durations view where the chart draws something that reads as a value but
isn't. Panel A labels every overflow sample with no collision detection, so on a large board the
overflow lane becomes an unreadable blob of overprinted text. Panel B clamps its bars at 45 minutes
with no scale break, so a 78-minute day renders as a 45-minute bar.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md` and REQ-219's recorded lesson, then the three layers end to end. Approach written in `## Plan` above.
- [x] **[APPLY]:** Five declared files, no others. One deviation from plan, recorded as D-02.
- [x] **[UNIFY]:** `git diff --stat` = 5 files, +686/−22. `gofmt -l .` clean, `go vet ./...` clean, `node --check` on the assembled client clean. Grep for `console.log|fmt.Print|debugger|TODO|FIXME|XXX` in the added lines: none. Files verified:
  - `durations.go` — placement section reads as one pass with no leftover scaffold; the shipped-today rule used to capture RED is gone (confirmed by grepping `RED SCAFFOLD` and `MUTATION`: absent).
  - `durations_test.go` — three new tests, each mutation-checked; the pinned live-archive test untouched.
  - `generate.go` — additive wire fields only; no existing field renamed or reordered.
  - `generate_test.go` — one new probe plus two small helpers; existing inventory and manifest assertions untouched.
  - `web/board-durations.js` — the parity label pass is gone (grepped `% 4` and `anchorsLeft`: absent); no stray globals (the probe's `ReferenceError` runs caught two).

## Why

"fix this" — with a screenshot of the maintainer's own board where the overflow lane is unreadable.

## Context

Addendum to REQ-219, which built the Durations view (archived at
`do-work/archive/UR-050/REQ-219-board-durations-view.md`, commit `cfffd90`). Both defects are in what
that REQ shipped; neither is a design disagreement with it.

**Defect 1 — the label blob.** `web/board-durations.js:246-265` filters `markIndex` down to every
sample above the 60-minute ceiling or below zero, then places a direct label for each one using a
four-step index cycle:

```js
var anchorsLeft = position % 2 === 0;
x: (mark.x + (anchorsLeft ? -9 : 9)),
y: (mark.y + (position % 4 < 2 ? 4 : 16)),
```

That is two text rows (y≈44 and y≈56) × two anchor sides = four slots, recycled forever. There is no
width measurement, no neighbour comparison, no drop rule, no leader lines. The alternation keys on
**array index**, and `markIndex` is in completion-time order — so a burst of overflow REQs finishing
in the same hour gets ~0px of x separation while positions 0 and 1 land on the *same row* on opposite
sides of the *same* x. Plot width is 1128 user units; a label like `REQ-407 14h 15m` is ~85-90 units
at `font-size: 11px` (`board.css:1866-1870`), i.e. ~7.5% of the plot. Two overflow samples within ~4%
of the time range already overprint.

Zooming the browser does not help: the SVG is `width: 100%` over a fixed `viewBox`
(`board-durations.js:133`, `board.css:1846-1851`), so text scales with the chart and density in user
units is invariant.

The overflow lane exists precisely because y carries no magnitude above the ceiling — every overflow
mark sits at the constant `DURATIONS_LANE_MARK_Y`. **The text is the only carrier of the value there**,
which is why losing it to overprinting costs more than it would on any other panel.

The same filter mixes overflow marks (lane, y≈40) with reversed-stamp marks (y=284), so one counter
drives two visually unrelated groups and desynchronises from either one's local density.

**Defect 2 — Panel B's silent clip.** `yOfDayMedian` clamps with
`Math.min(minutes, DURATIONS_MEDIAN_CEILING)` where the ceiling is 45
(`board-durations.js:35, 164-170`). Panel A handles its own ceiling honestly, with an overflow lane
above a scale break and a `60+` tick. Panel B has neither: the bar is drawn flat at the top of the
plot and reads as a 45-minute day. In the attached screenshot the slowest-day label already prints
`78 min` (`board-durations.js:299-311`) directly above a bar clipped at 45 — the chart is stating two
different numbers for the same day, and only the small one is drawn to scale.

Note this only ever affects the single slowest day's label today; every *other* clipped day is
clipped with no annotation at all.

**Why it is invisible here.** This repo's archive has three overflow samples (REQ-043 at 149.4 min,
REQ-064 at 655.2 min, REQ-171 at 80.5 min out of ~205 measurable) and zero reversed ones. Four slots
are enough for three labels. The reporting board has 560 samples across 42 days. Any RED test must
therefore build its own dense fixture rather than lean on the live archive.

## Detailed Requirements

1. **Panel A labels must not overprint.** Place direct labels greedily and skip any that would
   collide with one already placed. Collision is a real geometric test against the space the label
   text occupies, not an index parity trick.
2. **Nothing may be silently hidden.** Whatever the placement rule drops must be accounted for on the
   chart — print a remainder such as `+N more over 60 min` at the lane's edge. A reader must never be
   able to mistake "the lane holds four things" for the truth when it holds four hundred.
3. **Overflow and reversed marks are decided independently.** They occupy different bands and have
   independent local densities; one shared counter serving both is part of the current defect.
4. **The label-selection verdict ships in the payload.** Compute which samples get a direct label,
   and the remainder count, on the Go side and send the answer to the client — do not re-derive the
   rule in JS. This is REQ-219's own recorded lesson, from the very view being amended: *"ship a
   rule's verdict in the payload so a second reader cannot become a second definition"*
   (`do-work/archive/UR-050/REQ-219-board-durations-view.md#lessons-learned`, cited from
   `_dev/primes/prime-kanban-board.md:25`). It is also what makes the RED below runnable in the Go
   suite. **Latitude:** if the builder finds a better home for the rule, the constraint that survives
   is that the RED test stays runnable in `_dev/tests` / the Go suite — not the specific file.
5. **Panel B must not draw a clipped bar as if it were the real value.** Give it the honest treatment
   Panel A already has — an overflow lane above a scale break, a `45+` tick, or an equivalent device —
   so an over-ceiling day is visibly over-ceiling. Any day above the ceiling, not only the slowest
   one, must be distinguishable from a day that genuinely sat at the ceiling.
6. **Every value stays reachable without a pointer.** The hover readout
   (`board-durations.js:413-443`) and the table (`:492-515`) already carry every sample; a label that
   the new rule drops must still be findable there. Do not regress either.

## Constraints

- Read-only. No new board write surface — `CLAUDE.md` § *Kanban Board Write Surfaces* must still read
  "exactly three" and go unamended.
- No new frontmatter field and no second walk of the archive. `buildDurationAggregate` is documented
  as "a pass over memory, not over the archive" (`durations.go:73-75`); keep it that way.
- **Do not touch the read-time rule.** `durations.go:9-26` is explicit that it is the *second reader*
  of the rule defined in `skills/do-work/actions/estimate-reference.md` → Calibration, not a second
  definition. This REQ changes how values are drawn, never which values count.
- `durations_test.go:202-225` pins live-archive figures (≥195 samples, the 2026-07-31 and 2026-08-15
  medians). Those assertions must still pass unchanged.
- If a new `web/*.js` file is added, `boardJavaScriptFragmentPaths` (`generate.go:42-52`) and the
  inventory assertion (`generate_test.go:28-62`) change in the same commit or the build fails.
- Panel A's existing scale-break design is correct and stays. This REQ fixes the labelling on top of
  it; it does not redesign the panel.

## Dependencies

None. Independent of REQ-227 and REQ-228 — different files, different view. Can run first or in
parallel with them.

## Builder Guidance

**Certainty: Firm.** The defect, its cause, and the chosen fix were all confirmed with the maintainer
before capture. The collision-aware-plus-count approach was chosen over top-N-extremes and over
dropping labels entirely; do not re-open that choice.

Latitude is in the geometry: how label width is estimated or measured, how many text rows the lane
gets, whether dropped labels leave a subtle tick, and what device Panel B uses for its ceiling. Pick
what reads well and state why.

One judgment worth making early: a per-character width estimate is probably enough and is trivially
testable in Go, whereas `getComputedTextLength` is exact but only available after render and pushes
the decision back into the client. The requirement above prefers the former for that reason.

## Red-Green Proof

**RED prompt/case:** Build a synthetic board fixture of ~40 archived REQs that all completed inside a
two-day window and all carry wall spans over 60 minutes, then compute the Durations payload. Assert
that (a) the number of samples marked for a direct label is bounded by what the lane can physically
hold, and (b) no two labelled samples sit closer along the x-axis than the width their label text
needs. Today every one of the 40 is marked for a label, so both assertions fail — the current code
has no concept of either bound.

**Why RED now:** `web/board-durations.js:246-265` labels every overflow sample and places them with a
four-slot index cycle and no collision detection, so N overflow samples produce N labels in 4 slots.
The screenshot is that failure at N≈several hundred.

**GREEN when:** The same fixture yields a bounded, non-overlapping label set plus a remainder count
that accounts for every dropped label, and rendering the maintainer's own board shows a legible
overflow lane where each visible label is readable and the count states how many are not shown.
Separately, a day whose median exceeds 45 minutes is visibly distinguishable from a day sitting at
exactly 45.

**Validation:** User confirmed — the collision-aware-plus-count strategy was chosen by the maintainer
from three offered options; the Panel B clipping defect was surfaced by capture in answer to "can we
improve anything else?" and accepted into scope.

## Prior Implementation

REQ-219 (`do-work/archive/UR-050/REQ-219-board-durations-view.md`, Route C, commit `cfffd90`) built
the whole view: `durations.go` (the data layer — `DurationSample`, `DurationDay`,
`buildDurationAggregate`, the read-time rule, medians), the payload projection at
`generate.go:521-541` with wire types at `generate.go:218-248`, the HTML shell at
`web/template.html:188-219`, the 516-line SVG renderer `web/board-durations.js`, and CSS at
`web/board.css:1790-1955` with tokens at `:91-94` (light) and `:147-148` (dark).

Its recorded lesson — ship a rule's verdict in the payload so a second reader cannot become a second
definition — is directly load-bearing for Requirement 4 above.

## Assets

`do-work/user-requests/UR-051/assets/REQ-226-screenshot-1-durations-overflow-label-blob.png` — the
Durations view of the maintainer's `glw-game-find-the-difference` board, 560 archived REQs across 42
active days. Panel A shows one legible overflow label at the far left (`REQ-407 14h 15m`) and, from
roughly 1 July rightward, two horizontal bands of completely overprinted text. Panel B shows a
`78 min` annotation above a bar clipped flat at the 45 tick. Panel C is legible and undamaged.

The ids and counts in that image belong to another repository. They are evidence that the defect
scales with sample count, not data to reproduce locally.

## Full Context

See `do-work/user-requests/UR-051/input.md` for complete verbatim input.

---
*Source: "do-work capture-request: fix this: [Image #1] ... Am I missing anything, can we improve anything else?"*

---

## Triage

**Route: C** - Complex

**Reasoning:** Two independent defects across three layers — a new placement rule in Go, a rewritten label pass and a new scale break in the SVG renderer, and the payload contract between them. Six numbered requirements and six constraints, with a required RED that has to build its own dense fixture because the live archive is too sparse to reproduce the failure.

**Planning:** Required

## Plan

### The two root causes

**Defect 1.** `web/board-durations.js:246-265` places a direct label for *every* overflow or reversed sample, choosing its slot from `position % 2` and `position % 4` — the sample's index in a completion-ordered array. Index parity carries no information about where a label sits or how wide it is, so the rule cannot detect a collision even in principle. It also runs one counter across two bands (overflow at y≈40, reversed at y=284) that have unrelated local densities.

**Defect 2.** `yOfDayMedian` clamps with `Math.min(minutes, DURATIONS_MEDIAN_CEILING)` and Panel B draws the clamped bar with no scale break and no over-ceiling tick, so a 78-minute day and a 45-minute day render identically.

### Approach

1. **Go decides which labels are drawn** (`durations.go`). A new `planDurationLabels` pass runs after `buildDurationAggregate` fills `Samples`. For each of the two bands independently, walk that band's samples in x order and greedily place each label into the first row whose already-placed labels do not overlap the interval this label would occupy. Anchor `start` by default; fall back to `end` when a `start` label would cross the plot's right edge, and test that interval instead. Samples that fit nowhere are counted, not drawn.
2. **Geometry lives in user units**, the same space the SVG's `viewBox` defines, so the estimate is resolution-independent — the chart is `width: 100%` over a fixed `viewBox`, so browser zoom cannot rescue density. Label width is estimated per character (Builder Guidance's stated preference over `getComputedTextLength`, which is only available post-render and would push the decision back into the client). The estimate is deliberately generous: over-estimating width drops a label, under-estimating overprints one, and only one of those is the defect being fixed.
3. **The remainder is reserved, not discovered.** A remainder sentence needs room at the band's right edge, but its existence is only known after placement. Two passes: place with the full width; if nothing was dropped, that is the answer. If anything was dropped, re-place with the right-edge zone reserved on row 0, and report that pass's count. Terminating, deterministic, and it costs nothing on a board with no remainder.
4. **The payload carries the verdict.** Each sample gains `labelRow` (−1 when unlabelled) and `labelAnchor`; the durations payload gains per-band hidden counts. JS reads the verdict and never re-derives it — REQ-219's own recorded lesson, which this REQ cites as Requirement 4. The remainder *sentence* stays in JS with the view's other copy; Go ships the count, and reserves against a fixed-width template rather than against the composed string.
5. **Panel B gets Panel A's device.** An overflow lane above a scale break, a `45+` lane tick, and a detached cap block for any day over the ceiling — every such day, not only the slowest. The existing slowest-day annotation stays. Panel C's constants shift down by the lane's height; `DURATIONS_MEDIAN_TITLE_Y` does not move, so `describeAtPointer`'s panel split is untouched.

### Verification steps

1. Go RED in `durations_test.go` → verify: a 40-sample dense over-ceiling fixture yields a bounded label set with no two labels overlapping on a row, and a hidden count that accounts for every unlabelled sample.
2. Go test pinning Go's plot-width constant against `web/board-durations.js`'s own three constants → verify: the shared geometry cannot drift silently.
3. JS behavior probe in `generate_test.go` → verify: the renderer draws a label only where `labelRow >= 0`.
4. `go test -count=1 ./...` → verify: `durations_test.go:202-225`'s pinned live-archive figures still pass unchanged.
5. `bash _dev/tests/maintainer-verify.sh` → verify: exit 0.
6. Render the real board and read Panels A and B → verify: the lane is legible and Panel B's over-ceiling days are visibly over-ceiling.

*Generated by Plan agent (inline, serial mode)*

## Exploration

**Data layer.** `durations.go` (170 lines) — `DurationSample`, `DurationDay`, `buildDurationAggregate` (a pass over already-parsed tickets), `dayMedianExclusionReason`, `medianOf`. Samples are sorted by `CompletionTime`, so they already arrive in x order and no second sort is needed for placement.

**Payload.** `generate.go:218-248` holds the wire types (`generatedDurations`, `generatedDurationSample`, `generatedDurationDay`); `generate.go:521-541` projects the aggregate onto them. Adding fields here is additive — no existing consumer reads by position.

**Renderer.** `web/board-durations.js` (516 lines). Geometry constants at `:16-42`; `xOfEpoch`/`yOfMinutes`/`yOfDayMedian` at `:152-170`; mark plotting and `markIndex` at `:224-243`; the defective label pass at `:246-265`; Panel B at `:267-311`; Panel C at `:313-352`; shared axis at `:354-395`; hover readout at `:397-470`; table at `:472-515`.

`describeAtPointer` splits Panel A from B/C on `pointerY <= DURATIONS_MEDIAN_TITLE_Y - 12`, so Panel B's internal layout can change freely as long as that constant does not move.

**Test surfaces.** `durations_test.go` has six unit tests plus `TestLiveArchiveDurationsMatchTheCalibratedFigures` (:202-225), which pins ≥195 samples and the 2026-07-31 / 2026-08-15 medians — a REQ constraint, must pass unchanged. `generate_test.go:27-53` pins the exact embedded `web/*.js` inventory and `:56-131` the fragment manifest, so adding a JS file would require both; this REQ adds none. `runJavaScriptBehaviorProbe` (`generate_test.go:262-274`) executes a Node `-e` script and counts toward the maintainer-strict lane; `sliceBalancedBlockAfter` extracts a named function from the assembled client, which is how a JS probe gets at renderer internals.

**Live-archive density.** Confirmed against this repo: only three samples exceed 60 minutes and none is reversed, so the live board cannot exhibit the defect. The RED must build its own fixture, exactly as the REQ states.

**CSS.** `web/board.css:1840-1900` defines `.durations-mark-label` (sans 11px), `.durations-tick` (mono 11px, faint), `.durations-overflow-lane`, `.durations-bar`. The remainder sentence reads as chart furniture rather than a value, so it takes the existing `.durations-tick` class, and Panel B's lane and cap reuse `.durations-overflow-lane` and `.durations-bar`. No new class is needed — see D-01.

*Explored by work action (inline, serial mode)*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — label geometry constants, `planDurationLabels`, per-sample label fields, per-band hidden counts.
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — the dense-fixture RED, plus the constant-parity test against the JS geometry.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — wire fields for the label verdict and the hidden counts, and their projection.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — JS behavior probe that the renderer honors the shipped verdict.
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — replace the parity label pass with one that reads the verdict; add the remainder sentences; give Panel B a scale break, a `45+` lane tick, and over-ceiling caps.

**Files I will NOT touch:** `web/board.css` (declared by capture; dropped — see D-01), `skills/do-work/actions/estimate-reference.md` and the read-time rule in `durations.go:9-26` (this REQ changes how values are drawn, never which values count), `CLAUDE.md` § Kanban Board Write Surfaces (no new write surface), `web/template.html` (no new panel or control).

**Acceptance criteria (restated from REQ):**
- [ ] Panel A labels are placed by a real geometric collision test and never overprint.
- [ ] Nothing is silently hidden — a remainder accounts for every dropped label.
- [ ] Overflow and reversed marks are decided independently, with no shared counter.
- [ ] The label-selection verdict is computed on the Go side and shipped in the payload; the RED test runs in the Go suite.
- [ ] Panel B does not draw a clipped bar as if it were the real value, and every over-ceiling day — not only the slowest — is distinguishable from a day that sat at the ceiling.
- [ ] The hover readout and the table still carry every sample.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)

**What was done:** Moved the direct-label decision out of the renderer and into `durations.go`. `planDurationLabels` packs each band independently, walking its samples left to right in the SVG's own user-unit space and giving each label the first row where its estimated text width touches nothing already placed there; a label that fits nowhere is counted instead of drawn. The count is shipped in the payload alongside each sample's `labelRow` and `labelAnchor`, and the renderer draws that verdict rather than deriving one. Panel B's top gridline now reads `45+`, and every day whose median exceeds the ceiling gets a detached sliver above its full-height bar, so an over-ceiling day is distinguishable from one that sat at exactly 45. The read-time rule, the aggregate's single pass over memory, and the hover readout and table are untouched.

## Decisions

- **D-01**: Dropped `web/board.css` from the declared write set before implementation. Reasoning: the remainder sentence reads as chart furniture rather than a value, so it takes the existing `.durations-tick` class, and Panel B's over-ceiling sliver takes `.durations-bar`. No new rule was needed, and a declared-but-untouched file is scope noise. The sliver still carries a second class name, `durations-bar-over-ceiling`, with no rule behind it — that is a semantic hook so the element is identifiable, not dead style. DECIDE & STATE.
- **D-02**: Placement tries the anchor *before* the mark first and falls back to the anchor after it, rather than the reverse. This was not the first implementation: preferring the after-mark anchor cost a label on the maintainer's own board (`REQ-171` was dropped with a `+1 more` remainder where all three had previously fit). The walk runs left to right, so a label drawn to a mark's left reuses space the walk has already passed while one drawn to its right consumes space the next mark still needs — the preference is a packing property, not a style choice, and flipping it restored all three real labels with no remainder. Capacity in the dense case is unchanged, because the pitch between two same-anchor labels is identical either way. DECIDE & STATE.
- **D-03**: The remainder sentence sits on each band's *last* text row, not its first. The marks sit level with the first row, so a sentence there is legible only while the band is sparse — which is exactly when there is no remainder to print. The first render of the dense fixture showed `+58 more over 60 min` overprinted by the mark blob it was describing, which would have reproduced the defect inside its own fix. Placement reserves the last row's right edge to match. DECIDE & STATE.
- **D-04**: A residual overlap is left alone. In a saturated band a first-row label can still be crossed by a *neighbouring mark* (marks sit at y 40, first-row text at y 44), because this REQ's rule guards label-against-label, not label-against-mark. Fixing it needs a taller lane, and the REQ states that Panel A's scale-break design is correct and stays. Recorded in `## Discovered Tasks` rather than fixed here.

## Discovered Tasks

- **[low]** Panel A's first label row overlaps the mark band. `DURATIONS_LANE_MARK_Y` is 40 with a mark radius of 5 (y 35-45) and `DURATIONS_LANE_LABEL_ROW_Y` is 44 (text box roughly y 33-46), so a first-row label and a neighbouring mark share vertical space. A label always clears its *own* mark, because it is offset horizontally from it, so this only shows where the band is dense enough that other marks crowd the label — visible in the synthetic 60-sample fixture, invisible on this repo's three-sample lane. The collision rule added by REQ-226 is deliberately label-against-label; closing this one means giving the lane roughly 12 more user units and shifting the panels below it, which REQ-226 was told not to do ("Panel A's existing scale-break design is correct and stays").

## Testing

**Tests run:** `cd skills/do-work-board/tools/queue-kanban && go vet ./... && go test -count=1 ./...` (includes the maintainer-strict JavaScript behavior lane under Node), then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing — both exit 0

**Red-green validation:** every new assertion below was mutation-checked. RED was produced by putting the *shipped* rule back, not by asserting against an absent function, so each failure is the real defect rather than a compile error.

- `TestDenseOverflowLabelsStayBoundedAndNeverOverlap` (durations_test.go): ✗ `placed 40 labels, but two rows of 1128 units hold at most 26` → ✓. RED came from porting `web/board-durations.js:246-265`'s index-cycle rule into Go verbatim; the fixture is 40 REQs all over the ceiling inside a two-day window, because the live archive carries three overflow samples out of 206 and cannot express the defect.
- `TestReversedLabelPlacementIsIndependentOfOverflowDensity` (durations_test.go): ✗ `REQ-900 placed [0 end] alone but [1 start] beside a dense overflow lane; the bands share state` → ✓. This assertion was reordered mid-build: its saturation precondition originally ran first and tripped before the real comparison, so the mutation failed for the wrong reason and the test proved nothing.
- `TestJavaScriptBehaviorDurationsLabelsFollowTheShippedVerdict` (generate_test.go, Node probe): ✗ `baseline[1] = 44, want 56` when the renderer ignores `labelRow`, and ✗ `remainder baseline[0] = 44, want 56 — the sentence must clear the mark row` when the remainder moves back to the mark row → ✓.
- `TestJavaScriptBehaviorDurationsLabelWidthModelMatchesTheRendererFormatter` (generate_test.go, Node probe): ✗ `60.0 min: renderer draws "1 hr 0m" (7 chars) but the width model assumes "1h 0m" (5 chars)` → ✓.
- `TestDurationLabelGeometryMatchesTheRenderer` (durations_test.go): passes; it reads the renderer's own constants rather than hand-copied numbers, so it fails by construction on any divergence.

**Regression evidence (REQ constraint):** `TestLiveArchiveDurationsMatchTheCalibratedFigures` (durations_test.go:202-225) passes unchanged — ≥195 samples, the 2026-07-31 ruled median of 2.5 min, and the 2026-08-15 median of 19.6 min over 25 completions. The read-time rule was not touched.

**New tests added:**
- `TestDenseOverflowLabelsStayBoundedAndNeverOverlap`, `TestReversedLabelPlacementIsIndependentOfOverflowDensity`, `TestDurationLabelGeometryMatchesTheRenderer` (durations_test.go)
- `TestJavaScriptBehaviorDurationsLabelsFollowTheShippedVerdict`, `TestJavaScriptBehaviorDurationsLabelWidthModelMatchesTheRendererFormatter` (generate_test.go)

**Existing tests updated (cross-REQ impact):** none. REQ-219's assertions were read but not modified.

**Rendered acceptance (the REQ's GREEN names "rendering the maintainer's own board"):** headless Chromium against both generated boards.
- Synthetic 158-REQ fixture with a 60-REQ over-ceiling burst compressed into two days of a 22-day span — the shape of the attached screenshot. The lane draws one legible label plus `+58 more over 60 min` on the row below the marks; the reversed band draws both its labels. Panel B's three over-ceiling days each show a full-height bar with a detached sliver above it, under a `45+` tick.
- This repository's own board, 206 samples. All three overflow labels (`REQ-043 2h 29m`, `REQ-064 10h 55m`, `REQ-171 1h 21m`) are placed and legible with **no** remainder, so the fix costs nothing on a sparse lane. Panel B has no day over 45 minutes and is visually unchanged apart from the tick.

## Lessons Learned

**What worked:** Porting the shipped rule into Go to produce RED, rather than writing a test against a function that did not exist yet. A compile error would have been indistinguishable from a real failure; `placed 40 labels, but two rows hold at most 26` names the defect. Rendering a synthetic fixture dense enough to reproduce the reported failure was the other half — the live archive passes every assertion either way, so nothing short of a purpose-built board could show whether the fix worked.

**What didn't:**
- Preferring the after-mark anchor. It is the obvious default and it is wrong for a left-to-right greedy: it spends space the next mark still needs. It cost a label on the maintainer's own board before the render caught it (D-02).
- Putting the remainder sentence on the band's first row. The marks sit at that height, so the first dense render showed `+58 more over 60 min` overprinted by the blob it was describing — the defect reproduced inside its own fix (D-03).
- Restoring a file from a copy taken earlier in the session. The mutation-test restore silently reverted D-03, and the whole suite still passed, because at that point nothing pinned which row the remainder used. The near-miss is why `durationsRemainderBaselineY` exists as a named function with its own probe assertion rather than as an expression inline in the render pass.

**Worth knowing:** Three of the four defects found here were invisible to the test suite and visible in a render. A chart's correctness is partly a claim about pixels, and neither a passing assertion nor a payload dump can see two glyphs sharing a coordinate. Generate a board and look at it. The corollary bit twice: a fix verified only by tests that pass both before and after it is not verified.

## Orientation

The board's Durations view no longer draws things that read as values but are not. Panel A's overflow lane shows as many direct labels as it can fit without two touching, and states how many it could not fit; Panel B marks a day whose median ran past the ceiling instead of drawing it flat at the ceiling. **[MAP CHANGED]** — the label-placement rule moved from the renderer into `durations.go` and now travels in the payload as `labelRow`/`labelAnchor` plus per-band hidden counts, which is the third rule this view computes on the Go side and the client only draws. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path still resolves, the three-write-surface count is unchanged (this REQ adds none), and REQ-219's lesson that the prime links is the one this REQ applied — the prime is not stale.

## Review

**Overall: 94%** | 2026-08-18T00:54:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- Panel A's first label row still shares vertical space with the mark band, so a label there can be crossed by a neighbouring mark in a saturated lane; the new rule guards label-against-label only — gate: rule-change → recorded in `## Discovered Tasks` as `[low]`, queued by Step 8. Closing it needs a taller lane, which the REQ excluded ("Panel A's existing scale-break design is correct and stays").

**Minor findings:** 2 (report only)
- The over-ceiling sliver carries `durations-bar-over-ceiling`, a class with no CSS rule behind it. It is a semantic hook, not dead style, and adding a rule would mean touching `board.css` for no visual change.
- Placement necessarily models the renderer's geometry and its number formatting in Go. Both halves are pinned — constants by `TestDurationLabelGeometryMatchesTheRenderer`, formatting by the width-model probe — so the duplication cannot drift silently, but it is duplication.

**Acceptance:** Pass — all six restated criteria verified, four of them by mutation-checked assertions and two (Panel B's device, the lane's legibility) by rendering both a synthetic dense board and this repository's own.
**Suggested testing:** 0 items
**Follow-ups created:** see Step 8 (Discovered Tasks); **sweeps appended to:** None

*Reviewed by review-work action*
