## Review: REQ-707

**Approve with follow-ups**: organize does what the REQ asks, and every scratch case I ran kept links open and committed cleanly. In two edge cases the action's commit step does not work: a hand-made catalog file, and a catalog pair that Git ignores.
Route C | merge range `95146560..429cdb97` (builder commits `e45a8fc2`, `6fa71772`)

### What's built
- `do-work-cli ai-report-organize [--dry-run]` moves finished-month reports into `ai-reports/YYYY-MM_N-reports/` (`_1-report` for one report). It recounts and renames stale or plain `YYYY-MM` folders, empties current-month folders, rewrites repository path mentions (rule A) and relative links inside reports (rule B), refuses when it would touch dirty paths, regenerates the catalog and prints a `commit paths:` block.
- The new `### Organize form` in `ai-report.md` runs the verb and commits by exact path. The help line, the routing words, the guide section and the command-line guide sentence are in place.
- Missing: the guide's `## Input` command block does not list `organize` (F3).

### Decisions / risks for you
- F1: when the catalog file is hand-made, organize has already moved the reports, but the action tells the agent to stop. The move stays uncommitted, and a rerun prints `nothing to organize`, so the list of paths to commit is lost. Smallest fix: refuse during planning, before any write, when the catalog pair has no generator marker. Then step 2 is true as written. Value: no half-done state. Risk: about 10 lines of code. The other fix is prose only: step 2 commits the printed paths first when the result is `findings`.

### Findings

**Important:**
- F1 `skills/do-work-toolbox/actions/ai-report.md:59`: step 2 says "A refused or failed result names its paths (... a hand-made catalog file ...); report it and stop". A hand-made catalog yields outcome `findings` (exit 1, D-22) only after every move and link edit was applied. The agent stops, the moves stay uncommitted, and a rerun prints `nothing to organize`, so the `commit paths:` list is never printed again. I reproduced this in scratch repo `e4`. Fix: check the generator marker of the catalog pair while planning and refuse before any write, or make step 2 commit the printed paths on a `findings` result. — impact-user-visible → report only

**Minor:**
- F2 `report_organize.go:143-148`: the catalog pair goes into `commit paths:` without the D-17 "Git knows a file here" filter. When `.gitignore` lists `ai-reports/catalog.json` and `ai-reports/index.html` (or all of `ai-reports/`), `git add -A --` exits 1 ("paths are ignored"), and `git commit --` fails ("pathspec did not match any file(s) known to git"), so nothing is committed. I reproduced this in repos `e5` and `e5b`. When all of `ai-reports/` is ignored, rule B never runs either, so links between reports break. The guide's "Not covered: Git-ignored files" line documents that part. Fix: add the catalog paths only when Git does not ignore them. — impact-user-visible → report only
- F3 `skills/do-work-toolbox/docs/ai-report-guide.md:68-82`: the guide's `## Input` command block lists every ai-report form except `organize`. This file was in scope. Fix: add `do-work-toolbox ai-report organize [--dry-run]` with a one-line description. — impact-user-visible → report only
- F5 Restatement: `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:71` says "Existing artifacts are immutable: never delete, truncate, merge into, rename, migrate, or overwrite them", and the `report_index.go:21` comment says "published bundles are immutable". Neither names organize as the one exception: it renames month folders, moves bundles and rewrites bytes inside them. `skills/do-work-toolbox/actions/tidy-repo.md:49` calls stale paths in dated history records "correct history", but organize rewrites them (D-07). Fix: add one clause naming `ai-report organize` as the only command that moves bundles and rewrites their links. — impact-negligible → report only

**Nit:**
- F4 `report_organize.go:290-296`: for a month that is not finished, the planner checks whether the top-level name is taken by what exists now. It does not check names planned by other months. So `2026-11/same` plus `2026-12/same` both plan `ai-reports/same`, and the second move stops the run part way (exit 2, first move done). The final-boundary check prevented any overwrite. Repo `e10`. Fix: track planned top-level destinations across months and skip the later month. — impact-negligible → report only
- F6 `report_organize.go:26,31`: two link gaps are missing from the guide's "Not covered" list. An unquoted `href=../x/index.html` is not rewritten. A bare bundle-folder mention glued to a byte that is not a stop byte (`ai-reports/<bundle>!`, or a Unicode apostrophe) is not rewritten. Repos `e6` and `s1`. Fix: name "unquoted attributes" in the guide, or accept. — impact-negligible → report only
- F7 REQ record: `## Plan` holds D-01 to D-07, but Exploration and the hand-back cite "D-01 to D-13". D-14 cites "D-12's escalation", which is D-06. The plan validation says duplicate names are "covered by D-04", but that is D-02. Fix: correct the citations. — impact-negligible → report only

