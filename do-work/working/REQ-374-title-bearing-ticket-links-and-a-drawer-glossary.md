---
id: REQ-374
title: 'Title-bearing ticket links and a drawer glossary'
status: pending
created_at: 2026-08-26T13:02:24Z
user_request: UR-074
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-375, REQ-376]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-core.js
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/browser_probe_test.go
---

# Title-Bearing Ticket Links And A Drawer Glossary

## What

Every REQ and UR id the board's detail drawer renders carries its title. The first mention in a
body expands inline; later mentions stay bare. A glossary at the end of the body lists every id the
body referenced. Meta rows that are already reference lists always carry titles.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

A REQ body that reads `Read prime-cms.md, REQ-1679/REQ-1108 lessons` tells the reader nothing about
what those two requests were. The ids are already clickable links, so the information is one hop
away — but a hop per id, and you cannot skim a body that is half numbers. In the user's words:
"these numbers become too cryptic and unless I actually go and hunt them down I wouldn't know what
they are about but with the title I can get some idea."

## Context

The data is already there. `requestsById` and `userRequestsById` (`web/board.js:18-19`) carry a
`title` on every record, and `LoadBoard` walks `queue/`, `working/` **and** `archive/`
(`model.go:237`), so an id as old as REQ-1108 resolves. Nothing new is needed from Go.

Today's linkification lives in `web/board-detail.js:1-181`:

| Symbol | Line | Today |
|---|---|---|
| `requestIdByReqSegment` | 20 | indexes compound ids so a short `REQ-031` resolves; ambiguous → `null`, never guessed |
| `resolveTicketMention` | 32 | `requestsById` → `userRequestsById` → segment index |
| `makeTicketLink` | 48 | builds the anchor; **link text is the bare id, no title, no tooltip** |
| `bodyMentionPattern` | 80 | URL / repo file path / ticket id, in that priority |
| `buildLinkifiedFragment` | 92 | walks matches; already knows whether it is inside a code span |
| `linkifyDetailBody` | 150 | drops a duplicate H1, retargets autolinks, walks text nodes |

`makeTicketLink` has exactly five call sites, all in the same file, so one choke point serves the
whole change: the body linkifier (line 121), `makeDependencyDetailList` (line 234),
`makeTicketLinkList` (line 253 — serves `Blocked by`, `Unblocks`, `Overlapping write sets`), and
the `User request` meta row (line 299).

`board-core.js` is the shared-utility fragment and runs first in `boardJavaScriptFragmentPaths`
(`generate.go:43`). It already holds a ticket-lookup helper, `describeRequestStatus` (line 274), so
it is where a resolver shared with the clipboard belongs.

The adjacent precedent for a mention getting extra treatment at render time is the repo-file-mention
work (REQ-200 and REQ-207): the same "decide it on the client, from what the Go build already
verified" shape.

## Detailed Requirements

- **Move the resolver into `board-core.js`.** `requestIdByReqSegment` and `resolveTicketMention`
  move out of `board-detail.js:20-44` unchanged, including the ambiguous-segment `null` guard. Add
  `ticketTitleFor(detailKind, detailId)` returning the full title or `""`, and
  `shortTicketTitle(fullTitle)` returning the inline form. One definition, because REQ-375 consumes
  the same resolver and two copies would drift.
- **`shortTicketTitle` cuts on a word boundary at 60 characters** with an ellipsis. REQ-1685's own
  title is 88 characters; expanded in full mid-sentence it swamps the prose it sits in. The
  untruncated title always rides in the anchor's `title=` attribute and in the glossary, so nothing
  is lost.
- **`makeTicketLink` gains an expand flag.** When expanding, the anchor holds a mono
  `.ticket-link-id` span and a prose-font, one-step-dimmer `.ticket-link-title` span, and
  `anchor.title` carries the untruncated title. Extend the existing `.ticket-link` rules at
  `web/board.css:1722-1745`; do not replace them.
- **First mention expands, the rest stay bare.** `linkifyDetailBody` keeps a per-render set of
  already-expanded ids, reset on every drawer open so reopening a ticket looks identical.
