# REQ-425 Builder Brief

Worktree: `/home/user/skill-do-work-worktrees/worktree-agent-REQ-425-trailing-window-bound-assumptions`
Branch (operative name): `worktree-agent-REQ-425-trailing-window-bound-assumptions`
Hand-back file (absolute, main tree): `/home/user/skill-do-work/do-work/runs/work-2026-08-29-213539/REQ-425-handback.md`

---

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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
