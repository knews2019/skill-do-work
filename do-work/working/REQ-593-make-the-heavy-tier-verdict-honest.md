---
id: REQ-593
status: claimed
domain: testing
created_at: 2026-09-06T02:02:41Z
user_request: UR-105
review_generated: true
impact: impact-critical
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
maintenance: false
depends_on: []
related: [REQ-552, REQ-554, REQ-585, REQ-592]
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-06T02:03:38Z
  basis:
    - Route B
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
route: B
write_set: [skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go, skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go, skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go, skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go, _dev/tests/update-script-behavior.sh]
title: '[impact-critical] Make the heavy tier report a red lane as red, and fix two fixtures that cannot pass under the lane environment'
claimed_at: 2026-09-06T02:03:03Z
---

# Make the Heavy Tier's Verdict Honest

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Three defects found by running the six heavy lanes for the first time in this checkout. Together they
mean the heavy tier can report success for a lane that ran and failed — the same false-green shape
REQ-592 removed from the fast gate, in the tier that is supposed to be stricter.

**H1 — a red lane is reported as skipped.** `runOneLane` tees a lane's combined output through a
watcher that marks the lane "did not run" as soon as any line starts with `SKIP:`, and `Skipped` is
computed without consulting the exit status. The `do-work-cli-integrations` lane runs the
heavyverification package's own tests, one of which executes a fixture lane whose script prints
literally `SKIP: no browser is available`; `handleRunHeavyVerification` hardcodes the lane output
writer to `os.Stderr`, so the fixture's line reaches the real lane's stderr and the parent watcher
cannot tell it from a real announcement. In the run that found this, the lane exited 1 with a failing
test and was recorded as `HEAVY-RUN-LANE-SKIPPED`.

**H2 — one fixture cannot pass under the lane's own environment.**
`TestRecoveryRefusesFalseLegacyCheckpointAbsence` is the only test in `state_apply_test.go` that
commits in a fixture repository without first calling `configureStateGit`, and every lane argv in
`_dev/tests/heavy-lanes.json` forces `GIT_CONFIG_NOSYSTEM=1` and `GIT_CONFIG_GLOBAL=/dev/null`. With
no identity to inherit, `git commit` fails with `Author identity unknown` on any machine where git
cannot auto-detect an address from the hostname.

**H3 — two assertion helpers report the writer's death as a failed match.**
`_dev/tests/update-script-behavior.sh` runs under `set -uo pipefail` and both matchers are
`printf '%s' "$probe_output" | grep -Eq -- "$pattern"`. `grep -q` exits at the first match; on a long
probe output the `printf` still has chunks to write and dies on SIGPIPE; `pipefail` makes that the
pipeline's status. `assert_output_matches` then reports "output did not match" for output that
plainly contains the pattern. `assert_output_lacks` has the same bug in the more dangerous
direction: a pattern that should have been flagged is silently swallowed.

## Context

Found while draining the heavy lanes to finalize REQ-552 and REQ-554. All three predate this work
run: H1 and H2 both entered in commit `7bd3464` (REQ-585), and H3's matchers have never been edited
since the repository import. Nothing in REQ-486, REQ-552, REQ-554 or REQ-592 touches the affected
files. H1 masked H2 — the false SKIP suppressed the red finding for the lane that was actually failing.

## Requirements

- A lane that ran and exited non-zero is reported red, whatever it printed. An exit status outranks a
  skip announcement.
- A test's own fixture output cannot be mistaken for the enclosing lane's announcement.
- Every test in `state_apply_test.go` that commits sets its own identity, so no test depends on the
  host's git configuration.
- Both assertion helpers in `update-script-behavior.sh` report grep's verdict, never the writer's
  exit status.
- Each fix carries an assertion that fails without it.

## Red-Green Proof