### Requirements Checklist

- [x] `organize` form in `ai-report.md` with routing and help lines, backed by one verb (`CommandAIReportOrganize`): delivered
- [x] Month from the `yyyy-mm-dd` prefix. Finished means before the current local month. Future months count as current, and current-month reports stay flat: delivered
- [x] Finished month goes into `YYYY-MM_N-reports` (`_1-report` for N=1), with N counted after the move: delivered (repos `s1`, `s2`, `m1`)
- [x] Recount on every run, plain `YYYY-MM` renamed, current-month folder emptied: delivered (second test, repos `s2`, `m1`)
- [x] Every link fixed: repository mentions (`:19`, `blob/main/`, `file://`, `../ai-reports/`, `?query#fragment`, `blocked_by`, `## Reports`), relative links between reports, catalog, binary never edited: delivered. The documented limits are listed in the guide, and F6 names two more.
- [x] Catalog regenerated through `reportIndexWrite`: delivered
- [x] Second run changes nothing: delivered (every scratch repo printed `nothing to organize`)
- [x] `--dry-run` prints moves, renames and edits and writes nothing: delivered (test snapshot)
- [x] Undated folder, loose file and symlink stay in place and are listed: delivered
- [~] Commit by exact paths as Step 8 says, never push: delivered for tracked catalogs. It fails for a hand-made catalog (F1) and for an ignored catalog pair (F2).
- [~] Guide explains the layout and the command: delivered. The Input block omits the form (F3).
- [x] Constraints: no change to any writer, to `report_index.go` or to install/update code. Only organize moves reports.

### Acceptance Testing

**Result: Partial** (implementation and integration stages)
- `go test -run '^(TestAIReportOrganize.*|TestHandlersRegisterCanonicalToolboxCommands)$'`: 3 PASS. `-run 'AIReport|Handlers|Architecture|ReportIndex'`: ok. `go vet` 0, `gofmt -l` empty.
- CLI built from ROOT at `429cdb97` and run on 13 scratch repositories under `/tmp/rev707.TeguyD/repos/`. The exact commit recipe passed in each case below: `git add -A -- <paths>`, then `git commit -- <paths>`. After it, `git status` was empty, a staged unrelated file stayed staged and out of the commit, and the second run printed `nothing to organize`.
  - `s1`: existing `2026-09_2-reports` plus a loose report. The keeper was renamed to `_3-reports`, and `?x=1#sec` and `#top` were kept. `../../docs/x.md` became `../../../docs/x.md`. A doc's `../ai-reports/...`, `:19`, `blob/main/` and `file://` mentions were rewritten. `../2026-09_2-reports/` kept its trailing slash.
  - `s2`: organized repo plus a new loose September report. Recounted to `_4-reports`. Links into the renamed folder from the current report, `docs/`, `blocked_by` and `## Reports` were all rewritten.
  - `m1`: misfiled August report inside `2026-09/`, plus an empty untracked `2026-10_7-reports`. Moves and renames composed correctly, links crossing months were recomputed, and the D-17 filter kept the empty keeper out of the commit paths.
  - `d16`: rule A and rule B interaction (D-16). A broken `../ai-reports/X` stayed untouched, `/ai-reports/X` was rewritten, and repo-rooted relative links became opening relative links.
  - `e1` duplicate name (month skipped and listed), `e2` counted name exists as a symlink (skipped, nothing written), `d1` untracked file in a keeper (refused, `AI-REPORT-ORGANIZE-DIRTY`).
  - Failures: `e4` hand-made catalog (F1), `e5`/`e5b` ignored catalog (F2), `e10` collision of current and future months (F4), `e6` unquoted href (F6).
  - Scale: 270 bundles, 3,271 files and a 200 MB binary. The run took 0.9 s.
