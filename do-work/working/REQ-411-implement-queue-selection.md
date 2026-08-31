---
id: REQ-411
title: 'Implement dependency-aware queue selection and actionable summaries'
status: claimed
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
---

# Implement Dependency-Aware Queue Selection and Actionable Summaries

## What
Move deterministic queue selection and readiness decisions into canonical Go commands.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
- [ ] One typed read-only command resolves default, explicit REQ, and UR targets with exact provenance-sensitive override semantics.
- [ ] Canonical dependency readiness, cycles, wave depth, fan-out bounds, assignment/claim filters, negligible filtering, blocked probes, and estimates produce stable selected or excluded records.
- [ ] Every skipped/refused item has an actionable reason and exact next argv/Just/verification command in text and JSON.
- [ ] Work delegates deterministic Step 1 selection without moving claim/state/release ownership into this REQ.
- [ ] The cheap-work selector consumes canonical readiness while preserving all specialized vetoes and diagnostics.
- [ ] Mixed RED/GREEN fixtures and contract regressions prove text/JSON parity, deterministic ordering, timeout reuse, and no repository mutation.
