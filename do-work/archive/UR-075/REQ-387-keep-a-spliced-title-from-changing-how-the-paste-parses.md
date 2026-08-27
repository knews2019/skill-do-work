---
id: REQ-387
title: '[impact-user-visible] Keep a spliced title from changing how the pasted Markdown parses'
status: completed
created_at: 2026-08-26T23:07:00Z
status_changed_at: '2026-08-27T13:09:29Z'
user_request: UR-075
addendum_to: REQ-383
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-382]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379, REQ-383]
batch: ticket-id-autocomplete
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/citations_test.go
claimed_at: '2026-08-27T13:09:29Z'
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: '2026-08-27T13:09:29Z'
completed_at: '2026-08-27T13:20:36Z'
commit:
kb_status: pending
---

# Keep A Spliced Title From Changing How The Pasted Markdown Parses

## What

An expanded title is inserted into the document's Markdown verbatim, after a 60-character cut. Two
characters in a title can change how the pasted file parses: a pipe inside a table row, and a backtick
the cut leaves unbalanced.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the board prime and lessons plus always-on and testing rules. First exercise the shipped annotator and actual Go Markdown renderer with table pipes, preceding backslashes, single/double code delimiters crossing the cut, and later author code. Sanitize only the clipboard short title at the existing insertion point; preserve the drawer helper and full appendix. Run regression, qualification, independent review and canonical verification.
- [x] **[APPLY]:** Implemented within the declared source/test write set; parent owns release and lifecycle metadata.
- [x] **[UNIFY]:** Parent reviewed the complete diff for the files in Implementation Summary, checked scope and debug artifacts, and ran request-specific checks. Verification and review results are recorded below.

## Why

The feature's contract is that a paste saves back as a valid file and reads as the same document plus
titles. Two title characters break the second half:

- **A pipe inside a table row.** A mention in a GFM table cell gets its title spliced into the cell,
  and a title containing `|` adds a column: `| REQ-501 (Split the row | keep the pipe) | a row |`
  parses as three cells, not two.
- **A backtick the cut unbalances.** `shortTicketTitle` cuts at 60 characters on a word boundary. A
  title whose backticked span straddles that boundary contributes one backtick to the paste, which
  opens a code span that runs to the next stray backtick — swallowing the prose between them.

## Context

**Latent today, and worth fixing before it is not.** Six real titles carry backticks and none
currently produces an unbalanced cut; no real title carries a pipe. Both are one ordinary title away
from being live, and the failure is silent — the paste looks fine until someone re-renders it.

Both surfaces splice titles, but only the CLIPBOARD writes Markdown. The drawer builds DOM nodes and
sets `textContent`, so a pipe or a backtick there is inert. This is a clipboard-only defect even
though `shortTicketTitle` (`web/board-core.js`) is shared.

## Detailed Requirements

- **The cut must not leave a code span open.** Either extend the cut to the span's close, pull it back
  to the span's open, or strip the stray backtick — whichever reads best, but the pasted title must
  carry an even number of backticks.
- **A title spliced inside a table row must not add a cell.** The Go index already knows the block
  each mention sits in, so the two candidate approaches are: escape `|` as `\|` when the mention is
  inside a table, or suppress the expansion there and let the appendix carry the title. Prefer
  whichever keeps the table readable.
- **The appendix keeps the full untruncated title** and is not subject to either rule — it is a list,
  not prose, and it is where a reader looks up what the cut removed.

## Constraints

- **Do not change `shortTicketTitle`'s behaviour for the drawer.** It is shared, and the drawer needs
  no escaping; a change there must be additive or clipboard-side.
- **No new board write surface.**
- Do not widen into REQ-385, REQ-386 or REQ-388.

## Dependencies

`depends_on: [REQ-382]`, which shares `generate_test.go` with this REQ. Transitively this also orders
it after REQ-381, the other `generate_test.go` writer. Its own two files —
`web/board-clipboard.js` and `web/board-core.js` — are claimed by nothing else in the queue, so this
REQ is last only because of the test file.

**This edge serializes a shared file; it is not a need for the other's output.** `write_set` gates nothing — root `CLAUDE.md` § Glossary calls it "never a safety guarantee" — so only `depends_on` keeps two writers of one file apart under `do-work run --fan-out`. The whole batch is one chain, **REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387**, because `citations.go` alone is claimed by four of the six. That is ONE valid total order: reordering the queue means recomputing every edge, since a chain is only correct as a whole. `queue-kanban verify`'s `ungated-write-set-overlap` probe reports any pair this misses.

## Red-Green Proof

**RED prompt/case:** A board where REQ-501 is titled `Split the row | keep the pipe` and a body
contains the table row `| REQ-501 | a row |`. Copy it: the pasted row has three cells. Second case: a
title whose backticked span crosses the 60-character cut — the pasted body opens a code span that
never closes cleanly.

