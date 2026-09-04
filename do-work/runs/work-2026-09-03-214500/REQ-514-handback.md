# REQ-514 Builder Handback

## Branch and commit

- Branch: `worktree-agent-REQ-514-refusals-never-name-themselves-as-the-fix`
- Commit: `7611d7b8c76abfb06ad1492af7d7771d18f177fa`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-514-refusals-never-name-themselves-as-the-fix`

## Implementation

`NormalizeResult` now detects a refusal finding whose `next_argv` resolves to the command currently being rendered. It removes that false remedy while preserving same-command `verification_argv`. When every refusal blocker is REQ-owned, the normalized outcome becomes `findings`, making the request set-aside evidence instead of a global refusal loop. Unowned/shared blockers remain `refused` but no longer advertise the invoking command as their own fix.

The ambiguous finalization-discovery refusal now names `uncommitted-inventory` as its distinct resolving command and retains `recover-finalization --discover` only as verification.

## Exact file manifest

- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`

## RED / GREEN evidence

- RED: `TestRefusalRemediesNeverNameTheInvokingCommand` observed the REQ-owned `FINALIZATION-LIFECYCLE-APPLY` refusal remain `refused` with `recover-finalization` in `next_argv`.
- RED: `TestDiscoveryRefusalNamesInventoryAsTheResolvingVerb` observed the ambiguous discovery refusal use `recover-finalization` for both remedy and verification.
- GREEN: `go test -count=1 ./internal/resultmodel ./internal/commandruntime ./internal/finalization ./internal/requeststate` passed (`resultmodel` 0.290s, `commandruntime` 0.995s, `finalization` 49.741s, `requeststate` 5.321s).
- GREEN: `go test -count=1 ./internal/gittransaction -run TestCancelledCommitThatLandsReportsCommittedRisk` passed in 1.125s after the first all-package run hit that test's pre-existing timing flake (`hook.pid` was not recorded).
- GREEN: `go vet ./...` passed.
- GREEN: `git diff --check` passed.
- Broader run: `go test -count=1 ./...` passed every package except the one transient `gittransaction` timing failure above; the exact failing test passed on immediate isolated rerun.

## Folded lifecycle fixture

The existing `TestRecoverFinalizationResumesJournalAfterLifecycleInterruption` already exercises the required claim-state fixture, one-line implementation change, complete manifest, lifecycle-phase interruption, `recover-finalization`, and terminal clean recovery. It remained green with this change, so no duplicate fixture was added.

## Required reads

- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (complete)

Applied lesson: normalize typed evidence at the output boundary so every command producer inherits the invariant, while correcting the one discovery producer that has a truthful alternate resolver.

## P-A-U

- PLAN: enforce the invariant centrally; keep verification independent; distinguish REQ-owned set-asides from global blockers; correct ambiguous discovery's resolver.
- APPLY: added command-verb extraction and normalization, changed discovery remedy, and added result-model, runtime, and finalization tests.
- UNIFY: reviewed all five changed files, ran focused and broad Go tests, reran the unrelated flaky test, ran `go vet ./...`, and checked whitespace. No debug artifacts remain.

## Decisions

- Same-command `verification_argv` is explicitly preserved.
- A self-referential REQ-owned refusal becomes `outcome=findings` only when no unowned refusal blocker remains.
- A self-referential global refusal stays `outcome=refused`; only its false `next_argv` is removed.
- Global options before the command verb (`--format`, `--repo-root`) are ignored when comparing verbs.

## Scope contradiction

The declared eight-file scope named `finalization_apply.go` and `finalization_pipeline_dirt_test.go`, but the truthful producer of the ambiguous recovery remedy is `internal/finalization/finalization_discovery.go`. The implementation therefore changed that file instead. No changes were needed in `finalization_apply.go`, `finalization_pipeline_dirt_test.go`, `requeststate`, or command-runtime production code. The owner must expand/reconcile Scope before integration.

## Discovered tasks / integration seams

- No new follow-up request is required.
- Review the normalization boundary against any consumer that reads an unrendered `CommandResult`; this request changes the serialized/rendered contract, matching the runtime path.
- Preserve the existing lifecycle recovery fixture rather than adding a second near-duplicate.

## Bounded remediation handback

- Branch: `worktree-agent-REQ-514-remediation`
- Commit: `b63f911f3957f35fe03b8c31ab34af95bb741974`
- Exact remediation files: `internal/resultmodel/result_model.go`, `internal/resultmodel/result_model_test.go`, `internal/finalization/finalization_discovery.go`, `internal/finalization/finalization_recovery_test.go`, `internal/finalization/finalization_commands_test.go`, `internal/cleanup/cleanup_apply.go`, and `internal/cleanup/cleanup_apply_test.go`, all below `skills/do-work/tools/do-work-cli/`.
- F1 closed: refusal findings are normalized independently of aggregate outcome; the cited cleanup scratch refusal now emits no raw self-remedy; an `OutcomeFindings` regression pins the output boundary.
- F2 closed: singular and ordered discovery finalization records now use inventory collection as their resolving command, with recovery discovery retained as verification.
- F3 closed: the lifecycle recovery fixture now begins with a pending queue REQ and executes the real committed claim handler before the one-line implementation change and interrupted completion/recovery path.
- Verification: focused resultmodel, commandruntime, finalization, requeststate, and cleanup packages passed; `go vet ./...`, `git diff --check`, and `go test -count=1 ./...` passed. Full-module package timings included finalization 41.360s and publication 23.936s; the repository wrapper owns per-file timing enforcement.
