---
id: REQ-325
title: "Stop the report-image interruption path orphaning its backend"
status: pending
status_changed_at: 2026-08-23T11:42:00Z
created_at: 2026-08-23T02:09:42Z
user_request: UR-065
addendum_to: REQ-321
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-toolbox/scripts/generate-report-image.sh
  - _dev/tests/prescribed-shell-cases/generate-report-image.sh
---

# Stop the Report-Image Interruption Path Orphaning Its Backend

## What

An interrupted `generate-report-image.sh` leaves its image backend running. In the repo's own
test suite that hangs the canonical gate indefinitely; in a consumer's checkout it leaves a
process behind after an interrupted `do-work-toolbox ai-report`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`bash _dev/tests/maintainer-verify.sh` is the only proof this project accepts before a
hand-back. It can hang forever, and it did so twice on 2026-08-22/23 in two independent
sessions:

- REQ-320's review agent sat roughly 35 minutes in that lane and had to kill the stub by hand
  before its gate would finish.
- REQ-321's gate stalled 9+ minutes on the identical process; the run completed normally the
  moment the stub was killed.

## Context

Both times the same three processes were alive together:

- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` — the case script
- `skills/do-work-toolbox/scripts/generate-report-image.sh` — the shipped wrapper
- the fixture's stub `imagegen`, spinning in `while :; do sleep 0.1; done`

The case (`_dev/tests/prescribed-shell-cases/generate-report-image.sh:79-92`) starts the
wrapper in the background, waits for the fixture's ready marker, then TERMs **the wrapper**
and `wait`s on it. The stub traps TERM and exits 143 — but only if it receives one.

The signal-forwarding machinery exists and did not fire. The wrapper carries
`trap 'exit 143' TERM` (`skills/do-work-toolbox/scripts/generate-report-image.sh:106`) and a
helper that signals the backend's process group (`:63-74`) before `wait`ing on it (`:89`).
Which side fails to deliver — the wrapper never running its trap because it is blocked in
`wait`, the process-group kill targeting a group the stub is not in, or the case's `wait`
returning while the stub is orphaned — was **not** established. Establishing it is this REQ's
first job.

## Detailed Requirements

- Determine which side drops the signal, with evidence rather than inference. The three
  candidates above are a starting list, not a conclusion.
- An interrupted wrapper terminates its backend before exiting. This is the shipped behaviour
  and it is the half that matters to consumers.
- The case cannot hang. Whatever the outcome, a stuck backend must make the probe **fail
  with a diagnostic**, never wait forever — a test that can hang is worse than one that fails,
  because the gate's exit status is the only thing anyone reads.
- The fix holds under the load that surfaced it: both occurrences involved two gates running
  concurrently on the same machine.

## Constraints

- `_dev/primes/prime-shell-commands.md` governs any shell that ships. Read it first.
- Do not weaken the assertion the case exists to make: it checks that an interruption cleans
  the invocation-private staging file and leaves the previous target untouched. A fix that
  stops the hang by not exercising the interruption is not a fix.
- A timeout is a backstop, not the repair. If the only change is wrapping the case in
  `timeout`, the orphaned-backend defect is still shipped to consumers.

## Builder Guidance

**Certainty: Firm that it is broken, exploratory on the cause.** Two independent reproductions
with the same three live processes is not a flake. But the mechanism is genuinely unknown —
resist the first plausible story, and get evidence that distinguishes the candidates before
changing anything.

Scope cue: the wrapper's interruption path and the probe that exercises it. Not a rewrite of
either script.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-321: the canonical gate can
  hang indefinitely on `generate-report-image`'s interruption case, and the same shipped path
  would orphan the image backend after an interrupted `ai-report`. Should I process this as a
  new task? → Confirmed: Yes, add to queue
  - [2026-08-23] User approved via clarify: the hang blocks the canonical gate (two
    reproductions on 2026-08-22/23) and the shipped path leaves a stray process after an
    interrupted `ai-report`, so the fix is wanted as its own REQ. Nothing put out of scope.

## Red-Green Proof

**RED prompt/case:** Run `bash _dev/tests/prescribed-shell-cases/generate-report-image.sh`
and watch for a `bash .../image-interrupt-bin/imagegen` process outliving the case. Today it
survives, and the run does not return.

**Why RED now:** The TERM is sent to the wrapper; the stub never receives one and its
`trap "exit 143" TERM` never fires.

**GREEN when:** the case completes on its own with no `imagegen` process left behind, its
existing staging and old-target assertions still passing; and a deliberately unkillable
backend makes the case fail with a diagnostic naming what was still running, rather than
hanging.

**Validation:** Inferred during capture — a Discovered Task from REQ-321's work, not a user
request. The two reproductions are recorded in REQ-321's `## Discovered Tasks`.

## Assets

None. Reproduced by running the gate; the process list is recorded in REQ-321.

---
*Source: Discovered Task, REQ-321 (UR-065) — found twice while running the canonical gate.*
