#!/usr/bin/env bash
# REQ-731 pre-build baseline: the existing ai-report-index and organize --check tests pass at base.
set -euo pipefail
top=$(git rev-parse --show-toplevel)
cd "$top/skills/do-work/tools/do-work-cli"
output=$(go test -count=1 -v -run '^(TestAIReportIndex|TestReportOrganizeHintPrintsOnlyWhenWorkRemains)' ./internal/toolboxcommands/ 2>&1) || { printf '%s\n' "$output"; exit 1; }
grep -q -- '--- PASS: TestAIReportIndexCatalogsEveryBundleNamingStyleOnce' <<<"$output"
grep -q -- '--- PASS: TestReportOrganizeHintPrintsOnlyWhenWorkRemains' <<<"$output"
if grep -q -- '--- SKIP' <<<"$output"; then printf '%s\n' "$output"; echo "a skip is not a pass" >&2; exit 1; fi
echo "REQ-731 preflight baseline: PASS"
