# Hand-back — REQ-086

**Branch:** `worktree-agent-REQ-086-in-progress-record-unstated`
**Commit:** `0e04b4d`

## File manifest

- `actions/cleanup.md` (modified) — Pass 0 step 5 now states that moving a terminal REQ out of
  `working/` drops its `## In Progress (interrupted)` entry, citing the canonical home rather than
  restating the rule (requirement 4).
- `actions/forensics.md` (modified) — Check 1's suggested manual reset gained the same clause, in the
  sentence that already describes the move back to `do-work/queue/`.
- `docs/work-guide.md` (modified) — both sentences corrected. Line 66 now says the checkpoint is
  written at claim time and refreshed at session end, and says what that buys the user (a crashed run
  picks its own work back up). Line 119 no longer implies the file appears only "before stopping".

## Integration seams

**None.** All three files are leaf prose edits in this REQ's declared `write_set`; no shared registry,
template, or export barrel is involved.

## Notes for the owner

- All three sites were judged to need the clause. The Constraints allowed skipping any site where
  REQ-077's canonical list was already enough on its own; none qualified — Pass 0 and Check 1 each
  describe a move without mentioning the entry, and a reader following either procedure would leave the
  entry behind.
- Requirement 4 honored: each site *cites* `actions/work-reference.md` → **In-Progress Record (Step 2)**
  and none restates the trigger condition.
- No version bump and no `CHANGELOG.md` entry, per `CLAUDE.md` § Before Every Commit's scope clause —
  those are the integrator's, and a builder bumping either would race its sibling.
