# REQ-502 builder brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-502-remove-enriched-checkpoint-entries-in-cleanup-mover`

Branch: `worktree-agent-REQ-502-remove-enriched-checkpoint-entries-in-cleanup-mover`

Main-tree hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-03-210100/REQ-502-handback.md`

Do not read or write the worktree's tracked `do-work/` snapshot. Do not write any main-tree path except the hand-back above. Commit project changes on the builder branch.

## Request

Remove the alternate cleanup mover's header-only checkpoint deletion. Every cleanup-package departure from `do-work/working/` must remove the complete matching own-label checkpoint entry: its `- REQ-NNN:` header and immediately following nonblank indented continuation lines. Preserve foreign-label entries and continuation bytes exactly. Locate the real `## In Progress (interrupted)` section by a whole heading line; inline or backticked mentions elsewhere are not entries. Reuse or align with canonical request-state semantics without creating a drifting second definition.

The named instance is `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` function `ownedCheckpointRemoval`. Extend `TestWorkingArchiveRemovesOnlyThisCheckoutCheckpointEntry` with enriched own and foreign entries. The test must fail before the fix for the correct assertion, then pass after it. Check for the same header-only pattern elsewhere and report, but do not fix unrelated sites.

Required context: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, its full `lessons-do-work-cli.md`, and the `alternate-writer-contract-drift` lesson. The original completed request is `do-work/archive/REQ-489-remove-whole-checkpoint-entries-on-departure.md` in the main tree. Follow the bug-fix spec and the general, coding-guardrails, backend, testing, and communication-style crew files.

## Hand-back format

Write the branch name, commit hash, full file manifest, RED then GREEN commands/results, lesson-read evidence, integration seams, a `## Decisions` section if any, and a `## Discovered Tasks` section if any. Do not edit the main-tree REQ.
