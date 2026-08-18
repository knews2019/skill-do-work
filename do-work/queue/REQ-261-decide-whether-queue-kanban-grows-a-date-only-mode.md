---
id: REQ-261
title: Decide whether queue-kanban grows a date-only output mode
status: pending-answers
created_at: 2026-08-18T19:30:47Z
user_request: UR-055
addendum_to: REQ-253
domain: general
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work-reference.md
---

# Decide Whether queue-kanban Grows a Date-Only Output Mode

## What

The Timestamp rule's date-only paragraph carries its own tripwire: "adding a date-only mode for a single consumer would spend the skill's narrow compiled-tooling exception … **revisit if a second consumer appears**." REQ-253 made ui-review's report header the second named date-only consumer, so the tripwire has tripped. The question is the maintainer's: add a `queue-kanban` date-only subcommand now that two consumers exist, or amend the sentence to state a higher threshold (the POSIX/PowerShell one-liners already cover the floor).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Discovered by REQ-253's builder ([low]; the tripwire sentence was left verbatim as a deliberate maintainer call). Note the board's pinned write-surface count is unaffected either way — `now`-style output is read-only.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-253: the date-only paragraph's "revisit if a second consumer appears" condition is now true. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending') — and the builder decides between the subcommand and a re-stated threshold.
  Also: No, discard it — two consumers on the shell one-liner is still fine.
