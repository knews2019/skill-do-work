---
id: UR-051
title: Fix the Durations chart's silent overplotting and add a forecasting Gantt timeline
created_at: 2026-08-17T23:51:17Z
requests: [REQ-226, REQ-227, REQ-228]
word_count: 84
---

# Fix the Durations Chart's Silent Overplotting and Add a Forecasting Gantt Timeline

## Summary

Two requests against the board's Durations view, plus an invitation to name anything else worth
fixing on the same surface.

The first is a defect visible in the attached screenshot: on a 560-REQ board, Panel A's overflow
lane collapses into an unreadable blob of overprinted text. The cause is at
`skills/do-work-board/tools/queue-kanban/web/board-durations.js:246-265` — every sample above the
60-minute ceiling gets a direct label, placed by a four-step index cycle (two text rows × two
anchor sides) with **no collision detection at all**. Two overflow REQs completing within ~4% of the
time range already overprint, because the alternation keys on array index rather than on proximity.
This repo's own archive has only three overflow samples, which is why the defect is invisible here
and total on the reporting board.

The second is a new view: a zoomable, scrollable Gantt of per-REQ timing, with each REQ drawn as one
two-segment bar — `created_at`→`claimed_at` (the wait) and `claimed_at`→`completed_at` (the work) —
extended forward with projected bars for the remaining queue so the reader can see when the queue
empties and in what order.

Both stamps the Gantt needs are already parsed onto every ticket
(`skills/do-work-board/tools/queue-kanban/model.go:76-78`) and already shipped to the browser for
every REQ including pending ones (`generate.go:440-487`). No new frontmatter field and no second
walk of the archive are required.

The third REQ answers the maintainer's closing question directly. Two things were found on the same
surface that the request did not name: Panel B clips at 45 minutes with no scale break (so a
78-minute day renders as a 45-minute bar, visible in the same screenshot), and in-flight REQs are
invisible everywhere on the board's time views — a REQ running right now contributes to nothing.

## Extracted Requests

| REQ | Title | Covers |
|---|---|---|
| REQ-226 | Stop the Durations chart from silently overprinting and clipping | The screenshot defect (collision-aware labels + `+N more` count) and the unrequested Panel B clipping defect found on the same surface |
| REQ-227 | Add the Timeline view with two-segment REQ bars | The Gantt itself — wait/work segments, in-flight bars, zoom and scroll, fifth view tab |
| REQ-228 | Project the remaining queue onto the Timeline | The forecast — serial median-based projection, dependency-then-queue order, "queue empties ~X" |

## Batch Constraints

- **Nothing here adds a board write surface.** All three REQs are read-only rendering and
  aggregation. `CLAUDE.md` § *Kanban Board Write Surfaces* states the board has exactly three write
  surfaces and none touches pipeline state; that sentence must still be true and unamended when this
  batch lands. A forecast is a derived display, never a written field.
- **No new frontmatter field.** `created_at`, `claimed_at`, and `completed_at` already exist on every
  ticket and already reach the browser. A REQ in this batch that proposes a new schema field has
  misread the data layer.
- **Parser lock-step** (`_dev/primes/prime-kanban-board.md:13`). Any new payload type must keep
  `model.go` aligned with the status vocabulary in `skills/do-work/actions/work-reference.md`.
- **The JS fragment inventory is pinned by a test.** `boardJavaScriptFragmentPaths` at
  `generate.go:42-52` and the assertion at `generate_test.go:28-62` must change in the same commit as
  any new `web/*.js` file, or the build fails.
- **The live-archive figures in `durations_test.go:202-225` are pinned** (≥195 samples, the
  2026-07-31 and 2026-08-15 medians). Nothing in this batch may change what those assert — this is a
  rendering and forecasting batch, not a re-definition of the read-time rule.
- **The read-time rule has one definition**, in `skills/do-work/actions/estimate-reference.md` →
  Calibration, and `durations.go:9-26` is explicit that it is that rule's *second reader, not a
  second definition*. Any forecast that needs to exclude paused or reversed spans must consume the
  existing `DayMedianExclusion` verdict, never re-derive the thresholds.
- **A forecast must never read as a fact.** Projected geometry is visually distinct from measured
  geometry at a glance, and the projection's assumptions are stated on the view rather than only in a
  REQ file.
- **Release ceremony belongs to the integrating commit only** — the shared version bump across
  `VERSION`, `skills/do-work/VERSION`, and `skills/do-work/actions/version.md`, plus a `CHANGELOG.md`
  entry titled for what shipped and mirrored byte-identically into `skills/do-work/CHANGELOG.md`
  (`CLAUDE.md` § Before Every Commit).

