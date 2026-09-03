---
id: REQ-527
title: 'Teach cleanup Pass 5 that merged is not finished'
status: pending
priority: now
created_at: 2026-09-03T02:00:00Z
user_request: UR-086
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-458]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-458]
review_generated: true
addendum_to: REQ-458
write_set:
  - skills/do-work/actions/cleanup.md
  - skills/do-work/docs/cleanup-guide.md
---

# Teach Cleanup Pass 5 That Merged Is Not Finished

## What

REQ-458 stopped `verify` from advertising an active builder's worktree as a fixable leftover. `cleanup` Pass 5 still decides on merged-ness alone, so the two now contradict each other in the opposite direction: verify protects the worktree, cleanup removes it.

## Instances

- `skills/do-work/actions/cleanup.md:140` — "Both commands succeeding means the work was merged and the leftover was pure residue: report `Removed merged worktree <name>`." For a clean worktree belonging to a run that has not reached review, `git worktree remove` and `git branch -d` both succeed, so Pass 5 removes an in-flight builder's worktree with no consent gate.
- `skills/do-work/docs/cleanup-guide.md:29` — "Already-merged leftovers are removed automatically." Same premise, user-facing.
- `skills/do-work/docs/cleanup-guide.md:56` — same contract restated.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding F2** — `impact-user-visible` — from REQ-458's independent review (Approve with follow-ups, 73%), found by its restatement sweep. The builder's own sweep missed it; the reviewer swept the whole repository for consumers of the redefined contract and found three.
- REQ-458's Scope declared `cleanup.md` as "will NOT touch" on the premise that "Pass 5's consent-gated behavior is already correct". That premise holds for **dirty** worktrees, where non-forced `git worktree remove` refuses. It is false for **clean in-flight** ones, where both commands succeed.

## Detailed Requirements

- Pass 5 must not mechanically remove a `worktree-agent-*` worktree whose REQ is still in `do-work/working/`. Both commands succeeding is evidence the work is merged, never evidence the run is over.
- Use the same evidence REQ-458 established — the leftover name's `REQ-NNN` id and that REQ's tree section — so verify and cleanup agree on what "finished" means rather than each deciding separately.
- Introduce no liveness signal: no heartbeat, lock, PID probe, mtime heuristic, claim registry, or time threshold. REQ-073's ban applies here exactly as it does to verify.
- Keep removing genuinely finished, clean, merged residue mechanically; this must not turn Pass 5 into a prompt for the ordinary case.
- Keep `cleanup.md` and `cleanup-guide.md` consistent with each other and with `verify`'s corrected categories.

## Constraints

- Never `-D`, never `--force`. The existing non-forced commands stay non-forced.
- Do not relocate the consent gate; extend what routes to it.

## Dependencies

Depends on REQ-458, which establishes the evidence and the corrected verify categories this must agree with.

## Red-Green Proof

**RED prompt/case:** Run cleanup Pass 5 against a tree holding a `worktree-agent-REQ-NNN-*` worktree that is clean and merged while `do-work/working/REQ-NNN-*.md` still exists. Today both commands succeed and the worktree is reported as removed residue.
**Why RED now:** Pass 5 routes on the two commands' exit status alone, which answers "is it merged" and not "is the run over".
**GREEN when:** That worktree is left in place and reported as belonging to an unfinished run; a clean merged worktree whose REQ is archived is still removed mechanically; and `cleanup-guide.md` no longer says already-merged leftovers are removed automatically without qualification.

---
*Source: REQ-458 independent review finding F2.*
