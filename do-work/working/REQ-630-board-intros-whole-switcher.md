---
id: REQ-630
title: 'Board intros describe the board by what its switcher covers'
status: claimed
created_at: 2026-10-02T21:51:57Z
user_request: UR-134
domain: general
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: ["REQ-631"]
batch: ur-133-leftovers
write_set: ["skills/do-work-board/docs/board-guide.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
claimed_at: 2026-10-02T21:52:48Z
---
# Board Intros Describe the Board by What Its Switcher Covers

## What
Rewrite the opening line of `skills/do-work-board/docs/board-guide.md` and of `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` so neither names only some of the board's pages.

## Why
The REQ-628 review (finding M1) found both intros still read like "a Kanban board ... plus a queue activity calendar and a testing track", three of six pages. The lesson "describe the board by what its switcher covers, not by naming some of its tabs" (REQ-232, `_dev/primes/lessons-kanban-board.md`) says this drifts every time a page is added.

## Detailed Requirements
- Each intro describes the board by what its switcher covers (for example "the Board view plus the other pages in its switcher"), or names all six pages; prefer the wording that does not need editing when a page is added.
- Keep the rest of each opening paragraph (read-only toward the pipeline, the testing record) unchanged.
- Grep shipped files under `skills/do-work-board/` for other partial page lists in intros and report any in Discovered Tasks.

## Constraints
Docs only. Shipped files change, so this is a release.

## Dependencies
None. Independent of REQ-631.

## Red-Green Proof
**RED prompt/case:** Line 3 of `board-guide.md` and of `prime-do-kanban.md` names a calendar and a testing track as the board's extras, leaving out Activity, Timeline and Durations.
**Why RED now:** Both intros predate the six-page switcher.
**GREEN when:** Neither intro names a partial set of pages.
**Validation:** Inferred during capture

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-134/input.md` for complete verbatim input.

*Source: ok, do board intros and the recover message*
