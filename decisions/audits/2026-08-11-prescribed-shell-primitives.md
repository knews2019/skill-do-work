# Prescribed Shell Primitive Inventory

**Ran:** 2026-08-11 · **Scope:** shipped `skills/` tree · **Method:** grep-based inventory followed by reading every match; paths, not line numbers, are the durable coordinates.

## Disposition

| Primitive | Canonical shipped home | Former rationale sites | Divergence / execution-only disposition |
|---|---|---|---|
| Per-file untracked inventory | `skills/do-work/docs/prescribed-shell-primitives.md` → Per-file untracked inventory; executable source remains `tools/checks/uncommitted-inventory.sh` | `skills/do-work/actions/commit.md`, `skills/do-work-toolbox/actions/inspect.md`, `skills/do-work-toolbox/actions/stray-check.md` | REQ-121's NUL/rename/secret-aware variant wins. `work.md`, `tidy-repo.md`, preflight, and queue-kanban uses are narrower executions, not fallback copies. |
| Merge-aware commit diff | core guide → Merge-aware commit diff | `skills/do-work/actions/review-work.md`, `skills/do-work-toolbox/actions/ai-report.md`, three template/instruction sites in `skills/do-work-toolbox/actions/present-work.md` | The fixed form is quoted `'<commit>^2'` plus `git show --first-parent -m`; worktree orchestration keeps its explicit `<pre>..<merge_hash>` range. |
| Commit file listing | core guide → Commit file listing | explanatory copy in `skills/do-work-toolbox/actions/ai-report.md` | `skills/do-work/tools/checks/blanked-req-scan.sh` is executable use/comment, not a second policy home. `git diff-tree ... -r -m` wins over message-bearing `git show --name-only`. |
| Local Git ignore | core guide → Local Git ignore | long copies in the core, knowledge, and toolbox `crew-members/background-agents.md`; rationale in board/install actions | Commands in board/Just/install/setup-memory remain at their execution sites. The fixed variant uses `git rev-parse --git-path`, `**/`, and a separate tracked-file check when “never committed” is the requirement. |
| Atomic download publication | core guide → Atomic download publication | explanatory and Red-Flag copies in `skills/do-work-toolbox/actions/install.md` | Download commands in install tables and shipped updater/installer scripts are required uses. The fixed temp-download, rename-on-success, cleanup-plus-`false` shape wins. |
| Raw text before shell quoting | core guide → Raw text before shell quoting | `skills/do-work-knowledge/actions/memory-reference.md` | `memory.md` already points at its reference file; query-specific tokenization remains there, while the injection rationale moves to the guide. |
| Diff output filtering | core guide → Diff output filtering | no duplicate shipped action statement found | The maintainer trap is now shipped once. Existing diffs keep caller-specific filters; no command required rewriting in this pass. |
| State across command blocks | core guide → State across command blocks | general rationale in `work.md`, `work-reference.md`, `commit.md`, and toolbox `inspect.md` | Deterministic re-derivation and literal merge endpoints remain at callers; only the repeated general explanation moves. |

## Verification searches

The regression probe `_dev/tests/prescribed-shell-canonicalization.sh` owns the executable search. It requires the guide headings and former-site pointers, then rejects the known high-risk rationale phrases outside the guide. `_dev/tests/action-shell-blocks.sh` separately parses every surviving fenced shell block.

Execution count is not duplication count. A command that must run in three callers remains three commands; this audit consolidates the shared failure-mode statement so a future semantic correction has one source.

