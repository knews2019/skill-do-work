---
id: REQ-089
title: The board drawer's Copy button omits the ticket's frontmatter, so the paste carries no status, domain or timestamps
status: completed
claimed_at: 2026-08-04T00:17:30Z
completed_at: 2026-08-04T00:23:03Z
kb_status: pending
created_at: 2026-08-03T22:49:36Z
user_request: UR-017
domain: frontend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
write_set: [tools/queue-kanban/model.go, tools/queue-kanban/generate.go, tools/queue-kanban/generate_test.go, tools/queue-kanban/serve_test.go, tools/queue-kanban/web/board.js, docs/board-guide.md]
---

# The board drawer's Copy button omits the ticket's frontmatter, so the paste carries no status, domain or timestamps

## What

Clicking **Copy** in the board's card drawer puts the `# REQ-NNNN: <title>` heading plus the rendered
body's source Markdown on the clipboard. None of the metadata the drawer displays directly above the
body comes with it — status, domain, user request, created, claimed, tree.

Make the Copy payload the **ticket file exactly as it exists on disk**: the verbatim frontmatter
fence followed by the verbatim body. A paste then round-trips — it can be saved straight back as a
valid REQ (or UR) file.

The heading line is *not* the problem. It was added in 0.163.3 and the user confirmed it is present in
their paste; this REQ is scoped to the metadata.

## Why

The drawer shows the frontmatter as labelled rows, so a copy that drops all of it is surprising — the
user expected the panel they are looking at to be what lands on the clipboard. Verbatim frontmatter
was chosen over a human-readable metadata list specifically because it round-trips back into a file.

## Context

The frontmatter never reaches the clipboard because the payload is built from the body alone:

- `tools/queue-kanban/model.go:122` — `BodyMarkdown string // raw Markdown body after the closing frontmatter fence`
- `tools/queue-kanban/generate.go:333-347` — `buildGeneratedBoardMarkdownData` projects only
  `ticket.BodyMarkdown` / `userRequest.BodyMarkdown` into the id-keyed payload written to
  `board-markdown.js` (lazily loaded on first Copy click, so the initial page doesn't download source text).
- `tools/queue-kanban/web/board.js:2233-2268` — the click handler; `2137-2143` `rawMarkdownForDetail`,
  `2148-2179` `copyHeadingForDetail` / `copyTextWithHeading`.

The frontmatter *is* parsed — into struct fields that build the drawer's `<dl>` (`openRequestDetail`,
`board.js:1644-1758`) — it is just dropped from the Markdown map.

This is a feature request, not a bug against spec: `docs/board-guide.md:50` currently promises only
"the source Markdown under a `# REQ-042: <title>` heading."

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `tools/queue-kanban/prime-do-kanban.md`, `crew-members/general.md`,
  `coding-guardrails.md`, `testing.md`, and the four cited source locations. Approach exactly as the
  REQ traced it: one new field populated by exact-suffix removal at both parse sites, one changed
  projection, one branch in the click handler. Nothing widened into metadata formatting.
- [x] **[APPLY]:** Code written as planned, plus two scope extensions declared before the write (D-02,
  D-03).
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file.
  - `model.go` — checked both struct fields carry the exact-suffix contract in a doc comment, both
    parse sites compute it identically, and `splitFrontmatter`'s signature is untouched as required.
  - `generate.go` — checked the projection is the only behavior change and the doc comment now says
    the map holds file text.
  - `generate_test.go` — checked both new tests parse real files from disk rather than hand-building
    structs, which is the only way the "original bytes" claim is actually tested.
  - `serve_test.go` — checked the repointed assertion and its traceability comment.
  - `web/board.js` — checked the primary path returns `rawMarkdown` untouched, `copyTextWithHeading`
    is now reachable only from the two fallback branches, and `node --check` passes.
  - `docs/board-guide.md` — checked the rewritten sentence names the three rows that cannot carry.
  - No debug artifacts: no `console.log`, no `fmt.Println`, no `TODO`. Binary not staged (gitignored).

## Detailed Requirements

