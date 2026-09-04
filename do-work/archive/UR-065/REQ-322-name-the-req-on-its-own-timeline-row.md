---
id: REQ-322
title: "Name the REQ on its own timeline row"
status: completed
created_at: 2026-08-22T22:08:34Z
claimed_at: 2026-08-23T02:56:30Z
route: B
completed_at: 2026-08-23T03:37:18Z
commit: 1c42897
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-321]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-318, REQ-319, REQ-320, REQ-321, REQ-323, REQ-324]
batch: timeline-ux-audit
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-23T02:57:10Z
  basis:
    - Route B
    - 3-file write set
    - 5 acceptance criteria
    - dependency depth 4
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
kb_status: promoted
kb_entry: REQ-322-name-the-req-on-its-own-timeline-row.md
---

# Name the REQ on Its Own Timeline Row

## What

Show each REQ's title in the row's label column, and put its detail in a tooltip at the
pointer instead of only in a readout at the foot of the panel.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** The label face is monospace, so ONE measured advance answers every row — no
  per-row, per-frame measurement, and no width model. Measure it against the live DOM once
  per render with `getComputedTextLength`, verify the monospace assumption rather than
  assuming it, derive a character budget, and truncate. Widen the column to 184px. The
  tooltip is a native SVG `<title>` on the row group: no listener, no positioning code.
- [x] **[APPLY]:** Confined to the declared write set (`board.css` was declared and not
  needed — D-03). One defensive guard added after the suite caught a real gap (D-02).
- [x] **[UNIFY]:** `git diff --stat` reviewed; `node --check` clean; `go vet ./...` clean;
  `gofmt -l .` empty; no debug artifacts. Files verified — `web/board-timeline.js` (three new
  functions, the widened constant, the label and `<title>` calls, the host guard),
  `generate_test.go` (one probe added), `timeline_browser_probe_test.go` (two tests added).

## Why

The 104px label column holds `REQ-012` and nothing else. The title lives in a one-line
readout at the very bottom of the panel — below a five-sentence hint paragraph, some 700px
from the pointer — and in the collapsed table. Scanning three hundred rows tells the reader
nothing about what any of them are.

## Context

`renderVisibleRows` draws one `<text class="timeline-row-label">` per row holding `row.id`.
`timelineRowDescription` already composes the full sentence (id, title, route, status, both
spans, projection slot, anomaly reason) and is used for both the row's `aria-label` and the
foot readout. The pieces exist; they are just not where the eye is.

`TIMELINE_LABEL_WIDTH = 104` is subtracted from the plot width in `plotWidth()`, so widening
the label column narrows the plot. That trade is the substance of this REQ: enough width to
recognize a title, not so much that the chart loses its span.

## Detailed Requirements

- Each row shows its title beside its id, truncated to fit the label column.
- Truncate against a **measured** face, not an estimated character width. This module has
  been bitten twice: `REQ-292` (a width model returns the same number for every face, so
  slots never move and a wider face draws past them) and `REQ-241`/`REQ-242` (the same 12px
  face measured 12.0372 and 11.2300 on two Chromium builds). Record the browser and build
  beside any measured number, and take the larger where two disagree.
- Hovering a bar shows a tooltip at the pointer carrying what the foot readout carries.
  A native SVG `<title>` on the row group is the cheapest thing that works and needs no
  positioning code; a custom tooltip is fine if it earns the complexity.
- The foot readout stays as the `aria-live` region — it is the non-pointer path and it is
  not this REQ's to remove.
- A REQ with no title (nothing in `requestsById`) still renders its id and does not draw an
  empty tooltip.

## Constraints

- The plot must stay wide enough to be a chart. If the label column grows, say by how much
  and show the rendered result at a narrow viewport, not just at full width.
- Virtualization: the label is rebuilt per visible row on every scroll, so whatever measures
  the text must not measure per row per frame.
- Serial with the rest of the `timeline-ux-audit` batch.

## Builder Guidance

**Certainty: Exploratory.** The user asked for a more useful view, not for a title column —
this is capture's answer to "make it more useful UIUX", and the trade at its centre is a
judgment call: every pixel the label column gains, the chart loses. Pick a width, render it
at a full and a narrow viewport, and expect to move the number. Say what you picked and
what the render showed.

The tooltip and the title column are separable. If measured-face truncation turns out to
cost more than it returns, a native SVG `<title>` on the row group alone is an acceptable
smaller landing — say so rather than spending the REQ's budget on the harder half.

