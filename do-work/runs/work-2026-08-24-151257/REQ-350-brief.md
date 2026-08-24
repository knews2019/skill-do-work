# REQ-350 builder brief

Build REQ-350, “Narrow the Durations axis to a chosen window,” on Route B with TDD.

## Work location

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-350-narrow-the-durations-axis-window`
- Branch: `worktree-agent-REQ-350-narrow-the-durations-axis-window`
- Hand-back (the only permitted main-tree write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-151257/REQ-350-handback.md`

Never read or write any path under `do-work/` in the worktree or main checkout except the absolute hand-back path above. Never write the main working tree. Commit implementation and tests on the named branch.

## Governing context

Read in the worktree before coding:

- `_dev/primes/prime-kanban-board.md`
- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- `skills/do-work/crew-members/frontend.md`
- `skills/do-work/crew-members/testing.md`

The Durations axis currently always spans full sample history, wasting width on idle calendar. Add a separate Durations control for last 30 days, last 90 days, and all history; every panel and accessible surface must describe the same affine real-time domain.

## Required behavior

- Add a Durations-only topbar control using the existing control-group/button/active/`aria-pressed` pattern with distinct `data-durations-window="30|90|all"` attributes. It is hidden outside Durations and does not read or mutate the board's `viewState.windowHours`/`data-window-hours` value.
- Default to last 30 days. Anchor 30/90-day whole-UTC-day domains to the UTC day of the global latest completion, ending at the following midnight; static boards must not age into emptiness. `all` uses the current first-sample day through the day after the latest sample.
- Filter samples and precomputed days once through the half-open `[timeStart,timeEnd)` domain and use those projected arrays for summary/exclusions, Panel A marks and UR lane, colour context, Panel B/C geometry and annotations, hover indexes, and sample table.
- Keep epoch-to-x affine in real milliseconds. Idle gaps remain proportionally visible; never compress inactive days.
- Summary and SVG accessibility copy name Last 30 days / Last 90 days / All history, chosen-domain endpoints, and the projected sample count while preserving exclusion wording.
- Generate a board and inspect all three settings.

## Declared scope

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`

Existing control CSS is sufficient. No Go payload, aggregation, `board.js`, or browser-probe change is expected. If an integration seam outside scope is unavoidable, do not edit it; record its exact path and line in the hand-back.

## TDD and proof

Add RED behavior evidence first: template/default/hidden/distinct-state contract; 30/90/all transitions with exactly one Durations rerender and unchanged board recent-window state; a >90-day production-renderer fixture proving projected summary, mark/table counts, Panel B/C active-day counts, left-boundary inclusion, and all-history completeness; and an affine equal-day-spacing proof. Update the existing full-history day-bucket fixture to select `all` explicitly so its 400-day subject remains intact. Then implement, run focused and relevant canonical tests, and inspect each selection in a real browser. Record `location.href`, active/pressed button, summary, Panel A/B/C counts, table count, first/right axis labels, affine bar centres, console state, and desktop screenshots. Also check control wrapping/keyboard usability at 320, 768, and 1280 px.

## Hand-back format

Write the absolute hand-back file before finishing with:

- branch and commit hash;
- exact file manifest and test commands/results;
- RED and GREEN evidence;
- per-window rendered measurements, browser/build, `location.href`, and artifact paths;
- integration seams (or `none`);
- `## Decisions`;
- `## Discovered Tasks`;
- risks or unresolved issues.

The hand-back pattern makes fan-out failures survivable, not prevented.
