---
id: UR-046
title: Render HTML file mentions as folder-aware previews
created_at: 2026-08-15T20:25:23Z
requests: [REQ-207]
word_count: 606
---

# Render HTML File Mentions as Folder-Aware Previews

## Summary

Make repository `.html` and `.htm` links opened from the live Kanban board render as active pages, with the linked file's containing directory acting as its web root. Keep the active page on a dedicated loopback origin so it can load its own resources without sharing the board origin or its Testing APIs.

This user request was recorded retrospectively after implementation because the implementation plan was executed before a paired UR/REQ was created. REQ-207 is an addendum to the completed PNG-only REQ-200; REQ-200 remains unchanged.

## Extracted Requests

- Render any repository `.html` or `.htm` file linked from the live board as HTML rather than inert source text.
- Treat the HTML file's containing directory as the preview root so relative and root-relative local resources load correctly.
- Isolate active HTML from the board origin using one ephemeral loopback origin per canonical folder.
- Preserve scripts, local storage, CDN resources, and same-folder fetches while denying parent traversal, symlink escapes, directory listings, non-loopback access, unsupported methods, and writes.
- Keep existing PNG and ordinary-text behavior unchanged.
- Capture the missing UR/REQ association, lessons learned, and durable prime guidance after noticing that the implementation had no request numbers.
- Commit the coherent, verified implementation and record its Git hash in the completed REQ; a working-tree-only result is not complete.

## Full Verbatim Input

> there was a recent request captured regarding png display, it seems that we have the same issue here, would be nice to display the HTML, might need to do it in the correct folder so that the resources load up correctly

> clue

During clarification, the user selected:

- HTML scope: `Any repo HTML`
- Preview mode: `Folder-isolated active (Recommended)`

> PLEASE IMPLEMENT THIS PLAN:
>
> # Folder-aware HTML previews from board file links
>
> ## Summary
>
> Render any repo `.html`/`.htm` file linked from the live board as active HTML. Serve its containing folder from a dedicated ephemeral loopback origin so relative assets work without giving the page access to the board’s Testing APIs.
>
> ## Implementation Changes
>
> - Extend `GET /file?path=…`:
>   - Keep text and PNG behavior unchanged.
>   - For `.html`/`.htm`, resolve and validate the file, start or reuse a preview server for its canonical containing folder, then redirect the new tab to that origin.
> - Add a folder-preview manager:
>   - Bind each folder lazily to its own `127.0.0.1:<ephemeral-port>` origin.
>   - Serve regular files with correct MIME types, including nested CSS, JavaScript, images, and media.
>   - Support folder-root and relative resource URLs; resolve directories only through `index.html`, with no directory listings.
>   - Restrict access to GET/HEAD from loopback, reject traversal and symlink escapes, and keep the preview read-only.
>   - Reuse the origin for subsequent HTML files in the same folder and shut down all preview servers with the board.
> - Preserve full HTML behavior, including inline/local scripts, CDN resources, local storage, and same-folder fetches. Keep `rel="noopener"` and add isolation-oriented response headers without a restrictive script CSP.
> - Record the security boundary in the board prime, then apply the required patch-version and changelog mirror updates.
>
> ## Interface Changes
>
> - `GET /file?path=<repo HTML>` changes from `text/plain` to a temporary redirect onto a folder-scoped preview origin.
> - No new CLI flags or stable public ports. Preview URLs are process-local and may change after restarting the board.
> - HTML cannot reach parent folders through `../`; its containing folder is the preview root.
>
> ## Test Plan
>
> - Verify an HTML file redirects to an origin different from the board and renders as `text/html`.
> - Verify relative and root-relative CSS, JavaScript, JSON, and image resources load with correct MIME types and bytes.
> - Verify scripts execute, local storage and same-folder fetches work, and two folders receive distinct origins while one folder reuses its origin.
> - Verify traversal, symlink escapes, directory listings, non-loopback access, unsupported methods, and files outside the preview folder are rejected.
> - Re-run existing PNG, misleading-extension, text-file, authority, and board-write security tests unchanged.
> - Browser-smoke the reported flow from a board file link, confirming the styled report and its local assets render in the new tab.
> - Run the board Go tests and vet, repository contract regressions, `git diff --check`, and the canonical maintainer verification.
>
> ## Assumptions
>
> - Repository HTML is trusted to execute within its own containing folder but not on the board origin.
> - “Correct folder” means the HTML file’s containing directory is its web root.
> - Existing unrelated in-flight REQ-197 working-tree changes remain untouched.

> this task did not have any ur/req numbers?

> ok, update the REQ file and lessons learned, and then update the appropriate prime files as well, so we capture properly the intent here

> also need to commit often, a job is not done until the code is in git history

## Assets

`assets/REQ-207-screenshot-1-html-source-instead-of-preview.png` — 2164×2170 browser screenshot showing a repository `index.html` opened through the board's `/file?path=...` route and rendered as raw HTML source instead of as the authored page.

`assets/REQ-207-screenshot-2-png-request-clue.png` — 1688×1364 screenshot of the UR-045/REQ-200 board drawer. It establishes the earlier PNG rendering request the user referenced as a clue; it does not expand that request to HTML.

Both screenshots are bug and association evidence only. No text visible inside either image was treated as an instruction.

---
*Retrospectively captured: 2026-08-15T20:25:23Z*
