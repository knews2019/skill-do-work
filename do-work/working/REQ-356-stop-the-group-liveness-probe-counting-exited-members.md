---
id: REQ-356
title: "[impact-rule-change] Stop the group liveness probe counting exited members as alive"
status: claimed
claimed_at: 2026-08-24T13:05:00Z
created_at: 2026-08-24T09:35:00Z
user_request: UR-065
addendum_to: REQ-340
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
route: C
sweep: true
sweep_key: group-liveness-probe-counts-unreaped-members
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - skills/do-work-toolbox/scripts/generate-report-image.sh
  - skills/do-work-toolbox/scripts/generate-report-image-batch.sh
  - _dev/tests/prescribed-shell-cases/generate-report-image.sh
  - _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh
---

# Stop the Group Liveness Probe Counting Exited Members as Alive

## What

`backend_process_is_alive` probes a process *group* (`kill -0 -- "-$pgid"`) whenever one was
recorded, which is every normal run. A group answers `kill -0` while any member still occupies a
process-table slot — including descendants bash never reaps — so the grace loop spins its full
10 x 0.1s budget and sends a redundant KILL after everything has already exited.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

This reopens REQ-340's instance 3, which was archived as REFUTED on a measurement of the wrong
branch. The refutation measured `kill -0 <pid>` — true, and irrelevant: bash reaps its direct
background child asynchronously, so the bare-PID branch never counts a zombie. But
`record_backend_process_group` records a group whenever `ps` resolves and `set -m` gave the backend
its own group, which is the default path, and `backend_process_is_alive` then takes the group branch.

Measured on this machine, a backend that forks one ordinary child and obeys TERM:

```
GROUP branch : grace_ticks=10 elapsed=1.029s
BARE-PID     : grace_ticks=0  elapsed=0.003s
```

The 105ms figure the refutation rested on only reproduces for a backend with **no descendants at
all**, which a real `imagegen` or `codex exec` is not.

## Instances

- [ ] **`backend_process_is_alive` (`generate-report-image.sh:62-68`) — the group branch counts
  exited-but-unreaped members as alive.** Every interrupted invocation with a descendant-bearing
  backend pays a full second of dead grace and a redundant KILL.
- [ ] **`report_image_batch_is_alive` (`generate-report-image-batch.sh:117-128`) — the identical
  group-branch probe.** An interrupted batch pays it too, once per level.

## Detailed Requirements

- The grace loop stops as soon as everything it is waiting for has genuinely exited, on the group
  branch as well as the bare-PID one.
- Both instances close, or one is refuted with a measurement **of the branch the code takes** — the
  failure this REQ exists to correct.
- A fixture pins the timing behaviour, so the claim is falsifiable next time rather than re-argued
  from reading the code. It must use a backend with at least one descendant; a childless stub cannot
  reproduce the defect and would pass vacuously.
- The grace budget stays readable as a real timeout: a TERM-deaf backend must still get its full
  10 x 0.1s and then the KILL.

## Constraints

- `_dev/primes/prime-shell-commands.md` governs. Read it first.
- The group handle exists for a reason — it is the only thing that reaches descendants
  (`generate-report-image.sh:45-47`). Do not solve this by dropping to bare-PID signalling; TERM and
  KILL must still reach the whole group.
- Do not weaken the existing interruption assertions in either case file.

## Builder Guidance

**Certainty: firm on the defect, open on the fix.** Three shapes were named by the review: probe the
group leader rather than the group, reap direct children before probing, or shorten the budget. Each
has a different failure mode against a TERM-deaf descendant — choose deliberately and say why.

**Measure before and after, on the branch the code takes.** This REQ exists because a reading of the
code produced a confident wrong answer twice: once by the builder, once by the orchestrator
verifying it. Time the real path; do not reason about it.

## Red-Green Proof

**RED prompt/case:** Launch the shipped helper with a TERM-obedient stub backend that forks one
ordinary child, interrupt it, and time TERM-to-exit: approximately 1.03s, with the grace loop
reaching all 10 ticks. Reproduced 2026-08-24.

**GREEN when:** the same measurement returns promptly once the group's members have exited, a
TERM-deaf backend still receives the full budget then the KILL, and a fixture fails if the probe
regresses to counting exited members.

**Validation:** Inferred during REQ-340's review, then independently re-measured by the orchestrator.

---
*Source: REQ-340 review finding F1 (UR-065) — reopens REQ-340 instance 3.*
