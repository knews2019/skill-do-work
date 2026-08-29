---
id: REQ-418
title: 'Migrate toolbox commands and absorb audit-metrics into do-work-cli'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md]
tdd: true
suggested_spec:
depends_on: [REQ-417]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Migrate Toolbox Commands and Absorb Audit-Metrics into Do-Work-CLI

## What
Move deterministic toolbox utility behavior into the shared CLI and consolidate the separate audit-metrics module into it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement `do-work-note`, architecture-report preflight/publication, report-image lifecycle, portfolio publication, and last30days install/check.
- Preserve report/media process handling, publication atomicity, failure cleanup, and existing user-visible behavior.
- Absorb audit-metrics inventory, folders, churn, and hotspots behavior into `do-work-cli` with compatible flags/output plus JSON.
- Keep last30days’ Python prerequisite as a target-tool requirement, not a do-work implementation dependency.

## Constraints
- `queue-kanban` remains separate; only audit-metrics is consolidated.
- Characterization parity is required before retiring the separate audit-metrics source tree.

## Dependencies
Depends on REQ-417 (knowledge/store command migration precedes the toolbox family).

## Builder Guidance
Certainty level: Firm. Port existing audit-metrics tests rather than re-deriving its mature Git behavior.

## Red-Green Proof
**RED prompt/case:** Run current toolbox and audit-metrics fixture suites against absent shared CLI equivalents, including media cancellation and target-Python absence.
**Why RED now:** Toolbox behavior lives in shell and audit metrics lives in a separate Go module without the common result contract.
**GREEN when:** Shared commands match current status/output/effects, add actionable JSON, and preserve target-specific dependency reporting without Python/jq implementation branches.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
