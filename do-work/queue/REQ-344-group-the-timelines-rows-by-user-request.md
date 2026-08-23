---
id: REQ-344
title: "Group the Timeline's rows by user request"
status: pending
created_at: 2026-08-23T22:37:52Z
user_request: UR-068
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-342, REQ-343, REQ-345, REQ-346, REQ-347, REQ-348, REQ-349, REQ-350]
batch: durations-panel-improvement
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# Group the Timeline's Rows by User Request

## What

A UR has no measured duration anywhere on the board: Durations measures REQs only and the Timeline
rows them one per REQ, so the unit the maintainer actually plans in is invisible. Group the
Timeline's rows by UR, and make each UR's own span readable against the work inside it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Computed from this repository's own tickets: 65 URs have measurable spans, the median UR takes
0.6 hours end to end while the p90 takes 46.7, and the largest ones (UR-031 at 27 REQs, UR-042 at 22,
UR-055 at 18) spend 11 to 14 percent of their elapsed span in measured REQ work. That gap between
elapsed and worked is a real property of the pipeline and nothing on the board shows it.

**This REQ replaces the proposal's prompt A2, on the user's decision.** A2 proposed a fourth panel on
the Durations view: one stem per UR from summed REQ work up to elapsed calendar span, log hours,
radius by REQ count. The user chose the alternative the report named in its own Q1 — the Timeline
already draws REQ bars against calendar time, and a UR-grouped panel D would have sat one tab away
showing the same tickets against the same axis. **No panel D is added to Durations.**

## Detailed Requirements

- **Rows group under their UR.** A UR header row spans from its first REQ's start to its last REQ's
  end; its REQ rows sit under it.
- **The elapsed-versus-worked gap is readable from the group**, because that is the finding being
  closed, not the grouping itself. The UR header carries both its elapsed calendar span and the
  summed measured work of its REQs — in the row label, the drawer, or the hover, whichever the
  view's existing idiom supports.
- **The REQ count per UR is visible** on the header row.
- A REQ with no `user_request` still gets a row. Decide where ungrouped rows sit and state it.
- **Grouping is compatible with what the view already does**, or it does not ship:
  - Rows are virtualized (only rows inside the scrolled window get SVG nodes) — a group header must
    not force the whole archive into the DOM.
  - Rows are window-scoped (REQ-319: a row belongs on screen when its bar overlaps the visible time
    window) — decide whether a UR header spans the whole UR or only the part the window covers, and
    say which. A header claiming span the window excludes would repeat exactly the bug REQ-319 fixed.
  - The row list is one Tab stop with a roving index (REQ-338). Group headers must not reintroduce a
    Tab stop per row.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first — the render-evidence rule and
  the measured-face-per-browser rule both apply.
- `TimelineRow` carries no UR today. Prefer the same client-side join Durations uses
  (`boardData.requests[id].userRequestId`); add the field to the payload only if that proves awkward,
  and say which way it went.
- The payload decides the numbers: both spans arrive already measured against one `now`, signed, with
  the board's reversed-stamp verdict attached, and `web/board-timeline.js` never re-measures. A summed
  work figure is an aggregation of measured spans, not a new measurement — keep it that way.
- Generate a board and look at it.

## Builder Guidance

**Certainty: firm on the finding and on the view, open on the presentation.** The user's decision is
"UR grouping on the Timeline, no panel D" — the shape of the grouping is yours. Grouping that only
indents rows and does not surface the elapsed-versus-worked gap does not close F2; the gap is the
whole reason this REQ exists.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repository and open the Timeline. Every row is one
REQ, ordered by created_at; nothing states which UR a row belongs to, how long UR-031 took end to
end, or how much of that span was measured work.

**Why RED now:** `TimelineRow` carries `RequestId` and no UR, and the renderer draws one flat
newest-first row list.

**GREEN when:** the Timeline's rows are grouped under their UR, each group states its elapsed span,
its summed measured work and its REQ count, virtualization and window-scoping still hold, and the row
list is still one Tab stop.

**Validation:** User adjusted — the user replaced prompt A2's panel D with this during capture.

---
*Source: capture answer to the report's open question Q1, replacing prompt A2 (finding F2), `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html`.*
