---
id: REQ-489
title: 'Remove whole checkpoint entries when a REQ leaves working'
status: pending
created_at: 2026-09-01T19:46:57Z
user_request: UR-083
addendum_to: REQ-440
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
status_changed_at: 2026-09-01T21:05:20Z
---

# Remove Whole Checkpoint Entries When a REQ Leaves Working

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

When the canonical `complete` (and, by the same code path, `fail` and the other departures from `do-work/working/`) removes this checkout's entry from `do-work/CHECKPOINT.md`'s `## In Progress (interrupted)` list, it deletes only the `- REQ-NNN: ...` header line. The indented `Last known state:`, `Key files being modified:`, and `Known issues:` lines that Step 10 enrichment adds beneath it stay behind as an unattributed orphan block.

Remove the entire entry: the header line plus every following indented continuation line up to the next `- ` entry, blank line at the section boundary, or heading.

## Context

Found during REQ-440's archive on 2026-09-01: after `complete REQ-440` the checkpoint kept REQ-440's three detail lines with no header, and the same orphan block from REQ-418's earlier archive was still present above it. The orchestrator removed both by hand at Step 10. Root cause: `checkpointWithoutClaim` in `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` drops only lines matching `- REQ-NNN:` with the writer label and keeps every following indented line.

## Requirements

- Entry removal deletes the header and all of its indented continuation lines.
- A checkpoint whose entry has no continuation lines behaves exactly as today.
- Foreign-label entries and their continuation lines are untouched.

## Red-Green Proof
**RED prompt/case:** A requeststate test that builds a checkpoint with an enriched own-label entry (header plus three indented detail lines), runs the departure removal, and asserts the section contains none of the four lines.
**Why RED now:** Only the header line is removed; the three indented lines remain.
**GREEN when:** The test passes and an existing removal test for a bare one-line entry still passes.
**Validation:** Discovered task from REQ-440; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions
- [x] I discovered this out-of-scope task while working on REQ-440: the canonical archive command leaves orphaned detail lines in the session checkpoint when it removes a finished REQ's in-progress entry. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

<!-- D-XX counter: none used. Next decision: D-01. -->

---
*Source: discovered task recorded during REQ-440 (work action Step 8).*


## Answer Notes

- 2026-09-01 - [ ] I discovered this out-of-scope task while working on REQ-440: the canonical archive command leaves orphaned detail lines in the session checkpoint when it removes a finished REQ's in-progress entry. Should I process this as a new task?: Confirmed: Yes, add to queue
> ```
> Add this checkpoint-cleanup fix to the queue for normal implementation. No scope is excluded.
> ```
