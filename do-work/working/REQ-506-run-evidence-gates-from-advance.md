---
id: REQ-506
title: '[impact-rule-change] Run the evidence gates from advance'
status: claimed
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-505]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - _dev/tests/contracts/core-checks.sh
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go
  - skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
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
planning_at: 2026-09-04T17:36:43Z
exploration_at: 2026-09-04T17:41:25Z
dispatch_at: 2026-09-04T17:44:35Z
builder_handback_at: 2026-09-04T18:15:00Z
integration_at: 2026-09-04T18:15:22Z
status_changed_at: 2026-09-04T20:57:34Z
claimed_at: 2026-09-04T23:40:10Z
---

# Run the Evidence Gates From advance

## What
Steps 3.6 (estimate), 5.75 (pre-flight), 6.3 (qualify) and the mechanical half of 6.5 (test gate and baseline comparison) run from `advance`; the Qualification Anti-Rationalization Table and the Finding-Closure Ratchet stay as principles.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Compose current-phase typed command handlers through advance, bind subordinate evidence to the exact REQ and path, retain judgment gates in prose, and replace obsolete sentence predicates with public behavior tests.
- [x] **[APPLY]:** Implemented request-bound estimate, preflight, qualification, focused-test, and green-gate execution through existing typed handlers, plus the planned prose and contract migration.
- [x] **[UNIFY]:** Reviewed all 16 changed files, ran gofmt, focused and race tests, module tests, vet, contract regressions, and the direct maintainer gate; no debug artifacts or undeclared touches remained.

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
- `_dev/primes/lessons-action-files.md` — 4362 tokens; matched current pipeline/argv/evidence owner; partial family means no targeted form, exceeds shared 2000-token budget. Whole touch-conditional read remains additive.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens; matched current pipeline/argv/evidence owner; partial family means no targeted form, exceeds shared 2000-token budget. Whole touch-conditional read remains additive.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 7300 tokens; matched current pipeline/argv/evidence owner; partial family means no targeted form, exceeds shared 2000-token budget. Whole touch-conditional read remains additive.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Triage

**Route: C — Complex.** This rule-changing migration composes estimation, preflight, qualification, and repository-gate evidence into the lifecycle command while removing procedural action contracts and replacing them with public Go behavior tests.

**Planning:** Required.

## Plan

1. Add ordered typed evidence-gate records to the advance result, carrying exact request identity, path, evidence provenance, gate state, subordinate outcome, findings, changes, and tokenized next/verification commands.
2. Compose the current mechanical phase through existing typed command handlers rather than duplicate estimate, preflight, qualification, or baseline logic; reject missing or phase-irrelevant inputs and bind every result to the discovered request.
3. Collapse the four procedural action surfaces to an advance consumer loop while retaining qualification skepticism, finding closure, retry/failure judgment, TDD proof, heavy-lane planning, and repository-gate attribution in prose.
4. Replace stale sentence predicates with public RED/GREEN command behavior and run focused, race, module, contract, and canonical merged-state gates.

**Plan validation:** All three detailed requirements map to the four tasks, and every task maps back to the evidence-authority migration. The typed consumer contract explicitly carries per-record identity, provenance, state, outcome, findings, changes, and tokenized continuation. Scope must account for the post-capture contract-suite split and result projection; exploration will resolve those exact paths before dispatch.

## Exploration

The existing core-helper handler map can supply estimate, preflight, qualification, scope-drift, and blocked-probe results without a package cycle; the gate-evidence handler map can supply exact green-record checks. This matches recovery's existing in-process command-composition pattern and avoids exported one-line delegates. Preflight and qualification already return typed findings and changes, while estimator output needs structured projection rather than its compatibility text alone.

The current blocked probe is not yet a complete focused-test baseline API: it reports status but discards bounded diagnostics, and the baseline record is private to core helpers. The requested mechanical baseline comparison therefore needs the declared core-helper and platform probe files so launch state, timeout, status, and normalized diagnostic identity can be compared without turning semantic similarity judgment into code. The captured contract path is stale after the suite split; `core-checks.sh` is the current predicate owner, while the root dispatcher remains untouched.

