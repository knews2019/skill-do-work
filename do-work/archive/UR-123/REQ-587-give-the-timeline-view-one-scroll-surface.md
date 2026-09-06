---
id: REQ-587
title: 'Give the Timeline view one scroll surface, in the same style as the Activity view'
status: completed
created_at: 2026-09-05T12:50:59Z
user_request: UR-123
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: [REQ-585]
related: [REQ-585, REQ-586]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-09-05T17:54:34Z
estimate:
  p50_active_minutes: 50
  confidence: medium
  calculated_at: 2026-09-05T17:55:26Z
  basis:
    - Route C
    - 8-file write set
    - 1 new files
    - 2 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - browser evidence
    - cross-route regression gates
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
  - skills/do-work-board/tools/queue-kanban/timeline_scroll_browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
route: C
dispatch_at: 2026-09-05T18:13:36Z
builder_handback_at: 2026-09-05T19:05:52Z
review_at: 2026-09-05T19:43:13Z
integration_at: 2026-09-05T19:05:52Z
planning_at: 2026-09-05T18:09:17Z
commit: 8fad73b20055bbe66df91d423027867d780f3175
completed_at: 2026-09-06T03:36:10Z
release_at: 2026-09-06T03:36:10Z
---

# Give the Timeline View One Scroll Surface, in the Same Style as the Activity View

## What

The Timeline view scrolls twice, like the Activity view did before REQ-585 (give the Activity view one scroll surface): the chart rows live in `.timeline-scroll`, a box capped at `max-height: 58vh` with `overflow-y: auto`, inside `.board-main`, which is the board's scroll container. Leave exactly one scroll surface on this view, in the same style as the Activity fix: the board scrolls, the rows are content, and what the reader needs to keep in view (the time axis, as the column header is on Activity) sticks to the top edge.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Followed the plan's tasks T1 through T6 with its D-02 through D-13 taken as written, having re-verified each of its claims against source — all five Node-stub sites, the invalidation point, and the claim that no existing probe writes a scroll position on the chart. Took the orchestrator's override to pin the anchor requirement in the probe with its own display-list guard rather than checking it by hand.
- [x] **[APPLY]:** Layout, then the renderer, then the view switcher, then the four Node-stub files, then the new probe, then rebuild and measure. Scope stayed inside the declared list until the browser lane forced one further file, which was disclosed rather than absorbed.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file; `go vet ./...` and `gofmt -l .` clean; no print, log, debugger, TODO, FIXME or XXX anywhere in the diff. Per file: the old block comment no longer claims the chart is the scroll container and the chart carries exactly two declarations, neither an overflow; the axis paints the base token and carries the one-pixel rule the rows box used to own, so total height is unchanged; the three padding rules cover every arrangement of the four possible preceding siblings and leave exactly one carrying the inset; the Activity view's block is byte-identical. In the renderer, a grep for the old scroll reads returns nothing but one updated comment, and each of the twenty-four remaining references to the chart host was read individually to confirm it is a host, width, focus, hover, keyboard, wheel or pointer site rather than a scroll one. The view switcher's three lines sit after the panel toggle and before the first render. The five Node-stub lists each gained one entry and no assertion changed, which is why every row count and row id in that lane came back identical.

*Authored by the builder in `do-work/runs/work-2026-09-05-170806/REQ-587-handback.md` → `## P-A-U`; transcribed and checked here by the orchestrator, which is the only writer of this file in worktree dispatch mode. The orchestrator re-read the merged range itself before checking the UNIFY box.*

## Context

This is not the CSS-only change REQ-585 was. The Timeline rows are virtualized: `board-timeline.js` reads `scrollHost.scrollTop` and `scrollHost.clientHeight` (the element with id `timeline-scroll`) to decide which rows carry nodes, to keep the same REQ under the pointer when the window changes, to scroll a focused row into view, and to jump to the open work (roughly lines 1618, 1757 to 1797, 2077 to 2117, 2803 to 2806, 3186). The block comment above the timeline CSS (board.css near line 2579) states that `.timeline-scroll` is the scroll container by design. Two ways to get one scroll surface:

