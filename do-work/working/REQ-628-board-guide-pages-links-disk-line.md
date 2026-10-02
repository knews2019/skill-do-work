---
id: REQ-628
title: 'Board guide names all six pages, their links, and the disk line'
status: claimed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-02T20:48:08Z
created_at: 2026-10-02T19:49:03Z
user_request: UR-133
domain: general
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: ["REQ-629"]
batch: ur-132-follow-ups
write_set: ["skills/do-work-board/docs/board-guide.md", "skills/do-work-board/actions/board.md"]
claimed_at: 2026-10-02T20:47:38Z
---
# Board Guide Names All Six Pages, Their Links, and the Disk Line

## What
Update the board's user-facing docs so they match what shipped in 0.305.61 and 0.305.62: the page switcher has six pages, every page and Board lens has its own link, and the Testing toolbar shows the repo root's free disk space.

## Why
The REQ-626 review (Restatement Sweep, finding F3) found `skills/do-work-board/docs/board-guide.md` line 25 still calls the switcher "**Board / Calendar / Testing**", three of the six pages, and the guide never mentions the page links that were the user's stated purpose. The REQ-627 review (finding M1) found the guide's Testing section and `skills/do-work-board/actions/board.md` Step 6 describe the Testing toolbar without the new disk line. The user asked to capture this.

## Detailed Requirements
- `board-guide.md`: the toolbar paragraph names the switcher by what it covers (all six pages: Board, Activity, Calendar, Timeline, Durations, Testing), not by a partial list.
- `board-guide.md`: one short passage says each page and Board lens has its own link and lists the fragments (`#board`, `#board/by-ur`, `#board/urs-only`, `#activity`, `#calendar`, `#timeline`, `#durations`, `#testing`), that clicking a page updates the address bar without adding Back steps, that filters are not part of the link, and that links work on a static snapshot opened from disk.
- `board-guide.md` Testing view section: one sentence on the disk line (free of total, amber below 10 GiB, red below 3 GiB, "(at generation)" on a static snapshot, "not measured" when the platform cannot measure; the live board shows the value from page load until reload).
- `board.md` Step 6: one clause naming the disk line in the Testing toolbar, only if that step describes the toolbar's contents.
- Do not restate behaviour the code does not have; check each sentence against `web/board-controls.js` and `web/board-testing.js`.

## Constraints
Docs only. Shipped files change, so this is a release. Follow the lesson "describe the board by what its switcher covers, not by naming some of its tabs" (`_dev/primes/lessons-kanban-board.md`).

## Dependencies
None. Independent of REQ-629 (handoff resume must not reset claims).

## Builder Guidance
High certainty; wording is the builder's. Keep each addition short; this guide is read by people, not agents.

## Red-Green Proof
**RED prompt/case:** Read `skills/do-work-board/docs/board-guide.md`: the switcher is described as "Board / Calendar / Testing", nothing says a page can be opened by link, and the Testing section does not mention the disk line.
**Why RED now:** The docs predate 0.305.61 and 0.305.62.
**GREEN when:** The guide names all six pages, documents the page and lens links with their fragments, and the Testing section describes the disk line; `grep -n "Board / Calendar / Testing" skills/do-work-board/docs/board-guide.md` returns nothing.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the 2000 budget; `slugged: partial`). Matched: views and static output. The bullet this REQ cites (REQ-232, describe the board by what its switcher covers) carries no family marker, so no narrower entry exists; the REQ's Constraints name that bullet, and the builder reads it directly.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-133/input.md` for complete verbatim input.

*Source: capture both as REQs (finding F3 of the REQ-626 review and finding M1 of the REQ-627 review)*

---

## Triage

**Route: A** - Simple

**Reasoning:** Docs-only change that names its two files and lists every sentence to add; the only discovery is checking each sentence against `web/board-controls.js` and `web/board-testing.js`, which the builder does inline.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
