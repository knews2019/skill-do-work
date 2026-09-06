# REQ-593 Hand-back — making the heavy tier's verdict honest

**Branch:** `worktree-agent-REQ-593-heavy-tier-verdict`
**Head:** `cbd8f2b03c538ade2f50ac8c8b068101f4016ff3` (one commit on base `2693306`)
**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-593-heavy-tier-verdict`
Working tree is clean. Nothing outside this worktree was written except this file.

## What changed

Exactly the five files the request's `## Scope` declares. `+100 / -8`.

- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go` — `runOneLane` clears the
  skip announcement when the lane exited non-zero, so `Skipped` is false and the returned skip line is
  empty for a lane that ran and failed. `RunLanes` then falls through to `HEAVY-RUN-LANE-RED`.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` — the lane output
  writer is now a parameter of a new `runHeavyVerificationLanes`; `handleRunHeavyVerification` is a
  one-line wrapper that passes `os.Stderr`.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go` — the test helper
  passes `io.Discard`, plus a `skip-then-fail-lane` fixture (prints the skip line, exits 4) and
  `TestRunLanesReportsARedLaneThatAlsoPrintedASkipLine`.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` —
  `TestRecoveryRefusesFalseLegacyCheckpointAbsence` calls `configureStateGit(t, root)`, like its eleven
  siblings.
- `_dev/tests/update-script-behavior.sh` — both matchers use a herestring, plus
  `verify_output_matchers_read_greps_verdict`, a self-check that drives the two shipped matchers over a
  256 KiB probe output with the pattern on line 1.

## Decisions

**D1 — the lane output writer is an explicit parameter, not a package variable.** The candidate patch
proposed `var laneOutputTee io.Writer = os.Stderr` with the tests swapping it and restoring it in
`t.Cleanup`. I rejected that. A parameter needs no mutable package state, no save/restore bookkeeping in
every test, and cannot be left pointing at the wrong writer by a test that fails before its cleanup. The
registered `CommandHandler` signature is unchanged: `handleRunHeavyVerification` keeps it and forwards
`os.Stderr`. Cost is one extra function; benefit is that the seam is visible at every call site.

**D2 — the exit-status rule lives in `runOneLane`, not in the `RunLanes` switch.** Reordering the switch
so red is checked before skipped would produce the right finding but leave `HeavyLaneExecution.Skipped`
set to true on a lane that exited non-zero. That field is what is printed, stored as lane evidence, and
copied back by `reusedLaneExecution`, so the record itself has to be honest, not just the finding.

**D3 — H3's assertion is a deterministic self-check, not a replay of the race.** The exploration measured
2-3 false failures in 500 runs at 36 KB. That is too weak to serve as a lock-in. Measured here: at 36,000
bytes the pipeline form produced 0 false failures in 50 runs, at 200,000 bytes it produced 50 of 50, at
1,000,000 bytes 50 of 50. The self-check uses 262,144 bytes, well past the turn, and the herestring form
produced 0 of 50 at every size. The check also covers `assert_output_lacks` in its dangerous direction: it
asserts that a pattern the output contains IS flagged, which the pipeline form silently swallowed.

**D4 — the tar pipelines in the same file were left alone.** See Discovered Tasks F1. The request's
requirement names the two assertion helpers, and its Scope already deferred the same mechanism in two
other files as a follow-up. I kept that line rather than widening the diff.

## RED evidence (each fix, alone)

**H1 — a red lane reported as skipped.** With the `heavy_run.go` hunk reverted and everything else in
place, `go test -run TestRunLanesReportsARedLaneThatAlsoPrintedASkipLine ./internal/heavyverification/`
exits 1:

```
--- FAIL: TestRunLanesReportsARedLaneThatAlsoPrintedASkipLine
    heavy_run_test.go:167: a lane that ran and failed was recorded as skipped: []resultmodel.HeavyLaneExecution{
      resultmodel.HeavyLaneExecution{LaneID:"skip-then-fail-lane", ..., ExitStatus:4, Skipped:true, ...}}
```

`ExitStatus:4, Skipped:true` is exactly the misreport the heavy run at `6646ba51` produced.

**H2 — the fixture with no git identity.** With the `configureStateGit` line reverted, under the lane's
own environment (`GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null`):

```
--- FAIL: TestRecoveryRefusesFalseLegacyCheckpointAbsence
    state_apply_test.go:1055: git [commit -qm legacy fixture]: exit status 128: Author identity unknown
```

