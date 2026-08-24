---
id: REQ-367
title: 'Add a copy-all-REQs button to each Board column'
status: pending
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
---

# Add a Copy-All-REQs Button to Each Board Column

## What
On the Board view, each column (Pending, Claimed, Needs Input · Blocked, Recently Done) gets a copy-all button that puts every REQ the column lists on the clipboard cumulatively — cat-style concatenation of the same per-REQ payload the existing detail-drawer Copy button produces.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
