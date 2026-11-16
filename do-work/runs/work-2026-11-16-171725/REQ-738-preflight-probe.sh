#!/usr/bin/env bash
# Pre-build baseline for REQ-738 (the ai-report second-opinion launcher shows a short header and
# folds the full direction): the two tests that pin the launcher's header block today, run at the
# base revision. Every named test must PASS; a skip is not a pass.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"

fail() { printf 'REQ-738 pre-flight FAIL: %s\n' "$*" >&2; exit 1; }

header_tests='TestSecondOpinionDirectedRunNamesAndMarksPage TestSecondOpinionTreeTargetsNamesAndParentHeader'
header_pattern="^($(tr ' ' '|' <<<"$header_tests"))\$"
header_output="$(go test -count=1 -v -run "$header_pattern" ./internal/toolboxcommands/ 2>&1)" \
  || fail "header tests failed: $header_output"
for test_name in $header_tests; do
  grep -qF -- "--- PASS: $test_name " <<<"$header_output" || fail "$test_name did not PASS (a skip is not a pass): $header_output"
done

printf 'REQ-738 pre-flight PASS\n'
