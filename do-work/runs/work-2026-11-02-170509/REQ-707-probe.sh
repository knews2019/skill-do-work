#!/usr/bin/env bash
# GREEN probe for REQ-707: ai-report organize files finished months into counted folders, keeps
# links open, and is documented as an ai-report form with help and routing lines.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
output="$(go test -count=1 -v -run '^(TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen|TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets|TestHandlersRegisterCanonicalToolboxCommands)$' ./internal/toolboxcommands/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
grep -q -- '--- PASS: TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen' <<<"$output"
grep -q -- '--- PASS: TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets' <<<"$output"
grep -q -- '--- PASS: TestHandlersRegisterCanonicalToolboxCommands' <<<"$output"
cd "$root"
action="$(cat skills/do-work-toolbox/actions/ai-report.md)"
grep -q 'ai-report-organize' <<<"$action"
grep -q '### Organize form' <<<"$action"
help="$(cat skills/do-work-toolbox/actions/help.md)"
grep -q 'ai-report organize' <<<"$help"
guide="$(cat skills/do-work-toolbox/docs/ai-report-guide.md)"
grep -q 'ai-report organize' <<<"$guide"
printf 'REQ-707 GREEN: organize verb, tests and prose present\n'
