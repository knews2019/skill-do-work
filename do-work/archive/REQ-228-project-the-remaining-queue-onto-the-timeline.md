---
id: REQ-228
title: Project the remaining queue onto the Timeline
status: completed
completed_at: 2026-08-18T01:35:36Z
claimed_at: 2026-08-18T01:22:09Z
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
route: C
estimate:
  p50_active_minutes: 55
  confidence: medium
  calculated_at: 2026-08-18T01:22:00Z
  basis:
    - Route C
    - 4-file write set
    - 2 subsystems involved
    - 9 acceptance criteria
    - dependency depth 1
    - browser evidence
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/web/template.html
---

# Project the Remaining Queue onto the Timeline

## What

Extend the Timeline forward: give every unstarted REQ a projected bar, chained serially in execution
order, so the view answers "when does the queue empty, and in what order" — with the projection
visually unmistakable as a forecast rather than a measurement.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`, the dependency graph and ready-set in `model.go`, the read-time rule in `durations.go`, and the view REQ-227 just built. Approach in `## Plan` below.
- [x] **[APPLY]:** Eight files, all declared. One addition beyond the plan — the `Now` control (D-03).
- [x] **[UNIFY]:** `git diff --stat` = 8 files, +936/−4. `gofmt -l .` clean, `go vet ./...` clean, `node --check` on the assembled client clean, zero page errors in every headless run. Grep for `console.log|fmt.Print|debugger|TODO|FIXME|XXX|alert(` in added lines: none. Files verified:
  - `timeline.go` — the projection is a pure computation; grepped for `os.Create|os.Write|os.Mkdir|os.Remove|os.Rename|.Write(`: no file operations, which is Requirement 9 checked rather than asserted.
  - `durations.go` — one field added and one line copying it; the read-time rule and `medianOf` untouched.
  - `timeline_test.go` / `generate_test.go` — four Go tests and one Node probe; every honesty guarantee mutation-checked.
  - `generate.go` — additive wire types and projection only.
  - `web/board-timeline.js`, `web/template.html`, `web/board.css` — projected drawing, forecast, excluded list, legend; both theme blocks carry the new token.
  - `CLAUDE.md` — untouched, and § Kanban Board Write Surfaces still reads "exactly three".

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

---

## Triage

**Route: C** - Complex

**Reasoning:** A forecast layer spanning a new Go computation, its payload, and new drawing in the view REQ-227 created, with nine numbered requirements of which four are honesty constraints rather than features, and three separate existing rules it must consume rather than restate.

## Plan

### What must be consumed rather than re-derived

Three rules already have single owners, and the REQ names all three as the failure modes to avoid.

- **Which spans count.** `dayMedianExclusionReason` (`durations.go:105-115`) is the read-time rule's one reader. The medians here are taken over `DurationAggregate.Samples` that the rule already kept, so no code in this REQ decides what a paused or reversed span is.
- **Which REQs are ready.** `UnmetDependencies` (`model.go:147`) and the `PendingReady` partition (`model.go:1408-1413`) are already computed against `depends_on`. The chain consumes them.
- **What order they run in.** `actions/work.md` Step 1 processes dependency-ready REQs in numeric id order. The chain uses exactly that rule, so the projection cannot contradict what `do-work run` would claim next.

### The computation, in `timeline.go`

`buildTimelineProjection(tickets, durationAggregate, now)`:

1. **Medians.** Walk the rule-included samples newest-first, take the most recent `timelineProjectionWindowSize` of them, and take a median per `effort_estimate` bucket. That needs each sample's bucket, so `DurationSample` gains an `EffortEstimate` field copied at aggregation (D-01) — reading the bucket off the sample the rule already classified is what keeps this a consumer.
2. **Confidence.** A bucket with fewer than `timelineProjectionMinimumSamples` of its own falls back to the window's overall median and says so. A window with fewer than that overall produces **no projection at all** — no bars, no end date, and a stated reason. Declining is the honest floor the REQ asks for.
3. **The chain head.** In-flight REQs finish at `claimed_at + their bucket's median`, floored at now; the chain starts at the latest of those. The `estimate:` block would be a better offset for exactly this bar, but the board does not parse it and adding a nested-block parse surface for one bar is the sophistication the REQ tells the builder to trade for a caveat (D-02).
4. **The chain.** Repeatedly take the lowest-id pending REQ whose dependencies are all either terminally resolved or already placed, and give it `[cursor, cursor + bucket median]`. That is work.md's rule, applied.
5. **Exclusions.** `blocked`, `pending-answers`, and every `blocked-*` status cannot be given an honest start time. Anything still unplaced when the loop stalls — a dependency cycle, a dangling dependency, or a dependency on an excluded REQ — is excluded too, with the reason naming which. Every exclusion is listed and none silently joins the chain.
6. **Queue end** is the last projected finish, and it is absent when nothing was placed.

