---
id: REQ-356
title: "[impact-rule-change] Stop the group liveness probe counting exited members as alive"
status: completed
claimed_at: 2026-08-24T14:17:16Z
completed_at: 2026-08-24T15:14:05Z
commit: 502162f
status_changed_at: 2026-08-24T15:14:05Z
route: C
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
sweep: true
sweep_key: group-liveness-probe-counts-unreaped-members
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - skills/do-work-toolbox/scripts/generate-report-image.sh
  - skills/do-work-toolbox/scripts/generate-report-image-batch.sh
  - _dev/tests/prescribed-shell-cases/generate-report-image.sh
  - _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-08-24T14:46:37Z
  basis:
    - Route C
    - 4-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
---

# Stop the Group Liveness Probe Counting Exited Members as Alive

## What

`backend_process_is_alive` probes a process *group* (`kill -0 -- "-$pgid"`) whenever one was
recorded, which is every normal run. A group answers `kill -0` while any member still occupies a
process-table slot — including descendants bash never reaps — so the grace loop spins its full
10 x 0.1s budget and sends a redundant KILL after everything has already exited.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-shell-commands` plus testing and coding guardrails. Preserve group TERM/KILL signalling; use `ps` group-state output only to decide whether a recorded group has a non-zombie member, falling back to `kill -0 -- -pgid` if `ps` fails. Add descendant-bearing interruption fixtures for both helpers and retain a TERM-deaf 10-tick KILL path.
- [x] **[APPLY]:** Replaced only group liveness probes in the two declared scripts and added the paired prescribed-shell fixtures. Bare-PID probes and all signalling paths remain unchanged.
- [x] **[UNIFY]:** Reviewed the four declared code/test files and this REQ evidence; `git diff --check`, `bash -n` on all four files, and both focused case files pass. The test copies each real helper to an invocation-private path only to substitute a controlled `ps` executable, while the backend and descendant remain real processes. Against a detached pre-fix worktree, the helper fails with 3 grace ticks and the batch with 7; current code is green at 0. The TERM-deaf batch fixture now waits for its recorded descendant to become non-live after KILL. ShellCheck reports only pre-existing dynamic-trap/fixture warnings. No debug artifacts or unrelated paths were added.

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

- [x] **`backend_process_is_alive` (`generate-report-image.sh:62-68`) — the group branch counts
  exited-but-unreaped members as alive.** Every interrupted invocation with a descendant-bearing
  backend pays a full second of dead grace and a redundant KILL.
- [x] **`report_image_batch_is_alive` (`generate-report-image-batch.sh:117-128`) — the identical
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

## Triage

**Route: C** — The defect spans two nested process-group owners, timing-sensitive interruption behavior, descendant safety, and repository-wide shell gates. It requires a planned lifecycle fix and branch-accurate regression evidence.

## Plan

Retain group TERM/KILL delivery. Change only group liveness classification: inspect the verified process group's states with `ps`, treat zombie-only groups as finished, and conservatively fall back to `kill -0` if inspection fails. Pin both the prompt-exit path and TERM-deaf full-budget/KILL path in the helper and batch case files.

## Exploration

Both scripts record a private process group and later use `kill -0 -- -PGID` as their group branch. The bare-PID branch is not the affected path. The safe discriminator is whether `ps` reports any non-`Z` member in the recorded group; leader-only probing could orphan a TERM-deaf descendant.

## Scope

**Files I will touch:**

- `skills/do-work-toolbox/scripts/generate-report-image.sh`
- `skills/do-work-toolbox/scripts/generate-report-image-batch.sh`
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh`
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh`

**Acceptance criteria:** zombie-only groups stop consuming grace ticks; descendant-bearing fixtures exercise the actual group branch; TERM-deaf descendants receive the full ten-tick budget at each ownership level and are killed; bare-PID and group signaling behavior remain unchanged.

## Implementation Summary

- `skills/do-work-toolbox/scripts/generate-report-image.sh` (modified): classifies a recorded group as alive only when `ps` reports a non-zombie member, with conservative group-probe fallback.
- `skills/do-work-toolbox/scripts/generate-report-image-batch.sh` (modified): applies the same member-state rule across recorded helper groups while retaining nested group signals.
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` (modified): adds descendant-bearing prompt-exit and TERM-deaf full-budget/KILL cases.
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` (modified): adds descendant-bearing nested batch cases, including both ten-tick budgets for a TERM-deaf tree.

## Decisions

- **D-01 — DECIDE & STATE:** Use group-member state inspection rather than leader-only liveness. A leader may exit while a TERM-deaf descendant remains; the group must still reach KILL after its readable grace budget.
- **D-02 — DECIDE & STATE:** On `ps` failure, preserve the existing `kill -0` group result. A status-tool failure must delay cleanup rather than risk leaving a live descendant behind.

## Testing

- RED was first reproduced on the real group branch with descendant-bearing processes. The initial timing fixtures exposed platform-dependent reaping, so they were replaced with invocation-private copies of the real helpers whose only substitution is a controlled `ps` command.
- The deterministic fixtures fail against a detached pre-fix checkout on the named regression: the helper consumes 3 grace ticks and the batch 7 when the controlled group is zombie-only. Current code consumes 0 helper ticks and at most 1 batch tick.
- The TERM-deaf helper and batch cases retain real live descendants, consume the full 10-tick helper budget, receive group KILL, and explicitly poll that the recorded descendant is no longer live.
- Focused prescribed-shell results: `generate-report-image.sh` 11/11 and `generate-report-image-batch.sh` 6/6.
- `bash -n` passed for both production scripts and both case files; `git diff --check` passed.
- `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh` passed, including 104 prescribed-shell behavior cases, strict JavaScript behavior, Go vet, and uncached tests. The strict browser lane reported its documented skip because no browser was available to that gate invocation; this shell-only change has no browser acceptance surface.

## Qualification

- `qualify.sh`: PASS; remaining judgment checks were reviewed by the orchestrator.
- `scope-drift.sh`: PASS; the Implementation Summary matches the four declared files.
- Requirements trace: both liveness probes use the same non-zombie group-member rule; the fallback and signalling branches are unchanged; both prompt-exit and TERM-deaf behavior are pinned.

## Review

The first independent review correctly rejected the timing-only fixtures because they could false-green on a pre-fix platform and the batch case did not prove the TERM-deaf descendant was gone. One remediation replaced timing variance with a controlled observation seam and added the missing survivor assertion. Re-review approved with no findings: Requirements 98%, Code Quality 94%, Test Adequacy 94%, Scope 100%, overall 96%, low risk, acceptance pass.

## Lessons Learned

A lifecycle regression test cannot rely on how quickly one kernel happens to reap a process. Keep the process tree real, but control the observation boundary that defines the branch under test; then assert both the timeout count and the descendant's terminal state.

## Orientation

Released in 0.236.41. The report-image helper and batch wrapper now stop grace promptly for zombie-only groups while retaining the full timeout and group-wide KILL for a genuinely live TERM-deaf descendant.
