---
id: REQ-708
title: 'Install, update and new reports recommend ai-report organize when finished months are still loose'
status: completed
created_at: 2026-11-02T16:55:00Z
user_request: UR-161
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-706, REQ-707]
batch: ai-report-month-folders
depends_on: [REQ-707]
required_lessons: ["skills/do-work/tools/lessons-do-work-update.md", "_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go", "skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go", "skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/actions/architecture-report.md", "skills/do-work-toolbox/actions/stakeholder-report.md", "skills/do-work/actions/version.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "README.md", "skills/do-work/actions/stakeholder-answers.md"]
claimed_at: 2026-11-02T18:40:42Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-11-02T18:41:01Z
  basis:
    - Route B
    - 10-file write set
    - 3 subsystems involved
    - 6 acceptance criteria
dispatch_at: 2026-11-02T18:44:55Z
builder_handback_at: 2026-11-02T18:53:40Z
integration_at: 2026-11-02T18:53:46Z
review_at: 2026-11-02T19:05:01Z
kb_status: pending
commit: 112545ad050f73d9e6c0b6a6dfc5e7721554a21b
completed_at: 2026-11-02T19:08:02Z
release_at: 2026-11-02T19:08:02Z
---
# Install, Update and New Reports Recommend ai-report Organize When Finished Months Are Still Loose
## What
When a project's `ai-reports/` still holds reports from a finished month at the top level, or a month folder without a count, `do-work-toolbox install` and `do-work update` print one line recommending `do-work-toolbox ai-report organize`. `ai-report index` and every new report print the same one-line hint. Nothing prints when there is nothing to organize.
## Why
The user wants the installed repository to be told to run the organize command, and each new month leaves last month's reports loose until someone runs it again.
## Detailed Requirements
- One shared check decides whether there is anything to organize, using REQ-707's finished-month rule. It is a program, not prose repeated in each action.
- `do-work-toolbox install` and `do-work update` print the recommendation only when the check finds something.
- `ai-report index` and every new report print the one-line hint under the same condition: ai-report in every form that writes a bundle, architecture-report, and stakeholder-report.
- None of them moves a report or rewrites a link. The user runs organize.
## Constraints
- The current month's reports stay directly under `ai-reports/`. Only finished months move into a month folder, because the count is final only when the month is over. Report writers (ai-report in every form, architecture-report, stakeholder-report) keep writing new bundles directly under `ai-reports/`: no write-side change.
- A finished-month folder is named `YYYY-MM_N-reports` (`YYYY-MM_1-report` for one report). A plain `YYYY-MM` folder, as a consumer's own script makes, is read as a month folder too.
- The bundle contract is unchanged: one folder holding its entry file, the same catalog fields, and the same `yyyy-mm-dd_hhmm_<slug>` names. A month folder is only a parent directory.
- A repository with no month folders behaves exactly as before.
- Only `ai-report organize` moves reports or rewrites links. A report run never does.
- The inbox note was observed against 0.307.0, so its line numbers may have drifted.
## Dependencies
Depends on REQ-707 (ai-report organize), which the recommendation names.
## Builder Guidance
User answers: "Print a hint, you run it" for the monthly trigger, and "Install and update" for where the recommendation appears.
## Red-Green Proof
**RED prompt/case:** A fixture repository with a loose finished-month report, `ai-reports/2026-09-27_2200_a/index.html`. Run the install recommendation path, `do-work update`, and `ai-report index`.
**Why RED now:** No command checks the reports layout, so none of them mentions organize.
**GREEN when:** Each prints one line naming `do-work-toolbox ai-report organize`. In a fixture holding only current-month reports and counted month folders, none prints it.
**Validation:** User confirmed the triggers and where it prints. The fixture shape was inferred during capture.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (8398 tokens, over the 2000 budget; `slugged: partial`). Matching reason: family `alternate-writer-contract-drift`: the hint must reach every report writer (ai-report forms, architecture-report, stakeholder-report) and the install and update paths.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (23234 tokens per `do-work/lessons-index.md`, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its row covers changing do-work-cli internals; families `cross-module-fact-handoff, closed-enumeration-for-a-condition` match this work.
## Full Context
See `do-work/user-requests/UR-161/input.md` for the complete verbatim input: the user's words, the six answers chosen during capture, and the inbox note it adopts in full. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Read the brief, REQ D-01 to D-10, crew rules, primes (action-files, shell-commands, releases, do-work-cli, do-work-update) and the required lessons. Tests first through the real handlers so RED needs no new names; then one constant plus one exported function over REQ-707's planner, the `--check` branch, the index handler append, two narration blocks in the suite handlers, and the five prose edits.
- [x] **[APPLY]:** (from the builder hand-back) Implemented as planned in the ten Scope files, one commit `ecc9e0b5`; the planner runs only after a success outcome in every caller.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff 202287a3 --stat`: 10 files, 160 insertions, 4 deletions. Probe 0, gofmt clean, `go vet` 0 on both packages, `contract-regressions.sh` 0, `shipped-package-reference-contract.sh` 0, `action-shell-blocks.sh` 0 (86 blocks, ShellCheck), `prescribed-shell-canonicalization.sh` 0, `suite-manifest-contract.sh` 0, `git diff --check` 0. All ten diffs read; no debug artifacts. Built binary: `--check` prints the hint, `--format json` is a plain success, the index appends the hint, `--bogus` gives the new usage, a symlinked `ai-reports` makes `--check` print nothing with exit 0.
*Source: "So actually, when this is installed, It should recommend to run on AI report organize a script that will do it. For the repository. Where this is installed in."*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fixed by the user's answers (print a one-line hint, the user runs organize; install and update recommend it) and REQ-707 already exposes the check (`planReportOrganizeLayout` + `hasWork()`). What needs discovery is where install, update, the catalog and each report writer finish their output; exploration pins those lines, so no separate plan is needed.

**Planning:** Not required

**Open questions (Step 3.5):** none in the REQ. One wording question is settled by judgment and recorded as D-01 below: the REQ names `do-work-toolbox install`, but that toolbox action installs companion skills only; the suite itself is installed by the bootstrap installer (`install-suite`), which is what the user meant ("when this is installed ... for the repository where this is installed in").

**Test policy for this run (user decision, 2026-11-02):** focused tests only. The gate argv for every evidence step is `bash do-work/runs/work-2026-11-02-170509/REQ-708-gate.sh`; `_dev/tests/maintainer-verify.sh` and heavy lanes are not run.

**Prior work this builds on:** REQ-706 (0.309.1, month folders read as containers) and REQ-707 (0.309.2, `ai-report organize`; its plan decision D-06 escalated to this REQ how install and update reach the check). This REQ closes UR-161 and is the run's last integration, so its review runs the wave-end restatement sweep over REQ-706 and REQ-707.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Lessons consult:** `required_lessons` refreshed to `skills/do-work/tools/lessons-do-work-update.md` (628 tokens, captured, read: updater transaction and managed-section lessons; none constrains a narration line after success) plus `_dev/primes/lessons-releases.md` (666 tokens, read; `canonical-link-outlives-its-target` matters at this REQ's UR-closing finalization, which moves REQ-706 and REQ-707 into `do-work/archive/UR-161/` and so must repoint their lesson links). The two capture-time drops stay dropped for budget; the orchestrator read their matching families by hand: `cross-module-fact-handoff` (one Go module computing a fact for another: the check stays one function) and `closed-enumeration-for-a-condition` (the writers that print the hint are keyed on "every form that writes a new bundle", not a hand list).

**Exploration source:** a background exploration agent read the code at main `29e6bfa7` while REQ-707 was building (draft: `do-work/runs/work-2026-11-02-170509/REQ-708-exploration-draft.md`); the orchestrator checked it against REQ-707 as released (0.309.2, `5354a142`): `planReportOrganizeLayout(repositoryRoot, currentMonth string) (reportOrganizePlan, error)` and `(reportOrganizePlan).hasWork() bool` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go`, both unexported, Git-free.

**Key code (paths under `skills/do-work/tools/do-work-cli/` unless rooted):**
- Install: `README.md` bootstrap → `tools/install-do-work-suite.sh` (copy at `skills/do-work/tools/install-do-work-suite.sh`) → `do-work-cli.sh --format text install-suite` → `internal/suiteinstall/suite_commands.go` `handleInstallSuite`. Update: the agent's `do-work update` (`skills/do-work/actions/version.md` step 4, `--format json`), `just do-work-update` / `run-do-work-update` (`skills/do-work-board/justfile.template`), and `skills/do-work/tools/do-work-update.sh` all reach `update-suite` → `handleUpdateSuite`. The toolbox `install` action (`skills/do-work-toolbox/actions/install.md`) installs companion skills only.
- Both handlers receive a `narration io.Writer` (stderr); stdout carries only the rendered result. The success lines are narration: `internal/suiteinstall/install_transaction.go` "Installed do-work suite v...", `update_transaction.go` "You're up to date" / "Updated to v...". `RunUpdate` calls `RunInstall` in-process, so the hint belongs in the handlers, not in `RunInstall`.
- Imports: `suiteinstall` imports neither `toolboxcommands` nor anything that imports it; `toolboxcommands` does not import `suiteinstall`: no cycle. Precedent for one command package importing another: `internal/doctor` is imported by `corehelpers`, `hookcommands` and `runstatus`.
- Catalog: `internal/toolboxcommands/report_index.go` `handleAIReportIndex` (no-arg branch → `reportIndexWrite`). Organize's own catalog step calls `reportIndexWrite` directly.
- Clock: `internal/toolboxcommands/note.go` `toolboxNow`.
- Report writers and their last-output rules: `skills/do-work-toolbox/actions/ai-report.md` Step 8 and Output Format (the bundle path line must be the last line; covers the default form, `--kind` and revise); `skills/do-work-toolbox/actions/architecture-report.md` Step 7 ("print at most eight lines"); `skills/do-work-toolbox/actions/stakeholder-report.md` Step 6 ("print `ai-reports/<slug>/index.html` and stop"; its callers take that last line). No other bundle writer exists (`ai-reports/` grep over `skills/**/*.md`, `reportsFolder` in Go).
- Tests: `internal/suiteinstall/suite_commands_test.go` (`runCommand` helper, `TestInstallNarrationNeverReachesStandardOutput`), `update_transaction_test.go` (local archive server); `_dev/tests/install-suite-behavior.sh` (about 31 s) and `_dev/tests/update-script-behavior.sh` (about 68 s) cover the launchers.

**Decisions:**
- **D-01** (DECIDE & STATE): "install" in the REQ means the suite install (`install-suite`, reached by the bootstrap installer), not the toolbox companion `install` action. Reason: the user's words ("when this is installed ... for the repository where this is installed in"); the companion action never installs the suite.
- **D-02** (DECIDE & STATE): one exported function in `toolboxcommands`, `ReportOrganizeHint(repositoryRoot string) string`, returns the hint line when `planReportOrganizeLayout(root, toolboxNow().Format("2006-01")).hasWork()` holds, else "". `suiteinstall` imports `toolboxcommands` and calls it (no cycle; `doctor` precedent; the binary already links both). Rejected: shell-launcher calls (the agent and Just update paths skip the launchers, five edits incl. every consumer's managed justfile) and a new leaf package (moves REQ-707's fresh planner, no second consumer).
- **D-03** (DECIDE & STATE): the hint text is one constant, the same everywhere: `ai-reports/ needs organizing by month: run do-work-toolbox ai-report organize (add --dry-run to preview).` (changed after review F5 to `Finished months in ai-reports/ need organizing: run do-work-toolbox ai-report organize (add --dry-run to preview).`, so the line never starts with a path) A check error (for example a symlinked `ai-reports`) prints nothing and never changes an outcome or exit status.
- **D-04** (DECIDE & STATE): install and update print the hint as one narration (stderr) line right after their success line. `version.md` step 4 gains one clause telling the agent to relay that line when the engine prints it.
- **D-05** (DECIDE & STATE): install prints on success when the install was not cancelled; update prints on every success (updated, already current, declined), because the installed toolbox then has organize.
- **D-06** (DECIDE & STATE, known limit): the first `do-work update` into this release prints nothing, because the old installed binary runs `update-suite`. Later updates and every fresh install print. No variant avoids it.
- **D-07** (DECIDE & STATE): the catalog hint is added in `handleAIReportIndex` on success, not in `reportIndexWrite`, so organize's own catalog step and `find` never print it.
- **D-08** (DECIDE & STATE): the report writers that print the hint are every form that writes a new bundle: ai-report default, `--kind` and revise, architecture-report, stakeholder-report. Second-opinion (a page inside an existing bundle), judge, find and organize do not.
- **D-09** (DECIDE & STATE): writers call `ai-report-organize --check` (prints the hint line or nothing, writes nothing, reads only `ai-reports/`) and print its line just before their last line, so callers that take the last line still get the bundle path; for architecture-report it is an allowed ninth line.
- **D-10** (DECIDE & STATE): `--check` is an action-facing seam: text only, not in `help.md`, plain success in JSON.

*Generated by the exploration agent, verified by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go` (modify) — hint constant, exported hint function, the check flag
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modify) — the catalog handler appends the hint
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go` (modify) — one focused hint test
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (modify) — install and update print the hint
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modify) — one focused install/update test
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — Step 8 prints the hint before the path line
- `skills/do-work-toolbox/actions/architecture-report.md` (modify) — Step 7 prints the hint
- `skills/do-work-toolbox/actions/stakeholder-report.md` (modify) — Step 6 prints the hint before the path line
- `skills/do-work/actions/version.md` (modify) — relay the engine's hint line
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — one sentence on when the hint appears
- `README.md` (modify, added after review F1, D-18) — the upgrade prompt relays the organize line
- `skills/do-work/actions/stakeholder-answers.md` (modify, added after review F4, D-18) — "fresh sibling" bundle wording

**Files I will NOT touch:** the shell launchers (`install-do-work-suite.sh` copies, `do-work-update.sh`), `skills/do-work-board/justfile.template`, `skills/do-work-toolbox/actions/help.md`, the toolbox companion install action, and the release files (finalization writes them).

**Acceptance criteria (restated from REQ):**
- [ ] One shared program check decides whether there is anything to organize, using REQ-707's finished-month rule.
- [ ] Suite install and `do-work update` print the recommendation only when the check finds something.
- [ ] `ai-report index` and every new report (ai-report forms that write a bundle, architecture-report, stakeholder-report) print the same one-line hint under the same condition.
- [ ] None of them moves a report or rewrites a link.
- [ ] RED/GREEN: with a loose finished-month report, install, update and index each print one line naming `do-work-toolbox ai-report organize`; with only current-month reports and counted month folders, none prints it.

## Pre-Flight

**Git:** ✓ clean tracked tree at the REQ-708 claim commit (after REQ-707's release 0.309.2); only this REQ's orchestrator sections and the run directory's new REQ-708 files were dirty. The user's untracked inbox note stays hidden from discovery by the run's temporary `.git/info/exclude` line.
**Tests baseline:** ✓ `bash do-work/runs/work-2026-11-02-170509/REQ-708-preflight-probe.sh` exit 0, launched (`do-work/working/baseline.json`). Focused repository gate `bash do-work/runs/work-2026-11-02-170509/REQ-708-gate.sh` exit 0 (2026-11-02T18:42:58Z to 18:43:53Z, 55 s, 1-minute load 5.0); green-gate record satisfied. The first gate draft also ran `_dev/tests/install-suite-behavior.sh`, which exits 2 as heavy-only; it was dropped under the user's focused-tests decision, and the Go `suiteinstall` package tests stand in for it. The GREEN probe `REQ-708-probe.sh` fails at base, as it must.
**Dependencies:** ✓ Go toolchain and ShellCheck present (the gate ran both)

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/actions/architecture-report.md` (modified)
- `skills/do-work-toolbox/actions/stakeholder-report.md` (modified)
- `skills/do-work/actions/version.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `README.md` (modified, review fix F1)
- `skills/do-work/actions/stakeholder-answers.md` (modified, review fix F4)

**What was done:** One constant hint line and one exported `ReportOrganizeHint(repositoryRoot)` over REQ-707's Git-free layout planner decide whether anything is left to organize. `ai-report-organize --check` prints that line or nothing and writes nothing; `ai-report index` appends it after a successful catalog write; suite install (when not cancelled) and every successful update write it as one narration line after the success line. The ai-report (default, `--kind`, revise), architecture-report and stakeholder-report actions run the check and print its line just before their final path line; `version.md` step 4 relays the engine's line; the guide says when the hint appears, including the first-update limit. Nothing moves a report or rewrites a link.

**Review fixes (builder commit `bf860693`, re-merged as `112545ad` with the same `<pre>`):** the README's paste-in upgrade prompt relays the organize line (F1); the writers and the guide say the hint prints while organize has work to do, which includes a recount or rename (F2); the install comment gives the real reason a declined install prints nothing (F3); `stakeholder-answers.md` and `stakeholder-report.md` drop the "fresh sibling" bundle wording (F4); the hint now starts with words, not a path: "Finished months in ai-reports/ need organizing: run do-work-toolbox ai-report organize (add --dry-run to preview)." (F5).

## Decisions

(from the builder hand-back; D-01 to D-10 are in `## Exploration`)

- **D-11** (DECIDE & STATE): tests assert the command name and line shape, not the full hint text; the exact text lives only in the constant.
- **D-12** (DECIDE & STATE): the update half of the suite test runs a real update (0.200.0 to 0.300.0), the path where `RunUpdate` calls `RunInstall` in-process, so "exactly once" catches a hint misplaced inside `RunInstall`.
- **D-13** (DECIDE & STATE): the suite handlers guard `narration != nil`, like `handleFetchArchive` and `RunUpdate`'s `narrate`; every caller runs the planner only after a success outcome.
- **D-14** (DECIDE & STATE): the suite test keeps the REQ's 2026-09-27 fixture date; the update reads the real clock, and September 2026 is finished for any clock from October 2026.
- **D-15** (DECIDE & STATE): the toolbox test uses a plain temporary directory with no `git init`, which also proves `--check` and the index hint need no Git.
- **D-16** (DECIDE & STATE): the D-06 first-update limit fits in one clause of the guide sentence, so it is stated there.
- **D-17** (DECIDE & STATE, orchestrator): heavy lanes are not run and no heavy hold is written: the user chose focused tests only (2026-11-02). `_dev/tests/install-suite-behavior.sh` is heavy-only (exit 2) and was dropped from the focused gate; the Go `suiteinstall` package tests stand in. `## Testing` records the lanes the planner would have selected.
- **D-18** (DECIDE & STATE, orchestrator scope judgment after the wave-end review): `README.md` and `skills/do-work/actions/stakeholder-answers.md` join the Scope and `write_set`. The README's paste-in upgrade prompt restates what the agent reports after `update-suite` and would drop the organize line (F1); `stakeholder-answers.md` keeps the "fresh sibling" bundle wording REQ-706 removed elsewhere (F4).

## Discovered Tasks

(from the builder hand-back; the builder's "impact-low" is restamped to the canonical `impact-negligible`)

- **impact-negligible** The revise form runs `ai-report index`, whose output now ends with the hint, and Step 8 also prints the `--check` line, so an agent that echoes the index output could show the hint twice in one revise reply. Nothing says to echo it. → report only
- **impact-negligible** `hasWork()` is true while organize would still refuse (dirty files, a hand-made catalog), so the hint keeps printing until the user fixes the refusal organize names. Matches D-03 ("would move, rename or remove"). → report only

## Qualification

**Gate records:** `advance --diff-range ee789713..e4c9d951` → `qualify` satisfied (no debug artifact, no output-primitive finding, P-A-U boxes ticked) and `scope-drift` satisfied with no finding: ten declared paths, ten touched.

**Requirement trace (read against `git diff ee789713..e4c9d951` and the merged files):**
- One shared program check: `ReportOrganizeHint` in `report_organize.go` wraps `planReportOrganizeLayout(root, toolboxNow().Format("2006-01")).hasWork()` (REQ-707's finished-month rule); every caller uses it, and the report writers reach it through `ai-report-organize --check`.
- Install and update only when the check finds something: `handleInstallSuite` (success and not cancelled) and `handleUpdateSuite` (any success) write the line to narration only when `ReportOrganizeHint` returns it. Pinned by `TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose` (exactly once in install narration, once in a real update, never on stdout) and the absence assertion in `TestInstallNarrationNeverReachesStandardOutput`.
- `ai-report index` and every new report: `handleAIReportIndex` appends the hint after a successful write (not in `reportIndexWrite`, so organize's own catalog step and `find` stay silent); `ai-report.md` Step 8 (default, `--kind`, revise), `architecture-report.md` Step 7 and `stakeholder-report.md` Step 6 run the check and print its line before their final path line, so callers that take the last line still get the path. Pinned by `TestReportOrganizeHintPrintsOnlyWhenWorkRemains` (fixture A prints, fixture B with only current-month reports and a counted month folder prints nothing; `--check` writes nothing) and by the probe's prose checks.
- Nothing moves or rewrites: `--check` and the hint run only the layout planner, which reads `ai-reports/` and writes nothing; the toolbox test runs without Git.
- `version.md` step 4 relays the engine's line for the agent's `do-work update` path (JSON output plus narration); the guide states when the hint appears and the first-update limit (D-06).

**Scope comparison:** declared 10, touched 10, no undeclared touch; the shell launchers, `justfile.template`, `README.md`, `help.md` and the companion install action are untouched.

**Judgment:** substantive and small (160 insertions, 4 deletions; 51 production lines), matches D-01 to D-17. Qualified.

## Testing

**Tests run:** focused repository gate `bash do-work/runs/work-2026-11-02-170509/REQ-708-gate.sh` at merge `e4c9d951` through `run-timed-command` (2026-11-02T18:54:28Z to 18:55:16Z, 48 s, 1-minute load 3.5): `go vet` and `go test -count=1` on `./internal/toolboxcommands/` (21.1 s) and `./internal/suiteinstall/` (6.4 s), then `_dev/tests/shipped-package-reference-contract.sh`, `contract-regressions.sh`, `action-shell-blocks.sh`, `prescribed-shell-canonicalization.sh`, `suite-manifest-contract.sh`. Then `advance --gate-arg bash --gate-arg do-work/runs/work-2026-11-02-170509/REQ-708-gate.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-11-02-170509/REQ-708-probe.sh`: green-gate, scope-drift and `run-blocked-check` satisfied (probe status 0: the three tests PASS, all three report writers run `ai-report-organize --check`, `version.md` names the recommendation).
**Result:** ✓ All passing (gate exit 0 on the first run; both packages pass in full; builder's runs: `suiteinstall` 31 PASS, 3 SKIP; toolbox `AIReport|Organize|Handlers` 10 PASS, 3 SKIP; every skip is a pre-existing heavy-only test)

**Red-green validation:**
- `TestReportOrganizeHintPrintsOnlyWhenWorkRemains` (`report_organize_test.go`): ✗ before (`--check` was a usage error; `--check on a loose finished month must print one line naming do-work-toolbox ai-report organize, got ""`; the index output did not end with the hint) → ✓ after. This is the `## Red-Green Proof` index case and its "none prints it" fixture (current-month report plus a counted month folder).
- `TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose` (`suite_commands_test.go`): ✗ before (`install narration names "do-work-toolbox ai-report organize" 0 times, want 1`; the same for update) → ✓ after. This is the Red-Green Proof install and update case with the loose `ai-reports/2026-09-27_2200_a/index.html`.
- `TestInstallNarrationNeverReachesStandardOutput`: the new absence assertion (no `ai-reports/` → no hint) passes before and after, by design.

**Repository gate after the review fixes:** the same focused gate argv at `112545ad` through `run-timed-command` (2026-11-02T19:06:19Z to 19:07:06Z, 47 s, 1-minute load 6.0) exit 0, recorded with `record-green-gate`; `REQ-708-probe.sh` exits 0 there. The fixes are wording, one comment and the hint constant; the tests assert the command name, not the text (D-11), and still pass.

**New tests added:**
- `TestReportOrganizeHintPrintsOnlyWhenWorkRemains` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go`
- `TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose` in `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go`

**Existing tests updated (cross-REQ impact):**
- `suite_commands_test.go` `TestInstallNarrationNeverReachesStandardOutput`: one absence assertion for the hint.

**Heavy verification plan:** heavy lanes not run: user chose focused tests only (2026-11-02). Re-planned on the cumulative range after the review fixes; the planner would have selected (D-17):
- Range: ee789713b97b853a342466d87786132f4dbaae48..112545ad050f73d9e6c0b6a6dfc5e7721554a21b
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go matched subtree skills/do-work/tools/do-work-cli (and 4 more paths)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work-toolbox/actions/ai-report.md matched subtree skills (and 10 more paths)
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — skills/do-work/actions/version.md matched exact path skills/do-work/actions/version.md (and 5 more paths)
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — README.md matched exact path README.md (and 5 more paths)

*Verified by work action*

## Review

**Overall: 95%** | 2026-11-02T19:05:01Z

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

**Re-check (`e4c9d951..112545ad`, recorded 2026-11-02T19:07:38Z):** F1, F3, F4, F5 fixed and F2 mostly fixed; overall 96%, acceptance Pass; no shipped file still quotes the old hint text. New nits: N1 the hint ("Finished months in ai-reports/ ...") and the writers' prose also fire when the only work is flattening a current-month folder such as `ai-reports/2026-11/` (the remainder of F2) — impact-negligible → report only. N2 D-03 lacked a note that F5 changed the text — fixed in this REQ by the orchestrator.

*Reviewed by review-work action*

**Orchestrator disposition:** a Pass at 95% needs no remediation, but F1 to F5 are one-clause fixes and this is the run's last member (it carries the wave-end sweep), so they are fixed now on the builder branch (scope extended, D-18), re-merged with the same `<pre>`, re-gated and re-checked by the same reviewer.

## Lessons Learned

**What worked:** Tests that go through the real command handlers needed no new names to compile, so RED was an assertion failure, and running a real update (0.200.0 to 0.300.0) in the suite test is the path where `RunUpdate` calls `RunInstall` in-process, which is what pins "exactly once". Putting the check behind one exported function and one `--check` flag kept every caller (index, install, update, three prose writers) a one-line call.
**What didn't:** The first hint text started with `ai-reports/`, the same prefix as stakeholder-report's final path line that callers read; the wave-end review caught it (F5). The first prose said "while a finished month's reports are still loose", which undersells when `hasWork()` is true (recounts, renames, and flattening a current-month folder; N1 still open as report only).
**Worth knowing:** The first `do-work update` into 0.309.3 prints nothing, because the old installed binary runs the update (D-06). A hint inside `RunInstall` would print twice on update; it lives in the two handlers. The agent's `do-work update` path renders JSON, so `version.md` step 4 and the README upgrade prompt both have to say "relay the line".

## Orientation

Now install, update, `ai-report index` and every new report tell you, in one line, when `ai-reports/` has finished months to organize; nothing runs organize for you. Lives in the toolbox report tooling (`toolboxcommands.ReportOrganizeHint`, `ai-report-organize --check`) and the suite installer/updater (`suiteinstall` handlers). Prime spot-check: `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md` and `_dev/primes/prime-releases.md` paths still exist and none describes the report layout or the updater's narration; not stale.

## Timing

Observed 2026-11-02T18:40:42Z to 2026-11-02T19:07:39Z: 26m 57s total, 25m 33s attributed across 7 events, 1m 24s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 10m 53s | 2 |
| builder-work | 8m 45s | 1 |
| exploration-preflight | 4m 13s | 1 |
| verification-gate | 1m 36s | 2 |
| handback-merge | 6s | 1 |

Slowest stage: review / independent review with wave-end sweep, 9m 32s, outcome success.
Slowest command: verification-gate / focused repository gate, 48s, exit 0, bash (2 argv tokens).
