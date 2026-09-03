---
id: REQ-509
title: '[impact-rule-change] Merge the Common Rationalizations tables into one crew member'
status: pending
priority: now
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

## Addendum (2026-09-03)

User added (via the maintainability audit `do-work/audits/audit-2026-09-03.md`, Finding 11, sweep_key `rationalization-tables-not-duplicated`; the maintainer said "capture the requests" over the audit's §Plan, which carried this line):

> ```
> capture-request: [audit-2026-09-03 R11 · sweep_key: rationalization-tables-not-duplicated · JUDGMENT · addendum to REQ-509] Before REQ-509 runs, rewrite its Why clause and RED case: the 23 Common Rationalizations tables hold 125 rows with zero near-identical cross-file pairs (Finding 11's Reproduce), so "largely repeat each other" is refuted; either restate the goal as one loading point with a RED case that asserts loading, or cancel. Lock-in: near-identical cross-file rows = 0.
> ```

- Measured at dc8a64e3: the 23 `## Common Rationalizations` tables under `skills/` hold 125 rows and zero near-identical cross-file pairs (similarity above 0.75 on the trigger or the reason cell); the highest pair scores 0.77 and guards two different boundaries (capture stops before running the queue; validate-feedback stops before capturing). Reproduce: the Finding 11 command in the audit report.
- The Why clause "largely repeat each other" is therefore not supported at the row level, and the current RED case (remove `capture.md`'s table, a predicate fails) passes before any merge work, so it proves nothing about repetition.
- Resolved conflict: "the tables largely repeat each other" → the measured goal is one loading point for the tables' principles, not deduplication. Before this REQ is claimed, its What/Why are read as "one crew member is the loading point; each action keeps only rows unique to its own step", and its RED case becomes: `work.md` Step 6 does not name the merged file and no predicate pins it.
- Lock-in to carry: the Finding 11 Reproduce command keeps printing `near_identical_cross_file_pairs 0`; red the moment a row is copied between two action files.
- [x] Keep REQ-509 as restated above, or cancel it because one loading point is not wanted? → Keep, restated as one loading point; the builder rewrites the Why clause and the RED case per the 2026-09-03 addendum before claiming (do-work verify-requests on UR-105, maintainer applied the recommended fix)
  Recommended: keep, restated (the maintainer approved the audit's plan line, which offered both).
  Also: cancel via `do-work abandon` and recapture if a different goal emerges.


## Answer Notes

- 2026-09-03 - [ ] Keep REQ-509 as restated above, or cancel it because one loading point is not wanted?: Keep, restated as one loading point; the builder rewrites the Why clause and the RED case per the 2026-09-03 addendum before claiming (do-work verify-requests on UR-105, maintainer applied the recommended fix)
