---
id: REQ-411
title: 'Implement dependency-aware queue selection and actionable summaries'
status: completed
claimed_at: 2026-08-31T17:49:48Z
route: C
created_at: 2026-08-29T20:28:26Z
status_changed_at: 2026-08-30T22:14:34Z
estimate:
  p50_active_minutes: 100
  confidence: low
  calculated_at: 2026-08-30T22:10:19Z
  basis:
    - Route C
    - 18-file write set
    - 8 new files
    - 7 subsystems involved
    - 6 acceptance criteria
    - dependency depth 5
    - persistence changes
    - cross-route regression gates
    - full-suite verification
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-410, REQ-428, REQ-429, REQ-435]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
write_set:
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_types_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/select-simple-reqs.sh
  - _dev/tests/select-simple-reqs-behavior.sh
  - skills/do-work/actions/run-simple-reqs.md
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - _dev/tests/contract-regressions.sh
batch: go-no-llm-command-platform
completed_at: 2026-08-31T19:38:00Z
commit: 6209227b
---

# Implement Dependency-Aware Queue Selection and Actionable Summaries

## What
Move deterministic queue selection and readiness decisions into canonical Go commands.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Freeze a read-only typed selector around the existing repository snapshot, dependency graph, result envelope, and process-group-safe blocked-probe boundary.
- [x] **[APPLY]:** Added the mixed real-binary RED first, then implemented typed selection, command registration, compatibility delegation, and contract coverage inside the frozen 18-file scope.
- [x] **[UNIFY]:** Reviewed all 18 files and passed focused/full CLI tests, vet, selector behavior, blocked-process-tree contracts, exact Go 1.25, shell syntax, and diff/scope audit; no debug artifacts remain.

## Detailed Requirements
- Implement `do-work-next` and shared selection for explicit REQ/UR targeting, dependency readiness, cycle detection, waves, assignments, negligible-impact filtering, blocked probes, and estimates.
- Compose stable queue summaries from the same typed repository model and return actionable reasons for every skipped or refused item.
- Preserve explicit targeting semantics and successful-dependency gates.
- Provide text/JSON output with exact next argv/Just recipes and verification commands.

## Constraints
- Keep `queue-kanban` as the separate board/UI binary; share behavior through contracts or compatible parsing, not binary consolidation.

## Dependencies
Depends on REQ-410 (doctor and normalized repository findings).

## Builder Guidance
Certainty level: Firm. Lock current selection semantics with fixtures before replacing action-side scans.

## Red-Green Proof
**RED prompt/case:** Select from a fixture containing targeted IDs, a dependency chain, a cycle, an assignment, negligible impact, blocked probes, and wave depth.
**Why RED now:** Selection logic is distributed between action prose and small shell helpers rather than one runnable command.
**GREEN when:** The CLI returns the same eligible set and stable reason for every excluded item in text and JSON, without LLM judgment for deterministic gates.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Triage

**Route: C** - Complex

**Reasoning:** This request introduces shared dependency-aware selection across explicit targets, UR expansion, waves, assignments, blocked probes, estimates, text/JSON output, and action/Just integration. Its existing estimate identifies an 18-file, seven-subsystem implementation with persistence and cross-route regression requirements.

**Planning:** Required

## Plan

1. Add a read-only `internal/nextselection` command family that consumes one typed repository snapshot, reuses the canonical dependency graph, and proves default/explicit-REQ/UR targeting, readiness, cycles, assignments, negligible filtering, blocked probes, estimates, waves, and fan-out with a mixed RED fixture.
2. Register `next` and return stable typed text/JSON selection records with an actionable reason for every exclusion plus exact next argv, Just recipe, and verification argv. Reuse the shipped process-group-safe blocked-check runner; do not reproduce its timeout mechanics.
3. Delegate deterministic work/simple-run selection to the command while leaving claims and all state/release transactions in later REQs, then run focused, full CLI/vet, selector, blocked-probe, contract, exact-Go, and canonical repository gates.

**Architecture decisions:** Selection is strictly read-only; `dependencygraph` remains the sole readiness/cycle authority; typed selection data may extend the shared result envelope but must not be hidden in prose evidence; queue-kanban remains separate.

**Plan validation:** All six acceptance themes map to the three tasks; no task is orphaned. The three-task shape stays within the Route C quality ceiling, though the explorer must shrink or justify every shared contract-file edit before scope freezes.

*Generated by Plan agent and validated by work action*

## Exploration

The existing typed repository snapshot, request projection, and dependency graph already own discovery, normalized selection fields, readiness, cycles, missing/ambiguous dependencies, and depth. The implementation needs no production change in those packages. The shared result envelope needs one additive typed selection/exclusion payload; blocked probes must reuse `scripts/run-blocked-check.sh`; and the simple selector retains only its extra effort/maintenance/security/critical policy above canonical readiness.

