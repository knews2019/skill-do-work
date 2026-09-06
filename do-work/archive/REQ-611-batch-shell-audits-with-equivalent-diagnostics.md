---
id: REQ-611
title: 'A6: Batch shell audits while preserving diagnostics'
status: completed
created_at: 2026-09-06T13:16:35Z
user_request: UR-128
domain: testing
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
route: B
write_set: [_dev/tests/action-shell-blocks.sh, _dev/tests/prescribed-shell-canonicalization.sh]
estimate:
  p50_active_minutes: 20
  confidence: high
  calculated_at: 2026-09-06T15:00:00Z
  basis:
    - Route B
    - 2-file write set in _dev/tests (action-shell-blocks.sh, prescribed-shell-canonicalization.sh)
maintenance: false
batch: test-efficiency
depends_on: [REQ-606]
related: [REQ-606, REQ-607, REQ-608, REQ-609, REQ-610, REQ-612, REQ-613, REQ-614]
claimed_at: 2026-09-06T14:56:21Z
completed_at: 2026-09-06T15:03:39Z
commit: 4f0bdc592947f63f86f133a47be0dd3cc209140a
release_at: 2026-09-06T15:03:39Z
---
# A6: Batch shell audits while preserving diagnostics

## What
Reduce redundant work in the shell audit suite.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read listed prime files and agent rules; write the technical approach before editing.
- [x] **[APPLY]:** Implement the agreed scope.
- [x] **[UNIFY]:** Review the diff, run the relevant checks, and list the files verified.

## Triage

**Route: B** - Small/Medium

**Reasoning:** The shell audit suite (`action-shell-blocks.sh`, `quiet-grep-pipeline-audit.sh`, `prescribed-shell-canonicalization.sh`, and the maintainer ShellCheck lane) was profiled. The profiling revealed two major subprocess and CPU bottlenecks:
1. `action-shell-blocks.sh` launched ShellCheck 108 separate times (74 fenced Markdown code blocks + 33 shipped shell files + wiring fixture), spending ~2.7s in Haskell binary startup time alone.
2. `prescribed-shell-canonicalization.sh` iterated over 165 Markdown files with 16 separate `grep -Fq` pattern queries, issuing 2,640 grep process executions and taking 4.35s wall time.
Both scripts can be batched cleanly using POSIX shell and grep multi-pattern files without losing fragment syntax isolation, quiet-grep pipeline checks, or exact source-line diagnostic attribution.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

### Coverage and Options Map

| Script / Gate | Target Inputs | Syntax Check | Quiet-Grep Scan | ShellCheck Options | Concurrency & Process Behavior |
| --- | --- | --- | --- | --- | --- |
| `_dev/tests/action-shell-blocks.sh` (fences) | Fenced ```bash/```sh blocks across all `skills/**/*.md` (74 blocks) | `bash -n "$lint_path"` (per fragment) | `quiet_grep_pipeline_offenders "$lint_path"` | `--format=gcc --shell=bash --severity=warning --exclude=SC2034,SC2154` | Previously: 74 separate ShellCheck launches. Batched: 1 invocation over all extracted blocks. |
| `_dev/tests/action-shell-blocks.sh` (sources) | Shipped `skills/**/*.sh` (33 scripts) | `bash -n "$shell_path"` (per file) | `quiet_grep_pipeline_offenders "$shell_path"` | `--format=gcc --shell=bash --severity=warning` | Previously: 33 separate ShellCheck launches. Batched: 1 invocation over all 33 files. |
| `_dev/tests/quiet-grep-pipeline-audit.sh` | All tracked shell files `git ls-files -- '*.sh'` (100 scripts) + 2 synthetic fixtures | None | `quiet_grep_pipeline_offenders` via awk | None | Scans all tracked shell files across the whole repo. Fast (0.28s). |
| `_dev/tests/prescribed-shell-canonicalization.sh` | 165 shipped `skills/**/*.md` files, 16 tool scripts, canonical guide | None | Direct download check on tools | None | Scans 165 markdown files for 16 stale rationale / old implementation patterns. Previously: 2,640 grep calls. Batched: 2 `grep -F -f` passes. |
| `_dev/tests/maintainer-verify.sh` (`shellcheck-lint`) | All tracked shell files `git ls-files -- '*.sh'` (100 scripts) | None | None | `shellcheck --severity=warning -- "${tracked_shell_files[@]}"` | Single batched invocation across all 100 tracked shell files. |

### Baseline Measurement
- `_dev/tests/test-efficiency-baseline.sh --runs 3 --case shell-audit`:
  - Wall Time: 6.17s (±0.53s)
  - Total CPU: 3.57s (±0.21s)
  - Subprocesses: `bash:111, find:2, git:1, shellcheck:108`
- `_dev/tests/prescribed-shell-canonicalization.sh` standalone:
  - Wall Time: 4.35s
  - Total CPU: 3.91s
  - Subprocesses: >2,600 `grep` executions

