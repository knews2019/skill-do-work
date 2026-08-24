---
id: REQ-368
title: 'Add a UR copy-all that copies the UR plus all its REQs'
status: pending
created_at: 2026-08-24T14:53:28Z
user_request: UR-072
domain: frontend
prime_files: ['_dev/primes/prime-kanban-board.md']
tdd: true
suggested_spec:
depends_on: []
related: [REQ-367]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
---

# Add a UR Copy-All That Copies the UR Plus All Its REQs

## What
The UR detail view gets a copy-all control: one click puts the UR's own payload plus every REQ grouped under it on the clipboard, concatenated cat-style. Companion to REQ-367's per-column copy-all.

## Context
- The UR detail drawer already has a per-record `Copy` button (same plumbing as REQ copies: `web/board-clipboard.js` + the raw-Markdown bundle), but it copies only the UR record itself.
- "All its REQs" is the board's grouped set — the REQ IDS row the UR detail already renders (e.g. UR-065 shows 12 grouped REQs spanning queue, working, and archive), not just the capture-time `requests:` array.
- Copy-all is a control alongside the existing Copy, not a replacement: plain Copy keeps copying just the UR.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- The payload is the UR first, then each grouped REQ's full payload (the same content that REQ's own Copy yields), in the order the UR detail lists them.
- Reuse the existing clipboard write + copied/failed feedback pattern from `web/board-clipboard.js`.

## Red-Green Proof
**RED prompt/case:** Open UR-065's detail on the generated board. Its Copy button yields only the UR record; assembling the UR plus its 12 grouped REQs takes 13 separate copies. A browser-probe test (style of `browser_probe_test.go`) that clicks a UR copy-all control and asserts the payload contains the UR followed by every grouped REQ id fails today because the control does not exist.
**Why RED now:** The UR drawer's only copy control copies the single UR record.
**GREEN when:** One click on the UR detail's copy-all puts the UR's payload followed by all grouped REQs' payloads, concatenated in listed order, on the clipboard (probe asserts via a stubbed clipboard).
**Validation:** Inferred during capture

## Assets
- `do-work/user-requests/UR-072/assets/REQ-368-screenshot-1-ur-detail-drawer.png` — UR-065's detail view ("Audit the timeline view and make it more useful"): header with the existing `Copy`/`Close` buttons, a summary card showing GROUPED REQS 12 with the REQ IDS list (REQ-318…REQ-356), then Summary, Extracted Requests table, and Audit Findings sections. Shows the Copy button being extended and the grouped-REQs set copy-all must cover.

---
*Source: "similar request for the UR where copy all, will copy not only the UR but also all it's reqs to the clipboard."*