- **Re-point the scroll host to `.board-main`.** The chart becomes full height as ordinary content; the visible-row math reads the board's scroll position minus the chart's offset inside it; the axis (`.timeline-axis`) becomes sticky under the board's top edge. This is the same reading as REQ-585's M1 and is the recommended shape. Wheel-to-zoom (Ctrl or Cmd plus wheel) and drag-to-pan must keep working on the chart.
- **Let the chart fill the remaining viewport and stop the board scrolling on this view** (REQ-585's M2 shape). Cheaper, but the forecast line, the excluded list, the hint paragraph and the "REQs in this window, as a table" details sit below the chart today and would need a new home, and the Timeline would have a scroll model no other view has.

The Activity mockups and the side-by-side comparison of the two shapes are in `ai-reports/2026-09-05_1520_activity-view-double-scroll-mockups/index.html`. REQ-585's merged diff shows how the padding move was scoped to one view; reuse that mechanism rather than inventing a second one. Keyboard behaviour named in the view's hint paragraph (Tab to the chart, arrows move between rows and pan, plus and minus zoom) must survive.

## Red-Green Proof
**RED prompt/case:** On the served board, open the Timeline with the All days window and run in the console: `const m=document.querySelector('.board-main'), t=document.getElementById('timeline-scroll'); [m.scrollHeight>m.clientHeight, t.scrollHeight>t.clientHeight]`.
**Why RED now:** It returns `[true, true]` and two scrollbars are visible (screenshot under Assets: the chart's own scrollbar at the right edge of the rows, the board's behind it).
**GREEN when:** Exactly one element scrolls, the axis stays visible while scrolling through the rows, and the rows still render correctly at every scroll position (the virtualized range follows the scroll: a row 400 px below the fold appears after scrolling there, and the row under the pointer stays put when the window chips change). Measured in a real engine and recorded in the hand-back; a browser probe in the existing `*_browser_probe_test.go` lane pins it.
**Validation:** Inferred during capture

## Open Questions

- [~] Which shape: re-point the scroll host to the board (recommended; same style as the Activity fix, axis sticky) or let the chart fill the remaining viewport with the board not scrolling on this view. → **D-01**: Builder chose: re-point the scroll host to `.board-main`, the recommended shape. Reasoning: the request asks for "the same style as the Activity view", and that is what REQ-585 did — it deleted the inner scroll box's `max-height` and `overflow`, let the already-sticky header re-pin to the board, and moved the board's top padding onto the view's first element scoped with `.board-main:has(> #view-activity.is-active)`. The archived record names REQ-587 explicitly and asks it to reuse that scoping rather than invent a second one. Claim-time exploration also priced the alternative: the forecast line, the excluded list, the hint paragraph, the aria-live readout and the "REQs in this window, as a table" details are all siblings of the chart inside the view, so the fill-the-viewport shape needs a new home for five elements and gives the Timeline a scroll model no other view has. **Value:** one scroll surface with the axis pinned, reached by the mechanism the sibling view already proved in a real engine, and no new layout concept in the board. **Risk:** the Timeline's rows are virtualized and the Activity table was not, so this shape carries work the Activity fix did not — five scroll-position reads must convert to the board's coordinate space, the shared `.board-main` scroll position is not reset between views, and the verify-findings strip sits outside the view panel where the padding-move recipe cannot reach it. All are recoverable in code; the reversal is deleting one CSS block and restoring the element reference.

<!-- D-XX counter: last used D-01. Next decision: D-02. -->

## Required Lessons — Dropped for Budget

Refreshed at claim time against `do-work/lessons-index.md`; both entries still miss the 2000-token budget in `actions/capture-reference.md` → Required Lessons Budget Contract, and both are `slugged: partial`, so the targeted `path#family-slug` form is not legal and neither can be narrowed to fit.

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 6224 tokens at claim time (5744 at capture). Matches on "Changing queue-kanban model, parser, UI, timeline, testing, or browser behavior".
- `_dev/primes/lessons-kanban-board.md` — 4959 tokens at claim time (4820 at capture). Matches on "Changing queue-kanban parsing, views, view scroll or sticky layout, static output, timeline behavior".

Dropped as stamped entries only. The touch-conditional Lessons Discipline rule is separate and fires here through `_dev/primes/prime-kanban-board.md`, which this REQ declares and whose Read-first entries name the timeline and the board's SVG surface — so the builder is told to read both satellites anyway, and `_dev/primes/lessons-kanban-board.md` carries `[family: sticky-pins-to-content-box]`, which is the single most relevant lesson in the repository to this change.

## Assets

- `do-work/user-requests/UR-123/assets/REQ-587-screenshot-1-timeline-double-scroll.png`: the live board at 127.0.0.1:8090, Timeline view, All days window, generated 12:38 UTC on 2026-09-05. Above the chart, the verify-findings strip with five worktree cards, the "How long each REQ waited, and how long it took" heading, the legend, and the period controls (Last day to All days, From 05/27/2026 to 09/07/2026). The chart lists REQs grouped by UR (UR-120 down to UR-114) with bars at the right edge near the now-line. Two scrollbars on the right: the chart box's own, starting at the axis, and the board's thin one behind it starting under the top bar.

*Source: "BTW: this view also have a scroll issue, add a REQ for it as well: - the style should be similar"*

---

## Triage

**Route: C** - Complex

**Reasoning:** The Activity view's version of this was a CSS-only change to one block. This one is not. The Timeline's rows are virtualized, so five separate sites read the scroll position to decide which rows carry nodes, to keep the same REQ under the pointer when the window changes, to scroll a focused row into view, and to jump to the open work — each has to move into the board's coordinate space, and one of them must also subtract the height of the now-sticky axis or a focused row lands underneath it. Beyond that the change reaches a Node test harness whose stub DOM resolves a fixed id map across four Go test files, a browser probe lane, and a findings strip that lives outside the view panel and so is not reachable by the padding-move recipe this REQ is told to reuse. Multiple systems, a real architectural choice between two shapes, and a written plan is worth having before any of it is written.

**Planning:** Required.

## Plan

*Generated by Plan agent. The full working plan, with every code block, lives in the run directory as `do-work/runs/work-2026-09-05-170806/REQ-587-plan.md` and was handed to the builder; this section is the durable record of what it decided and why.*

**Six corrections the plan verified against source before planning anything.** `.board-main` has nine direct children, not four — `#board-notes`, `#board-anomalies` and `#board-findings` can all precede the view panel, which widens the findings-strip risk into a general one. The Kanban panel's id is `view-board`. `#view-timeline` has no CSS rule of its own, so REQ-585's "fold the view's own padding into its first element" half has nothing to fold — only the board's 24px moves. `renderTimelineView` re-runs on a filter change and releases its listeners first, so moving the scroll listener onto a node that outlives the render leaks nothing. No existing browser probe writes `scrollTop` on `#timeline-scroll`. And the Node stub's host-resize helper only sizes `timeline-scroll`, whose `clientHeight` of 400 the board stub also carries, so the lane's row counts are unchanged by the host split.

**Files, in order.** `web/board.css` (the layout change: rewrite the block comment that calls `.timeline-scroll` the scroll container, strip it to `position: relative` + `cursor: grab`, give `.timeline-axis` its first rule, move the padding and the focus ring) → `web/board-timeline.js` (one new host lookup, one memo pair, six call sites) → `web/board-controls.js` (three guarded lines) → the four Go test files carrying the Node stub's id map → a new `timeline_scroll_browser_probe_test.go` → rebuild, serve, measure and screenshot.

**Decisions.**

- **D-02 — `.timeline-scroll` loses `max-height` and *both* `overflow` declarations. DECIDE & STATE.** Leaving `overflow-x: hidden` beside an `overflow-y: visible` makes the visible axis compute to `auto`, so the element stays a scroll container while measuring `scrollHeight == clientHeight` — the fix would look applied and be half-applied. The probe asserts the computed `overflow-y` is `visible`, not only the height comparison, because the height comparison alone cannot see this.
- **D-03 — the sticky axis paints `var(--bg-base)`, not `var(--surface-1)`. DECIDE & STATE.** The prime states the surface behind this board's SVG is `<body>`; copying the Activity header's token would paint a visibly tinted strip over an untinted chart, which is the defect REQ-321 and REQ-346 each paid for once. `z-index: 1` is scoped inside the board because `.board-main`'s `container-type: inline-size` establishes a stacking context.
- **D-04 — the 1px separator moves from the rows box's `border-top` to the axis's `border-bottom`. DECIDE & STATE.** A border on the rows box scrolls away exactly when the axis is doing its job. Same pixel, on the element that stays on screen; total height unchanged.
- **D-05 — the board's top padding moves onto the first *unhidden* child of `.board-main`, keyed on that condition rather than on a list. DECIDE & STATE.** Three strips can precede the view panel and none has a top margin, so zeroing the board's padding would pull whichever one is showing flush against the top bar. Naming `#board-findings` alone fixes the screenshot and breaks the day a reader has notes and no finding. Uses `margin-top`, not `padding-top`, because the strips have a border and a tint; the board's `overflow-y: auto` establishes a block formatting context so the margin cannot collapse through. The arithmetic is a no-op today in both arrangements.
- **D-06 — the chart's focus ring moves onto the sticky axis. ESCALATE, proceeding with this default.** On a full-height element the inward 2px ring puts its top and bottom edges thousands of pixels apart, both off screen, and "Tab to the chart" is shipped text. The axis is the part of the chart always in view. **Value:** the keyboard affordance the hint promises stays visible at any scroll position. **Risk:** a reader may read the ring as "the axis has focus" rather than "the chart has focus"; reversal is one rule. Escalated because it is a visible accessibility affordance and reasonable people would pick differently.
- **D-07 — entering the Timeline resets `.board-main`'s scroll position to 0. ESCALATE, proceeding with this default.** Nothing resets that position between views. Today the Timeline page is short enough that arriving from a scrolled Kanban clamps near the top; once the chart is the tall thing on the page, the reader would land arbitrarily deep into it. **Value:** arrival stays predictable and matches what the view does today. **Risk:** leaving the Timeline and returning loses your place, which today survives because the inner box keeps its own position; reversal is three lines. The alternative, if the maintainer values the preserved place, is to remember and restore the board's position on leave and entry — about six lines of new state in the view switcher.
- **D-08 — the two measured geometry numbers are cached once per render and invalidated in `renderAll` beside the existing plot-width invalidation. DECIDE & STATE.** The visible-row render *is* the scroll listener and the drag-pan calls it once per frame; two extra `getBoundingClientRect()` calls inside it would add two forced layouts to every frame of a drag. Neither number can change without a re-render.
- **D-09 — the visible-range viewport stays the board's `clientHeight` with no axis subtraction. DECIDE & STATE.** Not subtracting the axis renders a superset — at most one extra row hidden behind it. Not subtracting the rows offset from the scroll *position* is not optional: the range would walk off the bottom and the top rows go blank. One number has to be right; the other only has to be generous.
- **D-10 — only the `scroll` listener moves to the board. DECIDE & STATE.** It must move, because `#timeline-scroll` stops emitting scroll events and the virtualized range would freeze at whatever the last render computed — the single failure most likely to ship in a change that otherwise looks complete. Wheel-zoom, drag-pan, pointer capture, keyboard, focus, hover and the width observer all stay on the chart: the zoom anchor is a per-chart x origin the axis and the plot must share, a ctrl-wheel listener on the board would call preventDefault over every other view, and capturing pointers on the board would retarget synthesized clicks everywhere, which REQ-336/337 already paid for.
- **D-11 — the Node stub resolves the board through `getElementById`, never a class selector. DECIDE & STATE.** The stub's `querySelector` returns null unconditionally, so a class lookup would make the renderer hit its guard and return on every Node probe, silently turning several into assertions about nothing.
- **D-12 — the probe scopes "one scroll surface" to the two elements the REQ's RED expression names. DECIDE & STATE.** The table below the chart is a third scroll box inside a `<details>` that is closed by default and is deliberately still a scroll box.
- **D-13 — delete the dead visible-row-range helper in the same commit, if it costs nothing. DECIDE & STATE.** It has no callers; delete-before-you-add applies and this is the render path being touched. Dropped if removing it reddens anything.

**The five scroll-position conversions.** The anchor read and the anchor write both go through a board-relative position rather than the raw board scroll position, or every recorded offset is wrong by 250-400px and the reader is dropped several rows from the REQ they were reading. The grow-the-SVG-before-writing ordering that REQ-319 taught survives the move and gets stronger, because the board's scroll height also depends on the rows SVG. The scroll-a-row-into-view site must convert both directions *and* clear the sticky axis, or a row scrolled to the top lands underneath it and the down arrow appears to do nothing. Jump-to-open-work's zero fallback stops meaning "scroll the board to the very top" and starts meaning "bring the first row just under the axis".

**Testing.** One new browser probe copying the Activity probe's seven shape rules, with a 120-REQ synthetic fixture sized so the chart is more than two screenfuls taller than the fold, plus one REQ engineered to land in the anomalies strip so the padding rule is testable with a visible sibling above the view panel. Four guards run before any comparison — rows drew, every measured box has positive height, the board actually scrolls, the scroll landed — because a probe whose fixture is too short passes while measuring nothing. Seven assertions: the inner box no longer scrolls *and* its computed overflow is visible; the REQ's own RED expression verbatim; the axis pinned with no rows painted above it; the virtualized range following the scroll, proved by disjoint rendered-id sets before and after; the padding scoping in both directions; the strip not flush against the top bar; and the view's own inset surviving. Five existing browser probes and five Node/string probes re-run as regression. The focus ring is deliberately not asserted — `:focus-visible` after a programmatic focus pins the engine's heuristic, so its evidence is a screenshot.

### Plan validation (orchestrator)

- **Requirement coverage: one gap, named rather than dropped.** "The row under the pointer stays put when the window chips change" is the anchor behaviour, and no existing probe has an equivalent. The plan offers a manual check or a ~25-line eighth probe assertion. **Orchestrator decision: pin it in the probe.** It is the requirement most directly threatened by the conversion this REQ is about, and a hand check recorded in a hand-back cannot fail later. It needs its own guard that the anchor id is still in the display list after the window change.
- **No orphan tasks.** Every one of the six tasks traces to a requirement in the What, Context or GREEN condition.
- **Scope: six tasks, over the five-task threshold — flagged, not split.** No subset ships: the CSS alone freezes the virtualized rows and leaves the view worse than before, the JS alone leaves both surfaces scrolling, the Node-stub edit is forced by the JS and silently disarms four test files if skipped, and the probe is the REQ's own GREEN condition. The count reflects files with different failure modes rather than independent deliverables. Recorded as a warning for review, per Step 4.
- **Two escalated decisions carry forward.** D-06 and D-07 are user-visible choices with stated value and risk. The builder proceeds with both defaults; they are confirm-or-override items for the user, not blockers.

## Exploration

Explore agent, read-only, run before the claim; every line number was independently re-checked by the Plan agent against source and all held.

**What REQ-585 actually did**, from its archived record and merged diff (builder commit `94532c45`, merge `c08ac2b4`, two files, 298+/8-): it deleted exactly two declarations — `max-height: 70vh` and `overflow: auto` — from the Activity table's scroll box, added nothing for the sticky header because the `thead th` rule already carried `position: sticky; top: 0`, and moved the board's top padding onto the view's first element under `.board-main:has(> #view-activity.is-active) { padding-top: 0 }`. Its measured RED→GREEN: board 665/677 → 665/3657, table 530/3509 → 3509/3509, header 55.5px from the board's top edge → 0.0px, rows above it 0. The record's own "Worth knowing" names REQ-587 and asks it to reuse that scoping rather than invent a second one.

**The Timeline's current model.** `.timeline-scroll` is `max-height: 58vh; overflow-y: auto; overflow-x: hidden; border-top: 1px solid`, with a block comment above the timeline CSS declaring it the scroll container by design — a sentence this change makes false. `.timeline-axis` has no CSS rule of its own; it is a bare `<div>` inside `.timeline-chart`. The board is `overflow-y: auto; padding: 24px 28px 56px`, dropping to `18px 16px 48px` below the 760px breakpoint.

**Twenty-three sites in `board-timeline.js` reference the scroll host, and the variable serves two distinct roles** — the DOM host that owns the rows SVG, the listeners, focus and pointer capture, and the scroll surface whose position and height drive virtualization. Re-pointing splits them. Five are scroll-position sites: the top-visible-row anchor read, the anchor write in the window refresh, the visible-range computation inside the scroll listener itself, scroll-a-focused-row-into-view, and jump-to-open-work. One is the scroll listener registration. The remaining seventeen are host, width, focus, hover, keyboard, wheel and pointer sites that stay where they are. There is no `scrollIntoView` call anywhere in the file — the row site is a manual position write, which is why it can be re-pointed at all.

**Seven risks the REQ text does not name.** Leaving `overflow-x: hidden` beside an `overflow-y: visible` keeps the element a scroll container, so the fix can look applied and measure RED. The board's scroll position is shared across every view and nothing resets it, which the Timeline is the first virtualized view to care about. The verify-findings strip is visible on the Timeline, unlike Activity, and sits outside the view panel where the padding-move recipe cannot reach it. The sticky axis needs the base background token, not a surface token, because the prime states outright that the surface behind this board's SVG is `<body>`. The chart's inward focus ring stops meaning anything on a full-height element while "Tab to the chart" stays in the shipped hint text. The measured plot width grows by a scrollbar's worth once the inner box stops drawing one. And the table below the chart is a third scroll box, inside a `<details>` closed by default.

**The least obvious cost is the Node lane.** The timeline renderer is driven through a stub DOM whose `getElementById` is a fixed id map across five list literals in four Go test files, and whose `querySelector` returns null unconditionally. The moment the renderer looks up the board, every one of those literals must gain it or the renderer hits its guard and returns — the probes still pass, measuring nothing.

**Test surface.** Browser probes gate on an environment flag before anything else and skip without a browser, which is not a pass. REQ-585's `activity_scroll_browser_probe_test.go` is the file to copy, and its seven shape rules are listed in the Plan. Four existing timeline probes must be re-run as regression: the ResizeObserver/width invariant, the roving-tabindex-over-a-virtualized-list probe, and the two drag-pan and pointer-capture pins.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — the block comment, the timeline-scroll box stripped to two declarations, the first timeline-axis rule, the scoped padding move, the relocated focus ring
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — one board host lookup, the cached geometry pair, the five scroll-position conversions and the listener move
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modify) — the guarded scroll reset on entering the view
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the shared Node stub id map
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modify) — two local stub id lists
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go` (modify) — one local stub id list
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modify) — one local stub id list
- `skills/do-work-board/tools/queue-kanban/timeline_scroll_browser_probe_test.go` (new) — the one-scroll-surface probe with its eight assertions and four guards
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify) — **added to this list by the orchestrator during integration, not declared at Step 5.5.** Four existing browser probes read the old scroll model and went red on the change; three aim by scrolling and measuring in the same task, and one guarded vacuity against a number that is now its own height. See D-22.

**Files I will NOT touch:** `web/template.html` — the markup already has everything the change needs and the chart's `tabindex` is pinned by a Go assertion. The Activity view's CSS block and its probe. `web/board-detail.js` and `web/board-filters.js`. The heavy-lane manifest and any lock-in test. Nothing under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] Exactly one element scrolls on the Timeline view: the REQ's own RED expression returns `[true, false]`, and the inner box's computed `overflow-y` is `visible` rather than merely measuring equal heights
- [ ] The time axis stays visible while scrolling through the rows, pinned to the board's top edge with no rows painted above it
- [ ] The rows still render correctly at every scroll position: a row well below the fold appears after scrolling there, proved by disjoint rendered-row sets before and after
- [ ] The row under the pointer stays put when the window chips change, pinned in the probe rather than checked by hand
- [ ] Wheel-to-zoom, drag-to-pan and the keyboard behaviour named in the view's hint paragraph all still work, evidenced by the four existing probes plus a manual pass in a real engine
- [ ] Measured in a real engine, with the browser and build named, and recorded in the hand-back

## Pre-Flight

**Git:** ⚠ Three tracked files are uncommitted in the shared main tree and none of them is this REQ's: `skills/do-work-board/tools/queue-kanban/model.go`, `skills/do-work-board/tools/queue-kanban/dependency_test.go`, and `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go`. A sibling session in this same checkout is mid-edit, and two of those files sit in this REQ's own package. They are named here, left exactly as they are, and never staged by this REQ. This is the concrete reason the builder works in an isolated worktree rather than the main tree: the worktree is cut from a commit, so a neighbour's uncommitted bytes cannot reach the build, the tests, or the diff this REQ is judged by.

**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exited 0 at revision `762a056f`, run directly and unpiped from the project root with its exit status captured rather than inferred. The green record was written for that exact argv at `d6b6adb1`; the only commit between the two adds files under `ai-reports/`, and `git diff --name-only 762a056f HEAD` outside `do-work/` and `ai-reports/` is empty, so no source the gate measured has moved.

**One gate incident, self-inflicted and repaired.** The run before this one committed REQ-583's focused probe as run evidence without a shebang, which the gate's shellcheck lane reports as SC2148, and that took the canonical gate red. It was found because the gate was invoked through a pipe to `tail`, so the shell reported the pipeline's status and not the gate's — the failure was visible in the output and invisible in the exit code. The probe gained a shebang and `set -euo pipefail` at `762a056f`, shellcheck is clean on it, and every gate invocation from that point on captures the gate's own status. Recorded here rather than quietly fixed, because the pipe is the part worth not repeating.

**Tests baseline:** ✓ Green, but measured twice for a reason. In the shared main tree `go -C skills/do-work-board/tools/queue-kanban test ./...` exited 1 on its first run and 0 on the rerun, and the retry rule makes the rerun the recorded result. That instability is the sibling's uncommitted `model.go` and `dependency_test.go` being read mid-edit, not a real flake in this package. The honest baseline is the one taken off the shared tree: the same command in this REQ's own worktree at `d6b6adb1`, with nothing uncommitted, exits 0 in 42.4s. That is the number a later red in this package is attributable against.

**Dependencies:** ✓ Go toolchain present and the board module builds. Chrome is installed at `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome` but is not on `PATH` under any name the probe looks for, so every browser lane and browser probe needs `QUEUE_KANBAN_BROWSER` pointing at it — without that the lane reports skipped, and a skip is not a pass.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_scroll_browser_probe_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)

**What was done:** The Timeline's rows now scroll against the board instead of against a 58vh box of their own, and the time axis is pinned to the board's top edge the way the Activity view's column header is. The chart box lost its height cap and both overflow declarations, because keeping either one would have left it a scroll container while it measured as though it were not. The axis got its first CSS rule of its own — sticky, painted with the base background token rather than a surface token, and carrying the one-pixel rule the rows box used to own so the separator stays on screen instead of scrolling away with the content. The board's top padding moved onto whichever of its children the reader actually sees first, keyed on that condition rather than on a list of the three optional strips that can precede the view panel.

In the renderer, five sites that read or wrote a scroll position moved into the board's coordinate space: the top-visible-row anchor and the write that restores it after a window change, the virtualized visible-row range, scroll-a-focused-row-into-view, and jump-to-open-work. The two geometry numbers those conversions need are measured once per render and invalidated beside the existing plot-width cache, because the visible-row render is the scroll listener and two extra layout reads would land on every frame of a drag. Only the scroll listener itself moved to the board; wheel-zoom, drag-pan, pointer capture, keyboard, focus, hover and the width observer all stayed on the chart. Entering the view now resets the board's scroll position, which nothing did before and which the Timeline is the first virtualized view to need. The dead visible-row-range helper was deleted and the one comment that named it reworded.

Two defects the plan did not predict were found by measuring rather than reasoning and fixed on top of it: deleting the chart's focus ring handed the same full-height box to the engine's own default ring, which a screenshot caught and an explicit suppression fixed; and clamping the board-relative scroll position at zero made the anchor's write disagree with its read, so pressing any window chip while the board sat above the chart scrolled it down by the whole offset above the chart. The clamp moved to the one call site that wants it.

Merge range `93ec7792..8fad73b2`; builder branch head `564c81a3`. The range is byte-identical to the builder branch's own diff against its base: nine files, 937 insertions, 46 deletions.

## Decisions

D-01 was recorded at Step 3.5 (the shape), D-02 through D-13 by the Plan agent, D-14 through D-21 by the builder in `do-work/runs/work-2026-09-05-170806/REQ-587-handback.md` → `## Decisions`, and D-22 by the orchestrator at integration. The plan's decisions are in `## Plan` above. The builder's and the orchestrator's:

