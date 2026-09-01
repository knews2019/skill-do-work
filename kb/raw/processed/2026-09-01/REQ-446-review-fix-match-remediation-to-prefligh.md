---
source_type: req_lesson
req_id: REQ-446
req_path: do-work/archive/REQ-446-match-remediation-to-preflight-failure-kind.md
date: 2026-09-01
domain: general
module: skills/do-work/tools/do-work-cli
tags: [general, review, match, remediation, preflight]
---

# Lessons from REQ-446: Review fix: Match remediation to preflight failure kind

## What the REQ was about

Make every caller that projects a shared transaction-preflight failure choose recovery and verification commands from the actual failure kind. Done means the class cannot recur: simultaneous-state regressions must reject guidance that inspects a target when the selected blocker is repository-wide index state.

Fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this preflight-failure-kind remediation root cause. REQ-445 covers pathless structural cleanup findings, not failure-kind projection in shared preflight callers.

## Solution summary

- **Implementation commit:** `6f173a12`
- **Lifecycle result:** completed after one remediation and fresh passing re-review.
- **Shared seam:** doctor and cleanup both derive actionable preflight guidance from `gittransaction.BuildCommandResult`.

## Worth knowing

- A test name is not closure evidence: the fixture must keep every simultaneous blocker present through the one invocation whose precedence it claims to prove.
- Shared failure kinds need one canonical remediation renderer, while callers may layer domain identity around that evidence without replacing it.

## Back-reference

See `do-work/archive/REQ-446-match-remediation-to-preflight-failure-kind.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `6f173a12`.
