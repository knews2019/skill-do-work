---
id: REQ-425
title: 'Stop the Timeline''s trailing-window controls assuming now and a full screenful are inside the bounds'
status: claimed
created_at: 2026-08-30T11:40:00Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-30T12:23:07Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
claimed_at: 2026-08-30T12:23:07Z
user_request: UR-079
domain: frontend
prime_files: ['_dev/primes/prime-kanban-board.md']
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
addendum_to: REQ-390
review_generated: true
sweep: true
sweep_key: trailing-window-bound-assumptions
write_set: [skills/do-work-board/tools/queue-kanban/web/board-timeline.js, skills/do-work-board/tools/queue-kanban/generate_test.go]
---

# Stop the Timeline's Trailing-Window Controls Assuming Now and a Full Screenful Are Inside the Bounds

## What
REQ-390's trailing-window controls both assume the board's `now` and a whole screenful of window fit inside the timeline bounds. Neither holds at the edges, and each assumption produces a distinct user-visible defect. One root cause, so one REQ.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] **Chips collapse to the one-hour floor on a drained board.** `timelineTrailingWindow` computes `[now - N days, now]` and clamps both endpoints into the bounds. `timelineRange` (`timeline.go`) only stretches `rangeEnd` to `now` while some row has `WaitOpen` or `WorkOpen`, so once a queue is fully drained the bounds end at the last completion instant and `now` falls outside them. Every chip whose span starts after the padded bound end clamps to zero width and settles to `TIMELINE_MIN_SPAN_MS`.
- [ ] **The arrows stopped being inverses near the right bound.** `timelinePannedWindow` clamps a forward step at the bound so it moves only a partial screenful, while the following back step moves a full one. Press `›` then `‹` and the reader lands to the left of where they started. The deleted calendar-period test had pinned this property; REQ-390's replacement covers only the mid-range case.

## Detailed Requirements
- A trailing-window chip must produce a distinct, non-degenerate window on a board whose `now` sits outside the bounds, and must never silently share a window with another chip.
- The lit chip and the state readout must not claim a window the reader did not ask for.
- `‹` and `›` must be inverses wherever both are enabled, including at and near both bounds. Whichever way this resolves — refusing a step that cannot move a full screenful, or making the back step mirror the clamped forward step — the pair must round-trip.
- Both behaviours need a test that fails before the fix. The arrow property had one and lost it; do not leave it uncovered a second time.

## Constraints
- Do not reintroduce calendar-period arithmetic. REQ-390 deleted roughly 200 lines of it deliberately.
- Keep the control set declared in `template.html` alone; the lit chip and state readout stay derived from the DOM rather than from a list in JS.
- Preserve REQ-390's clamp-before-settle discipline and the `## Scope` boundary of two files unless the fix genuinely requires more, in which case flag it rather than widening silently.

## Builder Guidance
Certainty level: Firm on the defects, open on the remedy. The reviewer proposed anchoring the trailing span on `Math.min(Math.max(nowMs, boundStartMs), boundEndMs)` and deriving the clipped flag from that anchor, and refusing an arrow step that cannot move a full screenful. Both are starting points, not decisions — the second trades reachability of the last partial screenful for a round-tripping pair, which is a judgment call worth stating explicitly.

## Red-Green Proof
**RED prompt/case:** In the Node lane, drive `timelineTrailingWindow` on a board whose bounds end days before `now` (no open rows) and assert the five chips produce five distinct windows, none of them at the one-hour floor. Separately, drive `timelinePannedWindow` a partial screenful from the right bound and assert `+1` then `-1` returns to the original window.
**Why RED now:** Measured on the merged tree at `59105df`. Chip collapse: 3 days idle collapses "Last day"; 10 days puts "Last day" and "Last 7 days" on the same one-hour window; 40 days collapses three chips; 100 days collapses four of five. Arrow drift: -120.00h a partial screenful from the right bound, -168.00h flush against it, 0.00h mid-range.
**GREEN when:** Both cases pass, no chip degenerates to the zoom floor while the board has a range to show, and the arrow pair round-trips at both bounds.
**Validation:** Confirmed by the orchestrator's own measurement during REQ-390's review; three further Important findings from that review were refuted 3-0 and are deliberately not carried here.

---
*Source: REQ-390 review (UR-079). Sweep REQ: one root cause, two instances.*

---

## Triage

**Route: B** - Medium

**Reasoning:** Both defects are precisely located and independently measured, so nothing needs planning from scratch — but the arrow remedy is a real design call (refuse a step that cannot move a full screenful, versus mirror the clamped forward step), and the surrounding Timeline conventions need discovery before choosing. Exploration first, then build.

