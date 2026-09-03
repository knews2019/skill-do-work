---
id: UR-091
title: 'Add a same-UR verification-append destination to the Fold-First Rule'
created_at: 2026-09-01T12:11:03Z
requests: [REQ-484]
---

# Add a Same-UR Verification-Append Destination to the Fold-First Rule

Add a new destination to the Fold-First Rule (skills/do-work/actions/capture-reference.md § Fold-First Rule), between destination 2 (matching non-sweep REQ → convert) and destination 3 (prose backlog): the same-UR verification append.

When a finding comes from work or review inside UR-N, and UR-N has a pending, unclaimed (assigned_to absent) downstream member whose stated scope is verification, parity, or acceptance over the code the finding names (typically a serial spine's terminal REQ), the finding folds there as an acceptance-criteria append: a `## Folded From REQ-NNN / <source> (YYYY-MM-DD)` section restating the finding as acceptance criteria or named fixtures, citing the source REQ/review and the fold date — never the destination-2 sweep conversion; the destination keeps its write_set, estimate, frontmatter, and non-sweep shape.

Conditions, all required:
(a) the finding is goal-shaped, not defect-shaped — it restates, refines, or extends the destination's acceptance goals rather than reporting behavior broken in what currently ships. "Shipped behavior" means what currently ships to consumers: a divergence in not-yet-authoritative migration code, whose legacy implementation still ships until the destination REQ converts it, is goal-shaped relative to the migration's acceptance. A shipped-behavior defect never folds behind a gate regardless of impact — a fix must never wait behind a gate — and is minted standalone as today.
(b) the destination's dependency chain is alive — every depends_on member is terminal-successful or present in the queue/working set; a failed or cancelled member keeps the destination ineligible.
(c) the finding's judged impact is not impact-critical — critical keeps current behavior (own REQ, status: pending, prominent pierce line) at any depth.

Each accepted append is recorded in the destination body with the fold date and source, so a stalled chain carrying folds stays visible to review-work Step 10 and to any board surface that shows gated REQs. Update every flow that restates fold-destination eligibility — actions/review-work.md Step 10 and actions/work-reference.md's follow-up flows cite the Fold-First Rule by name; verify none restates the gate — and add a contract-regression predicate pinning conditions (a)–(c) if the builder judges the wording load-bearing.

Evidence (RED case): REQ-414's fresh re-review findings 3 and 4 were minted as standalone REQ-464/REQ-465 solely because REQ-420 — the UR's terminal parity REQ, whose differential testing is where their substance belongs — was dependency-gated and therefore no legal destination existed; the maintainer folded both into REQ-420 by hand as `## Folded From` acceptance criteria and cancelled them (commit 593c5145). GREEN when the same scan routes such a finding through the new destination automatically, appending without converting, and destination 2's conversion path remains unchanged for dependency-ready same-root-cause matches. This request replaces cancelled REQ-480, which loosened destination 2's conversion gate instead — the conversion mechanic and root-cause match test made it unable to fire on its own RED case.
