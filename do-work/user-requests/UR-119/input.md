---
id: UR-119
title: 'Three independent-review findings from the 2026-09-05 work run'
created_at: 2026-09-05T01:30:57Z
requests: [REQ-581, REQ-582, REQ-583]
word_count: 712
---

## Summary

Three findings raised by independent reviewers during the work run at `do-work/runs/work-2026-09-05-003420/`. Each one says a control does not do what it claims: a test that cannot fail on the leak it names, a check that cannot see one citation form, and new behaviour that no test pins. None of them is speculation — every one was measured by mutating the code and observing that the suite stayed green.

## Extracted Requests

| REQ | Finding | Reviewing request | Where |
| --- | --- | --- | --- |
| REQ-581 (make the descendant-cleanup tests fail on a real process-group leak) | 1 | independent review of REQ-506 (run the evidence gates from advance) | `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` |
| REQ-582 (make the shipped-package reference contract see the arrow-form section citation) | 2 | discovered task from the REQ-510 builder (sweep work-reference sections owned by CLI tests) | `_dev/tests/shipped-package-reference-contract.sh` |
| REQ-583 (pin the evidence-gate remedy redirection, layered guard and interrupted focused-test code) | 3, findings M1 and M2 | independent review of REQ-506 (run the evidence gates from advance) | `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go` |

## Batch Constraints

- Every REQ's RED is the mutation the reviewer performed: delete or neutralise the code, run the suite, observe it stay green. Reproduce that mutation rather than inventing a substitute proof.
- Findings 1 and 3 are separate root causes in separate packages and share no files. Finding 1 is pre-existing; finding 3 is new behaviour delivered by REQ-506 (run the evidence gates from advance) and therefore carries `addendum_to: REQ-506`.
- Represent, do not expand. These are measured findings recorded at the size they were measured.

## Capture Verification

Finding 2 asked for the claim to be verified before writing the request. It was, and it holds. `bash _dev/tests/shipped-package-reference-contract.sh` exits 0 on the tree at commit `a55f24ce` while two live shipped citations name sections that do not exist in their target: `skills/do-work/CHANGELOG.md:279` cites `actions/work-reference.md` → **Recovery Refusals (Step 1)** (that heading was renamed to `## Stuck Runs Hand Off to Judgment (any step)`), and `skills/do-work/actions/cleanup.md:48` cites `actions/work-reference.md` → **In-Progress Record (Step 2)** where the only heading is `## In-Progress Record (Step 1)` at line 482. The script contains no occurrence of the arrow character and no bold-name parsing at all; its `citation_shape` regex matches path-shaped tokens only, and anchors are checked solely from a `#fragment` inside the token. 75 arrow-form citations ship across the three skill packages.

## Full Verbatim Input

> ```
> Capture these as three requests under one user request. They came from three separate independent reviews in the run at `do-work/runs/work-2026-09-05-003420/`.
> 
> ### Finding 1 — the descendant-cleanup tests cannot fail on a real process-group leak
> Source: the independent review of REQ-506.
> 
> Three tests in `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` claim to prove that a probe's descendant process group is terminated: `TestBlockedProbeTimeoutKillsDescendantGroup`, `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits`, and `TestBlockedProbeInterruptionIsTypedAndReapsDescendants`.
> 
> The reviewer reduced `terminateOwnedProcessGroup` and `cleanupReapedProcessGroup` in `blocked_probe_unix.go` to no-ops — a genuine leak — and the first two tests still passed, taking 30.01s and 31.35s instead of 2.90s and 2.01s. The third failed, but on its own 5-second "interrupted probe did not return" bound, not on its descendant assertion.
> 
> The structural reason: a surviving descendant inherits the parent-owned diagnostic pipe, so `<-diagnosticDone` and `<-done` hold the runner until that descendant exits. By the time the poll loop runs, the descendant is always already gone. What the loop actually measures is how long init takes to reap a zombie after the runner returns — which is why it flaked at 1.13–1.95s against a 2-second budget under load.
> 
> So these are reaping-latency assertions, not termination proofs. A leak shows up only as the test taking 30 seconds instead of 3. This is pre-existing and was not introduced by REQ-506.
> 
> The reviewer's suggested fix, which you should record as the RED case: a descendant that closes its inherited stdout and stderr so it cannot hold the diagnostic pipe open, and that outlives the budget. Under a real leak that fixture must fail on the descendant assertion, not on a timeout.
> 
> Impact: this is a control that gives false confidence about process cleanup. Suggest impact-rule-change.
> 
> ### Finding 2 — `shipped-package-reference-contract.sh` cannot see one citation form
> Source: the REQ-510 builder's discovered task, surfaced during the work-reference sweep.
> 
> `_dev/tests/shipped-package-reference-contract.sh` does not detect citations written in the `→ **Named Heading**` form. During REQ-510's RED run the check passed while two live callers named a heading that had just been deleted.
> 
> Every sweep that deletes or renames a section in a shipped action or reference file depends on this check to catch dangling inbound references. For this citation form it gives false confidence. Verify the claim before writing the request, and record what you find either way.
> 
> Impact: affects every future prose sweep, not one request. Suggest impact-rule-change.
> 
> ### Finding 3 — two pieces of new behaviour in the evidence gates are unpinned
> Source: the independent review of REQ-506, findings M1 and M2.
> 
> Both were verified by mutation — the reviewer deleted the code and the package stayed green.
> 
> - `redirectHelperRemedies` and `advanceArgvCommandVerb` in `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go` (around lines 334-368) are new behaviour delivered for a prior review finding. Deleting both call sites at lines 168 and 211 leaves `go test ./internal/lifecycleadvance` fully green. The tests assert record-level `NextArgv` and never a finding's rewritten remedy.
> - The layered guard in `focusedGateState` at `evidence_gates.go:180` reads `subordinateState == AdvanceGateFailed || !focusedTest.Launched || focusedTest.TimedOut`. Deleting either half keeps the package green, because the eligibility guard in `internal/corehelpers/commands.go:565` already leaves those executions at `not_compared`. The layering is correct and was explicitly asked for by the remediation plan — the point is that nothing stops a future change removing it silently.
> 
> The reviewer also suggests a public case for an interrupted focused test, since that path's finding code changed from `BLOCKED-PROBE-LAUNCH-FAILED` to `BLOCKED-PROBE-FAILED`; both are `failure`, but nothing pins the move.
> 
> Impact: suggest impact-user-visible for the remedy redirection, and fold the guard and the interrupted-path pin into the same request — they are one root cause, new behaviour landing without a lock.
> 
> ## How to write them
> - Represent, don't expand. These are measured findings; state them at the size they are.
> - Each request needs a real `## Red-Green Proof` whose RED is the mutation the reviewer actually performed, because that is what makes each one testable.
> - Preserve the full detail above in the UR's verbatim section — do not summarise it away.
> - Set `prime_files`, `domain`, `impact`, `effort_estimate` and `depends_on` per the templates. Findings 1 and 3 touch `skills/do-work/tools/do-work-cli/`, so `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` is the relevant prime; finding 2 touches `_dev/tests/`, so `_dev/primes/prime-shell-commands.md` is.
> - Note in the UR that all three came from independent review during the 2026-09-05 work run, and name the reviewing request for each.
> - Capture does not execute. Write the files and stop. Do not start any of this work.
> ```
