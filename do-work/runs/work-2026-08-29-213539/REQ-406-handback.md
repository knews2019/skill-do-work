# REQ-406 Builder Hand-Back

## Branch

`worktree-agent-REQ-406-create-do-work-cli-foundation` (base `18c19f1`)

| Commit | Subject |
|---|---|
| `571bb0d` | Keep git stderr out of porcelain answers and scope the rollback index check |
| `737dd20` | Leave one exit-code authority and one rollback result type |
| `33ff757` | Turn Git transaction failures into complete, actionable findings |
| `84f77c4` | Run the do-work-cli module from the canonical gate |

Nothing is pushed. Nothing outside the worktree was written except this file.

## Implementation Summary

Ten files, exactly the Scope declaration's list. No file outside it was touched — `VERSION`, both changelogs, `actions/version.md`, the REQ file, `do-work-cli.sh` and the launcher probe are all unmodified (`git diff --stat 18c19f1..HEAD` confirms).

| File | Action | What it now does |
|---|---|---|
| `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` | modify | `runGit` captures stdout and stderr separately and returns only stdout; `indexIsEmpty` takes optional pathspecs and `rollbackFailure` scopes its check to the declared targets; `TransactionResult.ExitCode int` is replaced by `Outcome resultmodel.CommandOutcome`; the local `RollbackResult`/`RollbackStatus` types are deleted; `fail` becomes `failTransaction` taking an outcome plus the affected paths; `TransactionFailure` gains `Paths`; `TransactionResult` gains `CreatedPaths`. |
| `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` | modify | Two new lock-in tests (git stderr warning must not read as target dirt; unrelated staged work must not break a complete rollback); all 11 exit-code assertions now read the number through `resultmodel.ExitCode(result.Outcome)`; `newRepository` pins `GIT_CONFIG_GLOBAL`/`GIT_CONFIG_SYSTEM`/`GIT_TERMINAL_PROMPT`. |
| `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go` | new | `BuildCommandResult` maps a `TransactionResult` to the one typed `CommandResult`. `FindingCode` derives the code from the `FailureKind` string; a per-kind table supplies only severity, fixability, stop reason and remediation argv; an unmapped kind yields a complete `GIT-UNMAPPED-FAILURE` error finding. Change kinds are `created`/`modified` from `CreatedPaths`; a completed rollback reports no surviving changes. |
| `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go` | new | Code-derivation table; a completeness assertion over every failure kind, read back out of `git_transaction.go` rather than restated; the loud-fallback test for a synthetic unmapped kind; truthful change kinds; no surviving changes after a completed rollback. |
| `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` | modify | Owns `RollbackStatus` and its three constants, so `RollbackResult.Status` is typed. `ExitCode` carries a comment naming it the single authority. |
| `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` | modify | The parity test's fixture now carries a rollback; it asserts the typed status reaches JSON as a plain string, that `rollback.errors` is non-null, and that the text form names the status and the rollback action. |
| `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` | modify | `usageFinding` now carries a `VerificationArgv`. One line. |
| `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` | modify | A completeness test over the runtime's own findings, plus the four-scenario exit-code contract test that drives real Git fixtures through `Run` with `--repo-root`/`--format` and observes 0/1/3/4 in both renderings. |
| `_dev/tests/contract-regressions.sh` | modify | Registers `_dev/tests/do-work-cli-launcher-behavior.sh`, which was tracked, executable and passing while nothing invoked it. |
| `_dev/tests/maintainer-verify.sh` | modify | `do-work-cli go vet` and `do-work-cli uncached tests` lanes, plus all four self-test enumerations moved in lock-step. |

**What I did not build.** No command is registered — `main.go` is untouched and still delegates to a nil handler map (plan D6). The launcher and its probe are untouched (plan D1). The real binary was exercised by hand end to end but no shell probe execs it (plan D7).

## P-A-U

### [PLAN]

