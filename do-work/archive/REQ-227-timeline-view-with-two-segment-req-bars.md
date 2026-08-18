---
id: REQ-227
title: Add the Timeline view with two-segment REQ bars
status: completed
completed_at: 2026-08-18T01:19:09Z
commit: 17b9422
claimed_at: 2026-08-18T00:58:19Z
created_at: 2026-08-17T23:51:17Z
user_request: UR-051
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-226, REQ-228]
batch: board-timing-views
route: C
estimate:
  p50_active_minutes: 65
  confidence: low
  calculated_at: 2026-08-18T01:00:00Z
  basis:
    - Route C
    - 8-file write set
    - 3 new files
    - 3 subsystems involved
    - 9 acceptance criteria
    - browser evidence
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board-filters.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board.js
  - skills/do-work-board/tools/queue-kanban/durations_test.go
---

# Add the Timeline View with Two-Segment REQ Bars

## What

Add a fifth board view, **Timeline** — a zoomable, scrollable Gantt with one horizontal bar per REQ.
Each bar carries two coloured segments: `created_at`→`claimed_at` (the wait) and
`claimed_at`→`completed_at` (the work). REQs that are claimed but not finished draw as open bars
running to now.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`, the ticket model, the payload projection, and all five existing view fragments. Approach in `## Plan` below.
- [x] **[APPLY]:** Eleven files. Nine as declared; two extensions recorded as D-03 and D-04 and mirrored into `write_set` before the work landed.
- [x] **[UNIFY]:** `git diff --stat` = 11 files, +1374/−10. `gofmt -l .` clean, `go vet ./...` clean, `node --check` on the assembled client clean, zero page errors in every headless run. Grep for `console.log|fmt.Print|debugger|TODO|FIXME|XXX|alert(` in added lines: none. Files verified:
  - `timeline.go` — signed spans throughout, no clamping; the anomaly fields are copied and never recomputed.
  - `timeline_test.go` / `generate_test.go` — four Go tests and two Node probes; both probes mutation-checked.
  - `generate.go` — additive wire types plus the fragment-manifest entry; the manifest and inventory assertions changed in the same commit, as the REQ requires.
  - `web/board-timeline.js` — the scaffold placeholder is gone (grepped: 0 hits); every listener goes through the teardown registry (grepped `addEventListener`: 1 hit, inside the helper).
  - `web/board-controls.js`, `web/board-filters.js`, `web/board.js` — timeline added to the panel map, the lazy-render guard, the stale set, and the guard object; no existing view's wiring altered.
  - `web/template.html`, `web/board.css` — new tab, panel, and styles; both theme blocks carry the new tokens.
  - `CLAUDE.md` — untouched, and § Kanban Board Write Surfaces still reads "exactly three".

## Why

"Also make another one that is a zoomable and scrollable gant chart ... gant chart contains also
duration, which can be calculated from the captured time, start time and finish time (one bar with 2
colors) distance between captured and started is one, and started and finished there is another one."

## Context

**The data is already there.** `created_at`, `claimed_at`, and `completed_at` are parsed onto every
ticket as raw frontmatter text (`model.go:76-78`, assigned at `model.go:717-719`) and every ticket —
including pending and claimed ones — is already emitted into the client payload with those stamps
(`generate.go:440-487`, wire fields at `generate.go:116-182`). The board cards already read
`createdAt`/`claimedAt` client-side for the live state timer (`web/board-core.js:190-212`). So this
REQ needs **no new frontmatter field and no second walk** — only a new aggregation and a new view.

**What the existing time views deliberately do not cover.** The Durations view's aggregate skips any
ticket that is not terminal-success or is missing either stamp (`durations.go:78-85`), and it measures
only `completed_at − claimed_at`. Two consequences this REQ fixes:

- The **wait** — `created_at`→`claimed_at` — is measured nowhere on the board. On a queue with real
  backlog it is very likely the dominant share of calendar time, which makes this view a queue-health
  chart at least as much as a duration chart. Expect the finding; design the colour weighting so it
  is legible rather than surprising.
