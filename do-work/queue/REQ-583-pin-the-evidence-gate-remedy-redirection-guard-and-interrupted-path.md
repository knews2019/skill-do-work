---
id: REQ-583
title: 'Addendum: pin the evidence-gate remedy redirection, layered guard and interrupted focused-test code'
status: pending
created_at: 2026-09-05T01:30:57Z
user_request: UR-119
addendum_to: REQ-506
depends_on: [REQ-506]
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
write_set: [skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
---

# Addendum: Pin the Evidence-Gate Remedy Redirection, Layered Guard and Interrupted Focused-Test Code

## What
Three pieces of behaviour in `internal/lifecycleadvance/evidence_gates.go` have no test that fails when they are removed. Delete the remedy-redirection call sites, or either half of the layered guard in `focusedGateState`, and `go test ./internal/lifecycleadvance` stays fully green. Add the tests that hold each one in place, plus a public case for the interrupted focused test whose finding code changed.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
One root cause: new behaviour landed without a lock. Each of the three was delivered on purpose — the remedy redirection to answer a prior review finding, the layered guard because the remediation plan explicitly asked for it — and nothing stops a future change deleting any of them silently.

## Context
All three were verified by mutation. The reviewer deleted the code and the package stayed green. Line numbers below are from the main tree at commit `a55f24ce`; the reviewer's own numbers, taken in the REQ-506 worktree, are given where they differ.

**M1, the remedy redirection.** `redirectHelperRemedies` (line 343; reviewer: around 334-368) and its helper `advanceArgvCommandVerb` (line 357) rewrite a subordinate finding's remedy so it points at the continuation rather than at the subordinate command. Deleting both call sites — line 171 (reviewer: 168) and line 211 — leaves `go test ./internal/lifecycleadvance` fully green, because the existing tests assert record-level `NextArgv` and never a finding's rewritten remedy.

**M2, the layered guard.** `focusedGateState` (line 183) opens with `if subordinateState == resultmodel.AdvanceGateFailed || !focusedTest.Launched || focusedTest.TimedOut` (line 184; reviewer: 180). Deleting either half keeps the package green, because the eligibility guard `finishedOnItsOwn` in `internal/corehelpers/commands.go` (line 545, consumed at line 564; reviewer cites 565) already leaves those executions at `FocusedBaselineNotCompared`. The layering is correct and was explicitly asked for by the remediation plan — the finding is that nothing holds it there.

**The interrupted focused test.** That path's finding code changed from `BLOCKED-PROBE-LAUNCH-FAILED` to `BLOCKED-PROBE-FAILED` (`internal/corehelpers/commands.go` lines 550 and 554). Both are `failure` severity, so no existing assertion moved, and nothing pins the change. The one live assertion on either code is `evidence_gates_test.go:347`, which checks `BLOCKED-PROBE-LAUNCH-FAILED`.

This REQ amends REQ-506 (running the evidence gates from advance), which delivered all three and is in flight. It is not a defect report against that work — the behaviour is right, it is just unpinned.

## Red-Green Proof
**RED prompt/case:** Three mutations, run one at a time with `go test ./internal/lifecycleadvance`.
1. Delete the `redirectHelperRemedies(...)` call at `evidence_gates.go:171` and at `evidence_gates.go:211`.
2. Delete `subordinateState == resultmodel.AdvanceGateFailed ||` from the guard at `evidence_gates.go:184`, then restore it and instead delete `|| focusedTest.TimedOut`.
3. Change the interrupted focused-test path to emit `BLOCKED-PROBE-LAUNCH-FAILED` again in place of `BLOCKED-PROBE-FAILED`.

**Why RED now:** Every one of those mutations leaves the package green. There is no test that reads a finding's rewritten remedy, no test that reaches `focusedGateState` with an execution the eligibility guard has already excluded, and no public case for the interrupted focused test at all.
**GREEN when:** Each mutation, applied alone, fails at least one named test in `internal/lifecycleadvance`, and the failure message says which behaviour is gone. With all mutations reverted the package is green.
**Validation:** User confirmed — the mutations are the ones the reviewer performed, recorded as the RED case at the maintainer's instruction.

## Detailed Requirements
- Assert on a finding's rewritten remedy, not only on record-level `NextArgv`, so removing either `redirectHelperRemedies` call site fails.
- Cover both halves of the `focusedGateState` guard independently, so removing either one fails. Building the state needs an execution the eligibility guard in `internal/corehelpers/commands.go` already leaves at `FocusedBaselineNotCompared`; that is what makes the layering invisible to today's tests.
- Add a public case for an interrupted focused test that asserts the finding code is `BLOCKED-PROBE-FAILED`.
- State in each new test what it pins and which deletion it catches, so a later reader can tell it apart from a smoke test.

## Constraints
- Tests only. None of the three behaviours changes.
- `depends_on: [REQ-506]` because these tests pin code that request delivers; do not start before it completes.

## Builder Guidance
Apply one mutation at a time and keep it applied while writing the test that catches it. A test that passes under the mutation is the exact defect this REQ exists to close, and REQ-581 (making the descendant-cleanup tests fail on a real process-group leak) is the same failure in the same tool — worth reading before writing assertions here.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 7450 tokens, over the 2000-token budget in `actions/capture-reference.md` → Required Lessons Budget Contract, and `slugged: partial` so no targeted entry is legal. Matches on families `lifecycle-section-evidence`, `opaque-evidence-projection`, `smoke-vs-characterization` and `silent-skip-reads-as-red`.

## Full Context
See `do-work/user-requests/UR-119/input.md` for complete verbatim input.

*Source: independent review of REQ-506 (running the evidence gates from advance), findings M1 and M2, work run `do-work/runs/work-2026-09-05-003420/`.*
