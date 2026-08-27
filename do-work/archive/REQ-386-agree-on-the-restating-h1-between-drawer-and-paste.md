---
id: REQ-386
title: '[impact-user-visible] Make the drawer and the paste agree about a body H1 that restates the title'
status: completed
created_at: 2026-08-26T23:06:00Z
status_changed_at: '2026-08-27T12:28:08Z'
user_request: UR-075
addendum_to: REQ-383
review_generated: true
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-381]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-378, REQ-379, REQ-383]
batch: ticket-id-autocomplete
write_set:
- skills/do-work-board/tools/queue-kanban/citations.go
- skills/do-work-board/tools/queue-kanban/citations_test.go
- skills/do-work-board/tools/queue-kanban/generate.go
- skills/do-work-board/tools/queue-kanban/go.mod
- skills/do-work-board/tools/queue-kanban/go.sum
claimed_at: '2026-08-27T12:28:08Z'
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: '2026-08-27T12:28:08Z'
completed_at: '2026-08-27T12:45:11Z'
commit: 59577def
kb_status: pending
---

# Make The Drawer And The Paste Agree About A Body H1 That Restates The Title

## What

The drawer deletes a body's opening H1 when it restates the frontmatter title, then decides which
mention expands. The clipboard keeps that H1 and counts it as the first prose mention. Pick one rule
and apply it to both.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Suppress clipboard annotation in the opening H1 when its rendered text normalizes to the model title, preserving raw bytes and citation search. Pin the shared comparison and actual copy/save/rebuild round trip for all three captured records before implementation.
- [x] **[APPLY]:** Implemented the five declared source/module files, with the reviewed conditional heading-fragment exception documented below.
- [x] **[UNIFY]:** Reviewed all five source diffs listed in the manifest. gofmt, go vet, module verification and diff hygiene pass. Tests exercise actual copy/save/rebuild and the same normalization corpus. No debug artifacts, raw-source rewrites, or Timeline edits.

## Why

`linkifyDetailBody` (`web/board-detail.js`) removes the restating H1 before linkifying, so the drawer
expands the first mention in the BODY. The Go walk sees the whole file and expands the first mention
in the H1. Three real records hit this today:

| Record | Its restating H1 |
|---|---|
| REQ-041 | Confirm: three small board/pipeline hardening follow-ups from REQ-034 |
| REQ-042 | Confirm: three worktree-mode evidence-path hardening follow-ups from REQ-037's review |
| REQ-085 | Run REQ-073's Live Two-Builder Acceptance Test and Record What It Found |

Two consequences, and the second is the one that matters:

- The two surfaces annotate **different occurrences** of the same id. The drawer expands the prose
  mention and the paste leaves it bare, expanding the heading instead.
- **The paste stops round-tripping.** Once the H1 reads `… from REQ-034 (Capture-time decomposition
  nudge…)` it no longer equals the record's `title:` field, so on save-back the drawer's exact-match
  H1 removal stops firing and the reader sees the title twice.

## Context

The H1-stripping rule is old and deliberate: REQ/UR bodies conventionally open with an H1 restating
the frontmatter title, and showing it inside a drawer that already displays the title is duplication.
`copyTextWithHeading` implements the same de-duplication on the clipboard's rendered-text FALLBACK
path — so the clipboard already knows the rule, and only its primary path ignores it.

Not a regression: REQ-379's scanner annotated the H1 too. REQ-383 moved the decision into Go and kept
the behaviour.

## Detailed Requirements

- **Decide which surface is right and say why in a `## Decisions` entry.** The two candidate rules:
  suppress annotation inside a restating H1 (the paste keeps a heading identical to its title, and
  save-back stays clean), or stop stripping the H1 in the drawer (both show it, no rule needed). The
  first preserves the round trip and is the smaller change; the second removes a rule instead of
  adding one, which this repo generally prefers. Argue it rather than assuming.
- **Whichever is chosen, the H1 comparison must be the SAME comparison on both sides.** The drawer
  uses `normalizeHeadingText` (collapse whitespace, trim, lowercase). Go needs the same normalization
  or the two disagree on a heading that differs only in spacing.
- **The round trip is the acceptance test**: copy one of the three records above, save the paste as a
  file, rebuild the board from it, and confirm the drawer still strips its H1.

## Constraints

- **No new board write surface**, and no change to what the frontmatter fence carries.
- Do not widen into REQ-385, REQ-387 or REQ-388.

