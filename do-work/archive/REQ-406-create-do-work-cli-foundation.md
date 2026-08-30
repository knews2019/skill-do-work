---
id: REQ-406
title: 'Create the shared do-work-cli runtime and Git transaction foundation'
status: completed
created_at: 2026-08-29T20:28:26Z
route: C
write_set: [skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go, skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go, skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go, skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go, skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go, skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go, skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go, skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go, _dev/tests/contract-regressions.sh, _dev/tests/maintainer-verify.sh]
estimate:
  p50_active_minutes: 70
  confidence: low
  calculated_at: 2026-08-29T21:37:00Z
  basis:
    - Route C
    - 11-file write set
    - 9 new files
    - 3 subsystems involved
    - 7 acceptance criteria
    - cross-route regression gates
    - full-suite verification
claimed_at: 2026-08-29T21:35:39Z
completed_at: 2026-08-30T07:11:22Z
commit: 2ca25d7
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Create the Shared Do-Work-CLI Runtime and Git Transaction Foundation

## What
Create the suite-wide `do-work-cli` Go module and its shared execution contracts.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Put one `do-work-cli` module under the installed core package and require Go 1.26.1+.
- Support global `--repo-root` and `--format text|json` options.
- Build on demand when the binary is missing or older than its Go sources, using the standard library unless an existing dependency is demonstrably necessary.
- Define one typed result model that renders both text and stable JSON with schema version, command, outcome, repository root, findings, changes, skipped work, and rollback result.
- Give every finding a stable code, severity, affected IDs/paths, observed evidence, fixability, automation-stop reason, exact next argv/Just recipe, and verification command.
- Enforce exit codes 0–4 exactly as specified by the UR.
- Implement Git target preflight, dry-run support, optional exact-path commit, rollback, committed-path verification, and committed-state risk reporting.

## Constraints
- Mutations require Git; read-only commands remain usable outside Git.
- Dirty target paths are refused while unrelated dirty paths remain allowed; `--commit` requires an empty pre-existing index.
- Pre-commit rollback restores tracked targets and removes only invocation-created paths; post-commit failure reports `git revert <sha>` and never rewrites history.

## Dependencies
Foundation REQ; all later batch REQs depend on it directly or transitively.

## Builder Guidance
Certainty level: Firm. Establish narrow packages and behavioral contracts that later migrations can consume without parallel implementations.

## Red-Green Proof
**RED prompt/case:** Invoke the absent `do-work-cli` in text and JSON modes and exercise a mutation against dirty, clean, rollback, and commit fixtures.
**Why RED now:** No shared CLI, typed result schema, build-cache launcher, or common Git transaction layer exists.
**GREEN when:** The same typed result drives stable text/JSON, documented exit codes are observed, and Git transaction fixtures prove exact-path refusal, rollback, and commit behavior.
**Validation:** User confirmed via the supplied implementation plan.

## Partial Work Present (2026-08-30)

Work began before the user clarified that this stage is capture-only. Preserve the following uncommitted files for the future REQ-406 implementation:

