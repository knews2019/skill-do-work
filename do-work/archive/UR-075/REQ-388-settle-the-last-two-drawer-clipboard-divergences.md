---
id: REQ-388
title: '[impact-rule-change] Settle the last two drawer/clipboard divergences: fence info strings and ids inside paths'
status: completed
created_at: 2026-08-26T23:08:00Z
status_changed_at: '2026-08-27T12:45:24Z'
user_request: UR-075
addendum_to: REQ-383
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-386]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379, REQ-383]
batch: ticket-id-autocomplete
write_set:
- skills/do-work-board/tools/queue-kanban/citations.go
- skills/do-work-board/tools/queue-kanban/citations_test.go
- skills/do-work-board/tools/queue-kanban/web/board-detail.js
claimed_at: '2026-08-27T12:45:24Z'
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: '2026-08-27T12:45:24Z'
completed_at: '2026-08-27T12:54:41Z'
commit: 3ed11c17
kb_status: promoted
kb_entry: REQ-388-settle-the-last-two-drawer-clipboard-div.md
---

# Settle The Last Two Drawer/Clipboard Divergences

## What

Two places where the drawer's glossary and the paste's appendix list different ids for the same body.
Decide which surface is right in each case and make them agree.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Exclude fence metadata from clipboard entries and keep file-path runs opaque in both surfaces. Preserve citation-search independence. First add one regression comparing actual drawer and clipboard reference lists directly, including static/live and captured path cases.
- [x] **[APPLY]:** Removed the two divergent behaviors in the three declared files; no new scanner or dependency.
- [x] **[UNIFY]:** Reviewed all three source diffs and their tests. gofmt, go vet, Node syntax and diff hygiene pass. Parent independently reran the real static/live browser comparison. No debug artifacts or unrelated code changes.

## Why

REQ-383's stated rule is that the drawer and the paste say the same thing about the same body. After
its review fixes, two cases still break it in opposite directions.

**D-A — a fence info string.** `collectMentionSurfaces` registers a fenced block's `Info` segment as a
surface, so a resolved id in an info string (` ```yaml REQ-1679-template `) is emitted and lands in
the paste's appendix. goldmark renders the info string only as `class="language-…"`, never as a text
node, so the drawer's glossary has no matching entry. **The paste says more than the drawer.**

**D-B — an id nested in a repo-relative path.** The file-path alternative claims the whole run and Go
never re-scans inside it. `board-detail.js` deliberately rewinds `lastIndex` by one on a skipped match
so a REQ id inside a skipped path still links, so the drawer glossaries an id the appendix omits.
**The drawer says more than the paste.** Two real instances: REQ-112 and REQ-239. It disappears in
serve mode, where the path becomes a file link instead.

## Context

**D-B is a deliberate, documented difference, not an accident.** The clipboard's own comment in
REQ-379 said so: *"A URL or a repo-relative path is one opaque run. The drawer resumes INSIDE it so a
nested id can still link; a clipboard payload must not, or an expansion lands mid-path and the paste
carries a path that no longer names a file."* That reasoning is about the EXPANSION, which would
corrupt the path — it says nothing about the appendix line, which is additive and cannot corrupt
anything. So the likely resolution is asymmetric: never expand inside a path on either surface, but
let both record the glossary entry.

**D-A is a genuine question.** An id in an info string is arguably an illustration (the reader never
sees it) or arguably a real reference (the author typed it deliberately).

## Detailed Requirements

- **Record a `## Decisions` entry for each**, naming which surface changes and why. These are rule
  changes, not bug fixes, and the next person needs the reasoning more than the diff.
- **Whatever is decided, one test must drive BOTH surfaces over the same body** and compare their two
  reference lists directly. Two tests asserting two expectations is how these diverged.
- **An expansion must still never land inside a path** on either surface — that part of REQ-379's
  reasoning stands regardless of what the appendix does.

## Constraints

- **No new board write surface.**
- Do not widen into REQ-385, REQ-386 or REQ-387.
- **`depends_on: [REQ-386]`** because both edit `citations.go` and `citations_test.go`. The edge
  serializes a shared file; it is not a need for REQ-386's output. `write_set` gates nothing — only
  `depends_on` does — and this batch has already shipped that mistake twice.

## Dependencies

`depends_on: [REQ-386]`, which shares `citations.go` and `citations_test.go` with this REQ.
Transitively this also orders it after REQ-381 and REQ-385 — REQ-385 matters, because the two also
share `web/board-detail.js`.

**This edge serializes a shared file; it is not a need for the other's output.** `write_set` gates nothing — root `CLAUDE.md` § Glossary calls it "never a safety guarantee" — so only `depends_on` keeps two writers of one file apart under `do-work run --fan-out`. The whole batch is one chain, **REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387**, because `citations.go` alone is claimed by four of the six. That is ONE valid total order: reordering the queue means recomputing every edge, since a chain is only correct as a whole. `queue-kanban verify`'s `ungated-write-set-overlap` probe reports any pair this misses.

## Red-Green Proof

**RED prompt/case:** Open REQ-112 in a STATIC board, read the drawer's Referenced tickets list, then
Copy it and read the appendix — the drawer lists REQ-110 and the appendix does not. Second case: a
body containing ` ```yaml REQ-1679-template ` — the appendix lists REQ-1679 and the drawer's glossary
does not.

**Why RED now:** the two surfaces derive their reference lists from different inputs (rendered DOM vs
source bytes) with no test comparing them.

**GREEN when:** for both bodies the drawer glossary and the paste appendix list the same ids, and one
test asserts that equality rather than two tests asserting two expectations.

**Validation:** Both reproduced by adversarial review of REQ-383; the two real instances of D-B were
enumerated across all 453 documents.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-383's independent review, findings S3, S4 and C3.*

## Triage

**Route: B** — Two known projection divergences; parent traced fence-info surfaces and drawer retry behavior, and chose shared opaque path semantics.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/citations.go` (modify)
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modify)

