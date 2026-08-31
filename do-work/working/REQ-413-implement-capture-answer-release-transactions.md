---
id: REQ-413
title: 'Implement capture-file, answer, release, version, and changelog transactions'
status: claimed
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-412]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-08-31T22:07:13Z
route: C
estimate:
  p50_active_minutes: 105
  confidence: low
  calculated_at: 2026-08-31T22:08:35Z
  basis:
    - Route C
    - 24-file write set
    - 12 new files
    - 7 subsystems involved
    - 4 acceptance criteria
    - dependency depth 1
    - persistence changes
    - cross-route regression gates
    - full-suite verification
---

# Implement Capture-File, Answer, Release, Version, and Changelog Transactions

## What
Move deterministic publication and resolution phases for capture, answers, and releases into Go.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement structured `do-work-capture-files` publication with atomic UR/REQ/assets writes and reservation handling.
- Implement `do-work-answer` for applying answered questions with outside-text containment and status resolution.
- Implement parameterized `do-work-release` transactions for version files and changelog updates.
- Preserve exact input containment, linkage, timestamps, collision refusal, dry-run, optional commit, and rollback.

## Constraints
- LLM capture/answer/release actions may decide content but must delegate deterministic file and state mutations to the CLI.

## Dependencies
Depends on REQ-412 (shared lifecycle transaction behavior).

## Builder Guidance
Certainty level: Firm. Separate typed publication inputs from repository mutations so callers can inspect the full intended transaction before apply.

## Red-Green Proof
**RED prompt/case:** Publish a structured capture containing outside text and assets, apply an answer, and perform a parameterized release across collision and rollback fixtures.
**Why RED now:** These deterministic writes are currently assembled and mutated by action prose and shell helpers.
**GREEN when:** Each operation is atomic, byte-safe, collision-aware, text/JSON actionable, and leaves no partial files after a pre-commit failure.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Triage

**Route: C** — Complex

**Reasoning:** This request introduces three public mutation domains plus shared publication/release primitives, spans capture and answer actions, and must preserve atomic filesystem, Git, containment, version, changelog, and rollback contracts. The cross-action persistence surface requires explicit planning and exploration.

**Planning:** Required
