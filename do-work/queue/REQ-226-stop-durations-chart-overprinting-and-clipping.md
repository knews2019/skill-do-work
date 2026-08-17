---
id: REQ-226
title: Stop the Durations chart from silently overprinting and clipping
status: pending
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
write_set:
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/durations_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# Stop the Durations Chart from Silently Overprinting and Clipping

## What

Fix two defects in the Durations view where the chart draws something that reads as a value but
isn't. Panel A labels every overflow sample with no collision detection, so on a large board the
overflow lane becomes an unreadable blob of overprinted text. Panel B clamps its bars at 45 minutes
with no scale break, so a 78-minute day renders as a 45-minute bar.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
