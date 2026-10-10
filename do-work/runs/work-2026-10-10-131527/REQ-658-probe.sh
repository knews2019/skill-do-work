#!/usr/bin/env bash
# REQ-658 GREEN probe: the four finalize --auto-manifest tests pass (a skip is
# not a pass), and both action files name the new command. Checks whichever
# tree it runs in. Expected to FAIL at the base commit: the tests do not exist.
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
cli_root="$repo_root/skills/do-work/tools/do-work-cli"
test_names=(
  TestFinalizeAutoManifestEmitsAManifestFinalizeAccepts
  TestFinalizeAutoManifestRefusesAStagedPathBeforeAnyJournal
  TestFinalizeAutoManifestRefusesWithoutAMessageFile
  TestFinalizeAutoManifestRefusesAStaleReleaseVersionBeforeWriting
)
run_pattern="^($(IFS='|'; printf '%s' "${test_names[*]}"))\$"
status=0
output="$(go test -C "$cli_root" -count=1 -v -run "$run_pattern" ./internal/finalization/ 2>&1)" || status=$?
if [ "$status" -ne 0 ]; then
  printf '%s\n' "$output" >&2
  printf 'FAIL: go test exited %s\n' "$status" >&2
  exit 1
fi
for test_name in "${test_names[@]}"; do
  if ! grep -qF -- "--- PASS: $test_name " <<<"$output"; then
    printf '%s\n' "$output" >&2
    printf 'FAIL: no PASS line for %s\n' "$test_name" >&2
    exit 1
  fi
done
for action_file in skills/do-work/actions/work.md skills/do-work/actions/work-reference.md; do
  if ! grep -qF 'finalize --auto-manifest' "$repo_root/$action_file"; then
    printf 'FAIL: %s does not name finalize --auto-manifest\n' "$action_file" >&2
    exit 1
  fi
done
printf 'REQ-658 probe passed: %s tests, both action files name the command.\n' "${#test_names[@]}"
