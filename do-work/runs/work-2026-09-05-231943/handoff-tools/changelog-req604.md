## 0.305.15 — Unify Atomic-Download Occupancy and Handle Post-Publish Stat Error (2026-09-06)

The `atomic-download` core helper now enforces the same occupancy rule in live execution as it does during dry-run, refusing an existing regular file with exit 2 (`DOWNLOAD-TARGET-OCCUPIED`) rather than silently replacing it, and checks the error on post-rename `os.Stat` to prevent nil pointer panics under concurrent races.

- `internal/corehelpers/commands.go` checks destination occupancy before branching on dry-run, ensuring both dry-run and live runs refuse directories with exit 1 (`target is a directory`) and existing regular files with exit 2 (`target already exists`).
- `internal/corehelpers/commands.go` handles errors on post-rename stat, returning typed failure finding `DOWNLOAD-FAILED` instead of panicking on nil dereference.
- `prescribed-shell-primitives.md` documents the unified occupancy rule refusing existing destinations before fetching and the error-checked stat.
