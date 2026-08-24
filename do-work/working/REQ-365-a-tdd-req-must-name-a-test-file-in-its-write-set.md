---
id: REQ-365
title: "[impact-rule-change] A tdd REQ must name a test file in its write set"
status: claimed
claimed_at: 2026-08-24T20:17:59Z
status_changed_at: 2026-08-24T20:17:59Z
route: A
created_at: 2026-08-24T12:50:00Z
user_request: UR-069
addendum_to: REQ-346
domain: general
review_generated: true
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
impact: impact-rule-change
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-24T20:17:59Z
  basis:
    - trivial short-circuit
write_set:
  - skills/do-work/actions/capture.md
  - skills/do-work/actions/capture-reference.md
---

# A TDD REQ Must Name a Test File in Its Write Set

## What

REQ-346 was captured `tdd: true` with a `write_set` naming four files, none of them a test file. Its
GREEN was a behavioural claim, so it was **uncompletable as written**: the builder had to either
breach the declared scope or hand back an untested change. State the rule at capture so the shape
cannot be minted again.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The builder made the right call — it flagged the conflict, proceeded, and reported it, which is what
`CLAUDE.md` asks for when nobody is watching. But the conflict should not have existed. A REQ whose
`tdd: true` triggers the work loop's test-first evidence gate, while its write set forbids touching
any test, contradicts itself at capture time.

This is not a hypothetical: `write_set` is display-only by contract, so nothing broke. The cost is
that every such REQ spends a builder's judgment re-deciding the same question, and the resolution
lives in a commit message rather than in the queue.

## Detailed Requirements

- Capture's write_set guidance states that a `tdd: true` REQ names at least one test file, or omits
  `write_set` entirely (absence reads as unknown, which is honest; a set that excludes the tests the
  REQ requires is not).
- **Key it on the condition, not on this instance.** The rule is "a REQ whose completion requires
  writing a file class must not declare a write set that excludes it" — tests are today's case
  (`_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale).
- State it once where `write_set` is defined, and let the callers cite it. `capture-reference.md`'s
  Populating `write_set` guidance is the canonical home.
- Say what a builder should do when it meets the contradiction anyway, since queued REQs already
  carry it: flag, proceed, report — the behaviour REQ-346's builder chose.

## Constraints

- `_dev/primes/prime-action-files.md` governs. Read it first.
- Do not make `write_set` a gate. It is display-only at any builder count
  (`actions/work-reference.md` → Request File Schema) and this REQ must not change that — an
  invented set is worse than absence.
- `maintenance: true`: this narrows an instruction. Prefer the smallest addition that closes it.

## Red-Green Proof

**RED prompt/case:** `grep -n 'tdd' skills/do-work/actions/capture-reference.md` in the Populating
`write_set` section returns nothing — capture may mint a `tdd: true` REQ whose write set excludes
every test file, and did so for REQ-346.

**GREEN when:** the rule is stated once at the canonical home, keyed on the condition rather than on
tests specifically, and names what a builder does when it meets the contradiction in an
already-queued REQ.

**Validation:** Inferred during REQ-346's review.

---
*Source: REQ-346 review finding F2 (UR-069) — the generalisable half.*

## Triage

**Route: A** — This is a small, fully specified maintenance edit to two named instruction files: place the invariant at the canonical `write_set` definition and have the caller cite it without introducing enforcement or runtime behavior.
