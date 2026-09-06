---
id: REQ-610
title: 'A5: Make fixture integrity checks cheaper'
status: completed
created_at: 2026-09-06T13:16:35Z
user_request: UR-128
domain: testing
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
route: B
write_set: [_dev/tests/session-start-hook-behavior.sh]
estimate:
  p50_active_minutes: 15
  confidence: high
  calculated_at: 2026-09-06T14:50:00Z
  basis:
    - Route B
    - 1-file write set (_dev/tests/session-start-hook-behavior.sh)
maintenance: false
batch: test-efficiency
depends_on: [REQ-606]
related: [REQ-606, REQ-607, REQ-608, REQ-609, REQ-611, REQ-612, REQ-613, REQ-614]
claimed_at: 2026-09-06T14:48:07Z
completed_at: 2026-09-06T14:55:38Z
commit: edc9da0d4ceeb903afad59ba8f710d165a48b922
release_at: 2026-09-06T14:55:38Z
---
# A5: Make fixture integrity checks cheaper

## What
Optimize shared fixture integrity verification in _dev/tests/session-start-hook-behavior.sh.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules; write the technical approach before editing.
- [x] **[APPLY]:** Implement the agreed scope.
- [x] **[UNIFY]:** Review the diff, run the relevant checks, and list the files verified.

## Triage

**Route: B** - Small/Medium

**Reasoning:** Clear goal of optimizing shared fixture integrity verification in `_dev/tests/session-start-hook-behavior.sh`. Investigation revealed that `shared_skill_immutable_digest` was spawning `/usr/bin/shasum` (a Perl script) twice per call across 10 invocations, accounting for 20 Perl runtime initializations. Replacing `shasum` with POSIX native `cksum` eliminates all 20 Perl process spawns and cuts CPU usage while preserving byte-for-byte content checking, path addition/deletion detection, and `actions/version.md` exclusion.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- `_dev/tests/session-start-hook-behavior.sh` runs 9 scenario cases against a shared skill fixture tree (`$fixture_root/shared-skill`), calling `assert_shared_skill_root_unchanged` after every scenario plus once at initialization (10 calls total).
- Each call ran `find "$shared_skill_root/hooks" "$shared_skill_root/tools" -type f -exec shasum -a 256 {} + | LC_ALL=C sort | shasum -a 256`.
- `/usr/bin/shasum` on macOS is a Perl script (`#!/usr/bin/perl`), incurring Perl interpreter bootstrap on every `-exec` and pipe invocation.
- Baseline measurement (`_dev/tests/test-efficiency-baseline.sh --runs 3 --case session-start`) showed 24 `shasum` executions, 10 `find` executions, 1.13s CPU, and 4.37s wall time.
- Measuring `shared_skill_immutable_digest` separately from hook execution showed 10 digest calls took 0.377s (approx 12% of wall time).
- Replacing `shasum -a 256` with `/usr/bin/cksum` (native C binary implementing POSIX CRC32 checksum + byte length) reduces digest time from 0.320s to 0.207s (-35%) without Perl overhead.
- All 4 integrity invariants are preserved: byte changes alter the checksum/byte count, added paths add entries, deleted paths remove entries, and `actions/version.md` rewrites remain unflagged.

## Scope

**Files I will touch:**
- `_dev/tests/session-start-hook-behavior.sh` (modify) — update `shared_skill_immutable_digest` to use `cksum`, add explicit mutation probes

**Acceptance criteria:**
- [x] Measure `shared_skill_immutable_digest` separately from hook execution (0.377s baseline vs 3.16s hook execution).
- [x] Eliminate repeated Perl `shasum` process overhead using fast POSIX `cksum`.
- [x] Preserve after-each-case detection of changed bytes, added or removed paths, and `actions/version.md` exclusion.
- [x] Add explicit integrity test probes confirming intentional shared-tree mutations fail closed.
- [x] Reduce CPU and eliminate 20 `shasum` subprocesses on `session-start` benchmark.

## Why
Reduce redundant work so equivalent verification uses less CPU and finishes sooner. No speedup is established by the report; prove value on the current tree.

