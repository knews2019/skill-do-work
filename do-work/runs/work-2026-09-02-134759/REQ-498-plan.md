# REQ-498 Route C Plan — Complete Resumable Finalization

## Planning boundary

This is the remaining safe vertical slice after commit `761d8e6a` (`[REQ-498] journal finalization transaction foundation`). That commit already registered `finalize` and `recover-finalization`, added the Git-private journal, composed canonical request-state/release/exact-commit/provenance authorities, and proved one lifecycle interruption. This slice does not replace those authorities or redesign ordinary request-state/publication commands. It completes startup recovery, narrowly handles the legacy unjournaled tail described by UR-096, and delegates the two action entry points.

Single-releaser behavior remains unchanged. Builder scope excludes `CHANGELOG.md`, `VERSION`, shipped version mirrors, and `skills/do-work/actions/version.md`; the queue owner remains their sole writer at integration.

## Exact write set

Modify:

1. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go`
2. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go`
3. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go`
4. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go`
5. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
6. `skills/do-work/actions/work.md`
7. `skills/do-work/actions/work-reference.md`
8. `skills/do-work/actions/commit.md`
9. `_dev/tests/contract-regressions.sh`
10. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`

Create:

11. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
12. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`

No do-work queue/archive/checkpoint state belongs to the builder write set.

## Architecture and invariants

### One engine, two entry paths

- `finalize --manifest` remains the only planned-tail entry. It validates the complete manifest, persists the journal and copied release payload before tracked mutation, then advances the existing phase engine.
- `recover-finalization` remains backward compatible as journal replay only. `recover-finalization --discover` first replays journals oldest-first, then takes one protected inventory/repository/Git snapshot and associates only safe legacy unjournaled groups.
- Every accepted legacy group is converted into a normal Git-private journal at the already-applied lifecycle/release boundary before its first commit. From that point it uses the same primary-commit, provenance, verification, and cleanup code as a planned finalization. Discovery does not get a second commit implementation.
- A journal failure stops discovery. During discovery, safe groups already committed remain committed, but any ambiguous shared remainder makes the overall command refuse and prevents selection/claim.

### Strict manifest/provenance contract

Keep strict unknown-field rejection and make provenance intent explicit:

- `implementation_provenance_mode`: exactly `primary_commit` or `supplied_commit`.
- `implementation_hash`: forbidden for `primary_commit`; required for `supplied_commit`, 7–40 lowercase hexadecimal characters, resolvable as a commit, and an ancestor of current `HEAD`.
- Serial success and already-green success use `primary_commit`; worktree success uses `supplied_commit` with the held merge hash. Failure uses `primary_commit` for its exact failure/lifecycle commit but never writes a success `commit:` field.
- Retain the existing exact fields: `request_id`, `request_path`, `writer_label`, `transition`, `terminal_status` or `failure_error`/`failure_type`, `completed_at`, request/checkpoint SHA-256 preimages, `commit_paths`, `commit_message`, and the paired optional `release_manifest_path`/`release_at`.

The action creates this one manifest only after its semantic judgments, follow-up/report/lesson content, implementation summary, and optional release payload are ready, but before `complete`, `fail`, `release`, staging, or committing has mutated the tail. The finalizer must validate that `commit_paths` contains every canonical lifecycle and optional release target; exact action-authored/project paths may be added, never inferred by the CLI.

### Phase behavior and rollback

- Preserve durable journal phases `prepared`, `lifecycle_applied`, `release_applied`, `primary_committed`, and `metadata_committed`; expose terminal verification/cleanup progress in the typed result without adding a second persisted state machine. A cleanup failure leaves the metadata journal, so retry safely re-verifies and removes it.
- Fold the `release_at` edit into recorded release pre/post evidence rather than leaving it as an unbound side write. Its preimage is the archived lifecycle postimage; its postimage is those exact bytes with the release stamp. This makes interrupted release stamping and pre-primary rollback image-driven.
- Before a primary commit exists, a nested lifecycle/release/commit-hook failure attempts to restore only the recorded CLI-owned release then lifecycle preimages, empties only the finalizer's index entries, returns the journal to a replayable prepared boundary, and preserves implementation/action-authored edits. If exact rollback cannot prove an image, report committed-state risk/refusal and retain the journal.
- Once the primary commit exists, never roll back; always roll forward through success provenance, verification, and journal cleanup. Matching an already-created primary/metadata commit must not create a duplicate.

### Bounded legacy discovery

Discovery is intentionally limited to the UR-096 evidence classes; it is not a general dirty-tree auto-committer.

