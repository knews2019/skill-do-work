---
id: REQ-522
title: 'Opt queue-kanban tests into t.Parallel'
status: cancelled
created_at: 2026-09-02T21:27:16Z
user_request: UR-100
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-518, REQ-519, REQ-520, REQ-521, REQ-523]
batch: cheap-maintainer-gate
write_set: [skills/do-work-board/tools/queue-kanban/*_test.go]
completed_at: 2026-09-03T14:49:21Z
---

# Opt queue-kanban Tests into t.Parallel

## What

None of the 451 tests in `skills/do-work-board/tools/queue-kanban` calls `t.Parallel()`, so the package runs strictly serial on an 8-core machine: 146 seconds of test time is 160 seconds of wall-clock. Audit the tests for process-global state (working directory, environment variables, the re-exec helper pattern in `allocate_test.go`, shared fixture paths) and mark every test that has none as parallel.

The fold-first scan found no pending REQ touching queue-kanban test parallelism.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The board package is the single longest gate stage. Parallel tests are the standard Go lever and cost no coverage.

## Context

Tests that spawn subprocesses (`os.Args[0]` re-exec), read or set environment variables, or `chdir` cannot run in parallel with each other without care. `generate_test.go` is 10,112 lines and holds most of the package's tests; `t.TempDir()` fixtures are already per-test.

## Detailed Requirements

- Every top-level test that touches only its own `t.TempDir()` and no process-global state calls `t.Parallel()` as its first statement; subtests follow the same rule.
- Tests that must stay serial carry a one-line comment naming the global state they touch.
- `go test -count=1 -race ./...` passes at least once during the REQ to prove no shared-state race was introduced; record it in the Testing section.
- The Node probes may run in parallel with each other (each is its own Node process); their probe counter is already atomic.

## Constraints

- No test is deleted or weakened; no assertion changes.
- `go vet` and `gofmt` stay clean.

## Batch Constraints

- Done means, measured on the maintainer's machine: the full uncached gate under 3 minutes, and a REQ that touches only action Markdown or one Go module gets a fast lane under 60 seconds.
- The full gate is never waived for the integrating commit. The fast lane is a per-REQ check, never the release check.
- Mechanics stay in Go or in the gate script; no new prose that walks a shell sequence.
- Every REQ carries a behavior test or a self-test stage, never a sentence pin alone. `_dev/tests/contract-regressions.sh` does not grow past its current line count (8,417).
- Write sets overlap with REQ-469, REQ-470 and REQ-471 (gate-failure flow in `work.md`); overlap is declared, not a dependency.

## Dependencies

None. Independent of REQ-520 and REQ-521.

## Builder Guidance

Certainty level: Firm on the goal, latitude on which tests stay serial. Prefer marking fewer tests and keeping the race run green over marking everything.

## Red-Green Proof

**RED prompt/case:** `cd skills/do-work-board/tools/queue-kanban && time go test -count=1 ./...`.
**Why RED now:** About 160 seconds wall-clock for 146 seconds of summed test time; zero tests opt in.
**GREEN when:** The same command finishes under 60 seconds on the maintainer's machine with the same pass count, and `go test -count=1 -race ./...` is green.
**Validation:** User confirmed (verify-requests, 2026-09-03: the per-stage threshold stands as the GREEN target).

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5083 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing queue-kanban testing and browser behavior.
- `_dev/primes/lessons-kanban-board.md` — 4707 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing queue-kanban testing.

## Full Context

See `do-work/user-requests/UR-100/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-03 on `_dev/tests/maintainer-verify.sh` taking 6.5 minutes, item A5 of the analysis report's improvements, captured by UR-100.*

## Cancelled

- **When:** 2026-09-03T14:49:21Z
- **Why:** superseded by UR-104: REQ-538 owns the queue-kanban audit and t.Parallel opt-in (maintainer decision, 2026-09-03)
- **Decided by:** user, via `do-work abandon`
