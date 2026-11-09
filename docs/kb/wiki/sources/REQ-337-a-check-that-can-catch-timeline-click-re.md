---
title: "Lessons from REQ-337: A check that can catch Timeline click retargeting"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-01/REQ-337-a-check-that-can-catch-timeline-click-re.md]
related:
  - page: REQ-336-timeline-clicks-open-the-detail-drawer-a
    rel: depends-on
  - page: REQ-338-cut-the-timeline-row-list-to-one-tab-sto
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-337: A check that can catch Timeline click retargeting

Part of the [[concept-duration-estimation-and-breaks]] cluster.

## What the REQ was about

Give the Timeline probe lane a check that fails on the pre-REQ-336 behaviour (pointer capture on
pointerdown retargeting the synthesized click) and passes after REQ-336's fix. The existing
lock-in test passes against the broken build, so the lane currently cannot catch this class of
regression.

## Solution summary

Added `TestTimelinePointerCaptureWaitsForThePanEngage`, which asserts over the
generated page that the Timeline's `pointerdown` handler neither requests pointer capture itself nor
calls any function that does, while the `pointermove` handler still reaches a request. The
capturing-function set is derived from the page by a new `pointerCapturingFunctionNames` helper
rather than hand-listed, and a vacuity guard fails the test if the page requests capture nowhere.
`web/board-timeline.js` was mutated four ways to prove the check bites and then restored unchanged.

## What worked

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

## What didn't work

- Nothing was tried and abandoned here; the REQ's Builder Guidance pointed at the `typeof`-guard
  trap and the mutation-testing requirement up front, and both landed first time. Worth recording
  that the guidance is the reason: an unwarned attempt would very likely have matched the feature
  detect, exactly as REQ-333's did.

## Worth knowing

- **`setPointerCapture` appears twice at every call site** — once in the `typeof` feature detect and
  once in the call. Any text-level assertion about it must include the opening paren, or it matches
  the guard and passes with the call deleted. This has now cost two REQs; the paren is the whole fix.
- **A structural check needs both a vacuity guard and a mutation table to be worth anything.** The
  guard says the thing being asserted about exists; the table says each assertion can fail. Either
  alone leaves a check that reads as coverage.
- The lane's real limitation is one sentence: `--dump-dom` gives no protocol channel, so no trusted
  input, so no captured pointer. Every workaround in this area traces to it, which is why it is a
  discovered task rather than another local note.

## Back-reference

See `do-work/archive/UR-067/REQ-337-check-that-catches-timeline-click-retargeting.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `77ca13e`.
