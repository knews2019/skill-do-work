---
title: "Lessons from REQ-491: Add canonical repository-gate deferral lifecycle"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-04/REQ-491-add-canonical-repository-gate-deferral-l.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-491: Add canonical repository-gate deferral lifecycle

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Add one transactional `defer-gate --manifest` operation that converts an unrelated repository-gate failure into explicit repair work and safely returns the parent REQ to the dependency-gated queue. Extend the request model and selector so repair work and resumed parents run in the intended order without stopping unrelated work.

## Solution summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/flat-just-recipes-behavior.sh` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)
- `justfile` (modified)
- `skills/do-work-board/justfile.template` (modified)
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/command-line-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` (created)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` (created)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_manifest_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_types.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

## What worked

- Exact preimage-bound transactions plus race, alternate-consumer, and repository-topology fixtures exposed failures that happy-path publication tests could not. A fresh re-review caught two missing topology states after the initial remediation suite was green.

## What didn't work

- Treating “dirty target support” as a parent/checkpoint concern was too narrow; every existing mutation target, including a folded repair and a move destination, needs classification before planning is considered complete.

## Worth knowing

- `defer-gate` is a cross-reader contract: publication, Git rollback, request/schema projections, selector ordering, recipes, and board sweep parsing must move together. REQ-493 carries the remaining topology-class closure.

## Back-reference

See `do-work/archive/REQ-491-add-canonical-repository-gate-deferral-lifecycle.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0a5d4e44`.