- **Never expand inside `<code>`.** `buildLinkifiedFragment` already computes `insideCodeSpan`
  (`board-detail.js:107`) — reuse it. A backticked id keeps its bare mono link so a code run is
  never contaminated with prose.
- **Meta rows always carry the title**, independent of the body's first-mention set: they are
  reference lists, not prose. That covers `Depends on` (which keeps its existing status text),
  `Blocked by`, `Unblocks`, `Overlapping write sets` and `User request`.
- **Add a drawer glossary.** After linkification, append a `.detail-glossary` section listing every
  id the body resolved, each still a link: id, full untruncated title, and status from the existing
  `describeRequestStatus`. Omit the section entirely when the body referenced nothing.
- **Both id kinds, and both drawer kinds.** Everything above applies to UR ids exactly as it
  applies to REQ ids, and to the UR drawer's body exactly as to the REQ drawer's — both already
  route through `linkifyDetailBody` (`board-detail.js:418` and `:442`), so wiring the expansion
  or the glossary to REQ drawers alone would be a half-fix.
- **Unresolved ids stay bare** and stay out of the glossary — the same never-guess posture the
  ambiguous-segment rule already takes.

## Constraints

- **Never rewrite a REQ file.** The user was explicit: "don't change the rec file". Every expansion
  is computed at render time from board data.
- **No Go source change.** The payload already carries titles; this REQ touches the client and its
  tests only.
- **No new board write surface**, so the count in root `CLAUDE.md` § Kanban Board Write Surfaces is
  untouched.
- The expansion must not break the drawer's layout at its default width — an over-long expansion is
  the failure mode this REQ's truncation exists to prevent, and it is checked by eye, not only by
  assertion.

## Dependencies

None. REQ-375 depends on this one for the shared resolver.

## Builder Guidance

**Certainty level: Firm.** The four shape decisions (first-mention-only, 60-character cut, meta rows
always expanded, glossary present) were taken with the user during capture and are not open.

Two things are genuinely yours to decide: whether the title span sits inside the anchor or beside
it, and the exact glossary markup. Judge both by rendering the board and looking at it.

The mention-linkification path has **zero** direct test coverage today — `makeTicketLink`,
`resolveTicketMention`, `requestIdByReqSegment` and `bodyMentionPattern` appear in no `*_test.go`,
and the one adjacent test (`generate_test.go:905`) stubs `createTreeWalker` to return no nodes, so
the text-node loop never runs. Do not extend that gap. Reuse the Node-harness pattern from
`TestJavaScriptBehaviorDrawerHeadingDeduplication` for the pure functions and the existing Chromium
probe lane for the rendered result.

Per `_dev/primes/prime-kanban-board.md`: a rendering change's correctness is partly a claim about
pixels, so generate a board and look at it in **both** themes, and return `location.href` alongside
every browser measurement.

## Red-Green Proof

**RED prompt/case:** Open the drawer on a REQ whose body cites a resolvable id twice — for example
a body reading `Read REQ-1679 lessons` early and `the REQ-1679 note` later. Today both render as the
bare id `REQ-1679`, there is no title anywhere on the page, and no glossary exists.

**Why RED now:** `makeTicketLink` (`web/board-detail.js:48`) builds its anchor from
`linkText || detailId` and every caller passes only the raw mention text. No caller reads
`record.title`, and nothing appends a reference list to the body.

**GREEN when:** The first anchor contains a `.ticket-link-title` span carrying the request's title
(truncated on a word boundary at 60 characters, full text in the anchor's `title=`), the second
anchor for the same id does not, a backticked occurrence stays a bare mono link, the glossary
section lists the id exactly once with its full title and status, and an id with no matching board
record is neither expanded nor glossed.

**Validation:** User confirmed — the display shape was chosen from three presented options during
planning ("Visible expansion, first mention only"), and the user separately confirmed the glossary
and that the rule covers user requests as well as requests.

## Full Context

See `do-work/user-requests/UR-074/input.md` for complete verbatim input.

---
*Source: user request in session, 2026-08-26, prompted by REQ-1685's body citing REQ-1679/REQ-1108.*
