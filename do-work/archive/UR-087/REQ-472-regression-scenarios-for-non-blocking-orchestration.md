---
id: REQ-472
title: '[impact-rule-change] End-to-end regression scenarios for non-blocking orchestration'
status: cancelled
created_at: 2026-09-01T04:29:16Z
user_request: UR-087
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-469, REQ-470, REQ-471]
maintenance: false
related: [REQ-468, REQ-469, REQ-470, REQ-471]
batch: non-blocking-orchestration
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [_dev/tests/contract-regressions.sh]
completed_at: 2026-09-03T11:47:26Z
---

# End-to-End Regression Scenarios for Non-Blocking Orchestration

## What

Lock in the batch's acceptance scenarios as regression tests covering serial and fan-out execution, crash recovery, repeated runs, and UR closure with blocked members — beyond the same-commit contract updates each sibling REQ already carries.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- "Regression tests cover serial and fan-out execution, crash recovery, repeated runs, and UR closure with blocked members."
- Scenario list (each names the real failure it pins):
  - REQ-A hits an unrelated canonical-gate failure while REQ-B is pending: A becomes blocked and leaves `working/`; B is subsequently claimed and processed — in serial and fan-out modes.
  - The unrelated failures create or fold into `pending-answers` REQs instead of being stored only in `CHECKPOINT.md`.
  - A blocked REQ with implementation edits cannot affect another REQ's diff, qualification, tests, staging, or commit.
  - A blocked REQ's checkpoint entry is removed when it leaves `working/`.
  - When the gate becomes green, the blocked REQ resumes from its preserved implementation and can complete normally.
  - A gate failure caused by the current REQ still follows remediation and code-failure handling rather than being misclassified as unrelated.
  - Repeated runs over a queue containing gate-blocked REQs stay stable (no re-hold, no duplicate fold targets, no checkpoint residue).
  - A UR with gate-blocked members is not closed by UR-closure readers until those members resolve.

## Constraints

- Tests are prose-contract and behavioral lanes in the `_dev/tests/` style (contract-regressions lanes, or a focused suite alongside the existing behavioral suites); every new lock-in names the real failure it pins — no decorative smoke tests.
- The instruction pipeline is prose executed by agents, so most scenarios pin the instruction contract (mutation-tested predicates over the shipped Markdown) rather than simulating an agent; where a shipped script participates (checkpoint helpers, `run-blocked-check.sh`, cleanup scripts), behavioral fixture tests are preferred.
- `bash _dev/tests/maintainer-verify.sh` remains the canonical baseline and must exit zero.

## Dependencies

- REQ-469, REQ-470, REQ-471 — the behaviors these scenarios pin must exist first. (Each of those REQs already updates the contract lanes it touches in the same commit; this REQ adds the cross-cutting end-to-end scenarios.)

## Builder Guidance

Certainty: Firm on the scenario list (it is the spec's acceptance-test list verbatim); latitude on contract-lane vs behavioral-fixture form per scenario. Fold a scenario into an existing lane where one already extracts the relevant block rather than duplicating extraction.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** After REQ-469/470/471 land, reverting the set-aside language in `actions/work.md` back to "preserve the claimed REQ and its checkpoint and stop" leaves every `_dev/tests/` suite green.
**Why RED now:** The new orchestration behavior would be unpinned — a later edit could silently restore the session-stopping hold.
**GREEN when:** The scenario lanes exist and that reversion (and each scenario's mutation) turns `bash _dev/tests/contract-regressions.sh` (or the new suite) red, while the unmodified tree passes `bash _dev/tests/maintainer-verify.sh` with exit zero.
**Validation:** Inferred during capture (from the spec's acceptance tests)

## Full Context
See `do-work/user-requests/UR-087/input.md` for complete verbatim input.

---
*Source: UR-087 — "Regression tests cover serial and fan-out execution, crash recovery, repeated runs, and UR closure with blocked members."*

## Cancelled

- **When:** 2026-09-03T11:47:26Z
- **Why:** folded into REQ-469 section Folded From REQ-472 as acceptance criteria for REQ-469, REQ-470, and REQ-471 (maintainer decision, 2026-09-03 roadmap triage)
- **Decided by:** user, via `do-work abandon`
