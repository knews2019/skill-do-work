---
id: REQ-542
title: 'SessionStart launches the background gate runner'
status: pending
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: [REQ-541]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
write_set:
  - skills/do-work/hooks/session-start.sh
  - _dev/tests/gate-runner.sh
---

# SessionStart Launches the Background Gate Runner

## What

The SessionStart hook launches `_dev/tests/gate-runner.sh` in the background when it is not already running for this checkout, so every session has the gate as evidence attached to revisions instead of a step. The runner never passes `--heavy`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

- The runner exists since 0.266.9 and records green through `record-green-gate`; a pipeline claim that finds HEAD proven green skips its baseline. Nothing starts it today.
- Chosen over a `just gate-watch` recipe at capture (Recommended in the draft; the maintainer did not object).
- The hook must stay fail-soft: a missing runner, a non-maintainer checkout, or a runner already running prints one line and exits 0.

## Detailed Requirements

- Detect an existing runner for this repository root (a pid file under `$TMPDIR/do-work-gate-runs/` is enough) and do not start a second one.
- Only start in this repository (the runner lives under `_dev/`, which is never shipped); an installed consumer's hook must not look for it.
- Print one line naming the log directory when it starts, nothing else.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** Open a new session in this checkout and run `pgrep -f gate-runner.sh`.
**Why RED now:** nothing is running.
**GREEN when:** the runner is running after SessionStart, a second session starts no second runner, and a consumer checkout's hook output is unchanged.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes a shipped hook script.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.

