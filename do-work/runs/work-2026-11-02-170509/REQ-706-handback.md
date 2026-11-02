# Hand-back: REQ-706 (ai-report tools read report bundles inside month folders)

- Branch: `worktree-agent-REQ-706-ai-report-tools-read-month-folders`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-706-ai-report-tools-read-month-folders`
- Base commit: `782f083a` (main HEAD at dispatch; the brief's line numbers were at `d231e2fa`, and the six files are identical at both)
- Commits: `c004cba0` `[REQ-706] ai-report tools read report bundles inside month folders` (one commit)
- Worktree status after commit: clean.

## File Manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modified): new `reportMonthPattern` regexp and `reportMonthFolder(name string) bool` helper. `reportIndexWalk` reads the non-dot child directories of a month folder as bundles, one level only. `readReportBundle` now takes the bundle's path relative to `ai-reports/` (`<folder>` or `<month folder>/<folder>`) and reads the date, kind, linked ids and default title from `path.Base` of it. Both doc comments are rewritten.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified): `secondOpinionBundle` accepts a parent of `ai-reports` or `ai-reports/<month folder>`. New refusal text: `bundle must be a physical report bundle directly under ai-reports or under a month folder`. Both `validateNoLinkedAncestors` calls and the dot-name refusal are kept.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modified): new `TestAIReportIndexReadsBundlesInsideMonthFolders`.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified): new `TestSecondOpinionBundleInsideMonthFolder`.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): the catalog sentence (defines "month folder" at its first use), the revise step 1 target rule plus the clause that a new revision is still written directly under `ai-reports/`, and the second-opinion bundle rule.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): one new layout paragraph at the top of `## Report Catalog`. The catalog sentence now treats month folders as containers. The revise sentence says the new report is written directly under `ai-reports/`. The second-opinion paragraph accepts a report inside a month folder. No organize command is described.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** I read the brief, the REQ (Exploration D-01 to D-06 and Scope), the rules (general, coding-guardrails, shared-principles, communication-style, testing, backend), the primes (action-files, releases, do-work-cli § Traps), `lessons-releases.md`, and the `catalog-tie-order` bullet. The approach follows the brief: one helper with one regexp, the walk descends one level into month folders, and `readReportBundle` takes a relative path and uses `path.Base` for every name-derived field. I checked that sort, supersedes resolution, `reportPageLink` and `reportIndexFind` need no change, because they key on `path.Base` or trim the `ai-reports/` prefix. In second-opinion, one predicate change covers every entry route (plain bundle, `latest`, page target, `--finalize`), because all of them call `secondOpinionBundle` (the `guard-on-one-entry-path` trap). The git-status outside check already uses `--untracked-files=all` with a prefix match, so nested bundle paths work there unchanged. I wrote the RED tests before any production edit.
- [x] **[APPLY]:** The code matches the plan, and the scope is exactly the six files in `## Scope`. No other files were touched. Nothing under `do-work/` was staged except that this hand-back was written (not staged).
- [x] **[UNIFY]:** I reviewed each changed file in `git diff 782f083a`, and there are no debug artifacts. `git diff 782f083a --stat`:

```
 skills/do-work-toolbox/actions/ai-report.md        |  6 +-
 skills/do-work-toolbox/docs/ai-report-guide.md     |  8 ++-
 .../internal/toolboxcommands/report_index.go       | 39 +++++++++++--
 .../internal/toolboxcommands/report_index_test.go  | 65 ++++++++++++++++++++++
 .../toolboxcommands/report_second_opinion.go       |  6 +-
 .../toolboxcommands/report_second_opinion_test.go  | 37 ++++++++++++
 6 files changed, 147 insertions(+), 14 deletions(-)
```

