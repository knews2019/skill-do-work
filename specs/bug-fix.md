# Spec: Bug Fix

> Specification template for bug fixes — reproduce, isolate, fix, and prove.

## Output Structure

- **Fix** — minimal, targeted code change that addresses the root cause
- **Regression test** — a test that fails without the fix and passes with it
- **Root cause note** — appended to the REQ explaining *why* the bug existed, not just what was changed

## Quality Standards

- Root cause identified — not just the symptom patched. The fix addresses *why* the bug happened, not just *where* it manifested.
- Regression test proves the fix — a test that reproduces the original bug (fails on the old code) and verifies it's fixed (passes on the new code)
- No unrelated changes — the diff contains only the fix and its test. No drive-by refactoring, no style cleanup, no "while I'm here" improvements.
- Fix is minimal and targeted — the smallest change that correctly addresses the root cause
- Similar bugs checked — after finding the root cause, scan for the same pattern elsewhere in the codebase

## Implementation Checklist

1. Reproduce — confirm the bug exists and is reproducible. Document the reproduction steps.
2. Isolate root cause — trace from symptom to origin. Understand *why* it happens, not just *where*.
3. Write failing test — a test that captures the buggy behavior. Run it — confirm it fails for the right reason.
4. Implement minimal fix — the smallest correct change. Resist the urge to refactor surrounding code.
5. Verify test passes — the failing test now passes. All other tests still pass.
6. Check for similar bugs — search for the same pattern elsewhere. If found, note in Discovered Tasks (don't fix inline).
7. Document root cause — add a brief explanation to the REQ so reviewers and future developers understand what went wrong.

## Evolution Path

- **Simple**: Typo, config, or off-by-one fix — obvious root cause, single-line fix, straightforward test
- **Medium**: Logic error with test — incorrect conditional, wrong variable, missing null check. Requires understanding the control flow.
- **Complex**: Systemic issue requiring design change — race condition, architectural flaw, incorrect abstraction. Fix may touch multiple files and require a plan.

## Common Pitfalls

- Fixing symptoms not causes — adding a null check where the real problem is that the value should never be null. The null check hides the upstream bug.
- Missing regression test — without a test, the same bug can return in a future change. The test is the proof, not the fix.
- Scope creep into refactoring — a bug fix is not an excuse to rewrite the surrounding code. If the code needs refactoring, create a separate REQ.
- Not checking for similar bugs — if the same pattern caused this bug, it likely caused others. A 30-second grep can prevent future bug reports.
- Over-engineering the fix — a complex fix for a simple bug suggests misdiagnosed root cause. Step back and verify the root cause is correct.
