---
id: REQ-348
title: "Group the Timeline's rows by user request"
status: completed
claimed_at: 2026-08-24T15:00:04Z
completed_at: 2026-08-24T16:41:18Z
commit: 205be83
status_changed_at: 2026-08-24T16:41:18Z
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
- [x] **[APPLY]:** The isolated builder implemented client-side window-scoped UR groups, honest elapsed/worked/count metrics, mixed-height virtual display items, REQ-only focus mapping, accessible table grouping, explanatory copy, and theme-aware presentation in the five scoped files.
- [x] **[UNIFY]:** Builder reviewed all five files and passed pure/browser grouping, virtualization, keyboard, accessibility, visual contrast, canonical maintainer verification, and diff checks. Orchestrator merged over the completed Durations wave and re-ran the combined focused Timeline behavior/browser probes green; no integration seams or debug artifacts were present.

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

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified): joins window-listed rows to URs and accepted Durations work, builds stable groups and flattened display items, renders honest group spans/metrics/table rows, and maps REQ focus through virtualized headers.
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified): adds light/dark group-header, span, label, and metric styling against the actual body surface.
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified): updates Timeline heading/hint copy with window membership, endpoint, exclusion, fallback, and ungrouped semantics.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified): pins grouping order, unknown-last placement, metrics/endpoints/exclusions, post-window membership, flattening, and bounded mixed-height virtualization.
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified): proves rendered header/member ordering and counts, rebuilt groups after window changes, bounded SVG nodes, REQ-only roving keyboard behavior, and light/dark readability.

## Decisions

- **D-01 — Group after window filtering.** Every header count, metric, and span covers exactly the REQs the chosen Timeline window lists.
- **D-02 — Keep `TimelineRow` unchanged.** UR identity is joined client-side through `boardData.requests`, matching the Durations idiom and avoiding payload duplication.
- **D-03 — Reuse the canonical Durations exclusion verdict.** Summed work consumes `excludedReason` instead of recreating duration validity policy.
- **D-04 — Keep focus REQ-indexed.** Mapping REQ indices to flattened display indices allows non-focusable virtualized headers while retaining one Tab stop and existing Up/Down/Left/Right/Enter behavior.
- **D-05 — Use frozen payload `now` for running groups.** Static reports remain deterministic and never present an open UR as closed at a stale last completion.

## Discovered Tasks

- Recalibrate `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` for Chromium headless shell 147: under a forced full-browser run, swallowing `setPointerCapture` no longer caused the expected dependency and the drag still completed.
- `TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening` failed once under forced full-suite load but passes alone; focused REQ-348 browser probes and canonical verification are green.

## Testing

- Initial RED proved no Timeline grouping function existed. Model GREEN pins stable first-seen groups, newest-first/unknown-last members, endpoints, accepted-work exclusions, post-window grouping, flattening, and mixed-height virtualization.
- Builder canonical maintainer verification passed; focused forced Chromium proved header/member order, rebuilt windows, bounded nodes, one REQ Tab stop, keyboard behavior, explicit table associations, and light/dark readability. Measured header contrast was 17.26:1/6.11:1 in light and 15.83:1/6.88:1 in dark.
- First review produced targeted RED: a 09:00–13:00 member in a 10:00–12:00 window reported 240 rather than 120 minutes; isolated/mixed unresolved completions invented or partially published elapsed; member cells lacked exact group/column associations.
- Remediation GREEN proves exact clipped 10:00/12:00 endpoints and 120-minute elapsed/work, endpoint-unavailable suppression for isolated and mixed unresolved groups, and live per-cell own-group + column headers with REQ `th scope="row"`.
- Remediation maintainer verify passed (queue-kanban 61.420s, strict JavaScript 17.11s, audit-metrics green); post-merge full suite passed in 67.833s. One combined focused Chromium run transiently read stale Fit-all contents, then passed alone and as the identical full focused group; re-review reproduced this only once under repeated-count load.

## Qualification

- Cumulative merge range `5e0d9e9..205be83` passed mechanical qualification and exact five-file scope drift.
- Orchestrator judgment confirmed substantive behavior, canonical Durations exclusion data flow, window-bounded metrics, endpoint honesty, REQ-only focus mapping, and accessible header associations.

## Review

First review scored 86% and required remediation of three Important findings: unbounded metrics, unresolved endpoint copy, and invalid rowgroup semantics. Re-review approved with Requirements 100%, Code Quality 97%, Test Adequacy 92%, Scope 100%, overall 97%, low risk, acceptance pass. One non-blocking fixed-wait synchronization Minor was captured as REQ-369 without weakening this release.

## Lessons Learned

Window-scoped membership is not sufficient evidence for window-scoped metrics: intersect every contributing interval with the visible domain. Accessibility semantics likewise need explicit structural associations, not a visually plausible header row.

## Orientation

Released in 0.236.44. Timeline rows now group beneath their UR with window-clipped elapsed and accepted-work metrics, honest open/unresolved endpoints, bounded virtualization, REQ-only keyboard focus, and explicitly associated accessible table headers.
