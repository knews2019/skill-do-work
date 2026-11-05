---
id: REQ-711
title: 'Repair runtime fingerprint baseline failures'
status: claimed
route: C
created_at: 2026-11-05T16:34:59Z
user_request: UR-163
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: runtime-fingerprint-baseline
depends_on: []
related: [REQ-710]
claimed_at: 2026-11-05T16:35:25Z
---

# Repair runtime fingerprint baseline failures

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: 14816037a6fd3188551b9a540c8b7fbb01470e6bdfd681f43fa1931d3c4c6e97 affecting REQ-710 (found by REQ-710 / UR-163)

## Repository Gate Repair Intake

- **Parent:** REQ-710
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** 14816037a6fd3188551b9a540c8b7fbb01470e6bdfd681f43fa1931d3c4c6e97
- **Repair dependency:** REQ-711
- **Diagnostic evidence:** "baseline revision 40e9e4f6ee0f9ed19404d1515e7750f3a453a398"
- **Diagnostic evidence:** "canonical gate exited 1 twice before implementation"
- **Diagnostic evidence:** "TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes: default runtime must have a determinable fingerprint (GOOS, GOARCH, CGO_ENABLED, DO_WORK_TEST_DO_WORK_CLI_BINARY)"
- **Diagnostic evidence:** "TestShippedGitIsolationPreservesGenericLaneInheritance: shipped runtime probe did not isolate host Git configuration"
