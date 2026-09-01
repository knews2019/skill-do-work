# Re-Review: REQ-418 (post-remediation)

**Approve with follow-ups** — the sole remediation closes all nine Important findings on substance, and I reproduced eight of the nine original counterexamples myself against the merged binary and confirmed each now refuses, preserves, or reports correctly. One finding is partially closed, and the remediation introduced three new Important defects in `architecture.go`, two of them reproduced regressions against the pre-remediation build.

Route C | remediation `a43b2587`, merged `82534d36` | reviewed at HEAD `d924b2fc` (product paths byte-identical to `82534d36`)

## What's built

The remediation adds a single confined-publication family (`os.Root` ancestor validation, pinned parent handles, private staging, exclusive hard-link claims, owned-directory binding) used by every file publisher; identity-and-content-bound tracked rollback in `gittransaction`; one invocation context spanning generation, staging, commit, and rollback with owned-process-group escalation; snapshot-first portfolio publication; separated display/filesystem paths in architecture scan; a leading-only mutation-option region with `--` escape; and option-presence tracking that restores malformed-flag parity with the standalone audit oracle. Thirteen named `TestRemediation*` tests were committed, and all thirteen run and pass on Go 1.26, including under `-race`.

## Decisions / risks for you

- The nine safety and parity contracts the initial review named are met. The reproduced data-loss, outside-root-write, and surviving-descendant behaviors are gone. That was the substance of the initial `Request changes`, and it is genuinely resolved.
- What blocks a clean Approve is collateral: `architecture-report-preflight --publish` now spins in an unbounded CPU-burning loop when its bundle claim fails for any reason other than "already exists", and its `--commit` mode is dead in every argument position. Both worked before the remediation. Both live in the one file the remediation rewrote most heavily, and neither is covered by a test.
- Four of the five usage strings still advertise the trailing `[--dry-run|--commit]` form that the finding-8 fix deliberately stopped accepting, so the commands now misdocument themselves.

## Findings

### Original findings 1–9

**1. Linked ancestors escape the declared output root — CLOSED.**
`repositoryPath` (`portfolio.go:163-174`) now calls `validateNoLinkedAncestors` before returning, and it is the shared entry for architecture, portfolio, and report-image targets; `note.go:26` validates its fixed target directly. Publication runs through `rootedPublishFile` / `rootedPublishInOwnedDirectory` (`mutation.go:127-275`), which open a pinned parent handle beneath `os.OpenRoot` and never traverse an ordinary absolute path. Probe: real Git fixture with `reports` symlinked to an external directory, `architecture-report-preflight --publish draft.html reports/run` exits 2 with `confined path contains a symbolic link: reports`; the external directory was empty before and after. Same refusal reproduced for `do-work-note` and for `publish-portfolio-summary` in both `--canonical-only` and `--with-snapshot` modes through a symlinked `do-work` ancestor. Test: `TestRemediationRootedBoundaryRejectsLinkedAncestor`.

**2. last30days replaces state before establishing eligibility — CLOSED.**
`preflightLast30DaysTarget` (`last30days.go:250-269`) runs project-containment, linked-ancestor, and `git ls-files` tracked checks before the clone and before any mutation (`last30days.go:51-53`); `prepareLast30DaysExclude` (`last30days.go:271-332`) then establishes and verifies the local exclude with an identity-and-bytes-bound undo, before publication. `publishLast30Days` claims the target with an exclusive `os.Mkdir` and verifies a complete payload before returning. Probe: project with a tracked, dirty, incomplete `.claude/skills/last30days/SKILL.md` exits 1 with `tracked vendored path must be untracked before installation; existing bytes were preserved`; the pre-run SHA-256 `15432c0b…` is unchanged after the run, no `scripts/` subtree appeared, and `.git/info/exclude` was not written. The original run mutated the hash and left an untracked subtree. Test: `TestRemediationLast30DaysRefusesTrackedTargetBeforeCloneOrMutation`. One narrow sub-path remains, recorded under Minor below.

