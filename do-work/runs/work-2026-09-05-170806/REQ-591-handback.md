# Hand-back — REQ-591: reduce repeated setup and unaffected reruns in the fast gate

Branch `worktree-agent-REQ-591-reduce-repeated-setup-and-unaffected-reruns-in-the-fast-gate`,
five commits on top of `e0bdf8bf`. The first four are the original build; the fifth is the
remediation described in `## Remediation`.

| Commit | Task | What it does |
|---|---|---|
| `9ffb8745` | T1 | one shared skill root in the SessionStart hook probe |
| `415c8449` | T2 | the fast-stage evidence engine and its three commands |
| `61d08136` | T3+T4+T5 | the manifest, the gate wiring, the end-to-end probe |
| `f104a570` | T4 fix | SHLVL excluded from the decision child's sealed environment |
| `7fb039b0` | T5 fix | the reuse probe stops inheriting `DO_WORK_FAST_STAGE_REUSE`, and its failure message names the line it wanted |

## File manifest

| Path | Change | Why |
|---|---|---|
| `_dev/tests/session-start-hook-behavior.sh` | modified | nine per-case copies of the CLI module replaced by one shared skill root; banner input becomes a required argument of every case helper; one new assertion that the shared tree is unchanged after every case |
| `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` | new | working-tree seal, fingerprint, evidence store, reuse decision, and the three request entry points |
| `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go` | new | 26-case decision table, expiry, recording refusals, manifest strictness |
| `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` | modified | `decide-fast-stage`, `record-fast-stage`, `invalidate-fast-stage` registered beside the lane commands, with one shared argument parser |
| `_dev/tests/fast-stages.json` | new | the two reusable stages, their coverage, their toolchain probes |
| `_dev/tests/maintainer-verify.sh` | modified | `run_stage_with_evidence`, the two stage call sites, the environment exclusions, the gate wall-time line |
| `_dev/tests/fast-stage-reuse-behavior.sh` | new | nine-case end-to-end probe of the shipped wrapper |
| `_dev/tests/contracts/probe-lanes.sh` | modified | registers the new probe in the fast tier |

`git diff --name-only e0bdf8bf` returns exactly these eight paths, all inside the declared `## Scope`.
Verified explicitly. **`skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` was declared in
scope and is NOT touched**: `heavyverification.Handlers()` is already registered there
(`main.go:62`), so the three new commands need no registration edit.

## Measured Evidence

The builder left this section as a placeholder. The orchestrator took the whole-gate comparison
instead, and the remediation builder took the one measurement still missing (the SessionStart probe
in isolation) plus the post-remediation gate runs. Every number below names the load it was taken
under or the window it was gated on.

### Whole gate, four conditions (orchestrator)

Three interleaved repetitions per condition, `GOMAXPROCS=4`, every run gated on a checked quiet
window with zero concurrent gates and the load recorded on both sides. Wall / user / sys seconds.

| rep | baseline (`e0bdf8bf`) | branch, reuse off | branch, cold store | branch, warm store |
|---|---|---|---|---|
| 1 | 101.05 / 99.82 / 99.95 exit 0 | 20.15 exit **1** | 93.19 exit **1** | 72.09 (one stage reused) |
| 2 | 96.29 / 89.93 / 101.28 exit 0 | 20.48 exit **1** | 93.01 / 91.15 / 101.88 exit 0 | **21.33 / 13.27 / 9.24 exit 0** |
| 3 | 96.19 / 88.35 / 96.98 exit 0 | 20.23 exit **1** | 93.69 / 91.29 / 101.40 exit 0 | **21.16 / 13.29 / 9.22 exit 0** |

Reading: baseline about 96.3 s. Cold store about 93.3 s, which is T1's setup saving net of the new
probe's own 4 s inside the parallel batch. Warm store about 21.2 s with both Go stages reused. The
whole reuse-off column and rep 1 of the cold-store column are the two defects, both handled in
`## Remediation` below.

### Whole gate after the remediation (`7fb039b0`)

