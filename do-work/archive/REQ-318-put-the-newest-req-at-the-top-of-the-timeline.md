---
id: REQ-318
title: "Put the newest REQ at the top of the timeline"
status: completed
created_at: 2026-08-22T22:08:34Z
claimed_at: 2026-08-22T22:51:04Z
route: A
completed_at: 2026-08-22T23:19:44Z
commit: 0aa0b03
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-319, REQ-320, REQ-321, REQ-322, REQ-323, REQ-324]
batch: timeline-ux-audit
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-22T22:53:10Z
  basis:
    - trivial short-circuit
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Put the Newest REQ at the Top of the Timeline

## What

Reverse the Timeline view's row order so the most recent REQ is the first row under the
axis and the oldest is last.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reverse the comparator in `buildTimelineAggregate` — one ordering decision
  in one place — rather than reversing client-side; reverse the id tiebreak with it so one
  instant's rows read the same direction as the list around them; rewrite the ordering test
  to pin newest-first and add a four-digit fixture so the tiebreak cannot silently go
  lexical; update the two sentences that state the old order.
- [x] **[APPLY]:** Changes confined to the declared write set, plus two code comments the
  change falsified (`timeline.go`'s `TimelineAggregate` doc and `timelineFirstOpenRowIndex`'s
  "hundreds of rows of finished work" premise).
- [x] **[UNIFY]:** `git diff --stat` reviewed; `go vet ./...` clean; gofmt clean via the
  canonical gate; no debug artifacts in the diff. Files verified — `timeline.go` (comparator
  plus two doc comments), `timeline_test.go` (one test rewritten, one added),
  `web/board-timeline.js` (subhead string, chart aria-label, one stale comment).

## Why

On a 309-REQ queue the current work is 5,700 pixels below the fold. Every visit to the
view opens on REQ-001.

## Context

`buildTimelineAggregate` (`timeline.go`) sorts rows by `created_at` ascending with
`requestIdLess` as the tiebreak, and `board-timeline.js` renders them in payload order.
Reversing in Go keeps one ordering decision in one place and keeps the client a renderer.

Two things read the row order and must still mean what they say afterwards:

- `timelineFirstOpenRowIndex` / `timelineNowJump` scroll to the first still-open row. With
  newest first, that row sits near the top rather than near the bottom; the function still
  works, but confirm the *Now* button lands somewhere sensible rather than assuming it.
- The subhead sentence in `renderTimelineView` literally reads "in capture order, oldest at
  the top". It is false the moment this lands and is this REQ's to fix.

## Detailed Requirements

- Newest `created_at` first, oldest last.
- The tiebreak for two REQs captured in the same instant stays deterministic — reverse it
  with the sort so a build never swaps two rows.
- **Two sentences go false together, not one.** The subhead
  (`web/board-timeline.js:566`, "in capture order, oldest at the top") and the rows SVG
  `aria-label` (`:588`, "One horizontal bar per REQ in capture order") both state the old
  order. The second is the screen-reader description of the whole chart, so leaving it
  behind would tell a non-sighted reader the opposite of what the chart does. Both are this
  REQ's to rewrite.
- `timeline_test.go` currently pins oldest-first
  (`TestTimelineRowsAreChronologicalWithAStableTiebreak`). That assertion is in scope and
  must be rewritten to pin newest-first, not deleted. Name the change in the hand-back: a
  quietly edited test looks identical in a diff.

## Constraints

- Serial with the rest of the `timeline-ux-audit` batch — all of them write
  `web/board-timeline.js`.
- The Calendar view already reads newest-first; matching it is the point, not a coincidence.
- Generate a board and look at it before hand-back. A green suite is not evidence about
  what the first row is (`_dev/primes/prime-kanban-board.md`).

## Builder Guidance

