# Builder brief — REQ-581

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1201/worktree-agent-REQ-581-descendant-fixture`
- **Your branch (already checked out there):** `worktree-agent-REQ-581-descendant-fixture`
- **Route:** A
- **Base commit:** 961697bb

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION` — release paths owned by finalization.
- Any file outside the write set declared in the REQ below. If you need one, stop and report it in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). The orchestrator owns it and concurrent runs corrupt each other's timing budgets. Run only the focused tests named below.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/testing.md` (the REQ is `tdd: true`)

Also read every path in the request's `prime_files`, and the `lessons-<name>.md` satellite beside each prime whose Read-first or Traps entries your change touches.

## P-A-U phasing (mandatory, reported in the hand-back)

The REQ file is the orchestrator's, so report your P-A-U record under a `## P-A-U` heading in the hand-back instead of ticking boxes in the REQ:
- **[PLAN]** — brief technical approach, written before code.
- **[APPLY]** — code exactly as planned, strictly inside the declared write set.
- **[UNIFY]** — run `git diff --stat`, run the native linters (`gofmt -l .`, `go vet ./...` for Go changes, `node --check` for changed client files), verify no debug artifacts in added lines, and list each file you checked and what you checked.

## Focused tests

Every test-file invocation must finish in under 30 seconds. Use:
- `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./internal/nextselection/...`
- Then the whole module once: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./...`

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-120117/REQ-581-handback.md`** using that absolute path — it is the one main-tree path you may write, and you must never stage or commit it.

It must contain, each under its own `##` heading:
- `## Branch` — the branch name and the head commit you left on it.
- `## File manifest` — every source file created/modified/deleted with the verb, plus tests touched.
- `## P-A-U` — the three phases above.
- `## Test evidence` — every command you ran, its exit status, the RED observation (test name + failure text) and the GREEN observation.
- `## Lesson evidence` — each lesson satellite you read and any listed path that was missing.
- `## Decisions` — significant choices as `D-NN`, each with reasoning. Mark a reversible low-reach choice DECIDE & STATE; mark an irreversible, taste-dependent or contestable one ESCALATE and add `Value:` and `Risk:` lines.
- `## Discovered Tasks` — out-of-scope findings, each stamped with one of exactly these impact tokens: `impact-critical`, `impact-user-visible`, `impact-rule-change`, `impact-negligible`. Do not invent a token outside that set and do not fix the items inline.
- `## Integration seams` — any exact line that belongs in a file outside your write set, with where it goes. The orchestrator applies it.

Work the mutation first, exactly as the request's Builder Guidance says. Reduce `terminateOwnedProcessGroup` and `cleanupReapedProcessGroup` in `internal/nextselection/blocked_probe_unix.go` to no-op bodies, keep that mutation applied while you rewrite the fixture, and revert it only once all three tests fail on their own descendant assertion rather than on a timeout bound. That file is NOT in your write set — the mutation is a scratch experiment and must be reverted with `git checkout --` before you commit. Prove it is reverted: `git status --porcelain` must be empty of it, and your hand-back must say so.

Report the elapsed time of the three tests both mutated and unmutated. The request forbids raising the package's unmutated runtime out of its current range, and another request (REQ-574) just spent effort bringing this module's test files under the 30-second per-file budget, so a fixture that adds seconds is a regression even if it is correct.

Other builders are running Go tests on this machine at the same time. If a run comes in slower than you expect, re-run before concluding anything about timing.

---

# The request

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

---

## Triage

**Route: A** - Simple

**Reasoning:** One file, named. The reviewer of REQ-506 already ran the mutation and reported the numbers, so the RED case, the fixture shape and the exact defect are all in the request. `effort_estimate: effort-substantive` reflects the care the fixture needs, not unknown territory: nothing about the location or the pattern needs discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route A: direct to builder

*Skipped by work action*
