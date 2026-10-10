#!/usr/bin/env bash
# Pre-build baseline for REQ-659 (frontmatter set and req append-section): the existing tests in
# the three packages the build touches must be green at base. Covers the frontmatter get tests
# (corehelpers), the section scanner tests (requestmodel) and the two advance tests that guard the
# section-order refactor (lifecycleadvance). Names tests only, never a whole package.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"

run_named_tests() {
  local package_path="$1" test_pattern="$2" output
  shift 2
  output="$(go test -C "$module_dir" -count=1 -v -run "$test_pattern" "$package_path" 2>&1)" || {
    printf '%s\n' "$output"
    printf 'FAIL: go test %s -run %s exited non-zero\n' "$package_path" "$test_pattern"
    exit 1
  }
  local test_name
  for test_name in "$@"; do
    grep -q -- "--- PASS: $test_name " <<<"$output" || {
      printf '%s\n' "$output"
      printf 'FAIL: no --- PASS line for %s\n' "$test_name"
      exit 1
    }
  done
}

run_named_tests ./internal/corehelpers/ \
  '^(TestFrontmatterUsesSchemaNormalizationAndDecodedValues|TestFrontmatterMembershipUsesExitStatusAndEmptyStdout|TestFrontmatterRejectsInvalidMembershipArguments)$' \
  TestFrontmatterUsesSchemaNormalizationAndDecodedValues \
  TestFrontmatterMembershipUsesExitStatusAndEmptyStdout \
  TestFrontmatterRejectsInvalidMembershipArguments
run_named_tests ./internal/requestmodel/ \
  '^(TestVisibleSectionsRetainOffsetsAndProtectUnclosedRegions|TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule)$' \
  TestVisibleSectionsRetainOffsetsAndProtectUnclosedRegions \
  TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule
run_named_tests ./internal/lifecycleadvance/ \
  '^(TestAdvanceCommandPhaseMatrix|TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates)$' \
  TestAdvanceCommandPhaseMatrix \
  TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates
printf 'REQ-659 pre-build baseline: green\n'
