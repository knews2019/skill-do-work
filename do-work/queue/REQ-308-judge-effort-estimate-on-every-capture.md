---
id: REQ-308
title: "[impact-rule-change] Judge effort_estimate on every capture, as impact already is"
status: pending
created_at: 2026-08-20T22:00:52Z
user_request: UR-064
domain: general
impact: impact-rule-change
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/capture.md
- skills/do-work/actions/capture-reference.md
- _dev/tests/contract-regressions.sh
---

# Judge effort_estimate on Every Capture, As impact Already Is

## What

`skills/do-work/actions/capture.md` requires a judged `impact:` on every REQ it mints, in a rule
written to close exactly this hole: "**Judge it yourself and write a value** — an absent `impact:`
must not be the common case." The neighbouring field gets no such treatment.
`actions/capture-reference.md` says only that capture **MAY** set `effort_estimate`, and the
schema line in `actions/work-reference.md` repeats it.

Done means capture judges both fields by the same standard, with the same escape hatches: judge it,
or put the judgment to the user, or leave it absent because neither was possible — never a copied
default.

## Why

The asymmetry stopped being cosmetic when 0.220.0 shipped `do-work run-simple-reqs`, which selects
work for a cheaper model on `effort_estimate` (`skills/do-work/tools/select-simple-reqs.sh`).
Measured on the live queue at capture time: **14 of 22 pending REQs carry the field**, so 8 resolve
to `effort-substantive` by documented default and are invisible to that verb. They are not
invisible because someone judged them substantive. They are invisible because nobody judged them.

That is the same failure `impact:`'s rule already names, and the fix is already written next door —
this REQ applies it one field over rather than inventing anything.

Fold-first scan (`actions/capture-reference.md` → Fold-First Rule): no pending REQ in any UR,
sweep or not, shares this root cause. The nearest candidate, REQ-293 (`sweep_key:
impact-effort-lockin-checks-underpin`), is about lock-in checks pinning a spelling instead of a
property in `contract-regressions.sh` and `generate_test.go` — a different root cause with a
different fix site, and folding into it would bury this behind test-assertion work.

## Requirements

1. **Capture judges `effort_estimate` on every REQ it mints**, using the same three-way contract
   `impact:` already carries: judge and write a value, or ask the user when the answer is genuinely
   unclear, or leave the field absent because neither was possible. Never write a value that was
   not judged — asserting `effort-substantive` because it is the default is exactly as wrong as
   inventing `effort-mechanical`, which is the symmetric half of the rule `impact:` learned in
   REQ-294.
2. **State the judgment as a question a capture can actually answer**, the way `impact:`'s two
   questions are stated. Size is the axis: would a competent implementer finish this in one focused
   pass over a small, already-identified set of files?
3. **The two axes stay independent.** `effort_estimate` must not be derived from the `impact:`
   verdict in either direction; the schema already says so in both field descriptions and that
   sentence is the thing being enforced, not restated.
4. **Update every site that currently states the weaker rule** so no file still says capture MAY
   set it: `actions/capture.md`, `actions/capture-reference.md`, and the schema line in
   `actions/work-reference.md`.
5. **A lock-in check pins the property, not the wording** — capture's checklist must fail a capture
   that emits an unjudged `effort_estimate`. Follow REQ-293's lesson rather than asserting one
   phrase is present.

## Requirements Deliberately Excluded

- **No enum growth.** `effort-mechanical` | `effort-substantive` stays closed. REQ-228 ruled against
  t-shirt sizes and this REQ does not reopen it.
- **No backfill of existing REQs.** Absence keeps reading as `effort-substantive` for every record
  already written, and frozen fields on claimed or archived REQs are never rewritten. Raising the
  bar for new captures is the whole scope; a migration is a separate decision with a separate risk.
- **No change to `run-simple-reqs` or its selector.** They read the field correctly already; this
  REQ improves what the field contains. Touching the reader would be scope creep.

## Context

Discovered while building `do-work run-simple-reqs` (0.220.0, PR #152) and deliberately excluded
from that PR, which was scoped to the selector and would have grown a capture-semantics change on
top of it. The user was asked and chose to file it rather than fold it in.

Related reading before starting: `actions/capture.md`'s impact-assessment bullet is the model to
copy, and `actions/work-reference.md`'s `effort_estimate` schema line records every current reader
of the field.

## AI Execution State (P-A-U Loop)

- [ ] [PLAN] Technical approach recorded
- [ ] [APPLY] Changes made within declared scope
- [ ] [UNIFY] Diff reviewed, linters run, no debug artifacts
