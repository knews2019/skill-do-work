#!/usr/bin/env bash
# Pre-build focused baseline for REQ-707: toolbox registration and the catalog tests pass at base.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
output="$(go test -count=1 -v -run '^(TestHandlersRegisterCanonicalToolboxCommands|TestAIReportIndex.*)$' ./internal/toolboxcommands/ 2>&1)"
grep -q -- '--- PASS: TestHandlersRegisterCanonicalToolboxCommands' <<<"$output"
grep -q -- '--- PASS: TestAIReportIndexReadsBundlesInsideMonthFolders' <<<"$output"
if grep -q -- '--- FAIL' <<<"$output"; then printf '%s\n' "$output"; exit 1; fi
printf 'REQ-707 pre-flight baseline: registration and catalog tests pass\n'
