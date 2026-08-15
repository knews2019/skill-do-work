---
id: REQ-175
title: Align board question preprocessing with valid Markdown fences
status: pending
status_changed_at: 2026-08-15T07:41:57Z
created_at: 2026-08-11T20:22:26Z
user_request: UR-039
addendum_to: REQ-174
review_generated: true
effort_estimate: normal
domain: testing
prime_files: [skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-174]
write_set: [skills/do-work-board/tools/queue-kanban/render.go, skills/do-work-board/tools/queue-kanban/generate_test.go]
---

# Align Board Question Preprocessing with Valid Markdown Fences

## What

Make the board's Open Questions hard-break preprocessor distinguish real fenced code blocks from invalid backtick-info lookalikes, matching the Markdown renderer it feeds.

## Why

REQ-174 fixed the shipped reference classifier, but independent review found the board renderer still toggles fence state for every trimmed triple-backtick prefix. A line such as `````lang`invalid`` is prose to Goldmark, yet the preprocessor treats it as a fence opener and can suppress hard breaks before later `Recommended:` or `Also:` lines.

## Requirements

- Reproduce the invalid backtick-info opener in a board-rendering test and prove question-option lines retain their intended visual breaks.
- Preserve verbatim handling inside valid backtick and tilde fences.
- Keep preprocessing behavior aligned with the pinned Goldmark renderer; do not add a second broader Markdown parser.
- Limit the fix to the board renderer and its focused regression coverage.

## Red-Green Proof

**RED prompt/case:** Render Open Questions prose containing an invalid backtick-fence opener followed by `Recommended:` and `Also:` option lines.
**Why RED now:** `insertQuestionOptionHardBreaks` toggles its `insideFence` state on a prefix alone, so Goldmark-visible option prose is skipped by preprocessing.
**GREEN when:** The invalid opener remains ordinary prose, the option lines receive hard breaks, valid fenced content stays byte-verbatim, and board tests pass.
**Validation:** User approved this review-generated follow-up via `do-work clarify` on 2026-08-15.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-174: the board can mistake invalid backtick-info prose for a real fence and merge later question options. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to `pending`).
  Also: No, discard it.

---
*Source: Important adjacent finding from independent review of REQ-174*
