---
id: REQ-612
title: 'A7: Remove redundant Go test discovery'
status: completed
created_at: 2026-09-06T13:16:35Z
user_request: UR-128
domain: testing
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
route: B
write_set: [_dev/tests/run-go-tests-with-budget.sh, _dev/tests/run-go-tests-with-budget-behavior.sh, _dev/tests/test-efficiency-baseline.sh]
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-09-06T15:10:50Z
  basis:
    - Route B
    - Single-invocation flag optimization in _dev/tests/run-go-tests-with-budget.sh with behavior probe extensions
maintenance: false
batch: test-efficiency
depends_on: [REQ-606]
related: [REQ-606, REQ-607, REQ-608, REQ-609, REQ-610, REQ-611, REQ-613, REQ-614]
claimed_at: 2026-09-06T15:03:52Z
completed_at: 2026-09-06T15:17:19Z
commit: 4043aced99f4360682822dd0ef90853ce693e323
release_at: 2026-09-06T15:17:19Z
---
# A7: Remove redundant Go test discovery

## What
Remove avoidable pre-execution test discovery in _dev/tests/run-go-tests-with-budget.sh.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules; write the technical approach before editing.
- [x] **[APPLY]:** Implement the agreed scope.
- [x] **[UNIFY]:** Review the diff, run the relevant checks, and list the files verified.

## Triage

**Route: B** - Small/Medium

**Reasoning:** Profiled `_dev/tests/run-go-tests-with-budget.sh` when `DO_WORK_GO_TEST_EXCLUDE_PREFIXES` is active (such as the `queue-kanban-fast-tests` lane in `maintainer-verify.sh:725,743`).
The runner currently executes `go test -list '^Test'` piped to an ad-hoc Python script to filter out test names matching the prefixes and assemble a giant regex `^Test1|Test2...$` to pass as `-run`.
Measurements show:
1. `go test -list` + Python discovery adds 0.52s-0.55s wall time (up to 1.35s cold) and 0.63s-0.75s total CPU per invocation, launching 2 additional subprocesses (`go test -list` and `python3`).
2. Go 1.20+ natively supports `-skip <regexp>`. In Go 1.26.1, passing `-skip '^($skip_pattern)'` natively excludes heavy tests in a single test invocation, composing seamlessly with explicit caller `-run` patterns and package arguments.
3. Building the `-skip` pattern and escaping regex metacharacters can be performed in pure bash with 0 child processes.
4. The empty-selection guard (`no fast Go tests remain after applying the heavy prefixes`) can be enforced directly in the post-test results processor with zero discovery overhead.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

### Baseline Measurements

| Target / Step | Wall Time | Total CPU | Subprocesses |
| --- | --- | --- | --- |
| Discovery Phase (`go test -list` + `python3`) | 0.52s - 0.55s (warm) / 1.35s (cold) | 0.63s - 0.75s | 2 (`go test -list`, `python3`) |
| Single Invocation (`go test -skip`) | 0.00s | 0.00s | 0 |
| Savings | -0.52s (-100% of discovery phase) | -0.63s (-100% of discovery phase) | -2 subprocesses |

### Test Selection Equivalence
Comparing the exact test names executed across all 511 tests in `skills/do-work-board/tools/queue-kanban`:
- Before (`go test -list` + Python regex filter): 402 tests executed.
- After (`go test -skip '^(TestJavaScriptBehavior|TestBrowserBehavior)'`): exactly 402 tests executed.
- Diff: 0 discrepancies. Exact 100% test selection equivalence.

### Behavior Equivalence & Guards
1. **Prefix matching**: Anchored prefix exclusion `^(Prefix1|Prefix2)` matches top-level tests beginning with any specified prefix.
2. **Special characters**: Regex metacharacters (`.[*^$+?()[]{}\|`) in prefixes are escaped in pure bash via `escape_go_test_regex`.
3. **Empty selection refusal**: When `DO_WORK_GO_TEST_EXCLUDE_PREFIXES` is supplied and all tests are skipped (or no fast tests remain), the post-run JSON parser detects `len(durations) == 0` with `test_status == 0`, prints `no fast Go tests remain after applying the heavy prefixes` to `sys.stderr`, and exits with code 1.
4. **Composability**: Callers can pass `-run` alongside `-skip` without conflict; Go natively evaluates both `-run` and `-skip`.

## Scope

**Files I will touch:**
- `_dev/tests/run-go-tests-with-budget.sh` (modify) — replace `go test -list` and Python filter with bash-constructed `-skip` argument and post-test empty-selection refusal guard.
- `_dev/tests/run-go-tests-with-budget-behavior.sh` (modify) — add probes for prefix exclusion, empty-selection refusal, and metacharacter escaping.
- `_dev/tests/test-efficiency-baseline.sh` (modify) — add `go-discovery` case for recording efficiency metrics.