*Generated by Explore agent.*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `_dev/tests/contracts/core-checks.sh`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`

**Files I will NOT touch:** the root contract dispatcher, gate-evidence implementation, finalization, heavy verification, publication/deferral, queue selection-and-claim transaction, board code, or unrelated action readers.

**Acceptance criteria:**
- [ ] Advance executes each current mechanical evidence gate and returns exact REQ/path-bound typed outcome, provenance, findings, changes, and tokenized follow-up evidence.
- [ ] Focused-test execution distinguishes green, matching baseline red, new red, timeout, launch failure, and an unusable `launched: false` baseline using bounded diagnostic identity.
- [ ] Qualification uses an exact resolvable merge range, green-gate evidence distinguishes reusable records from a required direct run, and hostile or irrelevant phase input fails closed without interpolation.
- [ ] Estimate, test-command, warning, TDD, retry, attribution/deferral, heavy-lane, and finding-closure judgment stays in prose while four duplicated procedural surfaces are removed.
- [ ] Public Go tests replace the retired sentence predicates and all focused, race, module, contract, and canonical repository gates pass.

## Decisions

- **D-01 — DECIDE & STATE:** Expanded the captured scope to the current contract owner, typed result projection, core-helper baseline record, and platform probe implementation. The capture preceded the contract-suite split and assumed probe diagnostics already existed; these files are required to meet the request's explicit typed evidence and baseline-comparison behavior without duplicating command logic.
- **D-02 — DECIDE & STATE:** Added the CLI prime after focused GREEN because its lifecycle-advance phase table would otherwise keep teaching the retired read-only/next-command behavior. The prime is the required source index for every future change in this subsystem, so leaving it stale would split the contract immediately.
- **D-03 — DECIDE & STATE:** Separate canonical repository-gate tokens with repeated gate-argument flags while phase argv follows the command separator. This preserves both token channels without shell interpolation.
- **D-04 — DECIDE & STATE:** Treat a red as baseline-matching only when launch succeeded and command text, non-zero status, and bounded normalized diagnostic identity all agree. Uncertain semantic similarity remains action judgment and is conservatively reported as new red.
- **D-05 — DECIDE & STATE:** Capture Unix probe output through a parent-owned bounded pipe. This retains diagnostics without letting a background descendant keep the direct child wait open.
- **D-06 — DECIDE & STATE:** Point every mechanical continuation back through advance with explicit placeholders for judgment-owned inputs. The action no longer invokes subordinate helpers directly.
- **D-07 — DECIDE & STATE:** Retain retry limits, warning interpretation, TDD validity, gate attribution/deferral, heavy-lane planning, and finding closure in prose. Only deterministic command execution and exact evidence comparison moved into Go.

## Prior Implementation Summary

**Files changed:**

- `_dev/tests/contracts/core-checks.sh` (modified) — rejects restored evidence recipes and requires the retained judgment boundaries.
- `skills/do-work/actions/work-reference.md` (modified) — aligns architecture and gate evidence flow with request-bound advance execution.
- `skills/do-work/actions/work.md` (modified) — replaces four mechanical procedures with one evidence-gate consumer loop.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (modified) — returns focused-test evidence and conservative saved-baseline comparison.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go` (modified) — covers diagnostic comparison and unusable baseline refusal.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modified) — routes current mechanical phases through request-bound advance execution.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` (modified) — updates the public phase matrix and continuation contract.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go` (new) — composes typed handlers with exact phase, identity, input, and range binding.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go` (new) — covers all evidence gates, baseline states, invalid inputs, and direct green recording.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go` (modified) — exposes bounded normalized probe diagnostics without changing selection semantics.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` (modified) — verifies root normalization, truncation, and stable diagnostic identity.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go` (modified) — captures output through an owned pipe while preserving descendant cleanup.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go` (modified) — carries the diagnostic writer through the Windows probe path.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified) — adds normalized typed gate and focused-test records with compact rendering.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified) — covers normalization and rendering of nested evidence records.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — records advance's new request-bound evidence execution authority.

**What was done:** Advance now discovers and executes the current request's mechanical evidence boundary through existing command owners, returning ordered records bound to the exact request and path. Focused tests preserve bounded diagnostic identity and distinguish green, matching baseline red, new red, timeout, launch failure, missing baseline, and an unusable unlaunched baseline. Work prose now supplies inputs and semantic judgment through one named consumer loop.

## Discovered Tasks

None.

## Prior Qualification

**Passed.** The exact merge range contains 16 substantive implementation files and no undeclared path. All five acceptance criteria trace to public lifecycle gate behavior, bounded probe evidence, retained judgment prose, and replacement structural contracts; the checked P-A-U state matches the merged diff. Two declared core-helper check files were not needed because their existing handler logic was reused unchanged. New Go files are convention-discovered, and the relocated output warning is the retained anti-rationalization example rather than runtime debug instrumentation.

## Prior Testing

**Red-green validation:** The public advance estimate-gate case first failed with `ADVANCE-USAGE` because phase inputs after the command separator were not accepted. The identical focused test passed after request-bound gate composition was implemented, and the final public file covers every evidence phase, exact and invalid ranges, baseline states, hostile inputs, timeout, launch failure, and direct green recording.

**Merged-state verification:** Focused lifecycle, core-helper, selector, and result packages passed in 24.69s; the same packages under the race detector passed in 26.88s; static analysis passed in 0.31s; Windows selector-probe compilation passed in 1.33s; the full Go module passed in 48.15s; and contract regressions passed in 16.26s. The direct, unpiped canonical gate passed with 375 board tests and 674 CLI tests, with every measured test file below 30 seconds. Green-gate evidence was recorded at revision `5e767e5bc2b6edeb9a7c2b78589a5275f2e18f4b`.

## Prior Heavy Verification Plan

**Base revision:** `24ed2fdda549a0759cdc571562c9b782bfeb6251`

**Target revision:** `06367337dd82d97416e0d9d37872cc35b56ae7bc`

The planner marked coverage uncertain because `_dev/tests/contracts/core-checks.sh` is outside every declared path rule, so it conservatively selected all lanes:

- **queue-kanban-javascript** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`; reason: uncovered contract path makes coverage uncertain.
- **queue-kanban-browser** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`; reason: uncovered contract path makes coverage uncertain.
- **do-work-cli-integrations** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`; reason: uncovered contract path makes coverage uncertain.
- **staged-skills** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`; reason: uncovered contract path makes coverage uncertain.
- **updater** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`; reason: uncovered contract path makes coverage uncertain.
- **installer** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`; reason: uncovered contract path makes coverage uncertain.

## Open Questions

- [x] Heavy lanes at `06367337dd82d97416e0d9d37872cc35b56ae7bc`: the work loop runs them at queue exhaustion and records the result here → Confirmed: All 6 selected heavy lanes passed without skips at 06367337dd82d97416e0d9d37872cc35b56ae7bc.

## Prior Orientation

[MAP CHANGED] The lifecycle now has one public request-bound evidence path: advance discovers the current phase, executes its existing mechanical handler, and returns ordered typed gate records; the work action supplies inputs and retains semantic judgment.


## Answer Notes

- 2026-09-04 - [ ] Heavy lanes at `06367337dd82d97416e0d9d37872cc35b56ae7bc`: the work loop runs them at queue exhaustion and records the result here: Confirmed: All 6 selected heavy lanes passed without skips at 06367337dd82d97416e0d9d37872cc35b56ae7bc.
> ```
> Exact-revision heavy verification via do-work clarify. Stored base, target, selected lanes, argv and coverage reasons were checked against the recomputed plan and matched. Execution revision: 06367337dd82d97416e0d9d37872cc35b56ae7bc. Each selected lane ran in the detached checkout; no results were borrowed from another revision.
> All 6 selected heavy lanes passed without skips at 06367337dd82d97416e0d9d37872cc35b56ae7bc.
> Chromium engine: /Applications/Google Chrome.app/Contents/MacOS/Google Chrome, explicitly supplied through QUEUE_KANBAN_BROWSER.
> Scope: record verification only; implementation fixes, fresh review and archiving are left to do-work run. Date and answer timestamp follow skills/do-work/actions/work-reference.md, Timestamp rule and its date-only paragraph.
> ```

## Prior Heavy Verification Result

Target revision: `06367337dd82d97416e0d9d37872cc35b56ae7bc`
Execution revision: `06367337dd82d97416e0d9d37872cc35b56ae7bc`

- queue-kanban-javascript: exit 0, 8s — `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- queue-kanban-browser: exit 0, 96s — `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- do-work-cli-integrations: exit 0, 63s — `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: exit 0, 45s — `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: exit 0, 61s — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: exit 0, 26s — `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Prior Review

