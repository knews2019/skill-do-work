---
id: REQ-210
title: Recalculate P50 estimates in verify-requests after material repairs
status: pending
created_at: 2026-08-16T23:52:07Z
user_request: UR-047
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: [REQ-208]
maintenance: false
related: [REQ-208, REQ-209]
batch: p50-estimation
write_set: [skills/do-work/actions/verify-requests.md]
estimate:
  p50_active_minutes: 25
  confidence: high
  calculated_at: 2026-08-16T23:52:07Z
  basis:
    - single focused edit in one action file
    - hooks into existing Step 7 fix flow
    - no new formats or tooling
---

# Recalculate P50 Estimates in verify-requests After Material Repairs

## What

Wire the secondary estimation point: when capture-QA Step 7 applies fixes that materially change a REQ (added requirements, resolved Ambiguous gaps that became new requirements, tightened Red-Green Proof), recalculate its `estimate:` block using the REQ-208 estimator and persist the fresh values.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Spec lifecycle rule 2: "verify-requests should recalculate it after repairing or materially changing a REQ." Verify has already re-read the whole REQ at that point, so the marginal cost is roughly one script call.

## Context

This is the secondary wire point — work (REQ-209) is primary. A REQ that never passes through verify still gets its estimate at first selection for work.

## Detailed Requirements

- In capture-QA **Step 7 (Offer Fixes)**, after fixes are applied to a REQ: if the fixes materially changed scope (new/changed requirements, constraints, acceptance criteria, or Red-Green Proof — not cosmetic rewording), recalculate the REQ's `estimate:` block via the estimator (load the estimate reference, extract signals, run the script) and persist it with a fresh `calculated_at` and updated `basis`.
- A repaired REQ with **no prior estimate** gets one derived and persisted (spec lifecycle rule 5).
- REQs whose files were **not** modified by Step 7 keep their existing estimates untouched.
- The re-score step (Step 7.3) mentions the refreshed estimate in its output when one was recalculated.
- **Decision-revalidation mode stays read-only** — no estimation there; it changes no REQ content by contract.
- Never block: an estimation failure logs a note and the verify flow completes normally.
- Claimed/archived REQs are never touched (their estimates are frozen — batch constraint; verify only edits queue files anyway).

## Constraints

- Estimation must never block or require user clarification. (Batch constraint.)
- No P80 or calendar promises. (Batch constraint.)
- The freeze rule wins over recalculation: only `pending`/`pending-answers` queue REQs are ever recalculated.

## Dependencies

`depends_on: [REQ-208]` — needs the estimator script and reference file.

## Builder Guidance

**Firm:** recalculation only after material change via Step 7 fixes; read-only revalidation mode untouched; freeze rule respected.
**Builder latitude:** the exact wording of the "materially changed" test and where in Step 7's numbered flow the recalculation sentence lands.

## Red-Green Proof
**RED prompt/case:** Run `do-work verify-requests UR-NNN` on a UR whose REQ has an Important gap plus an existing `estimate:` block; accept the offered fix. Today: the REQ gains the missing requirement but its estimate (and `calculated_at`) is left stale.
**Why RED now:** `actions/verify-requests.md` Step 7 has no estimation hook.
**GREEN when:** After the same repair, the REQ's `estimate:` block carries a fresh `calculated_at` and a basis reflecting the enlarged scope; a verify run that changes nothing leaves every estimate byte-identical; decision-revalidation mode still modifies no files.
**Validation:** User confirmed — verify-as-secondary-wire-point was explicitly confirmed in the capture session.

## Full Context
See `do-work/user-requests/UR-047/input.md` for complete verbatim input.

---
*Source: UR-047 — "Add P50 active-duration estimation to do-work REQs" (lifecycle rules 2 and 5)*
