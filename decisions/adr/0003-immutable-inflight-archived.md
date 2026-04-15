---
id: 0003
title: Immutable Working and Archive Folders
status: accepted
decided: 2026-02-01
version: 0.6.0
topic: queue-model
supersedes: []
superseded_by: null
related:
  - adr: 0002
    rel: extends
  - adr: 0010
    rel: complements
---

# ADR-0003: Immutable Working and Archive Folders

## Context

Once a builder claimed a REQ and started working on it, users would sometimes edit the REQ file to "add one more thing." Once work was archived, users (or the agent) would occasionally reach back in and tweak historical REQs to reflect how things ended up.

Both patterns quietly broke core assumptions. In-flight edits created race conditions: the builder had already loaded an older version of the REQ and would produce code matching that older version, then the review would compare against the new version and flag a phantom gap. Archive edits destroyed the audit trail — the record of "what was actually asked" got rewritten to match "what actually happened," erasing the gap that makes review valuable.

The deeper issue was that mutability of in-flight and archived REQs made every downstream assumption about "this REQ says X" depend on the time of reading. Two actions reading at different moments could see different REQs, and nothing in the skill caught the divergence.

## Decision

Files inside `do-work/working/` and `do-work/archive/` are **immutable**. Once a UR has left the queue, neither users nor the agent modify its REQ files.

Follow-ups to an in-flight or archived REQ are handled through the `addendum_to:` frontmatter field: a new REQ is created in the queue referencing the original by id. The addendum chain preserves the intent trail — you can reconstruct the full sequence of asks by following `addendum_to` links from any REQ back to the original.

## Alternatives Considered

- **Allow mid-flight edits.** Users could change their minds as work progresses. Rejected — creates the race conditions described above and makes review incoherent.
- **Lockfiles on in-flight REQs.** Block edits with a file lock. Rejected — solves the technical race but not the conceptual problem that mutable REQs break the intent trail.
- **Versioned REQs.** Every edit creates a new version in-place. Rejected — addendum-as-new-REQ is conceptually cleaner and integrates with the existing queue model instead of adding a parallel version system.

## Consequences

- **Clean boundary per phase.** A REQ in `queue/` is mutable (not yet claimed). A REQ in `working/` or `archive/` is immutable (no longer editable).
- **Addendum chains preserve the intent trail.** Reviewers can see not just what was delivered, but how the ask evolved over time.
- **Capture handles follow-ups specifically.** The capture action checks the location of any REQ it might duplicate — in-queue gets deduplicated or merged, in-flight/archived gets a new addendum REQ.
- **Cost: more REQ files over a project's life.** Each tiny "one more thing" becomes its own REQ. The `cleanup` action's consolidation passes keep the resulting sprawl in check.

## References

- **CHANGELOG**: v0.6.0 — The Bouncer (2026-02-01)
- **Action files**: `actions/capture.md` (addendum handling), `actions/work.md` (immutability docs in folder section)
- **Related ADRs**: [[0002-ur-req-pairing]] (immutability protects the UR+REQ pairing after claim), [[0010-reqs-as-validated-intent]] (validated intent can only be validated once — after that, it's historical fact)