### Review: REQ-506

**Request changes — a failed or timed-out focused test can be reported as satisfied and make `advance` exit successfully.** The gate composition works on its ordinary paths, but its baseline override is not safe enough to authorize completion evidence.

Route C | saved implementation `06367337dd82d97416e0d9d37872cc35b56ae7bc` | reviewed range `24ed2fdda549a0759cdc571562c9b782bfeb6251..06367337dd82d97416e0d9d37872cc35b56ae7bc`.

### What's built

`advance` composes estimate, preflight, qualification/scope, focused-test comparison, and green-gate handlers into ordered request/path-bound records. Four procedural sections were replaced by a common consumer loop; qualification skepticism, finding closure, TDD, retries, and repository-gate attribution remain judgment. All 16 substantive changed paths match the Implementation Summary and declared scope. The two declared `corehelpers/checks*` paths were left unchanged because their existing handlers are reused, as Qualification already explains.

### Findings

**Important:**

- **F01 — Failed execution can clear the focused-test boundary.** `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go:169` unconditionally replaces a failed gate with `satisfied` when comparison reports `matching_red`. The comparison in `internal/corehelpers/commands.go` checks the saved baseline's launch state but does not require the current execution to have launched successfully or completed without timeout. `internal/nextselection/blocked_probe.go:107` additionally derives `launched`/`timed_out` from reserved numeric exit values, so an ordinary child exit 125/124 is indistinguishable from infrastructure failure/timeout. Public replay produced aggregate `outcome: success`, process exit 0, and a focused gate simultaneously containing `state: satisfied`, `outcome: failure`, `launched: false`, and `BLOCKED-PROBE-LAUNCH-FAILED`. A real timeout likewise exited successfully with `timed_out: true` and `matching_red`. A completion consumer is explicitly told that satisfied gates authorize durable evidence; this is a false success at that production boundary. Preserve actual process lifecycle facts and refuse baseline exclusion for unsuccessful launch, interruption, or timeout; never overwrite a failed subordinate result merely because diagnostic identity matches. — **impact-critical** → returned to the orchestrator for REQ-506's permitted remediation and, only if still unresolved afterward, critical Fold-First routing.
- **F02 — Some emitted missing-input continuations cannot be followed.** `internal/lifecycleadvance/evidence_gates.go:239` puts every missing input after `--`. At qualification, the emitted argv ends `--`, `<exact --diff-range <pre>..<merge>>`; substituting the requested range and flag yields `ADVANCE-GATE-INPUT-IRRELEVANT`, because qualification rejects a separator. The same template places missing `--gate-arg` input in the subordinate channel. The initial classifier's phase-specific argv is correct, but the missing-input record clears that continuation and substitutes the unusable one. Subordinate findings also retain direct helper remedies despite the action forbidding separate helper calls. Generate continuation tokens per phase/channel and preserve the advance boundary. — impact-user-visible → report only
- **F03 — A live reference still prescribes the retired qualification entry point.** `skills/do-work/actions/work-reference.md:502` (line 496 at the saved revision) says the merge range goes to `tools/checks/qualify.sh` through `DO_WORK_DIFF_RANGE`; the new consumer loop says not to call that helper separately and the composed gate requires `--diff-range`. This is an operative hand-back instruction, not a historical mention. Narrow it to the canonical advance input. — impact-rule-change → report only
- **F04 — The estimate short-circuit lost its operative condition during subtraction.** The removed Step 3.6 applied the floor estimate whenever effort was mechanical, or Route A had no heavy-evidence indicators. The replacement at `skills/do-work/actions/work.md:170` tells readers to extract nontrivial signals but omits that condition. `internal/lifecycleadvance/advance_commands.go:154` only emits `--trivial` when effort is mechanical **and** route is A. That narrower classifier condition predates this REQ; deleting the correct action instruction now leaves it as the only direct continuation, while `actions/estimate-reference.md:63` and `work-reference.md:164` still promise mechanical effort alone short-circuits. Preserve the original judgment condition through the new command boundary. — impact-rule-change → report only

