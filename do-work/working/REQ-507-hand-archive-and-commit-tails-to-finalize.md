---
id: REQ-507
title: '[impact-rule-change] Hand the archive and commit tails to finalize'
status: claimed
status_changed_at: 2026-09-04T20:57:34Z
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-506, REQ-584]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, _dev/tests/contracts/core-checks.sh, skills/do-work/tools/do-work-cli/prime-do-work-cli.md, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate.go, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate_test.go, skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go, skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go, skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go]
route: C
planning_at: 2026-09-04T18:26:45Z
exploration_at: 2026-09-04T18:32:49Z
dispatch_at: 2026-09-04T18:36:04Z
builder_handback_at: 2026-09-04T19:04:09Z
integration_at: 2026-09-04T19:04:37Z
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-09-04T18:22:01Z
  basis:
    - Route C
    - 8-file write set
    - 3 subsystems involved
    - 3 acceptance criteria
    - dependency depth 3
    - cross-route regression gates
    - full-suite verification
gate_deferred: 'true'
claimed_at: 2026-09-05T10:05:13Z
remediation_dispatch_at: 2026-09-05T10:10:08Z
remediation_at: 2026-09-05T10:15:14Z
---

# Hand the Archive and Commit Tails to finalize

## What
Step 8 (66 lines) and Step 9 (21 lines) reduce to: mint follow-ups by Fold-First (prose), then `advance` runs `finalize`. The Changelog Entry Procedure and the Commit and Metadata-Commit Procedure in `work-reference.md` leave prose except the changelog title and prose judgment.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Keep `finalization` as the transaction authority; make `advance` compose it from one request-bound manifest, consume typed finalization evidence in the work action, remove the displaced archive/commit recipes and predicates, and prove the four terminal paths through the public CLI seam.
- [x] **[APPLY]:** The isolated builder implemented the frozen 12-file plan with public RED-first coverage and committed it as 7faafb9d3fbf04b26a4599a5141b81f3b476c6fa.
- [x] **[UNIFY]:** The owner verified the clean builder branch, exact scoped diff, absence of builder-authored lifecycle files, formatting, tests, contracts, and the conflict-free no-ff integration.

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
- [x] Heavy lanes at `ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`: the work loop runs them at queue exhaustion and records the result here → Heavy verification failed: staged-skills exited 1. Return to remediation.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Triage

**Route: C — Complex.** This rule-changing migration removes two finalization procedures, binds advance to the existing journaled finalizer, and must preserve four distinct completion paths plus release and follow-up judgment boundaries.

**Planning:** Required.

## Plan

1. Add public RED coverage for serial, supplied-worktree, completed-with-issues, already-green/no-release, and refusal paths at `advance`.
2. Add strict request-bound finalization-manifest handling to `advance`, delegating the transaction to the existing `finalization` handler and preserving its typed result.
3. Reduce the work action and reference procedures to Fold-First, impact, release-content, terminal-state, lesson, and cleanup judgment; update the CLI prime to the same ownership boundary.
4. Replace retired prose predicates with structural ownership guards, then run focused, race, module, contract, vet, and repository-gate verification.

## Exploration

`lifecycleadvance` already imports `finalization`, so composing the terminal phase introduces no package cycle. The existing finalizer owns the required transaction and result fields, but safe outer request binding must happen during its single manifest decode; prechecking in `lifecycleadvance` and reopening the file would create a replacement race. Ordered finalization records already normalize in JSON, while the text renderer needs a production update to maintain typed text/JSON parity. The captured shell-test path is stale: active predicates now live in `_dev/tests/contracts/core-checks.sh`. None of the four required terminal paths is covered through public `advance`, and completed-with-issues lacks a full finalization transaction case.

## Scope

**Files I will touch:**

- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `_dev/tests/contracts/core-checks.sh`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate.go` (new)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`

**Acceptance:**

