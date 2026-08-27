---
id: REQ-381
title: 'Index cited ticket ids and let the filter box match them'
status: completed
created_at: 2026-08-26T13:24:45Z
user_request: UR-076
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-385]
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
- skills/do-work-board/tools/queue-kanban/serve.go
- skills/do-work-board/tools/queue-kanban/web/board-cards.js
claimed_at: '2026-08-27T12:11:18Z'
status_changed_at: '2026-08-27T12:11:18Z'
route: B
estimate:
  p50_active_minutes: 50
  confidence: medium
  basis:
  - Route B
  - 7-file write set
  - 2 subsystems involved
  - 8 acceptance criteria
  - browser evidence
  - performance instrumentation
  - cross-route regression gates
  - full-suite verification
  calculated_at: '2026-08-27T12:11:18Z'
completed_at: '2026-08-27T12:27:56Z'
commit:
kb_status: pending
---

# Index Cited Ticket Ids And Let The Filter Box Match Them

## What

Typing `REQ-1679` into the board's filter box returns every card whose body cites it, not just the
card whose title happens to contain that text. The citation set is computed once on the Go side and
shipped per request, so the filter reads an index rather than re-scanning bodies in the browser.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Extend the existing mention analysis with resolved unique citation IDs and frontmatter references; share that analysis across eager/raw projections in static/live generation. Test first, preserve raw offsets and existing filter channels, then add a small citation-only reason to REQ cards and UR headings and verify actual rendering.
- [x] **[APPLY]:** Built the seven declared sources; no board write surface, scanner, dependency, or panel added.
- [x] **[UNIFY]:** Reviewed every source diff and the file manifest below. gofmt, go vet, Node syntax and diff hygiene pass; parent independently reran the actual Chrome 141 acceptance test. No debug artifacts in source.

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

`depends_on: [REQ-385]`. REQ-383 is complete, so the edge that used to name it is gone; the edge now
names REQ-385, which rewrites the mention pattern this REQ builds its citation index on — so this one
is BOTH a shared-file edge (`citations.go`, `citations_test.go`) and a real content dependency.
It also shares `generate.go` with REQ-386 and `generate_test.go` with REQ-382 and REQ-387; all three
sit later in the chain.

**This edge serializes a shared file; it is not a need for the other's output.** `write_set` gates nothing — root `CLAUDE.md` § Glossary calls it "never a safety guarantee" — so only `depends_on` keeps two writers of one file apart under `do-work run --fan-out`. The whole batch is one chain, **REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387**, because `citations.go` alone is claimed by four of the six. That is ONE valid total order: reordering the queue means recomputing every edge, since a chain is only correct as a whole. `queue-kanban verify`'s `ungated-write-set-overlap` probe reports any pair this misses.

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

## Triage

**Route: B** — Known producer/filter seams require one shared analysis and small card presentation change; parent explored both static/live callers and existing filters.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/citations.go` (modify)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` (modify)
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/serve.go` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify)

**Acceptance criteria (restated from REQ):**
- All Detailed Requirements and the captured Red-Green Proof are satisfied.

## Exploration

The eager generator and live server each call buildGeneratedBoardData then buildGeneratedBoardMarkdownData on the same freshly built Board. Both projections must share the existing citation analysis at those two call sites; a second body scanner or repeated AST walk per build would violate the request. The raw-body parse remains distinct from HTML rendering because rendering preprocesses question-option hard breaks and changes offsets.

Both request cards (including Testing cards) and UR group headings are constructed in web/board-cards.js, so one additional source file covers the visible citation-only marker. The filter currently also matches a REQ's userRequestId; preserve that existing channel and do not label it citation-only. UR title hits currently bypass child search while retaining domain/status filters; keep that behavior, with an explanatory marker on citation-matched UR headings.

Frontmatter addendum_to is not stored on RequestTicket. The existing document split/parse/coercion helpers can expose the three requested reference fields without disk reads or a new ticket regexp. Body citation IDs should be resolved and deduped from the existing raw-source mention analysis, including quoted and authored-link-label references; unresolved and ambiguous IDs stay excluded. Clipboard annotation entries keep their current shape and exclusions.

Planned scope adds serve.go for sharing one analysis on live regeneration and web/board-cards.js for the visible match reason, beyond the five captured paths. Shared computation and its exact API remain a builder decision, but must not introduce long-lived mutable caches or a new scanner.

Search interpretation: keep existing partial own-id/title/userRequestId matches even when a short segment is ambiguous; never add a citation match for that ambiguous segment. Whole canonical IDs compare case-insensitively. If resolving a unique short compound alias for citation search, reuse resolveTicketMention rather than add an ID regexp. Clarify this in Decisions because the captured GREEN's shorthand about ambiguous matches cannot override preserving existing own-id partial matching.

## Implementation Summary

Citation search now reads eager, resolved, deduplicated `citedTicketIds` arrays on REQ and UR records. Static generation and live refresh share one raw-document analysis with the existing clipboard projection. No disk reads, mention regex, dependencies, panel or board write surface were added.

All seven modified source paths are under `skills/do-work-board/tools/queue-kanban/`:

- `skills/do-work-board/tools/queue-kanban/citations.go` (modified) — One existing body scan produces separate citation and annotation projections; existing YAML helpers resolve frontmatter references.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — Eager arrays and a short-lived shared analysis feed both generated payloads; raw wire shape is unchanged.
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified) — Cache rebuild uses the same shared analysis for eager and lazy data.
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` (modified) — Whole canonical IDs and unique short aliases search the eager set; exact ID priority and arbitrary text queries are safe.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified) — Citation-only reasons on REQ cards, Testing cards and UR headings, using existing badge CSS.
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modified) — Eager set, Go/JS wire agreement, ambiguity, aliases, frontmatter coercion, dependency alias, parser counter, and existing cursor guard.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — Static/live source agreement and refresh, actual parser-call count, real-browser reasons/layout, and updated isolated caller fixtures.

