---
id: REQ-584
title: 'Repository gate repair: add the missing shebang to the REQ-572 probe script'
status: pending
route: C
created_at: 2026-09-05T09:50:18Z
user_request: UR-098
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: do-work-runs-probe-script-missing-shebang
depends_on: []
related: [REQ-507]
---

# Repository gate repair: add the missing shebang to the REQ-572 probe script

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: shellcheck:sc2148:do-work-runs-req-572-probe:missing-shebang affecting REQ-507 (found by REQ-507 / UR-098)

## Repository Gate Repair Intake

- **Parent:** REQ-507
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** shellcheck:sc2148:do-work-runs-req-572-probe:missing-shebang
- **Repair dependency:** REQ-584
- **Diagnostic evidence:** "ShellCheck error SC2148 in do-work/runs/work-2026-09-04-232225/REQ-572-probe.sh line 1: Tips depend on target shell and yours is unknown. Add a shebang or a 'shell' directive."
- **Diagnostic evidence:** "Both direct canonical gate runs at 12d264c2 (detached worktree, clean tree) exited 1 on this one lint finding before any Go test ran. The probe file was committed by another session in 7ba3148a as REQ-572 run evidence and is outside REQ-507's implementation range 8e3dbf01..ad8bceb7, which touches no do-work/runs path."
- **Implementation base:** 8e3dbf01e0660424965d79acb2e386b6604e4780
- **Implementation merge:** ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9
