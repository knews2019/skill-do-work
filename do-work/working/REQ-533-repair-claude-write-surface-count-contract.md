---
id: REQ-533
title: 'Repair CLAUDE.md write-surface count contract'
status: claimed
route: C
created_at: 2026-09-03T12:06:41Z
user_request: UR-102
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: claude-write-surface-count-stale
depends_on: []
related: [REQ-532]
claimed_at: 2026-09-03T12:07:48Z
---

# Repair CLAUDE.md write-surface count contract

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: contract-regressions:claude-write-surface-count-stale affecting REQ-532 (found by REQ-532 / UR-102)

## Repository Gate Repair Intake

- **Parent:** REQ-532
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** contract-regressions:claude-write-surface-count-stale
- **Repair dependency:** REQ-533
- **Diagnostic evidence:** "CLAUDE.md must state the tool has exactly three write surfaces once next-req reserves ids; testing fields, next-version, and reservation markers are the complete set."
