# REQ-433 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-433-split-misplaced-ur-partial-merge-conflicts`
Branch: `worktree-agent-REQ-433-split-misplaced-ur-partial-merge-conflicts`
REQ: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-433-split-misplaced-ur-partial-merge-conflicts.md`
Prior implementation: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/REQ-409-implement-safe-cleanup.md`
Integration owner: main checkout; do not edit `do-work/`, release metadata, changelogs, or version files.

## Frozen write set

- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go`

## Required method

- Read the REQ, prior implementation, listed CLI prime and lessons, root/project instructions, and applicable general/testing/coding guardrails.
- Write the captured conflicting-`input.md` plus nonconflicting sibling-REQ regression first and observe a real RED before production changes.
- Split only the misplaced archived-UR directory path into deterministic per-file conflict domains. Preserve symlink refusal, destination non-overwrite, exact evidence, existing nested-do-work relocation behavior, and cleanup transaction mechanics.
- Stay exactly inside the two-file scope. Run the focused RED/GREEN test, the full cleanup package, full CLI module tests and vet, exact Go 1.25 compatibility, qualification/scope checks, and diff checks.
- Commit on the branch, keep the worktree clean, and write the handback to `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-433-handback.md`.

## Handback

Record the full commit hash, RED/GREEN evidence, every changed file, commands and exit codes, deterministic group/code choices, decisions, and discoveries. Do not create or edit queue files.
