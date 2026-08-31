---
id: UR-086
title: 'Correct active worktree leftover classification'
created_at: 2026-08-31T21:38:14Z
requests: [REQ-458]
word_count: 56
---

# Correct Active Worktree Leftover Classification

## Full Verbatim Input

> ````text
> The board labels active builder worktrees as merged-worktree-leftover [fixable]. REQ-412 demonstrated a merged branch with uncommitted implementation changes, while REQ-436 remained claimed before review/remediation completed. “Merged” alone does not establish “leftover” or safe mechanical removal. Preserve the existing no-liveness-signal decision, but classify dirty or unfinished-run worktrees as present and non-fixable. Surface-cost: N/A — direct classification correction.
> ````

---
*Captured: 2026-08-31T21:38:14Z*
