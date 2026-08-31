# REQ-411 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-411-implement-queue-selection`
Branch: `worktree-agent-REQ-411-implement-queue-selection`
REQ: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-411-implement-queue-selection.md`
Plan: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-411-plan.md`
Exploration: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-411-exploration.md`
Integration owner: main checkout; do not edit `do-work/`, release metadata, changelogs, or version files.

## Frozen write set

- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go`
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/select-simple-reqs.sh`
- `_dev/tests/select-simple-reqs-behavior.sh`
- `skills/do-work/actions/run-simple-reqs.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `_dev/tests/contract-regressions.sh`

## Required method

- Read the REQ, accepted plan, exploration, UR-081, action and CLI primes/lessons, root/project instructions, and applicable general/testing/coding guardrails.
- Start with the mixed command fixture and observe the captured `UNKNOWN-COMMAND`/missing-contract RED before production changes. Keep mutation-sensitive assertions for each major selector branch.
- Consume one typed snapshot and reuse `dependencygraph`; do not add repository/request/dependencygraph production schema or rescan the filesystem.
- Keep selection read-only. Do not implement claim, state transition, archive, release, or queue-kanban behavior.
- Preserve explicit-REQ versus UR/default token provenance through dependency/assignment/negligible overrides. Keep `--wave` depth separate from `--fan-out` bounding.
- Execute blocked probes only through the shipped process-group-safe runner with exact contained bytes; inject the runner in pure command tests.
- Extend the result envelope with typed selected/excluded records and prove text/JSON parity, stable ordering, actionable reasons, exact next argv/Just/verification commands, and non-null JSON arrays.
- Route the simple selector through canonical readiness while retaining its extra cheap-work vetoes and every existing diagnostic.
- Stay exactly inside the 18-file scope. Run focused nextselection/result tests, full CLI tests/vet, selector behavior, blocked-check prescribed cases, contract regressions, exact Go 1.25, qualification/scope checks, and diff checks.
- Commit on the branch, keep the worktree clean, and write the handback to `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-411-handback.md`.

## Handback

Record the full commit hash, RED/GREEN evidence, every changed file, commands and exit codes, selection schema/ordering, token-provenance table, blocked-probe invocation, integration seams, decisions, and discoveries. Do not create or edit queue files.
