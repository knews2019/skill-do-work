---
id: REQ-520
title: 'Run the board JavaScript probes once per gate run'
status: cancelled
created_at: 2026-09-02T21:27:16Z
user_request: UR-100
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-518, REQ-519, REQ-521, REQ-522, REQ-523]
batch: cheap-maintainer-gate
write_set: [skills/do-work-board/tools/queue-kanban/generate_test.go, skills/do-work-board/tools/queue-kanban/browser_probe_test.go, _dev/tests/maintainer-verify.sh]
completed_at: 2026-09-03T13:21:16Z
---

# Run the Board JavaScript Probes Once per Gate Run

## What

The 55 `TestJavaScriptBehavior*` probes execute three times per gate run: in the ordinary `go test ./...` pass (about 70 seconds), inside the two `RejectsZeroProbes` meta-tests that re-execute the whole probe set (43 seconds and 16 seconds), and again in the strict lane (51 seconds). Run them once: set the strict marker on the ordinary run whenever Node is present so `TestMain`'s zero-probe guard applies to that run, delete the separate strict lane from the gate, and narrow the two meta-tests to re-execute one cheap probe instead of the whole set.

The fold-first scan found no pending REQ touching the strict lanes.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

About 110 seconds of every 6.5-minute gate run re-proves that a skip cannot be mistaken for green. That property is worth one cheap test.

## Context

`generate_test.go` lines 234 to 265 hold the two lane tests; `TestMain` at line 198 holds the guard for both the JavaScript and browser markers. The gate's self-test asserts a `board-strict` stage exactly once; that stage list shrinks by one. The browser lane keeps its shape (it runs only when a Chrome-family binary is present) but its `RejectsZeroProbes` meta-test narrows the same way.

## Detailed Requirements

- `maintainer-verify.sh` sets the strict JavaScript marker on the ordinary queue-kanban test run when `node` is on PATH, and prints the existing SKIP line otherwise; the separate strict-lane invocation and its self-test stage are removed.
- `TestMain`'s guard is unchanged: a strict run in which no probe executed exits non-zero with the existing diagnostic.
- `TestMaintainerStrictJavaScriptBehaviorLaneRejectsZeroProbes` and `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes` re-execute one named probe each, not `^TestJavaScriptBehavior` or `^TestBrowserBehavior`.
- Each probe's behavior coverage is unchanged; no probe is deleted.

## Constraints

- Keep the zero-probe protection; only remove the repetition.
- The browser lane's selection rule (`QUEUE_KANBAN_BROWSER` or a Chrome-family binary on PATH) is unchanged.

## Batch Constraints

- Done means, measured on the maintainer's machine: the full uncached gate under 3 minutes, and a REQ that touches only action Markdown or one Go module gets a fast lane under 60 seconds.
- The full gate is never waived for the integrating commit. The fast lane is a per-REQ check, never the release check.
- Mechanics stay in Go or in the gate script; no new prose that walks a shell sequence.
- Every REQ carries a behavior test or a self-test stage, never a sentence pin alone. `_dev/tests/contract-regressions.sh` does not grow past its current line count (8,417).
- Write sets overlap with REQ-469, REQ-470 and REQ-471 (gate-failure flow in `work.md`); overlap is declared, not a dependency.

## Dependencies

None. Independent of REQ-521 and REQ-522; all three may run in parallel.

## Builder Guidance

Certainty level: Firm. Read the Kanban prime's versioning and test notes first.

## Red-Green Proof

**RED prompt/case:** `time bash _dev/tests/maintainer-verify.sh` and count how many times the queue-kanban probe set runs.
**Why RED now:** Three times; the queue-kanban stages take 211 seconds.
**GREEN when:** The probes run exactly once per gate run, the queue-kanban stages take under 100 seconds on the maintainer's machine, and a strict run under an empty PATH still fails with the zero-probe diagnostic (the narrowed meta-test proves it).
**Validation:** User confirmed (verify-requests, 2026-09-03: the per-stage threshold stands as the GREEN target).

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5083 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing queue-kanban testing and browser behavior.
- `_dev/primes/lessons-kanban-board.md` — 4707 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing queue-kanban testing.

## Full Context

See `do-work/user-requests/UR-100/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-03 on `_dev/tests/maintainer-verify.sh` taking 6.5 minutes, item A3 of the analysis report's improvements, captured by UR-100.*

## Cancelled

- **When:** 2026-09-03T13:21:16Z
- **Why:** one board test run now carries the strict JavaScript marker; the separate lane and its test are deleted: landed in place by the maintainer's step-2 gate batch, commit 5e0e166c (release 0.266.9); recapture only if the behavior regresses (maintainer decision, 2026-09-03)
- **Decided by:** user, via `do-work abandon`
