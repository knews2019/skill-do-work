---
id: REQ-200
title: Render PNG file mentions as images
status: pending
created_at: 2026-08-15T17:38:31Z
user_request: UR-045
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
---

# Render PNG File Mentions as Images

## What

Make the queue-kanban local file view display a referenced PNG as an image instead of rendering its binary bytes as page text.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
