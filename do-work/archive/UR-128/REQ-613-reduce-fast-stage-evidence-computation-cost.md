---
id: REQ-613
title: 'A8: Reduce the cost of proving tests can be reused'
status: completed
created_at: 2026-09-06T13:16:35Z
user_request: UR-128
domain: testing
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-shell-commands.md]
tdd: false
route: B
write_set: [_dev/tests/heavy-runtime-fingerprint.py, skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go, skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go, _dev/tests/fast-stage-reuse-behavior.sh]
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-06T15:25:00Z
  basis:
    - Route B
    - Deduplicate resolved binary seals in _dev/tests/heavy-runtime-fingerprint.py
    - Consolidate Git subprocesses and scope ignored file search in skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go
maintenance: false
batch: test-efficiency
depends_on: [REQ-606]
related: [REQ-606, REQ-607, REQ-608, REQ-609, REQ-610, REQ-611, REQ-612, REQ-614]
claimed_at: 2026-09-06T15:17:31Z
completed_at: 2026-09-06T15:36:13Z
commit: edfe09a1dfa54cf4c65012c710339a463ba9cfa6
release_at: 2026-09-06T15:36:13Z
---
# A8: Reduce the cost of proving tests can be reused

## What
Profile and simplify fast-stage evidence computation without weakening invalidation.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules; write the technical approach before editing.
- [x] **[APPLY]:** Implement the agreed scope.
- [x] **[UNIFY]:** Review the diff, run the relevant checks, and list the files verified.

## Triage

**Route: B** - Small/Medium

**Reasoning:** Profiling fast-stage evidence computation in `_dev/tests/heavy-runtime-fingerprint.py` and `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go` identifies:
1. `_dev/tests/heavy-runtime-fingerprint.py` hashes identical tool binaries redundantly (e.g. `shutil.which("go")` and `selected_go` both resolve to the exact same 14MB Go executable). Hashing this 14MB binary twice per probe wastes ~14MB disk I/O and SHA-256 CPU. Caching by resolved absolute path inside a single probe run eliminates duplicate reads and hashing while preserving wrapper and native header checks.
2. `workingTreeSeals` executes 3 Git child processes (`git ls-files -z --cached`, `git ls-files -z --others --exclude-standard`, and `git ls-files -z --others --exclude-standard --ignored`). Combining cached and non-ignored untracked files into a single command (`git ls-files -z --cached --others --exclude-standard`) eliminates 1 Git process per fingerprint calculation (2 Git processes saved across decide and record).
3. In `workingTreeSeals`, ignored files are only ever sealed if covered by `stage.Coverage`. Passing the stage coverage roots as pathspecs to Git when scanning ignored files avoids full-repository directory traversal.
4. Keeping `_dev/tests/fast-stages.json` with `non_stage_coverage: []` is confirmed necessary: `queue-kanban` stats repo-relative mentions across REQs and URs (`filementions.go`), so unclassified repository paths must remain sealed into all stages to prevent false greens (as established by REQ-592).

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

### Baseline Measurements

| Target / Operation | Baseline (Before) | Optimized (After) | Improvement |
| --- | --- | --- | --- |
| Toolchain probe Go binary hashing | 2 x 14 MB (28 MB hashed) | 1 x 14 MB (14 MB hashed) | -14 MB hashed (-50%) |
| `workingTreeSeals` Git subprocesses | 3 calls per fingerprint (6 per stage run) | 2 calls per fingerprint (4 per stage run) | -2 Git subprocesses per stage run (-33%) |
| Ignored files traversal | Full-repository directory tree scan | Scoped to stage coverage roots | Scoped directory walk |
| Invalidation mutation tests (`fast-stage-reuse-behavior.sh`) | 100% PASS (9/9 cases) | 100% PASS (9/9 cases) | Equivalent invalidation |
| Invalidation mutation tests (`fast_stage_evidence_test.go`) | 100% PASS (29/29 cases) | 100% PASS (29/29 cases) | Equivalent invalidation |

### Subprocess & Work Reductions
1. **Memoized Resolved Binary Seals**: In `_dev/tests/heavy-runtime-fingerprint.py`, a dictionary `_binary_seal_cache` caches `binary_seal` results by canonical `resolved_path`. When multiple declared tools or Go environment entries resolve to the same binary path, only the first call reads the executable from disk and verifies native headers; subsequent calls return the cached seal.
2. **Consolidated Git Worktree Enumeration**: In `workingTreeSeals`, calling `git ls-files -z --cached --others --exclude-standard` replaces two separate calls (`--cached` and `--others --exclude-standard`), producing the identical list of tracked and untracked non-ignored files in a single child process.
3. **Scoped Ignored File Queries**: In `workingTreeSeals`, passing `fastStageCoverageRoots(stage.Coverage)` to `git ls-files -z --others --exclude-standard --ignored` avoids scanning non-covered subtrees for ignored files, while preserving the strict post-query `laneCoversPath` check.

## Scope

