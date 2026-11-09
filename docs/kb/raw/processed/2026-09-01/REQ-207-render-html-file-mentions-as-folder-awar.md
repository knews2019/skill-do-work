---
source_type: req_lesson
req_id: REQ-207
req_path: do-work/archive/UR-046/REQ-207-render-html-file-mentions-as-folder-aware-previews.md
date: 2026-08-15
domain: backend
module: _dev/primes
tags: [backend, render, html, file, mentions]
---

# Lessons from REQ-207: Render HTML file mentions as folder-aware previews

## What the REQ was about

Render any repository `.html` or `.htm` file opened from a live-board file link as active HTML. Serve the HTML file's canonical containing folder from a dedicated ephemeral loopback origin so the page can resolve its local resources without executing on the board origin.

This is a retrospective addendum to REQ-200. That completed request intentionally allowlisted byte-detected PNGs while keeping HTML and SVG inert; it remains correct and immutable for its PNG scope. The user later made active HTML previewing a separate product requirement and asked that the missing UR/REQ history be repaired after implementation.

## Solution summary

Active HTML links now redirect from the board to a preview server rooted at the HTML file's canonical containing directory. The manager binds lazily to loopback ephemeral ports, reuses origins by folder, serves only contained regular files for GET/HEAD, resolves directories only to `index.html`, sets isolation-oriented headers, and closes every listener during board shutdown. The board `/file` route retains its existing PNG and inert-text behavior for all other files.

## What worked

Modeling an HTML preview as a small folder-hosting boundary, rather than as another MIME exception, solved the entire authored-page behavior at once: relative and root-relative assets, scripts, fetches, and storage all align with ordinary browser origin rules. Keying origins by canonical folder then made reuse and isolation the same design decision.

## What didn't work

Implementation started directly from the chat plan without first capturing a paired UR/REQ, then stopped at a verified working tree without committing. The technical result existed, but neither the durable intent nor the code itself had a complete Git-history handoff until the user noticed both omissions. A fully specified implementation request still needs request capture before code changes, and passing verification is not completion until the coherent change is committed and its hash recorded.

## Worth knowing

HTML cannot safely follow the PNG fix's shape. A content-type exception on `/file` would execute repository scripts on the board origin; an HTML-only secondary response would break the page's resource graph. The containing directory is both the functional web root and the enforceable least-authority filesystem boundary. Keep the board origin and every active repository preview origin distinct.

## Back-reference

See `do-work/archive/UR-046/REQ-207-render-html-file-mentions-as-folder-aware-previews.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `b3261e9`.
