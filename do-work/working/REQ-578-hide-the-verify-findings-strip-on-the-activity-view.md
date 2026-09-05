---
id: REQ-578
title: 'Hide the verify-findings strip on the Activity view'
status: claimed
created_at: 2026-09-04T23:58:59Z
user_request: UR-117
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
related: [REQ-573]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
claimed_at: 2026-09-05T12:00:56Z
---

# Hide the Verify-Findings Strip on the Activity View

## What

The Verify Findings strip (`#board-findings`, added by REQ-285) sits outside the view panels so it stays visible on every view. On the Activity view it pushes the transitions table down and is not what that view is for. Hide the strip while the Activity view is active and show it again when the reader switches to any other view. The strip's content and its behavior on the other views do not change.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- With the Activity view selected, `#board-findings` is hidden (the `hidden` attribute, matching how the strip already hides itself when there are no findings) even when findings exist.
- Switching to Board, Calendar, Durations, Timeline or Testing shows the strip again exactly as today; the "probe(s) could not run" disclosure under it follows the strip.
- Only the Verify Findings strip is affected. The completion-anomalies strip above it is not part of this request.
- The rule lives in the view-switching code (`board-controls.js`), not in the Activity renderer, so a re-render of the Activity table never touches the strip. Update the template comment that says the strip stays visible in every view.

## Red-Green Proof
**RED prompt/case:** In the Node behavior lane, render with two verify findings in `boardData`, switch the view to `activity`, and read `document.getElementById("board-findings").hidden`.
**Why RED now:** The strip is outside the view panels by design and nothing in the view switch touches it, so it stays visible on the Activity view (screenshot 3).
**GREEN when:** `hidden` is true while the Activity view is active and false again after switching back to the Board view with the same findings.
**Validation:** User request from the live board; proof inferred during capture.

## Builder Guidance

The user is certain about the outcome. Keep it to the view switch plus one test; do not restructure the strips.

## Assets

- `do-work/user-requests/UR-117/assets/REQ-578-screenshot-3-activity-view-with-verify-strip.png`: the Activity view at 24h with "175 transitions across 38 REQs in the last 24 hours" and rows for REQ-576, 575, 574, 572 (four rows: work merged, builder handed back, builder dispatched, captured), 506, 570, 573, 505 and others; above the table a Verify Findings strip with two cards (WORKTREE-MERGE-STATE-UNDETERMINED for a REQ-506 worktree, WORKTREE-PRESENT-RUN-IN-FLIGHT for the REQ-570 worktree) and a "1 probe(s) could not run" disclosure.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

*Source: "remove verify finding from this view"*