**RED now:** `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
exits 1 with exactly one failing test, `TestRecoveryRefusesFalseLegacyCheckpointAbsence`, and the
same lane run through `do-work-cli run-heavy-verification` reports it as **skipped**, not red.
For H3, a matcher fed a probe output larger than the pipe buffer reports a miss on a pattern the
output contains — reproducible deterministically by enlarging the fixture payload, and observed at
2-3 failures in 500 runs on a single pinned CPU with the real 36 KB output.

**GREEN when:** the lane runs, its real exit status is reported, all three fixes are in, and a new
assertion for each fails when its fix is reverted.

## Full Context

Full diagnoses, with blame, verified fixes and measured evidence, are in this run's directory as the
heavy-lane diagnosis; each was developed and proven in a scratch copy without modifying the
repository.

## Triage

**Route: B** — Explore then build.

**Reasoning:** Three separate defects with three separate mechanisms, in two languages, found by
running the heavy lanes rather than by reading. The outcomes are stated exactly and each fix has
already been developed and proven in a scratch copy; what needed discovery was the mechanism of each,
and that is done. Nothing left to plan.

**Planning:** Skipped.

**Why this jumped the queue.** It blocks every remaining finalization. The heavy drain is what
archives a completed request, and until the heavy tier can tell a red lane from a skipped one, no
heavy result in this checkout is evidence of anything.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Three independent diagnoses, each run against the heavy-lane log at revision `6646ba51`, each
developing and proving its fix in a scratch copy without modifying the repository. Full report in the
run directory as `do-work/runs/work-2026-09-05-231943/REQ-593-exploration.md`.

**All three predate this work run, and the first masked the second.** `git log -S` places H1 and H2 in
commit `7bd3464` (REQ-585); H3's matchers have never been edited since the repository import. The
heavy-run record shows the lane as `[executed: no_prior_evidence]`, meaning `6646ba51` is the first
time that lane has ever executed since `7bd3464` landed — which is why nobody had seen it.

**H1 is a false green in the tier that exists to be strict.** `runOneLane` marks a lane "did not run"
on any output line beginning with `SKIP:`, and computes `Skipped` without consulting the exit status.
The `do-work-cli-integrations` lane runs the heavyverification package's own tests, and
`TestRunLanesMarksSkipFromExplicitSkipLine` executes a fixture lane whose script prints exactly
`SKIP: no browser is available`. `handleRunHeavyVerification` hardcodes the lane output writer to
`os.Stderr`, so the fixture's line lands on the enclosing lane's stderr and the parent watcher cannot
tell it from a real announcement. In the run that found this, the lane exited 1 and was recorded as
`HEAVY-RUN-LANE-SKIPPED`. The exposure is wider than one fixture: `internal/corehelpers/checks.go`
prints user-facing lines beginning with `SKIP:` at three sites, and two shipped test scripts do too.

**H2 is a fixture that cannot pass under its own lane's environment**, on any machine rather than only
this one. An `awk` scan over `state_apply_test.go` for tests that commit without calling
`configureStateGit` returns exactly one name, and every lane argv forces `GIT_CONFIG_NOSYSTEM=1` and
`GIT_CONFIG_GLOBAL=/dev/null`. It survives on a workstation only because git can guess an address from
an FQDN hostname.

**H3 is a race, and it was measured rather than argued.** Both matchers in
`_dev/tests/update-script-behavior.sh` are `printf | grep -Eq` under `set -o pipefail`. Replaying the
exact 36 KB output from the failing run through the original matcher on a single pinned CPU produced
2-3 false failures in 500 runs; the herestring form produced 0 in 500. Forcing the pipe capacity to
4096 bytes reproduces it deterministically: grep exits 0 while the writer dies with a write error.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go` (modify) — a non-zero exit status clears the skip announcement, so a lane that ran and failed is red whatever it printed
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` (modify) — the hardcoded lane output writer becomes a package variable the package's own tests can redirect
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go` (modify) — fixture lane output no longer reaches the enclosing lane's stderr, plus a lock-in for a lane that prints a skip line and then fails
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modify) — the one test that commits without an identity gets the same helper its eleven siblings call
- `_dev/tests/update-script-behavior.sh` (modify) — both matchers read grep's verdict instead of the writer's exit status

**Files I will NOT touch:** `_dev/tests/heavy-lanes.json` — the lane argv is not what is wrong.
`_dev/tests/select-simple-reqs-behavior.sh` and `_dev/tests/p50-estimator-determinism.sh` carry the
same pipeline construct, but their probe outputs are a few hundred bytes, far below where the race can
open; converting them is a follow-up, not a fix. The three `SKIP:`-prefixed sites in
`internal/corehelpers/checks.go` — a stricter skip channel than "any line starting with SKIP:" is a
design change and its own request.

**Acceptance criteria (restated from REQ):**
- [ ] A lane that ran and exited non-zero is reported red, whatever it printed
- [ ] A test's own fixture output cannot be mistaken for the enclosing lane's announcement
- [ ] Every test in `state_apply_test.go` that commits sets its own identity
- [ ] Both matchers in `update-script-behavior.sh` report grep's verdict, never the writer's status
- [ ] Each fix carries an assertion that fails when the fix is reverted

## Pre-Flight

**Git:** ✓ Clean. Canonical `recover` reports `FINALIZATION-NONE`. Seven claims sit in
`do-work/working/`: REQ-583, REQ-587, REQ-591 and REQ-592 held at Step 7.7; REQ-486 in flight on the
board module; REQ-552 and REQ-554 merged, reviewed, remediated and waiting on the very heavy drain this
request unblocks. None of them touches this request's five files.

**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exited 0 at this REQ's claim revision
`2693306` — **76s wall**, exit status read directly from `$?`. The fast tier has never seen any of
these three defects, which is the point: two of the three files this request fixes run only at heavy
tier, and the third is a race that needs a long output to open.

**Tests baseline:** ✓ `go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/heavyverification/
./internal/requeststate/` exited 0, launched true — under the *inherited* git configuration. Under the
lane's own configuration it does not, and that is H2.

**The failing state is recorded rather than described.** At this same revision,
`env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
exits **1** with exactly one failing test, `TestRecoveryRefusesFalseLegacyCheckpointAbsence`, and the
same lane driven through `do-work-cli run-heavy-verification` reports it as **skipped in 29s** with a
`HEAVY-RUN-LANE-SKIPPED` finding. Those two facts together are the request.

**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just` 1.43.0, Node v22, Chromium present. `git`
2.43.0. The container hostname is `vm` with no domain, which is what makes H2 visible here and
invisible on a workstation — but H2 is not a container defect: the lane argv strips the configuration
that would otherwise hide it.

*Checked by work action*
