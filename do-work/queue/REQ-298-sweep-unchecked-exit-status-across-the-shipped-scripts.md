---
id: REQ-298
title: "[impact-rule-change] Review fix: sweep the unchecked-exit-status primitive across every shipped script"
status: pending-answers
created_at: 2026-08-19T19:43:58Z
user_request: UR-056
addendum_to: REQ-268
domain: general
review_generated: true
sweep: true
sweep_key: unchecked-exit-status-read-as-content
impact: impact-rule-change
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
depends_on: []
maintenance: false
write_set: []
---

# Review Fix: Sweep the Unchecked-Exit-Status Primitive Across Every Shipped Script

## What

REQ-268 closed five instances of one condition — a command or process substitution whose
exit status is discarded while only its content is judged, so a tool that never ran reads
as a tool that found nothing — and stated that condition in the headers of the two scripts
it touched. Its independent review then found the same primitive **inherited verbatim from
a third script that neither REQ-268 nor its Requirements named**:
`skills/do-work/tools/checks/record-commit-hash.sh` is the file REQ-268's own header cites
as its mandatory guard-style template, and it carries the pattern itself (for example
`head_blob_bytes="$(git cat-file -s … 2>/dev/null || true)"`). The copy direction was
template → copies, so fixing the copies and leaving the template means the next script
written from it re-imports the defect.

**Done means the class cannot recur:** the condition is stated once where shipped shell is
governed, and every shipped script that takes a substitution and judges only its content
either checks the status or says in place why the content alone is sufficient. Patching N
sites one at a time is what this REQ exists to avoid.

## Instances

- [ ] `skills/do-work/tools/checks/record-commit-hash.sh` — the template REQ-268's header
  cites; the review named its `git cat-file -s … || true` as the origin of the copied
  pattern. This is the one that keeps the class alive, because new guarded-edit scripts
  are written from it.
- [ ] Every other shipped script under `skills/*/scripts/` and `skills/*/tools/` — the
  instance list is a sample, not the scope. `_dev/primes/prime-shell-commands.md`
  § Closed Enumerations Go Stale applies: key the rule on the condition and mark this
  list illustrative.
- [ ] `_dev/primes/prime-shell-commands.md` itself — the condition belongs in the trap
  list so it is loaded before shell is written, rather than restated per script header.

## Context

Found during the independent review of REQ-268 (finding I1, Important, gate:
`impact-rule-change` — it changes a rule applying across every shipped script rather than
fixing one visible defect). REQ-268 fixed all five sites inside its own two-file write set,
including the three the review surfaced, so nothing named there is left open; what remains
is the spread beyond that write set, which is genuinely a different REQ.

Created `pending-answers` per the generation-≥2 cascade depth stop: REQ-268 itself carries
`review_generated: true`.

## Requirements

- The condition is stated once, where shipped shell is governed, rather than re-derived per
  script header.
- `record-commit-hash.sh`'s own substitutions either check their status or state why the
  content is sufficient on its own.
- A mechanical check that finds the shape, so a new script cannot reintroduce it silently —
  or, if no check can be written without false positives, a stated review convention with
  the reason recorded.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** A lock-in in `_dev/tests/prescribed-shell-scripts-behavior.sh` that
drives `record-commit-hash.sh` with the relevant git subcommand failing on PATH and asserts
it does not report success for a guard that never ran.
**Why RED now:** The pattern is present at the site the review named; the same fixture
shape reproduced all five of REQ-268's instances before they were fixed.
**GREEN when:** That case passes and the full suite still exits 0.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure
Ratchet (Step 6.5)**.

## Open Questions

- [ ] REQ-268's review found that the unchecked-exit-status pattern it fixed in two
  timestamp scripts was copied from `record-commit-hash.sh` — the script whose guard style
  those two are required to imitate — so the template still carries the defect and any new
  guarded-edit script written from it inherits it. Fixing that means touching the script
  that guards the last write every REQ passes through, plus a sweep of the other shipped
  scripts and a new rule in the shell prime. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  Why this is yours rather than mine: the fix is safe in itself, but the file it centres on
  is the one that once truncated six archived REQs to zero bytes when edited freely, and
  the sweep's real cost is the third requirement — deciding whether a mechanical check for
  this shape can be written without flagging the many places where discarding a status is
  correct. That is a judgment about how much guard machinery this repo wants, not a
  technical unknown.