Read the brief in full, `_dev/primes/prime-shell-commands.md`, and the four crew members. Read every preserved file before writing anything. Confirmed the baseline is green at `18c19f1` (`go vet`, `go test -count=1 ./...`, the launcher probe, ShellCheck all pass), so the work is a completion pass, not a rewrite.

Sequenced to keep RED honest and to respect risk R5 (the outcome collapse must land before the bridge that maps it):

1. RED the two Git defects against the current API, then fix them.
2. RED the runtime's incomplete usage finding, then fill it.
3. Collapse the exit-code and rollback-type duplication (compile-driven; the retargeted assertions are the safety net).
4. RED the bridge and the exit-code contract test together, then implement.
5. Wire the gate, then prove both new lanes and the probe registration can fail.

Two deviations from the plan, both stated below as decisions: I skipped the release ritual and the REQ file (brief corrections 1 and 2), and I put the exit-code contract test in `command_runtime_test.go` rather than a new file so the diff stays inside the Scope list.

### [APPLY]

Five commits' worth of work landed as four commits, one per coherent slice. Every changed path was staged explicitly; `git add -A` and `git add .` were never used.

One process error worth recording: partway through the collapse I ran `git checkout -- git_transaction.go` to undo a temporary falsification edit, which discarded the uncommitted collapse work in that file along with it. I re-applied it and switched to scratch-directory backups for every later falsification. No committed work was lost and the final state is verified by the gate.

### [UNIFY]

`git diff --stat 18c19f1..HEAD` — 10 files, 791 insertions, 64 deletions. Working tree clean; nothing uncommitted.

| File checked | What I checked |
|---|---|
| `git_transaction.go` | No `CombinedOutput` remains. `indexIsEmpty` has exactly three call sites: the `--commit` precondition and the post-commit check unscoped, the post-rollback check scoped to the declared targets. No `fail(` leftovers. The only `ExitCode` token in the package is `exec.ExitError.ExitCode()`. The construction names `OutcomeSuccess` explicitly. Re-read all three `targetIsDirty` consumers (preflight, `changedTargets`, `rollbackFailure`) after the stream split, per trap C5. |
| `git_transaction_test.go` | 11 assertions read through `resultmodel.ExitCode`; the two rollback-status assertions read `resultmodel.RollbackSucceeded`/`RollbackIncomplete`. The two new tests are the ones that were RED. `newRepository`'s `t.Setenv` calls are safe — no test in the package is parallel. |
| `transaction_findings.go` | Every emitted literal spells `do-work-cli`, never `do-work` followed by a space (trap C8, the compiled binary is walked by the live-surface scan). All eight templates supply severity, fixability, stop reason, a next argv and a verification argv. `gitPathArgv` omits the `--` separator when there are no paths, so no argv ends in a dangling separator. |
| `transaction_findings_test.go` | The kind list is derived from `git_transaction.go`, not restated. The regex is anchored on `FailureKind = "…"`, which is the constant block's actual shape. |
| `result_model.go` | `RollbackStatus` has exactly one definition in the module; `ExitCode` is the only function mapping an outcome to a number (`grep -rn 'func ExitCode' internal/` returns one line). |
| `result_model_test.go` | The added assertions fail if the status stops reaching the wire as a plain string or the text renderer stops naming it. |
| `command_runtime.go` | One added line. No other behavior touched. |
| `command_runtime_test.go` | The three pre-existing tests are unchanged. The new fixture helpers (`newFixtureRepository`, `commitFixture`, `runFixtureGit`, `writeFixtureFile`) do not collide with anything in the package. The parity assertion builds its expected line with the production renderer rather than joining argv by hand. |
| `contract-regressions.sh` | The block sits with the other `_dev/tests` invocations (before the staged-skills probe), uses the `[ ! -x ]` form the file uses for tracked-executable probes, increments `fail_count` rather than exiting, and its comment names what the probe covers. |
| `maintainer-verify.sh` | All four self-test surfaces moved together: the `*/skills/do-work/tools/do-work-cli)` shim arm asserting exact argv, the fixture `mkdir -p` entry, the success stage list plus both totals (9→11, 10→12), and both stages in the failure-injection loop. The lanes use the subshell-`cd` shape the queue-kanban lanes use, so cwd never leaks. |

