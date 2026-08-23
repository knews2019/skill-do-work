---
id: UR-068
title: Durations panel improvement proposal
created_at: 2026-08-23T22:37:52Z
requests: [REQ-342, REQ-343, REQ-344, REQ-345, REQ-346, REQ-347, REQ-348, REQ-349, REQ-350]
word_count: 1657
---

# Durations Panel Improvement Proposal

## Summary

An in-repo design proposal (`ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html`, written 2026-08-23, nothing shipped) records eight findings against the board's Durations view and offers eight capture-request prompts plus one bundled invocation. The user chose the bundled form — one UR, its changes sliced into REQs — and answered four questions during capture, two of which changed the proposal's shape.

The view's core problem is that it is the only board view that drops the REQ to UR link: its payload sample carries `id`, `route`, `completionTime`, `dayKey`, `wallMinutes` and nothing else, though the client already holds the join at `boardData.requests[id].userRequestId`. Behind that sit two structural problems measured on this repository's own archive (305 samples, 29 active days, 66 URs): two thirds of the axis is idle calendar, and 55 percent of the marks sit in the bottom quarter of panel A's linear 0-60 scale.

## Extracted Requests

| REQ | Title | Source in input |
| --- | --- | --- |
| REQ-342 | Name the UR on every Durations sample | Prompt A1 (finding F1) |
| REQ-343 | Colour Durations marks by a chosen channel | Capture answer to the report's Q3 |
| REQ-344 | Group the Timeline's rows by user request | Capture answer to the report's Q1, replacing prompt A2 (finding F2) |
| REQ-345 | Fix panel A's scale and density | Prompt A4 (finding F4) |
| REQ-346 | Narrow the Durations axis to a chosen window | Prompt A3 (finding F3) |
| REQ-347 | Retire the Durations in-lane labels | Prompt A5 (finding F5) |
| REQ-348 | State the Durations view's own headline numbers | Prompt A7 (finding F7) |
| REQ-349 | Hide the dead filter knobs while Durations is on screen | Prompt A6 (finding F6) |
| REQ-350 | Open the detail drawer from a Durations mark | Prompt A8 (finding F8) |

## Capture Decisions

Resolved interactively during capture (ask tool):