| Check (from worktree root) | Exit | Wall |
|---|---|---|
| `gofmt -l skills/do-work/tools/do-work-cli/internal` (printed nothing) | 0 | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands` | 0 | <1 s |
| GREEN: `go test ... -run '^(TestAIReportIndexReadsBundlesInsideMonthFolders\|TestSecondOpinionBundleInsideMonthFolder)$'` | 0 | 1 s |
| Selection: `go test ... -run 'SecondOpinion\|AIReportIndex'` | 0 | 15 s |
| `bash .../REQ-706-probe.sh` (before and after the commit) | 0 | <1 s |
| `bash _dev/tests/contract-regressions.sh` | 0 | 22 s |
| `bash _dev/tests/shipped-package-reference-contract.sh` | 0 | 2 s |
| `git diff --check` | 0 | <1 s |

What I checked in each file:
- `report_index.go`: for a flat bundle, `relative == folder`, and `path.Base` of a name without a slash is the name itself. So Path, Entry, Title, Date, Kind and LinkedIDs are byte-identical, and the existing catalog tests pass unchanged. Sort and supersedes resolution are untouched.
- `report_second_opinion.go`: when `filepath.Rel` errors, `relative` is empty and the predicate still refuses. A month folder is matched only when its parent is `ai-reports`. The existing `filepath.Join(bundle, "nested")` refusal still holds, because `2026-10-15_1200-example` is not a month-folder name.
- The two test files: each new test has a `// Pins:` comment and reuses the existing helpers (`runReportIndex`, `readReportCatalog`, `secondOpinionFixture`, `secondOpinionInvoke`, `secondOpinionRows`, `toolboxTestGit`).
- The two prose files: the probe's literal checks pass, and no sentence says a bundle must be a direct child of `ai-reports/`. I re-ran the sweep `git grep -n -i 'direct child\|directly under ...\|every directory under\|every bundle under' -- skills/`. Its only remaining hits are unrelated (maintainability-audit folders, `audit_inventory.go`) or the historical CHANGELOG entry for 0.305.x, which is release history and not edited.

Mutation check: I made the walk accept any non-dot child (I dropped `child.IsDir()`). `TestAIReportIndexReadsBundlesInsideMonthFolders` then failed on the loose file. I restored the file from a scratch copy and the test passed again.

## Proof Record

**RED** (both tests written first, no production edit; exit 1):

```
=== RUN   TestAIReportIndexReadsBundlesInsideMonthFolders/2026-09_1-report
    report_index_test.go:286: catalog paths=[ai-reports/2026-11-01_0900_b ai-reports/2026-09_1-report], want [ai-reports/2026-11-01_0900_b ai-reports/2026-09_1-report/2026-09-27_2200_a]
=== RUN   TestAIReportIndexReadsBundlesInsideMonthFolders/2026-09
    report_index_test.go:286: catalog paths=[ai-reports/2026-11-01_0900_b ai-reports/2026-09], want [ai-reports/2026-11-01_0900_b ai-reports/2026-09/2026-09-27_2200_a]
--- FAIL: TestAIReportIndexReadsBundlesInsideMonthFolders (0.10s)
=== RUN   TestSecondOpinionBundleInsideMonthFolder
    report_second_opinion_test.go:380: nested bundle refused: {... Outcome:failure ... Findings:[{Code:TOOLBOX-USAGE ... Evidence:[bundle must be a physical direct child of ai-reports] ...}] ...}
--- FAIL: TestSecondOpinionBundleInsideMonthFolder (0.09s)
FAIL
```

Both failures are the expected ones: the month folder is listed as a bundle and `a` is missing, and second-opinion refuses the nested bundle as not a direct child.

**GREEN** (after the production edit; exit 0, no skip):

```
--- PASS: TestAIReportIndexReadsBundlesInsideMonthFolders (0.08s)
    --- PASS: TestAIReportIndexReadsBundlesInsideMonthFolders/2026-09_1-report (0.07s)
    --- PASS: TestAIReportIndexReadsBundlesInsideMonthFolders/2026-09 (0.01s)
--- PASS: TestSecondOpinionBundleInsideMonthFolder (0.11s)
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands
```

**Package selection** `-run 'SecondOpinion|AIReportIndex'`: 36 top-level PASS, 0 FAIL, 0 SKIP, exit 0. This includes the four existing `TestAIReportIndex*` catalog tests, unchanged, which are the byte-identical proof for a repository with no month folders.

What the new tests cover:
- Catalog: there are exactly two bundles, `b` first, with paths `ai-reports/2026-11-01_0900_b` and `ai-reports/<month>/2026-09-27_2200_a`, for both `2026-09_1-report` and plain `2026-09`. No bundle is named after the month folder. A loose `loose-notes.md` inside the month folder is skipped. A symlinked month folder `2026-08_1-report` adds no bundle. `index.html` links `<month>/2026-09-27_2200_a/index.html`. `--find 2026-09-27_2200_a` lists the nested bundle. A flat `2026-11-02_1000_a-rev1` whose `ai-report-supersedes` names the nested folder resolves both `supersedes` and `superseded_by`.
- Second-opinion: (a) `--prepare` on `ai-reports/2026-09_3-reports/2026-09-27_2200_a` returns state `prepared`. (b) `latest` resolves the nested bundle over an older flat one. (c) The same bundle reached through a symlinked month folder is refused.

