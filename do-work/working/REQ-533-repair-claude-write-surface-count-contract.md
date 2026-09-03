---
id: REQ-533
title: 'Repair CLAUDE.md write-surface count contract'
status: claimed
created_at: 2026-09-03T12:06:41Z
user_request: UR-102
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-03T12:08:08Z
  basis:
    - Route C
    - 1-file write set
    - 1 acceptance criteria
    - full-suite verification
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: claude-write-surface-count-stale
depends_on: []
related: [REQ-532]
status_changed_at: 2026-09-03T12:38:36Z
claimed_at: 2026-09-03T12:38:59Z
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

---

## Repository Gate Repair No-Op

- **Expected diagnostic fingerprint:** contract-regressions:claude-write-surface-count-stale
- **Gate command:** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 0
- **Observed result:** green before implementation; repair already satisfied
- **Verified at:** 2026-09-03T12:29:35Z