**Minor / Nit:** None.

### Requirements checklist

- [x] Current-phase handlers run through `advance`, preserving exact discovered identity and ordered typed subordinate evidence.
- [ ] Missing-input findings provide usable advance continuations for every gate: F02.
- [ ] Focused-test evidence safely distinguishes launch failure, timeout, and baseline red: F01.
- [x] Qualification validates an exact resolvable range and scope comparison; hostile/mismatched/phase-irrelevant argv is rejected without interpolation.
- [x] Green-gate misses request a direct run; exact reported status is recorded through the existing owner.
- [x] Four old procedural headings and their work-action recipes are absent; Go public-command tests cover each replacement gate. At this saved base, `core-checks.sh` already contained no sentence predicates for those four gate procedures, so its actual delta adds removal guards rather than deleting nonexistent predicates. Retained scope checks exercise behavior and are not obsolete sentence predicates.
- [x] Qualification Anti-Rationalization Table and Finding-Closure Ratchet remain unchanged in this range; warning interpretation, TDD, retries, attribution, deferral, heavy testing, and substantive review remain prose.
- [ ] The floor-agent continuation and retained estimate-judgment contract are complete: F02–F04.
- [x] Scope, P-A-U completion, significant Decisions, four-part migration surfaces, and original UR-098 intent were checked against the diff. No unrelated finalization or queue-selection change was attributed to this REQ.

