# REQ-707 builder hand-back (ai-report organize files finished months into counted month folders and fixes every link)

- Branch: `worktree-agent-REQ-707-ai-report-organize-counted-month-folders`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-707-ai-report-organize-counted-month-folders`
- Base commit: `29e6bfa7` (main HEAD at dispatch)
- Commits, oldest first:
  - `e45a8fc2` [REQ-707] ai-report-organize verb: file finished months into counted month folders and fix every link
  - `6fa71772` [REQ-707] ai-report organize form, help line, routing words and guide section
- Worktree clean after the second commit. Nothing merged, rebased or pushed. Organize never ran against the main repository or the worktree's own `ai-reports/`; every run used `t.TempDir()` fixtures or scratch repositories under `/tmp/req707-scratch.Yt2e1z`.

## File manifest

Go files are in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/`.

- `report_organize.go` (new): the verb. Layout planner `planReportOrganizeLayout` + `hasWork` + `finalPath`, link planner `planReportOrganizeLinks` (rules A and B in `rewriteReportLinks`), dirty check `reportOrganizeDirtyPaths`, apply `applyReportOrganize` / `applyReportOrganizeStep`, handler `handleAIReportOrganize`.
- `report_organize_test.go` (new): the two tests from the plan plus four small test helpers.
- `commands.go` (modified): `CommandAIReportOrganize = "ai-report-organize"` and its handler entry.
- `commands_test.go` (modified): the registration test lists the new verb; count 11 to 12.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): the Input sentence names `organize`; new `### Organize form` before `### Revise form` (one paragraph, the launcher line, four steps).
- `skills/do-work-toolbox/actions/help.md` (modified): `  ai-report organize [--dry-run] File finished months into month folders; fix links`.
- `skills/do-work-toolbox/SKILL.md` (modified): `organize reports` added to the ai-report routing row.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): REQ-706's layout sentence gains "and `ai-report organize` puts them there"; new `## Organizing Reports by Month` (three paragraphs: layout and recount, the link promise and its limits, dry-run, dirty refusal, interrupted run, commit).
- `skills/do-work/docs/command-line-guide.md` (modified): the stale two-verb list becomes "The `ai-report-*` toolbox CLI verbs have no recipe".

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Build the REQ's `## Plan` files 1 to 9 as written, D-01 to D-15 settled. Tests first, calling `Handlers()["ai-report-organize"]` so RED is a test failure. The layout planner reads `ai-reports/` only (no Git), so REQ-708 can reuse `hasWork()`. The link planner takes its file set from `git ls-files --cached --others --exclude-standard`. Apply order: edits, mkdirs, moves, renames, removes, then `reportIndexWrite`. Lessons applied: `exact-basename-authority` (rule A and `finalPath` match only at a path segment boundary, so `X-rev1` never follows `X`); `final-boundary-identity` (right before each rename: no linked ancestor on either end, the source is a real directory, the destination is absent); `canonical-link-outlives-its-target` (rule A accepts `/` before `ai-reports/`, so `blob/main/ai-reports/...` and `file://` links are rewritten); `alternate-writer-contract-drift` (swept every code path that lists the reports folder).
- [x] **[APPLY]:** Wrote the two tests and the registration update, ran RED, then wrote `report_organize.go` and registered the verb; GREEN on the first implementation run. Then the action form, help line, routing words, guide section and the command-line guide line. Scope stayed at the nine files. Additions the plan did not name are listed under the anti-bloat check below.
- [x] **[UNIFY]:** `git diff --stat 29e6bfa7 HEAD` shows exactly the nine Scope files (below). I re-read every changed file in the diff. No debug output, no leftover TODO, no scratch file in the tree. I ran the built CLI by hand on four scratch repositories to cover cases the two tests do not: a GitHub `blob/main/` URL, a `file://` link with `:19`, the `my-ai-reports/` lookalike, a trailing dot, a misfiled bundle moved out of another month's folder, a stale `_9-reports` count, an empty untracked keeper folder, a duplicate bundle name across two folders of one month, dot files in month folders, a hand-made `index.html`, a symlinked `ai-reports`, a folder with no Git, and a bad flag. Each case behaved as the plan says. The action's commit recipe (`git add -A -- <paths>`, then `git commit -- <paths>`) left `git status` empty, and the second run printed `nothing to organize`.

