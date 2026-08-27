---
id: REQ-389
title: 'Addendum: mark spliced paste titles with a leading arrow'
status: pending
created_at: 2026-08-27T11:22:38Z
user_request: UR-078
addendum_to: REQ-379
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-387]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-379, REQ-383, REQ-387]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Addendum: Mark Spliced Paste Titles With A Leading Arrow

## What

The Copy buttons splice a ticket's title after the first body mention of its id, as
`REQ-374 (Show how long each done card took)`. Change the spliced form to
`REQ-374 (-> Show how long each done card took)` so a reader of the paste can tell the
parenthetical was inserted by the board, not written by the ticket's author.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

"so I know that this is a programatic expand" — a pasted body carries no styling, so
without a marker the parenthetical reads as the author's own words.

## Context

- The insertion text is built client-side in
  `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (`" (" + expandedTitle + ")"`,
  currently around line 199); the insertion positions come precomputed from the Go side (REQ-383).
- Scope is the in-body splice only. The drawer already marks the expansion visually (the title
  renders in its own quieter-styled span), and the Referenced requests appendix is self-evidently
  board-generated, so neither changes. Capture judged this from the user's example, which shows
  only the in-body form.
- The marker is ASCII `->` followed by a space, exactly as the user typed: `(-> Title)`. Do not
  substitute a typographic arrow.

## Prior Implementation

- REQ-379 (commit f1e2ce8) shipped the clipboard title splice and the Referenced requests appendix.
- REQ-383 (commit a3d4e4c) moved mention resolution into Go (`citations.go`) and deleted the
  client-side Markdown scanner; the page now inserts titles at positions the board computed.
  The literal inserted text remains a client-side concern in `board-clipboard.js`.

## Constraints

- This is a deliberate paste-only divergence from the drawer's rendering: the drawer signals the
  expansion with styling, the paste signals it with `->`. Work descending from the
  drawer/clipboard-divergence line (REQ-386/REQ-388) must not "settle" this marker away.
- Any test pinning the old `(<title>)` spliced form updates to the new `(-> <title>)` form.

## Dependencies

`depends_on: [REQ-387]` serializes the shared files (`web/board-clipboard.js`,
`generate_test.go`) behind the end of the queued ticket-id chain
(REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387). The edge orders writers of the
same files; it is not a need for REQ-387's output.

## Red-Green Proof

**RED prompt/case:** Copy a ticket whose body's first mention of REQ-374 gets a title spliced.
The paste reads `REQ-374 (Show how long each done card took)`.
**Why RED now:** `board-clipboard.js` builds the insertion as `" (" + expandedTitle + ")"` — no
marker distinguishes the inserted title from author-written parentheticals.
**GREEN when:** The same paste reads `REQ-374 (-> Show how long each done card took)`, and a
harness test pins the `(-> ` prefix on the spliced insertion.
**Validation:** User confirmed — the before/after forms are the user's own example, verbatim.

---
*Source: "Instead of: REQ-374 (Show how long each done card took). expand it like this REQ-374 (-> Show how long each done card took). so I know that this is a programatic expand"*