**3. Failed commit rollback overwrites another writer's tracked replacement — CLOSED.**
`captureTrackedPublications` (`git_transaction.go:487-500`) records each touched target's inode and SHA-256 immediately after the mutation, and `trackedPublicationStillOwned` (`git_transaction.go:887-897`) gates the `git restore` so a changed target is unstaged and preserved instead of overwritten. Any appended rollback error flips the status to `incomplete` (`git_transaction.go:880-881`). Probe: tracked `do-work/notes.md` at `fb1d689e…`, a `pre-commit` hook that replaces the file and exits 1, `do-work-note --commit "ours"`. Result: outcome `committed_state_risk`, exit 4, `rollback: incomplete`, `rollback error: tracked target changed after publication; preserved replacement: do-work/notes.md`, final file contents `SECOND WRITER REPLACEMENT` (`1b096838…`), index unstaged (` M`). The original reported `rollback: succeeded` and restored the HEAD bytes over the replacement. Test: `TestRemediationTrackedRollbackPreservesConcurrentReplacement`.

**4. Media cancellation does not own the process tree or the Git phase — CLOSED.**
`runOwnedProcess` (`report_image_process.go:38-66`) now TERMs the group, waits for both leader exit and group liveness through the grace period, then KILLs the group and reaps; `handleReportImage` threads one `imageSignalContext` through generation and `runTransactionContext` (`report_image.go:58-85`); `configureCancellableProcessGroup` (`git_transaction.go:1078-1110`) puts Git and its hooks in an owned cancellable group with TERM/grace/KILL and a `WaitDelay`; and `rollbackFailure` runs under `context.WithoutCancel` so cleanup survives cancellation. Probe A: fake `imagegen` with a TERM-obedient leader and a TERM-deaf `sleep 60` descendant; after TERM to the CLI the process exits 143 and the recorded descendant PID is dead (`kill -0` fails). The original left it alive. Probe B: fake backend plus a TERM-ignoring blocking `pre-commit` hook, `--commit`; TERM to the CLI kills the hook, the CLI exits 143, no commit exists, and `rollback: succeeded — removed created target out/img.png`. Probe C (the "if that commit completes" clause): a blocking `post-commit` hook so the commit object is already written; the result does not present as a clean interruption — it reports `committed_state_risk` / `GIT-ROLLBACK-INCOMPLETE` with `the rollback did not complete, so the worktree needs a person before any retry`. Tests: `TestRemediationCancellationReachesMediaGitCommitAndRollback`, `TestRemediationLeaderExitStillKillsTermDeafDescendant`, `TestRemediationOwnedProcessDoesNotLaunchAfterCancellation`.

**5. Batch media publication and status accounting — CLOSED.**
Failed workers are retained in `failedNames` and emitted as both exact text and typed findings (`report_image.go:212-215`), and `generated/` is claimed with `rootedMkdirExclusive` whose failure returns `REFUSING: generated/ appeared before publication` (`report_image.go:171-174`), replacing the check-then-`os.Rename` gap. Children publish through the pinned owned-directory handle, and the deferred cleanup removes the directory only while it is still the same inode (`report_image.go:176-186`). Probe: fake `imagegen` exiting 1, two-item batch. stdout is `MISSING: one.png → fall back to SVG/Mermaid for that section` and the same line for `two.png`; no `generated/` directory was created. The original produced zero stdout and stderr bytes.

**6. Portfolio snapshot-first recovery — RESIDUAL.**
The snapshot half is closed: `handlePortfolio` (`portfolio.go:58-107`) publishes the immutable snapshot in its own completed transaction before attempting canonical publication, so a canonical failure no longer rolls the snapshot back. Probe: `--with-snapshot` against a canonical path that is a directory exits 2 and `pf/snap.md` now exists containing the source bytes; the original left no snapshot. `--format json` shows both the canonical refusal (`GIT-INVALID-OPTIONS`) and `PORTFOLIO-SNAPSHOT-RETAINED`. Test: `TestRemediationPortfolioRetainsSnapshotWhenCanonicalIsDirectory`.