`git diff --stat`, every seven-file diff, `git diff --check`, `gofmt`, `go vet ./...`, and Node syntax checks were reviewed/run. No debug code, generated output, lifecycle, version or changelog edits are included in the source manifest. Ignored browser evidence lives under `build/req381-*`.

## Decisions

**D-01 — Separate sets and occurrences, one analysis (DECIDE & STATE).** `analyzeDocumentTicketMentions` collects the citation set during the existing pattern scan. The set includes resolved quoted mentions and authored-link labels, while Copy keeps its current exclusions, UTF-16 offsets and source bytes. URLs and repository paths still claim their runs through the existing mention pattern. Unknown and ambiguous references are excluded. Sorted IDs make the payload deterministic; empty sets serialize as arrays. Clipboard annotation fields are not widened.

**D-02 — Reuse frontmatter semantics (DECIDE & STATE).** The already loaded document is split and parsed with existing helpers. Dependency references use `resolveDependsOn`, including its legacy `dependencies` fallback; `related` and `addendum_to` use existing coercion. No additional filesystem access is needed. Synthetic URs receive an empty citation array and still have no raw-copy entry.

**D-03 — Keep the existing search channels (DECIDE & STATE).** Own ID, title and parent UR substring matches are preserved. Citation matches require a whole resolved ID, case-insensitively, or a uniquely resolving short compound alias. A canonical ID map normalizes casing before the existing resolver, preserving exact flat-ID precedence over compound aliases, including suffix letters. The map has no prototype because the query is arbitrary text. An ambiguous segment never produces a citation match, but may still produce the preexisting own-ID/title/parent substring hit. This clarifies the captured GREEN shorthand without breaking the preservation requirement.

**D-04 — Small reasons using existing CSS (DECIDE & STATE).** REQ/Testing cards say `cites REQ-1679` only when ordinary fields do not match. The button's accessible label includes that reason. UR headings use a compact `cites` badge with a full-ID title and accessible label, keeping the heading legible at 320px. A matched UR keeps its existing child-search bypass while domain/status filtering still applies. No second results panel is introduced.

**D-05 — Share only within a fresh build (DECIDE & STATE).** Both real callers construct one short-lived board analysis and send its results to eager and lazy projections. Existing standalone projection helpers remain available. No long-lived analysis cache is introduced. The integration test counts actual parser invocations: exactly two per nonempty fixture body (one HTML parse and one raw-source parse), including after lazy/eager live refresh. Its parser decorator is restored with test cleanup; this module has no parallel tests.

## Red-Green Proof

Before implementation, `go test -run 'TestGeneratedCitationIndex|TestJavaScriptBehaviorCitationIndex' -count=1 ./...` exited 1. The eager-record test reported missing `citedTicketIds`, including nil instead of an empty array. The JavaScript assertion failed with `false !== true`: “a body-only citation must match at first keystroke” for REQ-378 searching `req-1679`.

After implementation the same cases pass. Additional coverage verifies fenced and indented code, inline code, authored links, deduplication, all required frontmatter fields, ambiguous/unknown references, unique aliases, flat/compound suffix collisions, arbitrary `__proto__`/`constructor` queries, partial titles, preserved parent UR matching, and domain/status gates.

## Testing

All commands run from the queue-kanban module unless specified:

