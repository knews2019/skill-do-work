# REQ-708 Exploration Draft

Base read: main `29e6bfa7`; REQ-707 `## Plan` (its escalation is labelled D-06 in the plan, "D-12" in its D-14: same item).

## Key code (file:line)

Paths under `skills/do-work/tools/do-work-cli/` unless rooted.

- Install: README.md:9-24 bootstrap -> `tools/install-do-work-suite.sh:72-81` (same bytes as the `skills/do-work/tools/` copy) -> `do-work-cli.sh --format text install-suite` (`skills/do-work/tools/do-work-cli-preamble.sh:24`). No Just recipe installs. Toolbox `install` (`skills/do-work-toolbox/actions/install.md:3`) installs companion skills only.
- Update, all `update-suite`: agent `do-work update` -> `skills/do-work/actions/version.md:48` (`--format json`); `just do-work-update`/`run-do-work-update` -> `skills/do-work-board/justfile.template:63-69`; `skills/do-work/tools/do-work-update.sh:32`. Agent and Just paths bypass the shell launchers.
- `internal/suiteinstall/suite_commands.go:96` `handleInstallSuite`, `:158` `handleUpdateSuite` get `narration` (stderr, `cmd/do-work-cli/main.go:30`); stdout = result only (`:31-33`).
- Success lines are narration: `install_transaction.go:179`; `update_transaction.go:100` up to date, `:124` cancelled, `:154` updated. `RunUpdate` calls `RunInstall` in-process (`:116`), so the hint belongs in the handlers, not `RunInstall`.
- Text renderer (`internal/resultmodel/result_model.go:1162-1418`): findings print before changes, as 2+ lines.
- Imports: only `cmd/do-work-cli/main.go:22-23` imports `suiteinstall`/`toolboxcommands`; `toolboxcommands` imports atomicfile, commandruntime, gittransaction, ownedprocess, requestmodel, resultmodel, rootedfs (no `suiteinstall`): no cycle. Precedent: command package `doctor` (`internal/doctor/doctor_commands.go:21`) is imported by `corehelpers/commands.go:20`, `hookcommands/session_start.go:14`, `runstatus/run_status.go:30`.
- `internal/toolboxcommands/report_index.go:77-80` `handleAIReportIndex` -> `:264` `reportIndexWrite` -> `:312` `exactOutputResult` (`commands.go:84`). REQ-707 apply calls `reportIndexWrite` directly. An empty `ExactTextOutput` prints nothing (`result_model.go:1166`).
- Clock: `internal/toolboxcommands/note.go:14` `toolboxNow` (freeze pattern `note_test.go:36`).

## Prose locations

- `skills/do-work-toolbox/actions/ai-report.md:180` Step 8; `:186` path line is last; covers default, `--kind`, revise (`:184`); revise also runs index (`:59`).
- `skills/do-work-toolbox/actions/architecture-report.md:117` Step 7 "Print at most eight lines".
- `skills/do-work-toolbox/actions/stakeholder-report.md:51` Step 6 "Print `ai-reports/<slug>/index.html` and stop" (callers: `skills/do-work/actions/stakeholder-answers.md:65`, work.md Step 8).
- `skills/do-work/actions/version.md:50` step 4 "Render the typed result".
- Grep of `ai-reports/` in `skills/**/*.md` and `reportsFolder` in Go: no other bundle writer.

## Proposed design

One exported function in `toolboxcommands`, three thin callers:

```go
const reportOrganizeHintLine = "ai-reports/ needs organizing by month: run do-work-toolbox ai-report organize (add --dry-run to preview)."
// ReportOrganizeHint returns reportOrganizeHintLine when the layout plan has work, else "".
func ReportOrganizeHint(repositoryRoot string) string // planReportOrganizeLayout(root, toolboxNow().Format("2006-01")); error -> ""
```

1. `ai-report-organize --check`: layout plan only (no link plan, no Git, no dirty check); `exactOutputResult(hint+"\n" or "")`, writes nothing.
2. `handleAIReportIndex` no-arg branch: on success, append the hint line to `ExactTextOutput`.
3. `handleInstallSuite` / `handleUpdateSuite`: after the transaction, write the hint as one narration line.
4. Writers: one prescribed fenced line each, output printed just before the action's last line.

