---
id: REQ-706
title: 'Addendum: ai-report tools read report bundles inside month folders'
status: pending
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
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md"]
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
