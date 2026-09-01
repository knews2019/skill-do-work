---
title: "Lessons from REQ-087: The board and verify hand the user the POSIX-only timestamp command the rule just stopped prescribing"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-087-the-board-and-verify-hand-the-user-the-p.md]
related:
  - page: REQ-078-the-windows-timestamp-fallback-cannot-ru
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-087: The board and verify hand the user the POSIX-only timestamp command the rule just stopped prescribing

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

REQ-078 made `actions/work-reference.md`'s Timestamp rule the only place in `actions/` that spells a
command for obtaining a stamp, and gave that rule a Windows form that actually runs. Three sites in
`tools/queue-kanban/` still hand a user the bare POSIX command:

- `tools/queue-kanban/verify.go:287` — the `Remedy:` string on a future-dated-timestamp finding:
  "re-stamp it with `date -u +%Y-%m-%dT%H:%M:%SZ` (the Timestamp rule in actions/work-reference.md)".
- `tools/queue-kanban/web/board.js:154` — the claim-stopwatch tooltip.
- `tools/queue-kanban/web/board.js:553` — the future-stamp data-warning text.

## Solution summary

Built in a git worktree as Builder B of REQ-085's live fan-out acceptance test — branch `worktree-agent-REQ-087-posix-only-timestamp-command`, commit `202ff3e`, integrated by the `--no-ff` merge `5cfe1b5` (range `3ccbf36..5cfe1b5`). `verify.go`'s remedy keeps a command but swaps the POSIX-only `date -u +…` for `queue-kanban now` — the Timestamp rule's own option 1, which prints the right shape on every platform and whose "only if already built" precondition is satisfied by the fact that the reader is looking at that binary's output. The three warning/tooltip strings in `board.js` (×2) and `model.go` drop to the target shape plus a citation of the rule. The `future_timestamp_test.go` assertion that pinned the old literal was repointed.

## What worked

- Re-running the grep across `tools/` instead of trusting the REQ's three-site list found the fourth site. That is now three sweep REQs in a row where the inventory was a floor.

## What didn't work

- Nothing failed, but the first instinct — follow the REQ's suggested split literally and paste the Windows one-liner into `verify.go` — would have satisfied the requirement while relocating the exact problem requirement 2 warns about. The better answer was already in the rule: option 1 exists precisely so callers do not have to branch on platform.

## Worth knowing

- `queue-kanban now` is the right timestamp instruction to put in **the tool's own output**, and only there. The rule's "only if the binary is already built" caveat is what normally makes option 1 unsafe to recommend blindly — and it cannot fail in a string printed by that binary. Anywhere else, cite the rule instead.

## Back-reference

See `do-work/archive/UR-015/REQ-087-board-and-verify-hand-users-the-posix-only-timestamp-command.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5cfe1b5`.
