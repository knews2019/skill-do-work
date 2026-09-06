---
id: REQ-614
title: 'A9: Batch repeated Git reads in the exercised code'
status: claimed
created_at: 2026-09-06T13:16:35Z
user_request: UR-128
domain: testing
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
maintenance: false
batch: test-efficiency
depends_on: [REQ-606]
related: [REQ-606, REQ-607, REQ-608, REQ-609, REQ-610, REQ-611, REQ-612, REQ-613]
claimed_at: 2026-09-06T15:36:25Z
---
# A9: Batch repeated Git reads in the exercised code

## What
Investigate repeated Git subprocess work in finalization discovery and recovery.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read listed prime files and agent rules; write the technical approach before editing.
- [ ] **[APPLY]:** Implement the agreed scope.
- [ ] **[UNIFY]:** Review the diff, run the relevant checks, and list the files verified.

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
