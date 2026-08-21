---
id: REQ-313
title: "Count the breaks the timeline actually draws"
status: completed
created_at: 2026-08-21T08:25:57Z
status_changed_at: 2026-08-21T15:38:59Z
claimed_at: 2026-08-21T17:58:54Z
completed_at: 2026-08-21T18:16:11Z
commit: 0761a10
route: B
user_request: UR-062
addendum_to: REQ-304
domain: frontend
review_generated: true
impact: impact-user-visible
kb_status: pending
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
depends_on: []
maintenance: false
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-21T17:59:34Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
    - full-suite verification
---

# Count the Breaks the Timeline Actually Draws

## What

The Timeline view's summary line renders "N with broken stamps, drawn as breaks"
(`web/board-timeline.js:546-557`), and N counts rows where `row.anomaly` is true. That flag
carries `detectCompletionAnomaly`'s verdict, which `model.go:246-252` scopes to **completion
bookkeeping** — no completion instant at all on a terminal REQ. It is false for both spans the view
draws as breaks:

- a reversed **work** span, whose `completed_at` parses fine but precedes `claimed_at`
- a reversed **wait** span, which REQ-304 just started drawing as a break

So the view draws break markers and the sentence beneath them says zero.

## Why

The count is the only place a reader learns how many rows to distrust without scrolling the whole
Gantt. A count that says zero while breaks are on screen teaches the reader to ignore the sentence.

## Context

Discovered during REQ-304, which fixed the reversed-wait rendering and was explicitly forbidden by
its own Context from broadening the anomaly flag to fix the count. That constraint stands and is the
reason this is a separate REQ rather than one more line in that one:

- `timeline.go:22-24` states that `detectCompletionAnomaly` decides what is broken and this file
  consumes that verdict rather than restating it. Recomputing an anomaly inside the renderer would
  create the second definition that comment exists to prevent.
- REQ-280 closed the **detection** side — the `created_at <= claimed_at <= completed_at` ordering
  probe in `queue-kanban verify` and `forensics.md` Check 12. It did not touch this count, and
  the count is not detection.

The shape that satisfies both: count what the renderer **drew**, not what a flag says. The draw pass
already decides, per row, whether each segment became a break; the summary can total those decisions
instead of asking a flag that answers a different question.

## Detailed Requirements

- The summary line's count equals the number of rows the current render drew at least one break
  marker for — reversed wait, reversed work, or the existing completion anomaly.
- No new anomaly verdict, flag, or reason is added to the timeline row model, and
  `detectCompletionAnomaly` is not broadened.
- The count follows the filters, like every other number on that line.
- The sentence still reads correctly when the count is zero.

## Red-Green Proof

**RED prompt/case:** Render a timeline whose rows include one reversed wait and one reversed work
span, neither carrying `anomaly: true`, and read the summary text. Today it omits the break clause
entirely while two break markers are drawn.
**GREEN when:** The same render reports 2, and a render with no reversed span and no completion
anomaly still omits the clause.
**Validation:** Discovered task from REQ-304; apply `actions/work-reference.md` →
**Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-304: the Timeline view prints a line
  saying how many REQs have broken timestamps and are drawn as break markers, but it counts a flag
  that only covers one kind of breakage. Two other kinds now draw break markers and are not counted,
  so the board can show you breaks while the sentence beneath says none. Fixing it means counting
  what the view drew rather than asking the flag — the flag itself must not grow, because a second
  definition of "broken" inside the renderer is exactly what the file's own comment forbids. Should
  I process this as a new task? → Confirmed: Yes, add to queue

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  **Answered 2026-08-21** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the builder's recommendation via
  `do-work clarify`: make the summary count the rows for which the current view actually draws a
  break, while preserving the existing anomaly-model boundary. No additional scope was requested.

  Why this is yours rather than mine: the fix is small, but it decides that the renderer starts
  keeping a tally of its own drawing decisions, which is a new kind of state in a file that has so
  far only consumed verdicts. That is a design call about where "what is broken" is allowed to be
  answered, and REQ-304's constraint was written to stop me making it unilaterally.

---

## Triage

**Route: B** — Medium

**Reasoning:** The correction is localized to the Timeline renderer and its existing generated-client
DOM probe seam, with no model or schema change. It still needs RED/GREEN TDD, mutation evidence, and
rendered-text acceptance because the visible summary must stay aligned with three independent drawing
branches and active filters.

**Planning:** Not required for Route B; the REQ already fixes the architectural boundary and the
exploration phase will identify the exact state/tally seam.

**Estimate:** 30 active minutes (P50, medium confidence).

## Exploration

The filtered `rows` array is already the population every summary number describes. The summary is
computed before virtualized SVG rows are materialized, while the two marker branches live later in
`renderVisibleRows`: negative waits and negative work spans. Counting only materialized DOM nodes
would make the number change with scrolling, so the narrow seam is a row-deduplicated union over the
filtered data: existing `row.anomaly`, reversed wait, or reversed work.

The existing `row.anomaly` cause remains in that union for backwards compatibility even though some
non-reversed completion anomalies render as an open work bar rather than a `.timeline-segment-broken`
rect. This REQ explicitly says to preserve that existing cause; repairing the older wording/render
mismatch would change anomaly semantics outside scope.

The existing generated-client DOM harness can test the caller seam directly. A dedicated fixture
needs reversed wait, reversed work, both on one row (dedupe), a healthy row, and an anomaly-only row;
filter variants prove the count follows the filtered population and the healthy zero case omits the
clause.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — replace the anomaly-only
  summary count with the unique filtered-row union of the renderer's three existing break causes
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — add caller-seam RED/GREEN and
  mutation-resistant summary/marker/filter coverage using the existing Timeline DOM harness

