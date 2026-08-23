---
id: REQ-323
title: "Let a timeline bar be read against a date"
status: completed
created_at: 2026-08-22T22:08:34Z
claimed_at: 2026-08-23T03:58:20Z
completed_at: 2026-08-23T04:40:00Z
commit:
route: B
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-322]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-318, REQ-319, REQ-320, REQ-321, REQ-322, REQ-324]
batch: timeline-ux-audit
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-08-23T03:58:50Z
  basis:
    - Route B
    - 4-file write set
    - 6 acceptance criteria
    - dependency depth 5
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Let a Timeline Bar Be Read Against a Date

## What

Three additions that make a bar's position mean something: gridlines through the plot, a
drawn queue-end rule, and a minimum bar that stays readable when the window is wide.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`. Three parts, one file each side.
  Gridlines: extract the tick instants `renderAxis` already computes into
  `timelineAxisTickInstants` and draw from that list into the rows SVG, before the rows so
  paint order puts them behind every bar. Queue-end: a `queueEndRuleInstant()` predicate both
  halves read, rule in the rows SVG (same x scale as the bars, the reason `drawNowRule` lives
  there), label in the non-scrolling axis SVG so a caption cannot scroll away from the mark it
  names. Minimum bar: a pure `timelineCollapsedRowMark(segments, …)` over the segment list
  windowing already computes, returning the extent to draw one marker over or null.
- [x] **[APPLY]:** Five files. `web/board-timeline.js` (the three parts),
  `web/board.css` (three rule styles plus three legend swatches), `web/template.html` (the
  Vertical rules key group and the collapsed-marker entry), `generate_test.go` (two Node
  probes), `timeline_browser_probe_test.go` (the prominence probe plus its markup pin).
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file. `node --check web/board-timeline.js`
  clean; `go vet ./...` clean; `bash _dev/tests/maintainer-verify.sh` exit 0. No debug
  artifacts, no console statements, no leftover fixtures. One dead branch found and deleted
  during the pass (see D-03).

## Why

The axis is a 26px strip above a 58vh scroll container. A bar two hundred rows down has no
vertical reference at all, and at *Fit all* over three months every completed REQ is a
1.5px sliver in which the wait and the work segment occupy the same pixel. Meanwhile the
forecast paragraph names a queue-empty instant that the chart never marks, though the
*Now* button already knows where it is.

## Context

Three separate causes, one reading problem, so one REQ:

- `renderAxis` draws ticks into `axisSvg` only; the rows SVG gets `drawNowRule` and nothing
  else.
- `drawSegment` floors width at `Math.max(1.5, …)`, and both segments are floored
  independently, so a short REQ at a wide zoom is two adjacent 1.5px marks of different
  hues.
- `timelineNowJump` already computes `queueEndMs` from `projection.queueEnd` and frames it
  in the window; nothing draws it.

`drawNowRule` is inside the rows SVG on purpose — that is what guarantees it shares the
bars' x scale. Anything added here follows the same rule.

## Detailed Requirements

- **Gridlines.** A faint vertical line at each axis tick, running the full height of the
  rows area, drawn from the same tick instants `renderAxis` computes so the two can never
  disagree. Faint enough to sit behind the bars, not compete with them.
- **Queue-end rule.** When the projection is confident and its queue-empty instant falls
  inside the window, draw a labelled vertical rule at it, visually distinct from the
  now-line, and name it in the legend. When the projection declined, or the instant is
  outside the window, draw nothing.
- **Minimum bar.** Raise the floor so a bar is visible rather than technically present.
  When a whole row's span is narrower than a readable two-segment bar, draw one marker for
  the row instead of two adjacent slivers claiming a wait/work split the pixels cannot
  show. Pick the threshold from the rendered result, and state the number and where it came
  from.

## Constraints

- Gridlines must not cost the virtualization: they are a handful of nodes per render, not
  per row.
- The now-line stays the most prominent vertical mark. A queue-end rule that outshouts it
  moves the reader's eye to a forecast instead of to the present.
- Both themes.
- Serial with the rest of the `timeline-ux-audit` batch.
- **The screenshot is the test here more than anywhere else in this batch.** Every one of
  these three changes is a claim about pixels — a gridline faint enough to sit behind the
  bars, a rule that does not outshout the now-line, a bar wide enough to see. Generate a
  board, look at it in both themes, and attach what you saw
  (`_dev/primes/prime-kanban-board.md`); a green suite says nothing about any of them.

## Builder Guidance

**Certainty: Exploratory.** Three numbers have to be chosen — gridline weight, queue-end
rule prominence, minimum bar width — and none of them can be argued to a value in the
abstract. Pick each from the rendered result, state the number and where it came from, and
expect to be asked to move one.

Landing the three parts one at a time is fine; they share a reading problem, not an
implementation. If one of them turns out to need its own REQ, say so instead of stretching
this one.

## Dependencies

`depends_on: [REQ-322]` — **ordering, not logic.** REQ-323 does not need anything REQ-322
produces; it needs REQ-322 not to be editing `web/board-timeline.js` at the same time. Every
REQ in the `timeline-ux-audit` batch writes that one file, and `write_set` is display-only —
`actions/work.md` computes a `--fan-out` wave from `depends_on` alone and explicitly does not
read `write_set`, `batch`, or the Constraints prose. Without this edge the batch's stated
serial requirement was a sentence nothing enforced, and a `--fan-out` run would have
dispatched four concurrent builders into one 1,100-line file.

**The cost, stated rather than hidden:** a chain gates on terminal *success*, so a `failed`
REQ upstream leaves the rest dependency-blocked until someone edits the chain or resolves the
failure. That is the trade for making the metadata say what the prose says.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repo's archive, press *Fit all*, and scroll
two hundred rows down. Ask: what date does this bar start on? There is no vertical reference
anywhere in the plot — the axis is off-screen above. Scroll back up and read the forecast
sentence: it names a queue-empty instant, and no mark on the chart corresponds to it. Look
at any completed REQ at this zoom: a 1.5px grey tick with no distinguishable wait or work.

**Why RED now:** Ticks are drawn only into the axis SVG; nothing draws `projection.queueEnd`;
`drawSegment` floors each segment independently at 1.5px.

**GREEN when:** gridlines run from each axis tick through the rows area at every zoom
level; a labelled queue-end rule appears when the projection is confident and its instant is
in the window, and is absent otherwise; a sub-threshold row draws one legible marker rather
than two 1.5px slivers; and the legend names every vertical rule and marker the chart can
draw.

**Validation:** Inferred during capture — audit findings, not one of the user's four items.

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — thirty-odd rows of
one-to-three-pixel ticks under an axis they cannot be measured against.

---
*Source: audit finding, UR-065 — "audit the timeline view, and make it more useful UIUX."*

---

## Triage

**Route: B** - Medium

**Reasoning:** Three changes, one reading problem, and all three are numbers that cannot be
argued to a value in the abstract — the REQ says so itself. Where each mark goes is known
(`renderAxis` computes the tick instants, `drawNowRule` shows the pattern for a rule inside
the rows SVG, `drawSegment` holds the width floor); what each should look like has to come
from a render.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

---

## Decisions

**D-01 — Gridline weight: `--line-firm` at `opacity: 0.6`.** `--line-firm` is the axis tick's
own colour, so the gridline is visibly a descendant of the tick above it rather than a fourth
grey from nowhere. The 0.6 came from the render: at full strength on Chromium 1194 the
gridline's measured presence equals the axis tick's exactly (0.1640 light, 2.3811 dark), which
inverts the strip-is-the-measured-edge reading and makes the plot look like ruled paper with a
chart on it. `TestBrowserBehaviorTimelineVerticalRulesRankByProminence` now holds both ends —
above a 1.03:1 visibility floor, below the tick's presence — in both schemes.

**D-02 — Queue-end rule: `--timeline-projected`, dashed `2 4`, `opacity: 0.6`.** Violet because
it is the colour the projected segments already use, so the forecast reads as one thing;
dashed because presence alone can be matched by two marks of different meaning. 0.6 rather
than the 0.75 first shipped: at 0.75 the rule cleared the now-rule by 27% in the light palette,
a margin one palette tweak erases. At 0.6 it clears the test's 75%-of-the-now-rule ceiling with
room. The label lives in the axis SVG, which does not scroll; the rule lives in the rows SVG,
which shares the bars' x scale.

**D-03 — `TIMELINE_MIN_SPLIT_WIDTH = 7`, and `TIMELINE_MIN_SEGMENT_WIDTH` raised 1.5 → 3.** Seven
is what a two-segment bar physically needs: two floored 3px segments plus a pixel of boundary.
The test derives it (`2*segmentWidth + 1`) from the shipped constants rather than restating 7,
so lowering either constant fails rather than passes quietly. Three came from the Fit-all
render: 1.5px was technically present and practically invisible.

While writing the mutation checks I found the collapse function's `isFinite` guard was
unreachable — `timelineRowSegments` only ever emits the `-Infinity → Infinity` sentinel *alone*,
so the `segments.length < 2` guard always catches it first. Deleting the guard changed no test
result, which is what proved it dead. It is gone, and the surviving guard's comment now names
both cases it covers.

**D-04 — One extra file beyond the write set.** `timeline_browser_probe_test.go` was not
in `write_set`; the prominence measurement belongs in the browser lane, which lives in
`timeline_browser_probe_test.go`, so that file is the fifth. Nothing else in the batch writes
it.

## Evidence

- `bash _dev/tests/maintainer-verify.sh` → **exit 0** (`GOTOOLCHAIN=go1.26.1`,
  `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`). Strict JavaScript lane and strict browser
  lane both PASS.
- Rendered and inspected in **both** themes on Chromium 1194, at *Fit all* over three months
  and at the *Now* window: gridlines run the full plot height at every tick and sit behind the
  bars; the queue-end rule and its `queue empty` caption appear at 04:44 with the now-line at
  04:16 plainly the louder of the two; short REQs draw one legible marker where they used to
  draw two slivers. A *Day* window on 11 Jul draws no queue-end rule, which is the
  out-of-window case.
- Every new assertion mutation-verified. Split threshold 7→1, the one-segment guard removed,
  the marker anchored at its start instead of its extent, `renderAxis` given a private tick
  loop, the broken-stamp guard removed at the call site, the gridline at full strength, the
  gridline at 0.01, the queue-end rule made solid and full-strength, the now-line made dashed —
  each fails the intended assertion and only that one.

## Review Fixes

Three defects found reviewing my own diff after the first green gate.

**R-01 — A collapsed row upgraded a forecast to a fact.** The collapse excluded broken-stamp
rows and nothing else, so a narrow *pending* REQ — open wait plus hatched projection — drew one
SOLID status-coloured marker. The wait/work split is the claim the collapse withdraws;
measured-versus-projected is a different claim, and the one a reader trusts hardest. Rows
carrying a projection are now excluded for the same reason broken rows are, and the exclusion
is pinned by test in both halves: collapsible by width alone, spared by the caller.

**R-02 — A comment claimed a saving the code does not make.** "Measured once per render and
handed to every row … asking it per row would force one layout per row" was false:
`drawSegment` still calls `plotWidth()` on every call. The hoisted read serves the collapse
decision only, and hoisting `drawSegment`'s is REQ-324's stated scope. Comment corrected rather
than the scope widened.

**R-03 — The rows description promised a rule the chart draws only sometimes.** The accessible
name said the dashed violet rule marks the queue-empty instant, full stop. It now says when.

Gate re-run after the three fixes: `bash _dev/tests/maintainer-verify.sh` → **exit 0**. The new
exclusion mutation-verified (removing it fails the call-site pin and nothing else).

## Discovered Tasks

None new. The narrow-width chart degradation recorded on REQ-322 still stands, and REQ-325 is
still `pending-answers`.

## Lessons Learned

- **A prominence metric that ignores a channel the design uses will contradict the render.**
  The first version of the vertical-rule probe measured contrast times stroke width and
  reported the light-theme forecast rule as louder than the now-rule — the opposite of what the
  screenshot showed. Dash duty was the missing channel: a `2 4` rule inks a third of its length
  and a `3 3` rule inks half. The fix was the metric, not the design, and the way to tell which
  is which was having the screenshot first.
- **A guard that no mutation can break is either dead or untested — check which before
  assuming the second.** Removing the `isFinite` branch from the collapse function changed no
  test result. The instinct is to write a test that reaches it; the truth was that nothing
  can, because the sentinel segment is only ever emitted alone.
