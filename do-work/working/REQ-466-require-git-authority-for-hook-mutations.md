---
id: REQ-466
title: 'Review fix: Require Git authority for hook mutations'
status: claimed
created_at: 2026-09-01T04:01:15Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-415]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 50
  confidence: medium
  calculated_at: 2026-09-03T14:42:54Z
  basis:
    - Route C
    - 5-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - dependency depth 1
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
related: [REQ-415, REQ-463]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-415
claimed_at: 2026-09-03T14:42:07Z
route: C
planning_at: 2026-09-03T14:43:58Z
---

# Require Git Authority for Hook Mutations

## What

Make all public hook commands honor UR-081's platform rule that mutations require an actual Git worktree. SessionStart timestamp/reservation work and memory log/ledger appends must be read-only or refuse actionably outside Git, with one explicitly documented command contract and no partial mutation.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, whose root cause is the missing cross-hook Git mutation prerequisite. REQ-438 concerns mismatched roots inside `gittransaction`, while REQ-463 governs committed evidence and final eligibility for reservation deletion.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- Define one shared hook precondition that distinguishes an actual Git worktree from a plain directory and unavailable/ambiguous Git authority.
- Before any timestamp, reservation, memory-log, or usage-ledger mutation, require that precondition; preserve every target byte on refusal.
- Preserve exact hook protocol and nonblocking Stop behavior while returning actionable typed evidence for skipped/refused mutation.
- Add real-command fixtures for all three hook commands outside Git, with Git unavailable, and inside a valid worktree; assert status, exact protocol bytes, typed JSON, and filesystem non-effects/effects.

## Red-Green Proof

**RED prompt/case:** Run `session-start`, `memory-session-start`, and `memory-stop-capture` against writable non-Git roots containing eligible timestamp or memory targets.
**Why RED now:** The commands currently apply those writes successfully even though UR-081 requires Git for mutation and read-only behavior elsewhere.
**GREEN when:** Every non-Git or ambiguous-authority case remains byte-identical with actionable typed evidence, while valid Git worktrees retain the characterized hook protocols and effects.

## Folded From REQ-529

Hand triage 2026-09-03, maintainer approved: REQ-529 (Give the cancellation reason the same containment condition) is cancelled and its requirement lands here as an acceptance criterion, so UR-081's remaining review leftovers are one REQ.

- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:895-903` `cancellationReasonBlock` inlines a reason on the sole test that it contains no newline, so a one-line reason of `***` or `## Notes` is written straight into the document. Apply the same condition-complete Markdown-delimiter classifier REQ-460 introduced for answer summaries, as one shared predicate. `skills/do-work/actions/abandon.md:66` already asks for that judgment.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-415-rereview.md`.

---
*Source: REQ-415 fresh re-review residual finding 1.*

## Triage

**Route: C** - Complex

**Reasoning:** The change spans three public hook commands, a shared Git-authority precondition, Markdown delimiter safety, and real-command regression fixtures across success and refusal paths.

**Planning:** Required

## Plan

1. Characterize all three hook commands at the real command boundary, then add failing fixtures for writable non-Git roots and unavailable Git. Assert exact hook status/stdout/stderr, typed `findings`/`changes`/`skipped_work`/`protocol_output`, and byte-identical mutation targets; keep valid-worktree fixtures as the compatibility control.
2. Add one shared hook Git-authority precondition and call it before timestamp, reservation, memory-log, or usage-ledger writes. Preserve SessionStart failure propagation and Stop's nonblocking protocol while returning actionable typed evidence. Reuse the condition-complete outside-text delimiter predicate for cancellation reasons instead of adding a second example list.
3. Run focused hook/request-state tests, retained hook behavior suites, Go vet and module tests, then the canonical maintainer gate. Sweep sibling hook mutation paths for the same missing precondition and record any out-of-scope instance rather than expanding silently.

The command result remains the mutation consumer contract: every refusal must preserve request identity, Git-authority outcome, finding code/severity/evidence, exact target paths, skipped-work reason, and protocol bytes so the launcher and JSON reader observe one event without partial effects.

*Generated inline after the Plan agent service failed twice*

**Plan validation:** All four captured requirements and the folded cancellation-delimiter acceptance criterion map to the three tasks; no task is orphaned, the plan stays below the five-task warning threshold, and command-driven mutations name the identity, authority, outcome, and protocol fields their consumers require.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 3215 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes hook mutation authority, characterization coverage, and a condition-based Markdown classifier.
