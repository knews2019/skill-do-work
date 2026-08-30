---
id: REQ-433
title: 'Review fix: Split misplaced UR partial-merge conflicts by item'
status: pending
domain: general
created_at: 2026-08-30T20:35:44Z
user_request: UR-081
addendum_to: REQ-409
review_generated: true
depends_on: [REQ-432]
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Split Misplaced UR Partial-Merge Conflicts by Item

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Plan a misplaced archived UR's files in per-item conflict domains so occupied destinations do not block unrelated siblings from reaching the canonical archive tree.

## Context
Found during re-review of REQ-409. Pass 3b currently bundles an entire misplaced UR folder into one group, so a conflict on `input.md` suppresses a safe sibling REQ move. Fold-first scan found no pending REQ or sweep in any UR that shares this partial-merge root cause.

## Requirements
- Plan each misplaced UR file against its own destination and conflict evidence.
- Move every nonconflicting sibling while leaving only conflicting sources in place.
- Preserve deterministic ordering and never overwrite a destination.

## Red-Green Proof
**RED prompt/case:** Create a misplaced UR folder with conflicting `input.md` and a nonconflicting sibling REQ, then run cleanup; assert the REQ moves and only `input.md` is refused.
**Why RED now:** Pass 3b assigns all files in the folder to one conflict domain.
**GREEN when:** The named fixture performs the safe partial merge and reports exact evidence for the remaining conflict.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
