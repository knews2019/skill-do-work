---
id: REQ-482
title: 'Stack verify-findings cards full width so they stop reading as REQ cards'
status: pending
priority: later
created_at: 2026-09-01T11:55:00Z
user_request: UR-090
domain: ui-design
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
assigned_to: 'after-drain'
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T12:51:51Z
  basis:
    - trivial short-circuit
write_set: [skills/do-work-board/tools/queue-kanban/web/board.css, skills/do-work-board/tools/queue-kanban/web/template.html, skills/do-work-board/tools/queue-kanban/web/board-cards.js, skills/do-work-board/tools/queue-kanban/*_test.go]
---

# Stack Verify-Findings Cards Full Width So They Stop Reading as REQ Cards

## What
On the board, the Verify Findings strip lays its finding cards out side by side in narrow columns, so each finding looks like a REQ card even though it is not one. Render each finding card full width and stack the findings vertically, and make the strip visually distinct from REQ cards.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)
The current side-by-side boxes look like REQ cards, which they are not — the resemblance misreads.

## Context
- Capture clarification: the user confirmed the screenshot is the queue-kanban board web UI.
- Capture located the surface: `#board-findings-cards` in `skills/do-work-board/tools/queue-kanban/web/template.html` (line ~184) reuses the anomalies strip's `board-anomalies-cards` class, whose grid is `repeat(auto-fill, minmax(260px, 1fr))` in `skills/do-work-board/tools/queue-kanban/web/board.css` (line ~568). Finding cards are built by `renderVerifyFindingsStrip` in `web/board-cards.js`.

## Constraints
- The anomalies strip shares the `board-anomalies-cards` layout class and hosts real REQ cards (`makeRequestCard`) — those should keep their REQ-card look and multi-column layout. Give the findings strip its own layout rather than restyling the shared class.
- Board changes follow `_dev/primes/prime-kanban-board.md` (versioning, parser lock-step, build outputs).

## Builder Guidance
Firm: each finding card full width, findings stacked vertically. Latitude granted: the user asked to "basically improve the UI/UX" — polish the strip's overall look while making it read as findings, not REQ cards, is in scope beyond the minimal layout swap.

## Red-Green Proof
**RED prompt/case:** Open the board (Board view) with two or more verify findings present at desktop width — e.g. the state in the captured screenshot (CLAIM-NEEDS-ATTENTION and MERGED-WORKTREE-LEFTOVER). The finding cards render side by side in ~260px-min columns, visually similar to REQ cards. Runnable form: a browser/layout test in the existing `skills/do-work-board/tools/queue-kanban/*_browser_test.go` pattern asserting that with ≥2 findings, each `.board-finding`'s rendered width equals the strip's content width (single column) — fails today.
**Why RED now:** `#board-findings-cards` reuses the anomalies strip's multi-column grid class, so findings wrap into narrow card-shaped columns.
**GREEN when:** Each verify-finding card spans the full width of the strip and findings stack vertically; the anomalies strip's REQ cards keep their current multi-column layout; the layout test above passes.
**Validation:** User confirmed — full-width and stacked-vertically is stated verbatim in the input, and the screenshot shows the RED state.

## Assets
`do-work/user-requests/UR-090/assets/REQ-482-screenshot-1-verify-findings-side-by-side.png` — queue-kanban board, light theme, Board view, generated 2026-09-01 11:42 UTC. Top strip "VERIFY FINDINGS 2" shows two white finding cards side by side: left "CLAIM-NEEDS-ATTENTION" (REQ-418 claimed 3h0m, reported not judged dead, with remedy text), right "MERGED-WORKTREE-LEFTOVER" (cleanup can fix, worktree-agent-REQ-418-toolbox-migration exists). Below: PENDING 36 (REQ-437, REQ-438 cards), CLAIMED 1 (REQ-418), NEEDS INPUT · BLOCKED 0. The finding cards' size and shape closely resemble the REQ cards below them — the problem being reported.

---
*Source: do-work capture-request [Image #1] instead of boxes like this, make it a longer box (full width) and stacked vertically), basically improve the UI/UX at the moment it looks like it is a REQ card, which is not.*

## Addendum (2026-09-03, 23:45 local)

User added (23:35 local, applying the velocity report's triage table): deferred until the queue drain finishes. `assigned_to: 'after-drain'` makes the default scan skip and report this REQ; explicit targeting (`do-work run REQ-482`) clears it. Not a judgment on the change, only on when.