The exact 18-file boundary is eight new `internal/nextselection` files plus command registration, result rendering, work/simple-action delegation, selector parity, prime mapping, and contract tests. The exploration report records the interfaces, RED matrix, commands, and process/provenance risks.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify)
- `skills/do-work/tools/select-simple-reqs.sh` (modify)
- `_dev/tests/select-simple-reqs-behavior.sh` (modify)
- `skills/do-work/actions/run-simple-reqs.md` (modify)
- `skills/do-work/actions/work.md` (modify)
- `skills/do-work/actions/work-reference.md` (modify)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify)
- `_dev/tests/contract-regressions.sh` (modify)

**Files I will NOT touch:** repository/request/dependencygraph production schemas, cleanup/doctor implementations, queue-kanban, claim/archive/release state-transition code, queue state from the builder, version files, or changelogs.

**Acceptance criteria (restated from REQ):**
- [x] One typed read-only command resolves default, explicit REQ, and UR targets with exact provenance-sensitive override semantics.
- [x] Canonical dependency readiness, cycles, wave depth, fan-out bounds, assignment/claim filters, negligible filtering, blocked probes, and estimates produce stable selected or excluded records.
- [x] Every skipped/refused item has an actionable reason and exact next argv/Just/verification command in text and JSON.
- [x] Work delegates deterministic Step 1 selection without moving claim/state/release ownership into this REQ.
- [x] The cheap-work selector consumes canonical readiness while preserving all specialized vetoes and diagnostics.
- [x] Mixed RED/GREEN fixtures and contract regressions prove text/JSON parity, deterministic ordering, timeout reuse, and no repository mutation.

## Implementation Summary

**What was done:** Added a read-only typed `next` command that selects default, targeted, UR-expanded, wave, fan-out, and cheap-work candidates from one repository snapshot and the canonical dependency graph. It reports stable selected/excluded records, estimates, summaries, exact commands, exact request paths, original states, and per-record probe outcomes in text/JSON; delegates blocked probes to the shipped bounded runner; and gives the action sufficient exact-target evidence for its existing unblock transaction without another queue scan.

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go` (new) — selection options, provenance, candidate, probe, and identifier helpers.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_types_test.go` (new) — identifier and canonicalization coverage.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go` (new) — default, REQ, and UR target resolution with stable provenance/order.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go` (new) — explicit-over-UR provenance and expansion coverage.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (new) — canonical eligibility, dependency, probe, estimate, wave/fan-out, summary, and simple-policy projection.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (new) — mutation-sensitive mixed selection fixtures.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go` (new) — command parsing, discovery, shipped-runner resolution, and private probe execution.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go` (new) — real-binary RED/GREEN and command behavior coverage.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified) — registers the next-command handlers.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified) — adds normalized typed selected/excluded/summary fields, per-record probe states, exact request targets, and matching text rendering.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified) — pins non-null JSON, probe outcome rendering, and text/JSON selection parity.
- `skills/do-work/tools/select-simple-reqs.sh` (modified) — replaces shell-side queue parsing with a compatibility invocation of the simple selector mode.
- `_dev/tests/select-simple-reqs-behavior.sh` (modified) — keeps legacy run-set and diagnostic behavior pinned through delegation.
- `skills/do-work/actions/run-simple-reqs.md` (modified) — delegates directly to the canonical typed selector command.
- `skills/do-work/actions/work.md` (modified) — makes typed selector output Step 1's sole deterministic queue read while retaining action-owned mutations.
- `skills/do-work/actions/work-reference.md` (modified) — ties auto-wave computation and fan-out exclusions to typed selector records.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — maps the new package and selection ownership.
- `_dev/tests/contract-regressions.sh` (modified) — pins sole-authority delegation and retained blocked-runner safety.

**Integration range:** `dd221ea8..6209227b`

*Generated by work action from the builder hand-back*

## Decisions

### D-01: Keep the shared result schema version additive

**Decision:** DECIDE & STATE — retain schema version 1 and add normalized selected, excluded, and selection-summary fields.

**Reasoning:** Existing consumers remain compatible while new consumers receive typed queue evidence. Value: one renderer contract without a breaking migration. Risk: future removals or semantic changes still require an explicit schema-version decision.

### D-02: Preserve the cheap selector as a compatibility renderer

**Decision:** DECIDE & STATE — keep the shell path as a thin caller of `next --simple`, retaining its run-set and stderr warning surfaces.

**Reasoning:** Queue parsing and readiness move to Go without breaking existing callers. Value: immediate authority consolidation. Risk: the compatibility output must stay pinned until all callers consume typed JSON directly.

