---
id: REQ-490
title: 'Review fix: compute wave depth from satisfied duplicate records'
status: pending-heavy-testing
priority: now
created_at: 2026-09-01T19:54:39Z
user_request: UR-094
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
review_generated: false
claimed_at: 2026-09-03T22:25:28Z
route: A
dispatch_at: 2026-09-03T22:39:03Z
implementation_at: 2026-09-03T22:41:27Z
builder_handback_at: 2026-09-03T22:41:27Z
integration_at: 2026-09-03T22:42:04Z
testing_at: 2026-09-03T22:45:30Z
review_at: 2026-09-03T22:45:30Z
status_changed_at: 2026-09-03T22:45:30Z
commit: 70b3be19f15c0448fd23008ef16b5ff19881677c
write_set:
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-03T22:25:42Z
  basis:
    - trivial graph-depth correction and one lock-in test
---

# Compute Wave Depth From Satisfied Duplicate Records

## What

When every duplicate record of a dependency is terminal-successful, the dependency graph marks the dependent satisfied, but `queueDependencyDepth` in `next` re-derives depth from the ambiguous node's blank `RequestStatus` and assigns the dependent depth 1 instead of 0. `next --wave 0` then excludes otherwise-ready work with WAVE-MISMATCH. Make the wave-depth calculation honor the same duplicate-satisfied verdict the graph already computed.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this wave-depth-versus-graph-readiness disagreement root cause. REQ-452 (refuse ambiguous explicit request IDs) concerns explicit-target identity, not depth; REQ-488 (keep empty inline frontmatter lists empty) concerns list parsing.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A direct implementation: pin the duplicate-satisfied wave-0 failure, then consume the dependency graph's resolved edge verdict in queue-depth calculation.
- [x] **[APPLY]:** Added the exact RED/GREEN lock-in and changed only the declared selector helper and test file.
- [x] **[UNIFY]:** Reviewed both changed files; focused/full tests, vet, formatting, and diff hygiene pass with no debug, lifecycle, generated, or release drift.

## Context

Finding provenance, carried per the Finding-Closure Ratchet:

- **Source:** external review comment, severity P2, adjudicated by `do-work validate-feedback` on 2026-09-01 with verdict Accept.
- **Verbatim claim:** "[P2] Compute wave depth from satisfied duplicate records — dependency_graph.go:123-125. When every duplicate dependency record is terminal-successful, this branch marks the dependency satisfied, but its graph node remains ambiguous with an empty status. queueDependencyDepth consequently assigns the dependent depth 1 instead of 0, so next --wave 0 incorrectly excludes otherwise-ready work with WAVE-MISMATCH; the duplicate-status aggregation must also feed depth calculation."
- **Evidence:** `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go:70-73` blanks `RequestStatus` for any ID with multiple records; `:123` marks the dependent satisfied via `duplicateStatusesSatisfied`. `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go:361-365` recomputes depth from `RequestStatus` alone: a blank status is neither terminal-success (contributes -1) nor pending (recurse), so the dependency contributes 0 and the dependent lands at depth 1. `next_selection.go:192-194` then excludes it with WAVE-MISMATCH before the readiness switch runs.
- **Surface-cost:** N/A — direct bug fix; no guard, fallback, or validation layer is added.
- The duplicate-satisfied rule was added deliberately in commit 2c82ef12 ("Stabilize fallback archives and duplicate readiness"); the graph side is the intended truth and the depth side is the stale one.

## Requirements

- `queueDependencyDepth` treats a dependency as satisfied (contributing -1) when the graph already resolved it as met — for example when it is absent from the dependent's `UnmetDependencies` and not pending — instead of re-deriving from the dependency node's `RequestStatus`.
- A dependent whose only dependency is a duplicate-satisfied ID is selected by `next --wave 0` and reports `DependencyDepth: 0`.
- Behaviour for single-record dependencies (terminal-success, pending, missing, cyclic, filename-collision ambiguous) is unchanged; a genuinely ambiguous, unsatisfied dependency still yields the existing exclusions.
- Closing check: one lock-in test in `internal/nextselection` naming this failure (two terminal-success records for one dependency ID, dependent pending, `--wave 0`).

