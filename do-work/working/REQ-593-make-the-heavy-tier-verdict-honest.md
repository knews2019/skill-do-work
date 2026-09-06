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
dispatch_at: 2026-09-06T02:26:22Z
builder_handback_at: 2026-09-06T02:26:22Z
write_set: [skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go, skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go, skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go, skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go, _dev/tests/update-script-behavior.sh, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
title: '[impact-critical] Make the heavy tier report a red lane as red, and fix two fixtures that cannot pass under the lane environment'
claimed_at: 2026-09-06T02:03:03Z
---

# Make the Heavy Tier's Verdict Honest

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Both `prime_files` read, plus the crew rules and the three diagnoses. The candidate
  patch from the diagnosis was read and deliberately not taken as written — see D-01.
- [x] **[APPLY]:** Five files, exactly the declared `write_set`. One sixth file was needed to fix a
  flake found on the way and the builder stopped and reported it rather than editing it.
- [x] **[UNIFY]:** `git diff --stat` on the merge range reports five files, 100 insertions, 8 deletions.
  Linters: `gofmt -l` on both packages — empty; `go vet` on both — exit 0;
  `shellcheck --severity=warning` and `bash -n` on the shell file — exit 0. No debug artifacts. Per
  file: `heavy_run.go` (the exit-status rule, in the record writer rather than the finding chooser),
  `heavy_commands.go` (the writer becomes a parameter, the registered handler signature unchanged),
  `heavy_run_test.go` (the skip-then-fail lock-in and the writer redirection),
  `state_apply_test.go` (one added call), `update-script-behavior.sh` (two matcher bodies and one
  deterministic self-check).

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

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go`
- `_dev/tests/update-script-behavior.sh`

**What was done:** A lane's exit status now outranks anything it printed, and the rule sits where the
record is written rather than where the finding is chosen — `Skipped` is printed, stored as lane
evidence and copied back on reuse, so the record has to be honest and not merely the message. The lane
output writer became an explicit parameter, so the heavy verification package's own fixture lanes can
no longer leak their output onto the enclosing lane's stderr; that was proved end to end rather than by
unit assertion, with zero lines beginning with `SKIP:` where the pre-fix run leaked one. The one
fixture that committed without a git identity now sets one like its eleven siblings. And both matchers
in the updater probe read grep's verdict instead of the writer's exit status.

**The sizing was measured before the fix was chosen.** The pipeline form produces 0 false failures in
50 runs at 36 KB, and 50 in 50 at 200 KB — and in every broken run *both* directions were wrong: a
present pattern went unmatched, and the same pattern went unflagged by the negative matcher. That is
why the lock-in is a deterministic self-check at 256 KiB driving the shipped matchers, rather than a
replay of the 2-in-500 race the diagnosis started from.

Merge range `b8398be7..30bb2733`, five files, 100 insertions, 8 deletions. Builder branch head
`cbd8f2b0`.

## Decisions — implementation

- **D-01 — the lane output writer is a parameter, not package state. DECIDE & STATE.** The candidate
  patch from the diagnosis proposed a package-level `io.Writer` the tests would point at `io.Discard`.
  A parameter needs no save-and-restore in every test, cannot be left pointing at the wrong place by a
  test that fails before its cleanup, and leaves the registered handler signature unchanged.
- **D-02 — the exit-status rule lives in `runOneLane`, not in the `RunLanes` switch. DECIDE & STATE.**
  Reordering the switch would produce the right *finding* while leaving `Skipped` true on the stored
  record, and that field is printed, persisted as lane evidence and copied back by `reusedLaneExecution`.
  A correct message over a false record is the shape this whole request exists to remove.
- **D-03 — the shell lock-in is a deterministic self-check, not a replay of the race. DECIDE & STATE.**
  A 2-in-500 flake is too weak to lock anything in. The self-check drives the two shipped matchers
  through bash dynamic scope — not copies of them — and asserts both directions.
- **D-04 — nine sibling pipelines in the same file were reported rather than fixed. ESCALATE, and
  overturned by the orchestrator.** The builder's reasoning was that the requirement names the two
  assertion helpers. The counter-argument is stronger: the nine are the same mechanism in the same file
  under the same `pipefail`, seven of them are `tar tzf … | grep -q` over an archive listing, and one
  of them *fired during this very session*. Leaving a known impact-critical false-failure mechanism in
  a file already being edited is the failure this request is about. Taken as a remediation.

## Discovered Tasks

- **impact-critical, taken here as a remediation — nine more `writer | grep -q` pipelines** in
  `_dev/tests/update-script-behavior.sh` (lines 607, 610, 615, 627, 630, 675, 678, 722, 725), the same
  SIGPIPE mechanism the two named helpers had. One produced
  `FAIL: upstream fetcher: requested branch archive omitted its marker` during this session on an
  archive that does contain the marker.
- **impact-critical, taken here as a remediation — `RunLanes` applies no default lane timeout.** The
  1800-second default lives in `parseRunArguments`, on the CLI path only, so an in-process caller that
  omits `LaneTimeout` gets `time.NewTimer(0)` and a lane terminated while it is still starting. That is
  why `TestLaneMutationCannotPublishOrReuseSuccess` flakes under gate load, reporting
  `HEAVY-RUN-DIRTY-TREE` where it expects `HEAVY-RUN-REVISION-CHANGED` — a false verdict from the heavy
  runner, arriving by a different route than H1.
- **report only — `_dev/tests/test-duration-log.sh` seeds its header non-atomically**, so a concurrent
  append can leave a headerless git-ignored duration log that fails every subsequent lane until a human
  deletes it, and the failure message does not name that remedy. Seen once this session.

## Qualification

**Passed.** Read from the merge range `b8398be7..30bb273a`; canonical `qualify` and `scope-drift` both
satisfied. Five files, 100 insertions, 8 deletions — the declared set exactly.

- **Each of the three fixes was reverted alone, with the other two in place, and its assertion watched
  to fail.** That matters more here than usual: three unrelated mechanisms in one change is exactly the
  shape where one fix quietly covers for another's missing test. H1's revert produces
  `LaneID:"skip-then-fail-lane", ExitStatus:4, Skipped:true` — byte-for-byte the misreport the heavy run
  at `6646ba51` produced. H2's revert produces `Author identity unknown` under the lane's own stripped
  git configuration. H3's revert fails the new self-check by name.
- **The writer seam was proved end to end, not by unit assertion.** Zero lines beginning with `SKIP:`
  in both the direct lane run and the runner's stderr, where the pre-fix run leaked exactly one. A unit
  test could have asserted the parameter is threaded; only the lane run shows the leak is gone.
- **The shell fix was sized before it was chosen.** 0 false failures in 50 runs at 36 KB, 50 in 50 at
  200 KB, 50 in 50 at 1 MB — and in every broken run *both* matchers were wrong, the positive one
  missing a present pattern and the negative one failing to flag it. That measurement is what rejected
  a replay of the original 2-in-500 race as the lock-in and produced a deterministic 256 KiB self-check
  driving the shipped matchers rather than copies of them.
- **The lane the request exists to fix now reports its real disposition.** Direct:
  `exit 0, 795 tests, wall 27s`, zero `--- FAIL` lines. Through the runner:
  `lane do-work-cli-integrations: exit 0 in 26s [executed: no_prior_evidence]`, no findings. Before this
  change the same lane exited 1 with a failing test and was recorded as *skipped*.
- **One failure was found on the way, correctly not fixed, and correctly not dismissed.** The single
  canonical gate run exited 1 on `TestLaneMutationCannotPublishOrReuseSuccess/commit=true`, in a file
  outside the five-file scope. The builder did not call it a flake: it reproduces 0/5 in isolation, 0/6
  pinned to one CPU, and green in the whole-module and lane runs — and then it named the mechanism.
  `RunLanes` applies no default when `LaneTimeout` is zero, so an in-process caller gets
  `time.NewTimer(0)` and the lane is killed before its own `git commit`. That is a false verdict from
  the heavy runner under load, and it is taken as a remediation because the fix belongs in
  `heavy_run.go`, which is in scope.
- **This also corrects something recorded earlier in this run.** REQ-592's pre-flight attributed that
  same test's failure entirely to a global `commit.gpgsign` with an unusable key. That attribution was
  right about what it saw — with signing off it passes deterministically — but incomplete: there are two
  independent causes, and the second one only shows under gate contention.

### Remediation qualification (after the first review round)

**Passed, and the record above is corrected here rather than left stating half the change.** The
Implementation Summary and the Qualification above were both written at `aa603e9`, which lands before
the remediation. They say "five files, 100 insertions, 8 deletions — the declared set exactly", and
that is true only of the first range. Remediation merge range `32ecdba..407432c` adds three more
files, 89 insertions, 12 deletions — all inside the same declared five — and one of them is a
production behaviour change no requirement asked for. It is recorded as D-05 below rather than left
for a reviewer to discover.

- **The nine sibling pipelines are converted, and the guard is a class scanner rather than nine
  assertions.** Seven `tar` listings are captured once and read as herestrings; the two `find` sites
  were converted for a different defect the pipeline hid — a `find` that could not read the tree
  produced no output, which read as a clean tree, with its failing status absorbed by the `if`. Each of
  the nine was reverted alone and the scanner named that exact line.
- **D-05 — `RunLanes` now treats a zero or negative lane timeout as the default.** No requirement asked
  for it. It is here because the single canonical gate run during the first round failed
  `TestLaneMutationCannotPublishOrReuseSuccess`, and the builder root-caused it rather than calling it
  a flake: the 1800-second default lived on the CLI path only, so an in-process caller that omitted the
  field armed `time.NewTimer(0)` and the lane was killed before its own `git commit`. That is a false
  verdict from the heavy runner under load — the same failure class as H1 by a different route — and
  the fix belongs in a file already in scope. Reverting it makes the new lock-in report exit 124, the
  lane-timeout status.
- **One claim in this run's merge message overstates the result and is corrected here.** That message
  says "the SIGPIPE class is closed". It is closed in **one file of twenty-three**: the scanner reads
  `${BASH_SOURCE[0]}` and nothing else, and roughly 130 sibling sites remain across `_dev/tests/`,
  including nineteen in the gate's own `contracts/core-checks.sh` — which disagrees with itself, using
  the herestring form at one site and the pipeline at another. The class is closed *in the updater
  probe*. The generalization is a follow-up, and until it exists the claim should not be repeated.

Requirements traced: a lane that ran and exited non-zero reports red whatever it printed; a test's own
fixture output cannot be mistaken for the enclosing lane's announcement; every test in
`state_apply_test.go` that commits sets its own identity; both matchers report grep's verdict; and each
fix carries an assertion that fails when the fix is reverted.

*Checked by work action*

## Review

**Overall: 71%** | 2026-09-06T03:05:35Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 75% |
| Test Adequacy | 70% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Partial |

**Approve with follow-ups** — the three declared fixes are real and hold under mutation, but one path to a false green heavy verdict is still open, and one fix has no assertion that fails when it is reverted.

REQ-593 (making the heavy verification tier's verdict honest) removes the case where a lane that ran and failed was reported as skipped. That case is closed: the exit-status rule sits in the record writer at `heavy_run.go:322-333`, so the stored `Skipped` field is honest and not only the printed finding. All three reviewers reverted the rule independently and got the same failing test with the same `ExitStatus:4 / Skipped:true` misreport. The request's own purpose is not fully served, because a second route to a false green survives in the same file and has the same shape as the one the remediation just closed.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**

- `RunLanes` reports a successful run with zero lanes when `LaneIDs` is empty, which is a false green heavy verdict from the function this request exists to make honest. `selectManifestLanes` (`heavy_run.go:250-261`) treats an empty request as "match nothing" and returns an empty lane list with no error, so the loop runs zero times and the run returns `Lanes: []` with no findings and no error. Probe in a scratch copy of the module: `RunLanes(LaneRunRequest{RepositoryRoot: root, ManifestPath: "heavy-lanes.json", LaneOutputWriter: io.Discard})` returned `err=<nil> lanes=0 findings=0`, rendered as "heavy verification run: lanes=0" with outcome success. Only the CLI argument parser refuses this (`heavy_commands.go:412`), exactly the CLI-only-default shape the remediation removed for `LaneTimeout` with the reasoning "LaneRunRequest is a struct, so an in-process caller with no opinion leaves the field zero". The guard belongs beside the `LaneTimeout` guard inside `RunLanes`. Found by one reviewer of three and confirmed in the source during this synthesis — impact-rule-change → report only
- Acceptance criterion 5 is not met for the writer-seam fix, and the evidence offered in its place cannot tell fixed from unfixed. Reverting the seam (the test helper passing `os.Stderr` instead of `io.Discard`) leaves the whole `heavyverification` package green with zero `^SKIP:` lines. `_dev/tests/run-go-tests-with-budget.sh:41` runs `go test -json` into a file and re-prints the Output events only when the package fails, so the "0 SKIP: lines in the lane run" evidence is 0 either way. Compiling the same reverted package with `go test -c` and running the binary directly gives exit 0 with one `SKIP: no browser is available` on stderr, so the leak is real and simply invisible to the lane. Two reviewers called this criterion met on the strength of that lane evidence, one showed the surviving mutation. The mutation wins, because a mutation that survives is proof and a green lane that cannot discriminate is not — impact-negligible → report only
- The new source scanner `verify_no_check_reads_a_quiet_grep_pipeline` is the sole guard for nine of the eleven converted shell sites, and it is evadable by ordinary shell. Reverting any one of the seven `tar` sites or two `find` sites raises exactly one FAIL, the scanner's, so its regex at `_dev/tests/update-script-behavior.sh:157` is the whole contract for those nine. Five rewrites of the real site at line 658, each reintroducing the SIGPIPE defect, each produced `fail_count=0` with the scanner silent: the pipe at end of line with `grep -q` on the next line, `grep --quiet`, `grep --silent`, `| LC_ALL=C grep -q`, `| command grep -q`. This is the Closed Enumerations Go Stale trap from `_dev/primes/prime-shell-commands.md` applied to a regex. Two reviewers found it independently with different probe sets — impact-rule-change → report only
- Discarding `tar`'s exit status silently narrowed two assertions, and the stated safety argument covers only a fully unreadable archive. For an archive that is truncated but whose partial listing still contains the marker, the old pipeline form recorded a failure through pipefail and the new captured-listing form passes. Demonstrated: an archive built with the marker first, 200 pad files, last 200 bytes truncated. `tar tzf` exits 2, the first listing line is the marker, the old form fired `record_failure` and the new form passed. Affects the requested-branch marker at line 658 and the default-branch marker at line 707, the two archives with no readability assertion. The fallback archive is genuinely safe, it is asserted readable at line 619. One line each fixes it — impact-rule-change → report only
- The merge commit 407432c claims "the SIGPIPE class is closed", and it is closed in one file out of 23. The scanner reads `${BASH_SOURCE[0]}`, its own file and nothing else. `grep -rc -E '\|[[:space:]]*grep[[:space:]]+-[A-Za-z]*q' _dev/tests/` returns 23 non-empty files, roughly 130 sites, including the gate's own contract file `_dev/tests/contracts/core-checks.sh` (19 sites) and `_dev/tests/prescribed-shell-cases/qualify.sh` (31). `core-checks.sh` disagrees with itself, line 547 uses the herestring form while line 49 still pipes. The request's Scope section named two other files and promised a follow-up, and `do-work/queue/` holds only REQ-555, REQ-557 and REQ-558, none about this. One reviewer graded this impact-critical. Downgraded to impact-rule-change here, because these are maintainer test scripts and not a shipped production path, and because the rule fix is one generalization applied in many places. The severity stays Important, the missing follow-up is the actual gap — impact-rule-change → report only
- The prime file the request itself declares still states the pre-change rule. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md:61`, trap `[family: silent-skip-reads-as-red]`, ends "the runner keys `skipped` on that prefix, never on a lane-name list". After this change the runner keys `skipped` on that prefix and a zero exit status. This is the first file the next agent touching the heavy runner reads, and it now describes the contract that was just replaced. One sentence closes it — impact-rule-change → report only
- The request file records only the first half of the delivered change. `do-work/working/REQ-593-make-the-heavy-tier-verdict-honest.md` was last written at aa603e9, which lands before the two remediation commits 45c860a and 407432c. Its Implementation Summary and Qualification both state "five files, 100 insertions, 8 deletions — the declared set exactly", while `git diff --stat 32ecdba..407432c` adds 3 files, 89 insertions, 12 deletions. That later range includes a production behavior change no requirement asked for: `RunLanes` now substitutes `defaultLaneTimeoutSeconds` for a zero or negative `LaneTimeout`. It appears in no `## Decisions` entry, and the five acceptance boxes are still unchecked. A reviewer who reads the request instead of both hand-backs reviews half the change and is told the counts are exact — impact-rule-change → report only

**Minor findings:**

- The exit-status rule halves the SKIP leak rather than containing it. A lane that runs every check, exits 0, and has any nested producer print an unindented `SKIP:` line is still recorded `Skipped:true`, gets `HEAVY-RUN-LANE-SKIPPED`, and never reaches `recordSuccessfulLane`, so it stores no reuse evidence either. Probe: a lane that echoes progress, runs `sh -c "printf 'SKIP: nested tool has no engine\n'"`, then exits 0 is recorded `{LaneID:skip-lane ExitStatus:0 Skipped:true}`. Per `actions/work.md:451` a skipped lane holds the REQ claimed, so a persistent leak of this shape holds a request indefinitely on a lane that is actually green. The request declares the stricter channel out of scope, which is defensible, so this records the residue rather than reopening the decision — impact-user-visible → report only
- The red finding drops the announced skip line, so the operator loses the reason the lane gave. `runOneLane` clears `announcedSkipLine` and `laneRedFinding` (`heavy_run.go:418-427`) takes only the execution. Probe: a lane printing `SKIP: no browser is available` then exiting 2 yields evidence `["heavy lane skip-lane exited 2 after 0s"]` and the announcement appears nowhere. The verdict is now right and the announcement is lost, and those are not exclusive. This bites the shipped probe lanes, where post-probe work runs under `set -euo pipefail` after the probe script — impact-user-visible → report only
- The scanner false-positives on the shape inside a string, a here-doc body, or a trailing comment. The awk program skips only whole-line comments. Three probes on the real script each produced `fail_count=1` with correct code: a `printf` containing the shape as a hint, the shape inside a `cat <<'DOC'` body, and a trailing `# was: tar tzf a | grep -q m`. The file already contains three here-docs that write fixture sources, so a future fixture using `| grep -q` legitimately would break the updater lane — impact-negligible → report only
- `handleRunHeavyVerification` is now one line whose whole content is the `os.Stderr` choice, and nothing tests it. `heavy_commands_test.go:52` calls it with only `--manifest`, so argument parsing refuses before any lane runs. That untested line is the one that would swallow every lane diagnostic if someone changed it to `io.Discard`, and nothing stops the next in-package test from calling it with a real `--lane` and reopening the leak the request removed. A one-line note saying the package's own tests go through `runHeavyVerificationLanes` is the cheapest containment — impact-negligible → report only
- The git-identity fix pins one instance and adds no guard, where the shell fix got a whole-file scanner. The twelfth test now calls `configureStateGit` and the class is empirically closed today, both packages pass under `GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null`. But the outlier was found with an awk scan that was never shipped, and nothing stops a thirteenth committing test from repeating it on a machine where git cannot guess an address from the hostname. `newStateRepository` (`state_apply_test.go:706`) creates every one of these fixture repositories, so moving the two config calls into it would delete twelve call sites and close the class — impact-negligible → report only
- The duration-log item is correctly out of scope but has now been carried twice with no destination. `_dev/tests/test-duration-log.sh:12-29` seeds the header then validates, and `record_test_duration` appends with `>>`, creating a headerless file if it runs first. It produces a false red rather than the false green this request targets, and deferring it was right. It was written "report only" in the first hand-back and copied unchanged into the second, and no queue entry exists. The reachable path is narrower than the note suggests, since the seeding is `printf` then `ln`, which is atomic create-or-fail — impact-negligible → report only
- The release owed by this change was correctly deferred to the work action's finalize transaction, but no number or file set was stated anywhere. From the current state it is 0.304.1, a patch, since all three fixes are repairs with no new capability and no breaking change. Five files: `VERSION`, `skills/do-work/VERSION`, the `**Current version**:` line at `skills/do-work/actions/version.md:5`, a new top entry in `CHANGELOG.md`, and `skills/do-work/CHANGELOG.md` copied byte-identical. Not `skills/do-work-board/tools/queue-kanban/VERSION` (0.236.20, independently versioned and untouched). All five mirrors read 0.304.0 right now and the two changelogs are byte-identical, so nothing is broken, the release is simply still owed — impact-negligible → report only

**Nit findings:**

- The same `find ... -print -quit | grep -q .` ambiguity the remediation converted still exists, unnamed, in `_dev/tests/prescribed-shell-cases/atomic-download.sh` lines 14, 63 and 110. It is in neither hand-back nor the Scope deferral — impact-negligible → report only

**Requirements Checklist**

- [x] A lane that ran and exited non-zero is reported red, whatever it printed — delivered. The rule sits in the record writer, so the stored record and the finding agree. Verified by all three reviewers with an independent revert. Timeout (124) and interrupt (128+signal) both produce non-zero statuses, so both clear the skip. Real skips are unaffected, every shipped script that prints `SKIP:` exits 0 on that path.
- [x] A test's own fixture output cannot be mistaken for the enclosing lane's announcement — delivered as code. Every in-package caller now reaches the runner with a discarding writer, `io.Discard` at `heavy_run_test.go:85-92` or a nil writer that the tee already discards at `heavy_run.go:447`. The only `os.Stderr` left is the real CLI path. It carries no revert-failing assertion, which is counted against criterion 5 rather than here.
- [x] Every test in `state_apply_test.go` that commits sets its own identity — delivered. The one outlier calls `configureStateGit` at line 1049. A per-function scan found no other committing test in the package without it, and both `heavyverification` and `requeststate` pass under the lane's stripped git configuration.
- [x] Both matchers in `update-script-behavior.sh` report grep's verdict, never the writer's exit status — delivered. Both are herestring redirects with no writer process. Restoring the pipeline form and running the shipped self-check failed 5 times out of 5, so the lock-in is deterministic and not a replayed 2-in-500 race.
- [ ] Each fix carries an assertion that fails when the fix is reverted — partially delivered. Four fixes have one and each was reverted individually: the skip rule, the git identity, the two matchers, and the zero lane timeout. The nine converted pipelines are covered by the source scanner, with nine individual ablations recorded. The writer seam has none, and its substitute evidence returns the same result fixed or unfixed.

**Acceptance:** Partial — the three declared defects are fixed and verified end to end (the updater heavy lane exits 0 in 77s with zero FAIL lines, `gofmt -l` empty, `go vet` clean, `shellcheck --severity=warning` exit 0, both Go packages green under the lane's own stripped git configuration), but a lane run with an empty `LaneIDs` still returns a successful heavy verdict with zero lanes, and one fix's proof does not distinguish fixed from unfixed.

**Suggested testing:** 6 items
- Regression: call `RunLanes` with `LaneIDs` unset and assert it refuses, the same shape as the `LaneTimeout` guard the remediation added.
- Regression: add a scanner probe file with the line-continuation pipeline, `grep --quiet`, `grep --silent` and a prefixed `grep -q`, and assert each is flagged before widening the regex.
- Edge case: assert both marker archives readable the way the fallback archive already is at line 619, then re-run the truncated-archive probe and confirm it fails.
- Environment: run the heavy tier on a machine where git cannot derive an identity from the hostname, to confirm the git-identity class is closed in practice and not only under the lane's stripped configuration.
- Integration: run the full canonical gate, not only the updater lane. No reviewer started one.
- Manual: read `prime-do-work-cli.md` trap `[family: silent-skip-reads-as-red]` next to `heavy_run.go:322-333` and confirm the trap states the new two-part rule.

**Follow-ups created:** None (15 findings report only)

*Reviewed by review-work action*
