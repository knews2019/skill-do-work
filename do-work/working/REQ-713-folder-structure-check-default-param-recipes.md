---
id: REQ-713
title: 'folder-structure-check sees justfile recipes with default parameter values'
status: claimed
created_at: 2026-11-09T21:47:36Z
user_request: UR-165
domain: backend
prime_files: []
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
claimed_at: 2026-11-09T21:51:56Z
---
# folder-structure-check sees justfile recipes with default parameter values
## What
Fix `folder-structure-check` so a justfile recipe whose parameters have default values (for example `run-kanban $port="8090":`) is recognized as a recipe and its body lines are counted. Today the header regex `justRecipeHeader` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go` stops at the `=`, so the recipe is skipped.
## Detailed Requirements
- Recognize recipe headers with quoted or unquoted default parameter values, `*args`/`+args` variadics and `$`-exported parameters.
- Keep the existing exclusions: `name := value` variables and the set/alias/export/import/mod directives are never recipes.
- Add a case to `TestFolderStructureCheckReportsEachGapClass` with a multi-line recipe that has a quoted default parameter, and show it fails without the fix.
- Layout rule this enforces: a justfile recipe holds no logic and calls `cmd/<name>.sh` in one line.
## Context
Source: do-work audit-current-folder-structure on this repository (2026-11-09); the maintainer settled every gap through the ask tool and chose to capture the prompts now.
Recorded exceptions from the same audit (do not move them): VERSION, skills/ and suite/ stay at the root because the release tool and the updater archive layout need them there; _dev/ stays as this repository's maintainer-only folder.
## Red-Green Proof
**RED prompt/case:** Run `folder-structure-check` on this repository; the root justfile's `run-kanban $port="8090":` recipe has 3 body lines.
**Why RED now:** The output reports only `justfile: kanban-static`; `run-kanban` is missing because the header regex fails on the `=`.
**GREEN when:** The output reports `FOLDER-STRUCTURE-JUST-LOGIC justfile: run-kanban`, the new test passes, and `go test ./internal/toolboxcommands/` passes.
**Validation:** User confirmed
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24067 tokens, `slugged: partial`, so it cannot be narrowed; matched because the fix is in do-work-cli internals.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: Fix folder-structure-check so it sees justfile recipes whose parameters have default values. Today justRecipeHeader in skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go stops at the "=" inside a parameter like `run-kanban $port="8090":`, so that recipe's 3 body lines are never counted (layout rule: a justfile recipe holds no logic and calls cmd/<name>.sh in one line). Keep the existing exclusions for `name := value` variables and set/alias/export/import/mod directives. Add a case to TestFolderStructureCheckReportsEachGapClass with a multi-line recipe that has a quoted default parameter, and show the test fails without the fix. Verify: go test ./internal/toolboxcommands/ passes, and running the check on this repo reports `justfile: run-kanban`.*