- **D-14 — the chart's own focus ring is explicitly suppressed, not merely relocated. DECIDE & STATE.** Moving the ring to the axis meant deleting the chart's `:focus-visible` rule, which handed the same full-height box straight back to Chrome's default ring — a blue box from the axis down past the last row, exactly the shape the relocation exists to prevent. Every assertion passed while it was there, because a programmatic focus cannot answer a `:focus-visible` question and the rule's evidence is deliberately a screenshot. The screenshot is what found it. Both before and after images are kept.
- **D-15 — the board-relative scroll position stays signed, and the floor moves to the one call site that wants it. DECIDE & STATE.** Clamping it at zero inside the helper made the anchor's write disagree with its read: the anchor records a position minus a row top and the refresh adds it back, so a swallowed negative returns as a positive scroll. Measured consequence: pressing any window chip while the board sat above the chart scrolled it down by the whole offset above the chart — 721px in the fixture, 914px on this repository's own board — slamming the first row against the top bar and hiding the heading, legend and toolbar the reader had just used. Only the visible-range call needs the floor, where a negative position can only widen the range upward.
- **D-16 — four existing browser probes were updated rather than left red. DECIDE & STATE, with the scope widening flagged.** Three of them aim by scrolling and then measuring in the same task, so they were reading rows the previous board position drew; one guarded vacuity against the chart's client height, which is now its own full height, so it compared a number with itself and could never fire again. All four were confirmed green at the base revision by stashing the whole change and re-running, rather than assumed. The behaviour change is intentional, so the Cross-REQ Test-Break Rules say update and document.
- **D-17 — the anchor assertion widens the window rather than narrowing it. DECIDE & STATE.** Every row in the narrow window is also in the wide one, so the display-list guard the orchestrator asked for cannot fire in practice and the comparison always runs; a skip therefore reads as a fixture bug rather than a pass. Narrowing instead would have made a real regression indistinguishable from a legitimately dropped row.
- **D-18 — the anchor fixture cycles ten user requests over 300 days. DECIDE & STATE.** Rows are ordered newest first, so on a naive fixture a wider window only appends older rows below a mid-list anchor, its top never moves, and a forgot-to-convert-both-ways bug preserves the position by accident. Cycling the user request puts old members in every group, so widening inserts rows above the anchor and the board has to move to hold its screen position.
- **D-19 — the fixture's anomaly REQ is kept even though it also raises a data warning. DECIDE & STATE.** The warning was unplanned and makes the fixture strictly better: the warnings banner becomes the board's first visible child, and it is an element with no id that the template does not contain — the exact arrangement a rule naming the strips by id would have got wrong.
- **D-20 — the dead visible-row-range helper was deleted and the comment naming it reworded. DECIDE & STATE.** Deleting the function while leaving a comment naming it would strand a grep target on nothing.
- **D-21 — the probe restates the rows-offset formula rather than reading it. DECIDE & STATE.** A borderline call against the read-a-constant rule: the number is a fact about the laid-out page rather than a decision, it is used only to pick a scroll target and to report, and no assertion depends on it. Said so at the measurement.
- **D-22 — the orchestrator accepted the scope expansion and extended the declared list before integrating. DECIDE & STATE.** `timeline_browser_probe_test.go` was not in the Scope written at Step 5.5, and a builder that needs a file outside its boundary owes a report rather than a silent write; this one committed the change on its own branch and disclosed it in full, which is the unattended form of that report. The judgment is that the alternative — reverting four hunks and handing back a red heavy lane — would deliver a REQ whose own GREEN condition is a browser measurement while the browser lane is red. The `## Scope` list and `write_set` were extended to match rather than leaving the touch undeclared, which is what keeps scope drift measurable instead of merely tolerated.

