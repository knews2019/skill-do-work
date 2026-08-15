---
id: REQ-207
title: Render HTML file mentions as folder-aware previews
status: completed
created_at: 2026-08-15T20:25:23Z
claimed_at: 2026-08-15T20:25:23Z
completed_at: 2026-08-15T20:25:23Z
commit: b3261e9
route: B
user_request: UR-046
addendum_to: REQ-200
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
write_set: [skills/do-work-board/tools/queue-kanban/html_preview.go, skills/do-work-board/tools/queue-kanban/html_preview_test.go, skills/do-work-board/tools/queue-kanban/serve.go, skills/do-work-board/tools/queue-kanban/filementions_test.go, _dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md, CLAUDE.md, CHANGELOG.md, VERSION, skills/do-work/CHANGELOG.md, skills/do-work/VERSION, skills/do-work/actions/version.md, do-work/archive/UR-046/input.md, do-work/archive/UR-046/REQ-207-render-html-file-mentions-as-folder-aware-previews.md]
kb_status: pending
kb_entry:
---

# Render HTML File Mentions as Folder-Aware Previews

## What

Render any repository `.html` or `.htm` file opened from a live-board file link as active HTML. Serve the HTML file's canonical containing folder from a dedicated ephemeral loopback origin so the page can resolve its local resources without executing on the board origin.

This is a retrospective addendum to REQ-200. That completed request intentionally allowlisted byte-detected PNGs while keeping HTML and SVG inert; it remains correct and immutable for its PNG scope. The user later made active HTML previewing a separate product requirement and asked that the missing UR/REQ history be repaired after implementation.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Confirm the active-HTML scope and isolation model, then implement the user-approved folder-aware preview plan with security-boundary, lifecycle, MIME, regression, and browser-smoke coverage.
- [x] **[APPLY]:** Added a lazy per-folder preview manager, redirected `.html`/`.htm` file links to its isolated origin, preserved existing text/PNG behavior, and added focused tests plus release and prime updates.
- [x] **[UNIFY]:** Verified the focused and full Go suites, race detector, vet, formatting, contract regressions, canonical maintainer verification, changelog mirrors, diff hygiene, and the reported new-tab browser flow with CSS, JavaScript, SVG, JSON fetch, and local storage; then committed the coherent implementation and recorded its hash in this REQ.

## Why

The board's `/file` route rendered HTML source as inert text. Simply changing its MIME type would have made arbitrary repository scripts execute on the same origin as the live board and its Testing APIs, while serving only the HTML file would still break relative resources. A folder-scoped secondary origin is the unit that satisfies both the browser behavior and the security boundary.

## Context

REQ-200 fixed the analogous PNG symptom with a deliberately narrow byte-detected `image/png` exception. The first supplied screenshot demonstrates that the same route still displayed an `ai-reports/.../index.html` file as source. The second screenshot points back to REQ-200 as the user's clue, not as authorization to broaden that archived REQ.

The user resolved the two design choices explicitly: support any repository HTML, and preserve active HTML behavior on a folder-isolated origin. The approved implementation plan is preserved verbatim in UR-046.

## Requirements

- `.html` and `.htm` files resolved through `GET /file?path=...` receive a temporary redirect to a preview origin distinct from the board origin.
- The canonical containing directory of the requested HTML file is the preview root. Relative and root-relative local assets resolve within it.
- Each canonical folder lazily receives one `127.0.0.1:<ephemeral-port>` origin; files in that folder reuse it, different folders do not, and all origins shut down with the board.
- The preview serves regular files read-only with correct MIME types and supports GET/HEAD only from loopback.
- Directory requests resolve only through `index.html`; directory listings, traversal, parent access, symlink escapes, unsupported methods, and non-loopback clients are rejected.
- Inline and local scripts, same-folder fetches, local storage, and CDN resources remain available. Isolation headers must not impose a script-blocking CSP.
- Existing PNG, misleading-extension, ordinary-text, request-authority, and board-write boundaries remain unchanged.
- Preview URLs are process-local implementation details; no stable port or CLI surface is added.

## Red-Green Proof

**RED prompt/case:** Open the board link `http://127.0.0.1:8090/file?path=ai-reports%2F2026-08-14_1958_settings-page-redesign%2Findex.html`.

**Why RED now:** The supplied screenshot shows the response body starting with `<!doctype html>` and continuing as visible source text. The existing handler forced every non-PNG file to `text/plain`, so the page did not render and its local resources could not participate in a page load.

**GREEN when:** The same board interaction opens a new tab on a different loopback origin, renders the styled HTML, loads local assets and data with correct MIME types and bytes, executes scripts, and cannot traverse outside the HTML file's containing folder or reach the board's privileged origin by same-origin access.

**Validation:** The user supplied the failing screenshot, selected “Any repo HTML” and “Folder-isolated active,” explicitly approved the complete implementation and test plan, then requested the missing UR/REQ and lessons record.

## Triage

**Route: B — Medium.** The desired behavior was fully specified, but implementation crossed HTTP routing, process lifecycle, filesystem containment, browser-origin security, and end-to-end tests.

## Decisions