What is not closed is the second half of the finding's demand — reporting the canonical refusal. `portfolio.go:99-105` assigns `canonicalResult.ExactTextOutput` on the failure branch as well as the success branch, and `resultmodel.go:384-385` returns `ExactTextOutput` in place of the entire rendered block. In the default text renderer the user therefore sees exactly one line, `pf/snap.md`, and exit 2, with no statement of what failed or that a snapshot was deliberately retained. Every other command in the package guards this assignment on `OutcomeSuccess` (`note.go:55`, `architecture.go:172`, `report_image.go:89`, `report_image.go:210`, `portfolio.go:143`); this new branch is the only one that does not. — `impact-user-visible`

**7. Architecture scan misreads prior reports for absolute paths; failed bundles not durably occupied — CLOSED.**
`architectureScan` now carries `filesystemPath` and `displayPath` separately (`architecture.go:56-96`), so the watermark is read from the resolved filesystem path while the reported path keeps the caller's form. Probe: fixture with a prior report watermarked `fade70d`; the relative scan and the absolute scan both return `prior_hash=fade70d` and `prior_hash_resolves=yes`, differing only in the `prior_report` / `report_candidate` display prefix. The original returned `prior_hash=unreadable` for the absolute form. Bundle occupation: the directory is claimed by `rootedMkdirExclusive` outside the transaction and is not registered as a created directory (`architecture.go:155-171`), so a rollback removes only `index.html` and the claimed directory survives; `TestRemediationArchitectureFailedClaimRemainsOccupied` forces a publication failure through the `architectureAfterClaim` hook and asserts the bundle is still a directory afterward. (The new claim loop introduced a separate defect — see new finding N1.)

**8. Mutation flags consume positional data, no `--` escape — CLOSED.**
`parseMutationFlags` (`commands.go:49-75`) recognizes options only in a leading region that ends at the first non-option argument or at `--`. Probe: `do-work-note literal --commit` in a clean repo stores `- [2026-09-01] literal --commit` and the commit count is unchanged at 1; the original created a commit and stored only `literal`. `do-work-note -- --dry-run` reaches the dirty-target guard, confirming `--dry-run` was treated as data rather than as a mode. Test: `TestRemediationMutationOptionsHaveLeadingRegionAndDoubleDash`. The stale usage strings this change leaves behind are reported as new finding N3.

**9. Audit malformed-option parity and missing characterization proof — CLOSED.**
`parseAuditOptions` now tracks option presence separately from sentinel value. Probe: I built the standalone oracle from `skills/do-work-toolbox/tools/audit-metrics` and ran a status differential over this repository. All four originally failing cases now match — `inventory --since-window '12 months'`, `churn --watch-lines -1`, `inventory --watch-files -1`, `folders --watch-lines -1` all exit 2 in both tools. Six further malformed and wrong-subcommand probes also match, as do an unknown flag and an unknown subcommand. Valid output is byte-identical for all four modes (`inventory` 1610 bytes, `folders` 596, `churn` 518, `hotspots` 863) and for `churn --since-window "6 months"`. Durable proof is committed: `TestRemediationAuditStandaloneDifferentialAndTypedOrder` builds the oracle and compares exact bytes for all four modes over a fixture carrying a binary file, an unterminated file, and a rename, and asserts rename-normalized typed ordering; `TestRemediationAuditExplicitDefaultsRemainWrongSubcommandFlags` pins the four malformed cases. The test skips on Go 1.25 only because the retained oracle module declares Go 1.26, which the handback discloses; it ran here on Go 1.26.1.

### Original Minor findings

- **`architectureScan` ignores `os.ReadDir` errors — CLOSED.** `architecture.go:58-61` returns `ARCHITECTURE-PREFLIGHT-FAILED` with `reports directory is unreadable: …`. Probe: `--scan nonexistent-dir` exits 1 with that finding instead of a clean first-run result.
- **Unchecked P-A-U boxes and unsupported adversarial claims — PARTIAL.** The claims are now supported: thirteen named tests exist, and each pins a real counterexample rather than restating the handback. The REQ's three P-A-U checkboxes are still `- [ ]` in `do-work/working/REQ-418-migrate-toolbox-absorb-audit-metrics.md:26-28`. The handback assigns that write to the orchestrator under the no-lifecycle instruction, which is a defensible split, but at review time the boxes are unchecked.

