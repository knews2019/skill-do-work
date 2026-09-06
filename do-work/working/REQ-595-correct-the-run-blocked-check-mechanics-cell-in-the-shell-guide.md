---
id: REQ-595
status: claimed
domain: general
created_at: 2026-09-06T03:56:38Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
depends_on: [REQ-555]
related: [REQ-555]
write_set: [skills/do-work/docs/prescribed-shell-primitives.md]
title: 'Correct the run-blocked-check mechanics cell, which describes shell the Go command does not run'
claimed_at: 2026-09-06T04:38:49Z
---

# Correct the run-blocked-check Mechanics Cell in the Shell Guide

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Row 13 of the "Shipped executable homes" table in
`skills/do-work/docs/prescribed-shell-primitives.md` gives `run-blocked-check`'s owned mechanics as
"GNU timeout selection and isolated stock-Bash process-group timeout/cleanup". The mechanics moved into
Go, and the Go implementation selects no external `timeout` binary and builds no stock-Bash process
group. The cell describes shell that no longer runs.

## Why

Same finding class as REQ-555: the guide is the pointer target from sixteen shipped files and it
currently describes an implementation that does not exist. REQ-555 corrected the route column; this
cell is the same defect one column over, and was equally stale before that change, so it was recorded
rather than folded in.

## Context

Found during the independent three-lens review of REQ-555, which confirmed the cell is false against
`blocked_probe_unix.go` and scored it as an adjacent finding outside that request's declared class.

## Detailed Requirements

- Read the Go implementation of the `run-blocked-check` subcommand and write the mechanics it actually
  owns. Do not paraphrase the old cell.
- Check the other thirteen Mechanics cells against their implementations in the same pass, and correct
  any that are stale. The row that was checked is not evidence about the rows that were not.
- Scope is exactly this finding class: no behaviour change, no code change, no route-column change.

## Constraints

- Guide prose only. `_dev/tests/audit-lockins.sh` already pins the route column and the orchestration
  claim from REQ-555; do not weaken either.
- If a Mechanics cell cannot be verified against code, say so in the request rather than guessing.

## Dependencies

Depends on REQ-555, which rewrote the route column of the same table.

## Open Questions

None.
