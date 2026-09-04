---
id: REQ-340
title: "Addendum: finish the report-image interruption sweep"
status: completed
completed_at: 2026-08-24T09:38:00Z
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
kb_status: promoted
kb_entry: REQ-340-addendum-finish-the-report-image-interru.md
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
- [ ] **REOPENED — `terminate_backend_process`'s grace loop counts an unreaped zombie as alive**
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

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh`, `bash _dev/tests/prescribed-shell-scripts-behavior.sh`, `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Gate exit 0; `101 named script cases across 17 per-script files`

**Red-green validation:**
- New staging-leak case: ✗ against the pre-fix batch script (`FAIL: ai-report early-interrupted batch replay leaked invocation-private staging`) → ✓ after. The reviewer confirmed it fails for the right reason: exit 143 is also produced by the default TERM action, so the status assertion alone would pass vacuously — the leak assertion is what carries the case.
- Instance 2, stall-injected copies of both scripts: pre-fix `exit=143 backend_alive=yes` → post-fix `exit=143 backend_alive=no`

**Existing tests updated:** the interrupted case's bare `wait` gained a 10s watchdog (D-02) — no assertion changed.

*Verified by work action*

## Review

**Overall: 75%** | 2026-08-24T09:00:39Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 85% |
| Test Adequacy | 75% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings:**
- Instance 3's refutation measured the bare-PID branch while the code takes the process-group branch; the behaviour it denies reproduces on the shipped helper (teardown 1.03s vs the claimed 0.105s). — `impact-rule-change` → **REQ-356** created as a sweep.

**Minor findings:** 5 (report only) — the fixture's `mktemp` shim also intercepts the helper's own staging, contrary to its stated guarantee (latent: the case interrupts before any helper launches); a group-delivered interruption during staging creation still leaks, microseconds wide with the real `mktemp`; `RESTART-PROMPT.md:32` states a stale case count → prose backlog; the trap-installation triple is still duplicated at two sites plus the status literals at a third; the watchdog's `sleep 10` outlives its subshell.

**Acceptance:** Partial — instances 1 and 2 verified red→green by independent reproduction; instance 3's stated GREEN clause is measurably false.

**Verified and explicitly not findings:** the deferral cannot lose a signal (three signals inside the window → `exit=143 backend_alive=no staging_leak=no`); no early-return path exists between defer and resume; no instant has a default disposition; the PATH shim does not leak across cases; existing assertions unweakened; the watchdog is earned rather than drift — it converts a hang into a KILL whose 137 status the existing assertion reads as failure, the safe direction.

**Restatement sweep:** done. `prescribed-shell-primitives.md:100`'s contract is still true and better honored; no doc or comment states the old ordering or trap behaviour.

**Follow-ups created:** REQ-356

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Shimming `mktemp` on `PATH` to hold the batch inside the trap gap turned a signal
race into a deterministic event — the interruption became something the case *fires*, not something
it hopes to hit. Injecting a stall into a copy of the script (never the shipped file) proved instance
2 without shipping a fixture that could never fail reliably.

**What didn't:** Instance 3 was refuted from a reading of the code plus a measurement of the wrong
branch, and the orchestrator's independent check repeated the same mistake — a true measurement of
`kill -0 <pid>` when the code calls `kill -0 -- -<pgid>`. Two readers, same wrong answer, because
both reasoned about which branch mattered instead of timing the real path. The generalisable rule is
the one REQ-325's own Lessons already stated and this REQ re-proved: do not reason about bash signal
semantics from memory, and when a fix is declined on a measurement, measure the branch the code takes.

**Worth knowing:** A process group answers `kill -0` while any member holds a process-table slot,
including descendants nothing has reaped — which is why the group branch and the bare-PID branch give
opposite answers about the same "dead" backend. The batch's `report_image_batch_is_alive` carries the
identical shape.

## Orientation

An interrupted `do-work-toolbox ai-report` no longer leaves a staging directory in the user's report
folder, and no longer orphans a helper process and its image backend when the interruption lands
inside the launch window. Lives in `skills/do-work-toolbox/scripts/` — the report-image path. The
contract at `prescribed-shell-primitives.md:100` ("terminates, escalates, and reaps everything it
launched before staging is removed") is now true for interruptions that arrive early, not only late.
