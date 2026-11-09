---
title: "Lessons from REQ-170: Finding-closure ratchet and canonical earned-defense rubric"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-170-finding-closure-ratchet-and-canonical-ea.md]
related:
  - page: REQ-165-shell-block-lint-harness-for-shipped-act
    rel: complements
  - page: REQ-166-simplify-session-start-hook-and-fix-dead
    rel: complements
  - page: REQ-167-deduplicate-copy-pasted-shell-primitives
    rel: complements
  - page: REQ-168-delete-or-test-audit-of-defensive-code-i
    rel: complements
  - page: REQ-169-validate-feedback-flags-remedies-that-ad
    rel: complements
  - page: REQ-171-addendum-promote-prescribed-shell-primit
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-170: Finding-closure ratchet and canonical earned-defense rubric

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

Two small, single-home rules that make the review loop converge instead of plateau:

1. **Finding-closure ratchet** — a REQ that originates from a review or triage finding may only close with either a named regression test (fails before the fix, passes after) or a deletion of the surface the finding lived in — never a bare patch. Canonical text lives in `skills/do-work/actions/work-reference.md` beside the other pipeline contracts; `actions/review-work.md` enforces it at the gate (closure without test-or-deletion evidence bounces); `actions/capture.md` sharpens the GREEN criterion for finding-origin REQs to name the test or the deletion.
2. **Earned-defense rubric** — one canonical paragraph in `skills/do-work/crew-members/coding-guardrails.md`: defensive code must name the incident that earned it, and the fix must cost less surface than the risk it covers — user's wording preserved: *"what earned this, and is the fix still cheaper than the surface it added?"* The rubric already shipped inside `../do-work-toolbox/actions/validate-feedback.md` (REQ-169, commit 063bb88) — condense that to its triage-specific application (Surface-cost verdicts, Accept bar) plus a one-line citation of the canonical paragraph; review-work's gate cites it too. Toolbox citing core is the allowed reference direction (core is the required sibling).

## Solution summary

- The canonical closure rule lives in `skills/do-work/actions/work-reference.md`; capture and review only enforce/cite it.
- Producer compatibility is pinned in `_dev/tests/contract-regressions.sh` across core `review-work.md` and toolbox `code-review.md`.

## Worth knowing

- A universal consumer gate is incomplete until every real producer is enumerated and shown to emit compatible data; testing the first named caller only can leave a second package silently broken.
- Condition-driven inventory (exact marker fields matched to fenced templates) is safer than a filename list because it makes newly introduced producers fail closed.

## Back-reference

See `do-work/archive/UR-038/REQ-170-finding-closure-ratchet-and-rubric-home.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `cab5ba5`.
