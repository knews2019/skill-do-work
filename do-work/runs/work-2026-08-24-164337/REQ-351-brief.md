# REQ-351 builder brief

Build REQ-351, “Retire the Durations in-lane labels for a ranked longest-spans list,” on Route B with TDD.

## Work location

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-351-retire-the-durations-in-lane-labels`
- Branch: `worktree-agent-REQ-351-retire-the-durations-in-lane-labels`
- Hand-back (only permitted main-tree write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-164337/REQ-351-handback.md`

Never read or write any `do-work/` path except that absolute hand-back. Commit all implementation/test changes on the named branch.

## Governing context

Read `_dev/primes/prime-kanban-board.md` and general, coding-guardrails, frontend, testing, and communication-style crew rules before editing.

## Required behavior

- Replace variable-width direct SVG overflow/reversed labels with a complete HTML longest-spans list beside the chart for the current Durations window.
- Include every positive over-60-minute sample, ordered descending by `wallMinutes` with a deterministic tie-break. Each item exposes REQ id, UR id/name, duration, route, and title using the existing `boardData.requests` join.
- Keep the overflow lane, all lane marks, exact hover readout, UR brackets, window projection, full sample table, and a plain complete-count sentence. Do not use `+N more`; nothing is omitted.
- Keep all list items in the DOM. Use a chart/list wrapper with a bounded scrollable desktop aside (~280–360px) and stacked responsive layout below about 1000px. Verify 60+ rows and 320/768/1280 widths without new horizontal clipping.
- Remove the JavaScript label planner/leader/row geometry/fixed-point remainder reservation and browser text measurement that exist only for direct SVG labels.
- Remove the Go label-text/width model and label-only fixtures/tests. Preserve `durationLabelTimeRange`/`durationLabelPlotX` and any renderer constants/helpers still consumed by day-bucket, jitter, Panel B annotation, window, lane, hover, or formatter tests.
- Preserve `.durations-mark-label` if Panel B annotation still uses it; remove only label-specific leader styles.

## Declared scope

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/durations.go`
- `skills/do-work-board/tools/queue-kanban/durations_test.go`
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/web/board.css`

Do not edit payload generation: existing data suffices. If an outside-scope seam is unavoidable, record it rather than editing it.

## Test deletion accounting

Name every removed test in the hand-back and explain that its shipped SVG-placement rule ended. Exploration identified these label-only browser tests for removal: drawn labels never overlap, labels go to longest spans, clustered labels fill both rows, bands pack independently, remainder counts what was not drawn, label rows clear neighbours, no measured-face constants, and placement responds to rendered face, plus their label-placement harness/types. Preserve generic style slicing and the dense Panel A spread/hover/distribution proof.

Delete the three label-planner production-renderer tests in `generate_test.go` (`DurationsLabelRowsAndRemainders`, `LabelWidthModelMatchesTheRendererFormatter`, `ReserveMatchesTheSentenceDrawn`) and the Go rounded-label/label-fixture tests only when their entire claim is retired. Keep the generic span formatter test.

## RED/GREEN and rendered proof

RED must show no complete adjacent HTML list and bounded SVG labels omitting samples. GREEN on a generated 705+ sample fixture with at least 60 positive overflow samples must assert: one complete list; DOM count equals current-window over-ceiling count; descending/tie order; every row's five fields; list outside SVG but adjacent in one wrapper; no SVG REQ label text/leader lines; lane circles and exact overflow hover remain; sentence count equals list count. Inspect light/dark and 320/768/1280 responsive layouts, record browser/build, `location.href`, console, measurements, and screenshots.

Run focused tests, full module/canonical gates proportionately, syntax/vet/diff checks, commit, and write the hand-back with manifest, RED/GREEN, exact deleted-test accounting, visual evidence, seams, `## Decisions`, `## Discovered Tasks`, and risks.

The hand-back pattern makes fan-out failures survivable, not prevented.