```
 skills/do-work-toolbox/SKILL.md                    |   2 +-
 skills/do-work-toolbox/actions/ai-report.md        |  15 +-
 skills/do-work-toolbox/actions/help.md             |   1 +
 skills/do-work-toolbox/docs/ai-report-guide.md     |  10 +-
 skills/do-work/docs/command-line-guide.md          |   2 +-
 .../internal/toolboxcommands/commands.go           |   2 +
 .../internal/toolboxcommands/commands_test.go      |   6 +-
 .../internal/toolboxcommands/report_organize.go    | 612 +++++++++++++++++++++
 .../toolboxcommands/report_organize_test.go        | 270 +++++++++
 9 files changed, 913 insertions(+), 7 deletions(-)
```

Checks, run from the worktree root (wall time in seconds):

| Check | Exit | Time |
|---|---|---|
| `bash .../REQ-707-probe.sh` (before and after the commits) | 0 | 1 |
| `gofmt -l skills/do-work/tools/do-work-cli/internal` | prints nothing | <1 |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands` | 0 | <1 |
| `bash _dev/tests/contract-regressions.sh` | 0 | 21 |
| `bash _dev/tests/shipped-package-reference-contract.sh` | 0 | 1 |
| `bash _dev/tests/action-shell-blocks.sh` (83 fenced blocks, ShellCheck enabled) | 0 | 1 |
| `bash _dev/tests/prescribed-shell-canonicalization.sh` | 0 | 1 |
| `git diff --check` | 0 | <1 |
| `bash _dev/tests/staged-skills-contract.sh` (named in the plan's GREEN list, not in the brief's) | 2: "staged-skill contract probes are heavy-only; run `_dev/tests/maintainer-verify.sh --heavy`" | 0 |

I did not run the heavy lane or `maintainer-verify.sh`, as the brief requires.

## Proof record

**RED** (before any production edit): `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestAIReportOrganize.*|TestHandlersRegisterCanonicalToolboxCommands)$' ./internal/toolboxcommands` exited 1:

```
=== RUN   TestHandlersRegisterCanonicalToolboxCommands
    commands_test.go:16: canonical command "ai-report-organize" is not registered
    commands_test.go:20: registered 11 toolbox commands, want 12
--- FAIL: TestHandlersRegisterCanonicalToolboxCommands (0.00s)
=== RUN   TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen
    report_organize_test.go:102: ai-report-organize is not registered
--- FAIL: TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen (0.06s)
=== RUN   TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets
    report_organize_test.go:217: ai-report-organize is not registered
--- FAIL: TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets (0.06s)
FAIL
```

**GREEN:** the same command exited 0:

```
--- PASS: TestHandlersRegisterCanonicalToolboxCommands (0.00s)
--- PASS: TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen (0.18s)
--- PASS: TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets (0.13s)
PASS
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	0.847s
```

The two new tests take 0.31 s together.

**Package selection:** `go test -C skills/do-work/tools/do-work-cli -count=1 -run 'AIReport|Handlers|Architecture' ./internal/toolboxcommands` exited 0: 17 PASS, 0 FAIL, 3 SKIP. The 3 skips are the existing `TestAIReportJudge*` browser cases ("ai-report-judge browser cases are heavy-only"). They are not new.

## Decisions

The plan's D-01 to D-13 and the orchestrator's D-14 and D-15 are applied as written. Each decision below is DECIDE & STATE.

