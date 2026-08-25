---
id: REQ-373
title: "[impact-rule-change] Decide what the TDD-evidence gate accepts as test-first evidence"
status: completed
status_changed_at: 2026-08-25T08:47:30Z
claimed_at: 2026-08-25T08:38:00Z
completed_at: 2026-08-25T08:47:30Z
route: A
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
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-25T08:38:19Z
  basis:
    - trivial short-circuit
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
- [x] **[PLAN]:** Tighten `actions/work.md` so `tdd: true` accepts only a test in the project's
  existing automated harness, then make `actions/capture.md` route repeatable probe-only work to
  `tdd: false` plus strong Red-Green Proof. Keep `actions/capture-reference.md` byte-untouched because
  the confirmed answer preserves its existing harness-test invariant.
- [x] **[APPLY]:** Updated the work gate and capture assessment exactly as planned. The implementation
  changed only `actions/work.md` and `actions/capture.md`; the conditional target
  `actions/capture-reference.md` remained byte-untouched.
- [x] **[UNIFY]:** Reviewed `git diff --stat`, `git diff --check`, and every changed file. Verified
  `actions/work.md` states the accepted-evidence properties, `actions/capture.md` applies the same
  boundary at capture time, the REQ/checkpoint changes contain only lifecycle records, and no debug
  artifacts are present. `bash _dev/tests/maintainer-verify.sh` passed.

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

Certainty: **Firm** that the two texts disagree with observed practice. **The answer itself was the
user's and is now confirmed** in the answered Open Question above: require a project-harness test
for `tdd: true` and route probe-only work to `tdd: false` plus strong repeatable proof.

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

---

## Triage

**Route: A** - Simple

**Reasoning:** The user resolved the only policy choice, and the request names the two instruction
sites that must change. The confirmed harness-test boundary makes `capture-reference.md` an explicit
non-target.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

- `skills/do-work/actions/work.md` (modified) — clarified the accepted TDD-evidence properties.
- `skills/do-work/actions/capture.md` (modified) — aligned capture-time TDD classification.

The work gate now accepts a re-runnable test in the project's existing automated harness that fails
before implementation and passes after it, and classifies out-of-harness checks as regression proof
instead. Capture now sends repeatable probe-only work to `tdd: false` plus strong repeatable
before-and-after proof.

Verified `skills/do-work/actions/capture-reference.md` already carried the matching harness-test
invariant and left this conditional target byte-untouched, as required by the confirmed answer.

## Discovered Tasks

None.

## Testing

- Before/after policy proof: the prior Step 6.5 text required test-first evidence but did not state
  whether a repeatable out-of-harness check qualified; the revised text explicitly says it does not,
  and capture now routes that work to `tdd: false` plus strong repeatable proof.
- Conditional-target proof: `git diff --name-only -- skills/do-work/actions/capture-reference.md`
  returned no output, confirming the existing invariant remained byte-untouched.
- Canonical gate: `bash _dev/tests/maintainer-verify.sh` passed all contract, shell, Go, and strict
  JavaScript lanes. The optional strict browser lane was skipped because no browser was configured.

## Qualification

- Mechanical qualification passed with both implementation files present in the diff and all P-A-U
  phases completed.
- The changes trace every detailed requirement: accepted evidence is stated by observable property,
  capture uses the same boundary, no evidence gate was weakened, and the conditional reference
  target remained unchanged.
- Route A correctly has no `## Scope`; the implementation is limited to the two policy sites selected
  by the confirmed answer and introduces no new API, interface, type, field, or gate.

## Review

**Overall: 100%** | 2026-08-25T08:47:02Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None

**Minor findings:** 0
**Acceptance:** Pass — the work gate and capture assessment now return the same unambiguous result
for probe-only evidence, and the existing write-set invariant remains consistent.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

The restatement sweep checked live action, reference, crew, roadmap, and verification surfaces. No
live instruction treats repeatable out-of-harness evidence as sufficient for `tdd: true`; historical
changelog entries remain history.

*Reviewed by review-work action*

## Orientation

`tdd: true` now has one explicit boundary throughout the workflow: a re-runnable test in the
project's automated harness must go RED before implementation and GREEN afterward. Repeatable
probe-only work remains valid regression proof under `tdd: false`, with strong before-and-after
evidence.
