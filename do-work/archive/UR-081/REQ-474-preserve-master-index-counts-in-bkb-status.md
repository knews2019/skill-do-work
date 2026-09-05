---
id: REQ-474
title: 'Review fix: Preserve master-index counts in BKB status'
status: cancelled
created_at: 2026-09-01T05:59:27Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-416]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-416]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-416
completed_at: 2026-09-01T12:02:37Z
---

# Preserve Master-Index Counts in BKB Status

## What

Restore the characterized `bkb-status` meaning: article and topic-cluster counts come from `wiki/_master_index.md`, not recursive disk file counts. Status must expose malformed or inconsistent index data actionably rather than silently substituting another source.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, whose root cause is the BKB status count-source substitution.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- Parse the canonical master-index article/topic counts exactly as the pre-migration BKB status action defined them.
- Keep disk inventories available only where separately named; never relabel disk counts as master-index counts.
- Return actionable typed findings for missing, malformed, duplicate, or inconsistent count declarations.
- Add clean and malformed fixtures, including a master index declaring 17 articles/3 topic clusters while disk inventory differs; prove text/JSON parity and read-only behavior.

## Red-Green Proof

**RED prompt/case:** Run `bkb-status` on a KB whose `_master_index.md` declares 17 articles and 3 topic clusters but has no corresponding disk pages.
**Why RED now:** The command reports recursive disk counts `0/0`, changing the characterized snapshot meaning.
**GREEN when:** Status reports `17/3` from the master index and malformed index counts produce typed evidence without mutation.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-416-rereview.md`.

---
*Source: REQ-416 fresh re-review residual finding 3.*

## Cancelled

- **When:** 2026-09-01T12:02:37Z
- **Why:** folded into REQ-420 as bkb-status master-index parity fixtures (maintainer decision, 2026-09-01 queue analysis)
- **Decided by:** user, via `do-work abandon`