- **In-flight REQs are invisible everywhere.** A REQ claimed and running right now contributes to no
  panel of any view. Drawing it here is what later lets REQ-228's forecast start from something real
  instead of from an empty present.

**The tab mechanism is a settled pattern.** Five files, all small: a `<button data-view-target>` in
`web/template.html:70-83`, an entry in the `viewPanels` map and the lazy-render guard in
`web/board-controls.js:12-50`, a `<section id="view-timeline" class="view-panel" hidden>` beside the
existing four (`template.html:175-222`), the fragment registered in `boardJavaScriptFragmentPaths`
(`generate.go:42-52`), and CSS. `applyView()` also hides control groups that do not apply to the
active view (`board-controls.js:28-30`) — decide what Timeline shows there rather than inheriting
Board's knobs by accident.

**There is no zoom or pan anywhere in the client today.** Grepping `web/*.js` for
wheel/pointerdown/drag/zoom yields only the detail-panel resizer (`web/board-detail.js:487-543`).
Durations has `overflow-x: auto` with a `min-width: 720px` SVG (`board.css:1841-1851`), which is a
scrollbar below 720px and nothing on a normal screen. This REQ builds the first real zoom/pan on the
surface, so it also sets the pattern.

**Scale is the hard part.** The reporting board has 560 archived REQs. At any readable row height
that is several thousand pixels of vertical extent, so vertical scrolling is mandatory and row
virtualization or aggregation may be. The client is framework-free plain JS by design
(REQ-195 modularized it that way); no library is coming to solve this.

## Detailed Requirements

1. **One row per REQ, one bar, two segments.** Segment one spans `created_at`→`claimed_at`; segment
   two spans `claimed_at`→`completed_at`. Distinct colours, both keyed in a legend that names them in
   plain words (waiting / working), not as field names.
2. **In-flight REQs draw as open bars.** A REQ with `claimed_at` and no `completed_at` renders its
   wait segment plus a work segment running to now, visually marked as still open rather than
   finished.
3. **Pending REQs are visible as waiting.** A REQ with only `created_at` has an open wait segment
   running to now. It has no projected work segment — that is REQ-228's job, and this REQ must not
   invent one.
4. **Zoom and scroll.** The time axis zooms (in to at least single-day resolution, out to the whole
   history) and pans. Vertical scrolling covers the row set. Zoom state must not reset when the user
   switches tabs and comes back.
5. **Usable at 560+ rows.** State the approach — virtualize, aggregate, page, or a row height that
   makes plain scrolling fine — and demonstrate it against a fixture at that scale, not against this
   repo's ~230.
6. **Row ordering is stated on the view.** Chronological by `created_at` is the recommended default.
   Whatever is chosen, the reader can tell what the vertical order means without guessing.
7. **Broken stamps are shown, not swallowed.** A REQ whose `completed_at` precedes its `claimed_at`
   is the reversed-stamp anomaly class the board already surfaces (`model.go:1206-1241`, and REQ-213
   before it). Render it as visibly broken; never clamp it to zero and never drop the row.
   Consume the existing verdict rather than re-deriving the rule.
8. **"Now" has one source and is labelled.** A live board's now is the request instant; a static
   snapshot's now is frozen at generation time (the header already prints "Generated … · 21s ago").
   Draw a now-line, derive it from one place, and make a stale snapshot self-evident rather than
   quietly wrong.
9. **Hover/selection readout.** Reuse the Durations pattern (`board-durations.js:397-490`): a readout
   line naming the REQ id, title, route, both segment durations, and status. The detail drawer already
   opens from any element carrying `data-detail-kind` (`board-controls.js:149-155`) — wire rows into
   it rather than inventing a second detail path.

## Constraints

- Read-only. No new board write surface — `CLAUDE.md` § *Kanban Board Write Surfaces* must still read
  "exactly three" and go unamended when this lands.
- No new frontmatter field, no second walk of the archive. Derive from tickets already parsed, the
  way `durations.go:73-75` does.
