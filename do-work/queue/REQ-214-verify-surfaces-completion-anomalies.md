---
id: REQ-214
title: verify surfaces completion anomalies as findings
status: pending
created_at: 2026-08-17T08:25:24Z
user_request: UR-048
addendum_to: REQ-213
review_generated: true
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
depends_on: []
maintenance: false
effort_estimate: normal
write_set: [skills/do-work-board/tools/queue-kanban/verify.go, skills/do-work-board/tools/queue-kanban/verify_test.go]
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-08-17T08:25:24Z
  basis:
    - Route B
    - 2-file write set
---

# verify Surfaces Completion Anomalies as Findings

## What

`queue-kanban verify` reports `OK: no findings` on a tree carrying 10 flagged completion anomalies — its only board-warnings lift is the duplicate-request-id prefix filter (`verify.go:241-251`). Lift `CompletionAnomaly` tickets into verify findings so an anomalous archive fails the mechanical check instead of passing silently.

**Finding provenance (REQ-213 review, Important 1, gate: user-visible):** "verify is blind to completion anomalies — it reports OK: no findings on a tree carrying 10 flagged tickets, yet the REQ's What/Requirements text promises verify's never-silent line surfaces the new class." Verified hands-on by the reviewer. Resolution chosen: add the probe (the alternative — correcting REQ-213's contract text — would leave the user's "surfaced and fixed" intent unmet in the one mode built for mechanical checking).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- verify emits one finding per `CompletionAnomaly` ticket (id + reason), in verify's existing finding format and severity conventions; exit code reflects findings per verify's existing contract.
- Reporting only — verify stays read-only; repairs remain cleanup's job (prime: "verify reports and routes").
- Go test: a fixture tree with a reversed-span ticket makes verify report it and exit non-zero; a clean tree still reports OK.

## Red-Green Proof
**RED prompt/case:** `queue-kanban verify` from this repo root → `OK: no findings`, exit 0, while the summary lists 10 completion anomalies including REQ-091.
**Why RED now:** No verify probe reads `CompletionAnomaly`.
**GREEN when:** The same invocation lists the anomalous tickets as findings and exits non-zero; module tests green.
**Validation:** Review-generated (REQ-213 Important 1); user intent "surfaced and fixed" recorded at capture.

## Full Context
See `do-work/user-requests/UR-048/input.md`.