**Certainty: Firm.** The direction is the user's own words — "most recent REQ's should be
on top" — and is not open. Latitude on *where* the reversal lives: reversing the Go sort is
the recommendation because it keeps one ordering decision in one place, not a constraint.
Scope cue: this is a sort direction plus the sentences that describe it. Resist widening it
into the row-label or scroll work that REQ-322 and REQ-319 own.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repo's archive and open the Timeline tab.
The first row under the axis is `REQ-001`; the newest REQ is only reachable by scrolling to
the bottom of the list. In Go, `buildTimelineAggregate` returns `rows[0]` as the earliest
`created_at`, and `timeline_test.go` asserts exactly that.

**Why RED now:** The sort in `buildTimelineAggregate` is ascending by `created_at`.

**GREEN when:** `rows[0]` is the newest `created_at` and the last row is the oldest; two
rows sharing an instant keep a stable, reversed id order across repeated builds; the
generated board opens with recent work on screen; the subhead no longer says "oldest at the
top".

**Validation:** Inferred during capture. The requirement is the user's own words (item 1 of
the request); the RED/GREEN pair above is capture's, and the user has not seen it.

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — rows `REQ-001` through
`REQ-042` filling the visible list in ascending order.

---
*Source: "1.  most recent REQ's should be on top."*

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names its exact write set (`timeline.go`, `timeline_test.go`,
`web/board-timeline.js`), the change is a sort direction plus the two sentences that state
it, and the test that pins the old order is named. Nothing needs discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified — added post-review, D-01)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified — added post-review, D-01)

**What was done:** `buildTimelineAggregate`'s comparator now sorts rows newest `created_at`
first and breaks equal instants by descending numeric id through `requestIdLess`, so the
list reads one direction throughout. The two sentences stating the old order — the subhead
and the rows SVG `aria-label` — now say newest-first, and two code comments the reversal
falsified were corrected in the same change.

## Testing

**Tests run:** `go vet ./...` and `go test ./...` in
`skills/do-work-board/tools/queue-kanban`; `bash _dev/tests/maintainer-verify.sh` from the
repo root (`GOTOOLCHAIN=go1.26.1`, `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`).
**Result:** ✓ All passing — package `ok` in 50.1s; canonical gate exit 0, strict browser
behavior lane included.

**Red-green validation:**
- `TestTimelineRowsAreNewestFirstWithAStableTiebreak`: ✗ before implementation
  (`row order = [REQ-500 REQ-501 REQ-502 REQ-503], want [REQ-503 REQ-502 REQ-501 REQ-500]`)
  → ✓ after. This is the REQ's captured RED/GREEN pair expressed at the aggregate: the
  captured RED named `rows[0]` as the earliest `created_at`, which is what this asserts
  against.
- `TestTimelineNewestFirstTiebreakIsNumeric`: ✗ before implementation
  (`row order = [REQ-999 REQ-1000], want [REQ-1000 REQ-999]`) → ✓ after.

**Render evidence (the half a suite cannot give):** generated a static board from this
repo's own archive and read it in headless Chromium —
`file:///tmp/claude-0/-home-user-skill-do-work/32295d3b-538a-57cc-a4d8-2d453777559b/scratchpad/board318/probe-order.html`,
Chromium 1194 (Playwright build, `/opt/pw-browsers/chromium`). First drawn row labels:
`REQ-324, REQ-323, REQ-322, REQ-321, REQ-320, REQ-319`; last drawn `REQ-300`; 24 label nodes
for 316 rows, so virtualization still bounds the DOM. Table first `REQ-324`, last `REQ-001`,
316 rows. Subhead read back as "316 REQs in capture order, newest at the top." and the chart
`aria-label` as "…in capture order, newest first." The seven REQs of this batch share one
`created_at`, so the top of that list is the descending numeric tiebreak being exercised on
real data.