- `boardJavaScriptFragmentPaths` (`generate.go:42-52`) and the inventory assertion
  (`generate_test.go:28-62`) must change in the same commit as the new `web/board-timeline.js`, or
  the build fails. Fragment execution order is that literal array.
- Parser lock-step (`_dev/primes/prime-kanban-board.md:13`): the board buckets by the `status`
  vocabulary in `skills/do-work/actions/work-reference.md`; anything new stays aligned with it.
- Framework-free plain-JS client, both light and dark themes, and both delivery modes — the live
  server (`serve.go`) and the static snapshot (`generate.go`) serve the same embedded assets and both
  must work.
- Durations stays exactly as it is apart from REQ-226's fixes. This view is additive.

## Dependencies

None inbound. REQ-228 depends on this one — it adds projected bars to the view this REQ creates, and
shares `timeline.go` and `web/board-timeline.js` with it. Landing REQ-227 alone leaves a working,
useful Gantt of real history; that is the intended slice boundary.

## Builder Guidance

**Certainty: Firm on what, Mixed on how.** The two-segment bar, the fifth-tab placement, and
zoom-plus-scroll were all confirmed with the maintainer before capture. The rendering strategy,
zoom interaction model, row virtualization approach, and colour assignment are yours.

Worth deciding early and stating in the plan: SVG (consistent with Durations and Calendar, but 560
rows × 2 rects is a lot of nodes) versus canvas (scales, but loses the existing hover/`data-detail-kind`
plumbing and the accessibility story). The existing views' table-fallback convention —
"every value the chart shows is reachable without a pointer" (`board-durations.js:492-493`) — applies
here too and may steer that choice.

Do not let the zoom implementation quietly break the pointer-to-data mapping: Durations converts
pointer coordinates with a fixed `DURATIONS_VIEW_WIDTH / bounds.width` scale
(`board-durations.js:476-483`), which is exactly the assumption a zoom invalidates. Whatever
Timeline does here needs to be viewBox-aware from the start.

## Open Questions

- [~] Should the topbar's domain/status/search filters apply to the Timeline?
  Recommended: yes — filtering a Gantt by domain is genuinely useful, and `onFiltersChanged`
  (`web/board-filters.js:147-166`) already resets `renderedOnce` per view; add `timeline` to that
  reset. Note that Durations is *deliberately* excluded from filtering, so this is a real divergence
  to make on purpose rather than by copying the neighbour.
  Also: leave it unfiltered like Durations; or filter but show the filtered-out count.
- [~] Should archived-and-cancelled REQs appear?
  Recommended: yes, visually distinct — a cancelled REQ still consumed queue time, and hiding it
  makes the wait segments lie by omission.
  Also: hide them; or hide behind a toggle.

Both are deferred to the builder deliberately: they are judgment calls the codebase can inform, and
neither changes the shape of the view.

## Red-Green Proof

**RED prompt/case:** Build a board fixture with four REQs — one completed with all three stamps, one
claimed-not-completed, one pending with only `created_at`, and one with `completed_at` earlier than
`claimed_at` — then ask the board for its timeline payload. Assert the completed REQ yields both a
wait span and a work span matching the stamps, the claimed REQ yields a wait span plus an open work
span, the pending REQ yields a wait span and no work span, and the reversed one is flagged rather
than clamped or dropped. Today there is no timeline aggregation at all, so there is nothing to ask
and the test cannot compile.

**Why RED now:** The only duration aggregation on the board is `buildDurationAggregate`
(`durations.go:75-103`), which skips every non-terminal ticket and measures only
`completed_at − claimed_at`. Neither the wait segment nor any in-flight span exists anywhere in the
codebase.

**GREEN when:** That fixture produces the four expected bar shapes, and opening the board shows a
Timeline tab whose rows render both segments, whose axis zooms and pans, and which stays legible and
responsive on a 560-row fixture.

**Validation:** User confirmed — the two-segment bar and its exact segment boundaries were specified
by the maintainer verbatim; the fifth-tab placement was chosen from three offered options; in-flight
bars were surfaced by capture in answer to "am I missing anything?" and accepted.

## Full Context

