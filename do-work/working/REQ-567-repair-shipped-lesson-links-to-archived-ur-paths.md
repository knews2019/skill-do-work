---
id: REQ-567
title: 'Repair shipped lesson links to archived UR paths'
status: claimed
route: C
created_at: 2026-09-04T15:07:20Z
user_request: UR-098
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: shipped-lesson-links-obsolete-archive-paths
depends_on: []
related: [REQ-503]
claimed_at: 2026-09-04T15:08:27Z
---

# Repair shipped lesson links to archived UR paths

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: sha256:3af85b84722557f94ddfd466fc32136086fb5fed306e478bd344f689902472ff affecting REQ-503 (found by REQ-503 / UR-098)

## Repository Gate Repair Intake

- **Parent:** REQ-503
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** sha256:3af85b84722557f94ddfd466fc32136086fb5fed306e478bd344f689902472ff
- **Repair dependency:** REQ-567
- **Diagnostic evidence:** "shipped-package-reference-contract: obsolete archive link for REQ-491"
- **Diagnostic evidence:** "shipped-package-reference-contract: obsolete archive link for REQ-492"
- **Diagnostic evidence:** "shipped-package-reference-contract: obsolete archive link for REQ-493"
