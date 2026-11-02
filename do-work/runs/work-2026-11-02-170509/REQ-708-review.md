## Review: REQ-708

**Approve**: install, update, `ai-report index` and every new report now print one line recommending `ai-report organize`, and only when organize has work to do. The hint never changes an outcome or an exit status.
Route B | merge range `ee789713..e4c9d951` (builder commit `ecc9e0b5`, merge `e4c9d951`)

### What's built
- One check, `ReportOrganizeHint` (`report_organize.go:41`), reuses REQ-707's Git-free layout planner. Suite install (written, not declined) and every successful update print its line once to stderr, after the success line. `ai-report index` appends it to the catalog output. `ai-report-organize --check` gives the same line, or nothing, to the three report writers.
- The ai-report Step 8 (default, `--kind`, revise), architecture-report Step 7 and stakeholder-report Step 6 run `--check` and print its line just before the final path line. `version.md` step 4 tells the agent to relay the line. The guide says when the hint appears, including the first-update limit (D-06).
- Still missing: the README's paste-in agent upgrade prompt does not ask the agent to relay the line (F1).

### Decisions / risks for you
- None. `suiteinstall` now imports `toolboxcommands` (D-02). The binary already linked both packages and there is no import cycle.

### Findings

**Important:**
- None.

**Minor:**
- F1 `README.md:49`: the "Upgrade an existing installation with an AI agent" prompt runs `update-suite` directly. It asks the agent to report only the previous version, the resulting version, and whether the update completed or was cancelled. `version.md:50` step 4 now also relays the organize line, but this second agent-facing update instruction was not changed, so an agent following it will likely drop the stderr hint. This is the same family as the `alternate-writer-contract-drift` lesson. README was declared out of Scope, so this is not scope drift. Suggested fix: end the prompt's last sentence with "and relay the `ai-report organize` recommendation line if the command prints one". — impact-user-visible → report only

**Nit:**
- F2 `ai-report.md:195`, `architecture-report.md:119`, `stakeholder-report.md:51`, `ai-report-guide.md:110`: the prose says the hint prints "while a finished month's reports are still loose". The check actually fires on any organize step, which also includes these cases (reproduced with the built CLI): a stale count (`2026-09_2-reports` holding one report → rename), a plain `YYYY-MM` folder, or a current-month report inside a month folder (`2026-11/` → flatten). The hint line itself is accurate ("ai-reports/ needs organizing by month"). Suggested fix: "while organize has anything to move or rename". — impact-negligible → report only
- F3 `suite_commands.go:110`: the comment "A declined install leaves no toolbox to run organize with" is false for a declined reinstall over an existing suite. The behavior (stay silent on decline) is fine. Suggested fix: "A declined install changed nothing, so it stays silent." — impact-negligible → report only
- F4 Restatement (wave-end, bundle location): `stakeholder-answers.md:94` ("a fresh sibling bundle") and `stakeholder-report.md:62` ("every regeneration is a fresh sibling") still use the "sibling" wording REQ-706 F3 removed from `ai-report.md`. Once organize files the prior report into a month folder, the regenerated top-level bundle is no longer its sibling. Suggested fix: "a fresh bundle directly under `ai-reports/`". — impact-negligible → report only
- F5 `report_organize.go:35`: the hint text starts with `ai-reports/`. stakeholder-report's path line starts the same way, so a caller that finds the path by prefix rather than by position could pick up the hint. No shipped consumer does this, and publication refuses a `blocked_by` that does not match the report path. Suggested fix: start the line with words, for example "Reports in ai-reports/ need organizing by month: …". — impact-negligible → report only

### Requirements Checklist

