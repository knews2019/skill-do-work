# REQ-761 builder brief

Implement the request below in your T3-assigned checkout and branch. Record their absolute path and name in the handback. Read CLAUDE.md, the applicable primes and the implementation crew rules in skills/do-work/actions/work-implementation-reference.md, including required lesson entries. Consult do-work/lessons-index.md and the linked user request. Preserve sibling changes. Only this request is yours. Use the REQ's own validation, capture the baseline before edits, follow red/green where specified, run focused checks, review the complete diff, and commit the implementation with [REQ-761] at the beginning of the subject. Do not publish, merge or modify queue records, shared version files or CHANGELOG.

Write the completed handback to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-26-011150/REQ-761-handback.md in the main checkout. Provide operative branch, SHA, base, file manifest, changed behavior, focused commands and results, baseline/red evidence, decisions and discovered tasks (explicitly say when none), integration seams and any limitations. Append UTC start/end lines to /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-26-011150/REQ-761-progress.log for long operations. These two main-tree paths are the only shared writes. No other agent owns your implementation. Full gate and release are serial integration work. No project run-policy file.

---
id: REQ-761
title: 'Match board executable names exactly before restart'
status: claimed
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-11-25T23:27:30Z
  basis:
  - trivial short-circuit
route: A
created_at: 2026-11-25T23:09:19Z
user_request: UR-191
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
impact: impact-user-visible
effort_estimate: effort-mechanical
maintenance: false
claimed_at: 2026-11-25T23:27:19Z
---
# Match board executable names exactly before restart
## What
Use an exact executable basename match in `skills/do-work-board/tools/run-kanban.sh`. Only `queue-kanban` and the legacy `go-bin-queue-kanban` qualify for an automatic stop. Commit `0a3d12e7` documents the legacy cross-repository restart.

## Red-Green Proof
**RED prompt/case:** In `_dev/tests/contracts/queue-kanban.sh`, run the foreign-listener fixture with `queue-kanban-helper`. The current wildcard at `skills/do-work-board/tools/run-kanban.sh:21` accepts the lookalike. The validation replay returned exit 1 but recorded a signal to stub PID 3131.
**GREEN when:** Run `bash _dev/tests/contracts/queue-kanban.sh`. For each documented board name, assert that the fixture signals its listener once and starts the replacement successfully. For `queue-kanban-helper`, assert refusal (exit 1) and an empty signal log. Keep the existing foreign-listener case.
**Validation:** User confirmed.

## Context
Accepted finding F1 from [the proposal report](../../ai-reports/2026-11-26_0041_proposal-repository-fix-priorities/index.html#f1), rated high impact. The fix narrows the existing signal decision.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` (9,341 tokens) covers shipped shell and exported-function fixtures.
- `_dev/primes/lessons-kanban-board.md` (7,493 tokens) covers exact basename matching.
Both index entries have partial family coverage and exceed the 2,000-token budget as whole files.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: [Original request](../user-requests/UR-191/input.md)*

## Triage

Route A. The executable check, both supported basenames and the fixture assertions are specified.

## Plan

**Planning not required** — Route A, direct implementation.

Run the shell-contract fixture with stubs (never a real process). Make the two known executable names restart successfully with exactly one recorded signal, keep the existing foreign case, and prove the lookalike receives none.