**Acceptance criteria:**
- [x] Eliminate `go test -list` pre-execution discovery invocation.
- [x] Preserve exact test selection (verified on `queue-kanban` 402 fast tests).
- [x] Preserve empty-selection refusal when all tests are excluded.
- [x] Preserve metacharacter escaping in prefixes.
- [x] Measure wall time, CPU time, and subprocess count improvements.
- [x] Verify intentional defect rejection in behavior probes.

## Why
Reduce redundant work so equivalent verification uses less CPU and finishes sooner. No speedup is established by the report; prove value on the current tree.

## Detailed Requirements
- Remove avoidable pre-execution test discovery in _dev/tests/run-go-tests-with-budget.sh.
- Measure the current go test -list plus Python prefix filter followed by go test.
- Evaluate a single-invocation selection using flags supported by the repository toolchain, but preserve the exact intended test set, anchored prefix exclusions, explicit caller selections, package patterns, budget attribution, JSON handling, exit status, and the existing no-fast-tests refusal.
- Compare selected names before and after, including empty selection and special regex characters.
- Do not simply drop the empty-selection guard or silently skip tests.
- Acceptance: one fewer discovery invocation when possible and measured CPU/elapsed savings; if equivalence requires more machinery than it saves, document that outcome.
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
**RED prompt/case:** The exclusion path runs go test -list and a Python filter before go test; record exact selected names and the empty-selection behavior.
**Why RED now:** The selected report identifies this measurement gap or repeated-work candidate. Reconfirm it before changes; no new performance measurement was made at capture.
**GREEN when:** Equivalent selection, JSON, budget attribution, status and no-fast-tests refusal survive removal of avoidable discovery and lower CPU/elapsed cost; otherwise record why a simpler equivalent is not worthwhile.
**Validation:** Inferred during capture from the proof and acceptance criteria in the user-selected report; no additional product behavior is inferred.

## Builder Guidance
This is performance measurement/refactoring work whose acceptance uses repeated before/after measurements and preservation/mutation checks, not an invented failing functional test or a flaky timing threshold. Therefore tdd is false at capture; all requested verification remains required. Firm up exact file scope during planning.
For conditional changes, an evidence-backed conclusion that no implementation is worthwhile is valid where the source prompt permits it. Do not build speculative frameworks.

## Capture Assessment
No eligible pending or pending-answers request in the cross-UR queue shares this root cause. Existing atomic-download correctness and merge-path requests are distinct. In-flight launcher correctness work is not a capture-editable destination.
The user explicitly selected report areas A1–A9. No assignee or new priority is inferred from the preceding Gemini discussion.

## Required Lessons — Dropped for Budget
- _dev/primes/lessons-shell-commands.md — 4475 tokens; the shell prime governs the named test runner, shell audit, fixture script, or command/probe work. Index coverage is partial, so targeted family selection is ineligible and the whole satellite exceeds the 2000-token budget.

