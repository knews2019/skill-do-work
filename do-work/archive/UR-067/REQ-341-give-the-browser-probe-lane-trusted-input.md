---
id: REQ-341
title: "[impact-rule-change] Addendum: give the browser probe lane trusted input"
status: completed
completed_at: 2026-08-24T11:20:00Z
claimed_at: 2026-08-24T10:05:00Z
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
route: C
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
- [x] **[PLAN]:** Read `prime-kanban-board.md` and REQ-336's prototype. Chose `--remote-debugging-pipe` with hand-rolled JSON framing over a WebSocket client or a driver, because the REQ forbids a new module dependency.
- [x] **[APPLY]:** Two files, both inside the write set. `go.mod` and `go.sum` are untouched — verified by the orchestrator.
- [x] **[UNIFY]:** Audited by the orchestrator against `b09e3e2..1a5b745`: confirmed no dependency was added, REQ-345's landing-probe clauses survived the merge, and the converted probe catches a variable-routed capture that the old structural grep cannot see.

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

---

## Triage

**Route: C** - Complex

**Reasoning:** A new transport for the probe lane, with the mechanism open (the REQ named a prototype worth reading but left the choice), a hard no-new-dependency constraint, and existing probes that must keep working unchanged.

**Planning:** Required.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modify) — the transport
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify) — the converted probe

**Acceptance criteria (restated from REQ):**
- [x] The lane dispatches trusted input and a probe reads the result
- [x] Existing `--dump-dom` probes keep working unchanged
- [x] The SKIP path and the strict lane's refusal-to-skip both keep working; a transport failure fails loudly
- [x] REQ-337's structural check converted to a behavioural one, with the change explained in its comment
- [x] No new module dependency

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modified, ~470 new lines)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)

**What was done:** `--remote-debugging-pipe` carries the DevTools Protocol over a pair of inherited file descriptors as NUL-terminated JSON — `os/exec`'s `ExtraFiles` supplies fd 3 and fd 4, `encoding/json` frames the messages, and nothing is added to `go.mod`. `Input.dispatchMouseEvent` delivers trusted events, so the engine issues the `pointerId` and a `setPointerCapture` on it establishes a real capture that really retargets the synthesized click. `TestTimelinePointerCaptureWaitsForThePanEngage` became `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` — the rename puts it inside the strict lane, which it was not in before — with four trials, each pair proving the other can fail.

## Testing

**Tests run:** the converted probe, `-count=2 -run '^TestBrowserBehavior'`, the strict lane and its zero-probe guard, and the full gate
**Result:** ✓ Gate exit 0, run twice. All 20 browser probes pass on both iterations of the doubled run.

**A correction to this REQ's own stated RED.** The REQ predicted `setPointerCapture` on a synthetic `PointerEvent` throws `NotFoundError`. On Chromium 141.0.7390.37 it does not throw — it silently establishes nothing:

```
{"capture":"established, pointerId 1","hasCapture":false,"isTrusted":false}
```

Same consequence, different mechanism, and the quieter one: the retargeting path is unreachable and nothing says so. That is why the check had to be structural.

**Mutation evidence, independently re-run by the orchestrator.** REQ-336's defect reintroduced *routed through a local variable* — the exact residual REQ-337's structural check recorded in its own doc comment as something it could not catch:

```js
var deferredCapture = capturePanPointer;
deferredCapture();
```

The old structural grep over the pointerdown body returns **0** — it passes. The converted probe **fails, exit 1**. That is the REQ's whole purpose demonstrated rather than argued.

**Two residuals the orchestrator found while re-running that mutation, neither of which the builder claimed otherwise:**
- The failure surfaced on trial 3 (the drag/pan trial) rather than trial 1 (the click/drawer trial the builder documented), and it took 46s — a 45s wait expiring, so the message names a timeout rather than the assertion the trial exists for. A failing run of this probe is slow, and its first line is not the most informative one.
- A **deferred** variant of the same mutation (`setTimeout(function(){ deferredCapture(); }, 0)`) **passed**. One run cannot distinguish "a deferred capture does not reproduce the defect" from "the probe misses that shape", and the orchestrator did not resolve it. Recorded as an open question for the reviewer rather than asserted either way.

