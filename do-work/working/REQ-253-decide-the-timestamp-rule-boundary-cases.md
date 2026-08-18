---
id: REQ-253
title: Decide the Timestamp rule's two uncovered stamp shapes
status: claimed
created_at: 2026-08-18T13:56:12Z
claimed_at: 2026-08-18T19:12:47Z
route: A
status_changed_at: 2026-08-18T14:12:05Z
user_request: UR-055
addendum_to: REQ-244
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-249]
maintenance: true
write_set:
- skills/do-work/actions/work-reference.md
- skills/do-work-toolbox/actions/ui-review.md
- skills/do-work-knowledge/actions/memory.md
- skills/do-work-knowledge/actions/memory-reference.md
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T19:13:28Z
  basis:
    - trivial short-circuit
---

# Decide the Timestamp Rule's Two Uncovered Stamp Shapes

## What

Two clock-write shapes in shipped action files are governed by neither paragraph of the Timestamp rule. Both need a decision before they can be swept, and neither decision is the pipeline's to make.

## Context

REQ-244 sweeps every timestamp write site under the rule. Its builder correctly declined to convert these two, because doing so would have decided a question the REQ did not ask; its reviewer confirmed both as genuinely uncovered rather than as misses.

## Open Questions

- [x] `skills/do-work-toolbox/actions/ui-review.md:216` writes `**Date**: [today]` into a human-facing report header. It is a clock write with no rule behind it, but nothing says whether that header wants a UTC date — matching every other stamp we write — or a deliberately local one, which is arguably friendlier at the top of a report a person reads. Which should it be? Once you say, it takes one line in the rule's date-only paragraph plus a citation at the site. → Confirmed: UTC, matching every other stamp, so there is one answer and no site-by-site judgment.
  Recommended: UTC, matching every other stamp, so there is one answer and no site-by-site judgment.
  Also: local date, on the grounds that a report header is for a human in their own timezone — in which case the rule should say so explicitly, because otherwise the next sweep converts it.

- [x] `skills/do-work-knowledge/actions/memory.md:140` and `memory-reference.md:46,93,135` write a `## HH:MM UTC` heading — a time of day with no date. The Timestamp rule has a paragraph for full instants and a paragraph for date-only stamps, and this is neither, so nothing governs it and a sweep has no rule to apply. Should the rule grow a third paragraph covering time-of-day stamps, or should these sites be declared out of the rule's scope and marked so a future sweep walks past them? → Confirmed: declare them out of scope and mark them, since a heading inside a dated file already carries its date from context.
  Recommended: declare them out of scope and mark them, since a heading inside a dated file already carries its date from context.
  Also: add a third paragraph, if you would rather the rule cover every clock write in the tree without exception.

**Answered [2026-08-18]:** User confirmed both recommendations via `do-work clarify`. The report header date goes UTC (one line in the rule's date-only paragraph plus a citation at the site); the time-of-day headings are declared out of the Timestamp rule's scope and marked at their sites so future sweeps walk past them — the rule does not grow a third paragraph.

## Requirements

- Each answer lands in `skills/do-work/actions/work-reference.md`'s Timestamp rule, in the paragraph the answer implies.
- Every affected site is then either cited under the rule or explicitly marked out of scope — no site is left in the state that produced this REQ, where a sweep reaches it and has nothing to apply.

## Ordering Gate

- [~] **D-01: gated behind REQ-249 rather than left free-running.** Both REQs edit the same shipped files — `skills/do-work/actions/work-reference.md` (this REQ adds paragraphs to the Timestamp rule; REQ-249 rewrites citation paths throughout), plus `ui-review.md`, `memory.md` and `memory-reference.md`, which this REQ marks and REQ-249 sweeps. Running them together is a guaranteed textual conflict in at least four files.
  **Ordering chosen: REQ-249 first.** Its sweep establishes which citation form is the rule and updates `_dev/primes/prime-action-files.md` to say so. This REQ then adds its two or three new citations already conforming. The reverse order has this REQ writing citations in a form that REQ-249 is about to rewrite — correct in the end, but it means writing something known-wrong on purpose.
  **Value:** the one-line resume cannot schedule them together, and the second one to run inherits a settled convention rather than a contested one.
  **Risk:** this REQ is small and now waits on a large sweep. Accepted — the user has already answered both of this REQ's questions, so nothing is lost but ordering.

---

## Triage

**Route: A** - Simple

**Reasoning:** Both decisions are user-answered; the work is one line in the Timestamp rule's date-only paragraph, a citation at ui-review.md's header, and out-of-scope markers at the memory time-of-day headings.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
