# Prose Backlog

Prose-only discrepancies in this project's own text, one line each. Not a REQ and not in the pipeline — appended by the Fold-First Rule's prose-only destination (`actions/capture-reference.md`), drained by an ordinary REQ that fixes items and deletes their lines.

- [ ] `skills/do-work/scripts/repair-req-timestamps.sh:7`, `:22`, `:136` and `skills/do-work/actions/work-reference.md:285`: all four cite "forensics Check 11" for the future-dated-timestamp check; Check 11 is *Unrecognized Status Vocabulary* and the future-stamp check is `### 12. Future-Dated Timestamps` (`skills/do-work/actions/forensics.md:156`) — verified 2026-08-20. (found by REQ-272 / UR-063)
- [ ] `skills/do-work-board/tools/queue-kanban/open_work.go:23` and `testing.go:42`: both say the board tool has "two write surfaces"; it has three, and `frontmatter_cli.go:34` already says "exactly three" — verified 2026-08-20. (found by REQ-273 / UR-063)
- [ ] `skills/do-work-board/actions/board.md:140`: the verification checklist names four modes ("serve / static / summary / open-work") while the Input table at `:29` now lists five — `verify` was added in 0.222.3 and the checklist gloss was not updated. Verified 2026-08-21. (found by REQ-283 / UR-057)
