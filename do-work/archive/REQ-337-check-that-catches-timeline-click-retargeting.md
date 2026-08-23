---
id: REQ-337
title: "A check that can catch Timeline click retargeting"
status: completed
claimed_at: 2026-08-23T20:16:43Z
completed_at: 2026-08-23T20:26:47Z
commit:
kb_status: pending
created_at: 2026-08-23T18:30:26Z
user_request: UR-067
addendum_to: REQ-324
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: [REQ-336]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-23T20:16:43Z
  basis:
    - Route B
    - 1-file write set
    - 3 acceptance criteria
    - browser evidence
    - cross-route regression gates
    - full-suite verification
route: B
related: [REQ-336, REQ-338]
batch: timeline-click-regression
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# A Check That Can Catch Timeline Click Retargeting

## What

Give the Timeline probe lane a check that fails on the pre-REQ-336 behaviour (pointer capture on
pointerdown retargeting the synthesized click) and passes after REQ-336's fix. The existing
lock-in test passes against the broken build, so the lane currently cannot catch this class of
regression.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md` (REQ-245, REQ-252, REQ-324 lesson links),
  `general.md`, `coding-guardrails.md`, `communication-style.md`, `testing.md`.
  Structural, per the REQ's second accepted mechanism, in three assertions:
  1. Derive the set of named functions in the generated page whose own body requests capture,
     matching `setPointerCapture(` **with** the paren → verify: the set is non-empty (vacuity
     guard) and contains `capturePanPointer`.
  2. Assert the `pointerdown` body contains no direct request and calls none of that set →
     verify: mutations M1, M2 and M3 each turn it red.
  3. Assert the `pointermove` body still reaches a request → verify: M4 (capture deleted) turns it
     red, so the absence half cannot be satisfied by removing capture.
  Then confirm REQ-324's lock-in stays green through every mutation — that gap is the REQ's whole
  premise, and it should be measured rather than quoted.
- [x] **[APPLY]:** `timeline_browser_probe_test.go` only: `pointerCapturingFunctionNames` plus
  `TestTimelinePointerCaptureWaitsForThePanEngage`. The four mutations were applied to
  `web/board-timeline.js` and reverted; `diff -q` against the pre-mutation copy confirms it is
  byte-identical, and it is not in this REQ's staged set.
- [x] **[UNIFY]:** `git diff --stat` → 1 file, +109/-0, the declared one.
  - `timeline_browser_probe_test.go` — read the whole hunk: one helper, one test, both with the
    reasoning that makes each assertion non-vacuous. `gofmt -l` prints nothing; `go vet` clean via
    the gate.
  - `git status` shows `web/board-timeline.js` unmodified, so no mutation leaked into the tree.
  - No debug artifacts: no `t.Log`-only paths, no skipped assertion, no `TODO`/`FIXME` in the diff.
  - Native lint/tests: board suite with the browser lane enabled → ok; canonical gate → exit 0.

## Why

`TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving`
(`timeline_browser_probe_test.go:626`) is REQ-324's lock-in test, and it passed while every real
mouse click in the Timeline was broken. Reason: it dispatches synthetic PointerEvents with
`pointerId: 1`, an id the engine does not know, so `setPointerCapture` throws and capture is
never established — the probe never exercises the code path that breaks real clicks. A lock-in
test that cannot catch the regression it pins is not doing its job.

## Detailed Requirements

- The check must fail on the broken behaviour (capture established on pointerdown, click
  retargeted to `#timeline-scroll`) and pass on the fixed behaviour (REQ-336).
- Either drive real input, or assert the structural property directly (where the capture call
  sits relative to the engage threshold `TIMELINE_PAN_THRESHOLD_PX`).
- **Mutation-test whichever you pick: a guard no mutation can break is dead code.** At minimum,
  reintroducing capture-on-pointerdown (reverting REQ-336's fix) must turn the check red.

## Builder Guidance

Certainty: Firm on the goal, open on the mechanism (real input vs structural assertion — the
user left the choice to the builder). Ordering was decided at capture and user-confirmed at
verify (2026-08-23): this REQ depends on REQ-336 so the committed suite never carries a red
check; prove the RED side by mutation (reintroduce the broken capture placement locally and
watch the check fail), not by committing against broken HEAD.
`tdd: false` because the deliverable IS the check — the mutation evidence is its proof, and the
work loop's test-first gate has no separate implementation to precede. Cautionary precedent from
REQ-333's UNIFY notes: its structural assertion first matched the `typeof` guard beside the
capture call instead of the call itself, and only mutation testing exposed that — expect the
same trap here.

## Red-Green Proof

**RED prompt/case:** Check out the pre-REQ-336 behaviour (or mutate the fix back to
capture-on-pointerdown) and run the new check: it must fail. On current HEAD before this REQ,
`TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving` passes against that same broken
behaviour — that gap is what this REQ closes.
**Why RED now:** The lane's only click/pan lock-in test never establishes pointer capture
(synthetic `pointerId` throws), so click retargeting is invisible to it.
**GREEN when:** The new check passes on post-REQ-336 HEAD, demonstrably fails when
capture-on-pointerdown is reintroduced (mutation evidence recorded in the REQ), and
`bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** User confirmed (goal and both accepted mechanisms stated verbatim in the input)

## Prior Implementation

REQ-324 ("Give the timeline drag a movement threshold", archived, completed, commit 3486ab2)
added `TIMELINE_PAN_THRESHOLD_PX` (4px) so a press that does not move is still a click, and
pinned it with `TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving`
(`timeline_browser_probe_test.go:626`). REQ-333 (commit 36c4518) later noted in its UNIFY that a
synthetic `pointerId` cannot be captured in this lane, so it asserted the capture call
structurally and drove `lostpointercapture` directly — the workaround whose blind spot this REQ
closes.

## Dependencies

`depends_on: [REQ-336]` — the fix lands first; this check then pins it, proving RED by mutation.

## Notes

The input tagged this `[impact-high]`, which is not a canonical `impact:` value; normalize-and-warn
fired and the user confirmed `impact-user-visible` during capture (default tier, so the title
carries no tag).

---
*Source: UR-067 — see `do-work/user-requests/UR-067/input.md` for complete verbatim input.*

---

## Triage

**Route: B** - Medium

**Reasoning:** One file, one deliverable, and the mechanism choice is stated. But which structural property to assert — and how to make it survive mutation rather than matching a neighbouring line the way REQ-333's first attempt did — has to be derived from the shape of the code under test.

**Planning:** Not required

## Exploration

**Why the existing lock-in cannot be extended.** `TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving`
(`timeline_browser_probe_test.go:626`) drives the real board with synthetic `PointerEvent`s carrying
`pointerId: 1` (`:641`). `setPointerCapture` throws `NotFoundError` on an id the engine does not
know, so capture is never established and the retargeting path is unreachable from inside this
lane. REQ-333 hit the same wall and said so at `:2192-2196`.

**Why the lane cannot be given real input cheaply.** `browser_probe_test.go:144-155` launches the
engine with `--headless --dump-dom` and reads one result node out of the serialized DOM. There is
no protocol channel, so nothing in the lane can dispatch trusted input. Adding one means a CDP
transport (`--remote-debugging-pipe` plus a JSON framing) inside the test binary — a new capability,
not an extension. The REQ explicitly offers the structural route as the alternative, so that is the
one taken (D-01).

**What the structural property has to be, to survive mutation.** Three separate mutations
reintroduce the defect, and a check that only looks for the literal `setPointerCapture(` inside the
`pointerdown` body catches one of them:

1. a direct `scrollHost.setPointerCapture(...)` back on `pointerdown`;
2. a call to the existing `capturePanPointer()` from `pointerdown`;
3. a *new* named wrapper that captures, called from `pointerdown`.

REQ-336's retargeted REQ-333 assertion catches none of them — verified, not assumed: under
mutation 2 that test still passes (its two halves stay true). So the check has to resolve the
capturing functions **out of the page** rather than from a hand-written name list, which is the
board prime's *Closed Enumerations Go Stale* rule applied to a test.

**And the absence has to be paired with a presence.** "No capture on `pointerdown`" is trivially
satisfied by capturing nowhere, which is REQ-333's original bug — a drag released outside the chart
never telling the host it ended. The board prime's REQ-245 lesson is the same shape: asserting a
phrase is absent is not a guard. So the third assertion pins that the *engage* path still reaches a
capture request.

**The trap REQ-333's Builder Guidance warned about, confirmed.** Matching the bare name
`setPointerCapture` also matches the `typeof scrollHost.setPointerCapture === "function"` feature
detect that sits beside every call. The match includes the opening paren for exactly that reason.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify) — one new test
  plus one helper that derives the capturing-function set from the generated page.

**Files I will NOT touch:** `web/board-timeline.js` (REQ-336 delivered the fix; this REQ only pins
it — the mutations below are applied and reverted, never committed), `browser_probe_test.go` (a CDP
transport is out of scope, see D-01), `web/board-controls.js`.

**Acceptance criteria (restated from REQ):**
- [ ] The check fails on the broken behaviour and passes on the fixed behaviour
- [ ] Mutation-tested: reintroducing capture-on-pointerdown turns it red
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0

## Decisions

- **D-01 — ESCALATE. Structural assertion, not real input.** The REQ left the mechanism open. Real
  input would need a CDP transport inside the test binary (`--remote-debugging-pipe` plus JSON
  framing), because this lane is `--dump-dom` only; that is a new capability in a file this REQ does
  not declare, and the REQ names the structural route as an accepted alternative.
  **Value:** the check lands now, in the same batch as the fix, and its mutation evidence is
  stronger than a behavioural check's would be — four distinct reintroductions of the defect turn it
  red, including one through a wrapper that does not exist yet.
  **Risk:** it reads text. A capture routed through a variable, a method lookup on a computed name,
  or an `eval` would pass. Named in the test's own doc comment as a known residual, and a CDP lane
  is raised as a discovered task rather than lost. Reversible — a behavioural check can replace this
  one without touching the fix.
- **D-02 — DECIDE & STATE. The capturing-function set is derived from the page, not listed.** A
  hand-written `[]string{"capturePanPointer"}` would have passed mutation M3 (a fresh wrapper called
  from `pointerdown`), which is the exact failure mode the board prime's *Closed Enumerations Go
  Stale* rule describes. The helper walks every `function NAME(` in the generated page and keeps the
  ones whose own body requests capture.
- **D-03 — DECIDE & STATE. Anonymous function expressions are skipped by the helper.** They have no
  name, so no other handler can route a request through one; a `pointerdown` handler that captured
  inside its own nested anonymous function is caught by assertion (a) instead, because that body is
  inside the sliced `pointerdown` block.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)

**What was done:** Added `TestTimelinePointerCaptureWaitsForThePanEngage`, which asserts over the
generated page that the Timeline's `pointerdown` handler neither requests pointer capture itself nor
calls any function that does, while the `pointermove` handler still reaches a request. The
capturing-function set is derived from the page by a new `pointerCapturingFunctionNames` helper
rather than hand-listed, and a vacuity guard fails the test if the page requests capture nowhere.
`web/board-timeline.js` was mutated four ways to prove the check bites and then restored unchanged.

## Qualification

Passed — 1 file verified, 3 requirements traced, P-A-U confirmed.

- Mechanical: `tools/checks/qualify.sh` exit 0; `tools/checks/scope-drift.sh` exit 0.
- Substantive (check 2): +115 lines, all executable test logic and the reasoning behind each
  assertion. No skipped subtest, no assertion behind a condition that cannot be true.
- Requirements traced (check 3): "fails on broken, passes on fixed" and "mutation-tested" are the
  table in `## Testing`; the canonical gate exits 0.
- Data flows (check 6): the check reads the *generated* `index.html` through
  `generateLiveSiteInDir`, not the source file — so it measures what ships, and a page that failed
  to embed the module fails the vacuity guard rather than passing on an empty string.
- Vacuity: proved by mutation rather than argued. Every one of the four mutations turns the check
  red, so none of its three assertions is dead.
- Contamination check: the previous REQ (REQ-336) touched `web/board-timeline.js` and this same
  test file. The overlap is expected and declared — REQ-336's D-02 records it — and this REQ's diff
  adds a new test rather than editing REQ-336's hunk. `web/board-timeline.js` is unmodified in this
  REQ's tree.

## Testing

**Tests run:**
- `go test -count=1 ./...` in `skills/do-work-board/tools/queue-kanban` with the browser lane
  enabled
- `bash _dev/tests/maintainer-verify.sh` from the repo root against the final tree
- Four mutations of `web/board-timeline.js`, each applied, measured, and reverted

**Result:** ✓ the new check passes on fixed HEAD; the canonical gate exits 0.

**Red-green validation (mutation evidence — this is the deliverable's proof):**

| mutation of `web/board-timeline.js` | new check | REQ-324's lock-in |
|---|---|---|
| none — REQ-336's fix in place | ✓ PASS | ✓ PASS |
| M1 direct `setPointerCapture()` back on `pointerdown` | ✗ FAIL | ✓ PASS |
| M2 `capturePanPointer()` called from `pointerdown` | ✗ FAIL | ✓ PASS |
| M3 a fresh named wrapper that captures, called from `pointerdown` | ✗ FAIL | ✓ PASS |
| M4 the capture request removed entirely | ✗ FAIL | ✓ PASS |

The right-hand column is the REQ's premise, measured rather than quoted: REQ-324's lock-in passes
through every reintroduction of the defect, including the one that broke every click in the chart.
M3 is what the derived function set buys — a hand-written name list would have passed it. M4 is
what the third assertion buys.

Also measured: under M2, REQ-336's retargeted REQ-333 assertion still **passes** (both its halves
stay true), so the existing structural check cannot catch a capture routed through the wrapper. That
is what assertion (b) is for.

**New tests added:**
- `TestTimelinePointerCaptureWaitsForThePanEngage`
- `pointerCapturingFunctionNames` (helper, derives the capturing set from the generated page)

**Existing tests updated (cross-REQ impact):** none. REQ-324's and REQ-333's tests are untouched by
this REQ; REQ-336 already retargeted the one assertion that had to move.

**`tdd: false` as the REQ states** — the deliverable *is* the check, so there is no implementation
for a test to precede; the mutation table is its red-green evidence.

*Verified by work action*

## Review

**Overall: 96%** | 2026-08-23T20:25:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 1 (report only) — the check reads text, so a capture routed through a variable,
a computed method lookup, or an `eval` passes it. Stated in the test's own doc comment as a known
residual and raised as a discovered task (a CDP lane), so it is neither hidden nor lost.

**Acceptance:** Pass — passes on fixed HEAD, fails under all four reintroductions of the defect,
canonical gate exit 0.
**Suggested testing:** 2 items
**Follow-ups created:** None; **sweeps appended to:** None

### Requirements Checklist

- [x] Fails on the broken behaviour, passes on the fixed one — delivered (M1-M3 red, baseline green)
- [x] Either real input or the structural property asserted directly — delivered (structural, D-01)
- [x] Mutation-tested; reintroducing capture-on-pointerdown turns it red — delivered, four
      mutations, and the pairing assertion (c) proves the absence half is not satisfiable by
      deleting capture
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0 — delivered
- [x] RED proved by mutation, not by committing against broken HEAD — delivered; the mutations were
      reverted and `web/board-timeline.js` is unmodified in this REQ's diff
- [x] REQ-333's `typeof`-guard trap avoided — delivered; the match includes the opening paren, and
      the reason is written beside it

### Restatement Sweep

The diff adds a test; it redefines nothing. The one statement it touches is the claim that this lane
cannot establish pointer capture (`timeline_browser_probe_test.go:2192-2196`, REQ-333) — the new
test's doc comment restates it as the reason for going structural, and it is still true. Checked
that no shipped prose or prime states how many probes the Timeline lane has or what it covers:
`_dev/primes/prime-kanban-board.md` describes the lane's *posture* (render and measure) without
counting, and the runner's figure is a derived count. Nothing to update.

### Code Review Notes

- **Naming for reach:** two new identifiers — `pointerCapturingFunctionNames` and
  `TestTimelinePointerCaptureWaitsForThePanEngage`. Both multi-word, both say what they hold.
  `capturingNames`, `capturingName`, `declarationIndex` are short-lived locals within a screen of
  their use.
- **The helper's loop is bounded and cannot spin:** `searchOffset` advances past the declaration
  token on every iteration whether or not a name is found, and the two `break`s cover the
  no-more-matches and malformed-tail cases.
- **`sliceBalancedBlockAfter` is called on a suffix slice** (`pageSource[declarationIndex:]`) so its
  anchor search cannot re-find an earlier declaration of the same name and return the wrong body.
- **The `strings.ContainsAny` filter** rejects anonymous expressions and anything that is not a bare
  identifier, so a `function (event) {` cannot enter the set as a nameless entry and make assertion
  (b) match every `(` in the page.
- **Test Adequacy is 95%, not 100%,** for the residual in the Minor finding: the check is as strong
  as a text-level check can be here, and that ceiling is real.

### Acceptance Testing

**Result: Pass**
- Ran the check on fixed HEAD: PASS in 0.9s.
- Applied four mutations covering both spellings of the defect, a wrapper that does not exist in the
  tree, and the opposite failure (no capture at all). Each turns the check red with a message naming
  what it found.
- Confirmed `web/board-timeline.js` is byte-identical to its pre-mutation copy (`diff -q`) and shows
  no modification in `git status`, so nothing leaked into the tree.
- Full board suite with the browser lane enabled, and the canonical gate: both clean, so the new
  test does not disturb the 13 existing Timeline probes.

### Suggested Additional Testing

- Run the check against a build where the module is minified or the function keyword is elided
  (arrow functions). The board ships the JS verbatim today, so the helper's `function NAME(` shape
  holds; if that ever changes, the vacuity guard fires rather than the check silently passing —
  worth confirming once by hand when it does.
- A behavioural version of this check, once the lane can dispatch trusted input. That is the
  discovered task below; the CDP prototype used for REQ-336's RED is the starting point.

*Reviewed by review-work action*

## Discovered Tasks

- **The Timeline probe lane cannot dispatch trusted input, and that is why two REQs in this batch
  had to work around it.** `browser_probe_test.go:144-155` runs the engine with
  `--headless --dump-dom` and reads one result node, so every probe drives synthetic events. A
  synthetic `pointerId` cannot be captured, which is why REQ-324's lock-in passed through the whole
  click regression, why REQ-333 fell back to a structural assertion, why REQ-336's RED had to be
  reproduced outside the suite, and why this REQ's check is structural. Giving the lane a CDP
  transport (`--remote-debugging-pipe`, JSON framing, `Input.dispatchMouseEvent`) would let all four
  be behavioural. A working prototype exists — REQ-336's RED harness drove exactly these gestures
  over CDP from Node in about 150 lines, and its two hard-won measurement rules (scroll the target
  into the viewport before dispatching; measure "it panned" from the engage, not the axis text)
  belong in whatever lands. Judged `impact-rule-change`: it changes what the lane can be asked to
  prove, rather than fixing a user-visible defect.

## Lessons Learned

**What worked:**
- Deriving the capturing-function set from the page instead of naming it. Mutation M3 — a wrapper
  that does not exist in the tree — is the one a hand-written list would have passed, and it is also
  the most likely real regression, because "extract the capture into a helper and call it earlier"
  is a plausible refactor.
- Pairing the absence assertions with a presence one. "No capture on pointerdown" and "capture
  somewhere" are two different requirements, and satisfying the first by deleting the second
  reintroduces REQ-333's bug. M4 is the mutation that proves the pairing earns its place.
- Measuring the REQ's premise rather than quoting it. Running REQ-324's lock-in under every mutation
  turned "it passed while clicks were broken" from a claim in the REQ into a row in a table — and
  the same run showed REQ-336's retargeted assertion cannot catch M2 either, which is what justified
  assertion (b).

**What didn't:**
- Nothing was tried and abandoned here; the REQ's Builder Guidance pointed at the `typeof`-guard
  trap and the mutation-testing requirement up front, and both landed first time. Worth recording
  that the guidance is the reason: an unwarned attempt would very likely have matched the feature
  detect, exactly as REQ-333's did.

**Worth knowing:**
- **`setPointerCapture` appears twice at every call site** — once in the `typeof` feature detect and
  once in the call. Any text-level assertion about it must include the opening paren, or it matches
  the guard and passes with the call deleted. This has now cost two REQs; the paren is the whole fix.
- **A structural check needs both a vacuity guard and a mutation table to be worth anything.** The
  guard says the thing being asserted about exists; the table says each assertion can fail. Either
  alone leaves a check that reads as coverage.
- The lane's real limitation is one sentence: `--dump-dom` gives no protocol channel, so no trusted
  input, so no captured pointer. Every workaround in this area traces to it, which is why it is a
  discovered task rather than another local note.

## Orientation

The Timeline probe lane can now catch a pointer capture taken before the pan engages — the
regression that broke every mouse click in the chart and that the lane's existing click lock-in
passed straight through. Lives in the board's Timeline probe suite
(`skills/do-work-board/tools/queue-kanban`); no production code changed, so nothing about the
system's shape moved. `_dev/primes/prime-kanban-board.md` gains one lesson link and its referenced
paths all still resolve.
