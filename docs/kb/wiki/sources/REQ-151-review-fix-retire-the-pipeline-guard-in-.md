---
title: "Lessons from REQ-151: Review fix: Retire the pipeline guard in manual settings reconciliation"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-01/REQ-151-review-fix-retire-the-pipeline-guard-in-.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-151: Review fix: Retire the pipeline guard in manual settings reconciliation

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Make the suite installer's no-JSON-tool fallback explicitly remove only the retired pipeline-guard Stop hook while preserving every unrelated and custom hook entry.

## Solution summary

The manual fallback now mirrors automated semantics without mutating settings: remove only inner Stop-hook objects whose string command contains the retired guard path, preserve custom/unrelated neighbors, then merge current core hooks.

## What worked

- A mixed-wrapper manual fixture proved the installer leaves settings byte-exact and made custom-neighbor preservation an explicit fallback contract.
- Keeping the fallback as one shared instruction variable updated both preview and success output without touching automated reconciliation.

## What didn't work

- The first path wording stopped at the inner `hooks` array rather than its `[*]` objects, so an exact-string test froze an instruction that could still be misapplied destructively.

## Worth knowing

- Automated jq/Python paths already remove individual matching inner objects and prune only empty wrappers. REQ-155 is held for consent to make the manual JSON path and empty-wrapper wording equally exact.

## Back-reference

See `do-work/archive/UR-031/REQ-151-retire-pipeline-guard-manual-settings-reconciliation.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1ab5ed8`.
