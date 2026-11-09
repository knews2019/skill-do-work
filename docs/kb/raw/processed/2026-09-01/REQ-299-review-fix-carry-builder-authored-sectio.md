---
source_type: req_lesson
req_id: REQ-299
req_path: do-work/archive/UR-055/REQ-299-carry-builder-authored-decisions-past-step-8.md
date: 2026-08-21
domain: general
module: _dev/primes
tags: [general, review, carry, builder, authored]
---

# Lessons from REQ-299: Review fix: carry builder-authored sections past Step 8, starting with ## Decisions

## What the REQ was about

REQ-270 closed the case where a worktree builder's `## Discovered Tasks` never reached
Step 8, and keyed its new rule on the condition — but scoped that condition's home to
**Step 8's substeps** (`actions/work.md` → *Where a builder-authored section is read from*,
which opens "Some substeps **below**"). Its independent review then found a second
builder-authored section with the identical defect, read from **outside** Step 8, where the
new rule structurally cannot reach:

`## Decisions` is written by the builder (`actions/work.md` Step 6, "Log Decisions as
D-XX" — unqualified, so a worktree builder is told to write a file it may not write), and
read at two sites: `actions/review-work.md` Step 4's traceability check ("If a
`## Decisions` section exists in the REQ: verify that significant implementation
choices … are documented") and the end-of-run Decision Brief's HANDLED block
(`actions/work-reference.md` → **Decision Brief (hand-back format)**).

## Solution summary

**What changed in the map.** The rule for reading a section the builder authored is no
longer part of Step 8. It is a named section of `actions/work-reference.md` —
**Reading a Builder-Authored Section (any step)** — that any step, action, or report can
cite, and `actions/work.md` Step 8, `actions/review-work.md` Step 4, and the Decision
Brief all now cite it rather than restating it.

## Worth knowing

- **A rule scoped to a step cannot be inherited by readers at other steps.** REQ-270 wrote
  "Some substeps below" and closed one instance; the second instance was already in the file
  and outside that scope. When a rule's condition is "whenever anyone does X", its home has
  to be a place every doer of X reads, and its opening sentence has to name the condition
  rather than the location.
- **A check that needs a list can often be turned into a check that needs a mark.** The
  property "every section Step 6 tells the builder to author" resisted mechanical extraction
  until the requirement flipped: instead of the test knowing which sections those are, every
  section mention in the block states who writes it. The list moved into the prose, where it
  is co-located with what it describes and cannot be forgotten by a test author.
- **A guard passing is not a guard covering.** `shipped-package-reference-contract.sh`
  passed on planted dangling citations of three different same-package shapes. The only way
  to learn that was to plant them. Confirm a guard covers new text by breaking the text,
  never by reading the guard.
- **Mutation-test the pin, not just the fix.** M11 first came back green because the
  mutation removed one of two phrases the file carried. A mutation that does not actually
  remove the property proves nothing — check the mutation applied to everything the check
  can see.

## Back-reference

See `do-work/archive/UR-055/REQ-299-carry-builder-authored-decisions-past-step-8.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `3469b39`.