**Why RED now:** the title is concatenated into the Markdown with no escaping and no balance check.

**GREEN when:** both pasted bodies re-render as the same document they came from, and feeding each
back through the repo's own renderer produces the original block structure.

**Validation:** Reproduced by adversarial review of REQ-383 with a purpose-built fixture board driven
through the shipped `annotateClipboardPayload`. Measured against the tree: 6 of 483 real titles carry
a backtick, 0 currently break.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-383's independent review, finding C1.*

## Triage

**Route: B** — Known splice location; renderer-based regression will prove pipes and truncated code delimiters cannot alter the pasted body. Shared drawer helper is explicitly unchanged.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modify; added after the expected raw-Markdown parity assertion failed)
- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (modify)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify)

**Acceptance criteria (restated from REQ):**
- All Detailed Requirements and the captured Red-Green Proof are satisfied.

## Exploration

Only the clipboard writes Markdown; keep drawer shortTicketTitle semantics unchanged. A clipboard-specific title helper can avoid introducing table-position metadata by escaping literal pipes in every splice, but backslashes immediately before a pipe need coverage: naively adding one slash can leave an even slash count and expose the table delimiter. Balance/truncate or strip code-span delimiters in the copied short title without changing the full appendix.

Use the actual Go renderer on pasted output and compare table cell count and prose/code-span structure, including a subsequent author-written code span that a truncated backtick could swallow. Include title backslashes around pipes. Keep the original frontmatter/raw payload untouched. Existing extracted-function tests in citations_test.go need the new helper in their harness if annotateTicketMentions calls one; declare that extra test path before edits.

The captured board-core.js path may not need a production change if the helper lives wholly beside clipboard splicing. Declare the smaller actual source scope before implementation rather than modifying shared drawer code just to satisfy a captured write set. REQ-389 owns the later ASCII arrow change, so retain the current splice marker in this REQ.

### Final implementation scope

The short clipboard title can be made inert at its existing assignment without adding a helper. Remove backtick delimiters from the short splice, then escape backslashes before escaping pipes. This makes even multi-backtick cut boundaries harmless and preserves literal pipes in every block without teaching JavaScript which block it is in. Plain text is readable in the splice; the appendix retains the entire original title including backticks. This is the request's allowed stripping choice applied consistently, rather than a second Markdown parser.

Only board-clipboard.js and generate_test.go need edits. The captured board-core.js path is unnecessary because drawer shortening stays unchanged; citations_test.go needs no helper-import changes because there is no new dependency. The earlier exploration's helper option is superseded by this smaller plan. Keep the current parenthesis marker; REQ-389 owns the arrow.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (modified). Sanitizes only the existing shortened title before splicing: remove code-span backticks, then escape every ASCII punctuation character in one pass. Full appendix, offsets and drawer title helper stay unchanged. Six additions, one deletion.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified). Adds eighteen shipped-annotator-to-real-renderer cases for plain pipes, one/two/three preceding backslashes, single/double delimiters crossing the cut, balanced code containing a pipe, exposed emphasis, unmatched emphasis spanning later author prose, explicit links, and entities. Asserts literal short-title content, unchanged full appendix title, original block/inline tag structure, final table cell, subsequent author code and prose. Eighty-four additions.

- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modified). Keeps REQ-386's same-prose-occurrence browser assertion by comparing the copied body's rendered text instead of raw Markdown, which now deliberately escapes punctuation. Ten additions, one deletion.

## Decisions

- **D1 — Strip delimiters only from the short splice:** Backticks in the inserted short title are formatting, and the requirement explicitly permits stripping them. Removing all delimiter runs is deterministic even for multi-backtick code and avoids a browser Markdown parser. This also removes complete short-title code formatting; the full title in the appendix and the drawer's text retain every original backtick. The original 60-character word cut happens first, unchanged.
- **D2 — Escape every splice's ASCII punctuation:** The initial pipe/backslash-only escape left Markdown previously inside code spans active after removing backticks. Independent review found emphasis, links and entities could change rendered structure or later author prose. Escape the complete ASCII punctuation set in one pass, including both backslashes and pipes; inserted escapes are not rescanned. This keeps former code contents literal and avoids table-position metadata, helper dependencies or a second Markdown parser. Hex ranges are equivalent to the ASCII punctuation ranges and avoid an existing test extractor treating a regex brace as a function brace.

## Testing

