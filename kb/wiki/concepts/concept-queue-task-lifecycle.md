---
title: "Queue Orchestration and Task Lifecycle"
type: concept
topic_cluster: queue-orchestration-and-lifecycle
sources:
  - raw/processed/2026-09-01/REQ-001-code-review-split-actions-work-md-into-o.md
  - raw/processed/2026-09-01/REQ-012-add-do-work-note-command-for-lightweight.md
  - raw/processed/2026-09-01/REQ-014-add-crew-members-maintenance-md-codifyin.md
  - raw/processed/2026-09-01/REQ-060-no-mechanism-resolves-a-failed-req-so-a.md
  - raw/processed/2026-09-01/REQ-079-two-guards-pin-the-weaker-fingerprint-of.md
  - raw/processed/2026-09-01/REQ-080-the-capture-template-emits-a-stray-instr.md
  - raw/processed/2026-09-01/REQ-081-next-version-ignores-flags-placed-after.md
  - raw/processed/2026-09-01/REQ-091-the-hand-back-merge-fails-while-the-owne.md
  - raw/processed/2026-09-01/REQ-102-scope-work-md-step-10-preserve-rules-to-.md
  - raw/processed/2026-09-01/REQ-147-addendum-reserve-request-numbers-during-.md
  - raw/processed/2026-09-01/REQ-151-review-fix-retire-the-pipeline-guard-in-.md
  - raw/processed/2026-09-01/REQ-155-review-fix-correct-the-manual-stop-hook-.md
  - raw/processed/2026-09-01/REQ-193-keep-archived-urs-closed-during-standalo.md
  - raw/processed/2026-09-01/REQ-288-fix-the-three-unfiled-contradictions-in-.md
  - raw/processed/2026-09-01/REQ-291-browser-behavior-probe-lane-beside-the-n.md
  - raw/processed/2026-09-01/REQ-299-review-fix-carry-builder-authored-sectio.md
  - raw/processed/2026-09-01/REQ-389-addendum-mark-spliced-paste-titles-with-.md
related:
  - page: concept-completed-work-presentation
    rel: complements
  - page: concept-modular-suite-architecture
    rel: complements
  - page: concept-queue-task-lifecycle
    rel: complements
  - page: concept-queue-task-lifecycle
    rel: complements
  - page: concept-session-checkpoints-and-recovery
    rel: complements
  - page: concept-timestamp-and-metadata-governance
    rel: complements
  - page: concept-worktree-isolation-and-parallelism
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: high
---

# Queue Orchestration and Task Lifecycle

Architectural overview and synthesis for the Queue Orchestration and Task Lifecycle subsystem in the do-work suite.

## Key Principles & Synthesized Lessons

This cluster synthesizes evidence from 17 source documents:

- [[REQ-001-code-review-split-actions-work-md-into-o]] — Code review: split actions/work.md into orchestrator + reference companion
- [[REQ-012-add-do-work-note-command-for-lightweight]] — Add do-work note command for lightweight roadmap notes
- [[REQ-014-add-crew-members-maintenance-md-codifyin]] — Add crew-members/maintenance.md codifying delete-before-you-add
- [[REQ-060-no-mechanism-resolves-a-failed-req-so-a]] — No mechanism resolves a `failed` REQ, so a UR containing one can never close
- [[REQ-079-two-guards-pin-the-weaker-fingerprint-of]] — Two guards pin the weaker fingerprint of the premise they exist to retire
- [[REQ-080-the-capture-template-emits-a-stray-instr]] — The capture template emits a stray instruction line into every REQ it produces
- [[REQ-081-next-version-ignores-flags-placed-after]] — next-version ignores flags placed after the bump size and silently bumps the calling repo
- [[REQ-091-the-hand-back-merge-fails-while-the-owne]] — The hand-back merge fails while the owner's claim bookkeeping is staged, on any install that tracks do-work/
- [[REQ-102-scope-work-md-step-10-preserve-rules-to-]] — Scope work.md Step 10 preserve rules to every non-own entry, and pin both label-destruction paths
- [[REQ-147-addendum-reserve-request-numbers-during-]] — Addendum: reserve request numbers during allocation
- [[REQ-151-review-fix-retire-the-pipeline-guard-in-]] — Review fix: Retire the pipeline guard in manual settings reconciliation
- [[REQ-155-review-fix-correct-the-manual-stop-hook-]] — Review fix: Correct the manual Stop-hook object path
- [[REQ-193-keep-archived-urs-closed-during-standalo]] — Keep archived URs closed during standalone review
- [[REQ-288-fix-the-three-unfiled-contradictions-in-]] — Fix the three unfiled contradictions in clarify's Step 4
- [[REQ-291-browser-behavior-probe-lane-beside-the-n]] — Browser behavior probe lane beside the Node behavior lane
- [[REQ-299-review-fix-carry-builder-authored-sectio]] — Review fix: carry builder-authored sections past Step 8, starting with ## Decisions
- [[REQ-389-addendum-mark-spliced-paste-titles-with-]] — Addendum: mark spliced paste titles with a leading arrow

## Cross-References

See related system components and verification gates.
