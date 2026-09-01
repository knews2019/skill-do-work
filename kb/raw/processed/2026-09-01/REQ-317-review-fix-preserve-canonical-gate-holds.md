---
source_type: req_lesson
req_id: REQ-317
req_path: do-work/archive/UR-055/REQ-317-preserve-canonical-gate-holds-in-error-handling.md
date: 2026-08-21
domain: general
module: _dev/primes
tags: [general, review, preserve, canonical, gate]
---

# Lessons from REQ-317: Review fix: Preserve canonical-gate holds in error handling

## What the REQ was about

Reconcile the Error Handling table's generic repeated-test-failure archive rule with Step 6.5's new
canonical-gate exception. An unrelated or pre-existing failure of a project-declared canonical
repository gate must preserve the claimed REQ and checkpoint for resumption, never fall through to
the generic Code-failure archive path.

Done means the shipped action has one consistent outcome for this failure class and the semantic
contract detects an opposing downstream directive.

## Solution summary

Reconciled the downstream error-disposition table with Step 6.5 while retaining
the ordinary remediation/follow-up/archive route.

## What worked

- Extending REQ-309's existing semantic block made the downstream reader part of the same policy
  contract instead of creating a second checker with subtly different vocabulary.
- Extracting Step 6.5 and Error Handling independently proved that correct upstream text cannot mask
  a contradictory generic row later in the action.

## What didn't work

- REQ-309 originally tested only the newly added lane. That left a broader later directive free to
  reverse its failure disposition even though the local contract was completely green.

## Worth knowing

A preservation exception must narrow downstream catch-all error handlers, not
only be stated where the exception originates. Attributable current-diff failures still use the
ordinary remediation/Code/archive path; the hold is only for unrelated or pre-existing canonical
gate failures.

## Back-reference

See `do-work/archive/UR-055/REQ-317-preserve-canonical-gate-holds-in-error-handling.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a9259d7`.
