---
id: REQ-317
title: "Review fix: Preserve canonical-gate holds in error handling"
status: pending
created_at: 2026-08-21T16:17:46Z
user_request: UR-055
addendum_to: REQ-309
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
depends_on: []
maintenance: true
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
- skills/do-work/actions/work.md
- _dev/tests/contract-regressions.sh
---

# Preserve Canonical-Gate Holds in Error Handling

## AI Execution State (P-A-U Loop)

- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Reconcile the Error Handling table's generic repeated-test-failure archive rule with Step 6.5's new
canonical-gate exception. An unrelated or pre-existing failure of a project-declared canonical
repository gate must preserve the claimed REQ and checkpoint for resumption, never fall through to
the generic Code-failure archive path.

Done means the shipped action has one consistent outcome for this failure class and the semantic
contract detects an opposing downstream directive.

## Context

Found during the independent review of REQ-309. Step 6.5 now says to preserve and stop, while the
later Error Handling table still says every repeatedly failing test is archived as failed.

Fold-first scan: no pending REQ or sweep in any UR shares this root cause.

## Requirements

- Narrow the generic repeated-test-failure row so it cannot consume unrelated or pre-existing
  canonical repository gate failures.
- Extend the REQ-309 semantic regression to cover the downstream Error Handling restatement, not
  only the Step 6.5 lane.
- Preserve the existing three-attempt Code-failure path for focused tests and current-diff
  regressions.

## Red-Green Proof

**RED prompt/case:** Extend the REQ-309 canonical-gate semantic contract to inspect the Error
Handling section; the current broad `Tests fail repeatedly` archive row must make it fail.

**Why RED now:** The row archives all repeated test failures without excluding the canonical gate's
unrelated/pre-existing preserve-and-stop path.

**GREEN when:** The same semantic contract passes only when Step 6.5 and Error Handling agree that
this canonical-gate failure preserves the claim and checkpoint, while ordinary current-diff test
failures retain their existing remediation/archive behavior.

**Validation:** Review finding from REQ-309; apply `actions/work-reference.md` → Finding-Closure
Ratchet (Step 6.5).
