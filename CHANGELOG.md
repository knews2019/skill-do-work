# Changelog

What's new, what's better, what's different. Recent release entries are here, most recent first.

For the complete release history, read this file and the archives below from newest to oldest. Archives are tracked in git and excluded from the distribution tarball; these GitHub links also work from an installed copy.

- [0.121.1 through 0.303.6](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-09-05-up-to-v0.303.6.md)
- [0.110.0 through 0.121.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-07-13-up-to-v0.121.0.md)
- [0.65.0 through 0.109.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-07-07-up-to-v0.109.0.md)
- [0.50.0 through 0.64.1](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-04-13-up-to-v0.64.1.md)
- [0.1.0 through 0.49.0](https://github.com/knews2019/skill-do-work/blob/main/CHANGELOG-2026-04-07-up-to-v0.49.0.md)

## 0.305.40 — Checkpoint Summary Cleanup (2026-09-13)

Checkpoint refresh removes obsolete generated session summaries while preserving live claims, legacy recovery evidence, and authored notes.

## 0.305.39 — Commit Evidence and Publication Recovery (2026-09-13)

Finalization preserves failed commit evidence and checks committed content against the complete prepared change, including new files. Publication failures keep every intended destination visible to rollback while preserving foreign files.

- Preserve stale reservations when committed Git authority is unavailable, including unborn repositories.
- Require actual shipped changes for releases and show image-batch failure diagnostics without changing compatibility stdout.
- Keep portable atomic publication tests active on Windows while skipping Unix special-mode assertions.

## 0.305.38 — AGY Handoff and Completion Checks (2026-09-09)

AGY guidance now covers the handoff and verification failures observed during REQ-624.

- Bound the remaining assignment and check fresh progress before interrupting or resuming.
- Review command findings independently of AGY success, with linked execution and verification evidence.

## 0.305.37 — Shorter Workflow Instructions (2026-09-09)

Work and review instructions now point to existing procedure owners instead of repeating them.

- Consolidate duplicate merge, resume, review and scope guidance while retaining command contracts, decision points and recovery evidence.
- Preserve core routing and verify simple, complex and interrupted workflow paths.

## 0.305.36 — Verified Prose Reconciliation (2026-09-09)

Workflow instructions and board guidance now match their current behavior.

- Correct verified drift in cancellation, clarification, containment, board search and copy guidance, and process-ownership documentation.
- Remove stale counts and redundant claims while preserving runtime behavior and existing contracts.

## 0.305.35 — Bundled AGY Usage Prime (2026-09-07)

Agents can reference portable Antigravity CLI guidance directly from the installed toolbox.

- Ship the AGY usage prime with model selection, headless execution, background monitoring, continuation and Boost guidance.
- Link it from core and toolbox skill entrypoints, keeping project-specific pipeline details out of the shared reference.

## 0.305.34 — Restore Unmerged Badge and Prime Discovery Fixes (2026-09-07)

Long blocked conditions now fit their cards, and lesson capture finds relevant primes even when a request's prime list is empty or stale.

- Adapt the historical badge fix to the split board client, retaining full conditions in tooltips.
- Share prime discovery across work and standalone review while preserving deferred archival writes and lesson satellites.

## 0.305.33 — Compact Prime Routing and Recurring-Lesson Traps (2026-09-06)

Prime readers load less repeated detail while keeping access to the original guidance and incident history.

- Move oversized action, shell and CLI reference material into paired lesson satellites and refresh their routing-index estimates.
- Add missing Stakes, correct the updater’s equal-version outcome, and promote eight recurring lesson families into prime traps.

## 0.305.32 — Release Planning Ignores Unrelated Missing Manifests (2026-09-06)

A root-only release can now proceed when an independent component's tracked project manifest has an unstaged deletion.

- Require manifest reads only after establishing root, suite, or workspace ownership. Missing owned declarations still refuse release planning.
- Regression coverage checks unrelated Node, Rust, and Python manifest deletions, plus missing root and nested-workspace declarations.

## 0.305.31 — Finalization Refuses a Release Whose Implementation Ships Nothing (2026-09-06)

A release is a change to shipped files. Four entries below (0.305.17, 0.305.21, 0.305.22 and 0.305.23) bumped the version for changes under `_dev/tests` alone; the finalizer checked the release payload and never the implementation. Those four entries stand as history, and this entry closes the gap.

- `complete` with a release manifest now lists what the implementation changed (the supplied commit's first-parent diff, or the finalization commit's own paths) and refuses with `RELEASE-WITHOUT-SHIPPED-CHANGE` when nothing lies under a declared module root (`suite/modules.tsv`), `suite/` or `tools/` other than release metadata. A repository that declares no modules is a consumer project and is not guarded.
- Pinned by tests for a maintainer-only commit, metadata-only paths, a merge that brings in only maintainer files, primary-commit provenance, and a consumer repository.

## 0.305.30 — Atomic Download Repairs an Interrupted Target Again, and the Inventory Launcher Forwards Only What It Honors (2026-09-06)

An audit of the queue-drain releases 0.305.14 through 0.305.25 found two shipped defects and three stale sentences; this release fixes them directly.

- `atomic-download` (0.305.15) refused every existing regular file, which closed the installer's documented repair: `install.md` reads a zero-byte `SKILL.md` from an interrupted download as absent so a re-run can replace it, and this command is what that re-run calls. A non-empty regular file is still refused with its bytes intact; a zero-byte one is replaced. Dry-run and live agree, and the prescribed-shell cases now prove both, plus the failed-transfer case that a pre-seeded target had made vacuous.
- `atomic-download` read the byte count it reports from the target after the rename, so a failed stat there reported a published file as unpublished. The count is read from the private file before the rename; a transfer that leaves no private file publishes nothing and says so.
- `scripts/protected-inventory.sh` (0.305.14) crashed under macOS `/bin/bash` 3.2 with `set -u` whenever no global flag was given (an empty array is an unbound variable there); the guard landed in an untagged commit and ships with this entry. The launcher's `--format` sifting was dead code, since the hard-coded `--format text` always won; it is removed, the guide no longer claims it, and a missing `--repo-root` value now says so on stderr instead of exiting 2 silently. The prescribed-shell case gains the `--repo-root` from-outside-the-repository proof the original request named first.
- `docs/prescribed-shell-primitives.md` describes the zero-byte repair and the before-rename byte count. The REQ-603, REQ-604 and REQ-605 records carry strikethrough corrections for the claims this audit found false, and `lessons-do-work-cli.md` records the lesson.

## 0.305.29 — Keep Cached and Fresh Test Results Consistent (2026-09-06)

Fast-stage reuse now matches the Git configuration used by fresh tests, and forced runs revoke earlier passes before execution.

- Isolate global and system Git configuration for both evidence decisions and test execution.
- Revoke old evidence before forced runs and stop if revocation fails, so failures and interruptions cannot leave a stale pass reusable.
- Make finalization fixtures select their initial branch explicitly and cover the stale-result failures with regression checks.

## 0.305.28 — Correct Test Reuse and Integrity Checks; Remove Unused Paths (2026-09-06)

Covered ignored files now invalidate test reuse even when their directory names have Git pathspec meaning, and fixture-integrity self-tests retain their failures.

- Treat coverage roots as literal Git paths, with a regression for a colon-prefixed directory.
- Preserve integrity assertion failures while clearing only expected mutation-probe increments.
- Remove unreachable per-file ShellCheck code and unused finalization tracked-path caching and wrapper code while retaining batching and committed-image caching.

## 0.305.27 — Timeline Disclosure Geometry and Markdown Version Mirrors (2026-09-06)

Keep Timeline rows visible when findings disclosures change height, and prevent releases from leaving owned Markdown version mirrors stale.

- Measure the chart's current vertical offset and repaint on disclosure toggles, including nested findings.
- Recognize `version.md` in owned release roots during planning and recovery while keeping independent components excluded.
- Add browser coverage for disclosure expansion and collapse, plus release-planning and recovery regressions for root and workspace mirrors.

## 0.305.26 — Restore Single-Image Staging Claim to Adjacent to Target (2026-09-06)

Corrected the single-image generation command description in `skills/do-work-toolbox/actions/ai-report-reference.md` to reflect that `generateImage` stages its invocation-private file adjacent to the target output path (`filepath.Dir(outputPath)`), rather than in the system temporary directory.

## 0.305.25 — Batch Repeated Git Reads in Finalization and Request State (2026-09-06)

Profiled recovery and state planning to batch repeated Git reads and memoize historical commits during finalization discovery.

- Batched `existingUntrackedPaths` queries in `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` using a single `git --literal-pathspecs ls-files -z -- <paths...>` command, eliminating over 1,180 per-path subprocesses.
- Batched `existingDirtyTrackedPaths` queries in `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` using a single `git status --porcelain=v1 -z --untracked-files=no -- <paths...>` command, eliminating over 660 per-path subprocesses.
- Introduced `discoverySession` in `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` to memoize immutable HEAD file images and tracked release paths across discovery phases, keyed to the operation entry head commit to guarantee freshness across write operations.
- Reduced redundant `git show` subprocess calls from 645 to 274 (-57.5%) and overall Git subprocesses in `TestRecoverFinalization` from 2,950 to 2,446, reducing recovery wall time to 18.13s.
- Added unit tests verifying batched path classification, session memoization, negative lookups for absent files, and cache invalidation across commits.

## 0.305.24 — Reduce Fast-Stage Evidence Computation Cost (2026-09-06)

Profiled and simplified fast-stage evidence computation in `heavy-runtime-fingerprint.py` and `do-work-cli` without weakening invalidation.

- Added in-memory caching of resolved binary seals in `_dev/tests/heavy-runtime-fingerprint.py`, eliminating redundant disk reads and SHA-256 hashing of identical executables.
- Consolidated working tree enumeration in `fast_stage_evidence.go` into a single `git ls-files` subprocess with combined `--cached` and `--others` flags.
- Scoped ignored file queries to stage coverage directory roots using `fastStageCoverageRoots`, eliminating full-tree ignored-file traversals while strictly preserving `laneCoversPath` checks.
- Preallocated internal slices and maps for stage seals to avoid dynamic slice reallocations.
- Added comprehensive unit tests and end-to-end behavior probes verifying that ignored files under stage coverage force execution while ignored files outside coverage preserve reuse.

## 0.305.23 — Remove Redundant Go Test Discovery using Native -skip (2026-09-06)

Replaced pre-execution Go test discovery in `run-go-tests-with-budget.sh` with Go 1.20+ native `-skip`, eliminating subprocess startup while strictly preserving test selection and safety guards.

- Replaces `go test -list` and Python regex assembly with pure-bash `-skip "^($skip_pattern)"`, escaping regex metacharacters with zero child processes.
- Preserves exact 100% test selection equivalence across fast test suites (verified on `queue-kanban` 402 fast tests).
- Retains empty-selection refusal (`no fast Go tests remain after applying the heavy prefixes`) in the post-test results processor with zero discovery overhead.
- Added comprehensive behavior probes in `_dev/tests/run-go-tests-with-budget-behavior.sh` for heavy exclusion, empty-selection refusal, and metacharacter escaping.

## 0.305.22 — Batch Shell Audits while Preserving Diagnostics (2026-09-06)

Batched compatible ShellCheck invocations and Markdown pattern audits across the shell verification suite, eliminating repeated process startup while preserving fragment syntax isolation and exact line attribution.

- `_dev/tests/action-shell-blocks.sh` batches ShellCheck across 74 extracted Markdown fences and 33 shipped shell files, writing `.meta` companion files to map GCC-formatted diagnostics back to exact Markdown paths and line numbers, eliminating 106 ShellCheck Haskell process spawns (-33.5% wall time).
- `_dev/tests/prescribed-shell-canonicalization.sh` scans 165 Markdown files using multi-pattern `grep -F -f` passes, eliminating over 2,600 child grep process launches (-84.1% wall time).
- Verifies defect detection across fragment syntax errors, quiet-grep pipeline violations, ShellCheck warnings, and canonicalization guide restatements.

## 0.305.21 — Use Fast POSIX Checksums for Fixture Integrity Verification (2026-09-06)

Optimized shared fixture integrity checking in `_dev/tests/session-start-hook-behavior.sh` by replacing Perl-based `shasum` with native POSIX `cksum`.

- Eliminates 20 Perl interpreter startup invocations across scenario runs, cutting CPU usage by 11.5% (-0.13s).
- Adds explicit mutation probes in `_dev/tests/session-start-hook-behavior.sh` verifying that byte changes, added files, and deleted files fail closed while preserving expected `actions/version.md` rewrites.

## 0.305.20 — Copy Prepared Recovery States in Finalization Tests (2026-09-06)

Extended fixture reuse in `internal/finalization` tests to copy prepared semantic legacy and planned recovery baseline states, eliminating repetitive Git seed commits and repository construction.

- `internal/finalization/finalization_commands_test.go` defines thread-safe singletons and lazy initializers for semantic legacy and planned finalization template repositories, cleaning them up in `TestMain`.
- `internal/finalization/finalization_recovery_test.go` updates `seedSemanticLegacyTail` and `seedPlannedFinalization` to instantiate isolated copies via `os.CopyFS`, eliminating 69 Git subprocess executions and reducing cold wall time by 13.4% (-6.45s) and CPU by 14.0% (-5.38s).
- Adds `TestPreparedRecoveryTemplateIsolation` proving complete filesystem and Git history isolation across copies and underlying templates, with mutation tests confirming robust defect rejection.

## 0.305.19 — Run Inventory Data Matrix In-Process Without Subprocesses (2026-09-06)

Decoupled porcelain byte parsing from Git acquisition in `internal/corehelpers/inventory.go`, allowing synthetic inventory test matrices to run completely in-process without spawning Git and CLI subprocesses.

- `internal/corehelpers/inventory.go` extracts `parseInventoryBytes` to cleanly separate byte stream parsing from `gitOutput`.
- `internal/corehelpers/inventory_test.go` executes the 45-case porcelain status matrix and 10-case secret origin/ambiguity matrix in-process, eliminating 56 `do-work-cli` subprocess executions and reducing test wall time by 82.5% (-7.31s) and child CPU by 61.2% (-1.01s).
- Adds explicit contract tests for malformed short porcelain records, missing rename origins, record ordering, metadata exclusions, and cross-row secret promotion while retaining end-to-end launcher and real-Git coverage.

## 0.305.18 — Build Integration Test CLI Once per Test Binary in Suiteinstall (2026-09-06)

Integration tests in `suiteinstall` now compile the `do-work-cli` executable once per test binary using `sync.Once` and clean up the temporary directory in `TestMain`, eliminating redundant `go build` invocations across heavy signal handling tests.

- `internal/suiteinstall/suite_commands_test.go` builds the test CLI binary on demand in an isolated temporary directory via `sync.Once` and deletes it upon process completion in `TestMain`.
- Retains signal interruption, recovery, and post-verification exit status assertions while reducing wall time by ~1.15s and total CPU by ~1.00s.

## 0.305.17 — Establish Test Efficiency Baseline with Descendant CPU and Work Counts (2026-09-06)

Introduced opt-in test efficiency instrumentation and reproducible baseline measurement tooling to record process descendant CPU times, toolchain subprocess counts, and per-file Go test duration attribution across representative verification selections.

- `_dev/tests/test-duration-log.sh` provides opt-in measurement functions capturing child CPU via POSIX getrusage and toolchain subprocess executions via PATH shims.
- `_dev/tests/test-efficiency-baseline.sh` orchestrates multi-run benchmarking across inventory, finalization, session-start, shell-audit, and heavy CLI build cases, emitting a structured evidence table.
- `_dev/tests/test-efficiency-baseline-behavior.sh` verifies opt-in isolation, exit code preservation, and Go test event stream parsing.
- Repointed archived lesson satellite links in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `_dev/primes/lessons-action-files.md`, and `_dev/primes/lessons-shell-commands.md`.

## 0.305.16 — List Merge Commit Paths with Diff-Tree -m in Finalization Range Check (2026-09-06)

`matchingHeadCommit` in the finalization subsystem now passes `-m` to `git diff-tree` when verifying candidate commits in the `PreparedHead..HEAD` range, preventing merge commits from emitting empty path listings and bypassing exactness path constraints.

- `internal/finalization/finalization_apply.go` invokes `git diff-tree` with `-m` so candidate merge commits enumerate paths modified across each parent against `EffectiveCommitPaths`.
- `internal/finalization/finalization_apply_test.go` verifies rejection of merge commits containing foreign paths and acceptance of clean merges.

## 0.305.15 — Unify Atomic-Download Occupancy and Handle Post-Publish Stat Error (2026-09-06)

The `atomic-download` core helper now enforces the same occupancy rule in live execution as it does during dry-run, refusing an existing regular file with exit 2 (`DOWNLOAD-TARGET-OCCUPIED`) rather than silently replacing it, and checks the error on post-rename `os.Stat` to prevent nil pointer panics under concurrent races.

- `internal/corehelpers/commands.go` checks destination occupancy before branching on dry-run, ensuring both dry-run and live runs refuse directories with exit 1 (`target is a directory`) and existing regular files with exit 2 (`target already exists`).
- `internal/corehelpers/commands.go` handles errors on post-rename stat, returning typed failure finding `DOWNLOAD-FAILED` instead of panicking on nil dereference.
- `prescribed-shell-primitives.md` documents the unified occupancy rule refusing existing destinations before fetching and the error-checked stat.

## 0.305.14 — Protected Inventory Launcher Passes Global Flags and Preserves Diagnostic Text (2026-09-06)

The `protected-inventory` compatibility launcher now forwards global flags (`--repo-root`, `--format`) ahead of the command token so it can run from any working directory, and the compatibility shim preserves diagnostic text and findings on failures instead of discarding them.

- `protected-inventory.sh` parses and forwards global options before the command verb, allowing `--repo-root` to be passed before or after the mode (`start` or `associate`).
- `internal/corehelpers/inventory.go` compatibility shim preserves prepared diagnostic text (`NO-DO-WORK-DIR`, `PARSE-FAILED`) and walk error findings on failure, ensuring exit 2 causes are visible.
- `commit.md` and `inspect.md` document that exit 2 with a finding or diagnostic text is an error to report rather than a skip condition, clarify exit 2 causes (such as non-git repository, status failure, or quarantine write failure), and note that `associate` exits 2 as not-started after `start --dry-run`.
- `commit.md` clarifies that re-running inventory uses `start` (which replaces quarantine rather than appending, whereas `associate` unions).
- `prescribed-shell-primitives.md` documents global flag forwarding and `--repo-root` support in the launcher.

## 0.305.13 — Correct Stale Mechanism and Script Claims Across Shipped Callers (2026-09-06)

Documentation and action files across toolbox, board, and core now accurately describe underlying Go command implementations, staging paths, error conditions, and conflict resolution semantics.

- `ai-report-reference.md` documents that single-image and batch image generation stage in the system temporary directory rather than adjacent to the target, and removes nonexistent retained script fallback references.
- `install.md` documents temporary skill downloads landing in `SKILL.md.download.<random>`, specifies that local Git ignores are managed inline following the canonical contract rather than calling an external helper, and removes nonexistent compatibility script references.
- `present-work.md` clarifies that portfolio summary publication reads the source once and writes both outputs from that buffer, and removes nonexistent compatibility script references.
- `board.md` clarifies that adding a local git exclude on non-git projects exits zero with a `GIT-EXCLUDE-NOT-A-REPOSITORY` warning finding.
- `prescribed-shell-primitives.md` accurately describes inventory conflict resolution: a parseable `completed_at` beats an unparseable or missing timestamp regardless of walk order, and ties fall back to `working/` before `archive/`.
- `architecture-report.md` removes nonexistent compatibility script references.
- `work-reference.md` clarifies that `run-timed-command` connects both child streams directly to the CLI's stderr handle.

## 0.305.12 — Rollback Decides the Handle Once and Closes Nil-Handle Panic (2026-09-06)

Transaction rollback now settles the rooted filesystem handle once upon entering `rollbackFailure`. If `os.OpenRoot` fails, the transaction immediately executes `rollbackWithoutRoot` to perform Git unstaging, tracked restoration from HEAD, and preimage restoration by pathname while safely recording unavailable root errors for rooted mutations. The missing guard in `quarantineAndRollbackPrivate` that caused a nil-pointer dereference panic when rolling back identity-recorded private untracked files is eliminated, and all eight downstream nil checks across rollback loops 1–3 were removed and pinned at zero.

- `rollbackFailure` cleanly partitions into `rollbackWithRoot` and `rollbackWithoutRoot`, ensuring that rooted operations only execute under a valid handle.
- All eight defensive nil-root checks in `git_transaction.go` have been removed, and `_dev/tests/audit-lockins.sh` Finding 3 now pins the count at zero.
- Added `TestRollbackWithoutRootHandleUnstagesRestoresFromHeadAndReportsTheRest`, providing the package's first unit test verifying rollback behavior when the root handle cannot be opened.

## 0.305.11 — The Rest of the Shell Guide Accurately Describes Shipped Primitives, and Inspect Associates Files (2026-09-06)

Sixteen stale claims across five sections of the prescribed shell primitives guide, plus the association prose in commit and inspect, now truthfully describe what the Go implementations and syscalls do. The two prescribed blocks in inspect that called the protected inventory wrapper with positional arguments now pass --quarantine-name, so inspect associates files with requests rather than failing with an unknown option error.

- inspect.md now passes `--quarantine-name do-work-inspect-secret-quarantine` to `protected-inventory.sh` rather than positional arguments, enabling the action to associate modified files with requests.
- Both `commit` and `inspect` now accurately document that the protected inventory wrapper takes no `--repo-root` argument, uses the current directory as the repository root, drops untracked hidden files under `do-work/` as metadata, and prints no placeholder rows for unassociated candidates.
- The prescribed shell guide's Verified Exact Publication section accurately distinguishes syscall semantics (`rename(2)` silently replacing a regular file while refusing a directory, `link(2)` refusing all existing targets) from shell `mv` and `ln` nesting.
- Lifecycle timing claims now truthfully state that elapsed seconds are derived by the command, that child processes inherit the CLI's stderr handle directly, and that redaction applies to the executable name rather than caller parameters.
- Commit file listing guidance clarifies quoting rules and the necessity of `--root` for the initial commit, and documents that `publish-portfolio-summary` and report image batch generation run through Go commands rather than non-existent retained shell scripts.

## 0.305.10 — Every Shipped Shell Block Is Now Checked for the Pipe Shape That Hides a Failed Command (2026-09-06)

0.305.5 removed `producer | grep -q PATTERN` from the repository's own scripts because, under the shell setting they all run with, the producer's death reads as the pattern's absence. One shipped block still carried it: the probe in the memory guidance that asks a local model server whether an embedding model is pulled. It could not misfire in practice, because that listing is far smaller than the size where the shape starts to fail, but it is the block agents copy, and copied guidance is where the shape spreads from.

- The probe now captures the listing first and searches the captured text. A missing or stopped server is the "no local model" answer the surrounding guidance already gives, and the block says so where the reader sees it.
- The maintainer check that already syntax-checks every shell block inside the shipped guidance (74 blocks across 32 files) now scans each one for the shape and fails at the block's own line in the Markdown file. Until now that guard read only shell scripts, so a block in guidance could ship carrying what a script could not.
- The scanner has one home, shared by both checks, so the two cannot drift apart. The check's own fixture still carries the nineteen forbidden spellings and seven safe ones, and a one-shape wiring test runs on every ordinary invocation, so removing the scan from the guidance check fails loudly rather than passing quietly.
- The maintainer's shell guide now describes the trap in the place an author reads before writing shell: the condition, why it is wrong in both directions, the measured window, the fix, and the reader shapes the guard cannot see.

## 0.305.9 — A Maintainer-Only Test No Longer Ships, and a Skewed Claim Keeps Its Progress Figures Ticking (2026-09-06)

Two fixes from an external review of the open pull request, applied directly.

- One test in the shipped `do-work-cli` module read this repository's own maintainer tree and skipped only when a directory six levels above the package was missing. In an installed copy that directory is `<project>/.claude/_dev/tests`; a project that happens to have it ran the test and failed on a manifest that never ships. The test now lives in the export-ignored file with the other maintainer-tree tests, so an installed module's `go test ./...` cannot meet it.
- On the board, a request claimed with a timestamp ahead of the viewer's clock is shown as clock skew rather than counted. Its user-request summary was never subscribed to the one-second refresh, so when the clock caught up the card's stopwatch moved while the header and drawer figures stayed on clock skew until something else redrew them. The summary now subscribes while a claim is skewed and recovers on the next tick.

## 0.305.8 — Archived Requests Stay Archived When the Project Sits Under a Directory Named working (2026-09-06)

The commit action asks which in-flight request owns each changed file. That answer came from the wrong place: the tool decided whether a request was in flight by looking for `/working/` anywhere in the request file's absolute path, so a project checked out beneath any directory called `working` made every archived request look active. Archived requests that had been cancelled or blocked could then claim files they never owned.

- In-flight-ness now comes from which directory is being walked, `do-work/working` or `do-work/archive`, and never from the path's spelling. The rule for an ordinary checkout is unchanged: a working request counts whatever its status says, an archived one only when it finished successfully.
- A request file stored under `do-work/archive/working/` was caught by the same mistake and is now treated as what it is, an archived request.
- A test builds its checkout under a `do-work/working/` directory, so any future fix that reads the path instead of the walked root fails it.
- When two requests claim one file, the later completion time wins; on an equal or missing time the request found first stands, working before archive. That order is now stated in the code where it is decided.

## 0.305.7 — One Redundant Rollback Check Removed, Eight Load-Bearing Ones Kept and Pinned (2026-09-06)

A maintenance audit counted nine places where the transaction rollback tests a filesystem handle for nil and asked for eight of them to go. Tracing every path that can reach them showed the opposite: eight are the only thing between a failed handle and a crash in the middle of restoring the tree.

- When the rollback cannot open its rooted handle it records that and keeps going, so a failed transaction still unstages its paths and still reports what it could not restore. That one missing handle then reaches eleven places, and each decides for itself what to do without it. Removing any one of seven of those checks turns a reported incomplete rollback into a crash mid-rollback.
- Exactly one check re-tested a handle every caller had already settled. That one is gone, and its precondition is stated where the callers can read it.
- A maintainer check pins the count at eight in both directions and counts the checks themselves, not any line that happens to mention them, so a deleted check cannot be paid for with a comment.
- Found on the way and not yet fixed: one place the rollback hands that possibly-missing handle onward with no check at all, so a transaction with a recorded private untracked target can crash during rollback when the handle cannot be opened. That is queued as its own fix with the package's first test for the no-handle path.
- This change was already present in 0.305.5 and 0.305.6, which were released on top of it while it was held for review; it changes no behaviour a user sees.

## 0.305.6 — The Manual Fallbacks in the Shell Guide Now Give the Tool's Answer (2026-09-06)

Two procedures in the shell guide tell you what to do by hand when the inventory tool will not run. Both gave different answers from the tool. That matters more than an ordinary documentation slip, because a person only reaches for them when the tool is already unavailable.

- The pattern list for secret-shaped filenames was narrower than the code's, so a by-hand inventory left `credential.json` unquarantined where the tool excludes it. A maintainer check now derives the expected patterns from the code itself and fails if the guide does not name every one — it catches a pattern added to the code, a pattern dropped from the guide, and the code's own function being renamed away.
- The by-hand classification said a renamed file's new path is always "modified". The code checks for deletion first, so a rename whose new path is then deleted is a deletion — and the guide told you to run a diff on a file that is no longer there.
- The by-hand association filtered out every in-flight request, which is the exact failure the paragraph above it exists to prevent. It also used a shell glob that skips the request files sitting directly in the archive directory, failed outright on a project that has never archived anything, and could name a different owner from the tool when two requests tie, because the file search returned a different order.
- The rule for resolving that tie said an archived request outranks an in-flight one. Nothing in the code compares those; the latest completion time wins, and on a tie the in-flight one is what stands.
- The excluded-file tag now has a reading rule of its own. An ordinary new file is marked excluded whenever a secret is present in the same inventory, and until now the closest rule told you to read it.
- Two descriptions of publication were also corrected: only one of the two commands verifies what it wrote, and the other reports a byte count without checking anything.

## 0.305.5 — Maintainer Checks Stop Reporting a Passed Check When Their Input Died (2026-09-06)

A shell check written as `producer | grep -q PATTERN` reports the producer's death as the pattern's absence. Under the shell setting these checks all run with, that makes the check report the wrong answer — and it is wrong in both directions, so a check written to fail when it finds something can quietly pass instead. It is invisible below roughly 36 KB of output from the producer and certain above about 200 KB.

- 129 such checks across 22 files are now 0. Ninety-nine were replaying text already captured; the other thirty ran a real command, and each of those now checks whether that command succeeded rather than trusting what it printed.
- **Two of the 129 were not waiting to break — they were broken.** A check that the old runtime is gone reported it gone whenever its file search failed for any reason, which is what happens when one directory in the tree cannot be read. Both now report the failure.
- The guard that forbids the shape moved out of the one probe file it lived in and now covers every shell file the repository tracks, running on an ordinary commit rather than only in the heavy tier. It costs a fifth of a second.
- The guard's own test carries nineteen ways of writing the forbidden shape and seven safe shapes it must leave alone, each named, so removing part of the guard reports which spelling it stopped catching rather than that a number moved. An earlier version of the guard could be walked past five different ways, including by a blank line.
- The guard's header states what it cannot see — a reader that is not `grep`, a pipeline assembled while the script runs, and three shapes it can flag by mistake, all of which fail loudly rather than passing quietly.

## 0.305.4 — Six Helper Names Defined Fifteen Times Now Have One Home Each (2026-09-06)

Six small helpers were copied across the tool's internal packages, and three of the copies disagreed with each other. That is how a version comparison ended up returning opposite answers in two places and a safety check ended up switched off without anyone noticing.

- One shared package now holds the four helpers more than one package needs. It imports nothing else inside the tool, so any package can use it and no copy has a reason to come back.
- The two copies of the version comparison returned **opposite** results for the same pair. One place asked "is the new version greater" and the other asked "is the old one older", and each had written its check to match its own copy. There is one comparison now, and it reports separately whether the version could be read at all, so nothing can mistake "not a version number" for "the same version".
- One copy was not a helper at all. It ran a path validator and threw the error away, returning nothing for the whole list whenever any single path was rejected — which quietly switched off the check that a finalization lists every file it plans to commit. That check is back on and has a test that fails if it is switched off again.
- Two path resolvers shared a name and disagreed about whether a missing file is an error. They keep their behaviour and now have two names, because the one that treats absence as an error uses that error as its existence check.
- A version string a project writes by hand — `1.09.0`, `1.0.` or `1.0.x` — is now refused with a clear message where the looser comparison used to guess at an order. Measured across 441 version pairs: 242 change, and every one of them changes to a refusal rather than to a different silent answer.
- A maintainer check pins the count of these definitions in both directions, so a seventh copy fails the build and so does losing one.

## 0.305.3 — The Shell Guide Says What Each Command Actually Does (2026-09-06)

The guide's table of shipped executable homes described work that moved into Go years ago, or credited a command with a step something else performs. All fourteen rows were checked against the code that implements them; three were wrong.

- `run-blocked-check` was described as selecting a GNU `timeout` binary and building a Bash process group. It does neither: there is no `timeout` lookup anywhere in the tool, and the process group comes from the Go runtime. Its entry now names what it owns — the probe's process-group timeout and kill escalation, first-hand launch and timeout evidence rather than facts guessed from an exit code, a bounded diagnostic identity, and the baseline comparison a focused test run needs.
- `install-memory-hooks` was credited with verification and rollback. Both are still steps a person follows in the memory actions; the command stops at a backup and a rename. Its entry now names the per-event gating and the settings merge that keeps your own keys in order.
- `record-timing-event` was credited with the folded per-request summary, which a different command produces — and the prose lower on the same page already said so, so the table had been contradicting its own document.
- Two descriptions of how a report's images are published were also wrong. Publication is not a single rename: the command claims the output directory outright, writes each finished image into it, and deletes the whole directory if any write fails, so a half-published directory never survives. And the staging directory is a temporary one, not a directory beside the output.
- A maintainer check now fails when an entry in that table credits a command with shell machinery, since every entry in it is a Go command. Its own comment states the case it cannot catch: an entry that claims another command's work reads exactly like a true one.

## 0.305.2 — The Shell Guide's Table Points at the Command, Not at Its Launchers (2026-09-06)

The "Shipped executable homes" table sent readers to nine shell scripts for mechanics that live in Go. Each of those scripts is a few lines that hand the work straight to a `do-work-cli` subcommand, so the table was routing people to the wrong file.

- All nine rows now name the `do-work-cli` subcommand that owns the mechanic. Each subcommand was read from its launcher's own `exec` line, not guessed from the file name.
- The sentence claiming that a six-line launcher orchestrates two other check scripts is gone. It was false twice over: that launcher hands everything to the Go command, and the two scripts it was said to drive are themselves launchers over the same command that nothing in the Go tree ever starts.
- The guide now explains the difference between the two, because it matters. The table says where a mechanic is owned; the launcher is what the guide's own instructions and the shipped actions invoke. They are not interchangeable — six launchers translate a legacy positional call into the flags the subcommand needs, and the protected-inventory launcher additionally selects the tag-and-path output that two actions parse one row per file from. Change the command's flags or its output and you have changed the launcher's contract.
- A maintainer check keeps both shapes from returning, and it asks what a row means rather than how it is spelled: it finds the table by its own header row, and it opens each script a row names to decide whether that script is a launcher. A row naming a file that does not exist fails, and a table with no rows fails. The first version of this check could be walked past seven different ways, including by the table's own formatting.

## 0.305.1 — The Fast Gate's Skip Decision Reads the Queue It Builds From (2026-09-06)

The stage-reuse seal added in 0.305.0 inherited a rule that ignored the whole `do-work/` tree. One of the two stages builds the Kanban board from that tree, so editing a request could leave the gate reporting a pass from a run that never saw the change.

- The board stage now declares `do-work` as an input it reads, and the tool stage declares the single archived file it reads there. Neither is assumed: the board's own prune rules and its file-mention list were walked to confirm nothing else is read.
- A separate exclusion list keeps four churn paths from invalidating a stage that never reads their bytes. The load-bearing one is the gate's own test-duration log: the stage appends to it while running, so a seal covering it could never match again and the stage would never be skipped.
- The two tests that pinned the old behaviour were rewritten rather than deleted, and each names the failure it now catches. The "queue changed" case became two — one proving the stage that reads the tree runs, one proving the stage that does not still skips — so the fix cannot become "make every change invalidate everything", which would have removed the saving entirely.
- After a review found the shipped input list itself was read by no test, restoring it to its pre-fix content now fails with three named messages, deleting either stage's rule fails with the message naming that stage, and broadening an exclusion until it swallows a real input fails too. All four of those used to pass silently.
- A typo in an exclusion's kind used to decode cleanly, match nothing, and turn skipping off for that stage forever. It is now rejected when the file is read.

## 0.305.0 — The Fast Gate Skips a Stage Whose Inputs Have Not Moved (2026-09-06)

Running the fast gate twice in a row used to re-run everything, even when nothing a stage reads had changed since the last green result. It now records what each stage read and skips a stage whose complete inputs are unchanged, printing one line per stage saying what it did and why.

- The seal is over the working tree, not over committed content. The fast gate exists to run on a tree with uncommitted work, so sealing commits would report a false pass the moment anything was uncommitted.
- Records are kept in their own key space, keyed by stage and by working-tree root, so two sibling worktrees cannot invalidate each other's evidence.
- Three commands expose the decision, the recording and the invalidation, and the gate wraps its two Go test stages in them. Anything the engine cannot determine forces the stage to run — an unreadable or missing decision never reads as a skip.
- Separately, the SessionStart hook probe stopped copying the tool's module once per scenario and now builds one shared tree. The cost was never the copying: each copy created a new absolute path, which the Go toolchain treats as a different build and re-links from scratch. Each scenario still keeps its own writable state, the banner input is now a required argument so forgetting it is an error rather than a silent pass, and a new check proves the shared tree is unchanged after every case.

## 0.304.7 — The Timeline Scrolls With the Board Instead of Inside Its Own Box (2026-09-06)

The Timeline's rows used to scroll inside a box about half the window tall, so the page had two scrollbars and the reader had to find the right one. The rows now scroll with the board, and the time axis stays pinned to the top edge the way the Activity view's column header does.

- The chart no longer has a height cap or its own scroll area. Keeping either would have left it a scroll container while measuring as though it were not.
- The time axis has a sticky rule of its own, painted on the base background, and carries the one-pixel separator the rows box used to own, so the line stays on screen instead of scrolling away with the content.
- The board's top padding now sits on whichever child the reader sees first, decided by that condition rather than by a list of the optional strips that can come before the view.
- Everything that reads or writes a scroll position moved into the board's coordinate space: the top-visible-row anchor and its restore, the virtualized visible range, scrolling a focused row into view, and jump-to-open-work. The two geometry numbers those need are measured once per render and refreshed beside the existing width cache, because the visible-row render is the scroll listener and extra layout reads would land on every frame of a drag.
- Only the scroll listener moved. Wheel zoom, drag-pan, pointer capture, keyboard, focus, hover and the width observer all stay on the chart.
- Entering the Timeline now resets the board's scroll position, which nothing did before.
- Two defects found by measuring rather than reasoning are fixed: removing the chart's focus ring handed the browser's own default ring the full-height box, caught by a screenshot; and clamping the board-relative scroll at zero made the anchor's write disagree with its read, so pressing a window chip while the board sat above the chart jumped it down by the whole offset.

## 0.304.6 — Three Lifecycle Behaviours Are Now Held By Tests That Name What They Catch (2026-09-06)

Three things the lifecycle gate does could be deleted from the code and every test would still pass. Each one now has a test that fails when it is deleted, and each test says in its own comment which deletion it catches.

- When a finding's suggested command belongs to a helper rather than to the current step, the gate rewrites it to point at the step you are actually on. That rewrite is now checked at both places it happens, by reading the rewritten command itself rather than a summary line beside it.
- Three negative controls prove the rewrite is selective rather than blanket: a neighbouring finding whose command must stay empty, a git command whose own wording must survive untouched, and one finding where one field is rewritten and the other deliberately is not.
- The guard that decides whether a focused test run counts as evidence is pinned by a nine-row table keyed on the facts the code actually reads — did it launch, did it time out, what was the baseline — rather than on the exit numbers that happen to produce those facts today. Four of the nine rows prove the guard still admits a valid run instead of refusing everything.
- An interrupted focused test is pinned by a case that signals the tool from inside the probe, paired with an ordinary failing run at the same exit status, so the only difference between them is whether the interruption was observed.
- No behaviour changed. Each of the three behaviours was temporarily broken to watch its new test fail, then restored and checked byte-for-byte.

## 0.304.5 — The Heavy Test Tier Stops Reporting a Failed Run as a Passed One (2026-09-06)

The heavy verification tier could tell you a lane passed when it had failed, and could tell you a run succeeded when it had verified nothing. Three separate routes to a false green are closed, each with a test that fails when the fix is taken back out.

- A lane that ran and exited with an error is now reported red, whatever it printed while running. It used to be enough for the lane to print a line starting with `SKIP:` — including a line printed by a test's own fixture — for the whole lane to be recorded as skipped, which reads as success. The lane's own announcement is kept and shown as extra evidence on the failure.
- A run that names no lane at all is refused instead of returning success. Before, a caller inside the tool that forgot to list lanes got a clean verdict for verifying nothing.
- A lane timeout of zero now means "use the default" rather than "expire immediately". The thirty-minute default only existed on the command-line path, so a caller inside the tool that left the field empty had its lanes killed mid-run — which then surfaced as a failing test that looked like bad luck.
- A test fixture's output can no longer reach the process running it. The check for that swaps the real output stream rather than the file descriptor, so it fails if the fix is reverted; the earlier evidence for the same fix could not tell the difference.
- The maintainer check that forbids a particular fragile shell pipeline was rewritten. The old pattern missed five ordinary ways of writing the same thing; the new one matches what makes the pipeline fragile instead of one spelling of it, and states in its own source the two shapes no source scan can catch.
- Two archive checks were reading a file listing while throwing away the error that said the archive could not be read. Both now check readability first, so a truncated archive fails instead of passing.

## 0.304.4 — The User Request Progress Figures Say When They Are Guessing (2026-09-06)

An independent review of 0.304.0 found the new per-user-request progress strip printing a confident number in two cases where it did not have one, and found the browser probe that was supposed to prove the layout measuring a page two of its three widths never had.

- A user request whose members have all run past their estimates now reads `~0 min (2 over estimate)` instead of a bare `~0 min`. On a real board this was 4 of the 5 user requests with a live claim, one of them 9.4 hours against a saved 20-minute estimate, and the reader could not tell "almost done" from "every member has blown its estimate".
- A claim timestamp the board rejects is now counted as unknown remaining time rather than as zero time spent. Before, the member was quietly charged its full estimate and the figure rendered with no qualifier, even though the same rejected stamp already showed a warning on the Active figure.
- The browser probe measures the real layout again. It reuses one page across three widths, and two things leaked between measurements: a detail drawer left open, which is a grid column rather than an overlay at these widths, and the probe's own result block, which lands in the same grid and can be a 32,000px-wide line. The group column measured 273 / 259 / 579 CSS px; it now measures 273 / 697 / 1209 again, and the probe fails if the measured box does not match the width its own case claims.
- Three claims that were written as if a test enforced them are now enforced by tests: the order of the two refresh passes inside a tick, and the absence of `try`/`catch` and of a completion timestamp read in the summary code path.
- A misspelled field name is corrected so a plain-text search finds every use of it.

## 0.304.3 — The Debug-Artifact Rule Is Stated Once, Where It Is Enforced (2026-09-06)

Three shipped action files repeated the same debug-artifact and P-A-U honesty rule seven times between them. The rule is enforced in code, so every prose copy was a restatement that could drift away from what actually runs.

- Seven mentions become two. `work.md`, `review-work.md` and `work-reference.md` now point at the finding codes `QUALIFY-DEBUG-ARTIFACT`, `QUALIFY-PAU-UNCHECKED` and `QUALIFY-UNIFY-DISARMED`, which is where the rule is decided.
- Two mentions are kept on purpose: `review-work.md`'s standalone-review hygiene bullet is a read the canonical `qualify` never makes, and the P-A-U template payload is byte-identical across four shipped files.
- The request claimed nine sites and the tree has seven. The commit its reproduction line names is not in this clone, so the captured failure could not be replayed and no lines were invented to reach nine. One anchor had moved eight lines earlier in the same run, so every edit was located by text rather than by line number.
- A maintainer lock-in pins the count of remaining mentions with a floor as well as a ceiling, so losing one of the two protected mentions is no longer silently green. It counts matches rather than lines — the earlier version failed on a pure reflow — reads its scanner's exit status without a pipeline, and reports path, line and matched text for each site.
- A companion assertion checks every `QUALIFY-*` code named in the action files against the code that defines it, so the enumeration goes stale loudly instead of quietly.

## 0.304.2 — The Commit and Inspect Actions Stop Repeating the Same Shell Rules (2026-09-06)

Two shipped actions carried the same block of file-inventory rules word for word, so a correction to one left the other behind. That block now lives once, in the prescribed-shell guide, and both actions point at it.

- The inventory tag legend, the secret-shaped matching sentence, the four file-reading bullets, the association semantics and both manual fallbacks are stated once in a new section of `prescribed-shell-primitives.md`. `commit.md` and `inspect.md` reference it instead of restating it.
- What each action *does* with an excluded row stays where it was — a writing action lets only the deletion proceed, a read-only one inspects the path and the deletion state. That is caller policy the guide's own charter excludes, so collapsing it would have been a behaviour change dressed as deduplication.
- The by-hand association fallback now names the two directories it tells you to glob, checked against the roots the code really walks, and the quarantine sentence names the `git rev-parse --git-path` call that produces the failure, confirmed by running both modes outside a repository.
- A maintainer lock-in now asserts each moved passage is absent from each action on its own, so a return to either file is caught by itself. It was proven in both directions: the legend pasted back into one action gives four failures naming that file, the whole pre-move file restored gives eight. The earlier line-count version could not see a one-sided return at all, and a single deleted line could push its similarity score past the ceiling — 660 fuzz runs now trip nothing.
- The lock-in reads its scanner's exit status instead of a piped total, so a scanner that cannot run fails instead of reading clean.

## 0.304.1 — The Archive Audit and Report Publishing Stop Launching Coreutils (2026-09-06)

Two places in the tool started a `find` or a `cp` process to do work it already does in Go. Both now do it directly, so they behave the same on a machine where those programs are missing, older, or a different implementation.

- The archive timestamp audit checks that it can read the archive with its own walk instead of launching `find`, and still stops with a clear message naming the path it could not read.
- Publishing an architecture report stages the draft in memory instead of launching `cp`; the published report carries the same bytes as before.
- Two maintainer test cases used to drive those failures by putting a fake program on the path, which stopped working the moment the code stopped launching one. Both now provoke the same failure inside the process, and each names the failure it expects, so neither can quietly pass on something else going wrong.

## 0.304.0 — User Request Groups Fold, and Report Their Own Progress (2026-09-06)

A user request with dozens of requests under it filled the board with cards and still told you nothing about how far along it was. Its groups now fold, and its header answers the question the card wall never did.

- Every **By UR** group header folds its own card grid, and starts open. More than one can be folded at a time, the control is a real button that announces whether it is open, and **Details** beside it still opens the drawer. The **URs only** reading is unchanged, including its collapsed default.
- The By UR header and the user request's drawer both show the same five figures for the whole request: how many requests it groups, active time already spent, an approximate remaining time, and the successful and resolved percentages with their counts. Both read the request's complete membership, so filters change which cards you see and never move the numbers.
- Active time ticks with the rest of the board while the page is open, so the header can never drift from the stopwatch on a claimed card below it.
- Missing evidence is stated, never counted as zero. A refused span, a member whose work ended with nothing measurable, an unfinished member nobody has estimated, and a claim stamped ahead of your clock each get their own qualifier, and a request with no members reads `unavailable` instead of dividing by zero.
- The drawer's grouped REQ id list starts open, is height-capped so it can no longer push `input.md` and the body out of the panel, and folds away entirely with one click.
- The board now reads each request's saved `estimate.p50_active_minutes` for the remaining-time figure, falling back to the Timeline's median only while the Timeline has enough history to call it confident. The Timeline's own forecasting is unchanged.

## 0.303.10 — The Top Bar's Controls Stay On Screen at Narrower Widths (2026-09-05)

Making the identity one unwrapping line (0.299.0) also gave it a minimum width, and on a window between roughly 760px and 1000px wide that pushed the filters and view buttons past the right edge. The page hides horizontal overflow, so there was no scrollbar to reach them — measured at 800px, the bar's contents wanted 905px and the last 105px were simply gone.

- The bar now stacks the identity above the controls below 1000px, where it used to stack only below 760px, so both halves stay reachable across that whole band.
- A project directory name longer than the line can carry truncates with an ellipsis instead of pushing the controls off-screen, and the controls take whatever width the identity does not want.
- Above 1000px the bar is the same single line 0.299.0 delivered; a browser probe measures both widths, so neither can be fixed at the other's expense.

## 0.303.9 — A Blocked Request No Longer Strands the Ones Behind It (2026-09-05)

The last two findings of the same review that produced 0.303.5. Both are cases where one item's problem was silently applied to everything after it.

- A frozen queue continuation used to stop at the first member it could not claim, so a request awaiting answers hid every sibling behind it — including one the just-claimed prerequisite had made ready. Because the continuation is re-emitted unchanged, replaying it stopped in the same place and claimed nothing. The pass now skips that member, leaves it in the queue, still reports why it was excluded, and carries on down the list.
- An ordinary `### ` sub-heading in a request body was read as an opening fence, which hid the request's own `## Timing` section from replacement: the fold left the stale summary and appended a second heading. Headings enclose nothing, so they no longer open one.
- Both lessons are recorded in `lessons-do-work-cli.md`.

## 0.303.8 — A Verify Finding Names Five Paths, Then Counts the Rest (2026-09-05)

One long list could take the board over. A builder worktree whose `do-work/` was untracked put roughly 700 paths into a single finding, and every surface that prints a finding printed all of them: the board's findings strip filled the page, the terminal report put them on one line, and the shareable snapshot carried the lot.

- Three `queue-kanban verify` findings now name five entries and then say how many they did not name — uncommitted queue state in a builder worktree, queue state committed on a builder branch, and the live members of an archived user request.
- A list of five entries or fewer reads exactly as before, so a finding about a single file is unchanged.

## 0.303.7 — The Verify Findings Strip Is One Line Until You Open It (2026-09-05)

This morning's rows (0.303.2) were still too tall: every remedy printed in full, all the time, inside a card. A remedy is what you read after deciding to act, so it now sits behind a click, and the strip is one line until you open it.

- Closed, the strip is a single line: a small warning mark, the counts, and every subject with a coloured dot for its weight (amber: read this, green: cleanup can fix, grey: a probe that never ran), with a Show button at the right.
- Open, each finding is one row: dot, subject, category in plain words, the detail clipped at the line's end, and a chevron that opens that finding's "What to do" under the row.
- The board remembers whether you left the strip open or closed.

