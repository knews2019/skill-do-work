---
title: "Topic Index: Workflow Orchestration"
type: topic-index
status: reference
topic_cluster: workflow-orchestration
sources:
  - CHANGELOG.md
  - CHANGELOG.md (0.51.5 The Full Send)
  - CHANGELOG.md (0.56.0 The Clean Sweep)
  - CHANGELOG.md (0.56.1 The Safety Net)
  - CHANGELOG.md (0.56.2 The Tight Scope)
  - CHANGELOG.md (0.60.3 The Paved Path)
  - CHANGELOG.md (0.60.5 The Honest Mirror)
  - CLAUDE.md
  - README.md
  - actions/capture.md
  - actions/cleanup.md
  - actions/review-work.md
  - actions/work.md
  - docs/work-guide.md
related:
  - page: adr-004-canonicalize-pending-reqs-under-do-work-queue
    rel: complements
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: complements
  - page: adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles
    rel: complements
  - page: adr-014-considered-declined-autonomous-loop-until-done
    rel: complements
  - page: adr-015-load-maintenance-crew-via-req-marker
    rel: complements
  - page: adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser
    rel: complements
  - page: adr-022-one-broken-pipe-does-not-stop-the-factory
    rel: complements
created: 2026-04-15
updated: 2026-09-02
confidence: high
---

# Workflow Orchestration

How pending work is stored and how the work orchestrator coordinates queue processing.

## ADRs

- [[adr-004-canonicalize-pending-reqs-under-do-work-queue]] — [ADR-004](../records/adr-004-canonicalize-pending-reqs-under-do-work-queue.md): Treat `do-work/queue/` as the canonical home for pending REQ files and update every workflow around that assumption.
- [[adr-005-pipeline-is-stateful-and-resumable]] — [ADR-005](../records/adr-005-pipeline-is-stateful-and-resumable.md) (**superseded**): Historical contract for the stateful orchestrator retired by ADR-019 and REQ-145.
- [[adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles]] — [ADR-006](../records/adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles.md) (**superseded**): Historical bounded-continuation contract retired with the stateful orchestrator.
- [[adr-014-considered-declined-autonomous-loop-until-done]] — [ADR-014](../records/adr-014-considered-declined-autonomous-loop-until-done.md) (**declined**): Do not re-add the `ultracode-fable` / loop-until-done workflow — its model-agnostic capabilities already survive as canon, and the model-specific tier table is intentionally out of scope.
- [[adr-015-load-maintenance-crew-via-req-marker]] — [ADR-015](../records/adr-015-load-maintenance-crew-via-req-marker.md): Load `crew-members/maintenance.md` in work.md Step 6 via a `maintenance: true` REQ marker (set by capture for removal findings on the skill's own instructions) — marker-only, no description heuristic. Resolves ADR-014/REQ-014's deferred D-01 loader gap.
- [[adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser]] — [ADR-018](../records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md): Re-grain ownership from `one queue owner per checkout` to `one releaser per queue` — any checkout may capture, claim and build; exactly one runs the release tail. Advisory `assigned_to` field (reserve verb/status stay dead), a static checkpoint writer label that is not liveness machinery, capture-anywhere with fix-at-merge, and `--fan-out` auto-wave superseding "nothing computes the set". Partially reverses the 0.161.0 exclusive-session decision, which had no ADR of its own.
- [[adr-022-one-broken-pipe-does-not-stop-the-factory]] — [ADR-022](../records/adr-022-one-broken-pipe-does-not-stop-the-factory.md): Set aside failures local to one REQ and keep draining unrelated runnable work; stop only on shared-target dirt, with `do-work run-with-recovery` as the deliberate sole-authority resolution path.

## Cross-Cluster Links

- [[adr-005-pipeline-is-stateful-and-resumable]] and [[adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles]] are superseded alongside [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] in [[_index_pipeline-deliverables]].