See `do-work/user-requests/UR-051/input.md` for complete verbatim input.

---
*Source: "Also make another one that is a zoomable and scrollable gant chart ... (one bar with 2 colors) distance between captured and started is one, and started and finished there is another one."*

---

## Triage

**Route: C** - Complex

**Reasoning:** A fifth board view spanning a new Go aggregation, new payload wire types, a new JS fragment with the client's first zoom/pan implementation, new markup, and CSS — with nine numbered requirements, a scale target above anything the repo's own data reaches, and two deferred Open Questions.

**Planning:** Required

## Plan

### Where the answers already are

Every stamp this view needs is already parsed and already shipped: `CreatedAt`/`ClaimedAt`/`CompletedAt` on the ticket (`model.go:76-78`), projected verbatim to the client at `generate.go:161-163`. The reversed-stamp verdict is already computed once, by `detectCompletionAnomaly` (`model.go:1221-1241`), and lands on the ticket as `CompletionAnomaly`/`CompletionAnomalyReason`. So this REQ adds an aggregation and a view, and re-derives nothing.

### 1. `timeline.go` — the aggregation

`buildTimelineAggregate(tickets, now)` yields one `TimelineRow` per ticket with a parseable `created_at`, carrying both spans as **signed** minutes plus the flags that say what kind of span each is:

- `WaitMinutes` = `claimed_at − created_at`, or `now − created_at` with `WaitOpen` when never claimed.
- `WorkMinutes` = `completed_at − claimed_at`, or `now − claimed_at` with `WorkOpen` when in flight. `HasWork` is false for a REQ that was never claimed — R3 forbids inventing a work segment for it.
- `Anomaly`/`AnomalyReason` are **copied from the ticket**, not recomputed (R7).

Signed, never clamped. A negative span is drawn as a broken marker rather than a rect, which is a fact about SVG (a rect has no negative width) rather than a second reading of the anomaly rule.

### 2. `now` has one source (R8)

Open spans are measured against `board.GeneratedAt`, and the same instant is shipped as the payload's `now`, which is where the client draws the now-line. A live board regenerates per request so its now is the request instant; a snapshot's now is frozen at generation, and the header's existing "Generated … · N ago" is what makes that visible. One value, computed once, used by both the span arithmetic and the line.

### 3. Payload

`generatedTimeline` carries rows plus the range and `now`. Rows carry timing only — id, the three stamps, the two spans and their flags — because the client already holds title, route, status and domain in `requestsById`. Duplicating them would create a second copy to keep in step.

### 4. `web/board-timeline.js` — pixel-space SVG with virtualized rows

- **Pixel coordinates, no `viewBox`.** The REQ warns that Durations' fixed `DURATIONS_VIEW_WIDTH / bounds.width` conversion is exactly what a zoom invalidates. Drawing in CSS pixels with no viewBox scaling removes the conversion instead of maintaining it: pointer-to-data is `clientX − rect.left`, at every zoom level, permanently.
- **Row virtualization (R5).** A scroll container of fixed height holds a full-height SVG; only the rows inside the scrolled window get nodes, rebuilt on scroll. Node count is bounded by the viewport, not by the row count, so 560 rows and 5600 cost the same. Row height is fixed, which is what makes the visible slice a division rather than a measurement.
- **Zoom and pan (R4).** A `timelineViewState` object outside the `renderedOnce` guard holds the visible time window, so switching tabs and back preserves it. Zoom in/out/fit buttons, ⌘/Ctrl+wheel zoom anchored at the pointer, and drag-to-pan.
- **Reuses the existing plumbing (R9).** Rows carry `data-detail-kind="request"`, so the delegated listener at `board-controls.js:148-155` opens the drawer with no second detail path. A readout line mirrors Durations, and a `<details>` table carries every value without a pointer.

### 5. Wiring

A `<button data-view-target="timeline">` in the topbar, a `view-timeline` panel, the `viewPanels` map plus the lazy-render guard in `applyView`, `timeline` added to `onFiltersChanged`'s stale set, and the fragment registered in `boardJavaScriptFragmentPaths` with its inventory assertion updated in the same commit.