## Decisions

- **D-07** (DECIDE & STATE): if a month folder cannot be read, the whole walk fails with `AI-REPORT-INDEX-FAILED`, the same as an unreadable `ai-reports/`. It is not skipped. Skipping would silently drop a month of reports from the catalog and could make `latest` pick the wrong bundle.
- **D-08** (DECIDE & STATE): the helper is `reportMonthFolder(name string) bool`, the brief's example name, so REQ-707 (the `ai-report organize` command) can reuse it. The regexp is `reportMonthPattern`, not `reportMonthFolderPattern`, so that gofmt does not re-align every line of the existing `var` block.
- **D-09** (DECIDE & STATE): the second-opinion test adds one older flat bundle (`2026-08-30_0900_older`). Without it, the `latest` case would only prove "the only bundle", not "the newest across flat and nested". The test cost is two lines.
- **D-10** (DECIDE & STATE): in the guide, the layout paragraph sits at the top of `## Report Catalog`, where month folders are first read, not in `## What It Produces`. It also names the plain `YYYY-MM` form. In `ai-report.md`, "month folder" is defined in a parenthetical at its first use (the catalog sentence). The new-revision clause (D-05) sits in revise step 1, as the brief placed it.

## Discovered Tasks

- impact-negligible: supersedes resolution in `reportIndexWalk` keys on the bare folder name (`byFolder[path.Base(bundle.Path)]`). Folder names were unique under `ai-reports/` before month folders existed. Now the same bundle name can appear twice: in two month folders (`2026-09/` and `2026-09_3-reports/`), or flat and nested after a copy instead of a move. In that case an `ai-report-supersedes` meta resolves to whichever duplicate comes last in catalog order. REQ-707 (organize) should refuse to create a duplicate bundle name. → report only

## Lessons Read

- `_dev/primes/lessons-releases.md` (whole file). Its `canonical-link-outlives-its-target` family does not apply here, because nothing moves. It does apply to REQ-707.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: only the `[family: catalog-tie-order]` bullet (REQ-656). The sort is unchanged, so nested paths tie-break by path ascending like flat ones. Revise prose still follows `superseded_by` forward.
- Prime traps applied: `guard-on-one-entry-path` (one predicate in `secondOpinionBundle` covers every route) and `alternate-writer-contract-drift` (I swept every restatement of the direct-child rule under `skills/`).

## Anti-bloat Check

- Added beyond the brief's named items: the regexp variable `reportMonthPattern`. The brief asked for "one regexp" without a name, and this is that regexp.
- I added no new flags, options, files, config or environment variables. The fixture-writing closure `writeReportFiles` is local to the new catalog test, not a package-level helper. It exists because the test writes files twice (the base fixture, then the `-rev1` bundle).
- The production change is 39 + 6 diff lines in two Go files.

## Proposed CHANGELOG Entry

**Report Tools Find Reports Filed in Month Folders**

A project with many reports can file each finished month in its own folder, and the ai-report catalog, find, revise and second-opinion still find every report. Before this change, the catalog listed such a folder as one broken report, and second-opinion refused every report inside it.

- `ai-report index` and `find` read a month folder under `ai-reports/` (`YYYY-MM`, `YYYY-MM_1-report` or `YYYY-MM_<N>-reports`) as a container. They list each report inside it with its full path, and they never list the folder as a report. A repository without month folders gets a byte-identical catalog.
- A revision chain resolves across a filed report and a newer one directly under `ai-reports/`, and new revisions are still written directly under `ai-reports/`.
- `ai-report second-opinion` accepts a report inside a month folder, and `latest` can pick it. A symlinked month folder is still refused.
- The ai-report action and guide describe the layout: the current month's reports sit directly under `ai-reports/`, and finished months may sit in month folders.

## Proposed Lesson Bullet

None. Nothing failed during the build, and the one open risk (duplicate bundle names) is recorded above as a Discovered Task for REQ-707.

## Addendum (review fixes)

- Commit: `a7a2fe6b` `[REQ-706] review fixes: architecture pre-flight reads month folders, month folder never a bundle`, on the same branch and in the same worktree. The worktree is clean afterwards.
- Branch history: `782f083a` (base), then `c004cba0`, then `a7a2fe6b`.
- Scope: the six original files plus `architecture.go` and `architecture_test.go`, as the orchestrator decided (D-12). `ai-report-guide.md` is unchanged in this commit.

