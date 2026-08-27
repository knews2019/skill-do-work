---
id: REQ-382
title: 'Expand ticket ids written as Markdown links'
status: completed
created_at: 2026-08-26T17:03:51Z
user_request: UR-075
addendum_to: REQ-378
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-388]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-379, REQ-381, REQ-383]
batch: ticket-id-autocomplete
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-detail.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
claimed_at: '2026-08-27T12:54:55Z'
status_changed_at: '2026-08-27T12:54:55Z'
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
  - Route B
  - 2-file write set
  - 7 acceptance criteria
  - browser evidence
  - cross-route regression gates
  - full-suite verification
  calculated_at: '2026-08-27T12:54:55Z'
completed_at: '2026-08-27T13:09:12Z'
commit: 59caf025
kb_status: pending
---

# Expand Ticket Ids Written As Markdown Links

## What

A REQ body that writes a ticket as an explicit Markdown link — `[REQ-123](https://…)` — gets neither
a title nor a glossary entry. `linkifyDetailBody` skips any text node already inside an `<a>`, so the
one place an author has gone out of their way to mark a reference is the one place REQ-378 does not
reach.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add a two-pass regression before changing the drawer. Decorate known IDs within author-owned anchors without nested links or navigation overrides, preserve whole-ID matching and first-mention/glossary rules, and distinguish renderer autolinks. Clipboard and mention patterns stay unchanged.
- [x] **[APPLY]:** Implemented within the declared source/test write set; parent owns release and lifecycle metadata.
- [x] **[UNIFY]:** Parent reviewed the complete diff for the files in Implementation Summary, checked scope and debug artifacts, and ran request-specific checks. Verification and review results are recorded below.

## Why

Found by the Codex reviewer on PR #169 against REQ-378, verified, and deliberately deferred rather
than folded into that REQ's fourth remediation round. Two reasons, both recorded on the thread:
the corpus has **zero** instances today, and the fix touches the guard that keeps the walker from
re-entering the renderer's own autolinks — the highest-risk line in the change, traded against no
current benefit.

Deferred is not dismissed. The rule REQ-378 established is that a reader never meets a bare number,
and this is the one authoring shape where that rule does not hold.

## Context

`linkifyDetailBody` (`web/board-detail.js`) walks text nodes and skips any whose
`parentElement.closest("a")` is non-null. That guard does two jobs at once:

1. it stops the walker re-processing links the Markdown renderer itself produced (GFM autolinks,
   which the same function has just retargeted to `_blank`);
2. it stops a mention being wrapped twice.

Only the first is a real constraint. An author-written `[REQ-123](…)` is a text node inside an
anchor exactly like an autolink, and the guard cannot currently tell the two apart.

The measurement, taken at capture: `[…REQ-NNN…](…)` in a REQ or UR body — **0 matches across 373
REQs and 76 URs**. The pattern does occur in `_dev/primes/prime-kanban-board.md`, but primes are
never rendered in the drawer; only REQ and UR bodies reach `linkifyDetailBody`.

## Detailed Requirements

- **An anchor whose text carries a resolvable ticket id gains that ticket's title**, under REQ-378's
  existing rules: first mention expands, later mentions stay bare, the untruncated title rides in
  the tooltip, and `shortTicketTitle`'s 60-character word-boundary cut applies unchanged.
- **The author's `href` survives untouched.** This is the requirement that makes the REQ a design
  question rather than a mechanical edit: the author pointed that link somewhere deliberately, and
  the expansion is additive decoration on their anchor, never a replacement of their destination.
- **The glossary entry records the ticket, not the destination.** A reader looking up `REQ-123` in
  the glossary wants the REQ's title and status; where the author's link happened to point is not
  that.
- **A resolvable id inside an anchor earns a glossary line even when it does not expand** — the same
  rule REQ-378 already applies to a backticked mention.
- **Autolinks the renderer produced are still skipped.** The guard's real job is preserved; only its
  accidental second job is narrowed.
- **A dead id inside an anchor is not flagged.** An author who wrote a link made a deliberate
  reference to something; REQ-378's broken-reference flag is for ids that resolve to nothing in
  prose, and painting an author's own link in the blocked accent asserts more than is known.