## Qualification

**Passed.** Read from the merge range `93ec7792..8fad73b2`, and the REQ's own acceptance re-run against the merged state rather than accepted from the hand-back.

- **The merged state was verified, not just the builder's branch.** A detached worktree was cut at the merge commit and the seven browser probes that bear on this change were run there with the engine named explicitly: the new one-scroll-surface probe, the Activity view's equivalent, and the five existing timeline probes. All seven PASS, none SKIP, 31.2s total. That distinction matters here more than usual — a browser probe with no engine reports skipped and exits 0, so an unchecked green would have meant nothing.
- **Substantive, and the range matches the branch exactly.** Nine files, 937 insertions, 46 deletions, and `git diff --name-only` over the merge range is identical to the builder branch's own diff against its base. Eight files are the declared write set; the ninth is the one the builder disclosed and the orchestrator accepted (D-22).
- **One hand-back number was wrong and is corrected here.** The manifest reported "9 files, 320 insertions, 51 deletions". The actual figure, from both the range and the branch, is 937 insertions and 46 deletions — the 640-line new probe file is most of the difference, which suggests the builder read a `git diff --stat` taken before the probe was added. The file list is exactly right; only the counts were stale. Recorded rather than silently corrected, because a manifest number that does not match the diff is precisely what this step exists to catch.
- **The two defects the plan missed were found by measuring, which is the point.** The focus-ring defect passed every assertion — a programmatic focus cannot answer a `:focus-visible` question, so the rule's evidence is deliberately a screenshot, and it was the screenshot that caught the engine's own ring being handed the same full-height box the relocation exists to avoid. The clamp defect was found by pressing a window chip and watching the board move by the whole offset above the chart. Neither would have been caught by reasoning about the diff.
- **The Node lane's numbers are the strongest evidence the conversions are right.** Every row count and row id in that lane came back byte-identical with no expectation edited. That lane drives the renderer through a stub DOM, so identical output after the scroll host was split in two says the coordinate conversion is a no-op where it should be.
- **The one warning is the documented exception.** `QUALIFY-NEW-FILE-UNWIRED` fires on the new probe file because nothing statically references it. A Go `_test.go` file is discovered by the toolchain, not by a reference, and this is exactly the entry-point case the check hands to judgment. Not a finding.
- **Two scope-drift warnings were prose, not drift.** The checker read the backticked CSS selectors in the orchestrator's own Scope bullet as declared paths. The wording was fixed; no declaration and no code changed. Same class as the one seen on REQ-583 this run.

