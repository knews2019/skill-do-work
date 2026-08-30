---
id: REQ-432
title: 'Review fix: Enforce the commit guard for consumed scratch cleanup'
status: pending
domain: general
created_at: 2026-08-30T20:35:44Z
user_request: UR-081
addendum_to: REQ-409
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Enforce the Commit Guard for Consumed Scratch Cleanup

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Keep the narrow non-rollback consumed-scratch exception from bypassing cleanup's `--commit` empty-index precondition.

## Context
Found during re-review of REQ-409. With an unrelated staged file, `cleanup --commit` still deleted consumed untracked scratch and returned success because the scratch exception exempted the group's preflight failure. Fold-first scan found no pending REQ or sweep in any UR that shares this commit-guard root cause.

## Requirements
- Enforce the empty-index precondition before any `--commit` cleanup mutation, including non-rollback scratch deletion.
- Preserve the exact-inventory, rooted-containment, and consumed-manifest checks already applied to scratch.
- Return truthful refusal evidence without deleting scratch when the index is nonempty.

## Red-Green Proof
**RED prompt/case:** Stage an unrelated file, create an untracked consumed run, and invoke `cleanup --commit`; assert cleanup refuses and the run remains byte-for-byte present.
**Why RED now:** The spent-scratch exception currently bypasses any group preflight failure, including the commit guard.
**GREEN when:** The named fixture refuses before mutation while the same scratch remains eligible in a non-commit cleanup run.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
