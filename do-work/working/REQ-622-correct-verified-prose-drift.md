---
id: REQ-622
title: 'Correct verified prose drift'
status: claimed
created_at: 2026-09-08T21:49:05Z
user_request: UR-130
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-kanban-board.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: false
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: []
related: ["REQ-623", "REQ-624"]
batch: instruction-history-cleanup
write_set: ["do-work/prose-backlog.md", "skills/do-work/scripts/repair-req-timestamps.sh", "skills/do-work/actions/work-reference.md", "skills/do-work-board/tools/queue-kanban/open_work.go", "skills/do-work-board/tools/queue-kanban/testing.go", "skills/do-work-board/actions/board.md", "_dev/primes/prime-action-files.md", "skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go", "skills/do-work/actions/forensics.md", "skills/do-work-board/tools/queue-kanban/web/board-core.js", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md", "skills/do-work-board/docs/board-guide.md", "skills/do-work-board/tools/queue-kanban/citations.go", "skills/do-work/actions/abandon.md", "skills/do-work-board/tools/queue-kanban/generate.go", "README.md", "skills/do-work/actions/clarify.md", "_dev/tests/contract-regressions.sh", "skills/do-work/tools/do-work-cli/internal/publication/answer.go", "skills/do-work/actions/verify-requests.md", "_dev/lessons/validated-runtime-boundaries.md"]
claimed_at: 2026-09-08T21:53:21Z
---

# Correct Verified Prose Drift

## What
Reconcile the 19 open entries in do-work/prose-backlog.md against current code, canonical contracts, and recorded intent. Correct surviving prose discrepancies and record which entries are resolved, obsolete, or still need an intent decision. Prefer removing duplicate statements; preserve runtime behavior.

## Detailed Requirements
- Revalidate the 19 open entries preserved in `do-work/user-requests/UR-130/assets/prose-backlog-at-capture.md`, including the stale isolation lesson. The snapshot identifies this drain's items; it does not assert they remain valid today.
- Give every entry an evidence-backed disposition: corrected, already obsolete, or unresolved intent. Tick resolved/obsolete lines in the live backlog with evidence and `(drained by REQ-622)`; retain unresolved lines and explain the conflict.

## Constraints
Preserve runtime behavior. Current code, canonical prose, and recorded intent must be reconciled; a disagreement is not permission to silently change a rule. Remove redundant wording where that resolves the discrepancy.

## Dependencies
First in the approved batch. REQ-623 (Reduce orchestration instruction duplication) follows this request; REQ-624 (Addendum: Trim the shipped changelog to 50 entries) follows that request.

## Builder Guidance
The backlog is the bounded scope, not a mandate to rewrite surrounding mechanisms. No pending or pending-answers candidate in any UR shares this request's root cause; the seven existing queued requests address different release, commit-integrity, and diagnostics defects.

## Completion Proof
Account for all 19 source items and verify changed references and affected contracts against their canonical owners. Record the evidence behind each disposition; code-comment edits preserve executable behavior. An unresolved intent decision remains explicit rather than being reported as a fixed discrepancy.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 5879 tokens; action/prime prose and condition-preserving extraction. Partial slug coverage prevents targeted selection.
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens; timestamp script and prescribed containment wording. Partial slug coverage prevents targeted selection.
- `_dev/primes/lessons-kanban-board.md` — 5912 tokens; the board prime governs named board paths. Partial slug coverage prevents targeted selection.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 6888 tokens; the board's own prime governs the named comments, docs, and testing files. Partial slug coverage prevents targeted selection.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 15905 tokens; alternate answer-writer contract drift and publication comments. Partial slug coverage prevents targeted selection.

## Full Context
See `do-work/user-requests/UR-130/input.md` for the user instruction, adopted report, and batch decisions. The approved prompt excerpts are in `do-work/user-requests/UR-130/assets/rev-v2-approved-scope.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read the listed primes and agent rules; record a brief approach.
- [ ] **[APPLY]:** Make the scoped changes.
- [ ] **[UNIFY]:** Review every changed file and record the focused checks and outcomes.
