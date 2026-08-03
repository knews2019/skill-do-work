---
id: REQ-079
title: Two guards pin the weaker fingerprint of the premise they exist to retire
status: pending
created_at: 2026-08-03T16:53:42Z
user_request: UR-015
domain: testing
prime_files: []
tdd: true
depends_on: []
maintenance: true
addendum_to: REQ-075
---

# Two guards pin the weaker fingerprint of the premise they exist to retire

## What

REQ-075 (v0.166.2) established that a retired premise leaves two fingerprints — the thing it *said*
("one REQ at a time") and the thing it was *called* ("under the exclusive-session model") — and that
the second is the more dangerous of the two. Its own regression check pins only the first. Its regex is
also narrower than the class it names. And `actions/cleanup.md:31` still argues the safety of the
skill's one destructive pass from a premise REQ-071 spent an entire REQ falsifying.

No live defect today: the shipped tree is clean of both fingerprints. This is a durability gap — the
guards will not catch the recurrence they were written for.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

A guard that matches only the wording that was already removed is worse than no guard: the suite goes
green, and the green reads as coverage. REQ-075 wrote down the exact insight needed to avoid this —
in `tools/queue-kanban/prime-do-kanban.md:60`, where it calls the weak form "the more dangerous,
because the model it names is still true and only its relevance died" — and then did not encode it.
The lesson is filed; the assertion does not implement it.

## Context

**Finding F4a — no assertion covers the weak fingerprint.** `_dev/tests/contract-regressions.sh:137`
greps `one REQ( [a-z]+)? at a time` and filters to lines mentioning `write_set|overlaps`. Nothing
asserts against the weak form near `write_set`. Verified: `grep -n "exclusive-session" _dev/tests/contract-regressions.sh`
returns lines 184, 190, 191, 205, 213, 222 — all about REQ-069's removed machinery, none about
`write_set`.

**Finding F4b — the regex is narrower than the class.** The optional word slot sits *after* `REQ`, so
the pattern matches "one running REQ at a time" but not:

- `one active REQ at a time`
- `a single REQ at a time`
- `one builder at a time`
- `only one REQ is ever building`

Each of those reintroduces the premise and passes the suite. The same regex is reused verbatim in both
file-level negatives (lines ~153 and ~158), so the gap is triplicated.

