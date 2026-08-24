---
id: REQ-367
title: 'Add a copy-all-REQs button to each Board column'
status: completed
claimed_at: 2026-08-24T20:17:59Z
completed_at: 2026-08-24T21:02:17Z
commit: 66e8de545450ef8b07e810f7700c5265141aed1a
status_changed_at: 2026-08-24T21:02:17Z
route: B
created_at: 2026-08-24T14:53:28Z
user_request: UR-071
domain: frontend
prime_files: ['_dev/primes/prime-kanban-board.md']
tdd: true
suggested_spec:
depends_on: []
related: [REQ-368]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-24T20:17:59Z
  basis:
    - trivial short-circuit
write_set:
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go
---

# Add a Copy-All-REQs Button to Each Board Column

## What
On the Board view, each column (Pending, Claimed, Needs Input · Blocked, Recently Done) gets a copy-all button that puts every REQ the column lists on the clipboard cumulatively — cat-style concatenation of the same per-REQ payload the existing detail-drawer Copy button produces.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Traced column rendering, active filters/windows, raw Markdown loading, clipboard feedback, and browser seams, then froze the exact five-file boundary and DOM-order authority.
- [x] **[APPLY]:** Added four accessible controls, render-derived enable/count state, exact raw-payload concatenation through the shared writer, visible failure behavior, styling, and one generated-site Chromium probe.
- [x] **[UNIFY]:** Reviewed all five files and rendered output; focused browser, five mutations, sliced JavaScript, full module, canonical verification, syntax, formatting, diff, and artifact checks passed.

## Why (if provided)
Great for extracting all pending, claimed, blocked, or done REQs in one action instead of copying each card's ticket one at a time.

## Context
- Today each REQ can be copied individually from the detail drawer (`Copy` in the drawer header; `web/board-clipboard.js` + the raw-Markdown bundle in `web/board-markdown.js`). There is no bulk copy anywhere.
- The user's "each row" means each Board-view lane, per their own enumeration: pending, claimed, blocked, done.
- "cat all reqs": concatenate each REQ's full payload one after another, in the column's display order — the same content a per-REQ Copy yields, not just a list of ids.
- **Time filters apply for the done column** (user's words): with the 24h/48h/7d Recently Done window active, copy-all on that column copies only the REQs inside the selected window.

## Detailed Requirements
- One copy-all control per Board column, including a column with only one card; a column showing "Nothing here" has nothing to copy (disable or no-op with feedback — builder decides).
- The copied set matches what the column currently lists: search/domain/status filters and the done-column time window all apply (user confirmed at verify).
- Reuse the existing clipboard write + copied/failed feedback pattern from `web/board-clipboard.js` rather than inventing a second mechanism.

## Open Questions
- [x] Should copy-all respect active filters (search box, domain, status) beyond the done-column time window? → Yes — copy exactly the cards the column currently shows (user confirmed, 2026-08-24 verify)

## Red-Green Proof
**RED prompt/case:** Open the generated Board view with several pending REQs. There is no control on the Pending column header to copy them together — extracting them requires opening each card and clicking Copy N times. A browser-probe test (style of `browser_probe_test.go`) that looks for a column copy-all control and asserts the composed clipboard payload contains every listed pending REQ fails today because the control does not exist.
**Why RED now:** Copy exists only per-record in the detail drawer; no bulk copy surface exists on columns.
**GREEN when:** Clicking the Pending column's copy-all button puts every currently listed pending REQ's payload, concatenated in display order, on the clipboard (probe asserts via a stubbed clipboard). On Recently Done with 24h selected, the copied set equals the cards shown for that window.
**Validation:** Inferred during capture

## Assets
- `do-work/user-requests/UR-071/assets/REQ-367-screenshot-1-board-columns.png` — Board view of this repo's queue: four columns (PENDING 14, CLAIMED 1, NEEDS INPUT · BLOCKED 0 showing "Nothing here", RECENTLY DONE 24 with the 24h/48h/7d selector top right), REQ-356's detail drawer open on the right with the existing per-REQ `Copy` button in its header. Shows the column headers where the new button belongs and the per-REQ Copy being extended.

---
*Source: "in the kanban each REQ can be copied, now I want each row to have a copy all REQs button where each REQ will be added cumulatively (see bash like terminology: cat all reqs) to the clipboard, this is great for extracting all pending, claimed, blocked, done reqs (time filters apply for the done column)"*

## Triage

**Route: B** — The behavior is specified, but the absent capture-time write set must be discovered across the column header, filtered display model, raw-Markdown bundle, shared clipboard feedback, styling, and trusted browser coverage before implementation.

## Exploration

- The exact displayed order already exists in each column's rendered `.req-card` DOM, including Pending's Ready-then-Waiting grouping and all active filters/windows. Read those ids at click time instead of reconstructing a parallel filter model.
- Reuse the raw Markdown map and existing clipboard writer/fallback. Concatenate each visible card's existing payload byte-for-byte with `join("")`; fail visibly if any rendered id lacks raw content rather than copying an incomplete set.
- Add one semantic control to each of the four flat Board headers only. Column rendering owns enabled/disabled state and an accessible count label; Testing, By UR, and REQ-368's UR detail remain out of scope.
- A generated-site Chromium probe will stub the clipboard and cover all four controls, empty state, exact pending display order, a one-card column, filter narrowing, Recently Done 24h/48h membership, raw payload fidelity, URL/console guards, and mutation seams.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`
- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js`
- `skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go`

**Acceptance criteria:**

- Each flat Board column has one accessible copy-all control whose enabled state and count match its currently rendered REQ cards.
- Clicking copies every visible card's existing raw-Markdown payload, byte-for-byte and in display order, through the shared clipboard/fallback and feedback path.
- Search/domain/status filters and the Recently Done time window change both the visible set and copied payload; empty columns cannot produce a misleading successful copy.
- Trusted Chromium coverage proves exact content/order, filter/window behavior, empty and single-card states, URL provenance, console cleanliness, and non-vacuous mutation seams.

## Decisions

- **D-01 — Treat the rendered column DOM as the membership/order authority.** It is already the exact filtered/windowed display contract the user asked to copy, avoiding a second selection model that can drift.
- **D-02 — Fail closed on missing raw payloads.** Silently omitting one visible REQ would make a successful-looking bulk copy incomplete; the existing failed-feedback path makes the mismatch visible.
- **D-03 — Keep UR bulk copy separate.** This REQ adds controls only to the four flat Board columns; REQ-368 owns UR-plus-grouped-REQ composition.
- **D-04 — Disable empty columns.** Explicit disablement avoids a clipboard no-op or transient success-like feedback when the rendered set is empty.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified): adds one semantic Copy all control to each of the four flat Board column headers.
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified): styles compact header actions and their disabled, focus, copied, and failed states without changing card geometry.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified): derives each control's enabled/count label from the exact rendered filtered/windowed id slice while preserving Pending Ready-then-Waiting order and REQ-366 wording.
- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (modified): generalizes shared feedback, reads visible card ids in DOM order, lazy-loads raw Markdown, concatenates exact bytes with no invented separator, and fails closed on a missing payload.
- `skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go` (new): proves four controls, empty and single-card states, exact order/content, combined filters, Recently Done windows, success/failure feedback, URL/console guards, and five mutation seams in Chromium.

