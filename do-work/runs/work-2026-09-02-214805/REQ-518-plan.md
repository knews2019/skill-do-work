# REQ-518 Implementation Plan

## Approach and decisions

D1. Add a small `internal/gateevidence` package with two public CLI commands: `check-green-gate` reads durable evidence and `record-green-gate` writes it only when the caller supplies the direct gate status `0`. The action remains responsible for resolving and directly running the project-owned gate, preserving its output and fingerprint procedure; the Go command owns argv identity, repository identity, revision/history validation, and durable storage. This avoids wrapping the gate in a subprocess whose output would collide with the CLI's text/JSON result contract.

D2. Store one record per exact gate argv below the repository's Git-private path returned by `git rev-parse --git-path do-work-green-gates`. Key the filename by SHA-256 of the canonical JSON encoding of the argv array and store the full argv in the record as a collision/tamper check. Records are private `0600` JSON files in a `0700` directory and are created or replaced atomically through `internal/atomicfile`; they are not REQ frontmatter, pipeline state, or tracked project files, and therefore survive a session restart without becoming hand-edited workflow data.

D3. Bind every record to the canonical absolute Git common directory, resolved from the invocation repository, and to the full commit object named by `HEAD` after the successful gate returns. A check must re-derive the repository identity, argv digest, current `HEAD`, and stored commit. It reports a match only when the stored revision is `HEAD`, or when the stored revision is an ancestor of `HEAD` and every path changed by every intervening commit is below `_dev/gate-runs/`. Use NUL-delimited Git path output and inspect each intervening commit, not a net tree diff, so a non-log change that was later reverted cannot disappear from the proof. Missing evidence, a different argv/repository, a missing commit, non-ancestry, or any intervening non-log path is a non-match.

D4. Preserve the exact scope boundary in REQ-518 even though it exposes a temporary batch-order limitation: a Step 9 finalization/release commit after the Step 6.5 gate is a non-log commit and therefore invalidates the record. Do not exempt `do-work/`, version, or changelog paths merely to make the next baseline skip. REQ-519 is the dependent change that moves the full release gate to the integrating commit; once that lands, the next REQ can consume the immediately preceding green revision. REQ-518 still provides the requested skip whenever `HEAD` is already the recorded green revision, including manual green runs and log-only descendants.

D5. Treat the new result as a first-class typed projection rendered from one observation. Valid match and valid non-match checks both return top-level `outcome: success`; callers branch on `gate_evidence.matches`, not process status or rendered prose. Git launch failures, unreadable/unsafe/corrupt records, and atomic publication failures return `outcome: failure` and stop the workflow. A `record-green-gate` call with a non-zero supplied gate status returns `outcome: refused`, makes no change, and names the exact gate argv as its next action, never itself, preserving the REQ-514 self-referential-refusal invariant.

## Command and result contract

The two invocations are:

```text
do-work-cli --repo-root <project-root> --format json check-green-gate -- <gate argv...>
do-work-cli --repo-root <project-root> --format json record-green-gate --gate-exit-status 0 -- <gate argv...>
```

Both require a non-empty argv after the literal `--`; no shell string is parsed or re-split. `record-green-gate` accepts exactly one integer status and refuses every value other than zero without creating or replacing evidence.

Add `GateEvidenceResult` to `internal/resultmodel.CommandResult` as optional `gate_evidence`. Its stable fields are:

