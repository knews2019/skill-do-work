#!/usr/bin/env bash
# Fixture execution proofs for show-commit-diff.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# show-commit-diff: a real merge must expose first-parent file changes.
merge_repo="$fixture_root/merge-repo"
fixture_repo_init "$merge_repo"
printf 'base\n' > "$merge_repo/base.txt"
fixture_repo_commit_all "$merge_repo" base
git -C "$merge_repo" checkout -qb feature
printf 'feature\n' > "$merge_repo/feature.txt"
fixture_repo_commit_all "$merge_repo" feature
git -C "$merge_repo" checkout -q main
printf 'main\n' > "$merge_repo/main.txt"
fixture_repo_commit_all "$merge_repo" main
git -C "$merge_repo" merge --no-ff -qm merge feature
merge_commit="$(git -C "$merge_repo" rev-parse HEAD)"
merge_output="$(cd "$merge_repo" && "$core_scripts/show-commit-diff.sh" "$merge_commit" 2>&1)" || fail_case 'show-commit-diff real-merge case returned nonzero'
printf '%s' "$merge_output" | grep -q 'feature.txt' || fail_case 'show-commit-diff real-merge case hid the merged file'

prescribed_shell_finish
