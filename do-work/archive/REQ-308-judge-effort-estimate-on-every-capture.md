---
id: REQ-308
title: "[impact-rule-change] Judge effort_estimate on every capture, as impact already is"
status: completed
created_at: 2026-08-20T22:00:52Z
status_changed_at: 2026-08-21T08:45:00Z
claimed_at: 2026-08-21T08:45:00Z
completed_at: 2026-08-21T08:57:04Z
kb_status: pending
commit: 9bce005
user_request: UR-064
domain: general
impact: impact-rule-change
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: true
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-21T08:51:36Z
  basis:
    - trivial short-circuit
write_set:
- skills/do-work/actions/capture.md
- skills/do-work/actions/capture-reference.md
- skills/do-work/actions/work-reference.md
- skills/do-work-board/tools/queue-kanban/model.go
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

- [x] **[PLAN]:** Mirror capture's impact bullet one field over; make the two checklist lines one sentence apart from the field they name, so a lock-in check can pin the symmetry rather than a phrase; sweep every restatement of the weaker rule.
- [x] **[APPLY]:** Five files. Two beyond the REQ's declared write set, extended per D-01 and D-02 below and mirrored into `write_set` above.
- [x] **[UNIFY]:** `git diff --stat` reviewed (5 files); `bash -n` on the changed shell; `gofmt -l .` and `go vet ./...` clean for the Go comment change; no debug artifacts; every changed file re-read at its edit site.

## Scope

**Files I will touch:**
- `skills/do-work/actions/capture.md` (modify) — the effort-assessment bullet and the checklist line
- `skills/do-work/actions/capture-reference.md` (modify) — the frontmatter template comment
- `skills/do-work/actions/work-reference.md` (modify) — the schema line and the Schema Read Contract row
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — the parser comment restating the weaker rule
- `_dev/tests/contract-regressions.sh` (modify) — the symmetry check and the no-restatement sweep

**Files I will NOT touch:** `skills/do-work/tools/select-simple-reqs.sh` and `actions/run-simple-reqs.md` — the REQ excludes the reader; `actions/work-reference.md`'s Discovered Tasks Classification, which governs review-minted follow-ups rather than capture.

**Acceptance criteria (restated from REQ):**
- [x] Capture judges `effort_estimate` on every REQ it mints, with the same three-way contract as `impact:`
- [x] The judgment is stated as a question a capture can answer
- [x] The two axes stay independent — `effort_estimate` is never derived from `impact:`
- [x] Every site stating the weaker rule is updated
- [x] A lock-in check pins the property, not the wording
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture.md` (modified) — new **Effort assessment** bullet beside the impact one; a checklist line that is impact's sentence with the field name swapped
- `skills/do-work/actions/capture-reference.md` (modified) — the template comment now states the three-way standard instead of MAY
- `skills/do-work/actions/work-reference.md` (modified) — the `effort_estimate` schema line and the Schema Read Contract row
- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — the parser comment's restatement of the weaker rule
- `_dev/tests/contract-regressions.sh` (modified) — the symmetry check plus the shipped-restatement sweep

**What was done:** Capture's impact rule was applied one field over, unchanged in shape. The
size question is stated plainly — would a competent implementer finish this in one focused pass over
a small, already-identified set of files — and both directions of the never-write-an-unjudged-value
rule are named, so `effort-substantive` by default is called out as the same failure as an invented
`effort-mechanical`.

## Decisions

- **D-01 — Write set extended to `actions/work-reference.md`.** DECIDE & STATE. The REQ's
  Requirement 4 names the schema line there explicitly while its `write_set` omits it, so the field
  and the requirement disagreed. The requirement wins: a rule left stated weakly in the schema is
  the restatement that defeats the fix. Recorded here and mirrored into `write_set`, per the
  stop-and-report contract in `actions/work.md` Step 6.
- **D-02 — Write set extended to `queue-kanban/model.go`.** DECIDE & STATE. The new sweep found a
  fifth site the REQ did not list: the parser comment said "capture MAY set it". That comment's own
  closing sentence requires it to stay in lock-step with the Schema Read Contract and to change in
  the same commit, so leaving it would have broken a rule the file states about itself. Comment
  only — no parser behavior moved.
- **D-03 — The lock-in check pins that the two checklist lines are one sentence apart, plus the
  three alternatives.** DECIDE & STATE. Symmetry alone would accept both fields being weakened
  together, and a phrase assertion is what REQ-293's lesson rules out. The pair holds the property
  from both sides: mutation M3 weakened both lines identically and still failed.

## Testing

**Tests run:** `QUEUE_KANBAN_BROWSER=<chromium> bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 — `Maintainer verification passed.`

**Red-green validation:**
- **RED, on the untouched tree, both arms.**
  `actions/capture.md's checklist must carry exactly one line stating how` effort_estimate `is
  decided; found 0` and `a shipped file still says capture MAY set a field it must now judge`.
- **GREEN.** Both pass, and the whole suite exits 0.

