# REQ-518 Exploration

## Findings

F1. Command registration has one intended seam: `cmd/do-work-cli/main.go:main` imports each command family and merges its `Handlers()` map into the runtime map. Add `internal/gateevidence.Handlers()` there. `internal/commandruntime.CommandRuntime.Run` supplies the final `command` and `repository_root`, selects text/JSON once, and owns the outcome-to-exit-code mapping through `resultmodel.ExitCode`; the new handlers should return typed results and must not render or exit themselves.

F2. `internal/resultmodel/result_model.go:CommandResult` uses optional typed pointers for domain projections (`AuditMetrics`, `GateDeferral`, `Finalization`). `NormalizeResult` initializes nested slices, and `renderText` prints the same typed object used by JSON. `GateEvidence *GateEvidenceResult` should follow that pattern, normalize `GateCommand` to `[]`, and render before generic skipped/selection output. `internal/resultmodel/result_model_test.go:TestRenderersUseOneNormalizedResult` and the typed selection/deferral tests are the patterns for text/JSON parity.

F3. The reusable publication primitives are `internal/atomicfile.CreateExclusive` and `ReplaceExisting`. `CreateExclusive` writes/syncs a new regular file with the requested complete mode; `ReplaceExisting` rejects symlinks/special files, preserves the existing complete mode, detects observed identity/content changes, and atomically replaces from the same directory. Reuse the `publishGeneratedImage` branch shape in `internal/toolboxcommands/report_image.go`: `Lstat`, refuse non-regular targets, replace an existing regular file, create exclusively on `IsNotExist`, and surface every other error. Do not reuse `corehelpers.writePrivateAtomic`; it is unexported and lacks `atomicfile`'s identity checks.

F4. `corehelpers.gitOutput` and its argument parser/finding builders are package-private. A new `internal/gateevidence` package therefore needs a narrow local Git runner. Use argument arrays with `git -C <root>` and preserve command failure as an error; do not import `corehelpers` or duplicate its broad helper inventory. The existing `CommandRuntime` passes every token after the command unchanged, including `--`, so the handlers can require `--` and preserve exact gate argv without shell parsing.

F5. The planned state location needs one correction. In a linked worktree, `git rev-parse --git-path do-work-green-gates` resolves under `.git/worktrees/<name>/`, while `git rev-parse --path-format=absolute --git-common-dir` resolves to the shared repository Git directory. If the evidence is repository-wide and must be reusable between integration and linked worktrees, resolve the common directory and place `do-work-green-gates/<argv-digest>.json` beneath it. Bind the stored `repository_identity` to the canonicalized common-directory path. Add a linked-worktree test so a future switch back to per-worktree `--git-path` cannot silently split the cache.

F6. The history proof should reuse the shell prime's safe Git path primitive in Go form: first require `git merge-base --is-ancestor <recorded> HEAD`, enumerate all commits in `<recorded>..HEAD` with `git rev-list`, then inspect each with `git diff-tree --no-commit-id --name-only -r -m -z <commit>`. Accept empty commits and commits whose every emitted path is below `_dev/gate-runs/`; reject the first other path. Do not use a net `git diff`, because a later revert could erase evidence that a non-log commit occurred.

F7. The recorder cannot itself prove an external gate ran; `--gate-exit-status 0` is caller-supplied evidence from the action. Keep the provenance label accurate (`reported_direct_zero_exit` or equivalent), require exactly zero, and rely on `actions/work.md` to call it only immediately after the direct gate. Avoid a label implying the CLI observed the child process. Wrapping the gate would break the current direct-output/fingerprinting seam and is not recommended in this REQ.

F8. The current repository-gate action owner is `skills/do-work/actions/work.md` under `#### Repository-gate baseline, deferral, and resume`; the final mandatory run is Step 6.5 item 4. The full semantic owner is `skills/do-work/actions/work-reference.md` under `Repository Gate Deferral and Resumption`, specifically `Session state and baseline` and `Late attribution`. Change only those paragraphs/table rows: typed check before the baseline fallback, typed record after each direct zero, no check before Step 6.5, and no change to fingerprinting, deferral, saved-range, or repair-failure branches.

F9. `_dev/tests/contract-regressions.sh:repository_gate_defects` already extracts exactly the baseline/testing/reference sections and has mutation cases for `structured direct baseline` and `late base attribution`. Replace/strengthen those existing predicates and their existing mutation strings; do not add a new sentence predicate. The file is currently 8,440 lines, already above the REQ's captured 8,417 figure, so this change should be line-neutral or shrinking rather than attempting an unrelated 23-line cleanup. REQ-519 owns the absolute ratchet.

F10. The REQ-514 refusal invariant affects finding construction. `corehelpers.usageResult` currently points to the invoking command and is not a safe pattern for the new refusal. Build gate-evidence findings locally: the non-green record refusal may use the exact gate argv as `next_argv` and `check-green-gate` as `verification_argv`; unrecoverable invalid/unsafe record state should be `outcome: failure` or carry an empty next action, never point back to `record-green-gate`/`check-green-gate`.

F11. The black-box RED belongs in `cmd/do-work-cli/gate_evidence_integration_test.go`: build the real command and decode generic JSON so the test compiles before the result type exists. Before registration it should fail an assertion on `UNKNOWN-COMMAND`, then cover record/exact match, project-commit invalidation, log-only descendant match, and the shared linked-worktree location. Unit tests in `internal/gateevidence` should cover malformed identity/argv/revision records, non-ancestry, mixed history, non-zero refusal/non-mutation, unsafe target, and repeat replacement. Result rendering gets its own focused parity test.

F12. There is a known batch sequencing gap, not a reason to widen this REQ: current Step 9 finalization/release commits occur after Step 6.5 and are non-log changes, so they correctly invalidate a just-recorded revision. Preserve that behavior. REQ-519 moves the full release gate to the integrating commit; exempting `do-work/`, version, or changelog paths here would contradict REQ-518's explicit rule that only `_dev/gate-runs/`-only commits are transparent.

## Recommended scope

- Add: `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go`, `gate_commands.go`, and focused tests; `skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go`.
- Modify: `cmd/do-work-cli/main.go`, `internal/resultmodel/result_model.go`, `internal/resultmodel/result_model_test.go`, `prime-do-work-cli.md`, `skills/do-work/actions/work.md`, `skills/do-work/actions/work-reference.md`, and the existing repository-gate block in `_dev/tests/contract-regressions.sh`.
- Do not modify: `internal/corehelpers`, `internal/atomicfile`, request lifecycle/finalization packages, gate scripts, REQ files, or any additional action sections.

## Focused verification

1. RED/GREEN seam: `cd skills/do-work/tools/do-work-cli && go test ./cmd/do-work-cli -run TestGreenGateEvidenceLifecycle -count=1`.
2. Package/result behavior: `cd skills/do-work/tools/do-work-cli && go test ./internal/gateevidence ./internal/resultmodel -count=1`.
3. Race/static/module checks: `cd skills/do-work/tools/do-work-cli && go test -race ./internal/gateevidence && go vet ./... && go test -count=1 ./...`.
4. Root contracts: `bash _dev/tests/do-work-cli-go125-compatibility.sh` and `bash _dev/tests/contract-regressions.sh`.
5. Integrating gate only: unpiped `bash _dev/tests/maintainer-verify.sh` when the owner reaches the repository gate.
