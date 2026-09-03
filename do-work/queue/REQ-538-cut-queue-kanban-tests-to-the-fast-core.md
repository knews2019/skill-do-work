---
id: REQ-538
title: 'Cut queue-kanban tests to the fast core and parallelize the rest'
status: pending
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: false
suggested_spec:
depends_on: [REQ-537]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
write_set:
  - skills/do-work-board/tools/queue-kanban/*_test.go
  - _dev/tests/maintainer-verify.sh
---

# Cut queue-kanban Tests to the Fast Core and Parallelize the Rest

## What

Bring the queue-kanban fast run under 30 s wall. Delete the two `RejectsZeroProbes` meta-tests that re-execute the whole test binary (41.8 s and 16.1 s) and replace each with one cheap assertion of the guard; make the 55 JavaScript probes skip in the fast tier through one explicit environment knob honored in `lookupNodeForJavaScriptProbe`, never by hiding `node`; audit the remaining tests, delete duplicates and decorative ones, and opt the rest into `t.Parallel()` where they own no process-global state. Folds REQ-522.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

- 457 tests, 141 s of test time; after the two meta-tests and the probes the rest is about 35 s. 105 tests live in `generate_test.go` (10,093 lines), 61 in `verify_test.go` (14 real git repos), 8 in `board_live_test.go` walk the live `do-work/` tree per test.
- The three re-exec sites that inherit the strict marker (`browser_probe_test.go`, `citations_test.go`, `allocate_test.go`) clear it since 0.266.9 and must keep doing so.
- `TestMain` holds the zero-probe guards for both lanes; the guard function is what deserves a unit test, not a re-exec of the suite.

## Detailed Requirements

- Delete `TestMaintainerStrictJavaScriptBehaviorLaneRejectsZeroProbes` and `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes`; assert each guard through a child that runs exactly one named test, or by extracting the guard into a function under unit test.
- One knob, for example `QUEUE_KANBAN_JAVASCRIPT_PROBES=off`, makes every `TestJavaScriptBehavior*` probe `t.Skip`; REQ-537's default run sets it, `--heavy` does not.
- Audit: a test is deleted when another test already proves the same failure, when it pins a rendered string with no behavior behind it, or when its fixture cost is out of proportion to what it pins; list each in the commit body.
- `t.Parallel()` on every remaining test that owns no process-global state (no `os.Chdir`, no shared listener port, no live-tree walk); `go test -race -count=1 ./...` stays green.
- Under `--heavy` all 55 JavaScript probes still execute at least once.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** `cd skills/do-work-board/tools/queue-kanban && time go test -count=1 ./...` with `node` on PATH.
**Why RED now:** about 150 s wall; the two meta-tests alone are 58 s and no test opts into `t.Parallel`.
**GREEN when:** the same command with the fast-tier knob set finishes under 30 s wall, `go test -race -count=1 ./...` is green, and `--heavy` still reports the JavaScript probe count above zero.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5562 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes queue-kanban testing and browser behavior.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.

