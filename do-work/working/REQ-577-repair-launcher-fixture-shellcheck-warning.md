---
id: REQ-577
title: 'Repository gate repair: remove the launcher fixture single-iteration loop'
status: claimed
route: C
created_at: 2026-09-04T23:55:41Z
user_request: UR-098
domain: backend
tdd: 'true'
maintenance: 'false'
impact: impact-critical
effort_estimate: effort-substantive
repository_gate_repair: 'true'
sweep: 'true'
sweep_key: do-work-cli-launcher-single-iteration-loop
depends_on: []
related: [REQ-506]
claimed_at: 2026-09-04T23:58:14Z
---

# Repository gate repair: remove the launcher fixture single-iteration loop

## What

Repair the repository-gate failure recorded below so dependency-gated requests can resume.

## Instances

- [ ] repository gate: shellcheck:sc2043:do-work-cli-launcher-behavior:single-iteration-loop affecting REQ-506 (found by REQ-506 / UR-098)

## Repository Gate Repair Intake

- **Parent:** REQ-506
- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** shellcheck:sc2043:do-work-cli-launcher-behavior:single-iteration-loop
- **Repair dependency:** REQ-577
- **Diagnostic evidence:** "ShellCheck warning SC2043 in _dev/tests/do-work-cli-launcher-behavior.sh: for command_name in bash; do; This loop will only ever run once. Bad quoting or missing glob/expansion?"
- **Diagnostic evidence:** "The required second direct canonical gate run exited 1 on this warning. The launcher test bytes now belong to committed release7cceea12 and remain unchanged since that failure. First attempt had only Go per-file timing overruns; those are not the deciding fingerprint."
