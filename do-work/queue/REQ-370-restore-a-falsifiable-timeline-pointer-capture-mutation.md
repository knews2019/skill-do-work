---
id: REQ-370
title: "[impact-negligible] Review fix: restore a falsifiable Timeline pointer-capture mutation"
status: pending
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
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Review Fix: Restore a Falsifiable Timeline Pointer-Capture Mutation

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
