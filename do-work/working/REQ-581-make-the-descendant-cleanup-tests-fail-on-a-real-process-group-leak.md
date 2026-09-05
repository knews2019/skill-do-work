---
id: REQ-581
title: '[impact-rule-change] Make the descendant-cleanup tests fail on a real process-group leak'
status: claimed
created_at: 2026-09-05T01:30:57Z
user_request: UR-119
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
write_set: [skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go]
tdd: true
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
claimed_at: 2026-09-05T12:40:39Z
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-09-05T12:41:13Z
  basis:
    - Route A
    - 1-file write set
    - 5 acceptance criteria
route: A
dispatch_at: 2026-09-05T12:41:13Z
builder_handback_at: 2026-09-05T12:59:09Z
integration_at: 2026-09-05T12:59:09Z
review_at: 2026-09-05T14:23:29Z
heavy_verified_at: 2026-09-05T14:23:29Z
heavy_verified_revision: 7b2673b690a671ccb360c26b0c19c56ecc7356b5
commit: 92339213
---

# Make the Descendant-Cleanup Tests Fail on a Real Process-Group Leak

## What
Three tests in `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` claim to prove that a probe's descendant process group is terminated. Two of them pass on a tree with a genuine leak, and the third fails on a timeout rather than on its descendant assertion. Rewrite the descendant fixture so a leak fails the assertion that names it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Builder read the prime, the crew rules and the four matching lesson families, then settled a five-step plan: reproduce the reviewer's numbers under the leak mutation, rewrite the fixture on three points while the mutation stayed applied, raise the interruption test's return bound past the new leader hold, correct the budget comment, then revert and prove the runtime did not move. Recorded under `## P-A-U` in `do-work/runs/work-2026-09-05-120117/REQ-581-handback.md`.
- [x] **[APPLY]:** One commit on the builder branch (`f5b3faf4`), 59 insertions and 34 deletions in the single declared file.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `gofmt -l .` empty; `go vet ./...` clean; the scratch mutation reverted with `git checkout --` and proved absent from `git status --porcelain` and from the diff before committing.

## Why
A control that cannot fail is worse than no control: it reports the process-cleanup path as proven every time the suite runs, so the next change to that path ships unchecked.

## Context
The three tests are `TestBlockedProbeTimeoutKillsDescendantGroup`, `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits` and `TestBlockedProbeInterruptionIsTypedAndReapsDescendants`. All three end in `waitForDescendantToDisappear`, which polls `kill(pid, 0)` against `descendantReapBudget`.

The structural reason they cannot fail: a surviving descendant inherits the parent-owned diagnostic pipe, so `<-diagnosticDone` and `<-done` hold the runner until that descendant exits. By the time the poll loop runs, the descendant is always already gone. What the loop actually measures is how long init takes to reap a zombie after the runner returns — which is why it flaked at 1.13-1.95s against the earlier 2-second budget under load. A real leak shows up only as the test taking 30 seconds instead of 3.

`descendantReapBudget` currently carries a comment claiming "it proves the descendant does not survive, and it is not a latency assertion." That comment states the opposite of what the code does and needs to change with the fixture.

This defect is pre-existing. It was surfaced by the independent review of REQ-506 (running the evidence gates from advance) but was not introduced by it, so this is not a regression fix for that request.

REQ-574 (bringing do-work-cli test files under a 30s per-file budget) is in flight over the same package's runtime. That is a reason not to trade this fix for a slower suite, not a dependency.

## Red-Green Proof
**RED prompt/case:** Reduce `terminateOwnedProcessGroup` and `cleanupReapedProcessGroup` in `internal/nextselection/blocked_probe_unix.go` to no-op bodies — a genuine process-group leak — then run `go test ./internal/nextselection`.
**Why RED now:** On today's tests that mutated tree stays green where it matters. `TestBlockedProbeTimeoutKillsDescendantGroup` and `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits` both still pass, taking 30.01s and 31.35s instead of 2.90s and 2.01s. `TestBlockedProbeInterruptionIsTypedAndReapsDescendants` does fail, but on its own 5-second "interrupted probe did not return" bound, never reaching its descendant assertion. A leak is therefore visible only as elapsed time, which nothing asserts on.
**GREEN when:** With the same two functions reduced to no-ops, each of the three tests fails inside its own budget on the descendant assertion, and the failure message names the surviving process id. No test in the package exits on a timeout bound instead. With the two functions restored, all three pass in their present single-digit-second range.
**Validation:** User confirmed — the fixture shape below is the reviewer's own suggestion, recorded as the RED case at the maintainer's instruction.

