---
id: REQ-383
title: '[impact-rule-change] Resolve ticket mentions in Go, and delete the client-side Markdown scanner'
status: pending
created_at: 2026-08-26T19:10:32Z
status_changed_at: 2026-08-26T19:41:12Z
user_request: UR-075
addendum_to: REQ-379
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-379, REQ-381, REQ-382]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/citations.go
  - skills/do-work-board/tools/queue-kanban/citations_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
  - _dev/primes/prime-kanban-board.md
---

# Resolve Ticket Mentions In Go, And Delete The Client-Side Markdown Scanner

## What

One goldmark AST walk on the Go side emits, per body, the byte positions where a ticket mention may
safely be annotated — already resolved to board records. The client splices titles in at those
offsets and stops knowing anything about Markdown. The hand-rolled fence scanner is deleted.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One goldmark walk in `citations.go` classifies every byte of a body as prose, code span, or code block, from the AST's own segments; mentions are matched against the shared pattern, resolved with REQ-378's semantics, and emitted as UTF-16 offsets into the whole clipboard document. `board-clipboard.js` keeps the titles and the appendix and loses the scanner.
- [x] **[APPLY]:** `citations.go` and `citations_test.go` added; `generate.go`, `web/board-clipboard.js`, `generate_test.go` and the prime changed. No file outside the write set touched.
- [x] **[UNIFY]:** Six files changed, each read back in full. `gofmt -l` clean, `go vet ./...` clean, `node --check` on the client fragment clean, no debug artifacts. 29 mutations applied one at a time, all killed, zero survivors. Whole-tree differential: 533 real Copy payloads compared old-vs-new, 2 differ, both from one cause and both fixes (see Decisions D3).

## Why

REQ-379 shipped a hand-rolled Markdown scanner in the browser client, and every external finding
against it since has been one shape — the scanner disagreeing with CommonMark:

| Finding | Cause |
|---|---|
| A fence inside a blockquote not recognised | scanner ≠ CommonMark |
| A backtick in a backtick fence's info string opened a fence | scanner ≠ CommonMark |
| A fence opened as a list item missed | scanner ≠ CommonMark |
| A code span crossing a line break split | scanner ≠ CommonMark |
| A four-space indented code block structurally invisible | scanner ≠ CommonMark |
| A link reference definition rewritten | scanner ≠ CommonMark |

Six symptoms, one cause. Fixing them individually — which is what this REQ originally proposed — pays
full price for each and leaves the next one waiting. **The parser was always available; it was just
on the wrong side.** `render.go:25` already builds a goldmark renderer and already parses every one
of these bodies at generate time to produce the drawer's HTML.

## Context

**The precedent is already in the code this REQ touches.** `filementions.go` exists because the
client cannot stat the filesystem: Go scans bodies for file paths, checks each, and ships
`repoFileMentions` (`generate.go:104`) so the client looks up an answer instead of deriving one. This
REQ is the same move for the same reason — the client cannot parse Markdown, so Go computes the
positions and ships them.

**The mechanism is proven, not assumed.** A probe against goldmark v1.8.2 over a body containing every
failing case returned exact byte ranges for all of them:

```
BLOCK FencedCodeBlock   35..59   "quoted REQ-200 verbatim\n"     <- blockquoted fence
BLOCK FencedCodeBlock   79..101  "depends_on: [REQ-300]\n"       <- list-item fence
SPAN  CodeSpan          75..95   "REQ-400 across\nlines"         <- span crossing a newline
BLOCK CodeBlock        113..136  "indented REQ-500 block\n"      <- indented block
```

and ```lang`invalid produced **no node at all**, because goldmark treats it as prose exactly as
CommonMark says.

**Two API constraints found by that probe, both load-bearing:**

- **`Lines()` panics on inline nodes** (`goldmark/ast/inline.go:38`). A `CodeSpan` has no line
  segments; its extent comes from walking its child `*ast.Text` nodes and taking min start / max
  stop. Two code paths, and a builder who assumes one will crash rather than misbehave.
