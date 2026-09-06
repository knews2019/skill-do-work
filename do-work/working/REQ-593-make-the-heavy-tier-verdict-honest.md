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