### The view

Projected work attaches to the row the REQ already has: REQ-227 draws a pending REQ's open wait and deliberately no work segment, and this adds that segment at the REQ's scheduled slot. Projected bars use an SVG hatch pattern — distinct from both the solid measured bars and the dashed open ones, and themed through CSS rather than a hard-coded stroke. The panel gains an excluded-REQ list and a queue-end readout carrying its own assumptions.

### Verification steps

1. Go RED in `timeline_test.go` → verify: the REQ's own fixture — a known history with one paused and one reversed span, plus four pending REQs of which one depends on another and one is blocked — produces the stated chain, medians, exclusions and end instant.
2. A thin-history test → verify: below the minimum the projection declines rather than forecasting.
3. `go test -count=1 ./...` → verify: `durations_test.go:202-225`'s pinned figures are unchanged by the new `DurationSample` field.
4. Headless render → verify: projected bars are distinguishable from measured and open ones in both themes, the excluded list and readout render, and `CLAUDE.md`'s write-surface count is untouched.
5. `bash _dev/tests/maintainer-verify.sh` → verify: exit 0.

*Generated by Plan agent (inline, serial mode)*

## Exploration

**Already computed, ready to consume.** `RequestTicket.UnmetDependencies` (`model.go:147`), `BoardColumns.PendingReady`/`PendingWaiting` (`model.go:266-268`, populated at `model.go:1408-1413` purely from `len(UnmetDependencies)`), `DependencyGraph` with forward and reverse adjacency (`model.go:235-241`), and `EffortEstimate` normalized against the closed two-value enum whose comment forbids growing it (`model.go:966-971`, default `normal`).

**The read-time rule.** `dayMedianExclusionReason` (`durations.go:105-115`) returns `"reversed"`, `"paused"`, or `""`; `DurationSample.DayMedianExclusion` carries the verdict and `ExcludedFromDayMedian()` reads it. `medianOf` (`durations.go:174`) already exists and returns a false second result for "no samples", which is the exact distinction the confidence gate needs.

**Not available.** The work action's `estimate:` block is a nested YAML mapping and the board parses no nested blocks — `grep p50_active_minutes model.go` is empty. It is therefore not a usable start-offset without new parse machinery.

**The view to extend.** REQ-227's `web/board-timeline.js` draws each row's segments inside `renderVisibleRows`, with `drawSegment` clamping to the visible window. A pending REQ's row already exists and already carries its open wait; `hasWork: false` is exactly the gap this REQ fills.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modify) — the projection, its medians, chain, exclusions and confidence gate.
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modify) — the REQ's fixture, plus thin-history and exclusion tests.
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — carry each sample's `effort_estimate` bucket (D-01).
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — projection wire types and projection.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — a behavior probe for the projected-segment drawing.
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — projected segments, excluded list, queue-end readout.
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — the excluded-list and readout containers, and the legend entry.
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — the hatch pattern's styling.

**Files I will NOT touch:** `model.go` (the dependency and effort rules are consumed, not changed), `skills/do-work/actions/estimate-reference.md` (the read-time rule stays where it is), `CLAUDE.md` (no fourth write surface), any REQ file (nothing is persisted).

**Acceptance criteria (restated from REQ):**
- [ ] Projected REQs form a serial chain from now, one at a time, no overlap.
- [ ] Order is dependency-then-queue, consuming the existing graph and ready set.
- [ ] Each projected span is the rolling median for its own `effort_estimate` bucket, with excluded spans absent from both medians.
- [ ] The window and its sample size are stated on the view.
- [ ] Thin history is declined rather than forecast.
- [ ] Projected bars never read as measured ones, and the distinction is in the legend.
- [ ] `blocked`, `pending-answers` and `blocked-*` REQs are listed as unschedulable and excluded from the chain and the end estimate.
- [ ] A plain-language queue-end readout states its assumptions beside it.
- [ ] Nothing is written — no REQ file, no frontmatter, no pipeline state.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)

