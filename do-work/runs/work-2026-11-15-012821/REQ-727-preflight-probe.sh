#!/usr/bin/env bash
# Pre-build baseline for REQ-727 (ship a relation-map preset for ai-report second-opinion): the
# three preset tests in report_second_opinion_test.go, run at the base revision. The shipped-presets
# test reads the real toolbox folder, so it pins every preset that ships today. Every named test
# must PASS; a skip is not a pass.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"

fail() { printf 'REQ-727 pre-flight FAIL: %s\n' "$*" >&2; exit 1; }

preset_tests='TestSecondOpinionShippedPresetsLoadFromTheToolbox TestSecondOpinionPresetResolutionAndArgumentRules TestSecondOpinionPresetResolvesFromBundleOnly'
preset_pattern="^($(tr ' ' '|' <<<"$preset_tests"))\$"
preset_output="$(go test -count=1 -v -run "$preset_pattern" ./internal/toolboxcommands/ 2>&1)" \
  || fail "preset tests failed: $preset_output"
for test_name in $preset_tests; do
  grep -qF -- "--- PASS: $test_name " <<<"$preset_output" || fail "$test_name did not PASS (a skip is not a pass): $preset_output"
done

printf 'REQ-727 pre-flight PASS\n'
