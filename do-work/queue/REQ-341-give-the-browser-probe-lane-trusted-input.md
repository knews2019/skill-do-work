---
id: REQ-341
title: "[impact-rule-change] Addendum: give the browser probe lane trusted input"
status: pending
status_changed_at: 2026-08-23T22:32:23Z
created_at: 2026-08-23T20:25:04Z
user_request: UR-067
addendum_to: REQ-337
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Give the Browser Probe Lane Trusted Input

## What

The board's browser probe lane runs Chromium with `--headless --dump-dom` and reads one JSON result
node out of the serialized DOM (`browser_probe_test.go:144-155`). There is no protocol channel, so
no probe in the lane can dispatch trusted input — every one of them synthesizes events from page
script. Give the lane a way to drive real input, and convert the probes that need it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Four REQs in a row worked around this, and one shipped a defect because of it.

- **REQ-324's click/pan lock-in passed through the entire click regression.** It dispatches
  `PointerEvent`s with `pointerId: 1`; `setPointerCapture` throws on an id the engine does not know,
  so the code path that broke every real click was unreachable from the test.
- **REQ-333** could not drive a captured drag end-to-end and fell back to a structural assertion
  plus a directly-dispatched `lostpointercapture`.
- **REQ-336's RED had to be reproduced outside the suite**, over the DevTools Protocol, because the
  lane cannot produce a trusted click.
- **REQ-337's check is structural for the same reason** — accepted deliberately, with the residual
  written into the test's doc comment: a capture routed through a variable or a computed lookup
  would pass it.

A lane whose central limitation forces every interaction probe into a workaround is worth fixing
once instead of routing around a fifth time.

## Detailed Requirements

- The lane can dispatch trusted input to a page it rendered, and a probe can read the result.
- The existing `--dump-dom` probes keep working unchanged — this adds a capability, it does not
  replace the transport those probes use.
- The engine-missing SKIP path and the strict lane's refusal-to-skip both keep working, and a
  transport failure fails loudly rather than reading as a probe that measured nothing.
- Convert at least REQ-337's structural check to a behavioural one, and say in that test's comment
  what changed and why. Whether REQ-324's and REQ-333's probes also convert is this REQ's call.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first — in particular the
  render-evidence rule (return `location.href` alongside every measurement) and the
  measured-face-per-browser rule.
- No new module dependency. `--remote-debugging-pipe` with JSON framing over a pair of file
  descriptors needs nothing outside the standard library; a WebSocket client would need one.
- Do not weaken `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes` or the probe counter it
  guards.

## Builder Guidance

**Certainty: firm that the limitation is real and costly, open on the transport.** A working
prototype exists and is worth reading before choosing: REQ-336's RED harness drove press, nudge,
drag, and release-outside gestures over CDP `Input.dispatchMouseEvent` in about 150 lines of Node,
against a board from `generateLiveSiteInDir`. Two measurement rules it learned the hard way must
survive into whatever lands:

- **Scroll the target into the viewport before dispatching.** An element's
  `getBoundingClientRect` is not a clickable coordinate until the element is on screen; the first
  bar in a 332-REQ board sat at y≈1538 in a 900px viewport, and the press landed on `HTML`.
- **Measure "it panned" from the pan engaging, not from the axis text.** A 150px drag clamps at the
  window bound and leaves every label identical — REQ-324's lesson, met again.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-337: the browser probe lane cannot
  dispatch trusted input, which is why REQ-324's lock-in missed the click regression and why
  REQ-333, REQ-336 and REQ-337 each had to work around it. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  - [2026-08-23] User approved via clarify: the probe lane gains a real-input channel and
    REQ-337's structural check becomes behavioural. The user was offered a keep-it-queued-for-
    later option and did not take it, so this is wanted now rather than deferred. Nothing put
    out of scope; whether REQ-324's and REQ-333's probes also convert remains the REQ's own
    call.

## Red-Green Proof

**RED prompt/case:** In the lane today, a probe that calls `setPointerCapture(event.pointerId)` on a
synthetic `PointerEvent` throws `NotFoundError`, so no probe can observe capture retargeting a click.

**GREEN when:** a probe in this lane presses and releases on a Timeline bar with trusted input and
observes the detail drawer opening, fails when capture is taken on `pointerdown`, and the existing
`--dump-dom` probes plus the strict lane's zero-probe refusal all still pass.

**Validation:** Inferred during REQ-337's implementation and review — a Discovered Task, not a user
request.

---
*Source: Discovered Task, REQ-337 (UR-067).*