- **D-16 (Rule A also skips links Rule B kept or skipped):** rule A skips every relative-link span that rule B examined, not only the spans B rewrote. That covers links B kept because they still open, links that are broken now, and links that climb above the root. Why: with only the rewritten spans excluded, rule A would rewrite `../../ai-reports/X/index.html` inside a moving report and give it the wrong depth. That breaks the plan's own promise that "a link that still opens keeps its bytes, a broken link is never touched".
- **D-17 (commit paths only name paths Git knows):** commit paths list a step's path and its final path only when Git lists a file under that path. The link planner fills the new `CommitPaths` field. Why: `git add -A -- <path>` rejects the whole list when one pathspec matches nothing. I reproduced this in a scratch repository with an empty untracked keeper folder `ai-reports/2026-10/`. Removed month folders go through the same filter. Every file they held sits under a listed move, so nothing is lost. A mkdir target is not listed on its own because its moved bundles' final paths are. Not pinned by a test (the brief allows two tests).
- **D-18 (two finding codes):** `AI-REPORT-ORGANIZE-DIRTY` (outcome refused, exit 1) is for uncommitted changes. `AI-REPORT-ORGANIZE-FAILED` (outcome failure, exit 2) is for every other stop: no Git repository, an `ai-reports` that is a symlink or a file, a read error, and a step that failed part way (`changes` then lists what was done). The plan says organize "refuses" for the no-Git and symlink cases. I used one failure code to keep the surface small. Easy to change.
- **D-19 (left reasons):** the line format is `left <path> (<reason>)`. The reasons are `symbolic link`, `loose file`, `no date prefix` and `month folder holds other files`. A skipped month gets one line at the blocking path: `<why>; month YYYY-MM not organized`, where the blocker is the counted name in the way, the duplicate bundle's path, or an occupied target. A file in the way of a counted name also shows as `loose file`, so the same path gets two lines. A month folder whose month was skipped gets no extra line. The JSON `skipped_work` code is `AI-REPORT-ORGANIZE-LEFT`.
- **D-20 (dry-run output):** dry-run prints no catalog line (the catalog is only rebuilt on a real run) and no commit block. Its JSON has no `changes`, because nothing was written, and carries `skipped_work`.
- **D-21 (edit and mkdir records):** an edit's recorded change uses its final path, the same as the `links` line and the commit paths. A mkdir step prints no line and records no change, because the moves into that folder are listed.
- **D-22 (catalog I/O error):** a catalog result that is not success, including an I/O failure, is treated like the hand-made refusal: outcome `findings`, the moves stay, and the catalog pair is left out of the commit paths.
- **D-23 (one limit left out of the guide):** the guide does not list the plan's "a count change between an interrupted run and its rerun" limit. It gives the recovery instead: "commit or restore that with Git, then run it again". Why: the edge needs a failed run plus a changed month, and naming it would confuse readers more than help them.
- **D-24 (registration test uses the constant again):** during RED the registration test named the verb by its literal string, so RED was a test failure and not a compile error. After GREEN it uses `CommandAIReportOrganize`, like the other entries in that list.

## Discovered Tasks

- **impact-negligible** A hand-made `ai-reports/index.html` or `catalog.json` (no generator marker) is skipped by the link scan, because the plan treats the catalog pair as regenerated. A hand-made file is refused and never regenerated, so its links to moved reports break. Seen in a scratch repository. → report only
- **impact-negligible** A bundle inside a month folder can have a month-shaped name, such as `2026-11/2026-10/`. When organize moves it to the top level it becomes a month folder, and the catalog then reads its children as bundles. Very rare. Refusing a top-level target with a month-folder name would close it. → report only
- **impact-negligible** After a run fails part way, a rerun can stop with `AI-REPORT-ORGANIZE-DIRTY`. The first run's own uncommitted edits or moves sit under the paths that are still to move. This matches plan D-05 ("Git is the undo"), and the guide and the failure text say to commit or restore with Git first. The plan's phrase "a rerun finishes" holds only after that step. → report only

## Lessons read

