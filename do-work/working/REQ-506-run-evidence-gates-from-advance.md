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
planning_at: 2026-09-04T17:36:43Z
exploration_at: 2026-09-04T17:41:25Z
dispatch_at: 2026-09-04T17:44:35Z
builder_handback_at: 2026-09-04T18:15:00Z
integration_at: 2026-09-04T18:15:22Z
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
- `_dev/primes/lessons-action-files.md` — 4163 tokens, over the shared budget; `slugged: partial` so no complete targeted form. Matched on action pipeline and downstream-reader changes.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the shared budget; `slugged: partial` so no complete targeted form. Matched on prescribed argv and migration-parity behavior.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 6265 tokens, over the shared budget; `slugged: partial` so no complete targeted form. Matched on structured evidence, owned probes, command composition, and failure identity.

## Open Questions
None.

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

## Implementation Summary

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

## Orientation

[MAP CHANGED] The lifecycle now has one public request-bound evidence path: advance discovers the current phase, executes its existing mechanical handler, and returns ordered typed gate records; the work action supplies inputs and retains semantic judgment.
