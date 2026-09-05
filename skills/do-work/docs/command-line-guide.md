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

- Core lifecycle and maintenance: `do-work-cleanup`, `do-work-doctor`, `do-work-next`, `do-work-claim`, `do-work-complete`, `do-work-fail`, `do-work-cancel`, `do-work-unblock`, `do-work-answer`, `do-work-capture-files`, `do-work-defer-gate`, `do-work-release`, and `do-work-update`.
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

## Defer an unrelated repository-gate failure

`do-work-cli defer-gate --manifest <path>` owns the mechanical deferral boundary. Its strict manifest binds the exact claimed parent and checkpoint preimages, writer label, direct non-zero gate result, structured command argv, diagnostic fingerprint and evidence, repair identity, reservation path, and optional paired implementation base/merge commits.

The transaction either creates one fingerprint-keyed `sweep: true` repair REQ with a canonical `## Instances` checklist or folds the parent into the unique exact-field match. In the same atomic apply it returns the parent to `pending`, appends the repair to `depends_on`, removes only the exact writer claim, records `## Repository Gate Deferral` evidence, and moves the parent back to `do-work/queue/`. Refused, stale, staged, colliding, ambiguous, unsafe, prefix-only, side-branch merge, or partially published inputs do not leave a partial lifecycle. Optional merge evidence must be a non-empty base-to-merge range already contained in current `HEAD`. JSON and text results carry the same semantic `gate_deferral` evidence, including paths, dependency, sweep key, create/fold outcome, and any saved merge range.

The action runs the canonical gate before source work, chooses create versus exact-preimage fold, and consumes the typed result to schedule the repair and suppress the parent. A committed repair remains foldable after SessionStart removes its reservation; an existing marker must still match exactly. On late worktree failures the action compares the same fingerprint at the saved pre-merge revision and supplies the verified base/merge pair only for a matching unrelated failure. Resumption proves merge ancestry and rename-aware path history before reusing implementation, then reruns qualification, tests, the gate, and review. `defer-gate` itself never executes a gate, chooses a fingerprint, or inspects path drift, and it schedules nothing: queue selection and the mechanical evidence gates belong to `advance`.

Collision retry is caller-owned and typed: retry only a collision refused before mutation with zero changes, or a post-mutation collision whose rollback succeeded. Repair REQs never defer themselves; a red or unverifiable final gate fails the repair canonically and leaves parents dependency-gated while unrelated selection continues. If the repair gate is already green before implementation, the action records the exact reviewed no-op evidence, completes/archives canonically, and skips changelog/version release mutation.

## Validate an already-green repair

`validate-already-green-repair --request-path <working REQ> --writer <completion writer> [--at <RFC3339>]` is the read-only authority shared by the work TDD exception and review's no-diff exception. It extracts the repair intake, no-op fingerprint, gate argv, and recorded revision from the request; callers cannot supply those values. It verifies the persisted green record at the recorded past revision, checks project/release state, performs canonical completion as a dry run, and compares staged paths to that exact result rather than authorizing an archive directory prefix.

The typed `already_green_repair` object carries separate `tdd_allowed` and `review_allowed` decisions from one observation set. Review additionally requires a successful completion projection and an exact staged subset. It also carries canonical completion paths, staged/project paths, stable reason codes, offending paths, and the gate evidence. Actions stop on a missing, malformed, failed, or false decision and never reconstruct the predicate from prose. Reuse the same `--writer` and `--at` values in finalization because they can affect the exact completion paths.

## Update compatibility

Use `just do-work-update` for the canonical no-agent update. `just run-do-work-update` remains available for compatibility and invokes the same `update-suite` transaction. Both retain the reviewed diff, confirmation, verification, and recovery boundary.

## Natural-language boundary

Natural-language routes remain useful where the workflow requires interpretation, drafting, or consent. They delegate mechanical phases to the canonical launcher. If that launcher is missing, fails, or returns malformed typed output, the action stops with the canonical finding and does not fall back to free-form mutation.
