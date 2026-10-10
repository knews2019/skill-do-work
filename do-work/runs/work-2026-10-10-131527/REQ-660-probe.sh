#!/usr/bin/env bash
# GREEN checks for REQ-660 (integrator test gate): the five worktree lifecycle tests pass (a skip or a missing test is not a pass), the command is registered, and the lifecycle source never passes --force or -D to git. Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module="$root/skills/do-work/tools/do-work-cli"
source_file="$module/internal/cleanup/worktree_lifecycle.go"
probe_tests=(
  TestWorktreeLifecycleNewMergeCleanupLeavesNothingBehind
  TestWorktreeMergeRefusesBuilderCommitUnderDoWork
  TestWorktreeMergeReportsEmptyHandBackWithoutCommit
  TestWorktreeCleanupRefusesDirtyWorktreeBeforeAnySideEffect
  TestWorktreeNewAppendsNumericSuffixOnCollision
)
[ -f "$source_file" ] || { echo "missing $source_file"; exit 1; }
output="$(go test -C "$module" -count=1 -v -run "^($(IFS='|'; echo "${probe_tests[*]}"))\$" ./internal/cleanup/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
source_text="$(cat "$source_file")"
if grep -q -e '"--force"' -e '"-D"' <<<"$source_text"; then
  echo "worktree_lifecycle.go passes --force or -D to git"; exit 1
fi
help_output="$(go run -C "$module" ./cmd/do-work-cli --format json help 2>&1)" || { printf '%s\n' "$help_output"; exit 1; }
grep -q '"worktree"' <<<"$help_output" || { echo "do-work-cli help does not list the worktree command"; exit 1; }
echo "REQ-660 probe: all checks passed"