1. Run canonical protected inventory first. Never read/diff an `X` path. Any staged `X` produces `FINALIZATION-PROTECTED-STAGED`; any other pre-existing staged path produces `FINALIZATION-FOREIGN-STAGED`. Both stop before mutation. Unstaged unrelated `X` remains untouched.
2. Candidate owners are archived terminal REQs whose success provenance is blank (or a dirty newly archived failed REQ for which no success provenance is expected). A project path is owned only when it is in that REQ's exact `## Implementation Summary`, the path is currently dirty, and no other candidate summary names it. Do not call or reproduce the generic latest-completed-owner tie-break.
3. Shared paths require whole-diff semantic proof for the same candidate:
   - working deletion and archive addition parse to the same REQ identity and canonical terminal transition;
   - the checkpoint diff removes only that REQ's exact writer-labelled entry;
   - calibration additions contain only rows for that REQ and its recorded completion/estimate evidence;
   - active-to-archive UR movement has the same `user_request` identity and member closure;
   - new/edited follow-up REQs explicitly reference the originating REQ (`addendum_to`, `related`, or the canonical source record);
   - release metadata contains one matching REQ/title entry and one coherent old→new version mirrored across the changed release targets.
4. A path joins a group only when its entire current diff belongs to that one group. Competing owners, foreign hunks in a shared file, unmatched queue/lifecycle/release paths, or a dirty residual checkpoint get stable reason codes and exact blocked paths. Do not split hunks or guess.
5. Freeze all associations from one snapshot, order groups by `completed_at` then REQ id, journal each group before commit, and commit every unambiguous group. Preserve unrelated unstaged project changes. After safe groups, rescan only to decide whether shared/index state is safe; refuse before selection if ambiguous shared state remains.
6. A discovered serial group uses its canonical recovery primary commit as provenance and then the ordinary isolated metadata commit. Legacy worktree state without an asserted, durable merge hash is ambiguous and remains blocked; discovery never chooses a nearby/latest merge.

### Typed output consumed by actions

Add authoritative ordered `finalizations: []` records to `CommandResult`; keep the existing singular `finalization` projection for the one-record `finalize` compatibility surface, deriving both from the same record helper. Each record has exactly:

- `request_id`
- `archived_path`
- `journal_path`
- `phase` (`prepared`, `lifecycle_applied`, `release_applied`, `primary_committed`, `metadata_committed`, `verified`, or `cleanup_complete`)
- `resumed`
- `discovered`
- `commit_paths` (the exact group allowlist)
- `created_commits` (only hashes created by this invocation)
- `primary_commit` and `metadata_commit` (settled hashes whether newly created or recovered; metadata may be empty for failed/worktree cases)
- `blocked_paths`
- `reason_codes`
- `next_argv`
- `verification_argv`

Normalize every slice to `[]`, not `null`. `recover-finalization --discover` returns records in execution order. The work/commit actions authorize continuation only from command `outcome: success` with every record's `reason_codes` and `blocked_paths` empty; they do not parse prose findings or rediscover paths. Verification argv is the exact rooted JSON replay command.

## Ordered implementation tasks

1. **Capture RED first.** Add the REQ-494-shaped no-journal fixture and action-order assertions before production edits. Confirm `--discover` is rejected/no-op, the archived request stays without provenance, the dirty checkpoint blocks the following claim, and no automatic primary/metadata commit exists.
2. **Finish schema/result projection.** Add explicit provenance mode validation and ordered per-record finalization output. Centralize record construction so success, refusal, replay, and discovery cannot disagree on phase, hashes, blockers, or argv.
3. **Close phase gaps.** Bind `release_at` to exact release images; add pre-primary reverse convergence for lifecycle/release and hook/index failures; retain roll-forward-only behavior after primary commit; surface verification and cleanup phases.
4. **Implement bounded discovery.** Build the one-snapshot candidate/group classifier in `finalization_discovery.go`, apply protected/staged guards before content reads, enforce whole-path ownership and semantic shared evidence, persist discovered journals, then pass safe groups to `advanceJournal`.
5. **Orchestrate recovery.** Parse only optional `--discover`; replay journal paths oldest-first; run discovery only after replay is clean; aggregate all records instead of overwriting earlier results; stop with exact residual blockers.
6. **Delegate `do-work run`.** At the very start of Step 1, before reading `CHECKPOINT.md`, active-working crash recovery, `next`, or `claim`, invoke rooted JSON `recover-finalization --discover`. Report recovered records, continue automatically on typed success, and stop before selection on any blocker. Then retain the existing working-REQ crash recovery unchanged and explicitly distinguish it from finalization-tail recovery in `work-reference.md`.
7. **Delegate completion.** Replace Step 8's direct `complete`/`fail` plus Step 9's direct `release`, manual staging/commit, and separate hash-recording tail with creation of the exact manifest and one `finalize` call. Preserve action-owned terminal/failure/release judgment, follow-up authoring, lessons, worktree merge hash, already-green release omission, and the single-releaser rule. Consume exact returned record paths/hashes for reporting.
8. **Delegate `do-work commit`.** Before protected-inventory association, invoke `recover-finalization --discover`. On typed success, continue grouping only the remaining changes; on refusal stop. Remove the statement that this action never retries finalization tails, while retaining ordinary generic association for non-finalization leftovers.
9. **Lock the action boundaries.** Extend contract regressions to prove startup ordering, canonical command spelling, typed-record consumption, absence of direct final-tail lifecycle/release/staging/hash commands, commit-action recovery-before-association, and the updated prime ownership statement.
10. **Unify.** Review every changed file, run focused and full verification unpiped, and ensure no journal/payload/build artifact is left in Git status.

