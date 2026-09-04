---
id: REQ-535
title: 'Review fix: keep same-id filename collisions ambiguous in duplicate readiness'
status: completed
priority: now
created_at: 2026-09-03T12:20:21Z
user_request: UR-103
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-534, REQ-536, REQ-490]
batch: validate-feedback-2026-09-03
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set: [skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go, skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go]
claimed_at: 2026-09-04T12:47:52Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-04T12:48:00Z
  basis:
    - trivial short-circuit
completed_at: 2026-09-04T12:55:00Z
commit: 4d378988f56676efc90048903b1a8abda31f4039
---

# Keep Same-Id Filename Collisions Ambiguous in Duplicate Readiness

## What

`duplicateStatusesSatisfied` in the dependency graph is meant to resolve only exact duplicate records (the same REQ present twice, for example in `queue/` and `archive/`) while filename/frontmatter collisions stay ambiguous. It decides by comparing the collision's claim-path count with the frontmatter-indexed record count. When two files with different filename ids both declare the same frontmatter id and both are terminal-successful, the counts match, the collision is treated as a satisfied duplicate, and dependents become ready despite the mismatched identity.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this root cause. REQ-490 (compute wave depth from satisfied duplicate records) consumes the duplicate-satisfied verdict in the wave-depth calculation; it does not change which records count as duplicates, so the two are independent and listed as related only.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** In `duplicateStatusesSatisfied`, require that every requestFile's FilenameID and TypedRecord.RequestID equal the target requestID so filename/frontmatter collisions remain ambiguous. Add lock-in test in `dependency_graph_test.go`.
- [x] **[APPLY]:** Code written in `dependency_graph.go` and `dependency_graph_test.go`. Scope strictly limited to planned files.
- [x] **[UNIFY]:** Ran `go test ./...` and `git diff --stat`. Verified all tests pass and no debug artifacts in diff.

## Context

Finding provenance, carried per the Finding-Closure Ratchet:

- **Source:** external review comment, severity P1, adjudicated by `do-work validate-feedback` on 2026-09-03 with verdict Accept; full block preserved in UR-103 input.md.
- **Verbatim claim:** "[P1] Keep filename/frontmatter collisions ambiguous — dependency_graph.go:156-158. A filename/frontmatter collision can have the same number of claim paths and frontmatter-indexed records—for example, REQ-020-first.md and REQ-021-second.md both declaring id: REQ-021. If both statuses are successful, this count comparison returns true, causing BuildGraph to suppress the ambiguity and mark dependents ready despite the mismatched identity."
- **Evidence:** `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go:156-158` compares only `len(collision.ClaimPaths)` with `len(requestFiles)`. Reproduced with a scratch test: `REQ-020-first.md` and `REQ-021-second.md` both declaring `id: REQ-021`, both `completed`, and `REQ-030` depending on REQ-021 gave `IsReady=true`, `DependenciesSatisfied=true`, `AmbiguousTargets=[]`. The existing lock-in `TestFilenameFrontmatterCollisionMakesDependencyAmbiguous` (`dependency_graph_test.go:100`) covers only the shape where the two frontmatter ids differ. The duplicate-readiness rule landed in commit 2c82ef12 with no recorded rationale for mismatched filenames; the comment at line 148 promises collisions stay ambiguous.
- **Surface-cost:** N/A, direct bug fix; no guard or validation layer is added.

## Requirements

- `duplicateStatusesSatisfied` returns false when any record indexed under the id has a filename id that differs from its frontmatter id (a filename/frontmatter collision), regardless of statuses. Only exact duplicates (every record's filename id and frontmatter id both equal the id) can satisfy dependents.
- Exact duplicates that are all terminal-successful keep satisfying dependents as today; an exact duplicate with any non-successful status keeps blocking.
- Closing check: one lock-in test in `internal/dependencygraph` naming this failure (two files with different filename ids declaring the same frontmatter id, both `completed`, dependent stays ambiguous and unready).

## Red-Green Proof
**RED prompt/case:** A repository snapshot with `do-work/archive/REQ-020-first.md` (`id: REQ-021`, `status: completed`), `do-work/archive/REQ-021-second.md` (`id: REQ-021`, `status: completed`), and `do-work/queue/REQ-030-dependent.md` (`status: pending`, `depends_on: [REQ-021]`). Build the graph.
**Why RED now:** REQ-030 reports `IsReady=true` with no ambiguous target, although REQ-021's records disagree about their own identity.
**GREEN when:** REQ-030 reports `IsReady=false`, `AmbiguousTargets=[REQ-021]`, `UnmetDependencies=[REQ-021]`, and a `go test ./internal/dependencygraph/...` lock-in test pins the case.
**Validation:** Inferred during capture (RED case reproduced during the triage)

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 3124 tokens, over the 2000-token budget; `slugged: partial`, so the targeted `#collision-fixture-identity` form is not eligible. Matched on the collision-fixture-identity family. The owning prime is listed in `prime_files` instead.

## Assets
None.

---
*Source: `do-work validate-feedback` triage of 2026-09-03, Finding 4 (Accept); full block preserved in UR-103 input.md.*

---

## Triage

**Route: A** - Simple

**Reasoning:** Bug fix with well-defined reproduction, clear requirements, and specific target files in `internal/dependencygraph`.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Scope

**Declared write_set:**
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go`
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go`

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go` (modified)

**What was done:** Updated `duplicateStatusesSatisfied` to verify that each record's `FilenameID` and `TypedRecord.RequestID` both match `requestID`, ensuring filename/frontmatter collisions remain ambiguous. Added `TestSameIdFilenameFrontmatterCollisionKeepsDependencyAmbiguous` lock-in test.

## Testing

**Tests run:** `go test ./internal/dependencygraph/...` and `go test -count=1 ./...`
**Result:** ✓ All passing (all packages in `do-work-cli`)

**Red-green validation:**
- `TestSameIdFilenameFrontmatterCollisionKeepsDependencyAmbiguous`: ✗ failed before fix (REQ-030 reported `IsReady=true` through collided dependency) → ✓ passed after fix

**New tests added:**
- `TestSameIdFilenameFrontmatterCollisionKeepsDependencyAmbiguous` in `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph_test.go`

*Verified by work action*

## Review

**Overall: 100%** | 2026-09-04T12:51:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | low |
| Acceptance | Pass |

**Important findings:**
None

**Minor findings:**
None

**Acceptance:** Pass — Filename/frontmatter collisions where both files declare the same frontmatter ID now remain ambiguous, preventing dependents from prematurely reporting ready.
**Suggested testing:** 0 items
**Follow-ups created:** None (0 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Adding exact identity matching (`requestFile.FilenameID != requestID || requestFile.TypedRecord.RequestID != requestID`) inside `duplicateStatusesSatisfied` cleanly rejects mismatched filename claims without breaking exact duplicate handling.
**What didn't:** Relying solely on `len(collision.ClaimPaths) != len(requestFiles)` failed when two files with different filename IDs both claimed the same frontmatter ID.
**Worth knowing:** Collision entries can have identical counts between claim paths and indexed records when multiple files collide on the same target ID.

