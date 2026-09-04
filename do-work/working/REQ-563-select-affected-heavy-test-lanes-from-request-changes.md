---
id: REQ-563
title: 'Select affected heavy-test lanes from request changes'
status: claimed
created_at: 2026-09-03T22:58:23Z
user_request: UR-109
domain: testing
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-539]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-539, REQ-564]
batch: smart-heavy-verification
route: C
planning_at: 2026-09-03T23:42:30Z
exploration_at: 2026-09-03T23:42:30Z
estimate:
  p50_active_minutes: 70
  confidence: medium
  calculated_at: 2026-09-03T23:42:30Z
  basis:
    - Route C
    - typed CLI planning authority
    - shell lane execution integration
    - action and reference contract changes
    - rename-safe Git diff parsing
    - fail-closed coverage fallback
    - TDD and cross-route regression verification
write_set:
  - _dev/tests/heavy-lanes.json
  - _dev/tests/maintainer-verify.sh
  - _dev/tests/contracts/probe-lanes.sh
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go
  - skills/do-work/tools/do-work-cli/internal/heavyverification/
  - skills/do-work/tools/do-work-cli/internal/resultmodel/
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
claimed_at: 2026-09-03T23:40:37Z
---

# Select Affected Heavy-Test Lanes From Request Changes

## What

Replace the all-or-nothing heavy-test trigger with change-aware lane selection. Select only heavy lanes affected by a request's changed paths, explain every selection, and fall back to the complete heavy suite whenever coverage is uncertain.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Define stable heavy-lane identities in one repo-owned manifest, add a typed planner that derives rename-aware changed paths and explains lane selection, make uncertainty select all lanes, and wire an internal selected-lane shell entry point while preserving explicit `--heavy` as force-all.
- [x] **[APPLY]:** Implemented the target-revision-bound heavy lane planner, six stable repository lanes, typed selection command and maintainer verification entry points.
- [x] **[UNIFY]:** Verified Go focused tests, go vet, bash -n, ShellCheck, maintainer verify, and fast contract aggregates.

## Why

The broad heavy gate can rerun unrelated expensive lanes and block an otherwise valid request on timing failures outside the paths it changed.

## Context

REQ-539 (Cut the contract file to the incident core and split the aggregate into fast and heavy) is the prerequisite that gives the heavy aggregate explicit lane boundaries. This request uses those boundaries to make selection change-aware; REQ-564 (Reuse matching per-lane verification evidence for four hours) adds reuse after lane identity exists.

## Detailed Requirements

- Select only heavy lanes affected by the request's changed paths.
- Explain why each lane was selected.
- Fall back to the complete heavy suite whenever coverage is uncertain.
- Preserve `--heavy` as a force-all override that is not narrowed by changed-path selection.
- Expose stable lane identities that per-lane evidence can reference.

## Constraints

- Build this after REQ-539's aggregate split.
- A change-aware selection must never silently trade away coverage; uncertainty selects the complete heavy suite.
- Keep the behavior compatible with the existing explicit heavy-test entry point.

## Dependencies

- Requires REQ-539 (Cut the contract file to the incident core and split the aggregate into fast and heavy).
- REQ-564 (Reuse matching per-lane verification evidence for four hours) requires this request's stable lane-selection result.

## Builder Guidance

Firm on the observable selection, explanation, fallback, and force-all behavior. The representation of the changed-path-to-lane mapping is a builder decision; prefer a condition-based mapping that fails closed when it cannot prove coverage.

## Red-Green Proof

**RED prompt/case:** Run request verification for a change limited to one known heavy lane, then for a changed path that no lane mapping recognizes.
**Why RED now:** The current heavy trigger is all-or-nothing, so the known change launches unrelated heavy lanes and has no change-aware uncertainty fallback to inspect.
**GREEN when:** The known change selects only its affected heavy lane or lanes and reports the reason for each, while the unrecognized path selects the complete heavy suite with an explicit uncertainty reason; explicit `--heavy` still selects every heavy lane.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens; relevant to routing and evidence behavior in action files, but the bare partial-slugged satellite exceeds the 2000-token budget.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens; relevant to safe heavy-lane command selection, but the bare partial-slugged satellite exceeds the 2000-token budget.

