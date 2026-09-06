---
id: REQ-558
title: '[impact-negligible] Keep one nil-root guard in git_transaction.go and delete the other eight'
status: claimed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: backend
prime_files: []
tdd: true
estimate:
  p50_active_minutes: 30
  confidence: low
  calculated_at: 2026-09-06T05:51:12Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - transaction boundary with no nil-branch coverage
suggested_spec:
depends_on: [REQ-557]
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557]
batch: maintainability-audit-2026-09-03
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
route: B
write_set: [skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-06T05:50:44Z
---

# Keep one nil-root guard in git_transaction.go and delete the other eight

## What
`internal/gittransaction/git_transaction.go` tests one `*os.Root` value for nil nine times; no test exercises any nil branch, no introducing commit or REQ names a nil-root failure, three caller⇒callee pairs test the same value twice on one path, and two functions taking the same value never test it. Keep one nil test at the single point where nil is producible (`rollbackTransaction`'s `os.OpenRoot`), pass a non-nil root or an explicit no-root branch downward, and delete the other eight guards.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Defensive checks with no incident behind them are the agent-creep class; here they accreted per REQ on the top do-work-cli hotspot (17 commits, file CCN 392, `ExecuteTransaction` CCN 87) and none is covered by a test.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 3, sweep_key `nil-root-guards-git-transaction`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -25. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- Guards at `inspectCreatedObject` (b877eb69, REQ-457: the REQ never mentions a nil root), `rollbackDirtyTracked` (0a5d4e44, REQ-491: same), `trackedPublicationStillOwned` (a43b2587, no REQ id), `rootedOpenSnapshot` (01d920dd, no REQ id), and two `root != nil && privateStateStillOriginal(root, state)` sites (01d920dd).
- Redundant pairs verified by reading: the two `privateStateStillOriginal` callers, `trackedPublicationStillOwned`, and `rollbackDirtyTracked` all reach `rootedOpenSnapshot`, which tests the same value again.
- Inconsistency evidence: `privateStateStillOriginal` and `rootedCreateRegular` take the same `root *os.Root` and do not test it.
- Behaviour preserved on every rollback path: the package's transaction and rollback tests are the safety net and stay green unchanged.
- Reproduce at dc8a64e3 (prints 9 guards and `NO TEST covers any nil-root branch`): `rg -n 'root [=!]= nil' skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go; rg -l 'OpenRoot|rollback root is unavailable|rooted filesystem handle is unavailable' skills/do-work/tools/do-work-cli/internal/gittransaction/*_test.go || echo 'NO TEST covers any nil-root branch'`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- This file only; the rollback path is a transaction boundary, so every change is behaviour-preserving and proved by the existing tests, never by a new mock.
- Lock-in limit: nil-root guards in git_transaction.go: 1 after this REQ (today 9).

## Dependencies
Depends on REQ-557 so no other audit REQ writes under `internal/` while the transaction file is refactored. Last of the batch.

## Builder Guidance
Firm on one producible-nil point; latitude on whether downstream functions take a non-nil root or an explicit no-root branch.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints nine `root == nil` / `root != nil` sites and no test covers any of them.
**GREEN when:** It prints one site; `go test ./internal/gittransaction/` green unchanged; the lock-in pins nil-root guards in that file at 1.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "do-work-cli internals".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for nil-root-guards-git-transaction.*

## Triage

**Route: B** — Explore then build.

**Reasoning:** The request's baseline holds exactly at HEAD, which is unusual in this batch: the
reproduce command prints nine `root [=!]= nil` sites and `NO TEST covers any nil-root branch`. But the
change is inside a transaction boundary, and the claim that matters is not "there are nine guards" — it
is "nil is producible at exactly one of them". Proving that means tracing every path that can reach each
of the nine with a nil value, and no test exercises any of those branches, so the compiler and the
existing suite cannot settle it. That is discovery.

**Planning:** Skipped. One file plus one lock-in assertion; the work is whatever the trace establishes.

**Deleting a guard is not the same shape as deleting a duplicate.** REQ-557 removed copies that were
provably interchangeable. Here each guard is the only thing standing between a nil dereference and a
rollback path. A guard that is genuinely unreachable costs nothing to delete and a guard that is not
costs a panic during rollback, which is the worst moment available. The exploration's job is to tell
those two apart per site, and to say plainly where it cannot.

## Plan

**Planning not required** — Route B: one source file plus one lock-in assertion, and the edit set is
whatever the reachability trace establishes.

*Skipped by work action*
