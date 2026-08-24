---
id: REQ-370
title: "[impact-negligible] Review fix: restore a falsifiable Timeline pointer-capture mutation"
status: completed
claimed_at: 2026-08-24T21:49:13Z
completed_at: 2026-08-24T22:11:03Z
commit: 46ed690a3d68b836363ed3f093da19de1d8ec8ea
status_changed_at: 2026-08-24T22:11:03Z
route: C
created_at: 2026-08-24T17:35:53Z
user_request: UR-067
addendum_to: REQ-341
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-341]
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-08-24T21:49:13Z
  basis:
    - Route C
    - 1-file write set
    - 4 acceptance criteria
    - browser evidence
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Review Fix: Restore a Falsifiable Timeline Pointer-Capture Mutation

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reproduced the retained-Chromium RED three times, instrumented the trusted event boundary, and identified the host-targeted `pointerleave` as the alternate teardown path before selecting a symmetric test-only isolator.
- [x] **[APPLY]:** Added narrow host-leave isolation to both outside-release trials, direct capture and leave observations, cleanup fallbacks, and non-vacuity assertions without changing runtime code or existing product assertions.
- [x] **[UNIFY]:** Reviewed the exact one-file diff, removed temporary instrumentation, applied formatting, and passed repeated named-browser runs, both inverse mutations, strict browser, module vet/tests, and browser-enabled canonical verification.

## What

Restore a non-vacuous capture-present/capture-suppressed mutation pair for the Timeline's trusted
outside-release browser probe under current Chromium. Measure which boundary event now clears the
pan before choosing the replacement gesture; do not weaken the product assertion.

## Context

Found during independent review of REQ-351. The clean click and early-capture mutation remain useful,
but the outside-release trial ends with `panning=false` both when capture is present and when
`setPointerCapture` is swallowed, so one half of the mutation pair can no longer falsify the other.
This is distinct from REQ-369's fixed-delay table-rebuild race.

## Requirements

- Reproduce `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` under the currently
  retained Chromium and identify the boundary event that clears pan state in the outside-release
  trial.
- Replace or recalibrate only that gesture/mutation boundary so capture-present and
  capture-suppressed trials demonstrably depend on one another again.
- Preserve the clean click, early-capture drawer-blocking, trusted CDP input, and product assertions;
  do not turn a failed dependency assertion into an unconditional expected value.
- Run the probe repeatedly and prove the mutation pair can fail when either required mechanism is
  suppressed.

## Red-Green Proof

**RED prompt/case:** Run `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` with the
current retained Chromium. Trial 4 reports `panning=false` even after `setPointerCapture` is swallowed,
so it cannot falsify trial 3's outside-release capture dependency.

**Why RED now:** Current Chromium clears the selected gesture's pan state for another reason; the
test still names a capture dependency that its mutation no longer demonstrates.

**GREEN when:** the trusted browser probe passes with a capture-present/capture-suppressed pair whose
opposite outcomes are both observed, and targeted mutations prove either side can make the probe fail.

**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---
*Source: REQ-351 independent review; pre-existing Timeline probe failure outside REQ-351's write set.*

## Triage

**Route: C** — This is a trusted-input, browser-version-sensitive pointer lifecycle ratchet. The one test file is known, but the retained engine's actual boundary event must be instrumented, a genuinely opposite capture-present/suppressed gesture designed, both dependency directions mutated, and the clean-click/early-capture product assertions preserved.

## Plan

Measure the capture-suppressed event sequence before editing. Isolate only the confirmed alternate release boundary equally in both outside-release trials, then prove the pair differs solely on whether pointer capture retargets the outside release to the host. Keep the existing clean-click, early-capture, trusted-input, drawer, and pan assertions unchanged and replay inverse mutations for both required mechanisms.

## Exploration

