#!/usr/bin/env bash
# Pre-build focused baseline for REQ-708: organize, catalog and install-narration tests pass at base.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
output="$(go test -count=1 -v -run '^(TestAIReportOrganize.*|TestAIReportIndexReadsBundlesInsideMonthFolders|TestInstallNarrationNeverReachesStandardOutput)$' ./internal/toolboxcommands/ ./internal/suiteinstall/ 2>&1)"
grep -q -- '--- PASS: TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen' <<<"$output"
grep -q -- '--- PASS: TestInstallNarrationNeverReachesStandardOutput' <<<"$output"
if grep -q -- '--- FAIL' <<<"$output"; then printf '%s\n' "$output"; exit 1; fi
printf 'REQ-708 pre-flight baseline: organize, catalog and install narration tests pass\n'