Same worktree, load average 3.8 to 5.2 across the three runs, no other gate process running, exit
status read directly from `$?` and never through a pipe. Wall / user / sys seconds from
`/usr/bin/time -p`.

| condition | wall | user | sys | exit | per-stage lines |
|---|---|---|---|---|---|
| variable unset, cold store | 90.48 | 91.62 | 112.63 | **0** | both stages `EXECUTING (no_prior_evidence)` |
| `DO_WORK_FAST_STAGE_REUSE=off` | 90.63 | 91.65 | 113.00 | **0** | none, and that is correct: the wrapper short-circuits before it decides anything, which is what makes this the forced-execution control |
| variable unset, warm store | 22.26 | 13.29 | 9.78 | **0** | both stages `REUSED (fingerprint_match, recorded …; per-file budget verdict inherited from that run)` |

The reuse-off run is the condition that used to die at 20 s. It now runs the whole gate and the new
probe passes inside it (`Fast-stage evidence reuse probes passed.`, 4 s of its 30 s budget). The
third run shows reuse still fires end to end after the fix, so the remediation did not buy the
control condition by disabling the feature.

### The SessionStart hook probe in isolation (the measurement that was missing)

`_dev/tests/session-start-hook-behavior.sh` timed alone with `/usr/bin/time -p`, four repetitions
interleaved between two worktrees of this repository so a drifting load hits both sides equally:
base `e0bdf8bf` and branch `f104a570`. `DO_WORK_TEST_DURATION_LOG` pointed at a scratch file, so no
tracked duration log moved. Both sides exited 0 every time with the same final line,
`SessionStart hook behavior probes passed.`

| rep | load before base run | base wall | load before branch run | branch wall |
|---|---|---|---|---|
| 1 | 3.59 | 12.31 | 3.19 | 3.22 |
| 2 | 3.19 | 12.29 | 7.85 | 3.39 |
| 3 | 7.86 | 12.66 | 7.03 | 3.27 |
| 4 | 6.55 | 12.51 | 6.49 | 3.34 |

Mean 12.44 s to 3.31 s, a saving of 9.13 s per probe run, 73 %. The four repetitions span a load
range of 3.19 to 7.86 and the base column moves by 0.37 s across it, so what separates the two
columns is the probe's own work rather than the machine's.

**The CPU columns for this probe are recorded but not relied on, deliberately.** `/usr/bin/time -p`
reported base user 0.28-0.29 / sys 1.22-1.27 against branch user 0.47-0.48 / sys 0.67-0.69, which
does not scale with a 9 s wall difference. The tool does account for a child's CPU in general (a
plain 0.34 s Python loop under the same `bash` invocation reports 0.32 s user), and one
`go tool -n do-work-cli` in a fresh module path measures 1.26 s wall against 3.25 s user + 0.98 s
sys, so the toolchain work is real CPU that this probe's own process-tree accounting does not show.
The reason it does not reach the timed tree was not resolved here and is named as unavailable
rather than guessed at. Wall time is the honest comparison for this probe, and the whole-gate
tables above carry process-tree CPU for the gate as a whole.

## Remediation

Two defects were found on the merged tree. The first is reproducible and is fixed on the branch. The
second is an intermittent in a test this REQ never edited, and the finding below is what was ruled
out and why, not a claim that it is "probably flaky".

### Defect 1 — the reuse probe read the switch it exists to exercise. FIXED, commit `7fb039b0`

`_dev/tests/fast-stage-reuse-behavior.sh` sources the gate and calls the shipped
`run_stage_with_evidence` directly. That function short-circuits to a plain stage run whenever
`DO_WORK_FAST_STAGE_REUSE` is `off` (`_dev/tests/maintainer-verify.sh:130`), and the probe neither
set nor cleared that name, so it inherited whatever the caller had. D-11 makes
`DO_WORK_FAST_STAGE_REUSE=off` the forced-execution switch the measurement protocol runs the whole
gate under. In exactly the control condition this REQ's acceptance depends on, the probe therefore
switched off the wrapper it was testing, every one of its nine cases failed, and the gate died at
about 20 s. Reproduced three times out of three.

