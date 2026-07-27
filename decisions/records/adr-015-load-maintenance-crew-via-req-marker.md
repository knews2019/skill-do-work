---
title: "ADR-015: Load the Maintenance Crew via a `maintenance: true` REQ Marker"
type: architecture-decision-record
status: accepted
topic_cluster: workflow-orchestration
decided: 2026-06-30
sources:
  - external review finding (chatgpt-codex-connector, P2 — "Load maintenance rules when running removal work")
  - do-work/archive/UR-002/REQ-014-delete-before-you-add-rule.md (D-01 + tracked-not-fixed loader gap)
  - crew-members/maintenance.md
  - actions/work.md (Step 6 crew loader)
  - actions/work-reference.md (Schema Read Contract)
  - actions/capture.md
  - actions/quick-wins.md
related:
  - page: adr-014-considered-declined-autonomous-loop-until-done
    rel: complements
  - page: adr-003-always-load-karpathy-guardrails
    rel: complements
created: 2026-06-30
updated: 2026-06-30
confidence: high
---

# ADR-015: Load the Maintenance Crew via a `maintenance: true` REQ Marker

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-014-considered-declined-autonomous-loop-until-done]] (complements), [[adr-003-always-load-karpathy-guardrails]] (complements, in [[_index_skill-architecture]])

## Context

REQ-014 (0.98.0 "The Delete Key", commit `db4d661`) shipped `crew-members/maintenance.md` — the delete-before-you-add discipline for deliberate maintenance passes — but wired it as *referenced, not loaded*. Its decision **D-01** chose to point `actions/quick-wins.md` at the rule rather than have any action auto-load it, and the REQ's review tracked the consequence as an explicit not-fixed item: the real subtraction work runs through `actions/work.md` Step 6, whose crew loader (general+karpathy always; domain/testing/security/caveman conditional) had no hook for `maintenance.md`. So a removal REQ ran under `karpathy.md`'s *implementation-time* surgical-changes posture — which says **don't** delete adjacent dead code — exactly the opposite of what the maintenance pass intends.

The deferral was deliberate, not an oversight: REQ-014 recorded that adding `maintenance.md` to Step 6's always/domain/flag loader "would misfire on every implementation REQ," and left the real load to "a future `maintain` action." An external review (chatgpt-codex-connector, P2) independently flagged the same enforcement gap and asked for "a concrete load trigger **or REQ marker** for deliberate maintenance passes." This ADR records the decision to close the gap now, via the marker half of that suggestion, rather than wait for a dedicated `maintain` action that may never land.

## Decision

**Add a `maintenance: true` REQ frontmatter marker that gates a conditional load of `crew-members/maintenance.md` in `actions/work.md` Step 6.** The marker is set by `actions/capture.md` when a request removes or narrows the skill's *own* operating instructions (a drifting agent/action/crew/prime file) — e.g. acting on a `do-work quick-wins` removal finding for a redundant rule or over-broad config.

Two properties are load-bearing:

- **Marker-only, no description heuristic.** Unlike `security.md` (which also loads on a heuristic scan of the REQ description), `maintenance.md` loads *only* on the explicit marker. A heuristic would fire on ordinary implementation REQs — which routinely touch adjacent dead code — and load the opposite posture from the one karpathy wants. This is the precise misfire REQ-014 warned against, so the heuristic is deliberately omitted.
- **Scoped to instruction-like artifacts, not all removals.** The discipline governs the skill's own rules/config/instructions, not arbitrary application source. Plain dead-code removal in app code stays under karpathy's implementation-time rule with no marker. `actions/quick-wins.md`'s removal-findings rule was tightened in the same change to stop lumping app dead-code together with rule/config removal.

The marker follows the existing `tdd`/`caveman` pattern: a boolean field in the REQ schema, normalized via the **Schema Read Contract** (`actions/work-reference.md`), emitted canonically by capture, read at Step 6. `maintenance.md` loads *alongside* the normal crew (general, karpathy, domain), never instead of it.

## Alternatives

1. **Add `maintenance.md` to Step 6's always/domain/flag loader directly (the review's "global load trigger").** Rejected: it loads delete-before-you-add on every implementation REQ, contradicting karpathy's "don't delete adjacent code" — the exact misfire REQ-014 declined.

2. **Reuse `domain: maintenance`** (Step 6 already loads `crew-members/[domain].md` when the file exists). Rejected: it conflates two orthogonal axes — *subject area* (domain) vs *kind of pass* (maintenance vs feature) — and steals the domain slot, so the REQ can no longer also load its real technical domain. `maintenance.md` is explicitly "loaded alongside domain rules," so it must be an independent flag, exactly as `tdd` is independent of `domain`.

3. **Description heuristic, mirroring `security.md`.** Rejected: see "marker-only" above — the implementation-time/maintenance-time boundary is too easy to cross on a keyword match, and the cost of a wrong load here is a contradictory posture, not a cheap extra checklist.

4. **Do nothing — keep it referenced-not-loaded until a `maintain` action exists** (REQ-014's interim state). Rejected: a shipped guardrail that never fires in its target scenario is unenforced canon. The marker closes the gap with one schema field and one loader bullet, and a future `maintain` action simply sets the same marker — no rework.

## Consequences

A deliberate maintenance pass captured through the pipeline now actually loads the discipline written for it, without polluting ordinary implementation REQs. The skill gains one schema field (`maintenance`), one Step 6 loader bullet, one Schema Read Contract row, and a capture rule for setting it — the minimal surface that satisfies the replay case (a removal REQ that loaded the wrong posture before, the right one after). REQ-014's D-01 is resolved and its tracked not-fixed item closed; the review finding is accepted in its marker form and its global-trigger form stays declined for the reason above. A future dedicated `maintain` action, if it lands, reuses the marker rather than introducing a second load path.

## References

- [crew-members/maintenance.md](../../crew-members/maintenance.md) — the discipline + its JIT_CONTEXT trigger
- [actions/work.md](../../actions/work.md) — Step 6 conditional load (bullet 5a)
- [actions/work-reference.md](../../actions/work-reference.md) — Schema Read Contract `maintenance` row
- [actions/capture.md](../../actions/capture.md) — marker emission rule
- [actions/quick-wins.md](../../actions/quick-wins.md) — tightened removal-findings rule
- `do-work/archive/UR-002/REQ-014-delete-before-you-add-rule.md` — D-01 and the tracked loader gap this ADR closes