### New findings

**N1. `architecture-report-preflight --publish` spins forever when the bundle claim fails for any non-existence reason. — `impact-critical`**
The claim loop at `architecture.go:156-164` increments `sequence` on every error returned by `rootedMkdirExclusive`, without distinguishing "already exists" from a persistent failure, and has no iteration bound. Any durable error — a non-writable parent, a read-only filesystem, a parent removed concurrently — retries a new suffix forever. Probe: fixture with `chmod 555 reports`, then `architecture-report-preflight --publish draft.html reports/bundle`. The process was still running after 15 seconds at 98.6% CPU and had to be killed with `SIGKILL`. This is a regression: the pre-remediation binary built from `a7c975c5` exits 3 on the identical fixture with `GIT-MUTATION-FAILED … mkdir …/reports/bundle: permission denied` and `rollback: succeeded`. No committed test exercises the claim-loop error path.

**N2. `architecture-report-preflight --commit` is dead in every argument position. — `impact-user-visible`**
`architecture.go:144` runs the new preflight as `runTransaction(…, true /* dryRun */, commit, …)`, passing the caller's `commit` alongside a hardcoded `dryRun=true`, and the transaction rejects that pair. Probe: `architecture-report-preflight --commit --publish draft.html reports/bundle` exits 2 with `GIT-INVALID-OPTIONS: --dry-run and --commit cannot be combined` — an error naming an option the user never passed — and creates no commit. The trailing form fails too, with a usage error, per the finding-8 option-region change. The pre-remediation binary publishes the bundle and creates the commit on the same fixture (exit 0, commit count 1→2). The defect is specific to this command: `publish-portfolio-summary --commit` and `do-work-note --commit` both still commit correctly on the same binary. `grep -rn "commit" internal/toolboxcommands/architecture_test.go` returns only fixture setup, so nothing covers the mode.

**N3. Restatement sweep: four usage strings still advertise the trailing flag form the option-region fix stopped accepting. — `impact-user-visible`**
The finding-8 fix redefined where `--dry-run` and `--commit` are recognized. `note.go:23` states the new leading form correctly. Four others still state the old trailing form: `architecture.go:38`, `report_image.go:38`, `report_image.go:103`, `portfolio.go:20` (and `last30days.go:33` for `[--dry-run]`). Following the printed usage now fails or silently changes meaning. Probe: `publish-portfolio-summary --canonical-only src.md can.md --dry-run` exits 2 with `portfolio mode has wrong argument count`; `install-last30days check . --dry-run` consumes `--dry-run` as the `[source-repository]` positional and proceeds as a non-dry-run check. The retained shell scripts define no such flags at all, so these usage strings are the only contract statement for them and are the sole thing a reader or agent has to go on.

### New Minor findings

- `ensureLast30DaysExclude` (`last30days.go:334-359`) is dead code. The remediation replaced it with `prepareLast30DaysExclude`; `grep -rn ensureLast30DaysExclude` finds only the definition. `go vet` does not catch unused functions.
- `publishLast30Days` discards `restore()`'s error at `last30days.go:192` and `:196`, while `defer os.RemoveAll(backup)` at `:163` then deletes the directory holding the only copy of the previous tree, and the message at `:197` claims `previous tree restored` on a path where it was not. Reaching it needs a copy or verification failure on a stage already verified complete plus a concurrent replacement of the target, so I could not reproduce it; graded from code alone.
- `TestRemediationAuditComparatorRejectsStatusTextAndOrderMutations` (`audit_metrics_test.go`) asserts against a two-line `compare` closure defined inside the test body. It exercises no product code and cannot fail for any change to the CLI, so it does not support the handback's "mutation-check" claim. The real differential coverage in the sibling test does.

## Requirements checklist

