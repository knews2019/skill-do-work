---
id: REQ-467
title: '[impact-rule-change] Review fix: Reconcile SessionStart authority restatements'
status: pending
created_at: 2026-09-01T04:01:15Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-415]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-415, REQ-420, REQ-463]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-415
---

# Reconcile SessionStart Authority Restatements

## What

Update every live shipped instruction, board check, and remediation that still assigns SessionStart reservation cleanup or timestamp repair to the retired shell implementations. Current readers must identify the Go hook/core authority, while retained scripts are described only according to their actual compatibility role and unsafe legacy paths are not recommended.

The fold-first scan found REQ-420 adjacent but not dependency-ready, so it is not a legal fold destination. This is not prose-only: contradictory shipped guidance changes which executable an agent invokes and can direct readers to the pre-REQ-463 fail-open cleanup path.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `README.md`: SessionStart is described as running the old reservation-cleanup and timestamp-repair scripts. (found by REQ-415 / UR-081)
- [ ] `skills/do-work/actions/capture.md`: lifecycle guidance still assigns cleanup/repair authority to retained scripts. (found by REQ-415 / UR-081)
- [ ] `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`: board authority documentation points at the retired SessionStart implementations. (found by REQ-415 / UR-081)
- [ ] `skills/do-work-board/tools/queue-kanban/verify.go` and its remedies/tests: live findings direct readers to inspect or invoke the old scripts. (found by REQ-415 / UR-081)

## Requirements

- Sweep all live consumers and restatements of SessionStart reservation/timestamp ownership, including wording variants rather than only filenames.
- Make user/agent guidance identify the registered Go hook/core owner and the thin retained SessionStart launcher accurately.
- Replace board findings and recovery commands that recommend retired or fail-open script paths with runnable canonical CLI evidence and verification.
- Add contract tests that fail when live shipped guidance reintroduces the old ownership claim.

## Red-Green Proof

**RED prompt/case:** Run the shipped-reference/board contract sweep and inspect current guidance for `cleanup-req-reservations.sh` and `repair-req-timestamps.sh` as SessionStart authorities.
**Why RED now:** README, capture action, board prime, and board remedies contradict the migrated Go authority and can lead readers to unsafe legacy behavior.
**GREEN when:** Every live restatement and board action points to the canonical Go owner/launcher contract, with a regression preventing stale script authority from returning.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-415-rereview.md`.

---
*Source: REQ-415 fresh re-review residual finding 2.*
