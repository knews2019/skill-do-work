---
title: "Lessons from REQ-169: validate-feedback flags remedies that add unearned defensive surface"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-169-validate-feedback-flags-remedies-that-ad.md]
related:
  - page: concept-contract-verification-gates
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-169: validate-feedback flags remedies that add unearned defensive surface

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

Extend `skills/do-work-toolbox/actions/validate-feedback.md` so the triage applies the surface-cost rubric — user's words, verbatim: *"For each incident check what earned this, and is the fix still cheaper than the surface it added?"* — to every finding whose remedy would **add** defensive surface (a guard, fallback, retry, validation layer, rule, or warning apparatus). A remedy that can't name the incident earning it, or whose added surface costs more than the risk it covers, should not sail through as a plain **Accept** — it gets flagged (Push back or Discuss, with the rubric as the stated reasoning).

## Solution summary

Added a Step 4 remedy classifier using the user's surface-cost question and the explicit guard/fallback/retry/validation/rule/warning boundary. Step 5 now bars plain Accept for unearned or net-costly added defense, maps speculative cases to Push back and unresolved real trade-offs to Discuss, and leaves direct fixes/deletions/simplifications at N/A. The per-finding block, Rules, and checklist expose the result; five aggregate assertions pin wording, scope, verdict effect, and visibility.

## What worked

- Putting the rubric after premise verification separates “is the finding true?” from “is its proposed fix worth owning?”, which prevents a valid complaint from laundering an overbuilt remedy.
- Reusing existing verdicts kept the change small: the new evidence changes classification without creating a parallel status vocabulary.

## What didn't work

- A rubric phrased only as generic “consider complexity” would be untestable and easy to skip; the explicit surface kinds and visible result are what make it operational.

## Worth knowing

- `N/A` is a functional scope guard, not decorative output. It makes direct repair/deletion/simplification immune to the prospective skepticism pass while proving every finding was considered.

## Back-reference

See `do-work/archive/UR-037/REQ-169-validate-feedback-surface-cost-rubric.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `063bb88`.
