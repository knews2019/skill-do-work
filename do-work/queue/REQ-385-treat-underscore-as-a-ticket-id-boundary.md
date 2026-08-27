---
id: REQ-385
title: '[impact-user-visible] Treat an underscore as a ticket-id boundary on both surfaces'
status: pending
created_at: 2026-08-26T23:05:00Z
status_changed_at: 2026-08-26T23:05:00Z
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
  - skills/do-work-board/tools/queue-kanban/citations.go
  - skills/do-work-board/tools/queue-kanban/citations_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
---

# Treat An Underscore As A Ticket-Id Boundary On Both Surfaces

## What

`\b` counts `_` as a word character, so the mention pattern behaves differently around an underscore
than a reader expects. Change the boundary on both the Go and the client side together, in one
commit, so the agreement test stays green.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Two symptoms, one cause. Both were found by adversarial review of REQ-383 and reproduced.

**A ticket id in underscore emphasis is silently dropped.** `Fixed in _REQ-1679_ last week.` yields no
mention at all — no title, no appendix line — because `_REQ-1679_` fails the pattern's `\b` anchors in
the SOURCE bytes. The drawer scans the RENDERED text, where emphasis is already consumed, so it
resolves the id, expands the title and adds a glossary entry. `*REQ-1679*` works and `_REQ-1679_` does
not, for the same rendered output. That is a drawer/clipboard divergence, and REQ-383's stated rule is
that the two must say the same thing about the same body.

**A compound id followed by an underscore corrupts on paste.** RE2 explores all alternatives, so when
the compound alternative's trailing `\b` fails against a following word character, the shorter
`UR-\d+` alternative succeeds at the same start. `_tracked under UR-003-REQ-077_` therefore emits a
mention for the six characters `UR-003` with `Expand: true`, and the client splices the UR's title
into the middle of the compound id:

```
_tracked under UR-003 (Ship The Widget)-REQ-077_
```

The glossary then lists `UR-003` instead of `UR-003-REQ-077`.

## Context

**Not a regression.** The pattern is byte-identical to the one REQ-379 shipped and the one
`web/board-detail.js` still carries; REQ-383 moved block classification into Go and left match
semantics alone on purpose. This REQ is the follow-up that changes them.

**The corruption is latent; the dropped mention is not.** Zero of this board's 485 `id:` fields are
compound — every id in the tree is the flat `REQ-nnn` form, and the flat form degrades to NO match
rather than a corrupt splice, which is the benign direction. Underscore emphasis around an id needs no
special setup at all.

**Why both files move together.** `bodyTicketMentionPattern` (`citations.go`) and `bodyMentionPattern`
(`web/board-detail.js`) are pinned to each other by
`TestJavaScriptBehaviorTicketMentionPatternAndResolverAgreeWithGo`, which drives both over one corpus
in both directions. Changing one alone fails that test — correctly. This is why REQ-383 could not do
it: the drawer was explicitly out of its scope.

## Detailed Requirements

- **A ticket id's boundary is a non-alphanumeric character, not a non-word character.** `_` ends an
  id the way a space or a bracket does. RE2 has no lookaround, so the boundary cannot simply be
  rewritten as a lookahead — either capture and restore the boundary characters, or post-filter each
  match against the bytes on either side. The client must use whichever shape the Go side uses.
- **The compound alternative must win wherever it matches at all.** A trailing word character must not
  demote `UR-003-REQ-077` to `UR-003`; if the compound form cannot match cleanly, the correct answer
  is no match, never a shorter one at the same start.
- **Both surfaces change in one commit**, keeping the agreement test green throughout.

## Constraints

- **The agreement test is the gate, not a formality.** Its corpus must gain the underscore cases in
  the same commit; a fix that passes the current corpus has not been tested.
- **Do not widen into the other REQ-383 follow-ups.** REQ-386 (the restating H1), REQ-387 (unescaped
  titles) and REQ-388 (the remaining divergences) are separate.

## Dependencies

None. REQ-383 is complete; this changes the pattern it left alone.

## Red-Green Proof

**RED prompt/case:** Copy a ticket whose body contains `Fixed in _REQ-1679_ last week.` — the paste
carries no title and no appendix line, while the drawer beside it shows both. Second case: a body
containing `_tracked under UR-003-REQ-077_` on a board holding both `UR-003` and `UR-003-REQ-077` —
the paste reads `UR-003 (…)-REQ-077`.

**Why RED now:** `\b` treats `_` as a word character, so the first case fails the anchor and the
second falls through to the shorter alternative.

**GREEN when:** both cases behave as the drawer does, the agreement corpus carries them, and no
existing assertion needed rewriting.

**Validation:** Reproduced by adversarial review of REQ-383, both cases run against a scratch copy of
the package.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-383's independent review, findings S1 and S2.*