## Requirement mapping

- Strict `finalize --manifest` and `recover-finalization --discover`: tasks 2 and 5.
- Exact Git-private journal before lifecycle mutation and all phases: tasks 2–3 and existing foundation.
- Canonical lifecycle/release/protected-inventory/Git authorities: tasks 3–4; discovery classifies only, then calls the existing engine.
- Idempotence/no duplicate archive, calibration, release, version, or commit: tasks 3–5 plus interruption matrix below.
- Serial primary provenance and supplied worktree merge provenance: tasks 2 and 7.
- Recovery before working recovery/selection/claim and automatic continuation: task 6.
- Safe semantic legacy association without latest-owner guessing: task 4.
- Preserve unrelated unstaged work; block ambiguous shared/index/protected state: tasks 3–5.
- Typed phase/recovery/commit/blocker/reason/verification output: task 2 and exact field list above.
- Work and commit action delegation/crash-recovery documentation: tasks 6–9.
- Backward-compatible `complete`, `fail`, `release`, and single releaser: no changes to their command interfaces; tasks 6–9 only replace callers' tail instructions.

## Tests and captured red/green proof

### Focused Go behavior

In `finalization_recovery_test.go`:

- Table-drive interruption after prepared, lifecycle, release, primary commit, metadata commit, verification, and immediately before cleanup. Rerun recovery and assert one archive move, calibration row, changelog entry, version bump, primary commit, provenance write, and metadata commit.
- Cover successful serial, failed, already-green/no-release, and supplied-worktree-hash flows. Assert the archived success `commit:` is primary for serial and the supplied merge hash for worktree.
- Make a commit hook fail without bypassing it; assert CLI-owned pre-primary bytes/index roll back, implementation bytes remain, journal remains resumable, and a later normal retry succeeds.
- Corrupt one journal image/preimage and assert refusal with byte-identical repository/index state.
- Pre-stage an ordinary path and an `X` path separately; assert stable foreign/protected reason codes and no reads/commits.
- Create competing Implementation Summary owners, a checkpoint with a foreign hunk, a shared path split across two candidates, and an unmatched release/version change; assert exact blocked paths and no guessed owner.
- Create two independent safe legacy groups plus unrelated unstaged project work; assert both groups commit in stable order, each records canonical provenance, and unrelated work remains unstaged/uncommitted.
- Verify already-existing phase commits are recognized and a metadata-state journal is only verified/cleaned.

### Required REQ-494 RED/GREEN fixture

Fixture `HEAD` contains claimed `REQ-494` in `working/`, its exact writer checkpoint entry, pre-change project bytes, and pending `REQ-452`. The worktree (with no journal) contains the canonical terminal archived `REQ-494` with blank `commit:`, removal of its checkpoint entry, the uniquely listed project edit, and any exact already-green lifecycle/calibration bytes; unrelated unstaged `notes.txt` is also dirty.

- **RED capture before implementation:** `recover-finalization --discover` has no supported discovery path; `REQ-494` remains unprovenanced, no recovery commits appear, and subsequent `next`/`claim REQ-452` cannot safely mutate the common dirty checkpoint.
- **GREEN:** the same command emits one `discovered: true` record, creates exactly one primary and one metadata commit, records primary provenance, leaves `notes.txt` untouched, and cleans its journal. A second invocation creates no commits/effects. In the same test flow, canonical `next` selects `REQ-452` and canonical `claim` succeeds without a manual `do-work commit` step.

The action contract fixture additionally proves `work.md` executes that recovery command before its first checkpoint read/working-REQ recovery/selector call, which connects the CLI behavior to the actual run contract.

### Commands

Run, in order and without pipes:

1. `cd skills/do-work/tools/do-work-cli && go test -count=1 ./internal/finalization`
2. `cd skills/do-work/tools/do-work-cli && go test -race ./internal/finalization ./internal/gittransaction ./internal/requeststate ./internal/publication`
3. `cd skills/do-work/tools/do-work-cli && go vet ./...`
4. `cd skills/do-work/tools/do-work-cli && go test -count=1 ./...`
5. `bash _dev/tests/do-work-cli-go125-compatibility.sh`
6. `bash _dev/tests/contract-regressions.sh`
7. `bash _dev/tests/maintainer-verify.sh`

The final hand-back must include the initial failing RED assertion/output, the focused GREEN result, the full maintainer gate exit code, and `git diff --stat` plus a file-by-file diff review.
