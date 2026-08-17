---
id: REQ-210
title: Recalculate P50 estimates in verify-requests after material repairs
status: completed
created_at: 2026-08-16T23:52:07Z
claimed_at: 2026-08-17T00:24:18Z
route: A
completed_at: 2026-08-17T00:26:46Z
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
- [x] **[PLAN]:** prime-action-files + crew loaded. Two anchored edits to verify-requests.md (numbered item 4 in Step 7's fix flow: material-change test, refresh block + fresh calculated_at + basis, no-prior-estimate derivation, untouched-REQs-byte-untouched, pending-only freeze guard, never-block; plus the revalidation-mode read-only bullet gains "recalculate any estimate: block"), then the contract pin — reference edit strictly before the pin so the suite never sees a pin without its reference.
- [x] **[APPLY]:** Code written exactly as planned; two declared files touched.
- [x] **[UNIFY]:** `git diff --stat` reviewed — verify-requests.md (two anchored edits, clean), contract-regressions.sh (one pin line). No debug artifacts. Full verify FAIL set identical to the 41-failure environment baseline (mod version string in one pre-existing message).

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

## Triage
**Route: A** - Simple
**Reasoning:** Names its file, well-specified: one numbered item in an existing fix flow plus a read-only-mode bullet and a contract pin. No exploration needed — Step 7's structure was read during capture-session analysis.
**Planning:** Not required

## Plan
**Planning not required** - Route A: direct to builder
*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/verify-requests.md` (modified) — Step 7 gains numbered item 4: recalculate the `estimate:` block after a material repair (material-change test; fresh `calculated_at` + updated `basis`; derive-if-absent per lifecycle rule 5; unmodified REQs stay byte-untouched; only `pending`/`pending-answers` queue REQs ever recalculated — claimed/archived estimates frozen; refreshed figure mentioned with the re-score; never blocks). Decision-revalidation's read-only bullet now names estimates explicitly.
- `_dev/tests/contract-regressions.sh` (modified) — `tools/estimate-p50.sh|actions/verify-requests.md` pin added (missing-script + missing-reference lock-in)

**What was done:** Wired the secondary estimation point: verify-requests refreshes estimates exactly when it has materially changed a REQ, and the wiring is contract-locked.

## Decisions
<!-- D-XX counter: none used in Open Questions. -->
- **D-01** (DECIDE & STATE): Reference edit committed-before-pin ordering inside this REQ's own diff — prevents any window where the aggregate suite sees a pin without its reference (the sibling REQ-209 review runs the suite from the working tree concurrently).

## Qualification

Passed — `tools/checks/qualify.sh` exit 0 (run below); Route A: no Scope section, drift check skipped by contract. Judgment: both edits substantive and traced to the REQ's requirements; the pin fired live for REQ-209 and guards this reference identically.

## Testing

**Tests run:** full `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ FAIL set identical to the recorded 41-failure environment baseline (only delta: version string inside one pre-existing installer message). The new `tools/estimate-p50.sh|actions/verify-requests.md` pin passes — script exists + executable, verify-requests.md references it.

**Red-green validation:** The REQ's captured RED (repair leaves a stale estimate) is prose behavior: verified by inspection that Step 7 item 4 fires on material change, skips untouched REQs, and revalidation mode is explicitly fenced. The wiring itself has mechanical red-green: the pin would FAIL the aggregate suite if the reference were removed (demonstrated by REQ-209's resolver incident).

*Verified by work action*

## Review

**Overall: 94%** | quick scan (Route A calibration, run inline by the orchestrator)

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor (report only): the material-change test ("not cosmetic rewording") is judgment prose — two agents could disagree at the margin; acceptable, the freeze guard bounds the blast radius to pending queue files. All six REQ requirements delivered: material-change recalc ✓, derive-if-absent ✓, untouched-stay-untouched ✓, re-score mention ✓, revalidation read-only ✓, never-block ✓.
**Acceptance:** Pass.

*Reviewed by review-work action (quick scan)*

## Lessons Learned

**What worked:** Ordering the reference edit before the pin inside one diff — a concurrency-safe pattern worth repeating whenever a contract pin and its referenced prose land together while another suite consumer may run mid-edit.
**What didn't:** Nothing — the pin pattern was proven by REQ-209's incident an hour earlier.
**Worth knowing:** Estimate freeze precedence: work.md Step 3.6 owns the freeze rule; verify's recalc item defers to it by scoping recalculation to `pending`/`pending-answers` queue files only.

## Orientation

**Now you can** repair a REQ in verify-requests and get its price refreshed in the same pass — estimates stay honest as scope changes, without ever touching claimed or archived work. Lives in the verify-requests capture-QA flow; leaf change, no map delta.

---
*Source: UR-047 — "Add P50 active-duration estimation to do-work REQs" (lifecycle rules 2 and 5)*
