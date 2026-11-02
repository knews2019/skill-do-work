---
id: REQ-706
title: 'Addendum: ai-report tools read report bundles inside month folders'
status: completed
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
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go"]
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
integration_at: 2026-11-02T17:18:57Z
review_at: 2026-11-02T17:28:23Z
kb_status: pending
commit: ba902c6370168a0fc64a6c45c3e9328c61595c5e
completed_at: 2026-11-02T17:54:50Z
release_at: 2026-11-02T17:54:50Z
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
- [x] **[PLAN]:** (from the builder hand-back) Read the brief, REQ, crew rules, primes and the `catalog-tie-order` bullet. One helper with one regexp; the walk descends one level into month folders; `readReportBundle` takes a relative path and uses `path.Base` for every name-derived field; one predicate change in `secondOpinionBundle` covers every entry route (`guard-on-one-entry-path`). RED tests written before any production edit.
- [x] **[APPLY]:** (from the builder hand-back) Code matches the plan; exactly the six Scope files changed, one commit `c004cba0`.
- [x] **[UNIFY]:** (from the builder hand-back) `git diff 782f083a --stat`: 6 files, 147 insertions, 14 deletions. gofmt clean, `go vet ./internal/toolboxcommands` 0, GREEN tests 0, selection `SecondOpinion|AIReportIndex` 36 PASS / 0 FAIL / 0 SKIP, `REQ-706-probe.sh` 0, `contract-regressions.sh` 0, `shipped-package-reference-contract.sh` 0, `git diff --check` 0. Each file reviewed: flat bundles byte-identical, Rel error still refused, both new tests carry `// Pins:` comments, prose sweep of the direct-child rule clean. Mutation check: dropping `child.IsDir()` makes the catalog test fail on the loose file.
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
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` (modify, added after review F1, D-12) — the architecture-report pre-flight scan finds a prior report inside a month folder
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go` (modify, added after review F1, D-12) — one focused test for the nested prior report

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

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` (modified, review fix F1)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go` (modified, review fix F1)

**What was done:** A new `reportMonthFolder` helper (one regexp: `YYYY-MM`, `YYYY-MM_1-report`, `YYYY-MM_<N>-reports`) lets `reportIndexWalk` read the non-dot child directories of a month folder as bundles with path `ai-reports/<month folder>/<bundle>`, one level deep; `readReportBundle` takes the bundle's relative path and reads date, kind, linked ids and default title from the bundle's own folder name, so flat output is unchanged. `secondOpinionBundle` accepts a parent of `ai-reports` or a month folder, which covers plain, `latest`, page and `--finalize` targets. Two new tests pin the catalog, find, supersedes and second-opinion cases; the action's catalog, revise and second-opinion rules and the guide describe the month-folder layout and say new reports are still written directly under `ai-reports/`.

**Review fixes (builder commit `a7a2fe6b`, re-merged as `ba902c63` with the same `<pre>`):** the architecture-report pre-flight scan finds a prior report inside a month folder (F1, new `TestArchitectureScanFindsPriorReportInMonthFolder`); second-opinion refuses a month folder itself as a bundle (F2) with an explicitly grouped condition (F4); the action no longer calls a revision a "sibling" (F3); the month-folder rule reads `_<N>-report(s)` for any N (F5, D-13).

## Decisions

(from the builder hand-back)

- **D-07** (DECIDE & STATE): if a month folder cannot be read, the whole walk fails with `AI-REPORT-INDEX-FAILED`, the same as an unreadable `ai-reports/`. Skipping would silently drop a month of reports and could make `latest` pick the wrong bundle.
- **D-08** (DECIDE & STATE): the helper is `reportMonthFolder(name string) bool` so REQ-707 can reuse it; the regexp is `reportMonthPattern` so gofmt does not re-align the existing `var` block.
- **D-09** (DECIDE & STATE): the second-opinion test adds one older flat bundle, so the `latest` case proves "newest across flat and nested", not "the only bundle".
- **D-10** (DECIDE & STATE): the guide's layout paragraph sits at the top of `## Report Catalog`; in `ai-report.md` "month folder" is defined at its first use (the catalog sentence).
- **D-11** (DECIDE & STATE, orchestrator): heavy lanes are not run and no `## Heavy Verification Plan` hold is written: the user chose focused tests only (2026-11-02). `## Testing` records the lanes the planner would have selected, marked "heavy lanes not run: user chose focused tests only (2026-11-02)". Writing `heavy_verified_*` stamps would claim a verification that did not happen.
- **D-12** (DECIDE & STATE, orchestrator scope judgment after review F1): `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` and `architecture_test.go` join the Scope and `write_set`. The architecture-report pre-flight scan is a reader of the bundle layout the REQ redefines ("every reader of the bundle layout", alternate-writer-contract-drift); left alone, REQ-707's first organize run would make every later architecture report claim no prior baseline. Fixed on the builder branch with one focused test.
- **D-13** (DECIDE & STATE, orchestrator, review F5): the month-folder name rule accepts `_<N>-report` and `_<N>-reports` for any N (`(_[0-9]+-reports?)?`), not only the canonical `_1-report`. Readers stay tolerant of a hand-named or stale-count folder so its reports are never hidden, and REQ-707's organize renames it to the canonical counted name. Value: no report disappears from the catalog because a folder name is slightly off. Risk: a dated report bundle could never be named like `2026-07_2-report`, because bundle names start with a full `yyyy-mm-dd` date, so nothing is misread. Reversible.

