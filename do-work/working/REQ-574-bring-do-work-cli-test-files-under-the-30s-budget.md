---
id: REQ-574
title: 'Repository gate repair: bring do-work-cli test files under the 30s per-file budget'
status: claimed
route: C
created_at: 2026-09-04T23:50:46Z
user_request: UR-115
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: do-work-cli-test-file-budget
depends_on: []
related: [REQ-572]
claimed_at: 2026-09-04T23:59:43Z
---

# Repository gate repair: bring do-work-cli test files under the 30s per-file budget

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: go-test-file-budget:do-work-cli:publication-defer-gate-test affecting REQ-572 (found by REQ-572 / UR-115)

## Repository Gate Repair Intake

- **Parent:** REQ-572
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** go-test-file-budget:do-work-cli:publication-defer-gate-test
- **Repair dependency:** REQ-574
- **Diagnostic evidence:** "post-merge run at 4adcff4e (fbdcd35e merged): FAIL: internal/corehelpers/inventory_test.go accumulated 38.92s; internal/publication/defer_gate_test.go 37.01s; internal/finalization/finalization_recovery_test.go 35.65s; internal/finalization/finalization_req499_test.go 30.85s; each test file must finish under 30s; every test passed"
- **Diagnostic evidence:** "pre-build run at f6c43d22: the same four files over budget (38.61s, 37.93s, 35.73s, 33.72s)"
- **Diagnostic evidence:** "detached diagnostic worktree at base 7ad53bff (clean tree): FAIL: internal/publication/defer_gate_test.go accumulated 32.52s; each test file must finish under 30s; every test passed; queue-kanban 24s"
- **Implementation base:** 7ad53bff1d867f1453e1e7765e988dedb308e7e1
- **Implementation merge:** fbdcd35e0908aca6a01f554cc9b7fd7c85347a49
