---
id: REQ-548
title: 'Repair repository gate update-layout broken pipe'
status: claimed
route: C
created_at: 2026-09-03T18:35:22Z
user_request: UR-102
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: repository-gate-update-layout-broken-pipe
depends_on: []
related: [REQ-531]
claimed_at: 2026-09-03T18:37:57Z
---

# Repair repository gate update-layout broken pipe

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: sha256:efed40e27e755df5b2733ea1fdf3d3f228c42e8d859baefb1f89b35769ebda2e affecting REQ-531 (found by REQ-531 / UR-102)

## Repository Gate Repair Intake

- **Parent:** REQ-531
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** sha256:efed40e27e755df5b2733ea1fdf3d3f228c42e8d859baefb1f89b35769ebda2e
- **Repair dependency:** REQ-548
- **Diagnostic evidence:** "maintainer-verify: update-script behavior probes failed"
- **Diagnostic evidence:** "update-script-behavior: suite update identifies layout output did not match four-module suite"
- **Diagnostic evidence:** "update-script-behavior.sh: printf write error: Broken pipe"