**Files changed in `a7a2fe6b` (7 files, 60 lines added, 12 removed):**
- `internal/toolboxcommands/architecture.go` (F1): `architectureScan` first collects candidate folders. A month folder contributes its child directories (directories only, symlinks never followed). The scan then matches `architectureName` on the folder's base name. `prior_report` keeps the month folder in its repository-relative path. An unreadable month folder fails the pre-flight with the existing "reports directory is unreadable" text, the same as D-07.
- `internal/toolboxcommands/architecture_test.go` (F1): new `TestArchitectureScanFindsPriorReportInMonthFolder`. It has a nested `2026-10_1-report/2026-10-01_0900_architecture-report` and an older flat `2026-09-01_0900_architecture-report`, and `prior_report` must name the nested one.
- `internal/toolboxcommands/report_second_opinion.go` (F2 and F4): the condition is now `!(flatBundle || inMonthFolder)`, where `flatBundle` excludes a month-folder name. So the month folder itself is refused even when it holds an `index.html`.
- `internal/toolboxcommands/report_second_opinion_test.go` (F2): one new assertion. An `index.html` is written into `2026-09_3-reports/`, and the month folder itself is refused. The `// Pins:` comment is extended to cover this.
- `internal/toolboxcommands/report_index.go` (F5): `reportMonthPattern` is now `^[0-9]{4}-(0[1-9]|1[0-2])(_[0-9]+-reports?)?$`. The helper's doc comment names the canonical forms and says that plain `YYYY-MM` and hand-made `_<N>-report`/`-reports` forms with any N are read too.
- `internal/toolboxcommands/report_index_test.go` (F5): a `2026-09_2-report` case is added to the subtest list, with a one-line comment (2 lines in total).
- `skills/do-work-toolbox/actions/ai-report.md` (F3): line 52 now reads "a new bundle directly under `ai-reports/` that supersedes the prior one", and the checklist line now reads "a new `-rev<N>` bundle".

**RED** (the tests were changed before the production edit; exit 1):
- `TestArchitectureScanFindsPriorReportInMonthFolder`: `prior_report=ai-reports/2026-09-01_0900_architecture-report/index.html`. It picked the older flat report and missed the nested one.
- `TestAIReportIndexReadsBundlesInsideMonthFolders/2026-09_2-report`: `catalog paths=[ai-reports/2026-11-01_0900_b ai-reports/2026-09_2-report], want [... ai-reports/2026-09_2-report/2026-09-27_2200_a]`.
- `TestSecondOpinionBundleInsideMonthFolder`: `the month folder itself accepted as a bundle`.

**GREEN and checks** (run from the worktree root):

| Check | Exit | Wall |
|---|---|---|
| The three tests above, `-v` | 0 (all PASS, including 3 catalog subtests) | 1 s |
| `gofmt -l skills/do-work/tools/do-work-cli/internal` (printed nothing) | 0 | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands` | 0 | <1 s |
| `go test -C skills/do-work/tools/do-work-cli -count=1 -run 'SecondOpinion\|AIReportIndex\|Architecture' ./internal/toolboxcommands` | 0 (44 top-level PASS, 0 FAIL, 0 SKIP) | 17 s |
| `bash .../REQ-706-probe.sh` | 0 | <1 s |
| `bash _dev/tests/contract-regressions.sh` | 0 | 20 s |
| `bash _dev/tests/shipped-package-reference-contract.sh` | 0 | 1 s |
| `git diff --check` | 0 | <1 s |

**Decisions:**
- **D-14** (DECIDE & STATE): after F5, the month-folder definitions in `ai-report.md` and the guide stay as they are. Both already name only the canonical forms (`YYYY-MM_<N>-reports`, `YYYY-MM_1-report`) plus plain `YYYY-MM`, which is the state the orchestrator asked for. Neither states the exact suffix rule, so there was nothing to widen. Only the Go doc comment, which did state the rule, changed.
- **D-15** (DECIDE & STATE): the F5 test case uses `2026-09_2-report`, not `2026-07_2-report`. The fixture's nested report is dated 2026-09, so the month stays consistent. The suffix shape under test (`_<N>-report` with N > 1) is the same.
- **D-16** (DECIDE & STATE): F4 is met by the restructured condition `!(flatBundle || inMonthFolder)` that F2 needed, not by adding parentheses to the old clause. The grouping is explicit, and each half has a name.

**Discovered tasks:** none new. The earlier duplicate-basename note (supersedes resolution keys on `path.Base`) still stands for REQ-707. → report only
