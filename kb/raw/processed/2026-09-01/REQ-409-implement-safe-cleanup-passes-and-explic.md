---
source_type: req_lesson
req_id: REQ-409
req_path: do-work/archive/REQ-409-implement-safe-cleanup.md
date: 2026-08-30
domain: general
module: _dev/primes
tags: [general, implement, safe, cleanup, passes]
---

# Lessons from REQ-409: Implement safe cleanup passes and explicit destructive repairs

## What the REQ was about

Make cleanup a canonical no-LLM command that applies only provably safe repairs by default.

## Solution summary

Added the canonical `cleanup` command and its safe planning/application layer for Passes 0–4, link repointing, exact-target Git transactions, blank-record recovery, and worktree evidence/repairs. Registered the command and changed the natural-language action, guide, and CLI prime to delegate deterministic cleanup to it.

## What worked

Exact rooted evidence, transaction-result relabeling after outcome, and focused adversarial fixtures closed the original safety defects without broadening the write set.

## What didn't work

Modeling each move as locally safe was insufficient; safety also depends on explicit prerequisites between operation groups. Unit coverage of isolated groups missed refusal combinations until end-to-end re-review.

## Worth knowing

Cleanup planners need a dependency graph, not just deterministic ordering: every derived mutation must name the successful operation that makes it valid, while unrelated groups remain independently eligible.

## Back-reference

See `do-work/archive/REQ-409-implement-safe-cleanup.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a57bf51e`.
