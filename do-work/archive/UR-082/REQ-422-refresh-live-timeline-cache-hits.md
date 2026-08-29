---
id: REQ-422
title: 'Refresh time-derived timeline data on live cache hits'
status: completed
claimed_at: 2026-08-29T20:31:49Z
completed_at: 2026-08-29T20:46:08Z
route: B
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
- [x] **[PLAN]:** Extract Timeline construction around one explicit instant, inject a live clock, and prove an unchanged-tree hit advances every time-derived field without reparsing.
- [x] **[APPLY]:** Extracted explicit-instant Timeline construction, injected the live-server clock, and rebuilt the complete Timeline from cached parsed tickets on unchanged-tree hits.
- [x] **[UNIFY]:** Reviewed the complete diff, ran gofmt, the focused cache regression, the full uncached queue-kanban suite, and the canonical maintainer verifier; all required lanes passed.

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

---

## Triage

**Route: B** - Medium

**Reasoning:** The existing cache and timeline seams are clear, but correctness spans several derived fields and needs deterministic clock-driven coverage.

**Planning:** Not required; the user supplied an implementation plan.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — centralize Timeline construction
- `skills/do-work-board/tools/queue-kanban/serve.go` (modify) — inject clock and refresh Timeline on cache hits
- `skills/do-work-board/tools/queue-kanban/serve_test.go` (modify) — deterministic cache-hit regression
- `_dev/primes/prime-kanban-board.md` (modify) — concise cache guidance
- `_dev/primes/lessons-kanban-board.md` (modify) — linked detailed lesson

**Files I will NOT touch:** browser Timeline schema or rendering.

**Acceptance criteria (restated from REQ):**
- [x] One fresh instant drives `GeneratedAt` and the rebuilt Timeline.
- [x] Open spans, range end, and projection times advance on a clean cache hit.
- [x] Parsed files and Markdown remain cached.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified)
- `skills/do-work-board/tools/queue-kanban/serve_test.go` (modified)
- `_dev/primes/prime-kanban-board.md` (modified)

**What was done:** Timeline payload construction now accepts parsed tickets, duration history, and one explicit instant. The live cache-hit path shallow-copies cached board data, refreshes `GeneratedAt`, and rebuilds Timeline spans and projections from cached tickets without reparsing the tree or Markdown.

## Testing

**Tests run:** Focused live cache regression and adjacent serve tests; `go test -count=1 ./...`; `bash _dev/tests/maintainer-verify.sh`.

**Result:** All passed. The uncached module suite passed in 131.123s; the canonical verifier exited 0.

**Red-green validation:**
- `TestServeMtimeCacheRefreshesTimelineAgainstCurrentTime`: RED with Timeline `now` frozen at 10:00 instead of 10:30 → GREEN with `GeneratedAt`, `now`, open spans, range end, and projection timestamps advancing together.

**New tests added:**
- Deterministic injected-clock regression for an unchanged-tree cache hit with pending, claimed, and completed projection samples.

## Lessons Learned

A filesystem fingerprint answers whether parsed source state changed; it says nothing about derived values whose other input is wall time. Cache-hit code must inventory those volatile inputs and rebuild every dependent field from cached parsed state using one captured instant, or the response can claim a fresh generation time while its internal clock-dependent views remain frozen.
