---
id: REQ-507
title: '[impact-rule-change] Hand the archive and commit tails to finalize'
status: claimed
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-506]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/]
claimed_at: 2026-09-04T18:21:29Z
route: C
---

# Hand the Archive and Commit Tails to finalize

## What
Step 8 (66 lines) and Step 9 (21 lines) reduce to: mint follow-ups by Fold-First (prose), then `advance` runs `finalize`. The Changelog Entry Procedure and the Commit and Metadata-Commit Procedure in `work-reference.md` leave prose except the changelog title and prose judgment.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
REQ-498 made the tail a journaled CLI transaction; the prose still walks the pre-498 sequence.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- Fold-First minting, sweep consolidation and impact stamping stay prose (judgment).
- Archive, release payload validation, staging, commit, provenance and verification run inside `finalize` driven by `advance`.
- Delete the mechanical prose of Steps 8 and 9 and both reference procedures' mechanical parts, plus their predicates; add Go tests for serial, worktree, completed-with-issues and already-green paths if REQ-498 left any uncovered.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-506.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Delete the mechanical parts of Steps 8 and 9 and run the contract suite.
**Why RED now:** Predicates naming Steps 8 and 9 (19) fail; the changelog procedure has no behavior test beyond `release`.
**GREEN when:** Suite passes without those lanes; `finalize` tests cover the four paths; Step 8 prose is only Fold-First minting and Step 9 is one sentence.
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

**Route: C — Complex.** This rule-changing migration removes two finalization procedures, binds advance to the existing journaled finalizer, and must preserve four distinct completion paths plus release and follow-up judgment boundaries.

**Planning:** Required.
