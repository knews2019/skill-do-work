---
id: REQ-340
title: "Addendum: finish the report-image interruption sweep"
status: claimed
claimed_at: 2026-08-23T23:08:00Z
status_changed_at: 2026-08-23T22:32:23Z
created_at: 2026-08-23T19:30:00Z
user_request: UR-065
addendum_to: REQ-325
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
route: B
depends_on: []
maintenance: false
sweep: true
sweep_key: report-image-interruption-ownership
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-toolbox/scripts/generate-report-image-batch.sh
  - skills/do-work-toolbox/scripts/generate-report-image.sh
  - _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh
  - _dev/tests/prescribed-shell-cases/generate-report-image.sh
---

# Finish the Report-Image Interruption Sweep

## What

REQ-325 closed two interruption defects in the per-image helper. The prime's rule is to grep the
same primitive across every caller before calling the class closed, and the batch that drives that
helper carries both of them. A third, smaller instance sits in the helper REQ-325 did fix.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-shell-commands.md` and the archived REQ-325. Reused REQ-325's two fixes (trap-before-staging reorder; defer-and-re-raise across the publish window) and measured instance 3 before committing to it.
- [x] **[APPLY]:** Two files changed, both inside the declared write set. `generate-report-image.sh` and its case file were declared but deliberately left untouched — instance 3 is refuted, so there was nothing to fix.
- [x] **[UNIFY]:** Audited by the orchestrator against the merged range `dceac86..da5e605`: confirmed `mktemp -d` (line 169) now sits below the trap installations (161-163), the defer/resume pair brackets the launch (95, 114), and the zombie refutation reproduces independently (a bash background child stops answering `kill -0`; a genuine perl-forked zombie answers).

## Why

An interrupted `do-work-toolbox ai-report` is a real consumer path. Two of the three instances leave
something behind on it — a staging directory in the user's report folder, or a helper process and its
image backend still running — and `skills/do-work/docs/prescribed-shell-primitives.md:100` states as
contract that an interrupted batch "terminates, escalates, and reaps everything it launched *before*
staging is removed". Today that statement is true only for an interruption that arrives late enough.

## Instances

- [x] **CLOSED — `generate-report-image-batch.sh` creates its staging directory before any interruption trap
  exists.** `mktemp -d` at `:36`, HUP/INT/TERM traps at `:131-133`. A signal in between takes the
  default action, so the EXIT trap never runs and `.generated.staging.*` is left in the report
  directory. This is the batch's twin of REQ-325's D-03; the fix there was a two-line reorder.
  Judged `impact-user-visible`.
- [x] **CLOSED — `launch_report_image` has the publish-the-PID window REQ-325 closed, one array append
  wider.** `"$report_image_helper" … &` (`:67`), then `image_helper_pid=$!`, then the
  `image_generation_pids+=()` append — and `terminate_report_image_batch` reads only the array, so an
  interruption inside that window reaps nothing. REQ-325's fix was to defer HUP/INT/TERM across the
  window and re-raise after publication. Judged `impact-user-visible`.
- [~] **REFUTED — `terminate_backend_process`'s grace loop counts an unreaped zombie as alive**
  (`generate-report-image.sh:83-89`). `kill -0` on a child that has exited but has not been `wait`ed
  succeeds, so every interrupted invocation spins the full 10 × 0.1s budget and then sends a
  redundant KILL. Costs a second each time and makes the grace budget unreadable as a real timeout.
  Judged `impact-negligible`.

## Detailed Requirements

- Each instance is closed or explicitly refuted with evidence, and the checklist above records which.
- The batch's own fixture file gains cases for whichever instances land there — an interrupted batch
  must leave no staging directory and no launched helper or backend alive.
- Existing batch assertions are not weakened: the two current cases cover parallel launch, retained
  per-image statuses, and all-or-nothing directory publication.

## Constraints

- `_dev/primes/prime-shell-commands.md` governs any shell that ships. Read it first.
- Read REQ-325's `## Exploration`, `## Decisions` and `## Lessons Learned` before starting: it
  records which bash semantics were confirmed by probe (a trap does fire from inside `wait`; a group
  kill does reach the backend; a direct child is always reapable by KILL) and which fix shapes
  worked. Re-deriving those cost most of that REQ's time.
- The launch-window instance could not be pinned by a fixture in REQ-325 — the parent won the race
  160/160 times. Do not ship a stress case that has never failed; that REQ's D-02 states why.

## Open Questions

- [x] I discovered these out-of-scope tasks while working on REQ-325: the report-image *batch*
  carries both interruption defects REQ-325 fixed in the per-image helper, and the helper's grace
  loop miscounts a zombie child so every interrupted run burns an extra second. Should I process
  this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  - [2026-08-23] User approved via clarify: all three instances are wanted, including the
    helper's zombie-counting grace loop. The user was offered a consumer-visible-only subset
    (the batch's leaked staging directory and stray processes) and did not take it, so the
    third instance stays in scope.

## Red-Green Proof

**RED prompt/case:** TERM `generate-report-image-batch.sh` while it is between its `mktemp -d` and
its interruption traps — the `.generated.staging.*` directory survives in the report folder. Then,
with a stall injected into a copy between the helper launch and the array append, TERM the batch —
the helper and its backend outlive it.

**GREEN when:** neither reproduces, the batch's existing publication assertions still pass, and the
grace loop no longer spends its full budget on an already-exited child.

**Validation:** Inferred during REQ-325's implementation and review — Discovered Tasks, not a user
request.

---
*Source: Discovered Tasks, REQ-325 (UR-065) — the same primitive, swept across its caller.*

---

## Triage

**Route: B** - Medium

**Reasoning:** A sweep with three named instances at named line numbers, but each needed its own reproduction before a fix could be justified, and REQ-325's archived fixes had to be located and reused. Clear what, discovery needed on how each reproduces.

**Planning:** Not required.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/scripts/generate-report-image-batch.sh` (modify) — instances 1 and 2
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` (modify) — fixture case for instance 1, watchdog on the existing interrupted case

**Files I will NOT touch:** `skills/do-work-toolbox/scripts/generate-report-image.sh` and its case file — declared in the capture-seeded write set for instance 3, which is refuted, so nothing there is broken.

**Acceptance criteria (restated from REQ):**
- [x] Each instance closed or explicitly refuted with evidence
- [x] The batch's fixture gains cases for the instances that land there
- [x] Existing batch assertions not weakened (parallel launch, retained per-image statuses, all-or-nothing publication)

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/scripts/generate-report-image-batch.sh` (modified)
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` (modified)

**What was done:** Moved the batch's `mktemp -d` staging creation below its HUP/INT/TERM trap installations, closing the window where a signal took the default action and left `.generated.staging.*` in the user's report directory. Added a defer/re-raise pair around `launch_report_image`'s fork so an interruption between the background launch and the third array append can no longer orphan a helper and its backend. The three trap installations now call one named `interrupt_report_image_batch <status>` so the deferral re-raises exactly what the traps do. Instance 3 was measured and refuted: no change made.

## Verification of the Refutation

The orchestrator re-derived instance 3 independently rather than accepting the builder's measurement:

```
bash child: kill -0 fails (builder right)
genuine zombie 26710: kill -0 SUCCEEDS (premise's general claim holds)
```

The premise's general claim about zombies is true; its antecedent is not, because bash reaps a background child asynchronously in its SIGCHLD handler. The grace loop therefore exits on the first tick (~105ms measured end-to-end), not the tenth, and no interrupted invocation costs an extra second. Declining to ship a change here is the correct outcome.
