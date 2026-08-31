---
id: REQ-445
title: 'Review fix: Keep pathless refusal commands runnable'
status: claimed
domain: general
created_at: 2026-08-31T15:34:58Z
status_changed_at: 2026-08-31T19:24:17Z
user_request: UR-081
addendum_to: REQ-430
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
sweep: true
sweep_key: pathless-refusal-recovery-argv
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
claimed_at: 2026-08-31T22:07:41Z
---

# Review Fix: Keep Pathless Refusal Commands Runnable

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Make every structural cleanup refusal that has no concrete target path emit a runnable recovery command. Done means the class cannot recur: result-level coverage must reject empty or otherwise invalid path arguments for every applicable pathless refusal.

Fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this pathless-refusal recovery-command root cause.

## Context
Found during review of REQ-430 (Couple UR closure to terminal member archival). Duplicate operation-group codes are correctly refused before mutation, but the shared path-bearing refusal helper emits an empty Git pathspec and its suggested command exits 128.

## Requirements
- Add a valid no-path recovery-command form for structural cleanup refusals, or otherwise ensure the duplicate-group refusal emits runnable next and verification commands.
- Preserve exact path-targeted Git diagnostics for refusals that do have a concrete target path.
- Add a result-level ratchet that covers every applicable structural dependency refusal and fails on empty invalid path arguments.

## Instances
- [ ] `internal/cleanup/cleanup_apply.go`: duplicate group-code refusal calls the path-bearing finding helper with an empty path, producing `git status --short -- ''`. (found by REQ-430 / UR-081)

## Red-Green Proof
**RED prompt/case:** Exercise a cleanup plan with duplicate operation-group codes and execute or validate every emitted next command; require each command to be non-empty and runnable without an empty pathspec.
**Why RED now:** The duplicate-group finding appends an empty path argument to `git status --short --`, which exits 128.
**GREEN when:** The same result-level test passes and every applicable pathless structural refusal carries valid actionable command evidence.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions
- [x] Should I process this as a new task? The cleanup safety behavior already works, but users who hit a duplicate internal group identity receive a recovery command that immediately fails. → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to `pending` and repair the command contract with a regression test).
  Also: No, discard it (the safe refusal remains, but its recovery command stays unusable for this edge case).

  **Answered 2026-08-31** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the recommendation via `do-work clarify`: add the focused fix to the queue
  so every applicable pathless structural cleanup refusal emits runnable recovery commands
  with result-level regression coverage. Nothing from the captured scope was put out of scope.