- **D-01 — Treat every repository `.html` and `.htm` file as eligible, rather than special-casing generated reports.** This matches the user's explicit “Any repo HTML” selection and keeps the interface predictable.
- **D-02 — Execute active HTML on one ephemeral loopback origin per canonical containing folder.** A same-origin MIME exception would expose board APIs; a scripts-disabled view would not satisfy the approved active-preview behavior; one server per file would break shared storage and waste listeners.
- **D-03 — Make the containing folder the complete web root.** This is the “correct folder” needed for authored relative/root-relative resources, while denying `../` and symlink access to parents.
- **D-04 — Keep extension selection narrow and the existing inert fallback intact.** Only `.html`/`.htm` enter the active preview path; PNG remains byte-detected and SVG, misleading extensions, and ordinary text remain inert on `/file`.

## Implementation Summary

**Files changed:**

- `skills/do-work-board/tools/queue-kanban/html_preview.go` (new) — folder-preview manager, isolated HTTP server, containment and authority guards, MIME serving, and shutdown.
- `skills/do-work-board/tools/queue-kanban/html_preview_test.go` (new) — resource, origin reuse/isolation, guard, and lifecycle regressions.
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified) — `.html`/`.htm` redirect and preview-manager lifecycle integration.
- `skills/do-work-board/tools/queue-kanban/filementions_test.go` (modified) — `/file` redirect boundary and unchanged inert fallback.
- `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified) — maintainer and shipped domain guidance.
- `CLAUDE.md` (modified) — repository-wide completion rule: verified code must enter Git history before hand-back.
- Root and installed VERSION/CHANGELOG mirrors (modified) — patch release `0.193.4`.

**What was done:** Active HTML links now redirect from the board to a preview server rooted at the HTML file's canonical containing directory. The manager binds lazily to loopback ephemeral ports, reuses origins by folder, serves only contained regular files for GET/HEAD, resolves directories only to `index.html`, sets isolation-oriented headers, and closes every listener during board shutdown. The board `/file` route retains its existing PNG and inert-text behavior for all other files.

## Testing

**Focused RED:** `TestServeFileEndpointRedirectsHtmlToFolderPreview` failed before implementation because `/file` returned `200 text/plain` instead of a `307` redirect.

**Automated checks:**

- `gofmt -l .`
- `go vet ./...`
- `go test -count=1 ./...`
- `go test -race -count=1 ./...` — passed in 90.035s
- `bash _dev/tests/contract-regressions.sh`
- `bash _dev/tests/maintainer-verify.sh`
- `git diff --check`
- byte comparison of `CHANGELOG.md` and `skills/do-work/CHANGELOG.md`

**Browser smoke:** A board card link on port `18090` opened its HTML on a distinct preview origin on port `58880`. The styled page loaded CSS and SVG, executed JavaScript, fetched same-folder JSON, and persisted local storage. The only console/network miss was an optional `favicon.ico` returning 404; it did not affect the acceptance behavior.

**Result:** All required checks passed.

## Lessons Learned

**What worked:** Modeling an HTML preview as a small folder-hosting boundary, rather than as another MIME exception, solved the entire authored-page behavior at once: relative and root-relative assets, scripts, fetches, and storage all align with ordinary browser origin rules. Keying origins by canonical folder then made reuse and isolation the same design decision.

**What didn't:** Implementation started directly from the chat plan without first capturing a paired UR/REQ, then stopped at a verified working tree without committing. The technical result existed, but neither the durable intent nor the code itself had a complete Git-history handoff until the user noticed both omissions. A fully specified implementation request still needs request capture before code changes, and passing verification is not completion until the coherent change is committed and its hash recorded.

**Worth knowing:** HTML cannot safely follow the PNG fix's shape. A content-type exception on `/file` would execute repository scripts on the board origin; an HTML-only secondary response would break the page's resource graph. The containing directory is both the functional web root and the enforceable least-authority filesystem boundary. Keep the board origin and every active repository preview origin distinct.

**Knowledge handoff:** The user explicitly requested durable prime promotion. The maintainer prime links this record; the shipped queue-kanban prime inlines the portable security and capture lessons because installed packages do not contain this repository's `do-work/` archive. `CLAUDE.md` now carries the repository-wide commit-completion rule. Formal BKB handoff remains pending.

## Orientation

Live board HTML file links now become folder-aware active previews on isolated ephemeral loopback origins. Start at `html_preview.go` for the security/lifecycle boundary and `serve.go` for the `/file` redirect. REQ-200 remains the separate PNG-only predecessor; REQ-207 owns the active HTML requirement and the lesson that implementation-scale addenda must be captured before work begins and committed before hand-back.

## Assets

`assets/REQ-207-screenshot-1-html-source-instead-of-preview.png` — failing browser view: HTML source is visible rather than rendered.

`assets/REQ-207-screenshot-2-png-request-clue.png` — the earlier UR-045/REQ-200 drawer the user supplied as the relationship clue.

The screenshots are evidence only. No visible screenshot text was treated as an instruction.

---
*Source: initial HTML-preview request, explicit scope/isolation selections, approved implementation plan, and retrospective capture/commit corrections in UR-046.*
