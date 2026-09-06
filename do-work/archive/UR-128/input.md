---
id: UR-128
title: 'Test-efficiency improvements A1–A9'
created_at: 2026-09-06T13:16:35Z
requests: [REQ-606, REQ-607, REQ-608, REQ-609, REQ-610, REQ-611, REQ-612, REQ-613, REQ-614]
word_count: 8
---

## Summary
Capture all nine areas of the selected test-efficiency report as separate work items. Preserve measurement, efficiency, correctness, and conditional no-change criteria. This capture does not run the queue.

## Extracted Requests
| Report area | Request | Work |
| --- | --- | --- |
| A1 | REQ-606 | Establish an honest performance baseline |
| A2 | REQ-607 | Build the integration CLI once per tested source |
| A3 | REQ-608 | Run the inventory data matrix without subprocesses |
| A4 | REQ-609 | Copy prepared recovery states, not just empty repositories |
| A5 | REQ-610 | Make fixture integrity checks cheaper |
| A6 | REQ-611 | Batch shell audits while preserving diagnostics |
| A7 | REQ-612 | Remove redundant Go test discovery |
| A8 | REQ-613 | Reduce the cost of proving tests can be reused |
| A9 | REQ-614 | Batch repeated Git reads in the exercised code |

## Batch Constraints
- Preserve current scenarios and meaningful assertions; keep real process, signal, Git transaction, and launcher coverage wherever the boundary is what is tested.
- Improve actual work efficiency. Lower concurrency, larger budgets, moving checks into a slower tier, or skipped coverage are not accepted as speedups.
- Compare the same selected work before and after at fixed concurrency; report wall time, total CPU including child processes, and the relevant work counts. Use the baseline method and distinguish cold/warm builds and reused results.
- Prove that optimized tests still reject representative deliberately introduced defects. Do not introduce decorative tests that mirror implementation.
- Recheck the live tree: the report is historical evidence, not proof that a mechanism remains unchanged. Preserve ongoing work and already implemented optimizations.
- This invocation captures intent only; future implementation is a separate do-work run.
- Sequence: A1 first; A2 and A3 after A1; A4 after both A2 and A3. A5–A9 use A1 evidence. No additional dependency is inferred solely from shared files.
- The report distinguishes observed mechanisms from unmeasured speedups; retain that uncertainty. Conditional work may end with a supported no-change conclusion.
- All nine areas have distinct roots from current pending work; no fold targets were changed.

## Source Provenance
- Original report: ai-reports/2026-09-06_1606_test-efficiency-improvement-proposal/index.html
- Source SHA-256: e0d4cc3198ecc98a69525c5cc1e179e28f179e05192c5b58af6baa7d85e7a946
- Exact report bundle copied under do-work/user-requests/UR-128/assets/; its source excerpts and diagnostic images remain available.
- The following nine prompts are extracted verbatim from the displayed HTML pre elements; surrounding report context is preserved in the copied source.

## Full Verbatim Input
> ```
> do-work capture-request from A1 to A9, see file:///Users/t2/Desktop/e1-experimental-repos/skill-do-work2/ai-reports/2026-09-06\_1606\_test-efficiency-improvement-proposal/index.html
> ```

## Referenced Capture Prompts

### A1: Establish an honest performance baseline
> ```
> do-work capture-request: Establish a reproducible test-efficiency baseline for this repository. Extend or reuse the existing duration tooling under _dev/tests instead of building a new benchmark framework. Record source revision and working-tree changes, toolchain, test selection, cache condition, fixed concurrency, whole-selection wall time, total CPU including descendants, and subprocess counts by executable where supported. Distinguish per-file accumulated Go durations from package and gate wall time. Measure the inventory, finalization, session-start, shell-audit, and heavy CLI-build cases. Use at least three comparable runs and report medians and spread; separate cold and warm builds and reused-stage results. Instrumentation must be opt-in and preserve exit status, cancellation, and cleanup. Produce an evidence table that prioritizes removable work; state unsupported counters explicitly. Capture only; do not execute the queue.
> ```

### A2: Build the integration CLI once per tested source
> ```
> do-work capture-request: Eliminate redundant builds of the same CLI executable in integration tests. Start at skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go, where three test functions call buildTestCLIBinary. Build once per test binary with correct lifetime and cleanup. Audit other package-local CLI builders and the existing DO_WORK_TEST_DO_WORK_CLI_BINARY facility; extend sharing only where measured savings justify it. A supplied executable must match the tested source tree and build configuration, and direct standalone package runs must still work. Never share across different worktree contents or race/build flags. Keep actual subprocess signal, exit-status, JSON, and launcher assertions. Acceptance: fewer go build invocations and improved measured CPU and elapsed time for the same heavy selection; stale or wrong binaries cannot produce a green result. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

### A3: Run the inventory data matrix without subprocesses
> ```
> do-work capture-request: Make the inventory tests efficient by separating byte parsing from Git acquisition in skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go. Move the complete synthetic status matrix and secret/origin/ambiguity cases in inventory_test.go to direct in-process tests with explicit expected rows and findings. Cover malformed NUL records, missing rename origins, ordering, metadata exclusions, and cross-row secret promotion. Trace the current retained-command path and document which assertions are independent and which exercise the same implementation twice. Retain representative real-Git and end-to-end CLI/compatibility output checks, including staged-addition then deletion and secret rename behavior. Do not replace end-to-end coverage with a self-comparison. Acceptance: all original data cases remain traceable, synthetic matrix subprocess count falls substantially, and both CPU and elapsed time improve. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

