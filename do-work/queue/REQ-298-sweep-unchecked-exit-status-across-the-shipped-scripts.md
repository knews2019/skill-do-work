---
id: REQ-298
title: "[impact-rule-change] Review fix: sweep the unchecked-exit-status primitive across every shipped script"
status: pending
created_at: 2026-08-19T19:43:58Z
status_changed_at: 2026-08-20T08:22:51Z
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

## Reproduced During Clarify (2026-08-20)

The REQ was written from the review's static reading. Driving the real script proved the
severity is higher than "the template still carries the pattern", and narrowed which of the
four sites actually matters.

**One site fails OPEN, and it is the incident guard itself** —
`skills/do-work/tools/checks/record-commit-hash.sh:466`:

```bash
head_blob_bytes="$(git cat-file -s "HEAD:$tracked_full_name" 2>/dev/null || true)"
if [ -n "$head_blob_bytes" ] && [ "$head_blob_bytes" -gt 0 ] && ...
```

`|| true` collapses *there is no blob in HEAD* and *git could not answer* into one empty
string, and the `-n` test then skips the truncation floor for both. Against a 10,178-byte
archived REQ truncated to 85 bytes — the shape of the original incident — with everything
real except a failing `git cat-file -s`:

```
control (real git):        FAIL: ... is 85 bytes on disk but 10178 bytes in HEAD —
                           content was lost BEFORE this script ran.        exit 1
cat-file -s failing:       OK: ... all content guards passed.
                           Now stage and commit exactly this file: ...     exit 0
```

It wrote the hash into the remnant and printed the instruction to commit it. That is the
exact failure this script exists to prevent.

**Three sites fail SAFE** (`:189`, `:194`, `:605`): a dead command yields `0`, the
comparison mismatches, and the run stops with a false FAIL. Real instances of the condition,
but they cost a confusing message rather than a lost file. They need the same treatment for
the class's sake; they are not urgent.

## Instances

- [ ] `skills/do-work/tools/checks/record-commit-hash.sh:466` — the truncation floor, the
  one reproduced fail-open. Fix it the way `repair-req-timestamps.sh` was fixed under
  REQ-268: ask `git cat-file -e` whether the blob exists, treat an existing blob whose size
  will not read as a failure, and keep the genuine no-blob case as the skip it is. One probe
  in `_dev/tests/record-commit-hash-guards.sh`, which already has the fixture repo for it.
- [ ] `record-commit-hash.sh:189`, `:194`, `:605` — the three fail-safe sites. Same
  treatment for the class, or an in-place comment recording that each was judged and why
  the content alone is sufficient there. Either is fine; silence is not.
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

  > **Recorded dissent, overruled by the user at clarify (2026-08-20).** I argued this
  > requirement should be dropped: discarding an exit status is *correct* at most of the
  > sites it appears — `grep -c` on a no-match is the common one, and this very script has
  > two of those with comments explaining why — so a `|| true`-shaped check flags mostly
  > legitimate code and gets muted within a week. My recommendation was to state the
  > condition in `_dev/primes/prime-shell-commands.md` instead, where it is loaded before
  > shell is written. The user chose the full scope. **Build it as written**, and take the
  > escape hatch the requirement already contains honestly: if the check cannot be made to
  > separate the correct uses from the defects, say so with the evidence and fall back to
  > the stated convention — that is a finding, not a failure to deliver.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** A probe in `_dev/tests/record-commit-hash-guards.sh` — **not** the
prescribed-shell behavior suite — that drives `record-commit-hash.sh` against a
pre-truncated archived REQ with `git cat-file -s` failing on PATH, and asserts it does not
report `all content guards passed`. The guards file is this script's dedicated behavioral
suite and already builds the throwaway git repo the fixture needs; REQ-276 hit exactly this
misdirection in its own `write_set` and recorded it as D-01.
**Why RED now:** Reproduced — see **Reproduced During Clarify** above. The unfixed script
prints `OK: ... all content guards passed`, writes the hash into an 85-byte remnant of a
10,178-byte file, and exits 0.
**GREEN when:** That case passes and the full suite still exits 0.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure
Ratchet (Step 6.5)**.

## Open Questions

- [x] REQ-268's review found that the unchecked-exit-status pattern it fixed in two
  timestamp scripts was copied from `record-commit-hash.sh` — the script whose guard style
  those two are required to imitate — so the template still carries the defect and any new
  guarded-edit script written from it inherits it. Fixing that means touching the script
  that guards the last write every REQ passes through, plus a sweep of the other shipped
  scripts and a new rule in the shell prime. Should I process this as a new task?
  → **Yes, add to queue — the FULL scope as written**, including the repo-wide sweep and the
  mechanical check. Chosen over the narrower "just the reproduced bug" option I recommended.

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  Why this is yours rather than mine: the fix is safe in itself, but the file it centres on
  is the one that once truncated six archived REQs to zero bytes when edited freely, and
  the sweep's real cost is the third requirement — deciding whether a mechanical check for
  this shape can be written without flagging the many places where discarding a status is
  correct. That is a judgment about how much guard machinery this repo wants, not a
  technical unknown.

**Answered [2026-08-20]:** User approved via `do-work clarify`, at full scope, after asking
how the fix would work and being shown the reproduction below and the four sites.