**Now-button confirmation (the REQ's Context asked for it explicitly):** on the live
316-row board the first still-open row is row 0 under newest-first, so `timelineNowJump`
returns `scrollTop` 0 and the row-list movement is a no-op. That is the correct outcome —
the button's job is to put the reader on the first open REQ, and newest-first already does.
The second movement still earns its place for a reader who has scrolled back through
history, and for a board whose newest open REQ is not row 0.

**New tests added:**
- `TestTimelineNewestFirstTiebreakIsNumeric` — a four-digit fixture (REQ-999 vs REQ-1000),
  the only shape that can catch a descending tiebreak going lexical. Mutation-verified: it
  fails with the tiebreak branch deleted.

**Existing tests updated (cross-REQ impact):**
- `timeline_test.go` — `TestTimelineRowsAreChronologicalWithAStableTiebreak` (from REQ-227,
  which stated the original ordering) renamed to
  `TestTimelineRowsAreNewestFirstWithAStableTiebreak` and its `wantOrder` reversed. **This
  is a deliberate behavior change, not a test bent to fit the code**: the REQ named this
  assertion as in scope, and REQ-237's lesson is that a quietly edited test looks identical
  in a diff. Flagged here so it is not.

## Qualification

Passed — 3 files verified in the diff, 5 requirements traced, P-A-U confirmed.

- Mechanical (`tools/checks/qualify.sh`): OK.
- **Substantive:** the comparator diff is a real inversion (`Before`→`After`, tiebreak
  arguments swapped), not whitespace or import shuffling.
- **Requirements traced:** newest-first → the comparator and the rendered board; stable
  reversed tiebreak → the comparator plus `TestTimelineNewestFirstTiebreakIsNumeric`;
  subhead → `board-timeline.js:566`; chart `aria-label` → `:588`; the named test rewritten
  rather than deleted → `timeline_test.go`.
- **Flowing:** not applicable — no data-fetch path in this change.

## Decisions

- **D-01 — Extend the write set to the two files the reversal falsified outside it.**
  DECIDE & STATE. The review's restatement sweep found stale prose in
  `web/template.html` (the hint's "Scroll to move down the queue", now backwards) and
  `generate_test.go:3462` (a fixture comment calling three-closed-then-open "the board in
  miniature", which newest-first makes untrue). Neither file was in the captured write set.
  UR-065's batch constraint assigns prose a REQ invalidates to that REQ, and both edits are
  one comment each, so extending the set was cheaper and more honest than routing two
  sentences to the prose backlog for a later REQ to find. `write_set` updated in the same
  edit. Reversible: both are comments.

## Review Remediation (pre-archive)

The review returned **Pass / 88% / Approve** with two Important findings, and recommended
fixing both inside this REQ rather than queueing follow-ups, because both live in the
declared write set and total about four lines. Done, with evidence:

- **F1 — the rewritten order test had stopped pinning its own name.** Reversing `wantOrder`
  without touching the fixture made the tied pair's expectation coincide with input order,
  so a stable sort with the tiebreak branch **deleted** satisfied it. Fixed by feeding the
  tied pair in ascending id order (`REQ-500` before `REQ-501`) and saying in a comment why
  the input order is load-bearing. **Mutation-verified**: with the `Equal`/`requestIdLess`
  branch removed via `go test -overlay`, the test now fails —
  `row order = [REQ-503 REQ-502 REQ-500 REQ-501], want [REQ-503 REQ-502 REQ-501 REQ-500]` —
  where before the fix it passed.
- **F2 — two more comments stated the falsified premise.** `board-timeline.js`'s
  `timelineNowJump` header and the `timeline-zoom-now` handler both justified the row-list
  jump by "the oldest archived ones" being on screen. Both rewritten: under newest-first the
  jump is a no-op in the common case and the whole answer for a reader who scrolled back
  through history. The REQ's original claim of "two code comments" undercounted by two;
  the count is now four comments plus two user-visible sentences.
- **Minor — the Now-button confirmation the REQ asked for.** Recorded below in Testing; the
  review confirmed `scrollTop` 0 → 0 on the live 316-row board, which is the correct
  outcome rather than a lost behavior.
- **Nits accepted, not fixed:** `identifierLess`'s no-numeric-suffix branch now leads a tied
  instant instead of trailing it, which needs a malformed id sharing a `created_at` with a
  well-formed one — cosmetic, and the board's ids are `REQ-NNN` by construction.

## Review

**Acceptance: Pass — Overall 88% — Approve.** Independent review agent, orchestrated mode,
serial working-tree diff.

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 94% | 8 of 9 fully delivered; the Now-button confirmation the REQ demanded was not evidenced at review time (recorded since) |
| Code Quality | 82% | Correct, minimal comparator with useful why-comments; two comments in the same file left on the falsified premise |
| Test Adequacy | 78% | RED/GREEN reproduced exactly and the four-digit guard mutation-verified; the rewritten order test had silently lost the tiebreak-deletion guard its name claims |
| Scope | 100% | Exactly the three declared files at review time |
| Risk | Low | Every order consumer identified and checked |
| Acceptance | Pass | `maintainer-verify.sh` exit 0 run by the reviewer, plus an independent 316-row render |

**What the review verified independently rather than accepting:**

- Re-ran the canonical gate itself: exit 0 with the strict JavaScript and strict browser
  lanes both executing.
- Reproduced RED by overlaying the pre-change `timeline.go` from `HEAD` — both new tests
  fail with exactly the messages this REQ records.
- Reproduced the render at 316 REQs in headless Chromium (`HeadlessChrome/141.0.0.0`),
  returning `location.href` in the same evaluation as the measurements. Every figure matched.
- Checked the regression surface rather than assuming it: `timelineRange` (min/max, order
  independent), `timelineVisibleRowRange` (index arithmetic), `timelineFirstOpenRowIndex`,
  `timelineNowJump`, the forecast chain (sorts its own slices, never reads `aggregate.Rows`),
  Calendar and Durations (separate aggregates), and `grep`-confirmed that
  `web/board-timeline.js` is the only consumer of `timeline.rows`.
- Cleared a latent order-dependency at `board-timeline.js:540` as provably unreachable: the
  `rows[0].createdTime` fallback fires only when every row shares one instant with no claim
  or completion, in which case both orders give the same value.

**Findings:** two Important, both fixed in this REQ before archive (see *Review Remediation*
above); three Nits, one fixed (`template.html` hint), one fixed
(`generate_test.go` fixture comment), one accepted (`identifierLess`'s malformed-id branch).

**Suggested additional testing carried forward, not done here:** layout at a narrow width
(this review read text and node counts, not pixels); a filtered Timeline opening on the
newest match; the keyboard focus-restoration path after a window move now that the row at a
given index changed. REQ-319 and REQ-322 both re-enter this panel and are the natural place
for the first two.

*Reviewed by an independent review agent under `actions/review-work.md`, orchestrated mode.*

## Lessons Learned

**What worked:** Reversing the comparator in Go kept one ordering decision in one place and
left the client a pure renderer — the whole client-side change was three sentences of prose.
Generating a real board and reading it in headless Chromium caught nothing wrong but proved
the thing a suite cannot: the top of the list, the virtualized node count, and both rewritten
sentences, all read back from a page whose URL was returned in the same evaluation.

**What didn't:** Reversing a lock-in test's expectation without touching its fixture. The
tied pair's expected order then coincided with input order, so a stable sort with the
tiebreak *deleted* satisfied the assertion — the test kept its name and lost its property.
It passed review's eye and only fell to a deliberate mutation. **A reversed assertion needs
its fixture reversed too, or the sort's own stability answers it.**

**Worth knowing:** Row order is read by more of this module than it looks. Four consumers
had to be cleared (`timelineRange`, `timelineVisibleRowRange`, `timelineFirstOpenRowIndex`,
`timelineNowJump`) and all four were order-independent or self-correcting — but the *prose*
was not: six statements described the old order, and the REQ, written after a careful audit,
named two of them. When a change inverts a stated property, grep the property, not the file.

## Orientation

The board's Timeline view now opens on current work rather than on REQ-001 — one comparator
in `timeline.go`, consumed unchanged by the `web/board-timeline.js` renderer. Leaf change: no
new module, no data-flow change, no renamed concept, so no `[MAP CHANGED]`.
`_dev/primes/prime-kanban-board.md` spot-checked — every path it references still exists, and
nothing this change touched makes it stale.