- `skills/do-work/tools/do-work-cli/go.mod`
- `skills/do-work/tools/do-work-cli/.gitignore`
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go`
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go`
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go`
- `skills/do-work/tools/do-work-cli.sh`
- `_dev/tests/do-work-cli-launcher-behavior.sh`
- Ignored build output: `skills/do-work/tools/do-work-cli/do-work-cli`

**Evidence already observed:** The launcher fixture first failed because the launcher was absent. Go tests, `go vet`, the output-sensitive `gofmt -l` check, the launcher fixture, ShellCheck, and real launcher text/JSON smoke checks passed during the interrupted work.

**Verification before handoff:** After the final signal-trap ordering adjustment was restored, `go test -count=1 ./...`, `go vet ./...`, the output-sensitive `gofmt -l` check, the launcher fixture, and ShellCheck all passed. The full maintainer gate and REQ-406 acceptance suite have not run.

**State:** Partial, preserved in commit `329c55a9`, and not accepted. This REQ remains pending. Future implementation must inspect the present files and commit, run the remaining RED/GREEN checks and all final gates, and continue from the recorded evidence rather than starting over.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

---

## Triage

**Route: C** - Complex

**Reasoning:** A new suite-wide Go module with a typed result model, exit-code contract, on-demand build launcher, and a Git transaction layer — architectural work every later REQ in the batch consumes. A partial foundation is already preserved in commit `329c55a9` and must be inspected rather than recreated.

**Planning:** Required

## Plan

The partial foundation in commit 329c55a9 is real and currently green: `go vet ./...`, `go test -count=1 ./...`, `_dev/tests/do-work-cli-launcher-behavior.sh`, and ShellCheck all pass right now against the preserved files, and the launcher builds and runs the real binary end to end (verified: `--format json inspect` emits the typed JSON result and exits 2). So this REQ is a completion pass, not a rewrite: the launcher (requirement 3), the global-option parser (requirement 2), the typed result schema (requirement 4), and the Git transaction body (requirement 7) already exist and are tested. Four things are genuinely missing or wrong. First, exit codes have two authorities — `resultmodel.ExitCode` maps outcomes to 0-4 while `gittransaction` hardcodes 1/2/3/4 in ten `fail(...)` calls, and `resultmodel.RollbackResult` is duplicated by `gittransaction.RollbackResult`; collapsing both into resultmodel is a deletion that satisfies requirement 6 and the Builder Guidance ban on parallel implementations. Second, nothing converts a `TransactionResult` into a `CommandResult`, so the Git layer — the one subsystem this REQ owns that produces refusals — emits no findings at all, leaving requirement 5 unserved; a single mapping function with codes derived mechanically from `FailureKind` fixes that without a stale enumeration. Third, two confirmed defects in `git_transaction.go`: `runGit` uses `CombinedOutput()`, so any git stderr warning makes `targetIsDirty` see non-empty output and refuse a clean target (reproduced deterministically with a malformed `.gitattributes` line: stdout 0 bytes, stderr one warning), and `rollbackFailure` checks the whole index for emptiness, so a user's unrelated staged file turns a successful rollback into a spurious exit 4 — both violate the REQ's own constraint that unrelated dirt stays allowed. Fourth, the module is invisible to the canonical gate: `maintainer-verify.sh` has lanes for queue-kanban and audit-metrics but none for do-work-cli, and `contract-regressions.sh` never invokes the launcher probe, which the file's own comment calls dead weight that reads as coverage. Everything else in the preserved commit stays untouched: the launcher's signal-trap ordering, the `find -newer` rebuild gate, and the `.gitignore` that (unlike the root dotfiles) is not export-ignored and therefore ships with the module exactly as queue-kanban's does.

**Tasks:**

1. Collapse the exit-code and rollback-result duplication. In skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go, replace `TransactionResult.ExitCode int` with `Outcome resultmodel.CommandOutcome`, change `fail(...)` to take an outcome instead of an int, and delete every hardcoded 1/2/3/4 (dirty target and dirty index -> OutcomeRefused, not-git and invalid-options and missing-mutate -> OutcomeFailure, mutation/commit failure with clean rollback -> OutcomeRolledBack, incomplete rollback and committed-state risk -> OutcomeRisk, otherwise OutcomeSuccess). Delete `gittransaction.RollbackResult` and `RollbackStatus`, moving the status type and its three constants into resultmodel so `resultmodel.RollbackResult.Status` is typed; gittransaction uses that one struct. Update git_transaction_test.go assertions from `result.ExitCode != N` to `resultmodel.ExitCode(result.Outcome) != N` so the tests still pin the numeric contract through the single authority.
   - *Serves:* Enforce exit codes 0-4 exactly as specified by the UR (Detailed Requirements 6); Builder Guidance: no parallel implementations.
2. Fix the two confirmed Git-layer defects and pin them. (a) In `runGit`, capture stdout and stderr separately (`cmd.Stdout`/`cmd.Stderr` buffers, or `Output()` plus `*exec.ExitError.Stderr`) and return only stdout on success, folding stderr into the error text; today `CombinedOutput()` makes a git warning read as porcelain content so `targetIsDirty` refuses a clean target and `resolveRepositoryRoot` can return a garbage root. (b) Give `indexIsEmpty` a variadic pathspec parameter and call it as `indexIsEmpty(ctx, repositoryRoot, targetPaths...)` inside `rollbackFailure`, leaving the two whole-index calls (the `--commit` precondition and the post-commit check) unscoped; today an unrelated staged file turns a successful rollback into exit 4. Add two lock-in tests in git_transaction_test.go naming these failures: one writes `*.txt [attr]bogus` into `.gitattributes` so `git status` emits a stderr warning with empty stdout and asserts the clean target is still mutated (not refused), the other stages an unrelated file, forces a mutation failure, and asserts OutcomeRolledBack with RollbackSucceeded.
   - *Serves:* Implement Git target preflight ... rollback (Detailed Requirements 7); Constraint: dirty target paths are refused while unrelated dirty paths remain allowed.
3. Add the Git-failure-to-finding bridge and make every emitted finding complete. New file skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go with `BuildCommandResult(result TransactionResult) resultmodel.CommandResult`: outcome from task 1, changes from `ChangedPaths` (add a `CreatedPaths []string` field to TransactionResult, populated from the recorder, so a change's Kind is truthfully `created` or `modified`), rollback passed straight through, and one finding per failure whose Code is derived mechanically from the FailureKind string (uppercase, underscores to dashes, `GIT-` prefix) so a new kind cannot silently miss a hand-written case. A per-kind template table supplies severity, fixability, automation-stop reason, next argv (`git status --short -- <path>` for a dirty target, `git revert <sha>` for committed-state risk) and verification argv; an unmapped kind yields a `GIT-UNMAPPED-FAILURE` error finding rather than an incomplete one. In command_runtime.go, fill `usageFinding`'s empty `VerificationArgv` so the runtime's own findings satisfy the same contract. New transaction_findings_test.go asserts, for a fixture per failure kind, that the finding carries code, severity, evidence, fixability, stop reason, at least one of next argv / Just recipe, and a verification command.
   - *Serves:* One typed result model ... findings, changes, skipped work, rollback result (Detailed Requirements 4); every finding carries a stable code ... and verification command (Detailed Requirements 5).
4. Prove the documented exit codes and text/JSON parity through the runtime seam. In command_runtime_test.go add a table test that registers a handler which runs a real `gittransaction.ExecuteTransaction` against a temporary Git fixture and returns `BuildCommandResult(...)`, then drives `CommandRuntime.Run` with `--repo-root <fixture> --format json` and again with `--format text`, asserting the returned process exit code is 0 (clean mutation), 1 (dirty target refused), 3 (mutation failure with clean rollback) and 4 (committed-state risk from a failing PostCommitVerify), that the JSON decodes with `schema_version`, non-null findings/changes/skipped_work/rollback, and that the text rendering of the same result names the same finding code and next argv. This is the first place the global `--repo-root` and `--format` options are exercised against a real Git mutation rather than a stub.
   - *Serves:* Support global --repo-root and --format text|json (Detailed Requirements 2); Red-Green Proof GREEN: the same typed result drives stable text/JSON and documented exit codes are observed.
5. Wire the module into the canonical gate. In _dev/tests/contract-regressions.sh add a probe block beside the other `_dev/tests/*.sh` invocations that fails when `_dev/tests/do-work-cli-launcher-behavior.sh` is missing or not executable and otherwise runs it, with a comment naming what it covers (build-on-demand, argv passthrough, stale-output refusal, Go floor refusal, no leftover build temp). In _dev/tests/maintainer-verify.sh add `go vet ./...` and `go test -count=1 ./...` lanes for skills/do-work/tools/do-work-cli in `run_verification`, and update the self-test in lockstep or it goes red: a `*/skills/do-work/tools/do-work-cli)` case in `write_command_shim` recording `cli-vet`/`cli-test`, the fixture directory in `run_self_test`'s `mkdir -p` list, the two new stage names plus `expected_count` 9->11 and 10->12 in `assert_success_stages`, and the two stages appended to the failure-stage loop so a missing lane fails the self-test.
   - *Serves:* Put one do-work-cli module under the installed core package and require Go 1.26.1+ (Detailed Requirements 1); build on demand when the binary is missing or older than its Go sources (Detailed Requirements 3) — both only hold if the gate actually runs them.

**Files to touch:**

- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify) — Replace ExitCode int with a resultmodel outcome, drop the duplicate RollbackResult type, fix runGit's stderr-as-content bug, scope the post-rollback index check to target paths, expose CreatedPaths.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify) — Retarget existing exit-code assertions through resultmodel.ExitCode and add the two lock-in tests (git stderr warning must not read as dirt; unrelated staged work must not become exit 4).
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go` (new) — BuildCommandResult maps a TransactionResult to the one typed CommandResult, with finding codes derived from FailureKind and a per-kind template for severity, next argv, and verification argv.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go` (new) — Asserts every failure kind produces a complete finding and that an unmapped kind fails loudly instead of emitting an incomplete one.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — Own the RollbackStatus type and its constants so the rollback result has one definition, and keep ExitCode the single exit-code authority.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — Cover the typed rollback status in the JSON/text parity test that already pins schema_version and non-null collections.
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` (modify) — usageFinding currently ships no verification command, so the runtime's own findings violate the per-finding completeness requirement.
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` (modify) — Add the end-to-end table test that drives real Git fixtures through Run and observes exit codes 0/1/3/4 plus text and JSON from the same result.
- `_dev/tests/contract-regressions.sh` (modify) — The launcher probe exists but nothing invokes it; the file's own convention is that no _dev/tests script is auto-discovered.
- `_dev/tests/maintainer-verify.sh` (modify) — Add do-work-cli go vet and uncached test lanes and update the self-test shim, fixture directories, stage list, stage counts, and failure-stage loop in the same edit.
- `VERSION` (modify) — Before Every Commit ritual: bump the shared version for the integrating commit (0.244.9 is current).
- `skills/do-work/VERSION` (modify) — Same shared version bump.
- `skills/do-work/actions/version.md` (modify) — The **Current version**: line must match the bumped VERSION files.
- `CHANGELOG.md` (modify) — New top entry with a title that says what was delivered and a version strictly greater than the current first entry.
- `skills/do-work/CHANGELOG.md` (modify) — Byte-identical mirror of root CHANGELOG.md, enforced by _dev/tests/shipped-package-reference-contract.sh.
- `do-work/working/REQ-406-create-do-work-cli-foundation.md` (modify) — Tick the P-A-U boxes, record the implementation summary and evidence, and record the implementation commit hash in a separate bookkeeping commit.

**Decisions the builder must honour:**

- D1 Complete the preserved commit rather than restart. Every preserved file passes its own checks today (go vet, go test -count=1, the launcher probe, ShellCheck), so the diff should be additive and corrective. Do not regenerate the launcher: its EXIT-vs-signal trap split is exactly what the prime file prescribes and the REQ records it as the last thing fixed before the interruption.
- D2 One exit-code authority. Recommended: delete TransactionResult.ExitCode and derive the number from resultmodel.ExitCode(Outcome). The alternative (keep both and add a test that they agree) preserves the drift this REQ exists to prevent, and costs more code than the deletion.
- D3 Where the Git-to-result bridge lives. Recommended: gittransaction imports resultmodel (no cycle; commandruntime already imports resultmodel and does not import gittransaction). Putting the mapping in commandruntime instead would force every later command package to re-derive findings from raw FailureKind values.
- D4 Finding codes are derived, not enumerated. Recommended: compute the code from the FailureKind string and keep only severity/next-argv in a per-kind table, with an explicit GIT-UNMAPPED-FAILURE fallback. A hand-written switch over eight kinds is exactly the closed enumeration the shell prime warns goes stale.
- D5 CreatedPaths on TransactionResult. Recommended: add it, so a RecordedChange.Kind is truthfully created vs modified rather than a constant filler. If the builder judges the field unearned, emitting a single kind is acceptable — but then say so in the REQ rather than leaving Kind meaningless.
- D6 No command is registered in this REQ. main.go stays a three-line delegation with a nil handler map; REQ-407 onward register real commands. The exit-code contract is therefore observed at the CommandRuntime.Run seam that main delegates to entirely, not by exec-ing the binary.
- D7 Do not add a second shell probe that builds and runs the real binary. The launcher probe already covers absent-binary build, argv passthrough including --format json, stale-output refusal, and the Go floor; the new maintainer-verify lanes compile cmd/do-work-cli on every run. Revisit when REQ-407 gives the binary a command worth exec-ing.
- D8 Leave the staged-skills live-surface scan alone. The built binary is gitignored but still walked by collect_live_files (which exempts queue-kanban by name); verified empirically that the current binary does not trip the retired-trigger matcher, so no edit is earned today.
- D9 The release ritual belongs to this REQ's integrating commit: bump VERSION, skills/do-work/VERSION, and actions/version.md; add a CHANGELOG.md entry titled for what shipped; copy it byte-identically to skills/do-work/CHANGELOG.md. Commit the implementation with a blank commit: field, then record the hash in a separate bookkeeping commit.

**Testing approach:**

All commands from the repo root /home/user/skill-do-work. RED first, per task. Task 2: add the two lock-in tests before touching git_transaction.go and run `cd skills/do-work/tools/do-work-cli && go test -count=1 ./internal/gittransaction/` — the stderr-warning test must fail with the clean target refused (FailureDirtyTarget) and the unrelated-staged test must fail with an OutcomeRisk/exit-4 result; capture both failure messages as RED evidence. Task 4: add the exit-code table test before the bridge exists and confirm `go test -count=1 ./internal/commandruntime/` fails to build or fails on the missing outcomes. Task 5: after wiring, prove the gate lanes are not decorative — run `bash _dev/tests/maintainer-verify.sh --self-test` with `cli-vet` and `cli-test` in the failure-stage loop (it must fail for each injected stage failure), and temporarily `chmod -x _dev/tests/do-work-cli-launcher-behavior.sh`, confirm `bash _dev/tests/contract-regressions.sh` fails naming the missing probe, then `chmod +x` and re-run. GREEN, in order: `cd /home/user/skill-do-work/skills/do-work/tools/do-work-cli && go vet ./... && go test -count=1 ./...`; `bash /home/user/skill-do-work/_dev/tests/do-work-cli-launcher-behavior.sh`; `bash /home/user/skill-do-work/_dev/tests/maintainer-verify.sh --self-test`; `shellcheck --severity=warning -- _dev/tests/maintainer-verify.sh _dev/tests/contract-regressions.sh _dev/tests/do-work-cli-launcher-behavior.sh skills/do-work/tools/do-work-cli.sh`; formatting judged by output emptiness, never exit status: `"$(go env GOROOT)/bin/gofmt" -l -- $(git ls-files '*.go')` must print nothing. Final gate, unpiped, exit 0 is the only proof: `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium-1194/chrome-linux/chrome bash _dev/tests/maintainer-verify.sh` — never pipe it through tail or head, the pipeline status hides the failure. Baseline observed during planning: go vet, go test -count=1 ./..., the launcher probe, ShellCheck, and staged-skills-contract.sh all pass today with the built binary present on disk.

**Risks:**

- R1 maintainer-verify's self-test is a closed enumeration in four places (write_command_shim PWD case, run_self_test mkdir list, assert_success_stages stage loop plus expected_count 9->11 and 10->12, and the failure-stage loop). Miss one and the self-test either goes red for the wrong reason or counts the new lanes as absent; contract-regressions runs it, so a partial edit fails the whole gate.
- R2 Retyping TransactionResult.ExitCode as an outcome touches every assertion in git_transaction_test.go. Failures here are compile errors, not silent passes, but a builder tempted to keep both fields would reintroduce exactly the second exit-code authority this REQ removes.
- R3 The real-binary end-to-end path (launcher -> compiled binary -> text/JSON with a working command) stays unproven mechanically until REQ-407 registers a command; the launcher probe exercises argv and build policy against a fake go toolchain. Verified manually during planning that the real binary renders both formats and exits 2 for an unknown command.
- R4 The built binary skills/do-work/tools/do-work-cli/do-work-cli is gitignored but still read by staged-skills-contract's live-surface walk, which exempts only queue-kanban by name. It passes today; a future finding string containing the literal 'do-work ' with a trailing space would turn a compiled artifact into a false retired-trigger violation.
- R5 Five tasks is the cap and the release ritual plus REQ bookkeeping ride along with the integrating commit. If the builder splits task 3 (bridge) from task 1 (outcome collapse) into separate commits, task 1 must land first or BuildCommandResult has no outcome to map.
- R6 Switching runGit to stdout-only changes error strings. No current test matches git error text, but a builder adding one during task 2 should match on FailureKind, not on message text.

**Plan validation:** 5 tasks (at the quality cap, not over it); every requirement maps to a task and every task traces to a requirement — no coverage gaps, no orphan tasks.

*Generated by Plan agent*

## Exploration

**Key files:**

- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` — The whole Git transaction. 432 lines, package-local, imports nothing from the module today (so importing resultmodel introduces no cycle: commandruntime->resultmodel, gittransaction->resultmodel, nobody imports gittransaction). (14-25 FailureKind + 8 constants (invalid_options, not_git_repository, dirty_target, dirty_index, mutation_failed, rollback_incomplete, commit_failed, committed_state_risk). 27-39 RollbackStatus + RollbackResult (the duplicate of resultmodel.RollbackResult). 46-54 TransactionResult{ExitCode int, RepositoryRoot, ChangedPaths, CommitSHA, RevertArgv, Rollback, Failure}. 56-63 TransactionOptions{RepositoryRoot,TargetPaths,DryRun,Commit,CommitMessage,PostCommitVerify}. 65-94 MutationRecorder with unexported allowedPaths/creatablePaths/touchedPaths/createdPaths — createdPaths is the only source for a truthful created-vs-modified Kind and is not reachable from TransactionResult today. 102-203 ExecuteTransaction. 103 result construction (ExitCode implicitly 0 = success). 105,108,112,117,121,126,138,144,151 fail(...,2,...). 129,132,141 fail(...,1,...). 148/168-170/202 bare `return result` success paths. 189 `committedPaths, verifyErr := committedPaths(...)` — the local shadows the function for the rest of the scope. 288-294 targetIsDirty (len(output)>0 is the whole predicate). 296-307 indexIsEmpty — whole-index only, and the ONE git call that omits --literal-pathspecs. 309-321 changedTargets. 323-369 rollbackFailure; 358-362 the unscoped index check; 366 fail(...,4,FailureRollback); 368 fail(...,3,kind). 381-385 committedRisk -> fail(...,4). 387-391 fail(). 393-401 runGit — CombinedOutput() is the confirmed defect.)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` — 7 tests, all asserting result.ExitCode as a bare int. Fixture helpers the new lock-in tests must reuse. (20,36,52,68,78,93,118,147,180 the ExitCode assertions that must be retargeted through resultmodel.ExitCode(result.Outcome). 118,147 also assert RollbackSucceeded/RollbackIncomplete (become resultmodel.* after the type move). 197-204 newRepository (git init -q + local user.name/user.email; does NOT isolate GIT_CONFIG_GLOBAL/SYSTEM). 206-210 commitAll. 212-220 runFixtureGit. 222-231 writeFile. 233-245 readFile/readPath. 247-253 commitPaths.)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` — The single typed result model and the exit-code authority. Already satisfies the UR's 0-4 table exactly. (10 SchemaVersion=1. 19-28 CommandOutcome{success,findings,refused,failure,rolled_back,committed_state_risk}. 30-44 FindingSeverity/FindingFixability. 46-57 CommandFinding (all 10 fields requirement 5 names). 59-74 RecordedChange/SkippedWork/RollbackResult (Status is a bare string — the field task 1 retypes). 76-85 CommandResult. 87-102 ExitCode: success=0, findings|refused=1, failure=2, rolled_back=3, risk=4, default=2. 104-140 NormalizeResult (nil-slice -> empty). 142-156 RenderResult. 158-200 renderText. 202-215 joinArgv.)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` — Global-option parser and the Run seam where exit codes are observed without exec-ing the binary. (13-18 ExecutionContext / CommandHandler. 35-60 Run. 62-78 writeResult (returns resultmodel.ExitCode(result.Outcome)). 80-133 parseGlobalOptions (--repo-root, --repo-root=, --format, --format=; first non-dash arg is the command). 135-144 usageFinding — VerificationArgv is absent, which is the requirement-5 hole task 3 closes. 146-152 absolutePath.)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` — 3 tests; all handlers are stubs, no Git fixture anywhere yet. Task 4's table test goes here. (14-48 TestReadOnlyHandlerRunsOutsideGitWithGlobalOptions (the --repo-root/--format json path). 50-67 usage/UNKNOWN-COMMAND exit 2. 69-89 invalid global options.)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` — JSON/text parity plus the numeric exit-code table. The place a typed RollbackStatus gets covered. (9-62 TestRenderersUseOneNormalizedResult (schema_version, non-null findings/changes/skipped_work, rollback.actions non-null, text substrings). 64-81 TestOutcomeExitCodes — the 6-row table. 83-86 unknown-format rejection.)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli.sh` — The build-on-demand launcher. Do not regenerate: the EXIT-vs-signal trap split at 51-54 is exactly what prime-shell-commands.md prescribes. (7 minimum_go_version=1.26.1. 9-21 version_at_least (awk, 3 components). 23-28 rebuild gate: not-executable OR any *.go/go.mod/go.sum -newer the binary. 30-43 Go presence + floor refusal, exit 2. 45-54 cleanup on EXIT, exit 129/130/143 on HUP/INT/TERM. 55-61 mktemp in module dir, build, chmod, mv -f (temp-then-rename, the curl -o trap's remedy). 62 trap - EXIT HUP INT TERM. 65 exec.)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/go.mod` — module github.com/knews2019/skill-do-work/do-work-cli; go 1.26.1. No go.sum — stdlib only. Installed toolchain is go1.26.1 linux/amd64, exactly at the floor. (1-3)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/.gitignore` — /do-work-cli and /do-work-cli.build.* . Ships with the module: root .gitattributes export-ignores only /.gitignore (root-anchored), same as queue-kanban's. (1-2)
- `/home/user/skill-do-work/_dev/tests/do-work-cli-launcher-behavior.sh` — The launcher probe, mode 100755, tracked, passing — but invoked by nothing. Fake-go fixture is the model for any shell probe that must not depend on a real toolchain. (12-17 fixture tree + trap. 19-52 fake `go` honoring FAKE_GO_VERSION/FAKE_GO_FAIL_BUILD/FAKE_GO_BUILD_LOG, emitting a bash 'binary' that echoes argv. 61-69 argv passthrough incl. a spaced arg + exactly-one-build. 71-75 no rebuild when fresh. 77-83 `sleep 1` before touch (1-second mtime granularity — required). 85-94 failed rebuild must not run stale output. 96-104 Go floor refusal. 106-109 no leftover do-work-cli.build.* .)
- `/home/user/skill-do-work/_dev/tests/contract-regressions.sh` — 7224 lines. The aggregate that maintainer-verify calls. Nothing auto-discovers _dev/tests/*.sh — the file says so in its own words at 6550-6551. (6426-6433 maintainer-verify --self-test probe (the -x + elif form to copy). 6548-6563 record-commit-hash block, the canonical `# comment naming what it covers` + probe-var + [ ! -x ] + elif ! bash form. 6539-6546, 6568-6577, 6580-6590, 6592-6602, 6619-6629 the other probe blocks, all beside each other — put the do-work-cli block in this run. 7147-7191 a python self-check over this file's own late `assert_contains "justfile"` lines; it parses a strict 3-line regex, so do not reformat those. 7216-7222 fail_count gate.)
- `/home/user/skill-do-work/_dev/tests/maintainer-verify.sh` — The canonical gate. Its self-test is a closed enumeration in four places that must all move together. (10 minimum_go_version=go1.26.1. 53-183 write_command_shim; 94-121 the `case "$PWD"` with */skills/do-work-board/tools/queue-kanban and */skills/do-work-toolbox/tools/audit-metrics arms and `*) exit 64` — a new arm `*/skills/do-work/tools/do-work-cli)` goes here. 185-197 count_stage_occurrences. 199-233 assert_success_stages: 206 expected_count=9, 209 =10 under strict, 211-213 the 9-name loop, 220-225 board-strict special case. 235-379 run_self_test; 262-266 the mkdir -p fixture list (add skills/do-work/tools/do-work-cli or the lane's `cd` fails); 333-335 the failure-stage loop names (11 today). 381-527 run_verification; 448-465 gofmt lane driven by `git ls-files -z -- '*.go'` (untracked files are invisible to it); 467-468 aggregate; 470-479 queue-kanban vet/test — the exact shape to copy for do-work-cli; 515-524 audit-metrics vet/test; 494-513 browser-lane guard (QUEUE_KANBAN_BROWSER or a PATH candidate, explicit SKIP otherwise).)
- `/home/user/skill-do-work/_dev/tests/staged-skills-contract.sh` — Live-surface retired-trigger scan. It walks the compiled do-work-cli binary. (20-51 core_files require-list (do-work-cli.sh is not in it; it is a must-exist list, not exhaustive). 60-83 legacy_runtime_paths retirement check. 501 retired_command_head = re.compile(r"(?<![A-Za-z0-9_-])do-work "). 504-508 exempt_command_fragments. 510-517 collect_live_files — walks skills/**/* and exempts ONLY the file names CHANGELOG.md and queue-kanban, so skills/do-work/tools/do-work-cli/do-work-cli (gitignored but on disk) is read with errors='replace' and scanned. 715-721 the violation loop.)
- `/home/user/skill-do-work/do-work/user-requests/UR-081/input.md` — Verbatim source of the exit-code table and the mutation rules the REQ restates. (97 the exit codes: 0 clean/success, 1 findings or safely refused, 2 usage/precondition/runtime failure, 3 operation failure with successful rollback, 4 incomplete rollback or committed-state risk. 82-88 mutation rules. 93-96 the JSON field list. 121 assigns the maintainer-verify go vet / uncached do-work-cli lanes to the LAST REQ of the batch.)
- `/home/user/skill-do-work/do-work/queue/REQ-420-replace-shell-implementations-verify-parity.md` — Queued REQ whose Detailed Requirements already own the maintainer-verify wiring the plan puts in task 5. (34: 'Extend _dev/tests/maintainer-verify.sh with go vet and uncached do-work-cli tests, retain queue-kanban verification, and replace the separate audit-metrics lane.')
- `/home/user/skill-do-work/do-work/runs/work-2026-08-29-213539/manifest.md` — REQ-406 is dispatched as a worktree builder in a 2-REQ fan-out wave alongside REQ-390. (9-12: operative name worktree-agent-REQ-406-create-do-work-cli-foundation, hand-back REQ-406-handback.md, Landed pending.)
- `/home/user/skill-do-work/skills/do-work/actions/work-reference.md` — The serial-only rule that decides whether the builder may touch VERSION/CHANGELOG at all. (407: 'Serial-only — never parallelised, at any builder count: every do-work/ queue transition (claim, status flip, archive move); REQ id allocation; and the project-owned release version file plus CHANGELOG.md — one changelog entry per REQ, written by the owner at merge time.' 345 the red-flag framing. 983 in worktree dispatch the builder's implementation is committed on its branch and the orchestrator stages the changelog/version/REQ separately.)
- `/home/user/skill-do-work/_dev/primes/prime-shell-commands.md` — The REQ's only declared prime. Directly load-bearing traps for this work. (21 EXIT vs signal trap contract (the launcher already implements it — do not 'simplify' it). 22 a tool that reports on stdout while exiting zero makes an exit-status lane decorative, plus `local name; name="$(...)"` . 31-47 Unchecked Exit Status Reads as Content. 49-51 Closed Enumerations Go Stale (governs the finding-code table and the self-test stage lists). 53-58 every flag on a shipped script needs a non-test caller.)
- `/home/user/skill-do-work/_dev/primes/lessons-shell-commands.md` — Prior-failure index. Step 8 substep 7 appends here on archive. (REQ-423 (cleanup-only signal trap returns to the workflow), REQ-260 (a lane can bite correctly yet be silently removable; gofmt exit-status shape), REQ-187 (one maintainer command inventory, fixture-only self-test mode), REQ-234 (a derived count that almost reproduces the remembered figure), REQ-257 (a lock-in that enumerates only the spellings it already knew reads as coverage while unable to fail), REQ-250 (grep the same primitive across the file before calling a class closed), REQ-309 (the repo gate runs from root and is judged only by its direct exit status).)
- `/home/user/skill-do-work/skills/do-work/crew-members/coding-guardrails.md` — Always loads during implementation; § 5 is the canonical naming rule that governs any new identifier here. (118-144 § 5 Naming for Reach: two words minimum for anything with reach, idiomatic short locals exempt, and 'Precedence against §3' — the rule governs names you introduce, not existing ones.)
- `/home/user/skill-do-work/_dev/tests/record-commit-hash-guards.sh` — The house style for a standalone _dev/tests shell probe: header stating what it proves, who invokes it, and the exit-code contract; global-git isolation. (1-15 the header block. 32-35 export GIT_CONFIG_GLOBAL=/dev/null, GIT_CONFIG_SYSTEM=/dev/null, GIT_TERMINAL_PROMPT=0. 38-40 mktemp fixture + trap.)
- `/home/user/skill-do-work/skills/do-work-board/tools/queue-kanban` — Reference Go module: .gitignore holds only /queue-kanban, has go.sum, is the module the maintainer-verify lanes were written for. Its browser/JS lanes are the pattern for an opt-in strict lane, not needed here. (generate_test.go:234-265 TestMaintainerStrictJavaScriptBehaviorLane + its RejectsZeroProbes twin (a lane that re-execs the test binary with a marker env var and fails when zero probes ran). browser_probe_test.go:721,740 the browser equivalent.)
- `/home/user/skill-do-work/suite/modules.tsv` — Sole source/destination manifest — 4 rows, one per skills/ package. do-work-cli lives INSIDE skills/do-work, so no row is added and validate-suite-manifest.sh's four-module check stays untouched. (1-5)

**Patterns to follow:**

- Go naming in this module is already fully-spelled: repositoryRoot, commandArgs, outputFormat, mutationErr, failureKind. No single-letter or abbreviated identifiers anywhere except the conventional ctx/err/t. Match it — coding-guardrails § 5 requires two words for anything with reach.
- Error wrapping is `fmt.Errorf("<verb phrase>: %w", err)` (git_transaction.go:227,231,277,306,398). Sentinel-free; failures are classified by FailureKind, never by message text.
- JSON field names are snake_case struct tags on CamelCase fields (result_model.go:46-85). Every slice is nil-normalized to `[]` in NormalizeResult before render — a new collection field needs a corresponding nil-check there or JSON emits null.
- Go tests are plain table/scenario funcs with `t.Fatalf("... = %#v", result)` on the whole struct — no testify, no subtests, no golden files. Helpers take `t *testing.T` first and call `t.Helper()`.
- Git fixtures are built in `t.TempDir()` by newRepository/commitAll/runFixtureGit/writeFile (git_transaction_test.go:197-253). Reuse them; do not introduce a second fixture helper.
- _dev/tests probe scripts: `#!/usr/bin/env bash`, `set -euo pipefail` (or `set -uo pipefail` + a fail_count tally for multi-probe files), repo_root derived from BASH_SOURCE, mktemp fixture root under ${TMPDIR:-/tmp} with a `trap 'rm -rf ...' EXIT`, `echo/printf FAIL: ... >&2` then a nonzero exit, and a single success line on stdout at the end.
- contract-regressions.sh probe blocks are all the same four lines: a comment naming what the probe covers and why it cannot be a grep, `name_probe="$repo_root/_dev/tests/<file>.sh"`, `if [ ! -x ]` (or `! -f`) printing a FAIL that names the missing file and what loses coverage, `elif ! bash "$probe"` printing a FAIL that points at the lines above — each incrementing fail_count rather than exiting.
- maintainer-verify lanes are `printf 'maintainer-verify: <name>\n'` followed by a subshell `( cd "$repo_root/<module>" ; go vet ./... )` — the cd is inside the subshell so cwd never leaks (lines 470-479, 515-524).
- The maintainer-verify self-test shims commands by name and asserts the EXACT argv (`[ "$#" -eq 2 ] && [ "$1" = 'vet' ] && [ "$2" = './...' ]`, else `exit 64`), so a lane whose argv differs from its shim arm fails loudly rather than silently passing.
- Exit codes are derived in exactly one place (resultmodel.ExitCode) and every caller goes through it. Requirement 6 plus the Builder Guidance ban on parallel implementations make a second numeric authority the thing this REQ removes.
- Shipped shell must not cite _dev/ (CLAUDE.md); do-work-cli.sh currently cites nothing outside itself. Keep it that way.
- A rule keyed on a condition must not be restated as a hand-maintained list (prime-shell-commands.md § Closed Enumerations Go Stale) — this is the argument for deriving finding codes from FailureKind rather than switching over eight kinds.

**Test conventions:**

"Everything runs from the repo root /home/user/skill-do-work unless a cd is shown. Baselines I measured just now, all green on the preserved commit 329c55a9 with no working-tree changes.\n\nGo, per package: `cd /home/user/skill-do-work/skills/do-work/tools/do-work-cli && go vet ./... && go test -count=1 ./...` (installed toolchain is go1.26.1 linux/amd64, exactly at the module floor; current run: commandruntime 0.004s, gittransaction 0.558s, resultmodel 0.005s, cmd has no test files). Focused: `go test -count=1 -run 'TestName' -v ./internal/gittransaction/`.\n\ngofmt is judged by output emptiness, never exit status: `cd /home/user/skill-do-work && \"$(go env GOROOT)/bin/gofmt\" -l -- $(git ls-files '*.go')` must print nothing. Take gofmt from GOROOT, not PATH — maintainer-verify.sh:409-418 explains why.\n\nLauncher probe: `bash _dev/tests/do-work-cli-launcher-behavior.sh` — 0.6s, prints 'do-work-cli launcher behavior tests passed', exit 0 today. It never touches the real toolchain: it copies the launcher into a temp tree, puts a fake `go` on PATH via FAKE_GO_VERSION / FAKE_GO_FAIL_BUILD / FAKE_GO_BUILD_LOG, and counts build lines. Two `sleep 1` calls before `touch` are load-bearing (1-second mtime granularity vs the `find -newer` gate).\n\nShellCheck: `shellcheck --severity=warning -- _dev/tests/maintainer-verify.sh _dev/tests/contract-regressions.sh _dev/tests/do-work-cli-launcher-behavior.sh skills/do-work/tools/do-work-cli.sh`. The gate's own lane runs it over every tracked *.sh.\n\nGate self-test: `bash _dev/tests/maintainer-verify.sh --self-test` — 0.53s, prints 'Maintainer verification self-test passed.' It runs the script against a shimmed PATH four ways (all-success, no-Node, newer-tools, and one run per injected failure stage) and asserts an exact per-stage occurrence count plus an exact total. It is invoked by contract-regressions.sh:6426-6433, so a broken self-test fails the whole aggregate.\n\nAggregate: `bash _dev/tests/contract-regressions.sh` — ends with 'Contract regression checks passed.'\n\nFinal gate, unpiped (a pipe hides the status — CLAUDE.md § Verify): `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium-1194/chrome-linux/chrome bash _dev/tests/maintainer-verify.sh`. I ran it without the browser var: exit 0, ~7 minutes, dominated by the queue-kanban strict JavaScript lane at 82s, and it printed 'SKIP: no browser is available'. The chromium binary at that path exists and is executable, so setting the var adds the strict browser lane (~77s more per the checkpoint).\n\nBrowser/JS probes here are Go tests, not shell: a strict lane self-selects on `flag.Lookup(\"test.run\")` matching its own pattern and otherwise t.Skip's (generate_test.go:250-256), and each lane has a paired *RejectsZeroProbes test that re-execs the test binary with an empty PATH and requires it to FAIL — that is what stops a skipped lane from reading as green. There is no browser surface in do-work-cli, so no browser probe is earned; the shell-probe shape (fake-toolchain fixture, FAIL lines on stderr, one success line on stdout) is the relevant one.\n\nRED evidence I captured for the two Git defects, reproduce these before fixing:\n(a) In a fixture repo containing `*.txt [attr]bogus` in .gitattributes, `git status --porcelain=v1 -z -uall -- target.txt` writes 0 bytes to stdout and 60 bytes to stderr. Through runGit's CombinedOutput() a clean target refuses: observed `exit=1 called=false failure=&{Kind:dirty_target Reason:target path \"target.txt\" is already dirty}`.\n(b) With one unrelated file staged (`git add unrelated.txt`) and a forced mutation failure on a tracked target, observed `exit=4 rollbackStatus=incomplete errors=[Git index is not empty after rollback] failure=&{Kind:rollback_incomplete Reason:forced mutation failure}` — a successful rollback reported as exit 4.\nScoping proof for the fix: `git diff --cached --quiet --exit-code -- <path>` exits 0 for an unstaged path, 1 for a staged one, and with `--` and no pathspecs behaves exactly like the whole-index form, so a variadic with zero args is safe."

**Concerns and traps:**

- C1 SCOPE — the release ritual may not be the builder's to perform. REQ-406 is dispatched as a worktree builder (do-work/runs/work-2026-08-29-213539/manifest.md:12, operative name worktree-agent-REQ-406-create-do-work-cli-foundation) in a 2-REQ fan-out wave alongside REQ-390. CLAUDE.md § Before Every Commit says the ritual belongs to the integrating commit only and 'A builder committing on its own worktree-agent-* branch skips it entirely'; work-reference.md:407 makes VERSION/version.md/CHANGELOG.md and every do-work/ queue transition serial-only at any builder count. The plan's D9 and six of its filesToTouch entries (VERSION, skills/do-work/VERSION, skills/do-work/actions/version.md, CHANGELOG.md, skills/do-work/CHANGELOG.md, the REQ file) are the orchestrator's writes. A builder bumping 0.244.9 races REQ-390. Confirm which tree you are in before touching any of them.
- C2 SCOPE — task 5 is REQ-420's stated deliverable, verbatim. do-work/queue/REQ-420-replace-shell-implementations-verify-parity.md:34 and UR-081 input.md:121 both read 'Extend _dev/tests/maintainer-verify.sh with go vet and uncached do-work-cli tests, retain queue-kanban verification, and replace the separate audit-metrics lane.' REQ-406's own Detailed Requirements never mention the gate. Doing it now is defensible (an unwired module is unverified) but it is cross-REQ scope and should be stated in the REQ rather than done silently — and note REQ-420 also removes the audit-metrics lane, which task 5 must not pre-empt.
- C3 THE ZERO-VALUE FLIP — highest-risk mechanical detail in task 1. TransactionResult.ExitCode's zero value is 0 = success, which is what the three bare `return result` success paths rely on (git_transaction.go:148 dry-run, 168-170 no-commit/no-changes, 202 full success). resultmodel.CommandOutcome's zero value is "", and resultmodel.ExitCode("") returns 2. Unless line 103's construction explicitly sets Outcome: resultmodel.OutcomeSuccess, every success path silently becomes a usage failure. The existing tests at lines 52, 68 do catch it once retargeted — retarget them first.
- C4 FUNCTION SHADOWED BY A LOCAL. git_transaction.go:189 is `committedPaths, verifyErr := committedPaths(ctx, repositoryRoot, commitSHA)` — the local shadows the package function for the rest of ExecuteTransaction's scope. It compiles today because only one call exists. Any second call added below line 189 fails to build with 'cannot call non-function'. Rename the local (e.g. actualCommittedPaths) if you touch that region.
- C5 THE CombinedOutput BUG HAS THREE CONSUMERS, NOT ONE. targetIsDirty's `len(output) > 0` predicate (git_transaction.go:293) is reached from the preflight loop (124), from changedTargets (312, so every declared target reads as changed and an untouched-but-declared path then triggers a bogus rollback at 163-166), and from rollbackFailure (329). Fix runGit, then re-read all three. The plan also claims resolveRepositoryRoot can return a garbage root: plausible by the same mechanism, but I could not reproduce a stderr-emitting `git rev-parse --show-toplevel` (a broken ref file produced no warning). Treat the status path as the confirmed instance and the rev-parse path as the same class, unproven.
- C6 THE SELF-TEST IS A CLOSED ENUMERATION IN FOUR PLACES AND ALL FOUR MOVE TOGETHER: write_command_shim's `case "$PWD"` arms (maintainer-verify.sh:94-121, add `*/skills/do-work/tools/do-work-cli)` before the `*) exit 64`), run_self_test's mkdir -p fixture list (262-266 — omit it and the new lane's `cd` fails, and the failure looks like a lane bug), assert_success_stages (206 expected_count=9, 209 =10, and the 9-name loop at 211-213), and the failure-stage loop (333-335). Miss the count and the all-success run fails with 'success run recorded N stages; want M'; miss the failure loop and the new lanes are unfalsifiable — exactly the 'reads as coverage while unable to fail' shape prime-shell-commands.md warns about. contract-regressions.sh runs --self-test, so a partial edit reddens the whole aggregate.
- C7 UNTRACKED GO FILES ESCAPE THE GOFMT LANE. maintainer-verify.sh:450 feeds gofmt from `git ls-files -z -- '*.go'`. transaction_findings.go and transaction_findings_test.go are invisible to that lane until they are `git add`ed, while go vet/go test compile them from disk. Stage new Go files before treating a green gate as proof.
- C8 THE COMPILED BINARY IS SCANNED AS A LIVE SURFACE. staged-skills-contract.sh:510-517 collect_live_files walks skills/**/* and exempts only the file names CHANGELOG.md and queue-kanban; skills/do-work/tools/do-work-cli/do-work-cli is gitignored but present on disk and gets read with errors='replace' and matched against `(?<![A-Za-z0-9_-])do-work ` plus 22 retired install targets. It passes today (I ran the probe: PASS). Any new finding string, evidence template, or usage text containing the literal 'do-work ' followed by a retired target word turns a build artifact into a contract violation. Prefer 'do-work-cli' as the literal in every emitted string.
- C9 THE GO FIXTURES DO NOT ISOLATE GLOBAL GIT CONFIG. newRepository (git_transaction_test.go:197-204) sets only local user.name/user.email. The shell probes here export GIT_CONFIG_GLOBAL=/dev/null, GIT_CONFIG_SYSTEM=/dev/null, GIT_TERMINAL_PROMPT=0 for exactly this reason (record-commit-hash-guards.sh:32-35). A developer's global core.autocrlf, core.hooksPath, commit.gpgsign or init.templateDir can change what the new stderr-warning test measures — and that test's whole point is controlling what git writes to stderr. Consider setting those in newRepository, and note that widening the helper touches all seven existing tests.
- C10 DELETING gittransaction.RollbackStatus IS A LOCK-STEP EDIT ACROSS THREE FILES: the type and its three constants (git_transaction.go:27-33), every construction site (103, 324, 365), and the test assertions at git_transaction_test.go:118 and 147. Nothing outside the package references them, so the compiler catches all of it — but keeping both types 'for compatibility' reinstates the duplication the REQ exists to remove.
- C11 usageFinding IS THE RUNTIME'S ONLY FINDING PRODUCER AND IT IS INCOMPLETE. command_runtime.go:135-144 emits no VerificationArgv (and no NextJustRecipe). Requirement 5 applies to every finding, including the runtime's own, so any completeness assertion written for the gittransaction bridge will fail against the runtime unless this is filled in the same change.
- C12 CHANGING runGit's ERROR TEXT IS OBSERVABLE. Today the message is `git <args>: <err>: <combined output>`; stdout/stderr separation changes it. No current test matches git error text — keep it that way and assert on FailureKind, never on the message.
- C13 NOTHING AUTO-DISCOVERS _dev/tests/*.sh — contract-regressions.sh says so at 6550-6551, and a probe file nobody invokes is 'dead weight that reads as coverage.' The launcher probe has been exactly that since 329c55a9. When adding the block, put it beside the others (6539-6629), not at the end, and use `[ ! -x ]` since the file is tracked 100755.
- C14 THE LAUNCHER IS FINISHED WORK — DO NOT REGENERATE IT. Its EXIT-vs-signal trap split (do-work-cli.sh:51-54, exit 129/130/143) is the exact contract prime-shell-commands.md line 21 and REQ-423's lesson prescribe, and the REQ records it as the last thing repaired before the interruption. Its temp-then-rename build (55-61) is the remedy for the `curl -o` partial-file trap. Same for the `find -newer` rebuild gate at 23-28.
- C15 A FINDING-CODE TABLE IS A CLOSED ENUMERATION. Whatever shape the per-kind severity/next-argv table takes, it must fail loudly on an unmapped FailureKind rather than emit an incomplete finding — prime-shell-commands.md § Closed Enumerations Go Stale, and the REQ-257 lesson that a lock-in enumerating only the spellings it already knew reads as coverage while being unable to fail. A test that asserts the eight known kinds map is weaker than one that also asserts a synthetic ninth kind produces the loud fallback.
- C16 `fail` (git_transaction.go:387) is a single-word package-level name with reach, which coding-guardrails § 5 forbids for names you introduce. §3 protects it as existing code — but task 1 changes its signature. If you are rewriting the call site anyway, failTransaction costs nothing; if you leave it, that is a defensible §3 call. Do not leave it half-renamed.

*Generated by Explore agent*

## Scope

**Files I will touch:**

- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify) — stdout-only runGit, dirty-target and staged-index refusals, CreatedPaths
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify) — lock-in tests for the two refusals
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go` (new) — derive stable finding codes from FailureKind
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go` (new) — finding-code derivation and the unmapped fallback
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — single exit-code authority
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — exit-code table
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` (modify) — Git-to-result bridge at the Run seam
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` (modify) — exit-code contract at the Run seam
- `_dev/tests/contract-regressions.sh` (modify) — register the launcher probe so it cannot silently stop running
- `_dev/tests/maintainer-verify.sh` (modify) — add the cli-vet and cli-test gate lanes

**Files I will NOT touch:**

- `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` — serial-only files owned by the integrator (CLAUDE.md § Before Every Commit). The plan's D9 is correct that the ritual belongs to this REQ, and wrong that this builder performs it.
- `skills/do-work/tools/do-work-cli.sh` and `_dev/tests/do-work-cli-launcher-behavior.sh` — preserved and passing; do not regenerate (plan D1).
- Anything under `do-work/` — queue state belongs to the orchestrator alone.

**Acceptance criteria (restated from REQ):**

- [ ] One `do-work-cli` module under the installed core package requiring Go 1.26.1+, with global `--repo-root` and `--format text|json`.
- [ ] One typed result model renders both text and stable JSON carrying schema version, command, outcome, repository root, findings, changes, skipped work and rollback result.
- [ ] Every finding carries a stable code, severity, affected IDs/paths, evidence, fixability, automation-stop reason, exact next argv/Just recipe, and a verification command.
- [ ] Exit codes 0–4 are produced by exactly one authority; no second exit-code field survives.
- [ ] Git transaction fixtures prove exact-path refusal on dirty targets, refusal when the index is non-empty under `--commit`, pre-commit rollback, and post-commit `git revert <sha>` reporting that never rewrites history.
- [ ] The launcher builds on demand when the binary is absent or older than its sources, and the gate lanes that run all of this are proven non-decorative.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go` (new)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/maintainer-verify.sh` (modified)

**What was done:** Completed the preserved `do-work-cli` foundation rather than restarting it. Collapsed the module to a single exit-code authority by deleting `TransactionResult.ExitCode` and the duplicate `RollbackResult`/`RollbackStatus` pair, leaving `resultmodel.ExitCode` as the only outcome-to-number mapping. Added `BuildCommandResult`, which turns a Git transaction into the one typed `CommandResult` with finding codes derived from `FailureKind` and a loud `GIT-UNMAPPED-FAILURE` fallback rather than a hand-maintained switch. Fixed two real defects in the preserved Git layer: `runGit` used `CombinedOutput()`, so a git stderr warning read as porcelain content and refused a clean target; and `rollbackFailure` inspected the whole index, so an unrelated staged file turned a successful rollback into a reported risk. Registered the previously-uninvoked launcher probe in `contract-regressions.sh` and added `do-work-cli` vet and uncached-test lanes to `maintainer-verify.sh`, moving all four of its self-test enumerations in lock-step. No command is registered — `main.go` is untouched and still delegates to a nil handler map; REQ-407 onward register real commands.

**Builder branch:** `worktree-agent-REQ-406-create-do-work-cli-foundation` — `571bb0d`, `737dd20`, `33ff757`, `84f77c4`.
**Merge range:** `ad354e2..767f425`.
**Hand-back:** `do-work/runs/work-2026-08-29-213539/REQ-406-handback.md` (builder-authored `## Decisions`, `## Discovered Tasks`, `## P-A-U` and `## Testing` are read from there).

## Testing

**Tests run:**
- `cd skills/do-work/tools/do-work-cli && go test -count=1 ./...` — ✓ exit 0
- `bash _dev/tests/do-work-cli-launcher-behavior.sh` — ✓ exit 0
- `bash _dev/tests/maintainer-verify.sh` (canonical repository gate, optional browser lane in its default skipped state) — ✓ exit 0 on the **merged** tree, 7m32s, with `do-work-cli go vet`, `do-work-cli uncached tests` and `do-work-cli launcher behavior tests passed` all visibly running

**Result:** ✓ All passing. Judged by direct exit status, never piped.

**Red-green validation:** traced to the REQ's `## Red-Green Proof`, whose GREEN condition is that one typed result drives stable text and JSON, the documented exit codes are observed, and Git fixtures prove exact-path refusal, rollback and commit behaviour.

- `TestGitStderrWarningIsNotReadAsTargetDirt`: ✗ before — assertion failure, exit 1 with `dirty_target` on a *clean* target because `runGit` used `CombinedOutput()` → ✓ after
- `TestUnrelatedStagedWorkDoesNotBreakRollback`: ✗ before — assertion failure, exit 4 with a substantively-complete rollback reported `incomplete` because `rollbackFailure` inspected the whole index → ✓ after
- `TestRuntimeFindingsCarryCompleteRemediation`: ✗ before — four assertion failures, each finding carrying an empty `VerificationArgv` against requirement 5 → ✓ after
- `TestExitCodeContractThroughRealGitTransactions`: ✗ before (build failure — `BuildCommandResult` did not exist) → ✓ after, driving real Git fixtures through `Run` and observing exit 0 / 1 `GIT-DIRTY-TARGET` / 3 `GIT-MUTATION-FAILED` / 4 `GIT-COMMITTED-STATE-RISK` in both renderings

The bridge tests' RED is a build failure rather than an assertion failure, because the function under test is what did not exist. The assertion-level RED for the same requirement is `TestRuntimeFindingsCarryCompleteRemediation` above, so the `tdd: true` evidence bar is met by a runnable harness test observed failing before the change and passing after.

**Falsification — the gates and guards are not decorative:**
- Deleting `Outcome: resultmodel.OutcomeSuccess` from result construction reddens two pre-existing tests, so the retargeted assertions really do guard the zero-value flip.
- Deleting the `FailureDirtyIndex` template entry makes the coverage test name the kind that lost it, rather than passing on a plausible-looking finding.
- Deleting the `cli-vet` lane reddens `maintainer-verify.sh --self-test` with `cli-vet ran 0 times`.
- `chmod -x` on the launcher probe makes `contract-regressions.sh` fail naming the file.
- Everything temporarily broken was restored and re-verified.

**New tests added:**
- `TestGitStderrWarningIsNotReadAsTargetDirt`, `TestUnrelatedStagedWorkDoesNotBreakRollback` (lock-ins for the two fixed defects)
- `TestFindingCodeIsDerivedFromTheFailureKind`, `TestEveryDeclaredFailureKindProducesACompleteFinding`, `TestUnmappedFailureKindFailsLoudly`, `TestSuccessfulTransactionCarriesTruthfulChangeKinds`, `TestCompletedRollbackReportsNoSurvivingChanges`
- `TestExitCodeContractThroughRealGitTransactions`, `TestRuntimeFindingsCarryCompleteRemediation`

*Verified by work action*

## Review

**Acceptance: Pass.** The REQ's GREEN condition holds: one typed result drives stable text and JSON, the documented exit codes 0–4 are observed through real Git fixtures, and those fixtures prove exact-path refusal, rollback, and post-commit risk reporting. The canonical gate `bash _dev/tests/maintainer-verify.sh` exits 0 on the merged tree with the new module lanes running.

**Method.** Five independent reviewers, one per dimension (requirements traceability, correctness, shipped shell and the gate, simplicity and maintainer fit, test quality and the TDD claim). Every finding was then put to three skeptics prompted to refute it, each with a distinct lens: does the code really do this, can the consequence actually occur, and is it already a deliberate decision or out of scope.

**Dimension verdicts:** shell-and-gate pass; requirements, correctness, simplicity and tests each partial.

**Findings that survived verification: one, and it was fixed inside this REQ.**

- **Important — `TestExitCodeContractThroughRealGitTransactions` was flaky at 10%.** Raised independently by two dimensions with different framings, then confirmed empirically by the orchestrator: 4 failures in 40 runs. Each subtest ran its scenario twice in two different fixture repositories and asserted the text run's output against the JSON run's finding; for the commit case that line embeds a commit SHA, and two independent repositories produce different SHAs whenever their commits land in different seconds. The test had already entered the canonical gate, so it would have reddened every later REQ's gate one run in ten. **Remediated** in `fcf1cb5`: one transaction now drives both renderings, which is what the assertion's own comment always claimed. The second fixture is deleted rather than the SHA masked — net −6 lines. Re-measured independently on the merged tree at **0 failures in 100 runs**, where the old rate predicts about ten.

**Findings refuted 3–0, by observation:**

- *Rollback reports success while a created target survives on disk.* Refuted. `recorder.createdPaths` has exactly one writer, which refuses any path outside the declared targets or already present at preflight, and any failed removal forces `RollbackIncomplete` and exit 4. A nine-scenario probe confirmed every recorded creation is removed or named in an incomplete rollback. The described state needs a caller that writes a file and never records it — caller misuse of an enforcing API, and there is no such caller.
- *An incomplete pre-commit rollback is labelled `committed_state_risk` when nothing was committed.* Refuted. The mechanism is real but the mapping is a written requirement of the source UR that this REQ was told to preserve.
- *`repository_root` echoes the `--repo-root` argument rather than the resolved root.* Refuted on both fact and reachability.
- *The failure-kind staleness guard misses the idiomatic Go constant form.* Refuted — the headline is backwards: the idiomatic form is exactly the one the regex catches, and both named alternatives are independently blocked.
- *The requirement-5 completeness test injects the paths it asserts.* Refuted — false about the test it names, and the state it treats as a latent bug is correct by design.

**A note on evidence completeness.** The first verification pass was cut short by a session usage limit that killed 72 of its 92 agents. The workflow then counted an errored verifier as a refutation, which would have reported 24 unadjudicated findings as knocked down. That post-processing bug was fixed, the five **Important** findings were re-verified to completion, and two more were settled directly by the orchestrator. Minor and Nit findings go in this report only and were deliberately not adjudicated — they create no follow-up REQs.

**Minor and Nit findings, recorded not actioned:** the `<command>` placeholder in usage findings is not a runnable argv until a real command is registered; `RollbackResult.Status` has a fourth wire value `""` that no constant names; text rendering of changes, skipped work and rollback errors has no direct assertion; the bridge's RED was compile-error-only (the assertion-level RED for the same requirement is `TestRuntimeFindingsCarryCompleteRemediation`); no test observes a successful `--commit` transaction. Each is folded onto REQ-407, which registers the first real command and is where they become testable.

**Scope drift:** none. `tools/checks/scope-drift.sh` reports the Implementation Summary matches the Scope declaration exactly — ten files, no undeclared touch, no unused declaration.

*Reviewed by work action (five dimensions, three refutation lenses per finding)*

## Lessons Learned

**What worked:** Inspecting the preserved foundation rather than restarting it — the partial commit `329c55a9` was sound, and the diff stayed additive and corrective. Deriving finding codes from `FailureKind` with a loud unmapped fallback, instead of a hand-written switch, is the "state conditions, not lists" rule applied where it actually pays: a ninth failure kind cannot silently ship without a template. Proving each new gate lane falsifiable — deleting the lane and watching `--self-test` go red — caught that a gate you cannot make fail is not a gate.

**What didn't:** The exit-code contract test's first shape ran each scenario twice in two fixture repositories to compare renderings. It looked like a parity assertion and was not one: comparing run-specific data across two runs made it fail whenever two commits landed in different seconds. Two reviewers found it and a 40-run loop settled it. The lesson is narrower than "beware flaky tests": **a parity assertion must render one value twice, never run one scenario twice.** The moment a test needs two executions to compare two representations, the thing it is comparing is no longer the representation.

**Worth knowing:** `resultmodel.ExitCode` is now the only outcome-to-number mapping in the module — keep it that way; the second authority this REQ deleted had already drifted. `RollbackResult.Status` has a fourth wire value, the empty string, for results that never ran a Git transaction; a JSON consumer switching on it must handle `""` alongside the three constants. The success path at `git_transaction.go:161-166` detects an unrecorded change to a declared target but does not consult `state.existed`, so a file created without `RecordCreated` reports `succeeded` — harmless today because no command is registered, and a hardening note for whoever registers the first one.

## Orientation

The suite now has one Go command module underneath it. `do-work-cli` gives every future command a shared runtime: global `--repo-root` and `--format text|json`, one typed result that renders both forms from a single value, exit codes 0–4 decided in exactly one place, and a Git transaction layer that refuses to touch work already in flight, rolls back what it created, and reports `git revert <sha>` rather than ever rewriting history. It lives under the installed core package at `skills/do-work/tools/do-work-cli/`, reached by the on-demand build launcher beside it.

No command is registered yet — that is REQ-407's job, and this REQ deliberately stops at the seam. What changed for a reader is the shape of the ground the next fourteen REQs stand on: shell utilities that each invented their own output format and exit conventions now have one contract to migrate onto.

`[MAP CHANGED]` — this adds a module the suite did not have, and the canonical gate now runs it. `_dev/primes/prime-shell-commands.md` was spot-checked and its referenced paths still resolve; it stays accurate because it governs shipped shell, which this REQ did not change in character. A `prime-do-work-cli.md` will be earned once commands exist to route between — creating one now would index an empty room.

