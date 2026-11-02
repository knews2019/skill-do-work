---
id: REQ-708
title: 'Install, update and new reports recommend ai-report organize when finished months are still loose'
status: claimed
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
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go", "skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go", "skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/actions/architecture-report.md", "skills/do-work-toolbox/actions/stakeholder-report.md", "skills/do-work/actions/version.md", "skills/do-work-toolbox/docs/ai-report-guide.md"]
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
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
- **D-03** (DECIDE & STATE): the hint text is one constant, the same everywhere: `ai-reports/ needs organizing by month: run do-work-toolbox ai-report organize (add --dry-run to preview).` A check error (for example a symlinked `ai-reports`) prints nothing and never changes an outcome or exit status.
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

**Files I will NOT touch:** the shell launchers (`install-do-work-suite.sh` copies, `do-work-update.sh`), `skills/do-work-board/justfile.template`, `README.md`, `skills/do-work-toolbox/actions/help.md`, the toolbox companion install action, and the release files (finalization writes them).

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
