---
title: "Lessons from REQ-325: Stop the report-image interruption path orphaning its backend"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-325-stop-the-report-image-interruption-path-.md]
related:
  - page: concept-completed-work-presentation
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-325: Stop the report-image interruption path orphaning its backend

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

An interrupted `generate-report-image.sh` leaves its image backend running. In the repo's own
test suite that hangs the canonical gate indefinitely; in a consumer's checkout it leaves a
process behind after an interrupted `do-work-toolbox ai-report`.

## Solution summary

Closed the wrapper's publish-the-PID window by deferring HUP/INT/TERM across
both backend launches — the status is recorded, the PID and process group are registered, then the
interruption is re-raised so the EXIT trap has a backend to terminate — and moved the staging
`mktemp` below the trap block so an interruption can never leave the invocation-private file
behind. Replaced both interruption cases' unbounded `wait` with `wait_for_wrapper_or_fail`, a
10-second deadline that names the processes still alive at expiry, fails the case, and kills them.
Added two cases: an interruption fired the moment the staging file appears, and a deadline probe
run against a deliberately TERM-deaf stand-in that asserts the diagnostic rather than hanging.

## What worked

- Injecting a stall into a *copy* of the script under test, one candidate window at a time. Three
  candidate mechanisms were on the table and reasoning could not separate them; a `sleep` at each
  suspected point separated them in minutes and refuted two of the three.
- Writing the new fixture case before reading the wrapper for a second defect. The
  early-interruption case failed on its first run for a reason nobody had predicted (the staging
  file predating its own cleanup trap) — the case found a bug the analysis had walked past twice.

## What didn't work

- Trying to reproduce the reported hang by running the suite, solo and two-up, and by adding six
  spinning load generators. Zero reproductions across ~10 suite runs and 160 targeted rounds. The
  stall is real (two independent occurrences) but nothing here reaches it, so it stayed
  unattributed rather than being explained by the first plausible story.
- Reasoning about bash signal semantics from memory. Three separate confident conclusions —
  that a trap cannot fire from inside `wait`, that a group kill would miss the stub, that a
  TERM-deaf backend would wedge the wrapper — were all wrong, and each cost more time than the
  probe that refuted it. `wait` returns early on a trapped signal, group kills reach the stub, and
  a direct child is always reapable by KILL.
- Chasing a fixture that could pin the launch window. The parent won that race 160/160 times, so a
  stress case would have been a test that cannot fail — the exact shape this repo's prime warns
  reads as coverage while locking in nothing.

## Worth knowing

- **A handle published one command after the launch is not a handle a trap can rely on.** Traps run
  *between* commands, so `cmd & pid=$!` has a window where the trap sees no pid. Any cleanup keyed
  on such a variable needs the interruption deferred across the window, not just a trap installed
  before it.
- **A file created before its cleanup trap exists is a file no trap owns.** An EXIT trap does not
  run when a signal takes its default action, so the gap between `mktemp` and the first
  HUP/INT/TERM trap is a leak window regardless of how good the EXIT handler is. Create the
  artifact after the traps, not before.
- **A wait status cannot tell "stuck" from "interrupted" once a watchdog is involved** — both are
  nonzero. The deadline verdict has to come from a separate signal (here, a report file the
  watchdog writes), and the survivor list has to be captured before the kill or it is empty.
- The environment this ran in had none of the gate's three required tools at the right version
  (Go ≥ 1.26.1, ShellCheck ≥ 0.11.0, `just`). `go env -w GOTOOLCHAIN=go1.26.1` works where
  `go.dev` is blocked by network policy, because the toolchain also resolves as a module through
  `proxy.golang.org`. The gate's browser lane still skips unless `QUEUE_KANBAN_BROWSER` names an
  engine — here `/opt/pw-browsers/chromium-1194/chrome-linux/chrome`.

## Back-reference

See `do-work/archive/UR-065/REQ-325-stop-the-report-image-interruption-path-orphaning-its-backend.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `92413b9`.
