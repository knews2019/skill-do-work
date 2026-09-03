---
id: REQ-511
title: '[impact-rule-change] Review fix: Disable optional locks for read-only Git probes'
status: cancelled
domain: backend
created_at: 2026-09-02T17:48:52Z
user_request: UR-097
addendum_to: REQ-500
review_generated: true
impact: impact-rule-change
effort_estimate: effort-substantive
tdd: true
suggested_spec: bug-fix
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
sweep: true
sweep_key: read-only-git-probes-allow-index-writes
completed_at: 2026-09-03T11:47:24Z
---

# Review Fix: Disable Optional Locks for Read-Only Git Probes

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Make every Git subprocess that is contractually read-only disable optional locks, and prove the class cannot recur with exact index-byte tests rather than porcelain-only checks.

## Context

Found during the post-remediation re-review of REQ-500. Its consolidated archive inventory fixed visibility and performance, but `git status` refreshed `.git/index` after a clean archive mtime change even though the command reported no working-tree mutation. The fold-first scan found no pending request or sweep sharing this read-only Git-probe root cause.

## Instances

- [ ] `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go`: the archive inventory probe used by doctor and SessionStart runs `git status` without `GIT_OPTIONAL_LOCKS=0`.
- [ ] Sweep other CLI Git subprocesses whose declared authority is read-only and apply the same no-optional-lock contract wherever index refresh is possible.

## Requirements

- Set `GIT_OPTIONAL_LOCKS=0` for the doctor archive inventory and every other Git subprocess governed by the same byte-for-byte read-only contract.
- Preserve the caller's environment while overriding only this Git control; do not weaken error reporting or typed tail evidence.
- Add a regression that touches a clean tracked archive, runs doctor and SessionStart, and compares `.git/index` bytes before and after.
- Add a mechanical guard or centralized helper so newly added read-only Git probes cannot silently omit the environment control.
- Keep porcelain output, tail findings, and the exact SessionStart pointer unchanged.

## Red-Green Proof

**RED prompt/case:** Touch a clean committed blank-provenance archive without changing its bytes, hash `.git/index`, run SessionStart, and hash the index again.
**Why RED now:** The current read-only archive inventory lets Git refresh index metadata, so the hashes differ while `git status --short` remains empty.
**GREEN when:** The exact index bytes remain unchanged for doctor and SessionStart, all read-only Git probes satisfy the centralized guard, and existing diagnostic behavior tests pass.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

None.

---
*Source: REQ-500 post-remediation review finding.*

## Cancelled

- **When:** 2026-09-03T11:47:24Z
- **Why:** no observed failure from optional Git locks on read-only probes; recapture on evidence (maintainer decision, 2026-09-03 roadmap triage)
- **Decided by:** user, via `do-work abandon`
