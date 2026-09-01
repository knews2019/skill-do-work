---
id: REQ-466
title: 'Review fix: Require Git authority for hook mutations'
status: pending
created_at: 2026-09-01T04:01:15Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-415]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-415, REQ-463]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-415
---

# Require Git Authority for Hook Mutations

## What

Make all public hook commands honor UR-081's platform rule that mutations require an actual Git worktree. SessionStart timestamp/reservation work and memory log/ledger appends must be read-only or refuse actionably outside Git, with one explicitly documented command contract and no partial mutation.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, whose root cause is the missing cross-hook Git mutation prerequisite. REQ-438 concerns mismatched roots inside `gittransaction`, while REQ-463 governs committed evidence and final eligibility for reservation deletion.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- Define one shared hook precondition that distinguishes an actual Git worktree from a plain directory and unavailable/ambiguous Git authority.
- Before any timestamp, reservation, memory-log, or usage-ledger mutation, require that precondition; preserve every target byte on refusal.
- Preserve exact hook protocol and nonblocking Stop behavior while returning actionable typed evidence for skipped/refused mutation.
- Add real-command fixtures for all three hook commands outside Git, with Git unavailable, and inside a valid worktree; assert status, exact protocol bytes, typed JSON, and filesystem non-effects/effects.

## Red-Green Proof

**RED prompt/case:** Run `session-start`, `memory-session-start`, and `memory-stop-capture` against writable non-Git roots containing eligible timestamp or memory targets.
**Why RED now:** The commands currently apply those writes successfully even though UR-081 requires Git for mutation and read-only behavior elsewhere.
**GREEN when:** Every non-Git or ambiguous-authority case remains byte-identical with actionable typed evidence, while valid Git worktrees retain the characterized hook protocols and effects.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-415-rereview.md`.

---
*Source: REQ-415 fresh re-review residual finding 1.*
