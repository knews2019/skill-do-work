# REQ-507 Exploration — Finalization Behind `advance`

## Finding

The migration is structurally safe, but the plan understates two implementation seams. `lifecycleadvance` already imports `internal/finalization` from `recovery_commands.go`, and `finalization` does not import `lifecycleadvance`, so composing finalization from `advance` creates no package cycle. The existing finalizer also already returns every terminal field the action needs through top-level ordered `finalizations` plus the one-record compatibility projection.

The gaps are at the outer request boundary and the text renderer:

- `finalization.Handlers()["finalize"]` validates the manifest against the repository, but it cannot know which REQ the outer `advance REQ-NNN` selected. Checking the manifest in `lifecycleadvance` and then asking the handler to reopen it would introduce a replacement race between identity validation and mutation. The bound entry point must decode once, compare `request_id` and `request_path` with the selected advance record before journal creation or any other side effect, and continue from those same decoded bytes.
- `resultmodel.NormalizeResult` already normalizes ordered finalization records, but `renderText` never renders them. A plan test claiming text/JSON parity therefore requires a production edit to `internal/resultmodel/result_model.go`; a test-only change cannot establish that contract.

## Current ownership and exact seams

- `internal/lifecycleadvance/advance_commands.go`
  - `handleAdvance` discovers one non-queue request, calls `classifyAdvance`, and routes every extra argument to `executeAdvanceEvidenceGates`.
  - `classifyWorkingAdvance` currently ends at `agent judgment: prepare finalization manifest` after `## Orientation`.
  - Change that terminal classifier to mechanical `finalize`, return tokenized argv containing the exact discovered request path and one manifest placeholder, and dispatch that phase to a dedicated finalization composition function rather than the ordinary gate aggregator.
- `internal/lifecycleadvance/evidence_gates.go`
  - Its parser owns estimator/preflight/qualification/test inputs and its aggregator intentionally projects only gate records. Finalization has a richer typed result and rollback contract; do not squeeze it into `AdvanceGateRecord` or make this file a second finalization model.
- `internal/finalization/finalization_commands.go` and `finalization_prepare.go`
  - `handleFinalize` parses `--manifest`, then `prepareJournal` decodes/validates the manifest before creating or resuming a journal.
  - Expose a package entry point for bound in-process finalization and make the selected REQ id/path comparison immediately after that single decode, before index checks, journal-directory creation, lifecycle planning, or mutation. The existing public `finalize` handler should continue to call the same path without an outer binding.
- `internal/resultmodel/result_model.go`
  - `CommandResult` already has `Finalization` and `Finalizations`; `FinalizationResult` already carries request/archive/journal identity, phase/status, resumed/discovered state, commit paths and hashes, blockers/reasons, and recovery argv.
  - `NormalizeResult` already supplies non-null finalization collections. Only text rendering is missing if parity remains a planned acceptance condition.
- `internal/finalization/finalization_apply.go`
  - `finalizationSuccess`, `finalizationFailure`, and `finalizationRecord` already produce the authoritative terminal records. No change is indicated.
- `cmd/do-work-cli/main.go`
  - No registration change is required: `advance` is already registered through `lifecycleadvance.Handlers()`, while the direct finalization commands remain backward-compatible.

## Existing versus missing coverage

| Path | Existing direct-engine coverage | Public `advance` coverage | Needed here |
|---|---|---|---|
| Serial / primary commit | `TestFinalizeAcceptsWorkingRequestDirtWrittenByThePipeline`; planned-release phase tests also exercise release payloads | None | Add a real `advance REQ --finalization-manifest ...` transaction and assert archive, optional release, primary/metadata hashes, and clean effects. |
| Supplied worktree commit | `TestFinalizeSupportsSuppliedWorktreeProvenanceWithoutMetadataCommit` | None | Add public composition coverage asserting supplied implementation provenance, a tail primary commit, and no metadata commit. |
| Completed with issues | Manifest validation accepts the value; schema/dependency tests recognize it, but no finalization transaction test exercises it | None | Add the full public terminal transaction and assert `terminal_status` plus archived frontmatter remain `completed-with-issues`. |
| Already-green / no release | `TestRecoverFinalizationAlreadyGreenNoReleaseManifest` covers direct recovery without release postimages | None | Add public composition coverage proving no changelog/version/lockfile mutation or release stamp and lifecycle-only commit paths. |

All four paths are therefore uncovered at the new authority seam. The completed-with-issues path is additionally uncovered end to end in the finalization package itself, so the public test should not be reduced to a phase-only fixture.

Add refusal cases beside the matrix: missing manifest input, duplicate/unknown finalization option, outer request-path mismatch, manifest `request_id` mismatch, manifest `request_path` mismatch, malformed/irrelevant input, and a hostile token passed as data. Snapshot bytes/HEAD before mismatch cases and prove no request, checkpoint, archive, release, index, journal, or commit effect.

## Active prose and predicate state

