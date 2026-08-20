---
id: REQ-301
title: "Let qualify tell a moved line from an added one"
status: pending-answers
created_at: 2026-08-20T09:06:00Z
status_changed_at: 2026-08-20T09:06:00Z
user_request: UR-056
addendum_to: REQ-258
domain: general
impact: impact-user-visible
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work/tools/checks/qualify.sh
- _dev/tests/prescribed-shell-cases/qualify.sh
---

# Let Qualify Tell a Moved Line From an Added One

## What

`skills/do-work/tools/checks/qualify.sh` runs `git diff` with no rename or copy detection (`-M`/`-C`) and greps `^+` for debug artifacts. Relocated text therefore reads as newly added, so **every REQ that moves code fails the `[UNIFY]` debug-artifact audit on markers that already existed**.

REQ-258 hit this: it FAILed on four `TODO` strings that are deliberate fixture data in the REQ-254 `qualify` cases and are byte-identical in `git show HEAD:_dev/tests/prescribed-shell-scripts-behavior.sh`. The REQ recorded the override with that evidence and proceeded.

The failure mode is the dangerous direction. A FAIL that is *usually* a false positive on this class of REQ teaches builders to un-check `[UNIFY]` or wave the FAIL away — which is precisely the reflex the audit exists to prevent. A gate that cries wolf on a whole category of change is worse than one that is quiet, because it erodes the response to the true positives.

**Candidate fixes** (pick one, do not build both):
- Enable copy/rename detection for the debug-artifact scan (`git diff -C --find-copies-harder`) so moved hunks are not read as additions.
- Or subtract from the flagged set any line that appears verbatim in the pre-change tree.

The first is git doing the work; the second is cheaper to reason about but can mask a genuinely re-added marker elsewhere in the same file. Whichever ships needs a lock-in case in `_dev/tests/prescribed-shell-cases/qualify.sh` proving that a *moved* marker passes while a *fresh* one still FAILs — the existing REQ-254 case already pins the second half.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-258: `tools/checks/qualify.sh` has no rename/copy detection, so every code-relocation REQ gets a false `[UNIFY]` FAIL on pre-existing debug markers inside the moved text. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

## Context

Discovered during REQ-258 (splitting the prescribed-shell behavior suite into per-script files). See that REQ's `## Qualification` section for the worked example and the evidence used to override the FAIL.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