- [x] One shared program check using REQ-707's finished-month rule: delivered (`ReportOrganizeHint` wraps `planReportOrganizeLayout(...).hasWork()`; every caller uses it, the writers through `--check`).
- [x] Install prints the recommendation only when the check finds something: delivered. D-01 reads "install" as the suite install (`install-suite`), which matches the user's words "when this is installed ... for the repository". The toolbox `install` action installs companion skills only.
- [x] `do-work update` prints it only when the check finds something: delivered for the engine (text and JSON callers, the hint on stderr) and the version.md agent path. The README agent prompt is F1.
- [x] `ai-report index` prints the hint under the same condition: delivered (`report_index.go:80-89`, after a successful write only).
- [x] Every new report (ai-report default, `--kind`, revise, architecture-report, stakeholder-report) prints it: delivered in prose. A `--kind` report reaches it through `ai-report-reference.md:148` ("Steps 6 to 8 ... apply unchanged"). Second-opinion, judge, find and organize correctly do not (D-08).
- [x] None of them moves a report or rewrites a link: delivered. `--check` only reads `ai-reports/`, and the Go test snapshots the tree before and after.
- [x] Nothing prints when there is nothing to organize: delivered (organized, empty, symlinked, loose-file-only and file-instead-of-folder fixtures all print nothing and exit 0).
- [x] UR constraints: writers still write new bundles directly under `ai-reports/`, and only organize moves reports.

### Acceptance Testing

**Result: Pass** (Implementation and Integration stages)
- Focused Go tests at `e4c9d951`: `TestReportOrganizeHintPrintsOnlyWhenWorkRemains` and the five other AIReportIndex/AIReportOrganize tests pass (8 PASS). `TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose` and `TestInstallNarrationNeverReachesStandardOutput` pass. `go vet` passes on both packages, `gofmt -l` reports nothing, and `git diff --check` is clean.
- CLI built from the merge into `/tmp/req708-review/dwc`. Scratch fixtures in `/tmp/req708-fx.zxVn`, outside ROOT:
  - These print the hint and exit 0: a loose `2026-09-27_2200_a`, a stale count, and a current-month `2026-11/` folder.
  - These print nothing and exit 0: an organized tree (current-month report plus `2026-09_1-report`), a correctly counted `2026-10_1-report`, no `ai-reports/`, a symlinked `ai-reports`, a loose-file-only folder, and `ai-reports` as a regular file.
  - `ai-report-index`: text output ends with the hint on the loose fixture and has no hint on the organized fixture. JSON output is a plain success.
  - `--check --format json` is a plain success. `--check --dry-run` gives the new usage line with exit 2.
- Scratch tests in a copy of the module (`/tmp/req708-review/cli-copy`, never in ROOT) covered the paths the diff does not pin:
  - Declined install: no hint, exit 0.
  - Already-current update: one hint, after "You're up to date".
  - Declined update: one hint, after "Update cancelled".
  - Unreachable upstream: exit 2 with no hint.
  - So the hint is never doubled on update (it is in the handler, not `RunInstall`), never prints on a declined install, and never changes an exit status.
- Agent update path: `version.md` step 3 runs `--format json`, the JSON result stays on stdout, and the hint is on stderr. Step 4's new clause tells the agent to relay it, so the JSON parse is unaffected.
- New-bundle writers: every writer prints the `--check` line before its final path line, so a caller that takes the last line still gets the path. stakeholder-report's callers (`stakeholder-answers.md:65`, `clarify.md:162`, and work.md Step 8 through the Stakeholder REQ Template) take the printed path, and publication `answer.go:312-314` refuses a `blocked_by` that does not match the report path. architecture-report has no last-line contract, so the hint is its allowed ninth line.
- Double printing: the hint can show twice only in a revise reply, if the agent echoes the index output of revise steps 1 and 6 as well as Step 8's line. The builder already recorded this under Discovered Tasks, so it is not a new finding.

### Restatement sweep

This is the wave end of UR-161 (REQ-706, REQ-707, REQ-708). The sweep covered `skills/`, `suite/` and `README.md`. It left out CHANGELOG history, tests and testdata.

| Element | Agreeing | Stale |
| --- | --- | --- |
| Where a bundle may sit (REQ-706) | `ai-report.md:42/52/65/67/85`, `ai-report-guide.md:96/98/106/114/126`, `report_second_opinion.go:389-392`, `architecture.go:68-86`, `report_index.go:51/99-101`, `ai-report-reference.md:120`, `README.md:114` | `stakeholder-answers.md:94`, `stakeholder-report.md:62` "sibling" (F4) |
| Only organize moves reports or rewrites links; counted layout; form and verb lists (REQ-707) | `completed-work-presentation-reference.md:71`, `report_index.go:21`, `tidy-repo.md:49`, `ai-report.md:34/52`, `help.md:17`, toolbox `SKILL.md:25`, `ai-report-guide.md:76/110` | none |
| Where the organize recommendation prints (REQ-708) | `version.md:50`, `ai-report-guide.md:110`, `report_organize.go:33`, the three writer actions | `README.md:49` agent upgrade prompt (F1). The condition wording is narrower than the check (F2) |
| `ai-report-organize --check` | Go usage line `[--dry-run \| --check]`. `help.md:17`, `ai-report.md:52/55` and `ai-report-guide.md:76` show `[--dry-run]` only, by design (D-10, an action-facing seam) | none |
| Report writers' closing output | `ai-report.md:205` (path line still last), `architecture-report.md:117-119` (eight lines plus one), `stakeholder-report.md:51-57` (path line still last), `ai-report-reference.md:120` | none |