## Discovered Tasks

None.

## Testing

- The focused Chromium 1228 probe passed with exactly four controls, disabled empty Needs Input, Pending order REQ-102/101/103, one-card Claimed copy, combined filter narrowing, and 24h-to-48h Recently Done membership.
- Five independent mutations—newline separator, first-only truncation, sorted ids, unfiltered enablement, and hard-coded 48h window—each made the focused probe RED and were restored.
- Exact raw payloads, five successful writes, copied feedback, forced clipboard/fallback failure, probe URL, and zero console/window errors were asserted. Visual QA confirmed all controls fit a 1440×1000 generated board.
- Sliced JavaScript, full queue-kanban module, Node syntax, formatting, diff checks, and the builder canonical gate passed. The strict browser aggregate's only RED was the unrelated Timeline pointer-capture baseline tracked by REQ-370.
- On the combined REQ-366/367 main tree, the focused Chromium probe and full queue-kanban module passed again, preserving dependency-gated blocked cards in the copied Pending DOM order. The final canonical gate passed all contracts, Go tests, strict JavaScript, and audit metrics; its optional browser lane made the standard no-browser skip after the focused browser run.

## Qualification

- Exact merge range `4821b34..66e8de545450ef8b07e810f7700c5265141aed1a` passed mechanical qualification.
- Scope drift passed: the five changed files exactly match the frozen Scope and Implementation Summary.
- Orchestrator judgment confirmed substantive visible-DOM membership flow, complete filter/window/order tracing, exact raw-byte composition, fail-closed behavior, non-vacuous browser evidence, and no generated/debug artifacts.

## Review

Independent review approved at 10/10 with no blocking, major, or minor findings. It confirmed exact five-file scope, rendered-DOM membership/order authority, byte-exact concatenation, filter/window behavior, empty/failure feedback, accessibility, five mutation seams, and clean REQ-366 integration. Residual risk is limited to non-blocking test-depth gaps around exact aria-label wording and a missing-map mutation.

## Lessons Learned

For bulk actions over a filtered UI, the rendered DOM can be the safest membership and ordering contract when it already embodies every filter, subgroup, and time window. Recomputing the set elsewhere would create a second model that can silently drift.

## Orientation

Released in 0.236.56. Every flat Board column now copies the exact visible REQ payloads in display order with shared success/failure feedback.