Requirements traced: one scroll surface, proved by the REQ's own expression and by the computed overflow value the height comparison cannot see; the axis pinned with no rows above it; the virtualized range following the board, proved by disjoint rendered-id sets; the anchor holding its screen position across a window change, pinned in the probe as the orchestrator directed rather than checked by hand; wheel-zoom, drag-pan and the keyboard path exercised by hand in a real engine and by four existing probes; and every number measured on a named browser build.

*Checked by work action*
## Testing

**Tests run:** `go -C skills/do-work-board/tools/queue-kanban test -count=1 ./...`
**Result:** ✓ All passing. The builder measured 62.4s and 48.9s against its own clean baseline of 54.3s taken in the same worktree minutes before the first edit; the independent reviewer measured 42.6s in a separate worktree at the merge commit. The spread is machine load from sibling sessions, not this change: the new probe is gated behind the browser-probe flag and contributes nothing to this lane.

**Browser lane, the REQ's own GREEN condition.** Run at the merge revision `8fad73b2` in a detached worktree, with the engine named explicitly — Chrome is installed but not on `PATH` under any name the probe looks for, and without that variable every browser probe reports skipped, which is not a pass. The output was grepped for skip lines rather than trusted to the exit code.

- Orchestrator's post-merge verification, seven probes that bear on this change: all PASS, none SKIP, 31.2s. The new `TestBrowserBehaviorTimelineViewHasOneScrollSurface`, the Activity view's equivalent, and the five existing timeline probes.
- Reviewer's independent run, the full lane: 65 probes, 65 PASS, 0 SKIP, 131.8s.