- A reviewed and oriented working request advances to a mechanical `finalize` phase, requires exactly one request-bound manifest, and returns the existing finalizer's full typed outcome, findings, changes, rollback, and ordered records.
- Public CLI tests prove serial, supplied-worktree, completed-with-issues, already-green/no-release, missing-input, hostile-input, and identity-mismatch behavior, with refusals producing no mutation.
- Step 8/9 and their reference procedures retain Fold-First, consolidation, impact, terminal/failure, release-content, lesson, and cleanup judgment but no longer teach archive, staging, commit, provenance, or verification mechanics.
- Structural contracts and the CLI prime enforce the new ownership boundary; focused, race, module, contract, vet, and direct repository gates pass.

## Decisions

- **D-01 — Expand the captured scope to the current owners.** Replace the stale contract dispatcher and lifecycle directory entries with the exact 12-file set above, including the finalizer's same-decode binding seam, result text renderer, and CLI prime. Value: fail-closed request identity and truthful public contracts. Risk: broader but auditable implementation surface.
- **D-02 — Preserve the finalizer as the sole transaction engine.** `advance` composes and projects it; it does not duplicate journal, archive, release, Git, provenance, or rollback logic.
- **D-03 — Bind before all observable preparation.** The selected request ID and path are compared immediately after the finalizer's single manifest decode, before index inspection, journal lookup, planning, or mutation; only that typed mismatch maps to the dedicated refusal.

## Prior Evidence (4 September, revision ad8bceb7 — superseded by drift)

The sections below were recorded for the merged range `8e3dbf01..ad8bceb7` before later REQs reworked its files. They are history, not current evidence: the saved-range resume proof on 5 September found drift, so qualification, testing and review are redone against current `main` below.

### Implementation Summary

- `_dev/tests/contracts/core-checks.sh` (modified) — enforces the concise judgment boundary and forbids restored finalization-tail recipes within the affected sections.
- `skills/do-work/actions/work-reference.md` (modified) — reduces changelog and commit procedures to release-content, provenance, and typed-result judgment.
- `skills/do-work/actions/work.md` (modified) — reduces Step 8 to Fold-First and terminal/release judgment and Step 9 to the advance continuation and typed-success condition.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified) — exposes bound in-process finalization while preserving the direct command contract.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modified) — checks outer request identity from the single decoded manifest before any preparation effect.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modified) — classifies oriented work as mechanical finalization and dispatches its input to the composition seam.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` (modified) — updates the public phase matrix and exact continuation expectation.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate.go` (added) — strictly parses one manifest and projects the canonical finalizer result with active advance identity.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate_test.go` (added) — proves four terminal paths and no-mutation refusal behavior through the public command.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified) — renders every ordered finalization record in text from the existing typed model.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified) — proves normalized text/JSON field parity and actual multi-record ordering.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — records advance as the request-bound finalization composer without moving judgment into the CLI.

### Discovered Tasks

None.

### Qualification

Typed qualification passed for exact range `8e3dbf01e0660424965d79acb2e386b6604e4780..ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`: all 12 implementation paths match the declared Scope, every summary entry matches the merged diff, and scope-drift reported no findings.

### Testing

- **Red-green validation:** Public finalization and phase-matrix tests first failed because `advance` still stopped at agent judgment and rejected the manifest input; after implementation, the four-path/refusal matrix passed in 4.005s and the final focused lifecycle/result rerun passed in 11.022s and 0.343s.
- **Focused merged-state gate:** the advance-owned probe ran lifecycleadvance, finalization, and resultmodel at the merged tree with exit 0, baseline state green, and diagnostic SHA-256 `e8a829932bae3fc53990317ac388f37816b6e8ddc32fa469c245b7cdac251a06`.
- **Race:** lifecycleadvance, finalization, and resultmodel passed under `go test -race` in 13.737s, 42.619s, and 1.515s.
- **Contracts:** the aggregate contract regression suite passed in 15.3s; every fast test file remained under 30s.
- **Repository gate:** `bash _dev/tests/maintainer-verify.sh` exited 0 at `02b5a2a3fe1831b8dc8088d8c80165617f0ec29f`, covering ShellCheck, gofmt, contracts, vet, 375 board tests, and 677 CLI tests; the slowest CLI test file was 19.28s.

### Heavy Verification Plan

- Base revision: `8e3dbf01e0660424965d79acb2e386b6604e4780`
- Target revision: `ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`
- `queue-kanban-javascript`: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for the changed core contract owner.
- `queue-kanban-browser`: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for the changed core contract owner.
- `do-work-cli-integrations`: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for the changed core contract owner.
- `staged-skills`: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for the changed core contract owner.
- `updater`: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for the changed core contract owner.
- `installer`: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for the changed core contract owner.


### Answer Notes

- 2026-09-04 - [ ] Heavy lanes at `ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`: the work loop runs them at queue exhaustion and records the result here: Heavy verification failed: staged-skills exited 1. Return to remediation.
> ```
> Exact-revision heavy verification via do-work clarify. Stored base, target, selected lanes, argv and coverage reasons were checked against the recomputed plan and matched. Execution revision: ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9. Each selected lane ran in the detached checkout; no results were borrowed from another revision.
> Heavy verification failed: staged-skills exited 1. Return to remediation.
> Failure evidence: FAIL: core runtime must resolve actions/work.md through sibling do-work-board. The staged-skills assertion checks for a do-work-board/ reference in skills/do-work/actions/work.md. The other five lanes passed.
> Chromium engine: /Applications/Google Chrome.app/Contents/MacOS/Google Chrome, explicitly supplied through QUEUE_KANBAN_BROWSER.
> Scope: record verification only; implementation fixes, fresh review and archiving are left to do-work run. Date and answer timestamp follow skills/do-work/actions/work-reference.md, Timestamp rule and its date-only paragraph.
> ```

### Heavy Verification Result

Target revision: `ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`
Execution revision: `ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`

- queue-kanban-javascript: exit 0, 8s — `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- queue-kanban-browser: exit 0, 99s — `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- do-work-cli-integrations: exit 0, 62s — `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: exit 1, 29s — `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: exit 0, 56s — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: exit 0, 38s — `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Repository Gate Deferral

- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** shellcheck:sc2148:do-work-runs-req-572-probe:missing-shebang
- **Repair dependency:** REQ-584
- **Diagnostic evidence:** "ShellCheck error SC2148 in do-work/runs/work-2026-09-04-232225/REQ-572-probe.sh line 1: Tips depend on target shell and yours is unknown. Add a shebang or a 'shell' directive."
- **Diagnostic evidence:** "Both direct canonical gate runs at 12d264c2 (detached worktree, clean tree) exited 1 on this one lint finding before any Go test ran. The probe file was committed by another session in 7ba3148a as REQ-572 run evidence and is outside REQ-507's implementation range 8e3dbf01..ad8bceb7, which touches no do-work/runs path."
- **Implementation base:** 8e3dbf01e0660424965d79acb2e386b6604e4780
- **Implementation merge:** ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9

## Saved-Range Resume Proof (2026-09-05T10:06:14Z)

- Base `8e3dbf01e0660424965d79acb2e386b6604e4780` and merge `ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9` both resolve; base != merge; base is an ancestor of merge; merge is an ancestor of `HEAD` (`1012e5e2`).
- Protected paths (`git diff --name-status -M base..merge`, no renames): 12 project files.
- **Drift:** 8 of the 12 have commit history after the merge — `work.md` (17 commits), `work-reference.md` (16), `prime-do-work-cli.md` (8), `finalization_commands.go` (4), `result_model.go` (3), `core-checks.sh` (2), `advance_commands.go` (1), `result_model_test.go` (1) — from REQ-504, 505, 506, 510, 515, 547, 562, 564, 569 and 570. None has a current staged, unstaged, untracked, deleted, type-changed or renamed state.
- **Result:** reuse rejected. The saved pair is deleted from the frontmatter, every prior qualification/testing/review verdict is treated as stale (kept above under Prior Evidence), and the request returns to Step 6 as a remediation pass: a builder re-verifies each acceptance criterion against current `main` and implements only what no longer holds (brief: `do-work/runs/work-2026-09-05-094707/REQ-507-brief.md`).