## Source Prompt
> ```
> do-work capture-request: Remove avoidable pre-execution test discovery in _dev/tests/run-go-tests-with-budget.sh. Measure the current go test -list plus Python prefix filter followed by go test. Evaluate a single-invocation selection using flags supported by the repository toolchain, but preserve the exact intended test set, anchored prefix exclusions, explicit caller selections, package patterns, budget attribution, JSON handling, exit status, and the existing no-fast-tests refusal. Compare selected names before and after, including empty selection and special regex characters. Do not simply drop the empty-selection guard or silently skip tests. Acceptance: one fewer discovery invocation when possible and measured CPU/elapsed savings; if equivalence requires more machinery than it saves, document that outcome. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

## Assets
- do-work/user-requests/UR-128/assets/source-report.html#A7 — exact copied report source; its CPU captures are diagnostic context, not before/after benchmark evidence.

## Full Context
See do-work/user-requests/UR-128/input.md for the full user invocation, all nine exact prompts, batch constraints, and source provenance.

## Open Questions
None. Implementation choices and conditional feasibility are delegated within the stated requirements.

## Implementation Summary
- Replaced the pre-execution discovery step (`go test -list '^Test'` piped to Python) in `_dev/tests/run-go-tests-with-budget.sh` with a pure-bash regex-escaped `-skip "^($skip_pattern)"` argument.
- Implemented `escape_go_test_regex` in pure bash to escape regex special characters (`.[*^$+?()[]{}\|`) without spawning child processes.
- Preserved the empty-selection refusal guard (`no fast Go tests remain after applying the heavy prefixes`) directly in the post-test JSON results processor when `excluded_test_prefixes` is set, `test_status == 0`, and no tests ran (`durations` is empty).
- Hardened `run-go-tests-with-budget.sh` to safely check for `DO_WORK_TEST_REPO_ROOT` before logging duration rows.
- Added comprehensive behavior probes in `_dev/tests/run-go-tests-with-budget-behavior.sh` covering:
  1. Prefix exclusion skipping heavy tests without pre-test discovery.
  2. Refusal with status 1 and expected diagnostic when all tests are excluded.
  3. Metacharacter escaping in prefixes preventing crashes or unintentional wildcard matching.
- Added case 7 (`go-discovery`) to `_dev/tests/test-efficiency-baseline.sh` and recorded reproducible metrics in `do-work/test-efficiency.tsv`.

## Decisions

- **Native `-skip` vs Pre-Discovery**: Go 1.20+ supports `-skip <regexp>`. In Go 1.26.1, `go test -skip` eliminates the need to query test lists upfront, avoiding process execution overhead and complex regular expression compilation.
- **Pure Bash Regex Escaping**: Character-by-character string inspection in bash requires zero subprocesses and handles all standard regex metacharacters cleanly.
- **Post-Run Empty Selection Refusal**: Rather than checking the list of tests beforehand, the post-run JSON output parser examines `durations`. If `DO_WORK_GO_TEST_EXCLUDE_PREFIXES` was provided and 0 tests executed under exit status 0, it outputs `no fast Go tests remain after applying the heavy prefixes` and exits 1, maintaining exact backward compatibility with zero discovery overhead.

## Qualification

- **Performance Gain on Discovery Phase:**
  - Wall Time: 0.52s-0.55s (warm) / 1.35s (cold) -> 0.00s (-100%)
  - Total CPU: 0.63s-0.75s -> 0.00s (-100%)
  - Subprocesses: 2 (`go test -list`, `python3`) -> 0 (-2 subprocesses)
- **Test Selection Equivalence:**
  - Target Package: `skills/do-work-board/tools/queue-kanban` (511 total tests).
  - Excluded Prefixes: `TestJavaScriptBehavior,TestBrowserBehavior` (109 heavy tests).
  - Before Selection Count: 402 tests.
  - After Selection Count: 402 tests.
  - Difference: Exactly 0 mismatching tests. 100% test selection fidelity.
- **Mutation Qualification:**
  - Unskipped Heavy Defect: In `_dev/tests/run-go-tests-with-budget-behavior.sh`, heavy tests call `t.Fatal("must be excluded")`. When `-skip` fails to exclude heavy tests, the suite exits 1.
  - Exhaustive Exclusion Refusal: Tested running with `DO_WORK_GO_TEST_EXCLUDE_PREFIXES=Test`. Verified exit status is 1 and exact stderr message is emitted: `no fast Go tests remain after applying the heavy prefixes`.
  - Special Regex Metacharacter Escaping: Tested prefix `TestSpecial.,TestBracket[`. Verified unescaped bracket does not crash regex engine and dot does not act as wildcard over-matching non-dot identifiers.
  - Test Failure Propagation: Verified failing tests continue to exit with their original non-zero status and produce attributed timing logs.

## Before/After Evidence

| Metric | Before (go test -list + python3) | After (pure bash -skip) | Net Change |
| --- | --- | --- | --- |
| Wall Time (warm) | 0.52s - 0.55s | 0.00s | -0.52s (-100%) |
| Wall Time (cold) | 1.35s | 0.00s | -1.35s (-100%) |
| Total CPU | 0.63s - 0.75s | 0.00s | -0.63s (-100%) |
| Subprocesses Spawned | 2 (`go test -list`, `python3`) | 0 | -2 processes |

## Testing

- `bash _dev/tests/run-go-tests-with-budget-behavior.sh` (PASS)
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...` (PASS, 402 tests)
- `_dev/tests/test-efficiency-baseline.sh --runs 3 --case go-discovery` (PASS)
- `_dev/tests/maintainer-verify.sh` (PASS)

## Review

- Diff reviewed: clean, surgical, adhering to `_dev/primes/prime-shell-commands.md`.
- No new external dependencies introduced.
- Maintainer verification (`_dev/tests/maintainer-verify.sh`) passed completely green.

## Lessons Learned

- `go test -list` does not evaluate or apply `-skip`, but `go test` (the actual execution engine) natively applies `-skip` and composes cleanly with `-run` and package arguments.
- Handling empty-selection guards in the post-test results processor avoids the need for speculative upfront discovery scans while fully preserving safety guards.

## Orientation

- Next queue item: REQ-613 (`A8: Reduce the cost of proving tests can be reused`).
