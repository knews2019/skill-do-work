---
id: REQ-509
title: '[impact-rule-change] Merge the Common Rationalizations tables into one crew member'
status: pending
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-508]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-mechanical
write_set: [skills/do-work/crew-members/, skills/do-work/actions/, _dev/tests/contract-regressions.sh]
---

# Merge the Common Rationalizations Tables Into One Crew Member

## What
The ten `## Common Rationalizations` tables across action files merge into one crew member loaded at implementation and review; each action keeps only rows unique to its own step.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The tables are principle-shaped already and largely repeat each other; ten copies drift.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- One crew member with the union of rows, deduplicated, each keyed on the condition it guards.
- An action keeps a row only when it is specific to that action; the rest point at the crew member.
- Loading order in `work.md` Step 6 names the new file; delete predicates that pinned individual tables and add one predicate on the merged file.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-508.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Remove the table from `capture.md` and run the contract suite.
**Why RED now:** A predicate that quotes that table fails; no single file carries the union.
**GREEN when:** Suite passes with the merged file's predicate; ten tables become one plus action-specific rows.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*
