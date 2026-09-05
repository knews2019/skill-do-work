---
id: REQ-589
title: 'Addendum: render the verify findings strip as the M4 slim band, one line closed and one row per finding open'
status: claimed
created_at: 2026-09-05T18:19:11Z
user_request: UR-125
addendum_to: REQ-588
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-588, REQ-579, REQ-578]
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-05T18:20:44Z
  basis:
  - Route A
  - 4-file write set
  - 6 acceptance criteria
route: A
dispatch_at: 2026-09-05T18:21:42Z
builder_handback_at: 2026-09-05T18:44:37Z
claimed_at: 2026-09-05T18:20:43Z
---

# Addendum: Render the Verify Findings Strip as the M4 Slim Band, One Line Closed and One Row per Finding Open

## What

REQ-588 (the M1 rows, release 0.303.2) made each finding a chip, a detail line and a remedy line; the user's verdict was that it is not visually nice and huge. Replace the strip's rendering with mock-up M4 from `ai-reports/2026-09-05_1800_REQ-588-verify-findings-slim-band-gallery/`: a slim band that is one line when closed and one row per finding when open, with each remedy behind its row's chevron. The mock-up page `mockups/m4-closed.html` (and `m4-open.html`, `m4-open-remedy.html`) is the specification: its `<style>` block is the CSS to ship and its markup is the structure to render.

## Prior Implementation

REQ-588 shipped in commit 707ffb6c (merge ab251f24). `renderVerifyFindingsStrip` in `web/board-cards.js` groups findings by subject, prints a `.board-findings-subject` heading per group, then one `.board-findings-row` per finding (chip span, then a text span holding detail, optional "cleanup can fix" tag and remedy as blocks); skipped probes are muted rows with a "not checked" chip. `web/template.html` holds the strip section: a header with the title, `#board-findings-count` and a hint, then `#board-findings-rows` wrapping two `display: contents` hosts `#board-findings-cards` and `#board-findings-skipped-list`, which `applyView` in `web/board-controls.js` (REQ-578) reads to decide whether the strip has content on the Activity view. `web/board.css` holds the `.board-findings-*` rules around lines 607–700. The Node lane pins the list shape (`TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList`), the M1 row (`TestJavaScriptBehaviorVerifyFindingRemedyIsItsOwnLineAfterTheDetail`) and the Activity hide rule (`TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip`).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The user's words: "neither of these is visually nice and they are huge". A remedy is read after deciding to act, so it belongs behind a click; what stays visible must be short enough to scan.

## Requirements

- **D1, the band.** The strip is a slim band: no card border box, a 3 px amber left edge (`--accent-pending`), 8 px radius, `--surface-2` background, 5 px vertical padding. The header line is a small warning glyph plus a quiet uppercase VERIFY label, then the counts ("3 findings · 1 probe not checked"). The hint sentence is no longer rendered.
- **D2, closed state.** The whole strip is a `<details>` whose summary is the header line: label, counts, then every finding's subject in the mono face preceded by its weight dot (amber for an ordinary finding, green for `fixable`, grey for a skipped probe), the skipped probes summarised as "N probe(s) not checked", and a Show button at the right. Closed height is one line (about 34 px).
- **D3, open state.** Opening swaps the subject list for one row per finding, in the producer's order grouped by subject as today: dot, subject (mono, semibold), category, the detail clipped with an ellipsis at the line's end, the "cleanup can fix" tag as a small green pill when `fixable`, and a chevron. Skipped probes are rows too, with a grey dot and the category "not checked". Each row is a `<details>`; opening it shows the remedy under the row in an inset block labelled "What to do:" with a 2 px amber left rule, and lets the detail wrap. Rows have a hover background and a focus-visible ring. The Show button reads Hide with its chevron flipped while open.
- **D4, remembered.** The strip's open/closed state persists per browser in `localStorage` under one key, best-effort like the detail-panel width in `web/board-detail.js`; the default is closed. Row (remedy) state is not persisted.
- **D5, category as words.** The category token renders lowercase with hyphens shown as spaces (a mechanical transform, no list in the client), in faint ink. No uppercase chip remains.
- **D6, nothing else moves.** Producer, payload, grouping by subject (exact match on the payload field), the two host ids and REQ-578's Activity hide rule, and hide-when-empty stay as they are. No `.board-request*` classes.
- Board changes follow `_dev/primes/prime-kanban-board.md`; embedded web assets reach consumers on the next build.

## Red-Green Proof
**RED prompt/case:** Render the board with three findings under three subjects and one skipped probe (the gallery's data) and look at the strip; in the Node lane, inspect `#board-findings`.
**Why RED now:** The strip has no `details` element, no Show control, no stored state; every remedy is always visible; the header carries the hint sentence and the category is an uppercase chip.
**GREEN when:** The Node lane asserts: the strip's rows sit inside a `details` element that is closed by default and whose summary names every subject; each finding row is a `details` whose summary holds dot, subject, category and detail, and whose content holds the remedy; the category text is the token lowercased with spaces; the skipped probe is a row with the "not checked" category; opening the strip and reloading with the stored key set renders it open; REQ-578's Activity hide test and the hide-when-empty test still pass. The board screenshot matches `mockups/m4-closed.html`, `m4-open.html` and `m4-open-remedy.html` in both themes.
**Validation:** User picked M4 from the gallery ("ok, M4 is good") after two rounds of mock-ups; the gallery page is the approved design.

## Assets

- `ai-reports/2026-09-05_1800_REQ-588-verify-findings-slim-band-gallery/` (committed at 2a795ba0): the approved M4 pages and their captures. Screenshot 3 (not saved; the attachment cache expires) showed that gallery's M4 section: the live frame closed, then State 1 closed and State 2 open in light and dark.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (6477 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4959 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

*Source: "neither of these is visually nice and they are huge, please provide better mocks, it's fine to be colapsible as well" / "I want all of the options in the mockup, don't make me imagine what would be, also make it beautiful and professional" / "ok, M4 is good"*

---

## Triage

**Route: A** - Simple

**Reasoning:** The approved mock-up pages are the specification (markup and CSS), the request names the four files, and the renderer change is a restructuring of one function plus its template section. Substantive in size, but nothing to explore or decide. `effort_estimate: effort-substantive`.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