- The captured `_dev/tests/contract-regressions.sh` path is stale. It is a 76-line dispatcher; `_dev/tests/contracts/core-checks.sh` owns the active action predicates.
- The measured “19 Step 8/9 predicates” no longer exist in the current contract owner. Current final contract checks cover the REQ-506 mechanical evidence loop only. Add structural finalization-ownership guards in `core-checks.sh`; do not edit the dispatcher or pretend to delete predicates that are no longer present.
- `work.md` still contains the direct `finalize --manifest` shell block and detailed manifest/provenance/verification mechanics in Step 9. Step 8 still contains a long mechanical preparation tail among its judgment substeps.
- `work-reference.md` still contains the large Changelog Entry Procedure and the direct Commit & Metadata-Commit recipe. The latter has its own `finalize --manifest` block.
- `prime-do-work-cli.md` still calls the tested tail “finalization-manifest judgment” and says finalization remains judgment until the action supplies a manifest.
- References from `actions/commit.md`, `actions/review-work.md`, and `docs/standing-preferences.md` to the named “Commit Phase” remain semantically valid if that heading is retained; they do not need edits merely because `advance` becomes the invoker. Board source comments about changelog precedence are release-domain readers, not alternate tail writers.

The new shell predicates should isolate the Step 8/9 and reference-procedure sections, require `advance` plus ordered typed finalization consumption, and forbid a direct finalizer invocation or reintroduced archive/release/stage/commit/provenance recipe there. Do not globally ban `git commit` in `work.md`: hand-back integration legitimately owns a merge commit earlier in the action.

## Judgment/mechanics boundary

Keep in prose:

- Fold-First consolidation or minting for user and stakeholder questions, cold-reader wording, audience routing, and cycle judgment.
- Sweep/discovered-task consolidation and impact stamping, including whether a finding is critical enough to queue.
- Terminal and failure classification (`completed`, `completed-with-issues`, intent/spec/code/environment, and external-precondition blocked flip).
- Release judgment: whether a release exists, affirmative project ownership, source/bump/mirror choice, changelog title and human prose, and the exact payload content supplied to the manifest. Already-green remains a judged no-release case.
- Deferred-lesson relevance/family promotion and the retained operative worktree cleanup decision after typed success.

Move behind `advance`/`finalization`:

- Request-bound manifest admission and validation; archive/checkpoint/UR/calibration mutation; release payload validation/publication; exact staging/commit allowlist; primary/metadata commit creation; provenance stamping; terminal verification; rollback; journal cleanup; and the typed success/refusal projection.
- The action should continue only when the command outcome is success, exactly one ordered record matches the selected REQ id/path, `phase == cleanup_complete`, and `blocked_paths`/`reason_codes` are empty. It may report archive and settled/created hashes from that record, but must not infer success from Git status, display text, archive presence, or singular `finalization`.

## Proposed exact write set

1. `skills/do-work/actions/work.md`
2. `skills/do-work/actions/work-reference.md`
3. `_dev/tests/contracts/core-checks.sh`
4. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
5. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go`
6. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate.go` (new)
7. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate_test.go` (new)
8. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go`
9. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go`
10. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go`
11. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
12. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`

Items 9–10 are required by the same-read request-binding rule; the current finalization handler has no safe outer-identity seam. Items 11–12 are required if the plan retains text/JSON parity. `internal/lifecycleadvance/evidence_gates.go`, `internal/finalization/finalization_apply.go`, `internal/finalization/finalization_types.go`, command registration, and the root contract dispatcher should remain untouched unless a RED test exposes a separate defect.

## Public behavior contract

With orientation present and no manifest input, `advance` reports mechanical phase `finalize`, typed `needs_input`, and exact tokenized continuation argv. With one bound manifest, the same public command runs finalization in process and returns:

- the pre-mutation `advance` identity for the selected request;
- the subordinate global outcome, findings, changes, and rollback unchanged;
- ordered top-level `finalizations`, with singular compatibility only for exactly one record;
- no shell interpolation and no second lifecycle/archive/commit implementation.

Irrelevant finalization input at another phase refuses. A valid manifest for another live REQ must refuse before any durable or Git effect.

## Verification commands

From `skills/do-work/tools/do-work-cli/`:

- `go test -count=1 ./internal/lifecycleadvance ./internal/finalization ./internal/resultmodel`
- `go test -race ./internal/lifecycleadvance ./internal/finalization ./internal/resultmodel`
- `go vet ./...`
- `go test -count=1 ./...`

From the repository root:

- `bash _dev/tests/contract-regressions.sh`
- `bash _dev/tests/maintainer-verify.sh` (direct and unpiped at the integration gate)

The public lifecycle tests already build the real CLI once through `advanceCLIBinary`; extend that seam rather than testing only a package-local helper. Keep each test-file invocation below the fast-tier 30-second budget.

## Risks and controls

- **Manifest replacement race:** avoid validate-then-reopen across packages; bind the selected request during the finalizer's single decode before journal creation.
- **Opaque projection:** copy subordinate outcome/findings/changes/rollback and ordered records from one result; do not rebuild a reduced aggregate that drops the actionable blocker.
- **False coverage:** registration/phase tests alone can stay green while mutations diverge. Assert status, ordered evidence, paths, commits, release/no-release effects, and byte-identical refusal cases through the public CLI.
- **Predicate overreach:** section-scope the negative action checks so legitimate earlier merge commits are not forbidden.
- **Stale prose sweep:** update the CLI prime in the same change; keep named Commit Phase callers unless their actual behavior statement becomes false.
- **Scope creep:** do not rewrite finalization transaction internals or the changelog policy beyond deleting displaced mechanics; the existing engine remains the authority.
