---
id: REQ-336
title: "[impact-critical] Timeline clicks open the detail drawer again"
status: completed
claimed_at: 2026-08-23T19:41:18Z
completed_at: 2026-08-23T20:09:24Z
commit: 4527a50
kb_status: promoted
kb_entry: REQ-336-timeline-clicks-open-the-detail-drawer-a.md
created_at: 2026-08-23T18:30:26Z
user_request: UR-067
addendum_to: REQ-333
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-23T19:41:18Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - browser evidence
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
route: B
related: [REQ-337, REQ-338]
batch: timeline-click-regression
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Timeline Clicks Open the Detail Drawer Again

## What

Mouse clicks on Timeline (Gantt) view REQ bars no longer open the detail drawer — every mouse
click inside `#timeline-scroll` is swallowed. Restore clicking while keeping the pan behaviour
REQ-333 and REQ-324 established.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md` (and its REQ-324, REQ-323, REQ-322,
  REQ-321 lesson links), `general.md`, `coding-guardrails.md`, `communication-style.md`.
  Fix B, at the one seam the pan state machine offers:
  1. Reproduce RED with trusted input over CDP — the repo's probe lane cannot dispatch it →
     verify: the click targets `#timeline-scroll` and `closest("[data-detail-kind]")` is NONE.
  2. Delete the `setPointerCapture` block from the `pointerdown` handler; call it from the
     `pointermove` handler on the false→true engage transition only, via a named
     `capturePanPointer` carrying the relocated rationale → verify: all four GREEN conditions with
     real input.
  3. Drop `pointerdown`'s `panState.pointerId = undefined` clobber on capture failure — the
     teardown's `hasPointerCapture` gate already answers the question → verify: condition 4 still
     holds.
  4. Probe the risk the move introduces (a press that leaves the host sub-threshold leaving
     `panState` armed, so an unbuttoned re-entry could pan) before deciding whether it needs a
     guard → verify: measured, not assumed.
- [x] **[APPLY]:** `web/board-timeline.js` only, plus the mid-build scope extension in D-02.
  Capture moved into `capturePanPointer()`, called from the engage transition; the `pointerId`
  clobber dropped; `lostpointercapture` kept in the teardown list.
