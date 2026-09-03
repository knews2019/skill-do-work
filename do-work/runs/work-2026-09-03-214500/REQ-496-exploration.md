# REQ-496 Exploration — Shared Executable Already-Green Repair Validator

## Scope and authoritative inputs

This exploration covers `REQ-496-add-shared-executable-already-green-repair-validator.md` only. It was derived from the authoritative working request, its durable brief, `CLAUDE.md`, the action-file and CLI primes/lessons, the `work`, `review-work`, and `work-reference` action contracts, current CLI packages and tests, the REQ-494 archive, and the REQ-494 contract-regression fixture. It proposes no lifecycle, release, board, or source mutation.

The request requires one executable authority for both decisions made around an already-green repository-gate repair:

1. whether `work` may bypass ordinary RED-first TDD; and
2. whether `review-work` may accept the intentionally empty implementation diff.

That authority must compare the no-op claim to actual repair-intake and recorded-gate evidence, derive allowed staged paths from the exact successful canonical-completion result, and fail closed for neighboring or malformed states.

## Root cause

The shipped action files and the executable regression test do not share an authority today.

- `skills/do-work/actions/work.md`, Step 6.5, restates the exception predicate in prose and lets the action decide whether RED-first TDD may be bypassed.
- `skills/do-work/actions/review-work.md`, Step 1, independently restates the predicate and lets review decide whether an empty implementation diff is valid.
- `_dev/tests/contract-regressions.sh` contains a Python `action_decisions(repository_root, request_path)` oracle. It reparses the request, runs the gate argv, examines Git state, and returns separate `tdd_allowed` and `review_allowed` booleans without invoking shipped action authority.
- The fixture self-authors `gate-red-fingerprint` as both expectation and evidence instead of comparing the no-op expectation to a canonical `## Repository Gate Repair Intake` section.
- The fixture runs the gate directly rather than proving a matching past-revision record through `check-green-gate --at-revision`.
- Its review-stage allowlist accepts the request path, checkpoint/calibration paths, or *any* path with prefix `do-work/archive/`. Therefore an unrelated staged archive can be accepted.
- The prose is mutation-tested through `already_green_noop_defects`, but those phrase checks do not prove that either action calls an executable validator or consumes one typed result.

The result is three drifting predicates: work prose, review prose, and the contract fixture's parallel Python implementation. REQ-496 must delete the decision oracle, centralize observation and projection, and make both consumers depend on that single executable result.

## Existing authorities and integration seams

### Repair intake

`skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` is the canonical writer for `## Repository Gate Repair Intake`. The validator must read, not reproduce or replace, this authority. Relevant fields are:

- `Parent`
- `Gate command (argv JSON)`
- `Direct exit status` (the intake failure)
- `Diagnostic fingerprint`
- `Repair dependency`
- diagnostic evidence and optional implementation-base/merge context

Folded repair requests can contain more than one intake section. Every section used for the repair must agree on the fingerprint and gate argv. Conflicting or structurally ambiguous intake is ineligible.

### Already-green no-op evidence

`skills/do-work/actions/work-reference.md`, “Already-green repair no-op completion,” defines the durable evidence shape:

- `## Repository Gate Repair No-Op`
- `Expected diagnostic fingerprint`
- `Gate command (argv JSON)`
- `Direct exit status: 0`
- `Recorded green revision` produced by `record-green-gate`
- `Observed result`
- `Verified at`
- the exact implementation-summary and qualification statements

The expected fingerprint must match the intake fingerprint. It must never be accepted from a command-line argument or because the request repeats the same self-asserted token in two non-authoritative places. Gate argv and recorded revision must likewise be extracted from the request evidence, not supplied as trusted validator inputs.

### Recorded gate evidence

`skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go` already has the correct past-revision comparison implementation, currently private as `checkGreenGateAtRevision`. It returns typed `GateEvidenceResult`; `internal/gateevidence/gate_commands.go` exposes it through `check-green-gate --at-revision`.

Export this function as `CheckGreenGateAtRevision` and reuse it. The validator must not run the repository gate at current `HEAD`: REQ-518 established that the evidence must target the recorded past revision. Eligibility requires the extracted gate argv and recorded revision to return `Matches: true`.

### Canonical completion and exact staged paths

`internal/requeststate` already owns completion planning and result projection:

