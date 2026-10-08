# REQ-643 Review Delta (1e6d9909..5b2c3e64)

Verdict: Pass. Score: 94.

Scope: the review fix for REQ-643 (the board shows an earmarked pending REQ under Pending → Earmarked instead of Ready). 8 files, comments and docs plus one tooltip string. Checked against `model.go` `bucketColumns` (lines 1723-1731).

## Checks

- C1 F1–F7 fixed and true against the code: board-cards.js tooltip and comment (F1), board.md:120 (F2), prime-do-kanban.md:15/39/48 (F3), generate.go:184 (F4), board.css:1037/1302 (F5), work-reference.md:120/248 and model.go:224 (F6), work-guide.md:119 (F7). Each restatement scopes Earmarked to a pending REQ with no unmet dependency, which matches the code.
- C2 No restatement says an assigned REQ waiting on a dependency is Earmarked. board.md:120 orders the rule correctly ("an unmet dependency sends a REQ to Waiting, otherwise a non-empty `assigned_to` sends it to Earmarked").
- C3 The badge class, label `assigned`, and visible text (`truncateBadgeText(request.assignedTo, 18)`) are unchanged. Only `title` changed.
- C4 The model.go `AssignedTo` comment (lines 183-192) and the work-reference.md:114 `assigned_to:` line agree. Neither changed in this delta. Both last changed together in 3393d267.
- C5 work-guide.md:119 keeps REQ-644's closing sentence (earmark vs. operator-blocked) and adds "shows it under Pending → Earmarked, not Ready, with an `assigned` badge". It reads well.
- C6 prime-do-kanban.md has no "do not restate" rule. The delta edits Read first, Traps and Stakes in place and adds no new copy of a rule. The model.go:1705 edit also repairs a misplaced `//` inside the comment.

## Findings

- F8 `web/board-cards.js:238-239` tooltip: "The board only groups it under Pending → Earmarked" is shown on every assigned card. For an assigned pending REQ that waits on a dependency, the card is under Waiting, so the sentence overstates there. Possible wording: "...a pending card with no unmet dependency sits under Pending → Earmarked". — impact-user-visible → report only
- F9 `model.go:1706` comment line is about 120 characters after the rewrap. Cosmetic only. — impact-negligible → report only
