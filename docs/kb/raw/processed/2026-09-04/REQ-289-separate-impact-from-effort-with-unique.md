---
source_type: req_lesson
req_id: REQ-289
req_path: do-work/archive/UR-060/REQ-289-separate-impact-from-effort-with-unique-tokens.md
date: 2026-08-19
domain: general
module: _dev/primes
tags: [general, separate, impact, from]
---

# Lessons from REQ-289: Separate impact from effort, with unique greppable tokens on both axes

## What the REQ was about

`effort_estimate` has two writers with two different meanings. Capture sets it as a size judgment;
review MUST-stamp it from an impact gate. Split the two axes into two fields, and give every value
on both axes a token that is unique repo-wide and findable by plain-text search.

## Solution summary

**Files changed:**

## What worked

- Mutation-testing the lock-in checks instead of accepting "the suite is green" — three of the four
- proved real on the first try, and the exercise is what made the review's counter-finding legible
- rather than arguable.
- Verifying backwards compatibility through the actual parser on purpose-built fixtures. The claim
- "the aliases carry legacy REQs" is unfalsifiable by grep and trivial to settle with a real parse.
- Extending the write set at Step 5.5 from the *planner's* survey rather than discovering the gaps
- mid-build. Nine of the twenty-one files were added before the builder started; two of them would
- have failed the build and one would have shipped a silent bug.

## What didn't work

- **The orchestrator's mutation test for Check A was self-confirming and it took the reviewer to see
- it.** The mutation used the verb "stamping" — the single literal the check greps. A green mutation
- test proved only that the check catches its own phrasing. A mutation must be written in words the
- check's author did not choose, or it tests the tester's vocabulary rather than the property.
- The plan's literal wording for Check A was unimplementable and the builder had to redesign it
- mid-build (D-03). Specifying a check by its *regex* rather than by the property it must hold pushed
- the design decision to the wrong moment.

## Worth knowing

- `timeline.go` compared `effort_estimate` against a bare `"trivial"` string, and `timeline_test.go`
- built its fixtures from the same literal. A rename would have left both compiling, both passing,
- and every REQ silently bucketed into the substantive median. When renaming an enum, grep the
- *value* across the whole module, not the constant name — the tests are as likely to hold the
- literal as the code is.
- Three separate vocabularies described one axis before this REQ (`gate:`, the `[critical]`/
- `[normal]`/`[low]` ladder, and `effort_estimate`'s double duty). The third consumer
- (`maintainability-audit-reference.md`) already stated this REQ's thesis in its own words — "trivial
- as a gate value routes the review flow and never doubles as a severity" — while carrying the
- conflation. A file arguing against a problem is not evidence it avoided it.
- The `impact:` default must stay `impact-user-visible`. Absence must never be mistakable for the
- user's stop signal, because REQ-290's filter acts on that signal.

## Back-reference

See `do-work/archive/UR-060/REQ-289-separate-impact-from-effort-with-unique-tokens.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `2ea7be5`.