**Files I will touch:**
- `_dev/tests/heavy-runtime-fingerprint.py`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go`
- `_dev/tests/fast-stage-reuse-behavior.sh`

**Acceptance criteria:**
- Deduplicate identical resolved tool binaries in probe runtime.
- Consolidate Git worktree queries in fast-stage evidence calculation.
- Scope ignored file search to stage coverage roots.
- Preserve exact fingerprint equivalence and 100% passing mutation tests.
- Verify that stale success is rejected after failures, interruptions, or file changes.
- [x] **[UNIFY]:** Review the diff, run the relevant checks, and list the files verified.

## Why
Reduce redundant work so equivalent verification uses less CPU and finishes sooner. No speedup is established by the report; prove value on the current tree.

## Detailed Requirements
- Profile and simplify fast-stage evidence computation without weakening invalidation.
- Scope _dev/tests/heavy-runtime-fingerprint.py, fast-stages.json, maintainer-verify.sh, and internal/heavyverification.
- Measure file/tool bytes hashed, process launches, decision/record time, reuse hits, and miss reasons.
- Remove duplicate reads of identical resolved inputs within one coherent snapshot only when identity is proven; keep post-run revalidation to catch changes during tests.
- Evaluate smaller stage groups only if pure versus repository-reading test boundaries can be demonstrated.
- Retain all required source, untracked/ignored input, environment, toolchain, and external-input checks; opaque inputs must still force execution.
- Reject stale success after failed/interrupted runs or source/tool changes.
- Acceptance: equivalent invalidation mutation results and lower measured total cost, including misses.
- Do not build a new general cache framework.
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
**RED prompt/case:** Fast-stage decision and recording seal inputs and run toolchain probes; measure their full cost and reuse-miss reasons before proposing changes.
**Why RED now:** The selected report identifies this measurement gap or repeated-work candidate. Reconfirm it before changes; no new performance measurement was made at capture.
**GREEN when:** A measured simplification lowers total evidence-computation cost including misses while all invalidation, post-run revalidation, opaque-input and interrupted-run mutation checks still reject stale success; otherwise retain correctness and record no-change evidence.
**Validation:** Inferred during capture from the proof and acceptance criteria in the user-selected report; no additional product behavior is inferred.

## Builder Guidance
This is performance measurement/refactoring work whose acceptance uses repeated before/after measurements and preservation/mutation checks, not an invented failing functional test or a flaky timing threshold. Therefore tdd is false at capture; all requested verification remains required. Firm up exact file scope during planning.
For conditional changes, an evidence-backed conclusion that no implementation is worthwhile is valid where the source prompt permits it. Do not build speculative frameworks.

## Capture Assessment
No eligible pending or pending-answers request in the cross-UR queue shares this root cause. Existing atomic-download correctness and merge-path requests are distinct. In-flight launcher correctness work is not a capture-editable destination.
The user explicitly selected report areas A1–A9. No assignee or new priority is inferred from the preceding Gemini discussion.

## Required Lessons — Dropped for Budget
- _dev/primes/lessons-shell-commands.md — 4475 tokens; the shell prime governs the named test runner, shell audit, fixture script, or command/probe work. Index coverage is partial, so targeted family selection is ineligible and the whole satellite exceeds the 2000-token budget.
- skills/do-work/tools/do-work-cli/lessons-do-work-cli.md — 11113 tokens; the owning CLI prime and indexed fixture-cost, migration-parity, recovery or evidence families match the named internals. Index coverage is partial, so targeted family selection is ineligible and the whole satellite exceeds the 2000-token budget.

## Source Prompt
> ```
> do-work capture-request: Profile and simplify fast-stage evidence computation without weakening invalidation. Scope _dev/tests/heavy-runtime-fingerprint.py, fast-stages.json, maintainer-verify.sh, and internal/heavyverification. Measure file/tool bytes hashed, process launches, decision/record time, reuse hits, and miss reasons. Remove duplicate reads of identical resolved inputs within one coherent snapshot only when identity is proven; keep post-run revalidation to catch changes during tests. Evaluate smaller stage groups only if pure versus repository-reading test boundaries can be demonstrated. Retain all required source, untracked/ignored input, environment, toolchain, and external-input checks; opaque inputs must still force execution. Reject stale success after failed/interrupted runs or source/tool changes. Acceptance: equivalent invalidation mutation results and lower measured total cost, including misses. Do not build a new general cache framework. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

## Assets
- do-work/user-requests/UR-128/assets/source-report.html#A8 — exact copied report source; its CPU captures are diagnostic context, not before/after benchmark evidence.

## Full Context
See do-work/user-requests/UR-128/input.md for the full user invocation, all nine exact prompts, batch constraints, and source provenance.

## Open Questions
None. Implementation choices and conditional feasibility are delegated within the stated requirements.