## Decisions (orchestrator, remediation pass)

- **D-07 — The merged range `8e3dbf01..ad8bceb7` stays this REQ's implementation.** The remediation builder re-verified all four acceptance criteria against current `main` with per-criterion evidence and returned no commits (hand-back: `do-work/runs/work-2026-09-05-094707/REQ-507-handback.md`, promoted into this record). The saved-range proof rejected *blind* reuse because later REQs reworked 8 of the 12 files; the verification pass is what replaces that reuse. Qualification therefore reads the real implementation range, every downstream verdict (tests, gate, review, heavy lanes) is renewed at the current tree, and provenance records the merge hash `ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`. Value: no fabricated re-implementation of work that is already on `main` and tested. Risk: a criterion weaker than its wording in a way neither tests nor predicates detect — mitigated by the builder tracing each criterion to live code and prose lines, and by the fresh independent review below.
- **D-08 (builder D-04) — Report no gap rather than manufacture one.** All four criteria trace to live artifacts and both focused suites pass; a cosmetic rewrite of `work.md`/`work-reference.md` would only add drift risk.
- **D-09 (builder D-05) — The red staged-skills lane was the assertion's fault, not this REQ's.** The lane asserted `work.md` must cite the sibling board package, an invariant this REQ deliberately ended; REQ-547 deleted that assertion (`c5dff3db`) with an attributing comment. The lane is re-run in this run's heavy drain, not assumed.
- **D-10 (builder D-06) — No new predicate for the CLI prime's finalization sentence.** No incident behind it; recorded as a report-only discovered task.

## Implementation Summary

Implementation landed on `main` in the merged range `8e3dbf01e0660424965d79acb2e386b6604e4780..ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9` (4 September); the 5 September remediation pass re-verified every acceptance criterion against the current tree and changed nothing.

- `_dev/tests/contracts/core-checks.sh` (modified) — pins the judgment-only shape of Step 8/9 and the two reference procedures and forbids restored finalization-tail recipes.
- `skills/do-work/actions/work-reference.md` (modified) — Changelog Entry and Commit & Metadata-Commit procedures reduced to release-content, provenance and typed-result judgment.
- `skills/do-work/actions/work.md` (modified) — Step 8 reduced to Fold-First, sweep, terminal, release/lesson and finalization-intent judgment; Step 9 to the the advance command continuation and its typed-success condition.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified) — exposes bound in-process finalization while preserving the direct command contract.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modified) — compares outer request identity inside the single manifest decode before any preparation effect.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modified) — classifies oriented work as the mechanical finalize phase and dispatches its input to the composition seam.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` (modified) — public phase matrix and exact continuation expectation.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate.go` (new) — parses exactly one manifest and projects the finalizer's result with the advance identity.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate_test.go` (new) — four terminal paths and seven no-mutation refusals through the public advance command.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified) — renders every ordered finalization record in text.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified) — text/JSON parity and multi-record ordering.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — records the advance command as the request-bound finalization composer.

**Remediation pass (5 September):** builder branch `worktree-agent-REQ-507-hand-archive-and-commit-tails-to-finalize` at `1012e5e2`, no commits, `git status` empty; per-criterion evidence C-1 to C-4 in the hand-back; focused suites `go test ./internal/lifecycleadvance ./internal/finalization ./internal/resultmodel` exit 0 and `_dev/tests/contracts/core-checks.sh` exit 0. Worktree and branch removed without force.

## Discovered Tasks

