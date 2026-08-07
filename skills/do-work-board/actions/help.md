# Board Help Action

> Prints the board package menu and stops without building, serving, generating, or editing anything.

```text
do-work-board — queue visualization and testing workflow

  do-work-board board                 Live board at http://localhost:8090
  do-work-board board --port 8091     Live board on another local port
  do-work-board static                Self-contained HTML snapshot
  do-work-board summary               Queue column and warning counts
  do-work-board cli                   In-flight terminal digest
  do-work-board help                  Show this menu

The full-suite installer adds equivalent Just recipes:
  just run-kanban [port]
  just run-kanban-cli
  just kanban-static
  just kanban-summary
```

The board needs the Go toolchain. Missing Go is reported as a board-only limitation and never blocks core queue work.
