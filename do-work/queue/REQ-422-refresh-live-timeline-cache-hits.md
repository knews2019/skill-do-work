---
id: REQ-422
title: 'Refresh time-derived timeline data on live cache hits'
status: pending
created_at: 2026-08-29T20:26:10Z
user_request: UR-082
domain: backend
prime_files: ['_dev/primes/prime-kanban-board.md']
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-421, REQ-423, REQ-424]
batch: accepted-review-fixes
write_set: ['skills/do-work-board/tools/queue-kanban/generate.go', 'skills/do-work-board/tools/queue-kanban/serve.go', 'skills/do-work-board/tools/queue-kanban/serve_test.go', '_dev/primes/prime-kanban-board.md', '_dev/primes/lessons-kanban-board.md']
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
---

# Refresh Time-Derived Timeline Data on Live Cache Hits

## What
Make an unchanged-tree live-board cache hit refresh the complete time-derived Timeline payload from cached parsed tickets, not only the top-level generation timestamp.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Extract Timeline payload construction into one helper accepting tickets, duration history, and an explicit instant.
- Add an injectable live-server clock.
- On unchanged-tree cache hits, set `GeneratedAt` and rebuild all Timeline fields from the cached parsed board at the same fresh instant.
- Do not reparse files or Markdown on the cache-hit path.

## Constraints
- Preserve the current live tree fingerprint cache.
- Make no public JSON schema or API changes.
- Add concise cache-invalidation guidance to the Kanban prime and a detailed linked lesson entry.

## Dependencies
None. It ships in the same release batch as REQ-421, REQ-423, and REQ-424.

## Builder Guidance
Certainty: Firm. Use one captured instant per refresh so `GeneratedAt`, Timeline `now`, open spans, range bounds, and projection timestamps remain coherent.

## Context
No pending or unassigned queue candidate shares this root cause. Provenance: accepted review finding `[P2] Refresh time-derived timeline data on live cache hits` against `skills/do-work-board/tools/queue-kanban/generate.go:781`. The review states that unchanged fingerprints currently refresh only `GeneratedAt`, leaving `timeline.now`, open spans, and forecast times frozen.

## Red-Green Proof
**RED prompt/case:** Seed a live server cache at `t0`, advance an injected clock to `t1` without changing the tree, then refresh.
**Why RED now:** Only the top-level timestamp advances; the cached Timeline remains anchored at `t0`.
**GREEN when:** `GeneratedAt`, Timeline `now`, open wait/work spans, range end, and projection timestamps all advance coherently to `t1` while the cached parsed board is reused.
**Validation:** User accepted the review finding and supplied the implementation plan.

## Full Context
See `do-work/user-requests/UR-082/input.md` for the approved plan and batch constraints.

---
*Source: accepted review finding [P2] on live Timeline cache freshness, followed by the user-approved plan.*
