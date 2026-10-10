#!/usr/bin/env bash
# REQ-693 pre-build baseline: the existing tests of every worktree dirt reader this REQ
# touches (worktree lifecycle, cleanup Pass 5, handoff survey, queue-kanban verify) pass
# at the base commit. Reads paths from whichever tree it runs in.
set -euo pipefail
repository_root="$(git rev-parse --show-toplevel)"
cli_module="$repository_root/skills/do-work/tools/do-work-cli"
board_module="$repository_root/skills/do-work-board/tools/queue-kanban"

expect_passes() {
  local label="$1" output="$2"
  shift 2
  local test_name
  for test_name in "$@"; do
    if ! grep -q -- "--- PASS: $test_name " <<<"$output"; then
      printf 'FAIL: %s did not report --- PASS: %s\n%s\n' "$label" "$test_name" "$output" >&2
      exit 1
    fi
  done
}

cleanup_tests=(
  TestWorktreeLifecycleNewMergeCleanupLeavesNothingBehind
  TestWorktreeCleanupRefusesDirtyWorktreeBeforeAnySideEffect
  TestMergedCleanBuilderWorktreeIsAutomaticButUnmergedNeedsExactConsent
  TestMergedCleanBuilderWorktreeRequiresSettledRequestEvidence
)
cleanup_pattern="^($(IFS='|'; echo "${cleanup_tests[*]}"))\$"
cleanup_output="$(cd "$cli_module" && go test -count=1 -v -run "$cleanup_pattern" ./internal/cleanup/ 2>&1)"
expect_passes cleanup "$cleanup_output" "${cleanup_tests[@]}"

handoff_tests=(TestHandoffSurveyNamesIntegrationAndDirtyCheckout TestHandoffStatusRowsAreTypedPerPathAndPreserveHostileNames)
handoff_pattern="^($(IFS='|'; echo "${handoff_tests[*]}"))\$"
handoff_output="$(cd "$cli_module" && go test -count=1 -v -run "$handoff_pattern" ./internal/corehelpers/ 2>&1)"
expect_passes handoff "$handoff_output" "${handoff_tests[@]}"

board_tests=(
  TestVerifyDoesNotAdvertiseADirtyMergedWorktreeAsFixable
  TestVerifyClassifiesWorktreeLeftoversByMergeState
  TestVerifyDoesNotAdvertiseAMergedWorktreeWithAnUnreadableStatusAsFixable
)
board_pattern="^($(IFS='|'; echo "${board_tests[*]}"))\$"
board_output="$(cd "$board_module" && go test -count=1 -v -run "$board_pattern" . 2>&1)"
expect_passes queue-kanban "$board_output" "${board_tests[@]}"

echo "REQ-693 preflight baseline: all named reader tests pass"
