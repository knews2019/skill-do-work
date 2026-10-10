#!/usr/bin/env bash
# Pre-flight baseline for REQ-690 (run before the build, at main HEAD): the existing tests that
# guard what the build touches pass today. Board: the open-work text digest, the column
# partition, the earmarked placement and the request-activity correlation. Core CLI: the
# result model's text/JSON parity and exit codes (it gains a sibling field) and the
# frozen-estimate read (may be exported). Focused names only, no whole package.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
board_dir="$root/skills/do-work-board/tools/queue-kanban"
cli_dir="$root/skills/do-work/tools/do-work-cli"

board_output="$(go test -C "$board_dir" -count=1 -v -run '^(TestOpenWork|TestBucketColumns$|TestAssignedPendingRequestIsEarmarkedNotReady$|TestSyntheticColumnBucketing$|TestRequestActivity)' . 2>&1)" || {
  printf '%s\n' "$board_output"
  exit 1
}
grep -q -- '--- PASS: TestBucketColumns' <<<"$board_output" || { printf 'missing PASS: TestBucketColumns\n%s\n' "$board_output"; exit 1; }
grep -q -- '--- PASS: TestOpenWorkDigestListsClaimedRequestsWithTitles' <<<"$board_output" || { printf 'missing PASS: open-work digest\n%s\n' "$board_output"; exit 1; }
grep -q -- '--- PASS: TestAssignedPendingRequestIsEarmarkedNotReady' <<<"$board_output" || { printf 'missing PASS: earmarked placement\n%s\n' "$board_output"; exit 1; }
if grep -q -- '--- SKIP' <<<"$board_output"; then printf 'a board test skipped\n%s\n' "$board_output"; exit 1; fi

cli_output="$(go test -C "$cli_dir" -count=1 -v -run '^(TestSimpleSelectionRetainsSpecializedVetoesAndFrozenEstimate|TestRenderersUseOneNormalizedResult|TestExactTextAndAuditJSONShareOneResult|TestOutcomeExitCodes)$' ./internal/resultmodel/ ./internal/nextselection/ 2>&1)" || {
  printf '%s\n' "$cli_output"
  exit 1
}
for test_name in TestSimpleSelectionRetainsSpecializedVetoesAndFrozenEstimate TestRenderersUseOneNormalizedResult TestExactTextAndAuditJSONShareOneResult TestOutcomeExitCodes; do
  grep -q -- "--- PASS: $test_name" <<<"$cli_output" || { printf 'missing PASS: %s\n%s\n' "$test_name" "$cli_output"; exit 1; }
done
if grep -q -- '--- SKIP' <<<"$cli_output"; then printf 'a core test skipped\n%s\n' "$cli_output"; exit 1; fi
printf 'REQ-690 pre-flight baseline: PASS\n'