- **The clipboard half is already DONE — do not rebuild it.** This REQ was written carrying REQ-379's
  finding F6 (`[REQ-100]: https://example.com/x` rewritten into `[REQ-100 (Alpha ticket)]: …`,
  orphaning every `[REQ-100]` reference in the paste). REQ-383 delivered it: the clipboard no longer
  decides anything about Markdown, and `collectDocumentTicketMentions` drops any mention goldmark
  keeps no prose text for, link reference definitions included, pinned by
  `TestCollectDocumentTicketMentionsClassifiesEveryQuotedConstruct`. **This REQ is the DRAWER only.**

## Constraints

- **Never wrap a mention twice.** The regression this guards against is a fragment nested inside a
  fragment; a test must prove a body survives two passes unchanged.
- **No change to the mention pattern.** `bodyMentionPattern` (`web/board-detail.js`) stays in
  lock-step with `bodyTicketMentionPattern` in `citations.go`, which REQ-383 pinned with
  `TestJavaScriptBehaviorTicketMentionPatternAndResolverAgreeWithGo`; it also still mirrors
  `repoFileMentionPattern` in `filementions.go` on its file-path alternative. This REQ changes where
  the walker looks, not what it matches — so neither pin should have to move.
- **No new board write surface**, and no Go source change outside the test file.

## Dependencies

`depends_on: [REQ-388]`, which shares `web/board-detail.js` with this REQ. Transitively this also
orders it after REQ-381 (`generate_test.go`, and the citation index this REQ's drawer work reads)
and after REQ-385 (`web/board-detail.js` again).

REQ-378 (archived) established the rules this extends, and REQ-383 (archived) delivered the clipboard
half this REQ was originally carrying.

**This edge serializes a shared file; it is not a need for the other's output.** `write_set` gates nothing — root `CLAUDE.md` § Glossary calls it "never a safety guarantee" — so only `depends_on` keeps two writers of one file apart under `do-work run --fan-out`. The whole batch is one chain, **REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387**, because `citations.go` alone is claimed by four of the six. That is ONE valid total order: reordering the queue means recomputing every edge, since a chain is only correct as a whole. `queue-kanban verify`'s `ungated-write-set-overlap` probe reports any pair this misses.

## Builder Guidance

**Certainty level: Mixed.** The requirement is firm; **how** the two anchor kinds are told apart is
yours to decide and to write down as a `## Decisions` entry.

The obvious discriminator — did the Markdown renderer make this anchor, or did the author — is not
directly observable after rendering. Candidates worth weighing: marking renderer-produced autolinks
during the existing retarget pass so the walker can recognise them later; comparing the anchor's
text against its `href` (an autolink's text *is* its href); or handling anchors in a separate pass
with its own rules rather than widening the text-node walk. Prefer whichever makes the
double-wrapping regression impossible rather than merely untested.

Start with the two-pass regression test, not the fix: it is what makes the candidates comparable and
it is the failure that would be expensive to ship.

## Red-Green Proof

**RED prompt/case:** Give a REQ body the line `See [REQ-1108](https://example.com/spec) for the
shape.` and open its drawer. The link renders as the bare text `REQ-1108`, carries no title and no
tooltip, and the glossary has no entry for it — while the same id written as plain prose in the same
body expands and is glossed.

**Why RED now:** `linkifyDetailBody` returns early for any text node whose `parentElement.closest("a")`
is non-null, so the mention is never offered to `buildLinkifiedFragment` at all.

**GREEN when:** That anchor shows `REQ-1108` plus its title, its `href` still points at
`https://example.com/spec`, the glossary lists REQ-1108 once with its own title and status, a second
`[REQ-1108](…)` in the same body stays bare, a renderer-produced autolink is still skipped, a dead id
inside an anchor is not flagged, and running the linkifier twice over the same body produces
identical DOM. The clipboard clause that used to close this list was delivered by REQ-383 and its
proof lives there.

**Validation:** Inferred during capture, from a verified reviewer finding. The user chose to capture
rather than build it.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: Codex reviewer P2 on PR #169, thread discussion_r3864240041, verified against the code and the corpus.*

## Triage

**Route: B** — Known anchor-skip guard; parent explored authored-anchor identity, DOM idempotence, existing helper callers and destination preservation.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modify)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify)

**Acceptance criteria (restated from REQ):**
- All Detailed Requirements and the captured Red-Green Proof are satisfied.

