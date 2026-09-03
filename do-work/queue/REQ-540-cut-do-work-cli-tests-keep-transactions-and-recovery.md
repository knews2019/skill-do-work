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
depends_on: [REQ-539]
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

