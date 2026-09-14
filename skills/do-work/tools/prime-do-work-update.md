# Prime: do-work update shortcut

`tools/do-work-cli.sh … update-suite` is the singular update implementation behind canonical `just do-work-update`, compatibility `just run-do-work-update`, and the natural-language version action. It updates only a four-skill do-work suite installed inside the invoking project; the install transaction owns module/configuration reconciliation, and `actions/version.md` owns the agent's judgment and rendering contract.

The five public shell entry points are compatibility launchers over the `do-work-cli` command, which owns the logic. **Updating requires Go 1.24.0 or newer, or a host the prebuilt binary covers**; the `do-work-cli.sh` launcher builds with the toolchain when one is usable and otherwise fetches the checksum-verified prebuilt binary for the installed version before it runs anything.

## Read first

- `tools/do-work-cli.sh` — the canonical launcher: repository-root argv, toolchain validation, and dispatch to `update-suite`; `tools/do-work-update.sh` remains a compatibility launcher only.
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

- `do-work-cli/internal/suiteinstall/update_transaction.go` + `do-work-cli/internal/suiteinstall/install_transaction.go` — project-local update boundary
  Req: require the project Git root and an installed skill inside it; validate the archive, skip equal versions successfully, reject older upstream versions, and obtain the install transaction's single confirmation before writes.
  Value: callers can update the suite while retaining consumer configuration outside owned spans; declining confirmation is a successful no-change result.
  Risk: accepting the diff discards dirty managed-module content. The install transaction must snapshot and restore managed originals and the Git index on failure, then verify the installed version; queue, KB and application paths must stay outside its plan.

## Lessons

See [`lessons-do-work-update.md`](lessons-do-work-update.md) — read it before changing what **Read first** or **Traps** name above.
