---
id: REQ-089
title: The board drawer's Copy button omits the ticket's frontmatter, so the paste carries no status, domain or timestamps
status: pending
created_at: 2026-08-03T22:49:36Z
user_request: UR-017
domain: frontend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
write_set: [tools/queue-kanban/model.go, tools/queue-kanban/generate.go, tools/queue-kanban/generate_test.go, tools/queue-kanban/web/board.js, docs/board-guide.md]
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

## Full Context

See `do-work/user-requests/UR-017/input.md` for the complete verbatim input, the screenshot
description, and the two capture clarifications.

---
*Source: "when I hit copy the status, domain, the header itself does not copy over [Image #1] was this already captured but not yet implemented, or do we need to use do-work capture-request to capture it now?"*
