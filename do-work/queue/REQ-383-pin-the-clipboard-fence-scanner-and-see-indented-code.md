---
id: REQ-383
title: '[impact-negligible] Pin the clipboard fence scanner, and teach it to see indented code blocks'
status: pending
created_at: 2026-08-26T19:10:32Z
user_request: UR-075
addendum_to: REQ-379
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-382]
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-379, REQ-382]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Pin The Clipboard Fence Scanner, And Teach It To See Indented Code Blocks

## What

Eight branches of REQ-379's clipboard fence scanner are correct as shipped and pinned by nothing, so
the next edit drops them silently. Separately, the scanner is structurally blind to one Markdown
shape: a four-space indented code block, which has no fence line for a fence scanner to find.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-379's independent review ran 24 mutations against the new scanner and **8 survived**. Every
survivor is correct behaviour — the reviewer verified each one directly under Node — but an
unpinned branch is one edit from silently disappearing, and these sit on the most intricate code the
REQ added.

That REQ also shipped a builder-reported mutation that **survived its first attempt** because its
anchor landed on a line the fixture never exercised. A vacuous mutation is worse than no mutation:
it reports coverage that does not exist. This REQ exists so that class of gap stops being carried
forward.

## Context

The eight unpinned branches, from REQ-379's review:

| Mutation that survives | Branch left unpinned |
|---|---|
| `while (scanOffset < 3 …)` → `< 0` | the up-to-three-space fence indent |
| `runLength >= openFenceRun.runLength` → `===` | a ```` ```` ````-fenced block quoting a ```` ``` ````-fenced one |
| `infoText.trim() === ""` → `true` | a closer carrying an info string must not close |
| `runEndOffset - runStartOffset === runLength` → `>=` | the equal-length inline code span rule |
| drop `missingNewline` | the setext-H2 guard for a payload with no trailing newline |
| drop the BOM skip | BOM tolerance, mirroring `frontmatter.go` |
| inline-span `flagMissingIds` `true` → `false` | REQ-378 **D-07**'s second half — an inline code span *does* flag a dead id |
| gate `recordReferencedTicket` on `expandTitles` | an id seen only inside a fenced block still earns its appendix line |

**The last two are cross-surface agreement claims.** The code comments assert twice that the drawer
and the paste "say the same thing about the same body", and nothing checks it. That is the shape
`_dev/primes/prime-kanban-board.md` records from REQ-248: one domain read by two languages needs a
test that fails whichever side drifts alone.

**The setext branch is live, not theoretical.** `copyTextWithHeading` returns
`headingLine + "\n\n" + remainingBody` with no guaranteed trailing newline, and the Go payload
preserves a file that has none. Every current fixture ends in `\n`, which is why nothing reaches it.

**The indented-code-block gap is structural.** A four-space indented code block has no fence line at
all, so `codeFenceRunFor` can never see it — closing this needs indented-block tracking in
`annotateMarkdownBody`, a different mechanism from the fence state machine. The same applies to a
fence indented four or more spaces under a list item. Measured at capture: **zero instances in
`do-work/`**, which is why this is captured rather than urgent.

## Detailed Requirements

- **Pin all eight branches.** Each is one or two extra lines in the existing host document in
  `TestJavaScriptBehaviorClipboardAnnotatesBodiesAndAppendsOneGlossary`. Every assertion must name
  the behaviour it protects, not the mutation that found it.
- **Prove each new assertion bites** by running the mutation from the table and watching that
  assertion fail by name. An assertion added without its mutation run is the vacuous-coverage failure
  this REQ was written to stop.
- **Verify each mutation is non-vacuous** — confirm the test input actually reaches the mutated line.
  A mutation on a branch no fixture exercises proves nothing and reports coverage that is not there.
- **Skip four-space indented code blocks** when annotating, matching what the drawer does via
  `closest("pre")`. An indented block is quoted content and its ids are illustrations.
- **A fence indented four or more spaces at top level is an indented code block, not a fence** —
  CommonMark's rule, and the reason the existing three-space bound exists.
- **The two cross-surface claims get agreement tests**, failing whichever side drifts alone, rather
  than two independent assertions that can pass while disagreeing.

## Constraints

- **Behaviour must not change for the eight pinned branches.** They are correct today; this REQ makes
  them provable, and a test that forces a behaviour change means the branch was misread.
- **No new dependency between the drawer and the clipboard.** They share `board-core.js`'s resolver
  and nothing else; the agreement tests compare outputs, they do not couple the two surfaces.
- Keep the Go payload untouched. The two round-trip tests must stay green **and unmodified**, as in
  REQ-379.

## Dependencies

`depends_on: [REQ-382]`, which follows REQ-381 and REQ-379. The edge serializes `generate_test.go`,
which every REQ in this batch writes; `write_set` is display-only and never a scheduling gate, so the
ordering has to live in `depends_on`. This REQ needs nothing from REQ-382's output.

## Builder Guidance

**Certainty level: Firm on the pinning, Mixed on the indented-block work.**

The eight assertions are mechanical — the table names each branch and the fixture already exists.

The indented-block tracking is the real work and the only place judgment is needed: `annotateMarkdownBody`
currently reasons line-by-line about fences, and an indented code block is defined by its
indentation plus what precedes it (a blank line, and not being inside a list item's continuation).
Decide how much of CommonMark's block structure is worth modelling here and write that boundary down
in a `## Decisions` entry — the honest answer may be "handle the top-level case and state that a
list-nested one is out of scope", which is fine if it is said rather than left implied.

Read `_dev/primes/prime-kanban-board.md` first. REQ-248 (agreement assertions) and REQ-293 (choose the
mutation before looking at the pattern, or you will choose the one the pattern already catches) both
bear directly on this REQ.

## Red-Green Proof

**RED prompt/case:** Apply any mutation from the Context table to `web/board-clipboard.js` and run
`go test -run 'TestJavaScriptBehaviorClipboard' ./...`. It passes. Eight distinct behaviour changes
are invisible to the suite. Separately, a body containing a four-space indented code block whose
lines mention `REQ-1679` has that id expanded, where the drawer showing the same body leaves it
alone.

**Why RED now:** No fixture exercises those eight branches, and `codeFenceRunFor` is a fence
scanner — an indented code block has no fence line for it to match.

**GREEN when:** Each of the eight mutations fails a named assertion; each mutation is demonstrated to
reach a line the fixture exercises; an id inside a four-space indented code block is neither expanded
nor glossed; and the drawer and clipboard are shown to agree on the same body by a test that fails if
either side alone changes.

**Validation:** Inferred during capture, from a verified independent review finding.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-379's independent review, findings F2, F3 and the structural half of F5.*
