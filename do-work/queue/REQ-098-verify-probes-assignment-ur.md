---
id: REQ-098
title: "Verify probes: assigned-elsewhere-claimed-here and ur-archived-with-live-member"
status: pending
created_at: 2026-08-04T19:44:17Z
user_request: UR-018
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-097]
maintenance: false
related: [REQ-097]
batch: parallel-building
write_set: [tools/queue-kanban/verify.go, tools/queue-kanban/verify_test.go]
---

# Verify Probes: Assignment Drift and UR Closure Drift

## What

Two new read-only probes in `queue-kanban verify`, extending the report-and-route contract (verify never repairs — fixes belong to `actions/cleanup.md`):

1. `assigned-elsewhere-claimed-here` — a REQ carrying `assigned_to` is sitting in `do-work/working/` (someone claimed work earmarked for another session without clearing the field).
2. `ur-archived-with-live-member` — an archived UR has a member REQ (by `user_request:` scan) still in `queue/` or `working/` (a silent-merge class also reachable from a botched cleanup).

## Detailed Requirements

- Implement in `tools/queue-kanban/verify.go`, ~30 lines each, following the existing probe structure (finding code, message, routing suggestion). Tests in `verify_test.go` with fixture trees for both the firing and non-firing case.
- Probe 2's membership scan uses `user_request:` frontmatter across `queue/`, `working/`, `archive/` root and `archive/UR-NNN/` — the same closure predicate `actions/work.md` Step 8 evaluates (the UR `requests:` array is capture-time-only, never the closure predicate).
- Route findings to `actions/cleanup.md` in the report text, matching existing probes' phrasing.
- `go test ./...` green.

## Red-Green Proof

**RED prompt/case:** A fixture with an `assigned_to` REQ in `working/`, and one with an archived UR whose member REQ sits in `queue/` — `queue-kanban verify` today reports neither.
**Why RED now:** No probe covers assignment drift (the field is new) or archived-UR/live-member drift.
**GREEN when:** Both probes fire on their fixtures, stay silent on clean trees, and the full Go test suite passes.
**Validation:** User confirmed (approved plan, Phase 2 item 7).

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 2).

---
*Source: approved plan, Phase 2*
