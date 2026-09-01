---
id: REQ-325
title: "Stop the report-image interruption path orphaning its backend"
status: completed
claimed_at: 2026-08-23T18:57:07Z
completed_at: 2026-08-23T19:29:47Z
commit: 92413b9
kb_status: promoted
kb_entry: REQ-325-stop-the-report-image-interruption-path-.md
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
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-23T18:57:07Z
  basis:
    - Route B
    - 2-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - async lifecycle behavior
    - full-suite verification
route: B
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
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` and `_dev/lessons/validated-runtime-boundaries.md`
  (§ *A timeout owns the process tree it starts*), plus `general.md`, `coding-guardrails.md`,
  `communication-style.md`, `testing.md`.
  1. **Wrapper — close the publish-the-PID window (E4).** A trapped signal runs between commands,
     so a TERM landing between `imagegen … &` and `backend_process_id=$!` reaches cleanup with no
     handle. Defer HUP/INT/TERM across both launch windows (record the status instead of exiting),
     restore the direct traps once the PID and process group are recorded, then exit through the
     normal path. → verify: the wrapper's stall-injected probe stops orphaning; the five existing
     cases still pass; a TERM that lands *before* any backend is launched still exits 143 with the
     staging file gone (new case — the deferral must not swallow that).
  2. **Case — no unbounded wait.** Both interruption cases do `kill -TERM` then a bare
     `wait "$helper_pid"`. Proved RED: that shape blocks forever against a child that does not
     finish on TERM (`timeout -s KILL 8` had to kill it), which is how a stuck backend wedges the
     gate. Replace with a deadline that names the processes still alive at expiry, then kills them.
     → verify: a self-check case drives the deadline against a deliberately TERM-deaf stand-in and
     asserts the diagnostic instead of hanging.
  3. Not doing: bounding the wrapper's own closing `wait`. E3 shows its KILL escalation reaps even
     a TERM-ignoring backend, so there is no reachable hang to earn that surface.
- [x] **[APPLY]:** Both declared files only. Wrapper: `deferred_interrupt_status` state plus
  `defer_interrupts_across_backend_launch` / `resume_interrupts_after_backend_launch`, wrapped
  around both launch sites; the staging `mktemp` moved below the trap block. Case file:
  `wait_for_wrapper_or_fail` deadline helper, adopted by both existing interruption cases, plus
  the two new cases the plan named.
- [x] **[UNIFY]:** `git diff --stat` → 2 files, +147/-5, both inside the declared write set.
  - `skills/do-work-toolbox/scripts/generate-report-image.sh` — read the whole diff: 4 hunks, every
    line traceable to the launch-window fix or the staging-file move. `bash -n` clean.
    ShellCheck `--severity=warning` clean (the gate's own 73-file lint lane passed).
  - `_dev/tests/prescribed-shell-cases/generate-report-image.sh` — read the whole diff: the helper,
    two call-site swaps, two appended cases. `bash -n` clean, ShellCheck clean.
  - No debug artifacts: the only `grep -E 'console\.log|debugger|TODO|XXX|FIXME'` hit in the diff is
    `mktemp`'s `XXXXXX` template. The stall-injection probes live in the scratchpad, not the repo.
  - Native lint/gate: `bash _dev/tests/maintainer-verify.sh` → exit 0.

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

---

## Triage

**Route: B** - Medium

**Reasoning:** The files are named and the failure is reproduced twice, but the REQ's own first requirement is to establish *which* side drops the signal with evidence — the cause is genuinely unknown. That is exploration, not planning: no architecture decision is pending, only a diagnosis inside two known scripts.

**Planning:** Not required

## Exploration

**Files that matter**

- `skills/do-work-toolbox/scripts/generate-report-image.sh` — the shipped wrapper. Signal path:
  `trap 'exit 143' TERM` (:106) → `trap cleanup_report_image_paths EXIT` (:103) →
  `terminate_backend_process` (:78-90) → `signal_backend_process` (:69-76), which signals the
  backend's process group when `record_backend_process_group` (:48-59) could prove one, and the
  bare PID otherwise.
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` — five fixture cases. Two exercise
  interruption: the staging/old-target case (:71-98, the one the REQ names) and the process-tree
  case (:131-193, which additionally records the backend's descendant and asserts both die).
- `_dev/tests/prescribed-shell-harness.sh` — per-case-file fixture root, `background_process_ids`
  reaping, `fail_case` tally. Each case file is its own process (`prescribed-shell-scripts-behavior.sh:18`).
- Callers: the suite reaches this case only through `contract-regressions.sh:5092` →
  `staged-skills-contract.sh:183` → the runner. No caller captures the runner's output in a
  command substitution, so a fd-holding orphan is not part of the chain.

**Evidence gathered (the REQ's first job)**

- **E1 — the case passes and does not hang on this machine.** One run: `5 cases, 0 failures`, no
  surviving `imagegen`. Three rounds of two concurrent runs: all six exit 0, no survivors.
- **E2 — TERM arriving while the wrapper is inside `record_backend_process_group` is handled.**
  A wrapper copy stalled 0.3s/1s/2s inside that function still exits 143 with no survivor:
  bash defers the trap until the foreground command returns, it does not lose it.
- **E3 — the stub does receive TERM on the normal path.** With a stub that records its trap
  firing: `stub_got_term=yes`, no survivor. With a stub that ignores TERM (`trap "" TERM`), the
  wrapper's KILL escalation still reaps it: `stub_alive=no`.
- **E4 (the defect) — the wrapper drops the signal in the window between forking the backend and
  publishing its PID.** A wrapper copy stalled between `imagegen … &` and
  `backend_process_id=$!` orphans the backend, and *every existing assertion still passes*:

  ```
  stall=0    wrapper_status=143 ORPHANED_BACKEND=no  stage_leaked=no old_target=stable
  stall=0.5  wrapper_status=143 ORPHANED_BACKEND=yes stage_leaked=no old_target=stable
  stall=2    wrapper_status=143 ORPHANED_BACKEND=yes stage_leaked=no old_target=stable
  ```

  Mechanism: `trap 'exit 143' TERM` is armed before the launch, so a TERM landing between the
  `&` and the assignment runs `exit 143` → the EXIT trap → `terminate_backend_process`, which
  reads an *empty* `backend_process_id` and returns 0 at its first line. Nothing is ever
  signalled. The window is one shell assignment wide, which is why it needs the parent to be
  descheduled right after a fork — exactly the condition two concurrent gates on four cores
  produce. The same window exists verbatim in the agentic branch (`:151-152`).

**Which side drops the signal:** the wrapper. Both other candidates the REQ listed are refuted —
the wrapper's trap does fire from inside `wait` (E2/E3), and the process-group kill does reach the
stub (E3).

**Not established:** the 35-minute hang itself. It did not reproduce here (E1), and E4 produces an
orphan with a *prompt* exit, not a stall, so it does not by itself explain a wedged gate. What is
structurally true is that nothing on this path is bounded: `terminate_backend_process`'s closing
`wait` (:89) and the case's `wait "$interrupt_helper_pid"` (:92) both block indefinitely, so any
backend that outlives its signals wedges first the shipped script and then the gate, with no
diagnostic. That is the REQ's third requirement and it is fixed on its own evidence rather than on
a story about the observed stall.

**Secondary finding:** the grace loop (:82-85) treats an unreaped zombie as alive, so the normal
interruption path always burns its full 1s budget and then sends a pointless KILL.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/scripts/generate-report-image.sh` (modify) — close the publish-the-PID
  window on both backend launches; bound the closing reap and make an unreapable backend print a
  diagnostic naming it instead of blocking forever.
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` (modify) — lock in the orphan
  (E4), lock in the diagnostic-instead-of-hang behaviour, and bound both interruption cases'
  `wait` calls so the probe can fail but never wedge.

**Files I will NOT touch:** `_dev/tests/prescribed-shell-harness.sh` (its per-file reaping already
covers the recorded PIDs), `_dev/tests/prescribed-shell-scripts-behavior.sh` (the runner counts
cases from the case files; no runner change needed), `skills/do-work-toolbox/actions/ai-report.md`.

**Acceptance criteria (restated from REQ):**
- [ ] Which side drops the signal is established with evidence, not inference.
- [ ] An interrupted wrapper terminates its backend before exiting — including a TERM that lands
      in the launch window.
- [ ] The case cannot hang: a stuck backend fails the probe with a diagnostic naming what was
      still running.
- [ ] The existing assertions are not weakened — the interruption is still exercised, the
      invocation-private staging file is still checked gone, and the old target still untouched.
- [ ] The fix holds with two gates running concurrently.

## Decisions

- **D-01 — DECIDE & STATE. The wrapper's own closing `wait` was left unbounded.** Requirement 3
  says a stuck backend must not wait forever, and `terminate_backend_process`'s `wait` (:90) is
  unbounded. It was left alone because E3 shows the KILL escalation reaps even a backend that
  ignores TERM, and `wait` only ever blocks on a live direct child — so there is no reachable hang
  to earn that surface (`coding-guardrails.md` § Simplicity First, *Earned defense*). The bound
  went where the hang actually is: the case's `wait` on the wrapper.
- **D-02 — DECIDE & STATE. The launch-window fix ships without a fixture that can fail.** The
  window is one shell assignment wide and the parent always won it here — 100 rounds with the stub
  TERMing its own parent, and 60 of those under six spinning load generators, produced zero
  orphans. It is only observable with a stall injected into a copy of the wrapper (E4). Shipping a
  stress case that has never failed would be the "lock-in that cannot fail" this repo's prime
  warns about, so the live coverage is the early-interruption case plus the five pre-existing ones;
  E4's reproduction command is recorded in `## Exploration` instead.
- **D-03 — DECIDE & STATE. The staging-file move went in even though the REQ did not name it.**
  The new early-interruption case failed on its first run with `leaked private staging`: `mktemp`
  ran ~80 lines above the EXIT trap, so any interruption in between took the default action and
  left the invocation-private file behind. That is the same interruption path and the very
  assertion the existing case exists to make, so it is this REQ's bug rather than a discovered
  task. Fix is a two-line reorder, not new machinery.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/scripts/generate-report-image.sh` (modified)
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` (modified)

**What was done:** Closed the wrapper's publish-the-PID window by deferring HUP/INT/TERM across
both backend launches — the status is recorded, the PID and process group are registered, then the
interruption is re-raised so the EXIT trap has a backend to terminate — and moved the staging
`mktemp` below the trap block so an interruption can never leave the invocation-private file
behind. Replaced both interruption cases' unbounded `wait` with `wait_for_wrapper_or_fail`, a
10-second deadline that names the processes still alive at expiry, fails the case, and kills them.
Added two cases: an interruption fired the moment the staging file appears, and a deadline probe
run against a deliberately TERM-deaf stand-in that asserts the diagnostic rather than hanging.

## Discovered Tasks

- **The case-count regex in `prescribed_shell_finish` misses two of this file's cases.** It counts
  `^# [a-z0-9][a-z0-9-]*: `, so `# generate-report-image caller contract: …` and
  `# generate-report-image, interrupted directly: …` are invisible — a space and a comma before the
  colon. The file reports 7 cases and contains 9. REQ-234 replaced a hand-maintained literal with a
  derived count precisely so the figure would stop being a remembered number; an undercount is the
  same untruth with a different cause. Worth a repo-wide check — other case files may have headers
  in the same shape.
- **`terminate_backend_process`'s grace loop counts an unreaped zombie as alive.** `kill -0` on a
  child that has exited but not been `wait`ed succeeds, so the normal interruption path spins the
  full 10 × 0.1s budget and then sends a redundant KILL before reaping. Costs a second on every
  interrupted invocation and makes the grace budget unreadable as a real timeout.
- **The same two windows exist in `generate-report-image-batch.sh`, which drives this helper.**
  The prime's rule is to grep a fixed primitive across every caller before calling the class
  closed, and both instances are there. (1) Its interruption traps are installed at `:131-133`
  while its staging directory is created at `:36` — a TERM anywhere in between takes the default
  action, so the EXIT trap never runs and `.generated.staging.*` is left in the report directory
  (the batch's twin of D-03). (2) `launch_report_image` (`:64-69`) has the same publish-the-PID
  window this REQ closed, one array append wider: `"$report_image_helper" … &` then
  `image_helper_pid=$!` then the `image_generation_pids+=()` append, with `terminate_report_image_batch`
  reading only the array. Out of this REQ's write set.


## Qualification

Passed — 2 files verified, 5 requirements traced, P-A-U confirmed.

- Mechanical: `tools/checks/qualify.sh` exit 0; `tools/checks/scope-drift.sh` exit 0 (Implementation
  Summary matches the Scope declaration exactly — no undeclared touch, no unused declaration).
- Substantive (check 2): both files are modifications, +147/-5; read the full diff of each. No
  placeholder bodies, no whitespace-only or import-shuffle hunks.
- Requirements traced (check 3): the diagnosis to `## Exploration` E1-E4 (the wrapper drops it, both
  other candidates refuted); "terminates its backend before exiting" to the deferral around both
  launch sites, red→green under an injected stall; "cannot hang / diagnostic" to
  `wait_for_wrapper_or_fail` and the deadline case; "assertion not weakened" to the untouched
  staging/old-target assertions in both interruption cases; "holds with two gates concurrently" to
  three rounds of two concurrent case-file runs, all 7-cases-0-failures with no survivors.
- Data flows (check 6): nothing stubbed. The deferral re-raises a real status, and the deadline
  helper's diagnostic is read from a file the watchdog writes before it kills anything, so the
  named PIDs are the ones that were actually alive rather than a post-kill re-scan.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-cases/generate-report-image.sh` (focused, ×3, plus
3 rounds of two concurrent runs) and `bash _dev/tests/maintainer-verify.sh` (the project's declared
canonical repository gate, from the repo root, against the final tree)
**Result:** ✓ focused `generate-report-image: 7 cases, 0 failures.` every run, no surviving
`imagegen` or wrapper process after any of them; canonical gate exit 0.

**Red-green validation:**
- early-interruption case (`_dev/tests/prescribed-shell-cases/generate-report-image.sh`):
  ✗ `FAIL: generate-report-image early-interruption case leaked private staging` on its first run
  against the pre-move wrapper → ✓ after the staging `mktemp` moved below the trap block. This is
  the RED that found D-03.
- unbounded-wait shape (the captured RED for the deadline): `kill -TERM` then a bare
  `wait "$pid"` on a TERM-deaf child ✗ never returned — `timeout -s KILL 8` had to kill it (outer
  status 137) → ✓ `wait_for_wrapper_or_fail` returns and prints
  `… did not finish within 1s; still alive: <pid>`, which the new deadline case asserts.
- launch-window orphan (the REQ's `## Red-Green Proof` mechanism, exercised against a stall-injected
  copy of the wrapper because the fixture cannot force a one-assignment window — D-02):
  ✗ `ORPHANED_BACKEND=yes` at stalls of 0.5s and 2s before the fix → ✓ `ORPHANED_BACKEND=no` at the
  same stalls after, with `wrapper_status=143`, `stage_leaked=no`, `old_target=stable` throughout.

**Captured proof, traced:** the REQ's GREEN asked for the case to complete with no `imagegen` left
behind (✓, checked after every run) and for a deliberately unkillable backend to make the case fail
with a diagnostic naming what was still running rather than hanging. The second clause is met by
the deadline case, adapted in one way: the stand-in is TERM-deaf rather than unkillable, because a
direct child can always be reaped by KILL, so a genuinely unkillable backend is not constructible
in a fixture. The behaviour under test — deadline expires, survivors named, probe fails instead of
blocking — is the one the clause asks for.

**New tests added:**
- `generate-report-image: an interruption arriving as early as the invocation is observable …`
- `generate-report-image: every interruption case above waits with a deadline …`

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` (interruption case, from REQ-220's
  ownership work; process-tree case, from REQ-204/REQ-220): both now wait through
  `wait_for_wrapper_or_fail` instead of a bare `wait`. Their assertions are unchanged — this
  replaces only the blocking primitive.

*Verified by work action*

## Review

**Overall: 96%** | 2026-08-23T19:26:39Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 1 (report only) — the launch-window fix ships without a live fixture that can fail against the unfixed wrapper (D-02 states why; a stress case that has never failed would be worse). The early-interruption case exercises the surrounding path, so the fix is not uncovered, only not pinned at its exact window.

**Acceptance:** Pass — focused case green on 3 solo runs and 3 concurrent pairs with no surviving processes; canonical gate exit 0; the pre-fix defect reproduces on demand and is gone after.
**Suggested testing:** 3 items
**Follow-ups created:** None; **sweeps appended to:** None

### Requirements Checklist

- [x] Which side drops the signal, established with evidence — delivered (E1-E4: the wrapper; the deferred-trap and process-group-kill candidates both refuted by probe)
- [x] An interrupted wrapper terminates its backend before exiting — delivered (deferral at both launch sites; `ORPHANED_BACKEND=yes` → `no` under an injected stall)
- [x] The case cannot hang; a stuck backend fails it with a diagnostic — delivered (`wait_for_wrapper_or_fail`, 10s deadline, survivors named from a pre-kill snapshot)
- [x] The interruption assertions are not weakened — delivered (both existing cases keep their staging and old-target checks verbatim; only the blocking primitive changed)
- [x] The fix holds with two gates running concurrently — delivered (3 rounds × 2 concurrent runs, all green)
- [x] A timeout is a backstop, not the repair — delivered (the repair is the deferral plus the staging-file move; the deadline is declared as the backstop)

### Restatement Sweep

The diff changes *when* the staging file is created and *when* an interruption is acted on — both facts other text states. Swept every consumer:

- `skills/do-work/docs/prescribed-shell-primitives.md:18` ("launched-process-tree ownership, verified exact invocation-private publication") and `:100` ("an interrupted batch terminates, escalates, and reaps everything it launched *before* staging is removed") — both still true, and the change moves the per-image helper further into compliance rather than away from it. No edit needed.
- `skills/do-work-toolbox/actions/ai-report-reference.md:34,57` — describe the helper's arguments and absolute-path contract, untouched by this diff.
- `_dev/lessons/validated-runtime-boundaries.md` § *A timeout owns the process tree it starts* — unchanged and now better honored.
- `decisions/audits/2026-08-11-defensive-surface.md:21` names `_dev/tests/prescribed-shell-scripts-behavior.sh` as this script's covering probe. That path is now a 35-line runner and the cases live in `prescribed-shell-cases/` (a REQ-300-class stale pointer in a dated audit record). Left alone deliberately: dated audit history is a record of what was true then, not a live restatement.

No stale restatement found.

### Code Review Notes

- **Naming for reach:** three new identifiers with reach — `defer_interrupts_across_backend_launch`, `resume_interrupts_after_backend_launch`, `deferred_interrupt_status` in the wrapper, `wait_for_wrapper_or_fail` and `wrapper_wait_status` in the case file. All multi-word and greppable.
- **Earned defense:** two helpers of four lines each, both called twice; the deadline helper replaces two copies of a blocking primitive rather than adding a layer. No speculative flags, no configurability beyond the deadline seconds the two call sites actually pass differently (10 and 1).
- **Surgical:** every hunk traces to the REQ. `terminate_backend_process`, `record_backend_process_group`, `publish_report_image`, and all five pre-existing cases are untouched.
- **Correctness of the deadline helper:** the timeout verdict reads the watchdog's report file, not the wait status — a wrapper the watchdog killed also exits nonzero, so status alone could not distinguish stuck from interrupted. Survivors are captured before the kill, so the diagnostic names what was actually alive.
- **`set -u` safety:** `deferred_interrupt_status` is initialised with the other state variables, so the EXIT trap can read it even on the earliest failure path.

### Acceptance Testing

**Result: Pass**
- Focused case file: 3 solo runs and 3 rounds of two concurrent runs, all `7 cases, 0 failures`, `ps` clean after each.
- Canonical gate from the repo root against the final tree: exit 0, `generate-report-image: 7 cases, 0 failures` inside it, and the lane that used to stall completed normally.
- Original defect verified gone: the stall-injected reproduction orphans a backend before the change and does not after.
- Adjacent consumer exercised: `generate-report-image-batch: 2 cases, 0 failures` — the batch drives this helper, and its cases pass unchanged.

### Suggested Additional Testing

- Environment: re-run the focused case on macOS. The deferral relies on bash running a trap between commands and on `set -m` group creation; bash 3.2 is the platform this repo has been bitten on before (REQ-216).
- Edge case: a real `imagegen` interrupted mid-write, to confirm the moved `mktemp` did not change what a consumer sees when the backend is genuinely producing bytes.
- Concurrency at higher width than two gates — the reported stall involved two, and two is what was reproduced against.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Injecting a stall into a *copy* of the script under test, one candidate window at a time. Three
  candidate mechanisms were on the table and reasoning could not separate them; a `sleep` at each
  suspected point separated them in minutes and refuted two of the three.
- Writing the new fixture case before reading the wrapper for a second defect. The
  early-interruption case failed on its first run for a reason nobody had predicted (the staging
  file predating its own cleanup trap) — the case found a bug the analysis had walked past twice.

**What didn't:**
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

**Worth knowing:**
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

## Orientation

Interrupting a report-image generation now always takes its backend with it, and the probe that
proves it can no longer wedge the gate. Lives in the do-work-toolbox report-image helper and its
prescribed-shell case file — the interruption path only; backend selection, publication, and the
agentic opt-in are unchanged. No new module, no contract change, so the system's shape is the same;
`_dev/primes/prime-shell-commands.md` gains one lesson link and its referenced paths all still
resolve.
