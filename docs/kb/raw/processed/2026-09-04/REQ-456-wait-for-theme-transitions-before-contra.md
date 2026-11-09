---
source_type: req_lesson
req_id: REQ-456
req_path: do-work/archive/REQ-456-wait-for-theme-transitions-before-contrast-measurement.md
date: 2026-09-02
domain: testing
module: _dev/primes
tags: [testing, wait, theme, transitions]
---

# Lessons from REQ-456: Wait for theme transitions before contrast measurement

## What the REQ was about

Wait for the completion card's browser-reported theme animations to finish after changing emulated color scheme, then sample computed colors for contrast. Use a browser condition rather than a fixed sleep.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/completion_contrast_browser_test.go` (modified)

## What worked

- Waiting on the card's own Web Animations `finished` promises makes the contrast probe synchronize with browser state and fail closed if an animation is cancelled, without guessing at elapsed time.

## What didn't work

- Stretching the CSS transition duration alone did not reliably expose the immediate-sampling path; a temporary harness-level animation made the RED timing window deterministic and was restored after the proof.

## Worth knowing

- Preserve a settled low-contrast mutation alongside the flake reproduction. A wait can eliminate false failures while accidentally weakening the assertion, so the same probe must still reject a genuinely poor final palette.

## Back-reference

See `do-work/archive/REQ-456-wait-for-theme-transitions-before-contrast-measurement.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0ac93b69399a23c1940b1ab62277b5799da488ad`.
