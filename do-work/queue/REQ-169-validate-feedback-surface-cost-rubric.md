---
id: REQ-169
title: validate-feedback flags remedies that add unearned defensive surface
status: pending
created_at: 2026-08-11T12:00:13Z
user_request: UR-037
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-168]
write_set: [skills/do-work-toolbox/actions/validate-feedback.md]
---

# validate-feedback Flags Remedies That Add Unearned Defensive Surface

## What

Extend `skills/do-work-toolbox/actions/validate-feedback.md` so the triage applies the surface-cost rubric — user's words, verbatim: *"For each incident check what earned this, and is the fix still cheaper than the surface it added?"* — to every finding whose remedy would **add** defensive surface (a guard, fallback, retry, validation layer, rule, or warning apparatus). A remedy that can't name the incident earning it, or whose added surface costs more than the risk it covers, should not sail through as a plain **Accept** — it gets flagged (Push back or Discuss, with the rubric as the stated reasoning).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Companion to the UR-036 stabilization batch: REQ-168 removes unearned defensive surface already shipped; this REQ stops new unearned surface from entering through accepted review findings — the front door where "add more defense" remedies arrive. Untested/unearned defensive code is negative-value (the session-start hook incident is the exemplar: the defense was the bug).

## Context

- Anchor points in the current file: Step 4 (Verify Each Item Against the Code) is where the surface-cost question gets asked per finding; Step 5 (Recommend a Verdict per Item) is where it shifts the verdict; the Output Format's per-finding block is where the flag surfaces to the reader. The builder chooses the exact placement — the contract is that the rubric is applied and visible, not where the prose lands.
- Same rubric, two homes: REQ-168's audit applies it retrospectively to shipped code; this applies it prospectively to incoming findings. Keep the wording consistent between the two so the discipline reads as one rule.
- Scope guard: the rubric applies to remedies that *add* surface. Findings that fix bugs, delete code, or simplify are untouched by this check — do not let it become a generic skepticism pass that pushes back on everything.
- The action's existing honesty rules still govern: this is not a license to push back to reduce work; a defensive remedy that *does* name its incident and is cheaper than its surface is a legitimate Accept.
- Per the action-file conventions, an added Rationalizations row (if the builder adds one) must pass the "can I name the specific failure this prevents?" test — the session-start pipefail incident (46 lines of defense around 2 lines of logic; the defense was the defect) is the traceable origin.
- `_dev/tests/contract-regressions.sh` runs against action-file edits — keep the file within its existing contracts.

## Red-Green Proof

**RED prompt/case:** Paste a finding into `do-work-toolbox validate-feedback` whose remedy proposes speculative defense — e.g. "wrap the version parse in a retry and add a fallback config path, just in case." Today's Step 4/5 verify the *claim's premise* but never question the *remedy's added surface*; the rubric appears nowhere in the action, so the item can land as a plain Accept.
**Why RED now:** The triage adversarially verifies whether a finding is true, but has no check on whether its fix is worth the surface it adds — precisely how unearned defensive layers accreted (UR-036 diagnosis).
**GREEN when:** The same pasted finding produces a report where the surface-adding remedy is flagged with the rubric — verdict Push back or Discuss unless the finding names the incident earning the defense — and the per-finding block shows the rubric-based reasoning. Bug-fix/deletion/simplification findings triage exactly as before.
**Validation:** User confirmed (rubric wording supplied verbatim by the user; interpretation — flag surface-adding remedies in triage — inferred from the UR-036 discussion).

## Full Context

See `do-work/user-requests/UR-037/input.md` for complete verbatim input.

---
*Source: "also do-work validate-feedback should also flag the following: 'For each incident check what earned this, and is the fix still cheaper than the surface it added?'" (UR-037)*
