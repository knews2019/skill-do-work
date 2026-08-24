# REQ-353 builder brief

Build REQ-353, “Hide the dead filter knobs while Durations is on screen,” on Route A.

## Work location

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-353-hide-dead-filter-knobs-on-durations`
- Branch: `worktree-agent-REQ-353-hide-dead-filter-knobs-on-durations`
- Hand-back (the only permitted main-tree write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-164337/REQ-353-handback.md`

Never read or write any path under `do-work/` in the worktree or main checkout except the absolute hand-back path above. Never write the main working tree. Commit implementation on the named branch.

## Governing context

Read `_dev/primes/prime-kanban-board.md` plus the general, coding-guardrails, frontend, testing, and communication-style crew rules before editing.

## Required behavior and scope

- Only modify `skills/do-work-board/tools/queue-kanban/web/board-controls.js`.
- In the existing `applyView` visibility path, hide the text search/filter, domain select, and status select whenever `viewState.view === "durations"`.
- Restore all three when switching to every other view where they apply.
- Reuse the same visibility mechanism already used for lens/recent-window/dead controls; do not add a parallel state path.
- Do not touch `web/board-filters.js` and do not make Durations respond to shared filters.
- Generate a board and inspect Durations → another view → Durations, including keyboard-reachable state and no console errors.

This request explicitly constrains implementation to one file. Do not expand scope to add a test file; use existing behavior/canonical tests plus direct generated-page proof.

## Proof and hand-back

Run Node syntax, relevant existing focused tests, the module suite proportionately, and `git diff --check`. Record a pre-fix generated-page observation showing the three controls visible on Durations and post-fix observation showing all hidden then restored. Include browser/build, `location.href`, exact element ids/hidden states, and any screenshot path.

Write the hand-back with branch/commit, exact file manifest, commands/results, RED/GREEN observation, integration seams (`none` if none), `## Decisions`, `## Discovered Tasks`, and residual risk.

The hand-back pattern makes fan-out failures survivable, not prevented.
