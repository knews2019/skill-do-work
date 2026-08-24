---
id: REQ-345
title: "[impact-critical] Stop adding a queued REQ from failing the timeline landing probe"
status: completed
completed_at: 2026-08-24T10:42:00Z
claimed_at: 2026-08-23T23:05:00Z
created_at: 2026-08-23T22:35:07Z
user_request: UR-068
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec: bug-fix
route: C
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
---

# Stop Adding a Queued REQ From Failing the Timeline Landing Probe

## What

Adding pending REQs to the queue makes `TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable`
fail, which red-lights the canonical gate. Capturing three REQs did it; removing those three files
makes it pass again. Either the step-forward arrow's enablement is wrong for this data or the probe's
premise is, and this REQ is to establish which and fix that side.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-kanban-board.md`. Reproduced the probe's exact five clicks and read what is drawn in the next period BEFORE changing anything, per the REQ's first requirement.
- [x] **[APPLY]:** One file. `web/board-timeline.js` was in the write set and deliberately left untouched — the measurement said the board is correct.
- [x] **[UNIFY]:** Audited by the orchestrator against `bfafab1..1ee071d`: confirmed only the probe file changed, then mutated `board-timeline.js` in an isolated worktree to drop REQ-329's guard and confirmed the rewritten assertion still fails with exit 1. Restored byte-identical.

## Why

`bash _dev/tests/maintainer-verify.sh` is the only proof this project accepts before a hand-back, and
right now an ordinary `do-work capture-request` turns it red. That blocks every commit behind an
unrelated action: the person who captures a request has no reason to be debugging a browser probe,
and the obvious way out — deleting the captured REQ — throws away the user's request.

It is `impact-critical` for that reason rather than for anything the board shows a reader: the
project's only accepted gate can be failed by adding data the queue exists to hold.

## Context

**Reproduced, both directions, on 2026-08-23 at about 22:40 UTC:**

- With `REQ-342`, `REQ-343` and `REQ-344` present in `do-work/queue/`, the gate fails:

      timeline_browser_probe_test.go:1905: the step-forward arrow is enabled on the current week
        (2026-08-17 00:00 UTC → 2026-08-24 00:00 UTC), whose next period exists only inside the
        cosmetic bound padding
      timeline_browser_probe_test.go:1909: pressing the step-forward arrow moved the window from
        2026-08-17 00:00 UTC → 2026-08-24 00:00 UTC to 2026-08-24 00:00 UTC → 2026-08-25 18:24 UTC,
        past everything drawn

- Move those three files out of the tree and the same test passes. Move them back and it fails again.

**What is NOT established.** Whether the board or the probe is wrong. The generated payload's latest
stamps ran to `2026-08-24T00:58:49Z` — past the current week's end — which would mean the next period
holds real forecast marks and the arrow is right to be enabled. But an independent CDP probe that
clicked `[data-timeline-period="week"]` after `timeline-zoom-now` landed on a fourteen-hour window
(`2026-08-23 16:42 → 2026-08-24 06:58`) with the arrow already disabled and the step a no-op — it did
not reproduce the test's own week-long window at all. So the payload reading is a hypothesis, not a
finding, and the first job here is to reproduce the probe's exact window before changing anything.

The probe reaches that state through five clicks in sequence (`:1802-1812`): `timeline-zoom-now`,
`timeline-zoom-in`, `[data-timeline-period="week"]`, then `timeline-period-next`. The `zoom-in` in the
middle is the step the independent probe omitted and is the likeliest reason the windows differed.

## Detailed Requirements

- Establish which side is wrong, with evidence, before changing either. The three candidates are the
  arrow's enablement rule, the period button's window computation, and the probe's premise.
- Adding or removing an ordinary pending REQ must not change whether the gate passes.
- Whatever the verdict, the assertion must still be able to fail. If the premise becomes
  data-dependent, **read the condition rather than restating it** — the board prime's REQ-322 lesson
  ("a constant a decision turns on must be READ by the test, never restated beside it") is the whole
  risk here, because the cheap fix is to relax the assertion until it cannot fire.
- The fix must hold across the period boundary, not just away from it. This surfaced roughly an hour
  before the week rolled over, and a fix verified at midday proves nothing about that.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first, including the
  render-evidence rule: return `location.href` alongside every measurement.
- Do not weaken or delete the assertion to get green — the two clauses it holds ("a step past
  everything drawn does not happen, and says so first") are REQ-329's contract.
- This probe drives the repo's own live queue through `generateLiveSiteInDir`, so its inputs move
  under it. Consider whether that is the deeper defect, but do not rewrite the whole probe here.

## Builder Guidance

**Certainty: firm that the trigger is real and reproducible; the mechanism is unestablished.** Do not
inherit the hypothesis in `## Context` — it is written down so it can be checked, not believed. Start
by reproducing the probe's exact window (all five clicks, in order) and reading what is actually drawn
in the next period.

Time is an input here. Record the wall-clock instant with every measurement, and re-run at least once
close to a period boundary; a green run at an arbitrary hour is not evidence.

## Open Questions

None — the diagnosis is the REQ's first requirement, not a question for the user.

## Red-Green Proof

**RED prompt/case:** With three pending REQs freshly captured into `do-work/queue/`, run
`bash _dev/tests/maintainer-verify.sh`: it fails in
`TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable` with the two lines quoted above.
Remove the three files and it passes.

**Why RED now:** Unestablished — either the step-forward arrow enables on a period that holds nothing
drawn, or the probe asserts the next period is always empty when the queue's forecast can reach into
it.

