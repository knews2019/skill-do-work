# Review: REQ-418

**Request changes** — all seven toolbox entry points are registered and the ordinary audit output is compatible, but the implementation does not preserve the required confinement, publication, rollback, cancellation, and diagnostic contracts. Several synthetic fixtures reproduced data loss or work outside the repository.

Route C | builder range `b713fa8b..a7c975c5` | reviewed integration `94560fde`

## What's built

- Seven shared-CLI handlers cover notes, architecture scan/publication, single and batch report images, portfolio publication, last30days, and four audit-metrics modes.
- The result model carries exact compatibility text, typed audit data, and conventional media interruption statuses.
- The standalone audit module and all retained toolbox scripts remain intact for later shim/removal work.
- The builder changed exactly the frozen 30 product files, and those 30 files are byte-equivalent at the reviewed integration commit.

## Decision / risk

Do not accept or release REQ-418 yet. The highest-risk counterexamples overwrite a dirty tracked last30days tree and another writer's tracked replacement during rollback; linked ancestors also let nominally repository-confined commands create objects outside the repository. Media cancellation can return while a TERM-resistant descendant remains alive. These are correctness and ownership-boundary failures, not test-only concerns.

## Findings

### Important

1. **Repository path validation is lexical, so linked ancestors escape the declared output root.** `repositoryPath` accepts a path after `filepath.Rel` without resolving or rooted-opening its ancestors (`portfolio.go:108-120`), while architecture publication performs `MkdirAll`, `Mkdir`, and file creation through ordinary absolute paths (`architecture.go:116-156`). The same helper feeds portfolio and report-image targets, and the fixed note target is also reached through ordinary path traversal. In a real-Git fixture where `reports` linked to an external directory, `architecture-report-preflight --publish draft.html reports/run` created the external `run/` directory. The command exited 4 only after rollback removed the file but could not identity-record/remove the linked directory. A user can therefore affect paths outside the repository despite the explicit containment check, and cleanup cannot restore the claimed boundary. Remediate the family once with rooted, no-follow ancestor resolution/publication rather than adding leaf-only checks per command. — `impact-critical`

2. **last30days replaces state before establishing that the target is eligible and does not roll publication back when later checks fail.** Installation calls `publishLast30Days` before `ensureLast30DaysExclude` (`last30days.go:44-68`); publication moves the old tree into a temporary backup, renames the staged tree, and then unconditionally deletes the backup on return (`last30days.go:121-154`). A fixture with a dirty tracked, incomplete `.claude/skills/last30days/SKILL.md` exited 1 at the tracked-path guard, but its pre-run hash changed from `29e5d7...` to the cloned source's `d7b4cd...`, and a new untracked `scripts/` subtree remained. Similar post-publication exclude failures have no restoration path, and target reappearance can make the ignored restoration fail before the deferred backup deletion. Move eligibility and dirty-target checks before mutation and make tree publication plus exclude verification one identity-bound transaction. — `impact-critical`

3. **Failed commit rollback overwrites another writer's tracked replacement.** Toolbox mutators use the shared transaction wrapper (`mutation.go:47-52`), whose tracked rollback asks only whether the target is dirty and then runs `git restore --source=HEAD --staged --worktree` (`git_transaction.go:763-778`); it does not bind the bytes/inode published by this invocation. In a note fixture, a blocking commit hook allowed a second writer to replace the mutated target, then failed. The command reported `rollback: succeeded`, but the final target hash was the original HEAD hash `29e5d7...`, not the second writer's replacement hash `cdf407...`. This is silent concurrent data loss. Rollback must compare-and-swap against the invocation's recorded publication identity/content and preserve/report any replacement. — `impact-critical`

4. **Media cancellation does not own the complete process tree or the later Git phase.** `runOwnedProcess` sends TERM to the group but returns as soon as the direct child exits, without checking whether live group members remain or applying KILL (`report_image_process.go:30-45`). A fixture whose backend parent obeyed TERM while a descendant ignored it returned status 143 with `rollback: succeeded`; the descendant was still alive afterward. Separately, every toolbox transaction is run with `context.Background()` (`mutation.go:47-51`). During a blocking media `git commit`, TERM was consumed by the media signal handler, but both CLI and hook remained alive; only after the hook was terminated did the CLI exit 143. If that commit completes, the command can commit state and still present itself as an interrupted invocation. Use one invocation context through generation, staging, commit, and rollback, and wait for process-group liveness—not only the leader—through grace/KILL/reap. — `impact-critical`