- [x] **[UNIFY]:** `git diff --stat` → 2 files, both in the (extended) write set.
  - `web/board-timeline.js` — read the whole hunk: one block moved, one function added, one
    two-line transition guard, one clobber removed. `node --check` clean.
  - `timeline_browser_probe_test.go` — read the whole hunk: one assertion retargeted into two
    halves plus its comment. `gofmt` clean (the gate's 59-file formatting lane passed).
  - No debug artifacts: no `console.log`, `debugger`, `TODO` or `alert` in the diff. The CDP probe
    lives in the scratchpad, not the repo.
  - Native lint/tests: `go vet ./...` clean via the gate; `go test -count=1 ./...` with the browser
    lane enabled → ok.

## Why

A primary interaction of a shipped view is broken: no mouse path from the Timeline to a REQ/UR's
details. It also breaks REQ-324's contract that a press which does not move is still a click.
User-judged `impact-critical` (broken production path).

## Context

**Cause (user-diagnosed):** the `setPointerCapture(downEvent.pointerId)` call in the pointerdown
handler of `web/board-timeline.js` (currently at line ~2580). Pointer capture retargets subsequent
pointer events AND the synthesized click to the capturing element (`#timeline-scroll`), so
`clickEvent.target.closest("[data-detail-kind]")` in the delegated click handler at
`web/board-controls.js:183` finds nothing and the drawer never opens. Introduced by REQ-333,
commit 36c4518, version 0.236.28.

**Blast radius (already audited — do not re-widen):** mouse clicks inside `#timeline-scroll`
only. Enter and Space on a focused row still open the drawer. A mouse click on a Board-view card
still opens it. Hover readout, ctrl+wheel zoom over plot and axis, and view switching are all
unaffected.

## Detailed Requirements

Reproduce RED first, with real browser input — not synthetic PointerEvents (the user's diagnosis
came from instrumented real input; synthetic PointerEvents with an engine-unknown `pointerId`
make `setPointerCapture` throw and never establish capture, so they cannot reproduce this).

Two candidate fixes, both already validated per-gesture on fresh pages by the user:

- **B** — call `setPointerCapture` only once the pan engages past `TIMELINE_PAN_THRESHOLD_PX`,
  not on pointerdown
- **C** — drop capture entirely, bind the pointerup/pointercancel release to window

**Prefer B**: it keeps REQ-333's release guarantee and makes "a drag is not a click" an explicit
property of the capture instead of an accident of node rebuilding. Take C only if B turns out to
leave a release path open.

## Builder Guidance

Certainty: Firm — cause, blast radius, and preferred fix are user-diagnosed and pre-validated.
`tdd: false` because the existing probe lane cannot drive the real input this RED requires
(REQ-333's UNIFY notes: a synthetic `pointerId` cannot be captured in this lane); the durable
automated check is REQ-337's deliverable — do not duplicate it here. Prove RED/GREEN with real
browser input against a generated board; what counts as real input in the probe lane is
REQ-337's call to make.

## Red-Green Proof

**RED prompt/case:** Generate a board, open the Timeline view, click a REQ bar with a real mouse
click. No drawer opens. Instrumented: pointerdown targets the `<rect>`, but the click event
targets `DIV#timeline-scroll`, so `closest("[data-detail-kind]")` returns NONE and the delegated
handler at `web/board-controls.js:183` never fires. A 2px sub-threshold nudge is swallowed too.
**Why RED now:** `setPointerCapture` on pointerdown retargets the synthesized click to
`#timeline-scroll` for every press, drag or not.
**GREEN when:** All four hold, with real browser input:
1. plain click on a bar → drawer opens with that REQ/UR
2. 2px nudge then release → drawer opens (REQ-324's contract)
3. 150px drag → pans, drawer does NOT open
4. drag released outside the chart → pan ends, no stuck pan state
**Validation:** User confirmed (RED, GREEN, cause, and candidates all stated verbatim in the input)

## Prior Implementation

REQ-333 ("Keep the timeline pointer and keyboard paths alive", archived, `completed-with-issues`,
commit 36c4518) added the `setPointerCapture` call to make a drag released outside the chart a
guaranteed release rather than a hope, fixing the latched `.is-panning` state. It wrote
`web/board-timeline.js` and `timeline_browser_probe_test.go`; the capture call is feature-guarded
and degrades to the ordinary release path when unavailable. That release guarantee is GREEN
condition 4 — keep it.

## Dependencies

None to start. REQ-337 (the lock-in check) depends on this REQ; REQ-338 also writes
`web/board-timeline.js`, declared in both write sets.

---
*Source: UR-067 — see `do-work/user-requests/UR-067/input.md` for complete verbatim input.*

---

## Triage

**Route: B** - Medium

**Reasoning:** Cause, fix and blast radius are user-diagnosed, so nothing needs planning. But the change lands inside a ~2,600-line frontend module whose pan state machine REQ-324 and REQ-333 both touched, and the exact seam for "capture only once the pan engages" has to be found in that state machine — that is exploration, not planning.

**Planning:** Not required

## Exploration

**The seam.** `web/board-timeline.js:2554-2630` holds the whole pan state machine: a
`pointerdown` handler that arms `panState` (`engaged: false`) and immediately calls
`scrollHost.setPointerCapture` (`:2578-2586`), a `pointermove` handler that engages the pan once
`timelinePanEngaged` clears `TIMELINE_PAN_THRESHOLD_PX` (`:2588-2609`), and one teardown bound to
`pointerup`/`pointercancel`/`pointerleave`/`lostpointercapture` (`:2610-2630`). Fix B's insertion
point is the false→true engage transition in the move handler — the same instant
`is-panning` goes on.

Two details the fix has to respect:

- **The teardown's release is already gated on `hasPointerCapture`** (`:2622-2623`), so it is
  correct whether or not capture was ever taken. The extra `panState.pointerId !== undefined`
  clause is only there to avoid passing `undefined`; nothing needs a separate "did we capture"
  flag.
- **`pointerdown` currently clobbers `panState.pointerId` to `undefined`** when
  `setPointerCapture` throws (`:2584`), overloading the field as both the id and a
  did-we-capture marker. Moving the capture out of `pointerdown` retires that overload.

**`lostpointercapture` is still needed in the teardown list** even though capture now starts
later: it is what fires when the engine takes a capture away mid-drag, which is only possible once
capture exists.

**Harness.** The repo's browser probe lane (`browser_probe_test.go:99-167`) runs a page under
`--dump-dom` and reads a JSON result node. It has no CDP channel, so it cannot dispatch trusted
input — which is exactly why this REQ carries `tdd: false` and why REQ-337 exists. For this REQ's
RED/GREEN I drove a headless Chromium over the DevTools Protocol directly
(`Input.dispatchMouseEvent`, so the events are trusted and the `pointerId` is engine-assigned)
against a real generated board, with `location.href` returned alongside every measurement per the
prime's render-evidence rule.

**RED, reproduced with real input** (Chromium 1194 / Playwright build, `--headless=new`, viewport
1400×900, first `[data-detail-kind]` in `#timeline-scroll` scrolled to centre — REQ-340's bar):

```
click         drawerOpen=False shownId='' panningSeen=False isPanningAfter=False clickTarget=DIV#timeline-scroll closest=NONE
nudge         drawerOpen=False shownId='' panningSeen=False isPanningAfter=False clickTarget=DIV#timeline-scroll closest=NONE
drag          drawerOpen=False shownId='' panningSeen=True  isPanningAfter=False clickTarget=None            closest=None
drag-outside  drawerOpen=False shownId='' panningSeen=True  isPanningAfter=False clickTarget=None            closest=None
```

`pointerdown` targets the `<rect>` while the click targets `DIV#timeline-scroll` and
`closest("[data-detail-kind]")` returns NONE — the user's diagnosis, confirmed. Conditions 3 and 4
already hold and are the ones the fix must not break.

**Two measurement traps found while building the probe, both worth writing down:**

- A bar's own coordinates are not clickable until the bar is on screen. The first
  `[data-detail-kind]` sat at y≈1538 in a 900px viewport, so the press landed on `HTML` and the
  probe reported a closed drawer for a click that never touched the chart. `scrollIntoView` plus a
  viewport assertion before dispatching is the fix.
- **"It panned" cannot be measured from the axis text.** A 150px rightward drag engages the pan and
  leaves the axis labels identical, because the window clamps at the bound — REQ-324's lesson,
  reproduced. The engage itself (`is-panning` observed during a `pointermove`) is the observable.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — move the pointer
  capture from `pointerdown` to the pan's engage transition.
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify) — **added
  mid-build, see D-02**: REQ-333's structural lock-in asserts the capture call sits in the
  `pointerdown` handler, so moving the call breaks it. Retarget that one assertion to the new
  instant; no new check (REQ-337 owns that).

**Files I will NOT touch:** `web/board-controls.js` (its delegated click handler is correct and is
the consumer this fix serves), `web/board-detail.js`.

**Acceptance criteria (restated from REQ):**
- [ ] Plain click on a bar → drawer opens with that REQ/UR
- [ ] 2px nudge then release → drawer opens (REQ-324's contract)
- [ ] 150px drag → pans, drawer does NOT open
- [ ] Drag released outside the chart → pan ends, no stuck pan state (REQ-333's guarantee)

## Decisions

- **D-01 — DECIDE & STATE. Fix B, as the REQ preferred, and the `pointerId` clobber went with it.**
  `pointerdown` used to set `panState.pointerId = undefined` when `setPointerCapture` threw,
  overloading the field as both the id and a did-we-capture marker. With capture taken at the
  engage that overload has no purpose: the teardown's release is gated on
  `hasPointerCapture(panState.pointerId)`, which answers false when nothing was ever captured, so
  keeping the real id is both simpler and more correct than blanking it. No separate `captured`
  flag was added for the same reason.
- **D-02 — ESCALATE (scope). `timeline_browser_probe_test.go` was added to the write set
  mid-build.** REQ-333 left a structural lock-in asserting `scrollHost.setPointerCapture(` appears
  inside the `pointerdown` handler's body (`timeline_browser_probe_test.go:2316`); moving the call
  is exactly what that assertion forbids, so it failed. The behaviour change is intentional and
  user-pre-validated, so the Cross-REQ Test-Break rule says update the test rather than the
  implementation. The assertion was retargeted, not deleted: it now checks the capture call inside
  `capturePanPointer` **and** that the move handler calls it, so deleting either half still fails.
  **Value:** REQ-333's release guarantee stays pinned at the instant it now happens, and the
  Timeline probe lane stays green so REQ-337 starts from a passing lane. **Risk:** the file is also
  in REQ-337's and REQ-338's write sets; a parallel run would collide. Serial here, so it does not,
  and both later REQs are told about it. Reversible — the assertion is eleven lines.
  No new check was added: REQ-337 owns the durable click-retargeting check and this REQ's Builder
  Guidance says not to duplicate it.
- **D-03 — DECIDE & STATE. No `buttons` guard added to the move handler.** Moving the capture
  reopens a window REQ-333 had closed: a press that leaves the host *before* the 4px threshold
  trips takes no capture, so the release can land elsewhere and leave `panState` armed — and the
  move handler does not check that a button is still held, so an unbuttoned re-entry could pan.
  Probed rather than assumed: `pointerleave` **does** reach the host on a buttoned exit in this
  engine, four times over, so the teardown runs and `panState` is cleared; the
  `hover-after-subthreshold-exit` gesture shows no pan (`panningSeen=False`). With no reproduction
  there is nothing to earn the guard (`coding-guardrails.md` § Simplicity First, *Earned defense*).
  Recorded rather than silent because the reasoning is what a future reader would otherwise redo.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)

**What was done:** Moved the Timeline pan's `setPointerCapture` out of the `pointerdown` handler
into a named `capturePanPointer` that the `pointermove` handler calls once, on the false→true pan
engage. A press that never travels 4px now takes no capture, so its synthesized click keeps its
`<rect>` target and the delegated `[data-detail-kind]` handler opens the drawer; a drag still
captures the moment it engages, so its release is still guaranteed wherever it lands. Retargeted
REQ-333's structural capture assertion to the new call site, in two halves so neither can be
deleted silently.

## Qualification

Passed — 2 files verified, 4 requirements traced, P-A-U confirmed.

- Mechanical: `tools/checks/qualify.sh` exit 0; `tools/checks/scope-drift.sh` exit 0 against the
  Scope list as extended by D-02.
- Substantive (check 2): +50/-21 across two files; read both hunks. The JS hunk moves a block,
  adds one named function and one two-line transition guard, and removes one clobber — no
  whitespace-only or import-shuffle content. The Go hunk replaces one assertion with two.
- Requirements traced (check 3): all four GREEN conditions measured with trusted input (table in
  `## Testing`); the fix's only behavioural lever is *when* capture is taken, which is the single
  cause the REQ names.
- Data flows (check 6): nothing stubbed. `capturePanPointer` reads the live `panState.pointerId`
  and calls the real DOM API behind its feature detect; the retargeted assertion reads the
  generated page rather than restating the source.
- Contamination check: the previous REQ in this session (REQ-325) touched
  `skills/do-work-toolbox/scripts/generate-report-image.sh` and its prescribed-shell case file.
  No overlap with this REQ's two files.

## Testing

**Tests run:**
- `go test -count=1 ./...` in `skills/do-work-board/tools/queue-kanban` with
  `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium-1194/chrome-linux/chrome`
- `bash _dev/tests/maintainer-verify.sh` (the project's declared canonical repository gate), from
  the repo root, against the final tree, with the browser lane enabled rather than skipped
- A CDP real-input probe against a freshly generated board, six gestures

**Result:** ✓ board suite `ok … 115.898s`; canonical gate exit 0 with the strict browser lane
running (not SKIPped); all four GREEN conditions hold.

**Red-green validation:** trusted input over `Input.dispatchMouseEvent`, Chromium 1194
(Playwright build), `--headless=new`, viewport 1400×900, against a board generated from this
repo's own queue. Target: the first `[data-detail-kind]` inside `#timeline-scroll` (REQ-340's bar),
scrolled to centre and asserted inside the viewport before dispatch. Drawer-open is read as
`#detail-drawer`'s `hidden === false` plus `#detail-id`'s text — the state `showDrawer` /
`closeDrawer` actually set — not as "some class looks open".

| gesture | before | after |
|---|---|---|
| click on a bar | ✗ drawer closed, click target `DIV#timeline-scroll`, `closest` NONE | ✓ drawer open showing `REQ-340`, click target `rect` |
| 2px nudge then release | ✗ drawer closed, same retargeting | ✓ drawer open showing `REQ-340` |
| 150px drag | ✓ pan engaged, drawer stays closed | ✓ unchanged |
| drag released outside the chart, pan engaged first | ✓ pan engaged, no `is-panning` left on | ✓ unchanged |

The last row is the gesture that actually tests REQ-333's guarantee; the first shape I wrote for it
left the host before the 4px threshold, so it never engaged a pan and proved nothing about the
release. Both shapes are kept in the probe and both pass.

**Existing tests updated (cross-REQ impact):**
- `timeline_browser_probe_test.go` — `TestBrowserBehaviorTimelinePointerAndKeyboardPathsStayAlive`
  (from REQ-333): its structural assertion that `scrollHost.setPointerCapture(` appears in the
  `pointerdown` handler is the thing this REQ deliberately changes. Retargeted to
  `capturePanPointer`'s body plus the move handler's call to it — two assertions, so deleting
  either half fails. REQ-333's contract (the capture is really requested, so the release is a fact)
  is unchanged; only its instant moved. Full reasoning in D-02.

**No new tests added here on purpose:** REQ-337 is the durable click-retargeting check and this
REQ's Builder Guidance forbids duplicating it. `tdd: false` for the reason the REQ states — the
probe lane runs under `--dump-dom` with no CDP channel, so it cannot dispatch the trusted input
this RED needs.

*Verified by work action*

## Review

**Overall: 97%** | 2026-08-23T20:06:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 1 (report only) — no new durable check lands in this REQ, so between this
commit and REQ-337's the click path is held only by a structural assertion about where the capture
call sits. That is deliberate (REQ-337 owns the behavioural check and this REQ's Builder Guidance
forbids duplicating it) and REQ-337 is the very next dependent in the queue.

**Acceptance:** Pass — all four GREEN conditions measured with trusted input; board suite ok with
the browser lane on; canonical gate exit 0.
**Suggested testing:** 3 items
**Follow-ups created:** None; **sweeps appended to:** None

### Requirements Checklist

- [x] Plain click on a bar → drawer opens with that REQ/UR — delivered (`shownId='REQ-340'`)
- [x] 2px nudge then release → drawer opens (REQ-324's contract) — delivered
- [x] 150px drag → pans, drawer does NOT open — delivered (`panningSeen=True`, drawer closed)
- [x] Drag released outside the chart → pan ends, no stuck pan state — delivered
      (`engage-then-outside`: pan engaged, `isPanningAfter=False`)
- [x] Fix B preferred over C — delivered; C was not needed because B leaves no release path open
      (D-03's probe)
- [x] RED reproduced with real browser input, not synthetic PointerEvents — delivered (CDP
      `Input.dispatchMouseEvent`)
- [x] Blast radius not re-widened — delivered; the diff changes one call's timing plus the comment
      that described the old single cause

### Restatement Sweep

The diff redefines *when* the Timeline takes pointer capture, and the prose around
`TIMELINE_PAN_THRESHOLD_PX` restated the consequence. Swept:

- **`web/board-timeline.js:124-129` was a half-truth and is fixed in this diff.** It said the
  first-pointermove re-render "is the whole defect" behind a lost click. REQ-336 proves a second,
  independent cause for the same outcome, and the threshold cannot protect against it — exactly the
  shape the board prime's REQ-231 lesson names ("a change that adds a second cause for an outcome
  makes every sentence naming the first one a half-truth"). Rewritten to name both causes and point
  at `capturePanPointer`. In the declared file, so it is this REQ's fix rather than a follow-up.
- `timeline_browser_probe_test.go:2192` explains that a synthetic `pointerId` cannot be captured in
  this lane — still true, and now the reason this REQ's RED needed CDP. Left as written.
- `web/board-detail.js:501,548` captures the pointer on the resizer's `pointerdown`. Not a
  restatement of this contract and not stale: that divider has no click consumer to retarget, so
  capturing on the press is correct there. Verified rather than assumed.
- No shipped prose (`skills/do-work-board/actions/board.md`, `docs/`) states when the Timeline
  captures; nothing to update.

### Code Review Notes

- **Naming for reach:** one new identifier, `capturePanPointer` — three words, greppable.
  `wasEngaged` is a short-lived local one line from its use.
- **Simplicity:** the fix is a move plus a two-line transition guard. No flag, no state field, no
  new abstraction; D-01 and D-03 each record something deliberately *not* added.
- **Surgical:** every changed line traces to the REQ. The pan maths, teardown list, threshold value
  and keyboard path are untouched.
- **Correctness of the transition guard:** `wasEngaged` is read before `timelinePanEngaged` updates
  the flag, so `capturePanPointer` runs exactly once per drag rather than on every move — and
  because `timelinePanEngaged` is sticky (`alreadyEngaged || …`), the flag cannot flap back and
  cause a second capture.
- **Failure path:** a `setPointerCapture` that throws now leaves `panState.pointerId` intact instead
  of blanking it. The teardown's `hasPointerCapture` gate returns false in that case, so the release
  is skipped exactly as before — verified by reading both branches, not inferred.

### Acceptance Testing

**Result: Pass**
- Six gestures with trusted input against a freshly generated board; the four GREEN conditions plus
  two shapes written to probe the fix's own risk (`drag-outside` sub-threshold and
  `hover-after-subthreshold-exit`), all as expected.
- Board Go suite with the strict browser lane enabled: ok.
- Canonical gate: exit 0, with the browser lane running rather than SKIPped.
- Regression check on the adjacent path the REQ names as unaffected: the keyboard half of
  `TestBrowserBehaviorTimelinePointerAndKeyboardPathsStayAlive` (three arrow presses, focus stays in
  the chart) passes, so Enter/Space and arrow panning are intact.

### Suggested Additional Testing

- Manual verification with a real hand on a trackpad: a tremor under 4px should still open the
  drawer, and a two-finger drag should still pan. The probe dispatches ideal geometry.
- Touch and pen input. Every gesture measured here was `pointerType: "mouse"`; capture semantics are
  the same in the spec but the threshold's feel is not.
- The `drag-outside` sub-threshold shape on a build where `pointerleave` really is suppressed while
  a button is held (REQ-333's premise, which did not reproduce in Chromium 1194 — see D-03). If such
  a build exists, that is where a stale armed `panState` would show up.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Building the CDP probe before touching the code. It cost more than reading the diff would have,
  and it is the only reason three separate claims in this area are now measured rather than
  believed: that the click is retargeted, that a 150px drag clamps without moving the window, and
  that `pointerleave` reaches the host on a buttoned exit.
- Writing the gesture that probes the *fix's own* new risk (`hover-after-subthreshold-exit`)
  before deciding whether to guard against it. It said no guard was needed, which is a cheaper
  answer than the guard.

**What didn't:**
- The first `drag-outside` gesture. It moved diagonally and left the host before the 4px threshold
  tripped, so it never engaged a pan — and it therefore proved nothing about REQ-333's
  release guarantee while looking like it did. A gesture aimed at a release path has to engage the
  thing being released first.
- Measuring "it panned" from the axis label text. A 150px drag pans and leaves the labels
  identical, because the window clamps at the bound. This is REQ-324's lesson arriving a second
  time, in a different disguise, in the same view.
- The first click attempt measured a closed drawer for a press that never reached the chart: the
  bar sat at y≈1538 in a 900px viewport. An element's `getBoundingClientRect` is not a clickable
  coordinate until the element is on screen.

**Worth knowing:**
- **Pointer capture retargets the synthesized click, not just the pointer events.** Any handler
  that captures on `pointerdown` inside a container with delegated click handling breaks every
  click in that container, and nothing about the pan logic makes that visible. `board-detail.js`'s
  resizer captures on `pointerdown` and is fine only because nothing delegates clicks through it.
- **This view has now eaten a click for two independent reasons** (a re-render that rebuilt the
  node, and capture retargeting). The movement threshold protects against the first only. A third
  cause would look identical from the outside, so a behavioural check on the click path is worth
  more here than any amount of reasoning about the pan state machine — which is REQ-337.
- REQ-333's stated premise, that Chromium suppresses boundary events while a button is held, did
  not reproduce in Chromium 1194: `pointerleave` reached the host four times on a buttoned exit.
  The capture is still worth taking at the engage, but a future reader should not treat that
  premise as established.
- The repo's browser probe lane runs under `--dump-dom` with no CDP channel, so it cannot dispatch
  trusted input at all. Anything needing real pointer semantics has to drive
  `Input.dispatchMouseEvent` itself; Node 22's global `WebSocket` is enough to speak the protocol
  with no dependencies.

## Orientation

Clicking a Timeline bar opens the detail drawer again, and dragging still pans without opening it.
Lives in the board's Timeline view (`skills/do-work-board/tools/queue-kanban`), in the pan state
machine only — the keyboard path, the drawer, and the delegated click handler are untouched. No new
module and no contract change, so the system's shape is unchanged; `_dev/primes/prime-kanban-board.md`
gains one lesson link and its referenced paths all still resolve.
