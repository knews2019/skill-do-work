#!/usr/bin/env bash
# GREEN probe for REQ-706: month folders are read as containers of bundles (catalog, find,
# second-opinion), and the ai-report prose no longer demands a direct child of ai-reports/.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
output="$(go test -count=1 -v -run '^(TestAIReportIndexReadsBundlesInsideMonthFolders|TestSecondOpinionBundleInsideMonthFolder|TestArchitectureScanFindsPriorReportInMonthFolder)$' ./internal/toolboxcommands/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
grep -q -- '--- PASS: TestAIReportIndexReadsBundlesInsideMonthFolders' <<<"$output"
grep -q -- '--- PASS: TestSecondOpinionBundleInsideMonthFolder' <<<"$output"
grep -q -- '--- PASS: TestArchitectureScanFindsPriorReportInMonthFolder' <<<"$output"
cd "$root"
action="$(cat skills/do-work-toolbox/actions/ai-report.md)"
if grep -q 'must be a directory directly under `ai-reports/`' <<<"$action"; then echo 'revise rule still demands a direct child'; exit 1; fi
if grep -q 'physical direct child of `ai-reports/`' <<<"$action"; then echo 'second-opinion rule still demands a direct child'; exit 1; fi
grep -q 'month folder' <<<"$action"
guide="$(cat skills/do-work-toolbox/docs/ai-report-guide.md)"
grep -q 'month folder' <<<"$guide"
printf 'REQ-706 GREEN: month folders read as bundle containers\n'