- Identity: `repository_identity`, `gate_command` (JSON array), `gate_command_sha256`, and absolute `record_path`.
- Provenance: `record_provenance` (`direct_zero_exit` in a stored record, `persisted_green_run` when read), `gate_exit_status`, `recorded_revision`, and `head_revision`.
- State: `state`, `matches`, `match_basis`, and `baseline_revision`. States cover `recorded`, `missing`, `exact_revision_match`, `gate_log_descendant_match`, `different_repository`, `different_argv`, `recorded_revision_missing`, `recorded_revision_not_ancestor`, and `invalidated_by_non_gate_log_commit`; malformed/unsafe evidence uses `invalid_record` with top-level failure. `match_basis` is `exact_revision`, `gate_log_only_descendant`, or `none`. `baseline_revision` equals current `HEAD` only on a match, so the work action has the exact attribution base to save even when the stored revision is an older log-only ancestor.
- Outcome: retain the canonical top-level `command`, `outcome`, `repository_root`, `findings`, `changes`, and `rollback`. Successful recording reports one `changes` row with kind `git-private`; checks are byte-for-byte read-only and report no change. Normalize `gate_command` to a non-null array and render all gate-evidence fields in text from the same typed object used for JSON.

Finding codes and actions should be stable and specific. Argument/Git/I/O failures carry exact retry or verification argv where a distinct safe action exists. A non-green record refusal uses `GATE-EVIDENCE-NOT-GREEN`, `next_argv` equal to the gate command, and a `check-green-gate` verification argv. No refusal may point back to `record-green-gate`.

## Ordered implementation changes

1. Add the public-seam RED test in `skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go`. Build/run the real CLI against a temporary Git repository, invoke the literal unregistered command names, and assert the typed contract for record then exact match, project-commit invalidation, and a subsequent `_dev/gate-runs/`-only commit match. Decode generic JSON in this first test so it compiles before `GateEvidenceResult` exists. Run it and retain the assertion failure showing `UNKNOWN-COMMAND`; this is the required RED, not a syntax/import failure.

2. Add `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go` for the record schema and deterministic mechanics: parse exact argv/status, resolve and canonicalize Git root/common-dir identity, resolve full `HEAD`, derive the argv hash and Git-private record path, atomically create/replace the record, validate stored schema/identity/argv/revision, prove ancestry, and enumerate intervening commit paths with NUL-safe Git output. Keep Git execution and record parsing in this package; do not rescan do-work request state or import unrelated lifecycle packages.

3. Add `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go` to expose `Handlers()` for `check-green-gate` and `record-green-gate`, translate the package observation into the typed command result, and emit actionable non-self-referential findings. Register those handlers in `cmd/do-work-cli/main.go` with one import and one handler loop, matching the existing command-registration boundary.

4. Extend `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` with the gate-evidence enums/result, normalization, and text rendering. Add focused parity tests in `internal/resultmodel/result_model_test.go` proving text and JSON expose the same identity, provenance, state, match, baseline, and outcome facts and that an empty gate argv normalizes to `[]`.

5. Add `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence_test.go` and `gate_commands_test.go`. Cover exact match, a new project commit, a log-only descendant, mixed log/project history, revision divergence/non-ancestry, a missing record, different argv, a valid record copied from a different repository, a missing recorded object, malformed and non-regular record targets, a non-zero record refusal with zero filesystem effects, repeated green recording replacing the prior revision atomically, text/JSON status parity, and every refusal's `next_argv` verb differing from the invoking command.

6. Update `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` as the low-noise package index: name `internal/gateevidence/` as the owner of Git-private green-gate evidence and add focused verification (`go test ./internal/gateevidence`). Keep implementation details in code and do not copy the record schema into the prime.

7. Update `skills/do-work/actions/work.md` Step 5.75. After resolving the canonical structured argv and before dispatch/source edits, invoke `check-green-gate` and consume only its typed `gate_evidence`. On `matches: true`, save `baseline_revision` as a green baseline and dispatch without running the gate. On a valid `matches: false`, run the existing direct baseline unchanged. A failed check stops safely. Immediately after any direct zero baseline, call `record-green-gate` and require typed success before dispatch; retain the existing red fingerprint/deferral branches untouched.