**Acceptance criteria (restated from REQ):**
- All Detailed Requirements and the captured Red-Green Proof are satisfied.

## Exploration

The Go occurrence projection currently exposes a fence's Info segment even though the rendered drawer has no corresponding text node. Remove that annotation surface while preserving the independently collected citation-search set from REQ-381.

The drawer retries skipped non-ticket runs, so a repo path can turn into a nested ticket link in a static board but not a live board. Prefer keeping URLs/paths opaque in both reference projections: this preserves the Go rule, matches the existing whole-run pattern priority, avoids guessing that a filename is a citation, and removes retry logic. This is an explicit choice against the captured likely alternative of recording nested path IDs in both lists; that alternative needs an additional interior extraction rule and changes the existing citation-search semantics. The REQ delegates this choice. Record the reason and prove static/live/reference-list agreement, not only lack of expansion.

Use one regression driving Go's actual occurrence index and the drawer's actual rendered DOM over the same bodies, directly comparing references. Cover fence-info-only, real REQ-112 and REQ-239 path cases, quoted/prose path text and a later real mention (so an opaque run does not consume its expansion). Preserve the REQ-385 whole-compound guard, which was separated from generic retry deliberately for this request. Check existing tests before changing any expected path behavior; request scope expansion if a prior test outside citations_test.go must be updated.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/citations.go` (modified) — omit invisible fence metadata from the clipboard projection, preserving search.
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modified) — directly compare actual drawer/Copy references over the same sources in static/live modes.
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified) — keep skipped path runs consumed so nested IDs cannot expand or enter the glossary.

Removed fence-info segments from the clipboard annotation surfaces, and removed the drawer's retry inside skipped path matches. Neither source rule changes the independent citation-search projection. Existing path navigation remains available in live mode. No new regular expression, dependency, board write surface, parser, or framework.

## Decisions

- **D-A — Clipboard changes:** Fence information is renderer metadata, not body text. Exclude it from clipboard references so the appendix agrees with the drawer. Leave the raw bytes untouched. A resolved info-string ID remains in the independently collected citation-search set because that feature indexes authored source, not only visible text.
- **D-B — Drawer changes:** A recognized URL or repository file path is one opaque run in both reference projections. Do not retry inside a skipped path, so neither surface expands nor glossaries a nested ID. This chooses against capture's likely alternative of extracting nested IDs solely for both glossaries: that would need an additional interior-extraction rule and would make a filename into a reference despite the existing pattern priority. The simpler policy preserves Go semantics and makes static and live drawers agree. Real standalone IDs appearing later still expand.

## Testing

- **Assertion RED:** Added `TestBrowserBehaviorFenceInfoAndPathReferencesAgreeAcrossSurfaces` before changing production. Actual Chrome for Testing 141.0.7390.37 exited 1 in 3.725s. Same-body reference equality failed for info-only (drawer empty, clipboard REQ-91679), static path-only (drawer REQ-91679, clipboard empty), prose path in both modes, and the full current REQ-112 body in static mode (drawer REQ-111/REQ-110/REQ-113, clipboard REQ-111/REQ-113). See ignored `build/do-work-run/REQ-388-red.log`.
- **GREEN:** That browser test plus `TestCollectDocumentTicketMentionsClassifiesEveryQuotedConstruct` exited 0 in 3.910s. All six bodies were measured in file:// static and production HTTP-handler live modes. Each measurement returned page URL and browser user agent, with no boot/interaction warning, error, or unhandled rejection. The probe opens the actual drawer and clicks the actual Copy button; it directly compares the rendered glossary IDs against parsed IDs in the real clipboard appendix, not independent expected lists. It accounts for Copy's existing self-reference exclusion when comparing external references.
- The path-only fixture remains unexpanded byte-for-byte, lists no nested ticket, and has zero file links in static mode / one in live mode. A later real mention expands exactly once in the drawer and clipboard. Fence-info-only still has its citation-search entry.
- **Real cases:** Full REQ-112 now yields REQ-111/REQ-113 on both surfaces in both modes. Full current REQ-239 yields REQ-233/REQ-240/REQ-236/REQ-235. The captured nested-file-path issue no longer reproduces in the current REQ-239 body: its backticked `worktree-agent-REQ-239-…` is not a file path. The test includes the whole current body without claiming it was RED, while synthetic path-only cases keep that regression discriminating.
- **Prior test intentionally changed:** REQ-383's `TestCollectDocumentTicketMentionsClassifiesEveryQuotedConstruct` no longer expects the one fence-info entry. Its invalid-backtick info case still expects prose expansion; all other quoted constructs remain unchanged.
- Focused Go/Node regression run exited 0 in 37.747s: `go test -count=1 -run '^(TestJavaScriptBehavior|TestCollectDocumentTicketMentions|TestGeneratedMentionsPreserveRestatingHeadingRoundTrip|TestBuildGeneratedBoardMarkdownData|TestBodyTicketMentionPattern|TestStaticAndLiveCitation)' -timeout 180s .`. This preserves the prior underscore/compound, H1, raw-source, and static/live citation tests.
- `go vet ./...`, `node --check web/board-detail.js`, and `git diff --check` exited 0. Canonical maintainer verification remains the parent's responsibility.

## Qualification / Orientation

Source frozen for parent review and canonical gate. Production changes are deletions and comments within the declared three-file scope. Reviewed both source diffs and the actual browser outputs. No changes to raw Markdown payload storage, H1 normalization/removal, Unicode boundary behavior, authored-link behavior, clipboard title formatting, Timeline code, version metadata, or lifecycle records. No layout/style changes or new interaction mechanism.

## Lessons Learned

When two projections differ, compare their final reference lists over the same source and exercise both static and live production paths. Merely finding the same regex candidates misses a drawer retry that exists only after a path fails to become a link. Keep annotation suppression independent from source-search citation collection.

## Discovered Tasks

None. REQ-389 will intentionally update the new test's `then REQ-91679 (Referenced title)` Copy assertion when changing the splice marker; this is already within its planned citations_test.go scope expansion.

## Parent Verification

Parent inspected the three actual diffs and reran `TestBrowserBehaviorFenceInfoAndPathReferencesAgreeAcrossSurfaces` with exact Chrome 141: exit 0, 3.862s. All twelve body/mode comparisons passed with page/browser provenance. The full optional browser suite remains unverified; the default canonical gate skips it.

## Qualification

The three-file manifest matches the declared scope; code changes are deletions plus explanatory comments. The only changed old expectation is the intentionally removed fence-info clipboard entry. Invalid-fence prose, underscored/compound IDs, H1 behavior, original bytes/offsets, and independent eager citations remain covered.

## Orientation

The drawer and copied appendix now list the same external references for fence metadata and file-path cases. Static boards no longer turn part of a file path into a ticket link, while live file navigation remains available.

## Review

**Independent review: Pass, 100%.** All scored dimensions 100%; no findings. Reviewer independently ran actual Chrome141 reference agreement (exit 0, 4.350s) and eight focused underscore/H1/index/quoted/link regressions (exit 0, 3.487s), inspected the assertion RED and actual production diff, and confirmed the intentional REQ-383 expectation change. Restatement sweep found no new stale contract. No follow-up required.

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