**Planning:** Not required

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js`
- `skills/do-work-board/tools/queue-kanban/web/template.html`

**Files I will NOT touch:**

- VERSION, skills/do-work/VERSION, skills/do-work/actions/version.md, CHANGELOG.md, skills/do-work/CHANGELOG.md — serial-only integrator files.
- Anything under do-work/ — queue state is the orchestrator's.

**Acceptance criteria (restated from REQ):**

- [ ] A chip produces a distinct, non-degenerate window on a board whose now sits outside the bounds, never silently sharing a window with another chip.
- [ ] The lit chip and the state readout never claim a window the reader did not ask for.
- [ ] `‹` and `›` are inverses wherever both are enabled, including at and near both bounds.
- [ ] Both behaviours have a test that fails before the fix.

## Implementation Summary

**Files changed (3), taken from the merge range:**

- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)

**What was done:** `timelineTrailingWindow` now hangs each chip's span off `now` clamped into the bounds, so a drained board reads each chip as the last N days of the *recorded range* rather than clamping to zero width and settling on the one-hour floor. It also returns `isClippedByBounds`, which the state readout consumes instead of recomputing the clamp — deleting a branch the new anchor makes unreachable. A new `timelineSteppedScreenfulWindow` makes the toolbar arrows all-or-nothing, so a step either moves a whole screenful and round-trips exactly or refuses; `timelinePannedWindow` is deliberately untouched, so the keyboard pan and drag keep their continuous clamp. The integrator folded in the builder's declared S1 seam: the timeline hint and toolbar comment in `template.html` both claimed the chips end the window at now, which stops being true once the bounds stop reaching now.

**Builder branch:** `worktree-agent-REQ-425-trailing-window-bound-assumptions` — `3e9f7f8` (RED only), `cbe9de0`, `52f8955`.
**Merge range:** `b2cbe87..04b8120`.
**Hand-back:** `do-work/runs/work-2026-08-29-213539/REQ-425-handback.md`.

## Testing

**Tests run (merged tree, merge range `b2cbe87..04b8120`):**
- `QUEUE_KANBAN_BROWSER=<Chrome 151> go test -count=1 -run '^TestBrowserBehaviorTimeline|^TestTimeline|^TestJavaScriptBehaviorTimeline' .` — ✓ exit 0, 119s
- `bash _dev/tests/maintainer-verify.sh` — ✓ exit 0, browser lane in its default skipped state

**Red-green validation:** RED was committed as its own commit (`3e9f7f8`) before any source change, reproducing the review's measurements exactly.

- Chip collapse: ✗ before — 3 days idle collapsed "Last day"; 10 days put "Last day" and "Last 7 days" on the same one-hour window; 40 days collapsed three chips; 100 days collapsed four of five → ✓ after
- Arrow inverse: ✗ before — round-trip drift −120.00h a partial screenful from the right bound, −168.00h flush against it, 0.00h mid-range → ✓ after

**Generalisation by sweep, not by the measured cases.** 141 board ages against the chip set read out of the shipped page, and 167 window positions in both directions, each with a vacuity guard that fatals if the sweep did not actually sweep. Seven mutations, seven kills — including one written specifically because the bounds assertion would otherwise have been unfalsifiable. Render evidence on a purpose-built 40-day-drained fixture: before, three of five chips landed on one shared dead window with both arrows dead; after, every chip lights itself and the arrow round trip returns exactly.

**Independently re-measured by the orchestrator on the merged tree**, using the same probes that found the defects:

```
drained board     idle=0d/3d/10d/40d/100d -> distinct=5/5, on-floor=0 at every value
arrow round trip  mid-range 0.00h | one span from bound 0.00h | flush at left bound 0.00h
                  partial screenful out -> step refused, window unmoved
                  flush at right bound  -> step refused, window unmoved
```

*Verified by work action*

## Review

**Acceptance: Pass.** Both acceptance criteria are met and were verified twice — by the builder's sweeps and mutations, and independently by the orchestrator using the original defect probes. Route B, so standard depth rather than a panel: two files, each behaviour pinned by a test that fails before the fix, and the outcomes confirmed directly rather than read from a claim.

**On the design call.** D-01 makes the toolbar arrows all-or-nothing: a step either moves a whole screenful and round-trips exactly, or refuses. The reviewer's suggested alternative — mirroring the clamped forward step — turned out **not implementable as a pure function of the window**, because many distinct windows map to the same bound-flush window, so the back step cannot know how far the forward step was clamped. That collapses the apparent two-way trade into refuse-or-keep-the-defect, and the builder said so rather than pretending to follow the suggestion. The reachability cost is smaller than the finding assumed: the last partial screenful stays reachable by drag, keyboard pan, Fit all, All days and the date fields — only reaching it with the arrows in one press is given up. `timelinePannedWindow` is deliberately untouched, so keyboard pan and drag keep their continuous clamp.

**Deletion over addition.** Returning `isClippedByBounds` let the state readout stop recomputing the clamp, removing a branch the new anchor makes unreachable.

**Carried forward, not fixed here:** the keyboard's `←` and `→` still carry the drift the toolbar arrows shed, by design — recorded in the hand-back's Discovered Tasks. A deliberate consequence of confining D-01 to the toolbar arrows.

*Reviewed by work action (Route B, standard depth, with direct orchestrator verification of both acceptance criteria)*

## Lessons Learned

**What worked:** Sweeping the parameter space instead of fixing the measured cases. The defects were reported at 3, 10, 40 and 100 days idle; the fix was proven across 141 board ages and 167 window positions, each sweep carrying a guard that fails if it did not actually sweep. That guard is the part worth copying — a sweep that silently covers nothing looks exactly like one that passes.

**What didn't:** The review's proposed remedy for the arrows could not be built. Mirroring a clamped step needs information the clamped window no longer carries. A proposed fix inside a finding is a hypothesis, and testing it before adopting it is the builder's job.

**Worth knowing:** `timelineRange` in `timeline.go` stretches `rangeEnd` to `now` only while a row has `WaitOpen` or `WorkOpen`, so a drained board is the case that breaks any Timeline maths assuming `now` is inside the bounds. Both defects came from that one assumption, which is why they were fixed as one sweep rather than two REQs.

## Orientation

The Timeline's window controls now behave the same on a finished board as on a live one. Where a queue has drained, each chip reads as the last N days of the recorded range instead of collapsing onto a shared one-hour window, and the toolbar arrows either move a full screenful or refuse, so forward-then-back always returns. It lives in the queue-kanban board's Timeline view, in `web/board-timeline.js`.

No `[MAP CHANGED]`: this closes two edge behaviours in the controls REQ-390 introduced and adds no new concept. `_dev/primes/prime-kanban-board.md` documents neither the trailing-window maths nor the step arrows, so it stays accurate.

