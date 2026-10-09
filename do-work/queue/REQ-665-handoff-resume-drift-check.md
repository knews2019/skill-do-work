---
id: REQ-665
title: 'Resuming from a handoff prints GO or NO-GO with one line per drift between the handoff and the live repo'
status: pending
created_at: 2026-10-09T21:16:07Z
user_request: UR-146
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-667]
related: [REQ-662, REQ-663, REQ-664, REQ-666, REQ-667]
batch: run-coordinate-mode
---
# Resuming From a Handoff Prints GO or NO-GO
## What
When a session starts from `do-work/RESTART-PROMPT.md`, or the user pastes a `Handoff:` line, check the handoff against the live repo before running anything, and print `GO` or `NO-GO` with one line per drift.
## Why
Report item A4 (UR-146 input, "What happened"): 3 sessions in two consumer repos (2026-09-27, 2026-09-28, 2026-10-02) asked by hand whether a handoff was on track before resuming it, for example "is this handoff on the right path? Handoff: do-work/RESTART-PROMPT.md (committed as …) Resume: do-work run --fan-out 1".
## Verified Facts (checked at 0.305.87)
- `skills/do-work/actions/restart-with-parallel-handoff.md:55`: the next session "re-checks all three conditions" before removing a REMOVABLE worktree. That is the only resume-time check.
- `restart-with-parallel-handoff.md:61-70`: the paste block holds one `advance REQ-NNN` line per claim and the resume command; the Reference section after `---` lists paths, merge ranges and worktree verdicts.
- Nothing compares the handoff's claimed REQs, branches and locks with the live repo.
## Detailed Requirements
1. Verify before running:
   - the handoff commit exists;
   - each `advance REQ-NNN` line names a REQ still in `do-work/working/`;
   - each worktree and `worktree-agent-*` branch the Reference lists exists in the state it names;
   - no lock is held by a dead process.
2. Print `GO` when nothing drifted, else `NO-GO` with one line per drift.
3. Wire the check into the resume path: the paste block that `restart-with-parallel-handoff.md` Step 4 writes runs it first.
4. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Read-only. The check never removes a worktree, branch or lock; REMOVABLE worktrees keep the existing three-condition rule (`:55`).
- Never `recover --take-over` (`restart-with-parallel-handoff.md:66`).
- No new REQ status and no new frontmatter field.
## Assumptions (recorded at capture, no questions asked)
- The checks are mechanical, so they become one do-work-cli command (name is the builder's choice) with Go tests; that is why `tdd: true`. If the builder keeps it prose, it records why and the REQ moves to `tdd: false`.
- The check applies to every handoff resume, not only coordinated runs; the report's A4 does not limit it to `--coordinate`. So it does not depend on REQ-662.
- The handoff commit is the commit that last changed `do-work/RESTART-PROMPT.md`, or the hash a pasted `Handoff:` line names.
- The Reference section is prose for humans. The command reads the parts it can parse (worktree paths, `worktree-agent-*` names, REQ ids) and reports a line it cannot parse as unreadable, never as GO. If parsing proves fragile, the builder may have Step 4 write a small machine-readable block, and records that choice; per the maintainer's earlier ruling, never grow a parser grammar for every way a model might format a path.
- "Locks" means the `full-gate.lock` files from REQ-667 (the coordinated-run rules); that is why this REQ depends on it.
## Dependencies
Depends on REQ-667 (the full-gate lock format the lock check reads).
## Builder Guidance
High certainty on the four checks and the GO / NO-GO output. Lower on how much of the Reference section to parse. Latitude: command name, output wording.
## Red-Green Proof
**RED prompt/case:** Write a handoff whose Reference names a `worktree-agent-REQ-N-x` worktree, delete that worktree, then resume.
**Why RED now:** no check exists; the session resumes without noticing the missing worktree.
**GREEN when:** the resume prints `NO-GO` with a drift line naming the missing worktree. An unchanged repo prints `GO`.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: a new do-work-cli command; families `silent-skip-reads-as-red` and `exact-basename-authority` fit the drift checks.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over budget; `slugged: partial`). Matching reason: changes the handoff action's paste-block contract.
## Full Context
See `do-work/user-requests/UR-146/input.md` for complete verbatim input (sections Request A4, What happened, Where the behaviour lives today, Proposed direction A4, Acceptance check). No queued candidate shares this root cause (the queue held REQ-654 to REQ-661, the ai-report and do-work-cli batches, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-coordinate-mode.md`, Request item A4: "A resume drift check when a session starts from a handoff: print go or no-go before running."*