## Discovered Tasks

(from the builder hand-back)

- **impact-negligible** Supersedes resolution in `reportIndexWalk` keys on the bare folder name. With month folders the same bundle name can exist twice (in `2026-09/` and `2026-09_3-reports/`, or flat and nested after a copy), and an `ai-report-supersedes` meta then resolves to whichever duplicate sorts last. REQ-707 (organize) should refuse to create a duplicate bundle name; its builder brief carries this. → report only

## Qualification

**Gate records:** `advance --diff-range 6e98d229..3edfd2c9` → `qualify` satisfied (no debug artifact, no output-primitive finding, P-A-U boxes ticked) and `scope-drift` satisfied with no finding: the six touched project paths are exactly the six declared in `## Scope`.

**Requirement trace (read against `git diff 6e98d229..3edfd2c9` and the merged files):**
- Month folder definition: `reportMonthPattern` = `^[0-9]{4}-(0[1-9]|1[0-2])(_1-report|_[0-9]+-reports)?$`, used only through `reportMonthFolder`. `reportIndexWalk` reads a matching direct child's non-dot child directories as bundles with `Path` `ai-reports/<month folder>/<bundle>` and never appends the month folder itself. Pinned by `TestAIReportIndexReadsBundlesInsideMonthFolders` for `2026-09_1-report` and plain `2026-09`.
- Flat bundles, loose files, symlinks: the flat branch calls `readReportBundle(root, name)` with `relative == folder`, so every field is byte-identical; `DirEntry.IsDir()` is false for symlinks at both levels. The test's loose file and symlinked `2026-08_1-report` add no bundle; the four existing catalog tests pass unchanged.
- Second-opinion: `secondOpinionBundle` accepts `parent == ai-reports` or a month-folder parent whose own parent is `ai-reports`; both `validateNoLinkedAncestors` calls stay, so a symlinked month folder is refused. Every entry route (bundle, `latest`, page target, `--finalize`) calls this one function. Pinned by `TestSecondOpinionBundleInsideMonthFolder` (prepared, `latest` picks the nested newest over an older flat one, symlinked month refused).
- Revise: the rule is prose (`ai-report.md` revise step 1) and now admits `ai-reports/<month folder>/<folder>` and a bare folder the catalog lists in either place, and states the new revision is written directly under `ai-reports/`. The supersedes chain across flat and nested is pinned in the catalog test (`2026-11-02_1000_a-rev1` supersedes the nested `a`).
- `find`: unchanged code matches `path.Base(bundle.Path)`; the test finds the nested bundle.
- Prose: `ai-report.md` :42, :54, :72 and the guide's catalog, revise and second-opinion paragraphs no longer say direct child; the builder's sweep and `REQ-706-probe.sh` confirm. No organize command is described (REQ-707).
- Byte-identical catalog without month folders: the flat code path is unchanged in output; existing catalog tests pass.

**Scope comparison:** declared 6, touched 6, no undeclared touch. The files the Scope said it would not touch (report writers, `SKILL.md`, release files) are untouched.

**Judgment:** substantive, small (147 insertions, 14 deletions, 45 production lines), matches D-01 to D-10. Qualified.

## Testing