### 6. Open Questions, decided

Both deferred items are answered in `## Decisions` (D-01, D-02): filters apply, and cancelled REQs are shown with distinct styling.

### Verification steps

1. Go RED in `timeline_test.go` → verify: the four-REQ fixture yields the four bar shapes (both spans, wait-plus-open-work, wait-only, flagged-reversed).
2. `go test -count=1 ./...` → verify: the fragment inventory and manifest assertions accept the new file, and nothing else regressed.
3. JS behavior probes under Node → verify: the zoom transform is its own inverse for pointer mapping, and the visible-row slice is correct at 560 rows.
4. Headless Chromium against a 560-row fixture → verify: the view renders, zooms, pans, and stays responsive at the stated scale.
5. `bash _dev/tests/maintainer-verify.sh` → verify: exit 0.

*Generated by Plan agent (inline, serial mode)*

## Exploration

**Ticket model.** `RequestTicket` (`model.go:76-78`) holds the three stamps as raw frontmatter text; `parseTimestamp` is the shared reader. `CompletionAnomaly`/`CompletionAnomalyReason` (`model.go:196-197`) are set by `buildBoard` from `detectCompletionAnomaly`, whose reversed-span branch is `model.go:1234-1240`. `isTerminalResolvedStatus` and `normalizeStatus` are the status vocabulary's readers.

**Payload.** `generatedBoardData` (`generate.go:69`) carries `GeneratedAt`; `generatedRequest` (`generate.go:161-163`) already ships `createdAt`/`claimedAt`/`completedAt` to every card, and `generatedCalendarEntry` (`generate.go:210-215`) is the precedent for a timing-only per-view record keyed by id.

**Client shell.** `web/board.js` declares `boardData`, `requestsById`, `generatedAtMs`, `viewState`, `filterState`, and `renderedOnce`, then splices the fragments at `/* INLINE_BOARD_FRAGMENTS */`. Fragments therefore share one closure — a new fragment can read `requestsById` and `generatedAtMs` directly, and must not redeclare them.

**View switching.** `applyView` (`board-controls.js:12-50`) owns the `viewPanels` map, the per-view control-group visibility, and the lazy-render guards. `onFiltersChanged` (`board-filters.js:147-166`) marks the lazily-rendered views stale. `requestMatchesFilters` (`board-filters.js:42-59`) is the shared predicate. Delegated detail opening is `board-controls.js:148-155`.

**No zoom exists yet.** Grepping `web/*.js` for wheel/pointerdown/drag/zoom finds only the detail-panel resizer (`board-detail.js:487-543`), confirming the REQ's note. Durations' `min-width: 720px` SVG under `overflow-x: auto` is a horizontal scrollbar on narrow screens, not a zoom.

**Scale.** This repository's board carries 232 REQs; the REQ's target is 560+. A synthetic fixture is required for both the scale test and the render check, exactly as REQ-226 needed one.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/timeline.go` (new) — the aggregation, signed spans, anomaly passthrough.
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (new) — the four-shape RED plus range and ordering tests.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — wire types, projection, fragment manifest entry.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — updated inventory/manifest expectations, JS behavior probes.
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (new) — the view.
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modify) — panel map, control visibility, lazy-render guard.
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` (modify) — add `timeline` to the stale set.
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — tab button and view panel.
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — the view's styles and colour tokens.
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modify) — the `renderedOnce` guard object and the `viewState.view` vocabulary comment (D-03).
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — generalize the renderer-constant reader the timeline probes reuse (D-04).

**Files I will NOT touch:** `durations.go` and `web/board-durations.js` (this view is additive; Durations stays as REQ-226 left it), `model.go` (no new parsed field and no re-derived anomaly rule), `CLAUDE.md` (no new write surface).