`version-guide.md` and the README "Updating" paragraph do not describe update output, so they restate nothing. Neither earlier member is unread: both REQ-706 and REQ-707 recorded a **Restatement sweep:** line.

### Anti-bloat count

Everything added is in the REQ: one constant (D-03), one exported function (D-02), one flag (D-09/D-10), two five-line handler blocks (D-04/D-05), one index-handler branch (D-07), two focused tests and one absence assertion. No new file, helper, option or decorative test. The two handler blocks repeat five lines, which is too little to need a shared helper.

### Domain Review (backend.md)

- Invalid input: `--check` combined with another flag is a usage error (exit 2).
- No secrets, authentication or data exposure surface.
- Idempotency: `--check` is read-only (the test compares tree snapshots).
- Adjacent tests still pass.

### Suggested Additional Testing

- Deployment: unassessed. The installer, updater, staged-skills and do-work-cli-integrations heavy lanes were not run (user decision: focused tests only). Run them before the next release if that policy allows.
- Live acceptance: unassessed. In a consumer repository with a loose September report, run the README bootstrap installer and check for one hint line after "Installed do-work suite". Then, on the release after this one (D-06), run `just do-work-update` and check for one hint after the success line.
- Agent prose paths: in a repository with a loose finished month, run `do-work-toolbox ai-report REQ-NNN`, `architecture-report`, and a stakeholder-answers regeneration. Check that the hint sits just before the final path line and that `blocked_by:` gets the path, not the hint.
- Agent update: run `do-work update` in Claude Code on a post-release install and check that the reply relays the stderr hint with the rendered result.
- Revise: run `ai-report revise latest` and check that the reply shows the hint once.

### Scores (on the record, not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every criterion delivered. D-01's reading of "install" is sound |
| Code Quality | 92% | Small, guarded, follows the narration contract. One inaccurate comment (F3) |
| Test Adequacy | 88% | Red-green on both tests. Decline, already-current and failure paths are not pinned (verified by this review in scratch) |
| Scope | 100% | 10 declared, 10 touched |
| Risk | Low | New `suiteinstall` → `toolboxcommands` import. Hint guarded by success and a nil check |
| Acceptance | Pass | Implementation and Integration stages |

### Follow-ups created
None (5 findings report only)

## Review

**Overall: 95%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 `README.md:49`: the paste-in AI-agent upgrade prompt runs `update-suite` but asks the agent to report only versions and completed/cancelled, so the organize hint that `version.md:50` step 4 now relays is likely dropped on that path. Add "and relay the `ai-report organize` recommendation line if the command prints one". — impact-user-visible → report only. Nits: F2 `ai-report.md:195`, `architecture-report.md:119`, `stakeholder-report.md:51`, `ai-report-guide.md:110` say the hint prints "while a finished month's reports are still loose", but it also fires on a stale count, a plain `YYYY-MM` folder or a current-month month folder (reproduced). — impact-negligible → report only. F3 `suite_commands.go:110`: the comment "a declined install leaves no toolbox" is false over an existing suite. — impact-negligible → report only. F4 `stakeholder-answers.md:94`, `stakeholder-report.md:62` keep the "fresh sibling" wording REQ-706 F3 removed elsewhere. — impact-negligible → report only. F5 `report_organize.go:35`: the hint starts with `ai-reports/`, like stakeholder-report's path line. — impact-negligible → report only
**Acceptance:** Pass — Implementation and Integration stages: focused Go tests pass at `e4c9d951`. A CLI built from the merge passed 10 scratch fixtures (hint only when organize has work, exit 0, read-only, JSON plain success). Scratch tests in a module copy confirmed no hint on a declined install, exactly one on an already-current or declined update, and none with exit 2 on a failed update. Deployment and Live acceptance unassessed.
**Restatement sweep:** redefined where the organize recommendation prints (suite install, `do-work update`, `ai-report index`, every new-bundle writer), the `ai-report-organize --check` seam, the report writers' closing output (hint before the final path line, architecture-report's ninth line) and version.md's relay. Inherited and swept: REQ-706 bundle location, REQ-707 organize-only mover, counted layout and form/verb lists. Stale: `README.md:49` (F1), `stakeholder-answers.md:94` and `stakeholder-report.md:62` "sibling" (F4). Agreeing sites are listed in the review report.
**Suggested testing:** 5 items
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*

