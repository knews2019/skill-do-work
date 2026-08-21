---
id: REQ-314
title: "[impact-rule-change] Judge effort_estimate on review-minted follow-ups too"
status: pending-answers
created_at: 2026-08-21T08:56:44Z
status_changed_at: 2026-08-21T08:56:44Z
user_request: UR-064
addendum_to: REQ-308
domain: general
review_generated: true
impact: impact-rule-change
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: true
---

# Judge effort_estimate on Review-Minted Follow-Ups Too

## What

REQ-308 made capture judge `effort_estimate` on every REQ it mints, by the same three-way contract
`impact:` already carried. The other writer of new REQs kept the weaker rule.
`actions/work-reference.md` → **Discovered Tasks Classification (Step 8)** and
`actions/review-work.md` Step 10 both tell that path to "write `effort-mechanical` only when you
have actually judged the fix small, and otherwise leave it absent to read as `effort-substantive`".

Half of that is right and stays: never invent `effort-mechanical`. The other half is permission not
to judge, and it is the rule capture just lost.

## Why

Review-minted follow-ups are a large share of the queue — this `do-work run` alone created three
(REQ-312, REQ-313, and this one) — so they are a large share of what `do-work run-simple-reqs`
cannot see. The measurement REQ-308 was built on was 14 of 22 pending REQs carrying the field; the
missing eight are exactly this population.

Leaving the two writers on different standards also reintroduces the asymmetry REQ-308 removed, one
level down: a reader who lands on Step 8's rule learns that leaving the field absent is fine.

## Requirements

- The review/discovered-task follow-up path judges `effort_estimate` by the same three-way contract
  capture now uses: judge it, or put the judgment to the user, or leave it absent because neither was
  possible — never a copied default, in either direction.
- Every site stating the weaker rule is updated, not just the one this REQ names. REQ-308's sweep
  found a fifth site its own REQ had not listed; run the same sweep.
- The two axes stay independent: `effort_estimate` is never derived from the finding's `impact:`
  token. That sentence already exists at both sites and is being enforced, not restated.
- The lock-in check pins the property, the way REQ-308's does — ideally by extending that check to
  cover this writer, rather than adding a second one that can drift from it.
- No backfill of existing REQs, and no enum growth. Both exclusions carry over from REQ-308 unchanged.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED prompt/case:** Read `actions/work-reference.md` → **Discovered Tasks Classification (Step 8)**
and `actions/review-work.md` Step 10. Both still say to leave `effort_estimate` absent when it is
not obviously mechanical, which is permission not to judge.
**GREEN when:** A check fails on a follow-up-minting rule that permits an unjudged
`effort_estimate`, and passes once both sites carry capture's contract.
**Validation:** Discovered task from REQ-308; apply `actions/work-reference.md` →
**Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-308: that REQ made request capture
  judge the size field on every new request instead of leaving it blank, because blank quietly reads
  as "big" and hides the request from the cheaper-model queue verb. But requests are also created a
  second way — automatically, when a code review finds something worth following up — and that path
  still says leaving the field blank is fine. Those follow-ups are a large share of the queue, so
  the problem REQ-308 fixed is mostly still there. Fixing it means applying the same rule one writer
  over. Should I process this as a new task?

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  Why this is yours rather than mine: REQ-308 deliberately scoped itself to capture, and you were
  the one who filed it that way rather than folding it into the pull request that discovered it. So
  the question of whether the same rule should reach the automatic path is a scope call you have
  already made once, in the other direction, and it should be yours to make again rather than mine
  to assume.