**Red-green validation** (traced to `## Red-Green Proof`), measured on Chrome 152.0.0.0 at 1600×900 on a rebuilt binary, with the pre-change build served side by side:

- The REQ's RED expression, verbatim on the served board: `[true, true]` before, **`[true, false]`** after.
- `#timeline-scroll` client/scroll: 521/14318 before, **14318/14318** after; computed `overflow-y` `auto` before, **`visible`** after; computed `max-height` `522px` before, `none` after. The computed-overflow row is the one the height comparison cannot see, which is why the probe asserts it.
- Axis pinned: top offset from the board's inner top edge **0.0px**, rows painted above it **0**, after scrolling 1500px of rows in.
- Virtualized range follows the board: 39 rendered row ids before the scroll, 44 after, **intersection empty**, 34 of them intersecting the visible rect afterwards.
- The anchor across a window change: display list 62 → 201 rows, board scroll 1275.0 → 1960.0 to follow it, anchor row back at **0.5px** of drift.
- Padding move verified at 320 / 768 / 1280 / 1600px, and in both directions — zero on the Timeline, restored on the Kanban view.
- Both themes: the axis background equals the body background exactly, light and dark. No tinted band over the untinted chart.

**New tests added:** `TestBrowserBehaviorTimelineViewHasOneScrollSurface` — four guards before any comparison (rows drew, every measured box has positive height, the board actually scrolls, the scroll landed) and eight assertions, including the anchor case the orchestrator directed be pinned rather than checked by hand.