## Detailed Requirements
- Replace the descendant fixture in all three tests with one that closes its inherited stdout and stderr before sleeping, so a surviving descendant cannot hold the parent-owned diagnostic pipe open and cannot delay the runner's return.
- Make the descendant outlive the assertion budget, so a leaked process is still alive when the poll loop looks for it.
- Assert on the descendant's liveness, not on how fast it disappears. Under the no-op mutation the failure must be the descendant assertion.
- `TestBlockedProbeInterruptionIsTypedAndReapsDescendants` must reach and fail its descendant assertion under the mutation instead of stopping at its 5-second "interrupted probe did not return" bound.
- Correct the `descendantReapBudget` comment to describe what the rewritten loop measures.

## Constraints
- Test-only change. `blocked_probe_unix.go` is correct; what is missing is any test able to detect its absence.
- Do not raise the package's unmutated runtime out of its current range.

## Builder Guidance
The reviewer already ran the mutation and reported the numbers. Reproduce that mutation as the first step, keep it applied while writing the fixture, and only revert it once all three tests fail on the descendant assertion.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 7450 tokens, over the 2000-token budget in `actions/capture-reference.md` → Required Lessons Budget Contract, and `slugged: partial` so no targeted entry is legal. Matches on families `reaped-by-its-own-parent`, `interruptible-blocking-io`, `smoke-vs-characterization` and `silent-skip-reads-as-red`.

## Full Context
See `do-work/user-requests/UR-119/input.md` for complete verbatim input.

*Source: independent review of REQ-506 (running the evidence gates from advance), work run `do-work/runs/work-2026-09-05-003420/`.*

---

## Triage

**Route: A** - Simple

**Reasoning:** One file, named. The reviewer of REQ-506 already ran the mutation and reported the numbers, so the RED case, the fixture shape and the exact defect are all in the request. `effort_estimate: effort-substantive` reflects the care the fixture needs, not unknown territory: nothing about the location or the pattern needs discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route A: direct to builder

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` (modified)

**What was done:** The descendant fixture now releases its inherited stdout and stderr before sleeping, so a surviving descendant cannot hold the runner's diagnostic pipe open and delay its return — which is the structural reason a leak used to be invisible. The descendant outlives the reap budget, so a leaked process is still alive when the poll looks for it. The probe leader holds itself for a bounded interval instead of waiting on the descendant, which keeps the timeout and interrupt branches firing while capping how long a broken runner can block, and the interruption test's return bound was raised past that hold so it reaches its own descendant assertion. Three copies of pid-file reading collapsed into one helper that tolerates a partially written file and registers its own cleanup kill. The budget comment no longer claims the loop proves something it did not. Merge range `1a06c3bc..92339213`; builder branch head `f5b3faf4`.

## Qualification

**Passed.** Read from the merge range `1a06c3bc..92339213`.

- Test-only change, one file, exactly the declared write set. `blocked_probe_unix.go` — the code under test — is untouched; the builder mutated it as a scratch experiment and proved the revert with both `git status --porcelain` and an empty diff before committing.
- The fixture change is the real fix, not a budget increase. Releasing the inherited descriptors is what makes a leaked descendant observable at all; making it outlive the budget is what makes the poll mean something. Neither is a timing tweak.
- The interruption test's raised bound is not a weakening: it is what lets that test reach the assertion it was written for instead of dying at an unrelated bound, which is exactly the defect the request names.
- Three near-duplicate pid-reading blocks, two of them with unchecked conversions, became one helper. That is a reduction, not an addition.

Requirements traced: all three tests fail under the leak mutation, each inside its own budget, each on the descendant assertion, each naming the surviving process id; no test exits on a timeout bound; the budget comment now describes what the loop measures.

*Checked by work action*

## Testing

**Red-green validation** (traced to `## Red-Green Proof`), four measured states on one machine in one session:
- Baseline, old fixture, unmutated: all three pass, 1.50s / 0.11s / 0.63s, package 2.624s.
- RED reproduction, old fixture, leak mutation applied: two tests still PASS at 30.01s each, the third fails at `blocked_probe_test.go:114: interrupted probe did not return` on its 5s bound. This reproduces the reported defect exactly — a leak visible only as elapsed time, which nothing asserts on.
- RED, new fixture, same mutation: all three FAIL inside their own budget on `descendant NNNNN survived 10s`, each naming the surviving process id. The typed-interruption assertion still passes under the mutation, so the test now separates a broken teardown from a broken result.
- GREEN, new fixture, mutation reverted: all three pass at 1.50s / 0.11s / 0.65s, package 2.610s against the 2.624s baseline. The runtime constraint holds; the change is timing-neutral.

