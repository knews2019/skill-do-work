---
id: REQ-345
title: "[impact-critical] Stop adding a queued REQ from failing the timeline landing probe"
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