### A4: Copy prepared recovery states, not just empty repositories
> ```
> do-work capture-request: Reduce repeated fixture setup in finalization recovery tests. Inspect seedSemanticLegacyTail and related repeated seed helpers in skills/do-work/tools/do-work-cli/internal/finalization. The initialized Git template already exists; extend reuse only to identical committed or interrupted baseline states that recur across cases. Each case must receive an independent filesystem copy with correct paths, permissions, index, HEAD, and journal semantics. Do not share mutable repositories or precompute the recovery behavior being asserted. Keep real Git commits, hooks, failure injection, rollback, and every durable-phase recovery case. Benchmark construction and copying separately; retain only templates that reduce total CPU and elapsed time. Acceptance includes isolation proof and representative mutation tests for wrong provenance, missing recovery steps, and foreign-file damage. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

### A5: Make fixture integrity checks cheaper
> ```
> do-work capture-request: Optimize shared fixture integrity verification in _dev/tests/session-start-hook-behavior.sh. Measure shared_skill_immutable_digest separately from hook execution. The shared module root already exists; do not recreate that optimization. Replace repeated process-heavy traversal/hash/sort plumbing with a simpler measured equivalent if worthwhile. Preserve after-each-case detection of changed bytes, added or removed paths, and the current version.md exclusion; explicitly test any additional file-type or permission guarantees claimed. Do not use mtime-only equality or move all checking to suite end. Acceptance: intentional shared-tree mutation still fails in the offending case, the same hooks/tools content is checked, and CPU and elapsed time improve. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

### A6: Batch shell audits while preserving diagnostics
> ```
> do-work capture-request: Reduce redundant work in the shell audit suite. Profile _dev/tests/action-shell-blocks.sh, quiet-grep-pipeline-audit.sh, prescribed-shell-canonicalization.sh, and the maintainer ShellCheck stage. Build a coverage/options map before removing duplicate scans. Batch compatible ShellCheck files or reuse extraction within one invocation while retaining source-versus-fence options, exact original path/line attribution, all scanner fixtures, and standalone script operation. Keep bash syntax isolation per fragment. Preserve enumeration of every currently covered file and fence, including unsafe quiet-grep forms. Acceptance: fewer process starts and repeated reads, unchanged diagnostic correctness on invalid fixtures, and improved CPU and elapsed time. Avoid a broad parser framework. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

### A7: Remove redundant Go test discovery
> ```
> do-work capture-request: Remove avoidable pre-execution test discovery in _dev/tests/run-go-tests-with-budget.sh. Measure the current go test -list plus Python prefix filter followed by go test. Evaluate a single-invocation selection using flags supported by the repository toolchain, but preserve the exact intended test set, anchored prefix exclusions, explicit caller selections, package patterns, budget attribution, JSON handling, exit status, and the existing no-fast-tests refusal. Compare selected names before and after, including empty selection and special regex characters. Do not simply drop the empty-selection guard or silently skip tests. Acceptance: one fewer discovery invocation when possible and measured CPU/elapsed savings; if equivalence requires more machinery than it saves, document that outcome. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

### A8: Reduce the cost of proving tests can be reused
> ```
> do-work capture-request: Profile and simplify fast-stage evidence computation without weakening invalidation. Scope _dev/tests/heavy-runtime-fingerprint.py, fast-stages.json, maintainer-verify.sh, and internal/heavyverification. Measure file/tool bytes hashed, process launches, decision/record time, reuse hits, and miss reasons. Remove duplicate reads of identical resolved inputs within one coherent snapshot only when identity is proven; keep post-run revalidation to catch changes during tests. Evaluate smaller stage groups only if pure versus repository-reading test boundaries can be demonstrated. Retain all required source, untracked/ignored input, environment, toolchain, and external-input checks; opaque inputs must still force execution. Reject stale success after failed/interrupted runs or source/tool changes. Acceptance: equivalent invalidation mutation results and lower measured total cost, including misses. Do not build a new general cache framework. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

### A9: Batch repeated Git reads in the exercised code
> ```
> do-work capture-request: Investigate repeated Git subprocess work in finalization discovery and recovery. Trace a representative slow recovery selection before changing production code; group commands by repository, revision, path, and intervening mutation. If duplicate immutable reads or compatible per-path reads dominate, implement the smallest operation-scoped batching or reuse. Use a fixed commit identity for historical data; never reuse working-tree/index/HEAD results across writes without valid invalidation. Preserve merge-aware diffs, unusual filenames, linked worktrees, hooks, cancellation, partial failures, recovery ordering, and foreign-file safety. Acceptance: fewer Git processes, lower CPU and elapsed time, identical outcomes, and mutation tests that catch stale reads after commits/index changes. If tracing shows no worthwhile repetition, finish with evidence and no speculative cache. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```