**Mutation-tested — five reverts, each caught:**
- M1 delete the effort checklist line → "found 0", caught
- M2 weaken only the effort line → "judges impact: and effort_estimate: by different standards", caught
- M3 **weaken both lines identically** → symmetry still holds, so the three-alternatives guard catches
  it. This is the mutation a symmetry-only check would have passed, and the reason the check has two
  halves.
- M4 restore "Capture MAY set it" in `capture-reference.md` → sweep, caught
- M5 restore "capture MAY set it" in `model.go` → sweep, caught

**New tests added:**
- `_dev/tests/contract-regressions.sh` — the symmetry check (capture's `impact:` and
  `effort_estimate:` checklist lines must be one sentence apart from the field they name), the
  three-alternatives guard on that sentence, and a repository sweep for any shipped file still
  saying capture MAY set or emit a field.

**Existing tests updated (cross-REQ impact):** none.

**Measured against the live queue.** The sweep the check runs found a fifth site the REQ did not
list — `queue-kanban/model.go` — which is what the sweep half exists for. Reading the four listed
sites would have shipped the fix with its own contradiction still in the tree.

*Verified by work action*

## Review — 2026-08-21T08:56:16Z

**Overall: 95%**

| Dimension | Score |
|---|---|
| Requirements Compliance | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope Discipline | 95% |
| Risk | None |

**Acceptance: Pass.**

**Requirements Compliance.** All six hold. The size question is stated as one answerable question
(Requirement 2); independence is stated in the bullet and in the schema line, in both directions
(Requirement 3); the excluded items stayed excluded — the enum is still two values, no existing REQ
was backfilled, and `select-simple-reqs.sh` and `run-simple-reqs.md` are untouched.

**Findings**

- **F1 — Minor, accepted.** The sweep greps for `apture MAY (set|emit)`, which catches the two
  spellings that existed and would miss a third phrasing of the same permission ("capture is free
  to set it"). Narrowing further would false-positive on ordinary prose; the symmetry half is the
  arm that holds the rule itself, and this one is the restatement net.
- **F2 — Minor, noted not fixed.** `actions/work-reference.md` → **Discovered Tasks Classification
  (Step 8)** still tells the review-minted follow-up path to "write `effort-mechanical` only when
  you have actually judged the fix small, and otherwise leave it absent". That is the weaker rule
  one writer over. The REQ scoped itself to capture deliberately and lists no follow-up path, so
  changing it here would widen the REQ; it is recorded in Discovered Tasks instead.
- **F3 — Minor, fixed in-REQ.** This REQ's own P-A-U boxes were captured without the bold
  `**[UNIFY]**` marker, which left qualify's box audit disarmed rather than passed — exactly the
  state REQ-264 made visible. Normalized to the template shape and re-qualified armed.

**Scope Discipline (95%).** Two files beyond the declared write set, both recorded as decisions
before the edit and both required by the REQ's own text or by the edited file's own lock-step rule.
The deduction is for the REQ shipping with a `write_set` that contradicted its Requirement 4, not
for the response to it.

## Discovered Tasks

- **impact-rule-change** The review-minted follow-up path still states the weaker rule for
  `effort_estimate`. `actions/work-reference.md` → **Discovered Tasks Classification (Step 8)**
  and `actions/review-work.md` Step 10 both tell that writer to "write `effort-mechanical` only
  when you have actually judged the fix small, and otherwise leave it absent" — permission not to
  judge, which is the rule capture just lost. Follow-ups minted by review are a large share of the
  queue (this run alone created three), so they are a large share of what
  `do-work run-simple-reqs` cannot see. The fix is the same one applied here, one writer over;
  the reason it was not done here is that this REQ named capture and only capture.

## Lessons Learned

- **Two rules that should be the same rule can be pinned by comparing them, not by quoting them.**
  Asserting a phrase is present is what REQ-293 ruled against; asserting that two sentences are
  identical apart from the field they name enforces "by the same standard" literally, and survives
  any rewording of either.
- **Symmetry alone is satisfiable by weakening both sides.** M3 did exactly that and stayed
  symmetric. A symmetry check needs a floor beside it — here, that the shared sentence still offers
  all three alternatives.
- **The sweep half of a check is what finds the site the REQ did not know about.** The REQ listed
  four files; the grep found a fifth, in a Go comment whose own closing sentence required it to
  change in the same commit. Enumerating sites in the REQ is a starting point, never the set.

## Orientation

**What changed in the map.** `effort_estimate` stopped being a field capture could skip. It is now
judged on every minted REQ by the same three-way contract as `impact:` — judge it, ask, or leave it
absent — and no shipped file says otherwise any more.

**What this makes true.** A REQ that reads as `effort-substantive` now means somebody looked at its
size and said so. `do-work run-simple-reqs` selects on that field, so the queue's small work stops
being invisible to the cheaper-model verb for want of a judgment nobody made. Nothing was backfilled
and the enum did not grow.

**Subsystem:** the do-work capture path and the shared request schema —
`actions/capture.md`, `actions/capture-reference.md`, `actions/work-reference.md`, and the board
parser that mirrors the schema. Prime: `_dev/primes/prime-action-files.md`.

**Discovered-task follow-up:** the review-minted follow-up path was queued as REQ-314 (`pending-answers`, `impact-rule-change`).
