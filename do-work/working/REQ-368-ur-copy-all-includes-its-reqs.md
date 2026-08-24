---
id: REQ-368
title: 'Add a UR copy-all that copies the UR plus all its REQs'
status: claimed
claimed_at: 2026-08-24T21:05:47Z
status_changed_at: 2026-08-24T21:05:47Z
route: B
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
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-24T21:05:47Z
  basis:
    - trivial short-circuit
write_set:
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go
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

## Open Questions
- [x] Separate Copy-all button, or repurpose the existing Copy? → Separate button beside Copy; plain Copy keeps copying just the UR record (user confirmed, 2026-08-24 verify)

## Red-Green Proof
**RED prompt/case:** Open UR-065's detail on the generated board. Its Copy button yields only the UR record; assembling the UR plus its 12 grouped REQs takes 13 separate copies. A browser-probe test (style of `browser_probe_test.go`) that clicks a UR copy-all control and asserts the payload contains the UR followed by every grouped REQ id fails today because the control does not exist.
**Why RED now:** The UR drawer's only copy control copies the single UR record.
**GREEN when:** One click on the UR detail's copy-all puts the UR's payload followed by all grouped REQs' payloads, concatenated in listed order, on the clipboard (probe asserts via a stubbed clipboard).
**Validation:** Inferred during capture

## Assets
- `do-work/user-requests/UR-072/assets/REQ-368-screenshot-1-ur-detail-drawer.png` — UR-065's detail view ("Audit the timeline view and make it more useful"): header with the existing `Copy`/`Close` buttons, a summary card showing GROUPED REQS 12 with the REQ IDS list (REQ-318…REQ-356), then Summary, Extracted Requests table, and Audit Findings sections. Shows the Copy button being extended and the grouped-REQs set copy-all must cover.

---
*Source: "similar request for the UR where copy all, will copy not only the UR but also all it's reqs to the clipboard."*

## Triage

**Route: B** — REQ-367 now supplies the shared exact-byte bulk clipboard path, but exploration must identify the UR drawer control and grouped-id ordering source, distinguish grouped membership from capture-time requests, and add trusted browser coverage before freezing scope.

## Exploration

- The drawer's REQ ids row and `userRequestsById[UR].requestIds` already share one sorted all-tree membership source derived from each REQ's `user_request`, spanning queue, working, and archive rather than the UR's capture-time `requests:` array.
- Add a hidden-by-default native control beside plain Copy and expose it only for UR details. Existing detail-copy focus and feedback styling applies, so no CSS change is needed.
- Snapshot the current UR id and grouped ids at click time, then compose exact raw UR bytes followed by `rawMarkdownForRequests` output with no separator. Missing UR or member raw data fails the whole operation visibly.
- A generated-site Chromium probe will distinguish plain Copy from Copy all, include grouped members beyond the capture array and across states, prove displayed order and exact payload, cover a no-member UR, REQ-detail hiding, missing-member failure, URL/console guards, and mutations.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js`
- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js`
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go`

**Acceptance criteria:**

- A separate Copy all control appears beside plain Copy only for UR details; plain Copy and REQ detail behavior remain unchanged.
- Copy all writes the UR's exact raw payload first, then every grouped REQ's exact raw payload in the same all-tree order displayed by the drawer, with no invented separator.
- Missing UR/member raw data fails closed through shared feedback, while a UR with no grouped REQs copies its own payload successfully.
- Trusted Chromium coverage proves cross-state/capture-array-independent membership, ordering/content, control visibility, no-member and failure states, URL provenance, console cleanliness, and mutation sensitivity.

## Decisions

- **D-01 — Reuse the drawer's all-tree grouped-id authority.** This keeps copied membership and ordering identical to the REQ ids the user sees and avoids the incomplete capture-time array.
- **D-02 — Snapshot detail identity at click time.** The async Markdown load operates on the UR and members the drawer showed when the user invoked the action.
- **D-03 — Reuse native detail styling.** The existing detail-copy styles already cover focus and feedback; adding CSS would duplicate a solved surface.
