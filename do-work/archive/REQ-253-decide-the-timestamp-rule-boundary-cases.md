---
id: REQ-253
title: Decide the Timestamp rule's two uncovered stamp shapes
status: completed
created_at: 2026-08-18T13:56:12Z
claimed_at: 2026-08-18T19:12:47Z
completed_at: 2026-08-18T19:31:15Z
kb_status: pending
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

---

## Implementation Summary

**What was done:** Both user-answered decisions implemented; the Timestamp rule did not grow a third paragraph. (1) The rule's date-only paragraph now names ui-review's report-header `**Date**` as a deliberately-UTC date-only consumer; the site cites the rule above the report fence and its placeholder reads `[today's UTC date]`. (2) One sentence in the same paragraph declares `## HH:MM UTC` time-of-day headings out of the rule's scope; the canonical statement lives in `memory-reference.md` § Daily-Log Entry Conventions, and all five write/template sites carry a uniform greppable marker ("time-of-day label, outside the Timestamp rule's scope"). The class was closed by grep over the corpus, not by the REQ's stale line list — the REQ's `memory.md:140` was a read-only checklist line; the real write sites (memory.md 50/84, memory-reference.md 46/93/135) are the marked set.

**Files changed (4, +9/−7):**
- `skills/do-work/actions/work-reference.md` (modified) — date-only paragraph: ui-review named as UTC consumer; out-of-scope sentence for time-of-day headings
- `skills/do-work-toolbox/actions/ui-review.md` (modified) — rule citation above the report fence; placeholder `[today's UTC date]`
- `skills/do-work-knowledge/actions/memory.md` (modified) — markers at the two heading write sites
- `skills/do-work-knowledge/actions/memory-reference.md` (modified) — canonical out-of-scope statement; markers at two more sites

*Integrated by orchestrator from builder hand-back; merge range `3114ca2..0d8d629`.*

## Decisions

Transcribed from the builder hand-back:

- **D-01 (DECIDE & STATE):** mark write/template sites by content class ("writes or templates `## HH:MM UTC`"), not the REQ's stale line list — the REQ's named line was a read site.
- **D-02 (DECIDE & STATE):** read-only mentions stay unmarked — a sweep has nothing to convert at a read, and the rule's walk-past sentence covers bare-occurrence greps; marking reads would teach that every mention needs annotation.
- **D-03 (DECIDE & STATE):** canonical marker lives once in Daily-Log Entry Conventions; site markers are terse uniform pointers.

## Qualification

Passed — 4 files in merge range `3114ca2..0d8d629` (+9/−7), both decisions traced to their sites, P-A-U audited (single scripted pass with uniqueness assertions; text-only diff, no debug artifacts). Orchestrator independently re-verified: the rule paragraph names ui-review and carries the out-of-scope sentence; the placeholder reads `[today's UTC date]`; five marker occurrences present; class grep over `skills/` shows only the rule, the marked sites, history, and the two read sites.

## Review

**Overall: 95%** | 2026-08-18T19:35:00Z (Route A quick scan, orchestrated inline)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

- [x] Both confirmed answers implemented exactly — UTC header (rule line + site citation + placeholder), out-of-scope declaration with marked sites, no third paragraph. Verified by reading the merged text, not the hand-back.
- [x] Citation depth follows REQ-249's literal rule and its RED was demonstrated against the real checker (wrong depth → named FAIL, exit 1); the reference contract passes on the merged tree.
- [x] Timestamp contract counts unchanged at 54/17 and the invariance is *explained*, not just observed: the site was already recognized date-only, and the new placeholder still matches the checker's date vocabulary.
- [x] D-01's line-list correction is the right call — the REQ inherited a read-site line number from REQ-244's occurrence grep; the class condition ("writes or templates the heading") is what closed it.

**Findings:** none Important. One Minor (report only): the rule's date-only paragraph is now a long single block carrying four distinct clauses (date-only consumers, no-tool-subcommand tripwire, local-date carve-out, time-of-day out-of-scope) — legible today, but the next addition should split it.

**Acceptance: Pass.** Follow-ups: none from findings; the two Discovered Tasks take the consent flow (REQ-261, REQ-262).

*Reviewed by review-work action (orchestrated inline, Route A depth; merge range `3114ca2..0d8d629`)*

## Lessons Learned

**What worked:** Closing the class by condition ("writes or templates the heading") instead of the REQ's line list — the listed line turned out to be a read site, and the real write sites were elsewhere. Uniqueness-asserted scripted edits made six small text changes safe in one pass.

**Worth knowing:** The `## HH:MM UTC` shape is defined once in memory-reference.md § Daily-Log Entry Conventions; the site markers are pointers to it. The date-only paragraph now carries a tripped tripwire — it says "revisit if a second consumer appears" and ui-review is that second consumer (REQ-261 asks the question).

## Orientation

Now every clock write in shipped action text is governed: instants by the Timestamp rule, date-only stamps by its date-only paragraph (ui-review's header settled UTC), and time-of-day headings explicitly outside the rule with marked write sites. Lives in the Timestamp-rule convention layer across three packages. Leaf change; map unchanged.
