---
id: REQ-509
title: '[impact-rule-change] Merge the Common Rationalizations tables into one crew member'
status: completed
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
dispatch_at: 2026-09-04T22:08:22Z
builder_handback_at: 2026-09-04T22:14:10Z
integration_at: 2026-09-04T22:14:41Z
commit: 2ba5b432658853690e8e5a6d20bd2dcc147e9ada
review_at: 2026-09-04T22:26:20Z
kb_status: pending
route: B
claimed_at: 2026-09-04T22:05:00Z
completed_at: 2026-09-04T22:26:40Z
release_at: 2026-09-04T22:26:40Z
---

# Merge the Common Rationalizations Tables Into One Crew Member

## What
One crew member becomes the loading point for shared principles currently spread across action-level `## Common Rationalizations` tables. Load it at implementation and review; each action keeps rows specific to its own step.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Declare five paths; write the loading/duplicate-row contract test first, observe RED, consolidate portable principles, then prove GREEN and preserved local boundaries.
- [x] **[APPLY]:** Five declared paths merged from isolated branch; nine portable rows become eight condition-keyed principles and both active readers load them.
- [x] **[UNIFY]:** Reviewed all five paths; focused core-checks, aggregate contracts, shipped references, warning-level ShellCheck, shell syntax and diff checks passed on the builder branch; no debug artifacts.

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
- [x] One condition-keyed shared principles crew member owns portable guidance, with action-specific rows retained.
- [x] Implementation and review actively load the shared file.
- [x] One contract predicate fails when either loader or shared file is removed; captures RED before source edits.
- [x] The audit comparison remains unchanged and reports zero cross-file near-identical action rows; a copied row makes it fail.
- [x] Existing action contracts and shipped references remain valid.

## Decisions

- **D-01 — preserve judgment ownership:** This move changes prose loading, so the chain-wide CLI-command clause is inapplicable; adding a command would contradict the explicit judgment-stays-prose constraint. No individual-table predicates remain to delete. Use the existing sourced core-checks owner instead of inflating its aggregate wrapper.
- **D-02 — resolve measured evidence:** The audit reproduction currently returns one pair (capture/validate-feedback reasons, 0.767857). Preserve both local authorization boundaries while narrowing capture’s reason; retain the audit threshold and comparison algorithm. The original zero claim is contradictory evidence, not a reason to weaken the new check.

- **D-03 — checkout location:** Keep the isolated builder under the repository’s Git-private scratch directory to honor the user’s repository boundary. It stays outside tracked working-tree scans; the builder treats its stale do-work snapshot as absent. Use the required `codex/` branch prefix and retain both exact branch and checkout path for cleanup.

## Pre-flight

Canonical advance records for this exact request/path satisfied baseline and green-gate evidence. `bash _dev/tests/contracts/core-checks.sh` passed; direct `bash _dev/tests/maintainer-verify.sh` passed at 45d9d3dbc50255eaf363699829a9730ed5fb4f0a (board 14s; CLI 47s; slowest CLI file 19.57s). No failing baseline exclusions. Existing report work from another session was committed separately and is outside this request’s range.

## Implementation Summary

- `skills/do-work/crew-members/shared-principles.md` (new): eight condition-keyed principles consolidate nine portable work/review rows. Detailed testing and discovery rules link to their existing crew owners. Scope applies only when the request requires a declaration; retry loading applies to the second test-failure repair attempt or later.
- `skills/do-work/actions/work.md` (modified): Step 6 always loads shared principles; removes five portable rows and points its retained lifecycle-specific table to the shared owner.
- `skills/do-work/actions/review-work.md` (modified): Step 2 reads the shared crew in both review modes; removes four portable rows while retaining score, findings and stale-restatement boundaries.
- `skills/do-work/actions/capture.md` (modified): narrows the queue-processing authorization reason to its capture-specific boundary. Its STOP instruction is unchanged.
- `_dev/tests/contracts/core-checks.sh` (modified): one merged contract checks shared-file existence, active implementation/review loads, and the audit's unchanged cross-file trigger-or-reason similarity rule. In-memory mutations remove, comment, negate and displace each loader, remove the shared file, and copy a row between files. No old predicates were deleted because individual Common Rationalizations table predicates are absent from the current owner.


