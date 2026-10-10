#!/usr/bin/env bash
# GREEN checks for REQ-690 (do-work status and do-work-cli run-status), for the integrator's
# test gate. Reads paths from the tree it runs in. Fails before the build: the action, the
# routing row, the runstatus package and the named tests do not exist yet.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
board_dir="$root/skills/do-work-board/tools/queue-kanban"
cli_dir="$root/skills/do-work/tools/do-work-cli"
status_dir="$cli_dir/internal/runstatus"
skill_file="$root/skills/do-work/SKILL.md"

[ -f "$root/skills/do-work/actions/status.md" ] || { printf 'missing skills/do-work/actions/status.md\n'; exit 1; }
[ -d "$status_dir" ] || { printf 'missing %s\n' "$status_dir"; exit 1; }

# Routing: one row routes status phrases to status.md, above the clarify row, and the bare
# word blocked stays on clarify.
skill_text="$(cat "$skill_file")"
status_row="$(grep -n -F './actions/status.md' <<<"$skill_text" || true)"
clarify_row="$(grep -n -F './actions/clarify.md' <<<"$skill_text" || true)"
[ -n "$status_row" ] || { printf 'SKILL.md has no row routing to ./actions/status.md\n'; exit 1; }
[ -n "$clarify_row" ] || { printf 'SKILL.md lost its clarify row\n'; exit 1; }
for phrase in '`status`' 'status and ETA' 'is it stuck' 'how do I unblock'; do
  grep -q -F -- "$phrase" <<<"$status_row" || { printf 'status row lacks %s\n' "$phrase"; exit 1; }
done
grep -q -F -- '`blocked`' <<<"$clarify_row" || { printf 'clarify row lost the blocked trigger\n'; exit 1; }
if grep -q -F -- '`blocked`' <<<"$status_row"; then printf 'status row took the bare blocked trigger\n'; exit 1; fi
[ "${status_row%%:*}" -lt "${clarify_row%%:*}" ] || { printf 'status row is not above the clarify row\n'; exit 1; }

# Read-only: no deletion, signal or process lookup in the new run-status Go files.
forbidden="$(grep -rnE 'os\.Remove|syscall\.Kill|FindProcess' "$status_dir" || true)"
[ -z "$forbidden" ] || { printf 'forbidden call in runstatus:\n%s\n' "$forbidden"; exit 1; }

unformatted="$(gofmt -l "$status_dir" "$board_dir"/*.go "$cli_dir/internal/resultmodel" "$cli_dir/cmd/do-work-cli")"
[ -z "$unformatted" ] || { printf 'gofmt -l lists:\n%s\n' "$unformatted"; exit 1; }
go vet -C "$cli_dir" ./internal/runstatus/ ./cmd/do-work-cli/

cli_output="$(go test -C "$cli_dir" -count=1 -v -run '^TestRunStatus' ./internal/runstatus/ 2>&1)" || { printf '%s\n' "$cli_output"; exit 1; }
for test_name in \
  TestRunStatusEtaCountsDownFromTheFrozenEstimate \
  TestRunStatusLandedHandbackIsClassC3 \
  TestRunStatusUnmetDependencyIsClassC5QuotingTheWaitingReason \
  TestRunStatusEarmarkedRequestIsClassC6QuotingTheEarmarkedReason \
  TestRunStatusStaleClaimNamesRecoveryOnlyInTheStopReason \
  TestRunStatusListsRunLocalLockFileAndLeavesIt \
  TestRunStatusJSONRowsUseTheCommandFindingShape \
  TestRunStatusWatchFitsTwentyLinesForEightOpenRequests \
  TestRunStatusReportsQueueRowsWithoutARunDirectory; do
  grep -q -- "--- PASS: $test_name" <<<"$cli_output" || { printf 'missing PASS: %s\n%s\n' "$test_name" "$cli_output"; exit 1; }
done
if grep -q -- '--- SKIP' <<<"$cli_output"; then printf 'a runstatus test skipped\n%s\n' "$cli_output"; exit 1; fi

board_output="$(go test -C "$board_dir" -count=1 -v -run '^(TestOpenWork|TestBucketColumns$|TestAssignedPendingRequestIsEarmarkedNotReady$)' . 2>&1)" || { printf '%s\n' "$board_output"; exit 1; }
for test_name in TestOpenWorkJSONPlacesWaitingAndEarmarkedTicketsWithTheirReasons TestOpenWorkDigestListsClaimedRequestsWithTitles TestBucketColumns TestAssignedPendingRequestIsEarmarkedNotReady; do
  grep -q -- "--- PASS: $test_name" <<<"$board_output" || { printf 'missing PASS: %s\n%s\n' "$test_name" "$board_output"; exit 1; }
done
if grep -q -- '--- SKIP' <<<"$board_output"; then printf 'a board test skipped\n%s\n' "$board_output"; exit 1; fi

# The command is registered in the shipped binary.
build_dir="$(mktemp -d)"
trap 'rm -rf "$build_dir"' EXIT
go build -C "$cli_dir" -o "$build_dir/do-work-cli" ./cmd/do-work-cli
help_output="$("$build_dir/do-work-cli" --format json help)"
grep -q -F '"run-status"' <<<"$help_output" || { printf 'run-status is not a registered command\n'; exit 1; }
printf 'REQ-690 GREEN probe: PASS\n'