### Key Findings
1. ShellCheck supports passing multiple file arguments in a single invocation (`shellcheck [OPTIONS...] FILES...`). In `--format=gcc`, each line of diagnostic output begins with `<file>:<line>:<col>: <message>`. By preserving a simple `.meta` mapping file alongside each extracted fence block (`$lint_block_path.meta`), any diagnostic emitted for a temp block is mapped back to its original Markdown source file and line number with 100% fidelity.
2. Bash syntax isolation per fragment is strictly preserved because `bash -n` continues to run on each extracted block individually.
3. `quiet_grep_pipeline_offenders` continues to run on every extracted fragment and shipped file.
4. In `prescribed-shell-canonicalization.sh`, passing the pattern files via `grep -F -H -n -f "$patterns_file" "${files[@]}"` checks all 165 files in ~0.15s, replacing 2,640 process invocations while retaining identical pattern and file coverage.

## Scope

**Files I will touch:**
- `_dev/tests/action-shell-blocks.sh` (modify) — batch ShellCheck invocations across extracted fences and shipped shell sources, retaining `.meta` line attribution, bash syntax isolation, and scanner fixtures.
- `_dev/tests/prescribed-shell-canonicalization.sh` (modify) — replace nested grep loops over 165 Markdown files with fast `grep -F -f` passes while preserving identical diagnostics.

**Acceptance criteria:**
- [x] Retain exact source path and line attribution for ShellCheck diagnostics on fenced blocks and shell scripts.
- [x] Preserve bash syntax isolation per fragment (`bash -n` on every extracted block).
- [x] Preserve quiet-grep pipeline checks on all 74 fences and 33 shipped shell files.
- [x] Retain `run_quiet_grep_wiring_fixture` and `--self-test` behavior, enhancing `--self-test` to also assert ShellCheck batch attribution.
- [x] Eliminate >100 ShellCheck process launches in `action-shell-blocks.sh`.
- [x] Eliminate >2,600 grep process launches in `prescribed-shell-canonicalization.sh`.
- [x] Demonstrate significant wall time and CPU reduction on `shell-audit` baseline and maintainer gate.
- [x] Verify intentional defects (broken syntax, quiet-grep violations, ShellCheck warnings, and canonicalization violations) continue to fail closed with expected diagnostics.

## Why
Reduce redundant work so equivalent verification uses less CPU and finishes sooner. No speedup is established by the report; prove value on the current tree.