## Exploration

The drawer's global anchor skip prevents both nested links and authored-label decoration. Keep the existing mention pattern untouched. Distinguish renderer autolinks by their URL-like label/destination, and generated ticket/file anchors by their existing classes/data. An author's href must remain byte-for-byte its attribute value; add inert spans/text to the existing anchor, not an inner anchor or data-detail navigation that overrides it.

Two passes need more than a 'processed' marker: the first pass's already expanded anchors must still account for expansion/glossary memory when the second pass reaches a later bare authored anchor. Test mixed ordering (plain before authored; authored before plain; code label before prose; multiple references in one anchor). Do not cache at the drawer-root level without accounting for innerHTML replacement when another ticket opens. Prefer per-mention metadata/reconstruction over a long-lived mutable body cache.

The source comment in citations.go saying both surfaces skip link labels will become dated when this REQ lands, but the request explicitly forbids Go source edits outside generate_test.go. Route that prose restatement to the existing backlog, not an undeclared source edit. Clipboard link/reference-definition protection stays unchanged.

Renderer autolinks include email and www labels, not only identical http URLs. Test those alongside explicit `[REQ-123](REQ-123)` (same label/href but not an autolink). Do not add data-detail-kind/data-detail-id to an author's anchor/child if the existing delegated handler would then override its destination. Decoration metadata must not hijack navigation. Existing generated ticket expansions must not have IDs in their title text rescanned on a second pass; skip owned subtrees while reconstructing first-mention/glossary state.

Existing citations_test.go probes extract buildLinkifiedFragment plus a fixed helper set. This REQ explicitly allows no Go changes outside generate_test.go: avoid gratuitous helper extraction from the ordinary fragment path that would force those probes to gain a new dependency. Keep default call behavior compatible, and localize authored-anchor work in the drawer caller or an optional branch that existing callers never execute.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified). Reuses existing ticket resolution, title shortening and glossary accounting for text inside authored anchors. Adds inert title-bearing spans with non-navigation identity metadata, skips renderer-shaped autolinks and unknown IDs in authored anchors, and reconstructs second-pass mention state from owned decorations while skipping inserted text. The mention pattern, clipboard source, href values and normal ticket navigation remain unchanged. Diff after review remediation: 52 additions, 10 deletions.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified). Adds TestBrowserBehaviorAuthoredTicketLinksPreserveDestinationsAndTwoPassDOM, through a generated board and the actual drawer caller. Pins authored/plain/code ordering, first expansion and later bare mentions, multiple IDs per anchor, title truncation and full tooltip, glossary title/status, dead IDs, URL/www/email/mailto autolinks, relative href equal to its ticket label, navigation delegation, two-pass identical DOM, state reconstruction and reuse after replacing the drawer root's HTML. Diff after review remediation: 135 additions, 1 deletion. The deletion replaces the prior exact-signature assertion with the new optional-argument signature.

## Decisions

- **D1 — Autolink discrimination:** Goldmark does not emit anchor provenance in the available DOM. Preserve anchors whose encoded label equals their absolute URI destination, whose www label corresponds to the renderer's `http://` destination, or whose email label corresponds to a `mailto:` destination. A same-label relative ticket link remains authored. Explicit Markdown links that deliberately use exactly the renderer-shaped URL/email label are conservatively treated as autolinks too; no destination or label is rewritten. This avoids introducing a Go rendering change, which the REQ expressly prohibits.
- **D2 — Decoration and repeated passes:** An authored anchor contains spans, never new anchors or `data-detail-kind` navigation attributes. Those spans retain resolved identity under `data-ticket-kind/id`. Each pass rebuilds its own expansion/glossary state in document order by replaying an already decorated node's original resolved ID, then skips its entire text subtree. Existing generated links use their existing identity attributes. Reusing the original fragment function avoids a second glossary implementation or a new helper dependency in the existing Go/JavaScript corpus probes; discarded replay fragments do not modify the DOM. No cache survives a drawer body's replacement.

## Testing

### Test-first RED

Before source edits, added the browser regression and ran:

```sh
QUEUE_KANBAN_BROWSER="$PWD/build/chrome-141/chrome-mac-arm64/Google Chrome for Testing.app/Contents/MacOS/Google Chrome for Testing" go -C skills/do-work-board/tools/queue-kanban test -count=1 -run '^TestBrowserBehaviorAuthoredTicketLinksPreserveDestinationsAndTwoPassDOM$' -v -timeout 60s .
```

