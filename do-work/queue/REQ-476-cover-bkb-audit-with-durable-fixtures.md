---
id: REQ-476
title: '[impact-negligible] Review fix: Cover BKB audit with durable fixtures'
status: pending
created_at: 2026-09-01T08:32:57Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-417]
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-417]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-417
sweep: true
sweep_key: memory-bkb-audit-tdd-coverage-missing
---

# Cover BKB Audit with Durable Fixtures

## What

Add the plan-required committed BKB audit matrix so the implemented mechanical probes cannot regress without a failing test. Done means every supported BKB discovery shape and evidence class is exercised through the real command, not only a reviewer-only fixture.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this missing BKB audit coverage root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `auditBKBEngine` has no committed Go coverage.
- [ ] `countBKBInboundReferences` has no committed Go coverage.

## Requirements

- Add real-command fixtures for explicit, default, absolute, and parent BKB discovery.
- Cover repository shape, committed history/authors, inbound references, absent and malformed ledger evidence, pre-ledger fairness, and exact classification boundaries.
- Prove text/JSON parity and byte-for-byte read-only behavior.
- Make the named audit functions execute under the committed test suite and retain focused coverage evidence.

## Red-Green Proof

**RED prompt/case:** Run focused coverage for `internal/knowledgecommands` and require the BKB audit engine plus inbound-reference counter to execute through real-command fixtures.
**Why RED now:** Both functions report 0.0% coverage even though reviewer-only acceptance fixtures pass.
**GREEN when:** The committed BKB audit matrix passes and the named functions are exercised with all required discovery/evidence rows.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/runs/work-2026-08-31-165510/REQ-417-rereview.md`.

---
*Source: REQ-417 fresh re-review residual finding 2.*
