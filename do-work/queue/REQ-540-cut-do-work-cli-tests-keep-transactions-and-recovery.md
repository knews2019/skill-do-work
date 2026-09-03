---
id: REQ-540
title: 'Cut do-work-cli tests: keep transactions and recovery, short-circuit the rest'
status: pending
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: testing
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
write_set:
  - skills/do-work/tools/do-work-cli/**/*_test.go
  - _dev/tests/maintainer-verify.sh
---

# Cut do-work-cli Tests: Keep Transactions and Recovery, Short-Circuit the Rest

## What

Keep the transaction and recovery suites (finalization, publication, gittransaction, requeststate, gateevidence); they are the product. Delete matrix tests that enumerate cases a table-driven test already proves once. Put the binary-building and signal tests behind `testing.Short()`, and have REQ-537's fast tier pass `-short` to this module.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

- 25 packages, 542 tests, 114 s CPU, 35 s wall. Slowest: `corehelpers.TestInventoryMatchesRetainedPorcelainXYMatrix` 12.6 s, `publication.TestDeferGateRollsBackUntrackedCreateAndFoldTopologies` 8.1 s, `finalization.TestRecoverFinalizationResumesEveryDurablePhaseExactlyOnce` 6.7 s, `corehelpers.TestAllSeventeenPublicCommandsRunInTextAndJSONWithStableStatusAndNoDryRunEffects` 2.5 s (builds the binary), `suiteinstall.TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation` 2.2 s.
- `finalization` is the wall-clock ceiling at 34.8 s: 39 fresh git repositories and three `go run` invocations inside tests.

## Detailed Requirements

- `testing.Short()` guards on tests that build the binary, run `go run`, or send process signals; the fast tier passes `-short` to `go test` for this module, `--heavy` does not.
- Delete matrix tests whose rows another test already proves; list each in the commit body.
- `finalization` under 20 s wall with `-short`; whole module under 25 s wall.
- No package test file grows; `go test -race -count=1 ./...` green under both tiers.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** `cd skills/do-work/tools/do-work-cli && time go test -short -count=1 ./...`.
**Why RED now:** `-short` changes nothing; 35 s wall with finalization at 34.8 s.
**GREEN when:** the `-short` run finishes under 25 s wall with exit 0, and the plain run (heavy) still executes the binary-building and signal tests.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 3124 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes do-work-cli test fixtures and migration parity.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.


## Addendum (2026-09-03)

User added (2026-09-03 21:29 local, in the test-budget session; 22:05 local, "update the batch REQs per A1-A3 via queued addenda" referring to the velocity report at `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/`, items A1 and A2):

> ```
> each test file should finish under 30 seconds (use the 80% value 20% effort principle until this is obtained)
>
> The rest of the test are accessible only when calling them with the --heavy parameter.
>
> the catch to call the --heavy parameter need to ask user for permission, and it should not block anything, meaning that where --heavy is required those tasks go into pending-testing status.
>
> Also because these tests tend to ballon, make sure to always measure the test duration, and adjust when the limits are reached.
> ```

- A1, dependency: `depends_on` changed from `[REQ-539]` to `[]`. This REQ edits do-work-cli test files and one argument line in `_dev/tests/maintainer-verify.sh`; REQ-539 edits the shell contract file and the aggregate. REQ-538 and this REQ are about 90% of the gate's seconds (board package 467 s, do-work-cli 111 s in-test of a 641 s run at 8d9d1bb) and can be built by two builders at once. REQ-539's `_dev/tests/*.sh` glob covers `maintainer-verify.sh` too; whichever of the two lands second rebases that one line.
- A2, per-package budget, Go side: every do-work-cli package finishes under 30 s wall in the fast tier (`-short`). The gate measures per-package wall time from `go test -json`, prints it, appends one row per package per run to the durations log REQ-539 introduces (or introduces the log itself if it lands first, same columns: run id, file or package, seconds, concurrent gate count), and fails the fast tier when a package exceeds 30 s. `--heavy` has no budget. The original "finalization under 20 s, module under 25 s" targets stay; 30 s per package is the enforced ceiling.
- GREEN additionally requires: every package's recorded fast-tier duration under 30 s in the proving run, and the durations log carrying one row per package.
- Coherence check: no contradiction with the original sections; the dependency change is the only frontmatter edit.
