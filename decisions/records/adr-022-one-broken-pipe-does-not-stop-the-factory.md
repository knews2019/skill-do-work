---
title: "ADR-022: One Broken Pipe Does Not Stop the Factory"
type: architecture-decision-record
status: accepted
topic_cluster: workflow-orchestration
decided: 2026-09-02
sources:
  - do-work/archive/REQ-494-complete-already-green-repository-gate-repair-lifecycle.md
  - do-work/archive/UR-096/REQ-498-make-orchestrator-finalization-resumable.md
  - do-work/archive/REQ-499-add-recover-finalization-assume-sole-releaser.md
  - do-work/archive/REQ-500-surface-unfinished-finalizations-in-doctor-and-session-start.md
  - do-work/archive/UR-097/REQ-501-add-run-with-recovery-and-record-one-broken-pipe-principle.md
  - skills/do-work/actions/run-with-recovery.md
related:
  - page: adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser
    rel: extends
created: 2026-09-02
updated: 2026-09-02
confidence: high
---

# ADR-022: One Broken Pipe Does Not Stop the Factory

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser]] (extends)

## Context

REQ-494 exposed a failure mode larger than its own implementation. Its session died after archive work had started but before the implementation, lifecycle metadata, and release were committed. Six hours later the next `do-work run` reached a different REQ and its claim refused on the dirty shared `do-work/CHECKPOINT.md`. The first broken tail had stopped the whole queue even though most remaining work did not depend on it.

The same pressure had already appeared in two solution families. REQ-468 through REQ-472 tried to isolate a REQ, set it aside, route its failure, update readers, and prove the continuing loop. REQ-491 and REQ-492 replaced that blocked-state design with canonical repository-gate deferral and dependency-backed repair. REQ-498 made finalization resumable; REQ-499 added the narrow sole-releaser assertion for ambiguous shared metadata; REQ-500 made unfinished finalizations visible. REQ-501 completes the user-facing recovery path and gives the common principle a stable home.

## Decision

**One broken pipe does not stop the rest of the factory. Whenever a failure is local to one REQ, record a typed finding, set that REQ aside, and continue draining unrelated runnable work.** The instances above are evidence for the condition, not a closed list of special cases.

The boundary is shared targets. A dirty checkpoint cannot be set aside with one REQ because every claim writes it. The same is true of any lifecycle or release target the next unit must mutate. Shared-target dirt may stop the loop until it is resolved mechanically or the user explicitly asserts ownership; the stop must name the verb that resolves it.

That explicit path is `do-work run-with-recovery`. It first invokes canonical finalization discovery with `--assume-sole-releaser`, then treats every `do-work/working/` claim as this checkout's interrupted work, and finally hands the original arguments to ordinary `run`. It is a separate verb, not a `run` flag, because the sole-writer assertion must be deliberate and must not leak into scripts or habits that rely on plain `run` preserving foreign claims.

The assertion does not weaken ADR-018's default. Plain `run` remains safe for a shared queue; `run-with-recovery` is correct only when the user knows this checkout is the sole writer and releaser for the invocation. Secret-classified and project paths remain outside the widened finalization authority.

## Alternatives

1. **Add `--assume-sole-releaser` directly to `do-work run`.** Declined. A flag is easy to retain in a script or copied command after its ownership premise stops being true. A named verb makes the changed safety boundary visible at routing, help, aliases, and next-step guidance.
2. **Let `claim` tolerate a dirty checkpoint and keep going.** Declined. The same guard also prevents a claim from writing through a half-resolved merge or another incomplete shared transaction. Weakening it would turn an explicit ownership decision into an implicit guess at every claim.
3. **Stop on every REQ failure.** Declined. It gives a local defect queue-wide reach and recreates the REQ-494 incident whenever a new failure shape appears.
4. **Automatically take authority whenever recovery finds stale work.** Declined. Age is not ownership; a long-running build can be live. The user chooses the authority verb.

## Consequences

Orchestration now has one reviewable rule for continuation: local failures are isolated and reported; shared-target dirt is resolved before more claims. The recovery verb gives a fresh session a prompt-free route through known ownership refusals without changing ordinary shared-queue behavior.

The cost is one additional public verb and its routing/help surface. That cost is deliberate: the name carries the safety assertion. Future orchestration changes are judged against the condition rather than added to an inventory of known failure codes.

## Instances

- REQ-468 through REQ-472 — per-REQ isolation and the superseded blocked set-aside family.
- REQ-491 and REQ-492 — canonical repository-gate deferral and its `do-work run` integration.
- REQ-498 — resumable finalization and legacy-tail discovery.
- REQ-499 — sole-releaser attribution for ambiguous shared finalization metadata.
- REQ-500 — unfinished-finalization diagnostics and session-start visibility.
- REQ-501 — `run-with-recovery`, authority crash recovery, and this decision record.

## References

- [run-with-recovery.md](../../skills/do-work/actions/run-with-recovery.md) — the deliberate authority verb
- [work-reference.md](../../skills/do-work/actions/work-reference.md) — Execution Model and Crash Recovery semantics
- [work.md](../../skills/do-work/actions/work.md) — the unchanged work pipeline that receives the handoff
- [adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md](./adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md) — the default shared-queue ownership model this decision extends