**Lane guards:** `TestMaintainerStrictBrowserBehaviorLane` PASS; `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes` PASS; the engine-missing SKIP path still skips with its message.

*Verified by work action*

## Decisions

<!-- D-XX counter: last used D-07. Next decision: D-08. -->

- **D-01 — `--remote-debugging-pipe` over a WebSocket client or a driver. DECIDE.** Value: real trusted input with zero additions to `go.mod`, which the REQ required. Risk: ~200 lines of hand-rolled protocol handshake, one more thing to keep working across Chromium versions. Reversible — the dump-dom transport is untouched beside it.
- **D-02 — A second transport, not a replacement. DECIDE.** The 19 existing probes measure rather than interact; converting them would be a large diff with no property gained.
- **D-03 — Enforce the render-evidence rule inside `evaluateInPage` rather than per probe. DECIDE.** Value: the next probe author cannot forget it, and it caught a real bug during development — the session was attaching to `about:blank`, which reports `readyState: complete` and answers every question confidently about the wrong document.
- **D-04 — Convert REQ-333's release-outside half (trials 3+4) but leave REQ-333's own probe alone. DECIDE.** Risk: the same property is now checked twice, once structurally. Cheap to remove later.
- **D-05 — Leave REQ-324's threshold probe on synthetic events. DECIDE.** It measures the pan threshold, latching and press-point continuity, none of which needs a real capture. Risk: it stays blind to a capture-shaped regression — which the new probe now covers.
- **D-06 — Settle each gesture on its own condition. DECIDE.** The engine synthesizes the click in a task *after* the pointerup, so settling on the pointerup read a drawer that had not been asked to open yet. That flaked once under the full suite's load; neither condition presumes an outcome, so one helper serves the trial expecting a drawer and the trial expecting none.

## Discovered Tasks

- `TestBrowserBehaviorTimelinePointerAndKeyboardPathsStayAlive` still asserts structurally that `capturePanPointer` calls `scrollHost.setPointerCapture(`. Trials 3+4 now pin that behaviourally; the two structural assertions could be dropped, or that probe's drag half rebuilt on the trusted transport so its `lostpointercapture` is real.
- `_dev/primes/prime-kanban-board.md` has no Lessons entry for this REQ, and the which-driver reasoning now lives only in `browser_probe_test.go`. Two lessons worth recording: a protocol target reports its new URL before the renderer has swapped documents, so attaching by "the first page target" measures `about:blank` confidently; and the engine synthesizes the click *after* the pointerup.
- Nothing pins that a lane probe cannot silently skip the trusted transport specifically — the probe counter is shared between both transports. A per-transport count would catch a future edit that quietly moves every probe back to `--dump-dom`.
- Five probes hardcode `"/probe.html"` in their href assertions instead of reading `browserProbePageFileName`. Same class as REQ-322's lesson about restating a constant beside the test.

## Open Questions

- [x] Does a **deferred** capture on pointerdown reproduce REQ-336's defect, and if so does the
  converted probe catch it? → **Answered by measurement in the remediation (`3148b1d`): yes to both.
  The probe had a gap; it was not a boundary of the defect.** The orchestrator's `setTimeout(…, 0)`
  variant passed because the press and release were dispatched 1ms apart and the timer callback had
  not run — `captureAt` was absent entirely. Given two frames of dwell, the capture is granted 16ms
  after the press, the pointerup retargets to the scroll host, and the drawer stays shut. **Any human
  click dwells far longer than 16ms, so that board was broken for every real user while the probe
  called it fine.** Closed by giving click gestures a real dwell (`clickTrustedMouseOnRow`). All four
  capture spellings — immediate, variable-routed, `setTimeout 0`, `requestAnimationFrame` — are now
  caught; the old structural grep saw three of them as zero hits. The reviewer swept further and
  found the sensitivity cut is exactly at two frames: `raf3` (capture at +47.9ms, after the release)
  is not caught. That boundary is stated in the test's comment. It is worth knowing the probe is less
  sensitive than a person — a human press is ~100ms and would catch `raf3`.

## Review

**Overall: 96%** | 2026-08-24T12:37:51Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** one, and it was a stale sentence in this REQ's own trail rather than anything in the code — D-07 recorded as unresolved after the remediation had answered it. Closed above.

