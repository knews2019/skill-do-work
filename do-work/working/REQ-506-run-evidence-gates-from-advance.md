---
id: REQ-506
title: '[impact-rule-change] Run the evidence gates from advance'
status: claimed
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-505]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/]
claimed_at: 2026-09-04T17:32:08Z
route: C
estimate:
  p50_active_minutes: 55
  confidence: low
  calculated_at: 2026-09-04T17:32:47Z
  basis:
    - Route C
    - 8-file write set
    - 2 new files
    - 3 subsystems involved
    - 3 acceptance criteria
    - dependency depth 2
    - cross-route regression gates
    - full-suite verification
---

# Run the Evidence Gates From advance

## What
Steps 3.6 (estimate), 5.75 (pre-flight), 6.3 (qualify) and the mechanical half of 6.5 (test gate and baseline comparison) run from `advance`; the Qualification Anti-Rationalization Table and the Finding-Closure Ratchet stay as principles.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Each gate already has a command (`estimate-p50`, `preflight`, `qualify`, `run-blocked-check`); the prose only says when to call it and what to paste.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- `advance` runs the gate for the current phase and reports missing evidence as typed findings; the agent's job is to satisfy the finding, not to know the command.
- Keep the anti-rationalization table and the Ratchet in prose, keyed on conditions.
- Delete the four steps' procedural prose and their predicates in the same commit; add Go tests per gate.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-505.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Delete Steps 3.6, 5.75, 6.3 and the gate half of 6.5 and run the contract suite.
**Why RED now:** Their predicates fail (Steps 5 and 6 carry 58 mentions between them); no Go test drives the gates in sequence.
**GREEN when:** Suite passes without those lanes; `advance` tests prove each gate refuses on missing evidence and passes on present evidence; the Ratchet and the table remain verbatim.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Triage

**Route: C — Complex.** This rule-changing migration composes estimation, preflight, qualification, and repository-gate evidence into the lifecycle command while removing procedural action contracts and replacing them with public Go behavior tests.

**Planning:** Required.