Decisions:
- **D-01 (DECIDE & STATE)** Seam = option (a): `suiteinstall` imports `toolboxcommands` and calls `ReportOrganizeHint`. Reasons: no cycle; `doctor` is the precedent; ~6 lines; the binary already links both. Rejected (b) shell `--check` after launchers: agent and Just update paths skip `do-work-update.sh`, so 5 edits (2 installer copies, updater, 2 Just recipes = every consumer's managed justfile) plus exit plumbing. Rejected (c) leaf package: moves REQ-707's fresh planner and `reportMonthFolder` (3 users), churn, no second consumer. `--check` stays, for the prose writers only.
- **D-02** Hint text is the single const above. It is the same everywhere and prints nothing when `hasWork()` is false. A check error (for example a symlinked `ai-reports`) prints nothing and never changes an outcome or exit code (one broken pipe).
- **D-03** Install/update print via narration (stderr), right after the success line. It is one line and the exact text. An info finding is 2+ lines mid-result. version.md step 4 gains one clause: "relay the `ai-report organize` line when the engine prints one" (the JSON path sees stderr, but "render the typed result" could drop it).
- **D-04** Install prints on success unless `Cancelled` (a declined fresh install has no toolbox). Update prints on every success, including up to date and declined: the running binary is the installed one, so organize exists.
- **D-05 (known limit)** The first update into this release prints nothing: the old installed binary runs `update-suite` (`skills/do-work/tools/prime-do-work-update.md` Traps: the running binary is the one the install replaces). Later updates and every bootstrap install print. Option (b) has the same limit.
- **D-06** Index hint lives in the handler, not `reportIndexWrite`, so organize and `--find` never print it; success only.
- **D-07** Forms: ai-report default, `--kind`, revise (Step 8). Excluded: second-opinion writes a page into an existing bundle, not a new bundle (REQ: "every form that writes a bundle"; Output Format `:184` already splits it). Also excluded: judge, find, organize itself.
- **D-08** Hint goes just before each writer's last line (path lines stay last); architecture-report allows it as a ninth line.
- **D-09** `--check` is a text-only, action-facing seam. It is not in help.md, and its JSON is a plain success (YAGNI).
- **D-10** The REQ's `do-work-toolbox install` means the suite install (`install-suite`). The companion `install` action is unchanged. Reason: the user's words "when this is installed ... for the repository where this is installed".

Prescribed line (same shape as `ai-report.md:45`, which passes `_dev/tests/action-shell-blocks.sh`, whose placeholder neutralizing is at `:194-195`):

```bash
<skill-root>/../do-work/tools/do-work-cli.sh --repo-root <project-root> --format text ai-report-organize --check
```

## Files to touch (write set)

1. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go`: const, `ReportOrganizeHint`, `--check` branch.
2. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`: handler appends the hint.
3. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go`: one test.
4. `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go`: import, two narration calls.
5. `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go`: one test plus one absence assert.
6. `skills/do-work-toolbox/actions/ai-report.md` (Step 8, after `:180`).
7. `skills/do-work-toolbox/actions/architecture-report.md` (Step 7 `:117`).
8. `skills/do-work-toolbox/actions/stakeholder-report.md` (Step 6 `:51`).
9. `skills/do-work/actions/version.md` (`:50`, one clause).
10. `skills/do-work-toolbox/docs/ai-report-guide.md`: one sentence in REQ-707's organizing section.

Not touched: installer/updater shell launchers, justfile.template, README, help.md, the toolbox install action.

## Tests

- **`TestReportOrganizeHintPrintsOnlyWhenWorkRemains`** (toolboxcommands; clock frozen at 2026-11-02). Fixture A, `ai-reports/2026-09-27_2200_a/index.html`: `ai-report-index` text output ends with the hint line, and `ai-report-organize --check` prints exactly that line and creates no file. Fixture B, `2026-11-01_0900_current/` plus `2026-09_1-report/2026-09-27_2200_a/`: neither prints it, and `--check` output is "". RED: `--check` is a usage error and index has no hint.
- **`TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose`** (suiteinstall, `runCommand` at `suite_commands_test.go:26`). Loose `ai-reports/2026-09-27_2200_a/index.html`: `install-suite --archive` narration names `do-work-toolbox ai-report organize` once, stdout never; then `update-suite --skill-root --upstream-url` (up-to-date, server `update_transaction_test.go:20`) names it once. Plus one absence assert in `TestInstallNarrationNeverReachesStandardOutput` (`:237`).
- Focused gate: `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^(TestReportOrganizeHint|TestInstallAndUpdateRecommendOrganize|TestInstallNarrationNeverReachesStandardOutput)' ./internal/toolboxcommands ./internal/suiteinstall`. Then full `./internal/suiteinstall` (~11 s recorded), `-run 'AIReport|Organize|Handlers' ./internal/toolboxcommands`, `gofmt -l`, `go vet`, `action-shell-blocks.sh` (~5 s), `shipped-package-reference-contract.sh` (~5 s), `staged-skills-contract.sh` (~35 s), `prescribed-shell-canonicalization.sh` (~3 s).
- Skip `install-suite-behavior.sh` (~31 s) and `update-script-behavior.sh` (~68 s): no `ai-reports/` fixture; the only pinned phrase (`install-suite-behavior.sh:272`) is unchanged.

## YAGNI

No `Cancelled` field on `UpdateResult`; no info finding or JSON field; no new verb, Just recipe or help line; no leaf package; no shell-launcher, justfile.template or README edits; no hint in second-opinion, judge, find or organize; no report count in the hint; no auto-run of organize; tests assert the command name, the exact text lives only in the const.