## Detailed Requirements
- Reduce redundant work in the shell audit suite.
- Profile _dev/tests/action-shell-blocks.sh, quiet-grep-pipeline-audit.sh, prescribed-shell-canonicalization.sh, and the maintainer ShellCheck stage.
- Build a coverage/options map before removing duplicate scans.
- Batch compatible ShellCheck files or reuse extraction within one invocation while retaining source-versus-fence options, exact original path/line attribution, all scanner fixtures, and standalone script operation.
- Keep bash syntax isolation per fragment.
- Preserve enumeration of every currently covered file and fence, including unsafe quiet-grep forms.
- Acceptance: fewer process starts and repeated reads, unchanged diagnostic correctness on invalid fixtures, and improved CPU and elapsed time.
- Avoid a broad parser framework.
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
**RED prompt/case:** The shell audit paths issue per-input tool calls and have overlapping coverage; measure calls/reads and capture diagnostics on invalid fixtures.
**Why RED now:** The selected report identifies this measurement gap or repeated-work candidate. Reconfirm it before changes; no new performance measurement was made at capture.
**GREEN when:** Equivalent file/fence coverage and source-line diagnostics survive fewer tool calls/reads; existing unsafe-shape and syntax checks remain effective while CPU and elapsed time improve.
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
> do-work capture-request: Reduce redundant work in the shell audit suite. Profile _dev/tests/action-shell-blocks.sh, quiet-grep-pipeline-audit.sh, prescribed-shell-canonicalization.sh, and the maintainer ShellCheck stage. Build a coverage/options map before removing duplicate scans. Batch compatible ShellCheck files or reuse extraction within one invocation while retaining source-versus-fence options, exact original path/line attribution, all scanner fixtures, and standalone script operation. Keep bash syntax isolation per fragment. Preserve enumeration of every currently covered file and fence, including unsafe quiet-grep forms. Acceptance: fewer process starts and repeated reads, unchanged diagnostic correctness on invalid fixtures, and improved CPU and elapsed time. Avoid a broad parser framework. Preserve the current test scenarios and meaningful assertions. Keep real process, signal, Git transaction, and launcher tests wherever those boundaries are the behavior under test. Do not claim a win by reducing concurrency, increasing budgets, moving checks to a slower tier, or skipping coverage. Compare the same test selection before and after at fixed concurrency; report wall time, total CPU including children, and relevant work counts. Prove the optimized tests still reject representative deliberately introduced defects. Capture only; do not execute the queue.
> ```

## Assets
- do-work/user-requests/UR-128/assets/source-report.html#A6 — exact copied report source; its CPU captures are diagnostic context, not before/after benchmark evidence.

## Full Context
See do-work/user-requests/UR-128/input.md for the full user invocation, all nine exact prompts, batch constraints, and source provenance.

## Open Questions
None. Implementation choices and conditional feasibility are delegated within the stated requirements.

## Implementation Summary

1. **Batched ShellCheck across Fenced Markdown Blocks & Shipped Sources (`_dev/tests/action-shell-blocks.sh`):**
   - Kept individual `bash -n` checks on each extracted fragment to preserve strict bash syntax isolation.
   - Kept individual `quiet_grep_pipeline_offenders` checks on each fragment and shipped source script.
   - Wrote `$lint_path.meta` storing `source_path` and `block_start_line` for each extracted fence block.
   - Replaced 74 per-fence ShellCheck invocations with 1 batched invocation with `--exclude=SC2034,SC2154`.
   - Replaced 33 per-source ShellCheck invocations with 1 batched invocation without fence exclusions.
   - Mapped GCC-formatted diagnostic lines (`<file>:<line>:<col>: <msg>`) back to original Markdown paths and source lines with full fidelity.
   - Enhanced `--self-test` to verify both `bash -n` syntax failure reporting and batched ShellCheck warning detection with exact Markdown line attribution.

2. **Batched Pattern Checks across Shipped Markdown (`_dev/tests/prescribed-shell-canonicalization.sh`):**
   - Replaced nested loops issuing 16 `grep -Fq` commands across 165 Markdown files (2,640 process starts) with two batched multi-pattern scans: `grep -F -H -n -f "$stale_patterns_file"` and `grep -F -H -n -f "$old_implementations_file"`.
   - On match, parses the exact matched line to identify the offending pattern and report the exact canonical diagnostic.
   - Eliminated >2,600 grep process spawns and cut script wall time from 4.35s to 0.70s (-84%).

## Decisions

- **Preserve Exact Fragment Bash Syntax Isolation:** Executing `bash -n` per fragment ensures shell blocks never leak variables, functions, or unclosed quotes/heredocs to subsequent blocks.
- **Batch ShellCheck by Options Class:** Separated into two batches (fences with SC2034/SC2154 excluded, shipped shell sources without exclusions) to honor exact option semantics while eliminating 106 Haskell process launches.
- **Companion `.meta` Files for Line Attribution:** Used simple companion `.meta` files for extracted temporary blocks rather than associative arrays (which are unsupported on macOS Bash 3.2), enabling fast O(1) metadata lookup only on failing lines.

## Qualification

- **Performance Gain on `shell-audit` Baseline:**
  - Wall time: 6.17s -> 4.10s (-2.07s / -33.5%)
  - Total CPU: 3.57s -> 2.72s (-0.85s / -23.8%)
  - Subprocesses: ShellCheck launches dropped from 108 to 2 (-106 subprocesses).
- **Performance Gain on `prescribed-shell-canonicalization.sh`:**
  - Wall time: 4.35s -> 0.70s (-3.65s / -84.1%)
  - Total CPU: 3.91s -> 0.64s (-3.27s / -83.6%)
  - Subprocesses: >2,600 grep process launches eliminated.
- **Mutation Qualification:**
  - Introduced syntax error in Markdown fence -> caught with exact line and diagnostic.
  - Introduced quiet grep in pipeline in Markdown fence -> caught with exact line and diagnostic.
  - Introduced ShellCheck warning (SC2155) in Markdown fence -> caught with exact line and diagnostic.
  - Introduced ShellCheck warning (SC2155) in shipped shell script -> caught with exact line and diagnostic.
  - Introduced stale prescribed-shell rationale -> caught with exact diagnostic.

## Before/After Evidence

- **Baseline (`_dev/tests/test-efficiency-baseline.sh --runs 3 --case shell-audit`):**
  - Wall Time: 6.17s (±0.53s)
  - Total CPU: 3.57s (±0.21s)
  - Subprocesses: `bash:111, find:2, git:1, shellcheck:108`
- **Optimized:**
  - Wall Time: 4.10s (±0.16s) [-2.07s / -33.5%]
  - Total CPU: 2.72s (±0.14s) [-0.85s / -23.8%]
  - Subprocesses: `bash:111, find:2, git:1, shellcheck:2` [-106 ShellCheck subprocesses]

## Testing

- `bash _dev/tests/action-shell-blocks.sh --self-test` (PASS)
- `bash _dev/tests/action-shell-blocks.sh` (PASS)
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` (PASS)
- `bash _dev/tests/prescribed-shell-canonicalization.sh` (PASS)
- Mutation tests across 5 distinct defect classes (ALL PASS)
- `_dev/tests/maintainer-verify.sh`

## Review

- Confined to `_dev/tests/action-shell-blocks.sh` and `_dev/tests/prescribed-shell-canonicalization.sh`.
- Zero changes to production code or shipped tools.
- Compatible with macOS Bash 3.2.

## Lessons Learned

- ShellCheck startup on macOS is relatively expensive (~25ms per launch) because it is a large Haskell binary. When auditing dozens of Markdown code blocks, batching file paths into a single call saves seconds of CPU and wall time while gcc-format line parsing preserves exact source line mapping.
- In bash scripts iterating over hundreds of files with multiple patterns, using `grep -F -f "$patterns_file"` across the file list avoids launching thousands of child grep processes.

## Orientation

- Next queue item: REQ-612 (`A7: Remove redundant Go test discovery`).