Linters, all after the final commit:

- `gofmt -l -- $(git ls-files '*.go')` from `$(go env GOROOT)/bin` — **no output** (judged by emptiness, not status). All new Go files were committed before this ran, so the tracked-file lane actually sees them (trap C7).
- `go vet ./...` in the module — exit 0.
- `shellcheck --severity=warning` over the four shell files — exit 0.
- Debug-artifact scan over every added line in the diff (`fmt.Print`, `println`, `TODO`, `FIXME`, `t.Skip`, `set -x`, `console.log`) — no matches.

Real-binary smoke check: deleted the built binary, ran `bash skills/do-work/tools/do-work-cli.sh --format json inspect` — the launcher rebuilt on demand and the binary emitted the typed JSON envelope and exited 2; the same invocation in text mode produced the same finding code, next step and verification command and also exited 2.

## Testing

All commands from the worktree root. Every exit code below is the command's own direct status; nothing was piped through `tail` or `head` when the status mattered.

### Task 2 — the two Git defects (RED → GREEN)

`cd skills/do-work/tools/do-work-cli && go test -count=1 -run 'TestGitStderrWarningIsNotReadAsTargetDirt|TestUnrelatedStagedWorkDoesNotBreakRollback' ./internal/gittransaction/`

RED, exit 1, before touching `git_transaction.go`:

```
--- FAIL: TestGitStderrWarningIsNotReadAsTargetDirt (0.02s)
    git_transaction_test.go:217: git stderr warning was treated as target dirt: result = gittransaction.TransactionResult{ExitCode:1, ... Failure:(*gittransaction.TransactionFailure)(0x...)}, called = false
--- FAIL: TestUnrelatedStagedWorkDoesNotBreakRollback (0.03s)
    git_transaction_test.go:244: unrelated staged work broke a complete rollback: result = gittransaction.TransactionResult{ExitCode:4, ... Rollback:gittransaction.RollbackResult{Status:"incomplete", Actions:[]string{"restored tracked.txt from HEAD"}, Errors:[]string{"Git index is not empty after rollback"}} ...}
FAIL
```

Both failures are the assertions, not build errors, and both reproduce the exact symptoms the exploration recorded: exit 1 with `dirty_target` on a clean target, and exit 4 with a `succeeded`-in-substance rollback reported as `incomplete`.

GREEN, exit 0:

```
--- PASS: TestGitStderrWarningIsNotReadAsTargetDirt (0.02s)
--- PASS: TestUnrelatedStagedWorkDoesNotBreakRollback (0.03s)
ok  	.../internal/gittransaction	0.056s
```

### Task 1 — the zero-value flip (falsification)

Deleting `Outcome: resultmodel.OutcomeSuccess` from the result construction (trap C3) reddens two pre-existing tests, so the retargeted assertions do guard it:

```
--- FAIL: TestDirtyTargetIsRefusedButUnrelatedDirtIsAllowed (0.03s)
    git_transaction_test.go:55: unrelated dirt result = ...{Outcome:"", ...}
--- FAIL: TestDryRunDoesNotMutateAndCannotCommit (0.01s)
    git_transaction_test.go:71: dry-run result = ...{Outcome:"", ...}, called = false
```

Restored; `go test -count=1 ./...` exit 0.

### Task 3 — the requirement-5 hole in the runtime (RED → GREEN)

`go test -count=1 -run TestRuntimeFindingsCarryCompleteRemediation ./internal/commandruntime/`

RED, exit 1, four cases, each an assertion failure:

