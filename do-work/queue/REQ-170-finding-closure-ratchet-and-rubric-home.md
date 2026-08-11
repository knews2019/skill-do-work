---
id: REQ-170
title: Finding-closure ratchet and canonical earned-defense rubric
status: pending
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
write_set: [skills/do-work/crew-members/coding-guardrails.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/review-work.md, skills/do-work/actions/capture.md, skills/do-work-toolbox/actions/validate-feedback.md]
---

# Finding-Closure Ratchet and Canonical Earned-Defense Rubric

## What

Two small, single-home rules that make the review loop converge instead of plateau:

1. **Finding-closure ratchet** — a REQ that originates from a review or triage finding may only close with either a named regression test (fails before the fix, passes after) or a deletion of the surface the finding lived in — never a bare patch. Canonical text lives in `skills/do-work/actions/work-reference.md` beside the other pipeline contracts; `actions/review-work.md` enforces it at the gate (closure without test-or-deletion evidence bounces); `actions/capture.md` sharpens the GREEN criterion for finding-origin REQs to name the test or the deletion.
2. **Earned-defense rubric** — one canonical paragraph in `skills/do-work/crew-members/coding-guardrails.md`: defensive code must name the incident that earned it, and the fix must cost less surface than the risk it covers — user's wording preserved: *"what earned this, and is the fix still cheaper than the surface it added?"* The rubric already shipped inside `../do-work-toolbox/actions/validate-feedback.md` (REQ-169, commit 063bb88) — condense that to its triage-specific application (Surface-cost verdicts, Accept bar) plus a one-line citation of the canonical paragraph; review-work's gate cites it too. Toolbox citing core is the allowed reference direction (core is the required sibling).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
