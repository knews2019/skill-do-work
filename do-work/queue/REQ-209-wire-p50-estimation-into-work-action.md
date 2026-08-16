---
id: REQ-209
title: Wire P50 estimation into the work action
status: pending
created_at: 2026-08-16T23:52:07Z
user_request: UR-047
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-208]
maintenance: false
related: [REQ-208, REQ-210]
batch: p50-estimation
write_set: [skills/do-work/actions/work.md, skills/do-work/docs/work-guide.md]
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
*Source: UR-047 — "Add P50 active-duration estimation to do-work REQs" (lifecycle rules 3–6 and Presentation sections)*
