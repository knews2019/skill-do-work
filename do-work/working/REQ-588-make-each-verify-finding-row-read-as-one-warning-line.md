---
id: REQ-588
title: 'Addendum: make each verify-finding row read as one warning line, not a paragraph'
status: claimed
created_at: 2026-09-05T14:45:41Z
user_request: UR-124
addendum_to: REQ-579
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-579]
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
status_changed_at: 2026-09-05T17:00:29Z
claimed_at: 2026-09-05T17:02:32Z
---

# Addendum: Make Each Verify-Finding Row Read as One Warning Line, Not a Paragraph

## What

REQ-579 (render verify findings and skipped probes as compact rows in one list) replaced the finding cards with rows, but each row now renders as one long grey paragraph: the detail sentence, an arrow and the remedy sentence flow together and wrap under the chip, with one orphaned word on the second line. The subject heading, the chip and the row text also sit at three unrelated sizes and weights. Change the row layout so each finding reads as one warning line with its remedy visibly separated, using the layout the user picks from the mock-up report in `ai-reports/2026-09-05_1445_REQ-588-verify-findings-row-mockups/`.

## Prior Implementation

REQ-579 shipped in commit b169396e (builder branch commit 1dc13ef7, board 0.236.20). The producer (`verify.go`) gained a `Subject` on `VerifyFinding`, mirrored as `subject` in the board payload by `generate.go`. The client (`renderVerifyFindingsStrip` in `web/board-cards.js`) groups findings by exact subject, prints a `.board-findings-subject` heading per group, then one `.board-findings-row` per finding: a `.board-findings-chip` span followed by a `.board-findings-text` span holding the detail, the optional "cleanup can fix" tag and the remedy prefixed with an arrow, all inline. Skipped probes are rows with a "not checked" chip and the muted class. `web/board.css` lays the row out as `display: flex; align-items: baseline` with the text span wrapping beside the chip. The Node behavior lane (`javascript_behavior_c_test.go`, the `board-findings` cases around lines 2600 and 3030) asserts the list shape, the subject heading, and which rows carry the muted class.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The user's words: "verify findings styling is broken". The strip exists to be scanned as a list of warnings; a row that wraps as a paragraph has to be read, not scanned, and the remedy (what to do) is the part that gets buried at the end of the wrap.

## Requirements

- **D1, remedy separated from detail.** The remedy never continues the detail sentence inline. It sits on its own line under the detail (M1, M2) or is revealed on demand (M3), per the answer to Q1 below. The arrow prefix goes away when the remedy has its own line; a line does not need a pointer to say it follows.
- **D2, chip and text on one grid.** The row is a two-column grid (chip, text) so wrapped text stays in the text column and never lands under the chip. In M2 the grid has a third leading column for the subject, and every chip in the strip aligns in one column.
- **D3, one type scale.** Subject heading, chip and row text use one scale: subject in the board's mono face at the row text size and semibold ink-base (it is an identifier, like a REQ id on a card); chip one step smaller, uppercase, letter-spaced; row text at the size REQ-579 set; remedy in ink-soft. No fourth size.
- **D4, the producer stops repeating the subject.** When `Subject` is set, the detail sentence no longer starts with that same name (today: "worktree-agent-REQ-573-activity-drawer exists — .git/...", under a heading that already says the name). Change the wording in `verify.go` for the subject-bearing findings and update `verify_test.go`. Grouping stays by the payload field, never by parsing text.
- **D5, everything else from REQ-579 stands.** One list, weight only from `fixable` and "not checked", grouping by subject, hide rules, no card classes. The Node behavior lane's existing assertions keep passing; the assertion for the chosen layout is added beside them (see Red-Green Proof).
- Board changes follow `_dev/primes/prime-kanban-board.md` (versioning, parser lock-step, build outputs).

## Open Questions

- [x] Which row layout from the mock-up report should ship? → M1, remedy on its own line under the detail, chip and text on a two-column grid, subject headings kept. Picked from the mock-up report ai-reports/2026-09-05_1445_REQ-588-verify-findings-row-mockups/ over M2 (subject as a row-label column) and M3 (remedy behind a toggle); both stay out of scope.
  Recommended: M1 (remedy on its own line under the detail, chip and text on a two-column grid) — the smallest change that fixes all three defects, and it keeps REQ-579's subject heading and row order untouched.
  Also: M2 (subject as a row label in a leading column, chips aligned in one column, remedy under the detail) — the most scannable, costs a small renderer change; M3 (one-line rows, remedy revealed by a per-row toggle) — the smallest strip, costs a click to read what to do.

## Red-Green Proof
**RED prompt/case:** Serve the board with two verify findings under two subjects (the state in Screenshot 1) and look at the strip, or render the same payload in the Node behavior lane and inspect one `.board-findings-row`.
**Why RED now:** The row's text span holds detail, arrow and remedy as inline siblings; on a 2000 px wide board the second finding wraps to a second line with one word on it, and the remedy has no separation from the detail.
**GREEN when:** For the chosen layout, the Node lane asserts the row's structure (M1/M2: the remedy is a block-level element after the detail with no arrow text; M2: the subject is a cell of the row, not a heading above it; M3: the remedy is inside a toggle element that is closed by default) and the board screenshot shows each finding as a chip, one detail line and one separated remedy line, with subject, chip and text at the scale in D3.
**Validation:** User confirmed the defect from Screenshot 1 ("verify findings styling is broken") and asked for mock-ups to choose the layout; the choice is recorded as the answer to Q1.

## Assets

- Screenshot 1 (not saved: the attachment cache had already expired when capture ran): the board at 127.0.0.1:8090, Board view, light theme, 14:25 UTC on 2026-09-05. Strip header "VERIFY FINDINGS 2 findings queue and process problems queue-kanban verify detects — each names what to do about it". Two subject headings (worktree-agent-REQ-573-activity-drawer, worktree-agent-REQ-582-arrow-citations), each with one row: an uppercase chip (UNMERGED-WORKTREE-LEFTOVER, WORKTREE-PRESENT-RUN-IN-FLIGHT), then detail, arrow and remedy as one grey paragraph wrapping to a second line with one orphaned word. The mock-up report's M0 page rebuilds this state from the same payload text.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (6224 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4959 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

*Source: "verify findings styling is broken, which req should fix it?" / "capture it and fix it" / "also it's make do-work ai-report with mock-ups so I have options to choose from"*


## Answer Notes

- 2026-09-05 - [ ] Which row layout from the mock-up report should ship?: M1, remedy on its own line under the detail, chip and text on a two-column grid, subject headings kept. Picked from the mock-up report ai-reports/2026-09-05_1445_REQ-588-verify-findings-row-mockups/ over M2 (subject as a row-label column) and M3 (remedy behind a toggle); both stay out of scope.
> ```
> M1 remedy under detail (Recommended)
> ```
