---
id: REQ-175
title: Align board question preprocessing with valid Markdown fences
status: pending-answers
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
**Validation:** Review-generated follow-up awaiting user confirmation.

## Open Questions

- [ ] The board currently mistakes an invalid backtick-info opener for a real code fence, which can merge later question options in the rendered card. Should I queue the focused renderer fix and regression test?
  Recommended: Yes, add to queue — the board will match Goldmark and preserve question-option line breaks.
  Also: No, discard it — REQ-174 remains complete, but this adjacent board edge case stays open.

---
*Source: Important adjacent finding from independent review of REQ-174*
