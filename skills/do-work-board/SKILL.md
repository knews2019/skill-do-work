---
name: do-work-board
description: Queue-kanban board, Testing workflow, queue activity calendar, summaries, and terminal digest for the modular do-work suite
argument-hint: "board [serve|static|summary|cli|verify] [--port N|--out DIR] | help"
---

# Do-Work Board Skill

This package owns the compiled queue-kanban application and its launcher. It reads the consuming project's `do-work/` records produced by sibling `do-work`, but it does not own queue orchestration.

## Routing

| Trigger | Route |
|---|---|
| empty, `help` | `./actions/help.md` |
| `board`, `kanban`, `queue board`, `show the board`, or a board mode | `./actions/board.md` |
| `verify`, `check invariants`, `probes` | `./actions/board.md` |

Pass `serve`, `static`, `summary`, `cli`, `verify`, `--port N`, and `--out DIR` through to the board action. An unknown command prints board help and stops.

## Ownership boundary

- Core queue actions and the schema contract live under sibling `../do-work/`.
- Reviews, inspection, and other repository utilities live under sibling `../do-work-toolbox/`.
- The full-suite installer owns reconciliation of `justfile.template`; this package owns the template bytes.
- `tools/queue-kanban/` is the only board implementation in this package and must build from its own Go module.

Read the selected action completely before running it. Serve mode is a long-running process and may be backgrounded when the harness supports it; all other modes run in the foreground.
