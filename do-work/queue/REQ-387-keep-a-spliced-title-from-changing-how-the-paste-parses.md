---
id: REQ-387
title: '[impact-user-visible] Keep a spliced title from changing how the pasted Markdown parses'
status: pending
created_at: 2026-08-26T23:07:00Z
status_changed_at: 2026-08-26T23:07:00Z
user_request: UR-075
addendum_to: REQ-383
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379, REQ-383]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/web/board-core.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Keep A Spliced Title From Changing How The Pasted Markdown Parses

## What

An expanded title is inserted into the document's Markdown verbatim, after a 60-character cut. Two
characters in a title can change how the pasted file parses: a pipe inside a table row, and a backtick
the cut leaves unbalanced.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The feature's contract is that a paste saves back as a valid file and reads as the same document plus
titles. Two title characters break the second half:

- **A pipe inside a table row.** A mention in a GFM table cell gets its title spliced into the cell,
  and a title containing `|` adds a column: `| REQ-501 (Split the row | keep the pipe) | a row |`
  parses as three cells, not two.
- **A backtick the cut unbalances.** `shortTicketTitle` cuts at 60 characters on a word boundary. A
  title whose backticked span straddles that boundary contributes one backtick to the paste, which
  opens a code span that runs to the next stray backtick — swallowing the prose between them.

## Context

**Latent today, and worth fixing before it is not.** Six real titles carry backticks and none
currently produces an unbalanced cut; no real title carries a pipe. Both are one ordinary title away
from being live, and the failure is silent — the paste looks fine until someone re-renders it.

Both surfaces splice titles, but only the CLIPBOARD writes Markdown. The drawer builds DOM nodes and
sets `textContent`, so a pipe or a backtick there is inert. This is a clipboard-only defect even
though `shortTicketTitle` (`web/board-core.js`) is shared.

## Detailed Requirements

- **The cut must not leave a code span open.** Either extend the cut to the span's close, pull it back
  to the span's open, or strip the stray backtick — whichever reads best, but the pasted title must
  carry an even number of backticks.
- **A title spliced inside a table row must not add a cell.** The Go index already knows the block
  each mention sits in, so the two candidate approaches are: escape `|` as `\|` when the mention is
  inside a table, or suppress the expansion there and let the appendix carry the title. Prefer
  whichever keeps the table readable.
- **The appendix keeps the full untruncated title** and is not subject to either rule — it is a list,
  not prose, and it is where a reader looks up what the cut removed.

## Constraints

- **Do not change `shortTicketTitle`'s behaviour for the drawer.** It is shared, and the drawer needs
  no escaping; a change there must be additive or clipboard-side.
- **No new board write surface.**
- Do not widen into REQ-385, REQ-386 or REQ-388.

## Dependencies

None. Disjoint write set from the other three follow-ups.

## Red-Green Proof

**RED prompt/case:** A board where REQ-501 is titled `Split the row | keep the pipe` and a body
contains the table row `| REQ-501 | a row |`. Copy it: the pasted row has three cells. Second case: a
title whose backticked span crosses the 60-character cut — the pasted body opens a code span that
never closes cleanly.

**Why RED now:** the title is concatenated into the Markdown with no escaping and no balance check.

**GREEN when:** both pasted bodies re-render as the same document they came from, and feeding each
back through the repo's own renderer produces the original block structure.

**Validation:** Reproduced by adversarial review of REQ-383 with a purpose-built fixture board driven
through the shipped `annotateClipboardPayload`. Measured against the tree: 6 of 483 real titles carry
a backtick, 0 currently break.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-383's independent review, finding C1.*
