#!/usr/bin/env bash
# REQ-731 GREEN probe: ai-report-index --check prints the catalog rebuild line only when the bundle
# is missing from ai-reports/catalog.json, and ai-report Step 8 plus the guide document it.
set -euo pipefail
top=$(git rev-parse --show-toplevel)
cd "$top/skills/do-work/tools/do-work-cli"
output=$(go test -count=1 -v -run '^(TestAIReportIndexCheckPrintsRebuildLineOnlyWhenBundleMissing|TestAIReportIndexCatalogsEveryBundleNamingStyleOnce|TestReportOrganizeHintPrintsOnlyWhenWorkRemains)$' ./internal/toolboxcommands/ 2>&1) || { printf '%s\n' "$output"; exit 1; }
grep -q -- '--- PASS: TestAIReportIndexCheckPrintsRebuildLineOnlyWhenBundleMissing' <<<"$output"
grep -q -- '--- PASS: TestAIReportIndexCatalogsEveryBundleNamingStyleOnce' <<<"$output"
grep -q -- '--- PASS: TestReportOrganizeHintPrintsOnlyWhenWorkRemains' <<<"$output"
if grep -q -- '--- SKIP' <<<"$output"; then printf '%s\n' "$output"; echo "a skip is not a pass" >&2; exit 1; fi
action=$(cat "$top/skills/do-work-toolbox/actions/ai-report.md")
grep -q -- 'ai-report-index --check' <<<"$action"
guide=$(sed -n '/^## Report Catalog/,/^## Organizing Reports by Month/p' "$top/skills/do-work-toolbox/docs/ai-report-guide.md")
grep -qi -- 'next rebuild' <<<"$guide"
echo "REQ-731 GREEN probe: PASS"
