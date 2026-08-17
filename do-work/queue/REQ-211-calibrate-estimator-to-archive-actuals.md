---
id: REQ-211
title: Calibrate estimator scoring table to archive actuals
status: pending
created_at: 2026-08-17T08:05:43Z
user_request: UR-048
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-212, REQ-213]
batch: estimator-calibration
write_set: [skills/do-work/tools/estimate-p50.sh, skills/do-work/actions/estimate-reference.md, _dev/tests/p50-estimator-determinism.sh]
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-17T08:05:43Z
  basis:
    - Route B
    - 4-file write set
    - (priced with the pre-calibration table this REQ replaces)
---

# Calibrate Estimator Scoring Table to Archive Actuals

## What

Re-fit `tools/estimate-p50.sh`'s scoring table to the measured archive distribution: bases become the per-route medians (A=5, B=10, C=20), the floor drops to 5, signal weights divide by ≈2.5, and the calibration's provenance is recorded in `estimate-reference.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- New table: bases A=5, B=10, C=20; floor 5 (`--trivial` prints 5); weights: write-set +1/file, new-files +2, subsystems +3 per extra, acceptance +1, deps-depth +2, browser +8, persistence +6, async +6, performance +4, regression +4, full-suite +4.
- Confidence rubric thresholds re-anchored to the new scale: high for trivial or Route A raw ≤ 10; low for Route C with write-set ≥ 15, subsystems ≥ 3, or raw ≥ 75 (≈ 2× the C p80); medium otherwise.
- Determinism, nearest-5 rounding, basis echo, flag surface, critical-path mode, and the no-P80 fence are unchanged.
- Re-pin every lock-in test to the new values; the suite's floor assertions move 10 → 5.
- `estimate-reference.md`: update the block-template comment (floor), the trivial section, the confidence rubric, and replace the pure-prior paragraph in Calibration Honesty with the provenance record: 188 samples kept of 190 (>4h/negative outlier rule), medians/p80 by route, date, and the wall≈active caveat with its conservative bias direction.
- Frozen estimates on already-archived REQs are not touched.

## Constraints

- Batch constraint: calibration provenance must be recorded where the next re-fit will look (estimate-reference.md), not only in the commit message.

## Builder Guidance

Firm on the table values and provenance recording; latitude on exact rubric-comment wording.

## Red-Green Proof
**RED prompt/case:** `tools/estimate-p50.sh --route C --write-set 12 --browser --persistence --full-suite` prints 125 (the uncalibrated prior); `--trivial` prints 10.
**Why RED now:** The table's scale was inherited from the spec's worked example, measured ≈2.5–5× above archive actuals (Route medians 4.7/9.2/21.4 vs bases 10/25/45).
**GREEN when:** The same invocations print 50 and 5; a bare `--route B` prints 10 (the measured B median); the re-pinned determinism suite passes; `maintainer-verify` FAIL set stays at the environment baseline.
**Validation:** User confirmed — "Yes, apply it," after reviewing the measured table.

## Full Context
See `do-work/user-requests/UR-048/input.md` for complete verbatim input.

---
*Source: UR-048 — calibration application*
