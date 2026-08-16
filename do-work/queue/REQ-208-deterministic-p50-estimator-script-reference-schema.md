---
id: REQ-208
title: Deterministic P50 estimator script, reference file, and schema
status: pending
created_at: 2026-08-16T23:52:07Z
user_request: UR-047
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-209, REQ-210]
batch: p50-estimation
write_set: [skills/do-work/tools/estimate-p50.sh, skills/do-work/actions/estimate-reference.md, skills/do-work/actions/work-reference.md, _dev/tests/p50-estimator-determinism.sh]
estimate:
  p50_active_minutes: 85
  confidence: medium
  calculated_at: 2026-08-17T00:04:00Z
  basis:
    - new shipped shell script with lock-in tests
    - critical-path graph mode added by verify-requests repair
    - new lazy-loaded reference action-companion file
    - schema amendment in work-reference.md incl. effort_estimate bridge
    - full-suite maintainer-verify qualification
---

# Deterministic P50 Estimator: Script, Reference File, and Schema

## What

Build the estimation foundation: a shipped deterministic script that maps extracted REQ signals to `p50_active_minutes`, a lazy-loaded reference file telling the estimating agent how to extract those signals, and the `estimate:` frontmatter block documented in the schema — with the `effort_estimate` comment amended in the same commit.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

An agent following prose is not deterministic and cannot be regression-tested; the spec requires "the same result for the same normalized REQ inputs" and automated determinism tests. Programs beat prose for anything mechanical (CLAUDE.md). The split: agent extracts judgment signals, script does the arithmetic.

## Context

Resolved during the capture session:

- The estimator ships in the **do-work skill's own tools** (alongside `tools/checks/`), never in queue-kanban — the board tool's three-write-surface rule stays intact.
- The reference file follows the existing companion-file pattern (`capture-reference.md`, `work-reference.md`): loaded lazily, only at estimation moments, by the actions wired in REQ-209/REQ-210.
- The board's yaml.v3 permissive-map frontmatter parser was verified to tolerate the nested `estimate:` block (unknown keys ignored on the happy path); no board change is needed for backwards compatibility.
- v1 is a **pure-prior** estimator: no historical actives data exists anywhere in the repo (only `claimed_at`/`completed_at` wall time, which the spec forbids using). Document this honestly in the reference file.

## Detailed Requirements

**The script** (POSIX shell, per `_dev/primes/prime-shell-commands.md`):
- Takes normalized signals as flags/arguments — the signal set is drawn from the spec's estimation inputs: route (A/B/C), write-set size and file types, number of runtime subsystems involved, new files/assets/dependencies, acceptance-criteria count, dependency depth and serialization, browser/responsive/visual/accessibility/screenshot requirements, persistence/migration/API/schema changes, asynchronous lifecycle/teardown/race/retry behavior, performance instrumentation, focused vs full-suite verification, lint/deploy/asset-integrity/cross-route regression requirements, independent review and expected remediation cost.
- Prints `p50_active_minutes`, `confidence` (low | medium | high), and echoes the dominant sizing factors for the `basis` list.
- Rounds estimates to the nearest five minutes.
- Enforces a reasonable minimum (the floor — also the `effort_estimate: trivial` short-circuit value).
- Deterministic: identical flags → identical output, always.
- Do not add P80 or other percentile outputs.
- **Critical-path mode** (added by verify-requests repair — spec acceptance requires automated dependency-graph coverage): a mode that takes per-REQ minutes plus `depends_on` edges (e.g. `REQ-209:60:REQ-208` triples) and prints total estimated effort and critical-path active minutes computed over the dependency graph (longest path, never a sum of parallel branches), both labeled. Deterministic and covered by the lock-in tests; REQ-209's multi-REQ presentation consumes this mode.

**The reference file** (`skills/do-work/actions/estimate-reference.md` or equivalent action companion):
- Signal-extraction guidance: how the estimating agent reads a REQ (route from triage, `write_set`, acceptance criteria, Red-Green Proof, Constraints) into the script's flags.
- The exact `estimate:` frontmatter block template (`p50_active_minutes`, `confidence`, `calculated_at` per the Timestamp rule, `basis` list).
- Confidence rubric (low/medium/high) and the trivial short-circuit rule.
- Documents that P50 means roughly a 50% chance of completing within the estimated active minutes, and that user wait / suspended time are excluded from the definition.
- Documents that v1 has no calibration data and why (`actual_active_minutes` deferred until pause-aware timing exists; never derive actives by subtracting `claimed_at` from `completed_at`).

**The schema** (`skills/do-work/actions/work-reference.md` → Request File Schema — Full Frontmatter):
- Add the optional `estimate:` block: backwards-compatible, display/forecast only, frozen once execution begins, never read by scheduling or gating logic. Existing REQs without it remain valid.
- Amend the `effort_estimate` comment ("a triage bit, not an estimation system") in the **same commit** so the two fields' relationship is stated where the fence was: `effort_estimate` stays the closed two-value chip; `trivial` ⇒ floor estimate via short-circuit; `estimate:` is the forecast block.

**Tests** (`_dev/tests/`):
- Lock-in tests covering: deterministic output (same flags twice → byte-identical), rounding to 5, floor enforcement, and representative signal sets for a small Route A, a focused Route B, an integrated Route C, and a browser-heavy QA REQ.
- Dependency-graph coverage (verify-requests repair): critical-path mode tests — a chain (critical path = sum along the chain) and a diamond/parallel graph (critical path = longest branch, not the total).
- Backwards-compatibility coverage (verify-requests repair): a legacy REQ without an `estimate:` block remains valid — the estimator and every reader treat the field as strictly additive/optional; asserted by a fixture without the block passing untouched plus the full existing suite staying green.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Constraints

- No P80 or other percentile fields; no calendar-time promises. (Batch constraint.)
- Nothing lands in queue-kanban; the three-write-surface sentence in CLAUDE.md must not need amending. (Batch constraint.)
- `actual_active_minutes` and history-based calibration are out of v1. (Batch constraint.)
- Script is a shipped capability, not an accelerator: plain POSIX shell, no compiled toolchain, works at the agent floor.

## Dependencies

Root of the batch — REQ-209 and REQ-210 wire actions to this script and reference file.

## Builder Guidance

**Firm:** determinism, rounding to 5, the floor, low/medium/high confidence, basis echo, no P80, the same-commit `effort_estimate` amendment.
**Builder latitude:** exact file names/paths, flag names, and the point values in the scoring table — calibrate against the archive's REQs by judgment (their wall times are not valid actives data, but their relative sizes are a sanity check for the table's ordering).

## Red-Green Proof
**RED prompt/case:** `skills/do-work/tools/estimate-p50.sh --route C --write-set 12 --browser --persistence --full-suite` — today the script does not exist; the invocation fails.
**Why RED now:** No estimator exists anywhere in the suite; `effort_estimate` is deliberately a two-value chip, not an estimation system.
**GREEN when:** The invocation prints a `p50_active_minutes` value that is a multiple of 5 and ≥ the floor, with confidence and basis; running it twice prints byte-identical output; the `_dev/tests/` determinism lock-in test passes and `maintainer-verify.sh` exits 0.
**Validation:** User confirmed — the script-based deterministic shape was agreed in the capture session.

## Full Context
See `do-work/user-requests/UR-047/input.md` for complete verbatim input.

---
*Source: UR-047 — "Add P50 active-duration estimation to do-work REQs" (estimator inputs, avoid-false-precision, and determinism sections)*
