---
id: REQ-209
title: Wire P50 estimation into the work action
status: claimed
created_at: 2026-08-16T23:52:07Z
claimed_at: 2026-08-17T00:18:23Z
route: B
user_request: UR-047
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-208]
maintenance: false
related: [REQ-208, REQ-210]
batch: p50-estimation
write_set: [skills/do-work/actions/work.md, skills/do-work/docs/work-guide.md, _dev/tests/contract-regressions.sh]
estimate:
  p50_active_minutes: 60
  confidence: medium
  calculated_at: 2026-08-16T23:52:07Z
  basis:
    - careful edits inside the large work.md orchestrator
    - new presentation formats incl. dependency-graph critical path
    - freeze + legacy-backfill + trivial short-circuit rules
    - docs update
---

# Wire P50 Estimation into the Work Action

## What

Make `actions/work.md` the primary estimation wire point: immediately after Step 3 (Triage) — when the route is real — ensure an `estimate:` block exists, print it before planning starts, freeze it once execution begins, and render dependency-aware totals for multi-REQ runs.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** prime-action-files + crew loaded. Three surgical edits to work.md: (a) new Step 3.6 between 3.5 and 3.7 — reuse-if-present/frozen, trivial short-circuit via `--trivial`, otherwise lazy-load estimate-reference.md → extract → run estimator → persist block + stamp, print the three-line format, hard never-block rule; (b) Step 1 estimate-summary paragraph after the queue status summary for any selected set >1, totals via graph mode, unestimated members excluded-and-counted; (c) checklist line. Plus work-guide.md (pipeline list + accumulates list + one paragraph on P50 meaning) and the contract pin.
- [x] **[APPLY]:** Code written exactly as planned; three declared files touched (work.md ×3 edits, work-guide.md ×2, contract-regressions.sh pin + resolver entry).
- [x] **[UNIFY]:** `git diff --stat` reviewed — work.md (Step 3.6 + estimate summary + checklist line, clean), work-guide.md (renumbered pipeline list + estimate-line section, clean), contract-regressions.sh (pin + `resolve_runtime_file` core-root entry — the bare `tools/` path fell through to repo root on first run, caught by the new pin's own FAIL and fixed). No debug artifacts. Full verify FAIL set identical to baseline (41, mod the version string inside one pre-existing message).

## Why (if provided)

Route is assigned at work Step 3, not at capture — so post-triage is the first moment the estimator's strongest signal exists. This single wire point also satisfies lifecycle rule 5 for free: any REQ without an estimate (legacy or newly captured) gets one derived and persisted the first time it is selected for work.

## Context

Resolved during the capture session: the spec's lifecycle rule 1 (capture-time estimation) is amended out of v1 — the interactive capture window stays fast, and "derive when next verified or selected" becomes the universal rule. The spec's "ensure an estimate exists immediately before claiming" is satisfied in spirit at the first post-claim moment the route exists; the estimate is a forecast, and printing it before planning is the user-visible contract.

## Detailed Requirements

- **Ensure-estimate step** after Step 3 (Triage), before Step 4 (Planning): if the REQ has no `estimate:` block, load the estimate reference (REQ-208), extract signals, run the estimator script, and persist the block to the REQ frontmatter (with `calculated_at` per the Timestamp rule). If a valid block already exists, use it as-is.
- **Print the estimate before planning starts**, per the spec's format:

  ```
  Starting REQ-1459 — Add SD-39 review links and QA gates
  Estimated active duration: approximately 125 minutes (P50, medium confidence)
  Dominant factors: Route C, browser evidence matrix, performance gate, storage auditing
  ```

- **Freeze the estimate once execution begins.** Do not rewrite it using knowledge gained during implementation. No later step in the pipeline may touch the block.
- **Trivial short-circuit:** `effort_estimate: trivial` (or Route A indicators per the reference's rubric) ⇒ emit the floor estimate without loading the reference file, keeping overhead near zero exactly where the estimate is worth the least.
- **Multi-REQ presentation** (any run whose selected set contains more than one REQ — targeted `UR-NNN` runs, multi-token runs, and default full-queue runs alike; per the spec: "a UR or queue containing multiple REQs"): show the P50 estimate per REQ, total estimated agent effort, and critical-path active time computed from the `depends_on` graph rather than summing parallel branches — both values clearly labeled, per the spec's example format. Compute the totals with the estimator tooling's critical-path mode (REQ-208), not ad-hoc arithmetic.
- **Never block:** estimation failure (script missing, unparseable signals) logs a note and proceeds without an estimate — it must never stop the claim, require user clarification, or fail the run.
- **Docs:** update `docs/work-guide.md` to explain the estimate line and that P50 means roughly a 50% chance of completing within the estimated active minutes — an informational forecast, not a deadline or execution budget; user wait and suspended time are excluded.

## Constraints

- No P80 or calendar-time promises. (Batch constraint.)
- Estimation must never block execution or require user clarification. (Batch constraint.)
- Board display is out of v1. (Batch constraint.)
- Action-file text must stay agent-floor compatible and platform-agnostic (`_dev/primes/prime-action-files.md`).

## Dependencies

`depends_on: [REQ-208]` — needs the estimator script, the reference file, and the schema block to exist.

## Builder Guidance

**Firm:** wire point (post-triage, pre-planning), freeze rule, non-blocking rule, both presentation formats, critical path from the dependency graph.
**Builder latitude:** exact step numbering (e.g., a "Step 3.6"), wording of the printed lines, where the trivial short-circuit check sits relative to the reference-file load.

## Red-Green Proof
**RED prompt/case:** Run `do-work run` against a pending REQ that has no `estimate:` block. Today: work claims and triages the REQ and proceeds straight to planning — no estimate is derived, persisted, or printed.
**Why RED now:** `actions/work.md` has no estimation step; the schema has no `estimate:` block.
**GREEN when:** The same run persists an `estimate:` block into the REQ and prints the "Estimated active duration: approximately N minutes (P50, X confidence)" line with dominant factors before the planning step; a targeted `UR-NNN` run prints per-REQ estimates plus labeled "Total estimated effort" and "Estimated critical path" values; a REQ that already carries an estimate is not recalculated.
**Validation:** User confirmed — wire-point ranking (work primary, verify secondary, capture dropped) was explicitly confirmed in the capture session.

## Full Context
See `do-work/user-requests/UR-047/input.md` for complete verbatim input.

---

## Triage
**Route: B** - Medium
**Reasoning:** Clear outcome — insert one estimation step, one multi-REQ summary paragraph, checklist line, and a docs section. The "where" is fully known (work.md's step structure is established); care is in wording, not architecture.
**Planning:** Not required

## Plan
**Planning not required** - Route B: Exploration-guided implementation
*Skipped by work action*

## Exploration
- Insertion point: between Step 3.5 (Open Questions) and Step 3.7 (Spec Loading) — numbering "Step 3.6" is free; estimate must print before Step 4 (Planning).
- Multi-REQ presentation home: Step 1, immediately after the queue status summary (it already renders per-run aggregate state); targeted-mode UR expansion announcement already lists resolved REQs in execution order.
- Freeze rule home: inside Step 3.6 itself (reuse-if-present is the same sentence), consistent with how Step 6.25 documents replace-not-append.
- Orchestrator Checklist block needs a Step 3.6 line (it enumerates every step).
- Wiring lock-in: contract-regressions.sh `hardened_check_scripts` pin table takes `"tools/estimate-p50.sh|actions/work.md"` — same missing-script + missing-reference assertions the checks scripts get.
- Docs: work-guide.md "Pipeline steps" numbered list and "What accumulates in the REQ file" section are the two places the estimate line belongs.
*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modified) — Step 3.6 ensure-estimate + print; Step 1 multi-REQ estimate summary; Orchestrator Checklist line
- `skills/do-work/docs/work-guide.md` (modified) — estimate line in pipeline steps + P50 meaning
- `_dev/tests/contract-regressions.sh` (modified) — estimator pin in `hardened_check_scripts`

**Files I will NOT touch:** `skills/do-work/actions/verify-requests.md` (REQ-210), `skills/do-work/tools/estimate-p50.sh` and `skills/do-work/actions/estimate-reference.md` (REQ-208, done), the board tool

**Acceptance criteria (restated from REQ):**
- [x] Ensure-estimate step after triage, before planning; reuse-if-present; frozen once execution begins
- [x] Printed format matches the spec example (approximately N minutes, P50, confidence, dominant factors)
- [x] Trivial short-circuit without reference-file load
- [x] Multi-REQ table + labeled total and critical path via the estimator's graph mode, for any selected set >1
- [x] Estimation never blocks and never asks the user
- [x] work-guide.md explains the estimate line and P50 meaning

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified) — new **Step 3.6: Estimate Active Duration** (reuse-frozen / trivial short-circuit / lazy-load reference + extract + run estimator + persist + print; hard never-block rule; legacy REQs estimated at first selection); **Estimate summary** paragraph in Step 1 after the queue status summary (per-REQ minutes + depends_on, labeled total effort and critical path via the estimator's graph mode, unestimated members excluded-and-counted, never gates selection); Orchestrator Checklist gains the Step 3.6 line
- `skills/do-work/docs/work-guide.md` (modified) — pipeline list renumbered with the Estimate step; "The estimate line" section explaining the printed format, P50 = ~50% chance within active minutes, active-time exclusions, forecast-not-deadline, and the two labeled totals; `estimate:` block added to the what-accumulates list
- `_dev/tests/contract-regressions.sh` (modified) — `tools/estimate-p50.sh|actions/work.md` added to `hardened_check_scripts` (missing-script + missing-reference lock-in); `tools/estimate-p50.sh` added to `resolve_runtime_file`'s core-root pattern so the pin resolves against `skills/do-work/`

**What was done:** Wired the primary estimation point into the work pipeline: estimates are ensured and printed post-triage/pre-planning where the route signal is real, frozen once execution begins, and multi-REQ runs render dependency-aware totals. The wiring itself is now contract-locked to the shipped script.

## Decisions
<!-- D-XX counter: none used in Open Questions. -->
- **D-01** (DECIDE & STATE): Step number 3.6 (between Open Questions and Spec Loading) rather than 3.9/4-renumber — printing lands before planning per the spec while leaving every existing step number stable for cross-references.
- **D-02** (DECIDE & STATE): The multi-REQ summary excludes unestimated members from both aggregates and says so, rather than force-estimating the whole set at scan time — keeps Step 1 cheap and honest; Step 3.6 fills the gaps as each REQ is claimed.
- **D-03** (DECIDE & STATE): Added the estimator to `resolve_runtime_file`'s core-root case — first verify run FAILed the new pin because bare `tools/` paths fall through to repo root (the resolver is a hand-maintained enumeration; Closed Enumerations Go Stale in action).

## Qualification

Passed — `tools/checks/qualify.sh` exit 0; `tools/checks/scope-drift.sh` exit 0 (Implementation Summary matches Scope). Judgment: all three files substantive, all six acceptance criteria traced, wiring is live prose + a live contract pin (proven by the pin's own FAIL firing before the resolver fix).

## Testing

**Tests run:** full `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh` (twice — the first run caught the resolver gap via the new pin's FAIL)
**Result:** ✓ FAIL set identical to the recorded 41-failure environment baseline (only text delta: the suite-manifest version string inside one pre-existing message, 0.193.6 → 0.194.0). The new `tools/estimate-p50.sh|actions/work.md` pin passes: script exists + executable, work.md references it. ShellCheck warning-severity clean.

**Red-green validation:**
- New contract pin: ✗ FAILed on first run (resolver sent the pin to repo root — a genuine missing-reference failure mode) → ✓ passes after the core-root entry. This is the wiring lock-in the REQ promised: deleting Step 3.6's pointer or the script now fails the aggregate suite.
- The REQ's captured RED (work run derives/prints no estimate) is prose behavior — verified by inspection: Step 3.6 sits before Step 4, prints the spec's exact format, and the checklist enumerates it.

*Verified by work action*

---
*Source: UR-047 — "Add P50 active-duration estimation to do-work REQs" (lifecycle rules 3–6 and Presentation sections)*
