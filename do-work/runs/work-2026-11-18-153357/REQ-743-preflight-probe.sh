#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"

run_tests() {
  local package="$1" test_name test_pattern='' output
  shift
  for test_name in "$@"; do
    test_pattern="${test_pattern:+$test_pattern|}$test_name"
  done
  if ! output="$(go test -count=1 -v -run "^($test_pattern)$" "$package" 2>&1)"; then
    printf '%s\n' "$output" >&2
    return 1
  fi
  for test_name in "$@"; do
    if ! grep -qF -- "--- PASS: $test_name (" <<<"$output"; then
      printf 'Expected PASS for %s\n%s\n' "$test_name" "$output" >&2
      return 1
    fi
    printf 'PASS %s\n' "$test_name"
  done
}

run_tests ./internal/toolboxcommands \
  TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded \
  TestAIReportIndexRefusesHandMadeCatalog
run_tests ./internal/suiteinstall \
  TestAnAvailableUpdateInstallsAndVerifiesTheNewVersion \
  TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose
