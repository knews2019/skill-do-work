---
id: REQ-373
title: "[impact-rule-change] Decide what the TDD-evidence gate accepts as test-first evidence"
status: pending
status_changed_at: 2026-08-25T08:19:12Z
created_at: 2026-08-24T23:37:06Z
user_request: UR-073
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-372]
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/capture.md
  - skills/do-work/actions/capture-reference.md
---

# Decide What the TDD-Evidence Gate Accepts as Test-First Evidence

## What

REQ-353 completed `tdd: true` without writing a test: its evidence was a generated-page RED/GREEN run
through Chromium. Step 6.5's stated bar is a failing **test** written first. Either the gate accepts
more than it says, or that REQ should have been `tdd: false`. Say which, in the text.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`actions/work.md` Step 6.5: "If the REQ has `tdd: true`, the `Red-green validation` section is
mandatory — the builder must show test-first evidence that they used RED/GREEN TDD (test written
before implementation, failed, then passed after). If this evidence is missing, treat it as a test
failure."

`do-work/archive/REQ-353-hide-dead-filter-knobs-on-durations.md` completed `tdd: true` anyway. Its
D-03 reads "**No new test file was added**; assembled-client tests, full module tests, and direct
generated-page RED/GREEN evidence cover the change," and its `## Testing` section records the RED and
GREEN as browser observations, not as a test that failed and then passed.

That REQ is archived and immutable, so this is not a request to redo it. What matters is the
mechanism: it was Route A with a one-file `write_set` naming no test, and the contradiction resolved
by lowering the evidence bar rather than by adding a test. If a runnable browser probe is acceptable
evidence, the gate should say so and stop reading as a test-file requirement; if it is not, capture's
TDD assessment should be sending this class of work to `tdd: false` plus a strong Red-Green Proof,
which is what its own heuristic already implies for behaviour "provable only by manual prompt/click/
visual inspection."

## Detailed Requirements

- Resolve the Open Question below, then make the two texts agree with the answer:
  - `actions/work.md` Step 6.5 — what the gate accepts, stated so a builder can tell before starting
    whether their planned evidence passes.
  - `actions/capture.md` Step 1 TDD assessment — only if the answer narrows what `tdd: true` should
    be captured for; leave it untouched otherwise.
  - `actions/capture-reference.md` → **Populating `write_set`**, the conditional completeness
    invariant — **conditional target.** It currently asserts that `tdd: true` "requires a runnable
    failing test to be written before it passes", so a declared set names at least one test file. If
    the answer accepts a committed re-runnable probe as evidence, that sentence contradicts the new
    policy and must be brought into line in the same change; if the answer keeps the harness-test
    bar, leave the invariant byte-untouched. Either way the invariant and Step 6.5 must state the
    same thing when this REQ closes.
- Whatever the answer, the accepted-evidence statement is keyed on a property of the evidence
  (runnable, fails before, passes after, re-runnable by another agent), never on a list of harness
  names.
- Do not weaken the gate into "any evidence counts". The failure this protects against is a
  behavioural claim standing in for a check nobody can re-run.

## Open Questions

- [x] Does a runnable, re-runnable browser/probe RED/GREEN satisfy `tdd: true`, or must the evidence
  be a test in the project's harness? → Confirmed: require a test in the project's harness
  Recommended: require a test in the harness, and route probe-only work to `tdd: false` plus a strong
  Red-Green Proof — that is what capture's own heuristic already says, and it keeps `tdd: true`
  meaning one thing.
  Also: accept a runnable probe when it is committed and re-runnable by another agent (then say so in
  Step 6.5 and name the property, not the tool); or accept it only for Route A REQs.

  **Answered 2026-08-25** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the builder's recommendation via `do-work clarify`: reserve `tdd: true` for work
  that writes a failing test in the project's normal automated test suite before implementation.
  Work supported only by a re-runnable browser or command-line probe stays `tdd: false` and must
  carry strong repeatable before-and-after proof. No additional scope was requested.

## Builder Guidance

Certainty: **Firm** that the two texts disagree with observed practice; **the answer itself is the
user's**, which is why this REQ is `pending-answers` rather than `pending`. Do not start the edit
before the Open Question is answered.

Scope cue: state what the gate accepts, in a property a builder can check before writing evidence.
No new field, no new gate, no rewrite of Step 6.5's surrounding flow.

Batch (UR-073): this REQ and REQ-372 overlap on `skills/do-work/actions/work.md` and, if the answer
triggers the conditional target above, on `skills/do-work/actions/capture-reference.md` as well. Run
them serially — REQ-372 first, since it moves a clause out of the same invariant paragraph.

## Constraints

- `_dev/primes/prime-action-files.md` governs. Read it first.
- `do-work/archive/` is immutable — REQ-353 is cited as evidence, never edited.
- No new gate, no new field, and no change to `write_set`'s display-only contract.

## Red-Green Proof

**RED prompt/case:** read Step 6.5's TDD verification paragraph in `skills/do-work/actions/work.md`
beside `do-work/archive/REQ-353-hide-dead-filter-knobs-on-durations.md` D-03. The rule demands a test
written first; the completed REQ records that no test was written, and it passed.

**Why RED now:** the gate's accepted-evidence standard is stated in one sentence about tests, and the
pipeline has at least one `tdd: true` completion that met it with browser evidence instead.

**GREEN when:** Step 6.5 states what evidence it accepts in terms a builder can check before writing
any of it, capture's TDD assessment agrees with that statement, and a REQ of REQ-353's shape reaches
one unambiguous answer under both.

**Validation:** Inferred during capture.

---
*Source: UR-073 finding F2 — found while reconciling this branch's capture against main's archive.*
