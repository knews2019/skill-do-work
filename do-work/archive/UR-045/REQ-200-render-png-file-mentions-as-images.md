---
id: REQ-200
title: Render PNG file mentions as images
status: completed
claimed_at: 2026-08-15T17:43:26Z
completed_at: 2026-08-15T17:58:49Z
route: A
created_at: 2026-08-15T17:38:31Z
user_request: UR-045
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
kb_status: pending
kb_entry:
---

# Render PNG File Mentions as Images

## What

Make the queue-kanban local file view display a referenced PNG as an image instead of rendering its binary bytes as page text.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add an end-to-end `/file` regression test with a valid PNG and confirm RED under the forced `text/plain` response; then allow only byte-detected PNG content to use `image/png`, preserving the text fallback and existing access guards for every other file.
- [x] **[APPLY]:** Added the failing PNG endpoint regression test first, then changed the handler to return `image/png` only for byte-detected PNG data while retaining `text/plain` for all other files.
- [x] **[UNIFY]:** Reviewed the scoped diff and `git diff --stat`; verified `serve.go` preserves every existing access guard and only selects the response type, verified `filementions_test.go` covers real PNG bytes plus a misleading `.png` extension, ran `gofmt`, `go vet ./...`, and `go test -count=1 ./...`, and found no debug artifacts in either changed file.

## Why

The referenced asset is an image and should be shown as one.

## Context

The supplied screenshot shows the local `/file?path=...png` route displaying PNG bytes, chunk labels, and embedded metadata as text in the browser rather than rendering the image.

## Red-Green Proof
**RED prompt/case:** Open a queue-kanban `/file?path=...png` link for a valid PNG asset; the browser displays the PNG's binary contents as text.
**Why RED now:** The attached screenshot shows the response beginning with a corrupted `PNG` signature and exposing binary chunk data across the page instead of displaying the asset.
**GREEN when:** Opening the same kind of local file link renders the referenced PNG as an image in the browser and does not display its binary bytes as text.
**Validation:** User confirmed through the screenshot and the explicit statement that it should have shown the image instead.

## Assets

`do-work/user-requests/UR-045/assets/REQ-200-screenshot-1-render-png-file-mention.png` — browser screenshot of the failing `/file` rendering. The address bar targets a `.png` file under a captured request's `assets/` directory, while the page body shows raw PNG data rather than an image.

---
*Source: "do-work capture-request: this should have shown the image instead" with attached screenshot.*

---

## Triage

**Route: A** - Simple

**Reasoning:** This is a focused bug with a concrete browser reproduction and an existing Go handler test surface.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01 — Detect PNG content from bytes and allowlist only `image/png`.** The user supplied a PNG regression, while the existing forced-text behavior deliberately prevents active HTML/SVG from executing in the board origin. Byte detection fixes the requested rendering without letting a misleading extension opt into executable content.

## Root Cause

The live `/file` handler intentionally hardcoded `Content-Type: text/plain; charset=utf-8` for every repository file. That security default also covered inert raster images, so browsers rendered valid PNG bytes as text instead of decoding the image.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified)
- `skills/do-work-board/tools/queue-kanban/filementions_test.go` (modified)

**What was done:** The live file endpoint now returns `image/png` only when Go's byte-level content detection identifies valid PNG data, retaining its inert text response for every other file. An end-to-end regression test proves valid PNG rendering, unchanged bytes, and the safe fallback for a mislabeled `.png` file.

## Qualification

Passed — 2 files verified, 1 requirement traced, P-A-U confirmed. The handler change is substantive and limited to response type selection after the existing containment/regular-file/size/read checks; the response body still flows from the resolved file bytes without placeholders or hardcoded output.

## Testing

**Tests run:** `go test -count=1 -run '^TestServeFileEndpointRendersPng$'`; focused `/file` endpoint suite; `gofmt -l .`; `go vet ./...`; `go test -count=1 ./...`
**Result:** ✓ All passing

**Red-green validation:**
- `TestServeFileEndpointRendersPng` in `filementions_test.go`: ✗ before implementation (`Content-Type = "text/plain; charset=utf-8", want image/png`) → ✓ after implementation

**New tests added:**
- `TestServeFileEndpointRendersPng` — proves valid PNG bytes receive `image/png`, remain byte-identical, and a `.png` filename containing active SVG markup stays inert `text/plain`.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-15T17:57:47Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None

**Minor findings:** 0
**Acceptance:** Pass — the built live server returned `image/png` and byte-identical content for the preserved screenshot, retained `text/plain` for Markdown, and kept the existing security headers and path guards.
**Suggested testing:** 1 optional item — visually open a large captured PNG in a browser; the exact asset's HTTP headers and bytes were already verified end-to-end.
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Byte-level PNG detection created the smallest safe exception to the file view's inert-text rule; a real encoded PNG plus a misleading `.png` fixture proved both sides of the boundary.

**What didn't:** Updating the handler alone left one test diagnostic restating the old “never the file's own type” contract; the restatement sweep caught it before release.

**Worth knowing:** The `/file` route applies `X-Content-Type-Options: nosniff` globally and deliberately keeps HTML/SVG and mislabeled files as `text/plain`. Any future inline format must be explicitly allowlisted with a regression test for both valid bytes and a misleading extension.

**Knowledge handoff:** Pending explicit user consent. No knowledge-base file was written automatically.

## Orientation

Captured PNG assets opened from the live Kanban board now render inline; the same `/file` boundary keeps every non-PNG response inert text.