- The symlink-swap boundary was verified by reading only: `validateNoLinkedAncestors` on both ends, a real-directory source and an absent destination are checked right before `rootedfs.Rename`. I did not race it.

### Suggested Additional Testing

- Deployment: unassessed. After `do-work update` in a consumer, run `ai-report organize --dry-run` through the installed launcher.
- Live acceptance: unassessed. In a real consumer repository with about 250 bundles, run `--dry-run`, then the real run, then `ai-report judge` on one moved bundle and one current-month report that links to it.
- Heavy lanes not run (user decision): `do-work-cli-integrations`, `staged-skills`, `updater`, `installer`.
- Interrupted-run recovery: force a failure after the edits (as in F4) and follow the guide's "commit or restore that with Git" text. Note that `git restore` alone leaves the moved untracked copies behind.
- macOS: a current-month folder holding only `.DS_Store` is listed on every run (by design, per D-02). Confirm a user accepts that.

### Anti-bloat count

There are no extra options, files or decorative tests. Beyond the plan's five functions there are six helpers: `reportOrganizeFailure`, `reportBundleMonth`, `rewriteReportLinks`, `reportNameByte`, `reportOrganizeDirtyPaths` and `applyReportOrganizeStep`. Each is used and each splits planned behavior. There is one new field, `CommitPaths` (D-17), and four test helpers. None of this is bloat.

### Scores (on the record, not the headline)

**Overall: 82%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | All delivered. The commit step fails in two edge cases (F1, F2). The guide's Input block is stale (F3). |
| Code Quality | 88% | Clean, reuses rooted helpers, safe apply order. Catalog paths are not filtered (F2). Cross-month collision is caught only at apply time (F4). |
| Test Adequacy | 85% | Two focused tests pin the proof, the recount and the refusals, with recorded red-green. Duplicate names, D-17, D-22 and the URL forms are not pinned (plan YAGNI). |
| Scope | 100% | 9 files declared, 9 touched |
| Risk | Low | Dirty check, final-boundary checks and Git as the undo. One part-way stop in a contrived case. |
| Acceptance | Partial | Implementation and integration pass. Hand-made and ignored catalogs break the commit step. |

Overall = (95 + 88 + 85 + 100) / 4 = 92, minus 10 for Partial = 82.

### Follow-ups created
None (7 findings report only)

## Review

**Overall: 82%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 88% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `skills/do-work-toolbox/actions/ai-report.md:59`: step 2 stops on "a hand-made catalog file", but organize has already applied every move and edit (outcome `findings`, D-22). The commit is skipped and a rerun prints `nothing to organize`, so the commit paths are lost. Fix: refuse a hand-made catalog pair during planning, before any write, or commit the printed paths on `findings`. — impact-user-visible → report only

**Minor findings:**
- F2 `report_organize.go:143-148`: the catalog pair skips the D-17 filter. When Git ignores it, `git add -A --` exits 1 and `git commit --` fails, so nothing is committed. — impact-user-visible → report only
- F3 `ai-report-guide.md:68-82`: the `## Input` command block omits `ai-report organize`. — impact-user-visible → report only
- F5 Restatement: `completed-work-presentation-reference.md:71` and `report_index.go:21` state that bundles are immutable and never renamed or migrated, without naming organize as the exception. `tidy-repo.md:49` calls stale history paths correct, which organize's D-07 contradicts. — impact-negligible → report only
- Nits: F4 `report_organize.go:290-296`: current and future month folders with the same undated child stop the run part way. F6: unquoted `href=` and glued bare-folder mentions are not rewritten and not in the guide. F7: the REQ cites D-08 to D-13, D-12 and D-04 for decisions numbered D-06 and D-02. — impact-negligible → report only
**Acceptance:** Partial — implementation and integration stages: 3 focused tests pass, and 13 scratch repos (recount, keeper rename, misfiled bundle, `?query#fragment`, `../../docs`, `../ai-reports/`, D-16, dirty refusal, duplicate, symlink target) committed cleanly with the exact recipe, but a hand-made or ignored catalog breaks the commit step.
**Restatement sweep:** redefined: who may move reports or rewrite their links (organize only), the counted month-folder layout now written by a command, the ai-report form list, the toolbox verb list. Stale: `ai-report-guide.md` Input block (F3), immutability statements in `completed-work-presentation-reference.md:71`, `report_index.go:21`, `tidy-repo.md:49` (F5). Agreeing: `help.md`, `SKILL.md`, `command-line-guide.md`, `ai-report.md` catalog/revise/second-opinion rules, `report_second_opinion.go`, `architecture.go`, `README.md`, `suite/modules.tsv`.
**Suggested testing:** 5 items
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*

