# REQ-574 plan: bring four do-work-cli test files under the 30 s per-file budget

Status: read-only plan. No files edited, no Go tests started by the planner. Timing evidence comes from the lead's background run (`go test -count=1 -v`, three packages concurrently, other gates on the machine) in `scratchpad/slow-tests-574.log` and from `do-work/test-durations.tsv` (one-gate history: defer_gate 22.9 s, inventory 18.7 s, recovery 22.3 s, req499 20.0 s in run 20260904T231224Z-80089).

## Plan

### Where the time goes

The budget wrapper (`_dev/tests/run-go-tests-with-budget.sh`) sums each top-level test's `Elapsed` (active wall time) and attributes it to the file whose `func TestX(` defines it. Three facts follow:

- `t.Parallel()` does not lower the sum. A parallel test's `Elapsed` excludes its paused time but still counts its own active wall; eight tests running side by side each still report their own ~1 s. It also cannot be used where `t.Setenv` is called (inventory matrix) or where package-level hooks are swapped (`afterPublicationMutation`, `beforeReservationMarkerOpen`, `afterFinalizationPhase`, `enumerateTrackedReleasePaths`). Not a lever here.
- Only fewer subprocess spawns per test, or fewer tests per file, move the number. Every heavy test here is git-spawn bound; nothing sleeps or retries.
- The gate runs `go test -short ./...` with package-level parallelism (8 packages at once on this 8-core machine), so per-test elapsed already includes cross-package contention. That is why 17 to 28 s at one gate becomes 31 to 60 s with two or three gates.

F1. `internal/corehelpers/inventory_test.go` (18.7 s one gate; 64.7 s loaded). 57 invocations of `runRetainedInventory` run `skills/do-work/tools/checks/uncommitted-inventory.sh`, a 5-line shim that execs `do-work-cli.sh`, which runs `go version` plus `go tool -C <module> -n do-work-cli` on every call before exec'ing the cached binary. The `go tool -n` staleness check (hash every module input) costs roughly 0.3 s alone and far more under load. 45 calls come from `TestInventoryMatchesRetainedPorcelainXYMatrix` (53.5 s loaded), 10 from `TestInventoryMatchesRetainedSecretOriginAndAmbiguityMatrix`, 1 each from the differential-comparator test and the real-repo AD test. The Go side of each pair is in-process and costs nothing. The "retained" script no longer holds a shell implementation; it is a launcher for the same Go code, so the differential now pins (i) the text-format projection of `uncommitted-inventory` against `inventoryRow`, and (ii) that the one-argument shim route works.

F2. `internal/publication/defer_gate_test.go` (22.9 s one gate; 66.9 s loaded). `newDeferGateRepository` is called 57 times (18 of them inside `TestDeferGateRollsBackUntrackedCreateAndFoldTopologies`, 26.6 s loaded) and spawns 7 git processes each (`init`, `config` x2, `add`, `commit`, `add`, `commit`): about 400 spawns, roughly 12 to 14 s of the 22.9 s. The rest is `BuildDeferGatePlan` plus `ApplyPlan` (transaction preflight spawns `ls-files` and `status` per target path, then per-path `diff --quiet`), about 0.15 to 0.2 s per case.

F3. `internal/finalization/finalization_recovery_test.go` (22.3 s one gate; 56.7 s loaded) and F4 `finalization_req499_test.go` (20.0 s; 32.1 s loaded). The fixture is cheap (`newFinalizationRepository` = 3 spawns, then 2 per seed). The cost is the production pipeline itself: legacy discovery runs `git diff --unified=0 HEAD -- <path>` and `git show HEAD:<path>` per dirty path (`finalization_discovery.go:291,420,1665`), and every phase transaction runs `ls-files` plus `status` per target path (`gittransaction/git_transaction.go:867-892`, `PreflightTargets`). With the 17-path `seedSemanticLegacyTail` fixture one recovery is 100+ spawns (2 to 3 s at one gate); `TestRecoverFinalizationResumesEveryDurablePhaseExactlyOnce` runs 7 finalize-then-recover-twice cycles (21.7 s loaded, about 9 s at one gate). Both files are grab-bags: `_req499_` is named after a REQ number and mixes sole-releaser attribution, release-mirror ownership, follow-up folds, release-replacement replay and pure version parsing; `_recovery_` mixes legacy-tail discovery, journal-phase replay and pure helper tests.