1. **Primary path emits the original file bytes** — `<verbatim frontmatter fence>` + `<verbatim body>`,
   byte-identical to the file on disk.

2. **No synthesized heading and no H1 de-duplication on the primary path.** `id:` and `title:` in the
   fence carry identity, which is the entire rationale the 0.163.3 heading exists for
   (`web/board.js:2145-2147`). Verbatim must mean verbatim or it does not round-trip.

3. **Add one field, populated by exact-suffix removal.** Add `FrontmatterMarkdown string` beside
   `BodyMarkdown` (`model.go:122`) on `RequestTicket` and on the UR struct. Populate at both parse
   sites — `parseRequestTicket` (~`model.go:551`) and the UR parser (~`model.go:620`):

   ```go
   rawContent := string(contentBytes)
   frontmatterMarkdown := ""
   if hasFrontmatter {
       frontmatterMarkdown = rawContent[:len(rawContent)-len(bodyText)]
   }
   ```

   Reconstructing the fence from `yamlText` instead re-serializes it and loses the original bytes
   (CRLF line endings, trailing spaces, key order, comments). `splitFrontmatter`'s signature stays
   untouched.

4. **Change the projection, not the wire shape.** In `buildGeneratedBoardMarkdownData`
   (`generate.go:333-347`) use `ticket.FrontmatterMarkdown + ticket.BodyMarkdown`. The payload stays
   `map[string]string`, so no JS payload-shape change is needed and the lazy `board-markdown.js` split
   keeps working — it only grows by the frontmatter. Update the function's doc comment: the map now
   holds file text, not body text.

5. **Keep the degraded path identified.** In the click handler (`board.js:2233-2268`) use
   `rawMarkdown` verbatim when it is non-null, and keep `copyTextWithHeading` **only** on the
   `renderedTextFallback` branch (stale or missing `board-markdown.js`), which has no frontmatter text
   available. Retire the now-wrong comment at `board.js:2145-2147` and state which shape belongs to
   which path.

6. **Do not synthesize a fence on the fallback path** from `requestsById[id]`. A fabricated fence that
   looks verbatim is worse than a heading — it would paste as a file whose values were reassembled
   from display state.

7. **Applies to URs too.** The UR parser is a separate code path (`model.go:620`) and must get the
   same treatment.

