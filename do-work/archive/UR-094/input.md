---
id: UR-094
title: 'Compute wave depth from satisfied duplicate records'
created_at: 2026-09-01T19:54:39Z
requests: [REQ-490]
word_count: 252
---

# Compute Wave Depth From Satisfied Duplicate Records

## Summary

Capture the one finding accepted by the `do-work validate-feedback` triage of an external review: the dependency graph marks a dependent satisfied when every duplicate record of its dependency is terminal-successful, but the wave-depth calculator in `next` re-derives depth from the ambiguous node's blank status and assigns depth 1, so `next --wave 0` excludes ready work with WAVE-MISMATCH.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-490 | Make `queueDependencyDepth` agree with the graph's duplicate-satisfied verdict so a dependent of a satisfied duplicate dependency lands at wave 0. |

## Batch Constraints

- Capture only; implementation belongs to a later `do-work run`.
- Finding provenance (verbatim claim, P2 severity, triage evidence, Surface-cost N/A) travels with the REQ.

## Full Verbatim Input

> ```
> capture the accepted one
> 
> ### Finding 1: Wave depth ignores satisfied duplicate records  ·  P2
> - **Source:** external review comment pasted to `do-work validate-feedback` (Full review comments), 2026-09-01
> - **Verbatim claim:** [P2] Compute wave depth from satisfied duplicate records — [prj]/.claude/skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go:123-125
>   When every duplicate dependency record is terminal-successful, this branch marks the dependency satisfied, but its graph node remains ambiguous with an empty status. queueDependencyDepth consequently assigns the dependent depth 1 instead of 0, so next --wave 0 incorrectly excludes otherwise-ready work with WAVE-MISMATCH; the duplicate-status aggregation must also feed depth calculation.
> - **Verdict:** Accept
> - **Evidence:** `skills/do-work/tools/do-work-cli/internal/dependencygraph/dependency_graph.go:70-73` blanks `RequestStatus` for any ID with multiple records; `:123` then marks the dependent satisfied via `duplicateStatusesSatisfied`. But `internal/nextselection/next_selection.go:361-365` recomputes depth from the node's `RequestStatus` alone: an empty status is neither terminal-success (depth -1) nor pending (recurse), so the dependency contributes 0 and the dependent lands at depth 1. `next_selection.go:192-194` then excludes it with WAVE-MISMATCH before the readiness switch runs.
> - **Reasoning:** Confirmed by reading both functions; the graph and the wave calculator disagree on the same input. The duplicate-satisfied rule was added deliberately in 2c82ef12 ("duplicate readiness"), so the graph side is the intended truth and the depth side is the stale one. No queued REQ covers this.
> - **Surface-cost:** N/A (direct bug fix).
> - **Remedy:** In `queueDependencyDepth`, treat a dependency as satisfied (-1) when it is absent from the dependent's `UnmetDependencies` and not pending, instead of re-deriving from `RequestStatus`; add a lock-in test with two terminal-success records for one ID and `--wave 0`.
> ```

---
*Captured: 2026-09-01T19:54:39Z*