```
command_runtime_test.go:120: Run(["--format=json" "unknown"]) finding names no verification command: resultmodel.CommandFinding{Code:"UNKNOWN-COMMAND", ..., VerificationArgv:[]string{}}
command_runtime_test.go:120: Run(["--format=json" "--repo-root"]) finding names no verification command: ...Code:"MISSING-OPTION-VALUE"...
command_runtime_test.go:120: Run(["--format=json" "--unknown" "value"]) finding names no verification command: ...Code:"UNKNOWN-GLOBAL-OPTION"...
command_runtime_test.go:120: Run(["--format=json"]) finding names no verification command: ...Code:"MISSING-COMMAND"...
```

GREEN after the one-line `VerificationArgv` addition, exit 0.

### Tasks 3 and 4 — the bridge and the exit-code contract (RED → GREEN)

RED for the bridge, `go test -count=1 ./internal/gittransaction/`:

```
internal/gittransaction/transaction_findings_test.go:38:13: undefined: FindingCode
internal/gittransaction/transaction_findings_test.go:53:20: undefined: BuildCommandResult
... (7 undefined-symbol errors)
FAIL	.../internal/gittransaction [build failed]
```

RED for the contract test, `go test -count=1 ./internal/commandruntime/`, exit 1:

```
internal/commandruntime/command_runtime_test.go:295:26: undefined: gittransaction.BuildCommandResult
FAIL	.../internal/commandruntime [build failed]
```

Both are build failures rather than assertion failures, because the function under test is what does not exist yet. The plan anticipated this for task 4 ("fails to build or fails on the missing outcomes"). The assertion-level RED for the same requirement is the `TestRuntimeFindingsCarryCompleteRemediation` evidence above.

GREEN after `transaction_findings.go` landed:

```
--- PASS: TestExitCodeContractThroughRealGitTransactions (0.25s)
    --- PASS: .../clean_mutation_succeeds (0.06s)                                    [exit 0]
    --- PASS: .../dirty_target_is_safely_refused (0.05s)                             [exit 1, GIT-DIRTY-TARGET]
    --- PASS: .../mutation_failure_rolls_back_cleanly (0.06s)                        [exit 3, GIT-MUTATION-FAILED]
    --- PASS: .../post-commit_verification_failure_reports_committed-state_risk (0.08s) [exit 4, GIT-COMMITTED-STATE-RISK]
--- PASS: TestRuntimeFindingsCarryCompleteRemediation (0.00s)
--- PASS: TestFindingCodeIsDerivedFromTheFailureKind (0.00s)
--- PASS: TestEveryDeclaredFailureKindProducesACompleteFinding (0.00s)
--- PASS: TestUnmappedFailureKindFailsLoudly (0.00s)
--- PASS: TestSuccessfulTransactionCarriesTruthfulChangeKinds (0.00s)
--- PASS: TestCompletedRollbackReportsNoSurvivingChanges (0.00s)
```

One intermediate failure was mine, not the code's: the parity assertion joined argv by hand and did not match the renderer's quoting of the `<command>` placeholder. Fixed by building the expected line with `resultmodel.RenderResult`, which is what makes it a parity assertion rather than a second renderer.

### Closed-enumeration falsification

Deleting the `FailureDirtyIndex` template entry (trap C15) — the coverage test names the kind that lost it rather than passing on a plausible-looking finding:

```
--- FAIL: TestEveryDeclaredFailureKindProducesACompleteFinding (0.00s)
    transaction_findings_test.go:59: kind "dirty_index" code = "GIT-UNMAPPED-FAILURE", want "GIT-DIRTY-INDEX"
```

Restored; package green.

### Task 5 — the gate lanes are not decorative

**A. `bash _dev/tests/maintainer-verify.sh --self-test`** — `Maintainer verification self-test passed.`, exit 0 with the new stage counts.

Falsified by deleting the production `do-work-cli go vet` lane:

```
FAIL: maintainer-verify self-test: cli-vet ran 0 times; want exactly once
SELFTEST_EXIT=1
```

Restored; self-test exit 0. The failure-injection loop already runs the whole script once per injected stage including `cli-vet` and `cli-test` and requires a nonzero exit from each, so both lanes are falsifiable in both directions.

**B. `chmod -x _dev/tests/do-work-cli-launcher-behavior.sh` then `bash _dev/tests/contract-regressions.sh`** — exit 1, naming the file and the coverage lost:

```
FAIL: _dev/tests/do-work-cli-launcher-behavior.sh is missing or not executable — the do-work-cli launcher has no behavioral coverage.
```

`chmod +x` restored (mode is back to the tracked 100755; `git diff --stat` on that path is empty), and `bash _dev/tests/contract-regressions.sh` then exits 0 with `do-work-cli launcher behavior tests passed` on line 12 and `Contract regression checks passed.` at the end.

### Final gate

`bash _dev/tests/maintainer-verify.sh` from the worktree root, unpiped, browser lane in its default skipped state:

```
maintainer-verify: aggregate contract suite
...
do-work-cli launcher behavior tests passed
...
Contract regression checks passed.
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	274.875s
maintainer-verify: queue-kanban strict JavaScript behavior lane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (82.12s)
SKIP: no browser is available; strict browser behavior lane was not run. Set QUEUE_KANBAN_BROWSER to name one.
maintainer-verify: do-work-cli go vet
maintainer-verify: do-work-cli uncached tests
?   	.../cmd/do-work-cli	[no test files]
ok  	.../internal/commandruntime	0.252s
ok  	.../internal/gittransaction	0.253s
ok  	.../internal/resultmodel	0.003s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	2.327s
Maintainer verification passed.

real	7m32.399s
GATE_EXIT=0
```

All three new gate surfaces visibly ran. `staged skills contract: PASS` inside the aggregate confirms the compiled binary — which carries every new finding string — does not trip the retired-trigger scan (risk R4, trap C8).

### Acceptance criteria

| Criterion | Evidence |
|---|---|
| One module under the installed core package, Go 1.26.1+, `--repo-root` and `--format text\|json` | `go.mod` unchanged at `go 1.26.1`; both options driven against real Git fixtures in `TestExitCodeContractThroughRealGitTransactions`; gate lanes compile and test the module on every run |
| One typed result renders text and stable JSON with schema version, command, outcome, root, findings, changes, skipped work, rollback | The contract test decodes JSON and asserts `schema_version` plus non-null `findings`/`changes`/`skipped_work`/`rollback`/`outcome`/`repository_root`, then asserts the text rendering names the same code and next step from the same result |
| Every finding carries code, severity, paths, evidence, fixability, stop reason, next argv/recipe, verification command | `TestEveryDeclaredFailureKindProducesACompleteFinding` over every declared kind, plus `TestRuntimeFindingsCarryCompleteRemediation` for the runtime's own findings |
| Exit codes 0–4 from exactly one authority; no second exit-code field survives | `resultmodel.ExitCode` is the only such function in the module; `TransactionResult.ExitCode` is deleted; the numeric table is observed at the `Run` seam for 0/1/3/4 and in `TestOutcomeExitCodes` for all six outcomes |
| Fixtures prove exact-path refusal, `--commit` index refusal, pre-commit rollback, post-commit `git revert <sha>` without rewriting history | The seven pre-existing `git_transaction_test.go` scenarios (all retargeted and passing) plus the two new lock-in tests |
| The launcher builds on demand, and the lanes that run all of this are proven non-decorative | Launcher probe now invoked by `contract-regressions.sh`; both falsifications above; real-binary rebuild-and-run smoke check |

## Integration Seams

These are yours to apply. I wrote none of them.

