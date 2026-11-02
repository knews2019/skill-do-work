---
id: REQ-706
title: 'Addendum: ai-report tools read report bundles inside month folders'
status: claimed
created_at: 2026-11-02T16:55:00Z
user_request: UR-161
addendum_to: REQ-655
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-707, REQ-708]
batch: ai-report-month-folders
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md"]
claimed_at: 2026-11-02T17:05:09Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-11-02T17:07:35Z
  basis:
    - Route B
    - 6-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
dispatch_at: 2026-11-02T17:11:16Z
builder_handback_at: 2026-11-02T17:18:49Z
---
# Addendum: ai-report Tools Read Report Bundles Inside Month Folders
## What
Teach the report catalog, the second-opinion bundle check, and the revise and second-opinion target rules to treat a month folder under `ai-reports/` as a container of report bundles, never as a bundle itself.
## Why
A project with hundreds of reports files them by month so the folder stays readable. Today the catalog lists each month folder as one broken bundle, and revise and second-opinion refuse any report inside one (inbox note in UR-161).
## Prior Implementation
REQ-655 (ai-report index and find build a derived catalog of every report bundle in any naming style), commit `d91312ae`, added the `ai-report-index` verb. Its `reportIndexWalk` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` reads one level of `ai-reports/` and treats every non-dot directory as a bundle. Supersedes resolution matches by `path.Base`, so it already works with nested paths when the walk returns them. The second-opinion form (REQ-698) refuses a bundle whose parent is not `ai-reports` (`report_second_opinion.go`, the `filepath.Dir(relative) != reportsFolder` check). `skills/do-work-toolbox/actions/ai-report.md` states the same direct-child rule for revise (Revise form step 1) and second-opinion.
## Detailed Requirements
- A month folder is a direct child of `ai-reports/` named exactly `YYYY-MM`, `YYYY-MM_1-report`, or `YYYY-MM_<N>-reports`. Its non-dot child directories are bundles, with `Path` = `ai-reports/<month folder>/<bundle>`. The month folder itself is never a bundle.
- Flat bundles keep working. Loose files and symbolic links are never bundles. A symlinked month folder stays refused (the existing linked-ancestor check).
- Second-opinion accepts a bundle whose parent is `ai-reports` or a month folder, with a regular `index.html`, and `latest` resolves nested bundles.
- Revise accepts `ai-reports/<folder>`, `ai-reports/<month folder>/<folder>`, or a bare `<folder>` found in either place, and resolves a supersedes chain across a flat and a nested revision.
- `ai-report find` matches nested bundles by folder name, as it does flat ones.
- The ai-report action text and the user guide stop saying a bundle must sit directly under `ai-reports/`.
## Constraints
- The current month's reports stay directly under `ai-reports/`. Only finished months move into a month folder, because the count is final only when the month is over. Report writers (ai-report in every form, architecture-report, stakeholder-report) keep writing new bundles directly under `ai-reports/`: no write-side change.
- A finished-month folder is named `YYYY-MM_N-reports` (`YYYY-MM_1-report` for one report). A plain `YYYY-MM` folder, as a consumer's own script makes, is read as a month folder too.
- The bundle contract is unchanged: one folder holding its entry file, the same catalog fields, and the same `yyyy-mm-dd_hhmm_<slug>` names. A month folder is only a parent directory.
- A repository with no month folders behaves exactly as before.
- Only `ai-report organize` moves reports or rewrites links. A report run never does.
- The inbox note was observed against 0.307.0, so its line numbers may have drifted.
## Dependencies
None. REQ-707 (ai-report organize) and REQ-708 (the organize recommendation) depend on this.
## Builder Guidance
The user asked for all three items of the inbox note's read-side direction, with the folder name carrying the report count. Keep the catalog output for a repository without month folders byte-identical.
## Red-Green Proof
**RED prompt/case:** A fixture with `ai-reports/2026-09_1-report/2026-09-27_2200_a/index.html` and `ai-reports/2026-11-01_0900_b/index.html`. Run `ai-report index`, then `ai-report second-opinion --prepare` on `ai-reports/2026-09_1-report/2026-09-27_2200_a`.
**Why RED now:** The walk reads one level, so the catalog lists `2026-09_1-report` as a bundle of kind `unknown` and misses `a`. Second-opinion refuses the nested bundle as not a direct child of `ai-reports`.
**GREEN when:** The catalog lists exactly two bundles, `b` first, with paths `ai-reports/2026-11-01_0900_b` and `ai-reports/2026-09_1-report/2026-09-27_2200_a`, and no bundle named after the month folder. The same holds for a plain `2026-09/` folder. Second-opinion passes the bundle check for the nested bundle. A symlinked month folder is still refused. A repository with no month folder gets a byte-identical catalog.
**Validation:** Inferred during capture from the inbox note's acceptance check, adjusted for the counted folder names the user chose.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (23234 tokens per `do-work/lessons-index.md`, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its row covers changing do-work-cli internals; families `catalog-tie-order, exact-basename-authority` match this work.
- `_dev/primes/lessons-action-files.md` as a whole satellite (8398 tokens, over the 2000 budget; `slugged: partial`). Matching reason: family `alternate-writer-contract-drift`: the direct-child bundle rule is restated in the revise form, the second-opinion form and the guide.
## Full Context
See `do-work/user-requests/UR-161/input.md` for the complete verbatim input: the user's words, the six answers chosen during capture, and the inbox note it adopts in full. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: "I actually want all of these." (the inbox note's read-side items 1 to 3), with "Count in the folder name".*

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fixed and the REQ's Prior Implementation names the two Go functions, but the change crosses the catalog walk, the second-opinion bundle check, two test files and two prose files, and the revise rule lives only in action prose. Exploration pins the exact lines, so no separate plan is needed.

**Planning:** Not required

**Open questions (Step 3.5):** none in the REQ. The user's design decisions are recorded in the REQ's Constraints (counted folder names, current month flat, readers only). Builder choices below start at D-01.

**Test policy for this run (user decision, 2026-11-02):** focused tests only. The repository gate argv for every evidence step is `bash do-work/runs/work-2026-11-02-170509/REQ-706-gate.sh` (go vet and go test on the changed package, plus the five contract scripts); `_dev/tests/maintainer-verify.sh` and heavy lanes are not run.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Lessons consult:** `required_lessons` refreshed to `_dev/primes/lessons-releases.md` (666 tokens, whole satellite, read; matches the REQ's `prime-releases.md` prime and the release this REQ carries). Its `canonical-link-outlives-its-target` family does not bite here (nothing moves), but it is the core risk of REQ-707. The two capture-time drops stay dropped for budget (section above). The orchestrator read their matching family bullets by hand: `catalog-tie-order` (REQ-656: the folder-name date is minute-precise, ties sort by path ascending; the walk's sort must stay the same for nested paths) and `alternate-writer-contract-drift` (sweep every restatement of the direct-child rule, not only the cited lines).

**Key code (main at d231e2fa, 0.309.0):**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go:79-117`, `reportIndexWalk`: one `os.ReadDir` of `ai-reports/`, every non-dot directory becomes `readReportBundle(repositoryRoot, entry.Name())`. `DirEntry.IsDir()` is false for a symbolic link, so a symlinked folder is already skipped. Sort: date descending, then path ascending. Supersedes resolution keys on `path.Base(bundle.Path)`, so nested paths resolve once the walk returns them.
- `report_index.go:119-165`, `readReportBundle(repositoryRoot, folder)`: builds `Path = ai-reports/<folder>` and reads `filepath.Join(repositoryRoot, reportsFolder, folder, ...)`. The date comes from `reportDatePattern` on the folder name and `reportKind(folder, ...)` and the linked-id scan read the folder name: all three must keep reading the bundle's own folder name, not the month prefix.
- `report_index.go:307-314`, `reportPageLink`: trims `ai-reports/` from the entry, so a nested entry becomes `<month folder>/<bundle>/index.html`, a correct link from `ai-reports/index.html`. No change needed.
- `report_index.go:376-407`, `reportIndexFind`: matches `path.Base(bundle.Path)`, so nested bundles match by folder name with no change once the walk returns them.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:374-426`, `secondOpinionBundle`: `latest` is `catalog[0].Path`; the check at :387 `filepath.Dir(relative) != reportsFolder` refuses a nested bundle. `validateNoLinkedAncestors` (:390, :420) refuses any path through a symbolic link, so a symlinked month folder stays refused. Every entry path (plain bundle, `latest`, page target at :174, `--finalize` staging page at :162) goes through this one function, so one change covers all forms (`guard-on-one-entry-path`).
- Tests: `report_index_test.go` (`TestAIReportIndexCatalogsEveryBundleNamingStyleOnce` at :100 and its fixture helpers), `report_second_opinion_test.go` (`secondOpinionFixture` at :21 builds `ai-reports/2026-10-15_1200-example`; `TestSecondOpinionLatestAndLinkedBundleRefusal` at :339 covers `latest` and the linked-ancestor refusal).

**Prose that restates the direct-child rule (the whole sweep, `git grep` over `skills/`):**
- `skills/do-work-toolbox/actions/ai-report.md:54` (revise step 1): "must be a directory directly under `ai-reports/`".
- `skills/do-work-toolbox/actions/ai-report.md:72` (second-opinion): "It must be a physical direct child of `ai-reports/`".
- `skills/do-work-toolbox/actions/ai-report.md:42` and `skills/do-work-toolbox/docs/ai-report-guide.md:95`: "catalogs every bundle under `ai-reports/`" / "lists every directory under `ai-reports/` as a report bundle": true today, false once month folders are containers.
- `skills/do-work-toolbox/docs/ai-report-guide.md:103` (revise): "writes a new sibling report": the new revision is written directly under `ai-reports/` even when the prior one sits in a month folder.
- Writers (`ai-report.md:101`, `ai-report-reference.md:120`, `architecture-report.md:121`, `stakeholder-report.md:23`) write `ai-reports/<slug>/`: unchanged by design (REQ constraint: no write-side change).

**Decisions:**
- **D-01** (DECIDE & STATE): a month folder is a direct child of `ai-reports/` whose name matches `^[0-9]{4}-(0[1-9]|1[0-2])(_1-report|_[0-9]+-reports)?$`: plain `YYYY-MM`, `YYYY-MM_1-report`, or `YYYY-MM_<N>-reports`. One exported-in-package helper in `report_index.go` (for example `reportMonthFolder(name string) bool`) is the single definition; second-opinion and REQ-707 reuse it. Reason: one rule, one place. A count that does not match the bundles inside is still a month folder here; recounting is REQ-707's job.
- **D-02** (DECIDE & STATE): only one level of month folders. A month folder inside a month folder is a bundle of that month (it is a directory), not a second container. Reason: the layout has one level; YAGNI.
- **D-03** (DECIDE & STATE): inside a month folder the rules are the top level's rules: non-dot directories are bundles, loose files and symbolic links are skipped. A symlinked month folder is skipped by the walk (as any symlink is) and refused by second-opinion's linked-ancestor check.
- **D-04** (DECIDE & STATE): second-opinion accepts a bundle whose parent is `ai-reports` or `ai-reports/<month folder>`; the refusal text becomes "bundle must be a physical report bundle directly under ai-reports or under a month folder".
- **D-05** (DECIDE & STATE): revise and second-opinion prose say `<dir>` is `ai-reports/<folder>`, `ai-reports/<month folder>/<folder>`, or a bare `<folder>` the catalog lists in either place; a new revision is still written directly under `ai-reports/`.
- **D-06** (DECIDE & STATE): tests are two new focused tests, named so the GREEN probe fails at base: `TestAIReportIndexReadsBundlesInsideMonthFolders` (counted and plain month folders, find, symlinked month folder skipped, supersedes across flat and nested) and `TestSecondOpinionBundleInsideMonthFolder` (nested bundle passes, `latest` resolves nested, symlinked month folder refused). The flat byte-identical promise is pinned by the existing catalog tests, which must pass unchanged.

*Generated by the orchestrator's exploration*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modify) — month-folder helper; the walk reads bundles inside month folders
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modify) — new month-folder catalog test
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify) — bundle check accepts a month-folder parent
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — new nested-bundle test
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — revise and second-opinion target rules, catalog sentence
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — catalog, revise and second-opinion sentences; the month-folder layout

**Files I will NOT touch:** the report writers (`skills/do-work-toolbox/actions/ai-report-reference.md`, `skills/do-work-toolbox/actions/architecture-report.md`, `skills/do-work-toolbox/actions/stakeholder-report.md`), `skills/do-work-toolbox/SKILL.md`, the board, any organize command (REQ-707), and the release files (finalization writes them).

**Acceptance criteria (restated from REQ):**
- [ ] A direct child of `ai-reports/` named `YYYY-MM`, `YYYY-MM_1-report` or `YYYY-MM_<N>-reports` is a month folder; its non-dot child directories are bundles with `Path` = `ai-reports/<month folder>/<bundle>`; the month folder is never a bundle.
- [ ] Flat bundles keep working; loose files and symbolic links are never bundles; a symlinked month folder stays refused.
- [ ] Second-opinion accepts a bundle whose parent is `ai-reports` or a month folder, with a regular `index.html`; `latest` resolves nested bundles.
- [ ] Revise accepts `ai-reports/<folder>`, `ai-reports/<month folder>/<folder>`, or a bare `<folder>` found in either place, and a supersedes chain resolves across a flat and a nested revision.
- [ ] `ai-report find` matches nested bundles by folder name.
- [ ] The action text and the guide no longer say a bundle must sit directly under `ai-reports/`.
- [ ] A repository with no month folders gets a byte-identical catalog.

## Pre-Flight

**Git:** ✓ clean tracked tree at d231e2fa (the REQ-706 claim commit, main after 0.309.0 and the UR-161 capture); only this run's `do-work/` artifacts (run directory, this REQ's orchestrator sections) were dirty. One untracked user file, `do-work/inbox/2026-10-28_do-work-upstream-suggestion-month-folder-report-layer.md`, is not this run's: it is hidden from discovery through a temporary `.git/info/exclude` line (removed at the end of the run) and never staged.
**Tests baseline:** ✓ `bash do-work/runs/work-2026-11-02-170509/REQ-706-preflight-probe.sh` exit 0, launched (`do-work/working/baseline.json`). Focused repository gate `bash do-work/runs/work-2026-11-02-170509/REQ-706-gate.sh` at d231e2fa exit 0 (2026-11-02T17:09:11Z to 17:09:53Z, 42 s); green-gate record satisfied.
**Dependencies:** ✓ Go toolchain and ShellCheck present (the gate ran both)

*Checked by work action*
