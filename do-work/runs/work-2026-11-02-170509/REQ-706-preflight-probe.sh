#!/usr/bin/env bash
# Pre-build focused baseline for REQ-706: the existing catalog and second-opinion bundle tests pass at base.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
output="$(go test -count=1 -v -run '^(TestAIReportIndex.*|TestSecondOpinionLatestAndLinkedBundleRefusal)$' ./internal/toolboxcommands/ 2>&1)"
grep -q -- '--- PASS: TestAIReportIndexCatalogsEveryBundleNamingStyleOnce' <<<"$output"
grep -q -- '--- PASS: TestSecondOpinionLatestAndLinkedBundleRefusal' <<<"$output"
if grep -q -- '--- FAIL' <<<"$output"; then printf '%s\n' "$output"; exit 1; fi
printf 'REQ-706 pre-flight baseline: catalog and second-opinion bundle tests pass\n'