**H3 — the matchers reading the writer's death.** With both matchers reverted to the `printf | grep`
form and the self-check kept, the full updater probe script exits 1 with:

```
FAIL: output matchers: on a 256 KiB probe output a matcher reported the writer's exit status instead of
grep's verdict
```

Isolated repeats of the self-check body: pipeline form broken 5 of 5, herestring form clean 5 of 5, and
in each broken run BOTH directions were wrong (`assert_output_matches` missed a present pattern,
`assert_output_lacks` failed to flag it).

## GREEN evidence

All commands run from the worktree root.

1. **The lane the request names, at lane level** —
   `env -u NODE_OPTIONS -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0..2 -u GIT_CONFIG_VALUE_0..2
   GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane
   do-work-cli-integrations` → **exit 0**, 795 tests, wall 27s, **0** lines matching `^--- FAIL`, **0**
   lines matching `^SKIP:`. That zero-SKIP count is the lane-level proof of the writer seam: the fixture
   lanes still print their announcements, they just no longer reach the enclosing lane's stderr.
   Log: `.../scratchpad/lane-integrations.log`

2. **The same lane through the runner** — `bash skills/do-work/tools/do-work-cli.sh
   run-heavy-verification --lane do-work-cli-integrations` → **exit 0**:

   ```
   run-heavy-verification: success
     lane do-work-cli-integrations: exit 0 in 26s [executed: no_prior_evidence]
   ```

   Real disposition, no `skipped`, no findings section, 0 `^SKIP:` lines on stderr.
   Logs: `.../scratchpad/heavy-run-e2e.out`, `.err`

3. **Updater probes direct** — `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/update-script-behavior.sh`
   → **exit 0**, "update-script behavior probes passed." Log: `.../scratchpad/green-h3.log`

