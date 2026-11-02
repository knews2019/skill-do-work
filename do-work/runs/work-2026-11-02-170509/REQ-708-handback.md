# Hand-back: REQ-708 (install, update and new reports recommend ai-report organize when finished months are still loose)

- Branch: `worktree-agent-REQ-708-recommend-ai-report-organize-on-install-update`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-708-recommend-ai-report-organize-on-install-update`
- Base commit: `202287a3` ([UR-161] run artifacts)
- Commits: `ecc9e0b5` [REQ-708] install, update, ai-report index and new reports recommend ai-report organize
- Worktree clean after the commit (`git status --short` empty). Nothing staged or written under `do-work/` except this file.

## File manifest

All (modified); exactly the ten `## Scope` files.

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go`: `reportOrganizeHintLine` constant (D-03 text verbatim), exported `ReportOrganizeHint(repositoryRoot string) string` (planner + `hasWork()`, error yields ""), `--check` branch (hint line or nothing, `exactOutputResult`, no write, exit 0), usage `ai-report-organize [--dry-run | --check]`.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`: `handleAIReportIndex` no-arg branch appends the hint line on success only. `reportIndexWrite` and `--find` unchanged.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go`: `TestReportOrganizeHintPrintsOnlyWhenWorkRemains`.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go`: imports `toolboxcommands`; `handleInstallSuite` writes the hint to narration after a successful, not-cancelled install; `handleUpdateSuite` after every successful update. Outcome and exit status untouched.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go`: `TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose`, plus one absence assertion in `TestInstallNarrationNeverReachesStandardOutput`.
- `skills/do-work-toolbox/actions/ai-report.md`: Step 8 runs `ai-report-organize --check` and prints its line just before the final bundle-path line (covers default, `--kind`, revise).
- `skills/do-work-toolbox/actions/architecture-report.md`: Step 7 runs the check; its line is the one allowed line beyond the eight.
- `skills/do-work-toolbox/actions/stakeholder-report.md`: Step 6 runs the check first and prints its line before the path line, because callers take the last line.
- `skills/do-work/actions/version.md`: step 4 gains one clause to relay the engine's `ai-report organize` line. The `**Current version**:` line is untouched.
- `skills/do-work-toolbox/docs/ai-report-guide.md`: one sentence in "Organizing Reports by Month": install, `do-work update` (with the D-06 first-update limit in one clause), `ai-report index` and every new report print a one-line hint; nothing runs organize for you.

## P-A-U

**[PLAN]** Read the brief, the REQ (D-01 to D-10), the rules (general, coding-guardrails, shared-principles, communication-style, testing, backend), the primes (action-files, shell-commands, releases, do-work-cli, do-work-update) and the required lessons. Approach: tests first through the real handlers (no compile dependency on new names), then one constant plus one exported function over REQ-707's planner, the `--check` branch, the index handler append, two narration blocks in the suite handlers, and four short prose additions with the same prescribed line shape as the existing `ai-report-index` block. No lesson contradicted the plan.

**[APPLY]** Implemented exactly as planned in the ten Scope files. The planner runs only after a success outcome in every caller.

**[UNIFY]**

```
 skills/do-work-toolbox/actions/ai-report.md        |  6 +++
 .../do-work-toolbox/actions/architecture-report.md |  6 +++
 .../do-work-toolbox/actions/stakeholder-report.md  |  6 +++
 skills/do-work-toolbox/docs/ai-report-guide.md     |  2 +-
 skills/do-work/actions/version.md                  |  2 +-
 .../internal/suiteinstall/suite_commands.go        | 14 ++++++
 .../internal/suiteinstall/suite_commands_test.go   | 40 +++++++++++++++++
 .../internal/toolboxcommands/report_index.go       | 11 ++++-
 .../internal/toolboxcommands/report_organize.go    | 26 ++++++++++-
 .../toolboxcommands/report_organize_test.go        | 51 ++++++++++++++++++++++
 10 files changed, 160 insertions(+), 4 deletions(-)
```

| Check (from the worktree root) | Exit | Wall |
|---|---|---|
| `bash .../REQ-708-probe.sh` | 0 (`REQ-708 GREEN: organize recommended by install, update, index and report writers`) | 1 s |
| `gofmt -l skills/do-work/tools/do-work-cli/internal` | 0, no output | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands ./internal/suiteinstall` | 0 | 1 s |
| `bash _dev/tests/contract-regressions.sh` | 0 | 29 s |
| `bash _dev/tests/shipped-package-reference-contract.sh` | 0 | 2 s |
| `bash _dev/tests/action-shell-blocks.sh` | 0 (86 fenced blocks, 33 shell files, ShellCheck enabled) | 2 s |
| `bash _dev/tests/prescribed-shell-canonicalization.sh` | 0 | 1 s |
| `bash _dev/tests/suite-manifest-contract.sh` | 0 | 2 s |
| `git diff 202287a3 --stat` | exactly the ten Scope files | - |
| `git diff 202287a3 --check` | 0, clean | - |