### Acceptance testing

**Result: Fail.** Independent tests ran in a detached checkout of the exact saved implementation. Existing focused packages passed uncached: lifecycleadvance 12.480s, corehelpers 14.765s, nextselection 3.541s, resultmodel 0.615s. These include the built public CLI's estimate, preflight, qualification, baseline-state, hostile-input, launch-failure, timeout, and green-record cases. The existing exact-revision record separately establishes six heavy lanes passing without skips; none was rerun or borrowed from another revision during this review.

Additional real-binary checks exposed the missing combinations:

1. A claimed Route A fixture at `test-gate` had focused probe `exit 125`, a usable saved baseline `{"test_command":"exit 125","exit_status":125,"launched":true}`, and empty baseline diagnostics. Restricting PATH to a private directory containing only `git` caused the actual `sh` launch to fail. `/usr/bin/true` was run directly and its exact zero exit supplied as canonical-gate evidence. `advance REQ-713 --request-path do-work/working/REQ-713-fixture.md --gate-arg /usr/bin/true --gate-exit-status 0 -- --probe-file focused.sh --timeout-seconds 2` exited **0**, with the contradictory focused record described in F01. No fake success was supplied for the focused probe.
2. Direct baseline `/bin/sleep 2; exit 124` completed with exit 124 and empty diagnostics. The identical focused probe with a one-second timeout was terminated by the runner; using the existing matching green-gate record, aggregate advance again exited **0**, focused state `satisfied`, `timed_out: true`, baseline state `matching_red`.
3. Ordinary immediate `exit 124` and `exit 125` probes reported timeout and unlaunched state respectively, showing the lifecycle fact/status collision without any missing executable.
4. Following the qualification missing-input continuation after replacing its placeholder with the requested `--diff-range` tokens returned `ADVANCE-GATE-INPUT-IRRELEVANT` before range validity was considered.

