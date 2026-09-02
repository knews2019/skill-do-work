---
id: REQ-502
title: 'Review fix: Remove enriched checkpoint entries in cleanup mover'
status: pending
domain: backend
created_at: 2026-09-02T14:26:49Z
user_request: UR-083
addendum_to: REQ-489
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
suggested_spec: bug-fix
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
sweep: true
sweep_key: checkpoint-section-blind-line-editing
---

# Review Fix: Remove Enriched Checkpoint Entries in Cleanup Mover

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Make every cleanup-package departure from `do-work/working/` remove the complete matching checkpoint entry, not only its header line. Done means the cleanup mover and canonical request-state mover share the same exact-section, whole-entry semantics, so this root cause cannot recur through an alternate departure path.

## Context

Independent review of REQ-489 found that `internal/cleanup.ownedCheckpointRemoval` still globally filters only the matching `- REQ-NNN:` header and leaves Step 10's indented detail lines orphaned. The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this root cause; REQ-489 itself is already claimed and therefore is not a fold candidate.

## Requirements

- Remove the matching own-label cleanup claim header and all immediately following nonblank indented continuation lines.
- Preserve foreign-label entries and their continuation bytes exactly.
- Locate the real `## In Progress (interrupted)` section by a whole heading line; inline or backticked mentions elsewhere are not entries.
- Reuse or align with the request-state semantics without creating a second drifting definition.

## Instances

- [ ] `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` `ownedCheckpointRemoval`: cleanup archives a terminal working REQ but leaves its enriched checkpoint continuation block unattributed (found by REQ-489 / UR-083)

## Red-Green Proof

**RED prompt/case:** Extend `TestWorkingArchiveRemovesOnlyThisCheckoutCheckpointEntry` with an enriched own-label entry and an enriched foreign-label entry, then build the cleanup plan.
**Why RED now:** `ownedCheckpointRemoval` skips only the own header line and appends all following continuation lines.
**GREEN when:** Cleanup removes the complete own entry, preserves the foreign header and continuation bytes exactly, and ignores inline/backticked heading mentions.
**Validation:** Independent review finding from REQ-489; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---
*Source: Important review finding from REQ-489; folded as the next same-root sweep after the claimed source REQ.*
