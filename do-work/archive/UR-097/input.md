---
id: UR-097
title: 'Recover with asserted authority and keep the queue draining'
created_at: 2026-09-02T13:31:12Z
requests: [REQ-499, REQ-500, REQ-501]
word_count: 213
---

# Recover With Asserted Authority and Keep the Queue Draining

## Summary

Add a `do-work run-with-recovery` verb: plain `run` under the user's assertion that this checkout is the queue's only writer, so every ownership refusal `run` makes is answered "mine" and the loop keeps draining. Record the running principle, "one broken pipe doesn't stop the rest of the factory from running", where the run loop reads it. Make the archived-but-uncommitted state visible.

Origin: REQ-494's session died between archive and commit; the next `run` refused its first claim on a dirty `do-work/CHECKPOINT.md`. UR-096 (REQ-498) makes the tail journaled and resumable for unambiguous state; this UR covers what stays ambiguous, the visibility gap, and the principle. Six messages from one maintainer session are concatenated verbatim below.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-499 | `recover-finalization --assume-sole-releaser`: attribute the ambiguous shared-metadata remainder to the single discovered tail |
| REQ-500 | Surface unfinished finalizations in doctor and the session-start banner |
| REQ-501 | `do-work run-with-recovery` action, router row, and the one-broken-pipe principle (ADR-022, Execution Model, CLAUDE.md) |

## Batch Constraints

- Preserve the single-releaser model; the assertion is a deliberate verb invocation, never a flag that leaks into scripted `run`.
- Never widen recovery to secret-classified (protected-inventory `X`) paths or to project paths; only shared lifecycle metadata may be attributed by assertion.
- REQ-499 and REQ-500 depend on REQ-498 and have disjoint write sets; REQ-501 depends on REQ-499.
- Ordering: REQ-498, then REQ-499 and REQ-500 in parallel, then REQ-501.

## Folded Requests

- REQ-489 (remove-whole-checkpoint-entries-on-departure) — the checkpoint-writer bug found during the audit: the section heading is matched as a substring, so a backticked mention of the heading in Session Notes captured the live REQ-453 claim entry

## Full Verbatim Input

> ```
> audit the code for the following: So, since the issue, from what I understand, was that there was perhaps a context compaction happening and that stopped the processing in a weird state. Can we have a run-with-recovery command or something along the lines? Basically, as long as I'm certain that the current agent has full responsibility, full authority on implementing all of the requests, I should be able to queue a next command that will, assumingly, will run in a new context. and in that context it should be able to pick up the work where it stopped the last time. Would that be a good plan? Also propose names, run with authority or I don't know, run with recovery, run with authority and recovery, something like that.
> 
> another suggestion would be do-work run-all-here
> 
> this is the running principle: "one broken pipe doesn't stop the rest of the factory from running." where can we make a good note of it so it is taken into consideration
> 
> Live build: Yes, my other session. Verb: New verb `run-with-recovery`. Delivery: Capture as UR + REQs.
> 
> explain to me what is the blocker here? can't you still audit the code?
> 
> your current goal is to capture the request for the run-with-recovery so it can implemented and tracked
> ```

---
*Captured: 2026-09-02T13:31:12Z*