**Fix.** The probe clears `DO_WORK_FAST_STAGE_REUSE`, and `MAINTAINER_VERIFY_SELFTEST_LOG` which
short-circuits the same wrapper, before it sources the gate. The comment at the call site states the
rule rather than the instance: a test of the wrapper must not read the caller's preference for the
wrapper. The two names are cleared, not pinned to a value, so the sealed environment the decision
child fingerprints stays consistent across the probe's own runs.

**Proof.** The probe alone exits 0 in all three environments (name unset, `DO_WORK_FAST_STAGE_REUSE=off`,
`MAINTAINER_VERIFY_SELFTEST_LOG` set). The whole gate exits 0 both with the name unset and with it
set to `off`, both statuses read directly from `$?`; the numbers are in `## Measured Evidence`.

**The failure message was a second, smaller defect.** It printed `ran=` and `status=`, which both
matched, beside `output=< >`. Three fields are asserted and only the output was wrong, so the
message showed two passing fields and an empty one and read as no failure at all. It now names the
line it wanted. Verified by deliberately breaking one expectation:

```
FAIL: unchanged inputs reuse ran=no (want no) status=0 (want 0) output=<maintainer-verify: stage
alpha-stage: REUSED (fingerprint_match, recorded 2026-09-05T22:05:57Z; per-file budget verdict
inherited from that run) > want-line=<DELIBERATELY WRONG LINE >
```

### Defect 2 — `TestLaneMutationCannotPublishOrReuseSuccess/commit=true`. NOT caused by this REQ

The failure, once in four full-gate runs, at `heavy_reuse_regression_test.go:174`: the run refused
with `HEAVY-RUN-DIRTY-TREE` naming `beta/fixture.txt` where the case wanted
`HEAVY-RUN-REVISION-CHANGED`. The orchestrator's differential already showed it passing in 6
isolated runs at the base revision, 6 at the merge revision, and twice for the whole package under
`-short` at the merge revision, so it appears only under a full gate's concurrency.

**What the failure's own shape says.** `verifyLaneRevision` checks the dirty tree *before* it checks
the revision (`heavy_run.go:188-200`). Reaching the dirty-tree refusal therefore means that at that
moment HEAD did not contain the lane's write. The lane script is three steps:
`echo broken > beta/fixture.txt`, then `git add`, then `git commit`. The write is a shell redirect
that needs no new process; both git steps must start one. So the only state that produces this exact
refusal is "the redirect landed and the lane's git steps did not complete". Demonstrated directly on
a copy of the fixture, with git made unreachable to the lane:

```
lane.sh: line 2: git: command not found
HEAD unchanged: yes
porcelain: [ M beta/fixture.txt]
```

That is the observed failure, produced with no evidence store involved at all. This machine's
`kern.maxprocperuid` is 3333, and a full gate plus the sibling sessions sharing this checkout spawn
short-lived processes in the thousands, which is the kind of condition that makes a `git` start fail
where an isolated run never does.

**What was checked in this REQ's own additions, and found not to interact.**

- **No shared parallelism.** `internal/heavyverification` contains no `t.Parallel()` at all, before
  or after this change, so the new fast-stage tests cannot run beside the lane tests inside the test
  binary. They run before them (file order) and finish before them.
- **No shared process state.** The package contains no `os.Setenv` and no `os.Chdir`. The two
  environment writes the new tests make are `t.Setenv`, which Go restores at the end of the subtest
  and refuses to combine with parallel tests.
- **No shared filesystem.** Every root in the new tests comes from `t.TempDir()`; the fixture
  template is built once per test function and copied with `os.CopyFS` into another `t.TempDir()`.
  Nothing in them resolves to this checkout or to another test's repository.
