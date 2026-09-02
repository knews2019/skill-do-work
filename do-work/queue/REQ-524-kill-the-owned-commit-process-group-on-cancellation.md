---
id: REQ-524
title: 'Kill the owned commit process group on cancellation'
status: pending
created_at: 2026-09-02T23:58:00Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-457]
---

# Kill the Owned Commit Process Group on Cancellation

## What

Make cancellation of a media Git transaction terminate the whole owned process group, not just the direct `git` child. A `pre-commit` hook that ignores `SIGTERM` currently outlives the cancelled transaction and keeps running after the command has returned.

## Instances

- `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback` fails with `media commit hook survived cancellation` (`report_image_process_test.go:85`). The transaction rolls back correctly and returns within the deadline; the orphaned hook process is the only unmet assertion.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Discovered by:** REQ-457's Step 5.75 pre-flight, as one of two pre-existing failures in the `bash _dev/tests/maintainer-verify.sh` baseline at `008f3d3`. Unrelated to REQ-457's created-path ownership invariant, so it is captured here rather than folded into that REQ.
- **Evidence:** the test writes a `pre-commit` hook that traps and ignores `TERM` then sleeps 30s, cancels the context once the hook is running, and asserts `syscall.Kill(pid, 0)` fails afterwards.

## Detailed Requirements

- Cancelling a Git-committing transaction must terminate every process the transaction launched, including grandchildren a hook spawns, not only the direct child.
- Escalate past a process that ignores the graceful signal; a hook trapping `SIGTERM` must not survive.
- Preserve the existing return-within-deadline behavior and the existing rollback-to-preimage behavior on cancellation.
- Do not leave a killed process group's exit status misreported as success.

## Constraints

- Keep the owned-process-group runner the single seam for launched subprocesses; do not add a second ad-hoc kill path.
- Preserve exact typed result and rollback contracts in `prime-do-work-cli.md`.

## Dependencies

No request prerequisite.

## Red-Green Proof

**RED prompt/case:** `go test ./internal/toolboxcommands/ -run TestRemediationCancellationReachesMediaGitCommitAndRollback` — currently fails with `media commit hook survived cancellation`.
**Why RED now:** cancellation signals only the direct `git` child, so a hook that ignores `SIGTERM` keeps running as an orphan.
**GREEN when:** that test passes without modification, the transaction still returns inside its deadline, and the target still rolls back to its preimage bytes.

---
*Source: REQ-457 pre-flight baseline, captured during the work run.*
