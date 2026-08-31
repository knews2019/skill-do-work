# Prime: do-work update shortcut

`tools/do-work-update.sh` is the terminal-facing implementation behind `just run-do-work-update`. It updates only a four-skill do-work suite installed inside the invoking project; the already-installed full-suite installer owns module/configuration reconciliation, and `actions/version.md` remains the canonical agent-driven update contract.

The five public shell entry points are compatibility launchers over the `do-work-cli` command, which owns the logic. **Updating requires Go 1.25.0 or newer**, which the `do-work-cli.sh` launcher enforces before it builds or runs anything.

## Read first

- `tools/do-work-update.sh` — the launcher: project-root argv, installed-skill-root derivation, delegation to `update-suite`.
- `tools/do-work-cli/internal/suiteinstall/update_transaction.go` — the update transaction: shared-install refusal, fetch, extract, manifest validation, version comparison, in-process install delegation, post-update verification.
- `tools/do-work-cli/internal/suiteinstall/install_transaction.go` — trusted installed module, managed-configuration (Just section, agent-instructions section, hook composition), verification, and exact recovery transaction. The install transaction is authoritative for the current surface set; any list here is illustrative and must not be read as closed.
- `tools/do-work-cli/internal/managedsection/managed_section.go` — byte-preserving managed-section replacement for whichever marker pair the caller names, plus the no-Just reserved-recipe collision scanner in `just_definitions.go`.
- `tools/do-work-cli/internal/settingshooks/settings_hooks.go` — order-preserving settings composition; a consumer's key order is part of the contract.
- `actions/version.md` — canonical safeguards the update path must keep aligned with.

## Do not edit

- `<project-root>/do-work/` — runtime queue, archives, and deliverables; the updater's explicit managed-path plan must never include it.
- Any consumer-owned bytes outside the installer's consent-gated managed marker spans — the condition is the rule, not the file. That covers a justfile outside the recipe markers and a `CLAUDE.md` outside the communication-style markers alike, and it covers whatever surface the installer manages next.

## Traps

- **The running binary sits inside a directory the install removes.** In the update flow the command lives at `.claude/skills/do-work/tools/do-work-cli/do-work-cli`, and the install replaces that whole module. Anything resolved relative to `os.Executable` must be read before the write phase, not after.
- **The launcher's staleness check is mtime-based.** It rebuilds whenever any `*.go`, `go.mod` or `go.sum` is newer than the binary, and `cp -R` copies in readdir order. A fixture that needs no rebuild must set mtimes explicitly rather than rely on copy order.
- **The built binary is not suite content.** Go embeds the build directory, so two projects that build it under different paths hold different bytes for identical sources. Never compare it as installed bytes; filter it out of a diff's output by path, never with `diff -x do-work-cli`, which would also match the module directory of the same name.

## Stakes

- `do-work-update.sh` → `update-suite` — project-local overwrite boundary
  Req: reject a skill root outside the invoking project, require `--project-root` to name the Git worktree root, refuse an upstream version that is not newer, review the diff and take the single confirmation through the install transaction (whose declined confirmation is a success with skipped work, not a failure), and verify the installed version afterward.
  Value: users can update without an agent turn while retaining the protection against clobbering a shared install or local customization.
  Risk: weakening any guard can overwrite user work or runtime queue data. The update transaction requires the project Git root and validates the suite archive before delegating to the install transaction in-process. That transaction reviews modules plus every owned configuration change, snapshots exact managed originals, and restores them plus the Git index on failure — a new managed surface earns a diff section, a backup, a recovery branch, and a post-write byte check in the same commit that adds it. Runtime, KB, application paths, unrelated settings, and any bytes outside a managed marker span must never enter that plan. Dirty module changes are named before the one confirmation; accepting discards them from both index and worktree before installation. `_dev/tests/update-script-behavior.sh` holds current-suite, hostile-manifest, dirty-consent, and forced-recovery behavior.

## Lessons

See [`lessons-do-work-update.md`](lessons-do-work-update.md) — read it before changing what **Read first** or **Traps** name above.
