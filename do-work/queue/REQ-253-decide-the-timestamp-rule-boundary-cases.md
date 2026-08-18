---
id: REQ-253
title: Decide the Timestamp rule's two uncovered stamp shapes
status: pending-answers
created_at: 2026-08-18T13:56:12Z
status_changed_at: 2026-08-18T13:56:12Z
user_request: UR-055
addendum_to: REQ-244
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
---

# Decide the Timestamp Rule's Two Uncovered Stamp Shapes

## What

Two clock-write shapes in shipped action files are governed by neither paragraph of the Timestamp rule. Both need a decision before they can be swept, and neither decision is the pipeline's to make.

## Context

REQ-244 sweeps every timestamp write site under the rule. Its builder correctly declined to convert these two, because doing so would have decided a question the REQ did not ask; its reviewer confirmed both as genuinely uncovered rather than as misses.

## Open Questions

- [ ] `skills/do-work-toolbox/actions/ui-review.md:216` writes `**Date**: [today]` into a human-facing report header. It is a clock write with no rule behind it, but nothing says whether that header wants a UTC date — matching every other stamp we write — or a deliberately local one, which is arguably friendlier at the top of a report a person reads. Which should it be? Once you say, it takes one line in the rule's date-only paragraph plus a citation at the site.
  Recommended: UTC, matching every other stamp, so there is one answer and no site-by-site judgment.
  Also: local date, on the grounds that a report header is for a human in their own timezone — in which case the rule should say so explicitly, because otherwise the next sweep converts it.

- [ ] `skills/do-work-knowledge/actions/memory.md:140` and `memory-reference.md:46,93,135` write a `## HH:MM UTC` heading — a time of day with no date. The Timestamp rule has a paragraph for full instants and a paragraph for date-only stamps, and this is neither, so nothing governs it and a sweep has no rule to apply. Should the rule grow a third paragraph covering time-of-day stamps, or should these sites be declared out of the rule's scope and marked so a future sweep walks past them?
  Recommended: declare them out of scope and mark them, since a heading inside a dated file already carries its date from context.
  Also: add a third paragraph, if you would rather the rule cover every clock write in the tree without exception.

## Requirements

- Each answer lands in `skills/do-work/actions/work-reference.md`'s Timestamp rule, in the paragraph the answer implies.
- Every affected site is then either cited under the rule or explicitly marked out of scope — no site is left in the state that produced this REQ, where a sweep reaches it and has nothing to apply.
