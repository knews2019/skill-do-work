## 0.305.20 — Copy Prepared Recovery States in Finalization Tests (2026-09-06)

Extended fixture reuse in `internal/finalization` tests to copy prepared semantic legacy and planned recovery baseline states, eliminating repetitive Git seed commits and repository construction.

- `internal/finalization/finalization_commands_test.go` defines thread-safe singletons and lazy initializers for semantic legacy and planned finalization template repositories, cleaning them up in `TestMain`.
- `internal/finalization/finalization_recovery_test.go` updates `seedSemanticLegacyTail` and `seedPlannedFinalization` to instantiate isolated copies via `os.CopyFS`, eliminating 69 Git subprocess executions and reducing cold wall time by 13.4% (-6.45s) and CPU by 14.0% (-5.38s).
- Adds `TestPreparedRecoveryTemplateIsolation` proving complete filesystem and Git history isolation across copies and underlying templates, with mutation tests confirming robust defect rejection.
