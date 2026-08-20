---
id: REQ-306
title: "[impact-rule-change] Reserve new-REQ creation for behavior changes, not prose discrepancies"
status: completed
completed_at: 2026-08-20T14:51:04Z
commit: f955d28
created_at: 2026-08-20T13:21:13Z
user_request: UR-063
domain: general
impact: impact-rule-change
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-307]
maintenance: true
write_set:
- skills/do-work/actions/capture.md
- skills/do-work/actions/capture-reference.md
---

# Reserve New-REQ Creation for Behavior Changes, Not Prose Discrepancies

## What

Capture will mint a REQ for a wrong cross-reference number as readily as for a broken checker. That
is the arrival rate, not the backlog: **145 of 298 REQs ever written were spawned by a prior REQ,
and four of the five deepest descent chains terminate in prose reconciliation** — a stale count, a
wrong number, a comment describing a superseded mechanism.

A prose-only discrepancy stops being its own REQ. It goes to the standing sweep (REQ-307), or it
gets folded into the next commit that touches that file for another reason. **New-REQ creation is
reserved for a change in behavior, in a checker's predicate, or in a rule's stated condition.**

**The judgment already exists one step downstream and must not be re-derived.**
`skills/do-work/actions/review-work.md`'s Sweep-consolidation block already routes an
`impact-negligible` finding — or any set sharing one root cause — into a sweep rather than its own
REQ, and already says "Done means the class cannot recur." This REQ applies that same judgment
**earlier, at capture**, to the class actually producing the volume. Cite that block as the canonical
statement; state here only what is new, which is the *timing* and the prose-only test.

**Done means the class cannot recur:** the boundary is stated once in capture's canonical home, keyed
on the condition rather than on a list of finding shapes, and capture has somewhere to send what it
no longer mints.

## Requirements

- **The boundary is stated as a condition, not a list.** "Prose-only" means the fix changes no
  behavior, no checker predicate, and no rule's stated condition. Any list of example shapes (stale
  count, wrong cross-reference, superseded-mechanism comment) is marked illustrative —
  `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale.
- **The rule names its destination.** A capture that declines to mint a REQ must say where the
  finding went: appended to REQ-307's `## Instances`, or folded into a named pending REQ whose
  declared scope already contains the file. Declining without a destination is how findings get
  lost, and is worse than the over-capture this fixes.
- **Three exemptions must survive explicitly**, or this rule will suppress work it must not:
  - An `impact-critical` finding is never deferred, at any depth, whatever its fix touches.
  - A **contradiction** between two shipped instructions is not prose reconciliation. Two rules that
    cannot both be followed change what an agent does, so they change behavior. REQ-288 is the live
    example and must remain a first-class REQ under this rule.
  - A prose discrepancy in a **user-facing** artifact contract (a changelog entry, a report template
    a consumer reads) is judged on its reader, not on its diff shape.
- **The rule applies to what capture *creates*, never to what already exists.** It must not be read
  as licence to cancel queued REQs retroactively; a queued REQ is user intent and only
  `do-work abandon` removes it.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Requirements Deliberately Excluded

- **No mechanical check.** A predicate that separates "prose-only" from "behavioral" would have to
  read intent from a diff, and the cost of a false positive is a suppressed real finding. State the
  condition where capture is governed and loaded; do not build a gate for it. If a future REQ wants
  one, the argument has to address that asymmetry first.
- **No retroactive sweep of the existing queue.** Out of scope by explicit instruction — the batch
  that created this REQ was capped at two REQs precisely so a rule about over-capture would not
  spawn more work than it prevents.

## Red-Green Proof

**RED prompt/case:** No mechanical check exists or is wanted (see Requirements Deliberately
Excluded), so the closure evidence is the *deletion surface* rather than a failing test: today
neither `skills/do-work/actions/capture.md` nor `capture-reference.md` states any condition under
which a finding is recorded without minting a REQ. `grep -n "prose-only\|standing sweep" ` over both
files returns nothing.
**Why RED now:** Capture's assessment steps judge domain, effort, impact, and TDD, and every path
ends in a REQ file. There is no branch that records a finding elsewhere, which is why the class has
no exit.
**GREEN when:** Both files state the boundary and name REQ-307 as the destination, the three
exemptions are present verbatim, and `maintainer-verify.sh` exits 0. Closure is judged against that
named surface per `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
**Validation:** Directed capture under UR-063; premise measured by the user across 298 REQs.

## Context

Captured by explicit instruction under UR-063 — one of exactly two REQs that batch was permitted to
create. `depends_on: [REQ-307]` is load-bearing rather than cosmetic: this REQ writes a rule whose
compliance path is "send it to the standing sweep", so building it before that sweep exists would
ship a forward reference.

Worth knowing for whoever builds this: the same session that captured it processed four REQs and
found a stated premise already stale in three of them (REQ-300's headline instance, REQ-271's
second defect, REQ-298's severity). That is the same phenomenon from the other side — prose
discrepancies accumulate faster than REQs about them can drain, which is the case for routing them
to a batch rather than a queue slot each.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
