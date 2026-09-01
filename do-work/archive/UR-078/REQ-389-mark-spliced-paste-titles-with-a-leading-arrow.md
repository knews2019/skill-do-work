---
id: REQ-389
title: 'Addendum: mark spliced paste titles with a leading arrow'
status: completed
created_at: 2026-08-27T11:22:38Z
user_request: UR-078
addendum_to: REQ-379
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-387]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-379, REQ-383, REQ-387]
batch: ticket-id-autocomplete
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/citations_test.go
claimed_at: '2026-08-27T13:20:59Z'
status_changed_at: '2026-08-27T13:20:59Z'
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: '2026-08-27T13:20:59Z'
completed_at: '2026-08-27T13:26:08Z'
commit: 4ed31496
kb_status: promoted
kb_entry: REQ-389-addendum-mark-spliced-paste-titles-with-.md
---

# Addendum: Mark Spliced Paste Titles With A Leading Arrow

## What

The Copy buttons splice a ticket's title after the first body mention of its id, as
`REQ-374 (Show how long each done card took)`. Change the spliced form to
`REQ-374 (-> Show how long each done card took)` so a reader of the paste can tell the
parenthetical was inserted by the board, not written by the ticket's author.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the board prime, lessons and implementation/testing rules. Update existing expected splice literals to the exact ASCII marker first and observe RED, then change the one production insertion string. Keep the drawer, safe-title encoding, full appendix, raw source payload and mention offsets unchanged. Run existing caller/renderer/browser probes, independent review and canonical verification.
- [x] **[APPLY]:** Implemented within the declared source/test write set; parent owns release and lifecycle metadata.
- [x] **[UNIFY]:** Parent reviewed the complete diff for the files in Implementation Summary, checked scope and debug artifacts, and ran request-specific checks. Verification and review results are recorded below.

## Why (if provided)

"so I know that this is a programatic expand" — a pasted body carries no styling, so
without a marker the parenthetical reads as the author's own words.

## Context