`_dev/primes/lessons-releases.md` (whole file). From `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, only the bullets the brief's grep returned: the `final-boundary-identity` bullets (lines 44, 45, 73, 74, 79, 80, 81, 91), the `exact-basename-authority` bullet (line 105) and the 0.309.1 `alternate-writer-contract-drift` bullet (line 124). Also the `## Traps` sections of `prime-do-work-cli.md`, `prime-action-files.md` and `prime-releases.md`, and the crew rules general, coding-guardrails, shared-principles, communication-style, testing and backend.

Sweep for `alternate-writer-contract-drift`: every non-test code path that lists the reports folder (`report_index.go` `reportIndexWalk`, `architecture.go` pre-flight, `report_second_opinion.go` bundle check) decides month folders with `reportMonthFolder`. That function accepts the counted names organize writes. No reader needs a change.

## Anti-bloat check

The YAGNI list is honoured: no `--commit`, no rollback or transaction code, no new resultmodel struct, no Just recipe, no exported API, no new package, no `url()`, `srcset` or script parsing, no percent-decoding, no symlink retargeting, no deeper nesting, no hint printing. There are two new tests plus the registration update, and no rewriter unit tests.

Added but not named by the plan, each with its reason:

- Types `reportOrganizeStep`, `reportOrganizeLeft`, `reportOrganizeEdit`, `reportOrganizePlan`: the plan's code-shape comment describes Steps, Left and Edits. These are those shapes.
- Field `reportOrganizePlan.CommitPaths`: D-17.
- `handleAIReportOrganize`: "its handler" in plan item 3.
- `reportOrganizeDirtyPaths`: the plan's dirty check.
- `applyReportOrganizeStep`: one step's identity checks plus the syscall, kept beside `applyReportOrganize` so the final-boundary checks sit right before `rootedfs.Rename`.
- `(plan) rewriteReportLinks`: rules A and B for one file, split out of `planReportOrganizeLinks` for reading.
- `reportBundleMonth`: the date-prefix month with a fallback to the month folder's month. Used twice.
- `reportNameByte`: rule A step 1's byte class `[A-Za-z0-9_.-]`.
- `reportOrganizeFailure`: builds the failure result. Used four times.
- Constant `reportMentionStopBytes` and regex variables `reportBundleDatePrefix`, `reportLinkScheme`, `reportLinkAttribute`, `reportMarkdownTarget`: rule A's stop set, the plan's date prefix, and rule B's scheme and token patterns.
- Finding code `AI-REPORT-ORGANIZE-FAILED` and skipped-work code `AI-REPORT-ORGANIZE-LEFT`: D-18 and D-19.
- Test helpers `runReportOrganize`, `writeOrganizeFiles`, `freezeOrganizeClock` and `snapshotOrganizeTree`. The last one exists because the existing `snapshotReportTree` walks `.git`, whose index `git status` may refresh, and records no directories, so a dry run that created an empty month folder would still pass.

## Proposed CHANGELOG entry

**Report Organizer Files Finished Months Into Counted Folders and Fixes Every Link**

A busy `ai-reports/` folder becomes readable by month without breaking a single link: `ai-report organize` files each finished month's reports into one folder whose name shows how many reports it holds, and fixes every link to a moved report.

- New `do-work-toolbox ai-report organize [--dry-run]`, backed by the do-work-cli verb `ai-report-organize`. Finished months move into `ai-reports/YYYY-MM_N-reports/` (`YYYY-MM_1-report` for one report). The current month stays directly under `ai-reports/`, and report writers are unchanged.
- Every run recounts. A stale count or a plain `YYYY-MM` folder is renamed, and a current-month folder is emptied back to the top. A second run changes nothing. Undated folders, loose files and symbolic links stay and are listed.
- Links: path mentions `ai-reports/<report>/...` in every text file Git lists (requests, docs, history, GitHub `blob/main/` and `file://` links) and relative `href`, `src` and Markdown links between reports are rewritten, so every link that opened before still opens. A link that was already broken, or still opens, keeps its bytes. The catalog is regenerated.
- Safety: `--dry-run` prints the whole plan and writes nothing. Organize refuses when a path it would move or rewrite has uncommitted changes. It never renames onto an existing path; that month is skipped and listed instead. It prints the exact `commit paths:` that the action commits. Nothing is pushed.
- Not covered: binary and Git-ignored files, symbolic links to a report, CSS `url()`, `srcset`, script-built links, percent-encoded paths, Markdown reference links, relative paths in plain prose, links to a removed current-month folder, and copies outside this checkout. A mention of another repository's report with the same name is rewritten too.