**What was done:** `buildTimelineProjection` forecasts the unstarted queue as a serial chain. Spans come from the median of the last 60 in-rule completions, split by `effort_estimate` bucket, with a bucket too thin to speak for itself borrowing the overall median and a window too thin overall producing no forecast at all. Order is work.md's own rule — among REQs whose dependencies have resolved or been placed, take the lowest id — so the chain cannot contradict what `do-work run` would claim next. The chain starts after work already in flight. Every REQ that cannot be given an honest start time is listed with a plain-words reason: the held statuses directly, and anything left unplaceable behind them. In the view, a projected span attaches to the row its REQ already has, drawn in a hatch pattern and its own hue, and the panel gained a queue-end sentence carrying its assumptions, an excluded list whose ids open the detail drawer, and a `Now` control. Nothing is written anywhere.

## Decisions

- **D-01**: `DurationSample` gained an `EffortEstimate` field. The medians must be split by bucket, and the samples are where the read-time rule has already decided what counts. Copying the bucket onto the sample lets the projection read both off one record; the alternative was re-walking the tickets and calling the classifier a second time, which is how a consumer turns into a second definition. The field is inert for Panel A and the day medians. DECIDE & STATE.
- **D-02**: The chain's head uses the in-flight REQ's bucket median, not its own `estimate.p50_active_minutes`. The REQ's guidance asks whether that block is the better offset for exactly this bar, and it would be — but the board parses no nested frontmatter blocks (`grep p50_active_minutes model.go` is empty), so using it means a new parse surface for one bar. The REQ's own instruction covers this case: "If you find yourself adding sophistication to make the number more defensible, prefer adding a caveat instead." The assumption is stated in the forecast sentence. DECIDE & STATE.
- **D-03**: Added a `Now` control beside zoom in/out/fit. Not in the REQ. Zoom anchors at the centre of the window, so on this repo's three-month axis the forecast at the far right needs a long drag to reach once you have zoomed in enough to read it — which was how the first attempt to screenshot the projected bars landed on 8 July. A forecast that has to be hunted for does not answer "when does the queue empty". Ten lines, reusing the existing clamp. ESCALATE. **Value:** the view's headline answer is one click away at any zoom. **Risk:** a fifth control in a toolbar that had four; reversible by deleting one button and one handler.
- **D-04**: Projected bars are hatched in their own hue rather than tinted. Reduced opacity alone is what a disabled control looks like, and the dashed outline was already taken by an *open measured* bar — a reader would have had to judge saturation to tell a measurement from a guess. The hatch lines carry a CSS class rather than a literal stroke, so both themes style them. DECIDE & STATE.
- **D-05**: A pending REQ's projected work attaches to the row it already has, rather than becoming a second row. REQ-227 draws that row's open wait and deliberately no work segment; this fills exactly that gap, so one REQ still means one row.

## Testing

**Tests run:** `cd skills/do-work-board/tools/queue-kanban && go vet ./... && go test -count=1 ./...` (with the maintainer-strict JavaScript behavior lane under Node), then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing — both exit 0

**Red-green validation:** every honesty guarantee was mutation-checked, because each is a way this view could quietly lie rather than visibly break.

- `TestTimelineProjectionChainsTheQueueSerially` (timeline_test.go): the REQ's own Red-Green Proof fixture — six normal and six trivial spans, plus a 900-minute paused span and a −120-minute reversed one, both normal-bucket and both enormous so neither could reach a median of 40 unnoticed. RED was `undefined: buildTimelineProjection`, the compile failure the REQ names. Mutation: dropping the dependency guard from the placement loop gives `chain[1] = REQ-402, want REQ-403`.
- `TestTimelineProjectionDeclinesOnThinHistory`: ✗ `two samples must not produce a confident forecast` when the minimum-sample gate is removed → ✓.
- `TestTimelineProjectionExcludesREQsBlockedBehindAnUnschedulableOne`: the subtle case — a pending REQ whose prerequisite is blocked, one with a dangling dependency, and one `pending-answers`. All four are named with reasons and none joins the chain, so the queue-end figure cannot describe a subset while looking complete.
- `TestTimelineProjectionStartsAfterWorkAlreadyInFlight`: the chain starts at the in-flight REQ's remaining median, not at now.
- `TestJavaScriptBehaviorTimelineForecastStatesItsAssumptions` (Node probe): asserts each stated assumption separately rather than as one blob, so any single regression names itself. Mutations: dropping the serial/no-parallelism/static-queue clause gives `does not state the serial assumption`; dropping the per-bucket sample count gives `does not state each bucket's sample count and median`. The declined path is asserted in the same probe — it must not say "Queue empties" and must carry its reason.

**Regression evidence (REQ constraint):** `TestLiveArchiveDurationsMatchTheCalibratedFigures` (durations_test.go:202-225) passes unchanged despite the new `DurationSample` field.

