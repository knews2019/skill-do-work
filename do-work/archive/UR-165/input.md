---
id: UR-165
title: 'folder-structure-check sees justfile recipes with default parameter values'
created_at: 2026-11-09T21:47:36Z
requests: [REQ-713]
word_count: 101
---
## Full Verbatim Input
> ```
> do-work capture-request: Fix folder-structure-check so it sees justfile recipes whose parameters have default values. Today justRecipeHeader in skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go stops at the "=" inside a parameter like `run-kanban $port="8090":`, so that recipe's 3 body lines are never counted (layout rule: a justfile recipe holds no logic and calls cmd/<name>.sh in one line). Keep the existing exclusions for `name := value` variables and set/alias/export/import/mod directives. Add a case to TestFolderStructureCheckReportsEachGapClass with a multi-line recipe that has a quoted default parameter, and show the test fails without the fix. Verify: go test ./internal/toolboxcommands/ passes, and running the check on this repo reports `justfile: run-kanban`.
> ```
