---
id: REQ-509
title: '[impact-rule-change] Merge the Common Rationalizations tables into one crew member'
status: claimed
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
write_set:
  - skills/do-work/crew-members/shared-principles.md
  - skills/do-work/actions/work.md
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/capture.md
  - _dev/tests/contracts/core-checks.sh
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-04T22:05:56Z
  basis:
    - trivial short-circuit
exploration_at: 2026-09-04T22:05:56Z
preflight_at: 2026-09-04T22:07:51Z
route: B
claimed_at: 2026-09-04T22:05:00Z
---

# Merge the Common Rationalizations Tables Into One Crew Member

## What
One crew member becomes the loading point for shared principles currently spread across action-level `## Common Rationalizations` tables. Load it at implementation and review; each action keeps rows specific to its own step.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Declare five paths; write the loading/duplicate-row contract test first, observe RED, consolidate portable principles, then prove GREEN and preserved local boundaries.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The audit found zero near-identical cross-file row pairs, so duplication is not the reason for this change. One loading point makes the shared principles available consistently during implementation and review while action-specific guidance stays with its action.

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
**RED prompt/case:** Assert that implementation and review load the shared-principles crew member through the contract suite.
**Why RED now:** Neither loading point names a shared-principles crew member, and no predicate protects that loading contract.
**GREEN when:** Both loading points name the crew member and the contract suite proves that removing either load fails; action-specific rows remain local, and the audit's near-identical cross-file row count remains zero.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4244 tokens; action-loading and downstream-reader match; partial family coverage prevents narrowing within the 2000-token budget.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens; existing shell test owner; partial family coverage prevents narrowing. Both satellites were read as additive touch-conditional context by exploration and remain required for the builder.

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

## Triage

**Route B** — the outcome is clear, but portable principles must be distinguished from action-specific boundaries across the suite. Exploration completed before this prioritized claim; no source was edited.


## Plan

Planning not required — Route B.

## Exploration


- The REQ's What/Why/RED already reflect the accepted one-loading-point goal. REQ-508 exists in `do-work/archive/`; the orchestrator must use canonical selection/claim evidence for readiness.
- Read both listed primes, maintenance/general crew contracts, the lessons budget contract and index, and both action/shell lesson satellites. Claim-time indexed selection is empty: action-files costs 4244 tokens and shell-commands 3385, both `slugged: partial`, above the 2000-token budget. Refresh the stale 3436-token drop; shell lessons match only if the implementation adds shell test logic. Touch-conditional satellite reading is additive to the budget, so both were read for this preparation.
- Inspected all 23 Common Rationalizations tables across four packages. They contain 126 rows now. Most encode action-specific boundaries or defaults, not portable instructions. Do not transplant consent, status, command ownership, or initialization rules into a global contract: doing so would change other actions' behavior.
- **Audit contradiction reproduced:** the exact Finding 11 command prints `tables 23 rows 126 near_identical_cross_file_pairs 1`, not zero. Its only pair is capture versus validate-feedback, reason similarity 0.767857: “Capture ≠ Execute — the user decides when to run the queue” versus “…what becomes work.” The audit acknowledges their different boundaries while giving a contradictory zero count. Preserve both local rows and narrow capture's reason to its actual boundary (for example, “Capturing intent does not authorize processing the resulting queue.”). This satisfies literal zero without changing the threshold, exclusions, or either behavior. Record the discrepancy and resolution in the REQ before implementation.
- **Four-part constraint:** this REQ explicitly migrates judgment/loading, not mechanics into Go. A new CLI command would contradict its purpose and the “judgment stays prose” constraint. Record that the CLI-command part is not applicable; deleted prose plus new loading/negative mutation proof are the applicable migration parts. Current table-preservation predicates are also absent: `contract-regressions.sh` is a 77-line aggregate, real core checks live in `_dev/tests/contracts/core-checks.sh`. Do not delete unrelated negative sentinels in `defensive-surface-audit.sh` or the qualification/finding-closure predicates in core-checks. Document that no old individual Common Rationalizations predicates remain to delete.


The full exploration and source-to-shared mapping are in `do-work/runs/work-2026-09-05-005615/REQ-509-exploration.md`.

## Scope

**Files I will touch:**
- `skills/do-work/crew-members/shared-principles.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/capture.md`
- `_dev/tests/contracts/core-checks.sh`

**Acceptance criteria (restated from REQ):**
- [ ] One condition-keyed shared principles crew member owns portable guidance, with action-specific rows retained.
- [ ] Implementation and review actively load the shared file.
- [ ] One contract predicate fails when either loader or shared file is removed; captures RED before source edits.
- [ ] The audit comparison remains unchanged and reports zero cross-file near-identical action rows; a copied row makes it fail.
- [ ] Existing action contracts and shipped references remain valid.

## Decisions

- **D-01 — preserve judgment ownership:** This move changes prose loading, so the chain-wide CLI-command clause is inapplicable; adding a command would contradict the explicit judgment-stays-prose constraint. No individual-table predicates remain to delete. Use the existing sourced core-checks owner instead of inflating its aggregate wrapper.
- **D-02 — resolve measured evidence:** The audit reproduction currently returns one pair (capture/validate-feedback reasons, 0.767857). Preserve both local authorization boundaries while narrowing capture’s reason; retain the audit threshold and comparison algorithm. The original zero claim is contradictory evidence, not a reason to weaken the new check.

- **D-03 — checkout location:** Keep the isolated builder under the repository’s Git-private scratch directory to honor the user’s repository boundary. It stays outside tracked working-tree scans; the builder treats its stale do-work snapshot as absent. Use the required `codex/` branch prefix and retain both exact branch and checkout path for cleanup.

## Pre-flight

Canonical advance records for this exact request/path satisfied baseline and green-gate evidence. `bash _dev/tests/contracts/core-checks.sh` passed; direct `bash _dev/tests/maintainer-verify.sh` passed at 45d9d3dbc50255eaf363699829a9730ed5fb4f0a (board 14s; CLI 47s; slowest CLI file 19.57s). No failing baseline exclusions. Existing report work from another session was committed separately and is outside this request’s range.

## Dispatch

Use the exact five-path Scope. See REQ-509-exploration.md in this main-tree run for detailed source-to-shared mapping. Own branch codex/worktree-agent-REQ-509-shared-principles, checkout .git/work-run-20260905/worktree-agent-REQ-509-shared-principles; hand back only through REQ-509-handback.md.