## Implementation Decisions

- **D-04 — one source predicate with small mutation fixtures (DECIDE & STATE):** Use the existing core-checks owner, with in-memory mutated readers and the original audit parser/comparison. Read actual suite files for the live check; mutation cases use only the rows needed to test their seam. This avoids repeated whole-suite similarity scans and leaves no fixture artifacts. Match active load directives inside implementation Step 6 and review Step 2 after stripping comments/fenced examples; do not let a mention elsewhere satisfy the load.
- **D-05 — preserve source conditions while consolidating (DECIDE & STATE):** The orchestrator's early read identified that Route A skips formal Scope. Key the shared declaration rule on requests requiring Scope, preserving that exception. Similarly preserve the retry rule's test-failure trigger instead of extending it to every kind of repair. Keep action-specific authorization/lifecycle/scoring rules local; the shared crew links detailed existing owners instead of creating a second authority.

D-01 and D-02 from the brief remain applicable: this request moves judgment, so no CLI command is appropriate; no current individual-table predicate exists to remove. The preexisting audit's contradictory zero claim is resolved by narrowing capture's reason, not by changing the comparison.


## Discovered Tasks

None.

## Qualification

Passed — five declared implementation files match the exact merge range `622a5e55de332984d7e180615a2c5c2b6c7ef2d7..2ba5b432658853690e8e5a6d20bd2dcc147e9ada`; canonical qualify and scope-drift records are satisfied. All five acceptance criteria trace to shared conditions, active loaders, mutation coverage and preserved local boundaries. The two reporter-output warnings are the contract suite’s explicit failure/pass diagnostics, not debug artifacts. P-A-U and test-first evidence were checked against the handback and diff.

## Testing

All commands ran from the isolated checkout.

| Command | Result | Duration |
|---|---|---|
| `bash _dev/tests/contracts/core-checks.sh` before source edits | **RED, exit 1**: shared file missing; work Step 6 active load missing; review Step 2 active load missing; `near_identical_cross_file_pairs 1` | 2.72s |
| Same command after implementation | **GREEN, exit 0**: `shared principles loads and mutations passed; near_identical_cross_file_pairs 0`; `core-checks contract probes passed.` | 2.84s |
| Same command at final commit after narrowing Scope/retry conditions | **GREEN, exit 0**, same output | 2.54s |
| `bash _dev/tests/contract-regressions.sh` | Exit 0, `Contract regression checks passed.` | 18.20s |
| `bash _dev/tests/shipped-package-reference-contract.sh` | Exit 0, `shipped package reference contract: PASS` | 0.99s |
| `bash -n _dev/tests/contracts/core-checks.sh` | Exit 0 | Below 1s |
| `shellcheck --severity=warning _dev/tests/contracts/core-checks.sh` | Exit 0 | Below 1s |

Aggregate test-file durations all remained below 30 seconds: core-checks 3s; maximum reported file duration 14s (session-start-hook-behavior). The aggregate normally skips staged-skills/updater/installer heavy probes; this is not heavy-lane evidence. The orchestrator will run the merged repository gate and choose affected heavy lanes canonically.

The exact Finding 11 Reproduce algorithm was also run unchanged after implementation: **`tables 23 rows 117 near_identical_cross_file_pairs 0`**. The test does not freeze counts or weaken its `> 0.75` trigger-or-reason comparison. The added shared table uses `## Principles`, so it does not become an action rationalization table.

The negative mutations exercise the same predicate as the live sources and cannot pass on an intact pointer elsewhere in the action. They prove the loading seam and duplicate-row regression. They do not prove that an LLM will invariably follow prose; semantic preservation still needs independent review.


**Merged verification:** canonical advance scope/focused/green-gate records all satisfied. Direct `bash _dev/tests/maintainer-verify.sh` passed at merge 2ba5b432658853690e8e5a6d20bd2dcc147e9ada; board module 16s, CLI module 50s; slowest CLI file 20.72s, below 30s. Captured RED/GREEN ordering is valid.

