#!/usr/bin/env bash
# Pre-flight baseline for REQ-691 (run before the build): the board's existing commit-to-REQ correlation tests,
# which the new request-commits helper extraction must keep green, pass today, and the shipped-prose contract
# tests that will read the new trace action files pass today. Reads paths relative to the tree it runs in.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
correlation_tests=(
  TestRequestActivityAttributesACommitThatTouchesTheRequestFileWithoutAPrefix
  TestRequestActivityAttributesEveryPrefixTokenInASubject
  TestRequestPathPatternMatchesOnlyRequestFilesAndRunArtifacts
)
output="$(go test -C "$root/skills/do-work-board/tools/queue-kanban" -count=1 -v -run "^($(IFS='|'; echo "${correlation_tests[*]}"))\$" . 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for correlation_test in "${correlation_tests[@]}"; do
  grep -q -- "^--- PASS: $correlation_test " <<<"$output" || { printf '%s\n' "$output"; echo "$correlation_test did not pass (skipped or missing)"; exit 1; }
done
bash "$root/_dev/tests/shipped-package-reference-contract.sh"
bash "$root/_dev/tests/action-shell-blocks.sh"
echo "REQ-691 pre-flight baseline: all checks passed"
