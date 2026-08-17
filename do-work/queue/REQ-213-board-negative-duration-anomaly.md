---
id: REQ-213
title: Board surfaces negative claimed→completed duration as a completion anomaly
status: pending
created_at: 2026-08-17T08:05:43Z
user_request: UR-048
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-211, REQ-212]
batch: estimator-calibration
write_set: [skills/do-work-board/tools/queue-kanban/model.go, skills/do-work-board/tools/queue-kanban/model_test.go]
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-17T08:05:43Z
  basis:
    - Route B
    - 3-file write set
    - (priced with the pre-calibration table)
---

# Board Surfaces Negative Claimed→Completed Duration as a Completion Anomaly

## What

`detectCompletionAnomaly` currently short-circuits whenever `completed_at` parses, so a REQ whose `completed_at` is earlier than its `claimed_at` (a real case: archived REQ-091) renders as normal. Flag that negative span as a completion anomaly so the always-visible strip and `verify` surface it for repair.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Fires only when **both** `claimed_at` and `completed_at` parse under the board's timestamp rules **and** completed is strictly earlier than claimed — a reversed span cannot be real for stamps written in order.
- Joins the existing `CompletionAnomaly`/`CompletionAnomalyReason` plumbing (anomaly strip + `verify`'s never-silent line) — no parallel channel, no new JSON fields.
- The reason names both raw values and the likely cause/fix in the established style (one stamp is usually local wall-clock written with a `Z` suffix; fix by rewriting the wrong stamp with the true UTC instant).
- Unparseable or absent stamps remain the existing checks' territory — this check must not double-report them.
- Go table-driven test: reversed span → flagged with both values in the reason; ordered span and absent claimed_at → not flagged by this check; existing anomaly tests unchanged.
- Read-only reporting: no write-surface change, no parser field additions (both stamps are already parsed fields).

## Constraints

- Batch constraint: the three-write-surface sentence must not need amending.
- Board versioning is folded into the skill — normal CHANGELOG entry + suite version bump, per the prime.

## Builder Guidance

Firm on strict-both-parse gating and reuse of the existing anomaly channel; latitude on exact reason wording.

## Red-Green Proof
**RED prompt/case:** A fixture ticket with `claimed_at: 2026-01-02T10:00:00Z`, `completed_at: 2026-01-01T10:00:00Z`, `status: completed` — `detectCompletionAnomaly` returns false today (completed_at parsed ⇒ short-circuit).
**Why RED now:** The frontmatter-parsed path returns before any span comparison; archive REQ-091 demonstrates the case reaching real data.
**GREEN when:** The new Go test shows that fixture flagged with a reason naming both stamps, the ordered-span fixture stays unflagged, and `go test -count=1 ./...` passes in the module.
**Validation:** User confirmed — "queue-kanban should report the negative duration anomaly so it can be surfaced and fixed."

## Full Context
See `do-work/user-requests/UR-048/input.md` for complete verbatim input.

---
*Source: UR-048 — negative-duration anomaly*