## Dependencies

`depends_on: [REQ-321]` — **ordering, not logic.** REQ-322 does not need anything REQ-321
produces; it needs REQ-321 not to be editing `web/board-timeline.js` at the same time. Every
REQ in the `timeline-ux-audit` batch writes that one file, and `write_set` is display-only —
`actions/work.md` computes a `--fan-out` wave from `depends_on` alone and explicitly does not
read `write_set`, `batch`, or the Constraints prose. Without this edge the batch's stated
serial requirement was a sentence nothing enforced, and a `--fan-out` run would have
dispatched four concurrent builders into one 1,100-line file.

**The cost, stated rather than hidden:** a chain gates on terminal *success*, so a `failed`
REQ upstream leaves the rest dependency-blocked until someone edits the chain or resolves the
failure. That is the trade for making the metadata say what the prose says.

## Red-Green Proof

**RED prompt/case:** Open the Timeline tab and, without moving the pointer or opening the
table, name any three REQs on screen. Impossible — every row reads `REQ-0NN` and nothing
else. Hover one bar and the description appears at the bottom of the panel, below the hint
paragraph.

**Why RED now:** The row label is `row.id` alone, and the only place a title is rendered is
the foot readout and the collapsed table.

**GREEN when:** each row shows a truncated title next to its id; a title too long for the
column is cut with an ellipsis at a measured boundary rather than overrunning into the plot;
hovering a bar raises a tooltip at the pointer carrying id, title, status and both spans;
and the foot readout still announces the same text for a screen reader.

**Validation:** Inferred during capture — this is an audit finding, not one of the user's
four items. The user asked for a full audit and for the view to be more useful.

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — a column of bare ids, and
the hovered row's description stranded at the foot of the page.

---
*Source: audit finding, UR-065 — "audit the timeline view, and make it more useful UIUX."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is exact but the central trade is a measurement, not a lookup: the
label column's width comes out of the plot's, and how much title fits depends on a face this
module has been burned by twice (`REQ-241`/`REQ-242` measured the same 12px face at 12.0372
and 11.2300 on two builds; `REQ-292` shipped a width model that returned the same number for
every face). Where the truncation boundary comes from has to be found before it is written.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The measurement question the REQ flagged has a much cheaper answer than it looks, and one
thing the REQ assumed turned out to be avoidable:

- **`.timeline-row-label` uses `var(--font-mono)`.** In a monospace face every glyph has the
  same advance, so a single measurement describes every row. That collapses the REQ's
  constraint — "whatever measures the text must not measure per row per frame" — from a
  design problem into one call per render. It is a property of the shipped CSS, not a
  licence: `timelineMeasureLabelAdvance` verifies it and refuses when it does not hold.
- **SVG has no `text-overflow: ellipsis`.** The alternatives were a `clipPath` (hard cut
  mid-glyph, no ellipsis) or a `foreignObject` (HTML layout inside the SVG, and a new class
  of quirk). Neither is needed once one advance is known.
- **The tooltip needs no code at all.** A native `<title>` child of the row group renders at
  the pointer, needs no listener and no positioning, and `timelineRowDescription` already
  composes exactly the sentence it should carry.

## Decisions

- **D-01 — `TIMELINE_LABEL_WIDTH` 104 → 184.** The trade the REQ centres on. At the measured
  6.0219 px/char that is a 28-character budget: seven for the id, two for the separator,
  nineteen to twenty-one of title. Picked from the render at 1400px and 760px, not from
  arithmetic. The plot loses 80px — at 1400px viewport that is 1266 → 1186 of chart, and at
  760px it is 650 → 570, both still a chart. Below about 20 characters the column stops
  earning its width; the browser test asserts that floor so the number cannot quietly drift
  under it.
- **D-02 — The advance measurement refuses in a host that cannot lay out text.** The suite
  caught this: three existing Node probes drive `renderTimelineView` against a DOM stub whose
  SVG node has no `removeChild`, and the first version threw and took the whole render with
  it. A stub genuinely cannot measure text, so the honest answer is the one that already
  existed for a proportional face — no advance, labels fall back to the id alone. Guarded up
  front rather than caught, because a throw here kills a render for a label that was never
  going to be measurable.
