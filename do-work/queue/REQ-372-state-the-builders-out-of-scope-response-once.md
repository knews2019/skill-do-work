---
id: REQ-372
title: "[impact-rule-change] State the builder's out-of-scope response once"
status: pending
created_at: 2026-08-24T23:37:06Z
user_request: UR-073
addendum_to: REQ-365
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
impact: impact-rule-change
effort_estimate: effort-mechanical
write_set:
  - skills/do-work/actions/capture-reference.md
  - skills/do-work/actions/work.md
---

# State the Builder's Out-of-Scope Response Once

## What

REQ-365's shipped invariant tells a builder to flag the contradiction and **proceed**;
`actions/work.md` tells a builder discovering an out-of-scope file to **stop and report to the
orchestrator**. Both are shipped, both describe the same moment, and they prescribe different
actions. State the response once, in the file the builder actually loads, and cite it from the other.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The two texts diverge on what the builder does next:

- `actions/capture-reference.md` → **Populating `write_set`**, conditional completeness invariant:
  "If a builder meets this contradiction in an already-queued REQ, flag it before editing, **proceed
  with the required file class despite the declared set**, and report the contradiction plus the
  actual files in the handback."
- `actions/work.md` § **Write only inside the declared scope**: "Discovering mid-build that you need
  a file outside it is a **stop-and-report to the orchestrator, never a silent write.** The
  orchestrator records the request and its resolution in the REQ trail as a `## Decisions` D-XX entry
  (it is a scope judgment) and extends both the Scope list and `write_set`."

Proceed-then-report and stop-then-wait are not the same instruction. A contradiction between two
shipped instructions changes what an agent does, which is why it earns a REQ rather than a prose
note (`actions/capture-reference.md` → Fold-First Rule → the prose-only test's second exemption).

The divergence is also in the wrong place: builders load `actions/work.md`, not the capture
reference. A builder-behaviour rule that only capture reads is a rule with no reader.

## Detailed Requirements

- Decide which behaviour is intended and state it **once**, in `actions/work.md` § "Write only inside
  the declared scope". The proceed-anyway case is a real one — REQ-346's builder proceeded, the
  orchestrator upheld it, and an unattended fan-out builder has no orchestrator to wait on — so if
  that is the intended response, it belongs in that bullet as a stated condition, not as a second
  rule elsewhere.
- Reduce the invariant's sentence in `actions/capture-reference.md` to a citation of that bullet.
  The invariant itself stays where it is: it is a capture-time rule about what a declared set must
  contain, and only its builder-response clause is being relocated.
- Key the surviving statement on the condition (a needed file outside the declared set), not on
  tests or on `tdd: true` — the invariant already establishes that tests are one instance
  (`_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale).
- Keep both files' existing guarantee intact: `write_set` stays display-only, and an honest absent
  set stays better than an invented one.

## Constraints

- `_dev/primes/prime-action-files.md` governs. Read it first.
- Delete before you add: the fix is one statement plus one pointer, not a third statement
  reconciling the two.
- Do not touch the invariant's capture-time half, and do not make `write_set` a gate
  (`actions/work-reference.md` → Request File Schema).

## Red-Green Proof

**RED prompt/case:**

```bash
grep -n "proceed with the required file class" skills/do-work/actions/capture-reference.md
grep -n "stop-and-report to the orchestrator" skills/do-work/actions/work.md
```

Both return a line today, each prescribing a different next action for the same moment.

**Why RED now:** REQ-365 added the builder-response clause to `capture-reference.md` without
reconciling it against the declared-scope bullet that already governed builders.

**GREEN when:** exactly one of the two files states what the builder does, that file is
`actions/work.md`, the other cites it by section name, and the stated response covers the unattended
case that made REQ-346's builder proceed.

**Validation:** Inferred during capture.

---
*Source: UR-073 finding F1 — reconciling this branch's capture against REQ-365 as shipped on main (`6265f1c`).*
