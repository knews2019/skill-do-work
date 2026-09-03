---
id: REQ-490
title: 'Review fix: compute wave depth from satisfied duplicate records'
status: claimed
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
