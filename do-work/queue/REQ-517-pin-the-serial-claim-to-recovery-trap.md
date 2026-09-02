---
id: REQ-517
title: 'Pin the serial claim-to-recovery trap'
status: pending
created_at: 2026-09-02T20:35:18Z
user_request: UR-099
domain: testing
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-513]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-513, REQ-514, REQ-515, REQ-516]
batch: recovery-never-traps
write_set: [skills/do-work/tools/do-work-cli/internal/finalization/]
---

# Pin the serial claim-to-recovery trap

## What

Add one lock-in test that runs the serial shape end to end in a fixture repository: claim, a one-line implementation change, complete, then `recover-finalization --discover`, asserting the terminal phase is `cleanup_complete`. Today that sequence stops at lifecycle apply, which is the real failure the test names.

The fold-first scan found no pending or pending-answers REQ that owns this sequence; REQ-472 (End-to-end regression scenarios for non-blocking orchestration) covers gate set-asides, not the claim-to-finalize tail.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The factory halt behind REQ-456 was never pinned. Route C REQs commit the claim first and so never exercise the serial path through the journaled finalizer. Without this test the next refactor of claim or finalize can reintroduce the trap silently.

## Context

`internal/finalization/finalization_recovery_test.go` and `finalization_req499_test.go` already build fixture repositories for recovery; reuse their helpers. The test is RED before REQ-513 lands and GREEN after, so it rides in the same run.

## Detailed Requirements

- One Go test in the finalization package: serial claim through `claim --commit`, a one-line tracked-file change, `complete`, then `recover-finalization --discover`.
- Assert terminal phase `cleanup_complete`, empty `blocked_paths` and `reason_codes`, and that the primary commit contains the checkpoint, archive, and implementation paths.
- A second case runs the pre-REQ-513 shape, claim without commit, and asserts the refusal names a verb other than `recover-finalization` once REQ-514 lands; until then it asserts only the refusal code.
- No new fixture helper if the existing ones suffice.

## Constraints

- Focused, not a smoke suite; two cases, each naming the failure it pins.
- Runs under `go test ./...` and the maintainer verify script without a real browser or network.

## Batch Constraints

- Judgment stays prose; mechanics stay in the Go CLI. No new prose that walks a shell sequence.
- A guard may still refuse. What it may not do is refuse for a REQ-scoped reason in a way that stops unrelated REQs, or name itself as the fix.
- Nothing here widens recovery to secret-classified or project paths; only dirt the pipeline itself wrote earlier in the run is in scope.
- Every REQ carries a behavior test on the command or a contract predicate on the action, never a sentence pin alone.

## Dependencies

Depends on REQ-513. Related to REQ-514 and REQ-516.

## Builder Guidance

Certainty level: Firm. Reuse the recovery test fixtures; do not build a new harness.

## Red-Green Proof

**RED prompt/case:** Run the new test against the tree before REQ-513.
**Why RED now:** The serial shape leaves the checkpoint dirty at complete time and `recover-finalization` refuses at lifecycle apply.
**GREEN when:** The test passes after REQ-513, and deleting the `--commit` from the claim step in the fixture makes it fail again with the refusal code.
**Validation:** Inferred during capture; the A1 through A5 list this REQ transcribes was confirmed by the maintainer in the same conversation.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on semantic recovery completeness and structured evidence projection in do-work-cli internals.

## Full Context

See `do-work/user-requests/UR-099/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-02, item A5 of "how can I update the orchestrator to not end up in a trap like this?", captured by UR-099.*