- **The two key spaces cannot touch, which is the specific interaction D-15 invites the question
  about.** The lane store is `<git-common-dir>/do-work-heavy-lanes/` keyed by `sha256(lane-id)`; the
  fast-stage store is `<git-common-dir>/do-work-fast-stages/` keyed by
  `sha256(stage-id + NUL + working-tree-root)`. Different directory, different digest input. More to
  the point, the "common directory" in the lane tests is their own `t.TempDir()` repository, never
  this checkout's, so the two stores are not even in the same tree during a test run.
- **The engine cannot dirty a tree or take a Git lock.** Its only three Git calls are
  `ls-files --cached`, `ls-files --others` and `rev-parse --git-common-dir`
  (`fast_stage_evidence.go:186,218,351`), all read-only, none of which writes the index, and all
  bound to the repository root the caller names.
- **The gate wiring adds no concurrency to that stage.** The decision runs before the stage and the
  recording after it, as separate serial child processes, and the stage itself receives the identical
  argv and budget runner it received before. The gate's stages are strictly serial, so no fast-stage
  command is ever in flight while `go test` runs.
- **The new shell probe runs elsewhere.** It runs Git only inside its own `mktemp -d` fixture and
  invokes the CLI with `--repo-root` pointing at that fixture, and it runs in the
  contract-regressions stage, which is a different serial stage from the two Go test stages.

**The one honest connection.** This REQ adds a test file to the same package, so the
`heavyverification` binary now spawns more Git subprocesses and lives longer, which raises a full
gate's total process pressure slightly. That is the only channel by which the change could make a
pre-existing spawn-failure more likely, and it is not a defect in the reuse logic. It is recorded as
a discovered task rather than hidden here.

**Not reproduced during this remediation.** The whole module suite (`go -C … test -count=1 ./...`,
no `-short`, 30 packages) exited 0, and three full gate runs exited 0.

## P-A-U

- **[PLAN]** — done by the Plan agent; this build followed its task order (T1 first, then T2 tests
  first, T3, T4, T5) and corrected three of its assumptions from evidence, recorded as D-12, D-13
  and D-14 below.
- **[APPLY]** — done. Five commits, eight files, all inside `## Scope`.
- **[UNIFY]** — done. Full gate, module tests, `go vet`, `gofmt -l`, ShellCheck on every shell file
  touched or added, plus eleven verified mutations (seven on the Go engine, four on the gate
  wrapper) proving the negative cases are live rather than decorative. Redone after the remediation:
  three full gate runs (variable unset, `=off`, warm store) all exit 0, `go -C … test -count=1 ./...`
  exit 0 across 30 packages, `go vet` and `gofmt -l` clean, `shellcheck --severity=warning` clean on
  the one shell file the remediation touched, and `git diff --name-only e0bdf8bf` still returns
  exactly the eight scoped paths.

## Decisions

D-01 through D-11 were made in the plan and are unchanged, with one exception noted at D-17.

**D-08 (from the plan) — carried forward, still ESCALATE, not downgraded.** A reused stage writes
no per-file duration rows and enforces no per-file budget for that run; it inherits the budget
verdict of the run whose evidence it reuses. Implemented as recommended, and the REUSED line prints
the inherited verdict:

```
maintainer-verify: stage do-work-cli-fast-tests: REUSED (fingerprint_match, recorded 2026-09-05T21:08:52Z; per-file budget verdict inherited from that run)
```

*Value:* the gate stops failing on other sessions' load for a stage whose inputs are provably
unchanged — five gate runs during this work run failed on exactly that. *Risk:* a real contention
problem goes unreported for up to four hours on a matching tree. *Reversal:* delete that stage's
`fingerprint` block in `_dev/tests/fast-stages.json` and it executes every time, with no code
change. It worked in practice; that is not a reason to reclassify it.

