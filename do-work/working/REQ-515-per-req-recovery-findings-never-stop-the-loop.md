---
id: REQ-515
title: '[impact-rule-change] Per-REQ recovery findings never stop the loop'
status: claimed
priority: now
created_at: 2026-09-02T20:35:18Z
user_request: UR-099
domain: general
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-514]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-513, REQ-514, REQ-516, REQ-517]
batch: recovery-never-traps
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/run-with-recovery.md, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/internal/finalization/]
claimed_at: 2026-09-04T18:15:54Z
---

# Per-REQ recovery findings never stop the loop

## What

Run Step 1 recovery per REQ. Each refused finalization or claim-recovery record becomes an exclusion with its reason code in the selector output, and selection continues with what remains. The only global stop left is a finding that owns no REQ, which is what shared-target dirt looks like.

The fold-first scan found REQ-469 (Replace the unrelated canonical-gate hold with a blocked set-aside) and REQ-504 (Collapse Step 10 and Crash Recovery prose into recovery) as neighbors: REQ-469 sets aside a gate failure inside a build, REQ-504 shortens the recovery prose once commands own it. Neither changes recovery's stop-versus-continue behavior, so this is a new REQ.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Both `run` and `rwr` put a global gate in front of the loop: Step 0.1 and Step 1 recovery are "recover everything or stop". REQ-456's stuck commit tail therefore parked 31 pending REQs. The maintainer's principle is that a failed REQ is set aside with a typed finding and the loop continues; only shared-target dirt may stop it.

## Context

`recover-finalization --discover` already returns an ordered `finalizations` list with one record per REQ. `actions/run-with-recovery.md` Step 0.1 says continue only when every record is terminal, and `actions/work.md` Step 1 has the same shape. The change is mostly on the action side, plus whatever the CLI needs to report per-record refusals as exclusions the selector understands.

## Detailed Requirements

- Step 1 in `actions/work.md` and Step 0.1 in `actions/run-with-recovery.md` iterate recovery records; a refused record excludes that REQ from this run's selection with its reason code and the loop continues.
- The composed exit summary lists set-aside REQs with their reason codes and resolving verbs.
- A finding with no owning REQ, such as dirt on a shared target that no REQ wrote, still stops the run, and it names a resolving verb per REQ-514.
- Contract predicates that pin "continue only if every record is clean" are replaced by predicates on the per-record wording, and the CLI carries a behavior test for a mixed result: one refused record, one clean record, selection proceeds.
- Serial and fan-out modes behave the same.

## Constraints

- Never widen what recovery accepts; this REQ changes what happens after a refusal, not whether it refuses.
- Keep the floor agent able to follow the loop with the command output plus the remaining prose.
- Coordinate wording with REQ-504 if both are in flight; the write sets overlap on `work.md` and `run-with-recovery.md`.

## Batch Constraints

- Judgment stays prose; mechanics stay in the Go CLI. No new prose that walks a shell sequence.
- A guard may still refuse. What it may not do is refuse for a REQ-scoped reason in a way that stops unrelated REQs, or name itself as the fix.
- Nothing here widens recovery to secret-classified or project paths; only dirt the pipeline itself wrote earlier in the run is in scope.
- Every REQ carries a behavior test on the command or a contract predicate on the action, never a sentence pin alone.

## Dependencies

Depends on REQ-514 for the set-aside finding shape. Related to REQ-469, REQ-472, and REQ-504.

## Builder Guidance

Certainty level: Firm on the behavior, latitude on how the exclusion is projected into the selector. Read `_dev/primes/prime-action-files.md` before touching an action file.

## Red-Green Proof

**RED prompt/case:** With REQ-456's journal at `prepared` and its checkpoint dirty, run `do-work run` on a queue with other claimable REQs.
**Why RED now:** Step 1 stops at the first refused finalization record and no other REQ is selected.
**GREEN when:** The same state reports REQ-456 as set aside with its reason code, selects the next claimable REQ, and the exit summary lists the set-aside with a resolving verb.
**Validation:** User confirmed (verify-requests, 2026-09-02).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3539 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing action routing and status contracts.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on semantic recovery completeness and structured evidence projection in do-work-cli internals.

## Full Context

See `do-work/user-requests/UR-099/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-02, item A3 of "how can I update the orchestrator to not end up in a trap like this?", captured by UR-099.*
