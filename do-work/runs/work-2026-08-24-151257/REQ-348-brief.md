# REQ-348 builder brief

Build REQ-348, “Group the Timeline's rows by user request,” on Route C with TDD.

## Work location

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-348-group-the-timelines-rows-by-user-request`
- Branch: `worktree-agent-REQ-348-group-the-timelines-rows-by-user-request`
- Hand-back (the only permitted main-tree write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-151257/REQ-348-handback.md`

Never read or write any path under `do-work/` in the worktree or main checkout except the absolute hand-back path above. Never write the main working tree. Commit the implementation and tests on the named branch.

## Governing context

Read in the worktree before coding:

- `_dev/primes/prime-kanban-board.md`
- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- `skills/do-work/crew-members/frontend.md`
- `skills/do-work/crew-members/testing.md`

The Timeline is the chosen place for UR grouping; do not add a fourth Durations panel. The board must make each window-scoped UR's elapsed-versus-worked gap and REQ count readable while retaining Timeline virtualization, focus, panning, and window behavior.

## Required behavior

- Group window-listed Timeline REQ rows beneath their `userRequestId`; rows without one remain under an explicit “No UR recorded” group placed last.
- Join UR identity client-side from `boardData.requests`. Do not add UR identity to `TimelineRow` unless the existing join genuinely cannot work.
- Join completed work to `boardData.durations.samples` and consume its existing `excludedReason` verdict. Do not duplicate the four-hour/negative-span rule.
- Group only after `timelineRowsInWindow`; preserve newest-first member order and first-seen group order. A header's metrics cover only listed members.
- A finished group spans earliest recorded REQ claim to latest resolved completion. An open group ends at the payload's frozen `now` and is visibly running. If no recorded claim exists, elapsed is unavailable; never silently substitute `created_at`.
- Show explicitly labelled elapsed, summed accepted work, listed REQ count, and exclusion/fallback detail in the existing visual/readout idiom.
- Flatten fixed-height headers and members for SVG virtualization. Headers are not Tab stops. Map the existing REQ-only roving focus index through display indices so exactly one REQ remains `tabindex="0"`; Up/Down skip headers, Left/Right panning and Enter activation keep working, and focus restoration survives virtualization.
- Group the accessible table too. Explain window scope, first-claim endpoint semantics, frozen-now running groups, unavailable no-claim spans, Durations' read-time work rule, and ungrouped-last placement.
- Add theme-aware styling and inspect both light and dark generated boards.

## Declared scope

- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js`
- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`

If an integration seam outside that scope is unavoidable, do not edit it. Record its exact path and line in the hand-back.

## TDD and proof

Add RED coverage first for stable grouping, unknown-last placement, count, closed/open endpoints, no-claim handling, Durations exclusions, and grouping after window membership. Extend the live generated-board probe to prove header/member order and counts, bounded viewport nodes, rebuilt groups after zoom/window changes, and the REQ-only keyboard contract. Then implement, run focused tests plus the relevant canonical lanes, generate a board, and inspect light/dark measurements and contrast.

## Hand-back format

Write the absolute hand-back file before finishing with:

- branch and commit hash;
- exact file manifest and test commands/results;
- RED and GREEN evidence;
- visual evidence/measurements and artifact paths;
- integration seams (or `none`);
- `## Decisions`;
- `## Discovered Tasks`;
- risks or unresolved issues.

The hand-back pattern makes fan-out failures survivable, not prevented.
