---
id: REQ-169
title: validate-feedback flags remedies that add unearned defensive surface
status: completed
created_at: 2026-08-11T12:00:13Z
user_request: UR-037
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-168]
write_set: [skills/do-work-toolbox/actions/validate-feedback.md, _dev/tests/contract-regressions.sh]
claimed_at: 2026-08-11T12:57:58Z
route: B
completed_at: 2026-08-11T13:03:16Z
commit:
kb_status: pending
kb_entry:
---

# validate-feedback Flags Remedies That Add Unearned Defensive Surface

## What

Extend `skills/do-work-toolbox/actions/validate-feedback.md` so the triage applies the surface-cost rubric — user's words, verbatim: *"For each incident check what earned this, and is the fix still cheaper than the surface it added?"* — to every finding whose remedy would **add** defensive surface (a guard, fallback, retry, validation layer, rule, or warning apparatus). A remedy that can't name the incident earning it, or whose added surface costs more than the risk it covers, should not sail through as a plain **Accept** — it gets flagged (Push back or Discuss, with the rubric as the stated reasoning).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route B planning skip recorded; exploration identified Step 4, Step 5, per-finding output, Rules, checklist, and aggregate assertions as the complete seam.
- [x] **[APPLY]:** Added the scoped surface-cost classification, verdict constraint, visible report result, and five focused aggregate assertions exactly within the declared files.
- [x] **[UNIFY]:** Reviewed the complete diff/stat for both scoped files; verified the exact user rubric, surface-kind list, non-applicability guard, verdict mapping, per-finding output, rule/checklist consistency, RED-to-GREEN assertions, shell-fence lint, Bash syntax, full contracts, and absence of debug/whitespace artifacts.

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

## Triage

**Route: B** - Moderate

**Reasoning:** The edit is localized, but the rubric must affect verification, verdict selection, report visibility, and regression coverage without broadening into skepticism of fixes/deletions/simplifications.

**Planning:** Not required

## Plan

**Planning not required** - Route B: the request names the action, anchor steps, output block, verdict constraint, scope guard, and exact RED/GREEN behavior.

*Skipped by work action*

## Exploration

- Step 4 already distinguishes claim provenance and verifies premise/code/history, so the new check belongs after those facts are known and before verdict selection.
- Step 5's exact verdict vocabulary is the enforcement seam: an added defensive layer can remain Accept only when it names the incident, explains why its surface is cheaper than the risk, and states replay/test evidence; otherwise Push back or Discuss fits the existing meanings.
- The per-finding Output Format needs a visible Surface-cost line. `N/A` preserves existing behavior for direct bug fixes, deletions, and simplifications; `Earned`/`Flagged` carries the new reasoning only for surface-adding remedies.
- A focused aggregate contract can pin the exact rubric question, the surface-kind boundary, the non-applicability guard, the no-plain-Accept rule, and the report line without adding a new test file for one prompt-only change.

*Generated by Exploration phase*

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/validate-feedback.md` (modified) — apply and expose the prospective surface-cost rubric
- `_dev/tests/contract-regressions.sh` (modified) — ratchet rubric wording, scope boundary, verdict effect, and output visibility

**Files I will NOT touch:** retrospective audit/removals from REQ-168, other findings actions, capture/execute routing, verdict vocabulary, or runtime code

**Acceptance criteria (restated from REQ):**
- [x] Every remedy that adds a guard, fallback, retry, validation layer, rule, or warning apparatus is asked what incident earned it and whether the fix is cheaper than its added surface.
- [x] An unearned or net-costly surface-adding remedy cannot receive a plain Accept; it becomes Push back or Discuss with rubric-based reasoning.
- [x] Direct bug fixes, deletions, and simplifications triage exactly as before and are visibly marked outside the rubric.
- [x] The per-finding report exposes the Surface-cost result so the reader can see why the verdict moved.
- [x] Aggregate contracts pass and pin the scope/verdict/output behavior.

## Decisions

- **D-01 — Classify the remedy, not the finding.** A valid bug report can still propose an overbuilt fix; the surface-cost check runs after premise verification and evaluates only the proposed long-lived defense.
- **D-02 — Keep verdict vocabulary unchanged.** `Push back` expresses speculative or dominated defense; `Discuss` expresses a real incident with unresolved cost; `Accept` remains valid when the incident, cost, and test earn it.
- **D-03 — Render N/A for non-surface remedies.** A visible per-finding result proves the rubric was considered while preventing it from penalizing deletions, simplifications, and direct fixes.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/validate-feedback.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Added a Step 4 remedy classifier using the user's surface-cost question and the explicit guard/fallback/retry/validation/rule/warning boundary. Step 5 now bars plain Accept for unearned or net-costly added defense, maps speculative cases to Push back and unresolved real trade-offs to Discuss, and leaves direct fixes/deletions/simplifications at N/A. The per-finding block, Rules, and checklist expose the result; five aggregate assertions pin wording, scope, verdict effect, and visibility.

## Qualification

Passed — mechanical qualification and scope drift verified both declared files and exact Implementation Summary agreement. Manual qualification traced the complete decision flow from verified finding → remedy classification → cost evidence → constrained verdict → visible report/checklist result, and confirmed the explicit N/A branch prevents behavioral spillover to fixes, deletions, and simplifications.

## Testing

**Tests run:**
- RED: `bash _dev/tests/contract-regressions.sh` after adding five assertions and before changing the action
- focused five-pattern contract check against `skills/do-work-toolbox/actions/validate-feedback.md`
- `_dev/tests/action-shell-blocks.sh`
- `bash -n _dev/tests/contract-regressions.sh`
- `skills/do-work/tools/checks/qualify.sh do-work/working/REQ-169-validate-feedback-surface-cost-rubric.md`
- `skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-169-validate-feedback-surface-cost-rubric.md`
- `bash _dev/tests/contract-regressions.sh`
- `git diff --check`

**Result:** ✓ All GREEN; aggregate contracts exited 0 and all existing defensive-audit, shell, hook, install, update, manifest, and package-reference probes remained green.

**Red-green validation:**
- RED: all five new assertions failed against the prior action—rubric wording, surface boundary, N/A scope guard, verdict constraint, and output line were absent.
- GREEN: the same assertions pass after the action edit, and the complete aggregate suite exits 0.

**New tests added:** None (the prompt-only behavior fits the existing aggregate action-contract seam).

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh` — five assertions cover exact rubric, scope, verdict effect, and per-finding visibility

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-11T13:02:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — surface-adding remedies receive the exact rubric and cannot pass as plain Accept without incident/cost/test evidence; other remedy classes remain unchanged and every report shows the decision.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Putting the rubric after premise verification separates “is the finding true?” from “is its proposed fix worth owning?”, which prevents a valid complaint from laundering an overbuilt remedy.
- Reusing existing verdicts kept the change small: the new evidence changes classification without creating a parallel status vocabulary.

**What didn't:**
- A rubric phrased only as generic “consider complexity” would be untestable and easy to skip; the explicit surface kinds and visible result are what make it operational.

**Worth knowing:**
- `N/A` is a functional scope guard, not decorative output. It makes direct repair/deletion/simplification immune to the prospective skepticism pass while proving every finding was considered.

## Orientation

[MAP CHANGED] `do-work-toolbox validate-feedback` now prices any proposed added defense at Step 4, constrains its Step 5 verdict, and exposes `Surface-cost: N/A / Earned / Flagged` in every finding block.