- **Offsets are body-relative.** The clipboard payload is frontmatter + body, so every offset shifts
  by the fence length — which is exactly the `bodyStartOffset` `splitFrontmatter` already computes
  (`frontmatter.go:28-68`). That arithmetic is the one place this can silently corrupt a paste, so it
  gets its own test rather than riding along in another.

**This also ends the lock-step problem.** `bodyMentionPattern` (`web/board-detail.js`) and
`repoFileMentionPattern` (`filementions.go`) each carry a comment obliging them to stay identical, and
that obligation has failed twice. Once Go resolves the mentions there is one pattern, in one
language, with nothing to drift against.

## Detailed Requirements

- **A new `citations.go` walks each body's goldmark AST once** and records the byte ranges of every
  construct whose text is quoted: fenced code blocks, indented code blocks, code spans, and link
  reference definitions. Blockquote and list nesting need no special handling — the parser has
  already resolved containers.
- **Emit resolved mention occurrences, not raw ranges.** Per body: an ordered list of
  `{offset, length, kind, id}` for each mention that (a) falls outside every quoted range and (b)
  resolves to a board record. The client then needs no resolver, no pattern and no block reasoning.
- **Resolution semantics match REQ-378 exactly:** compound id first, then the short-segment index,
  and an **ambiguous** segment resolves to nothing and is never guessed.
- **Ship it beside the payload, never inside it.** A new field on the generated data; `board-markdown.js`
  keeps returning the file's exact bytes.
- **Rewrite `board-clipboard.js` to splice at the shipped offsets** — descending order, so earlier
  offsets stay valid as text grows — and **delete** `codeFenceRunFor`, `codeFenceRunCloses`,
  `findMatchingBacktickRun`, `stripContainerPrefix`, the paragraph lookahead and the cross-line span
  state. Roughly 150 lines go.
- **Behaviour is unchanged for everything REQ-379 already got right.** Its whole clipboard test file
  should pass against the new implementation with its assertions intact — that is the acceptance
  signal, and an assertion that needs rewriting means a behaviour changed that should not have.
- **Record the new division in `_dev/primes/prime-kanban-board.md`:** Markdown knowledge lives in Go;
  the client splices.

## Constraints

- **The two Go round-trip tests must pass unmodified.** `TestBuildGeneratedBoardMarkdownDataKeepsExactSources`
  and `TestBuildGeneratedBoardMarkdownDataRoundTripsTheWholeFile` pin that the payload is the file's
  bytes. If either needs editing, the positions have leaked into the payload and the approach is
  wrong. This is the guarantee that keeps a paste saveable as a valid REQ.
- **The drawer is out of scope.** `board-detail.js` gets its block context from the DOM and is
  correct; do not touch it. This REQ is the clipboard surface only.
- **Do not widen into REQ-381's or REQ-382's work.** The citation index and the Markdown-link
  handling consume this walk; they are separate REQs that depend on it.
- Every new assertion must be proven to bite **and** its fixture shown to reach the mutated line.
  Three vacuous mutations shipped in this batch; that is the failure this REQ must not repeat.

## Dependencies

None. REQ-379 is archived, and this replaces the scanner it shipped. **REQ-381 and REQ-382 now depend
on this REQ** — both consume the same walk, so it runs first.

## Builder Guidance

**Certainty level: Firm on the design, Mixed on the emitted shape.**

The walk, the two API constraints and the offset arithmetic are settled above and probe-verified.

What is yours: the exact shape of the emitted index, and how much it shares with REQ-381's citation
list. They are close but not identical — this REQ needs *positions of annotatable occurrences*, and
REQ-381 needs *the set of ids a body cites*, which includes ids inside quoted text. Decide whether one
emitted structure serves both or whether they are two projections of one walk, and write that down as
a `## Decisions` entry, because REQ-381 is built directly on the answer.

