---
id: REQ-382
title: 'Expand ticket ids written as Markdown links'
status: pending
created_at: 2026-08-26T17:03:51Z
user_request: UR-075
addendum_to: REQ-378
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-381]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-379, REQ-381]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Expand Ticket Ids Written As Markdown Links

## What

A REQ body that writes a ticket as an explicit Markdown link — `[REQ-123](https://…)` — gets neither
a title nor a glossary entry. `linkifyDetailBody` skips any text node already inside an `<a>`, so the
one place an author has gone out of their way to mark a reference is the one place REQ-378 does not
reach.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Found by the Codex reviewer on PR #169 against REQ-378, verified, and deliberately deferred rather
than folded into that REQ's fourth remediation round. Two reasons, both recorded on the thread:
the corpus has **zero** instances today, and the fix touches the guard that keeps the walker from
re-entering the renderer's own autolinks — the highest-risk line in the change, traded against no
current benefit.

Deferred is not dismissed. The rule REQ-378 established is that a reader never meets a bare number,
and this is the one authoring shape where that rule does not hold.

## Context

`linkifyDetailBody` (`web/board-detail.js`) walks text nodes and skips any whose
`parentElement.closest("a")` is non-null. That guard does two jobs at once:

1. it stops the walker re-processing links the Markdown renderer itself produced (GFM autolinks,
   which the same function has just retargeted to `_blank`);
2. it stops a mention being wrapped twice.

Only the first is a real constraint. An author-written `[REQ-123](…)` is a text node inside an
anchor exactly like an autolink, and the guard cannot currently tell the two apart.

The measurement, taken at capture: `[…REQ-NNN…](…)` in a REQ or UR body — **0 matches across 373
REQs and 76 URs**. The pattern does occur in `_dev/primes/prime-kanban-board.md`, but primes are
never rendered in the drawer; only REQ and UR bodies reach `linkifyDetailBody`.

## Detailed Requirements

- **An anchor whose text carries a resolvable ticket id gains that ticket's title**, under REQ-378's
  existing rules: first mention expands, later mentions stay bare, the untruncated title rides in
  the tooltip, and `shortTicketTitle`'s 60-character word-boundary cut applies unchanged.
- **The author's `href` survives untouched.** This is the requirement that makes the REQ a design
  question rather than a mechanical edit: the author pointed that link somewhere deliberately, and
  the expansion is additive decoration on their anchor, never a replacement of their destination.
- **The glossary entry records the ticket, not the destination.** A reader looking up `REQ-123` in
  the glossary wants the REQ's title and status; where the author's link happened to point is not
  that.
- **A resolvable id inside an anchor earns a glossary line even when it does not expand** — the same
  rule REQ-378 already applies to a backticked mention.
- **Autolinks the renderer produced are still skipped.** The guard's real job is preserved; only its
  accidental second job is narrowed.
- **A dead id inside an anchor is not flagged.** An author who wrote a link made a deliberate
  reference to something; REQ-378's broken-reference flag is for ids that resolve to nothing in
  prose, and painting an author's own link in the blocked accent asserts more than is known.

- **The clipboard half: a link reference definition must not be rewritten.** Folded in from REQ-379's
  review (finding F6). `[REQ-100]: https://example.com/x` is a *definition*, and the clipboard's
  annotator currently turns it into `[REQ-100 (Alpha ticket)]: https://example.com/x`, which orphans
  every `[REQ-100]` reference in the pasted document — the paste silently loses links it had. Same
  root cause this REQ already owns (a ticket id written in Markdown link syntax), so it is one fix
  across both surfaces rather than two. Zero instances in `do-work/` today, as with the drawer half.

## Constraints

- **Never wrap a mention twice.** The regression this guards against is a fragment nested inside a
  fragment; a test must prove a body survives two passes unchanged.
- **No change to the mention pattern.** `bodyMentionPattern` stays in lock-step with
  `repoFileMentionPattern` in `filementions.go`; this REQ changes where the walker looks, not what
  it matches.
- **No new board write surface**, and no Go source change outside the test file.

## Dependencies

`depends_on: [REQ-381]`, which depends on REQ-379 — so the chain is REQ-379 → REQ-381 → REQ-382.

**The edge serializes a shared file; it is not a need for the others' output.** REQ-378 (archived)
established the rules this extends, and REQ-379 and REQ-381 change different source files. What all
three share is `generate_test.go`. Stating "do not fan these out" in prose enforces nothing —
`write_set` is display-only and "never a safety guarantee" (root `CLAUDE.md` § Glossary), so under
`do-work run --fan-out` a REQ with `depends_on: []` is a dependency root and gets dispatched
concurrently regardless of what its prose says. `depends_on` is the only field the work loop gates
on. This is the second time the same mistake was made in this batch; REQ-381 carries its edge for
exactly the same reason.

## Builder Guidance

**Certainty level: Mixed.** The requirement is firm; **how** the two anchor kinds are told apart is
yours to decide and to write down as a `## Decisions` entry.

The obvious discriminator — did the Markdown renderer make this anchor, or did the author — is not
directly observable after rendering. Candidates worth weighing: marking renderer-produced autolinks
during the existing retarget pass so the walker can recognise them later; comparing the anchor's
text against its `href` (an autolink's text *is* its href); or handling anchors in a separate pass
with its own rules rather than widening the text-node walk. Prefer whichever makes the
double-wrapping regression impossible rather than merely untested.

Start with the two-pass regression test, not the fix: it is what makes the candidates comparable and
it is the failure that would be expensive to ship.

## Red-Green Proof

**RED prompt/case:** Give a REQ body the line `See [REQ-1108](https://example.com/spec) for the
shape.` and open its drawer. The link renders as the bare text `REQ-1108`, carries no title and no
tooltip, and the glossary has no entry for it — while the same id written as plain prose in the same
body expands and is glossed.

**Why RED now:** `linkifyDetailBody` returns early for any text node whose `parentElement.closest("a")`
is non-null, so the mention is never offered to `buildLinkifiedFragment` at all.

**GREEN when:** That anchor shows `REQ-1108` plus its title, its `href` still points at
`https://example.com/spec`, the glossary lists REQ-1108 once with its own title and status, a second
`[REQ-1108](…)` in the same body stays bare, a renderer-produced autolink is still skipped, a dead id
inside an anchor is not flagged, running the linkifier twice over the same body produces identical
DOM, and — the folded clipboard half — a copied body containing `[REQ-100]: https://example.com/x`
keeps that definition line byte-identical, so every `[REQ-100]` reference in the paste still
resolves.

**Validation:** Inferred during capture, from a verified reviewer finding. The user chose to capture
rather than build it.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: Codex reviewer P2 on PR #169, thread discussion_r3864240041, verified against the code and the corpus.*