8. Update `skills/do-work/actions/work.md` Step 6.5 so the final direct gate remains mandatory and is never preceded by the cache check. After every zero final gate, call `record-green-gate` before continuing, which records the post-gate `HEAD` and therefore naturally records REQ-523's gate-log commit when that feature is present. A red/launch-failed gate follows the existing attribution/repair lifecycle and never writes green evidence; a record-command failure stops completion rather than claiming durable evidence was saved.

9. Update `skills/do-work/actions/work-reference.md` under Repository Gate Deferral and Resumption. Define the typed check/record semantics once in Session state and baseline, amend the exhaustive branch table to distinguish recorded-green and freshly-run green baselines, state that `baseline_revision` is the saved attribution base, and state that Step 6.5 always runs then records. Leave fingerprinting, repair non-recursion, detached-base attribution, saved-range drift, and deferral manifests unchanged. Mention the `_dev/gate-runs/`-only descendant rule as the sole transparent commit class, conditionally rather than as a growing list.

10. Adjust, without adding a new predicate or increasing file length, the existing repository-gate baseline predicates and mutation case in `_dev/tests/contract-regressions.sh`. Repurpose `structured direct baseline` to require the typed `check-green-gate` match/fallback branch, and strengthen the existing final-gate predicate to require `record-green-gate` only after a direct zero exit. Keep the behavioral matrix in Go; the shell contract test only ensures the action still delegates at the correct two seams.

## TDD red/green strategy

RED:

1. Land only the black-box CLI integration test described in step 1.
2. Run `go test ./cmd/do-work-cli -run TestGreenGateEvidenceLifecycle -count=1` from `skills/do-work/tools/do-work-cli`.
3. Confirm it compiles, invokes the real CLI, and fails its result assertion because `record-green-gate`/`check-green-gate` are unknown. Preserve that output as the REQ's red evidence.

GREEN:

1. Implement the minimum handler, record, result-model, and registration path needed for record plus exact match; rerun the focused test until that first case passes.
2. Add the moved-HEAD and log-only cases from the captured GREEN contract, then implement ancestry and per-commit path validation until all three pass.
3. Add the negative identity/provenance/state cases and atomic replacement/refusal tests, implementing only the validation each case requires.
4. Add result-renderer parity and action contract adjustments last, then run the focused package/action checks. Record the named test and its assertion failure before implementation plus its passing output after implementation in the REQ's Testing evidence.

The three captured acceptance cases remain primary: exact recorded `HEAD` matches, a project-changing commit does not match, and one `_dev/gate-runs/`-only commit above the record still matches. Additional tests protect the explicit repository/argv/ancestry constraints and the REQ-514 refusal invariant; they do not replace those three cases.

## Verification

Run from `skills/do-work/tools/do-work-cli` unless noted:

1. `gofmt -w` on only the new/changed Go files, then verify `gofmt -l` emits nothing for them.
2. `go test ./cmd/do-work-cli ./internal/gateevidence ./internal/resultmodel -count=1`.
3. `go test -race ./internal/gateevidence` to exercise record replacement/read boundaries.
4. `go vet ./...`.
5. `go test -count=1 ./...`.
6. From repository root, `bash _dev/tests/do-work-cli-go125-compatibility.sh`.
7. From repository root, `bash _dev/tests/contract-regressions.sh`; also confirm its line count did not increase and its existing mutation case fails when the check or post-green record directive is removed.
8. From repository root, run the unpiped `bash _dev/tests/maintainer-verify.sh` when the integrating workflow reaches the canonical gate.
9. Review `git diff --stat` and every changed file. Confirm no source/REQ files outside the declared write set changed, the Git-private fixture state is confined to temporary repositories, no debug artifacts remain, the action and reference agree on the same typed fields, and no alternate gate-success branch omits recording.

Acceptance is complete when the black-box test proves a pre-recorded current revision skips the baseline decision, Step 6.5 remains an unconditional direct gate, every direct green pipeline gate records its post-run `HEAD`, a project commit invalidates evidence, a log-only descendant preserves it, and all focused/module/repository checks exit zero.
