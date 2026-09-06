# Six reviewed fixes: upstream verification, 2026-09-06

All six fixes and their regression tests already exist in the owning repository. They shipped in **0.303.5**, release commit `35efb620bbfe920970010afb905acf4a43ccd2fb`. No replacement implementation or duplicate tests were needed.

Verification started at `6ecff1bc7aad6ef8bc80f0ecb515270b3d20c45f` (suite 0.305.31). A fresh `git fetch origin` resolved `origin/main` to `57e6431e24d7ccef5595960b81882190f4192ed2`; that revision is an ancestor of the tested checkout, which has eight additional commits. The six implementation commits and their release are already ancestors of upstream main.

After the heavy gate finished, concurrent commit `12e6a72b3f7f665b52eddfe75477e270bf53c095` added only an architecture-report HTML file. The fast retry continued on that descendant; implementation, regression tests and verification scripts were unchanged.

The authorized consumer handoff's `README.md`, `do-work.patch`, and `do-work-board.patch` were read from `/Users/t2/2code/g1w-game-find-the-difference/docs/handoffs/do-work-six-fixes-2026-09-05/`. Its v0.294.0 / `eff1679d` candidates were compared with current code and history. Consumer files were not edited, and no other repository was searched. Newer upstream finalizer and heavy-runner changes were preserved.

## Per-item evidence

All source paths below are relative to this repository. Each listed implementation commit includes its regression tests.

| Item | Verdict and implementation commit | Current behavior and regression evidence | Surface cost |
| --- | --- | --- | --- |
| F1 · P1 | Already done: `d6709da61aa33096ed93e24bd2c2e6b840df2505` | `heavyverification/heavy_run.go` checks tracked dirt and HEAD before reuse/execution decisions and after execution, before recording success. `TestLaneMutationCannotPublishOrReuseSuccess` covers dirty and committed fixture mutations, refusal before the later lane, and no published mutator success. `TestLaneQueueBookkeepingDoesNotInvalidateVerification` preserves queue-only dirty bookkeeping. | Earned by the mutating-lane replay: existing Git checks plus one helper; no new cache or persistent state. Within requested cost. |
| F2 · P2 | Already done: `b3f33798e2674ba016d0e39b5260693314eb05ea` | `publication/answer.go` rejects comment-bearing evidence. `TestStakeholderTerminalEvidenceRefusesMarkersForgedInCallerProse` independently hides each marker in multiline, inline-opening, closed and unclosed comments and requires zero planned mutations. Visible completion controls remain. | Earned by hidden-marker completion: four-line guard; no parser or schema migration. Within requested cost. |
| F3 · P2 | Already done: `4278fdb4a2ba84205103648eca00cf9adff987d6` | `lifecycleadvance/advance_commands.go` fully parses frozen queue continuations before dispatch by location. `TestExplicitQueueContinuationSurvivesFirstTargetLeavingQueue` replays the exact returned argv with the first target working and archived, claims the second once, excludes later arrivals, and checks duplicate replay. | Direct routing repair using the existing queue parser; no persistent state. Within requested cost. |
| F4 · P2 | Already done: `0b2362aebad2d910c62455def0df6e7794c731dc` | Working phases admit manifests through `finalization.FinalizeBound`; `prepareBoundJournal` enforces failure-only intent in its single identity-bound decode when completion evidence is incomplete. `TestAdvanceFinalizesEarlyFailureWithoutAdmittingEarlyCompletion` covers planner, implementation and pre-test failure archival/checkpoint cleanup, plus mutation-free refusals for premature completion, wrong identity, stale preimage and missing reason. | Direct repair through the canonical finalizer; completion gates and preimages retained. No second finalizer. Within requested cost. |
| F5 · P2 | Already done: `d10118be63524e114fa86022685df5def2f8e939` | `lifecycletiming` shares validation before launch and distinguishes recording errors. `TestInvalidTimingOptionsNeverLaunchChild` checks invalid request/run/category/operation inputs without side effects; usage failures use exit 2. `TestTimingRecordingFailurePreservesChildExitAndDoesNotReportLaunchFailure` exercises runtime exits 0, 7 and 127 after recording failure. | Existing validation moved/shared; two error categories, no retries or schema expansion. Within requested cost. |
| F6 · P2 | Already done: `725d17e4aa0cc1e078611989c16972c4e1772478` | Board `model.go` renders literal “status changed”; the destination-inference helper is gone. `TestActivityDoesNotLabelPendingRecoveryAsLaterCompletion` checks recovery at 00:11:11 and completion at 01:07:00, retaining distinct meanings. | Unsupported interpretation deleted; no history subsystem. Within requested cost. |

CLI paths above are under `skills/do-work/tools/do-work-cli/internal/`; board paths are under `skills/do-work-board/tools/queue-kanban/`.

## Fresh verification

Targeted checks passed in this upstream checkout: heavyverification and publication named regressions (`-count=1`), lifecycleadvance and finalization package suites, lifecycletiming package suite, and board Activity regressions (`-count=1`). No result here relies on the consumer report.

The first fast maintainer run passed all assertions, contracts, formatting, vet, 401 board tests and 814 CLI tests, but exited 1 on its time budget: `finalization_recovery_test.go` took 37.77s and `defer_gate_test.go` took 33.16s (limit under 30s per file). Other Go test processes were observed concurrently; contention is a possible contributor, not a proven cause.

The full heavy gate **passed**, exit 0, in 295s:

```sh
QUEUE_KANBAN_BROWSER='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' bash _dev/tests/maintainer-verify.sh --heavy
```

This ran ShellCheck, gofmt, contracts including staged-skills/updater/installer fixtures, both modules' `go vet`, 481 ordinary/strict JavaScript board tests, 35 strict browser tests on Chrome 152.0.7977.76, and 828 CLI tests with `DO_WORK_HEAVY_TESTS=1`. The Go stages used `-count=1`; the CLI stage completed in 53s. The heavy tier's documented time-budget exemption does not count as a fast-gate pass. Go was 1.26.1 on darwin/arm64.

The fast-gate retry, `bash _dev/tests/maintainer-verify.sh`, **passed**, exit 0, in 85s. It executed all 401 board and 814 CLI tests afresh and enforced the normal under-30s per-file budget. Its slowest CLI file, `finalization_recovery_test.go`, took 17.73s. Neither Go stage reused cached test success; the runner reported that reusable evidence was not recorded, which does not change these fresh test results.

All five maintainer-runtime tests absent from consumer installs **passed in the uncached upstream CLI suite**: `TestShippedLanesCoverRuntimeHelpersAndRefuseOpaqueBrowser`, `TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes`, `TestShippedRuntimeRefusesOpaqueOverrides`, `TestShippedGitIsolationPreservesGenericLaneInheritance`, and `TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint`. Current upstream keeps these in `heavy_maintainer_tree_test.go`, which reads the real `_dev/tests` manifest and runtime helper and contains no skip branches. Every item's regression result is therefore passing on current upstream, including the fixtures the consumer could not exercise.

This verification record is maintainer-only and requires no release/version bump. Unrelated queue-selection and timing-heading findings remain outside this change.

Local full logs are retained under `.git/six-fixes-verification/`: `maintainer-fast.log`, `maintainer-heavy.log`, and `maintainer-fast-retry.log`.