Each file checked: all ten diffs read in full. No debug artifacts. Prose: each new block is the brief's exact prescribed line; ai-report's Output Format last-line rule and stakeholder-report's "callers take the last line" stay true.

End-to-end on the built binary (scratch fixture `ai-reports/2026-09-27_2200_a/index.html`, real clock): `--format text ai-report-organize --check` printed the one hint line, exit 0; `--format json` printed a plain success, exit 0 (D-10); `ai-report-index` printed `catalog.json: 1 bundles`, `index.html`, then the hint line; `--bogus` gave the new usage text, exit 2; a symlinked `ai-reports` made `--check` print nothing with exit 0 (D-03 error path); only the index wrote files.

## Proof record

**RED** (tests written before any production edit; `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestReportOrganizeHintPrintsOnlyWhenWorkRemains|TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose|TestInstallNarrationNeverReachesStandardOutput)$' ./internal/toolboxcommands ./internal/suiteinstall`, 2 s):

```
report_organize_test.go:324: ai-report-organize --check must succeed with text output: {... Outcome:failure ... Code:TOOLBOX-USAGE ... Evidence:[Usage: ai-report-organize [--dry-run]] ...}
report_organize_test.go:326: --check on a loose finished month must print one line naming do-work-toolbox ai-report organize, got ""
report_organize_test.go:332: ai-report-index output must end with the --check line "":
--- FAIL: TestReportOrganizeHintPrintsOnlyWhenWorkRemains (0.04s)
--- PASS: TestInstallNarrationNeverReachesStandardOutput (0.25s)
suite_commands_test.go:289: install narration names "do-work-toolbox ai-report organize" 0 times, want 1:
suite_commands_test.go:302: update narration names "do-work-toolbox ai-report organize" 0 times, want 1:
--- FAIL: TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose (0.45s)
```

Every failure is an assertion, not a compile error. The absence assertion passes at RED by design (no `ai-reports/`, no hint, before and after).

**GREEN** (same command, 2 s): `--- PASS` for all three, no skip.

**Package runs:** `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/suiteinstall`: ok, 31 PASS, 3 SKIP (the pre-existing heavy-only built-binary signal tests, gated by `DO_WORK_HEAVY_TESTS`), 0 FAIL, 8 s. `-run 'AIReport|Organize|Handlers' ./internal/toolboxcommands`: ok, 10 PASS, 3 SKIP (pre-existing "ai-report-judge browser cases are heavy-only"), 0 FAIL, 1 s.

## Decisions

- **D-11** (DECIDE & STATE): the tests do not repeat the hint text. The toolbox test asserts that `--check` prints exactly one line naming `do-work-toolbox ai-report organize` and that the index output ends with that same line; the suite test counts the command name. The exact text lives only in the constant, and RED compiled without the new names.
- **D-12** (DECIDE & STATE): the update half of the suite test runs a real update (installed 0.200.0, upstream 0.300.0, asserts `Updated to v0.300.0`), not the already-current path. That is the path where `RunUpdate` calls `RunInstall` in-process, so "exactly once" catches a hint misplaced inside `RunInstall`.
- **D-13** (DECIDE & STATE): the suite handlers guard `narration != nil`, matching `handleFetchArchive` and `RunUpdate`'s own `narrate` helper. Every caller (index, install, update) runs the planner only after a success outcome, so a failure reads nothing extra.
- **D-14** (DECIDE & STATE): the suite test keeps the REQ's fixture date (2026-09-27). The update reads the real clock and `toolboxNow` cannot be frozen from another package; September 2026 is finished for any clock from October 2026 onwards.
- **D-15** (DECIDE & STATE): the toolbox test uses a plain `t.TempDir()` with no `git init`, which also proves `--check` and the index hint need no Git.
- **D-16** (DECIDE & STATE): the D-06 limit fits in one clause of the guide sentence ("except the update that first brings this hint, which the older updater runs"), so it is there.

## Discovered Tasks