## Red-Green Proof
**RED prompt/case:** A repository snapshot where REQ-1 has two records (for example one in `do-work/queue/` and one in `do-work/archive/`), both `status: completed`, and REQ-2 is `status: pending` with `depends_on: [REQ-1]`. Run `do-work-cli --format json next --wave 0`.
**Why RED now:** REQ-2 is excluded with `WAVE-MISMATCH: dependency depth is 1, not requested wave 0`, even though the graph reports REQ-2's dependencies satisfied and `IsReady` true.
**GREEN when:** The same snapshot selects REQ-2 at wave 0 and the projected `DependencyDepth` is 0; a `go test ./internal/nextselection/...` lock-in test pins the case.
**Validation:** Inferred during capture

## Assets
None.

---
*Source: `do-work validate-feedback` triage of 2026-09-01, Finding 1 (Accept); full block preserved in UR-094 input.md.*

## Triage

**Route: A** — Simple

**Reasoning:** The request names the stale depth helper, the already-authoritative graph verdict, the exact failing fixture, and one focused test seam. No design choice or cross-subsystem exploration remains.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`

**Acceptance criteria:** Reuse the graph's already-resolved satisfied edge when computing queue depth; select the exact two-completed-duplicate dependent at `--wave 0` with depth 0; preserve all single-record and genuinely ambiguous dependency behavior.

## Pre-Flight

**Git:** The shared wave baseline was clean at `b051879c` after claims, briefs, and Route C exploration artifacts were committed.

**Tests:** Direct canonical fast gate passed and was recorded at the shared wave baseline before any source dispatch.

**Dependencies:** None. The duplicate-satisfied graph behavior is already present and is the authority this selector-only fix consumes.

## Implementation Summary

**Builder commit:** `87b72c7789d58152f1eaf4253629ba7b12bfa65f`

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`

**What was done:** `queueDependencyDepth` now treats a dependency edge absent from the current node's authoritative `UnmetDependencies` as already satisfied. This includes terminal-success duplicate records whose aggregate graph verdict is satisfied even though their ambiguous node status is blank; pending unmet dependencies still recurse and other unresolved cases retain their prior depth behavior.

**Builder verification:** The new duplicate-satisfied fixture failed with REQ-312 excluded at depth 1 before the edit and passes with selection at wave 0/depth 0 afterward. Focused and full module tests, `go vet ./...`, formatting, and diff checks pass; durable evidence is in `do-work/runs/work-2026-09-03-214500/REQ-490-handback.md`.

## Review

**Verdict:** Pass. The exact two-file integration replaces only the stale status re-derivation with the graph's resolved edge verdict. Satisfied single-record and duplicate dependencies contribute no depth, pending unmet dependencies still recurse, and missing, cyclic, and genuinely ambiguous-unsatisfied paths retain their established exclusions. The lock-in names the reported failure and asserts both selection and projected depth.

## Qualification

Passed — mechanical qualification accepted the exact `4a04b850..70b3be19` integration range. The implementation is limited to the two declared next-selection files, P-A-U evidence is complete, and no builder-authored lifecycle, generated, or release path is present.

## Testing

**Tests run:** focused duplicate-satisfied and wave/fan-out tests; full `internal/nextselection`; `go vet ./...`; full `go test -count=1 ./...`; direct `bash _dev/tests/maintainer-verify.sh`; formatting and diff checks.

**Result:** All focused, module, vet, full-suite, and canonical fast-gate checks pass. The exact canonical gate record is stored at descendant `93ebbfd8`, which contains the `70b3be19` implementation plus only its settled request evidence.

**Red-green validation:** Before the production edit the new fixture excluded REQ-312 with `WAVE-MISMATCH` at depth 1; after the edit it selects REQ-312 at wave 0 with `DependencyDepth: 0`.

## Open Questions

- [ ] Run `bash _dev/tests/maintainer-verify.sh --heavy` at `70b3be19f15c0448fd23008ef16b5ff19881677c`; did it exit 0?
  Recommended: Yes
  Also: No — report the failing lane
