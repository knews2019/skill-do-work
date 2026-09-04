---
id: UR-099
title: 'Keep the factory running through recovery refusals'
created_at: 2026-09-02T20:35:18Z
requests: [REQ-513, REQ-514, REQ-515, REQ-516, REQ-517]
word_count: 481
---

# Keep the Factory Running Through Recovery Refusals

## Summary

When the maintainer runs `do-work`, the factory should keep working until everything that can be done is done. REQ-456 showed the opposite: a finished REQ's journaled finalizer refused its own uncommitted claim footprint, named itself as the fix, and both `run` and `run-with-recovery` stopped before selection, parking 31 pending REQs. The unblock was a hand-made claim commit, `cd9b01b0`. These five REQs remove the trigger, make the trap shape impossible, widen the escape hatch, and pin the sequence.

Origin: maintainer conversation of 2026-09-02 that diagnosed the REQ-456 stop against the principle "one broken pipe doesn't stop the rest of the factory" (UR-097, REQ-501).

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-513 | A1: claim commits its own footprint in every mode, serial included |
| REQ-514 | A2: a refusal never names its own command as `next_argv`; otherwise it sets the REQ aside |
| REQ-515 | A3: per-REQ recovery findings become selector exclusions; only REQ-less findings stop the run |
| REQ-516 | A4: the sole-authority assertion covers dirt the pipeline wrote earlier in the run; a remaining refusal names a third verb |
| REQ-517 | A5: one lock-in test for serial claim, complete, recover-finalization to `cleanup_complete` |

## Batch Constraints

- Order: REQ-513 first; REQ-517 rides with it. REQ-514 next; REQ-515 and REQ-516 depend on it and may run in parallel.
- Judgment stays prose; mechanics stay in the Go CLI. No new prose that walks a shell sequence.
- A guard may still refuse. What it may not do is refuse for a REQ-scoped reason in a way that stops unrelated REQs, or name itself as the fix.
- Nothing here widens recovery to secret-classified or project paths.
- Write sets overlap with REQ-504 and REQ-505 (UR-098) on `work.md`; overlap is declared, not a dependency.

## Full Verbatim Input

> ```
> [maintainer, message 1]
> basically when I run do-work I expect the factory to work and to figure out how to keep working until everything that can be done is done
> 
> [maintainer, message 2]
> how can I update the orchestrator to not end up in a trap like this?
> 
> [assistant reply the maintainer pointed at]
> The trap has one shape: a guard refuses, names itself as the fix, and a rule forbids anything else. You prevent it with two invariants the code enforces and one simplification, not with more prose in `work.md`.
> 
> A1, claim commits its own footprint in every mode. Serial mode defers the claim commit to Step 9 while worktree mode commits it at dispatch. That asymmetry is the whole bug: the finalizer only ever sees a dirty checkpoint in serial mode. Give `claim` a `--commit` that lands the queue-to-working move and the checkpoint line as one bookkeeping commit, and make Step 2 use it. Then complete never meets dirt it made, and the special-case acceptance from D2 is unnecessary. This is the delete-before-add answer.
> 
> A2, no self-referential refusal. A refusal must name a verb other than the one that produced it, or it is not allowed to refuse and must set the REQ aside instead. Enforce this in the result model, not in an action file: a test over every finding builder asserts `next_argv` differs from the invoking argv. The REQ-456 finding would have failed that test on the day it was written.
> 
> A3, REQ-scoped findings never stop the loop. Recovery in Step 1 runs per REQ. Each refused record becomes an exclusion with its reason code in the selector output, and selection continues with what remains. The only global stop is a finding with no owning REQ, which is what shared-target dirt looks like. The `finalizations` list is already ordered per record, so the CLI side is mostly there. The action side is the change: replace "continue only if every record is clean" with "exclude every record that is not".
> 
> A4, the escape hatch must be wider than the door. `rwr` today reuses the same Step 0.1 gate as `run` with one flag widened, so it can refuse for the exact reason `run` refused. Rule for that action: a sole-authority assertion covers every class of dirt the pipeline itself wrote earlier in the run. If a guard still refuses under `rwr`, it must name a third verb.
> 
> A5, pin the trap. One lock-in test: serial claim, build a one-line change, complete, then recover-finalization, asserting terminal `cleanup_complete`. Today that sequence fails, which is the real failure the test names.
> 
> Do A1 first. It removes the trigger and shrinks D2 to nothing. A2 and A3 make the next trap of this shape impossible rather than just this one. A4 is a sentence in `run-with-recovery.md`. A5 rides with A1.
> 
> [maintainer, message 3]
> do the unblock commit and capture A1-A5 as REQs
> ```

---
*Captured: 2026-09-02T20:35:18Z*