8. **Docs and release.** Rewrite the Copy sentence at `docs/board-guide.md:50`. Then, per the
   pre-commit ritual: bump `actions/version.md` (0.168.1 → **0.169.0** — minor, a behavior change to a
   shipped feature) and add a `CHANGELOG.md` entry with a descriptive title (e.g. "Board Copy Carries
   the Ticket's Frontmatter"); verify the title is not already used.

## Constraints

- **Three drawer rows cannot come through, by construction.** They are derived at parse/bucket time
  and documented as "never read from frontmatter" at `model.go:105-116`: **Tree** (`TreeSection`,
  inferred from which directory the file sits in), **Overlapping write sets** (`WriteSetOverlaps`,
  computed by `annotateWriteSetOverlap` after bucketing), and **Unblocks** (`Dependents`). Relative
  times ("6h ago") likewise won't carry — the raw stamps will. This is the accepted cost of the
  verbatim shape; TREE was one of the six rows in the user's screenshot, so do not treat its absence
  as an incomplete implementation.
- **The tool stays read-only here.** This adds no write surface, so the "exactly two write surfaces"
  rule for `tools/queue-kanban/` needs no amendment.
- **Do not commit the compiled binary** — `tools/queue-kanban/queue-kanban` is gitignored.

## Dependencies

None. Note that pending **REQ-087** also declares `tools/queue-kanban/web/board.js` in its
`write_set`, so the board will badge these two as overlapping. That is display-only — nothing
schedules on it — but the overlap is real. REQ-087 touches `board.js:154` and `board.js:553`
(timestamp command strings); this REQ touches the copy handler around `board.js:2145-2268`. Different
regions, so a merge conflict is unlikely.

## Builder Guidance

**Certainty: Firm.** The payload shape was confirmed with the user during capture, and the
implementation surface was traced before capture — one new field, one changed projection, one branch
in the click handler. Keep it that small; resist widening into a metadata-formatting feature.

## Red-Green Proof

**RED prompt/case:** In `tools/queue-kanban/generate_test.go`, build a board from a fixture REQ file
that has frontmatter, call `buildGeneratedBoardMarkdownData`, and assert the map entry for that REQ id
equals the fixture file's full bytes. Today it equals only the body — the assertion fails on the
missing `---\nid: …\n---\n\n` prefix.

**Why RED now:** `buildGeneratedBoardMarkdownData` (`generate.go:344`) assigns
`markdownData.Requests[ticket.RequestId] = ticket.BodyMarkdown`, and `BodyMarkdown` is by definition
the text *after* the closing fence (`model.go:122`). No field on the struct holds the frontmatter text
at all, so nothing in the process can currently produce it.

**GREEN when:** that assertion passes for both a REQ and a UR fixture, and — end to end — running
`do-work board`, opening a card drawer, clicking **Copy** and pasting into a scratch file produces a
file that `diff`s clean against the REQ file on disk.

**Validation:** User confirmed (payload shape chosen from presented options during capture; the
heading-line question confirmed the scope is metadata-only).

## Verification

1. `cd tools/queue-kanban && go test ./...`
2. `do-work board` in this repo → open any REQ drawer → **Copy** → paste into a scratch file →
   `diff` against the REQ file on disk; expect byte-identical.
3. Repeat on a UR drawer (separate parser).
4. Degraded path: load the board with `board-markdown.js` absent and confirm the paste still leads
   with `# REQ-NNN: <title>`.

## Assets

`do-work/user-requests/UR-017/assets/REQ-089-board-drawer-copy-metadata.png` — the board drawer for
REQ-1322 in the `g1w-game-find-the-difference` project, showing the six metadata rows (STATUS,
DOMAIN, USER REQUEST, CREATED, CLAIMED, TREE) that the Copy button does not carry. Fully described in
`do-work/user-requests/UR-017/input.md`.

## Triage

**Route: B** - Medium

**Reasoning:** The change is small and fully traced in the REQ (field, projection, handler branch), but
it spans five declared files across Go, JS and docs, and the "verbatim bytes" requirement makes *how*
the field is populated the whole correctness question. Worth an explicit scope declaration; not worth a
plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**The parse sites are symmetric but separate.** `parseRequestTicket` (`model.go:545`) and
`parseUserRequestTicket` (`model.go:631`) each call `splitFrontmatter(string(contentBytes))` and keep
only `yamlText` (parsed into fields) and `bodyText`. Neither retains the raw text, which is why
requirement 3 specifies exact-suffix removal — `rawFileText[:len(rawFileText)-len(bodyText)]` — rather
than re-serializing from `yamlText`.

**One projection point serves two surfaces.** `buildGeneratedBoardMarkdownData` is called by
`generateStaticSite` (`generate.go:193`) **and** by the live server (`serve.go:322`). The REQ's Context
named only the static path, so this is worth stating: changing the projection fixes `do-work board`
serve mode at the same time, with no second edit. It also means the live board's contract test
(`serve_test.go`) asserts against this function's output — see D-03.

**The heading machinery has exactly one caller.** `copyTextWithHeading` is invoked once, in the click
handler's `.then`, which today applies it to *both* branches. Requirement 5's split is therefore a
matter of moving the call into each branch rather than restructuring anything.

## Scope

**Files I will touch:**
- `tools/queue-kanban/model.go` (modify) — `FrontmatterMarkdown` on both structs, populated at both parse sites
- `tools/queue-kanban/generate.go` (modify) — the projection and its doc comment
- `tools/queue-kanban/generate_test.go` (modify) — the round-trip RED test and a fence-less case
- `tools/queue-kanban/web/board.js` (modify) — the click handler's branch split and the stale comment
- `docs/board-guide.md` (modify) — the Copy sentence
- `tools/queue-kanban/serve_test.go` (modify) — **added during the build (D-03)**: its assertion pins
  the pre-change payload shape of the shared projection

**Files I will NOT touch:**
- `tools/queue-kanban/serve.go` — it already calls the shared projection; no change needed
- `splitFrontmatter`'s signature — explicitly protected by requirement 3
- The drawer's `<dl>` rendering (`board.js:1644-1758`) — this REQ changes the clipboard, not the panel

**Acceptance criteria (restated from REQ):**
- [ ] The primary payload is the original file bytes, fence + body (req 1)
- [ ] No synthesized heading and no H1 de-duplication on the primary path (req 2)
- [ ] `FrontmatterMarkdown` added and populated by exact-suffix removal at both parse sites (req 3)
- [ ] Projection changed; wire shape still `map[string]string`; doc comment updated (req 4)
- [ ] The degraded path keeps the heading and is identified in a comment (req 5)
- [ ] No fence is synthesized on the fallback path (req 6)
- [ ] URs get the same treatment (req 7)
- [ ] `docs/board-guide.md` rewritten; version bumped minor; changelog entry added (req 8)

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model.go` (modified)
- `tools/queue-kanban/generate.go` (modified)
- `tools/queue-kanban/generate_test.go` (modified)
- `tools/queue-kanban/serve_test.go` (modified)
- `tools/queue-kanban/web/board.js` (modified)
- `docs/board-guide.md` (modified)

**What was done:** `RequestTicket` and `UserRequestTicket` gained `FrontmatterMarkdown`, populated at
both parse sites by removing the body as an exact suffix of the raw file text — so the fence keeps its
original key order, comments, spacing and line endings. `buildGeneratedBoardMarkdownData` now projects
`FrontmatterMarkdown + BodyMarkdown`, which fixes the static bundle and the live server together
because both call it. The click handler returns that payload verbatim on the primary path and applies
`copyTextWithHeading` only on the two fallback branches; the stale comment claiming "the stored
Markdown is the file body only" was replaced with one that names both shapes and says which path owns
each. The guide's Copy sentence was rewritten, including the three drawer rows that cannot carry.

## Decisions

- **D-01 (DECIDE & STATE)** — *The fallback heading call moved into both fallback branches rather than
  staying in a shared `.then`.* Requirement 5 asks that `copyTextWithHeading` apply only to the
  degraded path. The alternative — a flag threaded through the shared `.then` — would have kept one
  call site but made "which shape am I" a runtime question. Two direct calls make the branch/shape
  pairing readable at each return, which is what the replacement comment documents.
- **D-02 (DECIDE & STATE)** — *A fence-less-file test was added beyond the REQ's Red-Green Proof.*
  `frontmatterMarkdown` stays `""` when `hasFrontmatter` is false, so the concatenation degrades to the
  old behavior — but that is the branch where a naive `rawFileText[:len-len(body)]` outside the `if`
  would silently produce garbage. Cheap to pin, and it is the only path where the suffix arithmetic
  could be wrong.
- **D-03 (DECIDE & STATE, scope extension)** — *`tools/queue-kanban/serve_test.go` added to scope.*
  `TestServeLazyMarkdownEndpointReturnsExactSources` asserts the served payload equals the body only,
  and the live server shares the projection this REQ changes — so the REQ's own change breaks it by
  construction. This is a genuine discovery: the REQ's Context named only the static generator, and the
  test failure is what revealed the live board is fixed by the same edit. Repointed, intent preserved
  (it still asserts "exact sources"; the definition of exact widened to the whole file). Cross-REQ
  impact recorded under `## Testing`.

## Qualification

Passed — 6 files verified, 8 requirements traced, P-A-U confirmed.

- `tools/checks/qualify.sh` → OK. `tools/checks/scope-drift.sh` → OK, both set-differences empty.
- **Requirements traced:** 1 → proved end to end below, not just asserted. 2 → the primary branch
  returns `rawMarkdown` with no call to `copyTextWithHeading`; verified by reading and by the payload
  starting with `---`. 3 → both structs, both parse sites, exact-suffix form as specified;
  `splitFrontmatter`'s signature unchanged. 4 → `generatedBoardMarkdownData` is still
  `map[string]string`; only the value changed; doc comment rewritten. 5 → both fallback branches call
  `copyTextWithHeading`, and the replacement comment states which shape belongs to which path. 6 →
  no fence is constructed anywhere in `board.js`; `requestsById` is read only for the heading's title,
  exactly as before. 7 → the UR parser has the same population and the UR round-trip is proved below.
  8 → guide sentence rewritten, version bumped **minor** (0.168.6 → 0.169.0, a behavior change to a
  shipped feature), changelog title checked for reuse.
- **Substantive:** real logic in three source files, two new tests parsing real files, one repointed
  assertion, one rewritten user-facing sentence.
- **Flowing:** not hollow — the payload was extracted from the actually-generated `board-markdown.js`
  and diffed against disk (below), rather than trusting the unit tests.
- **Contamination check (Step 10):** the previous REQ in this session (REQ-087, via REQ-085's fan-out)
  touched `model.go` and `web/board.js`. The overlap is expected and unrelated: REQ-087 changed the
  future-timestamp warning string in `model.go` and two timestamp display strings in `board.js`, both
  far from this REQ's parse-site and copy-handler regions. Both REQs' tests pass together.

## Testing

**Tests run:** `cd tools/queue-kanban && go test ./...` (+ `go vet`, `gofmt`, `node --check web/board.js`)
**Result:** ✓ All passing

**Red-green validation:** traced to `## Red-Green Proof`, whose RED case was implemented as specified
(build a board from a fixture REQ file with frontmatter, call `buildGeneratedBoardMarkdownData`, assert
the map entry equals the file's full bytes).

- `TestBuildGeneratedBoardMarkdownDataRoundTripsTheWholeFile`: ✗ before implementation — got
  `"\n## What\n\n- [ ] keep formatting\n"`, want the same text prefixed with the whole
  `---\nid: REQ-4242\n…---\n` fence; the UR assertion failed identically → ✓ after. The fixture
  deliberately includes a comment line and an irregular `domain:   general` spacing, so a
  re-serialized fence would fail it even though a parsed-field comparison would pass.
- `TestBuildGeneratedBoardMarkdownDataHandlesAFenceLessFile`: ✓ — pins the empty-fence branch (D-02).

The RED was observed, not assumed: the struct fields were added first with no population, so both
assertions failed on content rather than on a build error.

**New tests added:**
- `TestBuildGeneratedBoardMarkdownDataRoundTripsTheWholeFile` — REQ and UR, parsed from real files
- `TestBuildGeneratedBoardMarkdownDataHandlesAFenceLessFile` — the empty-fence branch

**Existing tests updated (cross-REQ impact):** one, intentional.
`TestServeLazyMarkdownEndpointReturnsExactSources` (`serve_test.go`, from the REQ that introduced the
lazy Markdown endpoint) asserted the served payload for `REQ-0001` equals `"\n# REQ-0001\n"` — the body
alone. The live server calls the same `buildGeneratedBoardMarkdownData` this REQ changes, so widening
the payload changes the served bytes too. That is the intended behavior, not a side effect: the Copy
button on a live board should carry the same file text as on a generated one. The assertion was
**repointed, not deleted** — it now pins the full file text, so its original intent ("exact sources")
is asserted more strictly than before, with a comment tracing the change to REQ-089.
`TestBuildGeneratedBoardMarkdownDataKeepsExactSources` needed no change: its hand-built fixtures leave
`FrontmatterMarkdown` empty, so it still validates that nothing mangles the body.

**End-to-end acceptance (the REQ's GREEN condition, Verification steps 1–3):** ran
`queue-kanban generate --repo-root . --out /tmp/board-out`, parsed the real `board-markdown.js` the
board actually ships, extracted the payloads and diffed them against the files on disk:

```
REQ-089 → diff /tmp/copied-req-089.md do-work/…/REQ-089-….md   → byte-identical ✓
UR-017  → diff /tmp/copied-ur-017.md  do-work/user-requests/UR-017/input.md → byte-identical ✓
```

The REQ payload's first lines are the real fence (`---`, `id: REQ-089`, `title: …`, `status: claimed`,
`claimed_at: …`, `created_at: …`), confirming requirement 1 and requirement 2 together — the fence is
present and no `# REQ-089:` heading was prepended. Verification step 4 (the degraded path) was checked
by reading rather than by loading a stale bundle in a browser: both fallback branches call
`copyTextWithHeading`, which is unchanged, and no code path constructs a fence.

*Verified by work action*

## Lessons Learned

**What worked:**
- **Testing the round-trip from parsed files, not from hand-built structs.** The struct-level test that
  already existed passes both before and after the change, because a zero-value fence concatenates to
  nothing. Only a test that reads a real file can tell "the original bytes survived" from "a fence was
  produced" — and the fixture's comment line and irregular spacing are what make it discriminating.
- **Diffing the actually-generated bundle.** The unit tests prove the projection; extracting the
  payload from `board-markdown.js` and `diff`ing it against the file on disk proves the whole chain,
  including the JSON encoding hop that no Go test covers.

**What didn't:**
- Nothing failed outright, but the REQ's file list was one short. `serve.go` shares the projection, so
  its contract test broke — discovered by running the suite, not by reading the REQ's Context, which
  named only the static generator.

**Worth knowing:**
- **`buildGeneratedBoardMarkdownData` is the single projection point for both the static bundle and the
  live server** (`generate.go:193`, `serve.go:322`). Changing it changes `do-work board` in both modes
  at once. Convenient here; worth knowing before assuming a change is static-only.
- **Exact-suffix removal must stay inside the `hasFrontmatter` guard.** `splitFrontmatter` returns the
  whole text as `bodyText` when there is no fence, so the arithmetic yields `""` correctly — but only
  because the guard skips it. Moving the slice outside the `if` would still compute `""` today and
  become wrong the moment `splitFrontmatter`'s no-fence contract changes.
- **Three drawer rows can never round-trip**, by construction: `TreeSection`, `WriteSetOverlaps` and
  `Dependents` are derived at parse/bucket time and are documented as never read from frontmatter. The
  guide now says so, so their absence reads as designed rather than as an incomplete implementation.

## Orientation

**Now the board drawer's Copy button puts the ticket file itself on the clipboard** — frontmatter fence
and all — so a paste can be saved straight back as a valid REQ or UR instead of arriving as an
anonymous body with none of the status, domain or timestamps the drawer was displaying right above it.
Works the same on a generated board and a live one. Lives in the queue-kanban board tool
(`tools/queue-kanban/prime-do-kanban.md`).

`[MAP CHANGED]` — a new contract on a parsed field: `FrontmatterMarkdown` holds the **original fence
bytes** and is populated by exact-suffix removal, never re-serialized from the parsed fields. Anything
that later re-derives, normalizes, or rewrites that field breaks the round-trip guarantee this REQ
exists to provide, and the lazy `board-markdown.js` payload is now file text rather than body text.
Prime staleness spot-check: `tools/queue-kanban/prime-do-kanban.md`'s referenced paths all still exist;
its `## Lessons` gains this REQ's entry at Step 8.

## Review

**Overall: 96%** | 2026-08-04T00:22:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — the payload extracted from the actually-generated `board-markdown.js` diffs
byte-identical against both a REQ file and a UR file on disk.
**Suggested testing:** 2 items
**Follow-ups created:** None

**Restatement Sweep (MUST — run):** The diff redefines what the Copy payload *contains* (body text →
file text) and adds a field contract (`FrontmatterMarkdown` holds original bytes, never re-serialized).

- *Prose consumers:* `docs/board-guide.md:50` was the only place describing the Copy payload, and it is
  rewritten in this REQ — including the three rows that cannot carry, which pre-empts the obvious
  "the copy is incomplete" bug report. No other `.md` outside `do-work/` describes it.
- *Code consumers of the payload:* `buildGeneratedBoardMarkdownData` has exactly two callers,
  `generate.go:193` and `serve.go:322`. Both were swept; the second is what surfaced the `serve_test.go`
  break, which is the sweep working rather than a miss.
- *The wire shape:* unchanged (`map[string]string`), so `board.js`'s `rawMarkdownForDetail` and the
  lazy-load split needed no change — verified by reading, and by `node --check`.
- *`tools/queue-kanban/prime-do-kanban.md`:* greps clean for `BodyMarkdown` / copy / markdown — the
  prime never restated this contract, so nothing there is stale. Its `## Lessons` gains an entry at
  Step 8 (inline-only marker honored).
- *`CHANGELOG.md` history:* the 0.163.3 and lazy-payload entries describe what those releases did and
  remain accurate as history. Not restatements to update.

### Findings

**Minor — the fallback path's shape is asserted only by reading.** Requirements 5 and 6 are verified by
inspection (`copyTextWithHeading` is reachable only from the two fallback branches; nothing constructs
a fence) and by `node --check`, but there is no automated test for the degraded path — the repo has no
JS test harness, and adding one for this would be well outside the REQ. The risk is bounded: the
fallback code itself is unchanged, only its call sites moved. Worth knowing that the *only* proof the
heading still appears on a stale bundle is a human reading two branches.

### Notes on the dimensions

- **Requirements 100%** — all eight delivered. Requirement 2 deserves note: it is a *negative*
  requirement (no synthesized heading, no H1 de-dup on the primary path) and negatives are the easy
  ones to claim without evidence. It is proved positively here — the extracted payload's first line is
  `---`, not `# REQ-089:`, in the real generated bundle.
- **Code Quality 95%** — the field contract is documented where a later maintainer meets it (on both
  struct fields, not just one), and the replacement comment in `board.js` explains the two shapes and
  which path owns each, which is what requirement 5 actually asked for. Deducted slightly because
  `FrontmatterMarkdown + BodyMarkdown` is now computed in the projection rather than exposed as a
  single accessor on the ticket — fine at two call sites, but a third consumer wanting "the file text"
  would be tempted to re-concatenate rather than reuse.
- **Test Adequacy 96%** — the discriminating choice was parsing real files instead of hand-building
  structs: the pre-existing struct-level test passes both before *and* after the change, so it could
  never have caught this. The fixture's comment line and irregular `domain:   general` spacing mean a
  re-serialized fence fails even though a parsed-field comparison would pass — that is the actual
  requirement-3 risk, pinned. Deducted for the untested JS branch above.
- **Scope 100%** — `scope-drift.sh` clean. One file beyond the captured `write_set` (`serve_test.go`),
  declared before the write and logged as D-03, and it was a genuine discovery rather than drift: the
  REQ's Context named only the static generator.
- **Risk Low, not None** — this changes what a shipped, user-facing button puts on the clipboard, and
  the payload grew: `board-markdown.js` now carries every ticket's fence as well as its body. That file
  is lazily loaded on first Copy click, so the initial page weight is unaffected (the property the lazy
  split exists to protect), and a fence is a few hundred bytes against a body measured in kilobytes.
  Behavior is otherwise additive — nothing that previously copied stops copying.

### Suggested additional testing

- **A stale-bundle load in a real browser.** Generate a board, delete `board-markdown.js`, open it and
  click Copy: the paste should lead with `# REQ-NNN: <title>` and contain no `---` fence. This is the
  one requirement pair verified only by reading.
- **A CRLF or trailing-whitespace ticket file.** The exact-suffix approach was chosen precisely to
  preserve those, and the fixture covers a comment and irregular spacing but not line endings. A file
  written on Windows would be the honest test of the stated rationale.

*Reviewed by review-work action*

## Full Context

See `do-work/user-requests/UR-017/input.md` for the complete verbatim input, the screenshot
description, and the two capture clarifications.

---
*Source: "when I hit copy the status, domain, the header itself does not copy over [Image #1] was this already captured but not yet implemented, or do we need to use do-work capture-request to capture it now?"*
