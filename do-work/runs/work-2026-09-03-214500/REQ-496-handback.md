# REQ-496 Builder Handback

## Source result

- Branch: `codex/REQ-496-repair-validator`
- Commit: `5790b0519b75ed59d4458727e5d7dd6fd6b18e2c`
- Commit message: `[REQ-496] add shared already-green repair validator`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-496-add-shared-executable-already-green-repair-validator`
- Worktree status after commit: clean

## What changed

The commit adds `validate-already-green-repair` as the one executable decision authority for both the TDD bypass and no-diff review. It:

- strictly joins the canonical `## Repository Gate Repair Intake` evidence to the exact no-op, implementation-summary, and qualification shapes;
- sources fingerprint and gate argv from repair intake rather than accepting caller assertions;
- exports and reuses past-revision green-gate verification, and requires the no-op's recorded revision to equal the persisted record;
- observes repository and staged Git paths with NUL-safe, no-renames commands;
- derives allowed review staging from the exact `requeststate` canonical completion dry-run result, not an archive prefix;
- emits one typed `already_green_repair` result with separate `tdd_allowed` and `review_allowed` projections, gate evidence, exact completion/staged/project paths, reason codes, offending paths, writer, and planning time;
- makes `work.md` and `review-work.md` consume only their corresponding typed decision and removes their prose fallback;
- deletes the REQ-494 contract fixture's parallel `action_decisions()` oracle and replaces it with real CLI calls over real request/Git/gate-evidence/completion/metadata/selector state.

## Changed files

Exactly the 13 declared write-set paths are in the commit:

1. `_dev/tests/contract-regressions.sh`
2. `skills/do-work/actions/review-work.md`
3. `skills/do-work/actions/work-reference.md`
4. `skills/do-work/actions/work.md`
5. `skills/do-work/docs/command-line-guide.md`
6. `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
7. `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go`
8. `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go`
9. `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go` (new)
10. `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green_test.go` (new)
11. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
12. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
13. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`

No request/lifecycle file, run manifest, version, changelog, release state, board source, or other REQ was edited in the source commit.

## RED evidence

Tests were written before implementation. The focused command:

```sh
cd skills/do-work/tools/do-work-cli
go test -count=1 ./internal/repairvalidation ./internal/gateevidence ./internal/resultmodel
```

failed with exit 1. Literal failures included:

```text
internal/resultmodel/result_model_test.go:193:3: unknown field AlreadyGreenRepair in struct literal of type CommandResult
internal/resultmodel/result_model_test.go:193:24: undefined: AlreadyGreenRepairValidation
internal/repairvalidation/already_green_test.go:18:31: undefined: Validate
internal/repairvalidation/already_green_test.go:18:56: undefined: Options
FAIL github.com/knews2019/skill-do-work/do-work-cli/internal/repairvalidation [build failed]
FAIL github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel [build failed]
```

This is the requested literal failure-before proof that neither the shared authority nor its typed projection existed.

## GREEN evidence

All final commands passed on the committed content immediately before commit:

```sh
cd skills/do-work/tools/do-work-cli
go test -count=1 ./internal/repairvalidation ./internal/gateevidence ./internal/resultmodel
go vet ./...
go test -count=1 ./...

cd ../../../..
bash _dev/tests/contract-regressions.sh
git diff --check
```

Observed results:

- focused validator/gate-evidence/result-model packages: PASS (about 4 seconds for the focused invocation);
- `go vet ./...`: PASS;
- full do-work CLI module: PASS, including the new `internal/repairvalidation` package;
- contract regressions: PASS in about 21 seconds, with all listed standalone shell test files below the 30-second budget;
- diff whitespace check: PASS;
- branch/worktree status after commit: clean.

The contract now proves:

- real `record-green-gate` evidence is created at a revision before the durable no-op commit;
- canonical intake/no-op evidence returns both decisions true;
- ordinary, malformed history/summary/qualification, self-asserted fingerprint, missing revision, project-dirty, and release-mutated neighbors refuse;
- staging the exact working request path reported by canonical completion passes;
- an unrelated staged archive refuses review while preserving the stage-independent TDD projection;
- exact allowed staging plus an unrelated archive refuses with `REPAIR-STAGED-PATH-NOT-CANONICAL`;
- real `complete --commit` still produces the lifecycle commit plus exact archive metadata commit, leaves release files unchanged, records no `release_at`, cleans the repository, and makes the deferred parent selectable;
- corrupting either action's command invocation or typed decision consumption breaks the action contract.

## Heavy-test boundary

`bash _dev/tests/maintainer-verify.sh --heavy-surfaces` reports both `_dev/tests/*.sh` and CLI `*_test.go` as heavy surfaces, so this commit intersects that list. This builder did not run `--heavy` again; the user requested one manual heavy run before starting the chain. Integration should associate this handback with that one-off result if it covered this merged revision, or retain the ordinary permission/hold boundary if another heavy run would be needed.

## Risks and review notes

- The validator is deliberately read-only. It projects completion paths through canonical `requeststate` dry run; the contract's real completion tail remains the proof of mutation, commit metadata, and selector behavior.
- Negative eligibility is a typed successful observation with false decision fields, matching the existing `check-green-gate` model. Usage, discovery, and Git-observation failures are command failures. Both action consumers require typed command success plus their true boolean.
- Folded repair intake is accepted only when every canonical intake agrees on fingerprint and argv; malformed or conflicting intake fails closed.
- Exact Markdown headings/labels, whole request path, status/marker/TDD scalars, canonical UTC verification timestamp, recorded evidence, and no release metadata are required. Duplicate or neighboring shapes fail closed.
- Staged authorization is exact path membership only. It does not claim staged bytes are already the final lifecycle postimage; the finalization engine retains its own preimage/postimage and mutation guards.
- Git observations use `--porcelain=v1 -z --untracked-files=all --no-renames` and `diff --cached --name-only -z --no-renames`, preserving spaces and exposing both rename endpoints.
- The exported `gateevidence.CheckGreenGateAtRevision` preserves the former private wrapper for existing package tests; behavior is unchanged.

## Merge guidance

Cherry-pick or merge commit `5790b0519b75ed59d4458727e5d7dd6fd6b18e2c` into the integration lane. Expected overlap is concentrated in `_dev/tests/contract-regressions.sh`, the three already-green action/reference paragraphs, `cmd/do-work-cli/main.go`, and the shared result model. Preserve the new real validator fixture and do not restore `action_decisions()` or the `do-work/archive/` prefix allowlist when resolving conflicts.

After merge, rerun at minimum:

```sh
cd skills/do-work/tools/do-work-cli
go test -count=1 ./internal/repairvalidation ./internal/gateevidence ./internal/resultmodel
go vet ./...
go test -count=1 ./...
cd ../../../..
bash _dev/tests/contract-regressions.sh
git diff --check
```

No source remediation or out-of-scope follow-up is known at handback.