**Acceptance criteria (restated from REQ):**
- [ ] One row per REQ, one bar, two segments, both named in plain words in a legend.
- [ ] In-flight REQs draw as open bars running to now, marked as still open.
- [ ] Pending REQs show an open wait segment and no work segment.
- [ ] The time axis zooms and pans, rows scroll vertically, and zoom survives a tab switch.
- [ ] Usable at 560+ rows, demonstrated against a fixture at that scale.
- [ ] Row ordering is stated on the view.
- [ ] A reversed stamp renders as visibly broken, never clamped and never dropped, consuming the existing verdict.
- [ ] "Now" has one source, is drawn, and a stale snapshot is self-evident.
- [ ] A hover readout names id, title, route, both durations, and status; rows open the existing detail drawer.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/timeline.go` (new)
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (new)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)

**What was done:** Added a fifth board view. `buildTimelineAggregate` derives one row per REQ from stamps the board already parsed, carrying both spans as signed minutes with flags for which are still running, and copying the reversed-stamp verdict off the ticket rather than recomputing it. Open spans are measured against the board's generation instant, which is also shipped as the payload's `now` and is where the client draws the now-line, so the arithmetic and the line cannot disagree. The client fragment draws in CSS pixels with no `viewBox`, which removes the pointer-to-data conversion a zoom would otherwise invalidate, and virtualizes rows so node count follows the viewport rather than the archive. Zoom, pan and vertical scroll are wired, with the visible window held outside the render guard so it survives a tab switch. Rows carry `data-detail-kind` and open the existing drawer; a readout line and a full table carry every value without a pointer.

## Decisions

- **D-01**: The topbar's domain/status/search filters apply to the Timeline, diverging from Durations, which is deliberately unfiltered. Reasoning: a Gantt narrowed to one domain answers a straightforward question about a queue, whereas a durations distribution narrowed the same way is a different statistic wearing the same axes. This was Open Question one; the REQ recommended yes and named the divergence as one to make on purpose. `timeline` was added to `onFiltersChanged`'s stale set as that recommendation specified. ESCALATE. **Value:** the view is useful per-domain and per-search on a 560-REQ board where the unfiltered chart is a wall. **Risk:** a reader who has learned that Durations ignores filters may expect the same here; the divergence is stated in a code comment at both sites. Reversible in one line.
- **D-02**: Cancelled REQs are shown, dimmed rather than hidden. Reasoning: a cancelled REQ still occupied queue time, and hiding it would make the surrounding wait segments lie by omission. This was Open Question two; the REQ recommended exactly this. ESCALATE. **Value:** the wait picture stays honest. **Risk:** a queue with many cancellations gets visually busier; the dimming is one CSS rule to change. Reversible.
- **D-03**: Extended the write set to `web/board.js`. The `renderedOnce` guard object and the `viewState.view` vocabulary comment live in the client shell, so a fifth view cannot be added without touching it. Two lines, no behavior change to the other four views. DECIDE & STATE.
- **D-04**: Extended the write set to `durations_test.go` to generalize `durationsRendererConstant` into `rendererNumericConstant(t, assetPath, name)`. Both views' probes now read the shipped constants through one reader instead of each keeping a copy of the parser. The durations-specific wrapper stays, so its three call sites are untouched. DECIDE & STATE.
- **D-05**: `TIMELINE_MIN_SPAN_MS` is written as the literal `3600000` with a comment rather than `60 * 60 * 1000`. The constant reader that keeps tests honest parses plain literals only, and a constant a test cannot read is a constant a test ends up hand-copying — which is the drift the reader exists to prevent. DECIDE & STATE.
- **D-07**: Every listener this view attaches goes through a teardown registry that runs at the top of each render. Durations attaches its handlers to nodes it rebuilds every render, so they die with the old DOM and need no cleanup; this view binds to the scroll container and to `window`, which both outlive a render — and a filter change re-renders it. Copying the neighbour's habit left five filter changes with six live scroll handlers and every later scroll re-rendering the rows six times. Measured before and after by instrumenting `addEventListener`/`removeEventListener`: 10 live either way with the teardown, 10 → 60 without it. DECIDE & STATE.
- **D-06**: The now-line is one node inside the rows SVG rather than a positioned element over it. The first implementation used an absolutely-positioned span, which was silently erased by the render's own `textContent = ""` and would have needed the container's padding folded into its `left` by hand. Drawing it in the same coordinate space as the bars it crosses removes both problems. DECIDE & STATE.

## Discovered Tasks

- **[low]** The Timeline is the client's first zoom/pan surface, and its interaction affordances are discoverable only from the hint line under the chart ("Hold ⌘ or Ctrl and scroll to zoom the time axis; drag to pan"). Keyboard users get row focus and the detail drawer, but no keyboard path to zoom or pan — the zoom buttons are reachable, panning is not. Worth a follow-up that gives the chart arrow-key panning and `+`/`-` zoom when focused, and states the affordances in the panel's `aria-label` rather than only in adjacent prose. Not a regression: no other view has zoom or pan at all, so this adds an unkeyboarded capability rather than removing a keyboarded one.

## Testing

**Tests run:** `cd skills/do-work-board/tools/queue-kanban && go vet ./... && go test -count=1 ./...` (with the maintainer-strict JavaScript behavior lane under Node), then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing — both exit 0

**Red-green validation:**
- `TestTimelineAggregateProducesTheFourBarShapes` (timeline_test.go): ✗ `undefined: TimelineAggregate / TimelineRow / buildTimelineAggregate` → ✓. Compile failure is the RED the REQ itself names ("there is nothing to ask and the test cannot compile"), because the wait span and the in-flight span existed nowhere in the codebase.
- `TestJavaScriptBehaviorTimelineZoomHoldsTheAnchorInstant` (generate_test.go, Node probe): ✗ `the anchored instant drifted 489796875 ms over three zoom steps` when the transform anchors at the centre instead of the pointer → ✓. Also pins both clamps: zooming out settles exactly at the bound span, zooming in exactly at the one-hour floor.
- `TestJavaScriptBehaviorTimelineVirtualizesRowsAtScale` (generate_test.go, Node probe): ✗ `a 600px viewport rendered 560 of 560 rows; the slice must be viewport-bounded` when the slice is the row count → ✓.
- `TestTimelineRowsAreChronologicalWithAStableTiebreak`, `TestTimelineRangeReachesNowWhileAnyBarIsOpen`, `TestTimelineSkipsTicketsWithoutAParseableCreatedAt` (timeline_test.go): the stated ordering, the fitted range's treatment of open bars, and the drop rule for an unreadable `created_at`.

**New tests added:** four in `timeline_test.go`, two Node behavior probes in `generate_test.go`.

**Existing tests updated (cross-REQ impact):** `generate_test.go`'s embedded-JavaScript inventory and fragment-manifest assertions gained `web/board-timeline.js` — the REQ requires that in the same commit as the file, or the build fails. No assertion was weakened. `durations_test.go`'s constant reader was generalized (D-04) with its three call sites unchanged.

**Rendered acceptance.** Headless Chromium, no page or console errors in any run.
- **This repository's board** (226 rows): renders at fit, 44 row groups for 226 rows.
- **A 560-row synthetic fixture** (545 archived, 6 in flight, 9 never claimed, 10 reversed stamps), which is the scale the REQ names and above anything this repo reaches. **44 row groups and ~220 SVG nodes — identical to the 226-row board**, which is the virtualization claim measured rather than asserted. Zoom via button: 481 ms for two steps. Scroll to row 329: 311 ms.
- **Zoom survives a tab switch (R4):** after zooming in twice, switching to Board and back left the axis at `20 Apr` and the scroll position at row 329.
- **Every bar shape (R1, R2, R3, R7):** an in-flight REQ renders `timeline-segment-work is-open` with `aria-label` "REQ-2547 · Scale fixture 2547 · Route C · claimed · waited 4h 14m · worked 39d 3h so far"; nine never-claimed REQs render an open wait and no work rect; ten reversed stamps render the break marker.
- **Both themes:** light and dark rendered separately; the wait/work pair is legible on each.
- **Listener lifetime:** instrumented `addEventListener`/`removeEventListener` on the scroll host and window across five filter-triggered re-renders — live count stays at 10. Without the teardown it reaches 60. See D-07.

## Lessons Learned

**What worked:** Extracting the two rules that could be silently wrong — the zoom transform and the visible-row slice — as pure functions before writing any DOM code. Both are things no screenshot catches: a zoom that drifts the anchor looks fine in a still, and a slice that is subtly wrong shows blank strips only while scrolling fast. As pure functions they became two Node probes that fail loudly under mutation. Choosing pixel coordinates over a `viewBox` was the other one: the REQ warned that Durations' fixed conversion is exactly what a zoom invalidates, and removing the conversion is a stronger answer than maintaining it.

**What didn't:**
- Assuming `:root` was the light palette. It is the dark default here, with light under `@media (prefers-color-scheme: light)`, so the first render had the wait segment at full navy weight on white — the opposite of the "wait is the quieter hue" intent stated in its own comment.
- An absolutely-positioned now-line. It sat inside the scroll container, so the render's own `textContent = ""` erased it, and it would still have needed the container's padding folded into its `left` by hand. One node inside the rows SVG has neither problem.
- Copying Durations' listener habit. Durations attaches to nodes it rebuilds every render, so its handlers die with the old DOM; this view attaches to the scroll container and to `window`, which both outlive a render — and a filter change re-renders it. Five filter changes left six live scroll handlers, and every later scroll re-rendered the rows six times.

**Worth knowing:** The pattern behind the last two is the same. A view that re-renders into fresh nodes can be careless about cleanup; a view that binds to anything persistent cannot, and "the neighbour does it this way" is not evidence that the neighbour's constraints are yours. The board now has one of each, so the next view added here should decide which it is before copying either.

## Orientation

The board can now show what it never could: how long each REQ sat before anyone claimed it, and what is running right now. Both were invisible — the Durations view measures only claim-to-completion and skips every REQ that is not finished, so the wait was measured nowhere and an in-flight REQ appeared in no panel of any view. **[MAP CHANGED]** — a fifth view, a new `timeline.go` aggregation beside `durations.go`, a new payload section, and the client's first zoom/pan/virtualization surface, which sets the pattern for the next one. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path still resolves and the three-write-surface count is unchanged (this REQ adds none). One sibling document did go stale — `skills/do-work-board/actions/board.md` counts the board's views in prose — and that is REQ-232.

## Review

**Overall: 92%** | 2026-08-18T01:18:27Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 88% |
| Test Adequacy | 92% |
| Scope | 90% |
| Risk | Medium |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- `skills/do-work-board/actions/board.md` counts and part-lists the board's views in three places ("a third view next to Board / Calendar", and two two-item lists), all stale — one of them already wrong before this REQ — gate: user-visible → REQ-232 created `status: pending` as a new sweep (`sweep_key: shipped-prose-hand-counts-board-views`; REQ-227 carries no `review_generated` marker, so the generation-≥2 reroute does not apply).
- The view has no keyboard path to zoom or pan, and its affordances live only in adjacent prose — gate: trivial → recorded in `## Discovered Tasks` as `[low]`, queued by Step 8. Not a regression: no other view has zoom or pan at all.

**Minor findings:** 2 (report only)
- The write set grew twice during implementation (D-03, D-04). Both are real consequences of adding a fifth view to a shell that enumerates its views, and both were recorded and mirrored before the work landed rather than discovered at review — but a plan that had read `board.js` closely would have declared them up front.
- `generate.go` and `generate_test.go` are touched by both this REQ and REQ-226. Flagged by the Step 10 contamination check and cleared: they are the module's shared registries — the fragment manifest and its inventory assertion — which the REQ explicitly requires changing in the same commit as the new fragment.

**Acceptance:** Pass — all nine restated criteria verified: six by mutation-checked assertions or payload inspection, and the scale, zoom-persistence, bar-shape and theme criteria by instrumented headless renders of a 560-row fixture.
**Suggested testing:** 0 items
**Follow-ups created:** REQ-232, plus the Step 8 discovered task; **sweeps appended to:** None

*Reviewed by review-work action*
