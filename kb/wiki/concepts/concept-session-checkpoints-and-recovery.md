---
title: "Checkpoint and Crash Recovery"
type: concept
topic_cluster: checkpoint-and-crash-recovery
sources:
  - raw/processed/2026-09-04/REQ-457-record-cleanup-move-destinations-after-e.md
  - raw/processed/2026-09-04/REQ-489-remove-whole-checkpoint-entries-when-a-r.md
  - raw/processed/2026-09-01/REQ-035-represent-concurrent-claims-in-the-orche.md
  - raw/processed/2026-09-01/REQ-071-crash-recovery-must-respect-a-live-claim.md
  - raw/processed/2026-09-01/REQ-077-crash-recovery-s-own-crash-branch-is-unr.md
  - raw/processed/2026-09-01/REQ-086-the-in-progress-record-s-rule-is-unstate.md
  - raw/processed/2026-09-01/REQ-094-checkpoint-writer-label-crash-recovery-i.md
  - raw/processed/2026-09-01/REQ-095-two-clone-acceptance-run-checkpoint-pois.md
  - raw/processed/2026-09-01/REQ-104-label-less-checkpoint-entries-locally-mo.md
  - raw/processed/2026-09-01/REQ-108-review-fix-in-progress-record-still-enum.md
  - raw/processed/2026-09-01/REQ-109-work-md-session-start-note-still-enumera.md
  - raw/processed/2026-09-01/REQ-166-simplify-session-start-hook-and-fix-dead.md
  - raw/processed/2026-09-01/REQ-246-repair-detectably-wrong-queue-and-workin.md
  - raw/processed/2026-09-01/REQ-256-disclose-the-session-hook-s-queue-write-.md
  - raw/processed/2026-09-01/REQ-274-retire-the-the-sessionstart-hook-exits-n.md
related:
  - page: concept-queue-task-lifecycle
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: high
---

# Checkpoint and Crash Recovery

Architectural overview and synthesis for the Checkpoint and Crash Recovery subsystem in the do-work suite.

## Key Principles & Synthesized Lessons

This cluster synthesizes evidence from 15 source documents:

- [[REQ-035-represent-concurrent-claims-in-the-orche]] — Represent concurrent claims in the orchestrator lock and Crash Recovery gate
- [[REQ-071-crash-recovery-must-respect-a-live-claim]] — Crash recovery must respect a live claim before stripping and re-queueing
- [[REQ-077-crash-recovery-s-own-crash-branch-is-unr]] — Crash recovery's own-crash branch is unreachable, and its retired premise survives in the same file
- [[REQ-086-the-in-progress-record-s-rule-is-unstate]] — The in-progress record's rule is unstated at the two out-of-pipeline movers and contradicted in the user guide
- [[REQ-094-checkpoint-writer-label-crash-recovery-i]] — Checkpoint writer label — crash recovery ignores foreign entries
- [[REQ-095-two-clone-acceptance-run-checkpoint-pois]] — Two-clone acceptance run — checkpoint poisoning repro and claim-conflict evidence
- [[REQ-104-label-less-checkpoint-entries-locally-mo]] — Label-less checkpoint entries — "locally modified" is not evidence of authorship
- [[REQ-108-review-fix-in-progress-record-still-enum]] — Review fix: In-Progress Record still enumerates two recovery cases and owes no removal rule for a label-less entry
- [[REQ-109-work-md-session-start-note-still-enumera]] — work.md session-start note still enumerates the recovery case list and calls a label-less entry a foreign claim
- [[REQ-166-simplify-session-start-hook-and-fix-dead]] — Simplify session-start hook and fix dead fail-soft fallback
- [[REQ-246-repair-detectably-wrong-queue-and-workin]] — Repair detectably wrong queue and working timestamps from the session hook
- [[REQ-256-disclose-the-session-hook-s-queue-write-]] — Disclose the session hook's queue write surface in the docs
- [[REQ-274-retire-the-the-sessionstart-hook-exits-n]] — Retire the "the SessionStart hook exits nonzero" framing where it is still stated
- [[REQ-457-record-cleanup-move-destinations-after-e]] — Record cleanup move destinations after exclusive creation
- [[REQ-489-remove-whole-checkpoint-entries-when-a-r]] — Remove whole checkpoint entries when a REQ leaves working

## Cross-References

See related system components and verification gates.
