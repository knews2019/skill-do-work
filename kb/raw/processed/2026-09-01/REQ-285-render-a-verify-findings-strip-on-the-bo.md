---
source_type: req_lesson
req_id: REQ-285
req_path: do-work/archive/UR-058/REQ-285-render-a-verify-findings-strip-on-the-board.md
date: 2026-08-21
domain: frontend
module: _dev/primes
tags: [frontend, render, verify, findings, strip]
---

# Lessons from REQ-285: Render a verify-findings strip on the board

## What the REQ was about

Add a second always-visible strip to the board client, modeled on the existing completion-anomalies
strip, that renders `verifyFindings` and `verifySkipped` from the board payload.

## Solution summary

Added a second always-visible strip to the board client rendering `verifyFindings` and `verifySkipped` from the payload, modelled on the completion-anomalies strip and deliberately subordinate to it. Each finding is a card with the category as a badge, the detail as body text, the remedy as a muted second line, and a `cleanup can fix` marker only when the producer set `fixable`. Skipped probes render in a collapsed `<details>` footer on the same strip. Fixed the path-reduction regex the first render exposed.

## What worked

**What worked:** Actually opening the page. Two defects were invisible to every other check and both were obvious on sight — a phrase breaking across two lines, and a mangled path in the skipped footer. `node --check` passed, `go test` passed, the strict JS lane passed, and the page was still wrong. For a rendering REQ the screenshot is not a nicety, it is the test.

**What didn't:** The first browser probe used `[data-view]` and found no view buttons at all, silently reporting an empty list instead of failing. The real attribute is `data-view-target`. A probe that finds nothing and says so calmly is indistinguishable from a probe that found nothing because there was nothing — the check should have asserted it found five buttons before using them.

**Worth knowing:** The path-reduction defect (D-02) is the more interesting one. `reduceAbsolutePaths` ran two passes — replace the repo root with a relative path, then strip anything still absolute — and the second pass ate the first pass's output, because RE2 has no lookbehind and a bare `/` matches mid-token. Any two-pass text reduction where the second pass's pattern can match the first pass's product has this shape. The ordering is load-bearing and the boundary capture is what makes it safe; both are commented at the regex.

## Back-reference

See `do-work/archive/UR-058/REQ-285-render-a-verify-findings-strip-on-the-board.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `fed89c9`.
