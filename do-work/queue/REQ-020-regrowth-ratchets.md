---
id: REQ-020
title: "Regrowth ratchets: router word budget + new-action sibling-skill justification rule"
status: pending
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-015]
related: [REQ-015]
batch: harness-bloat-cleanup
maintenance: false
---

# Phase-4 ratchets

## What
1. Add a router word-count budget assertion to `_dev/tests/contract-regressions.sh`
   (or a sibling `_dev/tests/` script it sources): fail when `wc -w SKILL.md`
   exceeds the post-REQ-015 count + 10% headroom (hard number recorded in the
   test with a comment explaining the ratchet).
2. Add to CLAUDE.md (Action File Conventions): any NEW action file must state, in
   its description blockquote or an accompanying ADR, why it belongs in do-work
   rather than a sibling skill; reviewers reject additions without it.

## Why
Phase 4 of UR-003: the audit found bloat regrew through untested accretion
(five router enumerations, 162 changelog entries). Guards make regrowth visible
at commit time.

## Acceptance criteria
- [ ] Budget test fails on a synthetic 300-word SKILL.md addition (demonstrated,
      then reverted) and passes on the real file.
- [ ] CLAUDE.md rule added under Action File Conventions.

## Open Questions
(none)

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:**
- [ ] **[APPLY]:**
- [ ] **[UNIFY]:**
