---
id: REQ-423
title: 'Terminate archive fetching on interruption signals'
status: pending
created_at: 2026-08-29T20:26:10Z
user_request: UR-082
domain: backend
prime_files: ['_dev/primes/prime-shell-commands.md']
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-421, REQ-422, REQ-424]
batch: accepted-review-fixes
write_set: ['tools/fetch-upstream-archive.sh', 'skills/do-work/tools/fetch-upstream-archive.sh', '_dev/tests/update-script-behavior.sh', '_dev/primes/prime-shell-commands.md', '_dev/primes/lessons-shell-commands.md']
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
---

# Terminate Archive Fetching on Interruption Signals

## What
Make the archive fetcher stop after HUP, INT, or TERM instead of cleaning up and continuing into fallback success.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Reserve `EXIT` for cleanup in both byte-identical fetcher copies.
- Make HUP, INT, and TERM terminate with statuses 129, 130, and 143 respectively.
- Test every signal while a valid Git fallback exists.
- Assert no target archive is published and no success report is emitted after interruption.

## Constraints
- Keep root and shipped fetcher scripts byte-identical.
- Preserve existing HTTP and Git fallback behavior when no signal occurs.
- Add concise signal-preservation guidance to the shell prime and a detailed linked lesson entry.

## Dependencies
None. It shares fetcher and shell-test files with REQ-424 and is implemented as one shell slice in this batch.

## Builder Guidance
Certainty: Firm. Conventional `128 + signal` statuses are explicitly required.

## Context
No pending or unassigned queue candidate shares this root cause. Provenance: accepted review finding `[P2] Exit after handling fetch interruption signals` against `skills/do-work/tools/fetch-upstream-archive.sh:34`. The review states that a cleanup-only trap returns and permits fallback or false success after cancellation.

## Red-Green Proof
**RED prompt/case:** Inject HUP, INT, and TERM during HTTP fetching while a working Git fallback is configured.
**Why RED now:** Each cleanup trap returns, allowing execution to continue and possibly publish/report success.
**GREEN when:** Each case exits 129/130/143, publishes no target, and reports no success.
**Validation:** User accepted the review finding and supplied the implementation plan.

## Full Context
See `do-work/user-requests/UR-082/input.md` for the approved plan and batch constraints.

---
*Source: accepted review finding [P2] on fetch interruption, followed by the user-approved plan.*