**Requirement 9 checked, not asserted:** grepped `timeline.go` for `os.Create|os.Write|os.Mkdir|os.Remove|os.Rename|.Write(` — no file operations. `CLAUDE.md` is untouched and § Kanban Board Write Surfaces still reads "exactly three".

**New tests added:** four in `timeline_test.go`, one Node behavior probe in `generate_test.go`.

**Existing tests updated (cross-REQ impact):** none.

**Rendered acceptance.** Headless Chromium, no page or console errors in any run.
- **This repository's board:** "Queue empties around 2026-08-18 02:07 UTC — 2 REQs run one at a time from 2026-08-18 01:37 UTC. Assumes the median of the last 60 completed REQs (55 normal at 15 min, 5 trivial at 9 min), one REQ at a time, no parallel builders, and a queue that stops growing. Paused and reversed spans are excluded from both medians." Below it, three `pending-answers` REQs listed as "waiting on an answer from you", each id opening the detail drawer.
- **Projected bars, zoomed:** REQ-229 and REQ-232 draw hatched purple spans to the right of the now-line, chained one after the other, on the rows that already carry their measured open waits. The three excluded REQs draw their waits and no projected work. Verified in both themes.
- **The declined path, end to end:** a two-completion fixture renders "No end estimate: only 2 completed REQs inside the read-time rule; 5 are needed before a median means anything." with zero projected bars.

## Lessons Learned

**What worked:** Writing the fixture straight from the REQ's Red-Green Proof rather than paraphrasing it. Its two poisoned spans — a 900-minute paused one and a −120-minute reversed one, both in the same bucket as the figure under test — mean the median cannot come out at 40 if either leaks in. That is a better guard than asserting an exclusion count, because it fails on the consequence rather than on the bookkeeping. Asserting each stated assumption as its own check was the other one: a single "contains everything" assertion would have reported one failure for any of seven different regressions.

**What didn't:** Trying to screenshot the projected bars by zooming in nine times. Zoom anchors at the centre of the window, so it landed on 8 July while the forecast was at 18 August, and the capture came back with zero projected rects. The reflex was to fix the screenshot script; the actual finding was that a reader would have had exactly the same problem, which is what the `Now` control is (D-03).

**Worth knowing:** The estimator is deliberately crude and the honesty requirements are what make it publishable — but the crudeness is also load-bearing in a way worth stating: because the chain uses work.md's own ordering rule rather than a scheduler of its own, the forecast can be wrong about *timing* without ever being wrong about *order*. That is the property to preserve if anyone later makes the duration model cleverer. The moment the projection sorts by anything other than dependency-then-id, it starts disagreeing with what `do-work run` will actually claim next, and a forecast that predicts the wrong next REQ is worse than one that predicts the wrong hour.

## Orientation

The board's Timeline now runs forward as well as back: every unstarted REQ gets a projected bar in the order `do-work run` would actually claim it, and the view states when the queue empties along with everything that estimate assumes. Anything that cannot be given an honest start time — waiting on you, waiting on an external condition, or stuck behind something that is — is listed rather than folded in, so the end figure never quietly describes a subset. **[MAP CHANGED]** — `timeline.go` now holds a forward projection beside the measured aggregation, and it is the first thing in the board that produces a future instant at all. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path resolves, the three-write-surface count is unchanged and was re-verified rather than assumed, and REQ-219's payload-verdict lesson is the one this REQ applied for the third time.

## Review

**Overall: 94%** | 2026-08-18T01:35:28Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 95% |
| Scope | 95% |
| Risk | Medium |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None. The restatement sweep found no live restatement of the projection's rules: the read-time rule, the dependency rules and work.md's ordering rule are each consumed from their one owner, and no shipped prose describes the forecast. `skills/do-work-board/actions/board.md` is already queued as REQ-232 from REQ-227's sweep and needs no second entry.

**Minor findings:** 2 (report only)
- The forecast's headline is a single instant. It is hedged by the assumptions printed beside it and by the decline-on-thin-history gate, but a reader who quotes only the first clause quotes a point estimate. A range would be more honest and is exactly the sophistication the REQ tells the builder to trade for a caveat, so it is recorded rather than built.
- `durations.go` gained a field that only the Timeline reads. It is on `DurationSample` because that is the record the read-time rule has already classified, which is the right home, but it does mean a durations type now carries a timeline concern.

**Acceptance:** Pass — all nine restated criteria verified: five by mutation-checked Go assertions, one by a mutation-checked Node probe, one by grep (nothing is written), and the projected-bar distinction and declined path by headless renders in both themes.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*
