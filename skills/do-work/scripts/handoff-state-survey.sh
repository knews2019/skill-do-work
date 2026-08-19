#!/usr/bin/env bash
# Read-only git survey for actions/restart-with-parallel-handoff.md.
# Emits the facts the handoff classifies: recent history, every worktree, which
# builder branches are already merged, and the dirty state of every checkout.
# Prints only; changes nothing. Judgment (worktree verdicts, critical path,
# what may run concurrently) stays in the action.
set -uo pipefail

if [ "$#" -gt 1 ]; then
  printf 'Usage: %s [integration-branch]\n' "$0" >&2
  exit 2
fi

# The integration branch is an argument because a consumer repo may not call it
# "main". Detect rather than assume, and say which one was used.
integration_branch="${1:-}"
if [ -z "$integration_branch" ]; then
  for candidate_branch in main master trunk; do
    if git rev-parse --verify --quiet "refs/heads/$candidate_branch" >/dev/null; then
      integration_branch="$candidate_branch"
      break
    fi
  done
fi
if [ -z "$integration_branch" ] \
  || ! git rev-parse --verify --quiet "refs/heads/$integration_branch" >/dev/null; then
  printf 'No integration branch found. Pass one: %s <branch>\n' "$0" >&2
  exit 2
fi

printf '=== integration branch ===\n%s\n\n' "$integration_branch"

printf '=== recent history ===\n'
git log --oneline -15
printf '\n'

printf '=== worktrees ===\n'
git worktree list
printf '\n'

# --list with a glob does the filtering git-side. A `| grep worktree-agent` pipe
# would exit 1 on the ordinary no-builders case and, under pipefail, hand that
# status to any caller reading the exit code as a verdict.
printf '=== builder branches already merged into %s ===\n' "$integration_branch"
merged_builder_branches="$(git branch --merged "$integration_branch" --list 'worktree-agent-*')"
if [ -n "$merged_builder_branches" ]; then
  printf '%s\n' "$merged_builder_branches"
else
  printf '(none)\n'
fi
printf '\n'

printf '=== unmerged builder branches ===\n'
unmerged_builder_branches="$(git branch --no-merged "$integration_branch" --list 'worktree-agent-*')"
if [ -n "$unmerged_builder_branches" ]; then
  printf '%s\n' "$unmerged_builder_branches"
else
  printf '(none)\n'
fi
printf '\n'

# --untracked-files=all, never bare --porcelain: the short form collapses a wholly
# untracked directory into one row, and the handoff must name uncommitted files
# individually ("0 commits" reads as "nothing there" and is often wrong).
printf '=== dirty state, every checkout ===\n'
while IFS= read -r worktree_line; do
  case "$worktree_line" in
    'worktree '*) ;;
    *) continue ;;
  esac
  worktree_path="${worktree_line#worktree }"
  printf -- '--- %s\n' "$worktree_path"
  if [ ! -d "$worktree_path" ]; then
    printf '(path missing — prunable)\n'
    continue
  fi
  worktree_status="$(git -C "$worktree_path" status --short --untracked-files=all 2>&1)"
  if [ -n "$worktree_status" ]; then
    printf '%s\n' "$worktree_status"
  else
    printf '(clean)\n'
  fi
done <<EOF
$(git worktree list --porcelain)
EOF
