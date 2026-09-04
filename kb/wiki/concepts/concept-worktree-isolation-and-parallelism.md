---
title: "Worktree and Parallel Dispatch"
type: concept
topic_cluster: worktree-and-parallel-dispatch
sources:
  - raw/processed/2026-09-04/REQ-458-addendum-classify-active-worktrees-as-pr.md
  - raw/processed/2026-09-01/REQ-036-re-validate-write-set-disjointness-when.md
  - raw/processed/2026-09-01/REQ-037-place-the-worktree-merge-in-the-step-seq.md
  - raw/processed/2026-09-01/REQ-073-fan-out-dispatch-n-concurrent-builders-u.md
  - raw/processed/2026-09-01/REQ-075-five-files-still-explain-write-set-s-di.md
  - raw/processed/2026-09-01/REQ-082-the-fan-out-hand-back-file-has-no-legal.md
  - raw/processed/2026-09-01/REQ-083-verify-reports-every-builder-worktree-as.md
  - raw/processed/2026-09-01/REQ-085-run-req-073-s-live-two-builder-acceptanc.md
  - raw/processed/2026-09-01/REQ-092-actions-work-md-has-no-wave-selection-or.md
  - raw/processed/2026-09-01/REQ-096-execution-model-re-grain-claim-anywhere-.md
  - raw/processed/2026-09-01/REQ-099-automatic-wave-dispatch-the-work-loop-co.md
  - raw/processed/2026-09-01/REQ-101-docs-adr-multi-checkout-guide-and-the-se.md
  - raw/processed/2026-09-01/UR-018-parallel-building-batch-session-traps.md
related:
  - page: concept-queue-task-lifecycle
    rel: extends
created: 2026-09-01
updated: 2026-09-02
confidence: high
---

# Worktree and Parallel Dispatch

Architectural overview and synthesis for the Worktree and Parallel Dispatch subsystem in the do-work suite.

## Key Principles & Synthesized Lessons

This cluster synthesizes evidence from 12 source documents:

- [[REQ-036-re-validate-write-set-disjointness-when]] — Re-validate write-set disjointness when Step 5.5 firms the sets
- [[REQ-037-place-the-worktree-merge-in-the-step-seq]] — Place the worktree merge in the step sequence and re-point the evidence-consuming steps
- [[REQ-073-fan-out-dispatch-n-concurrent-builders-u]] — Fan-out dispatch: N concurrent builders under one queue owner
- [[REQ-075-five-files-still-explain-write-set-s-di]] — Five files still explain write_set's display-only status with a reason fan-out made false
- [[REQ-082-the-fan-out-hand-back-file-has-no-legal]] — The fan-out hand-back file has no legal write location
- [[REQ-083-verify-reports-every-builder-worktree-as]] — verify reports every builder worktree as a fixable orphan, including active and unmerged ones
- [[REQ-085-run-req-073-s-live-two-builder-acceptanc]] — Run REQ-073's live two-builder acceptance test and record what it found
- [[REQ-092-actions-work-md-has-no-wave-selection-or]] — actions/work.md has no wave-selection or launch-before-wait path, so documented fan-out concurrency cannot be reached by following it
- [[REQ-096-execution-model-re-grain-claim-anywhere-]] — Execution-model re-grain: claim anywhere, one releaser; dispatch widened to any tree
- [[REQ-099-automatic-wave-dispatch-the-work-loop-co]] — Automatic wave dispatch — the work loop computes and dispatches the ready set
- [[REQ-101-docs-adr-multi-checkout-guide-and-the-se]] — Docs + ADR — multi-checkout guide and the session-ownership decision record
- [[UR-018-parallel-building-batch-session-traps]] — Traps the parallel-building batch session already hit
- [[REQ-458-addendum-classify-active-worktrees-as-pr]] — Addendum: classify active worktrees as present and non-fixable

## Cross-References

See related system components and verification gates.
