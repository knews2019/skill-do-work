---
id: REQ-496
title: '[impact-critical] Review fix: Add shared executable already-green repair validator'
status: claimed
priority: now
domain: backend
created_at: 2026-09-02T04:53:21Z
user_request: UR-095
addendum_to: REQ-494
review_generated: true
impact: impact-critical
effort_estimate: effort-substantive
tdd: true
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
depends_on: [REQ-494]
related: [REQ-492]
sweep: true
sweep_key: already-green-repair-shared-validator-missing
claimed_at: 2026-09-03T22:25:28Z
route: C
planning_at: 2026-09-03T22:37:36Z
exploration_at: 2026-09-03T22:37:36Z
dispatch_at: 2026-09-03T22:39:03Z
implementation_at: 2026-09-03T22:57:28Z
builder_handback_at: 2026-09-03T22:57:28Z
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/docs/command-line-guide.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go
  - skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go
  - skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green_test.go
  - skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go
  - skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - _dev/tests/contract-regressions.sh
estimate:
  p50_active_minutes: 60
  confidence: low
  calculated_at: 2026-09-03T22:25:42Z
  basis:
    - Route C
    - shared executable authority across action, validation, completion, and selector seams
    - adversarial Git/staging fixtures and full-suite verification
---

# Review Fix: Add Shared Executable Already-Green Repair Validator

## What

Replace the duplicated prose/test decision for an already-green repository-gate repair with one executable validator consumed by both work and review.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Inventoried repair intake, durable no-op, past-revision gate evidence, canonical completion planning, Git observations, both action consumers, and the parallel contract oracle before freezing the 13-file authority boundary.
- [x] **[APPLY]:** Added the shared executable validator and typed projections, rewired both consumers, and replaced the test oracle with real CLI/Git/completion behavior inside the declared write set.
- [x] **[UNIFY]:** Reviewed all 13 declared paths; focused/full tests, vet, contract regressions, formatting, and diff checks pass with no lifecycle, release, board, generated, or debug drift.

## Requirements

- One executable validator is the sole decision authority for both TDD bypass and no-diff review.
- Match the expected fingerprint to actual repair intake evidence; never self-assert it.
- Validate staged paths against the exact successful canonical-completion result, not an archive prefix.
- Refuse an unrelated staged archive and every ordinary, malformed, nonempty, release-mutated, or over-staged neighbor.
- Drive real REQ/Git state through canonical completion, metadata, and selector behavior.

## Red-Green Proof

**RED prompt/case:** REQ-494's fixture can report eligibility from its own `action_decisions()` oracle while shipped guards or evidence are wrong.

**Why RED now:** Test and prose duplicate the decision instead of consuming one executable authority.

**GREEN when:** Removing or corrupting the shared validator breaks both action consumers; exact intake/result paths pass and an unrelated staged archive refuses.

**Validation:** REQ-494 re-review critical finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Instances

- [ ] `impact-critical` — TDD and review decisions still use a parallel test oracle.
- [ ] `impact-critical` — Fingerprint identity is not sourced from repair intake.
- [ ] `impact-critical` — Archive staging is prefix-authorized rather than exact-result-authorized.

## Triage

**Route: C** — Complex

**Reasoning:** This replaces duplicated action prose/test authority with one executable validator spanning repair intake, exact completion output, staged-path ownership, TDD bypass, no-diff review, metadata, and selector behavior. Planning and source exploration are required before declaring the write boundary.

**Planning:** Required

## Plan

1. Add one read-only `validate-already-green-repair` authority that strictly parses repair intake/no-op evidence, reuses past-revision gate evidence, observes NUL-safe Git state, and derives exact allowed lifecycle paths from canonical completion dry-run output.
2. Emit one typed result with separate `tdd_allowed` and `review_allowed` projections; export the existing gate-evidence comparison rather than reimplementing it.
3. Make `work` and `review-work` consume only their corresponding typed decision, document the contract, and remove the parallel Python decision oracle from contract regressions.
4. Pin malformed/self-asserted evidence, nonempty or release-mutated neighbors, unrelated archive staging, exact canonical staging, and the real completion/metadata/selector tail with RED/GREEN tests.

**Plan validation:** Every requirement maps to the shared validator, its two projections/consumers, exact completion-path authority, or the real lifecycle regression. No publication writer, finalization mutation, selector ordering, schema, board, or release surface is included.

*Generated from delegated exploration; full evidence: `do-work/runs/work-2026-09-03-214500/REQ-496-exploration.md`.*

## Exploration

The current action prose and `_dev/tests/contract-regressions.sh::action_decisions()` are three drifting predicates. The fixture self-asserts its fingerprint, reruns the gate at current HEAD, and prefix-authorizes `do-work/archive/`. The replacement must source identity from canonical repair intake, compare the recorded past revision through `gateevidence`, and treat the exact successful `requeststate` completion dry run as the only staging allowlist.

## Scope

**Files I will touch:** The 13 paths declared in `write_set`, comprising the two action consumers, durable reference/docs/prime, command registration, the new validator package and tests, exported gate-evidence seam, typed result rendering/tests, and the contract regression fixture.

**Acceptance criteria:** Both bypass decisions share one executable authority; intake fingerprint and recorded revision are independently verified; review staging is an exact subset of canonical completion output; malformed, ordinary, dirty, release-mutated, or over-staged neighbors refuse; deleting/corrupting either action consumer breaks the contract fixture.

## Pre-Flight

**Git:** The shared wave baseline was clean at `b051879c` after claims, briefs, and both Route C exploration artifacts were committed.

**Tests:** Direct canonical fast gate passed and was recorded at the shared wave baseline before source dispatch.

**Dependencies:** REQ-494 is completed. REQ-492 is related but not a prerequisite; existing repair-intake, gate-evidence, and canonical request-state authorities are available for composition.

## Implementation Summary

**Builder commit:** `5790b0519b75ed59d4458727e5d7dd6fd6b18e2c`

**Files changed:**
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/docs/command-line-guide.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
- `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green.go`
- `skills/do-work/tools/do-work-cli/internal/repairvalidation/already_green_test.go`
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go`
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `_dev/tests/contract-regressions.sh`

**What was done:** Added `validate-already-green-repair` as the sole executable authority for both `tdd_allowed` and `review_allowed`. It strictly joins repair intake to no-op evidence, verifies the extracted past green revision, observes project/staged paths NUL-safely, and authorizes staging only from exact canonical completion dry-run paths. Both actions now require their typed decision, and the contract fixture invokes the real CLI instead of reimplementing eligibility.

**Builder verification:** Focused RED failed on the missing validator/result symbols. Focused packages, `go vet ./...`, full module tests, contract regressions, and diff hygiene pass after implementation. The real fixture covers exact and over-staged paths, malformed/self-asserted/dirty/release-mutated neighbors, real completion commits/metadata, and selector behavior. Durable evidence is in `do-work/runs/work-2026-09-03-214500/REQ-496-handback.md`.