4. **Updater heavy lane** — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater` → **exit 0**, 0
   FAIL lines, 70s. Log: `.../scratchpad/lane-updater2.log`

5. **Both Go packages under the lane environment** — `DO_WORK_HEAVY_TESTS=1 go test -count=1
   ./internal/heavyverification/ ./internal/requeststate/` → **exit 0**.

6. **Whole do-work-cli module** — `go test -count=1 ./...` under the gate environment → **exit 0**.
   Log: `.../scratchpad/module-tests.log`

7. **Lint** — `gofmt -l` on both packages: empty. `go vet` on both packages: exit 0.
   `shellcheck --severity=warning _dev/tests/update-script-behavior.sh`: exit 0. `bash -n`: exit 0.

(`.../scratchpad` is
`/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad`.)

## The full canonical gate: one failure, and it is not this REQ's

The canonical gate was run once, as instructed:

```
env ... GIT_CONFIG_GLOBAL=<gate config> QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  bash _dev/tests/maintainer-verify.sh
```

**Exit 1, with exactly one failing test**, and it is not one of the five files:

```
heavy_reuse_regression_test.go:174: input mutation refusal: commit or stash beta/fixture.txt before
running heavy lanes; a lane result must be attributable to HEAD, want HEAVY-RUN-REVISION-CHANGED
--- FAIL: TestLaneMutationCannotPublishOrReuseSuccess/commit=true
```

Why it is not mine, and what it actually is:

- **Blame.** The failing call site entered in commit `d6709da` ("Refuse heavy-lane reuse when a lane
  changes another lane's inputs"), which predates this branch's base `2693306`. My diff does not touch
  `heavy_reuse_regression_test.go`, and the `heavy_run.go` hunk only sets `announcedSkipLine` and
  `Skipped` after the lane has already exited. The refusal code the test asserts is chosen in
  `verifyLaneRevision`, which I did not change.
- **It does not reproduce.** 5 repeats in isolation: green. 6 repeats pinned to one CPU (`taskset -c 0`):
  green. The whole do-work-cli module (`go test -count=1 ./...`, evidence 6 above): green. The whole
  heavyverification package inside the `do-work-cli-integrations` heavy lane (evidence 1): green, twice.
- **Mechanism.** `TestLaneMutationCannotPublishOrReuseSuccess` calls `RunLanes` with **no
  `LaneTimeout`**, so the field is zero. `RunLanes` applies no default — the 1800-second default lives in
  `parseRunArguments` in `heavy_commands.go`, on the CLI path only — so `runOneLane` does
  `time.NewTimer(0)` and the lane races its own termination bound from the instant it starts. The lane
  script is `echo broken > beta/fixture.txt` then `git add … && git commit`. `ownedprocess.TerminateGroup`
  gives a one-second grace, which the commit normally wins. Under gate-level contention (4 CPUs, several
  packages in parallel, other builders on the machine) it loses, the commit never happens, and the runner
  correctly reports `HEAVY-RUN-DIRTY-TREE` for a tree that really is dirty — while the test expected
  `HEAVY-RUN-REVISION-CHANGED`.
- **Remedy, not applied.** One line: give that `LaneRunRequest` a real `LaneTimeout` (its sibling at
  `heavy_evidence_test.go:103` uses `30 * time.Second`). The broader fix is for `RunLanes` to apply
  `defaultLaneTimeoutSeconds` when `LaneTimeout` is zero, so no caller can accidentally ask for an
  instant timeout. **That needs a sixth file, so I stopped and did not change it** — see Discovered Tasks
  F2.

**A second, unrelated local artifact.** The first `--heavy-lane updater` run failed with `FAIL: test
duration log has an invalid header: do-work/test-durations.tsv`. That file is git-ignored and untracked;
in this fresh worktree it had been created without its header row by an earlier test process. Deleting it
made the lane green (evidence 4). Nothing in the repository changed. See F3.

## Discovered Tasks

- `_dev/tests/update-script-behavior.sh` has nine more `writer | grep -q` pipelines under `set -o
  pipefail` (lines 607, 610, 615, 627, 630, 675, 678, 722, 725) — the same SIGPIPE mechanism this REQ
  fixed in the two assertion helpers, in a construct the requirement did not name. Seven are `tar tzf
  <archive> | grep -q <marker>`, and a repository archive listing is long enough for the writer to lose
  the race. One of them fired during this session: the ablation run recorded `FAIL: upstream fetcher:
  requested branch archive omitted its marker` on an archive that does contain the marker, and it did not
  recur in the green run. This is a false failure in the heavy updater lane, which is the exact class
  REQ-593 exists to remove. → impact-critical
- `RunLanes` in `internal/heavyverification/heavy_run.go` applies no default when `LaneRunRequest.LaneTimeout`
  is zero, so an in-process caller that omits it gets `time.NewTimer(0)` and a lane that is terminated
  while it is still starting. `TestLaneMutationCannotPublishOrReuseSuccess` in `heavy_reuse_regression_test.go`
  is such a caller and flakes under gate load (full diagnosis above). Fix is either the one-line
  `LaneTimeout` in that test or a zero-means-default in `RunLanes`; the second is the one that stops it
  happening again. → impact-critical
- `_dev/tests/test-duration-log.sh` seeds its header with a `printf` to a candidate file plus `ln … ||
  true`, then validates the header and hard-fails the whole lane if it is missing. Nothing makes that
  seeding atomic against a concurrent `record_test_duration` append, so a worktree can end up with a
  headerless duration log that fails every subsequent lane until a human deletes an ignored file. The
  failure message does not say that deleting the file is the remedy. → report only
- Changelog and version handling was not done. The commit changes shipped files under `skills/`, which
  `_dev/primes/prime-releases.md` calls a release, so a `CHANGELOG.md` entry, the byte-identical mirror at
  `skills/do-work/CHANGELOG.md`, and a `VERSION` bump are owed. All three are outside the declared
  five-file scope and belong to the work action's Step 9 finalize transaction, not to this build.
  → report only

## Acceptance criteria

- [x] A lane that ran and exited non-zero is reported red, whatever it printed — `heavy_run.go`, pinned by
      `TestRunLanesReportsARedLaneThatAlsoPrintedASkipLine`, red when reverted.
- [x] A test's own fixture output cannot be mistaken for the enclosing lane's announcement —
      `runHeavyVerificationLanes` + `io.Discard`, proved at lane level by 0 `^SKIP:` lines in evidence 1
      and 2.
- [x] Every test in `state_apply_test.go` that commits sets its own identity — the one outlier now calls
      `configureStateGit`; red when reverted.
- [x] Both matchers in `update-script-behavior.sh` report grep's verdict, never the writer's status —
      herestring form, pinned by `verify_output_matchers_read_greps_verdict`, red when reverted, both
      directions.
- [x] Each fix carries an assertion that fails when the fix is reverted — three ablations above, each run
      with the other two fixes in place.