## Proposed lesson bullet

For `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, if the orchestrator wants one (otherwise "none"):

- [family: exact-path-commit-list] <version> (REQ-707, 2026-11-02): **a printed exact-path commit list is a contract with `git add`, and `git add -A -- <paths>` rejects the whole list when one pathspec matches nothing.** An organize step can move or remove a directory that Git has never seen, such as an empty untracked month folder or a bundle that holds only ignored files. Its old and new paths then match nothing, and the action's commit fails. Filter the printed paths by Git's own listing (`ls-files --cached --others --exclude-standard`) taken before the moves, and prove completeness end to end: commit exactly the printed list, then require an empty `git status --porcelain --untracked-files=all`.

## The function REQ-708 should call

```go
plan, err := planReportOrganizeLayout(repositoryRoot, toolboxNow().Format("2006-01"))
// then:
plan.hasWork() // bool: true when organize would move, rename or remove anything
```

Signatures: `func planReportOrganizeLayout(repositoryRoot, currentMonth string) (reportOrganizePlan, error)` and `func (plan reportOrganizePlan) hasWork() bool`. Both are unexported in package `toolboxcommands` (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go`), per D-14. The planner reads only `ai-reports/`, with no Git and no other file scan. It returns an empty plan when `ai-reports/` is missing and an error when `ai-reports` is a symlink or a file. Months before `currentMonth` (`YYYY-MM`, local time) count as finished.

## Addendum (review fixes)

The review is `REQ-707-review.md` (82%, acceptance Partial). The orchestrator decided D-26, which adds three files to the scope: `completed-work-presentation-reference.md`, `tidy-repo.md` and `report_index.go`.

- **Commit:** `73319ce9` [REQ-707] review fixes F1-F6: refuse a hand-made catalog before any move, filter ignored catalog paths, skip cross-month name collisions, docs. It is on the same branch and worktree, on top of `6fa71772`. The worktree is clean.
- **Branch total now:** `git diff --stat 29e6bfa7 HEAD` shows 12 files changed, 999 insertions, 35 deletions. That is the 9 Scope files plus the 3 D-26 files.

### Files in `73319ce9`

- `report_index.go` (modified):
  - New helper `reportCatalogHandMadeFindings(repositoryRoot, command string) []resultmodel.CommandFinding`. It holds the generator-marker test, the same code and evidence as before. `reportIndexWrite` now calls it, so `ai-report index` behaves the same (`TestAIReportIndexRefusesHandMadeCatalog` still passes).
  - The comment at line 21 names organize as the one exception to bundle immutability (F5).
- `report_organize.go` (modified):
  - F1: after `hasWork()` and before the link plan, organize refuses a hand-made catalog pair. This applies on a dry run too.
  - F2: catalog paths enter `commit paths:` only when `git ls-files --cached --others --exclude-standard` lists them.
  - F4: a move target that an earlier month already planned skips and lists the later month.
  - The comment on the catalog-failure branch is updated.
- `report_organize_test.go` (modified): new assertions in `TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets`:
  - the fixture adds `2026-12/d` and a `.gitignore` that lists `ai-reports/index.html`;
  - a hand-made `ai-reports/index.html` is refused with `AI-REPORT-INDEX-HAND-MADE` and the tree is unchanged;
  - the December line `left ai-reports/d (another month's report moves here; month 2026-12 not organized)` appears and `2026-12/d` stays;
  - the commit block lists `ai-reports/catalog.json` and leaves out the ignored `ai-reports/index.html`.
