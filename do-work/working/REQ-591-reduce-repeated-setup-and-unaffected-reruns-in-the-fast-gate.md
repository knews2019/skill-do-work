---
id: REQ-591
title: 'Reduce repeated setup and unaffected reruns in the fast gate'
status: claimed
created_at: 2026-09-05T19:43:25Z
user_request: UR-127
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: []
related: [REQ-574]
claimed_at: 2026-09-05T20:04:13Z
---

# Reduce Repeated Setup and Unaffected Reruns in the Fast Gate

## What

Make this repository's routine verification faster and cheaper without losing the failure coverage that protects queue state and completed work. Start by removing repeated fixture/build setup, then reuse verification results only when the complete inputs of the relevant checks are unchanged.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read listed prime files and agent rules. Record the measured cost, selected scope, and coverage-preserving approach before implementation.
- [ ] **[APPLY]:** Implement the planned changes within the declared scope.
- [ ] **[UNIFY]:** Review the exact diff, run relevant checks, and record measured performance and retained failure coverage.

## Why

The user wants the queue to drain faster while maintaining good enough quality. The CPU investigation found several independent test pools competing for one machine. Avoiding repeated work lowers the verification cost per request; increasing concurrency alone does not establish an improvement.

## Context and Prior Work

- REQ-574 (reducing CLI fixture setup costs) is already completed, merged at `50569e88c8f1f5234cbdfaf0efaede671d72b13c`. It reused initialized Git repositories in finalization/publication tests and avoided repeated launcher/toolchain work in inventory cases. Its recorded whole-module comparison was 65s to 61s with 772 tests and every assertion retained. Build on that work instead of repeating it. Source: `do-work/archive/UR-115/REQ-574-bring-do-work-cli-test-files-under-the-30s-budget.md`.
- `_dev/tests/session-start-hook-behavior.sh` still copies the complete CLI module separately for each scenario. This is a concrete setup-reuse candidate, not a requirement to bypass the launcher path under test.
- `_dev/tests/maintainer-verify.sh` selects the ordinary board tests and the complete CLI module with `-short`; `_dev/tests/run-go-tests-with-budget.sh` forces `-count=1`. Existing fast/heavy separation is already implemented.
- Existing heavy-lane fingerprints and green-gate evidence provide starting points. Inspect their actual invalidation and caller contracts before extending them; a prior green at HEAD is not proof that current dirty or external inputs match.
- `do-work/test-durations.tsv` provides per-file observations. Summed file times are not complete-gate wall time or CPU time, and loaded-window observations cannot establish an uncontended speedup.
- Fold-first review found no eligible pending request in any UR with this root cause. The queued helper deduplication and shell-guide requests do not close the repeated fast-gate setup/rerun problem.

## Detailed Requirements

1. Reduce repeated setup in measured hot paths. Reuse immutable fixture material and the current build where that preserves the boundary being tested. Keep writable test state isolated. Retain explicit coverage of the real launcher, build failures, missing tools and installed layout where those are the behavior under test.
2. After the setup improvement, avoid repeating verification whose complete relevant inputs have not changed. Use the existing verification/evidence mechanisms where practical. Relevant inputs include source, transitive dependencies, fixtures, test/gate scripts, configuration and effective toolchain/runtime inputs, including uncommitted changes. Unknown impact, incomplete evidence or an unverifiable input must select the broader verification rather than produce a false green.
3. Make selective results reviewable: record which checks executed, which reused evidence and why. Preserve failure and interruption exit statuses. A skipped, failed or incomplete run cannot supply successful evidence. A board-only change may reuse unrelated CLI evidence only after proving that the relevant CLI inputs are unchanged.
4. Keep useful failure coverage. Preserve real checks for rollback, interrupted recovery, locking, concurrent writers and process cleanup. Do not silently remove assertions, replace a real boundary with a mock, or move essential coverage behind the heavy tier to improve the number. Any proposed consolidation must name the failure it still catches and its retained test.
5. Measure before and after on comparable revisions with a fixed worker limit, recorded toolchain/cache state and no competing expensive gate or synthetic load. Record complete-gate wall time and process-tree CPU cost (or explicitly identify unavailable measurements), alongside the existing per-file durations. Use the smallest bounded comparison that establishes a repeatable improvement; do not benchmark by repeatedly saturating the machine.

## Acceptance Criteria

- Repeated setup is measurably reduced in at least one identified hot path, with its behavior assertions retained and writable fixtures independent.
- A matching-input verification case reuses a recorded success; changing a relevant source/dependency, fixture, gate script or runtime input invalidates it. Uncertain impact and dirty relevant inputs cannot reuse stale evidence. Tests cover these positive and negative cases.
- The routine gate's end-to-end duration improves under the recorded comparable conditions; the evidence separates setup savings, avoided execution and contention. A one-off noisy sample, a renamed/split test file, a raised timeout or greater concurrency alone is not acceptance.
- The existing fast-tier per-file budget and heavy-tier policy remain satisfied. Required correctness checks pass; retained rollback/recovery/locking/cleanup tests remain capable of detecting their named failures. Reused results are visibly distinguishable from freshly executed ones.

## Constraints

- Repository scope is `skill-do-work2`. The other projects seen in process metadata are diagnostic context, not authorization to inspect or change their test configurations.
- No machine-wide scheduler, background reaper or watchdog in this request. Shared machine concurrency policy is adjacent work; use a controlled worker limit for the comparison here.
- Preserve existing heavy-test permission/defer behavior. This capture is not authorization to run an otherwise permission-gated heavy suite.
- Synthetic load is unnecessary for this optimization. Any deliberately backgrounded process still needs its own lifetime bound under the existing testing rule.
- Do not optimize to a fabricated percentage or treat the 30-second per-file ceiling as a whole-gate speed target.

## Builder Guidance

Implement and verify setup reuse before selective rerun behavior. Choose the smallest existing code path that supports conservative invalidation; do not introduce a second verification platform. Keep one coherent request, with small verified increments. Exact files and test cases are resolved during planning rather than guessed into `write_set` at capture.

## Red-Green Proof

**RED prompt/case:** Exercise the normal fast gate's current repetition on unchanged relevant inputs, then a fixture where only unrelated board inputs change. Record which expensive checks rerun. In the existing harness, add a test of the intended reuse decision plus counter-cases for relevant source, fixture, gate script, runtime and uncommitted-input changes.
**Why RED now:** Routine verification still executes broad uncached selections and repeats fixture setup; the proposed fast-gate input-aware reuse is not established by the existing heavy-lane or revision-only evidence rules.
**GREEN when:** The matching/unrelated-input cases reuse only valid success, every relevant or uncertain change executes the needed checks, retained failure scenarios still fail when broken, and comparable before/after evidence shows lower routine gate cost.
**Validation:** Inferred during capture from the recommendation the user asked to capture. The user approved the objective; exact regression fixtures and a defensible performance baseline are builder-resolved details.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — index cost 10543 tokens; exceeds the 2000-token budget and is `slugged: partial`, so targeted selection is not eligible. Matches `fixture-cost-is-subprocess-spawning`, `background-worker-self-bound`, `smoke-vs-characterization` and `shipped-module-test-self-containment`.
- `_dev/primes/lessons-shell-commands.md` — index cost 3385 tokens; exceeds the budget and is `slugged: partial`. Matches the shell/launcher fixture and migration-parity surface.

## Open Questions

None about user intent. The builder establishes the exact measured scope and comparison conditions.

## Full Context

See `do-work/user-requests/UR-127/input.md` for the capture instruction and prior conversation context.
