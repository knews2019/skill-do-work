---
id: REQ-374
title: 'Title-bearing ticket links and a drawer glossary'
status: claimed
created_at: 2026-08-26T13:02:24Z
claimed_at: 2026-08-26T13:20:45Z
user_request: UR-074
domain: frontend
route: B
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-26T13:21:43Z
  basis:
    - Route B
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - browser evidence
    - full-suite verification
related: [REQ-375, REQ-376]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-core.js
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
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
- **An unresolved id is flagged as broken, not left as prose.** An id matching the mention
  pattern that resolves to no board record renders as a non-link `span` in the blocked accent
  with a `title="Not found in this queue"` tooltip — the exact treatment `.repo-file-missing`
  (`web/board.css:1740`) already gives a file path the Go build could not find. A dead
  cross-reference is a typo or a reference to never-captured work, and today it is invisible
  until someone hunts for it. It carries no title (there is none) and gets no glossary line.
  **The never-guess rule is unchanged** — an ambiguous REQ segment (two cards sharing it)
  resolves to nothing and must stay plain prose, never flagged: ambiguous is not missing.

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

## Decisions

- **D-01 — Reverse the unresolved-id rule (orchestrator, pre-dispatch):** the captured REQ said an
  unresolved id "stays bare and stays out of the glossary". The user directed mid-run, before any
  code was written, that it show as broken instead — red with a tooltip. The requirement above is
  rewritten rather than contradicted by an addendum, because nothing had been built and shipping the
  reversed rule first would have meant writing it and undoing it in the same file. Recorded here so
  the change of intent is visible in the trail rather than looking like capture got it wrong.
  Reasoning: the failure it makes visible — a typo'd or never-captured id — is silent today and
  costs a manual hunt to notice. The board already owns this exact affordance for file paths, so the
  fix is one branch reusing an existing style, not a new concept.

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
record renders in the blocked accent with the not-found tooltip rather than as plain prose, while an
ambiguous segment stays plain.

**Validation:** User confirmed — the display shape was chosen from three presented options during
planning ("Visible expansion, first mention only"), and the user separately confirmed the glossary
and that the rule covers user requests as well as requests.

## Full Context

See `do-work/user-requests/UR-074/input.md` for complete verbatim input.

---
*Source: user request in session, 2026-08-26, prompted by REQ-1685's body citing REQ-1679/REQ-1108.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome and the files are both named in the REQ, but the existing render/copy seams, the CSS token choices, and the project's Node-harness and browser-probe test conventions need discovery before writing code. No architectural decision is open.

**Planning:** Not required

## Exploration

**Where the shared helpers go.** `board-core.js` is the first fragment in `generate.go`'s
`boardJavaScriptFragmentPaths` and already holds `describeRequestStatus` (line 274) under a
`// ---- dependency helpers ---` header. `requestsById` / `userRequestsById` are declared once in the
enclosing IIFE (`board.js:18-19`) and are in scope in every fragment, so the resolver can move there
unchanged. Section headers use the exact form `  // ---- name ---…` padded to ~78 columns.
`assembleBoardJavaScript` (`generate.go:854`) rejects a fragment not ending in exactly one newline.

**`describeRequestStatus` returns a bare status string** (`"pending"`, `"claimed"`, …) or the literal
`"not in tree"` — no id, no title. The glossary must compose its own line.

**Colour tokens.** Only four ink tokens exist, and the theme is a **`prefers-color-scheme` media
query only** — there is no `data-theme` attribute and no toggle. `:root` is the dark default;
`@media (prefers-color-scheme: light)` at `board.css:136` redefines the palette.

| token | dark | light |
|---|---|---|
| `--ink-base` (body prose) | `#c7cdd6` | `#39414f` |
| `--ink-soft` (one step dimmer) | `#99a1ad` | `#59636f` |
| `--ink-faint` (two steps) | `#6b7480` | `#6c7480` |

`--ink-soft` is the correct "one step dimmer" token and holds in both themes. `--ink-faint` is
already spoken for by `.detail-dep-status` and `.detail-meta dt`. There is no prose font token to
set: `--font-sans` is on `body` and inherits into the drawer — but `.ticket-link` sets `--font-mono`
(`board.css:1730`), so a title span **inside** the anchor inherits mono and must reset it.