- [x] Finding 1 — rooted, no-follow ancestor resolution applied once to the whole publisher family
- [x] Finding 2 — eligibility and dirty-target checks precede mutation; identity-bound exclude undo
- [x] Finding 3 — rollback compare-and-swaps against recorded publication identity and preserves replacements
- [x] Finding 4 — one invocation context through generation, commit, and rollback; group liveness through grace/KILL/reap
- [x] Finding 5 — exclusive directory claim and per-item retained failure diagnostics
- [~] Finding 6 — snapshot published and retained before canonical failure; the canonical refusal is not reported in text mode
- [x] Finding 7 — display and filesystem paths separated in scan; failed bundle claim stays occupied
- [x] Finding 8 — options parsed only in a leading region; `--` honored
- [x] Finding 9 — option presence tracked separately; oracle differential and malformed matrix committed and passing
- [x] Minor — `ReadDir` errors surfaced as incomplete evidence
- [~] Minor — adversarial claims now backed by tests; REQ P-A-U checkboxes still unchecked
- [x] Scope — 21 remediation paths, all under `do-work-cli`, within the declared 32-path ceiling
- [ ] No new user-visible regressions — N1 and N2 are reproduced regressions against `a7c975c5`

## Acceptance testing

**Result: Partial**

Required gates, run from the stated directories and judged on direct exit codes, never piped:

- `go test -count=1 ./...` in `skills/do-work/tools/do-work-cli` — exit 0
- `go vet ./...` in `skills/do-work/tools/do-work-cli` — exit 0
- `go test -race ./internal/toolboxcommands ./internal/gittransaction` — exit 0 (`toolboxcommands` 5.591s, `gittransaction` 5.917s)
- `bash _dev/tests/maintainer-verify.sh` from repo root — exit 0; the optional browser lane made its declared no-browser skip
- `go test -count=1 -v -run TestRemediation ./internal/toolboxcommands ./internal/gittransaction` — all 13 named tests RUN and PASS on Go 1.26.1; none skipped
- `GOOS=windows GOARCH=amd64 go build ./...` and `go vet ./internal/toolboxcommands` — exit 0 both

Hands-on probes against the binary built from HEAD, each described under its finding above: linked-ancestor escape (architecture, note, portfolio ×2) — refused, external directory untouched; dirty tracked last30days target — refused, hash unchanged, no subtree, exclude unwritten; concurrent tracked replacement during a failing commit hook — replacement survives, rollback honestly incomplete; TERM-deaf media descendant — dead after exit 143; TERM-resistant commit hook — killed, no commit, target restored; post-commit cancellation — reported as committed state risk, not as a clean interruption; all-failed batch — both MISSING lines emitted, no `generated/`; canonical-directory snapshot — retained with source bytes; absolute architecture scan — watermark now read; literal `--commit` note text — stored as data, no commit; audit malformed differential — 10/10 status matches against the standalone oracle; audit valid differential — byte-identical Markdown for all four modes.

Partial rather than Pass because two probes of the same command found reproduced regressions: `architecture-report-preflight --publish` hangs at 98.6% CPU indefinitely on a non-writable reports directory where the pre-remediation build exits 3 cleanly, and `--commit` no longer publishes or commits in any argument position where the pre-remediation build did.

Scope and hygiene: `git diff --name-only a7c975c5..a43b2587` is 21 paths, all under `skills/do-work/tools/do-work-cli/`; the full builder range `b713fa8b..a43b2587` is 32 paths, matching the declared ceiling; `git diff --check b713fa8b..a43b2587` exits 0; those 32 paths are byte-identical between `a43b2587` and HEAD, so reviewing at HEAD is equivalent to reviewing the stated merge `82534d36`.

## Suggested additional testing

- A claim-loop test that forces a persistent non-`EEXIST` error (read-only parent) and asserts the command returns a finding rather than looping; bound the loop by attempts and discriminate on `errors.Is(err, fs.ErrExist)`.
- An `architecture-report-preflight --publish … --commit` end-to-end test asserting the bundle is published and one commit is created — the absence of any `--commit` assertion in `architecture_test.go` is what let N2 through.
- A renderer-level assertion that `ExactTextOutput` is never set on a non-success outcome, which would have caught the portfolio failure-branch suppression and would pin the convention the other five call sites already follow.
- A usage-string contract test that parses each printed `Usage:` line with `parseMutationFlags` and asserts the documented invocation is accepted.
- A last30days test that fails the publication copy while the target has been replaced, asserting the backup is not deleted and the reported message does not claim restoration.