## Dependencies

`depends_on: [REQ-381]`, which shares `citations.go`, `citations_test.go` and `generate.go` with this
REQ. Transitively this also orders it after REQ-385.

An earlier draft of this section claimed a "disjoint write set from REQ-385", which this REQ's own
frontmatter contradicted on two files. That was the third write-set/`depends_on` mismatch in this
batch, and it is why the `ungated-write-set-overlap` probe now exists.

**This edge serializes a shared file; it is not a need for the other's output.** `write_set` gates nothing — root `CLAUDE.md` § Glossary calls it "never a safety guarantee" — so only `depends_on` keeps two writers of one file apart under `do-work run --fan-out`. The whole batch is one chain, **REQ-385 → REQ-381 → REQ-386 → REQ-388 → REQ-382 → REQ-387**, because `citations.go` alone is claimed by four of the six. That is ONE valid total order: reordering the queue means recomputing every edge, since a chain is only correct as a whole. `queue-kanban verify`'s `ungated-write-set-overlap` probe reports any pair this misses.

## Red-Green Proof

**RED prompt/case:** Open REQ-041 in the drawer and read where the title appears; then Copy it and
read the paste. The drawer expands the prose mention on the later line; the paste expands the one in
the H1 and leaves the prose mention bare.

**Why RED now:** `linkifyDetailBody` removes the H1 at `web/board-detail.js:208` before the mention
walk; the Go walk never sees that removal.

**GREEN when:** both surfaces annotate the same occurrence for all three records, and a pasted copy of
REQ-041 saved back to disk still has its H1 stripped by the drawer.

**Validation:** Reproduced by adversarial review of REQ-383 against the generated board in headless
Chromium; the three affected records were then enumerated across the whole tree.

## Full Context

See `do-work/user-requests/UR-075/input.md` for complete verbatim input.

---
*Source: REQ-383's independent review, finding C2.*

## Triage

**Route: B** — Known H1 suppression defect; reuse the existing raw AST and resolve heading normalization across drawer and clipboard.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/citations.go` (modify)
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify)

- `skills/do-work-board/tools/queue-kanban/go.mod` (modify)
- `skills/do-work-board/tools/queue-kanban/go.sum` (modify)

**Acceptance criteria (restated from REQ):**
- All Detailed Requirements and the captured Red-Green Proof are satisfied.

## Exploration

The three captured real records exist at archive/UR-007/REQ-041-confirm-board-hardening-followups.md, archive/UR-007/REQ-042-worktree-evidence-path-hardening-followups.md, and archive/UR-016/REQ-085-run-the-live-two-builder-acceptance-test.md. Their opening H1 repeats the frontmatter title; REQ-085 differs in letter case. Each has a later prose mention of the referenced ticket. The drawer removes its first rendered H1 using normalizeHeadingText, defined in board-clipboard.js as collapse whitespace, trim, lowercase; this helper also governs fallback Copy deduplication.

Prefer suppressing clipboard annotations in the matching first H1: preserve the existing duplicate-title removal and raw saved file, rather than reintroducing duplicate headings everywhere. Keep the heading bytes intact and leave the citation-search set available; change only which occurrence earns annotation/glossary in the clipboard projection. Reuse the raw-body AST from REQ-381's analysis. Match the rendered H1 text and normalization, with a shared test corpus for both languages, and test the captured real-record copy/save/rebuild path with actual drawer stripping.

### Scope amendment before dependency edits

The shared comparison probe exposed full Unicode casing differences (dotted I and contextual Greek sigma) between JavaScript and Go strings.ToLower. Add golang.org/x/text/cases with language.Und instead of a hand-maintained Unicode classifier or an English-only contract. This adds go.mod/go.sum to the declared scope before modification. It is required to preserve the requested comparison, not a new board feature or write surface. The frozen estimate is unchanged.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/citations.go` (modified) — suppress clipboard entries in the opening rendered H1 when it normalizes to the model title; preserve independent citation search.
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modified) — caller regression, shared Go/JavaScript normalization, and actual browser Copy → saved file → rebuilt drawer for the three captured records.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — pass the REQ or UR model title into the existing shared raw analysis.
- `skills/do-work-board/tools/queue-kanban/go.mod` (modified) — pin golang.org/x/text v0.41.0 for full default Unicode lowercase.
- `skills/do-work-board/tools/queue-kanban/go.sum` (modified) — module checksums.