**The reviewer verified the remediation's diagnosis independently rather than accepting it.** Reproduced outside the suite with its own CDP driver on a minimal page: press on `#a`, replace the host's children, release → `clicks: []`, `pressTargetSurvived: false`, 3 of 3; control gives `clicks: ["a"]`. The `mouseup` target reports the *new* node with the same id, so "survived" is the only usable discriminator — which is what the probe watches.

**The race is closed, not narrowed, and there was a third rebuild source.** `board-timeline.js:2424-2429` observes `scrollHost` with a `ResizeObserver` calling `requestFrameRender` → `renderAll`; the detail drawer taking or returning its grid column fires it, which is why trial 2 was the one that hung, and the first version had no settle after `closeDrawer()`/`resetWindow()` at all. Instrumented with a `MutationObserver` on the rows group: **0 rebuilds between press and release in 12 of 12 trials**. The other reachable sources (window `resize`, pointermove-while-engaged, `markTimelineTableStale`) are not triggered by these gestures. Under load: `-count=3` of the whole lane under 8 spinners, exit 0.

**The dwell adds exposure but no new failure mode.** Press-to-release measured 32.0-32.8ms against 1-2ms back-to-back, so the window a rebuild could land in is ~16x wider — but nothing the press does schedules a rebuild once the 8px margin removes the focus-scroll, and a future change that did would fail in ~6s naming the mechanism rather than at 45s naming the wait.

**The stated boundary is real and accurately stated.** Deferral sweep on the live board: `immediate` caught (capture at +0.0ms), `setTimeout 0` caught (+14.6), `raf1` caught (+14.4), `raf2` caught (+30.8), **`raf3` not caught** (capture refused at +47.9, after the release), `raf4` and `setTimeout 50` not caught. The cut is exactly at two frames, as the comment says. Worth knowing the probe is *less* sensitive than a person: a human press is ~100ms and would catch `raf3`.

**`evaluateInPage`'s promise resolution cannot mask a navigation** — the href is read inside the `.then`, at or after value production, and a cross-document navigation destroys the context so `Runtime.evaluate` fails loudly. Strictly stronger than a start-of-call check. **No caller ignores a false verdict** from `pageConditionHoldsWithin`: two callers, both `t.Fatalf`.

**Minor findings:** 4 (report only) — `waitForPageCondition` is used against its own contract at `:2951`, waiting on `window.timelineProbe.isPanning()`, a board property rather than an engine guarantee, so a regression that stops the pan engaging still burns the full 45s; the `NO WALL-CLOCK WAITS` header is now a half-truth since the new transport polls at 25ms; six `"/probe.html"` literals remain beside the constant; a failing *drag* trial is still slow, since only click trials got the 5s budget.

**Acceptance:** Pass — gate exit 0; converted probe 6/6 green in isolation at ~1.8s each; full browser lane `-count=3` under 8 CPU spinners exit 0 at 196s.

**Follow-ups created:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Falsifying the load hypothesis before acting on it. Six CPU spinners left the probe 12-for-12 green, and the failing capture showed the probe polling `Runtime.evaluate` every 25ms and getting answers throughout the 45s wait, with both mouse events delivered. That killed "starvation" and pointed at the one thing left: the click was never created.

**What didn't:** The first version settled nothing after `closeDrawer()`, so the drawer's own grid relayout was still in flight when the next trial took aim — which is why trial 2 was always the one that hung. Raising the deadline would have hidden it for months. A 45-second wait that expires is not a slow test, it is a test that has stopped knowing what it is waiting for.

**Worth knowing:** Chromium synthesizes a click on the nearest common inclusive ancestor of the mousedown and mouseup targets, so a mousedown target detached in between produces **no click at all** — not a click on the wrong element. Any probe that presses, lets the page work, and releases can hit this. And the probe is less sensitive than a human hand: it dwells two frames where a person dwells ~100ms, so a capture requested three frames after the press passes here and breaks for every real user.

## Orientation

The board's browser probe lane can drive real, trusted input — a press, a drag, a release the engine issues a `pointerId` for — over the DevTools Protocol on a pair of inherited file descriptors, with no new module dependency. Lives beside the existing `--dump-dom` transport, which is untouched. `[MAP CHANGED]` — this adds a second way for the test lane to talk to a rendered board, and probes that previously had to assert structure can now assert behaviour.
