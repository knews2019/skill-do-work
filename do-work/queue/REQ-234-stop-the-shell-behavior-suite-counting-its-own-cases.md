---
id: REQ-234
title: Stop the shell behavior suite counting its own cases
status: pending-answers
domain: general
created_at: 2026-08-18T01:44:18Z
user_request: UR-042
addendum_to: REQ-229
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
write_set:
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Discovered Task: Stop the Shell Behavior Suite Counting Its Own Cases

## What

`_dev/tests/prescribed-shell-scripts-behavior.sh` ends by printing how many cases it ran, from a hand-maintained literal. The number matches no count derivable from the file, so it reports a remembered figure as a measured one. Derive it or drop it.

## Context

Found while adding two cases in REQ-229. The closing line read `(45 named script cases)` while the file carried 40 case-header comments; REQ-229 bumped it to 47 because adding two cases in the existing style makes that correct under whatever convention produced 45, but it did not repair the underlying problem.

`_dev/primes/prime-shell-commands.md` § *Closed Enumerations Go Stale* is this exact pattern, and a summary line reporting a suite's own size is a bad place for it: a reader has every reason to take that number as measured.

## Requirements

- The closing line either reports a count computed from the file at run time, or reports no count at all. It must not carry a literal.
- If a count is computed, the thing being counted has a stated definition in the file — the shape that makes something "one case" — so the number and the file cannot disagree.
- No test case is added, removed, or weakened by this change.
- `bash _dev/tests/maintainer-verify.sh` still exits 0.

## Red-Green Proof

**RED prompt/case:** the literal at the end of `_dev/tests/prescribed-shell-scripts-behavior.sh`.
**Why RED now:** it is a hand-maintained number that already disagrees with every count derivable from the file.
**GREEN when:** that literal is absent — either replaced by a computed count with its counting rule stated, or removed along with the claim.
**Validation:** Discovered during REQ-229; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**, deletion branch.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-229: one of the maintainer test scripts finishes by announcing how many test cases it just ran, and that number is typed into the file by hand rather than counted. It is already out of step — it claimed 45 while the file held 40 of the obvious unit — so it reports something remembered as though it were measured. Nothing is broken: every test still runs and still passes, and the number appears only in a success message. The fix is either to count the cases at run time or to stop claiming a number. It is your call rather than mine because working out what the original number was counting means picking a definition of "one case" for this suite, and picking it wrong would silently change what that line reports to every future reader. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
