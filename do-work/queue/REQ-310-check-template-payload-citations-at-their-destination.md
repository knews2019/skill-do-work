---
id: REQ-310
title: "[impact-rule-change] Check a template payload's citations against where the payload lands"
status: pending
created_at: 2026-08-20T23:36:10Z
status_changed_at: 2026-08-21T00:00:00Z
user_request: UR-055
addendum_to: REQ-269
domain: general
impact: impact-rule-change
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set: []
---

# Check a Template Payload's Citations Against Where the Payload Lands

## What

REQ-269 split the fence exemption by its reason: a fenced payload is exempt because its text lands in *another* file, so it is not a citation from the file it is written in. That reasoning is sound and the exemption survives on it — but it only establishes that the payload is not checkable **from here**. It says nothing about whether the payload resolves **there**, and nothing checks that today.

A live example, found while confirming REQ-269's scope-drift was intended:

`skills/do-work-toolbox/actions/code-review.md:328` sits inside a ```markdown block that is a REQ template a reviewer copies into a new REQ file. It reads:

```
**Validation:** Review finding; apply `../do-work/actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
```

That payload lands in `do-work/queue/REQ-NNN-....md` in a consumer project. From there, `../do-work/actions/work-reference.md` resolves to `do-work/actions/work-reference.md`, which is not where the installed action lives — the installed path is under `.claude/skills/do-work/actions/`. So the citation is wrong at its destination, and both the pre-REQ-269 checker and the post-REQ-269 checker are structurally blind to it, for the same good reason.

`skills/do-work-toolbox/actions/validate-feedback.md:114` carries the same shape.

## Context

Discovered by REQ-269's orchestrator (REQ-269 `## Decisions` D-05) while verifying that its two declared-but-untouched files were correctly untouched. They were — this is the *next* question, and it is genuinely a different one: REQ-269 closed "is this token a citation from here", and this is "does a payload citation resolve where the payload goes".

Naming the shape honestly: this is the fourth marker. REQ-269 retired the backtick, the leading `../`, and the fence. Each retirement revealed the next boundary. This one is a boundary of a different kind — not a spelling standing in for a meaning, but a genuine gap in coverage that the correct exemption creates.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-269: a fenced template's payload is exempt from citation checking because it lands in another file, and nothing verifies it resolves at that destination. Two live templates carry a citation that is wrong where it lands. Should I process this as a new task? → Confirmed: Yes, add to queue
  *(2026-08-21)* User confirmed via `do-work clarify`: process as a new task. Nothing put out of scope — which of the three sizing options (hand-fix only, destination convention, installed-topology-only rule) to take is still open and belongs to whoever builds this REQ.

**Worth knowing before you decide**, because it changes the size a lot:

Checking a payload at its destination needs the destination to be **declared**, and today it is not. A fenced block does not say "this lands in a consumer REQ file" — a human infers it. So the honest options are:

1. **Fix the two sites, add nothing.** Smallest. Leaves the class open, which is exactly the posture REQ-269 spent its whole scope arguing against.
2. **Declare the destination where a template block is written** — an info-string convention (` ```markdown lands=do-work/queue/ `) or a marker comment — and have the checker resolve payload citations from there. Closes the class, costs a convention every template author must learn, and the checker can only check blocks that opted in.
3. **Treat installed-topology paths as the only correct form inside a template payload** and check that instead. No new convention, but it asserts a rule about template content that may not be true for every template.

Per `crew-members/maintenance.md` the deletion questions were asked: option 1 is the removal-shaped answer and it is genuinely on the table, because two hand-fixed sites may be the whole population. Counting the live template blocks that contain package-shaped paths is the first thing to do if you say yes — if it is two, option 1 is right and this REQ should close as a two-line fix.
