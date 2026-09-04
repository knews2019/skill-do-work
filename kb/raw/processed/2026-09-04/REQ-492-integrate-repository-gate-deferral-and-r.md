---
source_type: req_lesson
req_id: REQ-492
req_path: do-work/archive/REQ-492-integrate-repository-gate-deferral-and-resumption-into-do-work-run.md
date: 2026-09-02
domain: general
module: _dev/primes
tags: [general, integrate, repository, gate]
---

# Lessons from REQ-492: Integrate repository-gate deferral and resumption into do-work run

## What the REQ was about

Use the canonical lifecycle from REQ-491 inside `do-work run`: establish the repository baseline before implementation, classify late failures against the saved pre-merge revision, run repairs without recursion, and safely resume deferred implementations when their dependency completes.

## Solution summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)
- `skills/do-work/docs/command-line-guide.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

## What worked

- Fresh review across downstream actions caught orchestration contradictions that same-file lexical mutation tests could not; keeping REQ-493’s topology boundary explicit prevented accidental scope absorption.

## What didn't work

- Adding a special path only to `work.md` and its reference is insufficient when qualification, review, result projection, and commit validation are separate authorities. Prose predicates cannot substitute for executable typed-result fixtures.

## Worth knowing

- Any lifecycle exception must be swept through every downstream gate that can reject it, and retry contracts must be defined from actual wire values rather than desired semantic labels.

## Back-reference

See `do-work/archive/REQ-492-integrate-repository-gate-deferral-and-resumption-into-do-work-run.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f9eb65f7`.