**Files I will NOT touch:** `timeline.go`, `model.go`, CSS, templates, or payload schema. The REQ
forbids broadening `detectCompletionAnomaly` or adding a second model verdict; no geometry changes.

**Acceptance criteria (restated from REQ):**
- [x] The summary counts unique filtered rows having an existing completion anomaly, reversed wait,
  or reversed work, even when one row has more than one cause.
- [x] No anomaly verdict, payload field, reason, or model behavior changes.
- [x] Active filters change the count, and a healthy-only/zero-break view keeps the current sentence.
- [x] The captured reversed-wait plus reversed-work case fails before the fix and reports 2 after it.
- [x] The direct canonical repository gate exits 0.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/`; the clarified REQ-314 queue edit is preserved and excluded.
**Tests baseline:** ✓ `cd skills/do-work-board/tools/queue-kanban && go test -count=1 ./...`
**Dependencies:** ✓ Go 1.26.1, Node, Chromium headless shell, and ShellCheck are available.

*Checked by work action*

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the board prime and relevant lessons; use the existing real-renderer DOM
  harness to pin a row-deduplicated union over filtered rows, leaving model anomaly semantics intact.
- [x] **[APPLY]:** Added a full-renderer DOM summary probe RED-first, then replaced the anomaly-only
  tally with a filtered, row-deduplicated union of the three existing draw causes.
- [x] **[UNIFY]:** Re-read both write-set files; gofmt, Node syntax, focused renderer probes, the
  strict JavaScript lane, and the full board Go package pass; no scratch files remain.

## Decisions

- **D-01 — Count one filtered row when any existing break cause applies.** DECIDE & STATE. A single
  boolean predicate is the smallest representation of the renderer's row-level contract and makes
  a row with reversed wait and work contribute one, not two. Reversible: one local predicate.
- **D-02 — Gate reversed work on `row.hasWork`, matching the actual drawing branch.** DECIDE &
  STATE. The summary must count what this render can draw; the work marker is unreachable when the
  row has no work, even if malformed input carries a negative numeric field. This does not add or
  broaden any model verdict.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** The Timeline summary now counts unique rows in the already-filtered population
when any of the renderer's existing break causes applies: completion anomaly, reversed wait, or
drawable reversed work. A real-renderer DOM test covers every cause, multi-cause deduplication,
filtered counts, and the healthy zero-clause behavior without changing the row model or SVG geometry.

## Qualification

Passed — 2 implementation files verified, 5 acceptance criteria traced, P-A-U confirmed. Mechanical
qualification and Scope/Implementation Summary parity pass; the renderer change is substantive,
the test reaches the real caller, and no model/schema or unrelated board file changed.

## Testing

**Tests run:** `node --check skills/do-work-board/tools/queue-kanban/web/board-timeline.js`;
focused `go test` for the new summary and REQ-304 reversed-wait probes; strict JavaScript behavior
lane; `go test -count=1 ./...`; direct unpiped `bash _dev/tests/maintainer-verify.sh` with the
declared Chromium headless shell.
**Result:** ✓ All passing; canonical repository gate exited 0.

**Red-green validation:**
- `TestJavaScriptBehaviorTimelineSummaryCountsRowsDrawnAsBreaks`: ✗ before implementation — the
  reversed-wait plus reversed-work fixture expected `2 with broken stamps` but the anomaly-only
  renderer omitted the clause → ✓ after implementation.
- Seven restored mutations each made the focused test fail: anomaly-only revert; omission of wait,
  work, or existing-anomaly causes; counting causes instead of rows; using global instead of
  filtered population; and always printing the zero-count clause.

**New tests added:**
- `TestJavaScriptBehaviorTimelineSummaryCountsRowsDrawnAsBreaks` in `generate_test.go`, driving the
  real renderer across all three causes, multi-cause deduplication, filtering, and healthy zero state.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-21T18:15:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None.

**Minor findings:** 1 (the summary test pins every predicate and deduplication but does not directly
compare the count with emitted marker row IDs; current code's `row.hasWork` guard is correct and the
existing adjacent marker test keeps the risk low)
**Acceptance:** Pass — the count follows unique filtered rows across all three preserved causes.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Counting with one boolean predicate over the already-filtered rows makes deduplication automatic
  and makes the summary inherit the same population as every neighboring number.
- A real-renderer caller-seam test caught anomaly-only, global-population, and cause-count mutations
  without introducing a second implementation of timeline filtering.

**What didn't:**
- The test fixture proves summary semantics and the adjacent REQ-304 test proves marker emission, but
  it does not join those two observations by row ID. That leaves the `row.hasWork` guard defended by
  code review and production payload shape rather than a single end-to-end assertion.

**Worth knowing:** Do not tally only materialized SVG rows. Timeline rows are virtualized, so a DOM
count would change with scrolling; summary counts must use the filtered data population. The legacy
`row.anomaly` cause stays in the union even where an anomaly is not a reversed-span marker.

## Orientation

**[MAP CHANGED]** The queue-kanban Timeline summary now defines “drawn as breaks” as a unique-row
union across its preserved completion anomaly, reversed-wait, and reversed-work causes. The renderer
owns that display tally over filtered rows; the Go model remains the sole owner of anomaly verdicts.
This keeps the sentence below the Gantt aligned with what the reader sees without making scrolling
or filters change the semantic boundary.