## Scores

**Overall: 69%**

| Dimension | Score | Notes |
|---|---:|---|
| Requirements | 80% | Eight of nine findings fully closed and independently reproduced; finding 6 half-closed; the REQ's "preserve existing user-visible behavior" clause is violated by two architecture regressions. |
| Code Quality | 65% | The confinement, identity-bound rollback, and single-context cancellation designs are well factored and applied once across the family. Against that: an unbounded retry loop with no error discrimination, a failure branch that suppresses its own diagnostics, four stale usage strings, and one dead function. |
| Test Adequacy | 70% | Thirteen named tests, each pinning a real counterexample, all green including under `-race`, plus a genuine standalone-oracle differential. Both new defects sit in uncovered paths of the most heavily rewritten file, and one committed test is decorative. |
| Scope | 100% | Exactly 21 paths, all inside `do-work-cli`, within the declared 32-path ceiling; `git diff --check` clean; builder paths byte-equivalent at the reviewed integration. |
| Risk | Low | Every reproduced critical-class behavior from the initial review — outside-root writes, concurrent replacement loss, dirty-target loss, surviving media descendants — is gone and was re-probed. The remaining exposure is one command's publish path hanging or refusing: recoverable, no data loss, no security surface. The unbounded loop still carries `impact-critical` on its own line for routing, which is a different axis from this one. |
| Acceptance | Partial | All required gates exit 0 and the nine remediated contracts verify hands-on; two reproduced regressions in `architecture-report-preflight`. |

Computation: (80 + 65 + 70 + 100) / 4 = 78.75%, less the 10-point Acceptance = Partial penalty, gives 68.75%, reported as 69%. Verdict mapping: Acceptance Partial or Overall 50–74% maps to **Approve with follow-ups**; both conditions hold, so the two paths agree.

## Review record

**Original findings:** 9 Important — 8 CLOSED, 1 RESIDUAL (finding 6). 2 Minor — 1 CLOSED, 1 PARTIAL.

**Residual and new Important findings, each with its recorded impact token:**

- Finding 6 residual — portfolio's snapshot-first failure branch sets `ExactTextOutput` on a non-success outcome, so the canonical refusal and the snapshot-retained warning are invisible in the default text renderer (`portfolio.go:99-105`, `resultmodel.go:384-385`) — `impact-user-visible`
- N1 — unbounded CPU-spinning claim loop in `architecture-report-preflight --publish` on any persistent claim error (`architecture.go:156-164`); reproduced regression against `a7c975c5` — `impact-critical`
- N2 — `architecture-report-preflight --commit` rejected in every argument position because the preflight passes `dryRun=true` alongside the caller's `commit` (`architecture.go:144`); reproduced regression against `a7c975c5` — `impact-user-visible`
- N3 — four usage strings restate the superseded trailing-flag contract (`architecture.go:38`, `report_image.go:38`, `report_image.go:103`, `portfolio.go:20`; also `last30days.go:33`) — `impact-user-visible`

**Minor findings:** 4 new (dead `ensureLast30DaysExclude`; swallowed `restore()` error with a false "previous tree restored" message and deferred backup deletion; decorative audit comparator test; unchecked REQ P-A-U boxes) — report only.

**Acceptance:** Partial — every required gate exits 0 and all nine remediated contracts verify hands-on, but two reproduced regressions remain in `architecture-report-preflight`.

**Suggested testing:** 5 items.

**Termination:** the remediation allowance is exhausted and one residual plus three new Important findings stand, so REQ-418 terminates `completed-with-issues`; the four tokens above route the follow-ups, with N1's `impact-critical` piercing the depth stop.

*Independently re-reviewed against HEAD `d924b2fc`, whose product paths are byte-identical to the stated merge `82534d36`. No tracked file was edited by this review; every fixture was created outside the repository working tree.*
