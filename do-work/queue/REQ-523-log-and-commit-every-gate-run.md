---
id: REQ-523
title: 'Log and commit every maintainer gate run'
status: pending
created_at: 2026-09-02T21:27:16Z
user_request: UR-100
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-518, REQ-519, REQ-520, REQ-521, REQ-522]
batch: cheap-maintainer-gate
write_set: [_dev/tests/maintainer-verify.sh, _dev/gate-runs/, CLAUDE.md]
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-02T21:34:07Z
  basis:
    - trivial short-circuit
---

# Log and Commit Every Maintainer Gate Run

## What

Every run of `bash _dev/tests/maintainer-verify.sh` writes one tracked log file recording when it ran, how long it took, its exit status, which lane ran, and the complete output, and then commits that file by exact path as a bookkeeping commit. Capture-time answer Q5: the gate script commits its own log entry.

The fold-first scan found no pending REQ that logs gate runs.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The maintainer asked for a record of when the suite ran, how much time it took, and what it printed, committed so it travels with the history. It is also the measurement against which REQ-518 to REQ-522 prove their savings.

## Context

The gate runs both by hand and inside the work pipeline on a dirty tree, so the commit must stage the log path only. Repository-only history uses `_dev/`, which is export-ignored and never shipped. The self-test shims `git`, so the log write and commit need a fixture stage there too.

## Detailed Requirements

- One file per run under `_dev/gate-runs/`, named by the UTC start instant and lane (full or `--changed`), holding a small header (start, end, duration in seconds, exit status, lane, `HEAD` revision, hostname) followed by the complete captured stdout and stderr of the run.
- The header is written in a fixed line-per-field form a shell one-liner can grep; an index file listing one line per run (start, duration, status, lane, file) is appended in the same commit.
- The log is written and committed on every exit path, including a failing stage, an interrupted run, and the version-floor refusals.
- The commit stages exactly the new log file and the index, with message `[gate] <lane> <status> in <duration>s at <start>`; a dirty tree elsewhere is untouched. If the commit cannot be made (no repository, detached fixture, commit refused), the log file still exists and the gate's own exit status is unchanged.
- The self-test asserts a log file and commit exist after the all-success fixture and after one failing-stage fixture.
- `CLAUDE.md` § Verify gains one sentence naming the log location.
- The log commit happens before the green-revision record of REQ-518 (Run the full gate once per REQ) is written, so that record names the log commit itself; the two REQs share this rule and whichever lands second wires it. (verify-requests, 2026-09-03)

## Constraints

- The gate's exit status is never changed by logging or by a failed log commit.
- Output is captured by `tee`, never by re-running a stage.
- No hand-maintained counts anywhere in the log shape.

## Batch Constraints

- Done means, measured on the maintainer's machine: the full uncached gate under 3 minutes, and a REQ that touches only action Markdown or one Go module gets a fast lane under 60 seconds.
- The full gate is never waived for the integrating commit. The fast lane is a per-REQ check, never the release check.
- Mechanics stay in Go or in the gate script; no new prose that walks a shell sequence.
- Every REQ carries a behavior test or a self-test stage, never a sentence pin alone. `_dev/tests/contract-regressions.sh` does not grow past its current line count (8,417).
- Write sets overlap with REQ-469, REQ-470 and REQ-471 (gate-failure flow in `work.md`); overlap is declared, not a dependency.

## Dependencies

None. Independent of the other five; may run in parallel.

## Builder Guidance

Certainty level: Firm on the record's fields and the self-commit; latitude on file naming and the index format.

## Red-Green Proof

**RED prompt/case:** `bash _dev/tests/maintainer-verify.sh; ls _dev/gate-runs/; git log -1 --format=%s`.
**Why RED now:** No log directory exists and the last commit is unrelated to the run.
**GREEN when:** A new file under `_dev/gate-runs/` holds the run's header and full output, the index has one new line, the last commit is `[gate] full passed in <n>s at <start>` staging exactly those two paths, and the self-test fixtures for a green and a red run both find their log and commit.
**Validation:** User confirmed (capture-time answer Q5, 2026-09-03).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing shipped shell and prescribed command blocks.

## Full Context

See `do-work/user-requests/UR-100/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-03 on `_dev/tests/maintainer-verify.sh` taking 6.5 minutes, item A7 (gate run log, message 4) of the analysis report's improvements, captured by UR-100.*
