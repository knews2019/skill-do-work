# REQ-431 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-431-couple-document-rewrites-to-owning-moves`
Branch: `worktree-agent-REQ-431-couple-document-rewrites-to-owning-moves`
REQ: `do-work/working/REQ-431-couple-document-rewrites-to-owning-moves.md`
Exploration: `do-work/runs/work-2026-08-31-165510/REQ-431-exploration.md`
Integration owner: main checkout; do not edit `do-work/`, release metadata, changelogs, or version files.

## Frozen write set

- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go`
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go`

## Required method

- Read the REQ, exploration report, CLI prime and lessons, and applicable general/testing/coding guardrails.
- Add apply-level tests first for both refusal directions and all-safe composition; capture the intended RED before production changes.
- Implement per-group typed rewrite operations that retain only their owning moves and rewrite current document bytes at apply time.
- Preserve anchors, unchanged filename-only mentions, exact group preflight/rollback targets, unrelated safe progress, and the existing union transaction.
- Stay strictly inside the frozen write set. Run focused/package tests, `go vet ./...`, full module tests, exact Go 1.25 compatibility, and `git diff --check`.
- Review every changed file, commit the implementation, keep the worktree clean, and write the handback to the main checkout path `do-work/runs/work-2026-08-31-165510/REQ-431-handback.md`.

## Handback

Record branch and full commit hash, approach, RED/GREEN output, changed files, all checks with exit results, data-flow seams, decisions, and discovered tasks. Do not create or edit queue files.