- `requeststate.BuildPlan`
- `requeststate.ApplyPlan`
- `requeststate.StateOptions{Transition: Complete, RequestID, RequestPath, TerminalStatus: "completed", WriterLabel, DryRun: true}`
- `resultmodel.CommandResult.Changes` and `requeststate.StatePlan.TargetPaths`

The validator should build the dependency graph with `dependencygraph.BuildGraph`, discover the repository through `repositorymodel.DiscoverRepository`, and run canonical completion in dry-run mode with the exact request, writer, and time that the caller will use for completion. A successful dry-run result is the sole allowed-path projection. It can include the exact working source and archive destination, exact user-request moves when the UR closes, and checkpoint/calibration files only when the canonical plan actually includes them.

This removes the unsafe `do-work/archive/` prefix rule. Every staged path must be an exact member of the successful dry-run result's path set. The staged set may be empty. The projection is not proof that the later mutation succeeded; the contract test must still run real `complete --commit` and verify its lifecycle and metadata commits.

`WriterLabel` is significant because it affects checkpoint planning. The command therefore needs an explicit `--writer`, and actions must reuse the same writer for validation and completion. A supplied `--at` must likewise be the exact RFC3339 time reused by completion; otherwise the validator should use one captured time and expose it in its result.

### Git-state observations

The validator needs two independent, NUL-safe observations:

- project dirt from `git status --porcelain=v1 -z --untracked-files=all --no-renames`, rejecting any changed path outside `do-work/`; and
- staged paths from `git diff --cached --name-only -z --no-renames`, requiring each path to be an exact member of the canonical-completion dry-run path set.

It must also reject `release_at` in the request and release mutations such as `VERSION`, `CHANGELOG.md`, or generated release state. Rename source and destination must both remain visible through `--no-renames` and must be validated independently.

### Consumers

There are exactly two decision consumers:

| Consumer | Decision consumed | Required behavior |
| --- | --- | --- |
| `skills/do-work/actions/work.md` | `tdd_allowed` | Invoke the validator after writing/validating the no-op evidence. RED-first TDD can be bypassed only when the typed command result succeeds and this field is true. No prose fallback. |
| `skills/do-work/actions/review-work.md` | `review_allowed` | Invoke the validator fresh during review. The empty implementation diff can be accepted only when the typed command result succeeds and this field is true. Report typed blocker codes/paths. No manual predicate or prose fallback. |

`work-reference.md` defines the durable no-op evidence and invocation contract. `command-line-guide.md` documents operator-visible command behavior. No board/browser consumer exists for this decision.

## Proposed executable contract

Add a dedicated `internal/repairvalidation` package. This is a cross-domain read-only authority composing request evidence, recorded gate evidence, canonical completion planning, and Git state; it should not be hidden in generic `corehelpers` or folded into `gateevidence`.

Proposed command:

```text
<skill-root>/tools/do-work-cli.sh \
  --repo-root <project-root> \
  --format json \
  validate-already-green-repair \
  --request-path do-work/working/REQ-NNN-...md \
  --writer <exact-completion-writer> \
  [--at <RFC3339>]
```

Do not accept fingerprint, gate argv, recorded revision, or allowed lifecycle paths as arguments. Those are evidence or authority outputs, not caller choices.

Add one typed result, for example `AlreadyGreenRepairValidation`, to `resultmodel.CommandResult` as `AlreadyGreenRepair`. It should project at least:

- request ID and contained request path;
- `tdd_allowed` and `review_allowed`;
- intake and expected fingerprints;
- exact gate argv and recorded revision;
- gate-evidence match result (also attach the existing `GateEvidenceResult` to `CommandResult.GateEvidence`);
- canonical completion paths;
- staged paths and non-lifecycle project paths;
- stable blocker/reason codes and offending paths;
- the effective writer and timestamp used for planning.

The command should distinguish observation failure from ineligibility. Malformed input, unreadable state, Git failure, or completion-planning failure is a refused/failed command with typed reason. A well-observed negative can return a successful command envelope with the relevant boolean false, consistent with `check-green-gate` returning `Matches: false`. Each action must require both a successful typed result and its own true decision field.

### Eligibility calculation

Core eligibility, used by both decisions, requires all of the following:

1. The request is one exact contained regular file under `do-work/working/`, with one unambiguous `REQ-NNN` identity.
2. Request front matter has exact `status: claimed`, `repository_gate_repair: true`, and the repair exception's required TDD state; malformed, duplicated, or prefix-matching fields fail closed.
3. The no-op section, implementation summary, and qualification have their exact canonical shapes, with no duplicate labels or headings.
4. Every canonical repair-intake section agrees on a unique diagnostic fingerprint and nonempty JSON gate argv.
5. The no-op expected fingerprint equals the intake fingerprint, and its gate argv equals intake argv.
6. No-op evidence records exit status zero, the exact observed-green claim, a parseable canonical verification time, and a nonempty recorded green revision.
7. `gateevidence.CheckGreenGateAtRevision(repo, argv, recordedRevision)` succeeds with `Matches: true`.
8. There is no changed project path outside `do-work/`.
9. There is no `release_at` and no release mutation.

`tdd_allowed` equals this core result.

`review_allowed` additionally requires a successful canonical-completion dry-run and requires every staged path to be an exact member of that run's result paths. This separation preserves a useful distinction: a review-only staging defect can make review ineligible without retroactively changing whether the implementation was entitled to the TDD exception.

## Requirement-mapped implementation plan

1. **Create the shared authority.** Add `internal/repairvalidation/already_green.go` with strict section/label parsing, intake-to-no-op comparison, Git observations, past-revision gate-evidence reuse, canonical completion dry-run, exact staged-subset validation, typed blocker projection, and the command handler.
2. **Expose rather than duplicate recorded evidence logic.** Export `gateevidence.CheckGreenGateAtRevision` and route the existing command handler through it, preserving current behavior and tests.
3. **Project one typed result.** Extend `resultmodel.CommandResult` and text rendering so JSON and text expose the same decision, evidence, path, and blocker semantics.
4. **Register the command.** Register `validate-already-green-repair` in `cmd/do-work-cli/main.go` with the standard repo-root/format envelope and strict request-path, writer, and optional time flags.
5. **Replace both prose decisions.** Make `work.md` consume only `tdd_allowed`; make `review-work.md` consume only `review_allowed`. Update `work-reference.md` to say that the validator extracts evidence and internally derives completion paths. Neither action may reconstruct the predicate.
6. **Document and prime the authority.** Add command usage and refusal semantics to the CLI guide and add `internal/repairvalidation` ownership/testing guidance to the CLI prime.
7. **Replace the parallel contract oracle.** Delete `action_decisions()` from `_dev/tests/contract-regressions.sh`. Construct a real repair intake, record actual green evidence, write a no-op referencing its past revision, and call the real CLI validator for every positive and negative decision.
8. **Preserve the real lifecycle tail.** After the positive validation, continue to run real `complete --commit` and assert exact lifecycle commit paths, archive metadata commit, stored `commit:`, absent `release_at`, unchanged release files, clean repository, and selector behavior before and after repair completion.
9. **Make action wiring mutation-sensitive.** Contract checks must fail if either action's validator invocation is removed/corrupted or if it stops consuming its respective typed boolean. Phrase-level prose checks are not an authority substitute.

## Exact proposed write set

Implementation must stay within this set unless scope drift is reported and approved:

1. `skills/do-work/actions/work.md`
2. `skills/do-work/actions/review-work.md`
3. `skills/do-work/actions/work-reference.md`
4. `skills/do-work/docs/command-line-guide.md`
5. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
6. `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
7. `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go` (new)
8. `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green_test.go` (new)
9. `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go`
10. `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go`
11. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
12. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
13. `_dev/tests/contract-regressions.sh`

No separate binary integration test is required initially: the contract test must invoke the real launcher and therefore covers command registration and serialization, while package tests cover the validation matrix. Add another file only if implementation reveals a seam that cannot be tested within these named files; treat that as scope drift.

## Literal RED/GREEN verification plan

Write the package and contract tests first. Suggested unit cases:

- `TestValidateAlreadyGreenRepairUsesIntakeAndRecordedEvidence`
- `TestValidateAlreadyGreenRepairRejectsNeighboringShapes`
- `TestValidateAlreadyGreenRepairUsesCanonicalCompletionPathsForStageAuthority`
- `TestAlreadyGreenRepairCommandRendersTypedParity`

The neighboring-shape table must include ordinary non-repair requests; malformed or duplicate history, summary, qualification, intake, and no-op evidence; missing/invalid recorded revision; intake fingerprint mismatch; gate-argv mismatch; nonempty project diff; release mutation; unrelated staged archive; exact allowed lifecycle staging; and exact allowed staging plus one extra path.

RED evidence:

```sh
cd skills/do-work/tools/do-work-cli
go test -count=1 ./internal/repairvalidation ./internal/gateevidence ./internal/resultmodel
cd ../../../..
bash _dev/tests/contract-regressions.sh
```

The first must initially fail because the new validator/typed projection does not exist; the contract must fail because the command is unknown, the parallel `action_decisions()` oracle is forbidden, or the actions do not consume the shared result. Capture literal failures before implementation.

GREEN evidence:

```sh
cd skills/do-work/tools/do-work-cli
go test -count=1 ./internal/repairvalidation ./internal/gateevidence ./internal/resultmodel
go vet ./...
go test -count=1 ./...
cd ../../../..
bash _dev/tests/contract-regressions.sh
git diff --check
```

Run any prime-required maintainer verification after the focused and full module suites. No browser harness is applicable.

## Contract fixture details

The REQ-494 regression block should create an actual `## Repository Gate Repair Intake`, initialize and commit the fixture repository, invoke `record-green-gate` through `tools/do-work-cli.sh`, and capture its revision. Commit the no-op evidence *after* that record so the validator demonstrably checks a past revision rather than `HEAD`.

Each decision assertion must invoke `validate-already-green-repair` and inspect the typed JSON. Positive coverage must include an empty staged set and an exact subset of canonical completion paths. Negative coverage must include an unrelated `do-work/archive/...` path alone and an otherwise exact allowed set plus that extra path. These cases specifically prove removal of the prefix allowance.

The existing real completion assertions remain valuable and must not be weakened: exact lifecycle commit path set, exact archive metadata update, recorded lifecycle commit, absent release timestamp, unchanged version/changelog, clean working tree, and parent selector visibility before versus after repair completion.

## Risks and safeguards

- **Self-asserted evidence:** Never take fingerprint, argv, revision, or allowed paths from caller flags. Parse them from canonical evidence and authorities.
- **Current-HEAD re-execution:** The validator must not rerun the gate. It must compare recorded evidence at the extracted past revision.
- **Ambiguous Markdown:** Anchor exact level-two headings and exact labels; reject duplicate or prefix-like sections/fields rather than choosing the first.
- **Folded intake:** Multiple canonical intake sections are valid only when their fingerprint and argv agree. Conflicts fail closed.
- **Path broadening:** Use NUL-safe Git output, `--no-renames`, normalized repository-relative exact paths, and exact set membership. No directory-prefix or basename heuristic is acceptable.
- **Planner drift:** Obtain allowed paths from the same `requeststate` completion machinery as the real mutation, with identical writer/time inputs. Test checkpoint present/absent and user-request closure/calibration cases where the fixture can do so without widening scope.
- **Dry-run overclaim:** A dry-run is authoritative only for the proposed result paths. It does not prove mutation or commit behavior; the real completion integration remains mandatory.
- **Decision collapse:** Keep one observation/result with two projections. A stage-only blocker should affect `review_allowed` but need not make `tdd_allowed` false.
- **Prose fallback:** Both action files must treat absent, failed, malformed, or false typed output as refusal. There is no manual escape hatch.
- **Runtime growth:** The contract already builds/uses the CLI launcher; reuse its cached binary and keep the negative matrix within the fixture rather than creating many new repositories.

## Explicit exclusions

- No changes to lifecycle metadata, the run manifest, other REQs, versions, changelog, release lock/state, or release commands.
- No changes to the repair-intake writer or defer publication behavior.
- No changes to finalization mutation semantics, commit layout, dependency selection, or `next` ordering.
- No request schema or board/browser/UI changes.
- No new public Just recipe; this is an action-facing CLI command invoked through the existing launcher.
- No generalized Markdown parser, generic policy engine, or refactor of unrelated Git helpers.
- No separate reimplementation in contract tests or action prose.

## Scope-freeze conclusion

The narrow implementation is a new read-only `repairvalidation` command/package that reuses recorded gate evidence and canonical request-state planning, emits one typed result with `tdd_allowed` and `review_allowed`, and replaces all three current decision sites (work prose, review prose, test oracle). The 13-file write set above is sufficient; any need to change publication, lifecycle mutation, selector, schema, board, release, or additional files is scope drift and should stop implementation for review.