Read `_dev/primes/prime-kanban-board.md` first. REQ-248 (pin shared geometry with a both-directions
agreement assertion) applies to the Go↔JS offsets, and REQ-289 (grep the value, not the constant
name) applies to deleting the scanner cleanly.

## Red-Green Proof

**RED prompt/case:** Copy a ticket whose body contains a four-space indented code block mentioning
`REQ-1679`. The id is expanded, because a fence scanner has no fence line to find — the one failing
case the current implementation cannot reach at all. The same body copied after this REQ leaves it
byte-identical.

**Why RED now:** `codeFenceRunFor` matches fence *lines*; an indented code block has none, so no
amount of hardening the scanner sees it.

**GREEN when:** REQ-379's existing clipboard assertions all pass **unmodified**; an indented code
block, a blockquoted fence, a list-item fence, a code span crossing a newline and a link reference
definition each keep their exact text; ```lang`invalid is treated as prose and the ids under it
expand; the frontmatter fence of every document in a copy-all payload is byte-identical; and
`codeFenceRunFor`, `codeFenceRunCloses`, `findMatchingBacktickRun` and `stripContainerPrefix` no
longer exist in the generated page.

**Validation:** User confirmed — the design was put to the user with the probe output and they
directed this rewrite.

## Decisions

**D1 — The emitted shape is positions, not a citation set.** This REQ emits one entry per
annotatable OCCURRENCE — `{offset, length, kind, id, expand}` — because splicing needs a place, and
first-mention order needs occurrences. REQ-381 needs the opposite projection: the SET of ids a body
cites, including ids inside quoted text, with no positions at all. They are two projections of one
walk, not one structure serving both. REQ-381 should add its own field beside this one and reuse
`collectMentionSurfaces` and `ticketMentionResolver`; it must not widen these entries, because a
citation set has to include the fenced and code-span mentions that carry `expand: false` here AND
the ambiguous ones this index drops entirely.

**D2 — Two resolvers, pinned rather than merged.** Moving resolution to Go for the clipboard leaves
the drawer resolving in `board-core.js`, because the drawer works from a rendered DOM and has no
build step. Put to the user, who chose to pin: `TestJavaScriptBehaviorTicketMentionPatternAndResolverAgreeWithGo`
drives both the pattern and the resolver over one corpus in both directions, so whichever side
drifts alone fails. This is REQ-248's shape. Collapsing to one authority stays available later.

**D3 — Challenged and changed: the index ships INSIDE `board-markdown.js`, not beside it.** The
Detailed Requirements above say "ship it beside the payload, never inside it", to keep the payload
byte-exact. It is shipped as a sibling FIELD of the same `generatedBoardMarkdownData` struct
instead, which leaves both document maps untouched — the two round-trip tests pass unmodified — and
closes a skew window the split version opens: in serve mode `/board-markdown.js` re-walks the tree
on the Copy click, while `board-data.js` is whatever the page loaded with. A tree edit in between
would hand the client fresh text and stale offsets, and splice a title into the middle of a word.
Offsets and the text they measure now cannot arrive from different builds. Cost: `board-markdown.js`
grows 5% (6.50 MB → 6.84 MB on this repo's 375 REQs); the eager board payload is unchanged.

**D4 — An id inside a raw HTML block is no longer annotated.** Positive coverage by the parser's own
text nodes, rather than a list of quoted constructs, means a mention in an HTML block is skipped —
goldmark keeps no prose text there. This changed exactly one document in the tree (REQ-235, an id
inside an HTML comment) and it is a fix: the renderer emits `<!-- raw HTML omitted -->`, so the
drawer never showed that mention at all, and the old scanner was expanding a title into text no
reader sees.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-379's independent review (F2/F3/F5) plus five external findings, rewritten from "harden the scanner" to "delete it" on the user's direction.*