- `go test -run 'TestJavaScriptBehavior|TestCollectDocumentTicketMentions|TestGeneratedCitation|TestBuildGeneratedBoardMarkdown|TestGenerateSeparatesRawMarkdown|TestStaticAndLiveCitation|TestBodyTicketMention' -count=1 ./...` — exit 0, 32.395s. All existing Node behavior probes and focused raw/mention regressions pass.
- `go test -run 'TestStaticAndLiveCitation|TestGeneratedCitationIndex|TestJavaScriptBehaviorCitationIndex' -count=1 ./...` — exit 0 after final dependency helper change, 0.506s.
- `go test -run '^TestJavaScriptBehaviorCitationIndexAgreesWithGo$' -count=1 ./...` — exit 0 after null-prototype map correction, 0.446s.
- `go test -race -run '^TestStaticAndLiveCitationDataShareOneAnalysisAndRefreshTogether$' -count=1 ./...` — exit 0, 1.406s; the parser decorator and live refresh path have no reported races.
- `go vet ./...`, `node --check web/board-filters.js`, `node --check web/board-cards.js`, and repository `git diff --check` — exit 0.
- `TestBrowserBehaviorCitationSearchShowsReasonsAcrossViews` through the existing trusted CDP harness, Chrome 141 — exit 0. Captured-proof query entered using `Input.insertText`. Columns shows the target and citing card; Testing and By UR show their citation reasons before lazy Markdown has loaded. Light/dark × 320/768/1280 × three views gives 18 positive marker bounds, containment checks and empty mount/interaction error arrays. Ordinary title search and Clear remove stale markers. Programmatic focus proves focusability, not a Tab-navigation test.

Browser binary: `build/chrome-141/chrome-mac-arm64/Google Chrome for Testing.app/Contents/MacOS/Google Chrome for Testing`. Every measurement includes the actual page URL and browser version. No dump-dom transport was used.

Screenshots and full measurements: `build/req381-screenshots/`, with `acceptance.log` and 18 PNGs. A copied module under `build/req381-cdp-module/` adds capture-only code to the same browser test. Visually inspected `user-request-light-320.png`, `columns-dark-320.png`, `testing-light-768.png` and `columns-light-1280.png`: reasons are readable, contained, and distinguish the citation result from the target card. Fixture verify banners are expected from minimal in-memory records and are not browser errors.

Parent verification is recorded below.

## Prior-test impacts

The REQ-383 cursor guard now inspects `analyzeDocumentTicketMentions`, where the original cursor logic moved, retaining its two-read assertion. Existing By-UR caller tests include the new filter helper. The isolated done-card test supplies an empty `filterState`, matching the production environment. No prior behavioral assertion was removed, and raw-source/UTF-16 tests pass unchanged.

## Lessons Learned

- Preserve exact-record precedence before case-folded alias resolution. Sending the filter's lowercase needle directly to a case-sensitive exact resolver can incorrectly choose a compound alias.
- Search sets and annotatable occurrences have different exclusions. Derive them together, but collect citations before clipboard-only suppressions so later presentation changes cannot silently remove search results.
- A string-key lookup fed by arbitrary search text needs a null prototype or an own-property check; canonical ticket keys alone do not make the query safe.

## Discovered Tasks

None requiring an adjacent source change.

## Parent Verification

Parent inspected the seven source diffs and independently ran `TestBrowserBehaviorCitationSearchShowsReasonsAcrossViews` with the exact Chrome for Testing 141.0.7390.37 binary, exit 0 (0.851s). Viewed the 320px light By UR screenshot: marker and citing card are readable and contained. The full optional browser suite remains unverified; the default canonical gate skips it.

Performance: built the pre-change binary from 6887b518 and the final implementation binary separately. On the same repository tree, after one warmup each, seven interleaved pairs with alternating order measured median static generation 0.952309458s before and 0.956562333s after (+0.45%). Ranges were 0.901–1.015s and 0.913–1.013s. Concurrent canonical verification adds noise; this experiment detects no distinguishable build-time increase, rather than proving zero cost. Parser instrumentation independently verifies one raw analysis per loaded body and no new disk reads in the analysis.

## Qualification

The seven-file implementation matches the declared Scope and related ticket-chain seam. Raw source bytes/UTF-16 offsets and existing filters are preserved; the only earlier test adaptations supply new helper dependencies or follow the moved analysis, with no behavioral assertions removed. No unrelated implementation files changed.

## Orientation

The board can now answer “which cards cite this ticket?” before Copy data loads. Citation-only hits carry a small explanation, and static generation and live refresh use the same resolved index.

## Review

**Independent review: Pass, 98.75%.** Requirements 100%, code quality 95%, test adequacy 100%, scope 100%, risk Low. Reviewer independently ran the focused Go/JavaScript checks (exit 0, 1.269s) and exact Chrome 141 acceptance probe (exit 0, 0.765s), inspected the narrow render, and audited the real producer/consumer diff. No Important findings.

The restatement sweep found the shipped board guide still says “id or title”; this Minor prose-only omission was appended to the existing prose backlog. The frontmatter comment was clarified to “never annotated.” Qualification initially rejected an unnecessary success log in the new Node fixture; it was removed, assertions unchanged, then the focused wire test and qualification passed. No screen-reader, actual Tab-navigation, or full browser-suite pass is claimed.

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
