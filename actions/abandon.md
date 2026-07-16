# Abandon Action

> **Part of the do-work skill.** Marks a REQ as won't-do: sets `status: cancelled`, records the reason, and archives the file — so the decision shows up with finished work on the board instead of haunting the queue as a warning.

Cancelling is a first-class outcome, not a deletion. The REQ file survives with its full trail of intent plus a `## Cancelled` section explaining why — six months later, "we decided not to do this, and here's why" is exactly as valuable as "we did this." The canonical status vocabulary (including where `cancelled` sits relative to `failed` and the terminal-success set) is defined in `actions/work-reference.md`'s Schema Read Contract → Terminal-resolved status set.

## When to Use

**Use when:**
- A pending REQ is no longer wanted — priorities changed, the need evaporated, or another REQ superseded it
- A `pending-answers` or blocked REQ isn't worth unblocking — the open questions aren't worth answering
- The user says "abandon", "cancel", "won't do", "drop", or "we're not doing this" about a specific REQ

**Do NOT use when:**
- The work was attempted and didn't succeed — that's `failed`; `actions/work.md` Step 8's failure classification handles it and spawns follow-ups
- The user wants to deactivate a running *pipeline* — that's `do-work pipeline abandon` (`actions/pipeline.md`), which flips pipeline state, not REQ status
- The user wants to defer a REQ for later — leave it `pending`; the queue is the backlog, and sitting in it costs nothing

## Input

`$ARGUMENTS`: one or more REQ IDs (`REQ-NNN`), optionally followed by free-text — the cancellation reason.

- `do-work abandon REQ-042` — cancel one REQ; the action asks for a one-line reason
- `do-work abandon REQ-042 superseded by REQ-051` — cancel with the reason inline (everything after the last REQ ID is the reason, applied to every listed REQ)
- `do-work abandon` (no ID) — list cancellable REQs (everything in `do-work/queue/` and `do-work/working/` with a non-terminal status) and ask which; never guess a target

## Steps

### Step 1: Locate and Gate

For each REQ ID, glob `do-work/queue/REQ-NNN-*.md`, `do-work/queue/REQ-NNN.md`, `do-work/working/REQ-NNN-*.md`, and `do-work/archive/**/REQ-NNN*.md`. Then gate on what you find:

- **Not found anywhere** → report `REQ-NNN: not found` and skip it.
- **Only in archive** → report its archive path and status — it's already terminal; nothing to cancel.
- **Status `completed` or `completed-with-issues`** → refuse: finished work is history, not a cancellation target. If the user wants it undone, that's a new capture.
- **Status `failed`** → report that it's already terminal; `do-work cleanup` will archive it. Cancelling would erase the failure signal.
- **Status `claimed`** → warn that a work loop may be mid-flight on it (one orchestrator per queue) and require an explicit extra confirmation before proceeding.
- **Status `reserved`** → warn that it's allocated to another worktree/cloud session (`reserved_for` names it) and require an explicit extra confirmation — that session may be mid-build on it.
- **Any other status** (`pending`, `pending-answers`, `blocked-*`, or unrecognized) → cancellable; continue.

### Step 2: Confirm the Decision

Show the user what's about to be cancelled — ID, title, current status, owning UR — for every target in one prompt (use your environment's ask-user prompt). If no reason was given in `$ARGUMENTS`, ask for a one-line reason in the same prompt; accept "no reason" but never invent one. Do not write anything until the user confirms.

### Step 3: Write the Cancellation

For each confirmed REQ:

1. Frontmatter: set `status: cancelled` and stamp `completed_at: <now, UTC ISO-8601>` — that timestamp is what places the card in the board's recently-done window. Leave `claimed_at`/`route` and every other field untouched; they're history.
2. Append to the body:

   ```markdown
   ## Cancelled

   - **When:** 2026-07-06T16:45:00Z
   - **Why:** [the user's reason, verbatim — or "no reason given"]
   - **Decided by:** user, via `do-work abandon`
   ```

Always write the canonical value `cancelled` — never `canceled`, `abandoned`, or `wont-do` (those are read-side aliases only; write paths emit canonical values per the Schema Read Contract).