Current-presence check: the three F01 owner files and the F02 continuation implementation are byte-identical between the saved commit and main at review time. The cited F03/F04 prose is also still present. Later REQ-507 finalization extensions and concurrent release work do not resolve these findings. Existing REQ-503 classifier issues were not re-reported as new REQ-506 findings.

### Suggested additional testing

Add public regression combinations for a usable matching baseline plus real current launch failure/timeout, and retain negative controls for a normally exiting child whose status happens to be 124 or 125. Replay every emitted missing-input continuation after substituting only judgment-owned placeholder values. Re-review the remediation at its exact integrated range and renew whatever test evidence its changed source requires.

### Scores

**Overall: 50%** — percentage average 80.42%; Critical risk caps at 60%, then Acceptance Fail caps at 50%.

| Dimension | Score | Basis |
|---|---:|---|
| Requirements | 66.67% | Six of nine checklist items fully delivered; three partial or failed |
| Code Quality | 75% | Existing owners reused; execution-state override breaks their failure authority |
| Test Adequacy | 80% | Strong ordinary public matrix, but missing baseline × infrastructure combinations |
| Scope | 100% | All changed paths declared; two untouched declarations justified |
| Risk | Critical | Failed test execution can authorize completion evidence |
| Acceptance | Fail | Real public command reports success for failed launch and timeout |

### Review provenance and hand-off

Orchestrated, read-only preparation while REQ-506 remained queued. Read the full REQ and UR-098, review-work rubric, shared principles, relevant primes and lessons, general/testing/coding/maintenance/security guidance, and anti-slop reporting guidance. Self-validation checked both the ordinary red control and actual infrastructure failure, rather than inferring failure from reserved exit codes alone. No queue records, source, lifecycle, or commits were changed by this review. The orchestrator must re-resolve canonical claim identity, saved commit, and `resume_phase` evidence before consuming this report.

**Follow-ups created:** None. F01 is handed back for the current request's one remediation attempt; three noncritical findings remain report only. All review-owned fixtures, binaries, and detached checkout were removed after evidence was captured; no child process remains pending.

## Remediation Plan

The independent acceptance Fail activates the single remediation attempt. The five production execution owners remain byte-identical to the saved implementation at claim time. All historical verification stays under Prior headings; it cannot clear changed source.

### REQ-506 — Single-remediation plan

Repair the focused-test false success first. Include F02's continuation repair because it lives in the same `lifecycleadvance/evidence_gates.go` owner and is required by the original promise that an agent can follow `advance` evidence. Keep F03's stale reference and F04's estimate condition as report-only findings; neither is necessary to correct this execution boundary, and they would expand the concrete remediation owner set.

This is read-only preparation. No source or lifecycle changes and no new tests were run. The original review is `REQ-506-review.md` in this run directory. At dispatch, the orchestrator must re-resolve the exact claimed REQ/path and revision, settle its new baseline, record this one attempt, and give the builder an isolated checkout.

### Exact builder boundary

Eight paths, all already in REQ-506's original declared scope:

1. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go`
2. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go`
3. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go`
4. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go`
5. `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go`
6. `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go`
7. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go`
8. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go`

No new result-model field or enum is needed: existing `exit_status`, `launched`, `timed_out`, `baseline_state`, record `state`/`outcome`, findings, and argv can express the corrected result. Preserve the raw-status compatibility wrappers used by queue selection. Do not change `next_selection.go`, lifecycle classification, finalization, publication/release, action prose, or the estimator. Do not refactor process-group teardown. If implementation proves another path necessary, report the exact reason before expanding this boundary.

### Task 1 — Prove the false success and continuation failure before source changes

Add focused regressions at the existing public CLI seam in `evidence_gates_test.go`, using a correctly claimed fixture with an Implementation Summary and Qualification, plus a valid directly executed canonical-gate zero so a missing green record cannot mask the focused verdict. Suggested names describe the contract rather than the implementation:

