# REQ-352 builder brief

Build REQ-352, “State the Durations view's own headline numbers,” on Route B with TDD.

## Work location

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-352-state-the-durations-headline-numbers`
- Branch: `worktree-agent-REQ-352-state-the-durations-headline-numbers`
- Hand-back (only permitted main-tree write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-164337/REQ-352-handback.md`

Never read or write any `do-work/` path except that absolute hand-back. Commit all implementation/test changes on the named branch.

## Governing context

Read `_dev/primes/prime-kanban-board.md` and general, coding-guardrails, frontend, testing, and communication-style crew rules before editing.

## Required behavior

- Add a semantic four-item `<dl>` stat row above Panel A, using an auto-fitting responsive grid:
  - median REQ duration and p90 from all projected Panel A raw spans, including paused/reversed samples, visibly labelled “all plotted spans”;
  - active completion days as projected `days.length` against whole axis days;
  - projected REQs per active day as `samples.length / days.length`, one decimal.
- Tiles update under 30/90/all Durations windows. Preserve the existing summary exclusion-rule sentence byte-for-byte.
- Panel B rolling series consumes chronological `days.filter(day.hasMedian)` only. For index six onward, compute the R-7 median of that day plus the six previous drawn-day `medianMinutes`; place at real-calendar day noon. Idle and excluded-only days are skipped, never zero-filled.
- Draw rolling path after Panel B bars and add visible point markers. Six or fewer eligible days draw nothing; exactly seven draw one marker; eight draw two markers and a connecting path. State “trailing 7-active-day median” in visible Panel B title and SVG accessibility copy.
- Panel C retains peak/baseline and adds a zero tick plus exact `peakCount / 2` midpoint tick/gridline. Odd midpoint labels use one decimal.
- No Go payload/aggregation/control changes.

## Declared scope

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`

REQ-351 concurrently changes all five files. Stay isolated and do not inspect or merge its branch; the orchestrator will integrate serially and preserve both surfaces.

## TDD and proof

Drive the complete production renderer. Prove all four tiles change for 30/90/all, raw tiles include excluded spans, cadence includes excluded-only active days/all completions, and summary exclusion copy is unchanged. Pin 6/7/8 rolling cases, gapped/excluded-only/future-eighth data that distinguish active-day trailing behavior from calendar/zero-filled/centered alternatives, draw order after bars, and exact Panel C peak/midpoint/zero labels plus scaled midpoint gridline.

Live browser proof must return `location.href`, measure finite Panel B geometry, light/dark rolling-line contrast against `document.body`, tile/chart separation, Panel C label separation, window updates, 320/768/1280 layout, keyboard/semantic `<dl>` state, and zero console errors. Retain screenshots and visually inspect the generated repository board.

Run focused and full module/canonical tests proportionately, syntax/vet/diff checks, commit, and write the hand-back with manifest, RED/GREEN, rendered measurements/artifacts, integration seams, `## Decisions`, `## Discovered Tasks`, and risks.

The hand-back pattern makes fan-out failures survivable, not prevented.
