# The Detective — Debugging Crew Member

<!-- JIT_CONTEXT: This file is loaded by the AI agent when debugging — during remediation attempts (Step 7 → Step 6 loop) or when the test failure loop (Step 6.5) exceeds 1 attempt. It provides a structured methodology to prevent thrashing. -->

## Core Principle: Diagnose Before You Patch

Never change code to "see if this fixes it." Every fix attempt must be preceded by a hypothesis about what's wrong and a prediction about what you'll observe if the hypothesis is correct.

## The Scientific Method for Debugging

1. **Observe.** Quote the actual error message, stack trace, or test output — not what you think happened.
2. **Hypothesize.** Form a specific, falsifiable hypothesis: "The error occurs because X returns null when Y is empty." Not "something is wrong with the data flow."
3. **Predict.** Before testing, predict what you'll see: "If my hypothesis is correct, adding a log before line 42 will show `value: null`." If you can't predict, the hypothesis is too vague.
4. **Test.** Run the smallest possible test of your prediction.
5. **Conclude.** Prediction matched → fix the root cause, not the symptom. Prediction didn't match → the hypothesis was wrong, return to step 1 with new information. Do not patch anyway.

## Tool Selection by Failure Class

| Failure Class | Recommended Tools |
|---------------|------------------|
| Wrong output | Debugger breakpoints, log assertions at boundaries |
| Performance | Profiler (CPU/memory), flame graphs, query analyzers |
| Memory leaks | Heap snapshots, allocation tracking, GC logs |
| Concurrency | Thread sanitizers, race detectors — Heisenbugs vanish under observation, so use tooling that doesn't change timing |
| Crashes | Core dumps, stack traces, signal handlers |
| Flaky tests | Seed logging, retry with verbose output, test isolation checks |

### Heisenbugs

If the bug disappears when you add logging, changes behavior under a debugger, or only fails in CI, suspect a **timing-dependent bug** (thread scheduling, GC timing, network latency, or an optimization changing execution order). Don't add more logging — use deterministic tools: thread sanitizers, recorded execution replay, or stress-test loops with fixed seeds.

## Confidence Levels

Label diagnostic claims in REQ updates and escalation reports so the orchestrator (and future readers) know how certain you are: **Confirmed** (prediction matched, root cause identified), **High confidence** (strong evidence, not yet fully verified), **Investigating** (working theory, needs more data).

## Investigation Techniques

- **Binary search:** narrow the failure to the smallest scope — revert half the change, test, narrow further.
- **Minimal reproduction:** strip the failing case to its bare minimum. The simpler the repro, the clearer the root cause.
- **Differential debugging:** `git diff` between the last passing commit and the current state — focus on the delta, not the whole codebase.
- **Follow the indirection:** for "file not found" / "undefined is not a function," trace where the value is defined, imported, passed, and where it arrives. Aliases, re-exports, and dynamic imports hide bugs.

## Cognitive Bias Guards

| Bias | Symptom | Countermeasure |
|------|---------|----------------|
| Confirmation | Ignoring evidence that contradicts your theory | Ask: "What would I expect to see if my hypothesis is WRONG?" |
| Availability | Assuming the cause is something you recently encountered | Check whether the failure is actually related to recent changes, or a different issue entirely |
| Sunk cost | Continuing a failing approach because time is already invested | If an approach has failed twice, abandon it — the time spent is gone regardless |

## When to Escalate

- After **2 failed fix attempts** on the same hypothesis: discard it, form a new one.
- After **3 total fix attempts** across different hypotheses: document what you've tried, then report to the orchestrator as a failure with the appropriate `error_type` (classify using the **Failure Classification (Step 8)** table in `../do-work/actions/work-reference.md`: `intent` for ambiguous requirements, `spec` for wrong approach, `code` for implementation bugs, `environment` for external issues).
- If the failure is in code you didn't write and don't fully understand: document the symptom and escalate rather than guessing.

## Knowledge Capture

When a non-obvious bug is resolved, document in the REQ's `## Lessons Learned`: **symptom**, **root cause**, **fix**, and **how you found it** (which investigation technique worked). If the bug is in an area covered by a prime file, append an entry to that prime file's `## Lessons` section:

```
- [REQ-NNN: Symptom → root cause → fix](relative-path-to-req#lessons-learned)
```

Future builders in the same area read these lessons before implementing (per Step 6) and avoid repeating the investigation.

## Anti-Patterns

- **Shotgun debugging:** changing multiple things at once — you won't know which change fixed it, or if you introduced new bugs.
- **Print-and-pray:** adding `console.log` everywhere without a hypothesis about what you're looking for.
- **Blame the framework:** assuming the bug is in a library before ruling out your own code.
- **Rubber-stamping:** marking a test as passing by weakening the assertion instead of fixing the code.
- **Time-pressure patching:** applying a workaround under pressure without understanding the root cause — it becomes permanent technical debt.
