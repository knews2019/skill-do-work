---
id: REQ-309
title: "[impact-rule-change] Run the repo's canonical gate before hand-back, not only the changed area's tests"
status: pending-answers
created_at: 2026-08-20T23:18:07Z
status_changed_at: 2026-08-20T23:18:07Z
user_request: UR-055
addendum_to: REQ-262
domain: general
impact: impact-rule-change
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set: []
---

# Run the Repo's Canonical Gate Before Hand-Back

## What

REQ-283 archived with a Testing section listing four green checks — `go test ./...` and a fresh build in `skills/do-work-board/tools/queue-kanban`, plus `queue-kanban verify` returning `OK: no findings`. Every one of them was true. None of them was `bash _dev/tests/maintainer-verify.sh`, which the change had just turned red by adding a second `./actions/board.md` routing row that `_dev/tests/staged-skills-contract.sh` counts.

The gate stayed red across REQ-279, REQ-295 and REQ-283's own metadata commits, and REQ-262's run was the first to notice — because Step 5.75's pre-flight ran the command and reported the baseline failing.

`actions/work.md` Step 6.5 resolves test commands from the prime file's testing section, then falls back to generic detection per changed file. Both are **area-scoped by construction**: they answer "what tests cover the files I touched", which is the right question for a regression and the wrong question for a repo that also has one whole-repo gate its own `CLAUDE.md` calls "the canonical baseline pass/fail check before any hand-back."

## Context

Discovered by REQ-262's orchestrator while repairing the gate in order to be able to verify REQ-262 at all (REQ-262 `## Decisions` D-01). The repair itself shipped separately at version 0.222.5, commit `8e9cc46` — this REQ is about the process gap that let the breakage land and persist, not about the four files it broke.

Worth noting what did *not* fail: REQ-283's review scored it and passed it. The Restatement Sweep (`actions/review-work.md` Step 6) should in principle have caught a routing row whose count another file asserts on, but a routing table is not obviously "something other text restates" until you know the contract test counts its rows.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-262: the work pipeline's test resolution is area-scoped, so a change can pass Step 6.5 and Step 7 while leaving the repository's own canonical whole-repo gate red — which is what happened to `maintainer-verify.sh` at REQ-283 and stayed that way for three REQs. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — Step 5.75's pre-flight already surfaces a red baseline on the next run, so the gap is self-limiting to one REQ's blast radius plus however long until someone runs the queue again.

**Two things to decide if you say yes**, because they pull in different directions:

1. **Where the knowledge lives.** A prime file's testing section is the sanctioned home for project-specific test mapping, and this repo's `_dev/primes/` has no such section for the gate. Adding one is the smallest change and needs no action-file edit — but it only helps REQs whose `prime_files` list it. The alternative is a line in `actions/work.md` Step 6.5 that says a whole-repo gate, where the project declares one, always runs in addition to the area-scoped commands. That reaches every REQ and costs a rule.
2. **Whether it is a gate or a report.** Running the gate on every REQ is slower (it is a multi-minute suite here) and would have blocked REQ-279 and REQ-295 on a failure neither caused. Reporting it without blocking keeps those REQs moving but is exactly the posture that let three of them ship past a red gate without anyone noticing.

Per `crew-members/maintenance.md`, the deletion questions were asked first: the drift is not caused by a stale source, a bad example, or a too-broad tool. It is a genuinely absent instruction, which is why this is written as an addition candidate rather than a removal — and why it wants a replay case (REQ-283's diff, which must fail the new check) before anything is added.
