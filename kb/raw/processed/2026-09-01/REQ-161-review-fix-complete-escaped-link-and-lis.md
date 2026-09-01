---
source_type: req_lesson
req_id: REQ-161
req_path: do-work/archive/UR-031/REQ-161-complete-escaped-link-and-list-paragraph-classification.md
date: 2026-08-10
domain: general
module: skills/do-work/general
tags: [general, review, complete, escaped, link]
---

# Lessons from REQ-161: Review fix: Complete escaped-link and list-paragraph classification

## What the REQ was about

Complete the shared rendered-versus-ignored classification for the remaining escaped inline-link delimiter shapes and list-item paragraph continuations. Done means the approved Markdown class cannot recur through a different structural delimiter or block context; this remains a test-only release-guard correction and does not change the downstream publication policy.

## Solution summary

The test-only shipped-reference guard now hides complete escaped link-shaped regions from every target-discovery path and keeps four-column continuations rendered inside nonempty bullet and ordered-list paragraphs, without changing publication policy or downstream target resolution.

## What worked

- Exact production-helper RED cases closed the four approved defects without changing downstream publication or target-resolution policy.
- Pairing odd/even parity and rendered/code controls made the shared-mask intent directly testable while preserving offsets.

## What didn't work

- First-party URL controls could pass through the fallback even when structural relative-link extraction was still broken, so they overstated parity completeness.
- Bare list-fence controls did not exercise attached info strings or escaped opening brackets nested inside live labels; the independent matrix exposed both blind spots.

## Worth knowing

- Parser completeness needs paired relative and first-party targets for every delimiter form, plus context-sensitive controls that distinguish a structural escape from the same punctuation inside label content. REQ-163 records the remaining bounded variants for consent.

## Back-reference

See `do-work/archive/UR-031/REQ-161-complete-escaped-link-and-list-paragraph-classification.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ad3f8bd`.
