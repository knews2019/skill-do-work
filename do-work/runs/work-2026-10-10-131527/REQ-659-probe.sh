#!/usr/bin/env bash
# GREEN probe for REQ-659 (frontmatter set and req append-section), run by the integrator's test
# gate. Runs the five new focused tests plus the two advance tests that guard the shared
# section-order list, then checks gofmt on the touched Go packages. Fails at base: the five new
# tests do not exist yet, so their --- PASS lines are missing.
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
  '^(TestFrontmatterSetStampIsAppendOnly|TestFrontmatterSetWritesQuotedScalars|TestRequestAppendSectionLandsInCanonicalOrder|TestRequestAppendSectionRepeatIsNoOpAndConflictRefuses|TestRequestIDResolvesOneWorkingOrQueueFile)$' \
  TestFrontmatterSetStampIsAppendOnly \
  TestFrontmatterSetWritesQuotedScalars \
  TestRequestAppendSectionLandsInCanonicalOrder \
  TestRequestAppendSectionRepeatIsNoOpAndConflictRefuses \
  TestRequestIDResolvesOneWorkingOrQueueFile
run_named_tests ./internal/lifecycleadvance/ \
  '^(TestAdvanceCommandPhaseMatrix|TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates)$' \
  TestAdvanceCommandPhaseMatrix \
  TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates

unformatted="$(gofmt -l "$module_dir/internal/corehelpers" "$module_dir/internal/requestmodel" "$module_dir/internal/lifecycleadvance")"
[ -z "$unformatted" ] || { printf 'FAIL: gofmt -l lists:\n%s\n' "$unformatted"; exit 1; }
printf 'REQ-659 GREEN probe: pass\n'
