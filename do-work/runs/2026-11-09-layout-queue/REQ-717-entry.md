## 0.315.1 — Board Recipes Call Shipped Scripts in One Line (2026-11-10)

The `run-kanban` and `kanban-static` recipes are now one-line calls to scripts shipped with the board, so the suite's own recipes no longer show up as `FOLDER-STRUCTURE-JUST-LOGIC` gaps in `folder-structure-check`.

- The logic moved unchanged into `.claude/skills/do-work-board/tools/run-kanban.sh` and `kanban-static.sh`: the same port check, the same stop of a stale board on the port, the same messages. `do-work update` rewrites the managed recipe block and installs the scripts in the same step.
- `just` now prints the one-line script call when a recipe starts, in place of the long inline shell.
- This repository's dated changelog archives moved to `docs/changelog/`; the archive links at the top of this file point there. Installed copies never held the archives.
