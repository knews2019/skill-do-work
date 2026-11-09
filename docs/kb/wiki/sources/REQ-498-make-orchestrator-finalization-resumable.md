---
title: "Lessons from REQ-498: Make orchestrator finalization resumable"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-04/REQ-498-make-orchestrator-finalization-resumable.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-498: Make orchestrator finalization resumable

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Replace the crash-prone archive/release/commit tail with one CLI-owned, Git-private journaled finalization flow, and recover safe unfinished tails before selecting another REQ.

## Solution summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (created)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_journal.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (created)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

## What worked

- Reusing the journal phase engine, exact image identity, and typed result authority produced reliable phase replay, multi-group commit recognition, protected-state preservation, and actionable terminal evidence.

## What didn't work

- Inferring legacy ownership from only the dirty paths and current document metadata was too weak. A positive end-to-end fixture did not prove completeness: missing configured mirrors and pre-existing follow-up edits require negative preimage/required-set cases.

## Worth knowing

- Recovery association is safe only when it proves both sides of the boundary—every required member is present and every admitted byte belongs to the REQ. “All observed paths look coherent” is not equivalent to “the semantic set is complete.”

## Back-reference

See `do-work/archive/UR-096/REQ-498-make-orchestrator-finalization-resumable.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1249e856`.
