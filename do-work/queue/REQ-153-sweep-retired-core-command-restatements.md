---
id: REQ-153
title: "Review fix: Sweep retired core command restatements"
status: pending
domain: general
created_at: 2026-08-08T18:32:01Z
user_request: UR-031
addendum_to: REQ-146
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: retired-core-command-restatements
---

# Review Fix: Sweep Retired Core Command Restatements

## What
Remove or update every live restatement of the retired core moved-command and transition-era updater contracts, and add a shipped-surface guard so this class cannot recur.

## Context
Found during review of REQ-146. The runtime shim and routes were deleted, but several live templates, UI hints, hook messages, and updater-prime lessons still tell users or agents to use the retired contract.

## Requirements
- Replace each retired core invocation with the owning modular skill command or current guidance.
- Rewrite transition-era updater-prime lessons to describe only the permanent modular contract.
- Add a shipped-runtime/restatement regression that rejects retired core invocations and stale transition rules while excluding historical records and explicit negative-test fixtures.
- Preserve legitimate references to generic work, CI, or data pipelines and historical evidence.
- Verify every replacement points to an existing action in exactly one owning skill.

## Instances
- [ ] `skills/do-work-board/justfile.template`: replace `do-work install just-kanban` with the current board-owned invocation.
- [ ] `web/template.html`: replace the live `do-work board` hint with the current board-owned invocation.
- [ ] `web/board.js`: replace the live `do-work board` hint with the current board-owned invocation.
- [ ] `skills/do-work-knowledge/hooks/memory-session-start.sh`: replace `do-work memory recall` with the current knowledge-owned invocation.
- [ ] `skills/do-work/tools/prime-do-work-update.md`: remove transition-era lessons that staged skills remain export-ignored or marker-free legacy spans must fail.
