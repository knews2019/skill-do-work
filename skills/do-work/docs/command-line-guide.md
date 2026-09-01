# Deterministic Command-Line Interface

The installed core package owns one Go command platform for mechanical operations. Flat Just recipes expose its public interface without an LLM; natural-language actions call the same commands when they reach a deterministic phase.

## Discover and run recipes

Run this from the project root for the authoritative live inventory:

```bash
just --list
```

Recipes are flat. Invoke them directly and pass command arguments after the recipe name:

```bash
just do-work-doctor
just do-work-next REQ-042
just memory-recall "deployment decision"
just audit-metrics inventory --exclude-path do-work/
```

The managed section currently publishes these command families:

- Core lifecycle and maintenance: `do-work-cleanup`, `do-work-doctor`, `do-work-next`, `do-work-claim`, `do-work-complete`, `do-work-fail`, `do-work-cancel`, `do-work-unblock`, `do-work-answer`, `do-work-capture-files`, `do-work-release`, and `do-work-update`.
- Knowledge: `bkb-init`, `bkb-status`, `bkb-lint-structure`, `dream-scan`, the `interview-*` recipes, and the `memory-*` recipes.
- Toolbox: `do-work-note`, `architecture-report-preflight`, `generate-report-image`, `generate-report-image-batch`, `publish-portfolio-summary`, `install-last30days`, and `audit-metrics`.
- Board and compatibility: `run-kanban`, `run-kanban-cli`, `kanban-static`, `kanban-summary`, and the update compatibility alias `run-do-work-update`.

This grouped overview is orientation, not a second parser contract. `just --list` is the live inventory installed from the managed template.

## Call the CLI directly

Recipes resolve the project root and installed launcher for you. For automation that needs typed output, call the launcher directly:

```bash
.claude/skills/do-work/tools/do-work-cli.sh \
  --repo-root "$(git rev-parse --show-toplevel)" \
  --format json \
  doctor
```

Global options precede the command. `--format text` is the default for people; `--format json` returns the stable typed result used by actions and automation. Findings include exact next argv and verification argv where another step is available.

## Update compatibility

Use `just do-work-update` for the canonical no-agent update. `just run-do-work-update` remains available for compatibility and invokes the same `update-suite` transaction. Both retain the reviewed diff, confirmation, verification, and recovery boundary.

## Natural-language boundary

Natural-language routes remain useful where the workflow requires interpretation, drafting, or consent. They delegate mechanical phases to the canonical launcher. If that launcher is missing, fails, or returns malformed typed output, the action stops with the canonical finding and does not fall back to free-form mutation.
