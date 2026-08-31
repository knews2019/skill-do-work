---
id: REQ-433
title: 'Review fix: Split misplaced UR partial-merge conflicts by item'
status: claimed
claimed_at: 2026-08-31T17:49:48Z
route: A
domain: general
created_at: 2026-08-30T20:35:44Z
user_request: UR-081
addendum_to: REQ-409
review_generated: true
depends_on: [REQ-432]
impact: impact-user-visible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-31T17:49:48Z
  basis:
    - trivial short-circuit
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

---

## Triage

**Route: A** - Simple

**Reasoning:** The failure is isolated to one cleanup grouping boundary, with a concrete two-item regression and explicit non-overwrite behavior.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` (modify) — split a misplaced archived UR directory into independently preflighted per-file operation groups.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go` (modify) — add the exact conflicting-`input.md`/safe-sibling RED-GREEN regression and pin deterministic grouping.

**Files I will NOT touch:** cleanup application/transaction mechanics, repository-model schemas, action prose, queue state from the builder, release metadata, changelogs, or version files.

**Acceptance criteria (restated from REQ):**
- [ ] Each file under `do-work/archive/user-requests/UR-NNN/` is planned in its own conflict domain against its canonical destination.
- [ ] A conflicting `input.md` remains at the source while a nonconflicting sibling REQ moves to `do-work/archive/UR-NNN/`.
- [ ] Destination conflicts never overwrite existing content and report exact refusal evidence.
- [ ] Per-item groups and operations retain deterministic ordering.
