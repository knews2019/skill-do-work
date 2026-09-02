# REQ-498 Exploration — Remaining Resumable Finalization Slice

## Baseline

Inspected the tree at current `HEAD` after foundation commit `761d8e6a` (`[REQ-498] journal finalization transaction foundation`). The focused baseline currently passes:

```text
cd skills/do-work/tools/do-work-cli
go test -count=1 ./internal/finalization
ok github.com/knews2019/skill-do-work/do-work-cli/internal/finalization
```

The foundation already provides command registration, strict JSON decoding, Git-private per-REQ journals/payload copies, lifecycle/release postimage capture, exact-path commits, primary/metadata provenance, oldest-first journal replay, and one lifecycle-interruption test. `recover-finalization --discover` is still rejected by `handleRecoverFinalization`; actions still run direct `complete`/`fail`, `release`, staging/commit, and commit-hash recording.

## Authorities and patterns to reuse

- `internal/finalization/finalization_journal.go`: manifest decoding, journal identity/version/payload confinement, durable private writes, and oldest-first enumeration.
- `internal/requeststate`: `BuildPlan`, `PlannedPostimages`, `ApplyPlan`, and `RecordCommitProvenance`; do not reproduce lifecycle/archive/checkpoint/UR/calibration rules.
- `internal/publication`: `BuildReleasePlan` and `ApplyPlan`; do not reproduce release mutation rules.
- `internal/gittransaction/CommitExactPaths`: empty-index guard, exact staging, hooks, commit-path verification, and unstage-on-failure behavior.
- `internal/corehelpers.ReadProtectedInventory`: canonical M/A/D/X/XD classification. `AssociateProjectPaths` is usable only as parsing precedent: its generic latest-completed-owner tie-break is explicitly forbidden for discovery and skips all `do-work/` paths.
- `requestmodel.ParseDocument` and schema normalization are the existing identity/status/frontmatter authorities.
- `resultmodel.NormalizeResult` is where every newly added finalization slice must become `[]`, never `null`.
- Action contract checks belong near the lifecycle delegation assertions at the end of `_dev/tests/contract-regressions.sh`; assert ordered active directives, not mere token presence.

## Exact remaining production scope

The smallest coherent acceptance slice is still one vertical change: strict manifest/projection, phase rollback and terminal projection, bounded legacy discovery, recovery orchestration, and the two action entry points. Partial discovery without action startup ordering does not satisfy the captured RED/GREEN flow.

Modify:

1. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go`
2. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go`
3. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_journal.go`
4. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go`
5. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go`
6. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
7. `skills/do-work/actions/work.md`
8. `skills/do-work/actions/work-reference.md`
9. `skills/do-work/actions/commit.md`
10. `_dev/tests/contract-regressions.sh`
11. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`

Create:

12. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
13. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`

The supplied plan incorrectly omits `finalization_journal.go`. `validateManifest` and `readJournal` live there, so explicit `implementation_provenance_mode` validation cannot be strict for both incoming manifests and recovered journals without changing this file. No change to `cmd/do-work-cli/main.go`, `corehelpers/inventory.go`, request-state, publication, or Git-transaction packages is presently required; the foundation already exposed/composed those seams. Version/changelog files remain integrator-only.

## Required implementation details

- Add provenance mode exactly `primary_commit` or `supplied_commit`. Forbid a hash in primary mode; require 7–40 lowercase hex in supplied mode, resolve it as a commit, and require it to be an ancestor of `HEAD`. Preserve this intent in the journal and validate it again when reading a journal.
- Add authoritative ordered `finalizations: []` while retaining singular `finalization` for the one-record `finalize` compatibility surface. Construct both from one helper. Each record needs archived/journal paths, exact commit allowlist, hashes created this invocation versus settled primary/metadata hashes, blockers/reasons, next argv, and rooted JSON verification argv. Normalize every slice in `NormalizeResult`.
- Expose `verified` and `cleanup_complete` as result phases without inventing a second persisted state machine. A cleanup failure leaves the journal at metadata state; retry re-verifies and removes it.
- Fold `release_at` into the archived request's recorded release pre/post image. The current `stampReleaseAt` is an unjournaled side write after release convergence.
- On any pre-primary lifecycle/release/hook/index failure, restore only lifecycle/release images, empty only finalizer-owned staged entries, return the journal to `prepared`, and retain implementation/action-authored bytes. Refuse and retain the journal if exact image proof fails. After a primary commit, roll forward only.
- Strengthen recognition of an already-created primary commit. Current `matchingHeadCommit` accepts the same message plus any nonempty subset of allowed paths; that is weaker than exact transaction identity. Bind enough prepared-HEAD/diff evidence in the journal to reject an unrelated same-message/subset commit while still recovering the commit-created/journal-not-persisted interruption.
- Parse only optional `--discover`; replay journals first and stop discovery on replay failure. Freeze one protected inventory/repository/Git snapshot, reject staged X as `FINALIZATION-PROTECTED-STAGED` and any other pre-existing staged path as `FINALIZATION-FOREIGN-STAGED`, and never read X content.
- Candidate legacy owners are archived terminal requests with blank success provenance, plus the narrowly defined dirty newly archived failed case. Project ownership requires the exact Implementation Summary path, current dirt, and one candidate only. Do not call `AssociateProjectPaths` as the final decision because its latest-completion tie-break violates this REQ.
- Admit shared paths only through whole-diff REQ-specific semantic proof: lifecycle move/identity/transition, exact writer checkpoint removal, calibration-only row, UR membership closure/move, originating-REQ follow-up link, and coherent release entry/version mirrors. Never split hunks. Journal each accepted group at the already-applied lifecycle/release boundary, then use `advanceJournal`; discovery must not implement a second commit path.
- Process safe groups by `completed_at`, then REQ id, preserving unrelated unstaged project work. After them, refuse before selection if ambiguous shared/index state remains. Legacy worktree state without a durable asserted merge hash stays blocked.
- In `work.md`, recovery must be the first Step 1 operation, before the current first `CHECKPOINT.md` read, working-REQ crash recovery, `next`, or `claim`. Continue only on command success and empty blockers/reasons in every typed record. Keep ordinary working-REQ recovery otherwise unchanged and distinguish the two paths in `work-reference.md`.
- Replace Step 8's direct terminal call and Step 9's release/staging/primary/hash-record tail with one action-authored manifest and `finalize`. Semantic terminal/failure/release/follow-up/lesson judgment stays in the action. Serial/already-green use primary provenance; worktree uses the held merge hash in supplied mode. `commit.md` runs discovery before protected association, then groups only leftovers.

## RED/GREEN coverage

Extend the existing foundation test rather than replacing it, and put the larger matrix in the new recovery test file:

- Capture RED first with the REQ-494-shaped no-journal fixture: `--discover` is rejected/no-op, provenance stays blank, no commits are created, and dirty checkpoint state prevents the following REQ-452 claim.
- GREEN the same fixture: one discovered record, exactly one primary and metadata commit, primary provenance recorded, unrelated `notes.txt` untouched, journal cleaned, second invocation idempotent, then canonical `next` selects and canonical `claim` claims REQ-452.
- Table-drive interruptions after prepared, lifecycle, release, primary, metadata, verification, and pre-cleanup across serial success, failure, already-green/no-release, and supplied-worktree hash.
- Pin hook failure rollback, corrupt-image refusal with byte-identical repo/index, staged ordinary and X guards, competing owners, foreign checkpoint hunk, shared-path competition, unmatched release changes, two safe groups in stable order, unrelated unstaged preservation, and recognition of existing phase commits.
- Contract assertions must prove recovery precedes the first checkpoint read/working recovery/selector, `finalize` replaces direct final-tail lifecycle/release/staging/hash commands, actions consume typed records, commit recovery precedes association, and the CLI prime no longer says discovery remains caller-owned.

Required commands, unpiped:

1. `cd skills/do-work/tools/do-work-cli && go test -count=1 ./internal/finalization`
2. `cd skills/do-work/tools/do-work-cli && go test -race ./internal/finalization ./internal/gittransaction ./internal/requeststate ./internal/publication`
3. `cd skills/do-work/tools/do-work-cli && go vet ./...`
4. `cd skills/do-work/tools/do-work-cli && go test -count=1 ./...`
5. `bash _dev/tests/do-work-cli-go125-compatibility.sh`
6. `bash _dev/tests/contract-regressions.sh`
7. `bash _dev/tests/maintainer-verify.sh`

## Integration hazards

- REQ-489 is concurrently changing checkpoint insertion/removal in `requeststate/state_apply.go`. REQ-498 does not need to edit that file, but its lifecycle images and REQ-494 fixture exercise the same behavior. Integrate REQ-489 before REQ-498's merged-state verification, then rerun the finalization and request-state suites; otherwise a simple one-line fixture can hide orphaned continuation lines or false heading matches.
- Current `HEAD` also contains later REQ-453/UR-097 changes to targeted-run replay in `work.md` and `work-reference.md`; edit the live text, not the pre-`761d8e6a` shape, and preserve the frozen-ledger/fan-out semantics.
- The finalizer's exact commit allowlist includes action-authored follow-ups, reports, lessons, and project paths. Rollback must distinguish those preserved bytes from CLI-owned lifecycle/release bytes even though both are committed together.
- Recovery aggregation currently overwrites `aggregate.Finalization` on every journal. Adding discovery without fixing this first would silently discard earlier records and let actions continue on incomplete evidence.
- `finalizationSuccess` currently reports `metadata_committed` before deleting the journal, and cleanup failure reconstructs the same phase. Typed terminal phases must be computed carefully so retries remain truthful without persisting a new machine.
- `commandFailure` has no REQ-specific record. Discovery-wide staged/protected refusal therefore needs an explicit typed blocker projection (possibly an empty request id for a global refusal) so actions never parse prose.

## Scope judgment

No listed production file is safely removable from the coherent slice except that no new command-registration or shared-authority edits are needed. Conversely, `finalization_journal.go` is mandatory and missing from the supplied plan. Keep the discovery classifier narrow to the enumerated legacy evidence classes; a general dirty-tree association framework, generic hunk splitting, concurrent releaser coordination, or new public lifecycle interfaces would be scope expansion.
