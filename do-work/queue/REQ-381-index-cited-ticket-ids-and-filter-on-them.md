---
id: REQ-381
title: 'Index cited ticket ids and let the filter box match them'
status: pending
created_at: 2026-08-26T13:24:45Z
user_request: UR-076
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-383]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-378, REQ-379, REQ-383]
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

REQ-378 and REQ-379 make a cited id *readable*. They do not make it *findable*. "What else refers to
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

**REQ-383 now builds the walk this REQ consumes, which is why the dependency was re-pointed.** It
adds `citations.go` and one goldmark AST pass per body, resolving ticket mentions against board
records with REQ-378's exact semantics. This REQ takes the set of ids a body cites from that pass
rather than adding a second scanner. Read REQ-383's `## Decisions` entry before starting: it records
whether one emitted structure serves both needs or whether they are two projections of one walk, and
that answer is the seam this REQ builds on. The two are close but not identical — REQ-383 needs the
*positions of annotatable occurrences*, this REQ needs *every id a body cites*, including ids inside
quoted text, since a reference in a fenced example is still a reference for search.

**The pattern lock-step trap is narrowed and pinned — reuse it, do not re-open it.** REQ-383 left
exactly one definition of the syntax in Go: `bodyTicketMentionPattern` (`citations.go`) composes
`repoFileMentionPattern` (`filementions.go`) rather than restating it, pinned by
`TestBodyTicketMentionPatternComposesTheOneFilePathDefinition`. The drawer's `bodyMentionPattern`
(`web/board-detail.js`) survives because a browser cannot import a Go variable, and it is held by
`TestJavaScriptBehaviorTicketMentionPatternAndResolverAgreeWithGo`, a both-directions test over one
corpus. Take this REQ's ids from that walk. If the walk does not yet expose what this REQ needs,
extend it there rather than re-scanning here.

Note the resolver semantics REQ-378 establishes and this REQ must match: a short `REQ-031` may name a
compound card id (`UR-002-REQ-031`), and an **ambiguous** segment shared by two cards resolves to
nothing and is never guessed.

## Detailed Requirements

- **Compute the citation set on the Go side**, once per build, from each REQ's and UR's body
  markdown. Ship it on the generated record — an array of resolved ticket ids, deduped, order
  irrelevant. That is the EAGER `board-data.js` record, not beside REQ-383's `RequestMentions`:
  those ride in `board-markdown.js`, which is loaded only on a Copy click and which
  `TestGenerateSeparatesRawMarkdownForLazyCopy` keeps out of the initial paint. The filter runs at
  first keystroke, so its data has to be there already.
- **Reuse REQ-383's walk rather than repeating it.** `collectMentionSurfaces` and
  `ticketMentionResolver` (`skills/do-work-board/tools/queue-kanban/citations.go`) already classify
  every byte of a body and resolve every id; this REQ is a second projection of that one walk, per
  REQ-383's `## Decisions` D1. It must NOT widen REQ-383's per-occurrence entries — a citation set
  includes the quoted mentions those entries carry as `expand: false` AND the ambiguous ones they
  drop entirely.
- **Resolve exactly as the display does**: compound-id first, then the short-segment index, and an
  ambiguous segment resolves to nothing. An id resolving to no record is **not** a citation and does
  not enter the index — that case is REQ-378's broken-reference flag, a different feature.
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

- **Do not add another copy of the mention pattern.** REQ-383's `## Decisions` D2 and D5 settled
  where the definitions live and which tests hold them. Reuse `bodyTicketMentionPattern`; do not
  write a new regexp for ticket ids or repo paths, and do not add a second test covering the pair
  those two already cover. If this REQ introduces a NEW thing read by both languages — say, the
  citation array the filter matches on — that one needs its own both-directions pin, per REQ-248;
  the rule is one pin per shared thing, not one pin in total.
- **Never guess.** An ambiguous segment is not a citation.
- **No new board write surface.** The index is generated output, not a file the tool writes into the
  queue; root `CLAUDE.md` § Kanban Board Write Surfaces stays untouched.
- Keep the filter's existing behaviour intact: this adds a matching channel, it does not change how
  id or title matching works.

## Dependencies

`depends_on: [REQ-383]`.

**This edge is a real dependency, not a file-serialization gate.** REQ-383 builds `citations.go` and
the goldmark AST pass that resolves ticket mentions against board records; this REQ consumes that
pass for its citation index rather than adding a second scanner. Building this first would mean
writing a mention scanner that REQ-383 then deletes.

That is a change from the earlier ordering, which had this REQ depending on REQ-379 purely because
both wrote `generate_test.go`. That serialization still holds and still matters — `write_set` is
display-only and "never a safety guarantee" (root `CLAUDE.md` § Glossary), so a shared file must be
gated in `depends_on` or `--fan-out` will race it — but it is now the weaker of the two reasons.

REQ-382 depends on this REQ in turn, so the chain is REQ-383 → REQ-381 → REQ-382.

## Builder Guidance

**Certainty level: Mixed.** The requirement is firm; one thing is genuinely yours to decide.

**What marks a citation-only match** on the card. Keep it small; the temptation is a whole
"referenced by" panel and that is not what was asked for.

The pattern-agreement question that used to sit here was answered by REQ-383 (`## Decisions` D2) and
is no longer open.

Read `_dev/primes/prime-kanban-board.md` first. Two of its lessons bear directly here: REQ-248 on
pinning shared geometry with a both-directions agreement assertion, and REQ-289 on grepping the
*value* rather than the constant name when a token is read by two languages.

## Red-Green Proof

**RED prompt/case:** With a board generated from a tree where REQ-378's body cites `REQ-1679` and
REQ-378's own title contains no such text, type `REQ-1679` into the filter box. Today the REQ-378
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

See `do-work/user-requests/UR-076/input.md` for complete verbatim input.

---
*Source: user request in session, 2026-08-26, raised mid-run on REQ-378 from UR-075's adjacent-improvements list.*