- impact-low: the revise form runs `ai-report index` in its step 1 and step 6, and that output now ends with the hint line, while Step 8 also prints the `--check` line. If an agent echoes the index output, the hint can appear twice in one revise reply. Nothing says to echo it, so no edit was made. → report only
- impact-low: `hasWork()` is true when organize will still refuse (dirty files, hand-made catalog), so the hint keeps printing until the user fixes the refusal that organize names. This matches D-03 ("would move, rename or remove"). → report only

## Lessons read; anti-bloat check

Read in full: `skills/do-work/tools/lessons-do-work-update.md`, `_dev/primes/lessons-releases.md`. From `lessons-do-work-cli.md` only the `[family: cross-module-fact-handoff]` and `[family: closed-enumeration-for-a-condition]` bullets. The check stays one function (cross-module-fact-handoff). The writers are keyed on "every form that writes a new bundle", and ai-report's Step 8 covers all three forms by condition (closed-enumeration-for-a-condition).

Added beyond the names in the brief: none at package level. Test-local only: the closures `indexText` and `runOrganizeCheck` and the local constant `recommendation`, which keep each assertion to one line. No new verb, Just recipe, help line, JSON field, finding, result field, package or launcher edit. No hint in second-opinion, judge, find or organize. No report count in the hint. Organize never runs automatically.

## Proposed CHANGELOG entry

**Install, Update and New Reports Recommend Organizing Finished Months**

Each new month leaves last month's reports loose under `ai-reports/` until someone runs `ai-report organize`. The suite now says when that is due, and it never runs organize for you.

- While a finished month's reports are still loose (or a month folder's count is stale), suite install, `do-work update`, `ai-report index` and every new report (ai-report in each form that writes a bundle, architecture-report, stakeholder-report) print one line: `ai-reports/ needs organizing by month: run do-work-toolbox ai-report organize (add --dry-run to preview).`
- Nothing prints when there is nothing to organize. A check error prints nothing and never changes a result or exit status.
- Install and update print the line on stderr narration, so the JSON result on stdout is unchanged. Report writers print it just before their last line, so callers that read the bundle path from the last line still get it.
- New `ai-report-organize --check` prints the line or nothing and writes nothing. It is the seam the report actions call.
- Known limit: the first `do-work update` into this release prints nothing, because the older installed updater runs it. Later updates and every fresh install print.

## Proposed lesson bullet

For `skills/do-work/tools/lessons-do-work-update.md`:

- [family: in-process-install-narration] REQ-708: `RunUpdate` runs `RunInstall` in-process, so a narration line added inside `RunInstall` prints twice on an update. Add per-command narration in `handleInstallSuite` / `handleUpdateSuite`, and pin "exactly once" with a test that runs a real update, not the already-current path.

## Addendum (review fixes)

- Commit: `bf860693` [REQ-708] review fixes F1-F5 (on top of `ecc9e0b5`, same branch and worktree). Worktree clean afterwards. `git diff 202287a3 --stat`: 12 files (the ten Scope files plus D-18's `README.md` and `skills/do-work/actions/stakeholder-answers.md`), 163 insertions, 7 deletions.
- Files (all modified):
  - F1 `README.md`: the agent upgrade prompt also relays the `ai-report organize` recommendation line if the command prints one.
  - F2 `skills/do-work-toolbox/actions/ai-report.md`, `architecture-report.md`, `stakeholder-report.md` and `skills/do-work-toolbox/docs/ai-report-guide.md` all now say the same thing: the hint prints "while organize has work to do (a finished month's reports still loose, or a month folder to recount or rename)".
  - F3 `suite_commands.go`: the comment is corrected ("A declined install wrote nothing, so this run did not install organize"). Behavior is unchanged.
  - F4 `skills/do-work/actions/stakeholder-answers.md`, `stakeholder-report.md`: "fresh sibling" becomes "a new bundle directly under `ai-reports/`".
  - F5 `report_organize.go`: the constant is now `Finished months in ai-reports/ need organizing: run do-work-toolbox ai-report organize (add --dry-run to preview).` No shipped prose quoted the old text exactly. The proposed CHANGELOG bullet above still quotes the old text, so use the new one.
- Checks (worktree root): `gofmt -l` 0, no output. `go vet` on both packages 0. The focused three-test run: all `--- PASS`, ok in both packages. `REQ-708-probe.sh` 0. `contract-regressions.sh` 0 (21 s), `shipped-package-reference-contract.sh` 0, `action-shell-blocks.sh` 0, `prescribed-shell-canonicalization.sh` 0, `suite-manifest-contract.sh` 0. `git diff --check` clean.
- New decisions: none. The two test comments that say "still loose" describe their own fixtures, which are loose reports, so they stay.