### Step 4: Surface Dependents

Grep `do-work/queue/` and `do-work/working/` for REQs whose `depends_on` (or legacy `dependencies:`) lists a cancelled ID. A cancelled REQ does **not** satisfy dependency gating, so each dependent would sit blocked forever. For each dependent, ask the user to pick one:

- **Cascade** — abandon the dependent too (loop it back through Steps 1–3)
- **Re-point** — edit its `depends_on` to drop or replace the cancelled ID
- **Leave** — keep it blocked deliberately; it will show under blocked-by-dependencies until edited

Never cascade silently.

### Step 5: Archive

Move each cancelled REQ file out of the queue:

- If `do-work/archive/UR-NNN/` exists for its `user_request` → move it there.
- Otherwise → move it to `do-work/archive/` root (cleanup's Pass 2 consolidates later).
- **Collision guard:** if any `do-work/archive/**/REQ-NNN*.md` already exists, do NOT overwrite — leave the cancelled file in `do-work/queue/`, report the collision with both paths, and let the user resolve it (mirrors `actions/cleanup.md`'s duplicate handling).

### Step 6: Report

Summarize per REQ, note dependents and how each was dispositioned, and check the owning UR: if every sibling REQ is now terminally resolved (`completed`, `completed-with-issues`, or `cancelled`), say that `do-work cleanup` will close the UR.

## Output Format

```
Cancelled REQ-042 — [title]
  reason: superseded by REQ-051
  archived: do-work/archive/UR-012/REQ-042-slug.md
  dependents: REQ-047 re-pointed (depends_on: REQ-042 removed)

UR-012: all 3 REQs terminally resolved — `do-work cleanup` will close it.
```

## Rules

- **Never delete the REQ file.** Cancel + archive preserves the trail of intent; deletion destroys it.
- **Never cancel without confirmation** of the specific REQ IDs — this action removes items from the queue, and the queue is user intent.
- **Only the REQs the user named.** No opportunistic cancelling of stale-looking neighbors.
- **Write canonical `cancelled` only** — aliases are for reading hand-edited files, never for writing.
- Touch nothing beyond the target REQ files, their dependents' `depends_on` (when the user picks re-point), and the archive move.

## Common Rationalizations

| If you're thinking...                                   | STOP. Instead...                                            | Because...                                                                 |
| ------------------------------------------------------- | ----------------------------------------------------------- | -------------------------------------------------------------------------- |
| "`failed` is close enough for won't-do"                 | Use `cancelled`                                             | `failed` signals work that should have happened — it spawns follow-ups and holds the UR open; `cancelled` is the explicit no-follow-up decision |
| "Deleting the file is cleaner than archiving it"         | Set `cancelled`, append the reason, archive                  | The skill's primary value is the trail of intent — a recorded "no" included |
| "It's claimed but probably stale — cancel it quietly"    | Warn and get explicit confirmation first                     | Another orchestrator may be mid-flight on it; cancelling under it corrupts the run |
| "The queue is long — I'll cancel other stale REQs too"   | Cancel only the named REQs, mention candidates in the report | Staleness is the user's call; the queue is their backlog, not yours         |

## Red Flags

- A REQ file is gone from the repo after an abandon run — deletion happened instead of archival
- Frontmatter says `abandoned`, `canceled`, or `wont-do` — an alias leaked into a write path
- A cancelled REQ still sits in `do-work/queue/` with no reported archive collision — Step 5 was skipped
- A dependent REQ flipped to `cancelled` without the user choosing cascade
- The board shows the cancelled REQ under Needs input / Blocked — `completed_at` wasn't stamped or the status value drifted

## Verification Checklist

- [ ] Each cancelled REQ file lives under `do-work/archive/` (UR folder or root) — evidence: final file path in the report
- [ ] Frontmatter has `status: cancelled` + `completed_at`; body has a `## Cancelled` section carrying the user's reason verbatim
- [ ] Every dependent found in Step 4 was dispositioned by the user (cascade / re-point / leave) — evidence: one line per dependent in the report
- [ ] No file was deleted, and no file outside the named REQs (plus user-approved dependent edits) was modified
