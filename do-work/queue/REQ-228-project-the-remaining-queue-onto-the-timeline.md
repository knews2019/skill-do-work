---
id: REQ-228
title: Project the remaining queue onto the Timeline
status: pending
created_at: 2026-08-17T23:51:17Z
user_request: UR-051
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-227]
maintenance: false
related: [REQ-226, REQ-227]
batch: board-timing-views
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# Project the Remaining Queue onto the Timeline

## What

Extend the Timeline forward: give every unstarted REQ a projected bar, chained serially in execution
order, so the view answers "when does the queue empty, and in what order" — with the projection
visually unmistakable as a forecast rather than a measurement.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"...and includes also the timing for the remaining tasks so we know when it ends and in what order."

## Context

REQ-227 builds the Timeline with measured bars for real history and open bars for work in flight.
This REQ adds the forward half.

**The forecast model was chosen by the maintainer at capture**, from three offered options: serial,
median-based. One REQ at a time; no parallelism knob. A lane control for concurrent worktree builders
was offered and explicitly declined — worktree fan-out is real, but the serial model is the honest
floor and the knob is machinery this request does not need. Do not add it back.

**Execution order already has an owner.** `actions/work.md`'s selection scan honours `depends_on`,
not numeric id order, and the board already computes `PendingReady` — "pending with every `depends_on`
target at terminal success, actionable now" (`model.go:267`) — alongside the resolved dependency graph
(`model.go:226-240`: `DependencyEdge`, `DependencyGraph`, forward and reverse adjacency,
`UnmetDependencies` at `model.go:147`). The projection consumes those; it must not invent a second
ordering rule that can disagree with what `do-work run` will actually pick next.

**The duration basis already has an owner too.** `effort_estimate` is a deliberately closed two-value
enum, `trivial` | `normal`, and `model.go:966-971` is explicit about why: *"a triage bit, not an
estimation system. Do not grow it toward t-shirt sizes."* Absent reads as `normal`. That is the split
the median is taken over, and this REQ must not propose a third value or a per-REQ estimate field to
make the forecast prettier.

**And the exclusion rule has one definition.** Paused and reversed spans are already classified by
`dayMedianExclusionReason` (`durations.go:105-115`) against the rule defined once in
`skills/do-work/actions/estimate-reference.md` → Calibration. `durations.go:9-26` states plainly that
it is that rule's *second reader, not a second definition* — and REQ-219's recorded lesson is to ship
a rule's verdict in the payload so a second reader cannot become a second definition. A forecast that
re-derives "over four hours means paused" is exactly the failure both of those warn about. Consume
the verdict.

**There is also an `estimate:` block** written by the work action's ensure-estimate step
(`actions/work-reference.md:139-156`, visible on archived REQs as `p50_active_minutes`, `confidence`,
`basis`). It exists only on REQs that have already been claimed, so it cannot forecast an unstarted
REQ — but it is a legitimate calibration reference and worth a look before settling the median
window.

## Detailed Requirements

1. **Serial chain from now.** The first projected REQ starts when the currently-running work finishes
   (or now, if nothing is running); each subsequent one starts when its predecessor ends. One at a
   time, no overlap.
2. **Order is dependency-then-queue.** Consume the existing dependency graph and ready-set rather
   than re-deriving. A REQ whose prerequisites are unmet is scheduled after them, never before.
3. **Projected span = rolling median of recent completed work spans, split by `effort_estimate`.**
   Take the median over `trivial` and `normal` separately; a REQ with no `effort_estimate` is
   forecast as `normal`. Excluded spans (paused, reversed) do not feed the median — consume the
   existing verdict for that.
4. **State the window and the sample size on the view.** "Rolling" needs a definition — last N
   completions or last N active days, builder's choice — and the reader must be able to see which one
   and how many samples backed it, because a median over four samples deserves less trust than one
   over eighty.
5. **Degrade honestly on thin history.** Below a stated minimum sample count, say so and decline to
   draw a confident end date rather than projecting off two data points.
6. **A forecast must never read as a fact.** Projected bars are visually distinct from measured bars
   at a glance — hatching, reduced opacity, outline-only, whatever reads clearly in both themes — and
   the distinction is in the legend.