5. **Batch media publication and status accounting do not preserve the retained contract.** Backend errors are discarded (`report_image.go:143-150`), so a mixed or all-failed batch carries no per-item typed status or fallback guidance. An all-failed two-item fixture returned 0 with zero stdout and stderr bytes; the retained script emits `MISSING: <name> → fall back to SVG/Mermaid for that section` for every failed item. Publication also checks `generated/` and then calls `os.Rename` (`report_image.go:156-160`); on Unix, renaming a directory over an empty directory can replace the directory that appeared in that gap, so this is not the promised exclusive no-clobber boundary. Publish via an actually exclusive directory claim/identity check and retain every worker outcome in text/JSON. — `impact-user-visible`

6. **Portfolio does not implement snapshot-first recovery when canonical publication is unsafe or fails.** `runTransaction` preflights every target before the mutation, and target inspection rejects a canonical directory as an invalid regular-file target; therefore the snapshot is never created (`portfolio.go:47-93`, `git_transaction.go:615-630`). A `--with-snapshot` fixture with a canonical directory exited 2 and left no snapshot, whereas the retained script publishes and retains the immutable snapshot before reporting the canonical refusal. More generally, any failure after snapshot creation is passed to the transaction's all-target rollback, which removes the snapshot rather than preserving the named exception. Implement the documented snapshot-first partial-success state outside an all-or-nothing multi-target rollback and test every later canonical failure boundary. — `impact-user-visible`

7. **Architecture scan misreads prior reports for absolute scan paths, and failed bundles are not durably occupied.** Scan records an absolute `prior_report` when its argument is absolute, then prepends the repository root again when reading the watermark (`architecture.go:51-55,76-89`). Against the same prior HTML, a relative scan returned watermark `96b282...` while an absolute scan returned `prior_hash=unreadable`. Publication also registers the newly created bundle with the generic transaction; a subsequent copy/recording failure rolls that directory back, contrary to the requirement that a failed claimed bundle remain occupied and never be reused. Separate display paths from filesystem paths in scan, and preserve the exclusive bundle claim across failed publication while keeping partial HTML hidden. — `impact-user-visible`

8. **Mutation flags consume positional data and provide no `--` escape.** `parseMutationFlags` removes every argument equal to `--dry-run` or `--commit`, regardless of position (`commands.go:49-63`). Legal note text, style briefs, descriptions, filenames, project/source paths, or candidate paths equal to those strings cannot be passed literally. A fixture invoking `do-work-note literal --commit` created a Git commit and stored only `literal`, even though the final token may have been intended as note data. Parse options only in a defined option region and honor `--` so compatibility arguments remain lossless. — `impact-user-visible`

9. **The audit port matches ordinary output but not malformed-option parity, and the required characterization proof was not committed.** The parser uses sentinel values to represent both "absent" and explicitly supplied defaults (`audit_metrics.go:59-115`). Consequently `inventory --since-window '12 months'`, `churn --watch-lines -1`, `inventory --watch-files -1`, and `folders --watch-lines -1` all returned 0 in the new CLI while the standalone tool rejected each with status 2. The new package has only four small audit tests (52 source lines) compared with the standalone module's twelve mature fixture tests (438 lines), and there is no old/new differential harness, JSON-row parity harness, or mutation check promised by the frozen plan. Valid real-repository Markdown was byte-identical for inventory, folders, churn, and hotspots, so the algorithms appear sound on ordinary input; the acceptance failure is malformed behavior plus missing durable proof for rename/copy, excluded-live-source, shallow/unavailable, NUL/binary, threshold, ordering, and JSON cases. Track option presence separately and port/run the complete characterization matrix. — `impact-user-visible`

### Minor

- `architectureScan` ignores `os.ReadDir` errors (`architecture.go:58`), collapsing an unreadable reports directory into a clean first-run result instead of reporting incomplete evidence.
- The REQ's P-A-U checkboxes remain unchecked, and the handback's broad adversarial claims are not supported by the committed focused tests.

## Requirements checklist

