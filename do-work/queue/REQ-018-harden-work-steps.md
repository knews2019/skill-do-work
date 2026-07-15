---
id: REQ-018
title: "Harden work.md steps 2.0 + 5.75 (full) and 5.5 + 6.3 (mechanical parts) into tools/checks/"
status: pending
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: []
batch: harness-bloat-cleanup
maintenance: true
---

# Harden machine-checkable work.md steps

## What
Create shipped, platform-agnostic bash checks under `tools/checks/` and shrink the
corresponding work.md prose to a one-line pointer plus whatever judgment remains:

1. `archive-collision.sh REQ-NNN` — Step 2.0 verbatim: glob both archive patterns,
   exit code + collision path output. Prose → pointer.
2. `preflight.sh` — Step 5.75: git-clean (-uall), test-baseline run (command passed
   in or skipped), deps presence. Warnings to stdout, never blocking. Also emits
   `do-work/working/baseline.json` (test command + failing tests) — groundwork for
   the Step 6.5 probation test, without changing 6.5's prose beyond an optional
   "compare against baseline.json if present" sentence.
3. `scope-drift.sh <req-file>` — Step 5.5's review-time comparison only: set-diff of
   `## Scope` file list vs `## Implementation Summary` file list. Declaration prose stays.
4. `qualify.sh <req-file>` — Step 6.3 checks 1, 4, and the grep part of 5 (files
   exist/in-diff, P-A-U boxes vs debug-artifact grep, wiring grep with the
   documented exception list). Checks 2/3/6 stay prose; the script feeds them evidence.

All checks degrade gracefully (missing git, no tests → warn and continue) per the
design-for-the-floor rule. Steps 3.5, 3.7, 6.25, 7.5 are NOT converted (judgment —
audit HARDEN table). Update the Orchestrator Checklist lines and any Common
Rationalizations rows that cite the converted steps; add pointer-sync assertions to
`_dev/tests/contract-regressions.sh`.

## Why
Prose that transcribes mechanical shell logic drifts and taxes every read of
work.md (10,371 words). Audit HARDEN bucket.

## Acceptance criteria
- [ ] Four scripts, shellcheck-clean, each runnable standalone with usage text.
- [ ] work.md steps 2.0/5.75 reduced to pointer + intent; 5.5/6.3 keep judgment
      prose, mechanical parts point at scripts. Net work.md word count drops.
- [ ] Behavior parity: script output covers every condition the prose prescribed
      (map recorded in Implementation Summary).
- [ ] `_dev/tests/contract-regressions.sh` gains assertions that the pointers and
      scripts stay in sync; suite passes.

## Open Questions
(none — placement resolved at capture: tools/checks/, shipped)

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:**
- [ ] **[APPLY]:**
- [ ] **[UNIFY]:**