- `ai-report.md` (modified): step 2 adds "A refusal writes nothing. A `findings` result means the reports moved but the catalog could not be written: commit as step 3 says, report the finding, and rerun the catalog command once its cause is fixed."
- `ai-report-guide.md` (modified):
  - F3: the Input block line `do-work-toolbox ai-report organize [--dry-run]`.
  - F6: the Not covered list adds unquoted HTML attributes, and a mention followed directly by a character that can be part of a file name.
  - The refusal sentence now also names a hand-made `catalog.json` or `index.html`.
- `completed-work-presentation-reference.md` and `tidy-repo.md` (modified): one clause each naming `ai-report organize` as the one exception (F5).

### RED and GREEN

All three REDs come from `go test ... -run '^TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets$'`, each run before its fix:

1. Before the F1 fix: `report_organize_test.go:239: a hand-made catalog page must refuse before any write`. The result was `Outcome:failure`, with evidence `move "ai-reports/2026-12/d" -> "ai-reports/d": ai-reports/d already exists`. The old code went straight to applying the moves, and F4's collision stopped it part way. Exit 1.
2. After the F1 fix, before the F4 fix: `report_organize_test.go:249: organize: ... Outcome:failure ... AI-REPORT-ORGANIZE-FAILED` (the cross-month collision stopped the run part way). Exit 1.
3. After the F4 fix, before the F2 fix: `report_organize_test.go:293: commit paths must list the catalog and leave out the ignored catalog page`. The commit block listed `ai-reports/index.html`. Exit 1.

GREEN after all three fixes:

```
--- PASS: TestHandlersRegisterCanonicalToolboxCommands (0.00s)
--- PASS: TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen (0.23s)
--- PASS: TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets (0.18s)
```

### Checks, run from the worktree root

| Check | Exit | Time |
|---|---|---|
| `gofmt -l skills/do-work/tools/do-work-cli/internal` | prints nothing | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands` | 0 | <1 s |
| `go test -C skills/do-work/tools/do-work-cli -count=1 -run 'AIReport\|Handlers\|Architecture' ./internal/toolboxcommands` | 0 (17 PASS, 0 FAIL, 3 SKIP; the skips are the existing heavy-only `TestAIReportJudge*` browser cases) | 2 s |
| `REQ-707-probe.sh` | 0 | 1 s |
| `contract-regressions.sh` | 0 | 21 s |
| `shipped-package-reference-contract.sh` | 0 | 1 s |
| `action-shell-blocks.sh` | 0 | 2 s |
| `prescribed-shell-canonicalization.sh` | 0 | <1 s |
| `git diff --check` | 0 | <1 s |

### New decisions (DECIDE & STATE)

- **D-27:** the organize refusal for a hand-made catalog reuses the finding code `AI-REPORT-INDEX-HAND-MADE` and its evidence text. Only `next_argv` changes, to `ai-report-organize`. A new code would split one condition across two names.
- **D-28:** the hand-made check runs only when organize has work. A no-work run still prints `nothing to organize`, because it writes nothing, not even the catalog.
- **D-29:** after the catalog is written, the F2 filter asks Git once (`ls-files --cached --others --exclude-standard -- <catalog pair>`) rather than calling `git check-ignore`. A file that is tracked but matches an ignore rule stays listed, because `git add` accepts it. The `auditGit` call also cannot tell `check-ignore`'s "not ignored" exit status from a real error.
- **D-30:** a cross-month collision is reported on the target path with the reason `another month's report moves here; month YYYY-MM not organized`. The later month in sorted order waits. Running organize again after the user renames the clash finishes the job.
- **D-31:** F7 (decision numbers cited wrongly in the REQ record) is not mine to fix. The REQ file is outside my write boundary. My hand-back's "D-01 to D-13" came from the brief, and the real range is D-01 to D-07.

The earlier Discovered Task about a hand-made catalog's links is now moot: organize refuses before it scans or moves anything.