### Files to modify

Test-only. No production `.go` file changes.

1. `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go` (edit: resolve the launcher's binary once, exec it directly for the synthetic matrices).
2. `skills/do-work/tools/do-work-cli/internal/repositoryfixture/repository_fixture.go` (new, leaf package, standard library only: build a keyed template repository once per process, copy it per test in pure Go). It is imported only by `_test.go` files; keep it free of any command package import.
3. `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` (edit: `newDeferGateRepository` copies a template; two templates keyed by `secondParent`).
4. `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_rollback_test.go` (new, only if step 3 leaves the file above 12 s: move the three rollback-matrix tests and their helpers `assertDeferGateRollback`, `snapshotDeferGatePreimages`, `assertDeferGatePreimages`, `infoMode`, `deferGatePreimage`).
5. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go` (edit: `newFinalizationRepository` copies a template).
6. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` split into `finalization_legacy_discovery_test.go` and `finalization_journal_replay_test.go` (delete the old file; `git mv` one half so history follows).
7. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go` split into `finalization_sole_releaser_test.go`, `finalization_release_ownership_test.go`, and `finalization_release_replay_test.go` (delete the old file).
8. `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (one entry, same commit as the fix: the maintainer checks) and `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` § Verify (one line pointing at the budget command with the 15 s override).

### Ordered tasks

T0. RED. From the repository root, with nothing else running that you started:

```
DO_WORK_TEST_FILE_BUDGET_SECONDS=15 bash _dev/tests/run-go-tests-with-budget.sh \
  skills/do-work/tools/do-work-cli -short ./internal/publication ./internal/finalization ./internal/corehelpers
```

Expected: four `FAIL: <file> accumulated N s` lines (the files sit at 17 to 28 s at one gate, so the default 30 s budget is green alone; 15 s is the acceptance bar the lead set). Record the four sums in the REQ file. Also run `sysctl -n hw.ncpu` and note other gate processes (`pgrep -fl maintainer-verify | wc -l`) beside every measurement, because the numbers are load-sensitive.

T1. inventory_test.go (F1). Mechanism: keep the differential, remove the per-call toolchain check.

- Add a package-level `sync.Once` helper (name it `retainedInventoryBinary(t)`) that runs the exact resolution the launcher performs, `go tool -C <module-dir> -n do-work-cli`, once, and returns the printed executable path. This is the same cached executable the shim would exec, so the comparison target is unchanged. Resolve it before `t.Setenv("PATH", fakeRoot...)` installs the fake git: `go tool` may consult git for VCS stamping and the fake git exits 97 on unknown arguments. Simplest ordering: call the helper at the top of `runRetainedInventory`, before the `syntheticStatus != nil` block.
- In `runRetainedInventory`, when `syntheticStatus != nil` (the 55 matrix and comparator calls), exec the resolved binary with the argv the shim builds: `--repo-root <repository> --format text uncommitted-inventory`. When `syntheticStatus == nil` (the one real-repository AD test), keep executing `checks/uncommitted-inventory.sh` end to end so the shim's one-argument route stays pinned once.
- Do not touch the fake-git script, the status matrix, the expected classes, or `compareInventoryProjection`.
- Expected: 57 x ~0.3 s of `go version` + `go tool -n` becomes 1 x ~0.5 s plus 56 x ~0.03 s binary execs. File sum from 18.7 s to about 4 s at one gate.

T2. repositoryfixture package. Mechanism: pay git once per template, copy files per test.

- `func Template(t testing.TB, key string, build func(root string)) string` keeps a process-wide `map[key]*templateEntry` guarded by a mutex with a `sync.Once` per key. First call for a key runs `build` inside `os.MkdirTemp("", "do-work-fixture-"+key+"-*")`. Every call then copies the template tree into a fresh `t.TempDir()` with `filepath.WalkDir` + `os.MkdirAll` + `os.ReadFile`/`os.WriteFile`, preserving each file's permission bits (`info.Mode().Perm()`). Refuse symlinks and non-regular files in the template with `t.Fatalf`, since no template needs them and a silent skip would hide a broken fixture.
- Register `os.RemoveAll` of the template roots from a `TestMain` in each consuming package (one `TestMain` per package, 8 lines) or expose `repositoryfixture.Cleanup()` for them to call. A leaked 40-file temp directory per run is the alternative; do not leave it.
- Why not `git clone --shared`: it still spawns git, adds a remote and shares objects with the template, which the rollback tests would then have to reason about. A byte copy of `.git` yields a fully independent repository. The copied index carries stale stat data, so git re-hashes the handful of tracked files on first `status`; every test's before/after comparisons stay within one repository, so this is invisible to them.
- Expected: about 2 to 5 ms per copy against 200 to 250 ms per 7-spawn fixture.

T3. defer_gate_test.go (F2). Mechanism: two templates.

- `newDeferGateRepository(t, secondParent)` becomes `repositoryfixture.Template(t, "defer-gate-second-parent"/"defer-gate-single-parent", build)` followed by the existing `claimDeferParent(...)` file rewrite, which stays per test because it is plain file I/O. `build` is the current body minus `t.TempDir()`: init, config x2, write parent(s) and checkpoint, add, commit, write merge evidence, add, commit.
- The build closure cannot use `runGitFixture(t, ...)` with the outer `t` if the template outlives the first test; either run git via a small local function that returns an error and let `Template` fail the first caller, or have `Template` pass its own `testing.TB`. Choose the first: `build func(root string) error`.
- Expected: 57 x 7 = ~400 spawns removed, about 12 to 14 s. File sum from 22.9 s to about 8 to 10 s at one gate.
- Measure (T5 command). If the file is still above 12 s, do step T3b: move `TestDeferGateTrackedDirtyRepairFoldRollsBackExactPreimagesAfterEveryMutation`, `TestDeferGateRollsBackEveryMutationPositionToExactDirtyPreimages`, `TestDeferGateRollsBackUntrackedCreateAndFoldTopologies` and their rollback-only helpers into `defer_gate_rollback_test.go`. Family boundary: "an injected failure at every mutation index restores exact preimages and leaves no residue" versus "planning and successful publication of create and fold". Add a one-line file comment naming that boundary. Expected split: about 5 s and 4 s.
- Not chosen: reusing one repository across the `failureIndex` loop. It would cut 19 fixtures, but it changes what the tests prove (each index would then depend on the previous rollback being complete in ways the assertions do not check). Keep it as the maintainer's call, noted in the REQ file, not done here.

T4. finalization (F3, F4). Mechanism: template for the bare repository, then family split. The floor here is production spawn count; state that plainly in the REQ file.

- `newFinalizationRepository` copies template `"finalization-bare"` (init + user.name + user.email). Saves 3 spawns x ~70 repositories across the package, about 2 to 3 s per file, and every finalization test file benefits.
- Optional second template `"finalization-semantic-legacy-tail"` for the 13 `seedSemanticLegacyTail` callers (the 17-file seed commit): saves 2 more spawns each. Only if T5 shows a discovery file above 12 s; the post-commit dirty writes stay per test.
- Split `finalization_recovery_test.go`:
  - `finalization_legacy_discovery_test.go`: DiscoversLegacyNoJournalTail, AcceptsCoherentClaimOnlyTopology, DiscoveryRefusesExistingIndex, DiscoversCompleteSemanticLegacyTail, IgnoresFinderMetadataUnderDoWork, RefusesForeignSharedHunkByteIdentically, RefusesMultiHunkProjectPathByteIdentically, PreservesUnstagedProtectedAndDistinguishesStagedRefusals, DiscoveryRefusalRendersCompleteOrderedTypedEvidence, DiscoveryRefusalNamesInventoryAsTheResolvingVerb, ProcessesTwoSafeGroupsInStableOrder, StopsWhenTheRefusalOwnsNoRequest, plus `seedSemanticLegacyTail`, `seedSimpleDiscoveredTail`, `seedTwoSimpleDiscoveredTails`, `assertDiscoveryReason`, `semanticLegacyFixture`. Family: "an unjournaled tail is admitted or refused from repository evidence alone". Expected about 9 to 10 s at one gate.
  - `finalization_journal_replay_test.go`: ResumesEveryDurablePhaseExactlyOnce, FinalizeSupportsSuppliedWorktreeProvenanceWithoutMetadataCommit, FinalizeFailureReportsTerminalCleanupWithoutProvenanceMetadata, RefusesCorruptJournalImage, ReadJournalRejectsCorruptStoredImage, PreservesAndMatchesPrivateFileModeInLifecyclePostimages, SetsAsideRefusedRecordAndFinishesTheRest, StopsOnSharedDirtInsteadOfSettingOneRequestAside, StopsWhenRollbackLeavesResidue, ValidateManifestRequiresExplicitProvenanceMode, MatchingHeadCommit x2, VerifyFinalStateAllowsProvenanceAfterReleaseStamp, plus `seedPlannedFinalization`, `seedTwoPlannedFinalizations`. Family: "a journal written by finalize replays to the same terminal state from any phase". Expected about 11 to 12 s at one gate; ResumesEveryDurablePhaseExactlyOnce alone is about 9 s of it.
- Split `finalization_req499_test.go`:
  - `finalization_sole_releaser_test.go`: the four `AssumeSoleReleaser*` tests, `containsFinalizationPath`, `seedSoleReleaserLegacyTail`, plus the follow-up fold trio (RefusesForeignEditInTrackedFollowup, AcceptsExactTrackedFollowupFold, FollowupPathProvesTheEntireSingleNamedFold). Family: "which shared bytes one releaser may attribute to itself". About 3 to 4 s.
  - `finalization_release_ownership_test.go`: RefusesPartialConfiguredReleaseMirrors, RequiresWorkspaceMemberLockMirrors, ReleaseEnumerationFailureIsTypedAndFailClosed, RefusesReleaseMetadataWithoutProjectOwnership, RefusesCompleteInstalledSuiteWithoutMaintainerTopology, AcceptsSuiteMirrorsWithMaintainerTopology, AffirmativeReleaseOwnershipRequiresRootedWorkspaceChain, ReleaseVersionRecognizesOwnedCargoAndUVLockEntries, ReleaseVersionReadsOnlyProjectTOMLSections. Family: "which release mirrors the repository proves it owns". About 8 to 9 s; RequiresWorkspaceMemberLockMirrors is about 4.5 s of it.
  - `finalization_release_replay_test.go`: ReleaseReplacementPreservesModeZeroSentinel, ResumesAfterRealPreCommitHookFailure, AlreadyGreenNoReleaseManifest, PublicationPostimagesKeepExplicitReplacementMode, PublicRecoverFinalizationMovesURThenAllowsRealClaim (heavy-only already), `testCLIBinary`, `runPublicFinalizationCommand`, `seedPlannedReleaseFinalization`. Family: "an interrupted release publication resumes with its postimages and modes". About 4 to 5 s.
- Naming: every new file name says what it pins and is findable by plain-text search; no file is named after a REQ number any more. Put a one-line comment at the top of each stating the family. Keep `finalization_commands_test.go` as the shared-helper home.
- Do not mark any of these heavy-lane. None spawns the built CLI except the already-gated public test; they are unit tests of the discovery and phase engines.

T5. GREEN and measurement after each of T1, T3, T4:

```
DO_WORK_TEST_FILE_BUDGET_SECONDS=15 bash _dev/tests/run-go-tests-with-budget.sh \
  skills/do-work/tools/do-work-cli -short ./internal/publication ./internal/finalization ./internal/corehelpers
```

Then the ordinary gate, `bash _dev/tests/maintainer-verify.sh`, which appends one-gate numbers for every file to `do-work/test-durations.tsv`; paste the four (now nine) new rows into the REQ file with the concurrent-gate column.

T6. Same commit: lessons entry (family name suggestion: `per-call-toolchain-resolution-in-a-hot-loop`; the retained shim's `go tool -n` is correct for one user invocation and wrong inside a 57-iteration differential) and the prime § Verify line. Commit as a maintainer commit per `_dev/primes/prime-releases.md` (test-only under `skills/`, so it is a shipped-file change; check that prime for whether it is a release).

### Expected per-file times (one gate, this machine)

| File | Before | After |
|---|---|---|
| corehelpers/inventory_test.go | 18.7 s | ~4 s |
| publication/defer_gate_test.go | 22.9 s | ~8 to 10 s (T3); ~5 s + ~4 s if T3b splits rollback |
| finalization/finalization_legacy_discovery_test.go (from recovery) | 22.3 s combined | ~9 to 10 s |
| finalization/finalization_journal_replay_test.go (from recovery) | | ~11 to 12 s |
| finalization/finalization_sole_releaser_test.go (from req499) | 20.0 s combined | ~3 to 4 s |
| finalization/finalization_release_ownership_test.go (from req499) | | ~8 to 9 s |
| finalization/finalization_release_replay_test.go (from req499) | | ~4 to 5 s |

With two other gates (history ratio 1.6 to 2.2x) the worst file, journal replay, lands at 20 to 26 s. That is under 30 s but not comfortable, and it is the honest floor for a test-only change.

### Follow-up REQ to capture (not in this REQ)

The finalization floor is production fan-out: 2 spawns per target path in `gittransaction.PreflightTargets`/`inspectTargets` (`ls-files --error-unmatch` and `status --porcelain -z` per path) and per-path `git diff`/`git show` in `finalization_discovery.go`. One `git status --porcelain=v1 -z --untracked-files=all -- <all paths>` plus one `git ls-files -z -- <all paths>` per transaction would cut every transaction-heavy test by a third to a half with no output change, but it touches the `final-boundary-identity` family and needs its own differential (same outcomes, findings, paths on the existing fixtures). Capture it as a queued REQ with `related: [REQ-574]`.

### Testing approach

- RED is the wrapper at 15 s (T0). Do not lower the wrapper's default or pass `DO_WORK_TEST_ENFORCE_BUDGET=no`.
- Every moved test keeps its name and body byte-identical apart from the fixture call; `git diff --stat` on the move commits should show only the fixture lines changing. Diff the list of `func Test` names before and after (`grep -h '^func Test' internal/finalization/*_test.go | sort`) and require the same set; the wrapper's `tests=N` count must not drop.
- `go vet ./...` and `go test -count=1 -short ./...` from the module, then `DO_WORK_HEAVY_TESTS=1 go test -count=1 ./internal/finalization ./internal/corehelpers` once so the heavy-only tests still compile and run against the new helper locations.
- For T1, prove the differential still bites: temporarily change one expected class in `expectedOrdinaryInventoryClass` and confirm the matrix fails, then revert. Prove the shim route still runs once: `TestInventoryStagedAdditionDeletedFromWorktreeIsDeletion` must still fail if `checks/uncommitted-inventory.sh` is made non-executable for one run.
- For T2/T3, prove template copies are independent: `TestDeferGateFoldRefusesSymlinkReservationMarker` writes a symlink into its copy and `TestDeferGateRejectsDeferredMergeFromSideBranch` creates a branch; both must pass in the same run as every other test, in any `-shuffle=on` order. Run `go test -count=3 -shuffle=on ./internal/publication ./internal/finalization`.

### Risks

R1. Template leakage between tests. Mitigated by a byte copy per test (no shared object store, no shared index) and the `-shuffle=on -count=3` run. The template directory itself must never be handed to a test; `Template` returns only copies.

R2. Stale index stat data in copies. First `git status` in a copy re-hashes tracked files and may rewrite the index. All before/after comparisons in these tests are inside one copy, so no assertion can observe it. If a test ever compared index bytes across repositories it would need `git update-index --refresh` after copy; none does.

R3. `go tool -n` with the fake git on PATH. Resolve the binary before `t.Setenv("PATH", ...)`. Also resolve with the module directory as `-C` and never rely on the test's working directory being the package directory for anything but locating the module (`../..` from `internal/corehelpers`).

R4. Attribution. The wrapper maps a test to the file defining its `func Test`; helper-only files count nothing, and a test that moves takes its seconds with it. A new file that accidentally keeps a duplicate `func Test` name fails to compile, so a silent double count cannot happen.

R5. `t.Parallel()` temptation. Do not add it; it does not reduce the summed metric and conflicts with `t.Setenv` and the package-level hooks.

R6. Splits hide, they do not save. The finalization split lowers per-file sums without lowering gate wall time; the REQ file must say so and point at the follow-up REQ, so nobody reads the green gate as a speedup.

R7. `sync.Once` failure retry. A failed template build must fail every later caller with the same error, not retry silently; store the error beside the path.

R8. Heavy-lane classification is untouched. `DO_WORK_GO_TEST_EXCLUDE_PREFIXES` and `_dev/tests/heavy-lanes.json` are not edited; no test is reclassified.

### Consumer-field check

No production Go file, command, flag, typed result, finding code, JSON field, text rendering, hook, action file, or shell script under `skills/` changes, apart from the additive lessons and prime lines. The new `internal/repositoryfixture` package is imported only from `_test.go` files and builds no command. `go build ./...` output is byte-identical before and after; verify with `go tool -C skills/do-work/tools/do-work-cli -n do-work-cli` printing the same cache path before and after the change (the content hash covers only non-test inputs).

### Per-test timing table (heaviest 20, lead's loaded run)

Seconds are from the concurrent three-package `-v` run with other gates active; divide by roughly 2.5 for one-gate figures. Files in scope are marked.

| # | s | Test | File |
|---|---|---|---|
| 1 | 53.48 | TestInventoryMatchesRetainedPorcelainXYMatrix | corehelpers/inventory_test.go (F1) |
| 2 | 26.60 | TestDeferGateRollsBackUntrackedCreateAndFoldTopologies | publication/defer_gate_test.go (F2) |
| 3 | 21.70 | TestRecoverFinalizationResumesEveryDurablePhaseExactlyOnce | finalization/finalization_recovery_test.go (F3) |
| 4 | 11.00 | TestRecoverFinalizationRequiresWorkspaceMemberLockMirrors | finalization/finalization_req499_test.go (F4) |
| 5 | 7.30 | TestRecoverFinalizationReleaseReplacementPreservesModeZeroSentinel | finalization_req499_test.go (F4) |
| 6 | 7.23 | TestInventoryMatchesRetainedSecretOriginAndAmbiguityMatrix | inventory_test.go (F1) |
| 7 | 7.01 | TestDeferGateTrackedDirtyRepairFoldRollsBackExactPreimagesAfterEveryMutation | defer_gate_test.go (F2) |
| 8 | 6.11 | TestDeferGateRollsBackEveryMutationPositionToExactDirtyPreimages | defer_gate_test.go (F2) |
| 9 | 5.64 | TestRecoverFinalizationDiscoversCompleteSemanticLegacyTail | finalization_recovery_test.go (F3) |
| 10 | 4.70 | TestRecoverFinalizationIgnoresFinderMetadataUnderDoWork | finalization_recovery_test.go (F3) |
| 11 | 4.22 | TestRecoverFinalizationResumesJournalAfterLifecycleInterruption | finalization_pipeline_dirt_test.go (out of scope) |
| 12 | 3.72 | TestDeferGateClassifiesTrackedDirtyTrackedCleanAndUntrackedPreimagesIndependently | defer_gate_test.go (F2) |
| 13 | 3.05 | TestRecoverFinalizationProcessesTwoSafeGroupsInStableOrder | finalization_recovery_test.go (F3) |
| 14 | 2.91 | TestRecoverFinalizationSetsAsideRefusedRecordAndFinishesTheRest | finalization_recovery_test.go (F3) |
| 15 | 2.87 | TestFinalizeAcceptsWorkingRequestDirtWrittenByThePipeline | finalization_commands_test.go (out of scope) |
| 16 | 2.82 | TestRecoverFinalizationDiscoversLegacyNoJournalTail | finalization_recovery_test.go (F3) |
| 17 | 2.62 | TestDeferGateFoldWithoutReservationRequiresCleanCommittedRepair | defer_gate_test.go (F2) |
| 18 | 2.60 | TestRecoverFinalizationPreservesUnstagedProtectedAndDistinguishesStagedRefusals | finalization_recovery_test.go (F3) |
| 19 | 2.51 | TestRecoverFinalizationPreservesAndMatchesPrivateFileModeInLifecyclePostimages | finalization_recovery_test.go (F3) |
| 20 | 2.13 | TestDeferGateFoldRefusesStagedRepairPreimage | defer_gate_test.go (F2) |

Loaded per-file sums from the same run: defer_gate_test.go 66.85 s, inventory_test.go 64.70 s, finalization_recovery_test.go 56.71 s, finalization_req499_test.go 32.08 s. The next-heaviest files in these packages (req547 4.65 s, commands 4.48 s) are far from the budget and stay as they are.
