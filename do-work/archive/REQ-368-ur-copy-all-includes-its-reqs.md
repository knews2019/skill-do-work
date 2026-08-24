---
id: REQ-368
title: 'Add a UR copy-all that copies the UR plus all its REQs'
status: completed
claimed_at: 2026-08-24T21:05:47Z
completed_at: 2026-08-24T21:41:02Z
commit: 61dddb24ef1e412111fc14e442c8474940d4416b
status_changed_at: 2026-08-24T21:41:02Z
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
- [x] **[PLAN]:** Traced the UR drawer's all-tree grouped-id authority and merged clipboard path, then froze a four-file boundary with no CSS change and explicit plain-Copy separation.
- [x] **[APPLY]:** Added a UR-only Copy all control, current-detail identity reset/visibility, exact UR-plus-grouped-REQ composition, fail-closed semantics, shared feedback, and a generated-site Chromium probe.
- [x] **[UNIFY]:** Reviewed all four files and headed render geometry; RED/GREEN, five mutations, REQ-367 compatibility browser lane, full module, canonical, syntax, formatting, diff, and artifact checks passed.

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
- **D-04 — Fail synthesized or incomplete UR bundles closed.** A UR with no real raw payload or any missing member must produce visible failure and no partial write; a real zero-member UR remains a valid UR-only bundle.
- **D-05 — Keep asynchronous feedback bound to the same open UR.** Success or failure updates the control only if the drawer still shows the UR whose ids were snapshotted.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified): adds a hidden-by-default native Copy all control beside the existing drawer Copy action.
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified): resets both controls on drawer open and exposes Copy all only for UR details while preserving plain Copy and REQ behavior.
- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (modified): snapshots the current UR and grouped ids, composes exact UR-first bytes through the shared fail-closed raw Markdown path, and applies feedback only while the same UR remains open.
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` (new): proves cross-state all-tree membership beyond the capture array, displayed order, exact plain/all payloads, zero-member success, synthesized/missing-member failure, REQ hiding, feedback, URL/console guards, and five mutation seams.

## Discovered Tasks

None.

## Testing

- TDD RED failed because no Copy all control existed. GREEN Chromium proved controls exactly Copy/Copy all/Close and grouped queue, working, and archive REQs in displayed numeric order despite a stale one-member capture array.
- Exact plain-Copy UR bytes, exact UR-plus-three-REQ cat bytes, zero-member UR-only success, missing-member and synthesized-UR failure without partial writes, REQ-detail hiding, three successful feedback states, URL provenance, and zero application errors passed.
- Five mutations—empty grouped snapshot, REQ-first composition, newline separator, silent missing-member omission, and all-detail visibility—each made the focused probe RED and were restored.
- REQ-367 plus REQ-368 focused Chromium compatibility, full queue-kanban module, canonical verification, Node syntax, formatting, diff checks, and artifact checks passed. Headed Chrome QA confirmed UR-065's 12 members and non-overlapping Copy/Copy all/Close geometry at 1440×1000.
- On merged main, the combined REQ-367/368 Chromium compatibility lane, full queue-kanban module, and final canonical maintainer gate passed again; the canonical browser lane made its standard no-browser skip after the focused Chromium run.

## Qualification

- Exact merge range `a806fc3d0a78465294bfe655068e59db256ce864..61dddb24ef1e412111fc14e442c8474940d4416b` passed mechanical qualification; the new Go file is a test entry point under the static-reference exception.
- Scope drift passed: all four changed files exactly match the frozen Scope and Implementation Summary.
- Orchestrator judgment confirmed substantive all-tree grouping flow, complete plain/all/failure requirement tracing, exact raw-byte composition, non-vacuous browser evidence, and no generated/debug artifacts.

## Review

Independent review approved with no Important, Minor, or Nit findings. Correctness scored 100, tests/mutations 99, accessibility/layout 99, overall 99/100, low risk. It independently passed the combined Chromium lane three times, replayed all five RED mutations, measured minimum-drawer geometry, and confirmed exact bytes/order, failure atomicity, plain-Copy preservation, and REQ-367 compatibility.

## Lessons Learned

Grouped UI membership should come from the same derived all-tree relation the user sees, not an older capture-time list. Snapshot that authority before async work, then fail the entire bundle if any displayed member cannot be represented exactly.

## Orientation

Released in 0.236.58. UR details now offer Copy all for the UR's exact source followed by every currently grouped REQ across queue, working, and archive.
