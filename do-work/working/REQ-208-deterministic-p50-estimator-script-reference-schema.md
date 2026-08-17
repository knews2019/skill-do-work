---
id: REQ-208
title: Deterministic P50 estimator script, reference file, and schema
status: claimed
created_at: 2026-08-16T23:52:07Z
claimed_at: 2026-08-17T00:02:29Z
route: B
user_request: UR-047
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-209, REQ-210]
batch: p50-estimation
write_set: [skills/do-work/tools/estimate-p50.sh, skills/do-work/actions/estimate-reference.md, skills/do-work/actions/work-reference.md, _dev/tests/p50-estimator-determinism.sh, _dev/tests/contract-regressions.sh]
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
- [x] **[PLAN]:** Primes + crew loaded (prime-shell-commands, prime-action-files, general, coding-guardrails). Approach: (1) lock-in test suite first (TDD RED — script absent); (2) bash estimator `tools/estimate-p50.sh` with additive integer scoring calibrated so the spec's worked example (Route C, 12-file write set, browser, persistence, full-suite) lands at 125; nearest-5 rounding via floor((n+2)/5)*5; floor 10; deterministic confidence rubric; `critical-path` subcommand takes `ID:MIN[:dep,...]` triples, awk memoized longest-path with cycle detection; (3) companion `actions/estimate-reference.md`; (4) schema block + effort_estimate amendment in work-reference.md; (5) register suite in contract-regressions.sh.
- [x] **[APPLY]:** Code written exactly as planned; five declared files touched, nothing else.
- [x] **[UNIFY]:** `git diff --stat` reviewed — estimate-p50.sh (clean, shellcheck-0.11.0 warning-severity clean), p50-estimator-determinism.sh (clean, one grep `--` fix during RED→GREEN), estimate-reference.md (clean), work-reference.md (two anchored edits only), contract-regressions.sh (one probe block). No debug artifacts. Full maintainer-verify FAIL set byte-identical to the pre-change baseline (41 environment FAILs); new suite green inside the aggregate.

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

## Triage
**Route: B** - Medium
**Reasoning:** Clear outcome (script + reference + schema + tests) with all conventions discoverable — checks-script style, companion-file pattern, schema layout, and test-suite registration all exist as patterns to follow. Not architectural; no planning needed.
**Planning:** Not required

## Plan
**Planning not required** - Route B: Exploration-guided implementation
*Skipped by work action*

## Exploration
- Script style: `tools/checks/scope-drift.sh` — `#!/usr/bin/env bash`, header comment with usage/exit codes, `set -uo pipefail`, awk for parsing. ShellCheck 0.11.0 at `--severity=warning` lints every tracked `*.sh` (maintainer-verify stage 3).
- Companion-file pattern: `actions/capture-reference.md` — description blockquote naming the owner action, sections pointed to by name, lazy-load note in the header.
- Schema home: `actions/work-reference.md` → "Request File Schema — Full Frontmatter" YAML block; `effort_estimate` comment at its line is the one to amend (same commit).
- Test registration: `_dev/tests/contract-regressions.sh` invokes sibling suites explicitly with a missing-file FAIL guard (see the `update_script_probe` block) — no auto-discovery; my suite must be registered there.
- Board parser (`queue-kanban/frontmatter.go`) uses yaml.v3 into a permissive map — nested `estimate:` block is ignored harmlessly on cards; no board change needed.
- Baseline: `maintainer-verify.sh` exits 1 in this container with 41 pre-existing environment FAILs (missing `just`, tar/gzip exec, stat-mode probes) recorded in the session baseline log; ShellCheck lint and static contract checks are green at baseline.
*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/estimate-p50.sh` (new) — deterministic estimator + critical-path mode
- `skills/do-work/actions/estimate-reference.md` (new) — lazy-loaded signal-extraction + block-template companion
- `skills/do-work/actions/work-reference.md` (modified) — `estimate:` schema block + `effort_estimate` comment amendment
- `_dev/tests/p50-estimator-determinism.sh` (new) — lock-in suite
- `_dev/tests/contract-regressions.sh` (modified) — register the new suite

**Files I will NOT touch:** `skills/do-work-board/tools/queue-kanban/` (board tolerates the nested block; three-write-surface rule stays intact), `skills/do-work/actions/work.md` (REQ-209), `skills/do-work/actions/verify-requests.md` (REQ-210)

**Acceptance criteria (restated from REQ):**
- [x] Script prints p50_active_minutes (multiple of 5, ≥ floor), confidence, basis; identical flags → byte-identical output
- [x] Critical-path mode: total effort + longest-path minutes over depends_on graph, both labeled
- [x] Reference file documents extraction, block template, confidence rubric, P50 meaning, exclusions, no-calibration honesty
- [x] Schema documents optional backwards-compatible estimate: block; effort_estimate comment amended same commit
- [x] Lock-in tests: determinism, rounding, floor, Route A/B/C + browser-heavy pins, chain + diamond graphs, writes-nothing/backwards-compat
- [x] No P80 anywhere

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/estimate-p50.sh` (new) — deterministic estimator: additive integer scoring by route base + signals, nearest-5 rounding, 10-minute floor, deterministic confidence rubric, trivial short-circuit, and a `critical-path` subcommand (awk memoized longest-path with cycle rejection)
- `skills/do-work/actions/estimate-reference.md` (new) — lazy-loaded companion: P50 meaning + exclusions, `estimate:` block template, signal-extraction table → flags, confidence rubric, trivial short-circuit bridge, multi-REQ totals/critical-path presentation, calibration-honesty note
- `skills/do-work/actions/work-reference.md` (modified) — optional backwards-compatible `estimate:` block added to the Full Frontmatter schema; `effort_estimate` fence comment amended in the same commit to state the bridge
- `_dev/tests/p50-estimator-determinism.sh` (new) — lock-in suite: determinism, spec-example pin (125), Route A/B/C + browser-heavy pins, rounding, floor, trivial, low-confidence pin, chain + diamond + unknown-dep + cycle graphs, unrecognized-flag rejection, print-only/backwards-compat probe
- `_dev/tests/contract-regressions.sh` (modified) — registered the suite with the standard missing-file FAIL guard

