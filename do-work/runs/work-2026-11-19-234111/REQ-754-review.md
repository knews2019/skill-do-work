# REQ-754 independent review

Approve. The merged implementation lets second-opinion subprocesses run with unrelated uncommitted files, checks their temporary worktrees, and retains permitted report output in the real bundle. Session reviews keep report preservation checks. Local acceptance passed.

Reviewed range `02f3a5e1..d1156bc6`. Reviewed the full request, original UR-187 input, builder handback, all six changed files, and the existing publication/process/rooted-files helpers they call.

## Requirements

- [x] Each subprocess run creates a detached temporary worktree before probes or generation. Its working directory and prompt repository paths point into that worktree.
- [x] The current report is materialized from its snapshot. Replacing the checked-out bundle preserves uncommitted additions, edits, deletions and empty directories. Separate snapshots preserve inode-aware validation.
- [x] Outside-path checks run against the isolated checkout after subprocess calls. The original and existing bundle assets retain the preservation checks.
- [x] Exclusive copy-back admits the new staging leaf and new prefixed regular root assets. The existing publication tail creates the final page/reply and links the original only after its checks pass. Concurrent real-bundle changes refuse normal publication while a completed staging page can still be retained exclusively.
- [x] Deferred cleanup is installed before worktree setup. It follows owned-process reaping and uses a fresh bounded context after cancellation. Tests assert both absent registrations and absent temporary directories across success, refusals, setup failures, generation failure, cancellation, deadline and SIGINT. Explicit teardown errors are reported as failed runs.
- [x] Pre-existing outside changes in the real checkout cannot trigger the subprocess check.
- [x] Session prepare/finalize and the shared publication tail no longer inspect outside Git status. Existing session bundle checks remain.
- [x] The action and guide explain isolated outside_paths, retained output, cleanup and unrestricted absolute-path access. Existing backend permission arguments are unchanged.

## Findings

No material findings. No minor findings. No follow-ups created.

The test-file split is explicitly authorized in the builder brief and reflected in the final Scope. The replaced REQ-698 dirty-worktree refusal assertion is intentionally superseded by the required positive dirty-README regression. All P-A-U boxes are checked. The new Go helper is directly called by secondOpinionRun and its tests are discovered through Go's package/test conventions.

## Acceptance evidence

Independent command, run from the CLI module on the merged tree, exited 0 in 2.763 seconds.

```sh
go test ./internal/toolboxcommands -run '^(TestSecondOpinionWorktreeDirtyReadme|TestSecondOpinionSessionIgnoresOutsideDirtThroughPublication|TestSecondOpinionWorktreeSIGINT|TestSecondOpinionWorktreeSetupCancellation)$' -count=1
```

These public-handler and process tests exercised dirty tracked README publication, isolated outside-write refusal with untouched real README, retained refusal output, session publication, SIGINT, setup cancellation, descendant reaping, and disk/registration cleanup. `git diff --check 02f3a5e1..d1156bc6` passed. The builder's RED-before-source-change evidence and final focused/race/gate evidence were read and checked against the changed assertions. The parent owns the ongoing canonical integrated repository gate.

## Restatement sweep

Redefined the second-opinion clean-checkout prerequisite, subprocess working directory, outside_paths meaning, and session finalize coverage. Searched shipped Markdown, Go consumers, README and active docs for these terms and related instructions. The action and guide agree with implementation, the prompt boundary remains compatible with isolated repository paths, and result_model.go only declares the JSON field. Remaining old guarantees are dated historical changelog/lesson records. No stale active shipped instruction tells the caller to commit unrelated work.

## Limits

Acceptance covers local implementation and integration through native backend fixtures. Live paid/model sessions and a consumer installation were unassessed. The new behavior tests use the existing render-judge seam for successful publication, while the judge implementation is unchanged. Windows evidence in the handback is cross-build/vet evidence. Cleanup failure injection confirms an error is surfaced when Git refuses removal, with possible stale registration requiring recovery.

## Scores

| Dimension | Score |
| --- | --- |
| Requirements | 100% |
| Code quality | 100% |
| Test adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass for local implementation and integration |

Overall 100%, computed as the mean of the four percentage dimensions. No acceptance penalty or risk cap applies. The score describes the reviewed scope and observed fixtures.

Self-validation completed against the original request, frozen diff, source call paths and actual command result. Every requested behavior has implementation and regression evidence. No requested source changes remain from this review.

Reviewed at 2026-11-20T00:00:47Z