**Priority scheduling:** The user repeatedly prioritized running REQ-509. Its selected heavy lanes run immediately at this merged revision before review/finalization, so this priority request does not wait for the rest of the queue. No queued status is inferred from test output.

## Verification Incident

The initial selected-lane run at implementation revision 2ba5b432658853690e8e5a6d20bd2dcc147e9ada passed JavaScript (5s), browser (84s), CLI integrations (52s), staged skills (23s), and updater (54s), but installer exited 2 after 6s. All lanes executed; none skipped. A diagnostic replay revealed `skills/do-work/actions/version.md` still declared 0.281.0 while root/core declared 0.282.0. This was an omission in the orchestrator’s earlier release manifest, outside this request’s source diff.

Canonical `release --manifest … --commit` repaired the omitted 0.282.0 mirror and its changelog note at 2a5a698b0e2de7723d87183b74ed7bded5dadb79. The manifest authoring helper now includes that required mirror on every future release. The installer-only retry uses this repaired execution revision; the original failure remains in the saved evidence.

## Heavy Verification Plan

- Base revision: `622a5e55de332984d7e180615a2c5c2b6c7ef2d7`
- Target revision: `2ba5b432658853690e8e5a6d20bd2dcc147e9ada`
- `queue-kanban-javascript`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests.
- `queue-kanban-browser`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests.
- `do-work-cli-integrations`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests.
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests; skills/do-work/actions/capture.md matched subtree skills; skills/do-work/actions/review-work.md matched subtree skills; skills/do-work/actions/work.md matched subtree skills; skills/do-work/crew-members/shared-principles.md matched subtree skills.
- `updater`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests.
- `installer`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — _dev/tests/contracts/core-checks.sh matched subtree _dev/tests.

## Heavy Verification Result

Initial execution: `2ba5b432658853690e8e5a6d20bd2dcc147e9ada`. Every lane executed, with no skips.

| Lane | Exit | Seconds | Disposition |
|---|---:|---:|---|
| queue-kanban-javascript | 0 | 5 | executed |
| queue-kanban-browser | 0 | 84 | executed |
| do-work-cli-integrations | 0 | 52 | executed |
| staged-skills | 0 | 23 | executed |
| updater | 0 | 54 | executed |
| installer | 2 | 6 | executed |

Installer-only retry after the documented release-mirror repair: execution `2a5a698b0e2de7723d87183b74ed7bded5dadb79`, exit 0, 30 seconds, executed, unskipped. No REQ-509 implementation source changed. Passing evidence is complete across both explicitly recorded execution revisions; this is not an all-green six-lane run at the original merge.

## Review

**Approve** — the shared principles reach implementation and both review modes, with the original conditions and local action boundaries preserved. No actionable findings.
Route B | `2ba5b432658853690e8e5a6d20bd2dcc147e9ada`

### What's built

One shared crew file supplies eight condition-keyed principles to implementation and both review modes. Nine portable work/review rows move there; action-specific lifecycle, authorization and scoring boundaries remain local.

### Decisions / risks for you

None. The new predicate proves that the active readers receive the shared instructions; semantic review, rather than a prose-presence test, establishes that their meaning is preserved. It cannot guarantee an LLM follows every instruction.

### Findings

**Important:** None.
**Minor:** None.
**Nit:** None.

### Requirements Checklist

- [x] One loading point for shared principles, honoring the accepted addendum rather than the refuted duplication premise.
- [x] Portable source rows consolidated by condition. Completion-summary and independent claim verification combine without losing either obligation.
- [x] All 23 action/prompt tables inspected: remaining rows retain their action's specific operating boundary.
- [x] Work Step 6 and review Step 2 actively load the file. Review Step 2 serves both modes.
- [x] One merged contract proves the shared file and both loaders, with removed, commented, negated and displaced loader mutations, absent-file mutation and a copied-row mutation.
- [x] The exact audit reproduction reports zero near-identical cross-file pairs, preserving its trigger-or-reason comparison and threshold.
- [x] Five declared paths match the exact integration range; P-A-U is complete and significant decisions are documented.
- [x] Four-part migration constraint assessed: no CLI command is appropriate for this judgment/loading move, and no old individual-table predicates remain. D-01 records both exceptions; deleting unrelated checks or adding unused mechanics would contradict the request's purpose. Deleted local prose and new regression proof supply the applicable parts.

