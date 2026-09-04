# REQ-507 Builder Handback

Status: complete

Commit: `7faafb9d3fbf04b26a4599a5141b81f3b476c6fa` (`[REQ-507] route finalization through advance`)

Branch: `worktree-agent-REQ-507-hand-archive-and-commit-tails-to-finalize`

## Implementation Summary

- Changed the terminal working-request phase from action-owned manifest preparation to the mechanical `finalize` phase exposed through the exact public `advance` continuation.
- Added strict parsing for one non-empty action-authored manifest, optional exact outer request path, duplicate/unknown/hostile-input refusal, and no-mutation missing-input projection.
- Added `finalization.FinalizeBound`, which compares the selected request ID/path with the decoded manifest immediately after the finalizer's single manifest decode and before index inspection, journal creation, lifecycle planning, Git, or repository mutation. The existing finalizer remains the sole archive, release, commit, provenance, verification, rollback, and journal authority.
- Preserved the subordinate finalizer outcome, findings, changes, rollback, singular compatibility record, and ordered `finalizations`; `advance` adds only its bound request/phase identity.
- Rendered every ordered finalization record in text with identity, phase/status, archive/journal, resume/discovery state, commit paths/hashes, blockers/reasons, and next/verification/collection argv. JSON normalization remains the schema authority.
- Reduced Step 8/9 and the reference procedures to Fold-First, impact, terminal/failure, release-content, lesson, provenance, typed acceptance, and post-success cleanup judgment. Added section-scoped contract guards and updated the CLI prime.

## RED Evidence

Public tests were added before production changes.

- `go test -count=1 ./internal/lifecycleadvance -run 'TestAdvanceFinalization(RunsTerminalPathMatrix|RequiresOneBoundManifestWithoutMutation)'` failed because the old public `advance` endpoint remained at `agent judgment: prepare finalization manifest` and rejected/failed to execute the new `--finalization-manifest` continuation. This proved the missing advance-owned terminal composition rather than a fixture defect.
- `go test -count=1 ./internal/lifecycleadvance -run TestAdvanceCommandPhaseMatrix` failed because the ready/oriented working request projected the old agent-judgment phase instead of mechanical `finalize` and its exact continuation.
- `go test -count=1 ./internal/resultmodel -run TestFinalizationTextAndJSONCarryTheSameOrderedEvidence` failed because text output emitted no finalization evidence while JSON already carried the normalized record.

## GREEN Evidence

- Public terminal/refusal matrix: `go test -count=1 ./internal/lifecycleadvance -run 'TestAdvanceFinalization(RunsTerminalPathMatrix|RequiresOneBoundManifestWithoutMutation)'` — PASS, 4.005s.
- Final focused packages after the last parser/assertion edits: `go test -count=1 ./internal/lifecycleadvance ./internal/resultmodel` — PASS, 11.022s / 0.343s.
- Existing finalizer regression suite: `go test -count=1 ./internal/finalization` — PASS, 32.230s. The new public test group remains under the 30s fast-test budget.
- Race: `go test -race ./internal/lifecycleadvance ./internal/finalization ./internal/resultmodel` — PASS, 13.192s / 37.872s / 1.742s; final post-review race on the changed lifecycle/result packages — PASS, 20.816s / 2.719s.
- Full module: `go test -count=1 ./...` — PASS after final changes; `internal/finalization` was the longest package at 42.640s, while the repository budget harness later reported the slowest individual test file at 18.72s.
- Vet: `go vet ./...` — PASS.
- Formatting: changed-file `gofmt -l` empty; `git diff --check` — PASS.
- Contracts: `bash _dev/tests/contracts/core-checks.sh` — PASS; `bash _dev/tests/contract-regressions.sh` — PASS in 15.38s with every test file below 30s.
- Canonical gate: `bash _dev/tests/maintainer-verify.sh` — PASS. It included ShellCheck, gofmt, aggregate contracts, both Go-module vets/tests, and reported do-work-cli slowest test file 18.72s.

## Public Path Coverage

- Serial `primary_commit`: publishes a real optional release payload, asserts exact `VERSION` and `CHANGELOG.md` bytes, `release_at`, release files in the finalization allowlist, archive/provenance, matching created/settled hashes, and a clean tree.
- Supplied worktree: commits `implementation.txt` before finalization, supplies that resolvable hash, proves those implementation bytes exist at the supplied commit, proves `implementation.txt` is absent from the tail commit allowlist, asserts the archive records the supplied hash, and asserts no metadata commit.
- Completed with issues: runs the full public transaction and explicitly asserts both result `terminal_status` and archived frontmatter bytes retain `completed-with-issues`.
- Already green/no release: asserts byte-identical `VERSION` and `CHANGELOG.md`, no `release_at`, lifecycle-only finalization paths, provenance, and clean completion.
- Refusals/input: missing manifest, empty manifest, duplicate manifest, hostile/unknown token, outer path mismatch, manifest ID mismatch, and manifest path mismatch. Each captures repository bytes and HEAD before invocation and proves no mutation; the hostile token is also proven not to execute.
- Text/JSON parity: asserts all typed fields, singular compatibility behavior for one record, non-null collections, and actual two-record order in both text and JSON output.

## Exact Changed Files

1. `_dev/tests/contracts/core-checks.sh`
2. `skills/do-work/actions/work-reference.md`
3. `skills/do-work/actions/work.md`
4. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go`
5. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go`
6. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go`
7. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go`
8. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate.go`
9. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate_test.go`
10. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
11. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
12. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`

No file outside the canonical 12-file scope was changed on the builder branch, and nothing under `do-work/` was edited there.

## Decisions

- Kept direct `finalize` behavior compatible while exposing a separate bound in-process entry for `advance`; both use the same preparation/transaction path.
- Used a typed private binding error so only selected-request identity mismatch maps to `FINALIZATION-REQUEST-MISMATCH`; all other preparation failures retain the existing `FINALIZATION-PREPARE` semantics.
- Treated absent manifest input as the phase's typed missing-evidence projection, while empty, duplicate, malformed, irrelevant, mismatched, or hostile supplied input refuses before finalization.
- Kept negative prose predicates section-scoped so legitimate earlier handback merge mechanics are not globally banned.

## Discovered Tasks

- None.

## Clean-Status Proof

After commit:

```text
$ git status --short --branch
## worktree-agent-REQ-507-hand-archive-and-commit-tails-to-finalize
```

There were no staged, unstaged, or untracked files in the isolated worktree.
