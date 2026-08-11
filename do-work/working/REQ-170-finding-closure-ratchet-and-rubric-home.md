---
id: REQ-170
title: Finding-closure ratchet and canonical earned-defense rubric
status: claimed
created_at: 2026-08-11T13:55:58Z
user_request: UR-038
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-165, REQ-166, REQ-167, REQ-168, REQ-169]
batch: stabilization-audit
write_set: [skills/do-work/crew-members/coding-guardrails.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/review-work.md, skills/do-work/actions/capture.md, skills/do-work-toolbox/actions/validate-feedback.md, _dev/tests/contract-regressions.sh]
claimed_at: 2026-08-11T19:56:32Z
route: C
---

# Finding-Closure Ratchet and Canonical Earned-Defense Rubric

## What

Two small, single-home rules that make the review loop converge instead of plateau:

1. **Finding-closure ratchet** — a REQ that originates from a review or triage finding may only close with either a named regression test (fails before the fix, passes after) or a deletion of the surface the finding lived in — never a bare patch. Canonical text lives in `skills/do-work/actions/work-reference.md` beside the other pipeline contracts; `actions/review-work.md` enforces it at the gate (closure without test-or-deletion evidence bounces); `actions/capture.md` sharpens the GREEN criterion for finding-origin REQs to name the test or the deletion.
2. **Earned-defense rubric** — one canonical paragraph in `skills/do-work/crew-members/coding-guardrails.md`: defensive code must name the incident that earned it, and the fix must cost less surface than the risk it covers — user's wording preserved: *"what earned this, and is the fix still cheaper than the surface it added?"* The rubric already shipped inside `../do-work-toolbox/actions/validate-feedback.md` (REQ-169, commit 063bb88) — condense that to its triage-specific application (Surface-cost verdicts, Accept bar) plus a one-line citation of the canonical paragraph; review-work's gate cites it too. Toolbox citing core is the allowed reference direction (core is the required sibling).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the required core/toolbox instructions and existing REQ-169 contract assertions; selected one canonical home per rule and citation-sized caller hooks.
- [x] **[APPLY]:** Added the ratchet, rubric, capture/review enforcement, and provenance-preserving toolbox citation in exactly the five declared files.
- [x] **[UNIFY]:** Reviewed all five diffs, confirmed net +9 instruction lines, ran both contract suites and the focused canonical-home/caller check, audited changed paths/debug artifacts, and passed `git diff --check`.

## Why (if provided)

User goal (UR-036): stop reviews from returning 3–5 findings per pass. Per-item triage never sees the aggregate — forty individually-reasonable accepts each adding a small guard collectively drift the product toward complexity. The ratchet converts the entire findings stream into monotonic convergence (every accepted finding shrinks or hardens the product); the rubric's single home keeps the anti-complexity campaign from becoming its own complexity.

## Detailed Requirements