- **Q1 (queue shape)** — user chose **one bundled UR**. The eight prompts become REQs under this UR rather than eight separate URs, so `related`/`batch` carry the arc and the shared write set on `web/board-durations.js` is visible on the board.
- **Q2 (prompt A5, the label reversal)** — user chose **yes, capture it**. A5 reverses the separated-text-band direction settled on 2026-08-18; the explicit yes is recorded here because that decision was deliberate and worked. REQ-347 carries it.
- **Q3 (the report's own Q1, where panel D belongs)** — user chose **UR grouping on the Timeline instead**. Prompt A2 is therefore **not** captured as written: no panel D is added to Durations. REQ-344 puts the UR-level view on the Timeline, which already draws REQ bars against calendar time. Finding F2 (a UR has no measured duration at all) is still the thing being closed, just on the other view.
- **Q4 (the report's own Q3, how a mark says which UR it belongs to)** — user chose **lane plus colour-by toggle**. The lane alone was the proposal's recommendation; the user wants both channels. Sliced into REQ-342 (the lane, hover and table) and REQ-343 (the colour-by control), because the toggle brings its own caption rule and a second identity channel competing with route.

Decided by capture, from the report's own reasoning, without asking:

- **The report's Q2 (shared window state)** — REQ-346 reuses the topbar's recent-window control *pattern*, not its *value*. The board's window means "recently done" and the Durations window means "the axis domain"; conflating them would move the Durations axis when someone adjusts a board filter. Stated as a constraint on REQ-346.
- **The report's Q4 (how many lane rows)** — fixed rows plus an explicit counted remainder, which is the pattern the overflow lane already uses. Three rows placed 54 of 61 URs on 30 days of this repository's data, so the row count is the builder's call; the counted remainder is not. Stated as a requirement on REQ-342.

## Batch Constraints

- Board tool root: `skills/do-work-board/tools/queue-kanban/`. `_dev/primes/prime-kanban-board.md` governs every REQ here — read it first, in particular the render-evidence rule and the measured-face-per-browser rule.
- **Seven of the nine REQs write `web/board-durations.js`.** Declared in every affected `write_set`. Do not fan these out into parallel worktrees; they serialize on that file.
- **Every REQ here changes a chart.** The repository's standing rule applies to each: generate a board and look at it. A passing suite is not evidence about two glyphs sharing a coordinate.
- Keep one measure per scale, and keep the read-time exclusion rule (spans over four hours assumed paused, negative spans dropped) exactly as it is unless a REQ says otherwise. Panels A and B disagree about which samples count on purpose, and that rule lives in one place.
- The figures quoted throughout come from this repository's generated `board-data.js` on 2026-08-23 (305 samples, 29 active days, 66 URs, median 13.2 min, p90 42.7) and from the consuming project's board (692 samples, 47 active days). They are measurements, not targets — re-measure rather than hard-coding them.

## Full Verbatim Input

read ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html
ask your questions
then do-work capture-request
and do-work verify-requests

### The bundled invocation from the report (the input being captured)

do-work capture-request: Improve the board's Durations view so it answers UR-level questions and stays legible at 700 or more archived REQs. It is the only view that drops the REQ-to-UR link, so give its samples their UR: name the UR in the hover readout, add a UR column to the sample table, and draw a packed UR grouping lane under panel A with an explicit remainder for any UR that finds no row — the client already holds the join at boardData.requests[id].userRequestId. Add a panel D measuring URs themselves: one stem per UR at its last REQ completion, from summed REQ work up to elapsed calendar span, log hours, radius by REQ count, and say in the REQ how it differs from the Timeline's REQ Gantt. Fix panel A's density: square-root y scale with ticks at 0, 5, 15, 30, 45 and 60, deterministic within-day jitter bounded to the day slot, lower mark opacity, and a per-day median line with a p25-to-p75 ribbon. Add a window control for the last 30 days, last 90 days and all history that keeps the axis linear real time inside the window and states the active window in the summary. Replace the overflow lane's in-lane labels with a ranked longest-spans list beside the chart carrying REQ, UR, duration, route and title, and delete the label planner, its width model and the tests that exist only to place text inside the SVG. Add a stat tile row for median, p90, active days against axis span and REQs per active day, a rolling median line on panel B, and the missing zero and midpoint ticks on panel C. Hide the text filter, domain select and status select while Durations is on screen, since they change nothing there and the board's own rule is that the topbar never advertises dead knobs. Make a click on a mark open the detail drawer for the nearest REQ and give the view a keyboard path that does not require opening the collapsed sample table.

### The eight individual prompts from the report, verbatim

**A1 — Give the samples their UR.** do-work capture-request: The Durations view is the only board view that drops the REQ-to-UR link, so panel A's dot cloud cannot say which user request a mark belongs to — 66 URs are in this repository's 305 samples and a median of 2, up to 10, are active on the same day. Give the Durations view UR identity. The client payload already carries the join at boardData.requests[id].userRequestId, so prefer a client-side join and add the field to generatedDurationSample in generate.go only if that proves awkward. Name the UR in the hover readout beside the REQ id, add a UR column to the "Every sample, as a table" table, and draw a UR grouping lane under panel A: one horizontal bracket per UR spanning its samples' completion times, packed into a small fixed number of rows, with every UR that finds no row counted in an explicit remainder rather than silently dropped. Leave route colour as it is.

**A2 — Measure the UR itself.** do-work capture-request: Durations measures REQs only, so a UR's own duration is unmeasurable from the board. Add a panel D to the Durations view: one mark per UR on the shared calendar axis, positioned at the UR's last REQ completion, drawn as a stem from the summed measured REQ work at the bottom up to the elapsed calendar span from first REQ claim to last completion at the top, on a log hour scale, with the top mark's radius carrying the UR's REQ count. The gap between the two ends is the point: on this repository's archive the median UR spans 0.6 hours against a p90 of 46.7, and the largest URs (27, 22 and 18 REQs) spend 11 to 14 percent of their elapsed span in measured REQ work. Apply the same read-time exclusion rule the other panels use, keep one measure per scale, state in the panel subtitle what each end of a stem means, and say explicitly in the REQ how this differs from the Timeline view's REQ Gantt so the two views do not converge.

*(Not captured as written — see Capture Decisions Q3. REQ-344 carries the finding onto the Timeline instead.)*

**A4 — Fix panel A's scale and density.** do-work capture-request: Panel A of the Durations view wastes its vertical range and overplots. On this repository's archive the median REQ is 13.2 minutes with 55 percent under 15 and 78 percent under 30, so against a linear 0-to-60 scale more than half the marks compete for the bottom quarter of the panel, and at 8 to 13 SVG units per day a busy day (38 REQs here, 55 on the consuming project) is one column of overlapping 4-unit dots. Change panel A to a square-root y scale with ticks at 0, 5, 15, 30, 45 and 60; add a deterministic within-day x jitter bounded to the day's own slot so a mark can never cross a day boundary; lower mark opacity; and draw a per-day median line with a p25-to-p75 ribbon behind the marks so the trend is readable without hovering. Keep the 60-minute ceiling, the overflow lane and the read-time rule unchanged, and keep the jitter out of the hover's nearest-mark maths or state how it is compensated.

**A3 — Let the reader narrow the window.** do-work capture-request: The Durations axis spends most of its width on idle calendar: on this repository's archive 66 percent of the 86-day axis has no completions and 92 percent of the samples fall in the final 30 days, and on the consuming project a single April outlier holds the left half of the plot open by itself. Add a window control to the Durations view offering the last 30 days, the last 90 days and all history, and state the active window and the sample count it covers in the view's summary line. Keep the axis linear real time inside the window — the idle-gap cadence reading is why it is linear, so do not compress gaps — and reuse the topbar's existing recent-window control pattern rather than inventing a second one. Panels B and C follow the same window.

**A5 — Retire the in-lane labels.** do-work capture-request: The Durations overflow lane's direct labels have stopped paying for themselves at scale — the consuming project's board shows 5 labels beside "+55 more over 60 min" and this repository's shows 4 beside "+9 more", a 7 percent labelling rate — while the machinery behind them (browser-measured text widths, a greedy two-row packer, a remainder reserve run to a fixed point, and the geometry-agreement tests pinning durations.go's constants against web/board-durations.js) is the largest single piece of the view. Replace the in-lane labels with a compact ranked longest-spans list beside the chart, carrying REQ id, UR id, duration, route and title, longest first, with every over-ceiling sample present. Keep the lane, its marks, their hover and a plain remainder count. Delete the label planner, the label-text and width model in durations.go, and the constants and tests that exist only to place text inside the SVG, and name in the hand-back exactly which tests were removed and why they no longer describe a shipped rule. This deliberately reverses the separated-text-band direction chosen on 2026-08-18; that fix worked, and the reason to revisit it is the sample count, not a defect.

**A7 — Make the view state its own numbers.** do-work capture-request: The Durations view makes a reader hover to learn its own headline numbers. Add a small stat tile row above panel A carrying the median REQ duration, the p90, the active-day count against the axis span, and REQs per active day, each stated for the window actually plotted rather than for all history. Give panel B a rolling median line over its per-day bars so "is it getting slower" is answerable without reading 29 bars one at a time, and give panel C the zero and midpoint ticks it lacks — today only its peak is labelled. Derive every figure from the samples already in the payload, keep one measure per scale, and keep the existing summary sentence's exclusion-rule wording intact.

**A6 — Stop advertising dead knobs.** do-work capture-request: The topbar's text filter, domain select and status select stay visible and live while the Durations view is on screen and change nothing there. That is inconsistent with the board's own stated rule: applyView in web/board-controls.js hides the lens and recent-window groups on views where they do nothing, with the comment that the topbar should never advertise dead knobs, and onFiltersChanged in web/board-filters.js deliberately excludes Durations because a filtered distribution is a different statistic wearing the same axes. Hide the three filter controls while Durations is the active view, the same way the lens group is hidden, and restore them on view change. If a filtered durations distribution is wanted later, that is a separate request with its own caption rules, not a silent widening of this one.

**A8 — Give a mark somewhere to go.** do-work capture-request: A Durations mark is a dead end — the nearest-mark hover writes one line into a status paragraph and nothing opens. Make a click on panel A's hover surface open the detail drawer for the nearest REQ, the same drawer the board and timeline views open, and give the view a keyboard path to the same information: the sample table is the only keyboard route today and it sits inside a collapsed details element. Keep the hover readout exactly as it is, and keep the click's nearest-mark resolution consistent with the hover's so a reader never opens a REQ other than the one the readout names.

### Source

Full proposal, including the eight findings with file references, the two mockups, the cost/risk table and the four open questions: `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html`. Two authentic captures in that bundle's `screenshots/`.

---
*Captured: 2026-08-23T22:37:52Z*
