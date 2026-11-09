---
source_type: req_lesson
req_id: REQ-413
req_path: do-work/archive/REQ-413-implement-capture-answer-release-transactions.md
date: 2026-09-01
domain: general
module: _dev/primes
tags: [general, implement, capture, file, answer]
---

# Lessons from REQ-413: Implement capture-file, answer, release, version, and changelog transactions

## What the REQ was about

Move deterministic publication and resolution phases for capture, answers, and releases into Go.

## Solution summary

- **Implementation commit:** `db7bb7c8`
- **Lifecycle result:** `completed-with-issues` after one remediation and one fresh re-review.
- **Critical blocker:** REQ-457 owns transaction-created-path rollback identity.
- **Other routed work:** REQ-460 (delimiter containment), REQ-461 (release ownership), and REQ-419 (shell-safe Just recipe arguments).

## Worth knowing

- A rooted forward write is not a complete filesystem safety boundary: rollback must prove that the object it removes is the same object this invocation created, even after every later mutation point.
- Safety predicates stated as conditions cannot be implemented as finite example-prefix or directory-name lists; tests must generate representatives across the condition's classes.
- Human-runnable recovery recipes need shell-safe argument rendering independently of exact machine argv.

## Back-reference

See `do-work/archive/REQ-413-implement-capture-answer-release-transactions.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `db7bb7c8`.