**D-12 — the `do-work-cli` stage covers the board module too, so the plan's board-only reuse case is
not achievable here. DECIDE & STATE, evidence-driven.** The plan and the exploration assumed the two
modules' fast tests are input-independent. They are not:
`internal/publication/capture_files_test.go:15-18` runs `go run . next-req` **inside**
`skills/do-work-board/tools/queue-kanban`, and `internal/doctor/doctor_commands_test.go:159-161,187`
reads `verify.go`, `model.go` and `lessons-do-kanban.md` from that module. Neither is behind
`testing.Short()`, so both run in the fast `-short` stage. Declaring the CLI stage as covering only
its own module would have been a false green for exactly the case requirement 3 names. The board
module's tests were checked in the other direction and read nothing outside their own subtree, so
the honest separation is the reverse of the plan's: **a change confined to the do-work-cli module
leaves the board stage reusable**, and that is the direction the test table asserts. A board change
executes both stages.

**D-13 — the manifest declares only `do-work/` as non-stage coverage; `_dev/tests` and everything
else stay unclassified. DECIDE & STATE.** The plan's T3 would have declared `_dev/tests` as coverage
on both stages and then enumerated "coverage extras" — the maintainer-tree paths each module's tests
read — with a drift-guard test to keep the enumeration honest. Leaving those paths unclassified
reaches the same seal by construction: an unclassified path is sealed into *every* stage, so
`_dev/tests`, `skills/do-work/actions`, `README.md`, `VERSION` and the rest already force both
stages. That deletes the enumeration, the drift-guard test, and the possibility of the enumeration
going stale, and it is strictly more conservative. The only declarations that still matter are the
cross-module ones, which is what D-12 is about.

**D-14 — the gate/manifest agreement is checked at runtime by the commands, not by a test.
DECIDE & STATE.** The plan's T4.5 put the pinning assertions in
`heavy_maintainer_tree_test.go`, which is outside this REQ's `## Scope` and is the single
`.gitattributes` export-ignored Go test file, so the `shipped-module-test-self-containment` lesson
forbids putting a maintainer-tree read anywhere else. Instead all three commands take the caller's
own argv after `--` and refuse when it differs from the manifest's declared argv. That is checked on
every gate run rather than once per test run, and it made the gate/manifest disagreement in the
`## Discovered Tasks` note below visible in five minutes.

**D-15 — a fast-stage record is keyed by stage id AND working-tree root. DECIDE & STATE.** The
evidence store lives in the Git **common** directory, which every linked worktree shares. The heavy
lane can key on stage id alone because it refuses a dirty tree and attributes its result to a
revision. A fast green belongs to a working tree. Without the working-tree component two sibling
worktrees — the normal state of this repository — would revoke each other's records on every run,
because execution invalidates before it runs. The record carries `working_tree_root` and a mismatch
reads as `evidence_unusable`; mutation-verified (M7).

**D-16 — the fast record carries no revision, and the REUSED line names only the recorded
timestamp. DECIDE & STATE.** The plan's line format ended with an evidence revision. A fast green is
computed on a possibly-dirty tree; naming a commit beside it asserts something the record cannot
support. The decision line is four fields:
`<disposition> <reason> <fingerprint-or-dash> <recorded-at-or-dash>`.

**D-17 — `SHLVL` and `OLDPWD` join the three names D-07 removes from the decision child.
DECIDE & STATE, same class as D-07, earned by measurement.** With the plan's three exclusions, reuse
fired between two identically-nested shells and never fired between a terminal and a wrapper script:
the only differing sealed value was `SHLVL`, the caller's shell nesting depth. `OLDPWD` is the same
shape. Both are shell breadcrumbs that decide no assertion. `DO_WORK_TEST_ENFORCE_BUDGET` and
`DO_WORK_TEST_FILE_BUDGET_SECONDS` remain sealed, as D-07 requires. The exclusion is one greppable
line at the call site with the reason above it.

**D-18 — the gate's summary line reports wall time only, not executed/reused counts.
DECIDE & STATE.** The board stage runs inside a subshell that contains its `QUEUE_KANBAN_*` selector
assignments, and that subshell has to stay so the decision child inherits the selectors. A counter
incremented inside it is lost, so counts would cost a temporary file and an extra trap for
information the per-stage `EXECUTING`/`REUSED` lines already carry one line above.

