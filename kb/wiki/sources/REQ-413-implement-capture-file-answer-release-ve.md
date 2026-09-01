---
title: "Lessons from REQ-413: Implement capture-file, answer, release, version, and changelog transactions"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-413-implement-capture-file-answer-release-ve.md]
related:
  - page: concept-prescribed-shell-commands
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-413: Implement capture-file, answer, release, version, and changelog transactions

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Move deterministic publication and resolution phases for capture, answers, and releases into Go.

## Solution summary

- **Implementation commit:** `db7bb7c8`
- **Lifecycle result:** `completed-with-issues` after one remediation and one fresh re-review.
- **Critical blocker:** REQ-457 owns transaction-created-path rollback identity.
- **Other routed work:** REQ-460 (delimiter containment), REQ-461 (release ownership), and REQ-419 (shell-safe Just recipe arguments).

## What worked

- A rooted forward write is not a complete filesystem safety boundary: rollback must prove that the object it removes is the same object this invocation created, even after every later mutation point.
- Safety predicates stated as conditions cannot be implemented as finite example-prefix or directory-name lists; tests must generate representatives across the condition's classes.
- Human-runnable recovery recipes need shell-safe argument rendering independently of exact machine argv.

## Back-reference

See `do-work/archive/REQ-413-implement-capture-answer-release-transactions.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `db7bb7c8`.