**Drawer markup.** `#detail-body` is the last child of `<aside id="detail-drawer">`
(`template.html:520-554`); nothing follows it, and `.detail-drawer` is a `flex-direction: column`
container, so a sibling `<section>` after it becomes the next row. There is no `.detail-body` CSS
rule — the body is styled entirely by `.markdown-body` (`board.css:1777`).

**Clearing is the trap.** `drawerMeta.textContent = ""` runs on both opens (`board-detail.js:268`
REQ, `:434` UR) but `showDrawer` and `closeDrawer` clear no content. A new glossary section is
cleared by nothing, so a UR opened after a REQ would keep the previous ticket's glossary. It must be
emptied explicitly in **both** open functions.

**First-mention state has to be threaded.** `buildLinkifiedFragment` is stateless per text node and
resets `bodyMentionPattern.lastIndex` on entry; `linkifyDetailBody` is the only place that sees the
whole body, so both the seen-set and the glossary collection belong to it.

**Node-harness pattern.** `TestJavaScriptBehaviorDrawerHeadingDeduplication` (`generate_test.go:906`)
slices production functions out of the generated `index.html` with `sliceBalancedBlockAfter`
(`:1362`, brace-counting; the sliced block must contain no braces inside string literals), builds a
DOM stub prologue, and runs the whole thing through `runJavaScriptBehaviorProbe` (`:275`) which pipes
it to `node -` on **stdin** (not `-e` — a probe embedding the client exceeds the 128 KiB arg limit).
Richer `document.createElement` / fragment stubs exist around `:1700-1850`.

**Lane membership is a name prefix, nothing else.** Any test named `TestJavaScriptBehavior…` joins
the strict lane; any `TestBrowserBehavior…` joins the browser lane. `TestMain` fails a strict run
that executed zero probes, so a lane member must actually call the probe helper. A companion
structural test asserting required substrings appear in `index.html` is conventional beside each
behavior probe (`TestDrawerDropsOnlyAMatchingLeadingHeading`, `:887`).

**Browser probes** skip when no browser is found and `t.Fatalf` when `QUEUE_KANBAN_BROWSER` names
something unrunnable. Use `runBrowserBehaviorProbeInDirectory` — the drawer needs `index.html` beside
its `board-data.js`. Conventions worth copying from
`user_request_clipboard_browser_probe_test.go:44`: splice the probe script in before the client's
closing `})();`, drive real UI and poll with a capped `waitFor`, assert
`strings.HasSuffix(result.LocationHref, ...)`, and assert zero console errors.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (modify) — receive `resolveTicketMention` and the REQ-segment index from `board-detail.js`; add `ticketTitleFor` and `shortTicketTitle`
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modify) — title-bearing `makeTicketLink`, per-render first-mention set, broken-reference span, glossary build and its clearing in both open functions
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — the empty `<section id="detail-glossary">` after `#detail-body`
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — `.ticket-link-id`, `.ticket-link-title`, `.detail-glossary`, and reuse of the blocked accent for a broken ticket reference
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — `TestJavaScriptBehavior…` probes for the pure functions plus a structural companion
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modify) — `TestBrowserBehavior…` probe for the rendered drawer

**Files I will NOT touch:** any Go source outside the two test files (the payload already carries
titles), `board-clipboard.js` (REQ-375 owns the copy surface), `board-filters.js` and any citation
index (REQ-377 owns search), and every REQ file on disk.

**Acceptance criteria (restated from REQ):**
- [ ] The first mention of each resolvable id in a drawer body renders id + truncated title; later mentions render the bare id
- [ ] `shortTicketTitle` cuts on a word boundary at 60 characters; the untruncated title is in the anchor's `title=`
- [ ] An id inside a `<code>` span keeps its bare mono link and gains no title
- [ ] `Depends on`, `Blocked by`, `Unblocks`, `Overlapping write sets` and `User request` rows always carry titles
- [ ] A glossary section lists every id the body resolved, once each, with full title and status, and is absent when the body cited nothing
- [ ] The glossary is cleared when a different ticket is opened, REQ→UR and UR→REQ alike
- [ ] UR ids and UR drawer bodies behave exactly as REQ ids and REQ drawer bodies
- [ ] An id resolving to no board record renders as a non-link blocked-accent span with a not-found tooltip; an ambiguous segment stays plain prose
- [ ] One resolver serves the whole change, defined in `board-core.js`