**Finding F7 — `actions/cleanup.md:31` argues from the falsified premise.** It reads: "**Safe under the
exclusive-session model.** This pipeline assumes no other `do-work` session is running against this
checkout … so there is no live coexisting claim to protect — Pass 0 needs no lock and consults none."
REQ-071 exists *because* a live coexisting claim can be there. Pass 0's behavior is probably still
safe — it sweeps only terminal statuses, and integration is serial — but this is the skill's one
destructive pass reasoning from the premise the safety REQ removed, and its next sentence ("that
session's own REQ has already been moved out of `working/`") is singular where fan-out is plural.

## Detailed Requirements

1. **Add a weak-form assertion** for the `write_set` rule: no shipped file may justify `write_set`'s
   display-only status by naming the exclusive-session model. Keep REQ-075's per-class shape — a line
   sweep for prose, file-level negatives for the comment-carrying sources (`tools/queue-kanban/model.go`,
   `tools/queue-kanban/web/board.js`) whose comments wrap past a line-granularity grep.
2. **Widen the strong-form regex** to cover the class rather than one phrasing. Move the optional word
   slot so it can precede `REQ`, and cover the `single`/`builder` variants. Per the Closed Enumerations
   rule, the assertion's comment must state the trigger *condition* — "no shipped file may argue
   `write_set`'s display-only status from any builder count" — and mark the matched wordings as
   illustrative, not exhaustive.
3. **Do not let the widened pattern false-positive on the canonical section.** REQ-075's original
   comment records why granularity is one line: `actions/work-reference.md`'s Fan-Out Dispatch section
   says "integration … runs one REQ at a time" two lines below the advisory-`write_set` bullet, and
   that statement is **true**. Re-verify after widening that the suite still passes on a clean tree —
   a widened pattern is exactly the change that could start matching it.
4. **Correct `actions/cleanup.md:31`.** Replace the exclusive-session justification with the durable
   one: Pass 0 is safe because it sweeps only *terminal* statuses and integration is serial under one
   queue owner, not because no other claim can exist. Fix the singular in the second sentence — under
   fan-out several REQs sit in `working/` at once. Point at
   `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch rather than restating it.
5. **Sweep the remaining weak-form sites and rule on each.** `actions/work.md:549` and
   `actions/cleanup.md:40` also invoke the exclusive-session model. Some uses are legitimate — the
   model *is* still true, only its relevance to a given conclusion died. For each site, state whether
   the conclusion still follows from the premise, and correct only the ones where it does not. Do not
   sweep the phrase mechanically; that would delete true statements.
6. **Leave the badge's behavior alone.** No Go logic, no schema field, no board column. Same as
   REQ-075: the explanation changes, the code does not.

## Constraints

- **This REQ can only be proven by making the suite fail.** A guard change that is never observed
  failing is not evidence of anything. Every new or widened assertion must be demonstrated red by
  reintroducing the wording it targets, then green after reverting — record the observed failure text.
- Requirement 5 is where this REQ can do damage. `crew-members/maintenance.md`'s delete-before-you-add
  rule cuts both ways: removing a *true* premise because it pattern-matches a retired one is the
  failure mode here. Err toward leaving a site alone and saying why.
- `_dev/tests/contract-regressions.sh` is also touched by REQ-077 (the REQ-071 guard at ~line 364).
  Different block; if both are ever built concurrently the merge is the non-interference proof, not the
  overlaps badge.

## Dependencies

`addendum_to: REQ-075` — completes the guard REQ-075 specified in its lesson and under-implemented in
its assertion. Overlaps REQ-077 in one file (see Constraints). No `depends_on`: buildable immediately.

## Builder Guidance

**Certainty: Firm on requirements 1–4, genuinely open on 5.**

The regex and assertion work is mechanical and the findings were verified by grep; don't re-derive
them, but do re-check line numbers since REQ-077 may move them.

Requirement 5 needs judgment, and the right answer may well be "all three remaining sites are fine."
That is a valid outcome — record the per-site ruling rather than forcing a change to justify the REQ.
The value here is the two assertions, not a body count.

Scope note: this REQ is the guard, not another sweep. If requirement 5 turns up a *new* class of stale
reasoning beyond the exclusive-session premise, capture it as a discovered task rather than expanding
this REQ — REQ-075's scope doubling mid-build is the precedent to avoid repeating.

## Red-Green Proof

**RED case:** Add the line
`` `write_set` is display-only because only one active REQ builds at a time under the exclusive-session model. ``
to `actions/board.md` and run `bash _dev/tests/contract-regressions.sh`. It passes today — the strong
regex misses "one active REQ" (the optional slot is on the wrong side of `REQ`), and no assertion
covers the exclusive-session clause at all.

**Why RED now:** Both guards match one historical phrasing rather than the class, so the premise
REQ-075 retired can walk back in under any of four wordings, including the one REQ-075's own lesson
names as the more dangerous.

**GREEN when:** (1) The RED line above fails the suite and the failure message names the file and the
fix. (2) Each of the four wordings in Context fails the suite. (3) The suite still passes on an
otherwise clean tree, with Fan-Out Dispatch's true "integration runs one REQ at a time" sentence
untouched. (4) `actions/cleanup.md:31` no longer argues Pass 0's safety from the exclusive-session
premise. (5) The per-site rulings from requirement 5 are recorded, including the leave-alone ones.

**Validation:** Inferred during an adversarial audit; remediation plan reviewed and approved by the
user before capture.

## Full Context

See `do-work/user-requests/UR-015/input.md` for the audit's provenance and the findings it cleared.