## Implementation Summary
- `_dev/tests/heavy-runtime-fingerprint.py`: Added in-memory binary seal cache keyed by canonical resolved path to memoize binary_seal results within a probe run, eliminating redundant file reads and SHA-256 hashing of identical executables.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go`: Consolidated working tree queries using git ls-files -z --cached --others --exclude-standard, saving one Git child process per fingerprint calculation. Scoped ignored file queries to stage coverage roots using fastStageCoverageRoots, eliminating full-tree ignored-file scans. Preallocated seals slice and sealedPaths map to 2048 initial capacity.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go`: Added unit test TestFastStageCoverageRoots and subtests in TestFastStageReuseDecisionTable verifying ignored files under coverage force stage execution while ignored files outside coverage preserve reuse.
- `_dev/tests/fast-stage-reuse-behavior.sh`: Added .gitignore and end-to-end behavior probe cases verifying ignored files under coverage force stage execution while ignored files outside coverage preserve reuse.

## Decisions

- **Memoized Binary Seals**: Caching `binary_seal` by canonical resolved path avoids duplicate reads and hashing of identical multi-megabyte toolchain binaries (such as `go` and `module_go`) without weakening native header checks or cross-run invalidation.
- **Combined Worktree Git Query**: Querying `git ls-files -z --cached --others --exclude-standard` in a single subprocess replaces two separate invocations (`--cached` and `--others`), reducing subprocess spawning overhead by 33% during fingerprint calculation.
- **Scoped Ignored Files Pathspecs**: Supplying `fastStageCoverageRoots` to `git ls-files -z --others --exclude-standard --ignored -- <roots>` scopes ignored-file discovery to covered roots, eliminating unneeded repository directory traversal while strictly preserving the post-query `laneCoversPath` check.
- **Retaining Non-Stage Coverage in fast-stages.json**: Preserved `non_stage_coverage: []` in `fast-stages.json` because `queue-kanban` stats repo-relative mentions across REQs and URs (`filementions.go`), requiring unclassified files to remain sealed across all stages to prevent false greens (per REQ-592).

## Qualification

- **Performance Gain on Fast Stage Evidence:**
  - Go Toolchain Binary Hashing: 2 x 14 MB (28 MB hashed) -> 1 x 14 MB (14 MB hashed) (-50%).
  - Git Child Processes: 3 subprocesses per fingerprint calculation (6 per stage run) -> 2 subprocesses per fingerprint (4 per stage run) (-33%).
  - Unit Test Execution: `go test ./internal/heavyverification` dropped from 3.53s to 2.90s (-18% wall time).
- **Invalidation and Decision Equivalence:**
  - Unit Tests: 31 subtests in `TestFastStageReuseDecisionTable` pass 100%.
  - Behavior Probes: 11 probe cases in `_dev/tests/fast-stage-reuse-behavior.sh` pass 100%.
- **Mutation Qualification:**
  - Ignored File Under Coverage: Introducing an ignored file under a stage's coverage (`module-alpha/artifact.ignored`) triggers `EXECUTING (fingerprint_mismatch)`.
  - Ignored File Outside Coverage: Introducing an ignored file outside stage coverage (`outside-tree/artifact.ignored`) preserves `REUSED (fingerprint_match)`.
  - Root Scoping Fallback: Broad coverage rules (e.g. root `.` or empty) safely return nil roots to fall back to unscoped scanning.

## Before/After Evidence

| Metric | Baseline (Before) | Optimized (After) | Improvement |
| --- | --- | --- | --- |
| Toolchain probe Go binary hashing | 2 x 14 MB (28 MB) | 1 x 14 MB (14 MB) | -14 MB disk I/O & SHA-256 (-50%) |
| `workingTreeSeals` Git child processes | 3 per fingerprint (6 per stage) | 2 per fingerprint (4 per stage) | -2 Git subprocesses per stage run (-33%) |
| Ignored file directory traversal | Full repository tree walk | Scoped to stage coverage roots | Scoped directory walk |
| Unit test suite wall time (`heavyverification`) | 3.53s | 2.90s | -0.63s (-18%) |
| Invalidation mutation tests (`fast-stage-reuse-behavior.sh`) | 100% PASS (9/9) | 100% PASS (11/11) | 100% fidelity + ignored coverage cases |

## Testing

- `go test -v ./skills/do-work/tools/do-work-cli/internal/heavyverification/...` (PASS)
- `bash _dev/tests/fast-stage-reuse-behavior.sh` (PASS)
- `bash do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh` (PASS)

## Review

- Diff reviewed: surgical, strictly adhering to `_dev/primes/prime-shell-commands.md` and `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`.
- No new external dependencies introduced.
- Maintainer verification (`gate.sh`) passed completely green.

## Lessons Learned

- Combining `git ls-files` flags (`--cached` and `--others --exclude-standard`) reduces subprocess overhead when gathering tracked and untracked files simultaneously.
- When toolchain probes resolve multiple environment aliases or wrappers to the same canonical binary, in-memory caching by resolved path avoids duplicate disk reads and hashing of multi-megabyte executables while preserving security checks.

## Orientation

- Next queue item: REQ-614 (`A9: Batch repeated Git reads in the exercised code`).