**What was done:** Built the deterministic estimation foundation: agent extracts judgment signals, the shipped script does the arithmetic, so identical normalized signals always produce identical estimates. Calibrated so the spec's own worked example (Route C, 12-file write set, browser, persistence, full-suite) lands at exactly 125 minutes.

## Decisions
<!-- D-XX counter: none used in Open Questions. -->
- **D-01** (DECIDE & STATE): Script placed at `tools/estimate-p50.sh` (not `scripts/`) per the REQ's write_set hint; it is analytics tooling like `tools/checks/`, not a user-flow helper.
- **D-02** (DECIDE & STATE): Scoring table calibrated against the spec's worked presentation example so Route C + 12-file write set + browser + persistence + full-suite = 125 — makes the spec's own example a lock-in test.
- **D-03** (DECIDE & STATE): Critical-path input is self-contained `ID:MINUTES[:DEP,...]` triples — no file reads keeps the graph math deterministic and testable in isolation.
- **D-04** (DECIDE & STATE): Unknown dependency ids contribute zero minutes — mirrors work.md's `--wave` depth rule where archived dependencies are depth 0.

## Qualification

Passed — `tools/checks/qualify.sh` exit 0 (files exist in diff, P-A-U audited, wiring greps clean, not do-work-only); judgment checks: 5 files substantive, all six requirement groups traced to files, estimator computes real pinned values (no hollow paths).

## Testing

**Tests run:** `bash _dev/tests/p50-estimator-determinism.sh` (standalone), full `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ New suite: all probes passing (determinism, spec-example 125 pin, A/B/C + browser-heavy pins, rounding, floor, trivial, low-confidence, chain/diamond/unknown-dep/cycle graphs, flag rejection, print-only). Full verify: FAIL set byte-identical to the pre-change baseline — 41 pre-existing container-environment failures (missing `just`, tar/gzip exec, stat-mode probes), zero new; ShellCheck 0.11.0 warning-severity clean across all 46 tracked shell files including both new ones.

**Red-green validation:**
- `_dev/tests/p50-estimator-determinism.sh`: ✗ before implementation (FAIL: estimator missing, exit 1 — captured) → ✓ after (all probes passed, exit 0). Traces directly to the REQ's `## Red-Green Proof`: the RED invocation `--route C --write-set 12 --browser --persistence --full-suite` is the suite's determinism + spec-example probe.

*Verified by work action*

---
*Source: UR-047 — "Add P50 active-duration estimation to do-work REQs" (estimator inputs, avoid-false-precision, and determinism sections)*
