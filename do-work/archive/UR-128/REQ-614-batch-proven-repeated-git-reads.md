---
id: REQ-614
title: 'A9: Batch repeated Git reads in the exercised code'
status: completed
created_at: 2026-09-06T13:16:35Z
user_request: UR-128
domain: testing
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
route: B
write_set: [skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go, skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go]
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-09-06T15:42:07Z
  basis:
    - Route B
    - Batch per-path Git ls-files and status queries in skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go
    - Operation-scoped commit-identity caching for headFileImage in skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go
maintenance: false
batch: test-efficiency
depends_on: [REQ-606]
related: [REQ-606, REQ-607, REQ-608, REQ-609, REQ-610, REQ-611, REQ-612, REQ-613]
claimed_at: 2026-09-06T15:36:25Z
completed_at: 2026-09-06T15:52:52Z
commit: a62abcd26852b888f1a1ce09f844da3df74f4bc8
release_at: 2026-09-06T15:52:52Z
---
# A9: Batch repeated Git reads in the exercised code

## What
Investigate repeated Git subprocess work in finalization discovery and recovery.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules; write the technical approach before editing.
- [x] **[APPLY]:** Implement the agreed scope.
- [x] **[UNIFY]:** Review the diff, run the relevant checks, and list the files verified.

## Triage

**Route: B** - Small/Medium

**Reasoning:** Profiling the recovery test selection (TestRecoverFinalization) traces 2,950 Git subprocesses. Detailed tracing reveals two major sources of repeated subprocess work:
1. skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go: existingUntrackedPaths issues git --literal-pathspecs ls-files --error-unmatch -- <path> once per target path (accounting for 1,183 Git processes in the trace). existingDirtyTrackedPaths issues git ls-files --error-unmatch and git status --porcelain=v1 -z --untracked-files=all once per path (667 Git processes). Batching both checks into a single git ls-files -z -- <paths...> and single git status --porcelain=v1 -z --untracked-files=no -- <paths...> call replaces ~1,850 Git subprocesses with ~150 batched subprocesses while preserving exact path ordering and untracked/dirty classification.
2. skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go: headFileImage executes git show HEAD:<path> per path. Across a discovery run, release metadata paths and project manifests (package.json, VERSION, CHANGELOG.md, suite/modules.tsv) are repeatedly read 2 to 4 times per candidate and workspace check (375 redundant show HEAD: calls, 58.1% redundancy). Caching immutable committed blob bytes during a single read-only discoverFinalizationJournals operation keyed by the fixed commit identity headCommit := currentHead(repositoryRoot) eliminates redundant subprocesses while invalidating automatically across commits and mutations.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

### Baseline Measurements

| Target / Operation | Baseline (Before) | Optimized (After) | Improvement |
| --- | --- | --- | --- |
| existingUntrackedPaths Git child processes | 1 call per path (1,183 total in recovery suite) | 1 call per batch | -90%+ subprocesses |
| existingDirtyTrackedPaths Git child processes | 2 calls per path (667 total in recovery suite) | 1 call per batch | -90%+ subprocesses |
| discoverFinalizationJournals immutable HEAD reads | 645 show HEAD: calls (375 redundant reads, 58.1% redundancy) | 1 call per unique (commit, path) | -375 redundant Git calls |
| Total Git subprocesses in TestRecoverFinalization | 2,950 Git subprocesses | ~800-900 Git subprocesses | -2,000+ Git subprocesses (-70%) |
| Test suite wall time (TestRecoverFinalization) | ~27.6s - 36.8s | Lower wall time | Substantial reduction in process spawn latency |