## Maintainer Decisions at Capture

All three were answered interactively before any file was written.

1. **Forecast model: serial and median-based.** One REQ at a time in dependency-then-queue order.
   Each pending REQ's projected span is the rolling median of recent completed spans, split by
   `effort_estimate` (`trivial` | `normal`). Chained from now to produce a queue-end estimate.
   Projected bars drawn hatched or dimmed so a forecast never looks like a measurement. A
   configurable parallelism/lane knob was offered and **declined** — worktree fan-out is real, but the
   serial model is the honest floor and the knob is machinery this request does not need.
2. **Label fix: collision-aware placement plus a count.** Place labels greedily, skip any that would
   overlap one already placed, and print `+N more over 60 min` at the lane edge so nothing is
   silently hidden. A top-N-extremes rule and a drop-the-labels-entirely option were both offered and
   declined — the first leaves points near the cutoff unlabelled with no visible logic, the second
   loses the at-a-glance "REQ-407 took 14h" callout that makes the lane useful.
3. **Placement: a new fifth "Timeline" tab.** Board / Calendar / Durations / Timeline / Testing.
   Durations keeps its statistical job (spread, cadence); Timeline owns per-REQ sequencing and the
   forecast. Folding the Gantt in as a fourth Durations panel was declined because it needs its own
   zoom and a forward-extending axis, which fights the three shared-axis panels. Replacing Durations
   was declined because it discards working panels the maintainer did not ask to lose.
4. **Three REQs, not one**, sliced so the speculative half can be iterated or dropped without holding
   up a working Gantt. REQ-227 and REQ-228 unavoidably share `board-timeline.js` and `timeline.go`;
   that overlap is declared in both write sets rather than distorting the slice to avoid it.
5. **The two unrequested findings are in scope.** The maintainer's closing question ("Am I missing
   anything, can we improve anything else?") was answered with Panel B's silent clipping and the
   invisibility of in-flight REQs, and both were accepted into the batch — the clipping into REQ-226
   as a same-file same-class defect, the in-flight bars into REQ-227 where they also make REQ-228's
   forecast honest by giving the first projected bar something to start from.

## Provenance and Injection Check

Source: the maintainer, typing directly into this session, with one attached screenshot of their own
board (`glw-game-find-the-difference`, generated 2026-08-17 21:20 UTC). Per
`crew-members/prompt-injection.md` the input was read as data. **No injection detected** — the
request contains no instruction-like content beyond the ordinary imperatives of a feature request,
and the screenshot is a rendering of the maintainer's own queue.

Note that the counts and REQ ids visible in the screenshot (560 archived REQs, REQ-407) belong to
that other repository, not to this one. They are evidence that the defect scales with sample count;
they are not this repo's data and no REQ in this batch should try to reproduce them locally.

## Full Verbatim Input

do-work capture-request: fix this: [Image #1]

Also make another one that is a zoomable and scrollable gant chart, and includes also the timing for the remaining tasks so we know when it ends and in what order.

gant chart contains also duration, which can be calculated from the captured time, start time and finish time (one bar with 2 colors) distance between captured and started is one, and started and finished there is another one.

Am I missing anything, can we improve anything else?

### Attached Screenshot

`assets/REQ-226-screenshot-1-durations-overflow-label-blob.png` — the Durations view of the
maintainer's `glw-game-find-the-difference` board at 127.0.0.1:8090, generated 2026-08-17 21:20 UTC.
Header reads "560 archived REQs with both stamps, across 42 active days. Panel B excludes 14 spans
from its medians". Panel A's overflow lane at the `60+` tick shows one legible label at the far left
(`REQ-407 14h 15m`) and, from roughly 1 July rightward, two horizontal bands of completely
overprinted text where several hundred labels have been drawn on top of one another — individual
characters are visible but no id or duration is readable. Panel B shows a `78 min` annotation at the
top right sitting above a bar clipped flat at the 45 tick. Panel C is legible and undamaged.

### Follow-Up Instructions in the Same Session

After the design questions were answered, the maintainer interrupted the planning pass with
"stop when safe", then asked whether starting a fresh session with `do-work run` was sufficient
(it was not — nothing had been written to disk yet), then instructed: "capture it in this session",
and immediately after: "capture it in this session, then comit".

---
*Captured: 2026-08-17T23:51:17Z*
