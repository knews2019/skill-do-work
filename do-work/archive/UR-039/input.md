---
id: UR-039
title: Fix three accepted P2 review findings
created_at: 2026-08-11T17:00:04Z
requests: [REQ-172, REQ-173, REQ-174]
word_count: 2
---

# Fix three accepted P2 review findings

## Summary

Implement the three findings accepted during the preceding `do-work validate-feedback` review: make post-install screenshot source cleanup retry-safe, recognize a UTF-8 BOM during reserved-recipe classification, and align root fence classification with Goldmark/CommonMark.

## Extracted Requests

| Request | Intent |
|---|---|
| REQ-172 | Do not invalidate a verified screenshot installation when staged-source cleanup fails. |
| REQ-173 | Detect BOM-prefixed reserved Just recipes in the no-Just collision path. |
| REQ-174 | Do not mask rendered links after an invalid root backtick-fence opener. |

## Batch Constraints

- Preserve the earned no-clobber screenshot installation, no-Just collision scanner, and rendered-reference classifier.
- Prefer the direct simplifications accepted during feedback triage; do not add compensating recovery layers.
- Add focused regressions for each reproduced failure.

## Full Verbatim Input

fix accepted

---
*Captured: 2026-08-11T17:00:04Z*