### Acceptance Testing

**Result: Pass.**

- Independently ran `bash _dev/tests/contracts/core-checks.sh`: exit 0, 3.28s; active loads and all mutation cases passed.
- Replayed that same predicate in memory against the base revision's affected source bytes: RED with missing shared file, missing implementation load, missing review load and one near-identical pair. Merged bytes were GREEN. Builder handback separately records actual test-first RED before source edits.
- Executed the audit's exact Finding 11 reproduction unchanged: `tables 23 rows 117 near_identical_cross_file_pairs 0`.
- Independently mapped every removed row to the shared condition. The Scope condition preserves Route A's exemption; retry loading remains restricted to test-failure repairs from the second attempt. Discovery, testing and debugging links retain their canonical owners.
- Restatement sweep covered shared-file readers, remaining rationalization tables, Scope/Route A, test-repair attempts, discovered-task routing, completion evidence and acceptance/Klarna guidance. No conflicting restatement introduced by this change was found. Capture's STOP instruction is unchanged; its narrower reason preserves the original authorization boundary.
- Orchestrator reports the direct repository gate passed at the reviewed merge revision: board 16s, CLI 50s, slowest CLI test file 20.72s. Canonical scope, focused and green-gate records are satisfied.

- Read the canonical heavy-run results: at implementation revision `2ba5b432658853690e8e5a6d20bd2dcc147e9ada`, JavaScript (5s), browser (84s, explicit Chromium configuration), CLI integrations (52s), staged skills (23s) and updater (54s) all exited 0 and executed without skips. Installer initially exited 2 after 6s.
- The installer diagnostic identified a prior release bookkeeping defect: `skills/do-work/actions/version.md` displayed 0.281.0 while the required version was 0.282.0. The separate repair commit `2a5a698b0e2de7723d87183b74ed7bded5dadb79` changes only that mirror and the two changelog notes. Independently confirmed the five REQ-509 implementation paths are byte-identical across these revisions. The canonical installer-only retry at the repaired revision exited 0 in 30s, executed without skips. This is passing evidence for every selected lane across the stated revisions, not an all-six-lane green run at the original revision. The initial failure remains in the record and is not attributed to REQ-509's implementation.

### Suggested Additional Testing

None for this five-file change. Existing focused, mutation, source-meaning and package-level checks cover the applicable acceptance boundaries.

### Scores

**Overall: 100%**

| Dimension | Score | Notes |
|---|---|---|
| Requirements | 100% | Accepted one-loading-point goal and constraints satisfied |
| Code Quality | 100% | Conditions and canonical owners preserved |
| Test Adequacy | 100% | Actual test-first evidence, independent historical replay and negative mutations |
| Scope | 100% | Exactly five declared paths; decisions recorded |
| Risk | None | No unresolved regression attributable to this change |
| Acceptance | Pass | Loading, semantics and selected lane evidence verified |

Overall is the arithmetic mean of the four percentage dimensions; no qualitative modifier applies.

### Follow-ups created

None (0 findings report only).

Self-validation complete: checked both review modes, preserved source conditions, original RED replay, local-rule retention and the distinction between loading evidence and semantic compliance. No review-owned process or subagent remains pending.

## Lessons Learned

**What worked:** A test-first contract exercised both active readers and negative mutations. Independent source-to-shared review preserved local exceptions and the original audit algorithm.

**What didn’t:** An unconditional extraction initially broadened Scope to Route A; narrowing the condition preserved the existing exception. Earlier orchestrator release manifests omitted the version action mirror, so installer validation failed until the canonical release repair aligned it.

**Worth knowing:** Prose extraction must preserve the original condition, not merely the instruction. Canonical release validation checks the declared mirrors; the action must declare every mirror the suite installer requires, including `actions/version.md`.

## Orientation

[MAP CHANGED] Implementation and review now share one condition-based guidance source in the crew subsystem. Action-specific lifecycle and authority rules remain local, while tests prove both readers load the common principles. Relevant action and shell prime paths still resolve.
