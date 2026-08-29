---
id: REQ-420
title: 'Replace shell implementations with shims and prove whole-suite parity'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-419]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419]
batch: go-no-llm-command-platform
---

# Replace Shell Implementations with Shims and Prove Whole-Suite Parity

## What
Complete the migration by making every retained shell path a thin launcher and enforcing full suite parity mechanically.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Replace domain logic in all 41 shipped shell utilities and hooks with thin build-and-exec compatibility shims while preserving every path.
- Add a mechanical contract that retained `.sh` files are launchers only and contain no embedded Python or jq implementation branches.
- Require old shim and new subcommand parity for exit status, output, and filesystem effects on characterized fixtures.
- Remove retired audit-metrics sources after its tests and behavior move into `do-work-cli`.
- Extend `_dev/tests/maintainer-verify.sh` with `go vet` and uncached `do-work-cli` tests, retain queue-kanban verification, and replace the separate audit-metrics lane.
- Prove final installation/update without Python/jq, every Just command without an LLM, unchanged skill aliases, and actionable findings that avoid repository rescans.

## Constraints
- Keep target-specific Python checks only for Python targets such as Python project preflight and last30days.
- Run the complete canonical maintainer gate and whole-suite parity verification before acceptance.

## Dependencies
Depends on REQ-419 (complete command and recipe surface).

## Builder Guidance
Certainty level: Firm. Retire implementations only after their parity fixtures pass against the Go engine.

## Red-Green Proof
**RED prompt/case:** Run the shell-thinness contract and full parity suite while any shipped shell still contains domain logic or Python/jq branches, and run install with those tools absent.
**Why RED now:** All 41 scripts/hooks still contain or route to shell domain implementations, and audit-metrics has a separate verification lane.
**GREEN when:** The mechanical contract proves launcher-only shell, parity fixtures pass, audit-metrics is consolidated, Go/board maintainer lanes pass uncached, and final no-Python/no-jq acceptance succeeds.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
