# Builder Brief — REQ-430

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-430-couple-ur-closure-to-terminal-member-archival`
Branch: `worktree-agent-REQ-430-couple-ur-closure-to-terminal-member-archival`
Hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-430-handback.md`

## Ownership

Work only in the worktree above and commit on its branch. Treat all `do-work/` content in that worktree as stale and do not read or write it. The exact hand-back path is the sole allowed main-tree write. Never change version or changelog files.

## Request and frozen scope

Couple UR closure to successful archival of every concrete terminal member move required in the same cleanup plan. A refused member must leave both the member and active UR input in place, while unrelated safe groups still progress.

Only these files are writable:

- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go`

The exploration found that `BuildPlan` emits member moves and `CLOSE-UR` independently, while `ApplyPlan` admits groups after their own preflight. Add deterministic prerequisite group codes to operation groups, attach sorted unique concrete member-move codes to UR closure, and resolve direct/missing/transitive/duplicate/cyclic prerequisite refusal before forming eligible groups. Keep the existing single rollback-capable transaction and independent unrelated progress. Do not hard-code code-prefix relationships in the applier; the dependency seam should remain generic for REQ-431.

RED/GREEN: `TestURClosureWaitsForRequiredMemberArchival` must fail before production changes, then prove dirty-member refusal retains both inputs, names the blocker, and still applies an unrelated safe group. Also cover clean closure and stable prerequisite order.

Route B, `tdd: true`, `domain: general`.

## Required context

Read the general, coding-guardrails, communication-style, and testing crew files; then read `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and its lessons satellite. The full exploration is at the absolute main-tree path `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-430-exploration.md` and is read-only.

## Hand-back

Commit implementation on the branch. Write the hand-back with branch/hash, approach, exact RED/GREEN output, every changed file, checks and direct status, seams, and any `## Decisions` or `## Discovered Tasks`. Return only one status line once it exists.