- **D-03 — `web/board.css` was declared at capture and ended up untouched.** The label face,
  size and fill were already right; widening the column is a renderer constant, not a style.
  Removed from `write_set` rather than left claiming a file this REQ does not write. (REQ-319
  narrowed a declaration on reasoning and was wrong to; this one is narrowed on a sweep —
  `grep` for `timeline-row-label` returns one CSS rule and it needed no change.)

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)

**What was done:** Each timeline row now reads `REQ-042  Some title`, truncated with an
ellipsis to a budget derived from the label face measured against the live DOM once per
render. The id is never the part that gets cut — a title that fits whole always shows, one
that must be cut needs enough room to be worth cutting, and below that the id stands alone.
The label column widened from 104px to 184px to pay for it. Hovering a bar raises a native
SVG `<title>` carrying what the foot readout carries; the readout stays as the `aria-live`
path.

## Qualification

Passed — 3 files verified in the diff, 6 requirements traced, P-A-U confirmed.

- Mechanical (`tools/checks/qualify.sh`): OK.
- **Substantive:** three new functions, a changed constant, two new call sites and three new
  tests. Not whitespace.
- **Requirements traced:** title beside the id, truncated → `timelineRowLabelText`, pinned by
  eleven cases; truncated against a **measured** face, not an estimate →
  `timelineMeasureLabelAdvance` plus the browser test that proves it responds to the face;
  browser and build recorded → 6.0219 px/char, Chromium 1194 Playwright build, UA
  `HeadlessChrome/141.0.0.0`, and the test `t.Logf`s the number so a future reader sees it
  without re-deriving it; tooltip at the pointer → the native `<title>`, pinned by
  `TestTimelineRowTooltipMarkupMatchesTheProbe`; foot readout stays the aria-live region →
  untouched; a REQ with no title still renders its id → `noTitle` case, and the empty
  tooltip case is covered because `timelineRowDescription` always emits at least the id.
- **Flowing:** not applicable — no data-fetch path; `request.title` was already in the payload.
- **Contamination check:** REQ-321 touched `board.css`, `board-timeline.js`, `template.html`
  and the two probe files. This REQ touches `board-timeline.js` and the two test files —
  expected overlap on the shared surface, declared, and the reason the batch is serial.

## Testing

**Tests run:** `go vet ./...` and `go test -count=1 ./...` in
`skills/do-work-board/tools/queue-kanban`; `bash _dev/tests/maintainer-verify.sh` from the
repo root (`GOTOOLCHAIN=go1.26.1`, `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`).
**Result:** ✓ All passing — canonical gate exit 0, zero FAIL lines, strict JavaScript and
strict browser lanes included.

**Red-green validation:**
- `TestJavaScriptBehaviorTimelineRowLabelTruncation` — **it failed first, and it was right
  twice.** Against my initial rule it reported `a tight budget rendered "REQ-042  Colo…", want
  the id alone`. The rule had one threshold where it needed two: a title that *fits whole*
  should always show however short, and only a title that must be *cut* needs enough room to
  be worth cutting. Both are now separate branches and the test distinguishes them.
