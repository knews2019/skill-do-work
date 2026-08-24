---
id: REQ-359
title: "Review fix: suppress the duplicated status finding from the board strip"
status: pending
created_at: 2026-08-24T10:40:00Z
user_request: UR-068
addendum_to: REQ-343
domain: frontend
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/generate.go
---

# Suppress the Duplicated Status Finding From the Board Strip

## What

`boardRenderedVerifyCategories` exists to stop the board printing the same finding twice. REQ-343's
new `unrecognized-req-status` category meets that map's own stated criterion and was shipped
unsuppressed, so the class now reaches the page three times. Add it to the map, and amend the map's
doc comment in the same edit.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The map states its own criterion in its doc comment (`generate.go:524-535`): *"the findings the board
already shows by other means, so forwarding them would print the same prose a second or third time"*.
`unrecognized-req-status` meets it exactly. All three surfaces were confirmed during REQ-343's review:

1. `bucketColumns` (`model.go:1618-1620`) appends a status warning to `board.Warnings`, rendered by
   `renderWarningsBanner` (`web/board-cards.js:455-476`)
2. a per-card badge from `request.statusUnrecognized` (`web/board-cards.js:62`, and in
   `board-detail.js:268`, `board-calendar.js:146`, `board-timeline.js:1911`)
3. the findings strip, as of REQ-343

**`structurally-damaged-req` is genuinely new and must keep forwarding.** This REQ is about the
status class only.

REQ-343's builder escalated this rather than reaching into `generate.go`, which was the right call —
the file is outside its write set and a builder editing a serial-owned file races its siblings. The
recorded impact was `impact-negligible`; the review raised it, because a reader meets this at exactly
the moment they are debugging a broken status, and meets the same sentence three times.

## Detailed Requirements

- `unrecognized-req-status` is suppressed from the board's findings strip.
- `structurally-damaged-req` continues to forward.
- The map's doc comment is amended in the same commit if the criterion's wording needs it — a
  criterion that a shipped category visibly fails to obey is a half-truth, which is the prime's
  REQ-231 lesson applied to the map itself.
- `verify`'s own exit status and finding list are unchanged. This is a board-rendering change only:
  the mechanical check must still fail on an unrecognized status.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs. Read it first.
- Confirm by rendering, not by reading the code path. The prime's render-evidence rule applies — the
  review judged this from the code and said so.

## Red-Green Proof

**RED prompt/case:** Generate a board for a repository holding a REQ with a typo'd status. The same
finding appears in the warnings banner, on the card, and in the findings strip.

**GREEN when:** it appears in the banner and on the card but not in the strip; a
`structurally-damaged-req` finding still appears in the strip; and `verify` still exits non-zero for
the typo'd status.

**Validation:** Inferred during REQ-343's review; the three surfaces were confirmed against source.

---
*Source: REQ-343 review finding I2 (UR-068), folding the builder's own Discovered Task and raising its impact token.*
