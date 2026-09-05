---
id: REQ-581
title: '[impact-rule-change] Make the descendant-cleanup tests fail on a real process-group leak'
status: pending
created_at: 2026-09-05T01:30:57Z
user_request: UR-119
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
write_set: [skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go]
tdd: true
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
---

# Make the Descendant-Cleanup Tests Fail on a Real Process-Group Leak

## What
Three tests in `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` claim to prove that a probe's descendant process group is terminated. Two of them pass on a tree with a genuine leak, and the third fails on a timeout rather than on its descendant assertion. Rewrite the descendant fixture so a leak fails the assertion that names it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
