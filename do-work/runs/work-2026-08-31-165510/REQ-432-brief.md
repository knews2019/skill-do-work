# REQ-432 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-432-enforce-commit-guard-for-consumed-scratch`
Branch: `worktree-agent-REQ-432-enforce-commit-guard-for-consumed-scratch`
REQ: `do-work/working/REQ-432-enforce-commit-guard-for-consumed-scratch.md`
Integration owner: main checkout; do not edit `do-work/`, release metadata, changelogs, or version files.

## Goal

Enforce cleanup's empty-index precondition before every commit-mode mutation, including the narrow non-rollback consumed-scratch deletion path, while keeping non-commit scratch eligibility unchanged.

## Required method

- Read the REQ, CLI prime and lessons, and applicable general/testing/coding guardrails.
- Add the exact staged-unrelated-file plus untracked-consumed-run regression first; capture a real RED in which commit mode deletes scratch or succeeds incorrectly.
- Make the smallest coherent preflight correction. Preserve rooted containment, exact inventory, consumed-manifest checks, and non-commit cleanup behavior.
- Assert truthful refusal evidence and byte-for-byte scratch retention before mutation.
- Run focused/package/full module tests, `go vet ./...`, exact Go 1.25 compatibility, and `git diff --check`.
- Review the diff, commit on the branch, keep the worktree clean, and write the handback to main checkout path `do-work/runs/work-2026-08-31-165510/REQ-432-handback.md`.

## Handback

Record full commit, approach, RED/GREEN output, changed files, checks and exits, seams, decisions, and discoveries. Do not create or edit queue files.
