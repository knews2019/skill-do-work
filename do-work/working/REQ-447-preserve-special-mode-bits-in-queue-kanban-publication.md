---
id: REQ-447
title: '[impact-rule-change] Preserve special mode bits in queue-kanban publication'
status: claimed
created_at: 2026-08-31T20:30:00Z
status_changed_at: 2026-08-31T20:30:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: [REQ-436]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
addendum_to: REQ-436
sweep: true
sweep_key: preserve-special-mode-bits-in-file-publication
claimed_at: 2026-09-01T04:10:28Z
---

# Preserve Special Mode Bits in Queue-Kanban Publication

## What

Extend REQ-436's complete-mode publication contract to the separate queue-kanban atomic writer at `skills/do-work-board/tools/queue-kanban/atomic_write.go`. It still narrows the original regular target with `Mode().Perm()` before replacement.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Preserve ordinary permissions plus setuid, setgid, and sticky through queue-kanban's existing-target atomic replacement.
- Add RED/GREEN tests for each special bit using the real writer and full-mode assertions.
- Retain existing symlink/special-file refusal, changed-target detection, atomic replacement, and sync/close behavior.

## Instances

- [ ] `skills/do-work-board/tools/queue-kanban/atomic_write.go`: existing-target publication narrows the original complete mode with `Mode().Perm()`. (found by REQ-436 / UR-081)
- [ ] `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`: existing-untracked rollback snapshots only `Mode().Perm()` and recreates with `os.WriteFile`, stripping setuid, setgid, and sticky bits. (found by REQ-412 / UR-081)

## Red-Green Proof

**RED prompt/case:** Replace queue-kanban fixtures carrying `04640`, `02640`, and `01640`, then assert exact contents and full low-twelve-bit modes.
**Why RED now:** The writer passes `originalInfo.Mode().Perm()` to chmod, the same narrowing shape REQ-426 and REQ-436 proved strips all three special bits.
**GREEN when:** All three modes and contents survive the real publication seam without weakening its existing safety fixtures.

---
*Source: discovered while implementing REQ-436 (UR-081).*