- `impact-negligible` — `_dev/tests/contracts/core-checks.sh` pins the CLI prime's evidence-gate sentence but not the adjacent "request-bound composition of finalization" sentence. → report only
- `impact-negligible` — `_dev/tests/staged-skills-contract.sh` gates `assert_core_sibling_reference` on `[ -f ... ]`, so a renamed or removed action file silently drops its sibling assertion. Not introduced by this REQ. → report only

## Qualification

Typed qualification (advance `qualify` + `scope-drift`, merged range `8e3dbf01e0660424965d79acb2e386b6604e4780..ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`, run 5 September at `5cb094dd`): satisfied. All 12 summary paths match the range and the declared Scope; no undeclared touch, no unused declaration. Two `QUALIFY-NEW-FILE-UNWIRED` warnings on `finalization_gate.go` and `finalization_gate_test.go` are judged not dead code: both are files of the existing `lifecycleadvance` package, and `advance_commands.go` calls `executeAdvanceFinalization` from the first while the second is its `go test` file — Go wires package files by membership, which a filename grep cannot see. Orchestrator judgment: the diff is substantive (four terminal-path tests plus seven refusal tests exercise the new seam through the public command), every requirement traces to a criterion in the builder's C-1 to C-4 evidence, the live data flow (advance → FinalizeBound → ordered records → text/JSON renderers) was read at the current tree, and no debug artifact or unchecked P-A-U box remains.

## Testing

**Tests run:** `bash -c 'cd skills/do-work/tools/do-work-cli && go test -count=1 ./internal/lifecycleadvance ./internal/finalization ./internal/resultmodel'` (preflight baseline at `1012e5e2` and focused probe at `e6005c0e`, both `GOMAXPROCS=2`); `bash _dev/tests/contracts/core-checks.sh` (builder, worktree at `1012e5e2`).
**Result:** ✓ All passing — lifecycleadvance 25.3s, finalization 67.2s, resultmodel 0.4s (builder run); probe exit 0 recorded as the request-bound `run-blocked-check` record; core-checks exit 0.

**Red-green validation:** traced to the captured `## Red-Green Proof`. RED (4 September, builder): public finalization and phase-matrix tests failed while `advance` still stopped at agent judgment and rejected the manifest input. GREEN (4 September and re-run 5 September): `TestAdvanceFinalizationRunsTerminalPathMatrix` (serial, supplied worktree, completed-with-issues, already-green/no-release) and `TestAdvanceFinalizationRequiresOneBoundManifestWithoutMutation` (missing, duplicate, empty, hostile-token, outer-path, id and path mismatch — each asserting an unchanged tree digest and `HEAD`) pass; the contract suite passes with the Step 8/9 mechanical prose gone. No new RED was owed by the 5 September pass because it changed nothing.

**Repository gate:** `bash _dev/tests/maintainer-verify.sh` run directly and unpiped from the project root, `GOMAXPROCS=2`: pre-build exit 0 at `1012e5e2` (recorded green-gate); final exit 0, started at `5cb094dd` and finished after another session had committed `985fa736` and `e6005c0e` — `git diff --stat 5cb094dd e6005c0e` touches only two `do-work/working/` records, no project byte — so the status was returned to `advance` and recorded with `recorded_revision` `e6005c0e1a15cfae32a3016f250e7d540300a722`. Slowest CLI test file `internal/publication/defer_gate_test.go` 28.55s on the pre-build run (under the 30s budget). Neither run needed the one retry.

**Heavy verification plan:** (typed planner, `_dev/tests/heavy-lanes.json`)
- Range: `8e3dbf01e0660424965d79acb2e386b6604e4780`..`ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9`
- `queue-kanban-javascript`: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for: _dev/tests/contracts/core-checks.sh
- `queue-kanban-browser`: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for: _dev/tests/contracts/core-checks.sh
- `do-work-cli-integrations`: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for: _dev/tests/contracts/core-checks.sh
- `staged-skills`: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for: _dev/tests/contracts/core-checks.sh
- `updater`: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for: _dev/tests/contracts/core-checks.sh
- `installer`: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for: _dev/tests/contracts/core-checks.sh

*Verified by work action*
