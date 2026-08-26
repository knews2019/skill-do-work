---
id: REQ-377
title: 'Index cited ticket ids and let the filter box match them'
status: pending
created_at: 2026-08-26T13:24:45Z
user_request: UR-075
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-375]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-374, REQ-375]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/citations.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-filters.js
  - skills/do-work-board/tools/queue-kanban/citations_test.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Index Cited Ticket Ids And Let The Filter Box Match Them

## What

Typing `REQ-1679` into the board's filter box returns every card whose body cites it, not just the
card whose title happens to contain that text. The citation set is computed once on the Go side and
shipped per request, so the filter reads an index rather than re-scanning bodies in the browser.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-374 and REQ-375 make a cited id *readable*. They do not make it *findable*. "What else refers to
REQ-1679?" is the question a reader asks immediately after learning what REQ-1679 was, and today the
only way to answer it is to open cards one at a time. In the user's words: "titles are being added to
the display, they aren't being added to the index."

## Context

The filter is `web/board-filters.js` (173 lines): it populates the domain/status selects and holds
one predicate that decides whether a card is shown. The text box matches `id` and `title` only.

The board already has a precedent for exactly this shape — a build-time index the client reads
instead of re-deriving. `filementions.go` scans every body for repo-relative file paths, stats each
one, and ships a `repoFileMentions` map that `board-detail.js` consults to decide whether a path is a
link, a missing-file flag, or plain text. A `citedTicketIds` array per request is the same move: the
Go build knows every body, and the client should not re-scan them on every keystroke.

**The pattern lock-step is the trap here.** `bodyMentionPattern` (`web/board-detail.js:80`) and
`repoFileMentionPattern` (`filementions.go:17`) already carry a comment obliging them to stay
aligned, because a drift silently downgrades mentions to plain text rather than failing. A Go-side
ticket-id pattern is a **third** copy of that shape, and the same silent failure mode applies: an
under-matching index makes the filter quietly miss cards. Whatever the builder chooses, the agreement
must be pinned by a test that fails when either side alone changes, not by a comment.

Note the resolver semantics REQ-374 establishes and this REQ must match: a short `REQ-031` may name a
compound card id (`UR-002-REQ-031`), and an **ambiguous** segment shared by two cards resolves to
nothing and is never guessed.

## Detailed Requirements

- **Compute the citation set on the Go side**, once per build, from each REQ's and UR's body
  markdown. Ship it on the generated record — an array of resolved ticket ids, deduped, order
  irrelevant.
- **Resolve exactly as the display does**: compound-id first, then the short-segment index, and an
  ambiguous segment resolves to nothing. An id resolving to no record is **not** a citation and does
  not enter the index — that case is REQ-374's broken-reference flag, a different feature.
- **The filter box matches the citation index in addition to id and title.** Typing `REQ-1679`
  surfaces both the card whose id it is and every card citing it. Partial text keeps matching titles
  as it does today; the citation match is on a whole resolved id.
- **Say why a card matched.** A card surfaced only because it cites the searched id, with nothing in
  its own id or title matching, is otherwise indistinguishable from a title hit and reads as a bug.
  A minimal marker on the card — the builder chooses the form — is enough; do not build a second
  results panel.
- **Frontmatter references count as citations too**: `depends_on`, `related` and `addendum_to` name
  ids the reader means, and a card that depends on REQ-1679 should be findable by it.
- **The index costs no measurable build time.** It reads bodies the Go build has already loaded; do
  not re-read files from disk for it.

## Constraints

- **Do not add a third silently-drifting copy of the mention pattern.** Either share one definition
  across the Go and JS sides, or pin their agreement with a both-directions test that fails whichever
  side changes alone — the shape `_dev/primes/prime-kanban-board.md` records from REQ-248. A comment
  saying "keep these in sync" is what already failed twice.
- **Never guess.** An ambiguous segment is not a citation.
- **No new board write surface.** The index is generated output, not a file the tool writes into the
  queue; root `CLAUDE.md` § Kanban Board Write Surfaces stays untouched.
- Keep the filter's existing behaviour intact: this adds a matching channel, it does not change how
  id or title matching works.

## Dependencies

`depends_on: [REQ-375]`, which itself depends on REQ-374 — so the three run in the order
REQ-374 → REQ-375 → REQ-377.

**The edge exists to serialize a shared file, not because this REQ needs the other two's output.**
Its own work is independent: this REQ touches `board-filters.js`, `generate.go` and a new Go source
file, where they touch `board-core.js`, `board-detail.js`, `board-clipboard.js` and `board.css`. The
one file all three write is `generate_test.go`. An earlier draft declared that overlap in `write_set`
and stated the serial requirement in prose only — which enforces nothing: `write_set` is display-only
and "never a safety guarantee" (root `CLAUDE.md` § Glossary), so under `do-work run --fan-out` this
REQ and REQ-374 were both dependency roots and could have been dispatched concurrently into the same
test file. `depends_on` is the only field the work loop actually gates on, so the ordering is
declared there.

The alternative — dropping `generate_test.go` from this REQ's write set — was rejected: the filter
predicate is a pure function and belongs in the Node-harness lane that lives in that file.

## Builder Guidance

**Certainty level: Mixed.** The requirement is firm; two things are genuinely yours to decide.

First, **how the Go and JS mention patterns are kept in agreement** — one shared definition emitted
into the client, or two definitions plus an agreement test. Prefer whichever makes a drift *fail*
rather than degrade, and say why in a `## Decisions` entry.

Second, **what marks a citation-only match** on the card. Keep it small; the temptation is a whole
"referenced by" panel and that is not what was asked for.

Read `_dev/primes/prime-kanban-board.md` first. Two of its lessons bear directly here: REQ-248 on
pinning shared geometry with a both-directions agreement assertion, and REQ-289 on grepping the
*value* rather than the constant name when a token is read by two languages.

## Red-Green Proof

**RED prompt/case:** With a board generated from a tree where REQ-374's body cites `REQ-1679` and
REQ-374's own title contains no such text, type `REQ-1679` into the filter box. Today the REQ-374
card disappears — the predicate matches `id` and `title` only, so a card that cites the id is
filtered out exactly like one that never mentions it.

**Why RED now:** `web/board-filters.js` holds one predicate over `id` and `title`; nothing on the
generated request records what its body referenced, so there is no index for the predicate to
consult even if it wanted to.

**GREEN when:** Typing a whole resolved id into the filter box shows both the card that *is* that id
and every card that cites it in its body or in `depends_on` / `related` / `addendum_to`; a
citation-only match is visibly distinguishable from a title match; an ambiguous segment matches
nothing; and partial-text title matching behaves exactly as it does today.

**Validation:** User confirmed — raised verbatim as the "Search Limitation" and asked for directly.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: user request in session, 2026-08-26, raised mid-run on REQ-374 from UR-074's adjacent-improvements list.*
