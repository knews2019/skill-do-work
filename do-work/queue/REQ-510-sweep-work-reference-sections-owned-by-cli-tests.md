---
id: REQ-510
title: '[impact-rule-change] Sweep work-reference sections whose contract is now a CLI behavior test'
status: pending
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-509]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/work-reference.md, skills/do-work/actions/work.md, _dev/tests/contract-regressions.sh, skills/do-work/docs/]
---

# Sweep work-reference Sections Whose Contract Is Now a CLI Behavior Test

## What
Delete every `work-reference.md` section whose contract moved into a Go behavior test during REQ-503 to REQ-509, and fix the cross-references in `work.md` and docs. Keep the Execution Model, the schema read contract, the Fold-First and Ratchet homes, and the minimal templates.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
`work-reference.md` is 1,250 lines with 66 inbound references from `work.md`; after the chain, many sections restate what a command now guarantees.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- For each candidate section, name the Go test that now owns the contract before deleting; a section with no owning test stays.
- Fix every inbound reference in the same commit; `shipped-package-reference-contract.sh` must pass.
- Delete the matching predicates.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-509; last in the chain.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Delete the Commit and Metadata-Commit Procedure section and run the contract and shipped-reference suites.
**Why RED now:** Predicates and reference checks fail; inbound links dangle.
**GREEN when:** Both suites pass; `work-reference.md` is under 700 lines; every remaining section is judgment, schema, or a minimal template.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Addendum (2026-09-03, 23:45 local)

User added (23:35 local, "yes, on the nine cancelation", applying the velocity report's triage table):

- REQ-471 (flow and reader consistency plus documentation for the gate-blocked set-aside) was cancelled into this sweep. When deleting `work-reference.md` sections here, also sweep any surviving sentence that says an unrelated canonical-gate failure must preserve a claim and stop the session; the shipped behaviour is the deferral lifecycle (REQ-491 to REQ-494) plus the retry-once rule REQ-559 adds, and queue summaries must distinguish blocked work, pending user decisions and dependency-gated work. One sweep, no new section.
- Coherence check: no contradiction; this widens the sweep's search condition by one sentence family.