- `TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline`: use a valid launched baseline with probe `exit 125`, status 125, and empty diagnostic evidence. Run the built CLI with PATH containing only a fixture-owned `git` executable/symlink, so Git evidence still works and the real `sh` launch fails. Assert current `launched:false`, the launch failure finding, non-satisfied focused gate, and nonzero aggregate advance status. The unfixed revision reports aggregate success and `satisfied` despite subordinate `failure`.
- A subcase with direct baseline `/bin/sleep 2; exit 124` actually completing at status 124, then the identical focused command with one-second timeout. Assert `timed_out:true` and that timeout cannot become `matching_red`/`satisfied` or aggregate success. Use the existing platform gating and cleanup conventions; no externally timed signal is needed for this replay.
- `TestBlockedProbeEvidencePreservesOrdinaryReservedExitValues` in the probe tests: ordinary immediate `exit 124` and `exit 125` both launched, neither timed out, and their raw exit values remain intact. At the public focused seam, valid matching baselines for these ordinary child exits must still clear; they are not infrastructure failures merely because of their integer values.
- `TestAdvanceMissingInputContinuationsPreserveArgumentChannels`: request qualification with exact request path but no range, substitute a real resolvable range into the emitted continuation without moving flags or separators, and execute it. Do the analogous round trip for missing canonical-gate tokens while preflight/focused phase tokens are already present. The unfixed qualification continuation returns `ADVANCE-GATE-INPUT-IRRELEVANT`.

Run these tests before edits and retain the real failing assertions. Keep or extend the existing green, ordinary matching-red, new-red/status-mismatch, missing/malformed/unlaunched-baseline, hostile-token, and green-record controls. The fix must not turn every baseline red into a regression or every high exit status into an infrastructure failure.

### Task 2 — Carry observed execution facts through the existing owners

The three owners currently disagree:

- The platform `runOwnedProbe` knows whether `Start` happened and which completion/timer/signal branch won, but returns only `(int, error)`.
- `RunBlockedProbeEvidenceAtRoot` reconstructs launch and timeout from statuses 125/124, losing that distinction.
- `handleBlockedCheck` again classifies those integers, `compareFocusedBaseline` validates only the saved baseline's launch state, and `composeCoreGate` then promotes any matching-red result to satisfied even when the subordinate failed.

Return observed lifecycle facts from the private platform runner, preferably through the existing evidence type rather than a second public representation. Set launch facts at the actual launch boundary and timeout only in the timer branch; keep error/interruption identity and status intact. Populate normalized diagnostic data in the existing common wrapper. The Windows implementation remains a truthful unsupported launch, and `RunBlockedProbeAtRoot`/`RunBlockedProbe` keep returning the same raw status/error surface to existing queue consumers.

In `handleBlockedCheck`, classify execution from those facts plus the actual runner error, not reserved integer equality. Only a successfully launched, normally finished execution with no runner error and no timeout may be compared for green/matching-red exclusion. The existing `not_compared` state and execution finding can express an execution that never became eligible; do not invent a baseline value that suggests the saved record itself was invalid. Preserve saved-baseline parsing and bounded diagnostic equality for executions that are eligible.

In `composeCoreGate`, preserve failed subordinate authority. Satisfy the gate only for a valid current execution whose baseline state is green or matching red; a matching-red string cannot erase a failed outcome. Ordinary matching red may retain its subordinate `findings` outcome while its gate state is satisfied, as before. Preserve exact identity, affected paths, provenance, diagnostics, and the mandatory separate canonical-gate record.

### Task 3 — Restore usable continuation argv and validate the integrated behavior

Within `evidence_gates.go`, construct missing-input continuations in their owning channel: qualification `--diff-range` and repeated `--gate-arg` flags belong before `--`; estimator/preflight/probe arguments belong after it. Reuse the known phase-specific shape and preserve already supplied valid inputs, exact REQ/path, gate argv, timeout, and paired baseline paths. Placeholder substitution must be sufficient; callers must not repair syntax or know the leaf command. A continued green-record check must not manufacture a new direct-run zero attestation.