**Flake check:** `-count=3` over the three tests, no failure.


## Review

**Overall: 95%** | 2026-09-05T13:05:35Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Reviewer's own mutation run (independent of the builder's).** The module was copied to a scratch directory; the main checkout was never mutated. `cleanupReapedProcessGroup` and `terminateOwnedProcessGroup` in `internal/nextselection/blocked_probe_unix.go` were reduced to no-op bodies there, and the whole `./internal/nextselection` package was run:

```
--- FAIL: TestBlockedProbeTimeoutKillsDescendantGroup (14.02s)
    blocked_probe_test.go:67: descendant 65381 survived 10s
--- FAIL: TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits (10.03s)
    blocked_probe_test.go:78: descendant 65867 survived 10s
--- FAIL: TestBlockedProbeInterruptionIsTypedAndReapsDescendants (14.02s)
    blocked_probe_test.go:115: descendant 66223 survived 10s
```

All three fail on the descendant assertion, each inside its own budget, each naming the surviving process id. No test in the package exits on a timeout bound. Unmutated, the same scratch tree passes at 1.51s / 0.11s / 0.65s. The fix is real.

**Runtime constraint measured.** Two interleaved full-package rounds in the same scratch, old fixture against new: old 2.474s / 2.470s, new 2.491s / 2.491s. About +20ms (0.8%), inside noise. The unmutated runtime is not raised.

**Flake check under load.** `-count=5` over the three tests with six busy CPU loops running: 15/15 pass, per-test timings identical to the unloaded run (1.51s / 0.11s / 0.65s). No time-dependence observed.

**The raised bound is not a weakening.** `interruptedProbeReturnBudget` 5s → 15s bounds "the interrupted runner returns rather than hanging forever". The healthy return is 0.65s, so the bound is still ~23x the real value and still catches a hang. Under the leak mutation the leader holds 4s, so any bound at or under 5s fires there instead of at the descendant assertion — which is exactly the defect the request names. The only cost is that a genuine hang is reported after 15s instead of 5s.

**Restatement Sweep.** The diff redefines what `descendantReapBudget` measures. Swept `descendantReapBudget`, `waitForDescendantToDisappear`, "latency assertion", "reaping-latency" and `blocked_probe_test` across the repository. Every other occurrence is a historical record in `do-work/` (UR-119 input, REQ-506 archive, run briefs and hand-backs) describing the defect as it was found; those stay accurate as history. No shipped file, prime, lessons satellite, or changelog entry restates the old meaning. `prime-do-work-cli.md` line 76 (`[family: reaped-by-its-own-parent]`) is about the implementation contract and is unaffected. Sweep clean, no stale restatement.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- The SIGKILL escalation inside `cleanupReapedProcessGroup` (`blocked_probe_unix.go:92-94`) is still unreachable by any test. Reviewer mutation: delete only those three lines and the 100ms sleep, leaving the initial SIGTERM. `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits` stays green in 0.319s, because its descendant dies on the plain SIGTERM. The equivalent mutation on `terminateOwnedProcessGroup` is caught — T1 and T3 both fail on the descendant assertion — so only this one branch is uncovered. A `trap '' TERM` in that test's descendant, as T1 already carries, would close it for roughly 100ms of runtime. Narrower than this request's stated RED, so not a missed requirement. — impact-negligible → report only
- `waitForDescendantPid`'s `t.Cleanup` SIGKILL (`blocked_probe_test.go:144`) kills the recorded pid, which on this machine's `/bin/sh` is the subshell, not the `sleep 30` it forks. Verified twice: `ps` shows the subshell and a separate `sleep 30` child for both fixture shapes, and after the mutated run pid 64147 (`sleep 30`, ppid 1) was still alive on the host. So a failing run does leave one orphaned `sleep` per test for the remainder of its 30 seconds. The hand-back's D-05 claims the cleanup prevents this; it prevents most of it. Bounded and self-clearing, no effect on any assertion. — impact-negligible → report only
- `waitForDescendantPid` does not detect the partial-write case the hand-back's [APPLY] says it tolerates. A truncated file whose content is a valid integer prefix (`671` of `67123`) parses and passes the `pid > 0` guard, yielding a wrong pid. Requiring a trailing newline would deliver the claimed property. Practically unreachable — `echo $! > file` is one small `write()` — and the helper is still a real improvement over the two ignored `strconv.Atoi` errors it replaced. The overstatement is in the hand-back prose, not in the shipped comment. — impact-negligible → report only
- The timeout and interrupt branches now depend on the leader outliving the runner's own deadline by ~3s (`probeLeaderHoldSeconds = 4` against a 1s timeout in T1 and a ~0.1s interrupt in T3), where the old `wait` fixture gave ~29s. If a runner deadline ever fired more than 3s late the test would fail on its status assertion rather than flake silently, and 15 loaded repeats showed no drift at all, so this is a recorded margin change and not an observed defect. — impact-negligible → report only