**S1 — Version bump.** Current is `0.244.9`. Suggested next: `0.245.0` (minor: a new shared Go module and two new gate lanes). Three files: `VERSION`, `skills/do-work/VERSION`, and the `**Current version**:` line in `skills/do-work/actions/version.md`. `0.245.0` is strictly greater than the current first changelog entry.

**S2 — Changelog entry**, top of `CHANGELOG.md` below the header. Suggested wording, title verified not already used:

```markdown
## 0.245.0 — Shared Go Command Runtime for the Suite (2026-08-29)

The suite gets one Go command module underneath it: a typed result that renders identical text and JSON, a single exit-code authority, and a Git transaction layer that refuses to touch work you already have in flight.

- One `do-work-cli` module under the core package, built on demand by its launcher when the binary is missing or older than its sources.
- Every result carries a schema version, findings, changes, skipped work and a rollback result, in text or JSON from the same value, with exit codes 0-4 decided in exactly one place.
- Git mutations refuse a dirty target path while leaving your unrelated changes alone, roll back cleanly when they fail, and report `git revert <sha>` instead of ever rewriting history.
- Two fixed defects: a git warning on stderr no longer reads as a dirty target, and someone else's staged file no longer turns a successful rollback into a reported risk.
- The canonical gate now runs the module's vet, tests, and launcher probe — the probe had been passing for nobody since it was written.
```

**S3 — Changelog mirror.** Copy root `CHANGELOG.md` to `skills/do-work/CHANGELOG.md` byte-identically after S2 lands. `_dev/tests/shipped-package-reference-contract.sh` enforces it.

**S4 — REQ file.** Tick the three P-A-U boxes, record the implementation summary and the evidence above, then record the implementation hash in a separate bookkeeping commit. The implementation is `84f77c4` (the branch tip); if you squash, use the resulting hash.

## Decisions

**D-01 — Skipped the release ritual and the REQ file entirely.** DECIDE & STATE. Brief corrections 1 and 2, plus the Scope declaration, plus `work-reference.md:407` (serial-only files) all agree, and REQ-390 is building in the same wave — a version bump here would race it. The plan's D9 last sentence is the only thing that says otherwise and it loses to the Scope list.

**D-02 — Renamed `fail` to `failTransaction` and gave it the affected paths.** DECIDE & STATE. Trap C16: it is a single-word package-level name with reach, which coding-guardrails §5 forbids for names being introduced. §3 protects existing code, but this REQ rewrites the signature and all fourteen call sites anyway, so the rename costs nothing beyond the edit already required.

**D-03 — Added `TransactionFailure.Paths` rather than parsing the reason prose.** DECIDE & STATE. Requirement 5 wants affected paths on every finding. The only alternative was regexing the quoted path back out of `fmt.Sprintf("target path %q is already dirty", …)`, which is a parser over a message this repo's own convention says never to match on (exploration pattern: failures are classified by kind, never by text).

**D-04 — Left `rollback.status` empty for results that never ran a Git transaction.** DECIDE & STATE. Normalizing it to `not_needed` would make every read-only command's text output print a `rollback: not_needed` line, which reads as though a mutation was possible. The empty value already means "no rollback context" and the text renderer suppresses it. Noted as a discovered task for whoever adds the first read-only command, since a JSON consumer switching on the status must handle `""` plus the three constants.

**D-05 — Put the exit-code contract table test in `command_runtime_test.go`, not a new file.** DECIDE & STATE. I wrote it as `exit_code_contract_test.go` first and then merged it in, because a new file in that package is not on the Scope list and the difference is cosmetic — same package, same compilation unit, no behavioral difference. Staying inside the declared write set was worth more than a tidier file split.

**D-06 — The kind-coverage test reads the failure kinds out of `git_transaction.go` instead of restating them.** DECIDE & STATE. A hand-copied list of eight kinds in the test is the closed enumeration the shell prime warns about, and it would go stale the moment a ninth is added — which is precisely the failure the test exists to catch. The regex is anchored on `FailureKind = "…"`; if that declaration shape ever changes, the test fails loudly with "no FailureKind constants found" rather than silently covering nothing.