**D-19 — `decide-fast-stage` has no `--format json` variant. DECIDE & STATE.** The plan proposed one.
Its only consumer is a shell that reads one line; the typed twin would be a second contract with no
reader.

**D-20 — the three commands live in `heavy_commands.go`, not a new `fast_stage_commands.go`.
DECIDE & STATE.** The plan proposed two extra files; neither is in `## Scope`, and the commands are
six short handlers plus one shared parser beside the lane commands they are a sibling of.

**D-21 — the probe-timeout counter-case from the plan's test table is not implemented.
DECIDE & STATE.** It would spend five seconds of every gate run proving the behaviour of
`runFingerprintProbe`, which is shared code the heavy lane already covers. The probe-cannot-run case
is covered (`toolchain probe cannot run` → `fingerprint_uncertain`), and it exercises the same
error-to-uncertain path.

## Discovered Tasks

- **A pre-existing intermittent in a test this REQ never edited.**
  `TestLaneMutationCannotPublishOrReuseSuccess/commit=true`
  (`heavy_reuse_regression_test.go:152`) failed once in four full-gate runs with
  `HEAVY-RUN-DIRTY-TREE` naming `beta/fixture.txt` where it wanted `HEAVY-RUN-REVISION-CHANGED`. It
  passes in isolation at both the base and the merge revision, and under `-short` for the whole
  package, so it needs a full gate's concurrency to appear. The dirty-tree check runs before the
  revision check, so the refusal proves the lane's `git add && git commit` did not complete while
  the shell redirect that dirties the file did; the redirect needs no new process and both Git steps
  do. `## Remediation` lists what was ruled out in this REQ's own additions (no `t.Parallel()`, no
  `os.Setenv`/`os.Chdir`, `t.TempDir()` roots only, separate evidence directory and digest key,
  read-only Git calls, serial gate stages). Two concrete next steps: the test passes no
  `LaneOutputWriter`, so `laneSkipWatcher` scans the lane's output for a SKIP line and discards it
  (`heavy_run.go:425`) and the lane's own error text is lost exactly when it is needed, which is
  worth changing before the next hunt; and this host's `kern.maxprocperuid` is 3333, which a gate
  run beside sibling sessions can approach. → report only
- This REQ's own contribution to that pressure, stated rather than hidden: the new
  `fast_stage_evidence_test.go` lives in the same package, so the `heavyverification` binary spawns
  more Git subprocesses and runs longer during a gate. It adds no concurrency (the package has no
  `t.Parallel()`), only load. → report only
- A defect this work introduced and fixed inside its own scope, recorded because the shape is worth
  knowing: `run_stage_with_evidence` derives the manifest's module token by stripping `$repo_root/`
  from the module directory. When that strip fails, every run's argv comparison is refused and reuse
  is silently off forever behind a green gate. The wrapper now stops with a named error instead.
  → report only
- `internal/heavyverification` now holds two near-identical private-record readers (the
  Lstat/permission/SameFile/read dance, ~35 lines) and two copies of the whole-environment seal (~20
  lines). Extracting both needs `heavy_evidence.go`, which was outside this REQ's `## Scope`.
  → report only
- The package name `heavyverification` under-describes its contents now that it also owns the fast
  tier's evidence. A mechanical rename, no behaviour. → report only
- ShellCheck lint is the next honest reuse candidate: 3.7 s over a closed 92-file input set, all of
  it declarable. It is 3.6 % of the gate against three more CLI round trips, which is why it was
  named rather than taken. → report only
- `_dev/tests/contracts/core-checks.sh` builds six Git repositories inline and runs **serially**
  before the parallel probe batch, so its 5.25 s is on the lane's critical path twice over. A
  template-and-copy fixture would recover roughly 2 s. → report only
- `internal/lifecycleadvance` (13 `git init` sites, 41 s of file time) and `internal/heavyverification`
  (three repository builders, ~49 s) are the two densest remaining REQ-574-shaped targets, both with
  zero `t.Parallel()`. Deliberately left alone here so the before/after comparison could separate
  setup savings from avoided execution. → report only

