#!/usr/bin/env bash
# GREEN probe for REQ-694 (req append-section refuses a body that hides later sections; frontmatter
# set refuses status and id), run by the integrator's test gate. Runs the two new focused tests plus
# the existing canonical-order test (it pins the insert path whose re-check this REQ replaces), then
# checks gofmt and go vet on internal/corehelpers. A skip is not a pass. Fails at base: the two new
# tests do not exist yet, so their --- PASS lines are missing. Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
package_dir="$module_dir/internal/corehelpers"

unformatted="$(gofmt -l "$package_dir")"
[ -z "$unformatted" ] || { printf 'FAIL: gofmt -l lists:\n%s\n' "$unformatted"; exit 1; }
go vet -C "$module_dir" ./internal/corehelpers/

probe_tests=(
  TestRequestAppendSectionRefusesBodyThatHidesLaterSections
  TestFrontmatterSetRefusesLifecycleOwnedFields
  TestRequestAppendSectionLandsInCanonicalOrder
)
test_pattern="^($(IFS='|'; echo "${probe_tests[*]}"))\$"
output="$(go test -C "$module_dir" -count=1 -v -run "$test_pattern" ./internal/corehelpers/ 2>&1)" || {
  printf '%s\n' "$output"
  printf 'FAIL: go test -run %s exited non-zero\n' "$test_pattern"
  exit 1
}
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || {
    printf '%s\n' "$output"
    printf 'FAIL: no --- PASS line for %s (skipped or missing)\n' "$probe_test"
    exit 1
  }
done
printf 'REQ-694 GREEN probe: pass\n'