- The ratchet's canonical statement is written once, in work-reference; capture and review-work carry one-line pointers plus their local enforcement hook — not restatements.
- review-work's gate treats a finding-origin REQ without test-or-deletion evidence the way it treats missing TDD evidence: send back, don't waive.
- The rubric paragraph in coding-guardrails is ≤1 paragraph; it always loads during implementation, which is where a fix's actual shape is decided.
- Net-surface accounting: state in the diff/report roughly how many lines of instruction were added — the intent is a handful of lines total across four files.
- In scope for the single-home constraint: converting validate-feedback's shipped rubric restatement (REQ-169, commit 063bb88) into triage-specific application + citation, without weakening the Surface-cost verdict contract its checklist and `_dev/tests/contract-regressions.sh` additions now enforce.
- Explicitly out of scope: metric machinery (a findings-trend log fails the rubric's own test).
- Keep within existing file contracts — `_dev/tests/contract-regressions.sh` and `_dev/tests/shipped-package-reference-contract.sh` must stay green; shipped files must not cite `CLAUDE.md`/`AGENTS.md`.

## Builder Guidance

Certainty: Firm on both rules and the single-home constraint; exploratory on exact section placement within work-reference and review-work. Scope cue from the user's campaign: every added rule is itself surface — write the minimum that a capable model can apply, per the "state intent, not a directive rule" convention.

## Red-Green Proof

**RED prompt/case:** Today a finding-origin REQ (e.g. one captured from a validate-feedback Accept) can pass the review-work gate with a bare patch — no section of work-reference or review-work requires a regression test or a deletion at closure. And the earned-defense rubric's only shipped statement lives in toolbox's validate-feedback (REQ-169, commit 063bb88) — `coding-guardrails.md`, which always loads at implementation time where the fix's shape is decided, says nothing about it, so a builder never sees the rubric.
**Why RED now:** "Review findings become ratchet tests" lives in UR-036's batch-constraint prose — soft instruction nobody enforces; the rubric shipped into the triage gate only, not into the build path or the closure gate.
**GREEN when:** work-reference contains the ratchet's canonical section; review-work's gate/checklist bounces a finding-origin closure lacking test-or-deletion evidence; capture's GREEN criterion names the requirement for finding-origin REQs; coding-guardrails contains the one-paragraph rubric; all other mentions are one-line pointers; the contract test suites pass.
**Validation:** User confirmed (plan presented in full and captured on the user's explicit "capture").

## Full Context

See `do-work/user-requests/UR-038/input.md` for complete verbatim input.

---
*Source: "capture" — approving the stabilization plan v2 discussed in-session (UR-038)*

---

## Triage

**Route: C** - Complex

**Reasoning:** This changes a canonical cross-package rubric and adds a new closure gate spanning capture, implementation, review, and validation surfaces.

**Planning:** Required

## Plan

1. Add the one-paragraph earned-defense rubric to `skills/do-work/crew-members/coding-guardrails.md` § 2, preserving the user's exact question and keeping the five-guardrail taxonomy unchanged.
2. Add a compact canonical Finding-Closure Ratchet beside the Step 6.5 testing templates in `skills/do-work/actions/work-reference.md`; define finding provenance from existing durable markers/text and require either named fail-before/pass-after regression evidence or deletion of the finding surface.
3. Add capture's local enforcement hook in `skills/do-work/actions/capture.md`, and preserve accepted-feedback provenance in `skills/do-work-toolbox/actions/validate-feedback.md`, both by citation rather than restatement.
4. Enforce the canonical ratchet and earned-defense rubric in `skills/do-work/actions/review-work.md`: missing closure evidence is an Important finding and Acceptance Fail, independent of `tdd`, score, or unrelated green tests.
5. Condense validate-feedback's duplicated rubric into its triage-specific Surface-cost/Accept application while retaining its existing output and contract-test tokens.

**Architectural decisions:** Use the existing `review_generated` marker or explicit REQ/UR review/triage provenance instead of adding schema; keep the two normative rules in one home each; add no metrics or trend machinery.

**Requirement mapping:** The canonical section owns the ratchet; capture names its proof at intake; review hard-gates closure; the always-loaded guardrails own the defense rubric; validate-feedback retains only its local verdict behavior. All five declared files trace directly to those requirements.

**Testing approach:** Record the current missing-rule RED state, run focused canonical-home/restatement searches, then run `_dev/tests/contract-regressions.sh`, `_dev/tests/shipped-package-reference-contract.sh`, `git diff --check`, and a changed-path/net-line audit.

*Generated by Plan agent*

**Plan validation:** Every Detailed Requirement maps to a planned task and no task is orphaned. ⚠ Plan has 5 tasks — quality degrades past 3; keep each edit citation-sized and split only if implementation cannot remain within the declared five-file surface.

## Exploration

- The captured RED state is real: no current finding-origin closure contract identifies deletion as alternate proof or forces Acceptance Fail; capture does not preserve finding provenance; the always-loaded simplicity guardrail has no earned-defense rubric.
- `work-reference.md`'s Step 6.5 testing-template seam is the correct canonical home; the existing Testing shape already supports named fail-before/pass-after evidence without another schema field.
- Capture's hook must be independent of `tdd` and ordinary behavioral-proof inference, while review must set both an Important finding and Acceptance `Fail` so the existing verdict/remediation machinery actually bounces the REQ.
- `validate-feedback.md` must retain its Surface-cost boundary, N/A classification, Accept routing, output token, and existing REQ-169 regex wording while citing the new canonical core paragraph.
- Integration risk: `skills/do-work/actions/capture.md` already contains an unrelated uncommitted screenshot-publication change in the main tree; reconciliation must preserve it while landing only this REQ's Step 1 hook.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/crew-members/coding-guardrails.md` (modified) — canonical one-paragraph earned-defense rubric
- `skills/do-work/actions/work-reference.md` (modified) — canonical Finding-Closure Ratchet
- `skills/do-work/actions/review-work.md` (modified) — hard closure gate and rubric citation
- `skills/do-work/actions/capture.md` (modified) — finding-origin GREEN proof hook
- `skills/do-work-toolbox/actions/validate-feedback.md` (modified) — triage-specific Surface-cost application and provenance handoff
- `_dev/tests/contract-regressions.sh` (modified) — review-generated producer/closure-gate contract ratchet added during review remediation

**Files I will NOT touch:** test sources other than the narrow remediation ratchet below, `actions/work.md`, schemas/parsers/board code, metrics or trend logs, release metadata, or unrelated capture screenshot logic

**Review remediation scope:** The independent review exposed a producer/consumer gap in `review-work.md`; `_dev/tests/contract-regressions.sh` is added narrowly to pin the review-generated follow-up template's named GREEN/deletion proof contract. Other test sources remain out of scope.

**Acceptance criteria (restated from REQ):**
- [ ] The finding-closure ratchet is normative only in work-reference and requires a named fail-before/pass-after regression test or deletion of the named finding surface.
- [ ] Review bounces missing finding-origin evidence with an Important finding and Acceptance Fail, regardless of `tdd`, score, or unrelated passing tests.
- [ ] Capture makes a finding-origin GREEN criterion name the intended regression test or deletion.
- [ ] Coding guardrails contain one canonical earned-defense paragraph with the user's exact question.
- [ ] Review and validate-feedback cite the canonical rubric; validate-feedback retains its Surface-cost and Accept contracts without a competing restatement.
- [ ] Net instruction growth stays to a handful of lines; contract and shipped-reference suites pass.

## Pre-Flight

**Git:** ⚠ Pre-existing changes outside `do-work/` include the deferred release/capture work plus REQ-173 helper/test changes. Preserve them, exclude them from this REQ's branch, and reconcile only the overlapping `capture.md` hunk at integration.
**Tests baseline:** ⚠ `_dev/tests/contract-regressions.sh` already fails before REQ-170 because the uncommitted REQ-173 prime link points into repo-only `do-work/archive/`, which is absent from installed topology. All later probes in the baseline output pass.
**Dependencies:** ✓ No missing repository dependency directory detected.

*Checked by work action*

## Decisions

- **D-01 — Canonicalize closure beside Step 6.5.** The testing-template seam gives capture and review one compact definition without adding schema.
- **D-02 — Keep earned defense inside Simplicity First.** The exact user question fits in one paragraph without expanding the guardrail taxonomy.
- **D-03 — Preserve REQ-169's triage vocabulary.** Toolbox keeps its Surface-cost boundary, N/A carve-out, Accept routing, and output token while citing core.
- **D-04 — Carry accepted-finding provenance into capture.** Verbatim claim, severity/source, Evidence, and Surface-cost remain together so the closure proof can be named later.

## Implementation Summary

- `skills/do-work/actions/work-reference.md` (modified) — added the canonical Finding-Closure Ratchet.
- `skills/do-work/crew-members/coding-guardrails.md` (modified) — added the one-paragraph earned-defense rubric.
- `skills/do-work/actions/capture.md` (modified) — added the finding-origin GREEN proof hook.
- `skills/do-work/actions/review-work.md` (modified) — added closure-gate and earned-defense citations.
- `skills/do-work-toolbox/actions/validate-feedback.md` (modified) — preserved triage behavior while citing the canonical rubric and carrying provenance.
- Builder commit: `84c1b81df6f7c34920e5285b63b121355cdfb03f`; integrated by merge commit `035613e` over exact range `bd5ecf6..035613e`.
- No integration seams or discovered follow-up tasks remain.

## Qualification

- Exact diff range: `bd5ecf6..035613e`.
- Changed paths match the five-file Route C scope exactly.
- Named RED/GREEN caller contract: 9 missing contracts before implementation, all 14 checks passing after implementation.
- `_dev/tests/contract-regressions.sh`: pass on the builder branch.
- `_dev/tests/shipped-package-reference-contract.sh`: pass on the builder branch.
- `git diff --check`: pass; no debug or temporary artifacts found.
