# Fan-Out Run Manifest — work-2026-08-04-000254

Owner-written only. Builders never write this file.

Purpose: REQ-085's live two-builder acceptance test of REQ-073's fan-out dispatch.
Integration branch: `claude/do-work-run-ed81l6` @ `012469b` at run start.
Worktree parent: `/home/user/skill-do-work-worktrees/` (outside the repo working tree).

| REQ | Builder | Operative name | Handback file | Landed |
| --- | --- | --- | --- | --- |
| REQ-086 | A | `worktree-agent-REQ-086-in-progress-record-unstated` | `REQ-086-handback.md` | pending |
| REQ-087 | B | `worktree-agent-REQ-087-posix-only-timestamp-command` | `REQ-087-handback.md` | pending |

## Non-overlap justification (REQ-085 requirement 1)

Declared `write_set`s are disjoint:

- REQ-086 — `actions/cleanup.md`, `actions/forensics.md`, `docs/work-guide.md`
- REQ-087 — `tools/queue-kanban/verify.go`, `tools/queue-kanban/web/board.js`

Confirmed **by reading both REQs**, not by the badge: REQ-086 states a bookkeeping rule
(`CHECKPOINT.md`'s in-progress entry is dropped whenever a REQ leaves `working/`) at three consumer
sites; REQ-087 rewords three display strings that hand the user a POSIX-only `date` command. No shared
concept, no shared file, no shared section. `actions/forensics.md` is touched by REQ-086 at Check 1;
REQ-083/REQ-084 changed Check 14 earlier in this session, but that is committed main-tree history, not
a concurrent sibling.

## Dispatch mechanism used

Both builders driven **by hand in the owner's session** — separate worktrees, separate branches,
sequential builds, serial integration. REQ-073 requirement 7 leaves the dispatch mechanism unspecified.
See REQ-085's `## Testing` → *What this run did not cover* for what that does and does not exercise.