**D-07 — Both Go fixtures now pin `GIT_CONFIG_GLOBAL`/`GIT_CONFIG_SYSTEM`/`GIT_TERMINAL_PROMPT`.** DECIDE & STATE. Trap C9: the stderr-warning test's entire premise is controlling what git writes to which stream, and a developer's global `core.autocrlf`, `commit.gpgsign` or `init.templateDir` can add or suppress a warning. It matches what every shell probe in `_dev/tests` already does. Applied in the shared helper, so all nine tests in the package inherit it without their bodies changing.

**D-08 — A completed rollback reports no surviving changes.** DECIDE & STATE. `TransactionResult.ChangedPaths` is populated before the failure that triggers a rollback, so passing it through would describe a worktree that was undone. An incomplete rollback still reports them, because there some may genuinely persist.

**D-09 — Task 5 overlaps REQ-420's stated deliverable; I did the additive half only.** ESCALATE. `do-work/queue/REQ-420-...md:34` and `UR-081/input.md:121` both assign "extend `maintainer-verify.sh` with go vet and uncached do-work-cli tests, retain queue-kanban verification, and replace the separate audit-metrics lane" to REQ-420, the last REQ of the batch. I added the do-work-cli lanes and deliberately did **not** remove the audit-metrics lane, so REQ-420 still has its removal to do and nothing it owns was pre-empted destructively.
*Value:* an unwired module is an unverified module — without these lanes, REQ-406's own requirements 1 and 3 ("one module requiring Go 1.26.1+", "build on demand") hold only by inspection, and the launcher probe would have stayed dead weight through fourteen more REQs.
*Risk:* low and reversible. REQ-420 will find its first clause already satisfied, which reads as a scope surprise if nobody tells it. Worth a line in REQ-420 noting the lanes exist and only the audit-metrics removal remains. If you would rather this land with REQ-420 instead, `84f77c4` reverts cleanly on its own — the other three commits do not depend on it.

**D-10 — Duplicated a small Git fixture helper in the `commandruntime` test package.** DECIDE & STATE. The exploration's "do not introduce a second fixture helper" applies within a package; `gittransaction`'s helpers are unexported and Go gives no way to share them across the package boundary without exporting test-only surface. The duplicate is four short functions and it stays in the test file.

## Discovered Tasks

- **`maintainer-verify.sh`'s self-test EXIT trap reads a function-local.** `trap 'rm -rf -- "$self_test_root"' EXIT` is set inside `run_self_test`, where `self_test_root` is a `local`. On any self-test failure the function returns before `trap - EXIT`, so the trap fires at process exit with the variable out of scope and prints `_dev/tests/maintainer-verify.sh: line 1: self_test_root: unbound variable` after the real `FAIL:` line. Pre-existing, cosmetic, does not mask the failure or change the exit status — I saw it during the cli-vet falsification. The fixture directory is also left behind on that path.
- **REQ-420's first clause is now done.** Its remaining work is replacing the separate audit-metrics lane. See D-09.
- **`resultmodel.RollbackResult.Status` has a fourth wire value, `""`.** A consumer switching on it must handle the empty string alongside the three constants. See D-04 for why normalizing it is not free.
- **The `<command>` placeholder in usage and invalid-options findings is not a runnable argv.** `usageFinding` and the `invalid_options` template both emit `do-work-cli --format text <command>`, which tells a reader the shape but cannot be pasted. Once REQ-407 registers real commands, the runtime knows the command name and can thread it through.
- **`resolveRepositoryRoot`'s stderr exposure is the same class as the fixed `targetIsDirty` bug but was never reproduced.** The exploration could not make `git rev-parse --show-toplevel` emit a warning. It is fixed incidentally by the `runGit` change; there is no lock-in test for it because there is no known way to trigger it. Recorded so nobody assumes it was proved.