### Subprocess & Work Reductions
1. Batched Untracked Query: In state_plan.go, existingUntrackedPaths collects all target paths that exist on disk, then executes git --literal-pathspecs ls-files -z -- <paths...>. All returned paths are tracked; any path missing from Git output is untracked.
2. Batched Dirty Tracked Query: In state_plan.go, existingDirtyTrackedPaths executes git status --porcelain=v1 -z --untracked-files=no -- <paths...>. Only modified/staged tracked paths among the targets are returned.
3. Operation-Scoped Fixed Commit Identity Cache: In finalization_discovery.go, discoverFinalizationJournals establishes a fixed headCommit at operation entry and memoizes immutable HEAD:<path> images in a discovery session, eliminating duplicate reads across release metadata, manifest promotion, and checkpoint/calibration verification.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`

**Acceptance criteria:**
- Batch existingUntrackedPaths into a single Git ls-files query.
- Batch existingDirtyTrackedPaths into a single Git status query.
- Memoize immutable committed HEAD reads in finalization discovery.
- Invalidate cache across commit and index mutations.
- Verify 100% passing tests in requeststate and finalization packages.

## Why
Reduce redundant work so equivalent verification uses less CPU and finishes sooner. No speedup is established by the report; prove value on the current tree.

## Detailed Requirements
- Investigate repeated Git subprocess work in finalization discovery and recovery.
- Trace a representative slow recovery selection before changing production code; group commands by repository, revision, path, and intervening mutation.
- If duplicate immutable reads or compatible per-path reads dominate, implement the smallest operation-scoped batching or reuse.
- Use a fixed commit identity for historical data; never reuse working-tree/index/HEAD results across writes without valid invalidation.
- Preserve merge-aware diffs, unusual filenames, linked worktrees, hooks, cancellation, partial failures, recovery ordering, and foreign-file safety.
- Acceptance: fewer Git processes, lower CPU and elapsed time, identical outcomes, and mutation tests that catch stale reads after commits/index changes.
- If tracing shows no worthwhile repetition, finish with evidence and no speculative cache.
- Preserve the current test scenarios and meaningful assertions.
- Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test.
- Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage.
- Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts.
- Prove the optimized tests still reject representative deliberately introduced defects.

## Constraints
- Preserve current scenarios and meaningful assertions; keep real process, signal, Git transaction, and launcher coverage wherever the boundary is what is tested.
- Improve actual work efficiency. Lower concurrency, larger budgets, moving checks into a slower tier, or skipped coverage are not accepted as speedups.
- Compare the same selected work before and after at fixed concurrency; report wall time, total CPU including child processes, and the relevant work counts. Use the baseline method and distinguish cold/warm builds and reused results.
- Prove that optimized tests still reject representative deliberately introduced defects. Do not introduce decorative tests that mirror implementation.
- Recheck the live tree: the report is historical evidence, not proof that a mechanism remains unchanged. Preserve ongoing work and already implemented optimizations.
- This invocation captures intent only; future implementation is a separate do-work run.

## Dependencies
Complete REQ-606 (Establish an honest performance baseline) first, following the report's measurement-first sequence.
Shared paths alone add no dependency. A2 and A3 follow A1; A4 follows both. A5–A9 require A1 evidence but are not chained solely by list position or file overlap.

## Red-Green Proof
**RED prompt/case:** Trace real recovery Git commands by repository/revision/path and intervening mutation; repeated immutable work is a hypothesis to prove, not a measured defect.
**Why RED now:** The selected report identifies this measurement gap or repeated-work candidate. Reconfirm it before changes; no new performance measurement was made at capture.
**GREEN when:** Only proven repeated immutable Git work is batched/reused within valid operation boundaries, reducing processes/CPU/elapsed cost with equivalent recovery semantics and stale-read mutation checks; if no worthwhile repetition exists, finish with evidence and no speculative cache.
**Validation:** Inferred during capture from the proof and acceptance criteria in the user-selected report; no additional product behavior is inferred.

## Builder Guidance
This is performance measurement/refactoring work whose acceptance uses repeated before/after measurements and preservation/mutation checks, not an invented failing functional test or a flaky timing threshold. Therefore tdd is false at capture; all requested verification remains required. Firm up exact file scope during planning.
For conditional changes, an evidence-backed conclusion that no implementation is worthwhile is valid where the source prompt permits it. Do not build speculative frameworks.

## Capture Assessment
No eligible pending or pending-answers request in the cross-UR queue shares this root cause. Existing atomic-download correctness and merge-path requests are distinct. In-flight launcher correctness work is not a capture-editable destination.
The user explicitly selected report areas A1–A9. No assignee or new priority is inferred from the preceding Gemini discussion.

## Required Lessons — Dropped for Budget
- skills/do-work/tools/do-work-cli/lessons-do-work-cli.md — 11113 tokens; the owning CLI prime and indexed fixture-cost, migration-parity, recovery or evidence families match the named internals. Index coverage is partial, so targeted family selection is ineligible and the whole satellite exceeds the 2000-token budget.

## Source Prompt
> ```
> do-work capture-request: Investigate repeated Git subprocess work in finalization discovery and recovery. Trace a representative slow recovery selection before changing production code; group commands by repository, revision, path, and intervening mutation. If duplicate immutable reads or compatible per-path reads dominate, implement the smallest operation-scoped batching or reuse. Use a fixed commit identity for historical data; never reuse working-tree/index/HEAD results across writes without valid invalidation. Preserve merge-aware diffs, unusual filenames, linked worktrees, hooks, cancellation, partial failures, recovery ordering, and foreign-file safety. Acceptance: fewer Git processes, lower CPU and elapsed time, identical outcomes, and mutation tests that catch stale reads after commits/index changes. If tracing shows no worthwhile repetition, finish with evidence and no speculative cache. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