## Decisions

- **D1 — Preserve heading suppression.** The existing drawer avoids duplicating its title. Keeping the source H1 untouched and excluding its clipboard entries preserves save-back while letting the same visible prose occurrence expand. Removing the drawer rule would duplicate titles across every conventional record.
- **D2 — Use the existing renderer on the existing H1 node, with a conditional preprocessed fragment.** The normal path handles formatting, code spans, escapes, and entities using the existing AST. Review found that question-option preprocessing can append a visible backslash to an ATX H1 or destroy a setext underline. The matcher now applies the existing preprocessor to the prefix through two lines after the heading text. Only when that prefix changes does it parse the transformed fragment, preserving reference definitions from the original whole-body parser context. This also handles a model title that already matches the displayed added backslash. Raw HTML blocks become comments under the safe renderer, so they do not count as the first element, but their prefix remains available for the preprocessor's fence state. Strip renderer-owned tags before decoding entities. The original raw AST and offsets remain authoritative for the surface walk and citation search.
- **D3 — Use full Unicode casing.** A bounded comparison found `strings.ToLower` disagrees with JavaScript for dotted I and contextual Greek sigma. Parent approved and declared the two module files before edits. `cases.Lower(language.Und)` avoids maintaining a local Unicode classifier; its caser is local because it may be stateful. The module requires Go 1.25; this module's `go 1.26` directive and toolchain remain unchanged. See [official cases documentation](https://pkg.go.dev/golang.org/x/text/cases).

## Testing

- **Assertion RED:** `TestGeneratedMentionsPreserveRestatingHeadingRoundTrip` exited 1 before implementation. The H1 took the first expansion or retained a clipboard entry, and a title-only citation leaked into the clipboard index. Log: `build/do-work-run/REQ-386-red.log`.
- **Unicode RED:** `TestJavaScriptBehaviorHeadingNormalizationAgreesWithGo` exited 1 with dotted I and three contextual sigma mismatches under simple Go lowercase. Log: `build/do-work-run/REQ-386-unicode-red.log`.
- **Focused GREEN:** `go -C skills/do-work-board/tools/queue-kanban test -count=1 -run 'TestCollectDocumentTicketMentions|TestBuildGeneratedBoardMarkdownData|TestGeneratedMentionsPreserveRestatingHeadingRoundTrip|TestJavaScriptBehavior' .` exited 0, 34.038s. Log: `build/do-work-run/REQ-386-focused.log`.
- **Browser GREEN:** Exact Chrome for Testing 141.0.7390.37 ran `TestBrowserBehaviorRestatingHeadingsSurviveCopySaveRebuild`, the Go caller regression, and shared normalization together: exit 0, 4.075s. Log: `build/do-work-run/REQ-386-browser.log`.
- The browser used the real generated board, real card and Copy handlers, then saved all three copied payloads to temporary files, parsed those saved files, and rebuilt the board. Each drawer still starts at H2 and expands the same first visible prose occurrence. Captured REQ-041 and REQ-042 already had authored parentheses, so the assertion derives the board-added title from `.ticket-link-title` and compares its prose prefix. REQ-085 covers case differences. Frontmatter and the raw H1 remain unchanged.
- Twenty-one shared H1 cases cover REQ and UR model-title callers, title-only search independence, formatting/entities/escaping, comments and omitted HTML, ECMAScript whitespace (BOM and U+0085 difference), Unicode casing, nonmatching/not-first/H2/nested controls, and the preprocessing cases below. Actual browser stripping uses the same corpus. Browser measurements carry their page URL and browser identity; mount and interaction errors were empty.
- `go vet ./...` exited 0. `git diff --check -- skills/do-work-board/tools/queue-kanban` exited 0.
- Full canonical verification belongs to the parent; no full browser-suite pass or responsive geometry claim is made here.

## Prior Test Changes

The existing dependency citation fixture receives the new explicit empty model-title argument. No prior expected behavior was weakened. A later arrow-marker request must deliberately update this new browser test's `(Title)` body-splice expectation when that syntax changes.

## Qualification / Orientation

Reviewed all five source diffs. No browser Markdown scanner, new board write surface, extra body disk reads, raw payload rewrite, or timeline edits. The normal path uses one raw AST. The preprocessing exception reparses only the transformed prefix through the heading's following lines, carrying the original reference context. The title comes from the model for both REQ and UR documents, including records whose title is not stored in frontmatter. The existing UTF-16 occurrence cursor and eager citation sets remain intact.