**Nit findings:**
- `runBlockedProbeFixture` (`blocked_probe_test.go:167`) is a one-line passthrough to `RunBlockedProbeAtRoot` with a single caller. Pre-existing, already recorded by the builder under Discovered Tasks. — impact-negligible → report only

**Acceptance:** Pass — reviewer's own leak mutation makes all three tests fail on the descendant assertion inside their own budgets; unmutated they pass with the package runtime unchanged.
**Suggested testing:** 3 items — (1) run the three tests on a Linux host where `/bin/sh` is `dash`: `dash` exec-optimizes the last command in a subshell, so the recorded pid is the `sleep` itself rather than a subshell, and `exec 1>&- 2>&-` closure behavior should be confirmed there; (2) add `trap '' TERM` to `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits`'s descendant to cover the `cleanupReapedProcessGroup` SIGKILL escalation, then re-run the reviewer's narrow mutation to confirm it now fails; (3) exercise the three tests on a CI box under heavier load than this host to re-measure the 3s slack between `probeLeaderHoldSeconds` and the runner's own deadlines
**Follow-ups created:** None (5 findings report only)

*Reviewed by review-work action*

## Heavy Verification Result

- **Target revision:** 92339213
- **Execution revision:** 7b2673b690a671ccb360c26b0c19c56ecc7356b5
- **Run at:** 2026-09-05T14:23:29Z, from a detached worktree (the shared main tree carried other sessions' uncommitted work, which a lane result must not be attributed to)

| Lane | Exit | Wall | Disposition |
| --- | --- | --- | --- |
| `do-work-cli-integrations` | 0 | 78s | executed |
| `staged-skills` | 0 | 44s | executed |
| `updater` | 0 | 71s | executed |
| `installer` | 0 | 33s | executed |

Every lane this request selected was present in the run, exited 0, and none was skipped.

## Lessons Learned

**What worked:** Keeping the leak mutation applied while rewriting the fixture, so every edit was measured against the failure it had to detect rather than against a suite that was already green. The reviewer then reproduced the whole thing independently in a scratch copy, which is the only way a claim about a control that could not fail becomes evidence.

**What didn't:** The original fixture could not fail for a structural reason, not a sloppy one. A surviving descendant inherited the parent's diagnostic pipe, so the runner blocked until that descendant exited, and by the time the poll ran the process was always already gone. What the loop actually measured was how long init took to reap a zombie — and its comment claimed the opposite. A leak showed up only as the suite taking thirty seconds instead of three, which nothing asserted on.

**Worth knowing:** A test that has never been observed failing is not a control. This one had passed on every run for as long as it existed, on a tree with a genuine process-group leak. When a test guards a cleanup path, the cheapest honest check is to break the path on purpose and confirm the test notices — and to check *which* assertion fires, because a failure on an unrelated timeout bound is not the test doing its job.

## Orientation

The do-work-cli probe suite can now detect a process-group leak instead of reporting the cleanup path as proven on every run. Lives in the do-work-cli selection subsystem (`skills/do-work/tools/do-work-cli/prime-do-work-cli.md`), in the blocked-probe tests. Test-only: the cleanup code itself was already correct, and what was missing was any way to tell if it stopped being correct. Unmutated runtime is unchanged at 2.610s against a 2.624s baseline.