**Existing tests updated (cross-REQ impact):** `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` — four probes that read the old scroll model. Three aimed by scrolling and measuring in the same task and so were reading rows the previous board position drew; one guarded vacuity against the chart's client height, which is now its own full height. All four were confirmed green at the base revision before being changed. Intentional behaviour change, updated under the Cross-REQ Test-Break Rules; see D-16 and D-22.

**Canonical repository gate — red five times on machine load, then green twice.** This is recorded in full because the conclusion rests on it.

Five consecutive runs of `bash _dev/tests/maintainer-verify.sh` from the project root exited 1, every one of them **only** on the per-test-file wall-clock budget, and on four different files across the runs: `internal/finalization/finalization_recovery_test.go` (43.15s, then 36.52s, then 33.27s), `internal/finalization/finalization_req499_test.go` (35.21s, then 37.93s), `internal/publication/defer_gate_test.go` (31.63s, then 34.68s) and `_dev/tests/session-start-hook-behavior.sh` (40s). No correctness check failed in any run. None of those files is in this REQ's diff, which is confined to the board module.

Three pieces of evidence said machine, not tree. The same two Go files ran 18.59-25.02s in four earlier gate runs today on the same code. The shell probe that took 40s inside the gate takes **14.9s** run alone. And the load average during the failing runs was 22-40, with several sibling Claude sessions running Go suites and browser probes in this same checkout — the project's own recorded lesson for this is "one gate per machine".

The prescribed diagnostic settled it. A detached worktree at the saved pre-merge base `93ec7792` ran the gate directly: **exit 0**. A second detached worktree at the merge commit `8fad73b2` ran it back to back in the same quiet window: **exit 0**. Base green and merge green under identical conditions is the answer the attribution procedure asks for, and it rules out the current implementation as the cause.

The recorded green is a third run, at the current integration tip `5fc2f66a`, in an isolated detached worktree so a sibling's uncommitted edits in this very package could not reach it: **exit 0**, load 5-7. That is the revision the green record was written for, and `git diff --name-only` between the merge and it shows one gate-test file a sibling session changed, which is why the merge-revision green was not reused.

**A discovered task falls out of this and is reported, not fixed here:** those three Go test files have no headroom against the 30s budget. Measured alone at load ~13, `internal/finalization`'s top-level tests total 48.9s across the package and `internal/publication`'s total 25.0s, with single tests at 7.8s and 9.6s. They pass on a quiet machine and fail whenever it is busy, which makes the gate's verdict a function of who else is working. UR-127 was captured by another session while this run was in progress and covers the same ground.

**Heavy verification plan:** *(selected lanes are recorded below in `## Heavy Verification Plan`; held for the drain at queue exhaustion)*

*Verified by work action*

## Review

**Overall: 92%** | 2026-09-05T19:43:13Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 82% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- The D-15 fix has no regression test: restoring the plan's clamp in `rowsScrollTop` (`web/board-timeline.js:2082`) leaves the whole browser lane and the Node lane green, while the reintroduced defect measurably scrolls the board 0 → 636px and drags the view heading from +443.8px to -192.2px on a window-chip press from above the chart. The new probe's anchor case always starts from a positive rows-scroll position, so the negative path is never exercised — impact-user-visible → report only

**Minor findings:**
- D-16's recorded reason for the vacuity-guard hunk is inverted — the guard would always fire, not "never fire again"; reverting that hunk alone produces a hard failure at `timeline_browser_probe_test.go:3784`. The fix is correct; the written justification is not — impact-negligible → report only
- Restatement Sweep: `web/board-timeline.js:150` still calls `#timeline-scroll` "the scroll container" in the present tense, in the one file where that distinction is the point of the REQ — impact-negligible → report only
- Restatement Sweep, outside this REQ's Scope: `_dev/primes/lessons-kanban-board.md:71` states the padding recipe as "the view's first element", the narrower form REQ-587 proved insufficient once strips can precede the view panel; the satellite is the canonical home and should carry the generalised rule — impact-rule-change → report only
- Nit: the "no strips at all" arrangement of D-05 is verified by construction but pinned by no assertion — impact-negligible → report only
- Nit: `.timeline-scroll:focus-visible { outline: none }` is unconditional while its replacement ring depends on `:has()`, so an engine without `:has()` gets no focus indication; the dependency predates this REQ and `:has()` is baseline for this tool's engines — impact-negligible → report only

