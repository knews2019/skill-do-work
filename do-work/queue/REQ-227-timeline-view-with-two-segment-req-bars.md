---
id: REQ-227
title: Add the Timeline view with two-segment REQ bars
status: pending
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
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# Add the Timeline View with Two-Segment REQ Bars

## What

Add a fifth board view, **Timeline** — a zoomable, scrollable Gantt with one horizontal bar per REQ.
Each bar carries two coloured segments: `created_at`→`claimed_at` (the wait) and
`claimed_at`→`completed_at` (the work). REQs that are claimed but not finished draw as open bars
running to now.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
