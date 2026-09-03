---
id: REQ-505
title: '[impact-rule-change] Move selection and claim behind advance'
status: pending
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-504]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/]
---

# Move Selection and Claim Behind advance

## What
Steps 1 (95 lines), 2.0 and 2 become `advance` phases; `work.md` keeps one selection principle (what makes a REQ claimable, keyed on the condition) and nothing procedural.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Selection already runs through `next` and claim through `claim`; the 114 lines of prose restate their contracts.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- `advance` on an unclaimed queue runs `next` with the session's targeting tokens and fan-out bound, then `archive-collision` and `claim`, reporting each as a phase result.
- The targeted-run ledger contract (REQ-453) moves from prose into `advance` state; its lanes become Go tests.
- Delete the Step 1, 2.0 and 2 prose and their predicates in the same commit.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-504.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Delete Steps 1, 2.0 and 2 from `work.md` and run the contract suite.
**Why RED now:** Predicates naming Step 1 (16) and Step 2 (10) fail; the ledger contract has no Go test.
**GREEN when:** Suite passes without those lanes; `advance` tests cover default, targeted, fan-out and collision selection; `work.md` selection prose is one paragraph.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*