## Re-check (e4c9d951..112545ad)

**Overall: 96%**, Acceptance Pass (Implementation and Integration stages). Delta: builder commit `bf860693`, merge `112545ad`, 8 files, 9 lines changed. Cumulative range `ee789713..112545ad`: 12 files, +163/-7. Scope gained `README.md` and `stakeholder-answers.md` (D-18).

| Finding | Status | Evidence |
| --- | --- | --- |
| F1 README agent prompt | Fixed | `README.md:49` now ends "and relay the `ai-report organize` recommendation line if the command prints one" |
| F2 condition wording | Mostly fixed | The four prose sites now say "while organize has work to do (a finished month's reports still loose, or a month folder to recount or rename)". That covers stale counts and plain `YYYY-MM`. It still leaves out a current-month month folder, which organize flattens (see N1) |
| F3 install comment | Fixed | `suite_commands.go:110` "A declined install wrote nothing, so this run did not install organize". Accurate |
| F4 "sibling" wording | Fixed | `stakeholder-answers.md:94` and `stakeholder-report.md:62` now say "a new bundle directly under `ai-reports/`" |
| F5 hint prefix | Fixed | `report_organize.go:35` now starts with words: "Finished months in ai-reports/ need organizing: run do-work-toolbox ai-report organize (add --dry-run to preview)." |

Old hint text: no shipped file (`skills/`, `suite/`, `README.md`, Go, tests) still contains "needs organizing by month". The tests assert only the command name (D-11), so they need no change. The old text survives only in run records: `REQ-708 D-03` (`do-work/working/REQ-708-...md:111`) and `do-work/runs/work-2026-11-02-170509/REQ-708-exploration-draft.md:31`.

Checks run at `112545ad`:
- Focused tests pass: 8 toolboxcommands and 2 suiteinstall.
- `go vet` passes, `gofmt -l` reports nothing, and `git diff --check` is clean.
- A rebuilt CLI (`/tmp/req708-review/dwc2`, fixtures in `/tmp/req708-fx2.uIh1`) prints the new line on the loose fixture, both for `--check` and at the end of `ai-report-index`. It prints nothing on the organized fixture and exits 0 in all cases.
- The delta's only Go change is one comment and one string constant, so the earlier scratch results on decline, already-current and failed paths still hold.

New findings:
- N1 Nit `report_organize.go:35` and the four prose sites (`ai-report.md:195`, `architecture-report.md:119`, `stakeholder-report.md:51`, `ai-report-guide.md:110`): the new text says "Finished months in ai-reports/ need organizing". It also prints when the only work is flattening a current-month month folder (fixture `ai-reports/2026-11/2026-11-01_0900_cur`, reproduced). That is the layout the inbox note's consumer script made, so in that repository the line names finished months that do not exist. The prose list of cases leaves this case out too. Suggested fix: hint "Reports in ai-reports/ need organizing by month: …", and prose "while organize has work to do (for example …)". — impact-negligible → report only
- N2 Nit `do-work/working/REQ-708-recommend-ai-report-organize-on-install-update.md:111`: D-03 still quotes the old hint text as "the same everywhere", with no note that review F5 changed it. Suggested fix: append "(text changed after review F5, D-18)". — impact-negligible → report only

**Re-check record:** overall 96% (Requirements 100%, Code Quality 95%, Test Adequacy 88%, Scope 100%, Risk Low, Acceptance Pass). Status: F1, F3, F4 and F5 fixed. F2 mostly fixed, with the remainder recorded as N1. Follow-ups created: None (2 findings report only).