- The retained Chromium 1228 baseline reproduced RED 3/3: both outside-release trials ended with `panning=false`.
- In the capture-suppressed trial, document capture observed a `pointerleave` targeted at the rows SVG, then a host-targeted `pointerleave`; the product's host listener immediately cleared pan before later pointer movement and release targeted the page header.
- The host leave is legitimate defensive runtime behavior, but it makes the selected gesture unable to demonstrate capture independently. A document capture listener can isolate exactly the leave whose target is the Timeline scroll host while allowing nested and unrelated leaves through.
- Installing that isolator symmetrically means the capture-present trial must end through a capture-retargeted `pointerup`, while the capture-suppressed trial has no host release path and must remain panning.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`

**Acceptance criteria:**

- The current-browser boundary event that vacuously clears the capture-suppressed pan is directly measured and documented.
- Capture-present and capture-suppressed outside-release trials observe opposite pan outcomes under the same narrow boundary isolation.
- Capture delivery, capture suppression, and exercise of the isolated host leave are asserted directly so neither half can pass vacuously.
- Existing clean click, early-capture drawer blocking, trusted CDP input, product assertions, and runtime code remain unchanged; repeated probe, inverse mutations, strict lane, and canonical gates pass.

## Decisions

- **D-01 — Keep the product's multi-boundary teardown.** Current Chromium's host leave is valid defensive behavior; the test isolates it only while proving the separate capture mechanism.
- **D-02 — Apply the isolator to both paired trials.** This holds the gesture boundary constant so capture is the only difference between clean teardown and stranded pan.
- **D-03 — Observe mechanisms directly.** Count `gotpointercapture` and the exact isolated host leave, then reject missing or unexpected observations before accepting pan state.
- **D-04 — Preserve cleanup fallbacks.** Restore both the event listener and element method through explicit teardown plus test cleanup so a failing assertion cannot contaminate later browser work.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified): records actual capture and host-leave delivery, symmetrically isolates only the host-targeted current-Chromium boundary during paired outside-release trials, asserts both mechanisms are exercised, and returns the full drag outcome for non-vacuity checks.

## Discovered Tasks

None.

## Testing

- Baseline retained Chromium reproduced the vacuous capture-suppressed outcome 3/3.
- The final named probe passed 10/10 and another 3/3 after the boundary-counter precision adjustment; each run observed one capture in trial 3 and one isolated host leave with no capture in trial 4.
- Disabling boundary isolation returned RED on the non-vacuity assertion; disabling capture suppression returned RED after observing real capture.
- The strict retained-browser lane, Go vet, uncached full module tests, and browser-enabled canonical verification all passed in the builder worktree.
- On merged main, the named retained-Chromium probe passed another 5/5 with one capture event and one isolated leave per run; the complete strict browser lane passed in 16.69 seconds.
- The final browserless canonical gate passed all contract suites, queue-kanban vet/tests, strict JavaScript, and audit-metrics verification; its optional browser step skipped only after the explicit merged-main strict lane passed.

## Qualification

- Exact merge range `788db0ab70b7897cfb150aaf590d086a1f4ea86e..46ed690a3d68b836363ed3f093da19de1d8ec8ea` passed mechanical qualification.
- Scope drift passed: the sole changed test file exactly matches the declared Scope and Implementation Summary.
- Orchestrator judgment confirmed substantive symmetric boundary isolation, direct mechanism observations, preserved product assertions, and no runtime, generated, debug, or release artifacts in the implementation merge.

## Review

Independent review approved at 100/100 with no Important, Minor, or Nit findings and low test-only residual risk. It independently reproduced the base false/false outcome, repeated the final opposite outcomes, passed the strict lane, and proved three mutations RED: disabling the isolator, leaving capture unsuppressed, and removing the product capture call.

## Lessons Learned

A browser can add a valid alternate teardown path without changing the product contract, yet make a dependency probe vacuous. Measure the event order, isolate the confounding path symmetrically, and directly assert every mechanism the paired outcomes depend on.

## Orientation

Released in 0.236.60. The trusted Timeline outside-release probe now demonstrates capture independently under retained Chromium instead of passing both sides through the host's alternate leave teardown.