## Lesson evidence

Three entries are earned. None of the satellite files, nor `do-work/lessons-index.md`, is inside
this REQ's `## Scope` and `do-work/` is on the never-touch list, so they are written here for the
orchestrator to land rather than committed on the branch.

**For `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, extending the existing
`fixture-cost-is-subprocess-spawning` family:**

> **[family: fixture-cost-is-subprocess-spawning] REQ-591:** when a shell fixture copies a Go module
> per case, **the copy is not the cost — the new absolute path is**. The launcher resolves its module
> directory physically and then asks the Go toolchain for the tool binary, and a distinct module
> directory is a distinct build-cache action id, so every copied root pays a fresh link of a
> byte-identical binary. Nine copies in `session-start-hook-behavior.sh` cost 0.95 s of `cp` and
> about 10.5 s of relinking. Share one physical root and give each case only the state it mutates.
> The contrast that proves it: six Go test packages build the same command from the *same* source
> directory and their warm repeats cost 0.33 s. Measure the copy before rewriting it.

**For `_dev/primes/lessons-shell-commands.md`:**

> **[family: environment-seal-churn] REQ-591:** a fingerprint that seals the whole environment — the
> only honest choice, since a hand-maintained subset cannot prove an omitted selector stayed
> unchanged — is defeated by shell bookkeeping the caller never chose. `SHLVL` differs between a
> terminal and a wrapper script and `OLDPWD` differs with the caller's previous directory; either
> alone makes reuse never fire, silently, with a green gate. Exclude such names explicitly at the
> call site with the reason, keep every name that changes the verdict, and **prove reuse fires
> across two differently-nested shells** before believing a caching layer works.

> **[family: silent-capability-disable] REQ-591:** a fail-closed cache that cannot fire is
> indistinguishable from one that is working, because both produce a green gate. Every derived
> identifier a fail-closed comparison depends on — here a repository-relative module path from a
> prefix strip — needs a check that stops rather than falls through, or the feature is off forever
> and nothing says so.

The `[family: environment-seal-churn]` and `[family: silent-capability-disable]` families are new,
so no promotion into a prime's `## Traps` is mandatory yet; both are first occurrences. The
`fixture-cost-is-subprocess-spawning` family is now on its second REQ, which under
`crew-members/general.md` § Lessons Discipline makes promotion of one generalized line mandatory —
the generalized form is *"a fixture that reproduces a build input at a new path pays for the path,
not the copy"*.

## Integration seams

- **`_dev/tests/fast-stages.json` is the single declaration point.** Adding a reusable stage means
  one manifest entry plus one `run_stage_with_evidence <id> <module-dir> <args...>` call whose
  derived argv tokens match the manifest exactly, or the command refuses and the stage executes.
- **The evidence store** is `<git-common-dir>/do-work-fast-stages/`, 0700, records 0600, keyed by
  `sha256(stage-id + NUL + working-tree-root)`. It is deliberately a different directory, schema
  version and key space from `do-work-heavy-lanes/`; the two must never cross.
- **`DO_WORK_FAST_STAGE_REUSE=off`** is read only by `_dev/tests/maintainer-verify.sh`, which is
  export-ignored, so no shipped script gained an option. It is the forced-execution switch the
  measurement protocol needs.
- **The fast tier now registers twelve probes** rather than eleven; the new one finishes in about
  4 s inside the parallel batch and about 2.2 s alone.
- **The self-test bypass** is the same `MAINTAINER_VERIFY_SELFTEST_LOG` seam `run_budgeted_go_tests`
  already uses. `--self-test` still asserts nine stages exactly once, unchanged. A future change
  that removed the bypass would send a decision child into the self-test's command shim, whose
  `bash` case accepts only `contract-regressions.sh` and exits 64, so the fixture fails loudly.
- **`heavyverification.Handlers()`** is the only registration point the three commands needed;
  `main.go` already iterates it.
