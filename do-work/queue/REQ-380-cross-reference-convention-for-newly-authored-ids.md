---
id: REQ-380
title: '[impact-rule-change] Cross-Reference Convention for newly authored REQ and UR ids'
status: pending
created_at: 2026-08-26T13:02:24Z
user_request: UR-075
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work/actions/capture-reference.md
---

# Cross-Reference Convention For Newly Authored REQ And UR Ids

## What

One canonical section in `skills/do-work/actions/capture-reference.md` saying that any flow writing
a REQ or UR id into prose names it — `REQ-1679 (Admin can delete a card)` on first mention, bare
afterwards — while frontmatter fields keep bare ids.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-378 and REQ-379 patch the symptom at render and copy time. The cause is that the reference was
written as a bare number in the first place, which means it stays cryptic everywhere the board is
not: in a plain editor, in `git grep` output, in a pasted diff, in a review thread. The user's
framing is the general rule, not a board rule — "we shouldn't assume that people know what those
numbers mean."

## Context

`capture-reference.md` already carries two conventions in exactly this shape, both keyed on a
condition so their applicability lists cannot go stale:

- **REQ Title Convention** (line 9): "The condition is the rule — **any flow that mints a REQ
  carrying an `impact:` value follows this section** — so a new one inherits it without this list
  being re-counted."
- the **fold-before-mint contract** (line 24), with the same construction.

The convention this REQ writes down is also already in practice: that same file's own worked
example at line 221 reads
`- REQ-018 (durations-measured-face-constants-lack-provenance) — the second finding, that the face
numbers carry no provenance`. So this names an existing habit rather than inventing one, which is
why it is one section and not a campaign.

## Detailed Requirements

- **One new section in `capture-reference.md`**, written in that file's established
  condition-keyed form. Content:

  > **Cross-Reference Convention.** Any flow that writes a REQ or UR id into prose names it:
  > `REQ-1679 (Admin can delete a card)` on first mention in a document, bare afterwards. The
  > condition is the rule — capture, review findings, remediation notes, discovered tasks and
  > handoffs all inherit it without this list being re-counted. Frontmatter fields (`depends_on`,
  > `related`, `addendum_to`, `user_request`) stay bare ids: they are parsed, not read. Never
  > rewrite an existing file's references to add titles.

- **First mention only**, matching REQ-378's display rule — the two must not disagree about what
  the convention is.
- **Frontmatter stays bare.** State it explicitly; it is the part someone will get wrong.
- **No other action file is edited.** The flows that mint REQ prose already cite
  `capture-reference.md` as the canonical home for its conventions, so a second copy anywhere else
  is drift, not reinforcement.
- **Never retro-fit.** The convention governs newly authored references; existing REQ files are
  left exactly as they are.

## Constraints

- **Do not add a checker for this.** A prose convention enforced by a grep would fire on the
  section defining it — the failure `_dev/primes/prime-kanban-board.md` records from REQ-293, where
  "a prose-grep guard will fail on its own contract being stated". If a lock-in is wanted later it
  is its own REQ with that trap solved first.
- **One section, no list.** Do not enumerate the flows that inherit the rule; key it on the
  condition, per root `CLAUDE.md` § "State conditions, not lists".
- Keep it short. This is a convention statement, not a guide.

## Dependencies

None. Independent of REQ-378 and REQ-379 — it changes an action file, they change the board client,
and no file is shared.

## Builder Guidance

**Certainty level: Firm.** The wording above is the deliverable; adjust it for fit with the
surrounding sections' voice, not for scope.

Read `_dev/primes/prime-action-files.md` first — it owns the template, the earned-sections rule and
the cross-referencing conventions for anything under `skills/*/actions/`.

Resist growing this. The temptation is to also touch `review-work.md`, `work-reference.md` and the
handoff actions "for consistency"; that is exactly the drift the canonical-home pattern exists to
prevent.

## Red-Green Proof

**RED prompt/case:** Search `skills/do-work/actions/capture-reference.md` for a stated rule about
how a REQ or UR id is written into prose. There is none — the id-in-parens form appears only
incidentally, inside one worked example of the Folded Requests section, where a reader has no
reason to read it as a rule that applies anywhere else.

**Why RED now:** No section of any action file states the convention, so every flow that writes a
cross-reference decides the form for itself, and bare ids are the common outcome.

**GREEN when:** `capture-reference.md` carries a named Cross-Reference Convention section stating
the first-mention-with-title form, the frontmatter exception, and the never-retro-fit rule, keyed on
a condition rather than a list of flows — and no other action file gained a competing copy.

**Validation:** User confirmed — "Board + author-side convention" was selected from two presented
options during planning, with the alternative being board-only.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: user request in session, 2026-08-26 — "we shouldn't assume that people know what those numbers mean."*