Exit 1, with assertion failures: the first authored anchor was bare and had no tooltip; authored-only IDs were absent from the glossary; the second call returned an empty glossary. The captured reproduction is `build/do-work-run/REQ-382-red.log`.

### GREEN and regressions

- Initial implementation version of the same exact Chrome for Testing 141 browser command: exit 0, package 1.009s. Evidence: `build/do-work-run/REQ-382-green.log`. It reports `location.href` and `navigator.userAgent` in the same DOM measurement. The installed binary is 141.0.7390.37; its headless UA reports 141.0.0.0.
- The final fixture places `REQ-7777` within the displayed first 60 title characters, so the second-pass prohibition on scanning inserted title text is exercised, not merely a tooltip assertion. No glossary entry appears for that title-only ID or the autolink-only ID.
- Actual first and second passes returned the same six ordered glossary IDs, exact identical body HTML, zero nested anchors, zero dead-ID flags inside authored links, original raw href attributes (including `%2F`, query ampersand and relative fragment), and no board navigation interception. Full glossary title and REQ/UR status text survived.
- Generated-board mount and interaction error/warning/rejection arrays were empty. Across light/dark at widths 320, 768 and 1280, the authored anchor had positive rendered bounds and accepted programmatic focus. This does not claim viewport containment or trusted keyboard Tab behavior; no layout or CSS changed.
- `go -C skills/do-work-board/tools/queue-kanban test -count=1 -run 'TestJavaScriptBehavior|TestCollectDocumentTicketMentions|TestBuildGeneratedBoardMarkdownData|TestBuildGeneratedBoardDataUsesOneAnalysis|TestDrawerDrops' .`: exit 0, 38.047s. Evidence: `build/do-work-run/REQ-382-regression.log`. This suite was run before the final predicate generalization from HTTP URI labels to absolute URI labels; the final CDP run covers the added mailto URI case.
- After-remediation `go -C skills/do-work-board/tools/queue-kanban vet ./...`: exit 0.
- After-remediation `node --check skills/do-work-board/tools/queue-kanban/web/board-detail.js`: exit 0.
- After-remediation `git diff --check`: exit 0.

The existing REQ-378 `TestDrawerTicketMentionsCarryTitlesAndAGlossary` exact-signature assertion was updated to include `insideAuthoredAnchor`, the intentionally added optional argument. Its remaining wiring assertions are unchanged; the behavioral caller probes still call the legacy four-argument form. One browser test was added. No separate test framework, package, CSS or dependency was introduced. The optional full browser lane and the canonical maintainer gate were not run by this builder; parent owns the canonical gate.

## Review orientation

Inspect the authored-anchor optional branch in `makeTicketLink`/`buildLinkifiedFragment` and the owned-subtree replay in `linkifyDetailBody`. The existing global mention regexp remains byte-identical. URL/path candidates inside authored labels remain consumed whole; only actual resolvable ticket candidates are decorated. The browser fixture also verifies that the delegated click handler has not prevented authored navigation before the test itself cancels the event to keep the isolated page open.

## Discovered Tasks

- Two existing `citations.go` prose comments describing all drawer anchors as skipped become stale with this REQ. This REQ explicitly forbids Go source edits outside `generate_test.go`; parent already owns routing these comment-only restatements to the backlog. Clipboard link-label suppression remains intentional and unchanged.

## Lessons Learned

- Skipping generated DOM on a second pass prevents nesting but does not preserve first-mention/glossary memory. Reconstruct state in document order from the original mention identity while refusing to scan inserted title text; a drawer-root cache would become stale when the next ticket replaces its HTML.

## Review remediation

The independent reviewer found an Important acceptance failure: Goldmark encodes Unicode in autolink hrefs but preserves readable label text. Strict equality incorrectly decorated `<custom:REQ-1108/é>` and `<mailto:REQ-1108@éxample.com>`. This was fixed within the same two-file scope.

The persistent browser fixture gained seven cases: Unicode custom-scheme and mailto links, mixed Unicode plus `%2F` segments in both schemes, malformed `%zz` and bare `%` segments plus Unicode, and an already percent-encoded Unicode URI. Each renderer autolink is now individually asserted to contain zero ticket decorations; this catches incorrect bare spans even after another occurrence has already spent the title expansion.

