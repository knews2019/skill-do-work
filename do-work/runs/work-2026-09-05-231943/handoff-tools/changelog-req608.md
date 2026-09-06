## 0.305.19 — Run Inventory Data Matrix In-Process Without Subprocesses (2026-09-06)

Decoupled porcelain byte parsing from Git acquisition in `internal/corehelpers/inventory.go`, allowing synthetic inventory test matrices to run completely in-process without spawning Git and CLI subprocesses.

- `internal/corehelpers/inventory.go` extracts `parseInventoryBytes` to cleanly separate byte stream parsing from `gitOutput`.
- `internal/corehelpers/inventory_test.go` executes the 45-case porcelain status matrix and 10-case secret origin/ambiguity matrix in-process, eliminating 56 `do-work-cli` subprocess executions and reducing test wall time by 82.5% (-7.31s) and child CPU by 61.2% (-1.01s).
- Adds explicit contract tests for malformed short porcelain records, missing rename origins, record ordering, metadata exclusions, and cross-row secret promotion while retaining end-to-end launcher and real-Git coverage.