For subordinate retry/verification commands that would re-enter an evidence helper, translate them back to the same request-bound advance phase with its accepted inputs. Preserve actual diagnostic commands such as `git diff`; these inspect a finding and do not bypass a gate. Preserve the deliberate direct canonical-gate `next_argv` plus its advance verification path. Do not blanket-replace every subordinate command or introduce a second handler registry.

Run the new focused public regressions and all four existing owner packages: `go test -count=1 ./internal/lifecycleadvance ./internal/corehelpers ./internal/nextselection ./internal/resultmodel`. Exercise the existing interruption/process cleanup coverage, and compile the changed platform signature with `GOOS=windows GOARCH=amd64 go test -c ./internal/nextselection -o <owned temporary path>`. Run gofmt and the CLI prime's applicable vet/race checks. Let the orchestrator perform the canonical post-merge gate and select renewed heavy evidence for the actual integrated range; original saved green results cannot prove changed source.

Hand back the exact eight-path manifest, RED/GREEN outputs, preserved ordinary matching-red controls, argv round trips, and unresolved report-only findings. Independent re-review should verify that real failed launch and real timeout remain non-satisfied with a valid canonical green record present, and that ordinary child exits 124/125 are neither mislabeled nor blanket-rejected. No second remediation loop is authorized by this plan.

## Remediation Decisions

- **D-08 — execution evidence:** Carry observed launch/timeout/error facts across the existing private runner boundary; ordinary exit values 124/125 remain eligible baseline red. Include F02 continuation repair in the same existing owner. F03/F04 remain noncritical report-only findings.
- **D-09 — explicit resumed run:** After the shared-state refusal was reported, the user explicitly repeated run-to-completion. Preserve identified active foreign bytes and claims; continue only through canonical selection/claim for this run. Recovery remains refused and is not represented as successful or sole-writer authority. Canonical default selection and claim succeeded for this exact REQ at 7ad53bff1d867f1453e1e7765e988dedb308e7e1.
- **D-10 — revision attribution:** Preserve original base24ed2fdda549a0759cdc571562c9b782bfeb6251 and saved implementation06367337dd82d97416e0d9d37872cc35b56ae7bc as historical evidence; record a fresh supplemental repair range and use its merged hash for current evidence. REQ-570's canonical claim removed prior commit/heavy frontmatter; do not restore old success fields.
- **D-11 — shared baseline:** Preserve the other live run's existing baseline bytes while retaining this REQ's canonically generated baseline in Git-private request storage; pass explicit paired paths at the later focused gate. Do not stage or adopt another run's baseline.

- **D-12 — builder isolation:** Use the Git-private checkout `.git/work-run-20260905/worktree-agent-REQ-506-focused-evidence` and branch `codex/worktree-agent-REQ-506-focused-evidence` to honor the user repository boundary and configured branch prefix. Source base is4adcff4e1470177959d0ef814cb5d66f43b7cbfb; all lifecycle writes remain on main.

## Remediation Execution State

- [x] **[PLAN]:** Three tasks and eight paths; independent public RED first, smallest existing-owner fix, GREEN and caller-boundary verification.
- [ ] **[APPLY]:** Await successful preflight and isolated builder.
- [ ] **[UNIFY]:** Await actual diff and verification.

## Pre-Build Gate Observation

Canonical focused preflight passed (exit0, launchedtrue) and its own baseline is preserved in `.git/work-run-20260905/REQ-506/baseline.json`; the foreign shared baseline bytes were restored. The directly executed full gate first failed only three <30-second file budgets (44.35/43.39/40.18s; all728 CLI tests passed). Its one retry, with GOMAXPROCS=2 to reduce contention and the same assertions/budget, instead failed ShellCheck SC2043 in the active foreign `_dev/tests/do-work-cli-launcher-behavior.sh:56` edit. Neither attempt is green evidence. This source is outside the repair's eight-path boundary; preserve it and await its owner's source change before re-verification. Builder has not been dispatched. Exact attempt evidence: `.git/work-run-20260905/REQ-506/prebuild-gate-attempts.json`.