## Re-check (429cdb97..c52e7f54)

**Verdict: Approve.** All seven findings are closed. I found one new Nit. I would now give **95%** with **Acceptance: Pass** (implementation and integration stages).

Delta read in full: `73319ce9` (7 files, +96/−38). Changed functions: `handleAIReportOrganize` and `planReportOrganizeLayout` in `report_organize.go`, `reportIndexWrite` and the new `reportCatalogHandMadeFindings` in `report_index.go`, and the second organize test. The prose changes are in `ai-report.md` step 2, the guide, Collision-Safe Publication and `tidy-repo.md`. Scope grew by three files under D-26. The builder's D-27 to D-31 are in the hand-back addendum.

Checks: I built the CLI from ROOT at `c52e7f54` into `/tmp/rev707.TeguyD/dwc2`. `go test ./internal/toolboxcommands` passes (whole package, 19.9 s), `go vet` returns 0 and `gofmt -l` prints nothing.

| Finding | Status | Evidence |
|---|---|---|
| F1 hand-made catalog left the moves uncommitted | Fixed | `report_organize.go:103` calls the shared `reportCatalogHandMadeFindings` (`report_index.go:302`) before any write. Repo `r-e4`: dry run and real run are both `refused` (exit 1) with an empty `git status`. After moving the file aside, the rerun and the recipe commit cleanly and a second run prints `nothing to organize`. Step 2 now says a refusal writes nothing and a `findings` result still commits. The new test assertions pin it with a tree snapshot. |
| F2 ignored catalog pair broke the commit recipe | Fixed | `report_organize.go:152` lists only the catalog files Git would take. Repo `r-e5` (all of `ai-reports/` ignored): commit paths are `docs/n.md` only and the recipe commits. Repo `r-e5b` (only the pair ignored): the recipe commits. Repo `r-e5c` (catalog tracked but matching an ignore rule) is still listed and commits, as D-29 says. In `r-e5`, the link between reports still breaks. The guide documents that ("Git-ignored files"). |
| F3 guide Input block missing `organize` | Fixed | `ai-report-guide.md:76` |
| F4 cross-month collision stopped the run part way | Fixed | `plannedTargets` (`report_organize.go:318`). Repo `r-e10`: the current month's `same` moves, and month 2026-12 is skipped and listed. Exit 0, the recipe commits and the second run is stable. Pinned in the test. |
| F5 immutability restatements | Fixed | Collision-Safe Publication, `tidy-repo.md:49` and the `report_index.go:21` comment now name organize as the one exception. |
| F6 unquoted attributes and glued mentions | Fixed (documented) | Both are in the guide's "Not covered" list. |
| F7 wrong decision citations | Fixed | REQ D-14 now cites D-06, plan validation cites D-02, and the Decisions text says D-08 to D-13 are unused. |

**New finding:**
- N1 `report_organize.go:152-158`: two code paths could lose or swamp the commit paths. Neither is reachable today. (1) If `reportIndexWrite` succeeds with no changes, the argument list ends with `--` and no paths, so `git ls-files` lists every file into `commit paths:`. That happens only when `ai-reports/` is missing right after organize wrote into it. (2) If that `ls-files` call fails, the run ends as `failure` after the moves, with no commit paths printed. Fix: skip the call when `catalog.Changes` is empty, and on an `ls-files` error print the step commit paths before failing. — impact-negligible → report only

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Overall = (100 + 92 + 88 + 100) / 4 = 95.
**Acceptance:** Pass — implementation and integration stages: the whole package passes, and the reruns of e4, e5, e5b, e10 (plus e5c, e10b) commit cleanly with the action's exact recipe and leave a stable second run.
**Follow-ups created:** None (1 finding report only)