**Acceptance:** Pass — 65 browser probes run at the merge commit in a detached worktree with the engine named, 65 PASS, 0 SKIP; full package `ok` in 42.6s; vet, gofmt and build clean; the four probes D-16 updated confirmed green at base revision 93ec7792.
**Suggested testing:** 5 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Measuring instead of reasoning, at three separate points where reasoning would have shipped a defect. The plan was verified against source before a line was written and held everywhere it was checked, but two of its predictions were still wrong, and both were caught by looking at the page rather than at the diff. Deleting the chart's focus ring to move it onto the axis handed the same full-height box straight back to the engine's own default ring — every assertion passed while that was live, because a programmatic focus cannot answer a `:focus-visible` question, and it was the screenshot that found it. Clamping the board-relative scroll position at zero, which the plan asked for, made the anchor's write disagree with its read; it took pressing a window chip and watching the board move by the whole offset above the chart to see it.

**What didn't:** Two things that read as safe and were not. `overflow-x: hidden` left beside an `overflow-y: visible` keeps an element a scroll container — the CSS Overflow rule promotes the visible axis to `auto` — so the half-applied fix would have measured `scrollHeight == clientHeight` and looked finished. That is why the probe asserts the computed value and not only the heights. And the padding-move recipe this REQ was told to reuse names "the view's first element", which is right on the Activity view and insufficient here: three optional strips sit outside the view panel and can precede it, and on the probe fixture the first visible child turned out to be a data-warnings banner that is not in the template at all. A rule naming the strips by id would have been wrong the first time the queue's data changed.

**Worth knowing:** The fix for the clamp defect has no regression test, and the reviewer proved it — restoring the plan's clamp leaves the whole browser lane and the Node lane green while measurably scrolling the board 636px and dragging the view heading off screen on a window-chip press from above the chart. The new probe's anchor case always starts from a positive rows-scroll position, so the negative path is never exercised. That is the same shape as the three findings UR-119 was raised for: the code is right and nothing holds it there. It is reported rather than queued, because the impact rule only auto-queues critical findings, and it is the single item from this REQ most worth promoting by hand.

## Orientation

The Timeline reads like every other view now: the board scrolls, the chart is ordinary content inside it, and the date axis stays pinned to the top edge while the bars pass under it — so a reader with a pointer over the chart no longer gets a different scrollbar than a reader with a pointer beside it. Lives in the queue-kanban board's timeline view, governed by `_dev/primes/prime-kanban-board.md`. [MAP CHANGED] — which element is the Timeline's scroll surface has moved from the chart box to the board, and the virtualized row range is now computed from the board's position minus the chart's offset inside it, which is a contract every future change to that renderer inherits. The prime and its lesson satellite were spot-checked against the change: the prime is current, and the satellite's padding-move entry now states the narrower form this REQ proved insufficient, which is carried as a deferred lesson write.

## Heavy Verification Plan

- **Base revision:** `93ec779285c4c2f61ba9d3e9db3e0e88d05e5e28`
- **Target revision:** `8fad73b20055bbe66df91d423027867d780f3175` (the recorded `commit:`)
- **Changed paths in range:** all nine files of this REQ's diff, every one inside `skills/do-work-board/tools/queue-kanban`. No uncovered paths, planner not forced, not uncertain.

| Lane | Argv | Selection reason |
| --- | --- | --- |
| `queue-kanban-javascript` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` | every changed path matched subtree `skills/do-work-board/tools/queue-kanban` |
| `queue-kanban-browser` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` | same subtree |
| `staged-skills` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` | every changed path matched subtree `skills` |

Held at Step 7.7: the lanes are not run now and the queue loop is not held open. **The browser lane needs `QUEUE_KANBAN_BROWSER` pointing at the installed Chrome at drain time** — this machine has Chrome but not under any name the probe looks for, so without it the lane reports skipped, and a skip is not a pass. The recorded `commit:` above is what makes this REQ's source ready for anything that depends on it while it waits.

### Heavy verification result (run at drain, 2026-09-06)

**All six lanes ran and all six were green. Revision `a48b9eb6`, the tree quiet from the first command
to the last.** `bash _dev/tests/maintainer-verify.sh --heavy` printed `Maintainer verification passed.`
and exited **0**, gate wall **301s**.

**One deviation from the plan, stated.** The plan named six separate `--heavy-lane <id>` invocations.
What ran instead is the single `--heavy` gate, which executes the same six lanes' work in one process
at one revision. That is what the four held requests were waiting for — one run at the final revision
rather than four — and it removes the `HEAVY-RUN-REVISION-CHANGED` risk of interleaving six
invocations with four finalizations. The evidence below is per lane, not the gate's summary line,
because a skipped lane reports success.

| Lane | Its own evidence line | Result |
|---|---|---|
| `queue-kanban-javascript` | `module=…/queue-kanban wall=67s tests=481 slowest-file=generate_test.go:12.43s limit=none (heavy)` | 481 tests, green |
| `queue-kanban-browser` | `module=…/queue-kanban wall=102s tests=35 slowest-file=timeline_browser_probe_test.go:63.99s limit=none (heavy)` | 35 tests, green |
| `do-work-cli-integrations` | `module=…/do-work-cli wall=25s tests=798 slowest-file=internal/nextselection/blocked_probe_test.go:6.77s limit=none (heavy)` | 798 tests, green |
| `staged-skills` | `test-file duration: staged-skills-contract.sh 45s (limit none (heavy))` | green |
| `updater` | `test-file duration: update-script-behavior.sh 84s (limit none (heavy))` | green |
| `installer` | `test-file duration: install-suite-behavior.sh 28s (limit none (heavy))` | green |

**Zero `SKIP` lines in the whole run and zero `FAIL` lines.** The browser lane genuinely ran — 35 tests
and a 64-second `timeline_browser_probe_test.go` are not what a skipped lane prints — because
`QUEUE_KANBAN_BROWSER` pointed at `/opt/pw-browsers/chromium`, as every one of these four plans
required.

The run also needed a sanitized environment, which is worth recording for the next drain:
`NODE_OPTIONS` and the `GIT_CONFIG_COUNT` / `GIT_CONFIG_KEY_*` / `GIT_CONFIG_VALUE_*` triples unset,
and `GIT_CONFIG_GLOBAL` pointed at a config with `commit.gpgsign = false`. A heavy run refuses on an
opaque runtime extension or an opaque git configuration override, and an unusable global signing key
makes a fixture's own `git commit` fail inside the lane.
