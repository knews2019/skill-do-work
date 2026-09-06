---
id: REQ-598
status: pending
domain: backend
created_at: 2026-09-06T06:25:19Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
maintenance: false
depends_on: [REQ-558]
related: [REQ-558]
write_set: [skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go]
title: 'Close the live nil-handle panic in transaction rollback, and decide the handle once instead of eleven times'
---

# Close the Live Nil-Handle Panic in Transaction Rollback

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

`rollbackFailure` opens a rooted filesystem handle with `os.OpenRoot`. On failure it records the
error and **deliberately keeps going**, so a failed transaction still unstages its paths and still
returns a typed incomplete rollback. That nil handle then fans out to eleven consumers across four
loops. Eight of them test it. **One does not**: `rollbackFailure` calls
`quarantineAndRollbackPrivate(root, state, published)` with no check, and that function dereferences
the handle with `root.Mkdir`. Any transaction with an identity-recorded private untracked target
panics mid-rollback when that open fails.

## Why

A panic during rollback is the worst available failure: the transaction has already mutated the tree and
the process dies before restoring it. The tool's own report of an incomplete rollback — the thing the
keep-going decision at the open exists to produce — is lost.

Beyond the one missing guard, the eight that exist are a per-consumer answer to a question that should
be decided once. `privateStateStillOriginal` and `rootedCreateRegular` take the same handle and never
test it; they are safe only because guards at three call sites stand in front of them, so a future edit
to any of those sites breaks them silently.

## Context

Found by the three-agent reachability trace run for REQ-558, which asked for eight of the nine guards to
be deleted. The trace established the opposite — the guards are load-bearing — and two of the three
agents independently reproduced this panic. REQ-558 deleted the one genuinely unreachable guard and
pinned the rest at eight. The full trace is in
`do-work/runs/work-2026-09-05-231943/REQ-558-handback.md` and the workflow's per-guard reports.

## Detailed Requirements

- Close the panic. The minimum is a guard at the `quarantineAndRollbackPrivate` call site matching its
  eight siblings; the better answer is the one below.
- **Decide the handle once.** At the open site, either abandon the rooted half of rollback with a typed
  finding, or run it under a proven handle. Either makes all eight surviving guards genuinely dead, and
  subsumes the missing one. The keep-going behaviour that lets a failed transaction still unstage its
  paths must be preserved — that is why the open does not return today.
- **Write the package's first no-handle rollback test.** No test in the package reaches any no-handle
  branch today, which is why nine guards could sit there with one of them missing and nothing noticed.
  The test is the point of this request as much as the fix is.
- If the guards become dead, delete them and say so; if they stay, say why. Either way update the count
  the REQ-558 lock-in pins, in the same change.

## Constraints

- One file plus the lock-in count in `_dev/tests/audit-lockins.sh`, which REQ-558 pinned at 8 as a
  floor and a ceiling. Changing the guard count without changing the pin fails the gate, which is the
  pin working.
- This is a transaction boundary. Every change is behaviour-preserving on the paths that work today, and
  the package's existing transaction and rollback tests must stay green unchanged.
- Do not remove the keep-going behaviour at the open site without replacing what it delivers.

## Red-Green Proof

**RED case:** a transaction with an identity-recorded private untracked target, rolled back with the
rooted open forced to fail. Today: panic at `root.Mkdir`.
**GREEN when:** the same case returns a typed incomplete rollback naming the unusable handle, the
package's transaction and rollback tests are green unchanged, and a new test drives the no-handle branch.

## Open Questions

None.