## Review Remediation R1

The independent reviewer reproduced a parity defect when `Recommended:` immediately follows an H1: the drawer preprocessor changes the displayed heading before comparing its title, while the first implementation compared the original AST. Added shared-corpus cases and observed assertion RED for ATX and setext forms, plus titles matching the inserted backslash and headings with later reference definitions. Log: `build/do-work-run/REQ-386-preprocessor-red.log`.

The conditional fragment approach above fixes the comparison without copying the preprocessor's prefix vocabulary or changing raw offsets. An initial fragment bound covered ATX but missed setext's following option line; extending it through the second following line made both cases pass. Blank-line and existing-two-space controls still strip, and omitted HTML retains preprocessing fence state.

Final targeted Go, shared JavaScript normalization, and actual Chrome 141 browser run exited 0 in 4.011s. All twenty-one cases and all three real copy/save/rebuild paths passed with no browser errors. Log: `build/do-work-run/REQ-386-preprocessor-green.log`.

Final broad Go/JavaScript regressions exited 0 in 33.650s using the same focused selection recorded above. Log: `build/do-work-run/REQ-386-final-regressions.log`. Final `go vet ./...` and scoped `git diff --check` both exited 0. The earlier interim broad run in `REQ-386-final-focused.log` correctly failed the too-short setext prefix and is not final-pass evidence. Parent owns the repeated canonical gate.

## Lessons Learned

Rendered heading text is not the Markdown heading source. Reuse the renderer, account for its preprocessing, and explicitly match JavaScript whitespace and full lowercase before claiming two languages perform the same comparison. When reparsing a fragment, carry reference definitions from the document or the fragment can silently change a heading's text. A copy/save/rebuild test catches heading annotation that a single-surface title test misses.

## Discovered Tasks

None.

## Parent Verification

Parent independently ran the persistent actual Chrome 141 copy/save/rebuild test: exit 0, 3.282s. REQ-041, REQ-042 and REQ-085 all rebuilt with H2 as the first body element, the same first expanded prose title, and no captured errors. `go mod verify` reported all modules verified; diff hygiene passed.

## Review Remediation

Independent review found a narrow current-REQ divergence before acceptance: when an H1 is immediately followed by a Recommended: line, the drawer's existing question-option preprocessor appends a literal backslash to that H1. The raw AST comparison suppressed the H1 while the actual drawer retained it. Reviewer reproduced the failing Go assertion and observed the browser's retained heading. Builder is adding the shared corpus case first and reconciling this with the existing preprocessor; no prefix table or new body scanner is permitted.

The remediation scope permits a conditional transformed-heading fragment parse only when the existing question-option preprocessor changes the prefix through that heading. This is a demonstrated rendering-context exception to the initial no-extra-parse plan, not another citation scan. It must preserve later reference-link definitions and both cases where the transformed H1 is retained and where its added literal backslash matches the model title. A conservative blanket suppression guard would still violate parity.

## Qualification

Five-source implementation matches the amended Scope. Go module changes were declared before edits after a demonstrated Unicode mismatch. Existing raw-payload, UTF-16 and citation-search behavior is protected. The conditional fragment parse is confined to a heading affected by the preexisting question-option preprocessor and keeps original offset computation.

## Orientation

Saving a copied ticket back to disk no longer breaks duplicate-heading suppression. The drawer and Copy agree about which visible prose occurrence first receives the ticket title.

### Final remediation verification

Parent independently reran the final 21-case Go/JavaScript corpus and actual141 copy/save/rebuild test: exit 0, 3.920s. Builder's final broad Go/JavaScript regression run exited 0 (33.650s). Qualification passed after formatting the dependency name as prose rather than a manifest path. The earlier canonical pass predates the R1 correction and is not used as final acceptance evidence; the final canonical run is separate.

## Review

**Independent final review: Pass, 100%.** All four scored dimensions 100%; risk Low. R1 was an Important user-visible parity finding and is closed by the persisted 21-case corpus and independent actual141 rerun (exit 0, 4.532s). Reviewer confirmed unchanged raw bytes/UTF-16 offsets, citation independence, copied reference context, and actual round trips for the three captured records. Scope/P-A-U checked; restatement sweep found no newly stale contract. No open findings or follow-up REQs.

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
