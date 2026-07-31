---
id: REQ-051
title: Add a generate_test.go substring assertion covering the board overlap-badge render path
status: completed
route: A
claimed_at: 2026-07-29T13:55:51Z
commit: e7a1941
domain: testing
tdd: false
maintenance: false
prime_files: []
created_at: 2026-07-29T09:32:07Z
user_request: UR-007
addendum_to: REQ-034
depends_on: []
write_set: ["tools/queue-kanban/generate_test.go"]
---

# Add a generate_test.go substring assertion covering the board overlap-badge render path

## What

Add a substring assertion to `tools/queue-kanban/generate_test.go` verifying the badge render path survives into the generated board output — anchor on `badge-write-overlap` / `writeSetOverlaps` in the inlined board JS. Same style as the file's existing neighbors.

## Why

The overlap badge's entire frontend render path currently has zero automated coverage — the Go tests cover only the `annotateWriteSetOverlap` annotation. A refactor that drops the badge renderer would ship as a silent regression. Approved by the user via `do-work clarify` on 2026-07-29 (follow-up from REQ-034, surfaced in REQ-041).

## Acceptance

- A test in `tools/queue-kanban/generate_test.go` fails if the badge render tokens disappear from the generated output.
- `go test ./...` passes in `tools/queue-kanban/`.

## Triage

**Route:** A
**Reasoning:** Single Go test file (`generate_test.go`); add a substring assertion on the badge render tokens, matching existing neighbors. Route A.
**Rigor:** Standard main-context review (part of the parallel disjoint-write_set batch 051/052/054/057/058; single-file, no spec-cluster overlap).

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the REQ, `crew-members/general.md`, `crew-members/coding-guardrails.md`, `crew-members/testing.md` (domain: testing), and all of `tools/queue-kanban/generate_test.go` to pick the neighbor style. Located the render path: `web/board.js:535-553` (the `makeBadge("badge-write-overlap", "overlaps", …)` card badge gated on `request.writeSetOverlaps`), `web/board.js:1726-1729` (the "Overlapping write sets" drawer row), and `web/board.css:837` (the `.badge-write-overlap` rule). Confirmed via `generate.go:20,368-381` that `web/board.css` and `web/board.js` are both `//go:embed`-ed and inlined into `index.html` at fixed placeholders, so `generateLiveSite(t)` sees the tokens. Approach: one new test at the end of the file, matching `TestRecentlyDoneWindowDefaultsTo24h` (the nearest neighbor — also asserts on inlined board.js tokens) for the comment-then-named-token shape, and matching `TestGenerateWritesSelfContainedIndex` / `TestGenerateRendersColumnHeaders` for the `for … range []string{…}` marker-loop shape. Anchor on **code** tokens rather than rendered DOM, because the badge only renders when the live queue happens to contain overlapping REQs — a DOM assertion would be data-dependent and would silently `Skip`/pass on a clean queue.
- [x] **[APPLY]:** Added `TestGenerateInlinesWriteSetOverlapBadgeRenderPath` to `tools/queue-kanban/generate_test.go` (+28 lines, append-only). Four anchor tokens: `` `"badge-write-overlap"` `` (quoted form — occurs only in the inlined board.js `makeBadge` call, whereas the bare class name would also match the CSS rule and pass on a JS-only regression), `request.writeSetOverlaps` (the generated-payload field the badge gates on), `Overlapping write sets` (the drawer row), and `.badge-write-overlap` (the stylesheet rule — without it the badge renders unstyled). Scope strictly limited to the one write_set file; no production code, no other test file, no `web/` asset touched.
- [x] **[UNIFY]:** Red-green proved without editing any file outside the write_set: temporarily mutated this test's own first anchor to `"badge-write-overlap-REMOVED"` → `--- FAIL: TestGenerateInlinesWriteSetOverlapBadgeRenderPath (0.10s) … not found in the inlined board.js/board.css`, then restored → `--- PASS`. This proves the assertion actually reads generated output and fails when a render token is absent, so it is not vacuous. Full suite: `go test ./...` → `ok github.com/knews2019/skill-do-work/queue-kanban 1.206s`. `go vet ./...` → clean (no output). `gofmt -l .` → clean (no output). `git diff --stat` → `tools/queue-kanban/generate_test.go | 28 ++++`, 1 file changed, 28 insertions(+); `git status --porcelain --untracked-files=all -- tools/` shows only ` M tools/queue-kanban/generate_test.go`, confirming no stray artifacts (the throwaway green-copy lived in the session scratchpad, outside the repo). No frontmatter, git, or commit action taken.

## Implementation Summary

**What was done:**

1. **Added `TestGenerateInlinesWriteSetOverlapBadgeRenderPath`** to `tools/queue-kanban/generate_test.go`, appended after `TestRecentlyDoneWindowDefaultsTo24h`. It calls the existing `generateLiveSite(t)` helper (builds the board against the real `do-work/` tree with the stubbed git lookup, returns the generated `index.html`) and loops over four required substrings, failing with the missing token named.
2. **Anchored on code tokens from the inlined assets, not rendered DOM.** The overlap badge only appears on a card when the live queue contains REQs with overlapping `write_set` declarations, so any DOM-level assertion would pass or skip on a clean queue and provide no regression protection. The four tokens come from the `//go:embed`-ed `web/board.js` and `web/board.css` that `assembleStaticPage` inlines into `index.html`, so the assertion holds regardless of queue contents.
3. **Chose the quoted `"badge-write-overlap"` form deliberately** for the card-badge anchor. The bare class name also occurs in the `web/board.css` rule, so a bare-name assertion would still pass if the `makeBadge` call were deleted from `board.js` but the dead CSS rule left behind — exactly the silent regression this REQ exists to catch. The CSS rule is covered separately by the `.badge-write-overlap` token.
4. **Closed the coverage gap named in the Why.** `model_test.go` already covers the Go-side `annotateWriteSetOverlap` annotation (lines 502, 636-639); nothing covered whether the derived list still reached the page. This test covers the frontend half without duplicating the annotation coverage.
5. **No production code changed.** This is a test-only addition; the render path was already correct and the test is green against the current tree.

**Files changed:**

- `tools/queue-kanban/generate_test.go` (modified)

**Verification receipts:**

- `go test ./...` → `ok  github.com/knews2019/skill-do-work/queue-kanban  1.206s`
- `go vet ./...` → clean
- `gofmt -l .` → clean
- Red-green: mutated anchor → FAIL; restored anchor → PASS
