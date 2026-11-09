---
id: UR-169
title: 'Suite-managed justfile recipes call a script in one line'
created_at: 2026-11-09T21:47:36Z
requests: [REQ-717]
word_count: 151
---
## Full Verbatim Input
> ```
> do-work capture-request: Move the logic out of the suite-managed justfile recipes so each managed recipe is one line (layout rule 5: a justfile recipe holds no logic and calls a script in one line). In skills/do-work-board/justfile.template, run-kanban has 3 body lines (port validation, stop a stale queue-kanban listener, build and serve) and kanban-static has 2 (build and generate, then the .git/info/exclude write). Move that logic into scripts shipped with do-work-board (for example under skills/do-work-board/tools/), keep the exact behavior and messages, and make each recipe a single call. The installer rewrites the managed block in consumer justfiles, so also update this repo's root justfile. Depends on the REQ that fixes folder-structure-check for recipes with default parameter values. Verify: folder-structure-check reports no FOLDER-STRUCTURE-JUST-LOGIC on this repo, just run-kanban still replaces a stale instance on its port, just kanban-static still writes the build/ snapshot and the exclude line, and the board contract tests pass.
> ```
