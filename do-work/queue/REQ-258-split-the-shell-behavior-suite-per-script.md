---
id: REQ-258
title: Split the prescribed shell behavior suite per script
status: pending-answers
created_at: 2026-08-18T17:49:24Z
user_request: UR-056
addendum_to: REQ-246
domain: general
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Split the Prescribed Shell Behavior Suite Per Script

## What

`_dev/tests/prescribed-shell-scripts-behavior.sh` now carries 47 named cases, and the ten reservation-cleanup + timestamp-repair cases dominate its tail. If it keeps growing, per-script files may read and fail more legibly. Organizational only — no case changes.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Builder-discovered on REQ-246 (Discovered Tasks, third item), classified [low]. Touches only a test file but is a reorganization, not mechanical hygiene, so it takes the consent flow rather than the test-hygiene carve-out.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-246: the shell behavior suite is one growing file and could split per-script. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — one file is fine until it actually hurts.
