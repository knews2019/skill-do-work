# REQ-349 builder brief

Build REQ-349, “Fix panel A's scale and density,” on Route B with TDD.

## Work location

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-349-fix-panel-a-scale-and-density`
- Branch: `worktree-agent-REQ-349-fix-panel-a-scale-and-density`
- Hand-back (the only permitted main-tree write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-151257/REQ-349-handback.md`

Never read or write any path under `do-work/` in the worktree or main checkout except the absolute hand-back path above. Never write the main working tree. Commit implementation and tests on the named branch.

## Governing context

Read in the worktree before coding:

- `_dev/primes/prime-kanban-board.md`
- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- `skills/do-work/crew-members/frontend.md`
- `skills/do-work/crew-members/testing.md`

Panel A currently packs most marks into the bottom of a linear 0–60 scale and stacks busy-day marks at one x. The required dense-board target is at least 700 archived REQs, roughly 8 SVG units per day.

## Required behavior

- Use a square-root 0–60 minute y scale with ticks exactly at 0, 5, 15, 30, 45, and 60.
- Deterministically jitter each mark within its own UTC day slot. The same payload renders identically; no mark may cross a day boundary. Use stable identity/rank and keep useful spread within an 8-unit slot.
- Use the jittered x for both circle geometry and `markIndex.x`, so hover at a mark still names that exact REQ.
- Lower ordinary mark opacity without erasing route/UR/domain colour, critical reversed-stamp red, or unknown styling.
- From samples whose existing payload `excludedReason` is empty, compute sorted daily p25/median/p75 using a stated deterministic quantile rule. Draw a p25–p75 ribbon and median line behind the circles at day-centre x values.
- Keep the 60-minute ceiling, overflow lane, reversed band, label/leader behavior, and read-time exclusion rule unchanged. Apply jitter consistently to overflow/reversed marks and any dependent geometry.
- Generate and inspect dense light and dark boards; the ribbon must remain subordinate but readable and all marks/overlays remain bounded.

## Declared scope

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/durations_test.go`
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`

`generate_test.go::TestJavaScriptBehaviorDurationsDayBucketsStayInsideThePlot` currently requires exact unjittered x coordinates. Update it honestly to pin own-day bounds, repeat-render determinism, and non-degenerate busy-day spread. If styling can use SVG presentation attributes, keep CSS outside scope. If any integration seam outside scope is unavoidable, do not edit it; record its exact path and line in the hand-back.

## TDD and proof

Add RED evidence before implementation. Preserve Panel B/C day-bucket checks while changing only Panel A's invalid exact-x assertion. Pin exact tick labels and sqrt y positions, lower opacity, overlay draw order, one valid ribbon/median path, the unchanged ceiling/overflow/reversed rules, and deterministic within-day geometry. Add a real-browser probe with 700+ samples across about 47 days to measure actual coordinates/day bounds, busy-day spread, deterministic rerender, hover identity at jittered circle centres, finite bounded ribbon/line geometry, and the generated `location.href`. Run focused tests plus relevant canonical lanes and inspect light/dark output.

## Hand-back format

Write the absolute hand-back file before finishing with:

- branch and commit hash;
- exact file manifest and test commands/results;
- RED and GREEN evidence;
- rendered measurements, browser/build, `location.href`, and artifact paths;
- integration seams (or `none`);
- `## Decisions`;
- `## Discovered Tasks`;
- risks or unresolved issues.

The hand-back pattern makes fan-out failures survivable, not prevented.
