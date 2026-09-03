# REQ-561 builder brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-561-add-a-three-value-priority-field-the-selector-orders-by-and-the-board-shows`

Branch: `worktree-agent-REQ-561-add-a-three-value-priority-field-the-selector-orders-by-and-the-board-shows`

Main-tree hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-03-210100/REQ-561-handback.md`

Do not read or write the worktree's tracked `do-work/` snapshot. Do not write any main-tree path except the hand-back above. Commit project changes on the builder branch.

## Request

Add optional REQ frontmatter `priority` with closed values `now`, `next`, `later`. Absent and invalid values resolve to `next`, with invalid values warning under the Schema Read Contract. `nextselection` orders ordinary ready work by this value while gate-repair and deferred-parent classes remain above it, dependency readiness remains a hard gate, and queue order breaks ties. Project priority on selected and excluded typed records.

The board parses priority in lock-step, sorts the pending column by it inside the ready and waiting groups, and shows a small `now` or `later` tag but no default `next` tag. Capture documents and assesses the field only when the user's words rank work. The landing also stamps the current pending queue from the latest velocity report: build-now set `now`, deferred set `later`, all others untouched. REQ-530 is superseded and must be cancelled against this landing hash by the orchestrator after successful integration.

Required tests: table-driven schema/selector cases for absent, all values, invalid, and a `later` prerequisite required by a `now` dependent; board model sort; static and live client behavior. RED evidence must first show the selector ignores priority and the board has no tag/order. Generate and visually inspect the board in light and dark, including DOM evidence naming the exact URL.

Prime context: `_dev/primes/prime-action-files.md`, `_dev/primes/prime-kanban-board.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`, plus every lesson satellite those primes require. Follow general, coding-guardrails, backend, testing, and communication-style crews.

The orchestrator owns all main-tree queue stamping, the REQ-530 cancellation, version/changelog, and integration seams. The builder may change source/actions/tests on this branch and must hand queue-stamp paths and any seams back explicitly.

## Hand-back format

Write the branch name, commit hash, full file manifest, RED then GREEN commands/results, browser evidence, lesson-read evidence, integration seams, a `## Decisions` section if any, and a `## Discovered Tasks` section if any. Do not edit the main-tree REQ.