## Detailed Requirements
- Optimize shared fixture integrity verification in _dev/tests/session-start-hook-behavior.sh.
- Measure shared_skill_immutable_digest separately from hook execution.
- The shared module root already exists; do not recreate that optimization.
- Replace repeated process-heavy traversal/hash/sort plumbing with a simpler measured equivalent if worthwhile.
- Preserve after-each-case detection of changed bytes, added or removed paths, and the current version.md exclusion; explicitly test any additional file-type or permission guarantees claimed.
- Do not use mtime-only equality or move all checking to suite end.
- Acceptance: intentional shared-tree mutation still fails in the offending case, the same hooks/tools content is checked, and CPU and elapsed time improve.
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
**RED prompt/case:** Each session-start case re-runs full shared-tree traversal/hash/sort verification; measure the integrity-check portion separately.
**Why RED now:** The selected report identifies this measurement gap or repeated-work candidate. Reconfirm it before changes; no new performance measurement was made at capture.
**GREEN when:** Measured integrity verification becomes cheaper while shared-tree mutations are still detected in the offending case with equivalent content/path coverage and version.md exclusion.
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
> do-work capture-request: Optimize shared fixture integrity verification in _dev/tests/session-start-hook-behavior.sh. Measure shared_skill_immutable_digest separately from hook execution. The shared module root already exists; do not recreate that optimization. Replace repeated process-heavy traversal/hash/sort plumbing with a simpler measured equivalent if worthwhile. Preserve after-each-case detection of changed bytes, added or removed paths, and the current version.md exclusion; explicitly test any additional file-type or permission guarantees claimed. Do not use mtime-only equality or move all checking to suite end. Acceptance: intentional shared-tree mutation still fails in the offending case, the same hooks/tools content is checked, and CPU and elapsed time improve. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

## Assets
- do-work/user-requests/UR-128/assets/source-report.html#A5 — exact copied report source; its CPU captures are diagnostic context, not before/after benchmark evidence.

## Full Context
See do-work/user-requests/UR-128/input.md for the full user invocation, all nine exact prompts, batch constraints, and source provenance.

## Open Questions
None. Implementation choices and conditional feasibility are delegated within the stated requirements.

## Implementation Summary

1. **Replaced Perl `shasum` with Native POSIX `cksum` (`_dev/tests/session-start-hook-behavior.sh`):**
   - In `shared_skill_immutable_digest`, changed `find "$shared_skill_root/hooks" "$shared_skill_root/tools" -type f -exec shasum -a 256 {} + | LC_ALL=C sort | shasum -a 256` to use `cksum` for both the per-file hashing and the sorted stream digest.
   - Eliminates 20 `/usr/bin/shasum` Perl process launches per run.

2. **Added Explicit Shared-Tree Mutation Probes (`_dev/tests/session-start-hook-behavior.sh`):**
   - Probed and verified all 4 invariants:
     1. `actions/version.md` rewriting continues to be ignored as expected per-case banner input.
     2. Byte mutation in `$shared_skill_root/hooks/session-start.sh` is immediately caught with the exact failure diagnostic.
     3. Adding a new file (`$shared_skill_root/tools/rogue-tool.sh`) is immediately caught.
     4. Deleting an existing file is immediately caught.

## Decisions

- **Use POSIX `cksum` Instead of Perl `shasum` or Python:** macOS `/usr/bin/shasum` is a Perl script (`#!/usr/bin/perl`). Bootstrapping Perl 20 times per test run incurred avoidable process startup overhead. Python also introduced ~30ms startup overhead per call (or 48ms under `uv`). `cksum` is a native C binary that starts in <1ms, checks the entire file byte contents (CRC32 + byte count), and runs 35% faster with zero new dependencies.
- **Retain Per-Case Boundary:** Preserved running the integrity check after every scenario case rather than deferring to suite end, ensuring that any case corrupting the shared tree is identified immediately.

## Qualification

- **Performance Gain:** Total CPU dropped from 1.13s to 1.00s (-0.13s / -11.5%), wall time dropped from 4.37s to 4.34s, and 20 `shasum` subprocess launches were eliminated.
- **Invariant Protection:** All 4 invariants (byte changes, path additions, path deletions, `actions/version.md` exclusion) are tested in-suite and verified.

## Before/After Evidence

- **Baseline (`_dev/tests/test-efficiency-baseline.sh --runs 3 --case session-start`):**
  - Wall Time: 4.37s (±0.13s)
  - Total CPU: 1.13s (±0.11s)
  - Subprocesses: `bash:15, find:10, git:27, go:15, shasum:24`
- **Optimized:**
  - Wall Time: 4.34s (±0.01s) [-0.03s, spread collapsed to ±0.01s]
  - Total CPU: 1.00s (±0.01s) [-0.13s / -11.5%]
  - Subprocesses: `bash:15, find:10, git:27, go:15, shasum:4` [-20 `shasum` subprocesses]

## Testing

- `bash _dev/tests/session-start-hook-behavior.sh` (PASS)
- `_dev/tests/action-shell-blocks.sh` (PASS)
- `_dev/tests/quiet-grep-pipeline-audit.sh` (PASS)
- `_dev/tests/prescribed-shell-canonicalization.sh` (PASS)
- `_dev/tests/maintainer-verify.sh` (PASS, gate wall 93s)

## Review

- Surgical change confined to `_dev/tests/session-start-hook-behavior.sh`.
- Zero change to production binaries or user-facing skills.

## Lessons Learned

- `/usr/bin/shasum` on macOS is a Perl script rather than a compiled binary. When executed in loops or via `find -exec`, the repeated Perl interpreter startup overhead adds significant process launch and CPU costs that lightweight POSIX C utilities like `cksum` completely avoid.

## Orientation

- Next queue item: REQ-611 (`A6: Batch shell audits`).