- `TestBrowserBehaviorTimelineLabelAdvanceTracksTheFace` — **mutation-verified against exactly
  the defect it exists to prevent.** Replacing the refusal with an average (the shape of
  REQ-292's width model) makes it fail: `a proportional face produced an advance of 5.8359`.
  A plausible-looking number from a face one number cannot describe is precisely how that
  defect shipped last time.
- The same test refuses to run vacuously: it first asserts the proportional override actually
  took (`iiiiiiiiii` and `MMMMMMMMMM` measure differently in the override SVG), because if it
  had not, both faces would be the shipped one and the refusal would prove nothing.

**Measured, and recorded with its build:** the shipped 10px label face is **6.0219 px per
character** in Chromium 1194 (Playwright build at `/opt/pw-browsers/chromium`, UA
`HeadlessChrome/141.0.0.0`) — 28 characters in a 172px column. Per the prime's
per-browser rule this number is expected to differ elsewhere, which is why nothing pins it:
the tests assert a plausible range (4–9 px) and a useful-budget floor (≥20 characters), and
the code measures rather than assumes.

**Render evidence.** Generated a board from this repo's archive and drove it in headless
Chromium at 1400×900 and 760×900 —
`file:///tmp/claude-0/-home-user-skill-do-work/32295d3b-538a-57cc-a4d8-2d453777559b/scratchpad/board322/probe-label.html`,
`location.href` returned with the measurements.

| Viewport | Host width | Widest label | Labels overflowing the 184px column |
|---|---|---|---|
| 1400×900 | 1266px | 162.56px | **0** |
| 760×900 | 650px | 162.56px | **0** |

Labels read `REQ-325  Stop the report-im…`, `REQ-323  Let a timeline bar…`,
`REQ-321  Colour timeline ba…`. The tooltip on the first row reads
`REQ-325 · Stop the report-image interruption path orphaning its backend · no route recorded ·
pending-answers`. The widest label is 21px inside the column at both widths, and the budget is
viewport-independent because the column is fixed — which is what makes one measurement per
render correct rather than merely cheap.

**New tests added:**
- `TestJavaScriptBehaviorTimelineRowLabelTruncation` — eleven cases on the truncation rule.
- `TestBrowserBehaviorTimelineLabelAdvanceTracksTheFace` — the measurement responds to the
  rendered face and refuses a face it cannot describe.
- `TestTimelineRowTooltipMarkupMatchesTheProbe` — the renderer still writes the `<title>` and
  still routes the label through the truncation function.

**Existing tests updated (cross-REQ impact):** none. Three existing Node probes did fail
against the first version of the measurement, and the fix was in the code (D-02), not in them.

## Review

**Acceptance: Partial — Overall 73% — Approve with follow-ups.** Independent review agent,
orchestrated mode, driving a real board and reading the accessibility tree via CDP.

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 85% | 5 of 8 full; the measured-face claim partial for non-Latin, plot figures wrong, narrow case unexercised |
| Code Quality | 78% | Honest guards and correct clamping; two comments stated more than the code did; a Unicode gap |
| Test Adequacy | 70% | Real red-green and a real mutation on the measurement — but the whole suite passed with the widening reverted |
| Scope | 100% | Exactly the three declared files; `write_set` narrowed on a grep |
| Risk | Low | No security or data risk; an a11y regression and a narrow-width one |
| Acceptance | Partial | Works at ordinary widths; two user-visible defects |

**What it verified rather than accepted:** independently reproduced 6.0219 px/char on the same
build; measured Georgia and Arial against the guard's 0.5px tolerance and confirmed both fail
it by three orders of magnitude, so the tolerance is right in both directions; read the AX tree
with CDP; and mutation-tested the widening itself, which is how F3 was found.

## Review Remediation (pre-archive)

- **F3 — the worst of the four, because I wrote a false claim into a decision record.** D-01
  said "the browser test asserts that floor so the number cannot quietly drift under it." It
  did not. Both test files computed the budget from a hardcoded `172` with no link to
  `TIMELINE_LABEL_WIDTH`, so reverting the constant to 104 passed the entire suite — the
  floor was measuring a column the board had stopped using. This is REQ-265's lesson (grep the
  quantity, not the constant name) landing on me one batch after the prime records it. The
  probe now reads the constant through `timelineProbePreamble` and asserts a 20-cell floor
  against it. **Mutation-verified:** 184 → 104 now fails with `the shipped label column fits
  15 cells`. D-01 is corrected above rather than quietly amended.
- **F1 — every row announced its description twice.** The group carried an `aria-label` and,
  after this REQ, a `<title>` with the same 150-character sentence; the review's AX tree showed
  it as both the accessible name and the accessible description on all three hundred rows. The
  `aria-label` is gone: the `<title>` is the single source, it is now the group's first child
  per SVG's own guidance, and it has to exist regardless because it is the pointer tooltip.
- **F2 + M1 — the budget described only the face it sampled.** The guard proves the face is
  monospace *for Latin*, and the budget was then applied to arbitrary Unicode: on the same face
  中 draws 10px and 🙂 12.48px against a 6.02px cell, so a CJK title drew 36px past the column
  and into the plot. Fixed by counting **cells rather than characters** — non-ASCII counts as
  two, the East Asian Width convention, which over-estimates slightly and therefore cuts early
  instead of overflowing. Iterating code points rather than UTF-16 units fixes M1 in the same
  function: a cut can no longer split a surrogate pair into a fallback box.
  **Mutation-verified:** counting every code point as one cell fails with `counted 6 cells,
  want 12`.
