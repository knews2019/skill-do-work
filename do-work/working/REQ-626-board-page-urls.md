---
id: REQ-626
title: 'Give every board page and lens its own URL'
status: claimed
estimate:
  p50_active_minutes: 20
  confidence: medium
  basis:
  - Route B
  - 3-file write set
  - 7 acceptance criteria
  calculated_at: 2026-10-02T18:04:23Z
created_at: 2026-10-02T18:02:37Z
user_request: UR-132
domain: frontend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: ["REQ-627"]
batch: board-links-and-disk-space
integration_at: 2026-10-02T18:18:36Z
builder_handback_at: 2026-10-02T18:18:20Z
dispatch_at: 2026-10-02T18:09:36Z
claimed_at: 2026-10-02T18:03:38Z
status_changed_at: 2026-10-02T18:36:12Z
---

# Give Every Board Page and Lens Its Own URL

## What
Each board page (Board, Activity, Calendar, Timeline, Durations, Testing) and each Board lens (flat columns, By UR, URs only) gets its own URL that opens that page directly, and switching pages or lenses updates the address bar so the link can be copied and shared.

## Why
The board has many pages and the user wants to point people at a specific one by URL. Today `viewState.view` lives only in JavaScript (web/board-controls.js), so every link opens the Board page and the address bar never changes.

## Detailed Requirements
- A URL fragment names the page, and optionally the Board lens, for example `#board`, `#board/by-ur`, `#board/urs-only`, `#activity`, `#calendar`, `#timeline`, `#durations`, `#testing`. The exact spelling is the builder's, but it must be stable, lowercase, and readable.
- Opening the board with a fragment shows that page (and lens) directly, including the lazy first render the view switch does today.
- Clicking a page or lens control updates the fragment without reloading (`history.replaceState` or `location.hash`), and the browser Back button is not broken by it.
- An unknown or empty fragment falls back to the current default (Board, flat lens) without an error.
- The static snapshot (`do-work-board static`, a file opened from disk) behaves the same, since fragments work on `file://` URLs.
- Filters (search text, domain, status, recently-done window) stay out of the URL. User decision at capture.
- Any existing stored-view preference (localStorage, if present) yields to an explicit fragment.

## Constraints
No new dependency; plain JavaScript in the existing files. Keep the change in board-controls.js and template.html, plus tests. Shipped files change, so this is a release.

## Dependencies
None. Independent of REQ-627 (Show free disk space on the Testing page).

## Builder Guidance
High certainty on the behavior; builder latitude on fragment spelling and on whether lens is a path segment or a query-like suffix. A JavaScript behaviour test in the existing `javascript_behavior_*_test.go` harness is the expected proof; keep the test file under the 30 s per-file budget.

## Red-Green Proof
**RED prompt/case:** Open `http://127.0.0.1:8090/#timeline` (or any page URL). The Board page shows, not Timeline; click Timeline and the address bar still reads the bare URL.
**Why RED now:** The selected view is JavaScript state only; nothing reads or writes the URL.
**GREEN when:** `#timeline` opens the Timeline page directly; each of the six pages and each Board lens has a URL that opens it; clicking a page or lens control updates the address bar to that URL; an unknown fragment falls back to the Board page.
**Validation:** User confirmed

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (7095 tokens, over the 2000 budget; `slugged: partial`). Matched: changing queue-kanban UI and browser behaviour. Read it anyway if time allows.
- `_dev/primes/lessons-kanban-board.md` (5912 tokens, over budget; `slugged: partial`). Matched: views and view scroll.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.) Builder: five small functions in the controls file (fragment parse and format, guarded hash read, selection without render, guarded replaceState write); read once at wiring, write after clicks, hashchange applies known names only.
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.) Builder: declared files plus the granted one-line URL-check substitution in four more browser probe test files; board.js and template.html untouched.
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.) Builder: 10 files, 473 insertions, 19 deletions; gofmt empty, go vet exit 0; every substituted URL check reviewed line by line.

## Full Context
See `do-work/user-requests/UR-132/input.md` for complete verbatim input.

*Source: there are a lot of pages in that kanban (BTW: capture-request, I need a link for each page so I can point to them by URL)*

---