- **Remediation RED:** the extended test ran before the normalization change and exited 1, package 0.933s. Six autolinks had ticket decorations, the first gained a title, and the glossary incorrectly included the autolink-only REQ-5000. All expected raw hrefs already matched Goldmark's actual output. Evidence: `build/do-work-run/REQ-382-review-red.log`.
- **Remediation GREEN:** the same exact Chrome 141 command exited 0, package 0.856s. All nineteen href/text pairs passed; all eleven renderer autolinks had zero decorations; the two glossary passes still contained only the six intended IDs; existing DOM reuse, navigation and rendering checks passed. Evidence: `build/do-work-run/REQ-382-review-green.log`.
- **D1 refinement:** encode only the comparison label with `encodeURI`, then retain existing valid `%xx` segments rather than double-encoding them. Do not decode URI input: malformed percent sequences are valid input to `encodeURI` and are safely escaped, avoiding the `decodeURI` exception case. The original anchor label and href attribute are never assigned or normalized in the DOM. This keeps Unicode, escaped and malformed-percent autolinks equivalent to Goldmark's emitted destination while preserving authored relative links.
- The parent's first canonical run also exposed the old exact `makeTicketLink` signature assertion. That assertion was deliberately updated as documented above; all `TestDrawer` tests were added to the final focused regression command.

The parent owns final independent re-review and the fresh canonical gate; no post-remediation full-suite pass is claimed by this builder.

Final focused regression command after every remediation source/test edit:

```sh
go -C skills/do-work-board/tools/queue-kanban test -count=1 -run 'TestJavaScriptBehavior|TestCollectDocumentTicketMentions|TestBuildGeneratedBoardMarkdownData|TestBuildGeneratedBoardDataUsesOneAnalysis|TestDrawer' .
```

Exit 0, package 42.159s. This includes `TestDrawerTicketMentionsCarryTitlesAndAGlossary`, all existing JavaScript behavior probes and the shared citation/clipboard caller corpus. Evidence: `build/do-work-run/REQ-382-review-regression.log`. Final `go vet ./...`, JavaScript syntax and `git diff --check` also exited 0. Source remained frozen throughout these final checks.


## Parent Verification and Remediation

Parent independently reran the final authored-link browser test on Chrome for Testing 141.0.7390.37: exit 0, package 0.832s. All 19 href/text pairs, 11 excluded autolinks, same six ordered glossary IDs on both passes, and zero interaction errors passed. The earlier independent 1.085s pass preceded review remediation and is not substituted for this final check.

The first canonical gate exited 1 because the existing REQ-378 wiring test pinned the old four-argument helper signature. That assertion now pins the intentionally extended signature; its coverage was retained. The initial Important review finding (impact-user-visible) was percent-encoded renderer autolinks receiving title decorations. The final fixture reproduces it before the comparison-only encoding repair and passes afterward. Neither fix changes source hrefs, clipboard behavior, or mention patterns. Final canonical verification is recorded separately after completion.

## Qualification

Mechanical qualification exit 0. Parent verified substantive delivery against each Detailed Requirement, the two-file diff and generated-board caller; data flows retain source Markdown and href attributes, rebuild glossary state locally, and introduce no new write surface or dependency. Comment-only stale Go restatements were routed to the prose backlog under the explicit scope exclusion. The repeated-pass and encoded-autolink lessons are promoted to the maintainer board lesson index after archival.

## Review

**Overall: 99.5%** | 2026-08-27T13:06:58Z

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** Initial impact-user-visible encoded-autolink finding remediated and independently verified; none remain open.
**Minor findings:** One comment-only drift item covering two Go comments, preserved in the prose backlog under the explicit scope exclusion.
**Acceptance:** Pass — final shared Chrome141 browser probe exit 0 (0.836s), five adjacent tests exit 0 (4.968s), and isolated actual-renderer 24-URI-suffix matrix exit 0 (1.982s). Reviewer inspected both complete diffs and found no implementation follow-up. Parent separately reran the final browser test and owns the canonical gate.
**Suggested testing:** None required beyond the canonical gate; existing focus/bounds checks do not claim trusted keyboard or viewport containment.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action; independent reviewer.*

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
