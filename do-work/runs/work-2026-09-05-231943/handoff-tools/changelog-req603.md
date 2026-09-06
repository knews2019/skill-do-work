## 0.305.14 — Protected Inventory Launcher Passes Global Flags and Preserves Diagnostic Text (2026-09-06)

The `protected-inventory` compatibility launcher now forwards global flags (`--repo-root`, `--format`) ahead of the command token so it can run from any working directory, and the compatibility shim preserves diagnostic text and findings on failures instead of discarding them.

- `protected-inventory.sh` parses and forwards global options before the command verb, allowing `--repo-root` to be passed before or after the mode (`start` or `associate`).
- `internal/corehelpers/inventory.go` compatibility shim preserves prepared diagnostic text (`NO-DO-WORK-DIR`, `PARSE-FAILED`) and walk error findings on failure, ensuring exit 2 causes are visible.
- `commit.md` and `inspect.md` document that exit 2 with a finding or diagnostic text is an error to report rather than a skip condition, clarify exit 2 causes (such as non-git repository, status failure, or quarantine write failure), and note that `associate` exits 2 as not-started after `start --dry-run`.
- `commit.md` clarifies that re-running inventory uses `start` (which replaces quarantine rather than appending, whereas `associate` unions).
- `prescribed-shell-primitives.md` documents global flag forwarding and `--repo-root` support in the launcher.
