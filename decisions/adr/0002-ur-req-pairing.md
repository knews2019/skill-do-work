---
id: 0002
title: UR + REQ Pairing on Every Capture
status: accepted
decided: 2026-02-01
version: 0.4.0
topic: queue-model
supersedes: []
superseded_by: null
related:
  - adr: 0001
    rel: complements
  - adr: 0003
    rel: depends-on
  - adr: 0008
    rel: depends-on
---

# ADR-0002: UR + REQ Pairing on Every Capture

## Context

The earliest versions of the skill produced REQ files directly at a flat path. Each captured request became a single REQ — no grouping, no parent context.

This broke down as soon as anyone asked for something non-trivial. A user request like "redesign the settings page and wire up the new API for preferences" naturally decomposes into several REQs (redesign, API contract, integration). Filed flat, those REQs had no shared context — no place to store the umbrella ask, no way for the builder to see related work, no natural grouping for downstream actions like review or present.

The symptoms were telling: pipeline-style actions couldn't "operate on this UR" because there was no such thing as a UR. Reviewers couldn't ask "did we build *all* of what was asked?" because the asks were fragments without a whole. Cleanup routines had no group to close out.

## Decision

Every invocation of the capture action produces **two** coordinated artifacts:

1. A **User Request (UR) folder** with a numeric id (`UR-NNN`) holding the shared context for the request.
2. One or more **REQ files** inside that UR folder, each a unit of work the builder processes independently.

The UR folder and its REQ files are created together in the same capture pass — never a REQ without a UR, never a UR without at least one REQ. v0.8.0 elevated this pairing to a structural invariant by adding a "Required Outputs" section to the top of the capture action: UR + REQ pairing stated upfront as mandatory, visible even to skimming readers.

## Alternatives Considered

- **Flat REQ-only structure.** Simpler on disk, but loses the grouping that everything downstream relies on. Rejected — the friction from losing the grouping exceeded the friction of creating the UR folder.
- **Nested folders per topic (no UR concept).** Users would organize REQs under topic folders they invent. Rejected — too much cognitive overhead per capture, and "topic" is subjective.
- **A single long request file.** Capture writes one monolithic REQ that contains everything. Rejected — loses the unit-of-work granularity that lets the work action process pieces independently.

## Consequences

- **Every other action can assume a UR exists.** The pipeline action dispatches sub-steps by UR id. The review-work and present-work actions operate at UR granularity. The cleanup action closes URs, not loose REQs.
- **Natural container for multi-REQ work.** Related REQs share a folder, making it obvious which REQs move together.
- **Sample REQ/UR structure lives in the skill** (`actions/sample-archived-req.md`) to show agents and users what a completed pairing looks like.
- **Cost: two ids to track.** Users have to understand both `UR-NNN` and `REQ-NNN`. The skill mitigates this by exposing UR ids in user-facing reports and keeping REQ numbering sequential across the project.

## References

- **CHANGELOG**: v0.4.0 — The Organizer (2026-02-01) introduced the UR system; v0.8.0 — The Bright Light (2026-02-03) made the pairing structurally unmissable.
- **Action files**: `actions/capture.md` (Required Outputs), `actions/work.md` (folder moves), `actions/cleanup.md` (UR closure)
- **Related ADRs**: [[0001-capture-execute-boundary]] (two-phase workflow needs persistent artifacts), [[0003-immutable-inflight-archived]] (immutability protects the pairing after claim), [[0008-queue-canonical-path]] (pairing lives at `do-work/queue/UR-NNN/REQ-NNN.md`)