## Full Context

See `do-work/user-requests/UR-109/input.md` for complete verbatim input.

---
*Source: Replace the all-or-nothing heavy-test trigger with change-aware lane selection and reusable per-lane evidence. Select only heavy lanes affected by the request's changed paths, explain why each lane was selected, and fall back to the complete heavy suite whenever coverage is uncertain. Cache each successful lane result for at most four hours, but reuse it only when a deterministic fingerprint of its command, test inputs, fixtures, toolchain, and required environment still matches; time alone must never authorize reuse. Allow unaffected lanes to reuse evidence while affected lanes rerun, preserve --heavy as a force-all override, and record whether each lane was executed or reused. Build this after REQ-539's aggregate split.*

## Triage

**Route: C** — Complex

The request adds a new typed decision authority, maps Git changes to executable test lanes, changes the work/clarify contract, and must fail closed across renames and unknown paths.

## Plan

1. Add a versioned heavy-lane manifest with stable IDs, exact executable argv, and typed exact/subtree/suffix-under coverage rules.
2. Add `plan-heavy-verification` to compute a rename-aware, NUL-safe diff, select affected lanes with human-readable reasons, and select every lane when any changed path is not covered or when force-all is requested.
3. Add a lane-scoped maintainer entry point consumed only from the typed plan; keep `--heavy` as the unchanged force-all entry point.
4. Update the work and clarify guidance to plan first, request permission only for the selected commands, and preserve ordinary completion when no heavy lane is selected.
5. Prove one-lane selection, multi-lane selection, rename endpoint handling, unknown-path fallback, invalid-manifest refusal, and explicit force-all behavior with test-first fixtures.

## Exploration

REQ-539 exposes three expensive shell lanes—staged package, updater, and installer—inside `_dev/tests/contracts/probe-lanes.sh`; the maintainer gate also owns strict Node and browser behavior lanes. The current `--heavy-surfaces` interface emits broad globs and Step 6.5 converts any match into one complete `--heavy` run. Existing `gateevidence` is revision-wide and should remain separate from per-lane selection. The CLI already provides strict typed command registration and JSON results, so the new planner belongs in its own internal package rather than in shell pattern matching.

## Scope

The implementation will own the lane manifest, typed selection command/result, lane-scoped shell dispatch, action/reference wording, command registration, and focused tests. REQ-564 alone will own fingerprints, TTL, result persistence, and reuse decisions.

## Pre-Flight

The dependency REQ-539 is completed and archived. The exact current fast baseline launched successfully but recorded the pre-existing flaky `TestProtectedInventoryPersistsLaterXAndRequiresStartedState` failure; the builder must not attribute that known failure to this diff. No heavy test was launched.

## Implementation Summary

**Files changed:**
- `_dev/tests/heavy-lanes.json` (new)
- `_dev/tests/maintainer-verify.sh` (modified)
- `skills/do-work/actions/clarify.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification.go` (new)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modified)

**What was done:** Added target-revision-bound heavy lane planner with typed exact, subtree, and suffix-under coverage rules. Preserved both endpoints of rename-aware diffs, explained selected lanes, and fell back to all lanes for uncovered paths or force-all requests. Integrated at 6d9a28ed.

## Review

**Overall: 95%** | 2026-09-04T00:25:00Z

| Dimension | Score |
|---|---:|
| Requirements | 96% |
| Code Quality | 95% |
| Test Adequacy | 94% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** None.

**Acceptance:** Pass — all lane selection rules, fallback coverage, rename handling, and test-first fixtures verified.
