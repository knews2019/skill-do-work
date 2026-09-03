---
id: REQ-521
title: 'Run the aggregate contract sub-suites in parallel'
status: cancelled
created_at: 2026-09-02T21:27:16Z
user_request: UR-100
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-518, REQ-519, REQ-520, REQ-522, REQ-523]
batch: cheap-maintainer-gate
write_set: [_dev/tests/contract-regressions.sh]
completed_at: 2026-09-03T13:21:16Z
---

# Run the Aggregate Contract Sub-Suites in Parallel

## What

`_dev/tests/contract-regressions.sh` runs its 14 behavior sub-suites one after another (about 140 of its 149 seconds). Each is its own process with a private `mktemp` fixture root, so nothing forces the order. Launch them together, buffer each one's output to a file, print each buffer when it finishes, and fail the aggregate if any sub-suite failed.

The fold-first scan found no pending REQ touching the aggregate's dispatch.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Standalone the sub-suites sum to about 160 seconds; the longest single one (update-script-behavior) is 60 seconds. Running them together bounds the aggregate near the longest.

## Context

The inline predicates (about 9 seconds) stay serial and run first. The sub-suite block starts at the `maintainer_verify_probe` line (about line 7292). The self-test invoked there is itself one sub-suite and can join the parallel set.

## Detailed Requirements

- Each sub-suite's stdout and stderr go to its own buffer file; the aggregate prints buffers in a fixed order after all finish, so FAIL lines stay attributable and never interleave.
- Exit status: non-zero if any sub-suite is non-zero or missing, with the same FAIL line the aggregate prints today for a missing probe.
- Concurrency is bounded (a small fixed number or the CPU count); the bound lives in one variable at the top of the block.
- Sub-suites that share a fixture path or the repository's own working tree are identified and kept serial, with a comment naming why.
- Update-script-behavior's 62 fixtures are profiled once and the slowest fixture named in the REQ's Testing section; splitting it is a follow-up, not this REQ.

## Constraints

- No change to any sub-suite's own contents or counts.
- ShellCheck warning-clean; the gate lints this file.

## Batch Constraints

- Done means, measured on the maintainer's machine: the full uncached gate under 3 minutes, and a REQ that touches only action Markdown or one Go module gets a fast lane under 60 seconds.
- The full gate is never waived for the integrating commit. The fast lane is a per-REQ check, never the release check.
- Mechanics stay in Go or in the gate script; no new prose that walks a shell sequence.
- Every REQ carries a behavior test or a self-test stage, never a sentence pin alone. `_dev/tests/contract-regressions.sh` does not grow past its current line count (8,417).
- Write sets overlap with REQ-469, REQ-470 and REQ-471 (gate-failure flow in `work.md`); overlap is declared, not a dependency.

## Dependencies

None. Independent of REQ-520 and REQ-522.

## Builder Guidance

Certainty level: Firm on parallel dispatch with buffered output; latitude on the concurrency bound.

## Red-Green Proof

**RED prompt/case:** `time bash _dev/tests/contract-regressions.sh` on a green tree, then make one sub-suite fail (for example, remove a shipped script's executable bit in a scratch copy) and run again.
**Why RED now:** The green run takes about 149 seconds; the sub-suites start one at a time.
**GREEN when:** The green run takes under 75 seconds on the maintainer's machine, every sub-suite's summary line still appears, and the failing run exits non-zero with that sub-suite's FAIL lines printed contiguously.
**Validation:** User confirmed (verify-requests, 2026-09-03: the per-stage threshold stands as the GREEN target).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing shipped shell and prescribed command blocks.

## Full Context

See `do-work/user-requests/UR-100/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-03 on `_dev/tests/maintainer-verify.sh` taking 6.5 minutes, item A4 of the analysis report's improvements, captured by UR-100.*

## Cancelled

- **When:** 2026-09-03T13:21:16Z
- **Why:** _dev/tests/probe-batch.sh runs the behavioral sub-suites concurrently: landed in place by the maintainer's step-2 gate batch, commit 5e0e166c (release 0.266.9); recapture only if the behavior regresses (maintainer decision, 2026-09-03)
- **Decided by:** user, via `do-work abandon`
