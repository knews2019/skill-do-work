---
source_type: req_lesson
req_id: REQ-376
req_path: do-work/archive/UR-074/REQ-376-raise-the-done-lines-faint-text-to-readable-contrast.md
date: 2026-08-27
domain: ui-design
module: _dev/primes
tags: [ui-design, raise, done, line, faint]
---

# Lessons from REQ-376: Raise the done line''s faint text to readable contrast

## What the REQ was about

The done line's faint companions — `.relative-time` and `.elapsed-duration` — measure roughly 3.3:1 against `<body>` in both themes at 11px, under the 4.5:1 needed for text that size.

## Solution summary

Completion relative-time and elapsed-duration text now clears 4.5:1 against both the actual card and page backgrounds without changing pending or claimed text.

## Worth knowing

The completion-looking CSS class also appears on pending and claimed lines; scope a done-card change by status and keep a nonterminal control. The cards are opaque even though the chart SVG is transparent, so measure the actual card as well as the body. Time-sensitive browser fixtures must remain inside the production Recently Done window.

## Back-reference

See `do-work/archive/UR-074/REQ-376-raise-the-done-lines-faint-text-to-readable-contrast.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8dfdb24`.