- The insertion text is built client-side in
  `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (`" (" + expandedTitle + ")"`,
  currently around line 199); the insertion positions come precomputed from the Go side (REQ-383).
- Scope is the in-body splice only. The drawer already marks the expansion visually (the title
  renders in its own quieter-styled span), and the Referenced requests appendix is self-evidently
  board-generated, so neither changes. Capture judged this from the user's example, which shows
  only the in-body form.
- The marker is ASCII `->` followed by a space, exactly as the user typed: `(-> Title)`. Do not
  substitute a typographic arrow.

## Prior Implementation

- REQ-379 (commit f1e2ce8) shipped the clipboard title splice and the Referenced requests appendix.
- REQ-383 (commit a3d4e4c) moved mention resolution into Go (`citations.go`) and deleted the
  client-side Markdown scanner; the page now inserts titles at positions the board computed.
  The literal inserted text remains a client-side concern in `board-clipboard.js`.

## Constraints

- This is a deliberate paste-only divergence from the drawer's rendering: the drawer signals the
  expansion with styling, the paste signals it with `->`. Work descending from the
  drawer/clipboard-divergence line (REQ-386/REQ-388) must not "settle" this marker away.
- Any test pinning the old `(<title>)` spliced form updates to the new `(-> <title>)` form.

## Dependencies

`depends_on: [REQ-387]` serializes the shared files (`web/board-clipboard.js`,
`generate_test.go`) behind the end of the queued ticket-id chain
(REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387). The edge orders writers of the
same files; it is not a need for REQ-387's output.

## Red-Green Proof

**RED prompt/case:** Copy a ticket whose body's first mention of REQ-374 gets a title spliced.
The paste reads `REQ-374 (Show how long each done card took)`.
**Why RED now:** `board-clipboard.js` builds the insertion as `" (" + expandedTitle + ")"` — no
marker distinguishes the inserted title from author-written parentheticals.
**GREEN when:** The same paste reads `REQ-374 (-> Show how long each done card took)`, and a
harness test pins the `(-> ` prefix on the spliced insertion.
**Validation:** User confirmed — the before/after forms are the user's own example, verbatim.

---
*Source: "Instead of: REQ-374 (Show how long each done card took). expand it like this REQ-374 (-> Show how long each done card took). so I know that this is a programatic expand"*

## Triage

**Route: A** — Exact literal splice marker change; existing raw-paste and rendered-copy assertions update together, with no matching, drawer or appendix changes.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Exploration and Scope

The current insertion literal is in annotateTicketMentions. This is a Route A mechanical replacement, so no preflight overhead. Add citations_test.go to the captured two-file write set before edits because REQ-385's underscore corpus and REQ-386/388's actual Copy browser assertions deliberately pin the old in-body form. Only their expected spliced marker changes; exclusions and appendix expectations remain. The original author prose and the drawer never gain this marker.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (modified). Changes the single in-body insertion prefix to the exact ASCII arrow and space. Safe-title escaping, offsets and the appendix are unchanged. One line replaced.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified). Updates existing exact clipboard expectations and the actual-renderer title checks to the new marker. The serialized HTML expectation uses the renderer's escaped greater-than entity; copied Markdown still asserts literal ASCII. Broad negative exclusion assertions retain their original prefix so they reject either marker. Fourteen lines replaced.
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modified). Updates the existing underscore corpus and actual Copy/save/rebuild and path-reference browser expectations to the deliberate clipboard marker. Drawer and source expectations remain unchanged. Seven lines replaced.

## Decisions

**D1 — Clipboard only:** The request explicitly chooses ASCII rather than a typographic arrow and excludes both drawer styling and the full appendix. Insert the marker outside the safely escaped short title so it is exactly the requested Markdown text. Existing tests prove the intended representation change without adding a duplicate test framework or scanner.

## Testing

- Test-first RED: changed the existing clipboard caller's expected output before changing the source. `go -C skills/do-work-board/tools/queue-kanban test -count=1 -run '^TestJavaScriptBehaviorClipboardAnnotatesBodiesAndAppendsOneGlossary$' -v .` exited 1 (1.393s), showing the old unmarked parentheses against the new exact marker.
- GREEN: clipboard caller, 18 real-renderer safe-title cases, and underscore corpus together exited 0 (4.227s). The same exact appendix expectations remained unchanged.
- Actual Chrome for Testing 141.0.7390.37: Copy/save/rebuild heading and fence/path browser regressions exited 0 (7.649s), measuring real drawer and copied text in static/live contexts. The deliberate marker is present in the paste and absent from drawer title spans. URL and UA are reported with each measurement. No full optional browser-suite pass is claimed.
- JavaScript syntax and diff checks exited 0. Independent review and canonical results are recorded below after completion.

Intentional cross-REQ assertion updates: REQ-379/383 clipboard payload expectations; REQ-385 underscore copied body; REQ-386 same visible prose occurrence; REQ-388 path suppression versus later expansion; REQ-387 rendered short-title literal. Only the requested marker changes. No source Markdown, Go production, drawer, title shortening or appendix code changes.

## Review Orientation

Check the one-line production diff and all expectation replacements. Confirm the arrow is ASCII in the copied source and the greater-than HTML entity is only an HTML serialization expectation. The actual DOM/browser parity test reads rendered text and expects ASCII. Existing negative code/fence exclusion guards still reject any parenthetical expansion, where they were broad before.

## Discovered Tasks

None. Separately approved REQ-375 changes remain outside this commit; parent will process that request after this one.

## Lessons Learned

No new reusable lesson beyond the existing clipboard representation and actual-renderer lessons. This change deliberately distinguishes inserted text from author prose; do not later remove that marker to make clipboard text identical to drawer styling.


## Qualification

Mechanical qualification exit 0. Parent verified the one-line behavior change, exact ASCII bytes, existing raw source/appendix invariants and each intentionally updated caller expectation. No implementation scope drift or new dependency.

## Review

**Overall: 100%** | 2026-08-27T13:23:06Z

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.
**Minor findings:** None.
**Acceptance:** Pass — independent four caller/renderer/underscore/drawer tests exited 0 (5.031s); scope, syntax and exact marker checked. Greater-than is escaped only in serialized HTML, while copied source contains the requested ASCII marker. Parent actual-browser regressions also passed as recorded above.
**Suggested testing:** None beyond canonical gate.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action; independent reviewer.*

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