### D-03: Reproduce selection inputs in every verification command

**Decision:** DECIDE & STATE — attach one command-level verification argv carrying target provenance and all selection flags to every selected and excluded record.

**Reasoning:** Verification must replay the same mode rather than accidentally turn an exclusion into an explicit-target override. Value: pasteable, faithful diagnostics. Risk: every future selection flag must join this reconstruction in the same change.

## Qualification

Passed — all 18 declared files are substantive and present in `dd221ea8..6209227b`. The new package is registered through the CLI entry point, the result envelope carries typed selection and transition data, and both work callers delegate to the same command. The eight static-reference warnings are expected Go package files compiled by package/import membership; P-A-U, requirement tracing, exact scope, and debug-artifact checks passed.

## Testing

**Red-green validation:** The mixed real-binary fixture failed before production changes with exit 2 and `UNKNOWN-COMMAND` for `next`. The same fixture and mutation-sensitive unit cases pass after implementation, covering target provenance, dependencies/cycles, wave/fan-out separation, assignment/negligible filters, blocked probes, estimates, deterministic ordering, typed text/JSON output, and read-only bytes.

**Merged-state checks:**
- `go test -count=1 ./internal/nextselection ./internal/resultmodel` — PASS.
- `go test -count=1 ./...` and `go vet ./...` — PASS (full CLI module).
- `bash _dev/tests/select-simple-reqs-behavior.sh` — PASS.
- `bash _dev/tests/contract-regressions.sh` — PASS, including blocked-probe process-tree behavior.
- `bash _dev/tests/do-work-cli-go125-compatibility.sh` — PASS (exact Go 1.25.0).
- Qualification, scope-drift, shell syntax, and diff hygiene over `dd221ea8..6209227b` — PASS.

## Initial Review

**Overall: 50%** | 2026-08-31T19:04:56Z

| Dimension | Score |
|-----------|-------|
| Requirements | 75% |
| Code Quality | 80% |
| Test Adequacy | 75% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Successful blocked-probe identity is collapsed into an aggregate count; selected records expose neither a per-record probe outcome/original state nor an exact queue mutation target, so the action cannot perform its required unblock transaction from the sole typed queue read — impact-user-visible → returned for the one allowed remediation attempt.

**Minor findings:** None.
**Acceptance:** Fail — deterministic selection is broad and correct, but its typed action handoff loses required per-record transition evidence.
**Suggested testing:** Mixed ordinary and successfully probed records; exact mutation targets; distinct successful/failed/timed-out/missing probe evidence; action contract forbidding a second queue scan.
**Follow-ups created:** None pending remediation; **sweeps appended to:** None.

*Reviewed by review-work action*

## Remediation

**Initial review:** Acceptance failed at 50% because successful blocked-probe identity was collapsed into an aggregate count, leaving the action without a deterministic exact target for its owned unblock/history mutation.

**One allowed attempt:** Added exact repository-relative request paths, original statuses, typed per-record probe states, attempted/exit evidence, and an explicit unblock requirement to selected and excluded records. A mixed fixture with two ordinary and two successfully probed blocked REQs failed before the production change, then passed with exact targets preserved even across fan-out exclusion. Missing, failed, timed-out, launch-failed, and successful probes remain distinct. Work now validates and exact-reads only returned targets, never rescans or re-probes the queue.

**Remediation merge:** `6209227b` (branch commit `75aa3a24`), changing eight files inside the frozen 18-file scope.

## Review

**Overall: 98%** | 2026-08-31T19:34:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None. The initial blocked-probe handoff finding is closed.
**Minor findings:** None.
**Acceptance:** Pass — exact request paths, original states, and per-record probe/unblock evidence now survive every selection outcome, including fan-out exclusion, so the action can perform its owned transition without rescanning or re-probing.
**Suggested testing:** An optional direct `launch_failed` fixture could supplement the inspected production branches; it is not a release blocker.
**Follow-ups created:** None; **sweeps appended to:** `_dev/primes/lessons-action-files.md`.

*Re-reviewed by review-work action after the one allowed remediation attempt*

## Lessons Learned

A read-only decision command must preserve every identity-to-outcome association needed by the action that owns the later mutation. Aggregate counts are useful presentation evidence, but they cannot authorize an exact state transition; the typed result must retain the contained request path, original state, and outcome through all record-shape conversions.

## Orientation

Queue selection now has one typed, read-only authority in `do-work-cli next`. Future selection policy belongs in `internal/nextselection` with matching text/JSON and action-contract coverage. State mutations remain in the work action, which must consume exact returned targets and reject stale evidence rather than rescan the queue.