**GREEN when:** The gate passes with those three REQs in the queue and still passes without them; the
assertion still fails when the behaviour it guards is genuinely broken (shown by mutation, since
`tdd: false` here means the deliverable is a diagnosis plus a fix rather than a new test); and a run
close to a period boundary is among the evidence.

**Validation:** Inferred during capture — a defect found while capturing UR-068, not a user request.
The reproduction in both directions is recorded above; the mechanism is not.

## Assets

None. Gate output quoted in `## Context`.

---
*Source: found while capturing UR-068 — the capture's own three REQ files are the reproduction.*

---

## Triage

**Route: C** - Complex

**Reasoning:** The failure was reproducible but its mechanism was unestablished, and three candidate causes were named. The REQ's first requirement was a diagnosis, not a fix.

**Planning:** Required — the plan was the measurement.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify)

**Files I will NOT touch:** `web/board-timeline.js` — in the declared write set, but the measurement established the board is correct, so changing it to satisfy a wrong assertion would have been the cheap fix inverted (D-05).

**Acceptance criteria (restated from REQ):**
- [x] Establish which side is wrong, with evidence, before changing either
- [x] Adding or removing an ordinary pending REQ must not change whether the gate passes
- [x] The assertion must still be able to fail; read the condition rather than restating it
- [x] The fix holds across the period boundary, not only away from it

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)

**What was done:** Clause (3) was rewritten to read the step-forward arrow's own reported state and check the press against it — disabled means the press must not move the window, enabled means it must move it AND land on drawn segments. The unconditional "the next period is empty" claim moved to a new clause (6), anchored on a filtered archived REQ whose premise is read from the Fit-all extent rather than restated as a week number. `web/board-timeline.js` is unchanged.

## Testing

**Tests run:** the landing probe directly, and `GOTOOLCHAIN=go1.26.1 QUEUE_KANBAN_BROWSER=... bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Gate exit 0 with the strict browser lane run, not skipped

**Red-green validation:** the probe failed at 2026-08-23T23:15:55Z with the queue populated and passed at 23:25:26Z after the fix, same queue, same clock regime.

**Mutation evidence (three independent mutations, all caught):**
- drop `!stepLandsOffTheData(stepped)` from the availability rule → clause (6), both assertions. Independently reproduced by the orchestrator in an isolated worktree, exit 1.
- drop `movesTheWindow(stepped)` → clause (3)'s **enabled** branch
- measure off-data against the bound padding instead of the drawn extent → clause (6)

**Both branches proven reachable**, each caught in the regime where it applies.

*Verified by work action*

## Review

**Overall: 95%** | 2026-08-24T09:33:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None. The reviewer attacked the conclusion, both clause branches, clause (6)'s premise read and the new constant, and each held.

- **The conclusion is right, verified independently:** `timelineRowSegments` (`:812-814`) unconditionally adds a segment for `projectedRow`, `drawProjectedSegment` emits a real `rect.timeline-segment`, and `drawnExtent` is built from that same set. A projected bar is drawn data by every measure the probe uses.
- **No break the old assertion caught that the new pair misses.** In the disabled regime the new branch is byte-for-byte the old assertions. In the enabled regime the old clause asserted the opposite of correct behaviour, so it caught nothing there. Clause (6) is *stronger* in one respect: a mutation computing `drawnExtent` unfiltered passes the old assertion and fails (6).
- **Clause (6)'s premise read is sound.** `FilteredFit.EndMs` is production output settled through `timelineZoomedWindow` from the same `drawnExtent` the guard consults — nothing re-derived, so neither the REQ-305 trap nor REQ-322's is hit.
- **`readoutTruncationMs` is a real property, not dead slack.** `timelineFormatStamp` floors the readout to the minute; the allowance is applied in the conservative direction, making the premise guard harder to pass. Measured margin today is ~6 days, so it never changes an outcome. Reviewer's verdict on the ESCALATE: keep it.
- **Clock dependency:** clause (6) has none — it anchors on a completed archived REQ, so no now-line enters its extent and its anchor is a window midpoint. Clause (3) is deliberately data-dependent and branches, so the clock chooses the branch, never the verdict.

**Minor findings:** 3 (report only) — today's green run no longer discriminates the fix from the pre-fix probe, because the week rolled over and the current period is now cut short, so acceptance rests on the mutation table rather than the original RED; a stale count restatement at `:1751` → prose backlog; clause (3)'s retained heading line now overstates what clause (3) asserts, since that claim moved to clause (6).

**Acceptance:** Pass — gate exit 0 with the strict browser lane actually run; three independent mutations each caught with a precise message.

**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Reproducing the probe's exact five-click sequence and reading what was actually drawn in the next period, before touching either side. The `## Context` hypothesis (the payload's observed stamps reaching past the week end) turned out half right for the wrong reason — the reach comes from the *projection*, not from observed stamps — and only the measurement separated those.

**What didn't:** Nothing was tried and abandoned, but one thing is worth recording as a near-miss: the cheap fix here was to relax the assertion until it stopped firing, and the REQ named that trap explicitly. The assertion that shipped is stronger than the one it replaced.

**Worth knowing:** The pre-fix probe passes on its own after a week boundary, with no code change — one assertion's verdict inverted across midnight Sunday. Any assertion in this lane phrased about *where* something is, rather than about a relationship between measured quantities, can flip on the clock. That is the class, not the instance.

## Orientation

Adding an ordinary REQ to the queue no longer red-lights the project's only accepted gate. Lives in the board tool's browser probe lane. The Timeline's step-forward arrow behaviour is unchanged — this was a wrong assertion, not a wrong board.
