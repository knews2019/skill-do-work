---
id: REQ-485
title: 'Canonicalize REQ reservation marker filenames across allocation flows'
status: claimed
priority: now
created_at: 2026-09-01T12:11:03Z
user_request: UR-092
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
batch: go-no-llm-command-platform
claimed_at: 2026-09-03T21:43:41Z
---

# Canonicalize REQ Reservation Marker Filenames Across Allocation Flows

## What

Make every REQ-number reservation flow use one canonical marker filename so
exclusive-create collision detection actually collides. Today queue-kanban
`next-req` writes zero-padded `REQ-000482` while capture-files manifests carry
unpadded `REQ-482`; the two flows compute the same candidate from the same max-scan
and both succeed, defeating the guard entirely.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in
any UR sharing this root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Pick one canonical marker filename and enforce it in every writer: queue-kanban
  `allocate.go` (`next-req`), the capture-files marker-path guidance in
  `skills/do-work/actions/capture.md`, and any other flow that creates markers.
- Read-side acceptance of both legacy spellings everywhere markers are consumed:
  the allocation max-scan, `skills/do-work/scripts/cleanup-req-reservations.sh`'s
  committed-REQ and timeout reaping, and any capture-side duplicate scan — an
  existing marker in either spelling must block its number and must reap.
- Writers refuse or normalize a non-canonical marker path in a capture manifest
  rather than silently creating a second spelling.
- Lock-in test reproducing the 2026-09-01 collision shape: two flows reserving the
  same number through different spellings must collide, not both succeed.
- The board parser lock-step rule (`_dev/primes/prime-kanban-board.md`) governs if
  the board reads marker names anywhere; verify and follow it.

## Red-Green Proof

**RED prompt/case:** Reserve a number via `queue-kanban next-req`, then submit a
capture-files manifest whose reservation path is the unpadded spelling of the same
number.
**Why RED now:** Both succeed — reproduced 2026-09-01: `next-req` created
`REQ-000482` at 11:50Z and a concurrent capture created `REQ-482` and committed a
second REQ-482 file (78847fe4) at 11:56Z; only manual inspection prevented a
duplicate REQ id from being committed twice.
**GREEN when:** The second reservation is refused with an actionable finding, the
lock-in test pins the cross-spelling collision, and legacy-spelling markers still
count in the max-scan and still reap.
**Validation:** Observed collision, this session; evidence in the UR input.

## Full Context

See `do-work/user-requests/UR-092/input.md` for complete verbatim input.

---
*Source: UR-092 (Canonicalize REQ reservation marker filenames across allocation flows)*