- [x] All seven canonical toolbox command names are registered; `queue-kanban` remains separate.
- [~] `do-work-note` normalization, append bytes, dry-run, and ordinary commit work; literal flag-shaped text, confinement, concurrent rollback, and cancellation do not.
- [~] Architecture scan/publication works for ordinary relative paths and numeric suffixes; absolute scan resolution, linked ancestors, and failed-bundle occupation fail.
- [ ] Report/media process handling, exclusive publication, failure diagnostics, and complete cancellation cleanup — reproduced counterexamples remain.
- [~] Portfolio canonical/snapshot bytes are separate and ordinary suffixing works; snapshot-first recovery is absent.
- [~] last30days complete-payload checking, target Python 3.12+ probing, non-Git reporting, and ordinary installation work; dirty-target preservation and transaction rollback fail.
- [~] audit inventory/folders/churn/hotspots emit byte-compatible Markdown on ordinary fixtures and typed JSON; malformed flag parity and the required characterization matrix fail.
- [x] Python remains a target prerequisite only; no Python/jq implementation branch was introduced.
- [x] Windows toolbox cross-compilation succeeds and process launch fails closed there.
- [x] The standalone audit source and retained scripts remain unchanged; retirement is deferred to REQ-420.
- [x] Exact frozen 30-file builder scope and integration equivalence are satisfied.

## Acceptance testing

**Result: Fail**

- `bash _dev/tests/maintainer-verify.sh` — passed at current main; the optional browser lane made its declared no-browser skip.
- `go test -count=1 ./...`, `go vet ./...`, and `go test -race ./internal/toolboxcommands` in `do-work-cli` — passed.
- Standalone audit `go test -count=1 ./...` and `go vet ./...` — passed.
- `bash _dev/tests/do-work-cli-go125-compatibility.sh` — passed.
- `GOOS=windows GOARCH=amd64 go test -c ./internal/toolboxcommands` — passed.
- Valid-output differential — exact Markdown matched for all four audit modes with representative flags.
- Malformed audit differential — failed: four wrong-subcommand explicit-default/sentinel cases were accepted by the new CLI and rejected by the standalone tool.
- Adversarial fixtures — failed: linked-ancestor escape, dirty tracked last30days replacement, concurrent tracked replacement loss, TERM-resistant descendant survival, commit-phase cancellation, absolute architecture scan, canonical-directory snapshot retention, flag-shaped data, and batch failure diagnostics.
- Scope/hygiene — `git diff --check b713fa8b..a7c975c5` passed; exactly 30 paths changed, all belong to the frozen allowlist, and those paths have zero delta from `a7c975c5` to `94560fde`.

Passing repository gates do not change the verdict because they exercise the retained shell behavior and ordinary Go paths, not the failed shared-CLI boundaries above.

## Suggested additional testing

- Add rooted linked-ancestor fixtures for every file and directory publisher, including parent replacement between validation and publication.
- Add identity-bound concurrent-writer tests at mutation, stage, commit-hook failure, rollback quarantine, and final verification; add cancellation at Git add, commit, and post-commit verification.
- Port the retained media liveness matrix verbatim: leader-exits/descendant-lives, zombie-only group, TERM-deaf escalation, unrelated-group survival, mixed/all-failed statuses, and empty/nonempty final-boundary collisions.
- Port every standalone audit test and add executable old/new status+Markdown differential fixtures plus JSON row/order and malformed-flag mutation checks.
- Add last30days failures after backup, rename, exclude write, and verification, with dirty tracked/untracked targets and another writer's reappearing target.

## Scores

**Overall: 35%** (Acceptance Fail)

| Dimension | Score | Notes |
|---|---:|---|
| Requirements | 45% | Entry points and ordinary paths exist; multiple named safety and parity contracts fail. |
| Code Quality | 40% | Clear package split, but ownership boundaries rely on lexical checks, TOCTOU publication, and non-identity rollback. |
| Test Adequacy | 25% | Existing suites are green, but most of the frozen adversarial/TDD matrix is absent and direct probes fail. |
| Scope | 100% | Exact frozen 30-file range; reviewed integration is byte-equivalent on those paths. |
| Risk | Critical | Reproduced outside-root writes, concurrent replacement loss, dirty-target loss, and surviving media descendants. |
| Acceptance | Fail | Important requirements have executable counterexamples. |

## Review record

**Important findings:** 9, each carrying an impact token above. **Minor findings:** 2. **Follow-ups:** consolidate remediation by the nine root areas above; no tracked product file was edited by this review.

*Reviewed independently against `94560fde`; the later main commit changes only an unrelated generated report and does not alter the 30 reviewed product paths.*