7. **Unschedulable REQs are shown as unschedulable.** `blocked` (waiting on an external condition),
   `pending-answers` (waiting on the user), and `blocked-dependency-cycle` REQs cannot be given a
   start time honestly. List them, exclude them from the chain, and state that the end estimate
   excludes them. Silently folding them in would be the single easiest way to make this view lie.
8. **A queue-end readout.** One plain-language line — "queue empties ~<when>, N REQs" — with the
   assumptions it rests on stated next to it, not buried.
9. **Nothing is written.** The projection is derived at render time. No REQ file gains a projected
   date, no frontmatter field is added, no pipeline state is touched.

## Constraints

- Read-only. `CLAUDE.md` § *Kanban Board Write Surfaces* must still read "exactly three" and go
  unamended. A forecast is a derived display; the moment it persists anywhere it becomes a fourth
  write surface and this constraint is violated.
- No new frontmatter field. Not on REQs, not on URs. `effort_estimate` stays a two-value triage bit.
- Do not re-derive the read-time rule or the dependency rules; consume the existing verdicts.
- `durations_test.go:202-225` pins live-archive figures — unchanged by this REQ.
- The projection must not contradict what `do-work run` would actually claim next. If the two can
  disagree, the ordering source is wrong.
- Framework-free plain JS, both themes, both the live server and the static snapshot.

## Dependencies

`depends_on: [REQ-227]` — this REQ draws into the view REQ-227 creates and extends the aggregation
REQ-227 introduces. It cannot land first.

**Declared write-set overlap with REQ-227** on `timeline.go`, `timeline_test.go`,
`web/board-timeline.js`, and `web/board.css`. That overlap is unavoidable given the slice boundary
(measured half / projected half) and is declared rather than designed around, per the slicing
convention in `actions/capture-reference.md`. The dependency edge makes it serial in practice, so the
board's overlaps badge is informational here, not a collision warning.

## Builder Guidance

**Certainty: Firm on the model, Mixed on presentation.** Serial-and-median-based was chosen
deliberately over a parallelism knob and over showing no forecast at all. Do not re-open it. The
rolling window length, the minimum sample count, the hatching treatment, and the exact wording of the
queue-end readout are yours.

The honesty requirements above are the load-bearing part of this REQ, not decoration. A Gantt that
projects a confident end date is exactly the kind of artifact people screenshot and quote, and the
serial single-median model is a deliberately crude estimator — it ignores parallel builders, review
loops, blocked-then-unblocked churn, and the fact that the queue grows while it drains. The
mitigation is that the view says what it assumed, not that the estimate is made cleverer. If you find
yourself adding sophistication to make the number more defensible, prefer adding a caveat instead.

Worth checking early: whether the currently-running REQ's own `estimate.p50_active_minutes` is a
better start-offset for the chain's head than the global median. It is available only for claimed
REQs, which is exactly the one bar that needs it.

## Red-Green Proof

**RED prompt/case:** Build a fixture with a known completion history (a handful of `trivial` spans
and a handful of `normal` spans, plus one paused and one reversed span that must be excluded) and
four pending REQs — two `normal`, one `trivial`, and one `blocked` — where one of the `normal` ones
declares `depends_on` on the other. Ask the board for the projection. Assert the three schedulable
REQs are chained serially in dependency-respecting order, that each projected span equals the median
for its own `effort_estimate` bucket with the excluded spans absent from both medians, that the
blocked REQ carries no start time and is reported as excluded, and that the queue-end instant equals
the last projected finish. Today there is no projection of any kind, so the test cannot compile.

**Why RED now:** Nothing in the codebase forecasts anything. `buildDurationAggregate`
(`durations.go:75-103`) skips every non-terminal ticket outright, so pending REQs contribute to no
computation anywhere, and no code path produces a future instant.

**GREEN when:** That fixture produces exactly the expected chain, medians, exclusions, and end
instant; and the Timeline view shows projected bars visibly distinct from measured ones, a stated
window and sample size, an excluded-REQ list, and a queue-end readout with its assumptions beside it.

**Validation:** User confirmed — the serial median-based model was chosen by the maintainer from
three offered options, with the parallelism knob explicitly declined.

## Full Context

See `do-work/user-requests/UR-051/input.md` for complete verbatim input.

---
*Source: "...and includes also the timing for the remaining tasks so we know when it ends and in what order."*