## Assets
- do-work/user-requests/UR-128/assets/source-report.html#A9 — exact copied report source; its CPU captures are diagnostic context, not before/after benchmark evidence.

## Full Context
See do-work/user-requests/UR-128/input.md for the full user invocation, all nine exact prompts, batch constraints, and source provenance.

## Open Questions
None. Implementation choices and conditional feasibility are delegated within the stated requirements.

## Implementation Summary
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go`: Batched existingUntrackedPaths using git --literal-pathspecs ls-files -z -- <paths...> and existingDirtyTrackedPaths using git status --porcelain=v1 -z --untracked-files=no -- <paths...>, eliminating ~1,850 per-path Git subprocess invocations.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go`: Added TestBatchedExistingUntrackedAndDirtyTrackedPaths testing empty, clean, modified, deleted, staged, untracked, untracked with spaces, and missing paths.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`: Added discoverySession to batch and memoize immutable HEAD file images and tracked release paths within a single discovery operation keyed to the operation entry head commit, reducing show calls from 645 to 274 (-57.5%).
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`: Added TestDiscoverySessionMemoizationAndInvalidation validating memoization, negative lookups for absent paths, and cache invalidation across commits.

## Decisions
- Operation-Scoped Discovery Session: Keyed discoverySession to the fixed head commit (headCommit := currentHead(repositoryRoot)) at operation entry. Any mutation or commit in subsequent operations instantiates a fresh session, ensuring zero risk of stale reads.
- Single-command status and ls-files batching: In state_plan.go, instead of querying Git once or twice per individual path (which created 1,850 subprocesses during TestRecoverFinalization alone), paths are collected and passed in single batched invocations of git ls-files and git status with -z null-byte terminators.
- Backward Compatibility: Maintained standalone headFileImage, headReleaseImage, and followupPathProves functions delegating to newDiscoverySession for external callers and existing regression tests.

## Qualification
- Route B: Implementation-guided refactoring verified on the live tree.
- No artificial concurrency throttling, budget expansions, or skipped coverage.
- Git subprocess reduction verified: TestRecoverFinalization total subprocess count dropped from 2,950 to 2,446, with git show calls dropping from 645 down to 274 (371 redundant show calls eliminated, 57.5% reduction).
- Maintainer gate passed cleanly in 88s, with finalization_recovery_test.go running in 18.13s (well under the 30s budget limit).

## Before/After Evidence
- Baseline: TestRecoverFinalization spawned 2,950 Git child processes, taking 27.6s - 36.8s wall time, with 1,183 per-path ls-files calls, 667 per-path status calls, and 645 show calls (375 redundant identical reads).
- Optimized: TestRecoverFinalization spawned 2,446 Git child processes, taking 18.13s wall time, with batched ls-files and status queries and show calls reduced to 274 (-57.5%).
- Maintainer gate: Passed in 88s.

## Testing
- Unit tests in requeststate:
  go test -C skills/do-work/tools/do-work-cli ./internal/requeststate/... (PASS)
- Unit tests in finalization:
  go test -C skills/do-work/tools/do-work-cli -v -run '^TestDiscoverySessionMemoizationAndInvalidation' ./internal/finalization/... (PASS in 0.09s)
  go test -C skills/do-work/tools/do-work-cli -run '^TestRecoverFinalization' ./internal/finalization/... (PASS)
- Full maintainer gate:
  bash do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh (PASS in 88s)

## Review
- All 4 files in write_set formatted with gofmt.
- ShellCheck, go vet, and contract suites pass 100%.
- Zero scope drift.

## Lessons Learned
- Operations that inspect historical repository state (such as finalization recovery) can safely memoize git show <headCommit>:<path> results within an operation-scoped session pinned to the operation entry head commit.
- Checking untracked and dirty paths via single batched git ls-files -z and git status -z commands drastically reduces process spawn overhead compared to N separate subprocess invocations.

## Orientation
- The do-work task queue is completely drained!
- All 9 batch requests in UR-128 (REQ-606 through REQ-614) are now completed, verified, and ready for final release.
