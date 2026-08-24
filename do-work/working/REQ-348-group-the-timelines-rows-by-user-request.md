---
id: REQ-348
title: "Group the Timeline's rows by user request"
status: claimed
claimed_at: 2026-08-24T15:00:04Z
status_changed_at: 2026-08-24T15:00:04Z
route: C
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-346, REQ-347, REQ-349, REQ-350, REQ-351, REQ-352, REQ-353, REQ-354]
batch: durations-panel-improvement
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-08-23T22:37:52Z
  basis:
    - Route C
    - 3-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Group the Timeline's Rows by User Request

## What

A UR has no measured duration anywhere on the board: Durations measures REQs only and the Timeline
rows them one per REQ, so the unit the maintainer actually plans in is invisible. Group the
Timeline's rows by UR, and make each UR's own span readable against the work inside it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Model UR groups client-side after window filtering; flatten fixed-height group headers and REQ rows for virtualization while preserving a REQ-only roving focus index. Reuse Durations samples for canonical summed-work exclusions, render honest open/no-claim metrics, and cover the model plus live browser behavior before implementation.
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

- **Rows group under their UR.** Its REQ rows sit under it.
- **An open UR needs a stated endpoint, not a missing one.** "Last completion" is only defined for a
  UR whose work has finished. The current queue holds both open shapes: UR-068 and UR-069 have no
  completion at all (4 and 9 REQs, every one still open), and UR-065 and UR-067 have work still open
  *after* their last completion — where ending the group at that completion yields a stale span that
  can read shorter than the summed work inside it. Define the open-UR endpoint (the aggregate's
  frozen `now` is the obvious candidate, and is what the Timeline already measures open spans
  against) and mark such a group as still running rather than drawing it as a closed span.
- **The elapsed span runs from first REQ *claim* to last completion** for a finished UR, which is what
  the captured request asks for and what the figures quoted below were measured from. Not from
  `created_at`:
  Timeline rows distinguish `CreatedTime` from `ClaimedTime` and draw their wait segment from
  creation, so measuring the group from a row's visual left edge would fold queue wait into the
  number and produce a different statistic from the one this REQ reports. If a UR's earliest REQ was
  never claimed, state the fallback you chose on the view rather than silently substituting its
  creation stamp.
- **The elapsed-versus-worked gap is readable from the group**, because that is the finding being
  closed, not the grouping itself. The UR header carries both its elapsed calendar span and the
  summed measured work of its REQs — in the row label, the drawer, or the hover, whichever the
  view's existing idiom supports.
- **The REQ count per UR is visible** on the header row.
- **Say which figure is which.** Prompt A2 required the panel subtitle to state what each end of a
  stem meant; the same obligation survives the move to the Timeline. Elapsed span and summed work are
  two different measures on one row, and an unlabelled pair is a worse answer than no answer.
- **Name the exclusion rule the summed work uses.** Prompt A2 required the same read-time exclusion
  the Durations panels apply (spans over four hours assumed paused, negative spans dropped). The
  Timeline has its own convention — it draws every span the payload measured, with the board's
  reversed-stamp verdict attached — so the two rules can disagree about the same UR. Pick one, state
  it on the view, and make sure a reader comparing this figure with the Durations panels is told they
  are not the same statistic if they are not.
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

## Dependencies

None declared. **REQ-345 also writes `web/board-timeline.js`** and is `impact-critical`: adding
pending REQs to the queue currently fails
`TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable` and red-lights the canonical gate. It
lands first — not because this REQ needs it, but because no hand-back here can be proven green until
it does. Read its findings before touching the file; whichever of the two lands second rebases onto
the other.

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

## Triage

**Route: C** — Timeline grouping changes the row model, virtualization, keyboard behavior, aggregation semantics, and rendered presentation. Planning and browser-backed exploration are required.

## Plan

1. Join each window-listed Timeline row to its request's `userRequestId` and completed work to the existing Durations sample verdict. Group only after window membership is decided, preserve row and first-seen group order, and place `No UR recorded` last.
2. Compute window-scoped group count, earliest recorded claim, honest finished/open endpoint, elapsed availability, accepted-work sum, and exclusion count. Flatten headers plus members for virtual scrolling without admitting headers to the roving Tab sequence.
3. Render clipped group spans and explicitly labelled elapsed/worked/count copy, update the accessible table and explanatory text, and add theme-aware styling.
4. Add model and browser probes for ordering, endpoints, exclusions, virtualization, window changes, and keyboard behavior; generate and inspect light and dark boards.

## Exploration

- The existing request payload supplies UR identity, so `timeline.go` does not need a schema change.
- `boardData.durations.samples` already owns the read-time exclusion verdict; consuming `excludedReason` avoids duplicating the four-hour/negative-span policy.
- Window-listed REQs remain authoritative for summaries, table rows, projections, and group membership. Only the virtualized display list gains headers.
- A group without a recorded claim reports elapsed as unavailable; open groups end at the payload's frozen `now`; finished groups use their latest resolved completion.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js`
- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`
