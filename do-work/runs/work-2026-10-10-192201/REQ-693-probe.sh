#!/usr/bin/env bash
# REQ-693 GREEN probe: the two tests this REQ adds pass. One pins the missing-folder path
# of `worktree status` and `worktree cleanup`; the other pins that cleanup Pass 5 reads a
# worktree whose only untracked path is an owned do-work/worktree-links symlink as clean
# and removes it. At the base commit neither test exists, so this probe fails there.
set -euo pipefail
repository_root="$(git rev-parse --show-toplevel)"
cli_module="$repository_root/skills/do-work/tools/do-work-cli"

new_tests=(
  TestWorktreeStatusAndCleanupSurviveMissingWorktreeFolder
  TestPass5RemovesMergedWorktreeWhoseOnlyDirtIsAnOwnedLink
)
test_pattern="^($(IFS='|'; echo "${new_tests[*]}"))\$"
test_output="$(cd "$cli_module" && go test -count=1 -v -run "$test_pattern" ./internal/cleanup/ 2>&1)"
for test_name in "${new_tests[@]}"; do
  if ! grep -q -- "--- PASS: $test_name " <<<"$test_output"; then
    printf 'FAIL: no --- PASS: %s line\n%s\n' "$test_name" "$test_output" >&2
    exit 1
  fi
done
echo "REQ-693 GREEN: missing-folder and owned-link tests pass"