**Tests run:** focused repository gate `bash do-work/runs/work-2026-11-02-170509/REQ-706-gate.sh` at merge `3edfd2c9` through `run-timed-command` (2026-11-02T17:19:42Z to 17:20:25Z, 43 s, 1-minute load 3.3): `go vet ./internal/toolboxcommands/`, `go test ./internal/toolboxcommands/ -count=1` (18.1 s), then `_dev/tests/shipped-package-reference-contract.sh`, `contract-regressions.sh`, `action-shell-blocks.sh`, `prescribed-shell-canonicalization.sh`, `suite-manifest-contract.sh`. Then `advance --gate-arg bash --gate-arg do-work/runs/work-2026-11-02-170509/REQ-706-gate.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-11-02-170509/REQ-706-probe.sh`: green-gate, scope-drift and `run-blocked-check` satisfied (probe status 0: both new tests PASS, the direct-child sentences are gone, `month folder` is in the action and the guide).
**Result:** ✓ All passing (gate exit 0 on the first run; the `toolboxcommands` package passes in full; builder's selection `SecondOpinion|AIReportIndex` 36 PASS, 0 FAIL, 0 SKIP)

**Red-green validation:**
- `TestAIReportIndexReadsBundlesInsideMonthFolders` (`report_index_test.go`): ✗ before implementation (catalog paths `[ai-reports/2026-11-01_0900_b ai-reports/2026-09_1-report]`, the month folder listed as a bundle and `a` missing; same for plain `2026-09`) → ✓ after (PASS for both subtests at `c004cba0` and in the gate at `3edfd2c9`). This is the `## Red-Green Proof` RED case and GREEN condition (two bundles, `b` first, nested path, no month-folder bundle, plain `YYYY-MM` too).
- `TestSecondOpinionBundleInsideMonthFolder` (`report_second_opinion_test.go`): ✗ before (refused `bundle must be a physical direct child of ai-reports`) → ✓ after (prepared; `latest` resolves nested; symlinked month folder refused). This is the Red-Green Proof's second-opinion case and its symlink GREEN clause.
- The byte-identical GREEN clause is covered by the four existing `TestAIReportIndex*` catalog tests, unchanged and passing.

**Repository gate after the review fixes:** the same focused gate argv at `ba902c63` through `run-timed-command` (2026-11-02T17:31:46Z to 17:32:34Z, 48 s, 1-minute load 3.2) exit 0, recorded with `record-green-gate` (advance refuses gate input once the REQ is past review). `REQ-706-probe.sh`, extended with `TestArchitectureScanFindsPriorReportInMonthFolder`, exits 0 there. Builder RED for the fixes: the architecture scan picked the older flat report; the `2026-09_2-report` subtest listed the month folder as a bundle; second-opinion accepted the month folder itself.

**New tests added:**
- `TestAIReportIndexReadsBundlesInsideMonthFolders` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go`
- `TestSecondOpinionBundleInsideMonthFolder` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`
- `TestArchitectureScanFindsPriorReportInMonthFolder` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go` (review fix F1)

**Heavy verification plan:** heavy lanes not run: user chose focused tests only (2026-11-02). Re-planned on the cumulative range after the review fixes; the planner would have selected (D-11):
- Range: 6e98d229c67351a79d1c67daf22fa38439b1441b..ba902c6370168a0fc64a6c45c3e9328c61595c5e
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go matched subtree skills/do-work/tools/do-work-cli (and 5 more paths)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work-toolbox/actions/ai-report.md matched subtree skills (and 7 more paths)
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go matched subtree skills/do-work/tools/do-work-cli (and 5 more paths)
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go matched subtree skills/do-work/tools/do-work-cli (and 5 more paths)

*Verified by work action*

## Review

**Overall: 95%** | 2026-11-02T17:28:23Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1: `architecture.go:62-85` `architectureScan` reads only the direct children of `ai-reports/`, so a prior architecture report in a month folder is missed, `prior_report` is empty, and the next report claims no prior baseline (reproduced). Suggested fix: reuse `reportMonthFolder` there with a focused test, folded into REQ-707. — impact-user-visible → report only

**Minor findings:** F2: `report_second_opinion.go:387-389` accepts the month folder itself as a bundle when it holds a regular `index.html` (reproduced `prepared`), although the catalog never lists it; refuse a month-named direct child. — impact-negligible → report only. F3: `ai-report.md:52` and `:203` still call the revision a "sibling" bundle, which contradicts step 1's directly-under-`ai-reports/` rule; reword. — impact-negligible → report only. Nit F4: `report_second_opinion.go:389` mixes `||` and `&&` without parentheses. — impact-negligible → report only. Nit F5: `report_index.go:47` accepts `_N-reports` for any N but `_N-report` only for N=1, so a hand-named `_2-report` folder hides its reports. — impact-negligible → report only
**Acceptance:** Pass — Implementation and Integration stages: focused `AIReportIndex|SecondOpinion` tests green at `3edfd2c9`; a scratch CLI build confirmed the nested catalog, find, page links, second-opinion prepare on a nested bundle, a page target and `latest`, and a byte-identical flat catalog against a `6e98d229` build. Deployment and Live acceptance unassessed.
**Restatement sweep:** redefined where a report bundle may sit under `ai-reports/` (directly, or inside a month folder that is never itself a bundle). Swept `skills/`, `suite/` and `README.md` (no root `docs/` exists) for direct-child, every-directory-is-a-bundle, sibling and one-level readers: `ai-report.md:42/54/72` and `ai-report-guide.md:74/95/97/105/117` agree. Stale: `architecture.go:62-85` (F1) and `ai-report.md:52/203` "sibling" (F3).
**Suggested testing:** 4 items
**Follow-ups created:** None (5 findings report only)

**Re-check (`3edfd2c9..ba902c63`, 2026-11-02T17:34Z, recorded 2026-11-02T17:55:00Z):** F1 to F5 fixed; overall 97%, acceptance Pass. Second sweep of every `ReadDir`/`Glob`/`WalkDir` site in shipped Go found no other one-level reader of `ai-reports/`. One new Nit: N1 `report_second_opinion.go:392` refuses a month folder named as a bundle with the generic "directly under ai-reports or under a month folder" message, which does not say why. — impact-negligible → report only. (The reviewer's own append stalled; the orchestrator copied its section into `do-work/runs/work-2026-11-02-170509/REQ-706-review.md` from its transcript.)

*Reviewed by review-work action*

**Orchestrator disposition:** F1 to F5 are fixed now on the builder branch rather than left as report-only items, because REQ-707 (organize) is about to move finished months into month folders and F1 would then hit every architecture report. Scope extended (D-12), suffix rule widened (D-13); the fix is re-merged with the same `<pre>`, re-gated and re-checked by the same reviewer (see the re-check line below once recorded).

## Lessons Learned

**What worked:** One predicate in `secondOpinionBundle` covered every second-opinion entry route (plain, `latest`, page target, `--finalize`), because they all resolve through that one function. Writing both tests before the code made the RED output name the exact old behavior (month folder listed as a bundle, nested bundle refused).
**What didn't:** The capture and the exploration listed the readers of the bundle layout from the inbox note (catalog walk, second-opinion check, revise prose) and the sweep searched prose for the direct-child sentence. A third code reader, the architecture-report pre-flight scan (`architecture.go`), does its own `os.ReadDir` of the reports folder and was only found by the reviewer's reproduction. A layout change needs a grep for every code path that lists the folder (`ReadDir` of `ai-reports`), not only for the sentences that state the rule.
**Worth knowing:** Readers accept `_<N>-report(s)` for any N (D-13) so a stale or hand-made month folder never hides reports; only REQ-707's organize writes the canonical counted names. Supersedes resolution still keys on the bare folder name, so duplicate bundle names across month folders resolve to whichever sorts last (Discovered Tasks).

## Orientation

[MAP CHANGED] Now the report tools read reports filed in month folders: `ai-report index` and `find`, the revise and second-opinion targets, and the architecture-report pre-flight treat `ai-reports/<YYYY-MM...>/` as a container of report bundles, never as a bundle. Lives in the toolbox report tooling (`do-work-cli` `toolboxcommands`, the `ai-report` action and guide). Why it matters: the bundle layout contract now has two levels, and REQ-707's organize command writes that layout. Prime spot-check: `_dev/primes/prime-action-files.md` and `_dev/primes/prime-releases.md` do not describe the report layout and their paths still exist; not stale.

## Timing

Observed 2026-11-02T17:05:09Z to 2026-11-02T17:32:34Z: 27m 25s total, 22m 41s attributed across 6 events, 4m 44s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 7m 33s | 1 |
| review | 7m 17s | 1 |
| exploration-preflight | 6m 12s | 1 |
| verification-gate | 1m 31s | 2 |
| handback-merge | 8s | 1 |

Slowest stage: builder-work / builder worktree build, 7m 33s, outcome success.
Slowest command: verification-gate / focused repository gate after review fixes, 48s, exit 0, bash (2 argv tokens).