- **F4 — the label named nothing for the sixteen newest REQs.** `[impact-rule-change] ` is 21
  characters against a 19-cell title budget, so every review-minted REQ read
  `REQ-306  [impact-rule-chang…` — and REQ-318 put exactly those at the top of the list, 11 of
  the first 25 rows. The label now strips a leading classification tag. That tag exists so a
  human searching the board's title box finds the REQ
  (`actions/capture-reference.md` → REQ Title Convention); it is metadata, not the title's
  substance, and the full title including the tag stays in the tooltip and the table.
  **Mutation-verified:** not stripping it fails with `rendered "REQ-306  [impact-rule-chang…`.
- **Minors fixed:** the constant's comment claimed "about 30 characters … ~21 of title" where
  the measured budget is 28 cells and 19 of title (M4); the measurement comment claimed "not
  per frame" while sitting in the pan path, and now states the measured per-frame cost and why
  it stays there (M2); D-01's plot figures were host widths, 116px high at both ends (M3); the
  `<title>` is the group's first child (N1).

**Routed rather than fixed:**

- **M6 — the narrow-width case degrades, and this REQ made it 80px worse.** Dragging the detail
  drawer to its maximum pins the timeline host at 196px at any window width; bars now start at
  x=184, so 12px of each bar is visible where 92px was. The `Math.max(120, …)` plot floor bound
  before this REQ too — it is pre-existing — but what a 196px chart should do (collapse, hide,
  or keep a floor) is a decision, not a fix. Recorded as a Discovered Task below.
- **Not done:** a screen-reader pass, a non-Chromium engine, and a manual read of the first
  screen. The first two are outside this environment; the third is in Suggested Testing because
  the reviewer is right that pixel measurements are not legibility, and F4 is exactly what that
  gap concealed.

## Discovered Tasks

- **impact-negligible** The Timeline chart degrades to a 12px-wide plot when the detail drawer
  is dragged to its maximum, which pins the scroll host at 196px at any window width. The
  `Math.max(120, …)` plot floor already bound there before REQ-322; widening the label column
  to 184px took the visible bar width from 92px to 12px and pushes 62 rects past the SVG's
  right edge. Pre-existing and worsened rather than caused. What a 196px chart should do —
  collapse the label column, hide the chart, or keep a hard floor and clip — is a design
  decision, which is why this is not folded into REQ-322 as a fix. Found by REQ-322's review.

## Lessons Learned

**What worked:** The monospace property. One measured advance for the whole render, and a
guard that refuses a face it cannot describe — the review independently confirmed the 0.5px
tolerance is right in both directions (Georgia and Arial miss it by three orders of
magnitude). Building the refusal first meant the Unicode fix was a change of unit, not a
redesign.

**What didn't:**

- **I wrote a false claim into a decision record.** D-01 said a test asserted the column-width
  floor. No test did — both files restated `172` instead of reading `TIMELINE_LABEL_WIDTH`, so
  the number the whole REQ was about was pinned by nothing and reverting it passed everything.
  The prime records this exact class as REQ-265, one batch earlier. **A constant a decision
  turns on has to be read by the test, never restated beside it** — and a claim that a test
  exists is checkable in ten seconds, which is ten seconds I did not spend.
- **A verified assumption is only verified for what it sampled.** The guard proves the face is
  monospace using `i` and `M`, and I then applied its answer to arbitrary Unicode. 中 is 10px
  on the same 6.02px face. The guard was not wrong; its *scope* was narrower than the use.
- **Geometry is not legibility.** My render evidence measured widest-label pixels and overflow
  counts, both zero-defect, and quoted three labels that all happened to be non-review REQs.
  The sixteen newest REQs on this very board read `[impact-user-visib…` and named nothing —
  the exact failure the REQ exists to remove, on the first screen, invisible to every number I
  collected.

**Worth knowing:** the label cell is the face's *Latin* advance. Anything that puts new text
in that column — a different script, an emoji, a longer id — is measured in cells by
`timelineLabelCellCount`, not characters, and non-ASCII deliberately over-counts so the error
falls on the side of cutting early. And the `[impact-token] ` title convention and a
nineteen-cell budget cannot both have the front of the string: the label strips the tag, the
tooltip and table keep it.

## Orientation

A Timeline row now names itself: id, title, and a tooltip at the pointer carrying the full
detail. Leaf change in the board's frontend — no new module, no payload field. The one thing
worth knowing outside this view is that the label column is measured rather than assumed, and
the measurement refuses faces it cannot describe. `_dev/primes/prime-kanban-board.md`
spot-checked: every path it references still exists.
