---
id: REQ-020
title: "Regrowth ratchets: router word budget + new-action sibling-skill justification rule"
status: completed
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
claimed_at: 2026-07-15T18:50:00Z
completed_at: 2026-07-15T18:56:00Z
route: A
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-021]
related: [REQ-021]
batch: harness-bloat-cleanup
maintenance: false
commit: 5570ee4
---

# Phase-4 ratchets

## What
1. Add a router word-count budget assertion to `_dev/tests/contract-regressions.sh`
   (or a sibling `_dev/tests/` script it sources): fail when `wc -w SKILL.md`
   exceeds the post-REQ-021 count + 10% headroom (hard number recorded in the
   test with a comment explaining the ratchet).
2. Add to CLAUDE.md (Action File Conventions): any NEW action file must state, in
   its description blockquote or an accompanying ADR, why it belongs in do-work
   rather than a sibling skill; reviewers reject additions without it.

## Why
Phase 4 of UR-003: the audit found bloat regrew through untested accretion
(five router enumerations, 162 changelog entries). Guards make regrowth visible
at commit time.

## Acceptance criteria
- [x] Budget test fails on a synthetic 300-word SKILL.md addition (demonstrated,
      then reverted) and passes on the real file.
- [x] CLAUDE.md rule added under Action File Conventions.

## Open Questions
(none)

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Budget = measured post-diet 2,396 + ~10% = 2,650; assertion lives in the existing contract suite (single harness, runs everywhere the other contracts run).
- [x] **[APPLY]:** Word-budget block added to _dev/tests/contract-regressions.sh with an in-test comment explaining the ratchet and the merge-or-lazy-load escape hatch; sibling-skill justification rule added atop CLAUDE.md Action File Conventions.
- [x] **[UNIFY]:** Red-green demonstrated: +300 synthetic words → suite FAILs at 2,698 > 2,650 (exit 1); reverted → green. Real file passes at 2,396.

## Triage

Route A — two additive guards at named locations.

## Implementation Summary

**What was done:** Regression guards against router regrowth and unjustified action accretion.

Files changed:
- `_dev/tests/contract-regressions.sh` (modified) — router_word_budget=2650 assertion.
- `CLAUDE.md` (modified) — sibling-skill justification gate for new actions, with the routing-surface cost spelled out.
- `actions/version.md`, `CHANGELOG.md` — version 0.124.2 + entry.

## Testing

- Red-green: synthetic 300-word SKILL.md addition → FAIL (exit 1, message names the count and the budget); revert → suite green. Recorded in [UNIFY].

## Lessons Learned

**Worth knowing:** The budget comment in the test is deliberately load-bearing — it tells the person who hits the limit what the sanctioned fixes are (merge, lazy-load) so the ratchet teaches instead of just blocking.