- Preflight exited 0: clean outside do-work, existing clipboard caller regression passed.
- Added the test before source changes. Final RED command: `go -C skills/do-work-board/tools/queue-kanban test -count=1 -run '^TestJavaScriptBehaviorClipboardTitleSplicesPreserveMarkdownStructure$' -v .`, exit 1, package 1.813s. Actual Go rendering displaced table cells and matched cut code delimiters to subsequent author code. Backslash-title rendering also lost literal text. Output is in ignored build/do-work-run/REQ-387-red.log.
- Same command GREEN: exit 0, package 1.876s, all ten cases. One manually calculated double-backtick expected word boundary was corrected after checking the unchanged shared 60-character shortener; the first post-fix run failed only that fixture expectation, not structure. The earlier RED includes actual structural failures independently of that expectation.
- All JavaScript behavior, citation collection, raw Markdown payload and drawer tests: exit 0, 43.940s. No prior assertions changed.
- Actual Chrome for Testing 141.0.7390.37 clipboard browser regressions for Copy/save/rebuild heading parity and fence/path reference parity: exit 0, 7.307s. These exercise actual copied payloads and rendered drawers through persistent CDP in static and live contexts. Page URL and UA are recorded in measurements. They do not claim a full optional browser-suite pass.
- Final JavaScript syntax and diff checks: exit 0. Canonical gate and independent review are recorded below after completion.

## Review Orientation

Review delimiter removal and the one-pass punctuation escape at the existing assignment and the actual-renderer assertions. No browser DOM/CSS behavior changes, no Go production change, no shared drawer helper change, and no Markdown scanner. Check both cell contents and tag structure: GFM can discard surplus columns while retaining the original cell count. The Referenced requests appendix uses the full original title, not the escaped short splice.

## Discovered Tasks

None. Concurrent, unowned REQ-375/Timeline comment edits appeared after preflight; they are preserved and excluded from this REQ's staging. Neither this run nor the clarification task inferred authorization from those edits.

## Lessons Learned

A Markdown source delimiter count is not a parse oracle. GFM can silently discard surplus cells, and two backticks can be an unmatched delimiter. Test the real rendered structure, preserved neighboring cell contents and author code; escape pre-existing backslashes before inserting a literal pipe escape.


## Remediation

Initial independent review found an Important / impact-user-visible defect: stripping a code span's backticks exposed its emphasis/link/entity syntax. Added eight real-renderer cases before repairing the source. Review RED exited 1 (2.250s): emphasis changed tag structure, an unmatched star absorbed following author prose, a link became an anchor and an entity changed literal text. The final 18-case test passed with exit 0 (2.588s) after escaping the full ASCII punctuation set. The initial pipe-only plan above is historical and superseded by D2.

One intermediate equivalent regexp used a literal opening brace in a character range. The production asset passed node --check, but the legacy function extractor counts all literal braces and extracted too much code; its Node probes failed with Unexpected token '}'. The final hex-range spelling avoids changing that unrelated extractor or adding a helper dependency. The 48.874s broad failure during this interval is retained as failed evidence, not called a product GREEN.

### Scope amendment before test repair

The final clipboard browser regression failed because REQ-386 compared raw copied Markdown against the drawer's literal text; its real title contains hyphens and a plus, now escaped intentionally in the clipboard. Add citations_test.go to the scope before editing. Preserve the occurrence/parity assertion by rendering the copied body through the actual Go renderer and reading its DOM text before comparing with the drawer title. This tests the same displayed text while accepting the requested safe Markdown representation; do not change Go production code or the drawer.

## Qualification

Mechanical qualifier passed. Parent inspected all three final diffs, traced each change to safe clipboard splicing or the affected previous browser expectation, and confirmed source Markdown, offset resolution, shared drawer shortening and full appendix remain unchanged. Default canonical verification is rerun after remediation; no initial check substitutes for it.

### Final regression evidence

Final all-JavaScript/citation/raw/drawer regression command exited 0 (40.379s), after the punctuation-range repair. The existing REQ-386 browser expectation intentionally changes from raw Markdown to rendered text; its first-occurrence, preserved heading, title/frontmatter save, and rebuild checks remain. Earlier no-prior-assertion-change statements apply to the initial two-file patch, before this necessary test-only scope amendment.

## Review

**Overall: 100%** | 2026-08-27T13:17:18Z

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** Initial impact-user-visible exposed-Markdown finding remediated and independently verified; none remain open.
**Minor findings:** None.
**Acceptance:** Pass — independent final 18-case renderer and clipboard/raw/underscore regressions exit 0 (4.076s), plus both exact Chrome141 clipboard CDP tests exit 0 (6.716s) after inspecting the declared third-file amendment. Parent reran those final browser tests, exit 0 (6.640s). The occurrence assertion still compares actual drawer text and copied/rendered prose; H1 preservation, full frontmatter save and rebuild remain checked.
**Suggested testing:** None required beyond canonical gate.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action; independent reviewer.*

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
