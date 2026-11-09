---
source_type: req_lesson
req_id: REQ-200
req_path: do-work/archive/UR-045/REQ-200-render-png-file-mentions-as-images.md
date: 2026-08-15
domain: backend
module: _dev/primes
tags: [backend, render, file, mentions, images]
---

# Lessons from REQ-200: Render PNG file mentions as images

## What the REQ was about

Make the queue-kanban local file view display a referenced PNG as an image instead of rendering its binary bytes as page text.

## Solution summary

The live file endpoint now returns `image/png` only when Go's byte-level content detection identifies valid PNG data, retaining its inert text response for every other file. An end-to-end regression test proves valid PNG rendering, unchanged bytes, and the safe fallback for a mislabeled `.png` file.

## What worked

Byte-level PNG detection created the smallest safe exception to the file view's inert-text rule; a real encoded PNG plus a misleading `.png` fixture proved both sides of the boundary.

## What didn't work

Updating the handler alone left one test diagnostic restating the old “never the file's own type” contract; the restatement sweep caught it before release.

## Worth knowing

The `/file` route applies `X-Content-Type-Options: nosniff` globally and deliberately keeps HTML/SVG and mislabeled files as `text/plain`. Any future inline format must be explicitly allowlisted with a regression test for both valid bytes and a misleading extension.

## Back-reference

See `do-work/archive/UR-045/REQ-200-render-png-file-mentions-as-images.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `03b40a2`.
